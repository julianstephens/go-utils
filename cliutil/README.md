# CLI Utilities Package

The `cliutil` package provides helpers and utilities for building command-line interfaces (CLI) in Go projects. It offers reusable functions for common CLI patterns such as argument parsing, flag handling, interactive prompts, progress indicators, and colored output.

## Features

- **Argument Parsing**: Structured parsing of command-line arguments and flags
- **Colored Output**: Success, error, warning, and info message formatting
- **Progress Indicators**: Progress bars and spinners for long-running operations
- **Interactive Prompts**: User input with validation
- **Table Output**: Formatted table display
- **Flag Utilities**: Convenient flag handling functions
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
    "fmt"
    "github.com/julianstephens/go-utils/cliutil"
)

func main() {
    testArgs := []string{"--config", "app.yaml", "--verbose", "--output", "result.txt", "input1.txt", "input2.txt"}
    args := cliutil.ParseArgs(testArgs)
    
    fmt.Printf("Config file: %s\n", args.GetFlagWithDefault("config", "default.yaml"))
    fmt.Printf("Verbose mode: %t\n", args.HasFlag("verbose"))
    fmt.Printf("Output file: %s\n", args.GetFlag("output"))
    fmt.Printf("Input files: %v\n", args.Positional)
}
```

### Colored Output

```go
package main

import "github.com/julianstephens/go-utils/cliutil"

func main() {
    cliutil.PrintSuccess("Operation completed successfully!")
    cliutil.PrintError("Something went wrong")
    cliutil.PrintWarning("This is a warning message")
    cliutil.PrintInfo("Here's some information")
    cliutil.PrintColored("Custom blue message", cliutil.ColorBlue)
}
```

### Table Output

```go
package main

import "github.com/julianstephens/go-utils/cliutil"

func main() {
    tableData := [][]string{
        {"Name", "Age", "Role", "Department"},
        {"Alice", "30", "Engineer", "Development"},
        {"Bob", "28", "Designer", "UX/UI"},
        {"Charlie", "35", "Manager", "Operations"},
    }
    cliutil.PrintTable(tableData)
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
    total := 100
    pb := cliutil.NewProgressBarWithOptions(total, 40, "Processing files")
    
    for i := 0; i <= total; i++ {
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
    
    // Simulate work
    time.Sleep(2 * time.Second)
    spinner.UpdateMessage("Connecting to database...")
    time.Sleep(1 * time.Second)
    spinner.UpdateMessage("Finalizing...")
    time.Sleep(1 * time.Second)
    
    spinner.Stop()
    cliutil.PrintSuccess("Operation completed!")
}
```

### Interactive Prompts

```go
package main

import (
    "fmt"
    "github.com/julianstephens/go-utils/cliutil"
)

func main() {
    // Simple string prompt
    name := cliutil.PromptString("Enter your name: ")
    fmt.Printf("Hello, %s!\n", name)
    
    // Boolean prompt
    proceed := cliutil.PromptBool("Do you want to continue? (y/n): ")
    if proceed {
        cliutil.PrintSuccess("Continuing...")
    } else {
        cliutil.PrintWarning("Operation cancelled")
    }
    
    // String prompt with validation
    email := cliutil.PromptStringWithValidation("Enter your email: ", cliutil.ValidateEmail)
    cliutil.PrintSuccess(fmt.Sprintf("Email %s is valid!", email))
    
    // Choice prompt
    options := []string{"Create new project", "Open existing project", "Exit"}
    choice := cliutil.PromptChoice("What would you like to do?", options)
    fmt.Printf("You selected: %s\n", options[choice])
}
```

### Validation

```go
package main

import (
    "fmt"
    "github.com/julianstephens/go-utils/cliutil"
)

func main() {
    testEmails := []string{"valid@example.com", "invalid-email", "@example.com", "user@domain.org"}
    
    for _, email := range testEmails {
        if err := cliutil.ValidateEmail(email); err != nil {
            cliutil.PrintError(fmt.Sprintf("Invalid email '%s': %v", email, err))
        } else {
            cliutil.PrintSuccess(fmt.Sprintf("Valid email: %s", email))
        }
    }
}
```

### Flag Utilities

```go
package main

import (
    "fmt"
    "github.com/julianstephens/go-utils/cliutil"
)

func main() {
    sampleArgs := []string{"myapp", "--debug", "--port=8080", "--config", "prod.yaml", "serve"}
    
    if cliutil.HasFlag(sampleArgs, "--debug") {
        cliutil.PrintInfo("Debug mode is enabled")
    }
    
    port := cliutil.GetFlagValue(sampleArgs, "--port", "3000")
    cliutil.PrintInfo(fmt.Sprintf("Server will run on port: %s", port))
    
    configFile := cliutil.GetFlagValue(sampleArgs, "--config", "config.yaml")
    cliutil.PrintInfo(fmt.Sprintf("Using configuration file: %s", configFile))
}
```

## API Reference

### Types

#### Args
Represents parsed command-line arguments:
```go
type Args struct {
    Flags      map[string]string // Named flags with values
    BoolFlags  map[string]bool   // Boolean flags
    Positional []string          // Positional arguments
}
```

Methods:
- `HasFlag(name string) bool` - Check if flag exists
- `GetFlag(name string) string` - Get flag value
- `GetFlagWithDefault(name, defaultValue string) string` - Get flag with default

### Functions

#### Argument Parsing
- `ParseArgs(args []string) *Args` - Parse command line arguments

#### Output Functions
- `PrintSuccess(message string)` - Print success message in green
- `PrintError(message string)` - Print error message in red
- `PrintWarning(message string)` - Print warning message in yellow
- `PrintInfo(message string)` - Print info message in blue
- `PrintColored(message string, color Color)` - Print message in specified color
- `PrintTable(data [][]string)` - Print formatted table

#### Progress Indicators
- `NewProgressBar(total int) *ProgressBar` - Create new progress bar
- `NewProgressBarWithOptions(total, width int, message string) *ProgressBar` - Create progress bar with options
- `NewSpinner(message string) *Spinner` - Create new spinner

#### Interactive Prompts
- `PromptString(message string) string` - Prompt for string input
- `PromptBool(message string) bool` - Prompt for yes/no input
- `PromptStringWithValidation(message string, validator func(string) error) string` - Prompt with validation
- `PromptChoice(message string, options []string) int` - Prompt for choice selection

#### Validation
- `ValidateEmail(email string) error` - Validate email format

#### Flag Utilities
- `HasFlag(args []string, flag string) bool` - Check if flag exists in args
- `GetFlagValue(args []string, flag, defaultValue string) string` - Get flag value with default

## Color Constants

Available colors for `PrintColored`:
- `ColorRed`
- `ColorGreen`
- `ColorYellow`
- `ColorBlue`
- `ColorMagenta`
- `ColorCyan`
- `ColorWhite`
- `ColorReset`

## Thread Safety

Progress indicators and interactive prompts are designed to be thread-safe for concurrent CLI operations. However, colored output functions should be used from a single goroutine to prevent mixed output.

## Integration

This package integrates well with other go-utils packages like logger and config:

```go
// Use with logger for structured CLI logging
if cliutil.HasFlag(os.Args, "--verbose") {
    logger.SetLogLevel("debug")
}

// Use with config for CLI configuration
var cfg AppConfig
configFile := cliutil.GetFlagValue(os.Args, "--config", "config.yaml")
config.LoadFromFile(configFile, &cfg)
```