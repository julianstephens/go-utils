package middleware_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/julianstephens/go-utils/httputil/middleware"
)

// TestMiddlewareIntegration tests all middleware working together
func TestMiddlewareIntegration(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "[TEST] ", 0)

	// Create router with all middleware
	router := mux.NewRouter()
	router.Use(middleware.RequestID())
	router.Use(middleware.Logging(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.CORS(middleware.DefaultCORSConfig()))

	// Add test handlers
	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetRequestID(r.Context())
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok","request_id":"` + requestID + `"}`))
	}).Methods("GET")

	router.HandleFunc("/api/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}).Methods("GET")

	// Test successful request
	t.Run("successful request", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/health", nil)
		req.Header.Set("Origin", "https://example.com")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Check status
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		// Check CORS headers
		if w.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Error("CORS middleware should set Access-Control-Allow-Origin")
		}

		// Check Request ID header
		requestID := w.Header().Get(middleware.RequestIDHeader)
		if requestID == "" {
			t.Error("Request ID middleware should set X-Request-ID header")
		}

		// Check response content includes request ID
		body := w.Body.String()
		if !strings.Contains(body, requestID) {
			t.Error("Response should contain the request ID")
		}

		// Check logging output
		logOutput := buf.String()
		if !strings.Contains(logOutput, "GET") {
			t.Error("Logging middleware should log the request")
		}
		if !strings.Contains(logOutput, "/api/health") {
			t.Error("Logging middleware should log the path")
		}
		if !strings.Contains(logOutput, "200") {
			t.Error("Logging middleware should log the status code")
		}
	})

	// Test panic recovery
	t.Run("panic recovery", func(t *testing.T) {
		buf.Reset() // Clear previous logs

		req := httptest.NewRequest("GET", "/api/panic", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Check that panic was recovered
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d after panic, got %d", http.StatusInternalServerError, w.Code)
		}

		// Check error response
		body := w.Body.String()
		if !strings.Contains(body, "Internal Server Error") {
			t.Error("Recovery middleware should return 'Internal Server Error'")
		}

		// Check that panic was logged
		logOutput := buf.String()
		if !strings.Contains(logOutput, "Panic") {
			t.Error("Recovery middleware should log panics")
		}
		if !strings.Contains(logOutput, "test panic") {
			t.Error("Recovery middleware should log panic details")
		}

		// Request should still be logged even after panic
		if !strings.Contains(logOutput, "GET") {
			t.Error("Logging middleware should work even with panics")
		}
	})

	// Test OPTIONS preflight request
	t.Run("CORS preflight", func(t *testing.T) {
		// Test CORS middleware directly without router constraints
		corsHandler := middleware.CORS(
			middleware.DefaultCORSConfig(),
		)(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("should not reach here for OPTIONS"))
			}),
		)

		req := httptest.NewRequest("OPTIONS", "/api/health", nil)
		req.Header.Set("Origin", "https://example.com")
		req.Header.Set("Access-Control-Request-Method", "GET")
		w := httptest.NewRecorder()

		corsHandler.ServeHTTP(w, req)

		// Check that preflight request returns 204
		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status %d for OPTIONS request, got %d", http.StatusNoContent, w.Code)
		}

		// Check CORS headers are set
		if w.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Error("CORS middleware should set Access-Control-Allow-Origin for preflight")
		}

		if !strings.Contains(w.Header().Get("Access-Control-Allow-Methods"), "GET") {
			t.Error("CORS middleware should set Access-Control-Allow-Methods for preflight")
		}

		// Check that handler was not called
		body := w.Body.String()
		if body != "" {
			t.Error("OPTIONS request should not reach the handler")
		}
	})
}

// TestMiddlewareOrder tests that middleware is applied in the correct order
func TestMiddlewareOrder(t *testing.T) {
	var order []string

	// Create a custom middleware that tracks order
	trackingMiddleware := func(name string) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, name+"-before")
				next.ServeHTTP(w, r)
				order = append(order, name+"-after")
			})
		}
	}

	router := mux.NewRouter()
	router.Use(trackingMiddleware("first"))
	router.Use(middleware.RequestID())
	router.Use(trackingMiddleware("second"))
	router.Use(middleware.Logging(nil))
	router.Use(trackingMiddleware("third"))

	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Verify middleware order
	expectedOrder := []string{
		"first-before",
		"second-before",
		"third-before",
		"handler",
		"third-after",
		"second-after",
		"first-after",
	}

	if len(order) != len(expectedOrder) {
		t.Fatalf("Expected %d order entries, got %d: %v", len(expectedOrder), len(order), order)
	}

	for i, expected := range expectedOrder {
		if order[i] != expected {
			t.Errorf("Order[%d]: expected %s, got %s", i, expected, order[i])
		}
	}
}
