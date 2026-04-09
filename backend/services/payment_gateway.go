package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"teras-vps/backend/middleware"
)

// PaymentGatewayClient represents a generic payment gateway client
type PaymentGatewayClient struct {
	HTTPClient *http.Client
	APIKey     string
	BaseURL    string
	Logger     *middleware.ExternalAPILogger
}

// NewPaymentGatewayClient creates a new payment gateway client
func NewPaymentGatewayClient(apiKey, baseURL string) *PaymentGatewayClient {
	return &PaymentGatewayClient{
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		APIKey:  apiKey,
		BaseURL: baseURL,
		Logger:  middleware.NewExternalAPILogger("PaymentGateway"),
	}
}

// CreatePayment creates a payment
func (c *PaymentGatewayClient) CreatePayment(amount float64, currency string, description string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/payments", c.BaseURL)

	payload := map[string]interface{}{
		"amount":      amount,
		"currency":    currency,
		"description": description,
	}

	return c.doRequest("POST", url, payload)
}

// GetPayment retrieves payment details
func (c *PaymentGatewayClient) GetPayment(paymentID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/payments/%s", c.BaseURL, paymentID)

	return c.doRequest("GET", url, nil)
}

// RefundPayment refunds a payment
func (c *PaymentGatewayClient) RefundPayment(paymentID string, amount float64) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/payments/%s/refund", c.BaseURL, paymentID)

	payload := map[string]interface{}{
		"payment_id": paymentID,
		"amount":     amount,
	}

	return c.doRequest("POST", url, payload)
}

// doRequest makes a request to the payment gateway
func (c *PaymentGatewayClient) doRequest(method, url string, payload interface{}) (map[string]interface{}, error) {
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
		"Authorization": fmt.Sprintf("Bearer %s", c.APIKey),
		"Content-Type":  "application/json",
	}
	c.Logger.LogBefore(callID, method, url, headers, payloadData)

	startTime := time.Now()

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		c.Logger.LogAfter(callID, 0, nil, time.Since(startTime), err)
		return nil, err
	}

	// Set headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		c.Logger.LogAfter(callID, 0, nil, time.Since(startTime), err)
		return nil, err
	}
	defer resp.Body.Close()

	duration := time.Since(startTime)

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Logger.LogAfter(callID, resp.StatusCode, nil, duration, err)
		return nil, err
	}

	// Parse response
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		c.Logger.LogAfter(callID, resp.StatusCode, string(respBody), duration, err)
		return nil, err
	}

	// Log success
	c.Logger.LogAfter(callID, resp.StatusCode, result, duration, nil)

	if resp.StatusCode >= 400 {
		return result, fmt.Errorf("payment gateway request failed with status %d", resp.StatusCode)
	}

	return result, nil
}
