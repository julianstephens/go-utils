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
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
    
    "github.com/julianstephens/go-utils/httputil/response"
)

type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

func main() {
    // Create a responder with JSON encoder
    responder := &response.Responder{
        Encoder: &response.JSONEncoder{},
    }
    
    http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
        user := User{
            ID:        1,
            Name:      "John Doe",
            Email:     "john@example.com",
            CreatedAt: time.Now(),
        }
        
        // Send JSON response with 200 OK
        responder.OK(w, r, user)
    })
    
    http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        users := []User{
            {ID: 1, Name: "Alice", Email: "alice@example.com", CreatedAt: time.Now()},
            {ID: 2, Name: "Bob", Email: "bob@example.com", CreatedAt: time.Now()},
        }
        
        // Send JSON response with 200 OK
        responder.OK(w, r, users)
    })
    
    log.Println("Server starting on :8080")
    http.ListenAndServe(":8080", nil)
}
```

### Response with Status Codes

```go
package main

import (
    "net/http"
    "time"
    
    "github.com/julianstephens/go-utils/httputil/response"
)

type APIResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func main() {
    responder := &response.Responder{
        Encoder: &response.JSONEncoder{},
    }
    
    // Success responses
    http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
        newUser := map[string]interface{}{
            "id":   123,
            "name": "New User",
        }
        
        apiResp := APIResponse{
            Success: true,
            Message: "User created successfully",
            Data:    newUser,
        }
        
        responder.Created(w, r, apiResp)
    })
    
    // No content response
    http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
        // Simulate deletion
        responder.NoContent(w, r)
    })
    
    // Error responses
    http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
        apiResp := APIResponse{
            Success: false,
            Error:   "Something went wrong",
        }
        
        responder.InternalServerError(w, r, apiResp)
    })
    
    // Not found
    http.HandleFunc("/notfound", func(w http.ResponseWriter, r *http.Request) {
        apiResp := APIResponse{
            Success: false,
            Error:   "Resource not found",
        }
        
        responder.NotFound(w, r, apiResp)
    })
    
    http.ListenAndServe(":8080", nil)
}
```

### Response Hooks for Logging and Metrics

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
    
    "github.com/julianstephens/go-utils/httputil/response"
)

func main() {
    responder := &response.Responder{
        Encoder: &response.JSONEncoder{},
        
        // Before hook - called before encoding response
        Before: func(w http.ResponseWriter, r *http.Request, data interface{}) {
            log.Printf("Sending response for %s %s", r.Method, r.URL.Path)
            
            // Add custom headers
            w.Header().Set("X-API-Version", "v1.0")
            w.Header().Set("X-Response-Time", time.Now().Format(time.RFC3339))
        },
        
        // After hook - called after successful encoding
        After: func(w http.ResponseWriter, r *http.Request, data interface{}) {
            log.Printf("Successfully sent response for %s %s", r.Method, r.URL.Path)
            
            // Could collect metrics here
            // metrics.IncrementResponseCounter(r.Method, r.URL.Path, "success")
        },
        
        // Error hook - called when encoding fails
        OnError: func(w http.ResponseWriter, r *http.Request, err error) {
            log.Printf("Failed to encode response for %s %s: %v", r.Method, r.URL.Path, err)
            
            // Send fallback error response
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusInternalServerError)
            w.Write([]byte(`{"error": "Internal server error"}`))
            
            // Could collect error metrics here
            // metrics.IncrementResponseCounter(r.Method, r.URL.Path, "error")
        },
    }
    
    http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
        data := map[string]interface{}{
            "message": "Hello, World!",
            "time":    time.Now(),
        }
        
        responder.OK(w, r, data)
    })
    
    http.HandleFunc("/api/broken", func(w http.ResponseWriter, r *http.Request) {
        // This will cause encoding to fail (channels can't be JSON encoded)
        data := make(chan int)
        responder.OK(w, r, data)
    })
    
    log.Println("Server with hooks starting on :8080")
    http.ListenAndServe(":8080", nil)
}
```

### Custom Encoder

```go
package main

import (
    "encoding/xml"
    "fmt"
    "net/http"
    
    "github.com/julianstephens/go-utils/httputil/response"
)

// XMLEncoder implements custom XML encoding
type XMLEncoder struct{}

func (e *XMLEncoder) Encode(w http.ResponseWriter, v interface{}) error {
    w.Header().Set("Content-Type", "application/xml")
    
    encoder := xml.NewEncoder(w)
    encoder.Indent("", "  ")
    
    // Write XML declaration
    w.Write([]byte(xml.Header))
    
    return encoder.Encode(v)
}

type Product struct {
    XMLName     xml.Name `xml:"product"`
    ID          int      `xml:"id,attr"`
    Name        string   `xml:"name"`
    Description string   `xml:"description"`
    Price       float64  `xml:"price"`
}

func main() {
    // Create responder with XML encoder
    xmlResponder := &response.Responder{
        Encoder: &XMLEncoder{},
    }
    
    // Create responder with JSON encoder for comparison
    jsonResponder := &response.Responder{
        Encoder: &response.JSONEncoder{},
    }
    
    http.HandleFunc("/product.xml", func(w http.ResponseWriter, r *http.Request) {
        product := Product{
            ID:          1,
            Name:        "Laptop",
            Description: "High-performance laptop",
            Price:       999.99,
        }
        
        xmlResponder.OK(w, r, product)
    })
    
    http.HandleFunc("/product.json", func(w http.ResponseWriter, r *http.Request) {
        product := Product{
            ID:          1,
            Name:        "Laptop",
            Description: "High-performance laptop",
            Price:       999.99,
        }
        
        jsonResponder.OK(w, r, product)
    })
    
    fmt.Println("Server starting on :8080")
    fmt.Println("Try:")
    fmt.Println("  http://localhost:8080/product.xml")
    fmt.Println("  http://localhost:8080/product.json")
    
    http.ListenAndServe(":8080", nil)
}
```

