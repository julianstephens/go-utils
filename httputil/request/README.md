# HTTP Request Package

The `httputil/request` package provides HTTP request parsing utilities for JSON, form data, query parameters, and URL values. It simplifies common request handling patterns in HTTP services.

## Features

- **JSON Decoding**: Safe JSON request body parsing with validation
- **Form Processing**: Form data parsing and value extraction
- **Query Parameters**: Query string parameter handling
- **Type Conversion**: Automatic conversion to common types (int, bool, float64)
- **Validation**: Content-type validation and error handling
- **Safety**: Disallows unknown fields in JSON to prevent injection

## Installation

```bash
go get github.com/julianstephens/go-utils/httputil/request
```

## Usage

### JSON Request Handling

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    
    "github.com/julianstephens/go-utils/httputil/request"
)

type CreateUserRequest struct {
    Name     string   `json:"name"`
    Email    string   `json:"email"`
    Age      int      `json:"age"`
    Roles    []string `json:"roles"`
    Active   bool     `json:"active"`
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    
    // Decode JSON request body
    if err := request.DecodeJSON(r, &req); err != nil {
        if err == request.ErrInvalidContentType {
            http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
            return
        }
        http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
        return
    }
    
    // Validate required fields
    if req.Name == "" || req.Email == "" {
        http.Error(w, "Name and email are required", http.StatusBadRequest)
        return
    }
    
    // Process the request
    fmt.Printf("Creating user: %+v\n", req)
    
    // Send response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte(`{"message": "User created successfully"}`))
}

func main() {
    http.HandleFunc("/users", createUserHandler)
    
    log.Println("Server starting on :8080")
    log.Println("POST to /users with JSON body to test")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Form Data Processing

```go
package main

import (
    "fmt"
    "net/http"
    
    "github.com/julianstephens/go-utils/httputil/request"
)

func contactFormHandler(w http.ResponseWriter, r *http.Request) {
    // Parse form data
    if err := request.ParseForm(r); err != nil {
        http.Error(w, "Invalid form data", http.StatusBadRequest)
        return
    }
    
    // Extract form values
    name, hasName := request.FormValue(r, "name")
    email, hasEmail := request.FormValue(r, "email")
    message, hasMessage := request.FormValue(r, "message")
    
    // Validate required fields
    if !hasName || !hasEmail || !hasMessage {
        http.Error(w, "Name, email, and message are required", http.StatusBadRequest)
        return
    }
    
    // Convert optional fields
    newsletter, _ := request.FormBool(r, "newsletter")
    priority, _ := request.FormInt(r, "priority", 1) // default priority = 1
    
    fmt.Printf("Contact form submission:\n")
    fmt.Printf("  Name: %s\n", name)
    fmt.Printf("  Email: %s\n", email)
    fmt.Printf("  Message: %s\n", message)
    fmt.Printf("  Newsletter: %t\n", newsletter)
    fmt.Printf("  Priority: %d\n", priority)
    
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Thank you for your message!"))
}

func main() {
    http.HandleFunc("/contact", contactFormHandler)
    
    // Serve HTML form for testing
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        html := `
        <html>
        <body>
            <form method="POST" action="/contact">
                <label>Name: <input type="text" name="name" required></label><br>
                <label>Email: <input type="email" name="email" required></label><br>
                <label>Message: <textarea name="message" required></textarea></label><br>
                <label>Newsletter: <input type="checkbox" name="newsletter" value="true"></label><br>
                <label>Priority: <input type="number" name="priority" min="1" max="5" value="1"></label><br>
                <button type="submit">Send Message</button>
            </form>
        </body>
        </html>
        `
        w.Header().Set("Content-Type", "text/html")
        w.Write([]byte(html))
    })
    
    fmt.Println("Server starting on :8080")
    fmt.Println("Visit http://localhost:8080 for contact form")
    http.ListenAndServe(":8080", nil)
}
```

### Query Parameter Handling

```go
package main

import (
    "fmt"
    "net/http"
    "strconv"
    
    "github.com/julianstephens/go-utils/httputil/request"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
    // Extract query parameters
    query, hasQuery := request.QueryValue(r, "q")
    if !hasQuery {
        http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
        return
    }
    
    // Extract optional parameters with defaults
    page, err := request.QueryInt(r, "page", 1)
    if err != nil {
        http.Error(w, "Invalid page number", http.StatusBadRequest)
        return
    }
    
    limit, err := request.QueryInt(r, "limit", 10)
    if err != nil {
        http.Error(w, "Invalid limit", http.StatusBadRequest)
        return
    }
    
    // Boolean parameters
    exact, _ := request.QueryBool(r, "exact")
    
    // Float parameters
    minPrice, _ := request.QueryFloat64(r, "min_price", 0.0)
    
    // Validate ranges
    if page < 1 {
        http.Error(w, "Page must be >= 1", http.StatusBadRequest)
        return
    }
    
    if limit < 1 || limit > 100 {
        http.Error(w, "Limit must be between 1 and 100", http.StatusBadRequest)
        return
    }
    
    // Process search
    fmt.Printf("Search parameters:\n")
    fmt.Printf("  Query: %s\n", query)
    fmt.Printf("  Page: %d\n", page)
    fmt.Printf("  Limit: %d\n", limit)
    fmt.Printf("  Exact match: %t\n", exact)
    fmt.Printf("  Min price: %.2f\n", minPrice)
    
    // Calculate offset for pagination
    offset := (page - 1) * limit
    
    // Simulate search results
    results := make([]map[string]interface{}, 0)
    for i := 0; i < limit; i++ {
        result := map[string]interface{}{
            "id":    offset + i + 1,
            "title": fmt.Sprintf("Result %d for '%s'", offset+i+1, query),
            "price": minPrice + float64(i)*10.99,
        }
        results = append(results, result)
    }
    
    response := map[string]interface{}{
        "query":   query,
        "page":    page,
        "limit":   limit,
        "results": results,
        "total":   42, // Mock total
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    
    // Simple JSON encoding (in real app, use proper JSON encoder)
    fmt.Fprintf(w, `{
        "query": "%s",
        "page": %d,
        "limit": %d,
        "total": 42,
        "results": []
    }`, query, page, limit)
}

func main() {
    http.HandleFunc("/search", searchHandler)
    
    fmt.Println("Server starting on :8080")
    fmt.Println("Try: http://localhost:8080/search?q=golang&page=2&limit=5&exact=true&min_price=19.99")
    http.ListenAndServe(":8080", nil)
}
```

### Complete API Handler Example

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    
    "github.com/gorilla/mux"
    "github.com/julianstephens/go-utils/httputil/request"
)

