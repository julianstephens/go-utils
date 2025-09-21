package jsonutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// MarshalOptions holds configuration options for JSON marshaling.
type MarshalOptions struct {
	// Indent specifies the indentation string for pretty-printing.
	// Empty string means no indentation.
	Indent string
	// Prefix specifies the prefix for each line when indenting.
	Prefix string
	// EscapeHTML specifies whether HTML characters should be escaped.
	// Default is true to match json.Marshal behavior.
	EscapeHTML bool
}

// UnmarshalOptions holds configuration options for JSON unmarshaling.
type UnmarshalOptions struct {
	// DisallowUnknownFields causes the decoder to return an error when
	// the destination struct has fields that don't match any field in the JSON.
	DisallowUnknownFields bool
	// UseNumber causes the decoder to decode JSON numbers into json.Number
	// instead of float64.
	UseNumber bool
}

// EncoderOptions holds configuration options for stream encoding.
type EncoderOptions struct {
	// Indent specifies the indentation string for pretty-printing.
	Indent string
	// Prefix specifies the prefix for each line when indenting.
	Prefix string
	// EscapeHTML specifies whether HTML characters should be escaped.
	EscapeHTML bool
}

// DecoderOptions holds configuration options for stream decoding.
type DecoderOptions struct {
	// DisallowUnknownFields causes the decoder to return an error when
	// the destination struct has fields that don't match any field in the JSON.
	DisallowUnknownFields bool
	// UseNumber causes the decoder to decode JSON numbers into json.Number
	// instead of float64.
	UseNumber bool
}

// Marshal marshals the given value to JSON with enhanced error context.
// It provides the same functionality as json.Marshal but with better error messages.
func Marshal(v any) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("jsonutil: marshal failed: %w", err)
	}
	return data, nil
}

// MarshalIndent marshals the given value to JSON with indentation.
// It's equivalent to json.MarshalIndent but with enhanced error context.
func MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	data, err := json.MarshalIndent(v, prefix, indent)
	if err != nil {
		return nil, fmt.Errorf("jsonutil: marshal indent failed: %w", err)
	}
	return data, nil
}

// MarshalWithOptions marshals the given value to JSON using the provided options.
func MarshalWithOptions(v any, opts *MarshalOptions) ([]byte, error) {
	if opts == nil {
		return Marshal(v)
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(opts.EscapeHTML)

	if opts.Indent != "" {
		encoder.SetIndent(opts.Prefix, opts.Indent)
	}

	if err := encoder.Encode(v); err != nil {
		return nil, fmt.Errorf("jsonutil: marshal with options failed: %w", err)
	}

	// Remove the trailing newline that Encoder.Encode adds
	data := buf.Bytes()
	if len(data) > 0 && data[len(data)-1] == '\n' {
		data = data[:len(data)-1]
	}

	return data, nil
}

// Unmarshal unmarshals JSON data into the given value with enhanced error context.
func Unmarshal(data []byte, v any) error {
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("jsonutil: unmarshal failed: %w", err)
	}
	return nil
}

// UnmarshalStrict unmarshals JSON data with strict field matching.
// It returns an error if the JSON contains fields not present in the destination struct.
func UnmarshalStrict(data []byte, v any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("jsonutil: strict unmarshal failed: %w", err)
	}
	return nil
}

// UnmarshalWithOptions unmarshals JSON data using the provided options.
func UnmarshalWithOptions(data []byte, v any, opts *UnmarshalOptions) error {
	if opts == nil {
		return Unmarshal(data, v)
	}

	decoder := json.NewDecoder(bytes.NewReader(data))

	if opts.DisallowUnknownFields {
		decoder.DisallowUnknownFields()
	}

	if opts.UseNumber {
		decoder.UseNumber()
	}

	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("jsonutil: unmarshal with options failed: %w", err)
	}
	return nil
}

// EncodeWriter encodes the given value as JSON and writes it to the provided writer.
func EncodeWriter(w io.Writer, v any, opts *EncoderOptions) error {
	encoder := json.NewEncoder(w)

	if opts != nil {
		encoder.SetEscapeHTML(opts.EscapeHTML)
		if opts.Indent != "" {
			encoder.SetIndent(opts.Prefix, opts.Indent)
		}
	}

	if err := encoder.Encode(v); err != nil {
		return fmt.Errorf("jsonutil: encode to writer failed: %w", err)
	}
	return nil
}

// DecodeReader decodes JSON from the provided reader into the given value.
func DecodeReader(r io.Reader, v any, opts *DecoderOptions) error {
	decoder := json.NewDecoder(r)

	if opts != nil {
		if opts.DisallowUnknownFields {
			decoder.DisallowUnknownFields()
		}
		if opts.UseNumber {
			decoder.UseNumber()
		}
	}

	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("jsonutil: decode from reader failed: %w", err)
	}
	return nil
}

// DecodeReaderStrict decodes JSON from the provided reader with strict field matching.
func DecodeReaderStrict(r io.Reader, v any) error {
	opts := &DecoderOptions{DisallowUnknownFields: true}
	return DecodeReader(r, v, opts)
}

// Valid reports whether data is a valid JSON encoding.
func Valid(data []byte) bool {
	return json.Valid(data)
}

// Compact appends to dst the JSON-encoded src with insignificant space characters elided.
func Compact(dst *bytes.Buffer, src []byte) error {
	if err := json.Compact(dst, src); err != nil {
		return fmt.Errorf("jsonutil: compact failed: %w", err)
	}
	return nil
}

// Indent appends to dst an indented form of the JSON-encoded src.
func Indent(dst *bytes.Buffer, src []byte, prefix, indent string) error {
	if err := json.Indent(dst, src, prefix, indent); err != nil {
		return fmt.Errorf("jsonutil: indent failed: %w", err)
	}
	return nil
}

// HTMLEscape appends to dst the JSON-encoded src with HTML metacharacters escaped.
func HTMLEscape(dst *bytes.Buffer, src []byte) {
	json.HTMLEscape(dst, src)
}
