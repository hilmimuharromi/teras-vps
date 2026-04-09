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
	"teras-vps/backend/middleware"
	"time"
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
	Logger     *middleware.ExternalAPILogger
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
		Logger:     middleware.NewExternalAPILogger("Proxmox"),
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

	callID := middleware.GenerateExternalCallID()

	// Log request
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	c.Logger.LogBefore(callID, "POST", url, headers, payload)

	startTime := time.Now()

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonPayload)))
	if err != nil {
		c.Logger.LogAfter(callID, 0, nil, time.Since(startTime), err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		c.Logger.LogAfter(callID, 0, nil, time.Since(startTime), err)
		return err
	}
	defer resp.Body.Close()

	duration := time.Since(startTime)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(body))
		c.Logger.LogAfter(callID, resp.StatusCode, string(body), duration, err)
		return err
	}

	var result struct {
		Data struct {
			Ticket    string `json:"ticket"`
			CSRFToken string `json:"CSRFPreventionToken"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.Logger.LogAfter(callID, resp.StatusCode, nil, duration, err)
		return err
	}

	c.Ticket = result.Data.Ticket
	c.CSRFToken = result.Data.CSRFToken

	// Log success (mask sensitive ticket)
	safeResponse := map[string]string{
		"ticket":    "***MASKED***",
		"csrfToken": result.Data.CSRFToken,
	}
	c.Logger.LogAfter(callID, resp.StatusCode, safeResponse, duration, nil)

	return nil
}

// Request makes an authenticated request to Proxmox API
func (c *Client) Request(method, path string, payload interface{}) ([]byte, error) {
	return c.RequestWithRetry(method, path, payload, true)
}

// RequestWithRetry makes a request with optional retry on 401
func (c *Client) RequestWithRetry(method, path string, payload interface{}, retry bool) ([]byte, error) {
	url := c.Host + path

	var body io.Reader
	var payloadData interface{}
	if payload != nil {
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = strings.NewReader(string(jsonPayload))
		payloadData = payload
	}

	callID := middleware.GenerateExternalCallID()

	// Log request
	headers := map[string]string{
		"Cookie":              fmt.Sprintf("PVEAuthCookie=%s", c.Ticket),
		"CSRFPreventionToken": c.CSRFToken,
	}
	if payload != nil {
		headers["Content-Type"] = "application/json"
	}
	c.Logger.LogBefore(callID, method, url, headers, payloadData)

	startTime := time.Now()

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		c.Logger.LogAfter(callID, 0, nil, time.Since(startTime), err)
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
		c.Logger.LogAfter(callID, 0, nil, time.Since(startTime), err)
		return nil, err
	}
	defer resp.Body.Close()

	duration := time.Since(startTime)

	// Handle 401 - ticket expired, try to refresh
	if resp.StatusCode == 401 && retry {
		c.Logger.LogAfter(callID, 401, nil, duration, fmt.Errorf("ticket expired, refreshing"))

		// Re-login
		if loginErr := c.Login(); loginErr != nil {
			return nil, fmt.Errorf("failed to refresh Proxmox ticket: %w", loginErr)
		}

		// Retry request with fresh ticket
		return c.RequestWithRetry(method, path, payload, false)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Logger.LogAfter(callID, resp.StatusCode, nil, duration, err)
		return nil, err
	}

	if resp.StatusCode >= 400 {
		var errResp struct {
			Errors map[string]string `json:"errors"`
		}
		if json.Unmarshal(respBody, &errResp) == nil && len(errResp.Errors) > 0 {
			for _, errMsg := range errResp.Errors {
				err := errors.New(errMsg)
				c.Logger.LogAfter(callID, resp.StatusCode, string(respBody), duration, err)
				return nil, err
			}
		}
		err := fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
		c.Logger.LogAfter(callID, resp.StatusCode, string(respBody), duration, err)
		return nil, err
	}

	// Parse response
	var result struct {
		Data json.RawMessage `json:"data"`
	}

	var responseData []byte
	if err := json.Unmarshal(respBody, &result); err != nil {
		responseData = respBody
	} else {
		responseData = result.Data
	}

	// Log success
	c.Logger.LogAfter(callID, resp.StatusCode, responseData, duration, nil)

	return responseData, nil
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
