/*
Package config provides a reusable and idiomatic way to load, validate, and access
application configuration for Go projects.

This package supports loading configuration from environment variables and optionally
from YAML or JSON configuration files. It provides struct-based configuration with
support for default values, required fields, and comprehensive type validation.

Basic Usage:

The config package supports both environment variable and file-based configuration:

	package main

	import (
		"log"

		"github.com/julianstephens/go-utils/config"
	)

	type AppConfig struct {
		Port     int    `env:"PORT" default:"8080" required:"false"`
		Host     string `env:"HOST" default:"localhost"`
		Database string `env:"DATABASE_URL" required:"true"`
		Debug    bool   `env:"DEBUG" default:"false"`
	}

	func main() {
		var cfg AppConfig

		// Load from environment variables only
		if err := config.LoadFromEnv(&cfg); err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}

		// Or load from file with env override
		if err := config.LoadFromFileWithEnv(&cfg, "config.yaml"); err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}

		log.Printf("Server starting on %s:%d", cfg.Host, cfg.Port)
	}

File Loading:

The package supports both YAML and JSON configuration files:

	// YAML file (config.yaml)
	port: 3000
	host: "0.0.0.0"
	database_url: "postgres://user:pass@localhost/db"
	debug: true

	// JSON file (config.json)
	{
		"port": 3000,
		"host": "0.0.0.0",
		"database_url": "postgres://user:pass@localhost/db",
		"debug": true
	}

	// Load YAML
	if err := config.LoadFromFile(&cfg, "config.yaml"); err != nil {
		log.Fatal(err)
	}

	// Load JSON
	if err := config.LoadFromFile(&cfg, "config.json"); err != nil {
		log.Fatal(err)
	}

Struct Tags:

The config package uses struct tags to define configuration behavior:

  - `env:"ENV_VAR"` - specifies the environment variable name
  - `default:"value"` - sets a default value if not provided
  - `required:"true"` - marks field as required (fails if missing)
  - `json:"field_name"` - JSON field name (standard json tag)
  - `yaml:"field_name"` - YAML field name (standard yaml tag)

	type DatabaseConfig struct {
		Host     string `env:"DB_HOST" yaml:"host" json:"host" default:"localhost"`
		Port     int    `env:"DB_PORT" yaml:"port" json:"port" default:"5432"`
		Username string `env:"DB_USER" yaml:"username" json:"username" required:"true"`
		Password string `env:"DB_PASS" yaml:"password" json:"password" required:"true"`
		Database string `env:"DB_NAME" yaml:"database" json:"database" required:"true"`
		SSLMode  string `env:"DB_SSL_MODE" yaml:"ssl_mode" json:"ssl_mode" default:"disable"`
	}

Validation:

The package provides comprehensive validation for configuration:

	// Required fields validation
	type Config struct {
		APIKey string `env:"API_KEY" required:"true"`
	}

	var cfg Config
	if err := config.LoadFromEnv(&cfg); err != nil {
		// Will fail if API_KEY environment variable is not set
		log.Fatal(err)
	}

	// Type validation is automatic
	type NumericConfig struct {
		Port    int     `env:"PORT" default:"8080"`
		Timeout float64 `env:"TIMEOUT" default:"30.5"`
		Enabled bool    `env:"ENABLED" default:"true"`
	}

Supported Types:

The config package supports the following types:
  - string
  - int, int8, int16, int32, int64
  - uint, uint8, uint16, uint32, uint64
  - float32, float64
  - bool
  - Pointers to any of the above types
  - Slices of supported types (comma-separated for env vars)

Complex Example:

	type ServerConfig struct {
		// Server settings
		Host string `env:"SERVER_HOST" yaml:"server.host" json:"server_host" default:"localhost"`
		Port int    `env:"SERVER_PORT" yaml:"server.port" json:"server_port" default:"8080"`

		// Database settings
		DatabaseURL string `env:"DATABASE_URL" yaml:"database.url" json:"database_url" required:"true"`

		// Feature flags
		EnableMetrics bool `env:"ENABLE_METRICS" yaml:"features.metrics" json:"enable_metrics" default:"false"`
		EnableTracing bool `env:"ENABLE_TRACING" yaml:"features.tracing" json:"enable_tracing" default:"false"`

		// Optional settings
		LogLevel  string   `env:"LOG_LEVEL" yaml:"logging.level" json:"log_level" default:"info"`
		AdminIPs  []string `env:"ADMIN_IPS" yaml:"security.admin_ips" json:"admin_ips"`
		Timeout   *int     `env:"REQUEST_TIMEOUT" yaml:"server.timeout" json:"request_timeout"`
	}

	func loadConfig() (*ServerConfig, error) {
		var cfg ServerConfig

		// Try to load from file first, then override with env vars
		if err := config.LoadFromFileWithEnv(&cfg, "config.yaml"); err != nil {
			// If file doesn't exist, just load from env
			if err := config.LoadFromEnv(&cfg); err != nil {
				return nil, err
			}
		}

		return &cfg, nil
	}

Error Handling:

The config package provides detailed error messages for debugging:

	// Missing required field
	// Error: required field 'DATABASE_URL' is missing or empty

	// Invalid type conversion
	// Error: failed to parse field 'Port': strconv.Atoi: parsing "invalid": invalid syntax

	// File loading errors
	// Error: failed to read config file 'config.yaml': no such file or directory
*/
package config