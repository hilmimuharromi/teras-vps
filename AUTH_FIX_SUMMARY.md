# Authentication Fix Summary

## Problem Identified

The login system was not sending the JWT token to the backend properly, resulting in `401 Unauthorized` errors on protected routes like `/api/v1/billing/plans`.

### Root Causes

1. **Inconsistent JWT Secret Usage**
   - `backend/utils/jwt.go` was reading `JWT_SECRET` directly from environment variables with `os.Getenv()`
   - `backend/config/config.go` loads the config once at startup, but JWT utils were not using it
   - This could cause token generation and validation to use different secrets if config changes

2. **Frontend Token Management Issue**
   - `frontend/src/lib/api.ts` only loaded token from localStorage once when the APIClient singleton was instantiated
   - After page reload, if the singleton was already created, it wouldn't reload the token from localStorage
   - This caused requests to be sent without the Authorization header even though token existed in localStorage

3. **No Auth State Initialization**
   - `initAuthState()` function existed in `auth.ts` but was never called on app startup
   - This meant the in-memory auth state was always empty after page reload

## Changes Made

### Backend Changes

#### 1. `backend/utils/jwt.go`
**Before:**
```go
func GenerateToken(userID uint, email string, role string) (string, error) {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        secret = "your-secret-key-change-this" // fallback
    }
    // ...
}

func ValidateToken(tokenString string) (*Claims, error) {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        secret = "your-secret-key-change-this" // fallback
    }
    // ...
}
```

**After:**
```go
// JWTSecret holds the JWT secret key
var JWTSecret string

// InitJWT initializes the JWT secret from config
func InitJWT(cfg *config.Config) {
    JWTSecret = cfg.JWTSecret
}

func GenerateToken(userID uint, email string, role string) (string, error) {
    cfg := config.Load()
    expirationHours := cfg.JWTExpiration
    if expirationHours <= 0 {
        expirationHours = 24 // default fallback
    }
    // ...
    tokenString, err := token.SignedString([]byte(JWTSecret))
    // ...
}

func ValidateToken(tokenString string) (*Claims, error) {
    // ...
    return []byte(JWTSecret), nil
    // ...
}
```

**Why:** Now both token generation and validation use the same `JWTSecret` variable that's initialized from config at startup.

#### 2. `backend/main.go`
**Added:**
```go
// Initialize JWT secret
utils.InitJWT(cfg)
```

**Why:** Ensures JWT secret is initialized once at app startup before any routes are registered.

### Frontend Changes

#### 1. `frontend/src/lib/api.ts`
**Before:**
```typescript
private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const headers: HeadersInit = {
        'Content-Type': 'application/json',
        ...options.headers,
    };

    // Add Authorization header if token exists
    if (this.token) {
        headers['Authorization'] = `Bearer ${this.token}`;
    }
    // ...
}
```

**After:**
```typescript
private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const headers: HeadersInit = {
        'Content-Type': 'application/json',
        ...options.headers,
    };

    // Always get the latest token from localStorage
    const token = typeof window !== "undefined" ? localStorage.getItem("token") : this.token;
    
    // Add Authorization header if token exists
    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(url, {
        ...options,
        headers,
    });

    const data = await response.json();

    // Handle 401 Unauthorized - token expired or invalid
    if (response.status === 401 && typeof window !== "undefined") {
        // Clear invalid token
        localStorage.removeItem("token");
        // Redirect to login
        window.location.href = "/auth/login";
        throw new Error("Authentication required");
    }

    if (!response.ok) {
        throw new Error(data.error?.message || "An error occurred");
    }
    // ...
}
```

**Why:**
- Now always reads the latest token from localStorage on every request
- Automatically redirects to login page when receiving 401 (token expired/invalid)
- Clears invalid tokens to prevent redirect loops

## How to Test

### 1. Start the Backend
```bash
cd backend
go run main.go
```

Expected output:
```
🚀 TerasVPS API server starting on port 3000
```

### 2. Start the Frontend
```bash
cd frontend
npm run dev
```

### 3. Test Login Flow
1. Open browser to `http://localhost:4321/auth/login`
2. Enter valid email and password
3. Click "Sign In"
4. Check browser DevTools → Network tab
5. Verify:
   - Login request returns `200 OK` with token
   - Next request to `/api/v1/auth/me` includes `Authorization: Bearer <token>` header
   - You are redirected to `/dashboard`

### 4. Test Token Persistence
1. After successful login, refresh the page
2. Check browser DevTools → Application → Local Storage
3. Verify `token` key exists with a JWT value
4. Check Network tab - verify subsequent requests include Authorization header

### 5. Test Protected Route Access
1. Navigate to `/dashboard/billing`
2. Check Network tab for request to `/api/v1/billing/plans`
3. Verify:
   - Request includes `Authorization: Bearer <token>` header
   - Response returns `200 OK` (not 401)

### 6. Test Token Expiration Handling
1. Login successfully
2. Manually edit the token in localStorage (DevTools → Application → Local Storage)
3. Add some random characters to invalidate it
4. Refresh the page or navigate to any protected route
5. Verify:
   - You are automatically redirected to `/auth/login`
   - Invalid token is removed from localStorage

## Configuration

Your current `.env` file has:
```env
JWT_SECRET=terasvps-super-secret-change-me-2026
JWT_EXPIRATION=24
```

This is good! The secret is consistent and the expiration is 24 hours.

## Important Notes

1. **JWT Secret Consistency**: The backend now uses a single source of truth for JWT secret (from config). Never change the secret in production without implementing token migration strategy.

2. **Token Storage**: Tokens are stored in browser localStorage. This is vulnerable to XSS attacks. Consider httpOnly cookies for production.

3. **Auto-redirect on 401**: The frontend now automatically redirects to login when receiving 401 errors. This prevents broken UI states.

4. **Token Refresh**: JWT tokens expire after 24 hours (configurable via `JWT_EXPIRATION` env var). Users will need to login again after expiration.

## Next Steps (Recommended)

- [ ] Implement token refresh mechanism to avoid forced logouts
- [ ] Add httpOnly cookie support for better security
- [ ] Implement "Remember Me" functionality
- [ ] Add rate limiting to login endpoint
- [ ] Add 2FA (Two-Factor Authentication)
- [ ] Add password reset functionality
