# JSON Utilities Package

The `jsonutil` package provides enhanced JSON marshaling and unmarshaling with error context, formatting options, stream processing, and strict decoding support. It extends Go's standard `encoding/json` package with additional safety and convenience features.

## Features

- **Enhanced Marshaling**: Pretty-printing, HTML escaping control, and custom formatting
- **Strict Unmarshaling**: Disallow unknown fields and number type control
- **Stream Processing**: Encoder and decoder with custom options
- **Error Context**: Better error messages with additional context
- **Type Safety**: Strict type validation and conversion controls

## Installation

```bash
go get github.com/julianstephens/go-utils/jsonutil
```

## Usage

### Basic Marshaling and Unmarshaling

```go
package main

import (
    "fmt"
    "log"
    "github.com/julianstephens/go-utils/jsonutil"
)

type User struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Active   bool   `json:"active"`
    Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func main() {
    user := User{
        ID:     1,
        Name:   "John Doe",
        Email:  "john@example.com",
        Active: true,
        Metadata: map[string]interface{}{
            "department": "engineering",
            "level":      "senior",
        },
    }

    // Basic marshaling
    data, err := jsonutil.Marshal(user)
    if err != nil {
        log.Fatalf("Marshal failed: %v", err)
    }
    fmt.Printf("JSON: %s\n", data)

    // Basic unmarshaling
    var parsedUser User
    if err := jsonutil.Unmarshal(data, &parsedUser); err != nil {
        log.Fatalf("Unmarshal failed: %v", err)
    }
    fmt.Printf("Parsed: %+v\n", parsedUser)
}
```

### Pretty Printing

```go
package main

import (
    "fmt"
    "log"
    "github.com/julianstephens/go-utils/jsonutil"
)

func main() {
    user := User{
        ID:     1,
        Name:   "John Doe",
        Email:  "john@example.com",
        Active: true,
    }

    // Pretty print with default indentation
    prettyJSON, err := jsonutil.MarshalPretty(user)
    if err != nil {
        log.Fatalf("MarshalPretty failed: %v", err)
    }
    fmt.Printf("Pretty JSON:\n%s\n", prettyJSON)

    // Custom formatting options
    opts := jsonutil.MarshalOptions{
        Indent:     "    ", // 4 spaces
        Prefix:     "",
        EscapeHTML: false,
    }

    customJSON, err := jsonutil.MarshalWithOptions(user, opts)
    if err != nil {
        log.Fatalf("MarshalWithOptions failed: %v", err)
    }
    fmt.Printf("Custom formatted JSON:\n%s\n", customJSON)
}
```

### Strict Unmarshaling

```go
package main

import (
    "fmt"
    "log"
    "github.com/julianstephens/go-utils/jsonutil"
)

type StrictUser struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func main() {
    // JSON with extra field
    jsonWithExtra := `{
        "id": 1,
        "name": "John",
        "email": "john@example.com",
        "unknown_field": "value"
    }`

    // Regular unmarshal (ignores unknown fields)
    var user1 StrictUser
    if err := jsonutil.Unmarshal([]byte(jsonWithExtra), &user1); err != nil {
        log.Printf("Regular unmarshal failed: %v", err)
    } else {
        fmt.Printf("Regular unmarshal success: %+v\n", user1)
    }

    // Strict unmarshal (fails on unknown fields)
    opts := jsonutil.UnmarshalOptions{
        DisallowUnknownFields: true,
        UseNumber:            false,
    }

    var user2 StrictUser
    if err := jsonutil.UnmarshalWithOptions([]byte(jsonWithExtra), &user2, opts); err != nil {
        fmt.Printf("Strict unmarshal failed (expected): %v\n", err)
    } else {
        fmt.Printf("Strict unmarshal success: %+v\n", user2)
    }

    // JSON without extra fields works fine
    cleanJSON := `{"id": 1, "name": "John"}`
    var user3 StrictUser
    if err := jsonutil.UnmarshalWithOptions([]byte(cleanJSON), &user3, opts); err != nil {
        log.Printf("Clean unmarshal failed: %v", err)
    } else {
        fmt.Printf("Clean unmarshal success: %+v\n", user3)
    }
}
```

### Number Handling

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "github.com/julianstephens/go-utils/jsonutil"
)

type NumberData struct {
    ID     json.Number `json:"id"`
    Score  json.Number `json:"score"`
    Active bool        `json:"active"`
}

