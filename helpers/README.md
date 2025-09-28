# Helpers Package

The `helpers` package provides general utility functions including slice operations, conditional helpers, file system utilities, and struct manipulation. These are commonly used functions that help reduce code duplication across Go projects.

## Features

- **Slice Operations**: Contains, subset checking, and manipulation functions
- **Conditional Helpers**: Ternary operator implementation and default value handling
- **File Operations**: JSON file reading and writing utilities, file system checks
- **Type Utilities**: Generic helper functions for common operations

## Installation

```bash
go get github.com/julianstephens/go-utils/helpers
```

## Usage

### Slice Operations

```go
package main

import (
    "fmt"
    "github.com/julianstephens/go-utils/helpers"
)

func main() {
    mainSlice := []string{"apple", "banana", "cherry", "date", "elderberry"}
    subset := []string{"apple", "cherry"}
    
    // Check if all elements in subset are present in mainSlice
    hasAll := helpers.ContainsAll(mainSlice, subset)
    fmt.Printf("Contains all: %t\n", hasAll) // true
    
    // Check with missing elements
    missingSubset := []string{"apple", "grape"}
    hasAllMissing := helpers.ContainsAll(mainSlice, missingSubset)
    fmt.Printf("Contains all (with missing): %t\n", hasAllMissing) // false
    
    // Empty subset is always contained
    emptySubset := []string{}
    hasEmpty := helpers.ContainsAll(mainSlice, emptySubset)
    fmt.Printf("Contains empty: %t\n", hasEmpty) // true
}
```

### Conditional Helpers

```go
package main

import (
    "fmt"
    "github.com/julianstephens/go-utils/helpers"
)

func main() {
    // Ternary operator - If mimics cond ? vtrue : vfalse
    age := 25
    status := helpers.If(age >= 18, "adult", "minor")
    fmt.Printf("Status: %s\n", status) // "adult"
    
    // Works with any type
    score := 85
    grade := helpers.If(score >= 90, "A", helpers.If(score >= 80, "B", "C"))
    fmt.Printf("Grade: %s\n", grade) // "B"
    
    // With boolean values
    isLoggedIn := true
    message := helpers.If(isLoggedIn, "Welcome back!", "Please log in")
    fmt.Printf("Message: %s\n", message) // "Welcome back!"
    
    // With numeric values
    temperature := 22
    clothing := helpers.If(temperature > 20, "t-shirt", "jacket")
    fmt.Printf("Wear: %s\n", clothing) // "t-shirt"
}
```

### Default Value Handling

```go
package main

import (
    "fmt"
    "github.com/julianstephens/go-utils/helpers"
)

func main() {
    // Default returns defaultVal if val is the zero value for its type
    
    // String example
    emptyString := ""
    result := helpers.Default(emptyString, "default value")
    fmt.Printf("String result: %s\n", result) // "default value"
    
    nonEmptyString := "hello"
    result2 := helpers.Default(nonEmptyString, "default value")
    fmt.Printf("String result 2: %s\n", result2) // "hello"
    
    // Integer example
    zeroInt := 0
    intResult := helpers.Default(zeroInt, 42)
    fmt.Printf("Int result: %d\n", intResult) // 42
    
    nonZeroInt := 10
    intResult2 := helpers.Default(nonZeroInt, 42)
    fmt.Printf("Int result 2: %d\n", intResult2) // 10
    
    // Boolean example
    var zeroBool bool // false is zero value for bool
    boolResult := helpers.Default(zeroBool, true)
    fmt.Printf("Bool result: %t\n", boolResult) // true
}
```

### File Operations

```go
package main

import (
    "fmt"
    "log"
    "os"
    "github.com/julianstephens/go-utils/helpers"
)

type Config struct {
    Host string `json:"host"`
    Port int    `json:"port"`
    SSL  bool   `json:"ssl"`
}

func main() {
    config := Config{
        Host: "localhost",
        Port: 8080,
        SSL:  true,
    }
    
    // Write struct to JSON file
    filename := "/tmp/config.json"
    if err := helpers.WriteJSONFile(filename, config); err != nil {
        log.Fatalf("Failed to write JSON file: %v", err)
    }
    fmt.Println("Config written to file")
    
    // Read JSON file into struct
    var loadedConfig Config
    if err := helpers.ReadJSONFile(filename, &loadedConfig); err != nil {
        log.Fatalf("Failed to read JSON file: %v", err)
    }
    
    fmt.Printf("Loaded config: %+v\n", loadedConfig)
    // Output: Loaded config: {Host:localhost Port:8080 SSL:true}
    
    // Clean up
    os.Remove(filename)
}
```

### File System Utilities

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "github.com/julianstephens/go-utils/helpers"
)

func main() {
    // Check if files and directories exist
    fmt.Printf("Current directory exists: %v\n", helpers.Exists("."))
    fmt.Printf("go.mod exists: %v\n", helpers.Exists("go.mod"))
    fmt.Printf("Non-existent file: %v\n", helpers.Exists("does_not_exist.txt"))
    
    // Create a temporary file to test
    tempFile := filepath.Join(os.TempDir(), "test_file.txt")
    
    // Check before creation
    fmt.Printf("Temp file exists before creation: %v\n", helpers.Exists(tempFile))
    
    // Create the file using Ensure
    if err := helpers.Ensure(tempFile, false); err != nil {
        fmt.Printf("Error creating file: %v\n", err)
        return
    }
    
    // Check after creation
    fmt.Printf("Temp file exists after creation: %v\n", helpers.Exists(tempFile))
    
    // Clean up
    os.Remove(tempFile)
    fmt.Printf("Temp file exists after cleanup: %v\n", helpers.Exists(tempFile))
}
```

### Working with Different Types

```go
package main