type Product struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Price       float64   `json:"price"`
    InStock     bool      `json:"in_stock"`
    CreatedAt   time.Time `json:"created_at"`
}

type CreateProductRequest struct {
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    InStock     bool    `json:"in_stock"`
}

type UpdateProductRequest struct {
    Name        *string  `json:"name,omitempty"`
    Description *string  `json:"description,omitempty"`
    Price       *float64 `json:"price,omitempty"`
    InStock     *bool    `json:"in_stock,omitempty"`
}

// Mock database
var products = []Product{
    {ID: 1, Name: "Laptop", Description: "Gaming laptop", Price: 999.99, InStock: true, CreatedAt: time.Now()},
    {ID: 2, Name: "Mouse", Description: "Wireless mouse", Price: 29.99, InStock: false, CreatedAt: time.Now()},
}
var nextID = 3

func listProductsHandler(w http.ResponseWriter, r *http.Request) {
    // Parse query parameters for filtering and pagination
    page, _ := request.QueryInt(r, "page", 1)
    limit, _ := request.QueryInt(r, "limit", 10)
    inStock, hasInStock := request.QueryBool(r, "in_stock")
    minPrice, _ := request.QueryFloat64(r, "min_price", 0)
    maxPrice, hasMaxPrice := request.QueryFloat64(r, "max_price", 0)
    
    // Filter products
    filtered := make([]Product, 0)
    for _, product := range products {
        // Filter by stock status
        if hasInStock && product.InStock != inStock {
            continue
        }
        
        // Filter by price range
        if product.Price < minPrice {
            continue
        }
        if hasMaxPrice && product.Price > maxPrice {
            continue
        }
        
        filtered = append(filtered, product)
    }
    
    // Paginate (simple implementation)
    start := (page - 1) * limit
    end := start + limit
    if start > len(filtered) {
        start = len(filtered)
    }
    if end > len(filtered) {
        end = len(filtered)
    }
    
    result := filtered[start:end]
    
    response := map[string]interface{}{
        "products": result,
        "page":     page,
        "limit":    limit,
        "total":    len(filtered),
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func createProductHandler(w http.ResponseWriter, r *http.Request) {
    var req CreateProductRequest
    
    if err := request.DecodeJSON(r, &req); err != nil {
        http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
        return
    }
    
    // Validate required fields
    if req.Name == "" {
        http.Error(w, "Name is required", http.StatusBadRequest)
        return
    }
    if req.Price <= 0 {
        http.Error(w, "Price must be positive", http.StatusBadRequest)
        return
    }
    
    // Create product
    product := Product{
        ID:          nextID,
        Name:        req.Name,
        Description: req.Description,
        Price:       req.Price,
        InStock:     req.InStock,
        CreatedAt:   time.Now(),
    }
    nextID++
    
    products = append(products, product)
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(product)
}

func updateProductHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, _ := request.QueryInt(r, "id", 0) // This is actually from URL, but example purposes
    
    // Find product
    var productIndex = -1
    for i, product := range products {
        if product.ID == id {
            productIndex = i
            break
        }
    }
    
    if productIndex == -1 {
        http.Error(w, "Product not found", http.StatusNotFound)
        return
    }
    
    var req UpdateProductRequest
    if err := request.DecodeJSON(r, &req); err != nil {
        http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
        return
    }
    
    // Update product fields
    product := &products[productIndex]
    if req.Name != nil {
        product.Name = *req.Name
    }
    if req.Description != nil {
        product.Description = *req.Description
    }
    if req.Price != nil {
        if *req.Price <= 0 {
            http.Error(w, "Price must be positive", http.StatusBadRequest)
            return
        }
        product.Price = *req.Price
    }
    if req.InStock != nil {
        product.InStock = *req.InStock
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(product)
}

func main() {
    router := mux.NewRouter()
    
    api := router.PathPrefix("/api/v1").Subrouter()
    api.HandleFunc("/products", listProductsHandler).Methods("GET")
    api.HandleFunc("/products", createProductHandler).Methods("POST")
    api.HandleFunc("/products/{id}", updateProductHandler).Methods("PUT")
    
    fmt.Println("Product API server starting on :8080")
    fmt.Println("Endpoints:")
    fmt.Println("  GET  /api/v1/products?page=1&limit=10&in_stock=true&min_price=20&max_price=500")
    fmt.Println("  POST /api/v1/products (JSON body)")
    fmt.Println("  PUT  /api/v1/products/{id} (JSON body)")
    
    http.ListenAndServe(":8080", router)
}
```

## API Reference

### JSON Functions
- `DecodeJSON(r *http.Request, dst any) error` - Decode JSON request body into destination struct
- `ErrInvalidContentType` - Error returned for non-JSON content types

### Form Functions
- `ParseForm(r *http.Request) error` - Parse form data in request
- `FormValue(r *http.Request, key string) (string, bool)` - Get form value with existence check
- `FormInt(r *http.Request, key string, defaultValue int) (int, error)` - Get form value as int
- `FormBool(r *http.Request, key string) (bool, error)` - Get form value as bool
- `FormFloat64(r *http.Request, key string, defaultValue float64) (float64, error)` - Get form value as float64

### Query Parameter Functions
- `QueryValue(r *http.Request, key string) (string, bool)` - Get query parameter with existence check
- `QueryInt(r *http.Request, key string, defaultValue int) (int, error)` - Get query parameter as int
- `QueryBool(r *http.Request, key string) (bool, error)` - Get query parameter as bool
- `QueryFloat64(r *http.Request, key string, defaultValue float64) (float64, error)` - Get query parameter as float64

## Type Conversion

### Supported Conversions

The package automatically handles conversion from strings to:

- **int**: Parsed using `strconv.Atoi`
- **bool**: Supports "true"/"false", "1"/"0", "yes"/"no", "on"/"off"  
- **float64**: Parsed using `strconv.ParseFloat`

### Default Values

Functions that accept default values will return the default if:
- The parameter is not present
- The parameter is empty string
- Conversion fails (for type conversion functions)

## Error Handling

### JSON Errors
- `ErrInvalidContentType` - Content-Type is not application/json
- JSON syntax errors from `json.Decoder`
- Unknown field errors (when `DisallowUnknownFields` is enabled)

### Type Conversion Errors
- Invalid integer format
- Invalid float format  
- Invalid boolean format

### Form Parsing Errors
- Malformed form data
- URL encoding errors

## Validation Features

### JSON Validation
- Content-Type validation (must be application/json)
- Unknown field rejection (prevents JSON injection)
- Strict parsing with error details

### Parameter Validation
- Presence checking (distinguish between missing and empty)
- Type validation with clear error messages
- Range validation (implement in your handlers)

## Best Practices

1. **Always validate required fields** after parsing
2. **Use type-safe parameter extraction** instead of direct string access
3. **Handle conversion errors appropriately** - they indicate client errors
4. **Set reasonable default values** for optional parameters
5. **Validate parameter ranges** in your business logic
6. **Use Content-Type validation** to ensure API contracts
7. **Return clear error messages** for validation failures

## Security Considerations

1. **Unknown Field Protection**: JSON decoder disallows unknown fields to prevent injection
2. **Content-Type Validation**: Ensures requests match expected format
3. **Input Validation**: Always validate parsed data before use
4. **Size Limits**: Consider implementing request size limits
5. **Sanitization**: Sanitize string inputs as needed for your application

## Performance Notes

- JSON parsing uses streaming decoder (memory efficient)
- Form parsing is handled by Go's standard library (optimized)
- Query parameter access is O(1) after initial parsing
- Type conversions are lightweight using standard library functions

## Integration

Works well with other go-utils packages:

```go
// Use with logger for request logging
logger.WithField("user_id", userID).Info("Processing request")

// Use with validation from other packages
if err := request.DecodeJSON(r, &req); err != nil {
    logger.WithError(err).Error("JSON decode failed")
    return
}

// Use with httputil/response for consistent error handling
// (See response package for examples)
```