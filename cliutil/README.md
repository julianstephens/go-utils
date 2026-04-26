# CLI Utilities Package

The `cliutil` package provides helpers and utilities for building command-line interfaces (CLI) in Go projects. It offers reusable functions for common CLI patterns such as argument parsing, flag handling, interactive prompts, progress indicators, and formatted output.

Output formatting is driven by a `Formatter` interface and a `Console` type, making it easy to switch between colored, plain-text, JSON, or accessibility-friendly output — including full support for the `NO_COLOR` environment variable convention.

## Features

- **Formatter Interface**: Pluggable output formatting (colored, plain-text, JSON, accessible)
- **Console Type**: Combines a formatter with an `io.Writer`; fully injectable and testable
- **Argument Parsing**: Parse command-line arguments and flags
- **Progress Indicators**: Progress bars and spinners driven by the active formatter
- **Interactive Prompts**: User input with validation
- **Table Output**: Formatted table display
- **Flag Utilities**: Convenient flag handling
- **NO_COLOR Support**: Automatically respects the `NO_COLOR` environment variable

## Installation

```bash
go get github.com/julianstephens/go-utils/cliutil
```

## Usage

### Argument Parsing

```go
package main

import (
    "github.com/julianstephens/go-utils/cliutil"
)

func main() {
    args := cliutil.ParseArgs([]string{"--config", "app.yaml", "--verbose", "input.txt"})
    
    config := args.GetFlagWithDefault("config", "default.yaml")
    verbose := args.HasFlag("verbose")
    _ = args.Positional
}
```

### Output Formatters

The package ships four built-in formatters. Choose one when creating a `Console`, or swap
the formatter on the package-level default with `SetDefaultConsole`.

| Formatter | Description |
|---|---|
| `NewColoredFormatter()` | ANSI colors + Unicode symbols (✓ ✗ ! ℹ) — default |
| `NewPlainTextFormatter()` | No ANSI, ASCII progress chars, plain prefixes (OK/ERROR/WARN/INFO) |
| `NewAccessibleFormatter()` | No ANSI, no symbols, explicit prefixes ([SUCCESS]/[ERROR]/…) |
| `NewJSONFormatter()` | Newline-delimited JSON objects; suppresses animated progress/spinner |

### Console — the recommended API

Create a `Console` from any `Formatter` and `io.Writer`. All output methods are
mutex-protected, so `ProgressBar` and `Spinner` goroutines are safe to share one `Console`.

```go
package main

import (
    "os"
    "github.com/julianstephens/go-utils/cliutil"
)

func main() {
    // Colored output to stdout (default)
    c := cliutil.NewConsole(cliutil.NewColoredFormatter(), os.Stdout)
    c.Success("Deployed!")
    c.Error("Connection refused")
    c.Warning("Retrying...")
    c.Info("Listening on :8080")
    c.PrintColored("custom", cliutil.ColorMagenta)

    // Table
    c.PrintTable([][]string{
        {"Name", "Age", "Role"},
        {"Alice", "30", "Engineer"},
        {"Bob", "28", "Designer"},
    })
}
```

### Switching Formatters

```go
package main

import (
    "os"
    "github.com/julianstephens/go-utils/cliutil"
)

func main() {
    args := cliutil.ParseArgs(os.Args[1:])

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

    cliutil.PrintSuccess("done") // uses the formatter chosen above
}
```

### NO_COLOR environment variable

