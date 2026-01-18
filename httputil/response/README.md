# HTTP Response Package

The `httputil/response` package provides structured HTTP response handling with extensible encoders, hooks, and status code helpers. It enables consistent response formatting and provides flexibility for custom response processing.

## Features

- **Structured Response Handling**: Consistent response formatting across your API
- **Extensible Encoders**: Support for JSON, XML, or custom response formats
- **Response Hooks**: Before/after processing hooks for logging, metrics, etc.
- **Error Handling**: Centralized error response handling
- **Status Code Helpers**: Convenient functions for common HTTP status codes
- **Flexible Architecture**: Easy to extend and customize

## Installation

```bash
go get github.com/julianstephens/go-utils/httputil/response
```

## Usage

### Basic JSON Response Handling

```go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

responder := &response.Responder{Encoder: &response.JSONEncoder{}}

http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
    user := User{ID: 1, Name: "John Doe", Email: "john@example.com"}
    responder.OK(w, r, user)
})
```

### Response with Status Codes

```go
responder := &response.Responder{Encoder: &response.JSONEncoder{}}

http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{}{"id": 123, "name": "New Item"}
    responder.Created(w, r, data)
})

http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
    responder.NoContent(w, r)
})

http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
    responder.InternalServerError(w, r, map[string]string{"error": "Server error"})
})
```

### Response Hooks for Logging and Metrics

```go
responder := &response.Responder{
    Encoder: &response.JSONEncoder{},
    Before: func(w http.ResponseWriter, r *http.Request, data interface{}) {
        w.Header().Set("X-API-Version", "v1.0")
    },
    After: func(w http.ResponseWriter, r *http.Request, data interface{}) {
        log.Printf("Sent response for %s %s", r.Method, r.URL.Path)
    },
    OnError: func(w http.ResponseWriter, r *http.Request, err error) {
        log.Printf("Response error: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
    },
}

http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
    responder.OK(w, r, map[string]string{"message": "Hello"})
})
```

### Custom Encoder

```go
type XMLEncoder struct{}

func (e *XMLEncoder) Encode(w http.ResponseWriter, v interface{}) error {
    w.Header().Set("Content-Type", "application/xml")
    encoder := xml.NewEncoder(w)
    w.Write([]byte(xml.Header))
    return encoder.Encode(v)
}

xmlResponder := &response.Responder{Encoder: &XMLEncoder{}}
jsonResponder := &response.Responder{Encoder: &response.JSONEncoder{}}

http.HandleFunc("/data.xml", func(w http.ResponseWriter, r *http.Request) {
    xmlResponder.OK(w, r, data)
})
http.HandleFunc("/data.json", func(w http.ResponseWriter, r *http.Request) {
    jsonResponder.OK(w, r, data)
})
```

### Complete API with Structured Responses

```go
type APIResponse struct {
    Success   bool        `json:"success"`
    Data      interface{} `json:"data,omitempty"`
    Error     string      `json:"error,omitempty"`
    Timestamp time.Time   `json:"timestamp"`
}

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

var users = map[int]User{1: {ID: 1, Name: "Alice"}}
var nextID = 2

responder := &response.Responder{Encoder: &response.JSONEncoder{}}

http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        userList := make([]User, 0, len(users))
        for _, u := range users {
            userList = append(userList, u)
        }
        responder.OK(w, r, APIResponse{Success: true, Data: userList, Timestamp: time.Now()})
    } else if r.Method == "POST" {
        var u User
        json.NewDecoder(r.Body).Decode(&u)
        u.ID = nextID
        nextID++
        users[u.ID] = u
        responder.Created(w, r, APIResponse{Success: true, Data: u, Timestamp: time.Now()})
    }
})

http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
    responder.OK(w, r, APIResponse{Success: true, Data: map[string]string{"status": "healthy"}, Timestamp: time.Now()})
})
```

## API Reference

### Responder Struct

