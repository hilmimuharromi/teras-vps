# Complete Logging Guide - Quick Reference

## 🎯 What You Have Now

### 1. **Internal API Logging** (Purple/Color-coded)
Logs every request to your backend API with request/response bodies and call ID.

### 2. **External API Logging** (Magenta)
Logs every call to external services (Proxmox, Payment Gateways, etc.) with request/response bodies.

## 📊 What You'll See in Terminal

### Internal API Request (e.g., Login)
```
═══════════════════════════════════════════════════════════
✅ API CALL [550e8400] POST /api/v1/auth/login
═══════════════════════════════════════════════════════════
📥 REQUEST:
  • Time:      00:15:23.456
  • Client IP: 127.0.0.1
  • Method:    POST
  • Path:      /api/v1/auth/login
  • Body:
    {
      "email": "test@test.com",
      "password": "test12345"
    }

📤 RESPONSE:
  • Status:    200
  • Duration:  45.23 ms
  • Body:
    {
      "success": true,
      "data": { "token": "eyJhbGci..." }
    }
───────────────────────────────────────────────────────────────
```

### External API Request (e.g., Proxmox Login)
```
═══════════════════════════════════════════════════════════
🚀 EXTERNAL API REQUEST [ext-1712534400123]
═══════════════════════════════════════════════════════════
  Service:   Proxmox
  Method:    POST
  URL:       https://192.168.1.100:8006/api2/json/access/ticket
  Body:
    {
      "password": "***MASKED***",
      "username": "root@pam"
    }
───────────────────────────────────────────────────────────────

✅ EXTERNAL API RESPONSE [ext-1712534400123]
═══════════════════════════════════════════════════════════
  Service:   Proxmox
  Status:    200
  Duration:  245.67 ms
  Body:
    {
      "ticket": "***MASKED***",
      "csrfToken": "5f8a7b6c4d3e2f1a"
    }
───────────────────────────────────────────────────────────────
```

## 🎨 Color Legend

| Color | Meaning |
|-------|---------|
| 🟢 Green + ✅ | Success (2xx) |
| 🟡 Yellow + ⚠️ | Client Error (4xx) |
| 🔴 Red + ❌ | Server Error (5xx) |
| 🟣 Magenta | External API calls |
| 🔵 Blue | Information/Headers |

## 🔍 Quick Diagnosis

### Internal API Issues

**Problem**: 401 Unauthorized
```
⚠️ API CALL [550e8400] GET /api/v1/billing/plans
📥 REQUEST:
  (No Auth header!)  ← PROBLEM HERE
📤 RESPONSE:
  • Status: 401
```
**Solution**: Frontend not sending token

---

**Problem**: 400 Bad Request
```
⚠️ API CALL [550e8401] POST /api/v1/vms
📥 REQUEST:
  • Body:
    {
      "name": "",  ← Empty name!
      "plan_id": 0  ← Invalid plan!
    }
```
**Solution**: Check request body format

---

**Problem**: 500 Server Error
```
❌ API CALL [550e8402] POST /api/v1/vms
📤 RESPONSE:
  • Status: 500
  • Body:
    {
      "error": {
        "message": "Failed to create VM"
      }
    }
```
**Solution**: Check backend error logs

### External API Issues

**Problem**: Proxmox Authentication Failed
```
❌ EXTERNAL API REQUEST
  Service: Proxmox
  Body:
    {
      "password": "***MASKED***",
      "username": "root@pam"
    }

❌ EXTERNAL API RESPONSE
  Status: 401
  Error: login failed with status 401
```
**Solution**: Check Proxmox credentials in .env

---

**Problem**: VM Creation Failed
```
🚀 EXTERNAL API REQUEST
  Service: Proxmox
  Body:
    {
      "cores": 2,
      "memory": 2048,
      "name": "my-vps"
    }

❌ EXTERNAL API RESPONSE
  Status: 400
  Error: Invalid configuration
```
**Solution**: Check VM configuration parameters

---

**Problem**: Slow External Service
```
✅ EXTERNAL API RESPONSE
  Service: Proxmox
  Duration:  5678.90 ms  ← Too slow!
```
**Solution**: Check network or Proxmox server performance

## 🎯 Debug Workflows

### Workflow 1: Debug Authentication Issue

1. **Check Internal API Logs**
   ```
   Look for: "Auth:" line in request
   If missing → Frontend not sending token
   If present → Token might be invalid/expired
   ```

2. **Check External API Logs** (if using Proxmox auth)
   ```
   Look for: Proxmox login request
   Check if: Credentials are being sent
   Check response: Ticket received?
   ```

3. **Match with Call ID**
   ```
   Frontend console: 🔍 API [550e8400]
   Backend terminal: API CALL [550e8400]
   Match them to trace the issue!
   ```

