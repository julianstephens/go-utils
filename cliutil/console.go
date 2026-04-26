package cliutil

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// Console is the primary output surface for a CLI application. It combines a
// Formatter (which controls how text is styled) with an io.Writer (which
// controls where output is sent). All writes are mutex-protected so that
// ProgressBar/Spinner goroutines can safely share the same Console.
type Console struct {
	formatter Formatter
	out       io.Writer
	mu        sync.Mutex
}

// NewConsole creates a Console that formats output using formatter and writes
// to out. Use DefaultConsole() to get a pre-configured Console that writes to
// os.Stdout and respects the NO_COLOR environment variable.
func NewConsole(formatter Formatter, out io.Writer) *Console {
	return &Console{formatter: formatter, out: out}
}

// Formatter returns the Formatter used by this Console.
func (c *Console) Formatter() Formatter {
	return c.formatter
}

// SetFormatter replaces the Formatter on this Console. It is safe to call
// concurrently.
func (c *Console) SetFormatter(f Formatter) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.formatter = f
}

// write is the single internal path through which all Console output flows.
func (c *Console) write(s string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, err := fmt.Fprintln(c.out, s)
	return err
}

// writeTo writes to a specific writer without holding the Console mutex.
// Used by ProgressBar/Spinner which need fine-grained control over \r.
func (c *Console) writeTo(w io.Writer, s string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, err := fmt.Fprint(w, s)
	return err
}

// Success writes a success-level message.
func (c *Console) Success(message string) error {
	return c.write(c.formatter.Success(message))
}

// Error writes an error-level message.
func (c *Console) Error(message string) error {
	return c.write(c.formatter.Error(message))
}

// Warning writes a warning-level message.
func (c *Console) Warning(message string) error {
	return c.write(c.formatter.Warning(message))
}

// Info writes an info-level message.
func (c *Console) Info(message string) error {
	return c.write(c.formatter.Info(message))
}

// PrintColored writes text styled with the given Color.
func (c *Console) PrintColored(text string, color Color) error {
	return c.write(c.formatter.Colored(text, color))
}

// PrintTable writes data as a formatted table.
func (c *Console) PrintTable(data [][]string) error {
	rendered := c.formatter.Table(data)
	if rendered == "" {
		return nil
	}
	return c.write(rendered)
}

// --- package-level default Console ---

var (
	defaultConsoleMu sync.Mutex
	defaultConsole   *Console
)

func init() {
	var f Formatter
	if os.Getenv("NO_COLOR") != "" {
		f = NewPlainTextFormatter()
	} else {
		f = NewColoredFormatter()
	}
	defaultConsole = NewConsole(f, os.Stdout)
}

// DefaultConsole returns the package-level Console. It writes to os.Stdout
// and automatically uses PlainTextFormatter when the NO_COLOR environment
// variable is set, or ColoredFormatter otherwise.
func DefaultConsole() *Console {
	defaultConsoleMu.Lock()
	defer defaultConsoleMu.Unlock()
	return defaultConsole
}

// SetDefaultConsole replaces the package-level default Console. This is
// useful in main() to configure a global formatter (e.g. JSON output) once
// so all subsequent module-level Print* calls use it.
func SetDefaultConsole(c *Console) {
	defaultConsoleMu.Lock()
	defer defaultConsoleMu.Unlock()
	defaultConsole = c
}