```go
type Responder struct {
    Encoder Encoder     // Encoder for response data
    Before  BeforeFunc  // Hook called before encoding
    After   AfterFunc   // Hook called after successful encoding
    OnError OnErrorFunc // Hook called on encoding errors
}
```

### Status Code Methods

- `OK(w http.ResponseWriter, r *http.Request, data interface{})` - 200 OK
- `Created(w http.ResponseWriter, r *http.Request, data interface{})` - 201 Created
- `NoContent(w http.ResponseWriter, r *http.Request)` - 204 No Content
- `BadRequest(w http.ResponseWriter, r *http.Request, data interface{})` - 400 Bad Request
- `Unauthorized(w http.ResponseWriter, r *http.Request, data interface{})` - 401 Unauthorized
- `Forbidden(w http.ResponseWriter, r *http.Request, data interface{})` - 403 Forbidden
- `NotFound(w http.ResponseWriter, r *http.Request, data interface{})` - 404 Not Found
- `InternalServerError(w http.ResponseWriter, r *http.Request, data interface{})` - 500 Internal Server Error

### Custom Status

- `WriteWithStatus(w http.ResponseWriter, r *http.Request, data interface{}, statusCode int)` - Write with custom status code

### Encoders

#### JSONEncoder
```go
type JSONEncoder struct{}
func (e *JSONEncoder) Encode(w http.ResponseWriter, v interface{}) error
```

#### Custom Encoder Interface
```go
type Encoder interface {
    Encode(w http.ResponseWriter, v interface{}) error
}
```

### Hook Functions

```go
type BeforeFunc func(w http.ResponseWriter, r *http.Request, data interface{})
type AfterFunc func(w http.ResponseWriter, r *http.Request, data interface{})
type OnErrorFunc func(w http.ResponseWriter, r *http.Request, err error)
```

## Hook Use Cases

### Before Hook
- Add custom headers (CORS, caching, etc.)
- Log request processing start
- Modify response writer state
- Add response metadata

### After Hook  
- Log successful responses
- Collect metrics
- Update caches
- Trigger webhooks

### OnError Hook
- Log encoding errors
- Send fallback error responses
- Collect error metrics
- Alert monitoring systems

## Custom Encoders

Implement the `Encoder` interface for custom response formats:

```go
type CustomEncoder struct{}

func (e *CustomEncoder) Encode(w http.ResponseWriter, v interface{}) error {
    // Set appropriate Content-Type
    w.Header().Set("Content-Type", "your/content-type")
    
    // Encode data to writer
    return yourEncodingLogic(w, v)
}
```

## Error Handling

The package provides several layers of error handling:

1. **Encoding Errors**: Handled by OnError hook
2. **Status Code Errors**: Use appropriate status methods
3. **Fallback Responses**: Implement in OnError hook
4. **Logging**: Add to hooks for comprehensive error tracking

## Best Practices

1. **Use structured responses** for consistent API contracts
2. **Implement proper error handling** in OnError hook
3. **Add request tracing** using Before/After hooks
4. **Set appropriate Content-Type** in custom encoders
5. **Use status code methods** instead of raw WriteWithStatus when possible
6. **Handle encoding failures gracefully** with fallback responses
7. **Add metrics collection** in hooks for monitoring
8. **Include request context** in error responses for debugging

## Thread Safety

- Responder instances are safe for concurrent use
- Encoders should be stateless and thread-safe
- Hook functions should handle concurrent execution safely
- Custom encoders must be thread-safe

## Integration

Works well with other go-utils packages:

```go
// Use with logger for structured response logging
responder := &response.Responder{
    Encoder: &response.JSONEncoder{},
    Before: func(w http.ResponseWriter, r *http.Request, data interface{}) {
        logger.WithField("path", r.URL.Path).Info("Sending response")
    },
}

// Use with httputil/request for complete request/response handling
// Use with httputil/middleware for comprehensive API infrastructure
```