import (
    "fmt"
    "github.com/julianstephens/go-utils/helpers"
)

type Person struct {
    Name string
    Age  int
}

func main() {
    // Generic If works with custom types
    person1 := Person{Name: "Alice", Age: 30}
    person2 := Person{Name: "Bob", Age: 25}
    
    older := helpers.If(person1.Age > person2.Age, person1, person2)
    fmt.Printf("Older person: %+v\n", older) // {Name:Alice Age:30}
    
    // ContainsAll works with any comparable type
    numbers := []int{1, 2, 3, 4, 5}
    subset := []int{2, 4}
    
    hasNumbers := helpers.ContainsAll(numbers, subset)
    fmt.Printf("Contains numbers: %t\n", hasNumbers) // true
    
    // Default with custom types
    var emptyPerson Person // zero value
    defaultPerson := Person{Name: "Unknown", Age: 0}
    
    result := helpers.Default(emptyPerson, defaultPerson)
    fmt.Printf("Default person: %+v\n", result) // {Name:Unknown Age:0}
}
```

### Complex Example: Configuration Management

```go
package main

import (
    "fmt"
    "log"
    "os"
    "github.com/julianstephens/go-utils/helpers"
)

type AppConfig struct {
    Environment string   `json:"environment"`
    Debug       bool     `json:"debug"`
    Port        int      `json:"port"`
    Features    []string `json:"features"`
}

func loadConfig(filename string) AppConfig {
    var config AppConfig
    
    // Try to read from file
    if err := helpers.ReadJSONFile(filename, &config); err != nil {
        fmt.Printf("Could not read config file: %v\n", err)
        // Use defaults
        config = AppConfig{}
    }
    
    // Apply defaults for any zero values
    config.Environment = helpers.Default(config.Environment, "development")
    config.Port = helpers.Default(config.Port, 8080)
    
    // Set debug based on environment
    config.Debug = helpers.If(config.Environment == "development", true, config.Debug)
    
    // Ensure essential features are enabled
    essentialFeatures := []string{"auth", "logging"}
    if !helpers.ContainsAll(config.Features, essentialFeatures) {
        fmt.Println("Adding essential features")
        config.Features = append(config.Features, essentialFeatures...)
    }
    
    return config
}

func main() {
    // Create a sample config file
    sampleConfig := AppConfig{
        Environment: "production",
        Debug:       false,
        Port:        9000,
        Features:    []string{"metrics", "auth"},
    }
    
    configFile := "/tmp/app-config.json"
    if err := helpers.WriteJSONFile(configFile, sampleConfig); err != nil {
        log.Fatalf("Failed to write config: %v", err)
    }
    
    // Load and process config
    config := loadConfig(configFile)
    fmt.Printf("Final config: %+v\n", config)
    
    // Demonstrate conditional logic
    logLevel := helpers.If(config.Debug, "debug", "info")
    fmt.Printf("Log level: %s\n", logLevel)
    
    // Clean up
    os.Remove(configFile)
}
```

## API Reference

### Slice Operations
- `ContainsAll[T comparable](mainSlice, subset []T) bool` - Check if all elements in subset are present in mainSlice
- `Difference(a []string, b []string) []string` - Return elements in slice a that are not in slice b
- `DeleteElement[T any](slice []T, index int) []T` - Remove element at specified index from slice

### Conditional Helpers
- `If[T any](cond bool, vtrue T, vfalse T) T` - Ternary operator: returns vtrue if cond is true, otherwise vfalse
- `Default[T any](val T, defaultVal T) T` - Returns defaultVal if val is the zero value for its type, otherwise returns val

### File Operations
- `Exists(path string) bool` - Check if a file or directory exists at the given path
- `Ensure(path string, isDir bool) error` - Create a file or directory if it doesn't exist
- `ReadJSONFile(filename string, v interface{}) error` - Read JSON file and unmarshal into provided struct
- `WriteJSONFile(filename string, v interface{}) error` - Marshal struct to JSON and write to file

### Utility Functions
- `StringPtr(s string) *string` - Return a pointer to the given string
- `MustMarshalJson(v any) []byte` - Marshal value to JSON, panic on error
- `StructToMap(obj any) map[string]any` - Convert struct to map using reflection

## Type Support

### Generic Functions
All generic functions (`If`, `Default`, `ContainsAll`) work with any Go type that satisfies their constraints:

- `If[T any]` - Works with any type
- `Default[T any]` - Works with any type (uses reflection to check for zero values)
- `ContainsAll[T comparable]` - Works with any comparable type (strings, numbers, booleans, arrays, structs with comparable fields)

### Zero Values
The `Default` function recognizes zero values for all Go types:
- Numbers: `0`, `0.0`
- Strings: `""`
- Booleans: `false`
- Pointers: `nil`
- Slices, maps, channels: `nil`
- Structs: All fields are their zero values

## Thread Safety

- `If` and `Default` are pure functions and are thread-safe
- `ContainsAll` is thread-safe for read operations
- File operations (`ReadJSONFile`, `WriteJSONFile`) are not inherently thread-safe and should be synchronized if accessed concurrently

## Best Practices

1. **Use `If` for simple conditional assignments** instead of if-else blocks
2. **Use `Default` to provide fallback values** for potentially empty configuration fields
3. **Use `ContainsAll` for validation** to ensure required elements are present
4. **Handle file operation errors appropriately** - they can fail due to permissions, disk space, etc.
5. **Consider zero values carefully** when using `Default` - sometimes zero is a valid value

## Integration

Works well with other go-utils packages:

```go
// Use with config package for default values
serverPort := helpers.Default(cfg.Server.Port, 8080)

// Use with cliutil for conditional output
message := helpers.If(verbose, "Detailed info", "Basic info")
cliutil.PrintInfo(message)
```