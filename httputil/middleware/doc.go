/*
Package middleware provides common, reusable middleware for Go mux-based routers.

This package offers a collection of HTTP middleware that can be easily composed
and used with any HTTP router that follows the standard Go HTTP middleware pattern.

Available Middleware:

	• Logging - Logs HTTP requests and responses
	• Recovery - Recovers from panics and returns HTTP 500
	• CORS - Handles Cross-Origin Resource Sharing
	• RequestID - Injects unique request IDs into requests

Basic Usage:

	package main

	import (
		"log"
		"net/http"
		"os"

		"github.com/gorilla/mux"
		"github.com/julianstephens/go-utils/httputil/middleware"
	)

	func main() {
		logger := log.New(os.Stdout, "[HTTP] ", log.LstdFlags)

		// Create router
		router := mux.NewRouter()

		// Apply middleware
		router.Use(middleware.RequestID())
		router.Use(middleware.Logging(logger))
		router.Use(middleware.Recovery(logger))
		router.Use(middleware.CORS(middleware.DefaultCORSConfig()))

		// Add routes
		router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
			requestID := middleware.GetRequestID(r.Context())
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status":"ok","request_id":"` + requestID + `"}`))
		}).Methods("GET")

		log.Println("Server starting on :8080")
		log.Fatal(http.ListenAndServe(":8080", router))
	}

Logging Middleware:

The logging middleware logs each HTTP request with method, path, status code, and duration:

	logger := log.New(os.Stdout, "[API] ", log.LstdFlags)
	router.Use(middleware.Logging(logger))

Output example: [API] 2023/10/01 12:00:00 GET /api/users 200 1.234ms

Recovery Middleware:

The recovery middleware catches panics and returns a generic HTTP 500 error:

	logger := log.New(os.Stderr, "[ERROR] ", log.LstdFlags)
	router.Use(middleware.Recovery(logger))

	router.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("something went wrong") // Will be caught and logged
	})

CORS Middleware:

The CORS middleware handles Cross-Origin Resource Sharing with configurable options:

	// Using default configuration (allows all origins)
	router.Use(middleware.CORS(middleware.DefaultCORSConfig()))

	// Custom configuration
	corsConfig := middleware.CORSConfig{
		AllowedOrigins:   []string{"https://example.com", "https://app.example.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"X-Total-Count"},
		AllowCredentials: true,
		MaxAge:          3600,
	}
	router.Use(middleware.CORS(corsConfig))

	// Wildcard subdomain support
	corsConfig := middleware.CORSConfig{
		AllowedOrigins: []string{"*.example.com"},
	}

Request ID Middleware:

The request ID middleware injects a unique identifier into each request:

	router.Use(middleware.RequestID())

	router.HandleFunc("/api/endpoint", func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetRequestID(r.Context())
		log.Printf("Processing request %s", requestID)
		
		// Request ID is also automatically added to response headers as X-Request-ID
		w.Write([]byte("OK"))
	})

Middleware Composition:

Middleware should be applied in the correct order for optimal functionality:

	// Recommended order:
	router.Use(middleware.RequestID())    // First - adds ID to all requests
	router.Use(middleware.Logging(logger)) // Second - logs with request ID
	router.Use(middleware.Recovery(logger)) // Third - catches panics
	router.Use(middleware.CORS(corsConfig)) // Fourth - handles CORS

Custom Middleware:

All middleware in this package follows the standard Go HTTP middleware pattern,
making them compatible with other middleware libraries:

	func customMiddleware(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Pre-processing
			start := time.Now()
			
			next.ServeHTTP(w, r)
			
			// Post-processing
			duration := time.Since(start)
			log.Printf("Request took %v", duration)
		})
	}

	router.Use(customMiddleware)

The middleware can be easily integrated with existing applications and provides
a consistent, maintainable approach to common HTTP concerns.
*/
package middleware