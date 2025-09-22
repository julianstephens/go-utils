# Config Package

The `config` package provides a reusable and idiomatic way to load, validate, and access application configuration for Go projects. It supports loading configuration from environment variables and optionally from YAML or JSON configuration files with struct-based configuration, default values, required fields, and comprehensive type validation.

## Features

- **Environment Variables**: Load configuration from environment variables
- **File Support**: Load from YAML and JSON configuration files
- **Hierarchical Loading**: File-based defaults with environment variable overrides
- **Struct Tags**: Define configuration behavior using struct tags
- **Type Safety**: Support for all basic Go types, slices, and pointers
- **Validation**: Required field validation and type checking
- **Default Values**: Automatic default value assignment

## Installation

```bash
go get github.com/julianstephens/go-utils/config
```

## Usage

### Basic Environment Variable Loading

```go
package main

import (
    "fmt"
    "log"
    "os"
    "github.com/julianstephens/go-utils/config"
)

type AppConfig struct {
    Port     int    `env:"PORT" default:"8080"`
    Host     string `env:"HOST" default:"localhost"`
    Database string `env:"DATABASE_URL" required:"true"`
    Debug    bool   `env:"DEBUG" default:"false"`
}

func main() {
    // Set environment variables
    os.Setenv("DATABASE_URL", "postgres://user:pass@localhost/mydb")
    os.Setenv("PORT", "3000")
    os.Setenv("DEBUG", "true")

    var cfg AppConfig
    if err := config.LoadFromEnv(&cfg); err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    fmt.Printf("Server: %s:%d\n", cfg.Host, cfg.Port)
    fmt.Printf("Database: %s\n", cfg.Database)
    fmt.Printf("Debug: %t\n", cfg.Debug)
}
```

### Complex Configuration Structure

```go
package main

import (
    "fmt"
    "log"
    "os"
    "github.com/julianstephens/go-utils/config"
)

type AppConfig struct {
    // Server configuration
    Server struct {
        Host string `env:"SERVER_HOST" yaml:"host" json:"host" default:"localhost"`
        Port int    `env:"SERVER_PORT" yaml:"port" json:"port" default:"8080"`
    } `yaml:"server" json:"server"`

    // Database configuration
    Database struct {
        URL      string `env:"DATABASE_URL" yaml:"url" json:"url" required:"true"`
        MaxConns int    `env:"DB_MAX_CONNS" yaml:"max_conns" json:"max_conns" default:"10"`
        SSLMode  string `env:"DB_SSL_MODE" yaml:"ssl_mode" json:"ssl_mode" default:"disable"`
    } `yaml:"database" json:"database"`

    // Feature flags
    Features struct {
        EnableMetrics bool `env:"ENABLE_METRICS" yaml:"enable_metrics" json:"enable_metrics" default:"false"`
        EnableTracing bool `env:"ENABLE_TRACING" yaml:"enable_tracing" json:"enable_tracing" default:"false"`
    } `yaml:"features" json:"features"`

    // Application settings
    App struct {
        Name     string   `env:"APP_NAME" yaml:"name" json:"name" default:"my-app"`
        Debug    bool     `env:"DEBUG" yaml:"debug" json:"debug" default:"false"`
        LogLevel string   `env:"LOG_LEVEL" yaml:"log_level" json:"log_level" default:"info"`
        AdminIPs []string `env:"ADMIN_IPS" yaml:"admin_ips" json:"admin_ips"`
    } `yaml:"app" json:"app"`

    // Optional settings (pointers)
    Timeout   *int    `env:"REQUEST_TIMEOUT" yaml:"timeout" json:"timeout"`
    SecretKey *string `env:"SECRET_KEY" yaml:"secret_key" json:"secret_key"`
}

func main() {
    // Set some environment variables
    os.Setenv("DATABASE_URL", "postgres://user:pass@localhost/example_db")
    os.Setenv("SERVER_PORT", "3000")
    os.Setenv("DEBUG", "true")
    os.Setenv("ADMIN_IPS", "127.0.0.1,192.168.1.1,10.0.0.1")
    os.Setenv("REQUEST_TIMEOUT", "30")

    var cfg AppConfig
    if err := config.LoadFromEnv(&cfg); err != nil {
        log.Fatalf("Error loading config: %v", err)
    }

    fmt.Printf("Server: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
    fmt.Printf("Database URL: %s\n", cfg.Database.URL)
    fmt.Printf("Debug: %t\n", cfg.App.Debug)
    fmt.Printf("Admin IPs: %v\n", cfg.App.AdminIPs)
    if cfg.Timeout != nil {
        fmt.Printf("Timeout: %d seconds\n", *cfg.Timeout)
    }
}
```

