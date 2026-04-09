package middleware

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ExternalAPICall represents a single external API call
type ExternalAPICall struct {
	CallID         string        `json:"call_id"`
	Service        string        `json:"service"`
	Method         string        `json:"method"`
	URL            string        `json:"url"`
	RequestHeaders interface{}   `json:"request_headers,omitempty"`
	RequestBody    interface{}   `json:"request_body,omitempty"`
	ResponseStatus int           `json:"response_status"`
	ResponseBody   interface{}   `json:"response_body,omitempty"`
	Duration       time.Duration `json:"duration"`
	Error          string        `json:"error,omitempty"`
	Timestamp      time.Time     `json:"timestamp"`
}

// ExternalAPILogger logs external API calls
type ExternalAPILogger struct {
	serviceName string
}

// NewExternalAPILogger creates a new external API logger
func NewExternalAPILogger(serviceName string) *ExternalAPILogger {
	return &ExternalAPILogger{
		serviceName: serviceName,
	}
}

// LogBefore logs request details before making the call
func (l *ExternalAPILogger) LogBefore(callID, method, url string, headers map[string]string, body interface{}) {
	// Build request log
	logData := map[string]interface{}{
		"call_id":   callID,
		"service":   l.serviceName,
		"direction": "OUTBOUND",
		"method":    method,
		"url":       url,
		"timestamp": time.Now().Format("2006-01-02 15:04:05.000"),
	}

	// Add headers (with masking for sensitive ones)
	if len(headers) > 0 {
		safeHeaders := make(map[string]string)
		for key, value := range headers {
			// Mask sensitive headers
			if l.isSensitiveHeader(key) {
				safeHeaders[key] = l.maskValue(value)
			} else {
				safeHeaders[key] = value
			}
		}
		logData["headers"] = safeHeaders
	}

	// Add body (with masking for sensitive fields)
	if body != nil {
		logData["body"] = l.sanitizeBody(body)
	}

	// Print formatted log
	l.printOutboundLog(callID, logData)
}

// LogAfter logs response details after the call completes
func (l *ExternalAPILogger) LogAfter(callID string, status int, body interface{}, duration time.Duration, err error) {
	// Build response log
	logData := map[string]interface{}{
		"call_id":   callID,
		"service":   l.serviceName,
		"direction": "INBOUND",
		"status":    status,
		"duration":  fmt.Sprintf("%.2f ms (%d µs)", float64(duration.Microseconds())/1000.0, duration.Microseconds()),
	}

	// Add error if exists
	if err != nil {
		logData["error"] = err.Error()
	}

	// Add body (with masking for sensitive fields)
	if body != nil {
		logData["body"] = l.sanitizeBody(body)
	}

	// Print formatted log
	l.printInboundLog(callID, logData, err)
}

// LogComplete logs a complete API call (before + after in one)
func (l *ExternalAPILogger) LogComplete(callID, method, url string, headers map[string]string, reqBody interface{}, status int, resBody interface{}, duration time.Duration, err error) {
	// Build complete log
	logData := map[string]interface{}{
		"call_id":   callID,
		"service":   l.serviceName,
		"method":    method,
		"url":       url,
		"status":    status,
		"duration":  fmt.Sprintf("%.2f ms", float64(duration.Microseconds())/1000.0),
		"timestamp": time.Now().Format("2006-01-02 15:04:05.000"),
	}

	// Add request body if exists
	if reqBody != nil {
		logData["req_body"] = l.sanitizeBody(reqBody)
	}

	// Add response body if exists
	if resBody != nil {
		logData["res_body"] = l.sanitizeBody(resBody)
	}

	// Add error if exists
	if err != nil {
		logData["error"] = err.Error()
	}

	// Determine color based on status
	var statusColor string
	var statusEmoji string
	switch {
	case err != nil:
		statusColor = "\033[1;31m"
		statusEmoji = "❌"
	case status >= 500:
		statusColor = "\033[1;31m"
		statusEmoji = "❌"
	case status >= 400:
		statusColor = "\033[1;33m"
		statusEmoji = "⚠️"
	case status >= 200:
		statusColor = "\033[1;32m"
		statusEmoji = "✅"
	default:
		statusColor = "\033[0m"
		statusEmoji = "ℹ️"
	}

	// Print compact log
	fmt.Printf("\n%s%s EXTERNAL API CALL [%s] %s %s → %s%d%s (%s)%s\n",
		"\033[1;35m", statusEmoji, callID[:8], method, url,
		statusColor, status, "\033[0m",
		fmt.Sprintf("%.2fms", float64(duration.Microseconds())/1000.0), "\033[0m")

	// Print details
	fmt.Printf("  📥 Request:\n")
	if reqBody != nil {
		bodyStr := l.formatBody(reqBody)
		for _, line := range strings.Split(bodyStr, "\n") {
			fmt.Printf("    %s\n", line)
		}
	} else {
		fmt.Printf("    (no body)\n")
	}

	fmt.Printf("  📤 Response:\n")
	if resBody != nil {
		bodyStr := l.formatBody(resBody)
		for _, line := range strings.Split(bodyStr, "\n") {
			fmt.Printf("    %s\n", line)
		}
	} else {
		fmt.Printf("    (no body)\n")
	}

	if err != nil {
		fmt.Printf("  %s❌ Error: %s%s\n", "\033[1;31m", err.Error(), "\033[0m")
	}

	fmt.Printf("%s───────────────────────────────────────────────────────────────%s\n\n", "\033[1;35m", "\033[0m")
}

