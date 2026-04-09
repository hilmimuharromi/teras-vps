# Request/Response Logging - Implementation Summary

## ✅ What Was Added

### Backend Changes

#### 1. New Middleware: `backend/middleware/request_logger.go`
A comprehensive logging middleware that captures:
- ✅ **Unique API Call ID** (UUID v4) for every request
- ✅ **Request body** (pretty-printed JSON, truncated to 2KB)
- ✅ **Response body** (pretty-printed JSON, truncated to 2KB)
- ✅ **Request metadata** (IP, method, path, query params, user agent)
- ✅ **Authentication token** (partial view for security)
- ✅ **Response status** (color-coded: ✅ green, ⚠️ yellow, ❌ red)
- ✅ **Execution duration** (in milliseconds and microseconds)
- ✅ **Error information** (if any)

#### 2. Updated: `backend/main.go`
- Replaced default Fiber logger with custom `RequestLogger`
- Configured with detailed logging for development
- Skips `/health` endpoint to reduce noise

#### 3. Updated: `backend/go.mod`
- Added `github.com/google/uuid` dependency for UUID generation

### Frontend Changes

#### 1. New Helper: `frontend/src/lib/api-call-tracker.ts`
- `getAPICallID(response)` - Extract call ID from response
- `fetchWithCallID(url, options)` - Enhanced fetch with logging
- `apiTracker` - Singleton that tracks recent 50 API calls
- Browser console helpers: `window.trackAPI()`, `window.showAPICalls()`

#### 2. Updated: `frontend/src/lib/api.ts`
- Integrated with `apiTracker` to track all API calls
- Auto-logs call ID in development mode
- Call ID visible in browser console

## 📊 What You'll See Now

### Backend Terminal Output

**Before (old logger):**
```
00:10:19 | 401 |      60.958µs | 127.0.0.1 | GET | /api/v1/billing/plans | -
```

**After (new logger):**
```
═══════════════════════════════════════════════════════════
✅ API CALL [550e8400] GET /api/v1/billing/plans
═══════════════════════════════════════════════════════════
📥 REQUEST:
  • Time:      00:15:23.456
  • Client IP: 127.0.0.1
  • Method:    GET
  • Path:      /api/v1/billing/plans
  • Auth:      Bearer eyJhbGci...hlP8

📤 RESPONSE:
  • Status:    200
  • Duration:  12.45 ms (12450 µs)
  • Body:
    {
      "success": true,
      "data": {
        "plans": [...]
      }
    }
───────────────────────────────────────────────────────────────
```

### Browser Console Output (Development Mode)

```
🔍 API [550e8400] GET /api/v1/billing/plans → 200
🔍 API [550e8401] POST /api/v1/auth/login → 200
🔍 API [550e8402] GET /api/v1/vms → 200
```

## 🔧 How to Use

### 1. Backend - Watch Terminal

Just run your backend as usual:
```bash
cd backend
go run main.go
```

Every request will now show detailed logs!

### 2. Frontend - Browser Console

Open browser DevTools (F12) → Console tab

You'll see API calls logged automatically in development mode.

### 3. Track API Calls in Browser

In browser console, you can use:
```javascript
// Show recent API calls in table format
window.showAPICalls()

// Export all tracked calls
window.exportAPICalls()
```

### 4. Find Specific Request by Call ID

**From Backend Terminal:**
1. Note the call ID (e.g., `550e8400`)
2. Search logs: `grep "550e8400" backend.log`

**From Frontend:**
1. Open DevTools → Network tab
2. Click on any request
3. Look for `X-API-Call-ID` in Response Headers
4. Match with backend logs

## 🎯 Use Cases

### Debug Authentication Issues
```
📥 REQUEST:
  • Auth: (missing!)
  
📤 RESPONSE:
  • Status: 401
  • Body: { "error": { "code": "MISSING_TOKEN" } }
```
→ Frontend not sending token!

### Debug Bad Requests
```
📥 REQUEST:
  • Body: { "email": "invalid" }
  
📤 RESPONSE:
  • Status: 400
  • Body: { "error": { "code": "VALIDATION_ERROR" } }
```
→ Check request body format!

### Find Slow Requests
```
📤 RESPONSE:
  • Status: 200
  • Duration:  2345.67 ms (2345670 µs)
```
→ This is slow! Optimize database query!

### Match Frontend-Backend Issues
1. Frontend console shows: `🔍 API [550e8400] GET /api/v1/billing/plans → 401`
2. Backend terminal shows same call ID with full details
3. Now you can see exactly what was sent and received!

## ⚙️ Configuration

### Backend - Change Logging Behavior

In `backend/main.go`:

```go
// Disable response body logging (save space)
app.Use(middleware.RequestLogger(middleware.RequestLoggerConfig{
    LogRequestBody:  true,
    LogResponseBody: false,  // Don't log responses
    MaxBodySize:     1024,    // Reduce body size limit
    SkipPaths: []string{
        "/health",
        "/api/v1/public/*",  // Skip public endpoints
    },
}))
```

