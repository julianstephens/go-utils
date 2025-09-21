/*
Package jsonutil provides convenient, idiomatic helpers for working with JSON in Go projects.

This package offers enhanced JSON marshaling and unmarshaling functionality with features like:
- Error context and formatting options (indentation, HTML escaping control)
- Stream-based encoding/decoding for io.Reader and io.Writer
- Configurable strict decoding (error on unknown fields)
- Comprehensive error handling with meaningful context

Basic Usage:

	package main

	import (
		"fmt"
		"log"

		"github.com/julianstephens/go-utils/jsonutil"
	)

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	func main() {
		person := Person{Name: "Alice", Age: 30}

		// Marshal with indentation
		data, err := jsonutil.MarshalIndent(person, "", "  ")
		if err != nil {
			log.Fatalf("Marshal failed: %v", err)
		}
		fmt.Println(string(data))

		// Unmarshal with strict decoding
		var decoded Person
		err = jsonutil.UnmarshalStrict(data, &decoded)
		if err != nil {
			log.Fatalf("Unmarshal failed: %v", err)
		}
	}

Stream Processing:

	func processJSONStream(r io.Reader, w io.Writer) error {
		var data map[string]interface{}

		// Decode from reader
		if err := jsonutil.DecodeReader(r, &data); err != nil {
			return fmt.Errorf("decode failed: %w", err)
		}

		// Encode to writer with pretty printing
		opts := &jsonutil.EncoderOptions{
			Indent:     "  ",
			EscapeHTML: false,
		}
		return jsonutil.EncodeWriter(w, data, opts)
	}

Error Context:

All functions provide enhanced error context to help with debugging:

	// Instead of: "invalid character 'x' looking for beginning of value"
	// You get: "jsonutil: unmarshal failed: invalid character 'x' looking for beginning of value"

The package is designed for reuse across Go projects and integrates well with existing
JSON workflows while providing additional safety and convenience features.
*/
package jsonutil
