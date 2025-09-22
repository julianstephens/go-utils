package response

import (
	"context"
	"errors"
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

// WriteWithStatus writes a response with a specific HTTP status code.
// It sets the status code before calling the encoder.
func (r *Responder) WriteWithStatus(w http.ResponseWriter, req *http.Request, data any, statusCode int) {
	if r.Before != nil {
		r.Before(w, req, data)
	}

	if err := r.Encoder.Encode(w, data); err != nil {
		if r.OnError != nil {
			r.OnError(w, req, err)
		}
		return
	}
	if statusCode != 0 && statusCode != http.StatusOK {
		w.WriteHeader(statusCode)
	}

	if r.After != nil {
		r.After(w, req, data)
	}
}

// ErrorWithStatus handles error responses by calling the OnError hook.
// If no OnError hook is configured, it writes a basic 500 Internal Server Error response.
type contextKey string

const responseStatusKey contextKey = "response_status"

func (r *Responder) ErrorWithStatus(w http.ResponseWriter, req *http.Request, status int, err error) {
	if r.OnError != nil {
		// Pass status via request context for the default handler
		if req != nil && status != 0 {
			ctx := req.Context()
			ctx = context.WithValue(ctx, responseStatusKey, status)
			req = req.WithContext(ctx)
		}
		r.OnError(w, req, err)
		return
	}

	// Defensive: handle nil ResponseWriter
	if w == nil {
		return
	}

	// Defensive: handle nil error
	msg := "internal server error"
	if err != nil {
		msg = err.Error()
	}

	// Defensive: handle zero or invalid status
	if status == 0 {
		status = http.StatusInternalServerError
	}

	http.Error(w, msg, status)
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
func (r *Responder) BadRequest(w http.ResponseWriter, req *http.Request, data any) {
	r.ErrorWithStatus(w, req, http.StatusBadRequest, parseErrData(data))
}

// Unauthorized writes a response with HTTP 401 Unauthorized status.
func (r *Responder) Unauthorized(w http.ResponseWriter, req *http.Request, data any) {
	r.ErrorWithStatus(w, req, http.StatusUnauthorized, parseErrData(data))
}

// Forbidden writes a response with HTTP 403 Forbidden status.
func (r *Responder) Forbidden(w http.ResponseWriter, req *http.Request, data any) {
	r.ErrorWithStatus(w, req, http.StatusForbidden, parseErrData(data))
}

// NotFound writes a response with HTTP 404 Not Found status.
func (r *Responder) NotFound(w http.ResponseWriter, req *http.Request, data any) {
	r.ErrorWithStatus(w, req, http.StatusNotFound, parseErrData(data))
}

// InternalServerError writes a response with HTTP 500 Internal Server Error status.
func (r *Responder) InternalServerError(w http.ResponseWriter, req *http.Request, data any) {
	r.ErrorWithStatus(w, req, http.StatusInternalServerError, parseErrData(data))
}

func parseErrData(data any) error {
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
