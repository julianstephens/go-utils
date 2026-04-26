package cliutil_test

import (
	"bytes"
	"encoding/json"
	"strings"
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
	tst.AssertTrue(
		t,
		args.GetFlagWithDefault("nonexistent", "default.txt") == "default.txt",
		"should get default for nonexistent flag",
	)
}

func TestHasFlag(t *testing.T) {
	args := []string{"--verbose", "--output=file.txt", "input.txt"}

	tst.AssertTrue(t, cliutil.HasFlag(args, "--verbose"), "should find --verbose")
	tst.AssertTrue(t, cliutil.HasFlag(args, "--output"), "should find --output")
	tst.AssertFalse(t, cliutil.HasFlag(args, "--debug"), "should not find --debug")
}

func TestGetFlagValue(t *testing.T) {
	args := []string{"--output", "file.txt", "--config=config.yaml", "input.txt"}

	tst.AssertTrue(
		t,
		cliutil.GetFlagValue(args, "--output", "default") == "file.txt",
		"GetFlagValue --output should return file.txt",
	)
	tst.AssertTrue(
		t,
		cliutil.GetFlagValue(args, "--config", "default") == "config.yaml",
		"GetFlagValue --config should return config.yaml",
	)
	tst.AssertTrue(
		t,
		cliutil.GetFlagValue(args, "--nonexistent", "default") == "default",
		"GetFlagValue nonexistent should return default",
	)
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

// ---------------------------------------------------------------------------
// Formatter tests
// ---------------------------------------------------------------------------

func TestColoredFormatter(t *testing.T) {
	f := cliutil.NewColoredFormatter()

	tst.AssertTrue(t, strings.Contains(f.Success("ok"), "ok"), "Success should contain text")
	tst.AssertTrue(t, strings.Contains(f.Error("fail"), "fail"), "Error should contain text")
	tst.AssertTrue(t, strings.Contains(f.Warning("warn"), "warn"), "Warning should contain text")
	tst.AssertTrue(t, strings.Contains(f.Info("info"), "info"), "Info should contain text")

	// Should contain ANSI codes
	tst.AssertTrue(t, strings.Contains(f.Success("x"), "\033["), "Success should contain ANSI escape")
	tst.AssertTrue(t, strings.Contains(f.Error("x"), "\033["), "Error should contain ANSI escape")

	// Colored wrapper
	tst.AssertTrue(t, strings.Contains(f.Colored("hi", cliutil.ColorBlue), "hi"), "Colored should contain text")
	tst.AssertTrue(
		t,
		strings.Contains(f.Colored("hi", cliutil.ColorBlue), "\033["),
		"Colored should contain ANSI escape",
	)

	// Visual element characters
	tst.AssertTrue(t, f.ProgressCharFilled() != "", "filled char should not be empty")
	tst.AssertTrue(t, f.ProgressCharEmpty() != "", "empty char should not be empty")
	tst.AssertTrue(t, len(f.SpinnerFrames()) > 0, "spinner frames should not be empty")
}

func TestPlainTextFormatter(t *testing.T) {
	f := cliutil.NewPlainTextFormatter()

	// No ANSI codes
	tst.AssertFalse(t, strings.Contains(f.Success("ok"), "\033["), "PlainText Success should not contain ANSI")
	tst.AssertFalse(t, strings.Contains(f.Error("fail"), "\033["), "PlainText Error should not contain ANSI")
	tst.AssertFalse(t, strings.Contains(f.Warning("warn"), "\033["), "PlainText Warning should not contain ANSI")
	tst.AssertFalse(t, strings.Contains(f.Info("info"), "\033["), "PlainText Info should not contain ANSI")

	// Text is present
	tst.AssertTrue(t, strings.Contains(f.Success("ok"), "ok"), "PlainText Success should contain text")
	tst.AssertTrue(t, strings.Contains(f.Error("fail"), "fail"), "PlainText Error should contain text")

	// Colored strips color
	tst.AssertTrue(t, f.Colored("hi", cliutil.ColorBlue) == "hi", "PlainText Colored should return bare text")

	tst.AssertTrue(t, f.ProgressCharFilled() != "", "filled char should not be empty")
	tst.AssertTrue(t, f.ProgressCharEmpty() != "", "empty char should not be empty")
	tst.AssertTrue(t, len(f.SpinnerFrames()) > 0, "spinner frames should not be empty")
}

func TestAccessibleFormatter(t *testing.T) {
	f := cliutil.NewAccessibleFormatter()

	tst.AssertTrue(t, strings.HasPrefix(f.Success("x"), "[SUCCESS]"), "Accessible Success should have [SUCCESS] prefix")
	tst.AssertTrue(t, strings.HasPrefix(f.Error("x"), "[ERROR]"), "Accessible Error should have [ERROR] prefix")
	tst.AssertTrue(t, strings.HasPrefix(f.Warning("x"), "[WARNING]"), "Accessible Warning should have [WARNING] prefix")
	tst.AssertTrue(t, strings.HasPrefix(f.Info("x"), "[INFO]"), "Accessible Info should have [INFO] prefix")

	tst.AssertFalse(t, strings.Contains(f.Success("x"), "\033["), "Accessible should not contain ANSI codes")
	tst.AssertTrue(t, f.Colored("hi", cliutil.ColorRed) == "hi", "Accessible Colored should return bare text")
}

func TestJSONFormatter(t *testing.T) {
	f := cliutil.NewJSONFormatter()

	assertJSON := func(s, expectedType, expectedMsg string) {
		t.Helper()
		var m map[string]string
		if err := json.Unmarshal([]byte(s), &m); err != nil {
			t.Fatalf("not valid JSON: %q, err: %v", s, err)
		}
		tst.AssertTrue(t, m["type"] == expectedType, "type field should be "+expectedType)
		tst.AssertTrue(t, m["message"] == expectedMsg, "message field should match")
	}

	assertJSON(f.Success("done"), "success", "done")
	assertJSON(f.Error("oops"), "error", "oops")
	assertJSON(f.Warning("careful"), "warning", "careful")
	assertJSON(f.Info("fyi"), "info", "fyi")
	assertJSON(f.Colored("blue text", cliutil.ColorBlue), "colored", "blue text")

	// Table returns a JSON array
	data := [][]string{
		{"Name", "Age"},
		{"Alice", "30"},
		{"Bob", "25"},
	}
	tableJSON := f.Table(data)
	var rows []map[string]string
	tst.AssertNoError(t, json.Unmarshal([]byte(tableJSON), &rows))
	tst.AssertTrue(t, len(rows) == 2, "table should have 2 data rows")
	tst.AssertTrue(t, rows[0]["Name"] == "Alice", "first row Name should be Alice")
	tst.AssertTrue(t, rows[1]["Age"] == "25", "second row Age should be 25")

	// Animations are suppressed
	tst.AssertTrue(t, f.ProgressCharFilled() == "", "JSON formatter should suppress progress fill char")
	tst.AssertTrue(t, f.ProgressCharEmpty() == "", "JSON formatter should suppress progress empty char")
	tst.AssertTrue(t, len(f.SpinnerFrames()) == 0, "JSON formatter should suppress spinner frames")
}

func TestFormatterTable(t *testing.T) {
	formatters := []cliutil.Formatter{
		cliutil.NewColoredFormatter(),
		cliutil.NewPlainTextFormatter(),
		cliutil.NewAccessibleFormatter(),
	}
	data := [][]string{
		{"Name", "Age"},
		{"Alice", "30"},
	}

	for _, f := range formatters {
		result := f.Table(data)
		tst.AssertTrue(t, strings.Contains(result, "Alice"), "table should contain data")
		tst.AssertTrue(t, strings.Contains(result, "Name"), "table should contain headers")
		tst.AssertTrue(t, strings.Contains(result, "---"), "table should have separator")
	}
}

// ---------------------------------------------------------------------------
// Console integration tests
// ---------------------------------------------------------------------------

func newTestConsole(f cliutil.Formatter) (*cliutil.Console, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	return cliutil.NewConsole(f, buf), buf
}

func TestConsole_SemanticMethods(t *testing.T) {
	c, buf := newTestConsole(cliutil.NewPlainTextFormatter())

	tst.AssertNoError(t, c.Success("done"))
	tst.AssertTrue(t, strings.Contains(buf.String(), "done"), "buffer should contain success text")
	buf.Reset()

	tst.AssertNoError(t, c.Error("oops"))
	tst.AssertTrue(t, strings.Contains(buf.String(), "oops"), "buffer should contain error text")
	buf.Reset()

	tst.AssertNoError(t, c.Warning("warn"))
	tst.AssertTrue(t, strings.Contains(buf.String(), "warn"), "buffer should contain warning text")
	buf.Reset()

	tst.AssertNoError(t, c.Info("fyi"))
	tst.AssertTrue(t, strings.Contains(buf.String(), "fyi"), "buffer should contain info text")
}

func TestConsole_PrintColored(t *testing.T) {
	c, buf := newTestConsole(cliutil.NewPlainTextFormatter())
	tst.AssertNoError(t, c.PrintColored("hello", cliutil.ColorGreen))
	tst.AssertTrue(t, strings.Contains(buf.String(), "hello"), "buffer should contain colored text")
}

func TestConsole_PrintTable(t *testing.T) {
	c, buf := newTestConsole(cliutil.NewPlainTextFormatter())
	data := [][]string{
		{"Name", "City"},
		{"Alice", "NYC"},
	}
	tst.AssertNoError(t, c.PrintTable(data))
	out := buf.String()
	tst.AssertTrue(t, strings.Contains(out, "Alice"), "table output should contain Alice")
	tst.AssertTrue(t, strings.Contains(out, "Name"), "table output should contain header")

	// Empty table should not error
	buf.Reset()
	tst.AssertNoError(t, c.PrintTable([][]string{}))
	tst.AssertTrue(t, buf.Len() == 0, "empty table should produce no output")
}

func TestConsole_SetFormatter(t *testing.T) {
	c, buf := newTestConsole(cliutil.NewColoredFormatter())
	tst.AssertNoError(t, c.Success("start"))
	tst.AssertTrue(t, strings.Contains(buf.String(), "\033["), "colored formatter should emit ANSI")

	buf.Reset()
	c.SetFormatter(cliutil.NewPlainTextFormatter())
	tst.AssertNoError(t, c.Success("switched"))
	tst.AssertFalse(t, strings.Contains(buf.String(), "\033["), "plain formatter should not emit ANSI")
	tst.AssertTrue(t, strings.Contains(buf.String(), "switched"), "plain formatter should contain text")
}

func TestConsole_FormatterGetter(t *testing.T) {
	f := cliutil.NewPlainTextFormatter()
	c, _ := newTestConsole(f)
	tst.AssertTrue(t, c.Formatter() == f, "Formatter() should return the set formatter")
}

func TestConsole_ProgressBar(t *testing.T) {
	c, buf := newTestConsole(cliutil.NewPlainTextFormatter())
	pb := cliutil.NewProgressBarWithConsole(c, 100, 20, "Loading")

	pb.Update(50)
	pb.Finish()

	out := buf.String()
	tst.AssertTrue(t, strings.Contains(out, "Loading"), "progress bar should contain prefix")
	tst.AssertTrue(t, strings.Contains(out, "50.0%") || strings.Contains(out, "100.0%"),
		"progress bar should contain percentage")
}

func TestConsole_ProgressBarJSON(t *testing.T) {
	// JSONFormatter suppresses animated progress bar output
	c, buf := newTestConsole(cliutil.NewJSONFormatter())
	pb := cliutil.NewProgressBarWithConsole(c, 100, 20, "Loading")
	pb.Update(50)
	pb.Finish()
	tst.AssertTrue(t, buf.Len() == 0, "JSON formatter should suppress progress bar output")
}

func TestConsole_Spinner(t *testing.T) {
	c, buf := newTestConsole(cliutil.NewPlainTextFormatter())
	s := cliutil.NewSpinnerWithConsole(c, "Working")

	s.Start()
	time.Sleep(250 * time.Millisecond)
	s.Stop()

	tst.AssertTrue(t, buf.Len() > 0, "spinner should have written to buffer")
}

func TestConsole_SpinnerJSON(t *testing.T) {
	// JSONFormatter suppresses spinner output
	c, buf := newTestConsole(cliutil.NewJSONFormatter())
	s := cliutil.NewSpinnerWithConsole(c, "Working")
	s.Start()
	time.Sleep(100 * time.Millisecond)
	s.Stop()
	tst.AssertTrue(t, buf.Len() == 0, "JSON formatter should suppress spinner output")
}

func TestSetDefaultConsole(t *testing.T) {
	buf := &bytes.Buffer{}
	c := cliutil.NewConsole(cliutil.NewPlainTextFormatter(), buf)
	cliutil.SetDefaultConsole(c)

	cliutil.PrintSuccess("via default")
	tst.AssertTrue(t, strings.Contains(buf.String(), "via default"),
		"module-level PrintSuccess should use updated default Console")

	// Restore default to avoid polluting other tests
	cliutil.SetDefaultConsole(cliutil.NewConsole(cliutil.NewColoredFormatter(), bytes.NewBuffer(nil)))
}

func TestConsole_JSONOutput(t *testing.T) {
	c, buf := newTestConsole(cliutil.NewJSONFormatter())

	tst.AssertNoError(t, c.Success("all good"))
	var m map[string]string
	tst.AssertNoError(t, json.Unmarshal(bytes.TrimRight(buf.Bytes(), "\n"), &m))
	tst.AssertTrue(t, m["type"] == "success", "JSON type should be success")
	tst.AssertTrue(t, m["message"] == "all good", "JSON message should match")
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
