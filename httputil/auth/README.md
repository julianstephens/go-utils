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
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/julianstephens/go-utils/httputil/auth"
)

func main() {
    // Create JWT manager with secret key, expiration, and issuer
    manager := auth.NewJWTManager("my-secret-key", time.Hour*24, "my-app")
    
    // Generate basic token
    userID := "user123"
    roles := []string{"user", "admin"}
    
    token, err := manager.GenerateToken(userID, roles)
    if err != nil {
        log.Fatalf("Failed to generate token: %v", err)
    }
    
    fmt.Printf("Generated token: %s\n", token)
    
    // Validate and parse token
    claims, err := manager.ValidateToken(token)
    if err != nil {
        log.Fatalf("Failed to validate token: %v", err)
    }
    
    fmt.Printf("User ID: %s\n", claims.UserID)
    fmt.Printf("Roles: %v\n", claims.Roles)
    fmt.Printf("Issued at: %v\n", time.Unix(claims.IssuedAt, 0))
    fmt.Printf("Expires at: %v\n", time.Unix(claims.ExpiresAt, 0))
}
```

### Refresh Token Workflow

The package supports a secure refresh token workflow with separate long-lived refresh tokens:

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/julianstephens/go-utils/httputil/auth"
)

func main() {
    // Create JWT manager - access tokens expire in 15 minutes, refresh tokens in 7 days
    manager := auth.NewJWTManager("my-secret-key", time.Minute*15, "my-app")
    
    // Generate token pair (access + refresh token)
    userID := "user123"
    username := "john_doe"
    email := "john@example.com"
    roles := []string{"user", "admin"}
    
    tokenPair, err := manager.GenerateTokenPairWithUserInfo(userID, username, email, roles)
    if err != nil {
        log.Fatalf("Failed to generate token pair: %v", err)
    }
    
    fmt.Printf("Access Token: %s\n", tokenPair.AccessToken)
    fmt.Printf("Refresh Token: %s\n", tokenPair.RefreshToken)
    fmt.Printf("Token Type: %s\n", tokenPair.TokenType)
    fmt.Printf("Expires In: %d seconds\n", tokenPair.ExpiresIn)
    
    // Validate access token
    claims, err := manager.ValidateToken(tokenPair.AccessToken)
    if err != nil {
        log.Fatalf("Failed to validate access token: %v", err)
    }
    fmt.Printf("Access token valid for user: %s\n", claims.UserID)
    
    // When access token expires, exchange refresh token for new tokens
    newTokenPair, err := manager.ExchangeRefreshToken(tokenPair.RefreshToken)
    if err != nil {
        log.Fatalf("Failed to exchange refresh token: %v", err)
    }
    
    fmt.Printf("New Access Token: %s\n", newTokenPair.AccessToken)
    fmt.Printf("New Refresh Token: %s\n", newTokenPair.RefreshToken)
}
```

### HTTP Handler Integration

The package provides ready-to-use HTTP handlers for common authentication workflows:

```go
package main

import (
    "errors"  
    "log"
    "net/http"
    "time"
    
    "github.com/julianstephens/go-utils/httputil/auth"
)

func main() {
    // Setup JWT manager
    manager := auth.NewJWTManager("my-secret-key", time.Minute*15, "my-app")
    
    // Mock user authentication function
    authenticateUser := func(username, password string) (*auth.UserInfo, error) {
        // Replace with your user authentication logic
        if username == "demo" && password == "password" {
            return &auth.UserInfo{
                UserID:   "user123",
                Username: "demo",
                Email:    "demo@example.com",
                Roles:    []string{"user"},
            }, nil
        }
        return nil, errors.New("invalid credentials")
    }
    
    // Setup handlers
    http.HandleFunc("/auth/login", auth.AuthenticationHandler(manager, authenticateUser))
    http.HandleFunc("/auth/refresh", auth.RefreshTokenHandler(manager))
    
    // Example protected route using middleware
    http.Handle("/protected", jwtMiddleware(manager)(http.HandlerFunc(protectedHandler)))
    
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Access granted to protected resource"))
}

func jwtMiddleware(manager *auth.JWTManager) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            tokenString, err := auth.ExtractTokenFromHeader(authHeader)
            if err != nil {
                http.Error(w, "Authorization header required", http.StatusUnauthorized)
                return
            }

            _, err = manager.ValidateToken(tokenString)
            if err != nil {
                http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

### Cookie-Based Refresh Tokens

For web applications, you can use httpOnly cookies for secure refresh token storage:

```go
func loginWithCookies(manager *auth.JWTManager) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // ... authentication logic ...
        
        tokenPair, err := manager.GenerateTokenPairWithUserInfo(userID, username, email, roles)
        if err != nil {
            http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
            return
        }
        
        // Set refresh token as httpOnly cookie
        auth.SetRefreshTokenCookie(w, tokenPair.RefreshToken, time.Hour*24*7, true)
        
        // Return only access token in response
        response := map[string]interface{}{
            "access_token": tokenPair.AccessToken,
            "token_type":   tokenPair.TokenType,
            "expires_in":   tokenPair.ExpiresIn,
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    }
}

