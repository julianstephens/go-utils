package cliutil_test

import (
	"testing"
	"time"

	"github.com/julianstephens/go-utils/cliutil"
	tst "github.com/julianstephens/go-utils/tests"
	"github.com/julianstephens/go-utils/validator"
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

			tst.AssertDeepEqual(t, result.Flags, tt.expected.Flags)
			tst.AssertDeepEqual(t, result.BoolFlags, tt.expected.BoolFlags)
			tst.AssertDeepEqual(t, result.Positional, tt.expected.Positional)
		})
	}
}

func TestArgs_HasFlag(t *testing.T) {
	args := cliutil.ParseArgs([]string{"--verbose", "--debug"})

	tst.AssertTrue(t, args.HasFlag("verbose"), "verbose flag should be true")
	tst.AssertTrue(t, args.HasFlag("debug"), "debug flag should be true")
	tst.AssertFalse(t, args.HasFlag("quiet"), "quiet flag should be false")
}

func TestArgs_GetFlag(t *testing.T) {
	args := cliutil.ParseArgs([]string{"--output", "file.txt", "--config=config.yaml"})

	tst.AssertTrue(t, args.GetFlag("output") == "file.txt", "output flag should be file.txt")
	tst.AssertTrue(t, args.GetFlag("config") == "config.yaml", "config flag should be config.yaml")
	tst.AssertTrue(t, args.GetFlag("nonexistent") == "", "nonexistent flag should be empty")
}

func TestArgs_GetFlagWithDefault(t *testing.T) {
	args := cliutil.ParseArgs([]string{"--output", "file.txt"})

	tst.AssertTrue(t, args.GetFlagWithDefault("output", "default.txt") == "file.txt", "should get actual flag value")
	tst.AssertTrue(t, args.GetFlagWithDefault("nonexistent", "default.txt") == "default.txt", "should get default for nonexistent flag")
}

func TestHasFlag(t *testing.T) {
	args := []string{"--verbose", "--output=file.txt", "input.txt"}

	tst.AssertTrue(t, cliutil.HasFlag(args, "--verbose"), "should find --verbose")
	tst.AssertTrue(t, cliutil.HasFlag(args, "--output"), "should find --output")
	tst.AssertFalse(t, cliutil.HasFlag(args, "--debug"), "should not find --debug")
}

func TestGetFlagValue(t *testing.T) {
	args := []string{"--output", "file.txt", "--config=config.yaml", "input.txt"}

	tst.AssertTrue(t, cliutil.GetFlagValue(args, "--output", "default") == "file.txt", "GetFlagValue --output should return file.txt")
	tst.AssertTrue(t, cliutil.GetFlagValue(args, "--config", "default") == "config.yaml", "GetFlagValue --config should return config.yaml")
	tst.AssertTrue(t, cliutil.GetFlagValue(args, "--nonexistent", "default") == "default", "GetFlagValue nonexistent should return default")
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
		err := validator.ValidateNonEmpty(tt.input)
		if tt.shouldErr {
			tst.AssertNotNil(t, err, "expected error for input")
		} else {
			tst.AssertNoError(t, err)
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
		err := validator.ValidateEmail(tt.input)
		if tt.shouldErr {
			tst.AssertNotNil(t, err, "expected error for email")
		} else {
			tst.AssertNoError(t, err)
		}
	}
}

func TestNewProgressBar(t *testing.T) {
	pb := cliutil.NewProgressBar(100)
	tst.AssertNotNil(t, pb, "progress bar should be created")

	// Test update (we can't easily test the visual output, but we can test it doesn't panic)
	pb.Update(50)
	pb.Update(100)
	pb.Finish()
}

func TestNewProgressBarWithOptions(t *testing.T) {
	pb := cliutil.NewProgressBarWithOptions(200, 30, "Custom Progress")
	tst.AssertNotNil(t, pb, "progress bar with options should be created")

	pb.Update(100)
	pb.Finish()
}

func TestNewSpinner(t *testing.T) {
	spinner := cliutil.NewSpinner("Loading...")
	tst.AssertNotNil(t, spinner, "spinner should be created")

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
		tst.AssertTrue(t, string(color) != "", "color constant should not be empty")
		_ = i
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
		validator.ValidateEmail(email)
	}
}

// Test edge cases
func TestParseArgsEdgeCases(t *testing.T) {
	// Test flag at end without value
	args := cliutil.ParseArgs([]string{"input.txt", "--flag"})
	tst.AssertTrue(t, args.HasFlag("flag"), "flag should be parsed as boolean flag")

	// Test empty flag name
	args = cliutil.ParseArgs([]string{"--"})
	tst.AssertTrue(t, len(args.BoolFlags) == 1 && args.BoolFlags[""], "empty flag name should be handled")

	// Test single dash
	args = cliutil.ParseArgs([]string{"-"})
	tst.AssertTrue(t, len(args.Positional) == 1 && args.Positional[0] == "-", "single dash should be positional")
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
