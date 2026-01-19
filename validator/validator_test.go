package validator

import (
	"fmt"
	"strings"
	"testing"
)

func TestValidateNonEmpty(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"valid input", false},
		{"", true},
		{"   ", true},
		{"\t\n", true},
		{"a", false},
	}

	for _, tc := range tests {
		err := ValidateNonEmpty(tc.input)
		if tc.wantErr && err == nil {
			t.Fatalf("ValidateNonEmpty(%q) expected error, got nil", tc.input)
		}
		if !tc.wantErr && err != nil {
			t.Fatalf("ValidateNonEmpty(%q) unexpected error: %v", tc.input, err)
		}
	}
}

func TestValidateNonEmpty_Bytes(t *testing.T) {
	tests := []struct {
		input   []byte
		wantErr bool
	}{
		{[]byte("valid"), false},
		{[]byte(""), true},
		{[]byte{0}, false},
	}

	for _, tc := range tests {
		err := ValidateNonEmpty(tc.input)
		if tc.wantErr && err == nil {
			t.Fatalf("ValidateNonEmpty(%v) expected error, got nil", tc.input)
		}
		if !tc.wantErr && err != nil {
			t.Fatalf("ValidateNonEmpty(%v) unexpected error: %v", tc.input, err)
		}
	}
}

func TestValidateNonEmpty_Runes(t *testing.T) {
	tests := []struct {
		input   []rune
		wantErr bool
	}{
		{[]rune("valid"), false},
		{[]rune(""), true},
		{[]rune{'a'}, false},
	}

	for _, tc := range tests {
		err := ValidateNonEmpty(tc.input)
		if tc.wantErr && err == nil {
			t.Fatalf("ValidateNonEmpty(%v) expected error, got nil", tc.input)
		}
		if !tc.wantErr && err != nil {
			t.Fatalf("ValidateNonEmpty(%v) unexpected error: %v", tc.input, err)
		}
	}
}

func TestValidateNonEmpty_Map(t *testing.T) {
	tests := []struct {
		input   map[string]interface{}
		wantErr bool
	}{
		{map[string]interface{}{"key": "value"}, false},
		{map[string]interface{}{}, true},
		{nil, true}, // maps are nil by default
	}

	for _, tc := range tests {
		err := ValidateNonEmpty(tc.input)
		if tc.wantErr && err == nil {
			t.Fatalf("ValidateNonEmpty(%v) expected error, got nil", tc.input)
		}
		if !tc.wantErr && err != nil {
			t.Fatalf("ValidateNonEmpty(%v) unexpected error: %v", tc.input, err)
		}
	}
}

func TestValidateNonEmpty_Slice(t *testing.T) {
	tests := []struct {
		input   []interface{}
		wantErr bool
	}{
		{[]interface{}{"item"}, false},
		{[]interface{}{}, true},
		{nil, true}, // slices are nil by default
	}

	for _, tc := range tests {
		err := ValidateNonEmpty(tc.input)
		if tc.wantErr && err == nil {
			t.Fatalf("ValidateNonEmpty(%v) expected error, got nil", tc.input)
		}
		if !tc.wantErr && err != nil {
			t.Fatalf("ValidateNonEmpty(%v) unexpected error: %v", tc.input, err)
		}
	}
}

func TestFactoryFunctions(t *testing.T) {
	// Test Numbers factory
	intVal := Numbers[int]()
	if intVal == nil {
		t.Fatal("Numbers[int]() returned nil")
	}

	floatVal := Numbers[float64]()
	if floatVal == nil {
		t.Fatal("Numbers[float64]() returned nil")
	}

	// Test Strings factory
	stringVal := Strings[string]()
	if stringVal == nil {
		t.Fatal("Strings[string]() returned nil")
	}

	bytesVal := Strings[[]byte]()
	if bytesVal == nil {
		t.Fatal("Strings[[]byte]() returned nil")
	}

	// Test Parse factory
	parseVal := Parse()
	if parseVal == nil {
		t.Fatal("Parse() returned nil")
	}
}

