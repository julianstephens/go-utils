package response

import (
	"log"
	"net/http"
)

// DefaultBefore is a default Before hook that performs no operation.
func DefaultBefore(w http.ResponseWriter, r *http.Request, data any) {
	// No-op implementation
}

// DefaultAfter is a default After hook that performs no operation.
func DefaultAfter(w http.ResponseWriter, r *http.Request, data any) {
	// No-op implementation
}

// DefaultOnError is a default error handler that logs the error and returns a 500 status.
func DefaultOnError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Response encoding error: %v", err)
	if w == nil {
		return
	}
	msg := "internal server error"
	if err != nil {
		msg = err.Error()
	}
	// Use status from context if available (set by ErrorWithStatus)
	status := http.StatusInternalServerError
	if r != nil {
		if s, ok := r.Context().Value("response_status").(int); ok && s != 0 {
			status = s
		}
	}
	http.Error(w, msg, status)
}

// LoggingBefore is a Before hook that logs the incoming request.
func LoggingBefore(w http.ResponseWriter, r *http.Request, data any) {
	log.Printf("Handling request: %s %s", r.Method, r.URL.Path)
}

// LoggingAfter is an After hook that logs successful response completion.
func LoggingAfter(w http.ResponseWriter, r *http.Request, data any) {
	log.Printf("Response sent successfully for: %s %s", r.Method, r.URL.Path)
}

// LoggingOnError is an error handler that logs errors with request context.
func LoggingOnError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Error handling request %s %s: %v", r.Method, r.URL.Path, err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}