### Loading from YAML File

```go
package main

import (
    "fmt"
    "log"
    "os"
    "github.com/julianstephens/go-utils/config"
)

func main() {
    // Create a YAML configuration file
    yamlContent := `
server:
  host: "0.0.0.0"
  port: 8080
database:
  url: "postgres://localhost/yaml_db"
  max_conns: 25
  ssl_mode: "require"
features:
  enable_metrics: true
  enable_tracing: false
app:
  name: "yaml-config-app"
  debug: false
  log_level: "warn"
  admin_ips: ["admin1.local", "admin2.local"]
secret_key: "super-secret-key-from-file"
`

    // Write to file
    yamlFile := "/tmp/config.yaml"
    if err := os.WriteFile(yamlFile, []byte(yamlContent), 0644); err != nil {
        log.Fatalf("Failed to write YAML file: %v", err)
    }

    // Load configuration from file
    var cfg AppConfig
    if err := config.LoadFromFile(&cfg, yamlFile); err != nil {
        log.Fatalf("Failed to load from file: %v", err)
    }

    fmt.Printf("Loaded from YAML:\n")
    fmt.Printf("  Server: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
    fmt.Printf("  Database: %s (max_conns: %d)\n", cfg.Database.URL, cfg.Database.MaxConns)
    fmt.Printf("  Features: metrics=%t, tracing=%t\n", cfg.Features.EnableMetrics, cfg.Features.EnableTracing)
    
    // Clean up
    os.Remove(yamlFile)
}
```

### Hierarchical Configuration (File + Environment)

```go
package main

import (
    "fmt"
    "log"
    "os"
    "github.com/julianstephens/go-utils/config"
)

func main() {
    // Create base configuration file
    configContent := `
server:
  host: "localhost"
  port: 8080
database:
  url: "postgres://localhost/dev_db"
  max_conns: 10
app:
  name: "my-app"
  debug: false
  log_level: "info"
`

    configFile := "/tmp/app-config.yaml"
    if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
        log.Fatalf("Failed to write config file: %v", err)
    }

    // Set environment variables to override some settings
    os.Setenv("SERVER_PORT", "9000")                                    // Override port
    os.Setenv("DATABASE_URL", "postgres://prod:secret@db.prod/prod_db") // Override database
    os.Setenv("DEBUG", "true")                                          // Override debug
    os.Setenv("LOG_LEVEL", "debug")                                     // Override log level

    var cfg AppConfig
    // Load from file first, then override with environment variables
    if err := config.LoadFromFileWithEnv(&cfg, configFile); err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    fmt.Printf("Final configuration (file + env overrides):\n")
    fmt.Printf("  Server: %s:%d (port overridden by env)\n", cfg.Server.Host, cfg.Server.Port)
    fmt.Printf("  Database: %s (overridden by env)\n", cfg.Database.URL)
    fmt.Printf("  Debug: %t (overridden by env)\n", cfg.App.Debug)
    fmt.Printf("  Log Level: %s (overridden by env)\n", cfg.App.LogLevel)
    
    // Clean up
    os.Remove(configFile)
}
```

### JSON Configuration

