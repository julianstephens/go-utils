# HTTP Middleware Package

The `httputil/middleware` package provides common, reusable HTTP middleware for logging, recovery, CORS, and request ID injection. These middleware components help build robust and maintainable HTTP services.

## Features

- **Request Logging**: Comprehensive request/response logging
- **Panic Recovery**: Graceful panic recovery with error logging
- **CORS Support**: Cross-Origin Resource Sharing handling
- **Request ID**: Unique request identification and tracing
- **JWT Authentication**: Token validation and role-based access control
- **Configurable**: Flexible configuration options for all middleware

## Installation

```bash
go get github.com/julianstephens/go-utils/httputil/middleware
```

## Usage

### Basic Middleware Setup

```go
package main

import (
    "log"
    "net/http"
    "os"
    
    "github.com/gorilla/mux"
    "github.com/julianstephens/go-utils/httputil/middleware"
)

func main() {
    // Create loggers
    logger := log.New(os.Stdout, "[HTTP] ", log.LstdFlags)
    errorLogger := log.New(os.Stderr, "[ERROR] ", log.LstdFlags)
    
    // Create router
    router := mux.NewRouter()
    
    // Apply middleware in the correct order
    router.Use(middleware.RequestID())                          // First - adds ID to all requests
    router.Use(middleware.Logging(logger))                      // Second - logs with request ID
    router.Use(middleware.Recovery(errorLogger))                // Third - catches panics
    router.Use(middleware.CORS(middleware.DefaultCORSConfig())) // Fourth - handles CORS
    
    // Add routes
    router.HandleFunc("/api/health", healthHandler).Methods("GET")
    router.HandleFunc("/api/users", getUsersHandler).Methods("GET")
    
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    // Get request ID from context
    requestID := middleware.GetRequestID(r.Context())
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status": "ok", "request_id": "` + requestID + `"}`))
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"users": []}`))
}
```

### JWT Authentication Middleware

The JWT authentication middleware provides token validation and role-based access control using the `httputil/auth` package.

#### Basic JWT Authentication

```go
package main

import (
    "log"
    "net/http"
    "time"
    
    "github.com/gorilla/mux"
    "github.com/julianstephens/go-utils/httputil/auth"
    "github.com/julianstephens/go-utils/httputil/middleware"
)

func main() {
    // Create JWT manager
    jwtManager := auth.NewJWTManager("your-secret-key", time.Hour*24, "your-app")
    
    router := mux.NewRouter()
    
    // Public routes (no authentication required)
    router.HandleFunc("/api/login", loginHandler(jwtManager)).Methods("POST")
    router.HandleFunc("/api/public", publicHandler).Methods("GET")
    
    // Protected routes (authentication required)
    protected := router.PathPrefix("/api/protected").Subrouter()
    protected.Use(middleware.JWTAuth(jwtManager))
    protected.HandleFunc("/profile", profileHandler).Methods("GET")
    protected.HandleFunc("/data", dataHandler).Methods("GET")
    
    log.Println("Server with JWT auth starting on :8080")
    http.ListenAndServe(":8080", router)
}

func loginHandler(jwtManager *auth.JWTManager) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // In real application, validate credentials here
        userID := "user123"
        roles := []string{"user"}
        
        token, err := jwtManager.GenerateToken(userID, roles)
        if err != nil {
            http.Error(w, "Failed to generate token", http.StatusInternalServerError)
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"token": "` + token + `"}`))
    }
}

func publicHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"message": "This is a public endpoint"}`))
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
    // Get claims from context (injected by JWTAuth middleware)
    claims, ok := middleware.GetClaims(r.Context())
    if !ok {
        http.Error(w, "No claims found", http.StatusUnauthorized)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    response := `{
        "user_id": "` + claims.UserID + `",
        "username": "` + claims.Username + `",
        "email": "` + claims.Email + `",
        "roles": ` + `["` + claims.Roles[0] + `"]` + `
    }`
    w.Write([]byte(response))
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
    claims, _ := middleware.GetClaims(r.Context())
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"data": "secret data for ` + claims.UserID + `"}`))
}
```

#### Role-Based Access Control

```go
package main

import (
    "log"
    "net/http"
    "time"
    
    "github.com/gorilla/mux"
    "github.com/julianstephens/go-utils/httputil/auth"
    "github.com/julianstephens/go-utils/httputil/middleware"
)