func main() {
    jsonData := `{
        "id": 123456789012345,
        "score": 98.5,
        "active": true
    }`

    // Use json.Number to preserve precision
    opts := jsonutil.UnmarshalOptions{
        UseNumber:             true,
        DisallowUnknownFields: false,
    }

    var data NumberData
    if err := jsonutil.UnmarshalWithOptions([]byte(jsonData), &data, opts); err != nil {
        log.Fatalf("Unmarshal failed: %v", err)
    }

    fmt.Printf("ID as string: %s\n", string(data.ID))
    fmt.Printf("Score as string: %s\n", string(data.Score))

    // Convert to specific types when needed
    if id, err := data.ID.Int64(); err == nil {
        fmt.Printf("ID as int64: %d\n", id)
    }

    if score, err := data.Score.Float64(); err == nil {
        fmt.Printf("Score as float64: %f\n", score)
    }
}
```

### Stream Processing

```go
package main

import (
    "bytes"
    "fmt"
    "log"
    "github.com/julianstephens/go-utils/jsonutil"
)

func main() {
    users := []User{
        {ID: 1, Name: "Alice", Email: "alice@example.com", Active: true},
        {ID: 2, Name: "Bob", Email: "bob@example.com", Active: false},
        {ID: 3, Name: "Charlie", Email: "charlie@example.com", Active: true},
    }

    // Stream encoding
    var buf bytes.Buffer
    encoder := jsonutil.NewEncoderWithOptions(&buf, jsonutil.EncoderOptions{
        Indent:     "  ",
        EscapeHTML: false,
    })

    for _, user := range users {
        if err := encoder.Encode(user); err != nil {
            log.Fatalf("Encode failed: %v", err)
        }
    }

    fmt.Printf("Encoded stream:\n%s", buf.String())

    // Stream decoding
    decoder := jsonutil.NewDecoderWithOptions(&buf, jsonutil.UnmarshalOptions{
        DisallowUnknownFields: false,
        UseNumber:            false,
    })

    var decodedUsers []User
    for decoder.More() {
        var user User
        if err := decoder.Decode(&user); err != nil {
            log.Fatalf("Decode failed: %v", err)
        }
        decodedUsers = append(decodedUsers, user)
    }

    fmt.Printf("Decoded %d users\n", len(decodedUsers))
    for _, user := range decodedUsers {
        fmt.Printf("  %+v\n", user)
    }
}
```

### Configuration and Settings

```go
package main

import (
    "fmt"
    "log"
    "github.com/julianstephens/go-utils/jsonutil"
)

type AppConfig struct {
    Server struct {
        Host string `json:"host"`
        Port int    `json:"port"`
    } `json:"server"`
    Database struct {
        URL     string `json:"url"`
        Timeout int    `json:"timeout"`
    } `json:"database"`
    Features map[string]bool `json:"features"`
}

func main() {
    config := AppConfig{}
    config.Server.Host = "localhost"
    config.Server.Port = 8080
    config.Database.URL = "postgres://localhost/mydb"
    config.Database.Timeout = 30
    config.Features = map[string]bool{
        "auth":     true,
        "metrics":  true,
        "logging":  true,
        "caching":  false,
    }

    // Save configuration with pretty formatting
    configJSON, err := jsonutil.MarshalWithOptions(config, jsonutil.MarshalOptions{
        Indent:     "  ",
        EscapeHTML: false,
    })
    if err != nil {
        log.Fatalf("Failed to marshal config: %v", err)
    }

    fmt.Printf("Configuration:\n%s\n", configJSON)

    // Load configuration with strict validation
    var loadedConfig AppConfig
    strictOpts := jsonutil.UnmarshalOptions{
        DisallowUnknownFields: true, // Fail if config has unexpected fields
    }

    if err := jsonutil.UnmarshalWithOptions(configJSON, &loadedConfig, strictOpts); err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    fmt.Printf("Loaded config successfully: %+v\n", loadedConfig)
}
```

### API Response Processing

```go
package main

import (
    "fmt"
    "log"
    "github.com/julianstephens/go-utils/jsonutil"
)

type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data"`
    Error   *string     `json:"error,omitempty"`
    Meta    struct {
        Page       int `json:"page"`
        PerPage    int `json:"per_page"`
        Total      int `json:"total"`
        TotalPages int `json:"total_pages"`
    } `json:"meta,omitempty"`
}

