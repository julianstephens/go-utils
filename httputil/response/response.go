package response

import (
	"net/http"
)

// Encoder defines the interface for encoding response data to an HTTP response writer.
type Encoder interface {
	Encode(w http.ResponseWriter, v any) error
}

// BeforeFunc is called before encoding the response.
// It can be used for logging, metrics, or modifying the response writer.
type BeforeFunc func(w http.ResponseWriter, r *http.Request, data any)

// AfterFunc is called after successfully encoding the response.
// It can be used for cleanup, logging, or metrics collection.
type AfterFunc func(w http.ResponseWriter, r *http.Request, data any)

// OnErrorFunc is called when an error occurs during response processing.
// It should handle the error appropriately, typically by writing an error response.
type OnErrorFunc func(w http.ResponseWriter, r *http.Request, err error)

// Responder provides structured HTTP response handling with extensible hooks.
type Responder struct {
	Encoder Encoder     // Encoder for response data
	Before  BeforeFunc  // Hook called before encoding
	After   AfterFunc   // Hook called after successful encoding
	OnError OnErrorFunc // Hook called on encoding errors
}

// Write encodes and writes a successful response using the configured encoder and hooks.
// It calls Before hook, encodes the data, then calls After hook on success or OnError on failure.
func (r *Responder) Write(w http.ResponseWriter, req *http.Request, data any) {
	if r.Before != nil {
		r.Before(w, req, data)
	}

	if err := r.Encoder.Encode(w, data); err != nil {
		if r.OnError != nil {
			r.OnError(w, req, err)
		}
		return
	}

	if r.After != nil {
		r.After(w, req, data)
	}
}

// writeWithStatus writes a response with a specific HTTP status code.
// It sets the status code before calling the encoder.
func (r *Responder) writeWithStatus(w http.ResponseWriter, req *http.Request, data any, statusCode int) {
	if r.Before != nil {
		r.Before(w, req, data)
	}

	if err := r.Encoder.Encode(w, data); err != nil {
		if r.OnError != nil {
			r.OnError(w, req, err)
		}
		return
	}
	w.WriteHeader(statusCode)

	if r.After != nil {
		r.After(w, req, data)
	}
}

// Error handles error responses by calling the OnError hook.
// If no OnError hook is configured, it writes a basic 500 Internal Server Error response.
func (r *Responder) Error(w http.ResponseWriter, req *http.Request, err error) {
	if r.OnError != nil {
		r.OnError(w, req, err)
		return
	}

	// Default error handling if no OnError hook is provided
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

// OK writes a response with HTTP 200 OK status.
func (r *Responder) OK(w http.ResponseWriter, req *http.Request, data any) {
	r.writeWithStatus(w, req, data, http.StatusOK)
}

// Created writes a response with HTTP 201 Created status.
func (r *Responder) Created(w http.ResponseWriter, req *http.Request, data any) {
	r.writeWithStatus(w, req, data, http.StatusCreated)
}

// NoContent writes a response with HTTP 204 No Content status.
func (r *Responder) NoContent(w http.ResponseWriter, req *http.Request) {
	if r.Before != nil {
		r.Before(w, req, nil)
	}

	w.WriteHeader(http.StatusNoContent)
	if r.After != nil {
		r.After(w, req, nil)
	}
}

// BadRequest writes a response with HTTP 400 Bad Request status.
func (r *Responder) BadRequest(w http.ResponseWriter, req *http.Request, data any) {
	r.writeWithStatus(w, req, data, http.StatusBadRequest)
}

// Unauthorized writes a response with HTTP 401 Unauthorized status.
func (r *Responder) Unauthorized(w http.ResponseWriter, req *http.Request, data any) {
	r.writeWithStatus(w, req, data, http.StatusUnauthorized)
}

// Forbidden writes a response with HTTP 403 Forbidden status.
func (r *Responder) Forbidden(w http.ResponseWriter, req *http.Request, data any) {
	r.writeWithStatus(w, req, data, http.StatusForbidden)
}

// NotFound writes a response with HTTP 404 Not Found status.
func (r *Responder) NotFound(w http.ResponseWriter, req *http.Request, data any) {
	r.writeWithStatus(w, req, data, http.StatusNotFound)
}

// InternalServerError writes a response with HTTP 500 Internal Server Error status.
func (r *Responder) InternalServerError(w http.ResponseWriter, req *http.Request, data any) {
	r.writeWithStatus(w, req, data, http.StatusInternalServerError)
}