func refreshFromCookie(manager *auth.JWTManager) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Get refresh token from cookie
        refreshToken, err := auth.GetRefreshTokenFromCookie(r)
        if err != nil {
            http.Error(w, "Refresh token not found", http.StatusUnauthorized)
            return
        }
        
        // Exchange for new tokens
        tokenPair, err := manager.ExchangeRefreshToken(refreshToken)
        if err != nil {
            auth.ClearRefreshTokenCookie(w) // Clear invalid cookie
            http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
            return
        }
        
        // Set new refresh token cookie
        auth.SetRefreshTokenCookie(w, tokenPair.RefreshToken, time.Hour*24*7, true)
        
        // Return new access token
        response := map[string]interface{}{
            "access_token": tokenPair.AccessToken,
            "token_type":   tokenPair.TokenType,
            "expires_in":   tokenPair.ExpiresIn,
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    }
}
```

### Advanced Token Management

```go
// Custom refresh token duration and separate secret
manager := auth.NewJWTManagerWithRefreshConfig(
    "access-secret",           // Access token secret
    time.Minute*15,           // Access token duration (15 minutes)
    "my-app",                 // Issuer
    time.Hour*24*30,          // Refresh token duration (30 days)  
    "refresh-secret-key",     // Separate refresh token secret
)

tokenPair, err := manager.GenerateTokenPair("user123", []string{"user"})
if err != nil {
    log.Fatalf("Failed to generate tokens: %v", err)
}

// Token introspection
refreshClaims, err := manager.ValidateRefreshToken(tokenPair.RefreshToken)
if err != nil {
    log.Fatalf("Invalid refresh token: %v", err)
}
fmt.Printf("Refresh token ID: %s\n", refreshClaims.TokenID)
fmt.Printf("Token expires: %v\n", refreshClaims.ExpiresAt.Time)
```

### Migration from Legacy RefreshToken Method

If you're currently using the legacy `RefreshToken` method, here's how to migrate:

```go
// OLD: Legacy refresh (recreates access token from existing token)
oldAccessToken := "existing-access-token"
newAccessToken, err := manager.RefreshToken(oldAccessToken)

// NEW: Token pair workflow (exchange refresh token for new pair)
// 1. Generate token pair initially
tokenPair, err := manager.GenerateTokenPair("user123", []string{"user"})

// 2. Store refresh token securely (cookie, database, etc.)
// 3. Exchange refresh token when access token expires
newTokenPair, err := manager.ExchangeRefreshToken(tokenPair.RefreshToken)
```
    fmt.Printf("Issued at: %v\n", time.Unix(claims.IssuedAt, 0))
    fmt.Printf("Expires at: %v\n", time.Unix(claims.ExpiresAt, 0))
}
```

### Token with User Information

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/julianstephens/go-utils/httputil/auth"
)

