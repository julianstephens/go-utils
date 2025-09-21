package config_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/julianstephens/go-utils/config"
)

// Test configuration structs
type BasicConfig struct {
	Port     int    `env:"PORT" default:"8080"`
	Host     string `env:"HOST" default:"localhost"`
	Database string `env:"DATABASE_URL" required:"true"`
	Debug    bool   `env:"DEBUG" default:"false"`
}

type ComplexConfig struct {
	// Basic types
	StringField string  `env:"STRING_FIELD" default:"default_string"`
	IntField    int     `env:"INT_FIELD" default:"42"`
	BoolField   bool    `env:"BOOL_FIELD" default:"true"`
	FloatField  float64 `env:"FLOAT_FIELD" default:"3.14"`

	// Pointer types
	StringPtr *string `env:"STRING_PTR"`
	IntPtr    *int    `env:"INT_PTR" default:"100"`

	// Slice types
	StringSlice []string `env:"STRING_SLICE" default:"a,b,c"`
	IntSlice    []int    `env:"INT_SLICE" default:"1,2,3"`

	// Required fields
	RequiredField string `env:"REQUIRED_FIELD" required:"true"`
}

type FileConfig struct {
	Server struct {
		Host string `yaml:"host" json:"host"`
		Port int    `yaml:"port" json:"port"`
	} `yaml:"server" json:"server"`
	Database struct {
		URL      string `yaml:"url" json:"url"`
		MaxConns int    `yaml:"max_conns" json:"max_conns"`
	} `yaml:"database" json:"database"`
	Features struct {
		EnableMetrics bool `yaml:"enable_metrics" json:"enable_metrics"`
		EnableTracing bool `yaml:"enable_tracing" json:"enable_tracing"`
	} `yaml:"features" json:"features"`
}

