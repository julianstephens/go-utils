package cliutil_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/julianstephens/go-utils/cliutil"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected *cliutil.Args
	}{
		{
			name: "mixed arguments",
			args: []string{"--verbose", "--output", "file.txt", "-c", "config.json", "input.txt", "output.txt"},
			expected: &cliutil.Args{
				Flags:      map[string]string{"output": "file.txt", "c": "config.json"},
				BoolFlags:  map[string]bool{"verbose": true},
				Positional: []string{"input.txt", "output.txt"},
			},
		},
		{
			name: "flag with equals",
			args: []string{"--config=config.yaml", "--debug"},
			expected: &cliutil.Args{
				Flags:      map[string]string{"config": "config.yaml"},
				BoolFlags:  map[string]bool{"debug": true},
				Positional: []string{},
			},
		},
		{
			name: "only positional",
			args: []string{"file1.txt", "file2.txt"},
			expected: &cliutil.Args{
				Flags:      map[string]string{},
				BoolFlags:  map[string]bool{},
				Positional: []string{"file1.txt", "file2.txt"},
			},
		},
		{
			name: "empty args",
			args: []string{},
			expected: &cliutil.Args{
				Flags:      map[string]string{},
				BoolFlags:  map[string]bool{},
				Positional: []string{},
			},
		},
		{
			name: "short flags only",
			args: []string{"-v", "-c", "output.txt", "-h"},
			expected: &cliutil.Args{
				Flags:      map[string]string{"c": "output.txt"},
				BoolFlags:  map[string]bool{"v": true, "h": true},
				Positional: []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cliutil.ParseArgs(tt.args)

			if !reflect.DeepEqual(result.Flags, tt.expected.Flags) {
				t.Errorf("Expected flags %v, got %v", tt.expected.Flags, result.Flags)
			}
			if !reflect.DeepEqual(result.BoolFlags, tt.expected.BoolFlags) {
				t.Errorf("Expected bool flags %v, got %v", tt.expected.BoolFlags, result.BoolFlags)
			}
			if !reflect.DeepEqual(result.Positional, tt.expected.Positional) {
				t.Errorf("Expected positional %v, got %v", tt.expected.Positional, result.Positional)
			}
		})
	}
}

func TestArgs_HasFlag(t *testing.T) {
	args := cliutil.ParseArgs([]string{"--verbose", "--debug"})

	if !args.HasFlag("verbose") {
		t.Error("Expected verbose flag to be true")
	}
	if !args.HasFlag("debug") {
		t.Error("Expected debug flag to be true")
	}
	if args.HasFlag("quiet") {
		t.Error("Expected quiet flag to be false")
	}
}

func TestArgs_GetFlag(t *testing.T) {
	args := cliutil.ParseArgs([]string{"--output", "file.txt", "--config=config.yaml"})

	if args.GetFlag("output") != "file.txt" {
		t.Errorf("Expected output flag to be 'file.txt', got '%s'", args.GetFlag("output"))
	}
	if args.GetFlag("config") != "config.yaml" {
		t.Errorf("Expected config flag to be 'config.yaml', got '%s'", args.GetFlag("config"))
	}
	if args.GetFlag("nonexistent") != "" {
		t.Errorf("Expected nonexistent flag to be empty, got '%s'", args.GetFlag("nonexistent"))
	}
}

func TestArgs_GetFlagWithDefault(t *testing.T) {
	args := cliutil.ParseArgs([]string{"--output", "file.txt"})

	if args.GetFlagWithDefault("output", "default.txt") != "file.txt" {
		t.Error("Expected to get actual flag value, not default")
	}
	if args.GetFlagWithDefault("nonexistent", "default.txt") != "default.txt" {
		t.Error("Expected to get default value for nonexistent flag")
	}
}

func TestHasFlag(t *testing.T) {
	args := []string{"--verbose", "--output=file.txt", "input.txt"}

	if !cliutil.HasFlag(args, "--verbose") {
		t.Error("Expected to find --verbose flag")
	}
	if !cliutil.HasFlag(args, "--output") {
		t.Error("Expected to find --output flag")
	}
	if cliutil.HasFlag(args, "--debug") {
		t.Error("Expected not to find --debug flag")
	}
}

func TestGetFlagValue(t *testing.T) {
	args := []string{"--output", "file.txt", "--config=config.yaml", "input.txt"}

	if cliutil.GetFlagValue(args, "--output", "default") != "file.txt" {
		t.Error("Expected to get flag value for --output")
	}
	if cliutil.GetFlagValue(args, "--config", "default") != "config.yaml" {
		t.Error("Expected to get flag value for --config")
	}
	if cliutil.GetFlagValue(args, "--nonexistent", "default") != "default" {
		t.Error("Expected to get default value for nonexistent flag")
	}
}

func TestValidateNonEmpty(t *testing.T) {
	tests := []struct {
		input     string
		shouldErr bool
	}{
		{"valid input", false},
		{"", true},
		{"   ", true},
		{"\t\n", true},
		{"a", false},
	}

	for _, tt := range tests {
		err := cliutil.ValidateNonEmpty(tt.input)
		if tt.shouldErr && err == nil {
			t.Errorf("Expected error for input '%s'", tt.input)
		}
		if !tt.shouldErr && err != nil {
			t.Errorf("Expected no error for input '%s', got %v", tt.input, err)
		}
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		input     string
		shouldErr bool
	}{
		{"test@example.com", false},
		{"user@domain.org", false},
		{"invalid-email", true},
		{"", true},
		{"@domain.com", true},
		{"user@", true},
		{"user.domain", true},
	}

	for _, tt := range tests {
		err := cliutil.ValidateEmail(tt.input)
		if tt.shouldErr && err == nil {
			t.Errorf("Expected error for email '%s'", tt.input)
		}
		if !tt.shouldErr && err != nil {
			t.Errorf("Expected no error for email '%s', got %v", tt.input, err)
		}
	}
}

