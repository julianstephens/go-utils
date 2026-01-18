# HTTP Auth Package

The `httputil/auth` package provides JWT token creation, validation, and management with role-based access control and custom claims support. It's designed for secure authentication and authorization in HTTP services.

## Features

- **JWT Token Management**: Create, validate, and parse JWT tokens
- **Role-Based Access Control**: Built-in support for user roles
- **Custom Claims**: Add arbitrary data to tokens
- **User Information**: Convenient handling of username and email
- **Token Validation**: Comprehensive token verification
- **Security**: Configurable token expiration and issuer validation

## Installation

```bash
go get github.com/julianstephens/go-utils/httputil/auth
```

## Usage

### Basic JWT Manager Setup

```go
manager := auth.NewJWTManager("my-secret-key", time.Hour*24, "my-app")
token, _ := manager.GenerateToken("user123", []string{"user", "admin"})
claims, _ := manager.ValidateToken(token)
```

### Refresh token workflow

The package supports a secure refresh token workflow with separate long-lived refresh tokens. Use
`GenerateTokenPair*` helpers to create an access + refresh token pair and `ExchangeRefreshToken` to
exchange a valid refresh token for a new pair.

```go
manager := auth.NewJWTManager("my-secret-key", time.Minute*15, "my-app")
tokenPair, _ := manager.GenerateTokenPairWithUserInfo("user123", "john_doe", "john@example.com", []string{"user"})
newTokenPair, _ := manager.ExchangeRefreshToken(tokenPair.RefreshToken)
```

### HTTP Handler Integration

```go
manager := auth.NewJWTManager("my-secret-key", time.Minute*15, "my-app")
authenticateUser := func(username, password string) (*auth.UserInfo, error) {
    if username == "demo" && password == "password" {
        return &auth.UserInfo{UserID: "user123", Username: "demo", Email: "demo@example.com", Roles: []string{"user"}}, nil
    }
    return nil, errors.New("invalid credentials")
}

http.HandleFunc("/auth/login", auth.AuthenticationHandler(manager, authenticateUser))
http.HandleFunc("/auth/refresh", auth.RefreshTokenHandler(manager))
```

### Cookie-Based Refresh Tokens

```go
tokenPair, _ := manager.GenerateTokenPairWithUserInfo(userID, username, email, roles)
auth.SetRefreshTokenCookie(w, tokenPair.RefreshToken, time.Hour*24*7, true)
json.NewEncoder(w).Encode(map[string]interface{}{"access_token": tokenPair.AccessToken})

// Refresh: exchange cookie for new pair
refreshToken, _ := auth.GetRefreshTokenFromCookie(r)
newTokenPair, _ := manager.ExchangeRefreshToken(refreshToken)
auth.SetRefreshTokenCookie(w, newTokenPair.RefreshToken, time.Hour*24*7, true)
```

### Advanced Token Management

```go
manager := auth.NewJWTManagerWithRefreshConfig(
    "access-secret", time.Minute*15, "my-app",
    time.Hour*24*30, "refresh-secret-key",
)
tokenPair, _ := manager.GenerateTokenPair("user123", []string{"user"})
refreshClaims, _ := manager.ValidateRefreshToken(tokenPair.RefreshToken)
```

### Migration from Legacy RefreshToken Method

```go
// OLD: newAccessToken, _ := manager.RefreshToken(oldAccessToken)
// NEW: Exchange refresh token for new pair
newTokenPair, _ := manager.ExchangeRefreshToken(tokenPair.RefreshToken)
```

### Token with User Information

```go
manager := auth.NewJWTManager("secret", time.Hour*8, "user-service")
token, _ := manager.GenerateTokenWithUserInfo("user456", "john_doe", "john@example.com", []string{"user", "editor"})
claims, _ := manager.ValidateToken(token)
```

### Custom Claims

```go
manager := auth.NewJWTManager("secret", time.Hour*12, "api-service")
customClaims := map[string]interface{}{
    "department": "engineering",
    "level": 5,
    "permissions": []string{"read", "write", "admin"},
}
token, _ := manager.GenerateTokenWithClaims("user789", []string{"admin"}, customClaims)
claims, _ := manager.ValidateToken(token)
```

### Utilities

**Password helpers**: `HashPassword()`, `CheckPasswordHash()`

**Cookie helpers**: `SetRefreshTokenCookie()`, `GetRefreshTokenFromCookie()`, `ClearRefreshTokenCookie()`

## Error Types

Sentinel errors exported: `ErrInvalidToken`, `ErrTokenExpired`, `ErrInvalidClaims`, `ErrInvalidRefreshToken`, `ErrRefreshTokenExpired`


### HTTP Middleware Integration

```go
jwtManager := auth.NewJWTManager("my-secret-key", time.Hour*24, "web-app")

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        tokenString, err := auth.ExtractTokenFromHeader(r.Header.Get("Authorization"))
        if err != nil || jwtManager.ValidateToken(tokenString) != nil {
            http.Error(w, "Invalid or missing token", http.StatusUnauthorized)
            return
        }
        next(w, r)
    }
}

http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
    token, _ := jwtManager.GenerateTokenWithUserInfo("user123", "john_doe", "john@example.com", []string{"user"})
    fmt.Fprintf(w, `{"token": "%s"}`, token)
})
http.HandleFunc("/protected", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Access granted"))
}))
```