### Complete API with Structured Responses

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "time"
    
    "github.com/gorilla/mux"
    "github.com/julianstephens/go-utils/httputil/response"
)

type APIResponse struct {
    Success   bool        `json:"success"`
    Data      interface{} `json:"data,omitempty"`
    Error     string      `json:"error,omitempty"`
    Message   string      `json:"message,omitempty"`
    RequestID string      `json:"request_id,omitempty"`
    Timestamp time.Time   `json:"timestamp"`
}

type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

// Mock database
var users = map[int]User{
    1: {ID: 1, Name: "Alice", Email: "alice@example.com", CreatedAt: time.Now()},
    2: {ID: 2, Name: "Bob", Email: "bob@example.com", CreatedAt: time.Now()},
}
var nextID = 3

func createAPIResponse(success bool, data interface{}, message, error string) APIResponse {
    return APIResponse{
        Success:   success,
        Data:      data,
        Message:   message,
        Error:     error,
        Timestamp: time.Now(),
    }
}

func main() {
    responder := &response.Responder{
        Encoder: &response.JSONEncoder{},
        
        Before: func(w http.ResponseWriter, r *http.Request, data interface{}) {
            // Add CORS headers
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        },
        
        OnError: func(w http.ResponseWriter, r *http.Request, err error) {
            log.Printf("Response encoding error: %v", err)
            
            errorResp := createAPIResponse(false, nil, "", "Internal server error")
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusInternalServerError)
            json.NewEncoder(w).Encode(errorResp)
        },
    }
    
    router := mux.NewRouter()
    api := router.PathPrefix("/api/v1").Subrouter()
    
    // List users
    api.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        userList := make([]User, 0, len(users))
        for _, user := range users {
            userList = append(userList, user)
        }
        
        apiResp := createAPIResponse(true, userList, "Users retrieved successfully", "")
        responder.OK(w, r, apiResp)
    }).Methods("GET")
    
    // Get user by ID
    api.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
            apiResp := createAPIResponse(false, nil, "", "Invalid user ID")
            responder.BadRequest(w, r, apiResp)
            return
        }
        
        user, exists := users[id]
        if !exists {
            apiResp := createAPIResponse(false, nil, "", "User not found")
            responder.NotFound(w, r, apiResp)
            return
        }
        
        apiResp := createAPIResponse(true, user, "User retrieved successfully", "")
        responder.OK(w, r, apiResp)
    }).Methods("GET")
    
    // Create user
    api.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        var userData struct {
            Name  string `json:"name"`
            Email string `json:"email"`
        }
        
        if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
            apiResp := createAPIResponse(false, nil, "", "Invalid JSON")
            responder.BadRequest(w, r, apiResp)
            return
        }
        
        if userData.Name == "" || userData.Email == "" {
            apiResp := createAPIResponse(false, nil, "", "Name and email are required")
            responder.BadRequest(w, r, apiResp)
            return
        }
        
        user := User{
            ID:        nextID,
            Name:      userData.Name,
            Email:     userData.Email,
            CreatedAt: time.Now(),
        }
        nextID++
        
        users[user.ID] = user
        
        apiResp := createAPIResponse(true, user, "User created successfully", "")
        responder.Created(w, r, apiResp)
    }).Methods("POST")
    
    // Delete user
    api.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
            apiResp := createAPIResponse(false, nil, "", "Invalid user ID")
            responder.BadRequest(w, r, apiResp)
            return
        }
        
        if _, exists := users[id]; !exists {
            apiResp := createAPIResponse(false, nil, "", "User not found")
            responder.NotFound(w, r, apiResp)
            return
        }
        
        delete(users, id)
        responder.NoContent(w, r)
    }).Methods("DELETE")
    
    // Health check
    api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        health := map[string]interface{}{
            "status":    "healthy",
            "timestamp": time.Now(),
            "version":   "1.0.0",
        }
        
        apiResp := createAPIResponse(true, health, "Service is healthy", "")
        responder.OK(w, r, apiResp)
    }).Methods("GET")
    
    log.Println("API server starting on :8080")
    log.Println("Endpoints:")
    log.Println("  GET    /api/v1/health")
    log.Println("  GET    /api/v1/users")
    log.Println("  GET    /api/v1/users/{id}")
    log.Println("  POST   /api/v1/users")
    log.Println("  DELETE /api/v1/users/{id}")
    
    http.ListenAndServe(":8080", router)
}
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