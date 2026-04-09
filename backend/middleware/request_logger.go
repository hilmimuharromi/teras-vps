package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestLoggerConfig holds configuration for the request logger middleware
type RequestLoggerConfig struct {
	// SkipPaths is a list of paths to skip logging
	SkipPaths []string
	// LogRequestBody enables logging of request body (default: false for performance)
	LogRequestBody bool
	// LogResponseBody enables logging of response body (default: false for performance)
	LogResponseBody bool
	// MaxBodySize limits the body size to log (default: 1024 bytes)
	MaxBodySize int
}

// DefaultRequestLoggerConfig returns default config for request logger
func DefaultRequestLoggerConfig() RequestLoggerConfig {
	return RequestLoggerConfig{
		SkipPaths:       []string{"/health"},
		LogRequestBody:  true,
		LogResponseBody: true,
		MaxBodySize:     2048,
	}
}

// RequestLogger creates a middleware that logs detailed request/response info
func RequestLogger(config RequestLoggerConfig) fiber.Handler {
	// Build skip paths map for faster lookup
	skipMap := make(map[string]bool)
	for _, path := range config.SkipPaths {
		skipMap[path] = true
	}

	return func(c *fiber.Ctx) error {
		// Generate unique API call ID
		callID := uuid.New().String()
		c.Locals("call_id", callID)

		// Start timer
		start := time.Now()

		// Get request info
		method := c.Method()
		path := c.Path()
		query := c.Context().QueryArgs().String()

		// Skip if in skip list
		if skipMap[path] {
			return c.Next()
		}

		// Build request log
		logData := map[string]interface{}{
			"call_id":    callID,
			"timestamp":  start.Format("2006-01-02 15:04:05.000"),
			"method":     method,
			"path":       path,
			"client_ip":  c.IP(),
			"user_agent": c.Get("User-Agent"),
		}

		// Add query params if present
		if query != "" {
			logData["query"] = string(query)
		}

		// Add authorization info
		authHeader := c.Get("Authorization")
		if authHeader != "" {
			if strings.HasPrefix(authHeader, "Bearer ") {
				token := authHeader[7:]
				if len(token) > 20 {
					logData["auth_token"] = token[:10] + "..." + token[len(token)-10:]
				} else {
					logData["auth_token"] = token
				}
			} else {
				logData["auth"] = authHeader
			}
		}

		// Log request body if enabled and not GET/HEAD
		if config.LogRequestBody && method != fiber.MethodGet && method != fiber.MethodHead {
			body := c.Body()
			if len(body) > 0 {
				// Try to pretty print JSON
				var jsonData interface{}
				if err := json.Unmarshal(body, &jsonData); err == nil {
					if prettyJSON, err := json.MarshalIndent(jsonData, "", "  "); err == nil {
						if len(prettyJSON) > config.MaxBodySize {
							logData["req_body"] = string(prettyJSON[:config.MaxBodySize]) + "... [truncated]"
						} else {
							logData["req_body"] = string(prettyJSON)
						}
					}
				} else {
					// Not JSON or invalid JSON, log as string
					if len(body) > config.MaxBodySize {
						logData["req_body"] = string(body[:config.MaxBodySize]) + "... [truncated]"
					} else {
						logData["req_body"] = string(body)
					}
				}
			}
		}

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get response status
		status := c.Response().StatusCode()
		contentType := c.Get("Content-Type")

		// Add response info to log
		logData["status"] = status
		logData["duration_ms"] = fmt.Sprintf("%.2f", float64(duration.Microseconds())/1000.0)
		logData["duration_µs"] = duration.Microseconds()
		logData["content_type"] = contentType

		// Add response header with call ID
		c.Set("X-API-Call-ID", callID)

		// Log response body if enabled
		if config.LogResponseBody {
			responseBody := c.Response().Body()
			if len(responseBody) > 0 {
				// Try to pretty print JSON
				var jsonData interface{}
				if err := json.Unmarshal(responseBody, &jsonData); err == nil {
					if prettyJSON, err := json.MarshalIndent(jsonData, "", "  "); err == nil {
						if len(prettyJSON) > config.MaxBodySize {
							logData["res_body"] = string(prettyJSON[:config.MaxBodySize]) + "... [truncated]"
						} else {
							logData["res_body"] = string(prettyJSON)
						}
					}
				} else {
					// Not JSON or invalid JSON, log as string
					if len(responseBody) > config.MaxBodySize {
						logData["res_body"] = string(responseBody[:config.MaxBodySize]) + "... [truncated]"
					} else {
						logData["res_body"] = string(responseBody)
					}
				}
			}
		}

		// Add error info if exists
		if err != nil {
			logData["error"] = err.Error()
		}

		// Determine log level based on status
		var statusColor string
		var statusEmoji string
		switch {
		case status >= 500:
			statusColor = "\033[1;31m" // Red
			statusEmoji = "❌"
		case status >= 400:
			statusColor = "\033[1;33m" // Yellow
			statusEmoji = "⚠️"
		case status >= 300:
			statusColor = "\033[1;36m" // Cyan
			statusEmoji = "↗️"
		case status >= 200:
			statusColor = "\033[1;32m" // Green
			statusEmoji = "✅"
		default:
			statusColor = "\033[0m"
			statusEmoji = "ℹ️"
		}

		// Print formatted log
		fmt.Printf("\n%s═══════════════════════════════════════════════════════════%s\n", "\033[1;34m", "\033[0m")
		fmt.Printf("%s %s API CALL%s [%s] %s %s\n", statusEmoji, statusColor, "\033[0m", callID[:8], method, path)
		fmt.Printf("%s═══════════════════════════════════════════════════════════%s\n", "\033[1;34m", "\033[0m")

		// Print in sections
		fmt.Printf("\033[1m📥 REQUEST%s:\n", "\033[0m")
		fmt.Printf("  • Time:      %s\n", start.Format("15:04:05.000"))
		fmt.Printf("  • Client IP: %s\n", c.IP())
		fmt.Printf("  • Method:    %s\n", method)
		fmt.Printf("  • Path:      %s\n", path)
		if query, ok := logData["query"]; ok {
			fmt.Printf("  • Query:     %s\n", query)
		}
		if auth, ok := logData["auth_token"]; ok {
			fmt.Printf("  • Auth:      Bearer %s\n", auth)
		}
		if reqBody, ok := logData["req_body"]; ok {
			fmt.Printf("  • Body:\n")
			// Indent request body
			bodyStr := reqBody.(string)
			for _, line := range strings.Split(bodyStr, "\n") {
				fmt.Printf("    %s\n", line)
			}
		}

		fmt.Printf("\n\033[1m📤 RESPONSE%s:\n", "\033[0m")
		fmt.Printf("  • Status:    %s%d%s\n", statusColor, status, "\033[0m")
		fmt.Printf("  • Duration:  %s (%d µs)\n", logData["duration_ms"], duration.Microseconds())
		if resBody, ok := logData["res_body"]; ok {
			fmt.Printf("  • Body:\n")
			// Indent response body
			bodyStr := resBody.(string)
			for _, line := range strings.Split(bodyStr, "\n") {
				fmt.Printf("    %s\n", line)
			}
		}

		if err != nil {
			fmt.Printf("\n\033[1;31m❌ ERROR:%s\n", "\033[0m")
			fmt.Printf("  • %s\n", err.Error())
		}

		fmt.Printf("\033[1;34m───────────────────────────────────────────────────────────────\033[0m\n\n")

		return err
	}
}

