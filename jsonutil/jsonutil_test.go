package jsonutil_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/julianstephens/go-utils/jsonutil"
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
		input   interface{}
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
			input:   map[string]interface{}{"key": "value", "number": 42},
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
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !jsonutil.Valid(data) {
				t.Errorf("Marshal() produced invalid JSON: %s", string(data))
			}
		})
	}
}

func TestMarshalIndent(t *testing.T) {
	input := testStruct{Name: "Alice", Age: 30, Active: true}

	data, err := jsonutil.MarshalIndent(input, "", "  ")
	if err != nil {
		t.Fatalf("MarshalIndent() error = %v", err)
	}

	if !strings.Contains(string(data), "\n") {
		t.Error("MarshalIndent() should produce indented output")
	}

	if !jsonutil.Valid(data) {
		t.Error("MarshalIndent() produced invalid JSON")
	}
}

func TestMarshalWithOptions(t *testing.T) {
	input := map[string]interface{}{
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
			if err != nil {
				t.Fatalf("MarshalWithOptions() error = %v", err)
			}

			result := string(data)
			if !strings.Contains(result, tt.wantStr) {
				t.Errorf("MarshalWithOptions() result = %s, want to contain %s", result, tt.wantStr)
			}

			if !jsonutil.Valid(data) {
				t.Error("MarshalWithOptions() produced invalid JSON")
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	validJSON := `{"name":"Alice","age":30,"active":true}`
	invalidJSON := `{"name":"Alice","age":30,}`

	tests := []struct {
		name    string
		data    string
		target  interface{}
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
			target:  &map[string]interface{}{},
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
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnmarshalStrict(t *testing.T) {
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
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalStrict() error = %v, wantErr %v", err, tt.wantErr)
			}
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
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalWithOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
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
			if err != nil {
				t.Fatalf("EncodeWriter() error = %v", err)
			}

			if !jsonutil.Valid(buf.Bytes()) {
				t.Error("EncodeWriter() produced invalid JSON")
			}

			if tt.opts != nil && tt.opts.Indent != "" && !strings.Contains(buf.String(), "\n") {
				t.Error("EncodeWriter() should produce indented output when indent is specified")
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
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeReader() error = %v, wantErr %v", err, tt.wantErr)
			}
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
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeReaderStrict() error = %v, wantErr %v", err, tt.wantErr)
			}
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
			if result != tt.want {
				t.Errorf("Valid() = %v, want %v", result, tt.want)
			}
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
	if err != nil {
		t.Fatalf("Compact() error = %v", err)
	}

	result := buf.String()
	if strings.Contains(result, "\n") || strings.Contains(result, "  ") {
		t.Error("Compact() should remove whitespace")
	}

	if !jsonutil.Valid(buf.Bytes()) {
		t.Error("Compact() produced invalid JSON")
	}
}

func TestIndent(t *testing.T) {
	compactJSON := `{"name":"Alice","age":30}`

	var buf bytes.Buffer
	err := jsonutil.Indent(&buf, []byte(compactJSON), "", "  ")
	if err != nil {
		t.Fatalf("Indent() error = %v", err)
	}

	result := buf.String()
	if !strings.Contains(result, "\n") {
		t.Error("Indent() should add newlines")
	}

	if !jsonutil.Valid(buf.Bytes()) {
		t.Error("Indent() produced invalid JSON")
	}
}

func TestHTMLEscape(t *testing.T) {
	jsonWithHTML := `{"script":"<script>alert('test')</script>"}`

	var buf bytes.Buffer
	jsonutil.HTMLEscape(&buf, []byte(jsonWithHTML))

	result := buf.String()
	if !strings.Contains(result, "\\u003c") {
		t.Error("HTMLEscape() should escape HTML characters")
	}

	if !jsonutil.Valid(buf.Bytes()) {
		t.Error("HTMLEscape() produced invalid JSON")
	}
}

// Test error messages contain proper context
func TestErrorContext(t *testing.T) {
	invalidJSON := `{"name":"Alice","age":30,}`

	err := jsonutil.Unmarshal([]byte(invalidJSON), &testStruct{})
	if err == nil {
		t.Fatal("Expected error for invalid JSON")
	}

	errStr := err.Error()
	if !strings.Contains(errStr, "jsonutil:") {
		t.Errorf("Error should contain 'jsonutil:' prefix, got: %s", errStr)
	}
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(jsonData)
		var target testStruct
		err := jsonutil.DecodeReader(reader, &target, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}
