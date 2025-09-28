package response

import (
	"errors"
	"net/http"
)

// Encoder defines the interface for encoding response data to an HTTP response writer.
type Encoder interface {
	Encode(w http.ResponseWriter, v any, status int) error
}

// BeforeFunc is called before encoding the response.
// It can be used for logging, metrics, or modifying the response writer.
type BeforeFunc func(w http.ResponseWriter, r *http.Request, data any)

// AfterFunc is called after successfully encoding the response.
// It can be used for cleanup, logging, or metrics collection.
type AfterFunc func(w http.ResponseWriter, r *http.Request, data any)

// OnErrorFunc is called when an error occurs during response processing.
// It should handle the error appropriately, typically by writing an error response.
type OnErrorFunc func(w http.ResponseWriter, r *http.Request, err error, status int)

// Responder provides structured HTTP response handling with extensible hooks.
type Responder struct {
	Encoder Encoder     // Encoder for response data
	Before  BeforeFunc  // Hook called before encoding
	After   AfterFunc   // Hook called after successful encoding
	OnError OnErrorFunc // Hook called on encoding errors
}

// Error represents a structured error response.
type Error struct {
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

// WriteWithStatus writes a response with a specific HTTP status code.
// It sets the status code before calling the encoder.
func (r *Responder) WriteWithStatus(w http.ResponseWriter, req *http.Request, data any, statusCode int) {
	if r.Before != nil {
		r.Before(w, req, data)
	}

	if err := r.Encoder.Encode(w, data, statusCode); err != nil {
		if r.OnError != nil {
			r.OnError(w, req, err, statusCode)
		}
		return
	}

	if r.After != nil {
		r.After(w, req, data)
	}
}

// ErrorWithStatus handles error responses by calling the OnError hook.
func (r *Responder) ErrorWithStatus(w http.ResponseWriter, req *http.Request, status int, err error, details *map[string]any) {
	if r.OnError != nil {
		r.OnError(w, req, err, status)
		return
	}

	if w == nil {
		return
	}

	msg := http.StatusText(http.StatusInternalServerError)
	if err != nil {
		msg = err.Error()
	}

	if status == 0 {
		status = http.StatusInternalServerError
	}

	r.Encoder.Encode(w, Error{Message: msg, Details: map[string]any{
		"status": http.StatusText(status),
	}}, status)
}

// OK writes a response with HTTP 200 OK status.
func (r *Responder) OK(w http.ResponseWriter, req *http.Request, data any) {
	r.WriteWithStatus(w, req, data, http.StatusOK)
}

// Created writes a response with HTTP 201 Created status.
func (r *Responder) Created(w http.ResponseWriter, req *http.Request, data any) {
	r.WriteWithStatus(w, req, data, http.StatusCreated)
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
func (r *Responder) BadRequest(w http.ResponseWriter, req *http.Request, data any, details *map[string]any) {
	r.ErrorWithStatus(w, req, http.StatusBadRequest, ParseErrData(data), details)
}

// Unauthorized writes a response with HTTP 401 Unauthorized status.
func (r *Responder) Unauthorized(w http.ResponseWriter, req *http.Request, data any, details *map[string]any) {
	r.ErrorWithStatus(w, req, http.StatusUnauthorized, ParseErrData(data), details)
}

// Forbidden writes a response with HTTP 403 Forbidden status.
func (r *Responder) Forbidden(w http.ResponseWriter, req *http.Request, data any, details *map[string]any) {
	r.ErrorWithStatus(w, req, http.StatusForbidden, ParseErrData(data), details)
}

// NotFound writes a response with HTTP 404 Not Found status.
func (r *Responder) NotFound(w http.ResponseWriter, req *http.Request, data any, details *map[string]any) {
	r.ErrorWithStatus(w, req, http.StatusNotFound, ParseErrData(data), details)
}

// InternalServerError writes a response with HTTP 500 Internal Server Error status.
func (r *Responder) InternalServerError(w http.ResponseWriter, req *http.Request, data any, details *map[string]any) {
	r.ErrorWithStatus(w, req, http.StatusInternalServerError, ParseErrData(data), details)
}

func ParseErrData(data any) error {
	var err error
	if data != nil {
		if e, ok := data.(error); ok {
			err = e
		}
		if dataStr, ok := data.(string); ok {
			err = errors.New(dataStr)
		}
	}

	return err
}