func TestNewProgressBar(t *testing.T) {
	pb := cliutil.NewProgressBar(100)
	if pb == nil {
		t.Error("Expected progress bar to be created")
	}

	// Test update (we can't easily test the visual output, but we can test it doesn't panic)
	pb.Update(50)
	pb.Update(100)
	pb.Finish()
}

func TestNewProgressBarWithOptions(t *testing.T) {
	pb := cliutil.NewProgressBarWithOptions(200, 30, "Custom Progress")
	if pb == nil {
		t.Error("Expected progress bar to be created")
	}

	pb.Update(100)
	pb.Finish()
}

func TestNewSpinner(t *testing.T) {
	spinner := cliutil.NewSpinner("Loading...")
	if spinner == nil {
		t.Error("Expected spinner to be created")
	}

	// Test start and stop (we can't easily test the visual output, but we can test it doesn't panic)
	spinner.Start()
	time.Sleep(200 * time.Millisecond)
	spinner.UpdateMessage("Updated message")
	time.Sleep(200 * time.Millisecond)
	spinner.Stop()

	// Test double start/stop
	spinner.Start()
	spinner.Start() // Should not panic
	spinner.Stop()
	spinner.Stop() // Should not panic
}

// Test helper functions that don't require user input
func TestColorConstants(t *testing.T) {
	// Test that color constants are not empty
	colors := []cliutil.Color{
		cliutil.ColorReset,
		cliutil.ColorRed,
		cliutil.ColorGreen,
		cliutil.ColorYellow,
		cliutil.ColorBlue,
		cliutil.ColorMagenta,
		cliutil.ColorCyan,
		cliutil.ColorWhite,
		cliutil.ColorBold,
	}

	for i, color := range colors {
		if string(color) == "" {
			t.Errorf("Color constant %d should not be empty", i)
		}
	}
}

func TestPrintFunctions(t *testing.T) {
	// These functions write to stdout, so we can't easily test their output
	// But we can test that they don't panic
	cliutil.PrintColored("Test message", cliutil.ColorBlue)
	cliutil.PrintSuccess("Success message")
	cliutil.PrintError("Error message")
	cliutil.PrintWarning("Warning message")
	cliutil.PrintInfo("Info message")
}

func TestPrintTable(t *testing.T) {
	// Test empty table
	cliutil.PrintTable([][]string{})

	// Test normal table
	data := [][]string{
		{"Name", "Age", "City"},
		{"John", "30", "New York"},
		{"Jane", "25", "Boston"},
	}
	cliutil.PrintTable(data)

	// Test single row
	cliutil.PrintTable([][]string{{"Single", "Row"}})
}

// Benchmark tests
func BenchmarkParseArgs(b *testing.B) {
	args := []string{"--verbose", "--output", "file.txt", "-f", "config.json", "input.txt", "output.txt"}

	for i := 0; i < b.N; i++ {
		cliutil.ParseArgs(args)
	}
}

func BenchmarkHasFlag(b *testing.B) {
	args := []string{"--verbose", "--output=file.txt", "input.txt"}

	for i := 0; i < b.N; i++ {
		cliutil.HasFlag(args, "--verbose")
	}
}

func BenchmarkValidateEmail(b *testing.B) {
	email := "test@example.com"

	for i := 0; i < b.N; i++ {
		cliutil.ValidateEmail(email)
	}
}

// Test edge cases
func TestParseArgsEdgeCases(t *testing.T) {
	// Test flag at end without value
	args := cliutil.ParseArgs([]string{"input.txt", "--flag"})
	if !args.HasFlag("flag") {
		t.Error("Expected flag to be parsed as boolean flag")
	}

	// Test empty flag name
	args = cliutil.ParseArgs([]string{"--"})
	if len(args.BoolFlags) != 1 || !args.BoolFlags[""] {
		t.Error("Expected empty flag name to be handled")
	}

	// Test single dash
	args = cliutil.ParseArgs([]string{"-"})
	if len(args.Positional) != 1 || args.Positional[0] != "-" {
		t.Error("Expected single dash to be treated as positional argument")
	}
}

func TestProgressBarEdgeCases(t *testing.T) {
	// Test zero total
	pb := cliutil.NewProgressBar(0)
	pb.Update(0)
	pb.Finish()

	// Test negative values
	pb = cliutil.NewProgressBar(100)
	pb.Update(-1)  // Should not panic
	pb.Update(150) // Should not panic
	pb.Finish()
}

// Integration test to ensure all components work together
func TestIntegration(t *testing.T) {
	// Simulate a typical CLI workflow
	testArgs := []string{"--config", "test.yaml", "--verbose", "input.txt", "output.txt"}
	args := cliutil.ParseArgs(testArgs)

	// Check parsed arguments
	if !args.HasFlag("verbose") {
		t.Errorf("Expected verbose flag, got flags: %+v, boolFlags: %+v", args.Flags, args.BoolFlags)
	}

	config := args.GetFlagWithDefault("config", "default.yaml")
	if config != "test.yaml" {
		t.Errorf("Expected config to be 'test.yaml', got '%s'", config)
	}

	if len(args.Positional) != 2 {
		t.Errorf("Expected 2 positional arguments, got %d: %v", len(args.Positional), args.Positional)
	}

	// Test progress bar
	pb := cliutil.NewProgressBar(10)
	for i := 0; i <= 10; i++ {
		pb.Update(i)
	}
	pb.Finish()
}