### Frontend - Disable Console Logging

In `frontend/src/lib/api.ts`, remove or comment out:
```go
// Comment this out to stop console logging
// if (import.meta.env?.DEV) {
//   const callID = response.headers.get('X-API-Call-ID');
//   if (callID) {
//     console.log(`🔍 API [${callID.substring(0, 8)}] ...`);
//   }
// }
```

## 🚀 Production Recommendations

For production, use lighter logging:

```go
app.Use(middleware.RequestLogger(middleware.RequestLoggerConfig{
    LogRequestBody:  false,  // Don't log request bodies
    LogResponseBody: false,  // Don't log response bodies
    MaxBodySize:     256,    // Minimal body logging
    SkipPaths: []string{
        "/health",
        "/metrics",
    },
}))
```

Or use JSON mode for log aggregation (ELK, Datadog, etc.):

```go
app.Use(middleware.JSONRequestLogger(middleware.RequestLoggerConfig{
    LogRequestBody:  true,
    LogResponseBody: true,
    MaxBodySize:     2048,
}))
```

## 📁 Files Changed

### Backend
- ✅ `backend/middleware/request_logger.go` (NEW)
- ✅ `backend/main.go` (UPDATED - uses new middleware)
- ✅ `backend/go.mod` (UPDATED - added uuid dependency)

### Frontend
- ✅ `frontend/src/lib/api-call-tracker.ts` (NEW)
- ✅ `frontend/src/lib/api.ts` (UPDATED - integrated tracking)

### Documentation
- ✅ `docs/REQUEST_LOGGING.md` (NEW - comprehensive guide)
- ✅ `docs/LOG_QUICK_REFERENCE.md` (NEW - quick diagnosis guide)
- ✅ `LOGGING_IMPLEMENTATION_SUMMARY.md` (THIS FILE)

## 🐛 Troubleshooting

### Backend logs not showing
1. Restart backend server: `cd backend && go run main.go`
2. Check for compilation errors
3. Verify backend is actually running (`curl http://localhost:3000/health`)

### Console logs not showing in frontend
1. Make sure you're in development mode (`npm run dev`)
2. Check browser console is open (F12)
3. Verify `api.ts` is being used for API calls

### Call ID header missing
1. Check that backend middleware is running
2. Verify response headers in Network tab
3. If using production build, headers should still be present

### Performance issues with logging
1. Disable `LogResponseBody` for production
2. Reduce `MaxBodySize` to 512 or 256 bytes
3. Add more paths to `SkipPaths`
4. Use `JSONRequestLogger` for better performance

## 📊 Example Outputs

### Successful Authentication
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
  • Duration:  45.23 ms (45230 µs)
  • Body:
    {
      "success": true,
      "message": "Login successful",
      "data": {
        "token": "eyJhbGciOiJIUzI1NiIs...",
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

### Missing Token (Common Issue)
```
═══════════════════════════════════════════════════════════
⚠️ API CALL [550e8401] GET /api/v1/billing/plans
═══════════════════════════════════════════════════════════
📥 REQUEST:
  • Time:      00:15:24.789
  • Client IP: 127.0.0.1
  • Method:    GET
  • Path:      /api/v1/billing/plans
  (No Auth header!)

📤 RESPONSE:
  • Status:    401
  • Duration:  0.50 ms (500 µs)
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

### Invalid Token
```
═══════════════════════════════════════════════════════════
⚠️ API CALL [550e8402] GET /api/v1/billing/plans
═══════════════════════════════════════════════════════════
📥 REQUEST:
  • Time:      00:15:25.123
  • Client IP: 127.0.0.1
  • Method:    GET
  • Path:      /api/v1/billing/plans
  • Auth:      Bearer invalid_token_here

📤 RESPONSE:
  • Status:    401
  • Duration:  1.20 ms (1200 µs)
  • Body:
    {
      "success": false,
      "error": {
        "code": "INVALID_TOKEN",
        "message": "Invalid or expired token"
      }
    }
───────────────────────────────────────────────────────────────
```

## 🎉 Benefits

1. **Easier Debugging** - See exactly what's sent and received
2. **Faster Issue Resolution** - Call ID matches frontend to backend
3. **Better Monitoring** - Spot slow requests immediately
4. **Audit Trail** - Every request is logged with timestamp
5. **Development Speed** - No more guessing what went wrong

## 📚 Related Documentation

- `docs/REQUEST_LOGGING.md` - Full feature documentation
- `docs/LOG_QUICK_REFERENCE.md` - Quick diagnosis guide
- `AUTH_FIX_SUMMARY.md` - Authentication fix details
- `TROUBLESHOOTING.md` - General troubleshooting guide

---

**Ready to debug like a pro!** 🚀

Just run your backend and watch the beautiful detailed logs appear in terminal!