The package-level default `Console` automatically uses `PlainTextFormatter` when the
`NO_COLOR` environment variable is set (any non-empty value), following the
[no-color.org](https://no-color.org) convention. No code changes required.

### Backwards-compatible module-level functions

The original `Print*` functions remain and delegate to the package-level `Console`.
Existing code continues to work without changes.

```go
cliutil.PrintSuccess("Success!")
cliutil.PrintError("Error occurred")
cliutil.PrintWarning("Warning message")
cliutil.PrintInfo("Information")
cliutil.PrintColored("Custom message", cliutil.ColorBlue)
cliutil.PrintTable(data)
```

### Table Output

```go
package main

import "github.com/julianstephens/go-utils/cliutil"

func main() {
    c := cliutil.DefaultConsole()
    c.PrintTable([][]string{
        {"Name", "Age", "Role"},
        {"Alice", "30", "Engineer"},
        {"Bob", "28", "Designer"},
    })
}
```

### Progress Bar

```go
package main

import (
    "os"
    "time"
    "github.com/julianstephens/go-utils/cliutil"
)

func main() {
    c := cliutil.NewConsole(cliutil.NewColoredFormatter(), os.Stdout)

    // Using default Console (writes to os.Stdout)
    pb := cliutil.NewProgressBarWithOptions(100, 40, "Processing")
    // Or inject a Console directly
    pb = cliutil.NewProgressBarWithConsole(c, 100, 40, "Processing")

    for i := 0; i <= 100; i++ {
        pb.Update(i)
        time.Sleep(50 * time.Millisecond)
    }
    pb.Finish()
}
```

### Spinner

```go
package main

import (
    "os"
    "time"
    "github.com/julianstephens/go-utils/cliutil"
)

func main() {
    c := cliutil.NewConsole(cliutil.NewColoredFormatter(), os.Stdout)

    // Using default Console
    spinner := cliutil.NewSpinner("Loading...")
    // Or inject a Console directly
    spinner = cliutil.NewSpinnerWithConsole(c, "Loading...")

    spinner.Start()
    time.Sleep(2 * time.Second)
    spinner.UpdateMessage("Finishing...")
    time.Sleep(1 * time.Second)
    spinner.Stop()
}
```

### Interactive Prompts

```go
package main

import (
    "github.com/julianstephens/go-utils/cliutil"
)

func main() {
    name := cliutil.PromptString("Enter your name: ")
    
    proceed := cliutil.PromptBool("Continue? (y/n): ")
    
    options := []string{"Create", "Open", "Exit"}
    choice := cliutil.PromptChoice("What to do?", options)
    _ = choice
}
```

### Secure Password Prompts

Prompt for sensitive input with secure echo-disabled input.

```go
package main

import (
    "fmt"
    "github.com/julianstephens/go-utils/cliutil"
)

func main() {
    pass := cliutil.PromptPassword("Enter passphrase: ")
    fmt.Printf("Received %d characters\n", len(pass))

    validator := func(p string) error {
        if len(p) < 8 {
            return fmt.Errorf("must be at least 8 characters")
        }
        return nil
    }
    confirmed := cliutil.PromptPasswordWithValidation("Choose passphrase: ", validator)
    _ = confirmed
}
```

### Validation

```go
package main

import (
    "github.com/julianstephens/go-utils/cliutil"
    "github.com/julianstephens/go-utils/validator"
)

func main() {
    emails := []string{"valid@example.com", "invalid-email"}
    
    for _, email := range emails {
        if err := validator.ValidateEmail(email); err != nil {
            cliutil.PrintError("Invalid email: " + email)
        } else {
            cliutil.PrintSuccess("Valid email: " + email)
        }
    }
}
```

### Flag Utilities

```go
package main

import (
    "github.com/julianstephens/go-utils/cliutil"
)

func main() {
    args := []string{"myapp", "--debug", "--port=8080", "--config", "prod.yaml"}
    
    if cliutil.HasFlag(args, "--debug") {
        cliutil.PrintInfo("Debug enabled")
    }
    
    port := cliutil.GetFlagValue(args, "--port", "3000")
    _ = port
}
```

### Testing with a Console

Inject a `Console` backed by a `bytes.Buffer` to capture and assert on output in tests:

```go
func TestMyCommand(t *testing.T) {
    var buf bytes.Buffer
    c := cliutil.NewConsole(cliutil.NewPlainTextFormatter(), &buf)

    myCommand(c) // pass Console as a dependency

    if !strings.Contains(buf.String(), "expected output") {
        t.Fatalf("unexpected output: %q", buf.String())
    }
}
```

### Custom Formatter

Implement the `Formatter` interface to define your own output style:

```go
type MyFormatter struct{}

func (f *MyFormatter) Success(text string) string  { return "[+] " + text }
func (f *MyFormatter) Error(text string) string    { return "[-] " + text }
func (f *MyFormatter) Warning(text string) string  { return "[~] " + text }
func (f *MyFormatter) Info(text string) string     { return "[i] " + text }
func (f *MyFormatter) Colored(text string, _ cliutil.Color) string { return text }
func (f *MyFormatter) Table(data [][]string) string { /* ... */ return "" }
func (f *MyFormatter) ProgressCharFilled() string  { return ">" }
func (f *MyFormatter) ProgressCharEmpty() string   { return "." }
func (f *MyFormatter) SpinnerFrames() []string     { return []string{"-", "\\", "|", "/"} }

c := cliutil.NewConsole(&MyFormatter{}, os.Stdout)
```

## API Reference

### Formatter Interface

```go
type Formatter interface {
    Success(text string) string
    Error(text string) string
    Warning(text string) string
    Info(text string) string
    Colored(text string, color Color) string
    Table(data [][]string) string
    ProgressCharFilled() string
    ProgressCharEmpty() string
    SpinnerFrames() []string
}
```

### Console

- `NewConsole(formatter Formatter, out io.Writer) *Console` - Create a Console
- `DefaultConsole() *Console` - Package-level Console (respects `NO_COLOR`)
- `SetDefaultConsole(c *Console)` - Replace the package-level Console
- `(c *Console) Success(message string) error`
- `(c *Console) Error(message string) error`
- `(c *Console) Warning(message string) error`
- `(c *Console) Info(message string) error`
- `(c *Console) PrintColored(text string, color Color) error`
- `(c *Console) PrintTable(data [][]string) error`
- `(c *Console) Formatter() Formatter`
- `(c *Console) SetFormatter(f Formatter)`

### Built-in Formatters

- `NewColoredFormatter() *ColoredFormatter` - ANSI colors + Unicode symbols
- `NewPlainTextFormatter() *PlainTextFormatter` - No colors, ASCII chars
- `NewAccessibleFormatter() *AccessibleFormatter` - Screen-reader friendly
- `NewJSONFormatter() *JSONFormatter` - Newline-delimited JSON

### Args Type

```go
type Args struct {
    Flags      map[string]string
    BoolFlags  map[string]bool
    Positional []string
}
```

### Argument Parsing

- `ParseArgs(args []string) *Args` - Parse command-line arguments
- `(a *Args) HasFlag(name string) bool` - Check for a boolean flag
- `(a *Args) GetFlag(name string) string` - Get a named flag value
- `(a *Args) GetFlagWithDefault(name, defaultValue string) string` - Get a named flag with fallback
- `HasFlag(args []string, flag string) bool` - Check slice for a flag
- `GetFlagValue(args []string, flag, defaultValue string) string` - Get flag value from slice

### Module-level Output (delegates to DefaultConsole)

- `PrintSuccess(message string)`
- `PrintError(message string)`
- `PrintWarning(message string)`
- `PrintInfo(message string)`
- `PrintColored(message string, color Color)`
- `PrintTable(data [][]string)`

### Progress

- `NewProgressBar(total int) *ProgressBar` - Progress bar via default Console
- `NewProgressBarWithOptions(total, width int, prefix string) *ProgressBar` - With options
- `NewProgressBarWithConsole(c *Console, total, width int, prefix string) *ProgressBar` - Inject Console
- `NewSpinner(message string) *Spinner` - Spinner via default Console
- `NewSpinnerWithConsole(c *Console, message string) *Spinner` - Inject Console

### Interactive Input

- `PromptString(prompt string) string` - String input
- `PromptStringWithIO(prompt string, in io.Reader, out io.Writer) string` - With custom IO
- `PromptStringWithValidation(prompt string, validator func(string) error) string` - With validation
- `PromptBool(prompt string) bool` - Yes/no input
- `PromptChoice(prompt string, options []string) int` - Choice selection (0-based index)
- `PromptPassword(prompt string) string` - Secure password input (echo disabled)
- `PromptPasswordWithValidation(prompt string, validator func(string) error) string` - Password with validation
- `PromptPasswordWithIO(prompt string, in io.Reader, out io.Writer) string` - With custom IO
- `PromptPasswordWithValidationIO(prompt string, in io.Reader, out io.Writer, validator func(string) error) string`

## Colors

Available color constants: `ColorRed`, `ColorGreen`, `ColorYellow`, `ColorBlue`, `ColorMagenta`, `ColorCyan`, `ColorWhite`, `ColorBold`, `ColorReset`

Colors are only applied by formatters that support ANSI output (`ColoredFormatter`). Other formatters ignore the `color` parameter in `Colored()`.

## Thread Safety

All `Console` output methods are mutex-protected. `ProgressBar` and `Spinner`, which
write from background goroutines, safely share the same `Console`. The module-level
`Print*` functions are safe to call concurrently.

## Integration

Works well with other go-utils packages:

- **logger**: Set log level from CLI flags
- **config**: Load config file from CLI arguments
- **validator**: Validate user input from prompts
