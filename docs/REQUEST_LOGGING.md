# Request/Response Logging Middleware

## Overview

The new logging middleware captures detailed information about every API request and response, including:

- ✅ **Unique API Call ID** (UUID) for tracking
- ✅ **Request body** (formatted JSON)
- ✅ **Response body** (formatted JSON)
- ✅ **Request/Response headers**
- ✅ **Execution duration** (milliseconds & microseconds)
- ✅ **Client IP & User Agent**
- ✅ **Authentication token** (partial view for security)
- ✅ **Color-coded status** (✅ success, ⚠️ client error, ❌ server error)

## Features

### 1. **API Call ID**
Every request gets a unique UUID that is:
- Logged in the terminal
- Returned in response header `X-API-Call-ID`
- Can be used to trace specific requests in logs

### 2. **Detailed Request Logging**
```
📥 REQUEST:
  • Time:      00:15:23.456
  • Client IP: 127.0.0.1
  • Method:    POST
  • Path:      /api/v1/auth/login
  • Auth:      Bearer eyJhbGci...
  • Body:
    {
      "email": "test@test.com",
      "password": "*****"
    }
```

### 3. **Detailed Response Logging**
```
📤 RESPONSE:
  • Status:    200
  • Duration:  45.23 ms (45230 µs)
  • Body:
    {
      "success": true,
      "data": {
        "token": "eyJhbGci...",
        "user": {...}
      }
    }
```

### 4. **Color-Coded Output**
- ✅ **Green** (2xx): Success
- ⚠️ **Yellow** (4xx): Client error
- ❌ **Red** (5xx): Server error
- ↗️ **Cyan** (3xx): Redirect

### 5. **Smart Body Logging**
- Request bodies are only logged for POST, PUT, PATCH, DELETE
- Bodies are pretty-printed as JSON
- Bodies are truncated if > 2048 bytes (configurable)
- Non-JSON bodies are logged as plain text

## Usage

### Default Configuration (Already Enabled)

The middleware is already configured in `main.go` with sensible defaults:

```go
app.Use(middleware.RequestLogger(middleware.RequestLoggerConfig{
    LogRequestBody:  true,
    LogResponseBody: true,
    MaxBodySize:     2048,
    SkipPaths:       []string{"/health"},
}))
```

### Custom Configuration

You can customize the behavior:

```go
app.Use(middleware.RequestLogger(middleware.RequestLoggerConfig{
    LogRequestBody:  true,        // Log request bodies
    LogResponseBody: false,       // Don't log response bodies (save space)
    MaxBodySize:     4096,        // Increase truncation limit
    SkipPaths:       []string{    // Skip these paths
        "/health",
        "/api/v1/public/*",
    },
}))
```

### JSON Logging Mode

For production/ELK stack, use JSON output:

```go
app.Use(middleware.JSONRequestLogger(middleware.RequestLoggerConfig{
    LogRequestBody:  true,
    LogResponseBody: true,
    MaxBodySize:     2048,
}))
```

Output example:
```json
{
  "call_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2026-04-08T00:15:23.456Z",
  "method": "POST",
  "path": "/api/v1/auth/login",
  "client_ip": "127.0.0.1",
  "status": 200,
  "duration_ms": 45.23,
  "request_body": {
    "email": "test@test.com",
    "password": "test12345"
  },
  "response_body": {
    "success": true,
    "data": {...}
  }
}
```

## Example Terminal Output

```
═══════════════════════════════════════════════════════════
✅ API CALL [550e8400] POST /api/v1/auth/login
═══════════════════════════════════════════════════════════
📥 REQUEST:
  • Time:      00:15:23.456
  • Client IP: 127.0.0.1
  • Method:    POST
  • Path:      /api/v1/auth/login
  • Auth:      Bearer eyJhbGci...hlP8
  • Body:
    {
      "email": "test@test.com",
      "password": "test12345"
    }

📤 RESPONSE:
  • Status:    200
  • Duration:  45.23 ms (45230 µs)
  • Body:
    {
      "success": true,
      "message": "Login successful",
      "data": {
        "token": "eyJhbGci...",
        "user": {
          "id": 5,
          "email": "test@test.com",
          "username": "testuser",
          "role": "customer"
        }
      }
    }
───────────────────────────────────────────────────────────────
```

## Accessing Call ID in Handlers

You can access the API call ID in your route handlers:

