package jsonutil_test

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/julianstephens/go-utils/jsonutil"
)

// ExampleMarshal demonstrates basic JSON marshaling with error context.
func ExampleMarshal() {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	person := Person{Name: "Alice", Age: 30}
	data, err := jsonutil.Marshal(person)
	if err != nil {
		log.Fatalf("Marshal failed: %v", err)
	}

	fmt.Println(string(data))
	// Output: {"name":"Alice","age":30}
}

// ExampleMarshalIndent demonstrates pretty-printed JSON marshaling.
func ExampleMarshalIndent() {
	type Person struct {
		Name    string   `json:"name"`
		Age     int      `json:"age"`
		Hobbies []string `json:"hobbies"`
	}

	person := Person{
		Name:    "Bob",
		Age:     25,
		Hobbies: []string{"reading", "coding", "hiking"},
	}

	data, err := jsonutil.MarshalIndent(person, "", "  ")
	if err != nil {
		log.Fatalf("MarshalIndent failed: %v", err)
	}

	fmt.Println(string(data))
	// Output: {
	//   "name": "Bob",
	//   "age": 25,
	//   "hobbies": [
	//     "reading",
	//     "coding",
	//     "hiking"
	//   ]
	// }
}

// ExampleMarshalWithOptions demonstrates marshaling with custom options.
func ExampleMarshalWithOptions() {
	data := map[string]interface{}{
		"message": "Hello <script>alert('xss')</script>",
		"count":   42,
	}

	// Marshal without HTML escaping
	opts := &jsonutil.MarshalOptions{
		EscapeHTML: false,
		Indent:     "  ",
	}

	result, err := jsonutil.MarshalWithOptions(data, opts)
	if err != nil {
		log.Fatalf("MarshalWithOptions failed: %v", err)
	}

	fmt.Println(string(result))
	// Output: {
	//   "count": 42,
	//   "message": "Hello <script>alert('xss')</script>"
	// }
}

// ExampleUnmarshalStrict demonstrates strict JSON unmarshaling.
func ExampleUnmarshalStrict() {
	type Config struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}

	// This JSON has an unknown field that will cause strict unmarshaling to fail
	jsonData := `{"host":"localhost","port":8080,"unknown":"value"}`

	var config Config
	err := jsonutil.UnmarshalStrict([]byte(jsonData), &config)
	if err != nil {
		fmt.Printf("Strict unmarshal failed: %v\n", err)
		return
	}

	fmt.Printf("Config: %+v\n", config)
	// Output: Strict unmarshal failed: jsonutil: strict unmarshal failed: json: unknown field "unknown"
}

// ExampleDecodeReader demonstrates streaming JSON decoding.
func ExampleDecodeReader() {
	jsonStream := `{"name":"Charlie","age":35,"active":true}`
	reader := strings.NewReader(jsonStream)

	type User struct {
		Name   string `json:"name"`
		Age    int    `json:"age"`
		Active bool   `json:"active"`
	}

	var user User
	err := jsonutil.DecodeReader(reader, &user, nil)
	if err != nil {
		log.Fatalf("DecodeReader failed: %v", err)
	}

	fmt.Printf("User: %+v\n", user)
	// Output: User: {Name:Charlie Age:35 Active:true}
}

// ExampleEncodeWriter demonstrates streaming JSON encoding.
func ExampleEncodeWriter() {
	type Product struct {
		ID    int     `json:"id"`
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}

	product := Product{ID: 1, Name: "Laptop", Price: 999.99}

	var buf bytes.Buffer
	opts := &jsonutil.EncoderOptions{
		Indent:     "  ",
		EscapeHTML: false,
	}

	err := jsonutil.EncodeWriter(&buf, product, opts)
	if err != nil {
		log.Fatalf("EncodeWriter failed: %v", err)
	}

	fmt.Print(buf.String())
	// Output: {
	//   "id": 1,
	//   "name": "Laptop",
	//   "price": 999.99
	// }
}

// ExampleValid demonstrates JSON validation.
func ExampleValid() {
	validJSON := `{"name":"Alice","age":30}`
	invalidJSON := `{"name":"Alice","age":30,}`

	fmt.Printf("Valid JSON: %t\n", jsonutil.Valid([]byte(validJSON)))
	fmt.Printf("Invalid JSON: %t\n", jsonutil.Valid([]byte(invalidJSON)))
	// Output: Valid JSON: true
	// Invalid JSON: false
}

// ExampleCompact demonstrates JSON compacting.
func ExampleCompact() {
	indentedJSON := `{
    "name": "Alice",
    "age": 30,
    "hobbies": [
        "reading",
        "coding"
    ]
}`

	var buf bytes.Buffer
	err := jsonutil.Compact(&buf, []byte(indentedJSON))
	if err != nil {
		log.Fatalf("Compact failed: %v", err)
	}

	fmt.Println(buf.String())
	// Output: {"name":"Alice","age":30,"hobbies":["reading","coding"]}
}
