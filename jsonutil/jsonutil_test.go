package jsonutil_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/julianstephens/go-utils/jsonutil"
	tst "github.com/julianstephens/go-utils/tests"
)

type testStruct struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Email   string   `json:"email,omitempty"`
	Active  bool     `json:"active"`
	Balance *float64 `json:"balance,omitempty"`
}

type strictTestStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestMarshal(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantErr bool
	}{
		{
			name:    "valid struct",
			input:   testStruct{Name: "Alice", Age: 30, Active: true},
			wantErr: false,
		},
		{
			name:    "nil pointer",
			input:   (*testStruct)(nil),
			wantErr: false,
		},
		{
			name:    "map",
			input:   map[string]any{"key": "value", "number": 42},
			wantErr: false,
		},
		{
			name:    "slice",
			input:   []string{"a", "b", "c"},
			wantErr: false,
		},
		{
			name:    "invalid channel",
			input:   make(chan int),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := jsonutil.Marshal(tt.input)
			if (err != nil) != tt.wantErr {
				tst.AssertTrue(t, (err != nil) == tt.wantErr, "Marshal() error mismatch")
				return
			}
			if !tt.wantErr {
				tst.AssertTrue(t, jsonutil.Valid(data), "Marshal() produced invalid JSON")
			}
		})
	}
}

func TestMarshalIndent(t *testing.T) {
	input := testStruct{Name: "Alice", Age: 30, Active: true}

	data, err := jsonutil.MarshalIndent(input, "", "  ")
	tst.AssertNoError(t, err)
	tst.AssertTrue(t, strings.Contains(string(data), "\n"), "MarshalIndent() should produce indented output")
	tst.AssertTrue(t, jsonutil.Valid(data), "MarshalIndent() produced invalid JSON")
}

func TestMarshalWithOptions(t *testing.T) {
	input := map[string]any{
		"name":   "Alice & Bob",
		"script": "<script>alert('test')</script>",
		"age":    30,
	}

	tests := []struct {
		name    string
		opts    *jsonutil.MarshalOptions
		wantStr string
	}{
		{
			name:    "with HTML escaping",
			opts:    &jsonutil.MarshalOptions{EscapeHTML: true},
			wantStr: "\\u003cscript\\u003e",
		},
		{
			name:    "without HTML escaping",
			opts:    &jsonutil.MarshalOptions{EscapeHTML: false},
			wantStr: "<script>",
		},
		{
			name:    "with indentation",
			opts:    &jsonutil.MarshalOptions{Indent: "  ", EscapeHTML: false},
			wantStr: "\n",
		},
		{
			name:    "nil options",
			opts:    nil,
			wantStr: "\\u0026", // Default escaping behavior
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := jsonutil.MarshalWithOptions(input, tt.opts)
			tst.AssertNoError(t, err)
			result := string(data)
			tst.AssertTrue(
				t,
				strings.Contains(result, tt.wantStr),
				"MarshalWithOptions() result missing expected substring",
			)
			tst.AssertTrue(t, jsonutil.Valid(data), "MarshalWithOptions() produced invalid JSON")
		})
	}
}

func TestUnmarshal(t *testing.T) {
	validJSON := `{"name":"Alice","age":30,"active":true}`
	invalidJSON := `{"name":"Alice","age":30,}`

	tests := []struct {
		name    string
		data    string
		target  any
		wantErr bool
	}{
		{
			name:    "valid JSON into struct",
			data:    validJSON,
			target:  &testStruct{},
			wantErr: false,
		},
		{
			name:    "valid JSON into map",
			data:    validJSON,
			target:  &map[string]any{},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			data:    invalidJSON,
			target:  &testStruct{},
			wantErr: true,
		},
		{
			name:    "nil target",
			data:    validJSON,
			target:  nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := jsonutil.Unmarshal([]byte(tt.data), tt.target)
			tst.AssertTrue(t, (err != nil) == tt.wantErr, "Unmarshal() error mismatch")
		})
	}
}

