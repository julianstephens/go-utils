package tests

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"
)

// AssertDeepEqual asserts that two values are deeply equal.
func AssertDeepEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Expected %+v, got %+v", want, got)
	}
}

// AssertNoError fails the test if err is non-nil.
func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

// RequireNoError fails the test immediately if err is non-nil.
func RequireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

// AssertTrue asserts that cond is true.
func AssertTrue(t *testing.T, cond bool, msg string) {
	t.Helper()
	if !cond {
		t.Errorf("Assertion failed: %s", msg)
	}
}

// AssertFalse asserts that cond is false.
func AssertFalse(t *testing.T, cond bool, msg string) {
	t.Helper()
	if cond {
		t.Errorf("Assertion failed: %s", msg)
	}
}

// AssertNotNil asserts that v is not nil.
func AssertNotNil(t *testing.T, v interface{}, msg string) {
	t.Helper()
	if v == nil || (reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil()) {
		t.Errorf("Expected not nil: %s", msg)
	}
}

// AssertNil asserts that v is nil.
func AssertNil(t *testing.T, v interface{}, msg string) {
	t.Helper()
	if v == nil {
		return
	}
	rv := reflect.ValueOf(v)
	// Only some kinds support IsNil
	switch rv.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.Interface:
		if rv.IsNil() {
			return
		}
	}
	t.Errorf("Expected nil: %s", msg)
}

// AssertJSONEquals unmarshals gotJSON and compares deeply with want.
func AssertJSONEquals(t *testing.T, gotJSON string, want interface{}) {
	t.Helper()
	var got interface{}
	if err := json.Unmarshal([]byte(gotJSON), &got); err != nil {
		t.Fatalf("Invalid JSON: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("JSON mismatch. Expected %+v, got %+v", want, got)
	}
}

// AssertErrorContains checks that err is non-nil and its message contains want.
func AssertErrorContains(t *testing.T, err error, want string) {
	t.Helper()
	if err == nil {
		t.Fatalf("Expected error containing %q, but error was nil", want)
	}
	if !strings.Contains(err.Error(), want) {
		t.Errorf("Expected error to contain %q, got %q", want, err.Error())
	}
}

// AssertWithinDuration asserts two times are within delta of each other.
func AssertWithinDuration(t *testing.T, got, want time.Time, delta time.Duration) {
	t.Helper()
	diff := got.Sub(want)
	if diff < 0 {
		diff = -diff
	}
	if diff > delta {
		t.Errorf("Times differ by %v which is more than %v (got=%v, want=%v)", diff, delta, got, want)
	}
}

// AssertErrorIs asserts that err is non-nil and errors.Is(err, target) is true.
func AssertErrorIs(t *testing.T, err error, target error) {
	t.Helper()
	if err == nil {
		t.Fatalf("Expected error %v but got nil", target)
	}
	if !errors.Is(err, target) {
		t.Errorf("Expected error to be %v (via errors.Is), got %v", target, err)
	}
}

// AssertPanics asserts that f panics. Returns the recovered value.
func AssertPanics(t *testing.T, f func()) (recovered interface{}) {
	t.Helper()
	defer func() {
		recovered = recover()
		if recovered == nil {
			t.Errorf("expected panic but function completed normally")
		}
	}()
	f()
	return
}

// RequireDeepEqual fails immediately if got and want are not deeply equal.
func RequireDeepEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected %+v, got %+v", want, got)
	}
}

// AssertCloseTo asserts that two floats are within tolerance.
func AssertCloseTo(t *testing.T, got, want, tol float64) {
	t.Helper()
	diff := got - want
	if diff < 0 {
		diff = -diff
	}
	if diff > tol {
		t.Errorf("Values differ by %v which is more than %v (got=%v, want=%v)", diff, tol, got, want)
	}
}

// AssertGreaterThan asserts that got is greater than want.
func AssertGreaterThan[T interface{ ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | ~string }](t *testing.T, got, want T) {
	t.Helper()
	if !(got > want) {
		t.Errorf("Expected %v to be greater than %v", got, want)
	}
}

// AssertLessThan asserts that got is less than want.
func AssertLessThan[T interface{ ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | ~string }](t *testing.T, got, want T) {
	t.Helper()
	if !(got < want) {
		t.Errorf("Expected %v to be less than %v", got, want)
	}
}

// AssertGreaterThanOrEqual asserts that got is greater than or equal to want.
func AssertGreaterThanOrEqual[T interface{ ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | ~string }](t *testing.T, got, want T) {
	t.Helper()
	if !(got >= want) {
		t.Errorf("Expected %v to be greater than or equal to %v", got, want)
	}
}

// AssertLessThanOrEqual asserts that got is less than or equal to want.
func AssertLessThanOrEqual[T interface{ ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | ~string }](t *testing.T, got, want T) {
	t.Helper()
	if !(got <= want) {
		t.Errorf("Expected %v to be less than or equal to %v", got, want)
	}
}

// AssertEqual asserts that got is equal to want using == comparison.
func AssertEqual[T interface{ ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | ~string | ~bool }](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("Expected %v to equal %v", got, want)
	}
}