func main() {
    jwtManager := auth.NewJWTManager("your-secret-key", time.Hour*24, "your-app")
    
    router := mux.NewRouter()
    
    // Public endpoints
    router.HandleFunc("/api/login", loginHandler(jwtManager)).Methods("POST")
    
    // User endpoints (requires 'user' role)
    userRoutes := router.PathPrefix("/api/user").Subrouter()
    userRoutes.Use(middleware.RequireRoles(jwtManager, "user"))
    userRoutes.HandleFunc("/profile", userProfileHandler).Methods("GET")
    userRoutes.HandleFunc("/settings", userSettingsHandler).Methods("GET", "PUT")
    
    // Admin endpoints (requires 'admin' role)
    adminRoutes := router.PathPrefix("/api/admin").Subrouter()
    adminRoutes.Use(middleware.RequireRoles(jwtManager, "admin"))
    adminRoutes.HandleFunc("/users", adminUsersHandler).Methods("GET")
    adminRoutes.HandleFunc("/system", adminSystemHandler).Methods("GET")
    
    // Manager endpoints (requires 'manager' or 'admin' role)
    managerRoutes := router.PathPrefix("/api/manager").Subrouter()
    managerRoutes.Use(middleware.RequireRoles(jwtManager, "manager", "admin"))
    managerRoutes.HandleFunc("/reports", managerReportsHandler).Methods("GET")
    managerRoutes.HandleFunc("/team", managerTeamHandler).Methods("GET")
    
    log.Println("Server with role-based auth starting on :8080")
    log.Println("Login with POST /api/login to get tokens with different roles")
    http.ListenAndServe(":8080", router)
}

func loginHandler(jwtManager *auth.JWTManager) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Parse request to get user type
        userType := r.URL.Query().Get("type") // ?type=admin|manager|user
        
        var userID string
        var roles []string
        
        switch userType {
        case "admin":
            userID = "admin1"
            roles = []string{"admin", "manager", "user"}
        case "manager":
            userID = "manager1"
            roles = []string{"manager", "user"}
        default:
            userID = "user1"
            roles = []string{"user"}
        }
        
        token, err := jwtManager.GenerateToken(userID, roles)
        if err != nil {
            http.Error(w, "Failed to generate token", http.StatusInternalServerError)
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        response := `{
            "token": "` + token + `",
            "user_id": "` + userID + `",
            "roles": ["` + roles[0] + `"]
        }`
        w.Write([]byte(response))
    }
}

func userProfileHandler(w http.ResponseWriter, r *http.Request) {
    claims, _ := middleware.GetClaims(r.Context())
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"message": "User profile for ` + claims.UserID + `"}`))
}

func userSettingsHandler(w http.ResponseWriter, r *http.Request) {
    claims, _ := middleware.GetClaims(r.Context())
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"message": "User settings for ` + claims.UserID + `"}`))
}

func adminUsersHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"message": "Admin-only user management"}`))
}

func adminSystemHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"message": "Admin-only system management"}`))
}

