package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadFromEnv loads configuration from environment variables into the provided struct.
// The struct fields should be tagged with `env:"ENV_VAR_NAME"` to specify the environment variable name.
// Optional tags: `default:"value"` for default values and `required:"true"` for required fields.
//
// Supported types: string, int (all variants), uint (all variants), float32, float64, bool,
// slices of these types, and pointers to these types.
func LoadFromEnv(cfg interface{}) error {
	return loadFromEnv(cfg)
}

// LoadFromFile loads configuration from a YAML or JSON file into the provided struct.
// The file format is determined by the file extension (.yaml, .yml, or .json).
// The struct should use standard json/yaml tags for field mapping.
func LoadFromFile(cfg interface{}, filepath string) error {
	return loadFromFile(cfg, filepath)
}

// LoadFromFileWithEnv loads configuration from a file and then overrides with environment variables.
// This allows for a hierarchical configuration approach where files provide defaults
// and environment variables provide runtime overrides.
func LoadFromFileWithEnv(cfg interface{}, filepath string) error {
	// First load from file
	if err := loadFromFile(cfg, filepath); err != nil {
		return fmt.Errorf("failed to load from file: %w", err)
	}

	// Then override with environment variables
	if err := loadFromEnv(cfg); err != nil {
		return fmt.Errorf("failed to override with environment variables: %w", err)
	}

	return nil
}

// MustLoadFromEnv is like LoadFromEnv but panics on error.
// Useful for application initialization where configuration errors should be fatal.
func MustLoadFromEnv(cfg interface{}) {
	if err := LoadFromEnv(cfg); err != nil {
		panic(fmt.Sprintf("failed to load config from environment: %v", err))
	}
}

// MustLoadFromFile is like LoadFromFile but panics on error.
// Useful for application initialization where configuration errors should be fatal.
func MustLoadFromFile(cfg interface{}, filepath string) {
	if err := LoadFromFile(cfg, filepath); err != nil {
		panic(fmt.Sprintf("failed to load config from file '%s': %v", filepath, err))
	}
}

// MustLoadFromFileWithEnv is like LoadFromFileWithEnv but panics on error.
// Useful for application initialization where configuration errors should be fatal.
func MustLoadFromFileWithEnv(cfg interface{}, filepath string) {
	if err := LoadFromFileWithEnv(cfg, filepath); err != nil {
		panic(fmt.Sprintf("failed to load config from file '%s' with env overrides: %v", filepath, err))
	}
}

// loadFromEnv is the internal implementation for loading from environment variables
func loadFromEnv(cfg interface{}) error {
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("config must be a pointer to a struct")
	}

	return processStruct(v.Elem())
}

// processStruct recursively processes struct fields for environment variable loading
func processStruct(v reflect.Value) error {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Skip unexported fields
		if !field.CanSet() {
			continue
		}

		// Handle nested structs
		if field.Kind() == reflect.Struct {
			if err := processStruct(field); err != nil {
				return err
			}
			continue
		}

		envTag := fieldType.Tag.Get("env")
		if envTag == "" {
			continue
		}

		defaultVal := fieldType.Tag.Get("default")
		required := fieldType.Tag.Get("required") == "true"

		envVal := os.Getenv(envTag)

		// Handle required fields
		if required && envVal == "" {
			return fmt.Errorf("required field '%s' (env: %s) is missing or empty", fieldType.Name, envTag)
		}

		// Use default value if env var is not set
		if envVal == "" && defaultVal != "" {
			envVal = defaultVal
		}

		// Skip if no value available
		if envVal == "" {
			continue
		}

		// Set the field value
		if err := setFieldValue(field, envVal, fieldType.Name); err != nil {
			return fmt.Errorf("failed to set field '%s' (env: %s): %w", fieldType.Name, envTag, err)
		}
	}

	return nil
}

// loadFromFile is the internal implementation for loading from files
func loadFromFile(cfg interface{}, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read config file '%s': %w", filePath, err)
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return fmt.Errorf("failed to parse YAML config file '%s': %w", filePath, err)
		}
	case ".json":
		if err := json.Unmarshal(data, cfg); err != nil {
			return fmt.Errorf("failed to parse JSON config file '%s': %w", filePath, err)
		}
	default:
		return fmt.Errorf("unsupported config file format '%s', supported formats: .yaml, .yml, .json", ext)
	}

	return nil
}

// setFieldValue sets a struct field value from a string representation
func setFieldValue(field reflect.Value, value string, fieldName string) error {
	// Handle pointers
	if field.Kind() == reflect.Ptr {
		// Create new instance if nil
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		field = field.Elem()
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)

	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid boolean value '%s': %w", value, err)
		}
		field.SetBool(b)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 10, int(field.Type().Size()*8))
		if err != nil {
			return fmt.Errorf("invalid integer value '%s': %w", value, err)
		}
		field.SetInt(i)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(value, 10, int(field.Type().Size()*8))
		if err != nil {
			return fmt.Errorf("invalid unsigned integer value '%s': %w", value, err)
		}
		field.SetUint(u)

	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(value, int(field.Type().Size()*8))
		if err != nil {
			return fmt.Errorf("invalid float value '%s': %w", value, err)
		}
		field.SetFloat(f)

	case reflect.Slice:
		return setSliceValue(field, value, fieldName)

	default:
		return fmt.Errorf("unsupported field type %s for field %s", field.Kind(), fieldName)
	}

	return nil
}

// setSliceValue sets a slice field value from a comma-separated string
func setSliceValue(field reflect.Value, value string, fieldName string) error {
	if value == "" {
		return nil
	}

	// Split by comma and trim spaces
	parts := strings.Split(value, ",")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}

	sliceType := field.Type()

	// Create new slice
	newSlice := reflect.MakeSlice(sliceType, len(parts), len(parts))

	// Set each element
	for i, part := range parts {
		elem := newSlice.Index(i)
		if err := setFieldValue(elem, part, fmt.Sprintf("%s[%d]", fieldName, i)); err != nil {
			return err
		}
	}

	field.Set(newSlice)
	return nil
}
