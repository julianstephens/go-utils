package request_test

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"

	"github.com/julianstephens/go-utils/httputil/request"
	tst "github.com/julianstephens/go-utils/tests"
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
	tst.AssertNoError(t, err)
	tst.AssertTrue(t, ts.Name == "Alice" && ts.Age == 30, "decoded struct should match")
}

func TestDecodeJSON_InvalidContentType(t *testing.T) {
	body := []byte(`{"name":"Alice","age":30}`)
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "text/plain")

	var ts testStruct
	err := request.DecodeJSON(req, &ts)
	tst.AssertTrue(t, err == request.ErrInvalidContentType, "should return ErrInvalidContentType")
}

func TestQueryValue(t *testing.T) {
	req, _ := http.NewRequest("GET", "/?foo=bar", nil)
	val, ok := request.QueryValue(req, "foo")
	tst.AssertTrue(t, ok && val == "bar", "QueryValue should return bar")
	_, ok = request.QueryValue(req, "baz")
	tst.AssertFalse(t, ok, "QueryValue should be missing for baz")
}

func TestQueryInt(t *testing.T) {
	req, _ := http.NewRequest("GET", "/?num=42", nil)
	i, ok := request.QueryInt(req, "num")
	tst.AssertTrue(t, ok && i == 42, "QueryInt should return 42")
	_, ok = request.QueryInt(req, "bad")
	tst.AssertFalse(t, ok, "QueryInt should be missing for bad")
	req, _ = http.NewRequest("GET", "/?num=foo", nil)
	_, ok = request.QueryInt(req, "num")
	tst.AssertFalse(t, ok, "QueryInt should fail parsing for non-int")
}

func TestQueryBool(t *testing.T) {
	req, _ := http.NewRequest("GET", "/?flag=true", nil)
	b, ok := request.QueryBool(req, "flag")
	tst.AssertTrue(t, ok && b, "QueryBool should return true")
	req, _ = http.NewRequest("GET", "/?flag=notabool", nil)
	_, ok = request.QueryBool(req, "flag")
	tst.AssertFalse(t, ok, "QueryBool should fail parsing for invalid value")
}

func TestParseFormAndFormValue(t *testing.T) {
	form := url.Values{}
	form.Set("alpha", "beta")
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	err := request.ParseForm(req)
	tst.AssertNoError(t, err)
	val, ok := request.FormValue(req, "alpha")
	tst.AssertTrue(t, ok && val == "beta", "FormValue should return beta")
}

func TestParseQuery(t *testing.T) {
	vals, err := request.ParseQuery("x=1&y=2")
	tst.AssertNoError(t, err)
	tst.AssertTrue(t, vals.Get("x") == "1" && vals.Get("y") == "2", "ParseQuery values should match")
}
