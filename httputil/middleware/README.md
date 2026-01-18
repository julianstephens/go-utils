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
logger := log.New(os.Stdout, "[HTTP] ", log.LstdFlags)
router := mux.NewRouter()

router.Use(middleware.RequestID())
router.Use(middleware.Logging(logger))
router.Use(middleware.Recovery(log.New(os.Stderr, "[ERROR] ", log.LstdFlags)))
router.Use(middleware.CORS(middleware.DefaultCORSConfig()))

router.HandleFunc("/api/health", healthHandler).Methods("GET")
log.Fatal(http.ListenAndServe(":8080", router))
```

### JWT Authentication Middleware

The JWT authentication middleware provides token validation and role-based access control using the `httputil/auth` package.

#### Basic JWT Authentication

```go
jwtManager := auth.NewJWTManager("your-secret-key", time.Hour*24, "your-app")
router := mux.NewRouter()

router.HandleFunc("/api/login", loginHandler(jwtManager)).Methods("POST")

protected := router.PathPrefix("/api/protected").Subrouter()
protected.Use(middleware.JWTAuth(jwtManager))
protected.HandleFunc("/profile", profileHandler).Methods("GET")
```

#### Role-Based Access Control

```go
jwtManager := auth.NewJWTManager("your-secret-key", time.Hour*24, "your-app")
router := mux.NewRouter()

router.HandleFunc("/api/login", loginHandler(jwtManager)).Methods("POST")

userRoutes := router.PathPrefix("/api/user").Subrouter()
userRoutes.Use(middleware.RequireRoles(jwtManager, "user"))
userRoutes.HandleFunc("/profile", userProfileHandler).Methods("GET")

adminRoutes := router.PathPrefix("/api/admin").Subrouter()
adminRoutes.Use(middleware.RequireRoles(jwtManager, "admin"))
adminRoutes.HandleFunc("/users", adminUsersHandler).Methods("GET")
```

#### JWT with Full Middleware Stack

```go
jwtManager := auth.NewJWTManager("secret-key", time.Hour*24, "myapp")
router := mux.NewRouter()

router.Use(middleware.RequestID())
router.Use(middleware.Logging(logger))
router.Use(middleware.Recovery(errLog))
router.Use(middleware.CORS(middleware.DefaultCORSConfig()))

api := router.PathPrefix("/api/v1").Subrouter()
api.HandleFunc("/health", healthHandler).Methods("GET")
api.HandleFunc("/login", loginHandler(jwtManager)).Methods("POST")

protected := api.PathPrefix("/protected").Subrouter()
protected.Use(middleware.JWTAuth(jwtManager))
protected.HandleFunc("/profile", profileHandler).Methods("GET")

admin := api.PathPrefix("/admin").Subrouter()
admin.Use(middleware.RequireRoles(jwtManager, "admin"))
admin.HandleFunc("/users", adminHandler).Methods("GET")

http.ListenAndServe(":8080", router)
```

### Request ID Middleware

```go
router := mux.NewRouter()
router.Use(middleware.RequestID())
router.HandleFunc("/api/trace", traceHandler).Methods("GET")

func traceHandler(w http.ResponseWriter, r *http.Request) {
    requestID := middleware.GetRequestID(r.Context())
    w.Header().Set("X-Request-ID", requestID)
    w.Write([]byte(`{"message": "Request processed", "request_id": "` + requestID + `"}`))
}
```

### Logging Middleware

```go
logger := log.New(os.Stdout, "[API] ", log.LstdFlags|log.Lmicroseconds)
router := mux.NewRouter()

router.Use(middleware.RequestID())
router.Use(middleware.Logging(logger))

router.HandleFunc("/api/data", handler).Methods("GET")
http.ListenAndServe(":8080", router)
```

### Recovery Middleware

```go
errorLogger := log.New(os.Stderr, "[PANIC] ", log.LstdFlags|log.Lshortfile)
router := mux.NewRouter()
router.Use(middleware.RequestID())
router.Use(middleware.Recovery(errorLogger))

router.HandleFunc("/api/panic", panicHandler).Methods("GET")
router.HandleFunc("/api/safe", safeHandler).Methods("GET")
```

### CORS Middleware

```go
corsConfig := middleware.CORSConfig{
    AllowedOrigins: []string{"https://myapp.com"},
    AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders: []string{"Authorization", "Content-Type"},
    MaxAge: 86400,
}

router := mux.NewRouter()
router.Use(middleware.CORS(corsConfig))
router.HandleFunc("/api/data", dataHandler).Methods("GET", "POST")
```

### Combined Middleware Stack

```go
accessLogger := log.New(os.Stdout, "[ACCESS] ", log.LstdFlags)
errorLogger := log.New(os.Stderr, "[ERROR] ", log.LstdFlags)

router := mux.NewRouter()
router.Use(middleware.RequestID())
router.Use(middleware.Logging(accessLogger))
router.Use(middleware.Recovery(errorLogger))
router.Use(middleware.CORS(middleware.DefaultCORSConfig()))

api := router.PathPrefix("/api/v1").Subrouter()
api.HandleFunc("/users", listUsersHandler).Methods("GET")
api.HandleFunc("/users", createUserHandler).Methods("POST")
api.HandleFunc("/health", healthCheckHandler).Methods("GET")
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

- **Recovery**: Catches panics, logs errors, returns HTTP 500
- **JWT (401)**: Missing/invalid token returns `{"error": "authorization header required"}`
- **JWT (403)**: Insufficient role returns `{"error": "insufficient role"}`
- **CORS**: Handles preflight OPTIONS, validates origins, validates headers

## Best Practices

- Apply middleware in recommended order (RequestID → Logging → Recovery → CORS → JWT)
- Use request IDs for distributed tracing
- Configure CORS restrictively in production
- Log strategically - avoid exposing sensitive information
- Include request IDs in error responses

## Security Considerations

- CORS: Restrictive allowed origins in production
- Errors: Don't expose internal errors to clients
- Logging: Avoid logging passwords, tokens, sensitive headers
- Request IDs: Use cryptographically random format

## Performance

- Minimal overhead (< 1ms per request)
- Request ID generation: O(1) UUID v4
- Logging: Synchronous (consider async for high-traffic)
- CORS preflight: Cached by browsers

## JWT Claims Context

Access JWT claims in request context after authentication:

```go
func protectedHandler(w http.ResponseWriter, r *http.Request) {
    claims, ok := middleware.GetClaims(r.Context())
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    
    userID := claims.UserID
    roles := claims.Roles
    if claims.HasRole("admin") { /* admin logic */ }
}
```

## Integration

Works well with other go-utils packages:
- **logger**: Use request IDs for structured logging across requests
- **response**: Combine with response package for consistent JSON responses
- **auth**: JWTManager integrates seamlessly with JWTAuth middleware
- **httputil/request**: Parse incoming JSON with request package