# External API Logging

## Overview

All external API calls (Proxmox, Payment Gateways, third-party services) are now logged with detailed request/response information, similar to the internal API logging.

## What's Logged

### For Every External API Call:

✅ **Request Details:**
- Service name (Proxmox, PaymentGateway, etc.)
- HTTP method (GET, POST, PUT, DELETE)
- Full URL
- Request headers (with sensitive data masked)
- Request body (with passwords/tokens masked)
- Timestamp

✅ **Response Details:**
- HTTP status code
- Response body
- Duration (milliseconds & microseconds)
- Error information (if any)

✅ **Security Features:**
- Passwords are automatically masked: `***MASKED***`
- API keys are partially shown: `sk-1234...6789`
- Tokens are masked: `eyJhbGc...hlP8`
- Tickets/CSRF tokens are masked

## Example Output

### Proxmox Login
```
═══════════════════════════════════════════════════════════
🚀 EXTERNAL API REQUEST [ext-1234567890]
═══════════════════════════════════════════════════════════
  Service:   Proxmox
  Method:    POST
  URL:       https://192.168.1.100:8006/api2/json/access/ticket
  Timestamp: 2026-04-08 00:30:15.123
  Headers:
    • Content-Type: application/json
  Body:
    {
      "password": "***MASKED***",
      "username": "root@pam"
    }
───────────────────────────────────────────────────────────────

✅ EXTERNAL API RESPONSE [ext-1234567890]
═══════════════════════════════════════════════════════════
  Service:   Proxmox
  Status:    200
  Duration:  245.67 ms (245670 µs)
  Body:
    {
      "csrfToken": "5f8a7b6c4d3e2f1a",
      "ticket": "***MASKED***"
    }
───────────────────────────────────────────────────────────────
```

### Proxmox VM Creation
```
═══════════════════════════════════════════════════════════
🚀 EXTERNAL API REQUEST [ext-1234567891]
═══════════════════════════════════════════════════════════
  Service:   Proxmox
  Method:    POST
  URL:       https://192.168.1.100:8006/api2/json/nodes/proxmox/qemu
  Timestamp: 2026-04-08 00:30:16.456
  Headers:
    • Cookie: PVEAuthCookie=PVE1:root@pam:65432100...
    • CSRFPreventionToken: 5f8a7b6c4d3e2f1a
    • Content-Type: application/json
  Body:
    {
      "bootdisk": "scsi0",
      "cores": 2,
      "disk": 40,
      "memory": 2048,
      "name": "my-vps-server",
      "net0": "virtio,bridge=vmbr0",
      "ostype": "l26",
      "scsihw": "virtio-scsi-pci"
    }
───────────────────────────────────────────────────────────────

✅ EXTERNAL API RESPONSE [ext-1234567891]
═══════════════════════════════════════════════════════════
  Service:   Proxmox
  Status:    200
  Duration:  1234.56 ms (1234560 µs)
  Body:
    {
      "task_id": "UPID:proxmox:00001234:00012345:65432100:qmcreate:100:root@pam:"
    }
───────────────────────────────────────────────────────────────
```

### Failed Request
```
═══════════════════════════════════════════════════════════
🚀 EXTERNAL API REQUEST [ext-1234567892]
═══════════════════════════════════════════════════════════
  Service:   Proxmox
  Method:    GET
  URL:       https://192.168.1.100:8006/api2/json/nodes/proxmox/qemu/999/status/current
  Timestamp: 2026-04-08 00:30:17.789
  Headers:
    • Cookie: PVEAuthCookie=PVE1:root@pam:65432100...
    • CSRFPreventionToken: 5f8a7b6c4d3e2f1a
───────────────────────────────────────────────────────────────

❌ EXTERNAL API RESPONSE [ext-1234567892]
═══════════════════════════════════════════════════════════
  Service:   Proxmox
  Status:    404
  Duration:  89.12 ms (89120 µs)
  Body:
    {
      "data": null,
      "errors": {
        "vmid": "VM 999 not found"
      }
    }

  ❌ Error: VM 999 not found
───────────────────────────────────────────────────────────────
```

