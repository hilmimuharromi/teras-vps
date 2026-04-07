package proxmox

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Client represents Proxmox API client
type Client struct {
	HTTPClient *http.Client
	Host       string
	User       string
	Password   string
	Node       string
	Ticket     string // Authentication ticket
	CSRFToken  string // CSRF token
}

// NewClient creates a new Proxmox API client
func NewClient() (*Client, error) {
	host := os.Getenv("PROXMOX_HOST")
	user := os.Getenv("PROXMOX_USER")
	password := os.Getenv("PROXMOX_PASSWORD")
	node := os.Getenv("PROXMOX_NODE")

	if host == "" {
		return nil, errors.New("PROXMOX_HOST environment variable is required")
	}
	if user == "" {
		return nil, errors.New("PROXMOX_USER environment variable is required")
	}
	if password == "" {
		return nil, errors.New("PROXMOX_PASSWORD environment variable is required")
	}

	// Create HTTP client with insecure TLS (for self-signed certs)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	proxmox := &Client{
		HTTPClient: client,
		Host:       host,
		User:       user,
		Password:   password,
		Node:       node,
	}

	// Login to get ticket and CSRF token
	if err := proxmox.Login(); err != nil {
		return nil, fmt.Errorf("failed to login to Proxmox: %w", err)
	}

	return proxmox, nil
}

// Login authenticates with Proxmox API
func (c *Client) Login() error {
	url := fmt.Sprintf("%s/access/ticket", c.Host)

	payload := map[string]string{
		"username": c.User,
		"password": c.Password,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonPayload)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data struct {
			Ticket     string `json:"ticket"`
			CSRFToken string `json:"CSRFPreventionToken"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	c.Ticket = result.Data.Ticket
	c.CSRFToken = result.Data.CSRFToken

	return nil
}

// Request makes an authenticated request to Proxmox API
func (c *Client) Request(method, path string, payload interface{}) ([]byte, error) {
	url := c.Host + path

	var body io.Reader
	if payload != nil {
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = strings.NewReader(string(jsonPayload))
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// Set authentication headers
	req.Header.Set("Cookie", fmt.Sprintf("PVEAuthCookie=%s", c.Ticket))
	req.Header.Set("CSRFPreventionToken", c.CSRFToken)

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		var errResp struct {
			Errors map[string]string `json:"errors"`
		}
		if json.Unmarshal(respBody, &errResp) == nil && len(errResp.Errors) > 0 {
			for _, errMsg := range errResp.Errors {
				return nil, errors.New(errMsg)
			}
		}
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var result struct {
		Data json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return respBody, nil
	}

	return result.Data, nil
}

// GetNextVMID gets the next available VM ID
func (c *Client) GetNextVMID() (int, error) {
	path := fmt.Sprintf("/cluster/nextid")

	data, err := c.Request("GET", path, nil)
	if err != nil {
		return 0, err
	}

	var vmid string
	if err := json.Unmarshal(data, &vmid); err != nil {
		return 0, err
	}

	var vmidInt int
	if _, err := fmt.Sscanf(vmid, "%d", &vmidInt); err != nil {
		return 0, err
	}

	return vmidInt, nil
}

// RefreshTickets refreshes authentication tickets
func (c *Client) RefreshTickets() error {
	return c.Login()
}