func main() {
    manager := auth.NewJWTManager("secret", time.Hour*8, "user-service")
    
    // Generate token with user information
    userID := "user456"
    username := "john_doe"
    email := "john@example.com"
    roles := []string{"user", "editor"}
    
    token, err := manager.GenerateTokenWithUserInfo(userID, username, email, roles)
    if err != nil {
        log.Fatalf("Failed to generate token: %v", err)
    }
    
    fmt.Printf("Generated token with user info: %s\n", token)
    
    // Parse and access user information
    claims, err := manager.ValidateToken(token)
    if err != nil {
        log.Fatalf("Failed to validate token: %v", err)
    }
    
    fmt.Printf("User ID: %s\n", claims.UserID)
    fmt.Printf("Username: %s\n", claims.Username)
    fmt.Printf("Email: %s\n", claims.Email)
    fmt.Printf("Roles: %v\n", claims.Roles)
}
```

### Custom Claims

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/julianstephens/go-utils/httputil/auth"
)

func main() {
    manager := auth.NewJWTManager("secret", time.Hour*12, "api-service")
    
    userID := "user789"
    roles := []string{"admin", "manager"}
    
    // Create custom claims with additional data
    customClaims := map[string]interface{}{
        "username":    "admin_user",
        "email":       "admin@company.com",
        "department":  "engineering",
        "level":       5,
        "is_manager":  true,
        "permissions": []string{"read", "write", "delete", "admin"},
        "metadata": map[string]string{
            "region": "us-west",
            "tenant": "acme-corp",
        },
        "last_login": time.Now().Unix(),
    }
    
    // Generate token with custom claims
    token, err := manager.GenerateTokenWithClaims(userID, roles, customClaims)
    if err != nil {
        log.Fatalf("Failed to generate token: %v", err)
    }
    
    fmt.Printf("Token with custom claims generated\n")
    
    // Validate and access custom claims
    claims, err := manager.ValidateToken(token)
    if err != nil {
        log.Fatalf("Failed to validate token: %v", err)
    }
    
    fmt.Printf("User ID: %s\n", claims.UserID)
    fmt.Printf("Roles: %v\n", claims.Roles)
    
    // Access custom claims
    if department, ok := claims.CustomClaims["department"].(string); ok {
        fmt.Printf("Department: %s\n", department)
    }
    
    if level, ok := claims.CustomClaims["level"].(float64); ok {
        fmt.Printf("Level: %.0f\n", level)
    }
    
    if permissions, ok := claims.CustomClaims["permissions"].([]interface{}); ok {
        fmt.Printf("Permissions: %v\n", permissions)
    }
    
    if metadata, ok := claims.CustomClaims["metadata"].(map[string]interface{}); ok {
        fmt.Printf("Metadata: %v\n", metadata)
    }
}
```

### HTTP Middleware Integration

```go
package main

import (
    "fmt"
    "net/http"
    "strings"
    "time"
    
    "github.com/julianstephens/go-utils/httputil/auth"
)

var jwtManager *auth.JWTManager

func init() {
    jwtManager = auth.NewJWTManager("my-secret-key", time.Hour*24, "web-app")
}

// Authentication middleware
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Extract token from Authorization header
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }
        
        // Check for Bearer token format
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            http.Error(w, "Bearer token required", http.StatusUnauthorized)
            return
        }
        
        // Validate token
        claims, err := jwtManager.ValidateToken(tokenString)
        if err != nil {
            http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
            return
        }
        
        // Add claims to request context (you'd implement this)
        // ctx := context.WithValue(r.Context(), "claims", claims)
        // r = r.WithContext(ctx)
        
        fmt.Printf("Authenticated user: %s with roles: %v\n", claims.UserID, claims.Roles)
        next(w, r)
    }
}

// Role-based authorization middleware
func requireRole(requiredRole string) func(http.HandlerFunc) http.HandlerFunc {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return authMiddleware(func(w http.ResponseWriter, r *http.Request) {
            // In real implementation, you'd get claims from context
            // For demo, we'll validate token again
            authHeader := r.Header.Get("Authorization")
            tokenString := strings.TrimPrefix(authHeader, "Bearer ")
            
            claims, err := jwtManager.ValidateToken(tokenString)
            if err != nil {
                http.Error(w, "Invalid token", http.StatusUnauthorized)
                return
            }
            
            // Check if user has required role
            hasRole := false
            for _, role := range claims.Roles {
                if role == requiredRole {
                    hasRole = true
                    break
                }
            }
            
            if !hasRole {
                http.Error(w, fmt.Sprintf("Role '%s' required", requiredRole), http.StatusForbidden)
                return
            }
            
            next(w, r)
        })
    }
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    // In real app, you'd validate credentials
    userID := "user123"
    username := "john_doe"
    email := "john@example.com"
    roles := []string{"user", "admin"}
    
    token, err := jwtManager.GenerateTokenWithUserInfo(userID, username, email, roles)
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"token": "%s"}`, token)
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Access granted to protected resource"))
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Admin access granted"))
}

