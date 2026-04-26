package cliutil

// Formatter is the interface for all output formatters. Each method accepts
// plain text and returns a formatted string. Formatters are pure: they do not
// perform IO. The Console type is responsible for writing to an io.Writer.
type Formatter interface {
	// Semantic output methods
	Success(text string) string
	Error(text string) string
	Warning(text string) string
	Info(text string) string
	Colored(text string, color Color) string

	// Structured data
	Table(data [][]string) string

	// Visual element characters used by ProgressBar and Spinner.
	// Formatters that do not support terminal animations (e.g. JSONFormatter)
	// should return empty strings to suppress animated output.
	ProgressCharFilled() string
	ProgressCharEmpty() string
	SpinnerFrames() []string
}