func TestLoadFromEnv(t *testing.T) {
	// Clean environment
	cleanEnv := func() {
		os.Unsetenv("PORT")
		os.Unsetenv("HOST")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("DEBUG")
	}

	t.Run("default values", func(t *testing.T) {
		cleanEnv()
		os.Setenv("DATABASE_URL", "postgres://localhost/test")
		defer cleanEnv()

		var cfg BasicConfig
		err := config.LoadFromEnv(&cfg)
		if err != nil {
			t.Fatalf("LoadFromEnv failed: %v", err)
		}

		if cfg.Port != 8080 {
			t.Errorf("Expected Port to be 8080, got %d", cfg.Port)
		}
		if cfg.Host != "localhost" {
			t.Errorf("Expected Host to be 'localhost', got '%s'", cfg.Host)
		}
		if cfg.Database != "postgres://localhost/test" {
			t.Errorf("Expected Database to be 'postgres://localhost/test', got '%s'", cfg.Database)
		}
		if cfg.Debug != false {
			t.Errorf("Expected Debug to be false, got %t", cfg.Debug)
		}
	})

	t.Run("environment overrides", func(t *testing.T) {
		cleanEnv()
		os.Setenv("PORT", "3000")
		os.Setenv("HOST", "0.0.0.0")
		os.Setenv("DATABASE_URL", "postgres://prod/db")
		os.Setenv("DEBUG", "true")
		defer cleanEnv()

		var cfg BasicConfig
		err := config.LoadFromEnv(&cfg)
		if err != nil {
			t.Fatalf("LoadFromEnv failed: %v", err)
		}

		if cfg.Port != 3000 {
			t.Errorf("Expected Port to be 3000, got %d", cfg.Port)
		}
		if cfg.Host != "0.0.0.0" {
			t.Errorf("Expected Host to be '0.0.0.0', got '%s'", cfg.Host)
		}
		if cfg.Database != "postgres://prod/db" {
			t.Errorf("Expected Database to be 'postgres://prod/db', got '%s'", cfg.Database)
		}
		if cfg.Debug != true {
			t.Errorf("Expected Debug to be true, got %t", cfg.Debug)
		}
	})

	t.Run("required field missing", func(t *testing.T) {
		cleanEnv()
		defer cleanEnv()

		var cfg BasicConfig
		err := config.LoadFromEnv(&cfg)
		if err == nil {
			t.Fatal("Expected error for missing required field")
		}

		expectedError := "required field 'Database' (env: DATABASE_URL) is missing or empty"
		if err.Error() != expectedError {
			t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("invalid pointer", func(t *testing.T) {
		var cfg BasicConfig
		err := config.LoadFromEnv(cfg) // Not a pointer
		if err == nil {
			t.Fatal("Expected error for non-pointer argument")
		}

		expectedError := "config must be a pointer to a struct"
		if err.Error() != expectedError {
			t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
		}
	})
}

func TestComplexTypes(t *testing.T) {
	// Clean environment
	cleanEnv := func() {
		for _, env := range []string{
			"STRING_FIELD", "INT_FIELD", "BOOL_FIELD", "FLOAT_FIELD",
			"STRING_PTR", "INT_PTR", "STRING_SLICE", "INT_SLICE", "REQUIRED_FIELD",
		} {
			os.Unsetenv(env)
		}
	}

	t.Run("complex types with defaults", func(t *testing.T) {
		cleanEnv()
		os.Setenv("REQUIRED_FIELD", "required_value")
		defer cleanEnv()

		var cfg ComplexConfig
		err := config.LoadFromEnv(&cfg)
		if err != nil {
			t.Fatalf("LoadFromEnv failed: %v", err)
		}

		// Check basic types
		if cfg.StringField != "default_string" {
			t.Errorf("Expected StringField to be 'default_string', got '%s'", cfg.StringField)
		}
		if cfg.IntField != 42 {
			t.Errorf("Expected IntField to be 42, got %d", cfg.IntField)
		}
		if cfg.BoolField != true {
			t.Errorf("Expected BoolField to be true, got %t", cfg.BoolField)
		}
		if cfg.FloatField != 3.14 {
			t.Errorf("Expected FloatField to be 3.14, got %f", cfg.FloatField)
		}

		// Check pointer types
		if cfg.StringPtr != nil {
			t.Errorf("Expected StringPtr to be nil, got %v", cfg.StringPtr)
		}
		if cfg.IntPtr == nil || *cfg.IntPtr != 100 {
			t.Errorf("Expected IntPtr to be 100, got %v", cfg.IntPtr)
		}

		// Check slice types
		expectedStringSlice := []string{"a", "b", "c"}
		if !reflect.DeepEqual(cfg.StringSlice, expectedStringSlice) {
			t.Errorf("Expected StringSlice to be %v, got %v", expectedStringSlice, cfg.StringSlice)
		}

		expectedIntSlice := []int{1, 2, 3}
		if !reflect.DeepEqual(cfg.IntSlice, expectedIntSlice) {
			t.Errorf("Expected IntSlice to be %v, got %v", expectedIntSlice, cfg.IntSlice)
		}

		if cfg.RequiredField != "required_value" {
			t.Errorf("Expected RequiredField to be 'required_value', got '%s'", cfg.RequiredField)
		}
	})

	t.Run("complex types with overrides", func(t *testing.T) {
		cleanEnv()
		os.Setenv("STRING_FIELD", "override_string")
		os.Setenv("INT_FIELD", "999")
		os.Setenv("BOOL_FIELD", "false")
		os.Setenv("FLOAT_FIELD", "2.71")
		os.Setenv("STRING_PTR", "pointer_value")
		os.Setenv("INT_PTR", "200")
		os.Setenv("STRING_SLICE", "x,y,z")
		os.Setenv("INT_SLICE", "10,20,30")
		os.Setenv("REQUIRED_FIELD", "required_override")
		defer cleanEnv()

		var cfg ComplexConfig
		err := config.LoadFromEnv(&cfg)
		if err != nil {
			t.Fatalf("LoadFromEnv failed: %v", err)
		}

		if cfg.StringField != "override_string" {
			t.Errorf("Expected StringField to be 'override_string', got '%s'", cfg.StringField)
		}
		if cfg.IntField != 999 {
			t.Errorf("Expected IntField to be 999, got %d", cfg.IntField)
		}
		if cfg.BoolField != false {
			t.Errorf("Expected BoolField to be false, got %t", cfg.BoolField)
		}
		if cfg.FloatField != 2.71 {
			t.Errorf("Expected FloatField to be 2.71, got %f", cfg.FloatField)
		}

		if cfg.StringPtr == nil || *cfg.StringPtr != "pointer_value" {
			t.Errorf("Expected StringPtr to be 'pointer_value', got %v", cfg.StringPtr)
		}
		if cfg.IntPtr == nil || *cfg.IntPtr != 200 {
			t.Errorf("Expected IntPtr to be 200, got %v", cfg.IntPtr)
		}

		expectedStringSlice := []string{"x", "y", "z"}
		if !reflect.DeepEqual(cfg.StringSlice, expectedStringSlice) {
			t.Errorf("Expected StringSlice to be %v, got %v", expectedStringSlice, cfg.StringSlice)
		}

		expectedIntSlice := []int{10, 20, 30}
		if !reflect.DeepEqual(cfg.IntSlice, expectedIntSlice) {
			t.Errorf("Expected IntSlice to be %v, got %v", expectedIntSlice, cfg.IntSlice)
		}
	})
}

func TestLoadFromFile(t *testing.T) {
	// Create temporary directory for test files
	tempDir := t.TempDir()

	t.Run("YAML file", func(t *testing.T) {
		yamlContent := `
server:
  host: "127.0.0.1"
  port: 9000
database:
  url: "postgres://localhost/yaml_test"
  max_conns: 50
features:
  enable_metrics: true
  enable_tracing: false
`
		yamlFile := filepath.Join(tempDir, "config.yaml")
		if err := os.WriteFile(yamlFile, []byte(yamlContent), 0644); err != nil {
			t.Fatalf("Failed to create YAML test file: %v", err)
		}

		var cfg FileConfig
		err := config.LoadFromFile(&cfg, yamlFile)
		if err != nil {
			t.Fatalf("LoadFromFile failed: %v", err)
		}

		if cfg.Server.Host != "127.0.0.1" {
			t.Errorf("Expected Server.Host to be '127.0.0.1', got '%s'", cfg.Server.Host)
		}
		if cfg.Server.Port != 9000 {
			t.Errorf("Expected Server.Port to be 9000, got %d", cfg.Server.Port)
		}
		if cfg.Database.URL != "postgres://localhost/yaml_test" {
			t.Errorf("Expected Database.URL to be 'postgres://localhost/yaml_test', got '%s'", cfg.Database.URL)
		}
		if cfg.Database.MaxConns != 50 {
			t.Errorf("Expected Database.MaxConns to be 50, got %d", cfg.Database.MaxConns)
		}
		if cfg.Features.EnableMetrics != true {
			t.Errorf("Expected Features.EnableMetrics to be true, got %t", cfg.Features.EnableMetrics)
		}
		if cfg.Features.EnableTracing != false {
			t.Errorf("Expected Features.EnableTracing to be false, got %t", cfg.Features.EnableTracing)
		}
	})

	t.Run("JSON file", func(t *testing.T) {
		jsonContent := `{
  "server": {
    "host": "0.0.0.0",
    "port": 8080
  },
  "database": {
    "url": "postgres://localhost/json_test",
    "max_conns": 25
  },
  "features": {
    "enable_metrics": false,
    "enable_tracing": true
  }
}`
		jsonFile := filepath.Join(tempDir, "config.json")
		if err := os.WriteFile(jsonFile, []byte(jsonContent), 0644); err != nil {
			t.Fatalf("Failed to create JSON test file: %v", err)
		}

		var cfg FileConfig
		err := config.LoadFromFile(&cfg, jsonFile)
		if err != nil {
			t.Fatalf("LoadFromFile failed: %v", err)
		}

		if cfg.Server.Host != "0.0.0.0" {
			t.Errorf("Expected Server.Host to be '0.0.0.0', got '%s'", cfg.Server.Host)
		}
		if cfg.Server.Port != 8080 {
			t.Errorf("Expected Server.Port to be 8080, got %d", cfg.Server.Port)
		}
		if cfg.Database.URL != "postgres://localhost/json_test" {
			t.Errorf("Expected Database.URL to be 'postgres://localhost/json_test', got '%s'", cfg.Database.URL)
		}
		if cfg.Database.MaxConns != 25 {
			t.Errorf("Expected Database.MaxConns to be 25, got %d", cfg.Database.MaxConns)
		}
		if cfg.Features.EnableMetrics != false {
			t.Errorf("Expected Features.EnableMetrics to be false, got %t", cfg.Features.EnableMetrics)
		}
		if cfg.Features.EnableTracing != true {
			t.Errorf("Expected Features.EnableTracing to be true, got %t", cfg.Features.EnableTracing)
		}
	})

	t.Run("unsupported file format", func(t *testing.T) {
		unsupportedFile := filepath.Join(tempDir, "config.xml")
		if err := os.WriteFile(unsupportedFile, []byte("<config></config>"), 0644); err != nil {
			t.Fatalf("Failed to create unsupported test file: %v", err)
		}

		var cfg FileConfig
		err := config.LoadFromFile(&cfg, unsupportedFile)
		if err == nil {
			t.Fatal("Expected error for unsupported file format")
		}

		expectedError := "unsupported config file format '.xml', supported formats: .yaml, .yml, .json"
		if err.Error() != expectedError {
			t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("non-existent file", func(t *testing.T) {
		var cfg FileConfig
		err := config.LoadFromFile(&cfg, "non-existent.yaml")
		if err == nil {
			t.Fatal("Expected error for non-existent file")
		}
	})
}

func TestLoadFromFileWithEnv(t *testing.T) {
	tempDir := t.TempDir()

	// Create YAML file
	yamlContent := `
server:
  host: "file-host"
  port: 3000
database:
  url: "postgres://file-db/test"
  max_conns: 10
features:
  enable_metrics: false
  enable_tracing: false
`
	yamlFile := filepath.Join(tempDir, "config.yaml")
	if err := os.WriteFile(yamlFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create YAML test file: %v", err)
	}

	t.Run("file with env overrides", func(t *testing.T) {
		// Set environment variables to override some file values
		os.Setenv("SERVER_HOST", "env-host")
		os.Setenv("SERVER_PORT", "4000")
		defer func() {
			os.Unsetenv("SERVER_HOST")
			os.Unsetenv("SERVER_PORT")
		}()

		// Use a struct that supports both file loading and env overrides
		type MixedConfig struct {
			Server struct {
				Host string `yaml:"host" json:"host" env:"SERVER_HOST"`
				Port int    `yaml:"port" json:"port" env:"SERVER_PORT"`
			} `yaml:"server" json:"server"`
			Database struct {
				URL      string `yaml:"url" json:"url" env:"DATABASE_URL"`
				MaxConns int    `yaml:"max_conns" json:"max_conns" env:"MAX_CONNS"`
			} `yaml:"database" json:"database"`
		}

		var cfg MixedConfig
		err := config.LoadFromFileWithEnv(&cfg, yamlFile)
		if err != nil {
			t.Fatalf("LoadFromFileWithEnv failed: %v", err)
		}

		// Values overridden by env vars
		if cfg.Server.Host != "env-host" {
			t.Errorf("Expected Server.Host to be 'env-host' (from env), got '%s'", cfg.Server.Host)
		}
		if cfg.Server.Port != 4000 {
			t.Errorf("Expected Server.Port to be 4000 (from env), got %d", cfg.Server.Port)
		}

		// Values from file (no env override)
		if cfg.Database.URL != "postgres://file-db/test" {
			t.Errorf("Expected Database.URL to be 'postgres://file-db/test' (from file), got '%s'", cfg.Database.URL)
		}
		if cfg.Database.MaxConns != 10 {
			t.Errorf("Expected Database.MaxConns to be 10 (from file), got %d", cfg.Database.MaxConns)
		}
	})
}

func TestMustFunctions(t *testing.T) {
	t.Run("MustLoadFromEnv panics on error", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected MustLoadFromEnv to panic")
			}
		}()

		var cfg BasicConfig
		config.MustLoadFromEnv(&cfg) // Should panic due to missing required field
	})

	t.Run("MustLoadFromFile panics on error", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected MustLoadFromFile to panic")
			}
		}()

		var cfg FileConfig
		config.MustLoadFromFile(&cfg, "non-existent.yaml") // Should panic due to missing file
	})

	t.Run("MustLoadFromEnv succeeds", func(t *testing.T) {
		os.Setenv("DATABASE_URL", "postgres://localhost/test")
		defer os.Unsetenv("DATABASE_URL")

		var cfg BasicConfig
		config.MustLoadFromEnv(&cfg) // Should not panic

		if cfg.Database != "postgres://localhost/test" {
			t.Errorf("Expected Database to be 'postgres://localhost/test', got '%s'", cfg.Database)
		}
	})
}

func TestErrorHandling(t *testing.T) {
	t.Run("invalid integer", func(t *testing.T) {
		os.Setenv("PORT", "invalid")
		os.Setenv("DATABASE_URL", "postgres://localhost/test")
		defer func() {
			os.Unsetenv("PORT")
			os.Unsetenv("DATABASE_URL")
		}()

		var cfg BasicConfig
		err := config.LoadFromEnv(&cfg)
		if err == nil {
			t.Fatal("Expected error for invalid integer")
		}
	})

	t.Run("invalid boolean", func(t *testing.T) {
		os.Setenv("DEBUG", "maybe")
		os.Setenv("DATABASE_URL", "postgres://localhost/test")
		defer func() {
			os.Unsetenv("DEBUG")
			os.Unsetenv("DATABASE_URL")
		}()

		var cfg BasicConfig
		err := config.LoadFromEnv(&cfg)
		if err == nil {
			t.Fatal("Expected error for invalid boolean")
		}
	})
}