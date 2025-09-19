package request_test

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"

	"github.com/julianstephens/go-utils/httputil/request"
)

type testStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestDecodeJSON(t *testing.T) {
	body := []byte(`{"name":"Alice","age":30}`)
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	var ts testStruct
	err := request.DecodeJSON(req, &ts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ts.Name != "Alice" || ts.Age != 30 {
		t.Errorf("unexpected struct: %+v", ts)
	}
}

func TestDecodeJSON_InvalidContentType(t *testing.T) {
	body := []byte(`{"name":"Alice","age":30}`)
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "text/plain")

	var ts testStruct
	err := request.DecodeJSON(req, &ts)
	if err != request.ErrInvalidContentType {
		t.Errorf("expected ErrInvalidContentType, got %v", err)
	}
}

func TestQueryValue(t *testing.T) {
	req, _ := http.NewRequest("GET", "/?foo=bar", nil)
	val, ok := request.QueryValue(req, "foo")
	if !ok || val != "bar" {
		t.Errorf("expected bar, got %s (%v)", val, ok)
	}
	_, ok = request.QueryValue(req, "baz")
	if ok {
		t.Errorf("expected missing value")
	}
}

func TestQueryInt(t *testing.T) {
	req, _ := http.NewRequest("GET", "/?num=42", nil)
	i, ok := request.QueryInt(req, "num")
	if !ok || i != 42 {
		t.Errorf("expected 42, got %d (%v)", i, ok)
	}
	_, ok = request.QueryInt(req, "bad")
	if ok {
		t.Errorf("expected missing value for 'bad'")
	}
	req, _ = http.NewRequest("GET", "/?num=foo", nil)
	_, ok = request.QueryInt(req, "num")
	if ok {
		t.Errorf("expected parsing failure")
	}
}

func TestQueryBool(t *testing.T) {
	req, _ := http.NewRequest("GET", "/?flag=true", nil)
	b, ok := request.QueryBool(req, "flag")
	if !ok || !b {
		t.Errorf("expected true, got %v (%v)", b, ok)
	}
	req, _ = http.NewRequest("GET", "/?flag=notabool", nil)
	_, ok = request.QueryBool(req, "flag")
	if ok {
		t.Errorf("expected parsing failure")
	}
}

func TestParseFormAndFormValue(t *testing.T) {
	form := url.Values{}
	form.Set("alpha", "beta")
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	err := request.ParseForm(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	val, ok := request.FormValue(req, "alpha")
	if !ok || val != "beta" {
		t.Errorf("expected beta, got %s (%v)", val, ok)
	}
}

func TestParseQuery(t *testing.T) {
	vals, err := request.ParseQuery("x=1&y=2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vals.Get("x") != "1" || vals.Get("y") != "2" {
		t.Errorf("unexpected values: %v", vals)
	}
}
