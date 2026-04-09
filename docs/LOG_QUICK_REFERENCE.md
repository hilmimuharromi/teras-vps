# Quick Reference: Reading API Logs

## What You'll See Now

Every API request will show detailed information in your backend terminal:

### Example: Successful Login
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
        "token": "eyJhbGci...",
        "user": {
          "id": 5,
          "email": "test@test.com"
        }
      }
    }
───────────────────────────────────────────────────────────────
```

### Example: Missing Token (401)
```
═══════════════════════════════════════════════════════════
⚠️ API CALL [550e8401] GET /api/v1/billing/plans
═══════════════════════════════════════════════════════════
📥 REQUEST:
  • Time:      00:15:24.789
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

## Quick Diagnosis Guide

| What You See | What It Means | Action |
|--------------|---------------|--------|
| ✅ Green + 200 | Success | Nothing, all good! |
| ⚠️ Yellow + 401 | No token or invalid token | Check if frontend sends Authorization header |
| ⚠️ Yellow + 400 | Bad request | Check request body for missing/invalid fields |
| ⚠️ Yellow + 403 | Forbidden | User doesn't have permission |
| ❌ Red + 500 | Server error | Check error message in logs |
| No "Auth" line in request | Token not sent | Fix frontend to send token |
| Duration > 1000ms | Slow request | Check database queries or external APIs |

## Finding the Call ID

In browser DevTools → Network tab → Click any request:
```
Response Headers:
  X-API-Call-ID: 550e8400-e29b-41d4-a716-446655440000
```

Then search backend logs for `550e8400` to find that exact request!

## Common Patterns

### Token Being Sent Correctly
```
📥 REQUEST:
  • Auth: Bearer eyJhbGci...hlP8
```

### Token NOT Being Sent
```
📥 REQUEST:
  (No Auth line here)
```

### Response Contains Data
```
📤 RESPONSE:
  • Status: 200
  • Body:
    {
      "success": true,
      "data": {...}
    }
```

### Response is Error
```
📤 RESPONSE:
  • Status: 401
  • Body:
    {
      "success": false,
      "error": {
        "code": "MISSING_TOKEN"
      }
    }
```

## Tips

1. **Watch the Auth line** - Most important for debugging auth issues!
2. **Check Duration** - Helps spot performance problems
3. **Look at Body** - Confirms what was actually sent/received
4. **Use Call ID** - Match frontend errors with backend logs
