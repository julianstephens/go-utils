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

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"test@example.com", false},
		{"user@domain.org", false},
		{"invalid-email", true},
		{"", true},
		{"@domain.com", true},
		{"user@", true},
		{"user.domain", true},
		{"a@b", true},
	}

	for _, tc := range tests {
		err := ValidateEmail(tc.input)
		if tc.wantErr && err == nil {
			t.Fatalf("ValidateEmail(%q) expected error, got nil", tc.input)
		}
		if !tc.wantErr && err != nil {
			t.Fatalf("ValidateEmail(%q) unexpected error: %v", tc.input, err)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		pass    string
		wantErr bool
	}{
		{"Abcdef1g", false},     // valid: length 8, upper, lower, digit
		{"short1A", true},       // too short (7)
		{"alllowercase1", true}, // no uppercase
		{"ALLUPPERCASE1", true}, // no lowercase
		{"NoDigitsHere", true},  // no digit
		{"A1b2C3d4", false},     // valid
	}

	for _, tc := range tests {
		err := ValidatePassword(tc.pass)
		if tc.wantErr && err == nil {
			t.Fatalf("ValidatePassword(%q) expected error, got nil", tc.pass)
		}
		if !tc.wantErr && err != nil {
			t.Fatalf("ValidatePassword(%q) unexpected error: %v", tc.pass, err)
		}
	}
}

func TestValidateUUID(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"", true},
		{"not-a-uuid", true},
		{"550e8400-e29b-41d4-a716-446655440000", false},
		{"00000000-0000-0000-0000-000000000000", false},
	}

	for _, tc := range tests {
		err := ValidateUUID(tc.input)
		if tc.wantErr && err == nil {
			t.Fatalf("ValidateUUID(%q) expected error, got nil", tc.input)
		}
		if !tc.wantErr && err != nil {
			t.Fatalf("ValidateUUID(%q) unexpected error: %v", tc.input, err)
		}
	}
}
