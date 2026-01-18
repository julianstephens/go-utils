package helpers

import (
	"os"
	"reflect"
)

// If mimics the ternary operator s.t. cond ? vtrue : vfalse
func If[T any](cond bool, vtrue T, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

// Default returns defaultVal if val is the zero value for its type, otherwise returns val
func Default[T any](val T, defaultVal T) T {
	var zero T
	if reflect.DeepEqual(val, zero) {
		return defaultVal
	}
	return val
}

// Exists checks if the given file or directory path exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Ensure checks if the given path exists and creates it if not
func Ensure(path string, isDir bool) error {
	var f *os.File
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if isDir {
			if err = os.MkdirAll(path, os.ModePerm); err != nil {
				return err
			}
		} else {
			f, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
			if err != nil {
				return err
			}
			_ = f.Close()
		}
	}

	return nil
}

// StringPtr returns a pointer to the given string.
func StringPtr(s string) *string {
	return &s
}

// StructToMap converts a struct to a map[string]any using reflection.
func StructToMap(obj any) map[string]any {
	objValue := reflect.ValueOf(obj)
	if objValue.Kind() == reflect.Ptr {
		objValue = objValue.Elem()
	}

	if objValue.Kind() != reflect.Struct {
		return nil
	}

	result := make(map[string]any)
	objType := objValue.Type()

	for i := 0; i < objValue.NumField(); i++ {
		field := objType.Field(i)
		fieldValue := objValue.Field(i)
		result[field.Name] = fieldValue.Interface()
	}
	return result
}