### Workflow 2: Debug VM Creation Failure

1. **Check Frontend Console**
   ```
   Look for: Error message
   Note: API call ID
   ```

2. **Check Internal API Log**
   ```
   Find: API CALL [ID]
   See: What was received from frontend
   ```

3. **Check External API Log**
   ```
   Find: EXTERNAL API CALL (Proxmox)
   See: What was sent to Proxmox
   Check: Proxmox response and error
   ```

### Workflow 3: Find Slow Requests

1. **Scan Terminal for Duration**
   ```
   Look for: Duration > 1000ms
   Internal: 45.23 ms ✅
   External: 2345.67 ms ❌ ← Problem!
   ```

2. **Identify Bottleneck**
   ```
   If Internal is slow → Check database queries
   If External is slow → Check network/external service
   ```

## 💡 Pro Tips

### 1. **Search by Call ID**
```bash
# Find specific request in terminal history
grep "550e8400" terminal.log

# Find all calls from specific user
grep "192.168.1.100" terminal.log
```

### 2. **Watch for Patterns**
```
Multiple 401s in a row → Auth system issue
Multiple 500s → Backend bug
All external calls slow → Network issue
```

### 3. **Use Browser DevTools**
```
F12 → Network tab
Click request → Response Headers
Look for: X-API-Call-ID
Match with backend logs!
```

### 4. **Console Helpers** (in browser)
```javascript
// Show recent API calls
window.showAPICalls()

// Export API calls
window.exportAPICalls()
```

## 📁 File Locations

### Backend
- **Internal Logger**: `backend/middleware/request_logger.go`
- **External Logger**: `backend/middleware/external_api_logger.go`
- **Proxmox Client**: `backend/proxmox/client.go` (with logging)
- **Main Config**: `backend/main.go`

### Frontend
- **API Client**: `frontend/src/lib/api.ts` (with tracking)
- **Call Tracker**: `frontend/src/lib/api-call-tracker.ts`

### Documentation
- `docs/REQUEST_LOGGING.md` - Internal API logging guide
- `docs/EXTERNAL_API_LOGGING.md` - External API logging guide
- `docs/LOG_QUICK_REFERENCE.md` - Quick diagnosis
- `LOGGING_IMPLEMENTATION_SUMMARY.md` - Implementation details
- `EXTERNAL_API_LOGGING_SUMMARY.md` - External API summary
- `COMPLETE_LOGGING_GUIDE.md` - This file!

## 🚀 Quick Start

### 1. Start Backend
```bash
cd backend
./teras-vps-api
```

### 2. Start Frontend
```bash
cd frontend
npm run dev
```

### 3. Watch Logs Appear
- **Backend terminal**: Detailed internal + external logs
- **Browser console**: API call tracking (development mode)

### 4. Use Application
- Login, create VMs, view billing
- **Watch the magic happen in terminal!** ✨

## 🎉 Benefits

✅ **Complete Visibility** - See everything happening
✅ **Faster Debugging** - Know exactly what's sent/received
✅ **Performance Monitoring** - Track slow requests
✅ **Security** - Sensitive data masked automatically
✅ **Audit Trail** - Full history of all API calls
✅ **Developer Experience** - Beautiful colored output

## 🔧 Configuration

### Backend - Change Logging Level

In `backend/main.go`:

```go
// Reduce logging for production
app.Use(middleware.RequestLogger(middleware.RequestLoggerConfig{
    LogRequestBody:  true,
    LogResponseBody: false,  // Save space
    MaxBodySize:     1024,    // Smaller bodies
    SkipPaths:       []string{"/health", "/metrics"},
}))
```

### External API - Disable Logging (Not Recommended)

In `backend/proxmox/client.go`:

```go
// Comment out logger initialization
Logger: nil, // middleware.NewExternalAPILogger("Proxmox"),
```

## 📊 Statistics

With this logging setup, you now have:

- ✅ **Internal API logging** - All `/api/v1/*` endpoints
- ✅ **External API logging** - Proxmox, Payment Gateways, etc.
- ✅ **Frontend tracking** - Last 50 API calls tracked
- ✅ **Call ID correlation** - Match frontend to backend
- ✅ **Performance tracking** - Duration for all calls
- ✅ **Error logging** - Detailed error information
- ✅ **Security masking** - Sensitive data protected

## 🎓 Learning Path

1. **Start Simple** - Just watch the logs appear
2. **Learn Patterns** - Notice what normal looks like
3. **Debug Issues** - Use logs to find problems
4. **Optimize** - Use duration to find bottlenecks
5. **Add More** - Extend to new external services

---

**Happy Debugging!** 🚀

You now have world-class logging for your entire application!
