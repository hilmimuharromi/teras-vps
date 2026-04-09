# External API Logging - Summary

## ✅ What Was Implemented

Complete logging for all external API calls (Proxmox, Payment Gateways, third-party services) with detailed request/response information and automatic sensitive data masking.

## 📊 What You'll See

### Example: Proxmox Login
```
═══════════════════════════════════════════════════════════
🚀 EXTERNAL API REQUEST [ext-1712534400123456789]
═══════════════════════════════════════════════════════════
  Service:   Proxmox
  Method:    POST
  URL:       https://192.168.1.100:8006/api2/json/access/ticket
  Timestamp: 2026-04-08 00:30:15.123
  Body:
    {
      "password": "***MASKED***",
      "username": "root@pam"
    }
───────────────────────────────────────────────────────────────

✅ EXTERNAL API RESPONSE [ext-1712534400123456789]
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

## 🎯 Features

✅ **Automatic Logging** - All Proxmox API calls logged automatically
✅ **Request Details** - Method, URL, headers, body
✅ **Response Details** - Status, body, duration
✅ **Sensitive Data Masking** - Passwords, tokens, API keys automatically masked
✅ **Performance Tracking** - See how long external calls take
✅ **Error Logging** - Full error details when calls fail
✅ **Color Coded** - Purple for external, different from internal APIs

## 📁 Files Created/Modified

### New Files
- ✅ `backend/middleware/external_api_logger.go` - External API logger
- ✅ `backend/services/payment_gateway.go` - Payment gateway example with logging
- ✅ `docs/EXTERNAL_API_LOGGING.md` - Complete documentation

### Modified Files
- ✅ `backend/proxmox/client.go` - Added logging to all Proxmox calls

## 🚀 How to Use

### Backend - Just restart and it works!

```bash
cd backend
./teras-vps-api
```

All Proxmox API calls will now be logged automatically!

### What Gets Logged

**Proxmox Operations:**
- ✅ Login/Authentication
- ✅ Create VM
- ✅ Clone VM  
- ✅ Start/Stop/Reboot VM
- ✅ Delete VM
- ✅ Get VM Status
- ✅ Get VM Stats
- ✅ Update VM Config
- ✅ Resize Disk
- ✅ Backup operations

**Payment Gateway (when implemented):**
- ✅ Create Payment
- ✅ Get Payment Details
- ✅ Refund Payment

**Any External Service (easy to add):**
- Email services (SendGrid, Mailgun)
- SMS services (Twilio)
- Cloud providers (AWS, GCP)
- Webhook calls
- Third-party APIs

## 🔒 Security - Sensitive Data Masking

### Automatically Masked:

**Headers:**
- `Authorization` → `Bearer eyJhb...`
- `Cookie` → `PVEAuthCookie=PVE1...`
- `CSRFPreventionToken` → Masked
- `X-API-Key` → Masked

**Body Fields:**
- `password` → `***MASKED***`
- `api_key` → `***MASKED***`
- `token` → `***MASKED***`
- `secret` → `***MASKED***`
- Credit card numbers → `***MASKED***`

## 🐛 Debugging Examples

### Find Why VM Creation Failed
```
🚀 EXTERNAL API REQUEST
  Service:   Proxmox
  Method:    POST
  Body:
    {
      "cores": 2,
      "memory": 2048,
      ...
    }

❌ EXTERNAL API RESPONSE
  Status:    400
  Error:     Invalid value for parameter 'memory'
```
→ Memory value is invalid!

### Find Slow External Calls
```
✅ EXTERNAL API RESPONSE
  Service:   Proxmox
  Duration:  5678.90 ms  ← Too slow!
```
→ Proxmox is slow, check network or Proxmox server

### Check Authentication
```
🚀 EXTERNAL API REQUEST
  Headers:
    • Cookie: PVEAuthCookie=***MASKED***
```
→ If Cookie is missing, login failed!

## 📈 Benefits

1. **Complete Visibility** - See every external API call
2. **Faster Debugging** - Understand external service errors
3. **Performance Monitoring** - Track slow external calls
4. **Security** - Sensitive data masked automatically
5. **Audit Trail** - Full history of external interactions
6. **Consistency** - Same pattern for all external services

## 🔧 Adding Logging to New Services

### Easy 3-Step Process:

```go
// Step 1: Create logger
logger := middleware.NewExternalAPILogger("ServiceName")

// Step 2: Log before request
callID := middleware.GenerateExternalCallID()
logger.LogBefore(callID, "POST", url, headers, body)
startTime := time.Now()

// Step 3: Log after response
duration := time.Since(startTime)
logger.LogAfter(callID, statusCode, responseBody, duration, err)
```

## 📚 Documentation

- ✅ `docs/EXTERNAL_API_LOGGING.md` - Complete guide
- ✅ `docs/REQUEST_LOGGING.md` - Internal API logging
- ✅ `LOGGING_IMPLEMENTATION_SUMMARY.md` - Overall summary
- ✅ `docs/LOG_QUICK_REFERENCE.md` - Quick diagnosis

## 🎨 Log Colors

- **Purple/Magenta**: External API calls (Proxmox, Payment, etc.)
- **Green**: Internal API success (2xx)
- **Yellow**: Internal API client error (4xx)
- **Red**: Internal API server error (5xx)

## ✨ Example Output Comparison

### Before (No Logging):
```
VM creation failed: request failed with status 400
```

### After (With Logging):
```
═══════════════════════════════════════════════════════════
🚀 EXTERNAL API REQUEST [ext-1234567890]
═══════════════════════════════════════════════════════════
  Service:   Proxmox
  Method:    POST
  URL:       https://192.168.1.100:8006/api2/json/nodes/proxmox/qemu
  Body:
    {
      "cores": 2,
      "memory": 2048,
      "disk": 40,
      "name": "my-vps"
    }
───────────────────────────────────────────────────────────────

❌ EXTERNAL API RESPONSE [ext-1234567890]
═══════════════════════════════════════════════════════════
  Service:   Proxmox
  Status:    400
  Duration:  123.45 ms
  Body:
    {
      "errors": {
        "memory": "Invalid value: 2048 (must be multiple of 512)"
      }
    }

  ❌ Error: Invalid value for parameter 'memory'
───────────────────────────────────────────────────────────────
```

Much more helpful! 🎯

## 🚀 Ready to Use!

Just restart your backend and all external API calls will be logged automatically!

```bash
cd backend
./teras-vps-api
```

Now you'll see **every call to Proxmox, payment gateways, and external services** with full details in your terminal! 🎉
