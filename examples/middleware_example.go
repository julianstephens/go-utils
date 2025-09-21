package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/julianstephens/go-utils/httputil/middleware"
)

func main() {
	// Create loggers
	logger := log.New(os.Stdout, "[HTTP] ", log.LstdFlags)
	errorLogger := log.New(os.Stderr, "[ERROR] ", log.LstdFlags)

	// Create router
	router := mux.NewRouter()

	// Apply middleware in the correct order
	router.Use(middleware.RequestID())                          // First - adds ID to all requests
	router.Use(middleware.Logging(logger))                      // Second - logs with request ID
	router.Use(middleware.Recovery(errorLogger))                // Third - catches panics
	router.Use(middleware.CORS(middleware.DefaultCORSConfig())) // Fourth - handles CORS

	// Add routes
	router.HandleFunc("/api/health", healthHandler).Methods("GET")
	router.HandleFunc("/api/user/{id}", getUserHandler).Methods("GET")
	router.HandleFunc("/api/panic", panicHandler).Methods("GET")

	// Start server
	log.Println("Server starting on :8080")
	log.Println("Try these endpoints:")
	log.Println("  GET http://localhost:8080/api/health")
	log.Println("  GET http://localhost:8080/api/user/123")
	log.Println("  GET http://localhost:8080/api/panic (to test recovery)")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Get request ID from context
	requestID := middleware.GetRequestID(r.Context())

	response := map[string]interface{}{
		"status":     "ok",
		"request_id": requestID,
		"message":    "Service is healthy",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	requestID := middleware.GetRequestID(r.Context())

	response := map[string]interface{}{
		"user_id":    userID,
		"request_id": requestID,
		"name":       "John Doe",
		"email":      "john@example.com",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func panicHandler(w http.ResponseWriter, r *http.Request) {
	// This will panic and be caught by the recovery middleware
	panic("This is a test panic!")
}
