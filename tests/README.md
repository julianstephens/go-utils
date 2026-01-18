# tests (package)

This package contains shared test helper functions used across the repository. The helpers reduce test boilerplate and provide clearer assertions.

Import
------

Use an import alias (commonly `tst`) in your test files:

```go
import (
    "testing"
    tst "github.com/julianstephens/go-utils/tests"
)
```

Common helpers
--------------

- AssertDeepEqual / RequireDeepEqual: compare values using reflect.DeepEqual
- AssertNoError / RequireNoError: fail when an error is non-nil
- AssertTrue / AssertFalse: boolean assertions with messages
- AssertNotNil / AssertNil: nil checks that handle typed nils
- AssertJSONEquals: compare JSON payloads
- AssertErrorContains / AssertErrorIs: check error messages or wrap matches
- AssertWithinDuration: compare times with a tolerance
- AssertCloseTo: numeric closeness for floats
- AssertGreaterThan / AssertLessThan: compare ordered types
- AssertGreaterThanOrEqual / AssertLessThanOrEqual: compare ordered types with equality
- AssertEqual: compare equatable types (supports booleans)
- AssertPanics: assert that a function panics

All assertion functions support optional custom error messages:

```go
tst.AssertEqual(t, got, want)                           // without message
tst.AssertEqual(t, got, want, "custom error message")   // with message
```

HTTP helpers (in http_helpers.go)
---------------------------------

- AssertStatus: assert response status code
- AssertBodyContains / AssertBodyEquals: check response body
- NewRequestAndRecorder: convenience for creating *http.Request and *httptest.ResponseRecorder
- AssertHeaderEquals: check response header value

Misc
----

- Print: convenience for printing debug messages during tests
