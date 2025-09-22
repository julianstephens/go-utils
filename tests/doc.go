// Package tests provides shared test helpers used across the repository.
//
// The helpers are convenience assertions and small utilities used to
// reduce boilerplate in unit tests. Import them in your test files
// like this:
//
//	import (
//	    "testing"
//	    tst "github.com/julianstephens/go-utils/tests"
//	)
//
// Example usage:
//
//	func TestSomething(t *testing.T) {
//	    got := DoWork()
//	    want := SomeValue()
//	    tst.RequireDeepEqual(t, got, want)
//	}
package tests