func managerReportsHandler(w http.ResponseWriter, r *http.Request) {
    claims, _ := middleware.GetClaims(r.Context())
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"message": "Manager reports for ` + claims.UserID + `"}`))
}

func managerTeamHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"message": "Team management"}`))
}
```

#### JWT with Full Middleware Stack

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "time"
    
    "github.com/gorilla/mux"
    "github.com/julianstephens/go-utils/httputil/auth"
    "github.com/julianstephens/go-utils/httputil/middleware"
)

func main() {
    // Setup loggers
    accessLogger := log.New(os.Stdout, "[ACCESS] ", log.LstdFlags)
    errorLogger := log.New(os.Stderr, "[ERROR] ", log.LstdFlags)
    
    // Setup JWT manager
    jwtManager := auth.NewJWTManager("your-secret-key", time.Hour*24, "your-app")
    
    router := mux.NewRouter()
    
    // Apply base middleware stack
    router.Use(middleware.RequestID())
    router.Use(middleware.Logging(accessLogger))
    router.Use(middleware.Recovery(errorLogger))
    router.Use(middleware.CORS(middleware.DefaultCORSConfig()))
    
    // Public API
    api := router.PathPrefix("/api/v1").Subrouter()
    api.HandleFunc("/health", healthHandler).Methods("GET")
    api.HandleFunc("/login", loginHandler(jwtManager)).Methods("POST")
    
    // Protected API with authentication
    protected := api.PathPrefix("/protected").Subrouter()
    protected.Use(middleware.JWTAuth(jwtManager))
    protected.HandleFunc("/profile", profileHandler).Methods("GET")
    
    // Admin API with role requirement
    admin := api.PathPrefix("/admin").Subrouter()
    admin.Use(middleware.RequireRoles(jwtManager, "admin"))
    admin.HandleFunc("/users", adminHandler).Methods("GET")
    
    log.Println("Full-stack JWT API server starting on :8080")
    log.Println("Endpoints:")
    log.Println("  POST /api/v1/login       - Get JWT token")
    log.Println("  GET  /api/v1/health      - Public health check")
    log.Println("  GET  /api/v1/protected/profile - Requires valid JWT")
    log.Println("  GET  /api/v1/admin/users - Requires admin role")
    
    http.ListenAndServe(":8080", router)
}

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type APIResponse struct {
    Success   bool        `json:"success"`
    Data      interface{} `json:"data,omitempty"`
    Error     string      `json:"error,omitempty"`
    RequestID string      `json:"request_id"`
}

func respondJSON(w http.ResponseWriter, r *http.Request, status int, data interface{}, err string) {
    requestID := middleware.GetRequestID(r.Context())
    
    response := APIResponse{
        Success:   status < 400,
        RequestID: requestID,
    }
    
    if err != "" {
        response.Error = err
    } else {
        response.Data = data
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("X-Request-ID", requestID)
    w.WriteHeader(status)
    
    json.NewEncoder(w).Encode(response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    respondJSON(w, r, http.StatusOK, map[string]string{
        "status": "healthy",
        "version": "1.0.0",
    }, "")
}

func loginHandler(jwtManager *auth.JWTManager) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req LoginRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            respondJSON(w, r, http.StatusBadRequest, nil, "Invalid JSON")
            return
        }
        
        // In real app, validate credentials against database
        if req.Username == "" || req.Password == "" {
            respondJSON(w, r, http.StatusBadRequest, nil, "Username and password required")
            return
        }
        
        // Mock authentication
        var roles []string
        if req.Username == "admin" {
            roles = []string{"admin", "user"}
        } else {
            roles = []string{"user"}
        }
        
        token, err := jwtManager.GenerateTokenWithUserInfo(
            req.Username, 
            req.Username, 
            req.Username+"@example.com", 
            roles,
        )
        if err != nil {
            respondJSON(w, r, http.StatusInternalServerError, nil, "Failed to generate token")
            return
        }
        
        respondJSON(w, r, http.StatusOK, map[string]interface{}{
            "token": token,
            "user": map[string]interface{}{
                "username": req.Username,
                "email":    req.Username + "@example.com",
                "roles":    roles,
            },
        }, "")
    }
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
    claims, ok := middleware.GetClaims(r.Context())
    if !ok {
        respondJSON(w, r, http.StatusUnauthorized, nil, "No claims found")
        return
    }
    
    respondJSON(w, r, http.StatusOK, map[string]interface{}{
        "user_id":  claims.UserID,
        "username": claims.Username,
        "email":    claims.Email,
        "roles":    claims.Roles,
    }, "")
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
    claims, _ := middleware.GetClaims(r.Context())
    
    users := []map[string]interface{}{
        {"id": 1, "username": "user1", "email": "user1@example.com", "roles": []string{"user"}},
        {"id": 2, "username": "admin", "email": "admin@example.com", "roles": []string{"admin", "user"}},
    }
    
    respondJSON(w, r, http.StatusOK, map[string]interface{}{
        "users":      users,
        "accessed_by": claims.UserID,
    }, "")
}
```

### Request ID Middleware

```go
package main

import (
    "fmt"
    "net/http"
    
    "github.com/gorilla/mux"
    "github.com/julianstephens/go-utils/httputil/middleware"
)

func main() {
    router := mux.NewRouter()
    
    // Add request ID middleware
    router.Use(middleware.RequestID())
    
    router.HandleFunc("/api/trace", traceHandler).Methods("GET")
    
    http.ListenAndServe(":8080", router)
}

func traceHandler(w http.ResponseWriter, r *http.Request) {
    // Get request ID from context
    requestID := middleware.GetRequestID(r.Context())
    
    // Use request ID for logging or tracing
    fmt.Printf("Processing request %s\n", requestID)
    
    // Include in response headers for client-side tracing
    w.Header().Set("X-Request-ID", requestID)
    w.Header().Set("Content-Type", "application/json")
    
    response := fmt.Sprintf(`{
        "message": "Request processed",
        "request_id": "%s",
        "trace_info": "Use this ID for support requests"
    }`, requestID)
    
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(response))
}
```