func TestStringValidator_ParseAccess(t *testing.T) {
	// Test that StringValidator provides access to Parse validator
	strVal := Strings[string]()
	if strVal == nil {
		t.Fatal("Strings[string]() returned nil")
	}

	if strVal.Parse == nil {
		t.Fatal("StringValidator.Parse should not be nil")
	}

	// Test that we can call parse methods through the string validator
	if err := strVal.Parse.ValidateEmail("user@example.com"); err != nil {
		t.Errorf("Parse access through StringValidator failed: %v", err)
	}

	if err := strVal.Parse.ValidateEmail("invalid-email"); err == nil {
		t.Error("Parse access through StringValidator should fail for invalid email")
	}
}

func TestValidateMatchesField_Strings(t *testing.T) {
	tests := []struct {
		value1  string
		value2  string
		field   string
		wantErr bool
	}{
		{"password123", "password123", "password", false},
		{"password123", "password456", "password", true},
		{"", "", "field", false},
		{"x", "x", "confirmation", false},
		{"abc", "xyz", "test", true},
	}

	for _, tc := range tests {
		err := ValidateMatchesField(tc.value1, tc.value2, tc.field)
		if tc.wantErr && err == nil {
			t.Fatalf("ValidateMatchesField(%q, %q, %q) expected error, got nil", tc.value1, tc.value2, tc.field)
		}
		if !tc.wantErr && err != nil {
			t.Fatalf("ValidateMatchesField(%q, %q, %q) unexpected error: %v", tc.value1, tc.value2, tc.field, err)
		}
	}
}

func TestValidateMatchesField_Integers(t *testing.T) {
	tests := []struct {
		value1  int
		value2  int
		field   string
		wantErr bool
	}{
		{42, 42, "id", false},
		{42, 43, "id", true},
		{0, 0, "count", false},
		{-1, -1, "offset", false},
		{100, 200, "value", true},
	}

	for _, tc := range tests {
		err := ValidateMatchesField(tc.value1, tc.value2, tc.field)
		if tc.wantErr && err == nil {
			t.Fatalf("ValidateMatchesField(%d, %d, %q) expected error, got nil", tc.value1, tc.value2, tc.field)
		}
		if !tc.wantErr && err != nil {
			t.Fatalf("ValidateMatchesField(%d, %d, %q) unexpected error: %v", tc.value1, tc.value2, tc.field, err)
		}
	}
}

func TestValidateMatchesField_Bools(t *testing.T) {
	tests := []struct {
		value1  bool
		value2  bool
		field   string
		wantErr bool
	}{
		{true, true, "flag", false},
		{false, false, "flag", false},
		{true, false, "flag", true},
		{false, true, "flag", true},
	}

	for _, tc := range tests {
		err := ValidateMatchesField(tc.value1, tc.value2, tc.field)
		if tc.wantErr && err == nil {
			t.Fatalf("ValidateMatchesField(%v, %v, %q) expected error, got nil", tc.value1, tc.value2, tc.field)
		}
		if !tc.wantErr && err != nil {
			t.Fatalf("ValidateMatchesField(%v, %v, %q) unexpected error: %v", tc.value1, tc.value2, tc.field, err)
		}
	}
}

func TestCustomValidator_Single(t *testing.T) {
	cv := NewCustomValidator()
	cv.Add(func() error {
		return nil
	})

	err := cv.Validate()
	if err != nil {
		t.Fatalf("CustomValidator.Validate() expected nil, got %v", err)
	}
}

func TestCustomValidator_Multiple(t *testing.T) {
	cv := NewCustomValidator()
	cv.Add(func() error {
		return nil
	}).Add(func() error {
		return nil
	}).Add(func() error {
		return nil
	})

	err := cv.Validate()
	if err != nil {
		t.Fatalf("CustomValidator.Validate() expected nil, got %v", err)
	}
}

func TestCustomValidator_FirstFails(t *testing.T) {
	cv := NewCustomValidator()
	cv.Add(func() error {
		return fmt.Errorf("first validator failed")
	}).Add(func() error {
		return nil
	})

	err := cv.Validate()
	if err == nil {
		t.Fatal("CustomValidator.Validate() expected error, got nil")
	}
	if err.Error() != "first validator failed" {
		t.Fatalf("CustomValidator.Validate() got %q, want %q", err.Error(), "first validator failed")
	}
}