func TestUnmarshalStrict(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		target  any
		wantErr bool
	}{
		{
			name:    "exact fields match",
			data:    `{"name":"Alice","age":30}`,
			target:  &strictTestStruct{},
			wantErr: false,
		},
		{
			name:    "unknown field should fail",
			data:    `{"name":"Alice","age":30,"unknown":"value"}`,
			target:  &strictTestStruct{},
			wantErr: true,
		},
		{
			name:    "missing optional field is OK",
			data:    `{"name":"Alice"}`,
			target:  &strictTestStruct{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := jsonutil.UnmarshalStrict([]byte(tt.data), tt.target)
			tst.AssertTrue(t, (err != nil) == tt.wantErr, "UnmarshalStrict() error mismatch")
		})
	}
}

func TestUnmarshalWithOptions(t *testing.T) {
	jsonWithNumber := `{"name":"Alice","age":30,"balance":123.45}`
	jsonWithUnknown := `{"name":"Alice","age":30,"unknown":"value"}`

	tests := []struct {
		name    string
		data    string
		opts    *jsonutil.UnmarshalOptions
		target  interface{}
		wantErr bool
		check   func(t *testing.T, target interface{})
	}{
		{
			name:    "use number option",
			data:    jsonWithNumber,
			opts:    &jsonutil.UnmarshalOptions{UseNumber: true},
			target:  &map[string]interface{}{},
			wantErr: false,
			check: func(t *testing.T, target interface{}) {
				m := target.(*map[string]interface{})
				if _, ok := (*m)["balance"].(json.Number); !ok {
					t.Error("Expected balance to be json.Number")
				}
			},
		},
		{
			name:    "disallow unknown fields",
			data:    jsonWithUnknown,
			opts:    &jsonutil.UnmarshalOptions{DisallowUnknownFields: true},
			target:  &strictTestStruct{},
			wantErr: true,
		},
		{
			name:    "nil options",
			data:    jsonWithNumber,
			opts:    nil,
			target:  &map[string]interface{}{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := jsonutil.UnmarshalWithOptions([]byte(tt.data), tt.target, tt.opts)
			tst.AssertTrue(t, (err != nil) == tt.wantErr, "UnmarshalWithOptions() error mismatch")
			if !tt.wantErr && tt.check != nil {
				tt.check(t, tt.target)
			}
		})
	}
}

func TestEncodeWriter(t *testing.T) {
	input := testStruct{Name: "Alice", Age: 30, Active: true}

	tests := []struct {
		name string
		opts *jsonutil.EncoderOptions
	}{
		{
			name: "basic encoding",
			opts: nil,
		},
		{
			name: "with indentation",
			opts: &jsonutil.EncoderOptions{Indent: "  "},
		},
		{
			name: "without HTML escaping",
			opts: &jsonutil.EncoderOptions{EscapeHTML: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := jsonutil.EncodeWriter(&buf, input, tt.opts)
			tst.AssertNoError(t, err)
			tst.AssertTrue(t, jsonutil.Valid(buf.Bytes()), "EncodeWriter() produced invalid JSON")
			if tt.opts != nil && tt.opts.Indent != "" {
				tst.AssertTrue(
					t,
					strings.Contains(buf.String(), "\n"),
					"EncodeWriter() should produce indented output when indent is specified",
				)
			}
		})
	}
}

func TestDecodeReader(t *testing.T) {
	validJSON := `{"name":"Alice","age":30,"active":true}`
	invalidJSON := `{"name":"Alice","age":30,}`

	tests := []struct {
		name    string
		data    string
		opts    *jsonutil.DecoderOptions
		target  interface{}
		wantErr bool
	}{
		{
			name:    "valid JSON",
			data:    validJSON,
			opts:    nil,
			target:  &testStruct{},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			data:    invalidJSON,
			opts:    nil,
			target:  &testStruct{},
			wantErr: true,
		},
		{
			name:    "strict decoding with unknown field",
			data:    `{"name":"Alice","age":30,"unknown":"value"}`,
			opts:    &jsonutil.DecoderOptions{DisallowUnknownFields: true},
			target:  &strictTestStruct{},
			wantErr: true,
		},
		{
			name:    "use number option",
			data:    `{"balance":123.45}`,
			opts:    &jsonutil.DecoderOptions{UseNumber: true},
			target:  &map[string]interface{}{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.data)
			err := jsonutil.DecodeReader(reader, tt.target, tt.opts)
			tst.AssertTrue(t, (err != nil) == tt.wantErr, "DecodeReader() error mismatch")
		})
	}
}