// Helper functions

func (l *ExternalAPILogger) isSensitiveHeader(key string) bool {
	sensitiveHeaders := []string{
		"authorization",
		"cookie",
		"CSRFPreventionToken",
		"X-API-Key",
		"X-Secret",
		"X-Access-Token",
	}

	keyLower := strings.ToLower(key)
	for _, sensitive := range sensitiveHeaders {
		if keyLower == strings.ToLower(sensitive) {
			return true
		}
	}
	return false
}

func (l *ExternalAPILogger) maskValue(value string) string {
	if len(value) <= 10 {
		return "****"
	}
	return value[:5] + "..." + value[len(value)-5:]
}

func (l *ExternalAPILogger) sanitizeBody(body interface{}) interface{} {
	// Convert to JSON
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return body
	}

	// Parse back to map for sanitization
	var data map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		return body
	}

	// Mask sensitive fields
	sensitiveFields := []string{
		"password",
		"passwd",
		"secret",
		"token",
		"api_key",
		"apiKey",
		"credit_card",
		"card_number",
		"cvv",
	}

	for _, field := range sensitiveFields {
		if _, exists := data[field]; exists {
			data[field] = "***MASKED***"
		}
	}

	return data
}

func (l *ExternalAPILogger) formatBody(body interface{}) string {
	jsonBytes, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", body)
	}
	return string(jsonBytes)
}

func (l *ExternalAPILogger) printOutboundLog(callID string, logData map[string]interface{}) {
	fmt.Printf("\n%s═══════════════════════════════════════════════════════════%s\n", "\033[1;35m", "\033[0m")
	fmt.Printf("%s🚀 EXTERNAL API REQUEST [%s]%s\n", "\033[1;35m", callID[:8], "\033[0m")
	fmt.Printf("%s═══════════════════════════════════════════════════════════%s\n", "\033[1;35m", "\033[0m")
	fmt.Printf("  Service:   %s\n", l.serviceName)
	fmt.Printf("  Method:    %s\n", logData["method"])
	fmt.Printf("  URL:       %s\n", logData["url"])
	fmt.Printf("  Timestamp: %s\n", logData["timestamp"])

	if headers, ok := logData["headers"]; ok {
		fmt.Printf("  Headers:\n")
		headersMap := headers.(map[string]string)
		for key, value := range headersMap {
			fmt.Printf("    • %s: %s\n", key, value)
		}
	}

	if body, ok := logData["body"]; ok {
		fmt.Printf("  Body:\n")
		bodyStr := l.formatBody(body)
		for _, line := range strings.Split(bodyStr, "\n") {
			fmt.Printf("    %s\n", line)
		}
	}

	fmt.Printf("\033[1;35m───────────────────────────────────────────────────────────────\033[0m\n\n")
}

func (l *ExternalAPILogger) printInboundLog(callID string, logData map[string]interface{}, err error) {
	var statusColor string
	var statusEmoji string
	status := logData["status"].(int)

	switch {
	case err != nil:
		statusColor = "\033[1;31m"
		statusEmoji = "❌"
	case status >= 500:
		statusColor = "\033[1;31m"
		statusEmoji = "❌"
	case status >= 400:
		statusColor = "\033[1;33m"
		statusEmoji = "⚠️"
	case status >= 200:
		statusColor = "\033[1;32m"
		statusEmoji = "✅"
	default:
		statusColor = "\033[0m"
		statusEmoji = "ℹ️"
	}

	fmt.Printf("%s═══════════════════════════════════════════════════════════%s\n", "\033[1;35m", "\033[0m")
	fmt.Printf("%s%s EXTERNAL API RESPONSE [%s]%s\n", statusEmoji, statusColor, callID[:8], "\033[0m")
	fmt.Printf("%s═══════════════════════════════════════════════════════════%s\n", "\033[1;35m", "\033[0m")
	fmt.Printf("  Service:   %s\n", l.serviceName)
	fmt.Printf("  Status:    %s%d%s\n", statusColor, status, "\033[0m")
	fmt.Printf("  Duration:  %s\n", logData["duration"])

	if body, ok := logData["body"]; ok {
		fmt.Printf("  Body:\n")
		bodyStr := l.formatBody(body)
		for _, line := range strings.Split(bodyStr, "\n") {
			fmt.Printf("    %s\n", line)
		}
	}

	if err != nil {
		fmt.Printf("\n  %s❌ Error: %s%s\n", "\033[1;31m", err.Error(), "\033[0m")
	}

	fmt.Printf("\033[1;35m───────────────────────────────────────────────────────────────\033[0m\n\n")
}

// GenerateCallID generates a unique call ID for external API calls
func GenerateExternalCallID() string {
	// Use timestamp-based ID for simplicity
	return fmt.Sprintf("ext-%d", time.Now().UnixNano())
}