```go
func (c *BillingController) ListPlans(ctx *fiber.Ctx) error {
    callID := ctx.Locals("call_id").(string)
    
    // Use it in error logs, responses, etc.
    log.Printf("Fetching plans [call_id: %s]", callID)
    
    // Or return it in response
    return ctx.JSON(fiber.Map{
        "call_id": callID,
        "success": true,
        "data": fiber.Map{
            "plans": plans,
        },
    })
}
```

## Response Header

Every response includes the call ID:

```
HTTP/1.1 200 OK
Content-Type: application/json
X-API-Call-ID: 550e8400-e29b-41d4-a716-446655440000

{
  "success": true,
  "data": {...}
}
```

This allows frontend to report the call ID for debugging:

```javascript
const response = await fetch('/api/v1/billing/plans');
const callID = response.headers.get('X-API-Call-ID');
console.log(`API Call ID: ${callID}`);
```

## Performance Considerations

### Production Recommendations

1. **Disable response body logging** for high-traffic endpoints:
```go
LogResponseBody: false
```

2. **Reduce body size limit** to save memory:
```go
MaxBodySize: 512  // bytes
```

3. **Skip static/health endpoints**:
```go
SkipPaths: []string{"/health", "/favicon.ico", "/static/*"}
```

4. **Use JSON logging mode** for log aggregation:
```go
middleware.JSONRequestLogger(config)
```

### Current Configuration

The current setup logs everything except `/health` endpoint. This is perfect for development and debugging. For production, consider:

```go
middleware.RequestLogger(middleware.RequestLoggerConfig{
    LogRequestBody:  true,
    LogResponseBody: false,  // Save bandwidth
    MaxBodySize:     512,     // Reduce memory
    SkipPaths:       []string{
        "/health",
        "/metrics",
    },
})
```

## Debugging Tips

### 1. **Find a Specific Request**
Search terminal output by call ID:
```bash
# In terminal history or log file
grep "550e8400" server.log
```

### 2. **Track Error Responses**
Look for ⚠️ or ❌ emojis in logs:
```bash
# Find all 4xx errors
grep "⚠️" server.log

# Find all 5xx errors
grep "❌" server.log
```

### 3. **Slow Requests**
Find requests taking longer than expected:
```bash
# In JSON log mode, filter by duration
cat server.log | jq 'select(.duration_ms > 1000)'
```

### 4. **Authentication Issues**
Check if token is being sent:
```
Look for "Auth: Bearer ..." in request logs
If missing → frontend not sending token
If present but 401 → token invalid/expired
```

## Migration from Old Logger

The old Fiber logger (`github.com/gofiber/fiber/v2/middleware/logger`) has been replaced with this custom middleware.

**Old output:**
```
00:03:59 | 401 |      33.708µs | 127.0.0.1 | GET | /api/v1/billing/plans | -
```

**New output:**
```
═══════════════════════════════════════════════════════════
⚠️ API CALL [550e8400] GET /api/v1/billing/plans
═══════════════════════════════════════════════════════════
📥 REQUEST:
  • Time:      00:03:59.123
  • Client IP: 127.0.0.1
  • Method:    GET
  • Path:      /api/v1/billing/plans

📤 RESPONSE:
  • Status:    401
  • Duration:  33.71 ms (33708 µs)
  • Body:
    {
      "success": false,
      "error": {
        "code": "MISSING_TOKEN",
        "message": "Authorization header is required"
      }
    }
───────────────────────────────────────────────────────────────
```

Much more helpful for debugging! 🎯

## Files Modified

- `backend/middleware/request_logger.go` - New middleware file
- `backend/main.go` - Updated to use new middleware
- `backend/go.mod` - Added `github.com/google/uuid` dependency

## Environment Variables

No new environment variables needed. Configuration is done in code.

## Troubleshooting

### Issue: Logs not showing
- Check that backend was restarted after code changes
- Look for compilation errors in terminal

### Issue: Body not logged
- GET/HEAD requests don't log request body (by design)
- Empty bodies are not logged
- Check `LogRequestBody` and `LogResponseBody` config

### Issue: Call ID not in response
- Check response headers in browser DevTools → Network tab
- Look for `X-API-Call-ID` header

## Future Enhancements

Potential improvements:
- [ ] Log to file instead of stdout
- [ ] Integrate with logging services (Logrus, Zap, etc.)
- [ ] Add request tracing across services
- [ ] Add metrics collection (Prometheus, etc.)
- [ ] Add log rotation
- [ ] Add structured logging mode for production