func TestCustomValidator_SecondFails(t *testing.T) {
	cv := NewCustomValidator()
	cv.Add(func() error {
		return nil
	}).Add(func() error {
		return fmt.Errorf("second validator failed")
	}).Add(func() error {
		return nil
	})

	err := cv.Validate()
	if err == nil {
		t.Fatal("CustomValidator.Validate() expected error, got nil")
	}
	if err.Error() != "second validator failed" {
		t.Fatalf("CustomValidator.Validate() got %q, want %q", err.Error(), "second validator failed")
	}
}

func TestCustomValidator_StopsOnFirstError(t *testing.T) {
	callOrder := []int{}
	cv := NewCustomValidator()
	cv.Add(func() error {
		callOrder = append(callOrder, 1)
		return nil
	}).Add(func() error {
		callOrder = append(callOrder, 2)
		return fmt.Errorf("error at 2")
	}).Add(func() error {
		callOrder = append(callOrder, 3)
		return nil
	})

	err := cv.Validate()
	if err == nil {
		t.Fatal("CustomValidator.Validate() expected error, got nil")
	}

	// Verify that we stopped after validator 2, didn't call 3
	if len(callOrder) != 2 {
		t.Fatalf("Expected 2 validators called, got %d", len(callOrder))
	}
	if callOrder[0] != 1 || callOrder[1] != 2 {
		t.Fatalf("Expected call order [1, 2], got %v", callOrder)
	}
}

func TestCustomValidator_Empty(t *testing.T) {
	cv := NewCustomValidator()
	err := cv.Validate()
	if err != nil {
		t.Fatalf("CustomValidator.Validate() with no validators expected nil, got %v", err)
	}
}

func TestCustomValidator_Chaining(t *testing.T) {
	// Test that Add returns the receiver for chaining
	cv := NewCustomValidator()
	result := cv.Add(func() error { return nil })

	if result != cv {
		t.Fatal("CustomValidator.Add() should return receiver for chaining")
	}
}
func TestCustomValidator_SafeValidate_AllPass(t *testing.T) {
	cv := NewCustomValidator().
		Add(func() error { return nil }).
		Add(func() error { return nil })

	err := cv.SafeValidate()
	if err != nil {
		t.Fatalf("SafeValidate with all passing validators expected nil, got %v", err)
	}
}

func TestCustomValidator_SafeValidate_FirstFails(t *testing.T) {
	cv := NewCustomValidator().
		Add(func() error { return fmt.Errorf("first failed") }).
		Add(func() error { return nil })

	err := cv.SafeValidate()
	if err == nil {
		t.Fatal("SafeValidate with failing validator expected error, got nil")
	}
	if err.Error() != "first failed" {
		t.Fatalf("Expected 'first failed', got %q", err.Error())
	}
}

func TestCustomValidator_SafeValidate_PanicRecovery(t *testing.T) {
	cv := NewCustomValidator().
		Add(func() error {
			panic("test panic")
		}).
		Add(func() error { return nil })

	err := cv.SafeValidate()
	if err == nil {
		t.Fatal("SafeValidate with panic expected error, got nil")
	}
	// Error should indicate panic recovery occurred
	if !strings.Contains(err.Error(), "panic") {
		t.Fatalf("Expected panic in error message, got %q", err.Error())
	}
}

func TestCustomValidator_SafeValidate_StopsOnFirstError(t *testing.T) {
	callOrder := []int{}
	cv := NewCustomValidator().
		Add(func() error {
			callOrder = append(callOrder, 1)
			return nil
		}).
		Add(func() error {
			callOrder = append(callOrder, 2)
			return fmt.Errorf("error at 2")
		}).
		Add(func() error {
			callOrder = append(callOrder, 3)
			return nil
		})

	err := cv.SafeValidate()
	if err == nil {
		t.Fatal("SafeValidate with error expected error, got nil")
	}

	// Verify stopped at first error
	if len(callOrder) != 2 {
		t.Fatalf("Expected 2 validators called, got %d", len(callOrder))
	}
	if callOrder[0] != 1 || callOrder[1] != 2 {
		t.Fatalf("Expected call order [1, 2], got %v", callOrder)
	}
}