### Logging Middleware

```go
package main

import (
    "log"
    "net/http"
    "os"
    "time"
    
    "github.com/gorilla/mux"
    "github.com/julianstephens/go-utils/httputil/middleware"
)

func main() {
    // Create custom logger with specific format
    logger := log.New(os.Stdout, "[API] ", log.LstdFlags|log.Lmicroseconds)
    
    router := mux.NewRouter()
    
    // Add request ID first, then logging
    router.Use(middleware.RequestID())
    router.Use(middleware.Logging(logger))
    
    router.HandleFunc("/api/slow", slowHandler).Methods("GET")
    router.HandleFunc("/api/fast", fastHandler).Methods("GET")
    
    log.Println("Server with logging middleware starting on :8080")
    http.ListenAndServe(":8080", router)
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
    // Simulate slow operation
    time.Sleep(2 * time.Second)
    
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Slow operation completed"))
}

func fastHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Fast operation completed"))
}
```

### Recovery Middleware

```go
package main

import (
    "log"
    "net/http"
    "os"
    
    "github.com/gorilla/mux"
    "github.com/julianstephens/go-utils/httputil/middleware"
)

func main() {
    errorLogger := log.New(os.Stderr, "[PANIC] ", log.LstdFlags|log.Lshortfile)
    
    router := mux.NewRouter()
    
    // Add recovery middleware to catch panics
    router.Use(middleware.RequestID())
    router.Use(middleware.Recovery(errorLogger))
    
    router.HandleFunc("/api/panic", panicHandler).Methods("GET")
    router.HandleFunc("/api/safe", safeHandler).Methods("GET")
    
    log.Println("Server with recovery middleware starting on :8080")
    log.Println("Try GET /api/panic to test panic recovery")
    http.ListenAndServe(":8080", router)
}

func panicHandler(w http.ResponseWriter, r *http.Request) {
    // This will cause a panic
    var data map[string]string
    data["key"] = "value" // nil map assignment causes panic
}

func safeHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Safe handler - no panic here"))
}
```

### CORS Middleware

```go
package main

import (
    "net/http"
    
    "github.com/gorilla/mux"
    "github.com/julianstephens/go-utils/httputil/middleware"
)

func main() {
    router := mux.NewRouter()
    
    // Custom CORS configuration
    corsConfig := middleware.CORSConfig{
        AllowedOrigins:   []string{"https://myapp.com", "https://api.myapp.com"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Requested-With"},
        ExposedHeaders:   []string{"X-Total-Count", "X-Request-ID"},
        AllowCredentials: true,
        MaxAge:          86400, // 24 hours
    }
    
    router.Use(middleware.CORS(corsConfig))
    
    // API routes
    router.HandleFunc("/api/data", dataHandler).Methods("GET", "POST")
    router.HandleFunc("/api/upload", uploadHandler).Methods("POST")
    
    // Preflight requests are handled automatically by CORS middleware
    
    http.ListenAndServe(":8080", router)
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"data": "example"}`))
    case "POST":
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        w.Write([]byte(`{"message": "Data created"}`))
    }
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    // Handle file upload
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"message": "Upload successful"}`))
}
```

### Combined Middleware Stack

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "time"
    
    "github.com/gorilla/mux"
    "github.com/julianstephens/go-utils/httputil/middleware"
)

type APIResponse struct {
    Success   bool        `json:"success"`
    Data      interface{} `json:"data,omitempty"`
    Error     string      `json:"error,omitempty"`
    RequestID string      `json:"request_id"`
    Timestamp time.Time   `json:"timestamp"`
}

func main() {
    // Setup loggers
    accessLogger := log.New(os.Stdout, "[ACCESS] ", log.LstdFlags)
    errorLogger := log.New(os.Stderr, "[ERROR] ", log.LstdFlags)
    
    router := mux.NewRouter()
    
    // Apply full middleware stack
    router.Use(middleware.RequestID())
    router.Use(middleware.Logging(accessLogger))
    router.Use(middleware.Recovery(errorLogger))
    router.Use(middleware.CORS(middleware.DefaultCORSConfig()))
    
    // API routes
    api := router.PathPrefix("/api/v1").Subrouter()
    api.HandleFunc("/users", listUsersHandler).Methods("GET")
    api.HandleFunc("/users", createUserHandler).Methods("POST")
    api.HandleFunc("/users/{id}", getUserHandler).Methods("GET")
    api.HandleFunc("/health", healthCheckHandler).Methods("GET")
    api.HandleFunc("/error", errorHandler).Methods("GET")
    
    log.Println("Full-featured API server starting on :8080")
    log.Println("Available endpoints:")
    log.Println("  GET  /api/v1/health")
    log.Println("  GET  /api/v1/users")
    log.Println("  POST /api/v1/users")
    log.Println("  GET  /api/v1/users/{id}")
    log.Println("  GET  /api/v1/error (test error handling)")
    
    http.ListenAndServe(":8080", router)
}