## Services with Logging

### 1. **Proxmox API** ✅ (Implemented)
- Location: `backend/proxmox/client.go`
- Logs all VM operations:
  - Login/Authentication
  - Create VM
  - Clone VM
  - Start/Stop/Reboot VM
  - Delete VM
  - Get VM Status/Stats
  - VM configuration changes
  - Backup operations

### 2. **Payment Gateway** ✅ (Template Ready)
- Location: `backend/services/payment_gateway.go`
- Example implementation for:
  - Xendit
  - Midtrans
  - Stripe
  - Any payment gateway

### 3. **Any External Service** 📝 (Easy to Add)
Use the logger for any third-party API:
- Email services (SendGrid, Mailgun)
- SMS services (Twilio)
- Cloud providers (AWS, GCP)
- Monitoring services
- Webhook calls

## How to Use

### For Proxmox (Already Done)

The Proxmox client automatically logs all requests. No configuration needed!

```go
// In main.go or wherever you initialize Proxmox
proxmoxClient, err := proxmox.NewClient()
// Logger is automatically initialized
```

### For Payment Gateway

```go
// Initialize payment gateway
paymentClient := services.NewPaymentGatewayClient(
    os.Getenv("XENDIT_API_KEY"),
    "https://api.xendit.co",
)

// Use it - logging is automatic!
payment, err := paymentClient.CreatePayment(100000, "IDR", "VPS Starter Plan")
```

### For Custom External API

```go
import "teras-vps/backend/middleware"

// Create logger
logger := middleware.NewExternalAPILogger("SendGrid")

// Make request with logging
callID := middleware.GenerateExternalCallID()

logger.LogBefore(callID, "POST", url, headers, payload)
startTime := time.Now()

// Make HTTP request
resp, err := httpClient.Do(req)

duration := time.Since(startTime)
logger.LogAfter(callID, resp.StatusCode, responseBody, duration, err)
```

## Sensitive Data Masking

The logger automatically masks sensitive information:

### Masked Headers:
- `Authorization`
- `Cookie`
- `CSRFPreventionToken`
- `X-API-Key`
- `X-Secret`
- `X-Access-Token`

### Masked Body Fields:
- `password`
- `passwd`
- `secret`
- `token`
- `api_key` / `apiKey`
- `credit_card`
- `card_number`
- `cvv`

### Example of Masking:
```
Original:  {"password": "my-secret-123"}
Logged:    {"password": "***MASKED***"}

Original:  {"api_key": "sk-1234567890abcdef"}
Logged:    {"api_key": "***MASKED***"}

Original:  {"Authorization": "Bearer eyJhbGciOiJIUzI1NiIs..."}
Logged:    {"Authorization": "Bearer eyJhb..."}
```

## Color Coding

External API logs use **magenta/purple** color to distinguish from internal API logs:

- **Purple/Magenta**: External API calls (Proxmox, Payment, etc.)
- **Green**: Internal API success (2xx)
- **Yellow**: Internal API client error (4xx)
- **Red**: Internal API server error (5xx)

## Performance Tracking

All external API calls include duration information:

```
Duration:  245.67 ms (245670 µs)
```

This helps identify:
- Slow external services
- Network issues
- Timeout problems
- Performance bottlenecks

### Example: Identifying Slow Service
```
✅ EXTERNAL API RESPONSE [ext-123]
  Service:   Proxmox
  Status:    200
  Duration:  5678.90 ms  ← This is too slow!
```

## Debugging Tips

### 1. **Authentication Failures**
Check if credentials are being sent correctly:
```
🚀 EXTERNAL API REQUEST
  Body:
    {
      "password": "***MASKED***",  ← Check if present
      "username": "root@pam"
    }
```

