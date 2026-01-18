package validator_test

import (
	"testing"

	"github.com/julianstephens/go-utils/validator"
)

func TestParseValidator(t *testing.T) {
	v := validator.Parse()

	// Test ValidateEmail
	if err := v.ValidateEmail("user@example.com"); err != nil {
		t.Errorf("ValidateEmail('user@example.com') should pass, got error: %v", err)
	}
	if err := v.ValidateEmail("invalid-email"); err == nil {
		t.Error("ValidateEmail('invalid-email') should fail")
	}
	if err := v.ValidateEmail(""); err == nil {
		t.Error("ValidateEmail('') should fail")
	}

	// Test ValidatePassword
	if err := v.ValidatePassword("ValidPass123"); err != nil {
		t.Errorf("ValidatePassword('ValidPass123') should pass, got error: %v", err)
	}
	if err := v.ValidatePassword("short"); err == nil {
		t.Error("ValidatePassword('short') should fail")
	}
	if err := v.ValidatePassword("nouppercase123"); err == nil {
		t.Error("ValidatePassword('nouppercase123') should fail")
	}
	if err := v.ValidatePassword("NOLOWERCASE123"); err == nil {
		t.Error("ValidatePassword('NOLOWERCASE123') should fail")
	}
	if err := v.ValidatePassword("NoDigitsHere"); err == nil {
		t.Error("ValidatePassword('NoDigitsHere') should fail")
	}

	// Test ValidateUUID
	if err := v.ValidateUUID("550e8400-e29b-41d4-a716-446655440000"); err != nil {
		t.Errorf("ValidateUUID valid UUID should pass, got error: %v", err)
	}
	if err := v.ValidateUUID("not-a-uuid"); err == nil {
		t.Error("ValidateUUID('not-a-uuid') should fail")
	}
	if err := v.ValidateUUID(""); err == nil {
		t.Error("ValidateUUID('') should fail")
	}

	// Test ValidateBool
	if err := v.ValidateBool("true"); err != nil {
		t.Errorf("ValidateBool('true') should pass, got error: %v", err)
	}
	if err := v.ValidateBool("false"); err != nil {
		t.Errorf("ValidateBool('false') should pass, got error: %v", err)
	}
	if err := v.ValidateBool("1"); err != nil {
		t.Errorf("ValidateBool('1') should pass, got error: %v", err)
	}
	if err := v.ValidateBool("invalid"); err == nil {
		t.Error("ValidateBool('invalid') should fail")
	}

	// Test ValidateInt
	if err := v.ValidateInt("123"); err != nil {
		t.Errorf("ValidateInt('123') should pass, got error: %v", err)
	}
	if err := v.ValidateInt("12.34"); err == nil {
		t.Error("ValidateInt('12.34') should fail")
	}
	if err := v.ValidateInt("abc"); err == nil {
		t.Error("ValidateInt('abc') should fail")
	}

	// Test ValidateFloat
	if err := v.ValidateFloat("123.456"); err != nil {
		t.Errorf("ValidateFloat('123.456') should pass, got error: %v", err)
	}
	if err := v.ValidateFloat("abc"); err == nil {
		t.Error("ValidateFloat('abc') should fail")
	}

	// Test ValidateUint
	if err := v.ValidateUint("123"); err != nil {
		t.Errorf("ValidateUint('123') should pass, got error: %v", err)
	}
	if err := v.ValidateUint("-1"); err == nil {
		t.Error("ValidateUint('-1') should fail")
	}
	if err := v.ValidateUint("12.34"); err == nil {
		t.Error("ValidateUint('12.34') should fail")
	}

	// Test ValidatePositiveInt
	if err := v.ValidatePositiveInt("5"); err != nil {
		t.Errorf("ValidatePositiveInt('5') should pass, got error: %v", err)
	}
	if err := v.ValidatePositiveInt("0"); err == nil {
		t.Error("ValidatePositiveInt('0') should fail")
	}
	if err := v.ValidatePositiveInt("-1"); err == nil {
		t.Error("ValidatePositiveInt('-1') should fail")
	}

	// Test ValidateNonNegativeInt
	if err := v.ValidateNonNegativeInt("5"); err != nil {
		t.Errorf("ValidateNonNegativeInt('5') should pass, got error: %v", err)
	}
	if err := v.ValidateNonNegativeInt("0"); err != nil {
		t.Errorf("ValidateNonNegativeInt('0') should pass, got error: %v", err)
	}
	if err := v.ValidateNonNegativeInt("-1"); err == nil {
		t.Error("ValidateNonNegativeInt('-1') should fail")
	}

	// Test ValidatePositiveFloat
	if err := v.ValidatePositiveFloat("3.14"); err != nil {
		t.Errorf("ValidatePositiveFloat('3.14') should pass, got error: %v", err)
	}
	if err := v.ValidatePositiveFloat("0.0"); err == nil {
		t.Error("ValidatePositiveFloat('0.0') should fail")
	}
	if err := v.ValidatePositiveFloat("-1.5"); err == nil {
		t.Error("ValidatePositiveFloat('-1.5') should fail")
	}

	// Test ValidateURL
	if err := v.ValidateURL("http://example.com"); err != nil {
		t.Errorf("ValidateURL('http://example.com') should pass, got error: %v", err)
	}
	if err := v.ValidateURL("not a url"); err == nil {
		t.Error("ValidateURL('not a url') should fail")
	}

	// Test ValidateIPAddress
	if err := v.ValidateIPAddress("192.168.1.1"); err != nil {
		t.Errorf("ValidateIPAddress('192.168.1.1') should pass, got error: %v", err)
	}
	if err := v.ValidateIPAddress("2001:db8::1"); err != nil {
		t.Errorf("ValidateIPAddress('2001:db8::1') should pass, got error: %v", err)
	}
	if err := v.ValidateIPAddress("invalid.ip"); err == nil {
		t.Error("ValidateIPAddress('invalid.ip') should fail")
	}

	// Test ValidateIPv4
	if err := v.ValidateIPv4("192.168.1.1"); err != nil {
		t.Errorf("ValidateIPv4('192.168.1.1') should pass, got error: %v", err)
	}
	if err := v.ValidateIPv4("2001:db8::1"); err == nil {
		t.Error("ValidateIPv4('2001:db8::1') should fail")
	}

	// Test ValidateIPv6
	if err := v.ValidateIPv6("2001:db8::1"); err != nil {
		t.Errorf("ValidateIPv6('2001:db8::1') should pass, got error: %v", err)
	}
	if err := v.ValidateIPv6("192.168.1.1"); err == nil {
		t.Error("ValidateIPv6('192.168.1.1') should fail")
	}
}