func respondJSON(w http.ResponseWriter, r *http.Request, status int, data interface{}, err string) {
    requestID := middleware.GetRequestID(r.Context())
    
    response := APIResponse{
        Success:   status < 400,
        RequestID: requestID,
        Timestamp: time.Now(),
    }
    
    if err != "" {
        response.Error = err
    } else {
        response.Data = data
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("X-Request-ID", requestID)
    w.WriteHeader(status)
    
    json.NewEncoder(w).Encode(response)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
    respondJSON(w, r, http.StatusOK, map[string]string{
        "status": "healthy",
        "version": "1.0.0",
    }, "")
}

func listUsersHandler(w http.ResponseWriter, r *http.Request) {
    users := []map[string]interface{}{
        {"id": 1, "name": "Alice", "email": "alice@example.com"},
        {"id": 2, "name": "Bob", "email": "bob@example.com"},
    }
    
    respondJSON(w, r, http.StatusOK, users, "")
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
    var user map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        respondJSON(w, r, http.StatusBadRequest, nil, "Invalid JSON")
        return
    }
    
    // Simulate user creation
    user["id"] = 123
    user["created_at"] = time.Now()
    
    respondJSON(w, r, http.StatusCreated, user, "")
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["id"]
    
    user := map[string]interface{}{
        "id":    userID,
        "name":  "User " + userID,
        "email": "user" + userID + "@example.com",
    }
    
    respondJSON(w, r, http.StatusOK, user, "")
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
    respondJSON(w, r, http.StatusInternalServerError, nil, "This is a test error")
}
```

## Configuration

### CORS Configuration

```go
type CORSConfig struct {
    AllowedOrigins   []string // Allowed origins (e.g., ["https://example.com"])
    AllowedMethods   []string // Allowed HTTP methods
    AllowedHeaders   []string // Allowed headers
    ExposedHeaders   []string // Headers exposed to client
    AllowCredentials bool     // Allow credentials (cookies, auth headers)
    MaxAge          int      // Preflight cache duration in seconds
}

// Get default CORS config
config := middleware.DefaultCORSConfig()

// Custom CORS config
config := middleware.CORSConfig{
    AllowedOrigins: []string{"https://myapp.com"},
    AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
    AllowedHeaders: []string{"Authorization", "Content-Type"},
    AllowCredentials: true,
    MaxAge: 3600,
}
```

## API Reference

### Middleware Functions

- `RequestID() func(http.Handler) http.Handler` - Adds unique request ID to context
- `Logging(logger *log.Logger) func(http.Handler) http.Handler` - Logs HTTP requests/responses
- `Recovery(logger *log.Logger) func(http.Handler) http.Handler` - Recovers from panics
- `CORS(config CORSConfig) func(http.Handler) http.Handler` - Handles CORS headers
- `JWTAuth(manager *auth.JWTManager) func(http.Handler) http.Handler` - JWT token validation
- `RequireRoles(manager *auth.JWTManager, roles ...string) func(http.Handler) http.Handler` - Role-based access control

### Utility Functions

- `GetRequestID(ctx context.Context) string` - Extract request ID from context
- `GetClaims(ctx context.Context) (*auth.Claims, bool)` - Extract JWT claims from context
- `DefaultCORSConfig() CORSConfig` - Get default CORS configuration

### Request ID Context

The RequestID middleware adds a unique identifier to each request that can be accessed throughout the request lifecycle:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    requestID := middleware.GetRequestID(r.Context())
    // Use requestID for logging, tracing, etc.
}
```

## Middleware Order

Apply middleware in this recommended order for best results:

1. **RequestID** - First, to ensure all subsequent middleware has access to request ID
2. **Logging** - Second, to log all requests with their IDs
3. **Recovery** - Third, to catch panics from application handlers
4. **CORS** - Fourth, to handle cross-origin requests
5. **JWT Authentication** - Fifth, for protected routes
6. **Role-based Authorization** - Sixth, for role-specific routes
7. **Application handlers** - Last

