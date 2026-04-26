/*
Package cliutil provides helpers and utilities for building command-line interfaces (CLI) in Go projects.

Output formatting is driven by a Formatter interface and a Console type. A Console combines a
Formatter (which controls how text is styled) with an io.Writer (which controls where output is
sent). All writes are mutex-protected, so ProgressBar and Spinner goroutines are safe to share
one Console.

# Formatters

Four built-in formatters are provided:

  - ColoredFormatter  — ANSI colors and Unicode symbols (✓ ✗ ! ℹ); the default.
  - PlainTextFormatter — No ANSI codes, ASCII progress chars, plain prefixes (OK/ERROR/WARN/INFO).
  - AccessibleFormatter — No ANSI, no symbols, explicit prefixes ([SUCCESS]/[ERROR]/…).
  - JSONFormatter — Newline-delimited JSON objects; suppresses animated progress and spinner output.

The package-level default Console automatically uses PlainTextFormatter when the NO_COLOR
environment variable is set (any non-empty value), following the no-color.org convention.

# Basic Usage

	package main

	import (
		"os"

		"github.com/julianstephens/go-utils/cliutil"
	)

	func main() {
		// Parse command-line arguments
		args := cliutil.ParseArgs(os.Args[1:])

		// Choose a formatter based on flags
		var f cliutil.Formatter
		switch {
		case args.HasFlag("json"):
			f = cliutil.NewJSONFormatter()
		case args.HasFlag("no-color"):
			f = cliutil.NewPlainTextFormatter()
		default:
			f = cliutil.NewColoredFormatter()
		}

		// Replace the package-level default so all Print* calls use it
		cliutil.SetDefaultConsole(cliutil.NewConsole(f, os.Stdout))

		// Semantic output via Console
		c := cliutil.DefaultConsole()
		c.Success("Deployed successfully")
		c.Error("Connection refused")
		c.Warning("Retrying in 5s")
		c.Info("Listening on :8080")

		// Or use the module-level wrappers (delegate to DefaultConsole)
		cliutil.PrintSuccess("All done")
	}

# Console

Create a Console from any Formatter and io.Writer:

	c := cliutil.NewConsole(cliutil.NewColoredFormatter(), os.Stdout)
	c.Success("ok")
	c.PrintTable([][]string{
		{"Name", "Age"},
		{"Alice", "30"},
	})

Swap the formatter at runtime (e.g. after parsing --no-color):

	c.SetFormatter(cliutil.NewPlainTextFormatter())

# Testing

Inject a Console backed by a bytes.Buffer to capture and assert on output:

	func TestMyCommand(t *testing.T) {
		var buf bytes.Buffer
		c := cliutil.NewConsole(cliutil.NewPlainTextFormatter(), &buf)

		myCommand(c)

		if !strings.Contains(buf.String(), "expected output") {
			t.Fatalf("unexpected output: %q", buf.String())
		}
	}

# Argument Parsing

	args := cliutil.ParseArgs([]string{"--verbose", "--output", "file.txt", "input.txt"})
	verbose := args.HasFlag("verbose")
	output  := args.GetFlag("output")
	files   := args.Positional

# Flag Handling

	if cliutil.HasFlag(os.Args, "--debug") {
		// Enable debug mode
	}
	logLevel := cliutil.GetFlagValue(os.Args, "--log-level", "info")

# Interactive Prompts

	name    := cliutil.PromptString("Enter your name: ")
	proceed := cliutil.PromptBool("Continue? (y/n): ")
	option  := cliutil.PromptChoice("Select:", []string{"Create", "Open", "Exit"})

	// Secure password input (echo disabled on TTY)
	pass := cliutil.PromptPassword("Passphrase: ")
	pass  = cliutil.PromptPasswordWithValidation("Passphrase: ", func(p string) error {
		if len(p) < 8 {
			return fmt.Errorf("must be at least 8 characters")
		}
		return nil
	})

# Progress Indicators

	// Progress bar — inject a Console for full control
	c  := cliutil.NewConsole(cliutil.NewColoredFormatter(), os.Stdout)
	pb := cliutil.NewProgressBarWithConsole(c, total, 40, "Processing")
	for i := 0; i <= total; i++ {
		pb.Update(i)
	}
	pb.Finish()

	// Spinner — inject a Console
	s := cliutil.NewSpinnerWithConsole(c, "Loading...")
	s.Start()
	defer s.Stop()

# Custom Formatter

Implement the Formatter interface to define your own output style:

	type MyFormatter struct{}

	func (f *MyFormatter) Success(text string) string                   { return "[+] " + text }
	func (f *MyFormatter) Error(text string) string                     { return "[-] " + text }
	func (f *MyFormatter) Warning(text string) string                   { return "[~] " + text }
	func (f *MyFormatter) Info(text string) string                      { return "[i] " + text }
	func (f *MyFormatter) Colored(text string, _ cliutil.Color) string  { return text }
	func (f *MyFormatter) Table(data [][]string) string                 { return "" }
	func (f *MyFormatter) ProgressCharFilled() string                   { return ">" }
	func (f *MyFormatter) ProgressCharEmpty() string                    { return "." }
	func (f *MyFormatter) SpinnerFrames() []string                      { return []string{"-", "\\", "|", "/"} }

	c := cliutil.NewConsole(&MyFormatter{}, os.Stdout)

# Thread Safety

All Console output methods are mutex-protected. ProgressBar and Spinner, which write from
background goroutines, safely share the same Console. The module-level Print* functions
are safe to call concurrently via the package-level default Console.
*/
package cliutil