func main() {
    http.HandleFunc("/login", loginHandler)
    http.HandleFunc("/protected", authMiddleware(protectedHandler))
    http.HandleFunc("/admin", requireRole("admin")(adminHandler))
    
    fmt.Println("Server starting on :8080")
    fmt.Println("Endpoints:")
    fmt.Println("  POST /login - Get JWT token")
    fmt.Println("  GET /protected - Requires any valid token")
    fmt.Println("  GET /admin - Requires 'admin' role")
    fmt.Println()
    fmt.Println("Usage:")
    fmt.Println("  1. POST to /login to get token")
    fmt.Println("  2. Use token in Authorization: Bearer <token> header")
    
    http.ListenAndServe(":8080", nil)
}
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
    
    // Parse token without validation (useful for debugging)
    unverifiedClaims, err := manager.ParseToken(token)
    if err != nil {
        log.Fatalf("Failed to parse token: %v", err)
    }
    
    fmt.Printf("Unverified claims - Issuer: %s\n", unverifiedClaims.Issuer)
    
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
- `NewJWTManager(secretKey string, expiration time.Duration, issuer string) *JWTManager`
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

```go
type Claims struct {
    UserID       string                 `json:"user_id"`
    Username     string                 `json:"username,omitempty"`
    Email        string                 `json:"email,omitempty"`
    Roles        []string               `json:"roles"`
    CustomClaims map[string]interface{} `json:"custom_claims,omitempty"`
    
    // Standard JWT claims
    Issuer    string `json:"iss"`
    Subject   string `json:"sub"`
    Audience  string `json:"aud,omitempty"`
    ExpiresAt int64  `json:"exp"`
    IssuedAt  int64  `json:"iat"`
    NotBefore int64  `json:"nbf"`
}
```

### Utility Functions

- `ExtractTokenFromHeader(authHeader string) (string, error)` - Extract Bearer token from Authorization header
- `GenerateSecretKey() (string, error)` - Generate cryptographically secure secret key

## Error Types

The package provides specific error types for different validation failures:

- **Invalid token format**: Malformed JWT structure
- **Token expired**: Token past expiration time
- **Invalid signature**: Token signature verification failed
- **Invalid issuer**: Token issuer doesn't match expected value
- **Missing required claims**: Required fields missing from token

## Security Considerations

1. **Secret Key Management**: Use a strong, randomly generated secret key
2. **Token Expiration**: Set appropriate expiration times (shorter for sensitive operations)  
3. **HTTPS Only**: Always use HTTPS in production to protect tokens in transit
4. **Token Storage**: Store tokens securely on the client side
5. **Refresh Token Security**: 
   - Use separate secrets for access and refresh tokens when possible
   - Store refresh tokens securely (httpOnly cookies recommended for web apps)
   - Implement refresh token rotation (each exchange issues new refresh token)
   - Revoke refresh tokens when user logs out or suspicious activity detected
6. **Rate Limiting**: Implement rate limiting on authentication endpoints
7. **Logging**: Don't log full tokens, only token IDs or user IDs

### Refresh Token Best Practices

1. **Short Access Token Lifetime**: Keep access tokens short-lived (15 minutes or less)
2. **Longer Refresh Token Lifetime**: Refresh tokens can be longer-lived (days to weeks)
3. **Token Rotation**: Always issue new refresh tokens when exchanging (automatic in this implementation)
4. **Secure Storage**: Store refresh tokens in httpOnly cookies or secure storage
5. **Revocation**: Implement refresh token blacklisting/revocation for security incidents
6. **Separate Secrets**: Use different secrets for access and refresh tokens
7. **Monitoring**: Monitor refresh token usage patterns for anomalies

## Best Practices

1. **Use short expiration times** for sensitive operations
2. **Include minimal necessary information** in tokens
3. **Validate tokens on every request** to protected resources
4. **Use HTTPS** to prevent token interception
5. **Implement proper error handling** for authentication failures
6. **Log authentication events** for security monitoring
7. **Use environment variables** for secret keys
8. **Implement token refresh** for better user experience

## Integration

Works well with other go-utils packages:

```go
// Use with logger for authentication logging
logger.WithFields(map[string]interface{}{
    "user_id": claims.UserID,
    "roles":   claims.Roles,
}).Info("User authenticated")

// Use with cliutil for CLI authentication
if cliutil.HasFlag(os.Args, "--token") {
    token := cliutil.GetFlagValue(os.Args, "--token", "")
    claims, err := jwtManager.ValidateToken(token)
}
```