```go
router.Use(middleware.RequestID())         // 1st
router.Use(middleware.Logging(logger))     // 2nd  
router.Use(middleware.Recovery(errLog))    // 3rd
router.Use(middleware.CORS(corsConfig))    // 4th

// For protected routes only
protected.Use(middleware.JWTAuth(jwtMgr))           // 5th
adminRoutes.Use(middleware.RequireRoles(jwtMgr, "admin")) // 6th
```

## Logging Format

The logging middleware outputs structured log entries:

```
[HTTP] 2023/10/15 14:30:45 GET /api/users 200 1.234ms [req-abc123]
[HTTP] 2023/10/15 14:30:46 POST /api/users 201 0.567ms [req-def456]
```

Format: `METHOD PATH STATUS_CODE DURATION [REQUEST_ID]`

## Error Handling

### Recovery Middleware

When a panic occurs, the recovery middleware:

1. Catches the panic
2. Logs the error with stack trace
3. Returns HTTP 500 with generic error message
4. Prevents the server from crashing

### JWT Authentication Errors

JWT middleware returns structured JSON error responses:

**401 Unauthorized** - Missing or invalid token:
```json
{
  "error": "authorization header required"
}
```

```json
{
  "error": "invalid or expired token"
}
```

**403 Forbidden** - Valid token but insufficient role:
```json
{
  "error": "insufficient role"
}
```

### CORS Errors

CORS middleware handles:

- Preflight OPTIONS requests
- Invalid origin rejections
- Missing required headers
- Method not allowed errors

## Best Practices

1. **Apply middleware in correct order** as shown above
2. **Use request IDs for tracing** across distributed systems
3. **Configure CORS restrictively** for production
4. **Log errors appropriately** - avoid exposing sensitive information
5. **Include request ID in error responses** for debugging
6. **Use structured logging** for better log analysis
7. **Monitor panic frequency** to identify unstable code

## Security Considerations

1. **CORS Configuration**: Be restrictive with allowed origins in production
2. **Error Information**: Don't expose internal errors to clients
3. **Request ID Format**: Use cryptographically random IDs
4. **Logging**: Don't log sensitive data (passwords, tokens)
5. **Headers**: Be careful about exposing internal headers

## Performance Considerations

- Middleware adds minimal overhead (typically < 1ms per request)
- Request ID generation is fast (UUID v4)
- Logging is synchronous - consider async logging for high-traffic applications
- CORS preflight responses are cached by browsers

## JWT Claims Context

After successful JWT authentication, claims are available in the request context:

```go
func protectedHandler(w http.ResponseWriter, r *http.Request) {
    claims, ok := middleware.GetClaims(r.Context())
    if !ok {
        // This shouldn't happen if JWTAuth middleware is working correctly
        http.Error(w, "No claims found", http.StatusInternalServerError)
        return
    }
    
    // Access claim information
    userID := claims.UserID
    username := claims.Username
    email := claims.Email
    roles := claims.Roles
    
    // Check specific roles
    if claims.HasRole("admin") {
        // Admin-specific logic
    }
    
    // Check if user has any of multiple roles
    if claims.HasAnyRole("manager", "admin") {
        // Manager or admin logic
    }
    
    // Access custom claims
    if department, ok := claims.GetCustomClaimString("department"); ok {
        // Use department information
    }
}
```

## Integration

Works well with other go-utils packages:

```go
// Use with logger package for structured logging
logger.WithField("request_id", middleware.GetRequestID(ctx)).Info("Processing request")

// Use with response package for consistent JSON responses
responder := response.New()
if claims, ok := middleware.GetClaims(ctx); ok {
    responder.OK(w, r, map[string]string{"user": claims.UserID})
} else {
    responder.Unauthorized(w, r, "authentication required")
}

// Full integration example
router.Use(middleware.RequestID())
router.Use(middleware.Logging(logger))
router.Use(middleware.Recovery(errorLogger))
router.Use(middleware.CORS(middleware.DefaultCORSConfig()))

// JWT-protected routes
jwtManager := auth.NewJWTManager("secret", time.Hour*24, "myapp")
protectedRoutes := router.PathPrefix("/api/protected").Subrouter()
protectedRoutes.Use(middleware.JWTAuth(jwtManager))

// Role-protected routes  
adminRoutes := router.PathPrefix("/api/admin").Subrouter()
adminRoutes.Use(middleware.RequireRoles(jwtManager, "admin"))
```