### Advanced Token Management

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/julianstephens/go-utils/httputil/auth"
)

func main() {
    // Create manager with shorter expiration for testing
    manager := auth.NewJWTManager("secret", time.Minute*5, "test-app")
    
    userID := "test-user"
    roles := []string{"tester"}
    
    // Generate token
    token, err := manager.GenerateToken(userID, roles)
    if err != nil {
        log.Fatalf("Failed to generate token: %v", err)
    }
    
    fmt.Println("Token generated successfully")
    
    // Validate immediately (should work)
    claims, err := manager.ValidateToken(token)
    if err != nil {
        log.Fatalf("Failed to validate token: %v", err)
    }
    
    fmt.Printf("Token valid - User: %s, Roles: %v\n", claims.UserID, claims.Roles)
    fmt.Printf("Expires at: %v\n", time.Unix(claims.ExpiresAt, 0))
    
    // Check time until expiration
    expirationTime := time.Unix(claims.ExpiresAt, 0)
    timeUntilExpiry := time.Until(expirationTime)
    fmt.Printf("Time until expiry: %v\n", timeUntilExpiry)
    
    // Test with expired token (simulate by creating manager with past time)
    expiredManager := auth.NewJWTManager("secret", -time.Hour, "test-app")
    expiredToken, _ := expiredManager.GenerateToken(userID, roles)
    
    _, err = manager.ValidateToken(expiredToken)
    if err != nil {
        fmt.Printf("Expected error for expired token: %v\n", err)
    }
}
```

## API Reference

### JWTManager

#### Constructor
- `NewJWTManager(secretKey string, tokenDuration time.Duration, issuer string) *JWTManager`
- `NewJWTManagerWithRefreshConfig(secretKey string, tokenDuration time.Duration, issuer string, refreshTokenDuration time.Duration, refreshSecretKey string) *JWTManager`

#### Token Generation
- `GenerateToken(userID string, roles []string) (string, error)` - Generate basic access token
- `GenerateTokenWithUserInfo(userID, username, email string, roles []string) (string, error)` - Generate access token with user info
- `GenerateTokenWithClaims(userID string, roles []string, customClaims map[string]interface{}) (string, error)` - Generate access token with custom claims

#### Token Pair Generation (Access + Refresh)
- `GenerateTokenPair(userID string, roles []string) (*TokenPair, error)` - Generate token pair with basic claims
- `GenerateTokenPairWithUserInfo(userID, username, email string, roles []string) (*TokenPair, error)` - Generate token pair with user info
- `GenerateTokenPairWithClaims(userID string, roles []string, customClaims map[string]interface{}) (*TokenPair, error)` - Generate token pair with custom claims

#### Token Validation
- `ValidateToken(tokenString string) (*Claims, error)` - Validate and parse access token
- `ValidateRefreshToken(refreshTokenString string) (*RefreshClaims, error)` - Validate and parse refresh token
- `ExchangeRefreshToken(refreshTokenString string) (*TokenPair, error)` - Exchange valid refresh token for new token pair

#### Legacy Token Refresh
- `RefreshToken(tokenString string) (string, error)` - Legacy method: refresh access token (deprecated, use ExchangeRefreshToken instead)

### TokenPair
Structure returned when generating token pairs:
```go
type TokenPair struct {
    AccessToken  string `json:"access_token"`  // Short-lived access token
    RefreshToken string `json:"refresh_token"` // Long-lived refresh token  
    TokenType    string `json:"token_type"`    // Always "Bearer"
    ExpiresIn    int64  `json:"expires_in"`    // Access token expiration in seconds
}
```

### Claims Structure

The package defines `Claims` and `RefreshClaims` types. They embed `jwt.RegisteredClaims` and include convenience fields:

- `UserID`, `Username`, `Email`, `Roles`, and `CustomClaims` (map[string]any)
- `RefreshClaims` also includes a `TokenID` string

`Claims` exposes helpers like `HasRole`, `HasAnyRole`, `IsExpired`, and `Expiration()`.

### Utility Functions

- `ExtractTokenFromHeader(authHeader string) (string, error)` - Extract Bearer token from Authorization header

## Error Types

The package exports sentinel errors which callers can check with `errors.Is`:

- `ErrInvalidToken`
- `ErrTokenExpired`
- `ErrInvalidClaims`
- `ErrInvalidRefreshToken`
- `ErrRefreshTokenExpired`

## Best Practices

- Use strong, randomly generated secret keys
- Set short expiration times for access tokens (≤15 minutes)
- Always use HTTPS in production
- Store refresh tokens securely (httpOnly cookies preferred)
- Implement token rotation on refresh
- Use separate secrets for access and refresh tokens when possible
- Implement rate limiting on auth endpoints
- Monitor for anomalous token usage patterns

## Integration

Works well with other go-utils packages:

```go
logger.WithField("user_id", claims.UserID).Info("Authenticated")

// Use with httputil/middleware for complete auth infrastructure
```