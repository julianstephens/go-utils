package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/julianstephens/go-utils/config"
	tst "github.com/julianstephens/go-utils/tests"
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
		tst.AssertNoError(t, err)
		tst.AssertTrue(t, cfg.Port == 8080, "Port should be 8080")
		tst.AssertTrue(t, cfg.Host == "localhost", "Host should be localhost")
		tst.AssertTrue(t, cfg.Database == "postgres://localhost/test", "Database should match")
		tst.AssertFalse(t, cfg.Debug, "Debug should be false")
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
		tst.AssertNoError(t, err)
		tst.AssertTrue(t, cfg.Port == 3000, "Port should be 3000")
		tst.AssertTrue(t, cfg.Host == "0.0.0.0", "Host should be 0.0.0.0")
		tst.AssertTrue(t, cfg.Database == "postgres://prod/db", "Database should match")
		tst.AssertTrue(t, cfg.Debug, "Debug should be true")
	})

	t.Run("required field missing", func(t *testing.T) {
		cleanEnv()
		defer cleanEnv()

		var cfg BasicConfig
		err := config.LoadFromEnv(&cfg)
		tst.AssertNotNil(t, err, "expected error for missing required field")
		tst.AssertTrue(t, err.Error() == "required field 'Database' (env: DATABASE_URL) is missing or empty", "error message should match")
	})

	t.Run("invalid pointer", func(t *testing.T) {
		var cfg BasicConfig
		err := config.LoadFromEnv(cfg) // Not a pointer
		tst.AssertNotNil(t, err, "expected error for non-pointer argument")
		tst.AssertTrue(t, err.Error() == "config must be a pointer to a struct", "error message should match")
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
		tst.AssertNoError(t, err)
		// Check basic types
		tst.AssertTrue(t, cfg.StringField == "default_string", "StringField should have default")
		tst.AssertTrue(t, cfg.IntField == 42, "IntField should have default")
		tst.AssertTrue(t, cfg.BoolField, "BoolField should be true")
		tst.AssertTrue(t, cfg.FloatField == 3.14, "FloatField should be 3.14")

		// Check pointer types
		tst.AssertNil(t, cfg.StringPtr, "StringPtr should be nil")
		tst.AssertNotNil(t, cfg.IntPtr, "IntPtr should not be nil")
		tst.AssertTrue(t, *cfg.IntPtr == 100, "IntPtr value should be 100")

		// Check slice types
		tst.AssertDeepEqual(t, cfg.StringSlice, []string{"a", "b", "c"})
		tst.AssertDeepEqual(t, cfg.IntSlice, []int{1, 2, 3})

		tst.AssertTrue(t, cfg.RequiredField == "required_value", "RequiredField should match")
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
		tst.AssertNoError(t, err)
		tst.AssertTrue(t, cfg.StringField == "override_string", "StringField should be overridden")
		tst.AssertTrue(t, cfg.IntField == 999, "IntField should be overridden")
		tst.AssertFalse(t, cfg.BoolField, "BoolField should be false")
		tst.AssertTrue(t, cfg.FloatField == 2.71, "FloatField should be overridden")

		tst.AssertNotNil(t, cfg.StringPtr, "StringPtr should not be nil")
		tst.AssertTrue(t, *cfg.StringPtr == "pointer_value", "StringPtr value should match")
		tst.AssertNotNil(t, cfg.IntPtr, "IntPtr should not be nil")
		tst.AssertTrue(t, *cfg.IntPtr == 200, "IntPtr value should match")

		tst.AssertDeepEqual(t, cfg.StringSlice, []string{"x", "y", "z"})
		tst.AssertDeepEqual(t, cfg.IntSlice, []int{10, 20, 30})
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
		tst.AssertNoError(t, err)
		tst.AssertTrue(t, cfg.Server.Host == "127.0.0.1", "Server.Host should match")
		tst.AssertTrue(t, cfg.Server.Port == 9000, "Server.Port should match")
		tst.AssertTrue(t, cfg.Database.URL == "postgres://localhost/yaml_test", "Database.URL should match")
		tst.AssertTrue(t, cfg.Database.MaxConns == 50, "Database.MaxConns should match")
		tst.AssertTrue(t, cfg.Features.EnableMetrics, "EnableMetrics should be true")
		tst.AssertFalse(t, cfg.Features.EnableTracing, "EnableTracing should be false")
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
		tst.AssertNoError(t, err)
		tst.AssertTrue(t, cfg.Server.Host == "0.0.0.0", "Server.Host should match")
		tst.AssertTrue(t, cfg.Server.Port == 8080, "Server.Port should match")
		tst.AssertTrue(t, cfg.Database.URL == "postgres://localhost/json_test", "Database.URL should match")
		tst.AssertTrue(t, cfg.Database.MaxConns == 25, "Database.MaxConns should match")
		tst.AssertFalse(t, cfg.Features.EnableMetrics, "EnableMetrics should be false")
		tst.AssertTrue(t, cfg.Features.EnableTracing, "EnableTracing should be true")
	})

	t.Run("unsupported file format", func(t *testing.T) {
		unsupportedFile := filepath.Join(tempDir, "config.xml")
		if err := os.WriteFile(unsupportedFile, []byte("<config></config>"), 0644); err != nil {
			t.Fatalf("Failed to create unsupported test file: %v", err)
		}

		var cfg FileConfig
		err := config.LoadFromFile(&cfg, unsupportedFile)
		tst.AssertNotNil(t, err, "expected error for unsupported file format")
		tst.AssertTrue(t, err.Error() == "unsupported config file format '.xml', supported formats: .yaml, .yml, .json", "error message should match")
	})

	t.Run("non-existent file", func(t *testing.T) {
		var cfg FileConfig
		err := config.LoadFromFile(&cfg, "non-existent.yaml")
		tst.AssertNotNil(t, err, "expected error for non-existent file")
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
		tst.AssertNoError(t, err)
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
		tst.AssertNoError(t, err)
		// Values overridden by env vars
		tst.AssertTrue(t, cfg.Server.Host == "env-host", "Server.Host should be overridden by env")
		tst.AssertTrue(t, cfg.Server.Port == 4000, "Server.Port should be overridden by env")

		// Values from file (no env override)
		tst.AssertTrue(t, cfg.Database.URL == "postgres://file-db/test", "Database.URL should match file")
		tst.AssertTrue(t, cfg.Database.MaxConns == 10, "Database.MaxConns should match file")
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

		tst.AssertTrue(t, cfg.Database == "postgres://localhost/test", "Database should match")
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
		tst.AssertNotNil(t, err, "expected error for invalid integer")
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
		tst.AssertNotNil(t, err, "expected error for invalid boolean")
	})
}
