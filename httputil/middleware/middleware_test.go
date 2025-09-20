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
)

func TestLogging(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	handler := middleware.Logging(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	logOutput := buf.String()
	if !strings.Contains(logOutput, "GET") {
		t.Error("Log should contain HTTP method")
	}
	if !strings.Contains(logOutput, "/test") {
		t.Error("Log should contain request path")
	}
	if !strings.Contains(logOutput, "200") {
		t.Error("Log should contain status code")
	}
}

func TestLoggingWithNilLogger(t *testing.T) {
	// Should use default logger without panicking
	handler := middleware.Logging(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
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

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Internal Server Error") {
		t.Error("Response should contain 'Internal Server Error'")
	}

	logOutput := buf.String()
	if !strings.Contains(logOutput, "Panic") {
		t.Error("Log should contain panic message")
	}
	if !strings.Contains(logOutput, "test panic") {
		t.Error("Log should contain panic details")
	}
}

func TestRecoveryWithNilLogger(t *testing.T) {
	handler := middleware.Recovery(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestRecoveryNoPanic(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	handler := middleware.Recovery(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("no panic"))
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if body != "no panic" {
		t.Errorf("Expected body 'no panic', got %s", body)
	}

	logOutput := buf.String()
	if strings.Contains(logOutput, "Panic") {
		t.Error("Log should not contain panic message when no panic occurs")
	}
}

func TestCORS(t *testing.T) {
	config := middleware.CORSConfig{
		AllowedOrigins: []string{"https://example.com", "https://test.com"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		ExposedHeaders: []string{"X-Total-Count"},
		AllowCredentials: true,
		MaxAge:         3600,
	}

	handler := middleware.CORS(config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("cors test"))
	}))

	// Test with allowed origin
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "https://example.com" {
		t.Error("Should set Access-Control-Allow-Origin for allowed origin")
	}
	if w.Header().Get("Access-Control-Allow-Methods") != "GET, POST" {
		t.Error("Should set Access-Control-Allow-Methods")
	}
	if w.Header().Get("Access-Control-Allow-Headers") != "Content-Type, Authorization" {
		t.Error("Should set Access-Control-Allow-Headers")
	}
	if w.Header().Get("Access-Control-Expose-Headers") != "X-Total-Count" {
		t.Error("Should set Access-Control-Expose-Headers")
	}
	if w.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Error("Should set Access-Control-Allow-Credentials")
	}
	if w.Header().Get("Access-Control-Max-Age") != "3600" {
		t.Error("Should set Access-Control-Max-Age")
	}
}

func TestCORSPreflightRequest(t *testing.T) {
	config := middleware.DefaultCORSConfig()
	handler := middleware.CORS(config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("should not reach here"))
	}))

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d for OPTIONS request, got %d", http.StatusNoContent, w.Code)
	}

	body := w.Body.String()
	if body != "" {
		t.Error("OPTIONS request should have empty body")
	}
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

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("Should set Access-Control-Allow-Origin to * for wildcard config")
	}
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

	if w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Error("Should not set Access-Control-Allow-Origin for disallowed origin")
	}
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
	if capturedRequestID == "" {
		t.Error("Request ID should be set in context")
	}

	// Check that request ID was set in response header
	responseRequestID := w.Header().Get(middleware.RequestIDHeader)
	if responseRequestID == "" {
		t.Error("Request ID should be set in response header")
	}

	// Should be the same ID
	if capturedRequestID != responseRequestID {
		t.Error("Request ID in context and header should match")
	}
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
	if capturedRequestID != existingRequestID {
		t.Errorf("Should use existing request ID from header, got %s", capturedRequestID)
	}

	responseRequestID := w.Header().Get(middleware.RequestIDHeader)
	if responseRequestID != existingRequestID {
		t.Errorf("Should return existing request ID in header, got %s", responseRequestID)
	}
}

func TestGetRequestIDFromEmptyContext(t *testing.T) {
	ctx := context.Background()
	requestID := middleware.GetRequestID(ctx)
	if requestID != "" {
		t.Error("Should return empty string for context without request ID")
	}
}

func TestDefaultCORSConfig(t *testing.T) {
	config := middleware.DefaultCORSConfig()

	if len(config.AllowedOrigins) != 1 || config.AllowedOrigins[0] != "*" {
		t.Error("Default config should allow all origins")
	}

	expectedMethods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	if len(config.AllowedMethods) != len(expectedMethods) {
		t.Error("Default config should have standard HTTP methods")
	}

	if config.MaxAge != 86400 {
		t.Error("Default config should have 24 hour max age")
	}

	if config.AllowCredentials {
		t.Error("Default config should not allow credentials")
	}
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

	if w.Header().Get("Access-Control-Allow-Origin") != "https://api.example.com" {
		t.Error("Should allow subdomain matching wildcard pattern")
	}

	// Test domain that should not be allowed
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.Header.Set("Origin", "https://notexample.com")
	w2 := httptest.NewRecorder()

	handler.ServeHTTP(w2, req2)

	if w2.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Error("Should not allow domain that doesn't match wildcard pattern")
	}
}