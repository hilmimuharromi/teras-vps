# Proxmox Ticket Auto-Refresh Fix

## ❌ Error

```json
{
  "error": {
    "code": "PROXMOX_ERROR",
    "message": "Failed to get VM stats: request failed with status 401: "
  },
  "success": false
}
```

## 🔍 Root Cause

**Proxmox Authentication Ticket Expired**

Proxmox API tickets have a limited lifetime (usually 2 hours). After expiration:
- All API requests return `401 Unauthorized`
- VM stats, start/stop, and other operations fail
- User must restart backend to get fresh ticket

## ✅ Fix

### Auto-Refresh on 401

The Proxmox client now automatically:
1. **Detects 401 error** - Ticket expired
2. **Re-authenticates** - Calls `Login()` to get fresh ticket
3. **Retries request** - Repeats the original request with new ticket
4. **Prevents infinite loops** - Only retries once per request

### Implementation

**Before:**
```go
func (c *Client) Request(method, path string, payload interface{}) ([]byte, error) {
    // Make request
    resp, err := c.HTTPClient.Do(req)
    if resp.StatusCode == 401 {
        return nil, fmt.Errorf("request failed with status 401") // ❌ Fails
    }
    // ...
}
```

**After:**
```go
func (c *Client) Request(method, path string, payload interface{}) ([]byte, error) {
    return c.RequestWithRetry(method, path, payload, true) // ✅ Auto-retry
}

func (c *Client) RequestWithRetry(method, path string, payload interface{}, retry bool) ([]byte, error) {
    // Make request
    resp, err := c.HTTPClient.Do(req)
    
    // Handle 401 - ticket expired
    if resp.StatusCode == 401 && retry {
        // Re-login to get fresh ticket
        if loginErr := c.Login(); loginErr != nil {
            return nil, fmt.Errorf("failed to refresh ticket: %w", loginErr)
        }
        
        // Retry with fresh ticket (retry=false to prevent infinite loop)
        return c.RequestWithRetry(method, path, payload, false)
    }
    
    // ...
}
```

## 🎯 Benefits

1. **No Manual Restarts** - Backend auto-refreshes expired tickets
2. **Transparent to User** - No interruption, stats continue working
3. **Prevents Data Loss** - Operations don't fail due to expired tickets
4. **Safe Retry Logic** - Only retries once to prevent infinite loops

## 📊 Flow Diagram

```
Request → Proxmox API
           ↓
      Status 401?
           ↓
         YES
           ↓
    Re-login (get fresh ticket)
           ↓
    Retry Request (retry=false)
           ↓
      Success! ✅
```

## 🔧 Testing

### Simulate Expired Ticket:
1. Start backend
2. Wait for ticket to expire (or manually clear `c.Ticket`)
3. Request VM stats
4. **Expected:** Auto-refresh ticket and succeed

### Verify in Logs:
```
❌ EXTERNAL API RESPONSE [ext-123]
  Status:    401
  Service:   Proxmox
  
🚀 EXTERNAL API REQUEST [ext-124]  ← Fresh login
  Service:   Proxmox
  Method:    POST
  URL:       .../access/ticket

✅ EXTERNAL API RESPONSE [ext-125]  ← Retried request
  Status:    200
  Service:   Proxmox
```

## 📁 Files Changed

- ✅ `backend/proxmox/client.go` - Added retry logic

## ✨ Result

Proxmox operations now work continuously without manual restarts. Ticket expiration is handled automatically!

### Before:
```
❌ 401 Unauthorized
→ User must restart backend
```

### After:
```
✅ Auto-refresh ticket → Retry → Success!
```