// SimpleRequestLogger is a simplified version without config
func SimpleRequestLogger() fiber.Handler {
	return RequestLogger(DefaultRequestLoggerConfig())
}

// JSONRequestLogger outputs logs in JSON format for easy parsing
func JSONRequestLogger(config RequestLoggerConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		callID := uuid.New().String()
		c.Locals("call_id", callID)

		start := time.Now()
		method := c.Method()
		path := c.Path()

		// Build JSON log
		logEntry := map[string]interface{}{
			"call_id":   callID,
			"timestamp": start.Format(time.RFC3339),
			"method":    method,
			"path":      path,
			"client_ip": c.IP(),
		}

		// Log request body if enabled
		if config.LogRequestBody && method != fiber.MethodGet && method != fiber.MethodHead {
			body := c.Body()
			if len(body) > 0 {
				var jsonData interface{}
				if err := json.Unmarshal(body, &jsonData); err == nil {
					logEntry["request_body"] = jsonData
				}
			}
		}

		// Process request
		err := c.Next()

		duration := time.Since(start)
		status := c.Response().StatusCode()

		// Add response data
		logEntry["status"] = status
		logEntry["duration_ms"] = float64(duration.Microseconds()) / 1000.0
		logEntry["call_id"] = callID

		// Add response header
		c.Set("X-API-Call-ID", callID)

		// Log response body if enabled
		if config.LogResponseBody {
			responseBody := c.Response().Body()
			if len(responseBody) > 0 {
				var jsonData interface{}
				if err := json.Unmarshal(responseBody, &jsonData); err == nil {
					if len(responseBody) <= config.MaxBodySize {
						logEntry["response_body"] = jsonData
					} else {
						// Truncate and note it
						truncated := responseBody[:config.MaxBodySize]
						var truncatedData interface{}
						if err := json.Unmarshal(truncated, &truncatedData); err == nil {
							logEntry["response_body"] = truncatedData
							logEntry["response_body_truncated"] = true
						}
					}
				}
			}
		}

		if err != nil {
			logEntry["error"] = err.Error()
		}

		// Output as JSON
		jsonBytes, _ := json.Marshal(logEntry)
		fmt.Printf("\n%s\n\n", string(jsonBytes))

		return err
	}
}

// BodyReader is a custom io.ReadCloser that duplicates the body
type BodyReader struct {
	io.ReadCloser
	Body []byte
}

func (br *BodyReader) Read(p []byte) (int, error) {
	n, err := br.ReadCloser.Read(p)
	if n > 0 {
		br.Body = append(br.Body, p[:n]...)
	}
	return n, err
}

// ReadAllBody safely reads the entire body and returns it
func ReadAllBody(body io.ReadCloser) ([]byte, error) {
	defer body.Close()
	return io.ReadAll(body)
}
