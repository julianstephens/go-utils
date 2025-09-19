package request

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// ErrInvalidContentType is returned when the Content-Type is not application/json.
var ErrInvalidContentType = errors.New("invalid content-type, expected application/json")

// DecodeJSON decodes a JSON request body into the given destination structure.
// It returns ErrInvalidContentType if the Content-Type is not application/json.
func DecodeJSON(r *http.Request, dst any) error {
	ct := r.Header.Get("Content-Type")
	if ct != "" && ct != "application/json" && ct != "application/json; charset=utf-8" {
		return ErrInvalidContentType
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

// ParseForm parses form values into the request, returning an error if parsing fails.
func ParseForm(r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	return nil
}

// FormValue returns a string form value and a boolean indicating if it was present.
func FormValue(r *http.Request, key string) (string, bool) {
	val := r.FormValue(key)
	return val, val != ""
}

// QueryValue returns a string query param and a boolean indicating if it was present.
func QueryValue(r *http.Request, key string) (string, bool) {
	vals := r.URL.Query()[key]
	if len(vals) == 0 {
		return "", false
	}
	return vals[0], true
}

// QueryInt returns an int query param and a boolean indicating if it was present and valid.
func QueryInt(r *http.Request, key string) (int, bool) {
	vals := r.URL.Query()[key]
	if len(vals) == 0 {
		return 0, false
	}
	i, err := strconv.Atoi(vals[0])
	if err != nil {
		return 0, false
	}
	return i, true
}

// QueryBool returns a bool query param and a boolean indicating if it was present and valid.
func QueryBool(r *http.Request, key string) (bool, bool) {
	vals := r.URL.Query()[key]
	if len(vals) == 0 {
		return false, false
	}
	b, err := strconv.ParseBool(vals[0])
	if err != nil {
		return false, false
	}
	return b, true
}

// QueryValues returns all values for a query parameter.
func QueryValues(r *http.Request, key string) []string {
	return r.URL.Query()[key]
}

// FormValues returns all values for a form parameter.
func FormValues(r *http.Request, key string) []string {
	return r.Form[key]
}

// ParseQuery parses a raw query string into url.Values.
func ParseQuery(raw string) (url.Values, error) {
	return url.ParseQuery(raw)
}