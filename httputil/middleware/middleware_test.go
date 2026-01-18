package middleware_test

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julianstephens/go-utils/httputil/middleware"
	tst "github.com/julianstephens/go-utils/tests"
)

func TestLogging(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	handler := middleware.Logging(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("test response"))
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	tst.AssertStatus(t, w, http.StatusOK)

	logOutput := buf.String()
	tst.AssertTrue(t, strings.Contains(logOutput, "GET"), "Log should contain HTTP method")
	tst.AssertTrue(t, strings.Contains(logOutput, "/test"), "Log should contain request path")
	tst.AssertTrue(t, strings.Contains(logOutput, "200"), "Log should contain status code")
}

func TestLoggingWithNilLogger(t *testing.T) {
	// Should use default logger without panicking
	handler := middleware.Logging(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	tst.AssertStatus(t, w, http.StatusOK)
}

func TestRecovery(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	handler := middleware.Recovery(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	tst.AssertStatus(t, w, http.StatusInternalServerError)

	tst.AssertBodyContains(t, w, "Internal Server Error")
	logOutput := buf.String()
	tst.AssertTrue(t, strings.Contains(logOutput, "Panic"), "Log should contain panic message")
	tst.AssertTrue(t, strings.Contains(logOutput, "test panic"), "Log should contain panic details")
}

func TestRecoveryWithNilLogger(t *testing.T) {
	handler := middleware.Recovery(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	tst.AssertStatus(t, w, http.StatusInternalServerError)
}

func TestRecoveryNoPanic(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	handler := middleware.Recovery(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("no panic"))
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	tst.AssertStatus(t, w, http.StatusOK)

	tst.AssertBodyEquals(t, w, "no panic")
	logOutput := buf.String()
	tst.AssertFalse(
		t,
		strings.Contains(logOutput, "Panic"),
		"Log should not contain panic message when no panic occurs",
	)
}

func TestCORS(t *testing.T) {
	config := middleware.CORSConfig{
		AllowedOrigins:   []string{"https://example.com", "https://test.com"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"X-Total-Count"},
		AllowCredentials: true,
		MaxAge:           3600,
	}

	handler := middleware.CORS(config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("cors test"))
	}))

	// Test with allowed origin
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	tst.AssertTrue(
		t,
		w.Header().Get("Access-Control-Allow-Origin") == "https://example.com",
		"Should set Access-Control-Allow-Origin for allowed origin",
	)
	tst.AssertTrue(
		t,
		w.Header().Get("Access-Control-Allow-Methods") == "GET, POST",
		"Should set Access-Control-Allow-Methods",
	)
	tst.AssertTrue(
		t,
		w.Header().Get("Access-Control-Allow-Headers") == "Content-Type, Authorization",
		"Should set Access-Control-Allow-Headers",
	)
	tst.AssertTrue(
		t,
		w.Header().Get("Access-Control-Expose-Headers") == "X-Total-Count",
		"Should set Access-Control-Expose-Headers",
	)
	tst.AssertTrue(
		t,
		w.Header().Get("Access-Control-Allow-Credentials") == "true",
		"Should set Access-Control-Allow-Credentials",
	)
	tst.AssertTrue(t, w.Header().Get("Access-Control-Max-Age") == "3600", "Should set Access-Control-Max-Age")
}

func TestCORSPreflightRequest(t *testing.T) {
	config := middleware.DefaultCORSConfig()
	handler := middleware.CORS(config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("should not reach here"))
	}))

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	tst.AssertStatus(t, w, http.StatusNoContent)
	tst.AssertBodyEquals(t, w, "")
}

func TestCORSWildcardOrigin(t *testing.T) {
	config := middleware.DefaultCORSConfig() // Uses "*" by default
	handler := middleware.CORS(config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://any-origin.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	tst.AssertTrue(
		t,
		w.Header().Get("Access-Control-Allow-Origin") == "*",
		"Should set Access-Control-Allow-Origin to * for wildcard config",
	)
}

func TestCORSDisallowedOrigin(t *testing.T) {
	config := middleware.CORSConfig{
		AllowedOrigins: []string{"https://allowed.com"},
	}

	handler := middleware.CORS(config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://disallowed.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	tst.AssertTrue(
		t,
		w.Header().Get("Access-Control-Allow-Origin") == "",
		"Should not set Access-Control-Allow-Origin for disallowed origin",
	)
}

func TestRequestID(t *testing.T) {
	var capturedRequestID string
	handler := middleware.RequestID()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedRequestID = middleware.GetRequestID(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Check that request ID was generated and set in context
	tst.AssertTrue(t, capturedRequestID != "", "Request ID should be set in context")
	// Check that request ID was set in response header
	responseRequestID := w.Header().Get(middleware.RequestIDHeader)
	tst.AssertTrue(t, responseRequestID != "", "Request ID should be set in response header")
	// Should be the same ID
	tst.AssertTrue(t, capturedRequestID == responseRequestID, "Request ID in context and header should match")
}

func TestRequestIDFromHeader(t *testing.T) {
	existingRequestID := "existing-request-id-123"
	var capturedRequestID string

	handler := middleware.RequestID()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedRequestID = middleware.GetRequestID(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set(middleware.RequestIDHeader, existingRequestID)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Should use the existing request ID
	tst.AssertTrue(t, capturedRequestID == existingRequestID, "Should use existing request ID from header")
	responseRequestID := w.Header().Get(middleware.RequestIDHeader)
	tst.AssertTrue(t, responseRequestID == existingRequestID, "Should return existing request ID in header")
}

func TestGetRequestIDFromEmptyContext(t *testing.T) {
	ctx := context.Background()
	requestID := middleware.GetRequestID(ctx)
	tst.AssertTrue(t, requestID == "", "Should return empty string for context without request ID")
}

func TestDefaultCORSConfig(t *testing.T) {
	config := middleware.DefaultCORSConfig()

	tst.AssertTrue(
		t,
		len(config.AllowedOrigins) == 1 && config.AllowedOrigins[0] == "*",
		"Default config should allow all origins",
	)
	expectedMethods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	tst.AssertDeepEqual(t, len(config.AllowedMethods), len(expectedMethods))
	tst.AssertTrue(t, config.MaxAge == 86400, "Default config should have 24 hour max age")
	tst.AssertFalse(t, config.AllowCredentials, "Default config should not allow credentials")
}

func TestCORSWildcardSubdomain(t *testing.T) {
	config := middleware.CORSConfig{
		AllowedOrigins: []string{"*.example.com"},
	}

	handler := middleware.CORS(config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Test subdomain that should be allowed
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://api.example.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	tst.AssertTrue(
		t,
		w.Header().Get("Access-Control-Allow-Origin") == "https://api.example.com",
		"Should allow subdomain matching wildcard pattern",
	)

	// Test domain that should not be allowed
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.Header.Set("Origin", "https://notexample.com")
	w2 := httptest.NewRecorder()

	handler.ServeHTTP(w2, req2)

	tst.AssertTrue(
		t,
		w2.Header().Get("Access-Control-Allow-Origin") == "",
		"Should not allow domain that doesn't match wildcard pattern",
	)
}
