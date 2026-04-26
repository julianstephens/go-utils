package cliutil

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ColoredFormatter formats output with ANSI color codes and Unicode symbols.
// This mirrors the original hardcoded behavior of the package.
type ColoredFormatter struct{}

// NewColoredFormatter creates a formatter that applies ANSI colors and Unicode
// symbols (✓, ✗, !, ℹ) to output.
func NewColoredFormatter() *ColoredFormatter {
	return &ColoredFormatter{}
}

func (f *ColoredFormatter) Success(text string) string {
	return fmt.Sprintf("%s✓ %s%s", ColorGreen, text, ColorReset)
}

func (f *ColoredFormatter) Error(text string) string {
	return fmt.Sprintf("%s✗ %s%s", ColorRed, text, ColorReset)
}

func (f *ColoredFormatter) Warning(text string) string {
	return fmt.Sprintf("%s! %s%s", ColorYellow, text, ColorReset)
}

func (f *ColoredFormatter) Info(text string) string {
	return fmt.Sprintf("%sℹ %s%s", ColorBlue, text, ColorReset)
}

func (f *ColoredFormatter) Colored(text string, color Color) string {
	return fmt.Sprintf("%s%s%s", color, text, ColorReset)
}

func (f *ColoredFormatter) Table(data [][]string) string {
	return formatTable(data)
}

func (f *ColoredFormatter) ProgressCharFilled() string { return "█" }
func (f *ColoredFormatter) ProgressCharEmpty() string  { return "░" }
func (f *ColoredFormatter) SpinnerFrames() []string {
	return []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
}

// PlainTextFormatter formats output as plain text with no ANSI codes.
// Useful for piped output, log files, or terminals that don't support colors.
type PlainTextFormatter struct{}

// NewPlainTextFormatter creates a formatter that produces no ANSI codes.
func NewPlainTextFormatter() *PlainTextFormatter {
	return &PlainTextFormatter{}
}

func (f *PlainTextFormatter) Success(text string) string { return "OK: " + text }
func (f *PlainTextFormatter) Error(text string) string   { return "ERROR: " + text }
func (f *PlainTextFormatter) Warning(text string) string { return "WARN: " + text }
func (f *PlainTextFormatter) Info(text string) string    { return "INFO: " + text }

func (f *PlainTextFormatter) Colored(text string, _ Color) string { return text }

func (f *PlainTextFormatter) Table(data [][]string) string { return formatTable(data) }

func (f *PlainTextFormatter) ProgressCharFilled() string { return "#" }
func (f *PlainTextFormatter) ProgressCharEmpty() string  { return "-" }
func (f *PlainTextFormatter) SpinnerFrames() []string {
	return []string{"|", "/", "-", "\\"}
}

// AccessibleFormatter formats output for screen readers and accessibility
// tools. It uses clear, explicit text prefixes with no decorative symbols.
type AccessibleFormatter struct{}

// NewAccessibleFormatter creates a formatter optimized for accessibility.
func NewAccessibleFormatter() *AccessibleFormatter {
	return &AccessibleFormatter{}
}

func (f *AccessibleFormatter) Success(text string) string { return "[SUCCESS] " + text }
func (f *AccessibleFormatter) Error(text string) string   { return "[ERROR] " + text }
func (f *AccessibleFormatter) Warning(text string) string { return "[WARNING] " + text }
func (f *AccessibleFormatter) Info(text string) string    { return "[INFO] " + text }

func (f *AccessibleFormatter) Colored(text string, _ Color) string { return text }

func (f *AccessibleFormatter) Table(data [][]string) string { return formatTable(data) }

func (f *AccessibleFormatter) ProgressCharFilled() string { return "#" }
func (f *AccessibleFormatter) ProgressCharEmpty() string  { return "." }
func (f *AccessibleFormatter) SpinnerFrames() []string {
	return []string{"|", "/", "-", "\\"}
}

// jsonMessage is the JSON payload emitted by JSONFormatter for each output call.
type jsonMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// JSONFormatter formats output as newline-delimited JSON objects.
// Useful for machine-readable output, CI pipelines, and log aggregation.
type JSONFormatter struct{}

// NewJSONFormatter creates a formatter that emits newline-delimited JSON.
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

func (f *JSONFormatter) marshal(msgType, text string) string {
	b, _ := json.Marshal(jsonMessage{Type: msgType, Message: text})
	return string(b)
}

func (f *JSONFormatter) Success(text string) string { return f.marshal("success", text) }
func (f *JSONFormatter) Error(text string) string   { return f.marshal("error", text) }
func (f *JSONFormatter) Warning(text string) string { return f.marshal("warning", text) }
func (f *JSONFormatter) Info(text string) string    { return f.marshal("info", text) }

func (f *JSONFormatter) Colored(text string, _ Color) string {
	return f.marshal("colored", text)
}

func (f *JSONFormatter) Table(data [][]string) string {
	if len(data) == 0 {
		return ""
	}
	rows := make([]map[string]string, 0, len(data)-1)
	if len(data) < 2 {
		b, _ := json.Marshal(rows)
		return string(b)
	}
	headers := data[0]
	for _, row := range data[1:] {
		m := make(map[string]string, len(headers))
		for i, h := range headers {
			if i < len(row) {
				m[h] = row[i]
			}
		}
		rows = append(rows, m)
	}
	b, _ := json.Marshal(rows)
	return string(b)
}

// JSONFormatter suppresses terminal animations.
func (f *JSONFormatter) ProgressCharFilled() string { return "" }
func (f *JSONFormatter) ProgressCharEmpty() string  { return "" }
func (f *JSONFormatter) SpinnerFrames() []string    { return nil }

// formatTable renders a 2-D string slice as a plain-text table with a header
// separator. It is shared by all formatters that produce textual tables.
func formatTable(data [][]string) string {
	if len(data) == 0 {
		return ""
	}

	colWidths := make([]int, len(data[0]))
	for _, row := range data {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	var sb strings.Builder
	for i, row := range data {
		for j, cell := range row {
			if j < len(colWidths) {
				fmt.Fprintf(&sb, "%-*s", colWidths[j]+2, cell)
			}
		}
		sb.WriteByte('\n')
		if i == 0 {
			for j := range row {
				if j < len(colWidths) {
					sb.WriteString(strings.Repeat("-", colWidths[j]+2))
				}
			}
			sb.WriteByte('\n')
		}
	}
	return strings.TrimSuffix(sb.String(), "\n")
}