### 2. **VM Creation Issues**
Check what configuration was sent:
```
🚀 EXTERNAL API REQUEST
  Body:
    {
      "cores": 2,
      "memory": 2048,
      "disk": 40
      ...
    }

❌ EXTERNAL API RESPONSE
  Error: Invalid configuration...
```

### 3. **Payment Problems**
Verify payment payload:
```
🚀 EXTERNAL API REQUEST
  Service:   PaymentGateway
  Body:
    {
      "amount": 100000,
      "currency": "IDR",
      "description": "VPS Starter Plan"
    }

❌ EXTERNAL API RESPONSE
  Status: 400
  Error: Invalid amount
```

### 4. **Network Issues**
Look for requests that fail immediately:
```
❌ EXTERNAL API RESPONSE
  Status:    0  ← No response received
  Duration:  1.23 ms  ← Very fast = connection failed
  Error: dial tcp: connection refused
```

## Configuration

### Disable Logging (Not Recommended)

If you need to disable external API logging temporarily:

```go
// In proxmox/client.go, comment out logger initialization:
proxmox := &Client{
    // ...
    Logger: nil, // middleware.NewExternalAPILogger("Proxmox"),
}
```

Or add a check in Request method:

```go
func (c *Client) Request(method, path string, payload interface{}) ([]byte, error) {
    // Skip logging if logger is nil
    if c.Logger != nil {
        c.Logger.LogBefore(...)
    }
    // ... rest of code
}
```

### Enable JSON Logging

For production log aggregation (ELK, Datadog):

```go
// Create JSON logger instead
logger := middleware.NewExternalAPILoggerJSON("Proxmox")
```

## Adding Logging to New Services

### Step 1: Create Logger
```go
type MyExternalService struct {
    HTTPClient *http.Client
    APIKey     string
    Logger     *middleware.ExternalAPILogger
}

func NewMyService(apiKey string) *MyExternalService {
    return &MyExternalService{
        HTTPClient: &http.Client{Timeout: 30 * time.Second},
        APIKey:     apiKey,
        Logger:     middleware.NewExternalAPILogger("MyService"),
    }
}
```

### Step 2: Log Requests
```go
func (s *MyService) DoSomething(data interface{}) (interface{}, error) {
    callID := middleware.GenerateExternalCallID()
    
    headers := map[string]string{
        "Authorization": fmt.Sprintf("Bearer %s", s.APIKey),
        "Content-Type":  "application/json",
    }
    
    s.Logger.LogBefore(callID, "POST", url, headers, data)
    startTime := time.Now()
    
    // Make request...
}
```

### Step 3: Log Responses
```go
    // After getting response:
    duration := time.Since(startTime)
    
    if err != nil {
        s.Logger.LogAfter(callID, 0, nil, duration, err)
        return nil, err
    }
    
    s.Logger.LogAfter(callID, resp.StatusCode, responseBody, duration, nil)
```

## Files Modified

### Backend
- ✅ `backend/middleware/external_api_logger.go` (NEW)
- ✅ `backend/proxmox/client.go` (UPDATED - added logging)
- ✅ `backend/services/payment_gateway.go` (NEW - example)

## Benefits

1. **Visibility** - See all external API calls in real-time
2. **Debugging** - Understand why external services fail
3. **Performance** - Track slow external services
4. **Security** - Sensitive data is automatically masked
5. **Audit Trail** - Complete history of external API interactions
6. **Consistency** - Same logging pattern for all external services

## Related Documentation

- `docs/REQUEST_LOGGING.md` - Internal API logging
- `LOGGING_IMPLEMENTATION_SUMMARY.md` - Overall logging summary
- `docs/LOG_QUICK_REFERENCE.md` - Quick diagnosis guide

---

**Now you have complete visibility into all external API calls!** 🚀

Every call to Proxmox, payment gateways, or any external service will be logged with full details in your backend terminal.
