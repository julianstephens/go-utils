package cliutil

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPromptStringWithIO(t *testing.T) {
	in := bytes.NewBufferString("hello world\n")
	out := &bytes.Buffer{}
	res := PromptStringWithIO("Enter: ", in, out)
	if res != "hello world" {
		t.Fatalf("unexpected result: %q", res)
	}
	if out.String() != "Enter: " {
		t.Fatalf("prompt not written to out: %q", out.String())
	}
}

func TestPromptPasswordWithIO_NonTTY(t *testing.T) {
	// simulate non-tty by using a buffer as in
	in := bytes.NewBufferString("s3cr3t\n")
	out := &bytes.Buffer{}
	res := PromptPasswordWithIO("Password: ", in, out)
	if res != "s3cr3t" {
		t.Fatalf("unexpected password: %q", res)
	}
	if out.String() != "Password: " {
		t.Fatalf("prompt not written: %q", out.String())
	}
}

func TestPromptPasswordWithValidationIO_Confirmation(t *testing.T) {
	// Simulate entering password and confirmation sequence. Use separate
	// readers for the initial entry and the confirmation to avoid conflicts
	// in the test harness.
	in := bytes.NewBufferString("hunter2\n")
	confirmIn := bytes.NewBufferString("hunter2\n")
	out := &bytes.Buffer{}

	validator := func(p string) error {
		// Use a separate reader to simulate a confirmation prompt.
		confirm := PromptPasswordWithIO("Confirm: ", confirmIn, out)
		if p != confirm {
			return fmt.Errorf("mismatch")
		}
		return nil
	}

	res := PromptPasswordWithValidationIO("Enter: ", in, out, validator)
	if res != "hunter2" {
		t.Fatalf("unexpected result: %q", res)
	}
}