func TestDecodeReaderStrict(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		target  interface{}
		wantErr bool
	}{
		{
			name:    "exact fields match",
			data:    `{"name":"Alice","age":30}`,
			target:  &strictTestStruct{},
			wantErr: false,
		},
		{
			name:    "unknown field should fail",
			data:    `{"name":"Alice","age":30,"unknown":"value"}`,
			target:  &strictTestStruct{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.data)
			err := jsonutil.DecodeReaderStrict(reader, tt.target)
			tst.AssertTrue(t, (err != nil) == tt.wantErr, "DecodeReaderStrict() error mismatch")
		})
	}
}

func TestValid(t *testing.T) {
	tests := []struct {
		name string
		data string
		want bool
	}{
		{
			name: "valid JSON object",
			data: `{"name":"Alice","age":30}`,
			want: true,
		},
		{
			name: "valid JSON array",
			data: `[1,2,3]`,
			want: true,
		},
		{
			name: "invalid JSON",
			data: `{"name":"Alice","age":30,}`,
			want: false,
		},
		{
			name: "empty string",
			data: ``,
			want: false,
		},
		{
			name: "valid null",
			data: `null`,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := jsonutil.Valid([]byte(tt.data))
			tst.AssertDeepEqual(t, result, tt.want)
		})
	}
}

func TestCompact(t *testing.T) {
	indentedJSON := `{
    "name": "Alice",
    "age": 30
}`

	var buf bytes.Buffer
	err := jsonutil.Compact(&buf, []byte(indentedJSON))
	tst.AssertNoError(t, err)
	result := buf.String()
	tst.AssertFalse(
		t,
		strings.Contains(result, "\n") || strings.Contains(result, "  "),
		"Compact() should remove whitespace",
	)
	tst.AssertTrue(t, jsonutil.Valid(buf.Bytes()), "Compact() produced invalid JSON")
}

func TestIndent(t *testing.T) {
	compactJSON := `{"name":"Alice","age":30}`

	var buf bytes.Buffer
	err := jsonutil.Indent(&buf, []byte(compactJSON), "", "  ")
	tst.AssertNoError(t, err)
	result := buf.String()
	tst.AssertTrue(t, strings.Contains(result, "\n"), "Indent() should add newlines")
	tst.AssertTrue(t, jsonutil.Valid(buf.Bytes()), "Indent() produced invalid JSON")
}

func TestHTMLEscape(t *testing.T) {
	jsonWithHTML := `{"script":"<script>alert('test')</script>"}`

	var buf bytes.Buffer
	jsonutil.HTMLEscape(&buf, []byte(jsonWithHTML))

	result := buf.String()
	tst.AssertTrue(t, strings.Contains(result, "\\u003c"), "HTMLEscape() should escape HTML characters")
	tst.AssertTrue(t, jsonutil.Valid(buf.Bytes()), "HTMLEscape() produced invalid JSON")
}

// Test error messages contain proper context
func TestErrorContext(t *testing.T) {
	invalidJSON := `{"name":"Alice","age":30,}`

	err := jsonutil.Unmarshal([]byte(invalidJSON), &testStruct{})
	tst.AssertNotNil(t, err, "Expected error for invalid JSON")
	errStr := err.Error()
	tst.AssertTrue(t, strings.Contains(errStr, "jsonutil:"), "Error should contain 'jsonutil:' prefix")
}

// Benchmark tests
func BenchmarkMarshal(b *testing.B) {
	data := testStruct{Name: "Alice", Age: 30, Active: true}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := jsonutil.Marshal(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	jsonData := []byte(`{"name":"Alice","age":30,"active":true}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var target testStruct
		err := jsonutil.Unmarshal(jsonData, &target)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeWriter(b *testing.B) {
	data := testStruct{Name: "Alice", Age: 30, Active: true}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		err := jsonutil.EncodeWriter(&buf, data, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeReader(b *testing.B) {
	jsonData := `{"name":"Alice","age":30,"active":true}`

	for b.Loop() {
		reader := strings.NewReader(jsonData)
		var target testStruct
		err := jsonutil.DecodeReader(reader, &target, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}
