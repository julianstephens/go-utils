package tests

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"
)

// Ordered is a constraint for types that support comparison operators.
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | ~string
}

// Equatable is a constraint for types that support equality comparison.
type Equatable interface {
	Ordered | ~bool
}

// AssertDeepEqual asserts that two values are deeply equal.
func AssertDeepEqual(t *testing.T, got, want interface{}, msg ...string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		var errMsg string
		if len(msg) > 0 && msg[0] != "" {
			errMsg = ": " + msg[0]
		}
		t.Errorf("Expected %+v, got %+v%s", want, got, errMsg)
	}
}

// AssertNoError fails the test if err is non-nil.
func AssertNoError(t *testing.T, err error, msg ...string) {
	t.Helper()
	if err != nil {
		var errMsg string
		if len(msg) > 0 && msg[0] != "" {
			errMsg = ": " + msg[0]
		}
		t.Errorf("Unexpected error: %v%s", err, errMsg)
	}
}

// RequireNoError fails the test immediately if err is non-nil.
func RequireNoError(t *testing.T, err error, msg ...string) {
	t.Helper()
	if err != nil {
		var errMsg string
		if len(msg) > 0 && msg[0] != "" {
			errMsg = ": " + msg[0]
		}
		t.Fatalf("Unexpected error: %v%s", err, errMsg)
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
func AssertNotNil(t *testing.T, v interface{}, msg ...string) {
	t.Helper()
	if v == nil || (reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil()) {
		var errMsg string
		if len(msg) > 0 && msg[0] != "" {
			errMsg = msg[0]
		} else {
			errMsg = "value is nil"
		}
		t.Errorf("Expected not nil: %s", errMsg)
	}
}

// AssertNil asserts that v is nil.
func AssertNil(t *testing.T, v interface{}, msg ...string) {
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
	var errMsg string
	if len(msg) > 0 && msg[0] != "" {
		errMsg = msg[0]
	} else {
		errMsg = "value is not nil"
	}
	t.Errorf("Expected nil: %s", errMsg)
}

// AssertJSONEquals unmarshals gotJSON and compares deeply with want.
func AssertJSONEquals(t *testing.T, gotJSON string, want interface{}, msg ...string) {
	t.Helper()
	var got interface{}
	if err := json.Unmarshal([]byte(gotJSON), &got); err != nil {
		t.Fatalf("Invalid JSON: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		var errMsg string
		if len(msg) > 0 && msg[0] != "" {
			errMsg = ": " + msg[0]
		}
		t.Errorf("JSON mismatch. Expected %+v, got %+v%s", want, got, errMsg)
	}
}

// AssertErrorContains checks that err is non-nil and its message contains want.
func AssertErrorContains(t *testing.T, err error, want string, msg ...string) {
	t.Helper()
	if err == nil {
		var errMsg string
		if len(msg) > 0 && msg[0] != "" {
			errMsg = " (" + msg[0] + ")"
		}
		t.Fatalf("Expected error containing %q, but error was nil%s", want, errMsg)
	}
	if !strings.Contains(err.Error(), want) {
		var errMsg string
		if len(msg) > 0 && msg[0] != "" {
			errMsg = " (" + msg[0] + ")"
		}
		t.Errorf("Expected error to contain %q, got %q%s", want, err.Error(), errMsg)
	}
}

// AssertWithinDuration asserts two times are within delta of each other.
func AssertWithinDuration(t *testing.T, got, want time.Time, delta time.Duration, msg ...string) {
	t.Helper()
	diff := got.Sub(want)
	if diff < 0 {
		diff = -diff
	}
	if diff > delta {
		var errMsg string
		if len(msg) > 0 && msg[0] != "" {
			errMsg = " (" + msg[0] + ")"
		}
		t.Errorf("Times differ by %v which is more than %v (got=%v, want=%v)%s", diff, delta, got, want, errMsg)
	}
}

// AssertErrorIs asserts that err is non-nil and errors.Is(err, target) is true.
func AssertErrorIs(t *testing.T, err error, target error, msg ...string) {
	t.Helper()
	if err == nil {
		var errMsg string
		if len(msg) > 0 && msg[0] != "" {
			errMsg = " (" + msg[0] + ")"
		}
		t.Fatalf("Expected error %v but got nil%s", target, errMsg)
	}
	if !errors.Is(err, target) {
		var errMsg string
		if len(msg) > 0 && msg[0] != "" {
			errMsg = " (" + msg[0] + ")"
		}
		t.Errorf("Expected error to be %v (via errors.Is), got %v%s", target, err, errMsg)
	}
}

// AssertPanics asserts that f panics. Returns the recovered value.
func AssertPanics(t *testing.T, f func(), msg ...string) (recovered interface{}) {
	t.Helper()
	defer func() {
		recovered = recover()
		if recovered == nil {
			var errMsg string
			if len(msg) > 0 && msg[0] != "" {
				errMsg = ": " + msg[0]
			}
			t.Errorf("expected panic but function completed normally%s", errMsg)
		}
	}()
	f()
	return
}

// RequireDeepEqual fails immediately if got and want are not deeply equal.
func RequireDeepEqual(t *testing.T, got, want interface{}, msg ...string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		var errMsg string
		if len(msg) > 0 && msg[0] != "" {
			errMsg = ": " + msg[0]
		}
		t.Fatalf("Expected %+v, got %+v%s", want, got, errMsg)
	}
}

// AssertCloseTo asserts that two floats are within tolerance.
func AssertCloseTo(t *testing.T, got, want, tol float64, msg ...string) {
	t.Helper()
	diff := got - want
	if diff < 0 {
		diff = -diff
	}
	if diff > tol {
		var errMsg string
		if len(msg) > 0 && msg[0] != "" {
			errMsg = " (" + msg[0] + ")"
		}
		t.Errorf("Values differ by %v which is more than %v (got=%v, want=%v)%s", diff, tol, got, want, errMsg)
	}
}

// AssertGreaterThan asserts that got is greater than want.
func AssertGreaterThan[T Ordered](t *testing.T, got, want T, msg ...string) {
	t.Helper()
	if got <= want {
		var errMsg string
		if len(msg) > 0 && msg[0] != "" {
			errMsg = " (" + msg[0] + ")"
		}
		t.Errorf("Expected %v to be greater than %v%s", got, want, errMsg)
	}
}

// AssertLessThan asserts that got is less than want.
func AssertLessThan[T Ordered](t *testing.T, got, want T, msg ...string) {
	t.Helper()
	if got >= want {
		var errMsg string
		if len(msg) > 0 && msg[0] != "" {
			errMsg = " (" + msg[0] + ")"
		}
		t.Errorf("Expected %v to be less than %v%s", got, want, errMsg)
	}
}

// AssertGreaterThanOrEqual asserts that got is greater than or equal to want.
func AssertGreaterThanOrEqual[T Ordered](t *testing.T, got, want T, msg ...string) {
	t.Helper()
	if got < want {
		var errMsg string
		if len(msg) > 0 && msg[0] != "" {
			errMsg = " (" + msg[0] + ")"
		}
		t.Errorf("Expected %v to be greater than or equal to %v%s", got, want, errMsg)
	}
}

// AssertLessThanOrEqual asserts that got is less than or equal to want.
func AssertLessThanOrEqual[T Ordered](t *testing.T, got, want T, msg ...string) {
	t.Helper()
	if got > want {
		var errMsg string
		if len(msg) > 0 && msg[0] != "" {
			errMsg = " (" + msg[0] + ")"
		}
		t.Errorf("Expected %v to be less than or equal to %v%s", got, want, errMsg)
	}
}

// AssertEqual asserts that got is equal to want using == comparison.
func AssertEqual[T Equatable](t *testing.T, got, want T, msg ...string) {
	t.Helper()
	if got != want {
		var errMsg string
		if len(msg) > 0 && msg[0] != "" {
			errMsg = " (" + msg[0] + ")"
		}
		t.Errorf("Expected %v to equal %v%s", got, want, errMsg)
	}
}
