# 🔐 Troubleshooting Auth Issues

## Current Status
✅ **Backend JWT is working correctly** - Token generation and validation fixed
✅ **Frontend API client fixed** - Now reads token from localStorage on every request

## Quick Diagnostic

### 1. Open Debug Page
Go to: `http://localhost:4321/auth/debug`

This page will show:
- Whether token exists in localStorage
- JWT token details (expiry, user info)
- Test buttons to verify login and API calls

### 2. Check Browser Console
Open DevTools (F12) and check:

**Console Tab:**
- Look for any JavaScript errors
- Check if there are CORS errors

**Network Tab:**
- Filter by "billing" or "plans"
- Click on the request to `/api/v1/billing/plans`
- Check **Request Headers**:
  ```
  Authorization: Bearer eyJhbGci...
  ```
- If Authorization header is **missing** → Token not being sent
- If Authorization header is **present** but still 401 → Token invalid/expired

**Application Tab → Local Storage:**
- Check if `token` key exists
- If missing → Login didn't save token
- If exists → Check if token looks valid (should be 3 parts separated by dots)

### 3. Common Issues & Solutions

#### Issue: Token not in localStorage
**Symptom:** No `token` key in localStorage
**Solution:**
1. Go to login page
2. Login with credentials
3. Check if token appears after login

#### Issue: Token exists but still 401
**Symptom:** Token in localStorage but API returns 401
**Possible causes:**
1. **Token expired** - Check expiry on debug page
2. **Token corrupted** - Clear token and login again
3. **Backend restarted** - Old token signed with different secret

**Solution:** Clear token and login again

#### Issue: CORS errors
**Symptom:** Request blocked by CORS
**Solution:** Check that backend CORS is configured for your frontend URL in `main.go`:
```go
app.Use(cors.New(cors.Config{
    AllowOrigins:     cfg.FrontendURL,  // Should be http://localhost:4321
    AllowCredentials: true,
}))
```

#### Issue: Frontend not proxying to backend
**Symptom:** Requests going to wrong port
**Check:** 
1. `astro.config.mjs` has proxy setup for `/api` → `http://localhost:3000`
2. Backend is running on port 3000
3. Frontend dev server is running on port 4321

## Step-by-Step Fix

### Nuclear Option (Guaranteed to Work)
1. **Clear everything:**
   - Open DevTools → Application → Local Storage
   - Delete all keys
   - Or click "Clear Token" on debug page

2. **Restart both servers:**
   ```bash
   # Kill backend
   pkill -f "teras-vps-api"
   pkill -f "go run"
   
   # Kill frontend
   pkill -f "astro"
   ```

3. **Start backend:**
   ```bash
   cd backend
   go run main.go
   ```

4. **Start frontend (new terminal):**
   ```bash
   cd frontend
   npm run dev
   ```

5. **Login fresh:**
   - Go to `http://localhost:4321/auth/login`
   - Login with your credentials
   - Check if you can access `/dashboard/billing`

### Verify It Works
1. After login, check Network tab
2. Look for request to `/api/v1/auth/me`
3. Should see `Authorization: Bearer <token>` in request headers
4. Navigate to `/dashboard/billing`
5. Should see request to `/api/v1/billing/plans` with token
6. Should get `200 OK` response (not 401)

## Test Credentials
Use these test credentials (auto-created on backend start):
- **Email:** test@test.com
- **Password:** test12345

## Still Not Working?

### 1. Check Backend Logs
Look at terminal where backend is running. Should see:
```
🚀 TerasVPS API server starting on port 3000
```

And for each request:
```
2026/04/08 00:03:59 GET /api/v1/billing/plans
```

### 2. Manual API Test
Test directly with curl:
```bash
# Login
curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"test12345"}'

# Use token from response
TOKEN="<token-from-login>"

# Test billing plans
curl http://localhost:3000/api/v1/billing/plans \
  -H "Authorization: Bearer $TOKEN"
```

If this works but frontend doesn't → Issue is in frontend token handling

### 3. Check JWT Secret Consistency
```bash
# Check backend .env
cat backend/.env | grep JWT_SECRET

# Should be: terasvps-super-secret-change-me-2026
```

## Files Changed
- `backend/utils/jwt.go` - Fixed JWT secret handling
- `backend/main.go` - Added JWT initialization
- `frontend/src/lib/api.ts` - Fixed token reading from localStorage
- `frontend/src/pages/auth/debug.astro` - Debug page (new)