func processAPIResponse(jsonData []byte) error {
    var response APIResponse
    
    // Use strict unmarshaling for API responses to catch schema changes
    opts := jsonutil.UnmarshalOptions{
        DisallowUnknownFields: true,
        UseNumber:            false,
    }

    if err := jsonutil.UnmarshalWithOptions(jsonData, &response, opts); err != nil {
        return fmt.Errorf("failed to parse API response: %w", err)
    }

    if !response.Success {
        if response.Error != nil {
            return fmt.Errorf("API error: %s", *response.Error)
        }
        return fmt.Errorf("API request failed with no error message")
    }

    fmt.Printf("API call successful: %+v\n", response)
    return nil
}

func main() {
    // Successful response
    successJSON := `{
        "success": true,
        "data": {"message": "Operation completed"},
        "meta": {
            "page": 1,
            "per_page": 10,
            "total": 25,
            "total_pages": 3
        }
    }`

    if err := processAPIResponse([]byte(successJSON)); err != nil {
        log.Printf("Error: %v", err)
    }

    // Error response
    errorJSON := `{
        "success": false,
        "error": "Invalid authentication token"
    }`

    if err := processAPIResponse([]byte(errorJSON)); err != nil {
        log.Printf("Expected error: %v", err)
    }
}
```

## Configuration Options

### MarshalOptions

```go
type MarshalOptions struct {
    Indent     string // Indentation string (e.g., "  ", "\t")
    Prefix     string // Prefix for each line
    EscapeHTML bool   // Whether to escape HTML characters
}
```

### UnmarshalOptions

```go
type UnmarshalOptions struct {
    DisallowUnknownFields bool // Reject JSON with unknown fields
    UseNumber            bool // Use json.Number instead of float64
}
```

### EncoderOptions

```go
type EncoderOptions struct {
    Indent     string // Indentation for streaming
    Prefix     string // Line prefix for streaming  
    EscapeHTML bool   // HTML escaping control
}
```

## API Reference

### Basic Functions
- `Marshal(v interface{}) ([]byte, error)` - Standard JSON marshaling
- `MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)` - Marshal with custom indentation
- `MarshalWithOptions(v interface{}, opts *MarshalOptions) ([]byte, error)` - Marshal with custom options
- `Unmarshal(data []byte, v interface{}) error` - Standard JSON unmarshaling
- `UnmarshalStrict(data []byte, v interface{}) error` - Strict unmarshaling (disallow unknown fields)
- `UnmarshalWithOptions(data []byte, v interface{}, opts *UnmarshalOptions) error` - Unmarshal with options

### Stream Processing
- `EncodeWriter(w io.Writer, v interface{}, opts *EncoderOptions) error` - Encode directly to writer
- `DecodeReader(r io.Reader, v interface{}, opts *DecoderOptions) error` - Decode directly from reader
- `DecodeReaderStrict(r io.Reader, v interface{}) error` - Strict decode from reader

### Utility Functions
- `Valid(data []byte) bool` - Check if JSON is valid
- `Compact(dst *bytes.Buffer, src []byte) error` - Remove whitespace from JSON
- `Indent(dst *bytes.Buffer, src []byte, prefix, indent string) error` - Add indentation to JSON
- `HTMLEscape(dst *bytes.Buffer, src []byte)` - Escape HTML characters in JSON

## Error Handling

The package provides enhanced error context:

- **Marshal errors**: Include type information and field details
- **Unmarshal errors**: Include line/column information when possible
- **Validation errors**: Clear messages for unknown fields or type mismatches
- **Stream errors**: Context about position in stream

```go
// Example error handling
if err := jsonutil.Unmarshal(data, &result); err != nil {
    // Errors include helpful context
    log.Printf("JSON parsing failed: %v", err)
}
```

## Best Practices

1. **Use strict unmarshaling for APIs** to catch schema changes early
2. **Use json.Number for large integers** to avoid precision loss
3. **Enable HTML escaping for web contexts** to prevent XSS
4. **Use pretty printing for configuration files** to improve readability
5. **Handle streaming for large datasets** to reduce memory usage
6. **Validate JSON structure** with `DisallowUnknownFields` in production
7. **Use appropriate indentation** for different contexts (2 spaces for web, 4 for config files)

## Thread Safety

- All functions are thread-safe
- Encoders and decoders are not thread-safe and should not be shared between goroutines
- Multiple encoders/decoders can be used concurrently

## Integration

Works well with other go-utils packages:

```go
// Use with config package
configJSON, _ := jsonutil.MarshalPretty(cfg)
logger.WithField("config", string(configJSON)).Info("Configuration loaded")

// Use with httputil for API responses
response := APIResponse{Success: true, Data: result}
jsonData, _ := jsonutil.Marshal(response)
httputil.WriteJSON(w, jsonData)
```