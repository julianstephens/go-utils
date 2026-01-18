# CLI Utilities Package

The `cliutil` package provides helpers and utilities for building command-line interfaces (CLI) in Go projects. It offers reusable functions for common CLI patterns such as argument parsing, flag handling, interactive prompts, progress indicators, and colored output.

## Features

- **Argument Parsing**: Parse command-line arguments and flags
- **Colored Output**: Success, error, warning, and info formatting
- **Progress Indicators**: Progress bars and spinners
- **Interactive Prompts**: User input with validation
- **Table Output**: Formatted table display
- **Flag Utilities**: Convenient flag handling
- **Email Validation**: Built-in email validation

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

### Colored Output

```go
package main

import "github.com/julianstephens/go-utils/cliutil"

func main() {
    cliutil.PrintSuccess("Success!")
    cliutil.PrintError("Error occurred")
    cliutil.PrintWarning("Warning message")
    cliutil.PrintInfo("Information")
    cliutil.PrintColored("Custom message", cliutil.ColorBlue)
}
```

### Table Output

```go
package main

import "github.com/julianstephens/go-utils/cliutil"

func main() {
    data := [][]string{
        {"Name", "Age", "Role"},
        {"Alice", "30", "Engineer"},
        {"Bob", "28", "Designer"},
    }
    cliutil.PrintTable(data)
}
```

### Progress Bar

```go
package main

import (
    "time"
    "github.com/julianstephens/go-utils/cliutil"
)

func main() {
    pb := cliutil.NewProgressBarWithOptions(100, 40, "Processing")
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
    "time"
    "github.com/julianstephens/go-utils/cliutil"
)

func main() {
    spinner := cliutil.NewSpinner("Loading...")
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

## API Reference

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
- `HasFlag(args []string, flag string) bool` - Check if flag exists
- `GetFlagValue(args []string, flag, defaultValue string) string` - Get flag value

### Output Functions
- `PrintSuccess(message string)` - Print success (green)
- `PrintError(message string)` - Print error (red)
- `PrintWarning(message string)` - Print warning (yellow)
- `PrintInfo(message string)` - Print info (blue)
- `PrintColored(message string, color Color)` - Print with color
- `PrintTable(data [][]string)` - Print table

### Progress
- `NewProgressBar(total int) *ProgressBar` - Create progress bar
- `NewProgressBarWithOptions(total, width int, message string) *ProgressBar` - With options
- `NewSpinner(message string) *Spinner` - Create spinner

### Interactive Input
- `PromptString(message string) string` - String input
- `PromptBool(message string) bool` - Yes/no input
- `PromptStringWithValidation(message string, validator func(string) error) string` - With validation
- `PromptChoice(message string, options []string) int` - Choice selection
- `PromptPassword(prompt string) string` - Secure password input
- `PromptPasswordWithValidation(prompt string, validator func(string) error) string` - Password with validation

### Validation
- `ValidateEmail(email string) error` - Validate email format

## Colors

Available colors: `ColorRed`, `ColorGreen`, `ColorYellow`, `ColorBlue`, `ColorMagenta`, `ColorCyan`, `ColorWhite`, `ColorReset`

## Thread Safety

Progress indicators and interactive prompts are designed to be thread-safe for concurrent CLI operations. However, colored output functions should be used from a single goroutine to prevent mixed output.

## Integration

Works well with other go-utils packages:
- **logger**: Set log level from CLI flags
- **config**: Load config file from CLI arguments
- **validator**: Validate user input from prompts