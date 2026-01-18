package validator

import "testing"

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

	// Test New function
	validator := New()
	if validator == nil {
		t.Fatal("New() returned nil")
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
