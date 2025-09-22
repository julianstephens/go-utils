package response_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julianstephens/go-utils/httputil/response"
	testhelpers "github.com/julianstephens/go-utils/tests"
)

func TestResponder_NoContent(t *testing.T) {
	responder := response.New()
	req, w := testhelpers.NewRequestAndRecorder("GET", "/test")
	responder.NoContent(w, req)
	testhelpers.AssertStatus(t, w, 204)
	testhelpers.AssertBodyEquals(t, w, "")
}

func TestResponder_Created(t *testing.T) {
	responder := response.New()
	req, w := testhelpers.NewRequestAndRecorder("POST", "/test")
	data := map[string]string{"created": "yes"}
	responder.Created(w, req, data)
	testhelpers.AssertStatus(t, w, 201)
	testhelpers.AssertBodyContains(t, w, "created")
}

func TestResponder_Unauthorized(t *testing.T) {
	responder := response.New()
	req, w := testhelpers.NewRequestAndRecorder("GET", "/test")
	responder.Unauthorized(w, req, "unauthorized access")
	testhelpers.AssertStatus(t, w, 401)
	testhelpers.AssertBodyContains(t, w, "unauthorized access")
}

func TestResponder_Forbidden(t *testing.T) {
	responder := response.New()
	req, w := testhelpers.NewRequestAndRecorder("GET", "/test")
	responder.Forbidden(w, req, "forbidden access")
	testhelpers.AssertStatus(t, w, 403)
	testhelpers.AssertBodyContains(t, w, "forbidden access")
}

func TestResponder_NotFound(t *testing.T) {
	responder := response.New()
	req, w := testhelpers.NewRequestAndRecorder("GET", "/test")
	responder.NotFound(w, req, "not found")
	testhelpers.AssertStatus(t, w, 404)
	testhelpers.AssertBodyContains(t, w, "not found")
}

func TestResponder_InternalServerError(t *testing.T) {
	responder := response.New()
	req, w := testhelpers.NewRequestAndRecorder("GET", "/test")
	responder.InternalServerError(w, req, "internal error")
	testhelpers.AssertStatus(t, w, 500)
	testhelpers.AssertBodyContains(t, w, "internal error")
}

func TestParseErrData(t *testing.T) {
	err := response.ParseErrData(errors.New("err1"))
	testhelpers.AssertDeepEqual(t, err.Error(), "err1")
	err = response.ParseErrData("err2")
	testhelpers.AssertDeepEqual(t, err.Error(), "err2")
	err = response.ParseErrData(nil)
	testhelpers.AssertDeepEqual(t, err, nil)
}

func TestDefaultHooks(t *testing.T) {
	req, w := testhelpers.NewRequestAndRecorder("GET", "/")
	response.DefaultBefore(w, req, nil)
	response.DefaultAfter(w, req, nil)
	response.DefaultOnError(w, req, errors.New("fail"), 0)
	testhelpers.AssertStatus(t, w, 500)
	testhelpers.AssertBodyContains(t, w, "fail")
}

func TestLoggingHooks(t *testing.T) {
	req, w := testhelpers.NewRequestAndRecorder("GET", "/")
	response.LoggingBefore(w, req, nil)
	response.LoggingAfter(w, req, nil)
	response.LoggingOnError(w, req, errors.New("fail2"), 400)
	testhelpers.AssertStatus(t, w, 400)
	testhelpers.AssertBodyContains(t, w, "fail2")
}

func TestResponder_ErrorWithStatus(t *testing.T) {
	responder := response.New()
	req, w := testhelpers.NewRequestAndRecorder("GET", "/test")
	err := errors.New("test error")
	responder.ErrorWithStatus(w, req, 500, err)
	testhelpers.AssertStatus(t, w, 500)
}

func TestJSONEncoder(t *testing.T) {
	encoder := response.NewJSONEncoder()
	data := map[string]string{"key": "value"}
	_, w := testhelpers.NewRequestAndRecorder("GET", "/test")
	err := encoder.Encode(w, data, 200)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
	}
	testhelpers.AssertBodyContains(t, w, "key")
	testhelpers.AssertBodyContains(t, w, "value")
}

func TestJSONEncoderWithIndent(t *testing.T) {
	encoder := response.NewJSONEncoderWithIndent("  ")
	data := map[string]string{"key": "value"}
	_, w := testhelpers.NewRequestAndRecorder("GET", "/test")
	err := encoder.Encode(w, data, 200)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	testhelpers.AssertBodyContains(t, w, "  ")
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
		OnError: func(w http.ResponseWriter, r *http.Request, err error, status int) {
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
	req, w := testhelpers.NewRequestAndRecorder("GET", "/test")
	responder.OK(w, req, map[string]string{"status": "success"})
	testhelpers.AssertStatus(t, w, 200)
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
	}
	testhelpers.AssertBodyEquals(t, w, "{\"status\":\"success\"}\n")
}

func TestErrorWithStatus(t *testing.T) {
	tests := []struct {
		name   string
		status int
		err    error
	}{
		{name: "Basic 500 Internal Server Error", status: 500, err: errors.New("internal server error")},
		{name: "Nil Error", status: 500, err: nil},
		{name: "400 Bad Request with String Error", status: 400, err: errors.New("bad request error")},
		{name: "403 Forbidden with Custom Error", status: 403, err: errors.New("forbidden access")},
		{name: "404 Not Found with Nil Error", status: 404, err: nil},
		{name: "500 Internal Server Error with Wrapped Error", status: 500, err: errors.New("wrapped internal error")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := response.New()
			req, w := testhelpers.NewRequestAndRecorder("GET", "/test")
			r.ErrorWithStatus(w, req, tt.status, tt.err)
			testhelpers.AssertStatus(t, w, tt.status)
		})
	}
	// Special cases for nil ResponseWriter and/or Request
	t.Run("Nil ResponseWriter", func(t *testing.T) {
		r := response.New()
		req, _ := testhelpers.NewRequestAndRecorder("GET", "/test")
		defer func() { _ = recover() }()
		r.ErrorWithStatus(nil, req, 500, errors.New("internal server error"))
	})
	t.Run("Nil ResponseWriter and Request", func(t *testing.T) {
		r := response.New()
		defer func() { _ = recover() }()
		r.ErrorWithStatus(nil, nil, 500, errors.New("internal server error"))
	})

	// Additional scenario: ErrorWithStatus with string data
	t.Run("String data as error", func(t *testing.T) {
		r := response.New()
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		// Use BadRequest to trigger parseErrData with string
		r.BadRequest(w, req, "string error message")
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
		body := w.Body.String()
		if !strings.Contains(body, "string error message") {
			t.Errorf("Expected response body to contain 'string error message', got %s", body)
		}
	})
}
