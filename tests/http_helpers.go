package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// AssertStatus asserts that the response recorder has the expected status code.
func AssertStatus(t *testing.T, rr *httptest.ResponseRecorder, want int) {
	t.Helper()
	if rr.Code != want {
		t.Errorf("Expected status %d, got %d", want, rr.Code)
	}
}

// AssertBodyContains asserts that the response body contains the expected substring.
func AssertBodyContains(t *testing.T, rr *httptest.ResponseRecorder, want string) {
	t.Helper()
	if !contains(rr.Body.String(), want) {
		t.Errorf("Expected response body to contain '%s', got '%s'", want, rr.Body.String())
	}
}

// AssertBodyEquals asserts that the response body equals the expected string.
func AssertBodyEquals(t *testing.T, rr *httptest.ResponseRecorder, want string) {
	t.Helper()
	if rr.Body.String() != want {
		t.Errorf("Expected response body to be '%s', got '%s'", want, rr.Body.String())
	}
}

// NewRequestAndRecorder returns a new HTTP request and response recorder.
func NewRequestAndRecorder(method, target string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, target, nil)
	rr := httptest.NewRecorder()
	return req, rr
}

// AssertResponseJSON reads the recorder body and compares its JSON to want.
func AssertResponseJSON(t *testing.T, rr *httptest.ResponseRecorder, want interface{}) {
	t.Helper()
	var got interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("Invalid JSON response: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Response JSON mismatch. Expected %+v, got %+v", want, got)
	}
}

// AssertHeaderEquals asserts that a response header equals the expected value.
func AssertHeaderEquals(t *testing.T, rr *httptest.ResponseRecorder, key, want string) {
	t.Helper()
	got := rr.Header().Get(key)
	if got != want {
		t.Errorf("Expected header %s=%q, got %q", key, want, got)
	}
}
