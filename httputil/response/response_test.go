package response_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julianstephens/go-utils/httputil/response"
)

func TestResponder_WriteWithStatus(t *testing.T) {
	responder := response.New()

	// Test data
	data := map[string]string{"message": "test"}

	// Create test request and response recorder
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Call WriteWithStatus
	responder.WriteWithStatus(w, req, data, http.StatusOK)

	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	body := w.Body.String()
	if !strings.Contains(body, "test") {
		t.Errorf("Expected response body to contain 'test', got %s", body)
	}
}

func TestResponder_ErrorWithStatus(t *testing.T) {
	responder := response.New()

	// Create test request and response recorder
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Call ErrorWithStatus
	err := errors.New("test error")
	responder.ErrorWithStatus(w, req, http.StatusInternalServerError, err)

	// Verify error response
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestJSONEncoder(t *testing.T) {
	encoder := response.NewJSONEncoder()

	// Test data
	data := map[string]string{"key": "value"}

	// Create response recorder
	w := httptest.NewRecorder()

	// Encode
	err := encoder.Encode(w, data)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	// Verify JSON content
	body := w.Body.String()
	if !strings.Contains(body, "key") || !strings.Contains(body, "value") {
		t.Errorf("Expected JSON content, got %s", body)
	}
}

func TestJSONEncoderWithIndent(t *testing.T) {
	encoder := response.NewJSONEncoderWithIndent("  ")

	// Test data
	data := map[string]string{"key": "value"}

	// Create response recorder
	w := httptest.NewRecorder()

	// Encode
	err := encoder.Encode(w, data)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify indented JSON
	body := w.Body.String()
	if !strings.Contains(body, "  ") {
		t.Errorf("Expected indented JSON, got %s", body)
	}
}

func TestHooks(t *testing.T) {
	beforeCalled := false
	afterCalled := false
	errorCalled := false

	responder := &response.Responder{
		Encoder: response.NewJSONEncoder(),
		Before: func(w http.ResponseWriter, r *http.Request, data any) {
			beforeCalled = true
		},
		After: func(w http.ResponseWriter, r *http.Request, data any) {
			afterCalled = true
		},
		OnError: func(w http.ResponseWriter, r *http.Request, err error) {
			errorCalled = true
		},
	}

	// Test successful write
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	responder.WriteWithStatus(w, req, map[string]string{"test": "data"}, http.StatusOK)

	if !beforeCalled {
		t.Error("Before hook was not called")
	}
	if !afterCalled {
		t.Error("After hook was not called")
	}
	if errorCalled {
		t.Error("OnError hook should not have been called")
	}

	// Reset flags and test error
	beforeCalled = false
	afterCalled = false
	errorCalled = false

	w = httptest.NewRecorder()
	responder.ErrorWithStatus(w, req, http.StatusInternalServerError, errors.New("test error"))

	if !errorCalled {
		t.Error("OnError hook should have been called during ErrorWithStatus method")
	}
}

func TestNewConstructors(t *testing.T) {
	// Test New()
	r1 := response.New()
	if r1.Encoder == nil {
		t.Error("New() should provide a default encoder")
	}
	if r1.Before == nil || r1.After == nil || r1.OnError == nil {
		t.Error("New() should provide default hooks")
	}

	// Test NewWithLogging()
	r2 := response.NewWithLogging()
	if r2.Encoder == nil {
		t.Error("NewWithLogging() should provide a default encoder")
	}
	if r2.Before == nil || r2.After == nil || r2.OnError == nil {
		t.Error("NewWithLogging() should provide logging hooks")
	}

	// Test NewCustom()
	r3 := response.NewCustom(nil, nil, nil, nil)
	if r3.Encoder == nil {
		t.Error("NewCustom() should provide defaults when nil is passed")
	}
}

func TestWriteOK(t *testing.T) {
	responder := response.New()

	// Create test request and response recorder
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Call WriteOK
	responder.OK(w, req, map[string]string{"status": "success"})

	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	body := w.Body.String()
	if body != "{\"status\":\"success\"}\n" {
		t.Errorf("Expected response body to be '{\"status\":\"success\"}', got %s", body)
	}
}
