package main

import (
	"fmt"
	"log"
	"os"

	"github.com/julianstephens/go-utils/config"
)

// AppConfig demonstrates various configuration features
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
		Name     string   `env:"APP_NAME" yaml:"name" json:"name" default:"go-utils-example"`
		Debug    bool     `env:"DEBUG" yaml:"debug" json:"debug" default:"false"`
		LogLevel string   `env:"LOG_LEVEL" yaml:"log_level" json:"log_level" default:"info"`
		AdminIPs []string `env:"ADMIN_IPS" yaml:"admin_ips" json:"admin_ips"`
	} `yaml:"app" json:"app"`

	// Optional settings (pointers)
	Timeout   *int    `env:"REQUEST_TIMEOUT" yaml:"timeout" json:"timeout"`
	SecretKey *string `env:"SECRET_KEY" yaml:"secret_key" json:"secret_key"`
}

func runConfigExample() {
	fmt.Println("=== Config Package Example ===\n")

	// Example 1: Load from environment variables only
	fmt.Println("1. Loading from environment variables:")

	// Set some environment variables for demonstration
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost/example_db")
	os.Setenv("SERVER_PORT", "3000")
	os.Setenv("DEBUG", "true")
	os.Setenv("ADMIN_IPS", "127.0.0.1,192.168.1.1,10.0.0.1")
	os.Setenv("REQUEST_TIMEOUT", "30")

	var cfg1 AppConfig
	if err := config.LoadFromEnv(&cfg1); err != nil {
		log.Printf("Error loading from env: %v", err)
		return
	}

	fmt.Printf("  Server: %s:%d\n", cfg1.Server.Host, cfg1.Server.Port)
	fmt.Printf("  Database URL: %s\n", cfg1.Database.URL)
	fmt.Printf("  Debug: %t\n", cfg1.App.Debug)
	fmt.Printf("  Admin IPs: %v\n", cfg1.App.AdminIPs)
	if cfg1.Timeout != nil {
		fmt.Printf("  Timeout: %d seconds\n", *cfg1.Timeout)
	}

	// Example 2: Create and load from YAML file
	fmt.Println("\n2. Loading from YAML file:")

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

	// Write YAML file
	yamlFile := "/tmp/config_example.yaml"
	if err := os.WriteFile(yamlFile, []byte(yamlContent), 0644); err != nil {
		log.Printf("Error creating YAML file: %v", err)
		return
	}
	defer os.Remove(yamlFile)

	var cfg2 AppConfig
	if err := config.LoadFromFile(&cfg2, yamlFile); err != nil {
		log.Printf("Error loading from YAML: %v", err)
		return
	}

	fmt.Printf("  Server: %s:%d\n", cfg2.Server.Host, cfg2.Server.Port)
	fmt.Printf("  Database: %s (max_conns: %d, ssl: %s)\n",
		cfg2.Database.URL, cfg2.Database.MaxConns, cfg2.Database.SSLMode)
	fmt.Printf("  App: %s (debug: %t, log_level: %s)\n",
		cfg2.App.Name, cfg2.App.Debug, cfg2.App.LogLevel)
	fmt.Printf("  Features: metrics=%t, tracing=%t\n",
		cfg2.Features.EnableMetrics, cfg2.Features.EnableTracing)
	fmt.Printf("  Admin IPs: %v\n", cfg2.App.AdminIPs)
	if cfg2.SecretKey != nil {
		fmt.Printf("  Secret Key: %s\n", *cfg2.SecretKey)
	}

	// Example 3: Load from file with environment overrides
	fmt.Println("\n3. Loading from file with environment overrides:")

	// Override some values with environment variables
	os.Setenv("SERVER_HOST", "production.example.com")
	os.Setenv("SERVER_PORT", "443")
	os.Setenv("ENABLE_METRICS", "false")
	os.Setenv("LOG_LEVEL", "error")

	var cfg3 AppConfig
	if err := config.LoadFromFileWithEnv(&cfg3, yamlFile); err != nil {
		log.Printf("Error loading from file with env overrides: %v", err)
		return
	}

	fmt.Printf("  Server: %s:%d (overridden by env)\n", cfg3.Server.Host, cfg3.Server.Port)
	fmt.Printf("  Database: %s (from file)\n", cfg3.Database.URL)
	fmt.Printf("  Metrics enabled: %t (overridden by env)\n", cfg3.Features.EnableMetrics)
	fmt.Printf("  Log level: %s (overridden by env)\n", cfg3.App.LogLevel)

	// Example 4: Demonstrate error handling
	fmt.Println("\n4. Error handling demonstration:")

	// Clear required environment variable
	os.Unsetenv("DATABASE_URL")

	var cfg4 AppConfig
	if err := config.LoadFromEnv(&cfg4); err != nil {
		fmt.Printf("  Expected error (missing required field): %v\n", err)
	}

	// Example 5: Using Must functions (would panic on error)
	fmt.Println("\n5. Using Must functions:")

	// Restore required env var for Must function
	os.Setenv("DATABASE_URL", "postgres://localhost/must_example")

	var cfg5 AppConfig
	config.MustLoadFromEnv(&cfg5)
	fmt.Printf("  Successfully loaded with MustLoadFromEnv: %s:%d\n",
		cfg5.Server.Host, cfg5.Server.Port)

	// Clean up environment variables
	for _, env := range []string{
		"DATABASE_URL", "SERVER_PORT", "DEBUG", "ADMIN_IPS", "REQUEST_TIMEOUT",
		"SERVER_HOST", "ENABLE_METRICS", "LOG_LEVEL",
	} {
		os.Unsetenv(env)
	}
}

func main() {
	runConfigExample()
}
