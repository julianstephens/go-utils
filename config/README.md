# Config Package

The `config` package provides a reusable and idiomatic way to load, validate, and access application configuration for Go projects. It supports loading configuration from environment variables and optionally from YAML or JSON configuration files with struct-based configuration, default values, required fields, and comprehensive type validation.

## Features

- **Environment Variables**: Load from environment variables
- **File Support**: YAML and JSON file loading
- **Hierarchical Loading**: File defaults with environment overrides
- **Struct Tags**: Configuration via struct tags
- **Type Safety**: Support for all basic Go types, slices, and pointers
- **Validation**: Required field and type validation
- **Default Values**: Automatic defaults

## Installation

```bash
go get github.com/julianstephens/go-utils/config
```

## Usage

### Basic Environment Variable Loading

```go
package main

import (
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
    os.Setenv("DATABASE_URL", "postgres://user:pass@localhost/mydb")
    os.Setenv("PORT", "3000")
    os.Setenv("DEBUG", "true")

    var cfg AppConfig
    if err := config.LoadFromEnv(&cfg); err != nil {
        log.Fatal(err)
    }
}
```

### Complex Configuration Structure

```go
package main

import (
    "log"
    "os"
    "github.com/julianstephens/go-utils/config"
)

type AppConfig struct {
    Server struct {
        Host string `env:"SERVER_HOST" default:"localhost"`
        Port int    `env:"SERVER_PORT" default:"8080"`
    } `yaml:"server" json:"server"`

    Database struct {
        URL      string `env:"DATABASE_URL" required:"true"`
        MaxConns int    `env:"DB_MAX_CONNS" default:"10"`
    } `yaml:"database" json:"database"`

    App struct {
        Name     string   `env:"APP_NAME" default:"my-app"`
        Debug    bool     `env:"DEBUG" default:"false"`
        AdminIPs []string `env:"ADMIN_IPS"`
    } `yaml:"app" json:"app"`
}

func main() {
    os.Setenv("DATABASE_URL", "postgres://user:pass@localhost/db")
    os.Setenv("SERVER_PORT", "3000")
    os.Setenv("ADMIN_IPS", "127.0.0.1,192.168.1.1")

    var cfg AppConfig
    if err := config.LoadFromEnv(&cfg); err != nil {
        log.Fatal(err)
    }
}
```

### Loading from YAML File

```go
package main

import (
    "log"
    "os"
    "github.com/julianstephens/go-utils/config"
)

func main() {
    yamlContent := `
server:
  host: "0.0.0.0"
  port: 8080
database:
  url: "postgres://localhost/mydb"
  max_conns: 25
app:
  name: "my-app"
  debug: false
`

    yamlFile := "/tmp/config.yaml"
    os.WriteFile(yamlFile, []byte(yamlContent), 0644)

    var cfg AppConfig
    if err := config.LoadFromFile(&cfg, yamlFile); err != nil {
        log.Fatal(err)
    }

    os.Remove(yamlFile)
}
```

### Hierarchical Configuration (File + Environment)

```go
package main

import (
    "log"
    "os"
    "github.com/julianstephens/go-utils/config"
)

func main() {
    configContent := `
server:
  host: "localhost"
  port: 8080
database:
  url: "postgres://localhost/dev_db"
  max_conns: 10
`

    configFile := "/tmp/app-config.yaml"
    os.WriteFile(configFile, []byte(configContent), 0644)

    // Environment variables override file settings
    os.Setenv("SERVER_PORT", "9000")
    os.Setenv("DATABASE_URL", "postgres://prod:secret@db.prod/prod_db")

    var cfg AppConfig
    if err := config.LoadFromFileWithEnv(&cfg, configFile); err != nil {
        log.Fatal(err)
    }

    os.Remove(configFile)
}
```

### JSON Configuration

```go
package main

import (
    "log"
    "os"
    "github.com/julianstephens/go-utils/config"
)

func main() {
    jsonContent := `{
  "server": {"host": "0.0.0.0", "port": 8080},
  "database": {"url": "postgres://localhost/mydb", "max_conns": 20},
  "app": {"name": "my-app", "debug": true}
}`

    jsonFile := "/tmp/config.json"
    os.WriteFile(jsonFile, []byte(jsonContent), 0644)

    var cfg AppConfig
    if err := config.LoadFromFile(&cfg, jsonFile); err != nil {
        log.Fatal(err)
    }

    os.Remove(jsonFile)
}
```

## Struct Tags

### Available Tags
- `env:"ENV_VAR"` - Environment variable name
- `default:"value"` - Default value if not provided
- `required:"true"` - Marks field as required
- `json:"field_name"` - JSON field name
- `yaml:"field_name"` - YAML field name

### Example
```go
type DatabaseConfig struct {
    Host     string `env:"DB_HOST" default:"localhost"`
    Port     int    `env:"DB_PORT" default:"5432"`
    Username string `env:"DB_USER" required:"true"`
    Password string `env:"DB_PASS" required:"true"`
}
```

## Supported Types

### Basic Types
- Strings, booleans, integers (`int`, `int8`...`int64`)
- Unsigned integers (`uint`, `uint8`...`uint64`)
- Floats (`float32`, `float64`)

### Complex Types
- **Slices**: `[]string`, `[]int`, etc. (comma-separated in env vars)
- **Pointers**: `*string`, `*int`, etc. (optional fields)
- **Nested Structs**: Embedded configuration structures

### Environment Variable Parsing
- **Strings**: Direct value
- **Numbers**: Automatic parsing
- **Booleans**: "true"/"false", "1"/"0", "yes"/"no"
- **Slices**: Comma-separated ("item1,item2,item3")

## API Reference

### Loading Functions
- `LoadFromEnv(cfg interface{}) error` - Load from environment variables
- `LoadFromFile(cfg interface{}, filepath string) error` - Load from YAML/JSON
- `LoadFromFileWithEnv(cfg interface{}, filepath string) error` - File with env overrides
- `MustLoadFromEnv(cfg interface{})` - Load or panic
- `MustLoadFromFile(cfg interface{}, filepath string)` - Load or panic
- `MustLoadFromFileWithEnv(cfg interface{}, filepath string)` - Load or panic

### Error Handling
Provides detailed errors for missing required fields, type conversion issues, file errors, and invalid syntax.

## Best Practices

1. Use hierarchical loading for flexibility (file + env overrides)
2. Mark sensitive fields as required
3. Provide sensible defaults for non-critical config
4. Use nested structs for organization
5. Use pointers for optional fields
6. Validate configuration after loading

## Integration

Works well with other go-utils packages:
- **logger**: Set log level from config
- **cliutil**: Load config file from CLI flags