```go
package main

import (
    "fmt"
    "log"
    "os"
    "github.com/julianstephens/go-utils/config"
)

func main() {
    jsonContent := `{
  "server": {
    "host": "0.0.0.0",
    "port": 8080
  },
  "database": {
    "url": "postgres://localhost/json_db",
    "max_conns": 20,
    "ssl_mode": "require"
  },
  "app": {
    "name": "json-config-app",
    "debug": true,
    "log_level": "debug",
    "admin_ips": ["192.168.1.1", "10.0.0.1"]
  },
  "timeout": 45
}`

    jsonFile := "/tmp/config.json"
    if err := os.WriteFile(jsonFile, []byte(jsonContent), 0644); err != nil {
        log.Fatalf("Failed to write JSON file: %v", err)
    }

    var cfg AppConfig
    if err := config.LoadFromFile(&cfg, jsonFile); err != nil {
        log.Fatalf("Failed to load JSON config: %v", err)
    }

    fmt.Printf("Loaded from JSON:\n")
    fmt.Printf("  Server: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
    fmt.Printf("  Database: %s\n", cfg.Database.URL)
    fmt.Printf("  App: %s (debug: %t)\n", cfg.App.Name, cfg.App.Debug)
    
    // Clean up
    os.Remove(jsonFile)
}
```

## Struct Tags

The config package uses struct tags to define configuration behavior:

### Available Tags

- `env:"ENV_VAR"` - Specifies the environment variable name
- `default:"value"` - Sets a default value if not provided
- `required:"true"` - Marks field as required (fails if missing)
- `json:"field_name"` - JSON field name (standard json tag)
- `yaml:"field_name"` - YAML field name (standard yaml tag)

### Example with All Tags

```go
type DatabaseConfig struct {
    Host     string `env:"DB_HOST" yaml:"host" json:"host" default:"localhost"`
    Port     int    `env:"DB_PORT" yaml:"port" json:"port" default:"5432"`
    Username string `env:"DB_USER" yaml:"username" json:"username" required:"true"`
    Password string `env:"DB_PASS" yaml:"password" json:"password" required:"true"`
    Database string `env:"DB_NAME" yaml:"database" json:"database" required:"true"`
    SSLMode  string `env:"DB_SSL_MODE" yaml:"ssl_mode" json:"ssl_mode" default:"disable"`
}
```

## Supported Types

The config package supports the following Go types:

### Basic Types
- `string`
- `bool`
- All integer types: `int`, `int8`, `int16`, `int32`, `int64`
- All unsigned integer types: `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- Floating point types: `float32`, `float64`

### Complex Types
- **Slices**: `[]string`, `[]int`, etc. (comma-separated in env vars)
- **Pointers**: `*string`, `*int`, etc. (optional fields)
- **Nested Structs**: Embedded configuration structures

### Environment Variable Formats

- **Strings**: Direct value
- **Numbers**: Parsed automatically
- **Booleans**: "true"/"false", "1"/"0", "yes"/"no", "on"/"off"
- **Slices**: Comma-separated values (e.g., "item1,item2,item3")

## API Reference

### Functions

- `LoadFromEnv(cfg interface{}) error` - Load configuration from environment variables
- `LoadFromFile(cfg interface{}, filepath string) error` - Load from YAML or JSON file
- `LoadFromFileWithEnv(cfg interface{}, filepath string) error` - Load from file with env overrides
- `MustLoadFromEnv(cfg interface{})` - Load from env (panics on error)
- `MustLoadFromFile(cfg interface{}, filepath string)` - Load from file (panics on error)
- `MustLoadFromFileWithEnv(cfg interface{}, filepath string)` - Load with overrides (panics on error)

### Error Handling

The package provides detailed error messages for:
- Missing required fields
- Type conversion errors
- File reading errors
- Invalid YAML/JSON syntax
- Unsupported field types

## Best Practices

1. **Use hierarchical loading** for flexibility: file-based defaults with environment overrides
2. **Mark sensitive fields as required** to ensure they're explicitly set
3. **Provide sensible defaults** for non-critical configuration
4. **Use nested structs** to organize related configuration
5. **Use pointers for optional fields** to distinguish between zero values and unset values
6. **Validate configuration** after loading to ensure consistency

## Integration with Other Packages

The config package integrates well with other go-utils packages:

```go
// Use with logger
logger.SetLogLevel(cfg.App.LogLevel)

// Use with cliutil for CLI configuration
configFile := cliutil.GetFlagValue(os.Args, "--config", "config.yaml")
config.LoadFromFile(&cfg, configFile)
```