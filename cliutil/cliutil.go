package cliutil

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/term"
)

// Args represents parsed command-line arguments
type Args struct {
	// Flags contains named flags (e.g., --flag=value or --flag value)
	Flags map[string]string
	// BoolFlags contains boolean flags (e.g., --verbose)
	BoolFlags map[string]bool
	// Positional contains positional arguments
	Positional []string
}

// ParseArgs parses command-line arguments into a structured format
func ParseArgs(args []string) *Args {
	result := &Args{
		Flags:      make(map[string]string),
		BoolFlags:  make(map[string]bool),
		Positional: make([]string, 0),
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if after, ok := strings.CutPrefix(arg, "--"); ok {
			// Long flag
			flagName := after

			if strings.Contains(flagName, "=") {
				// --flag=value format
				parts := strings.SplitN(flagName, "=", 2)
				result.Flags[parts[0]] = parts[1]
			} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				// Check if next argument looks like a value (not starting with -)
				// But we need to be smarter about boolean flags
				nextArg := args[i+1]

				// If the next argument is a common boolean value, treat as boolean
				lowerNext := strings.ToLower(nextArg)
				if lowerNext == "true" || lowerNext == "false" {
					result.Flags[flagName] = nextArg
					i++ // skip next argument
				} else {
					// Check common boolean flag patterns
					isBoolFlag := isBooleanFlag(flagName)
					if isBoolFlag {
						result.BoolFlags[flagName] = true
					} else {
						// --flag value format
						result.Flags[flagName] = nextArg
						i++ // skip next argument
					}
				}
			} else {
				// Boolean flag
				result.BoolFlags[flagName] = true
			}
		} else if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			// Short flag
			flagName := strings.TrimPrefix(arg, "-")

			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				// Check if this looks like a boolean flag
				if isBooleanFlag(flagName) {
					result.BoolFlags[flagName] = true
				} else {
					// -f value format
					result.Flags[flagName] = args[i+1]
					i++ // skip next argument
				}
			} else {
				// Boolean flag
				result.BoolFlags[flagName] = true
			}
		} else {
			// Positional argument
			result.Positional = append(result.Positional, arg)
		}
	}

	return result
}

// isBooleanFlag checks if a flag name is commonly used as a boolean flag
var commonBoolFlags = map[string]bool{
	"verbose": true, "v": true,
	"debug": true, "d": true,
	"help": true, "h": true,
	"version": true,
	"quiet":   true, "q": true,
	"force": true, "f": true,
	"dry-run":     true,
	"interactive": true, "i": true,
	"recursive": true, "r": true,
	"all": true, "a": true,
}

func isBooleanFlag(flagName string) bool {
	return commonBoolFlags[flagName]
}

// HasFlag checks if a boolean flag is present
func (a *Args) HasFlag(name string) bool {
	return a.BoolFlags[name]
}

// GetFlag returns the value of a named flag, or empty string if not found
func (a *Args) GetFlag(name string) string {
	return a.Flags[name]
}

// GetFlagWithDefault returns the value of a named flag, or the default value if not found
func (a *Args) GetFlagWithDefault(name, defaultValue string) string {
	if value, exists := a.Flags[name]; exists {
		return value
	}
	return defaultValue
}

// HasFlag checks if a boolean flag is present in the arguments slice
func HasFlag(args []string, flag string) bool {
	for _, arg := range args {
		if arg == flag || strings.HasPrefix(arg, flag+"=") {
			return true
		}
	}
	return false
}

// GetFlagValue returns the value of a flag from arguments slice, or default if not found
func GetFlagValue(args []string, flag, defaultValue string) string {
	for i, arg := range args {
		if arg == flag && i+1 < len(args) {
			return args[i+1]
		}
		if after, ok := strings.CutPrefix(arg, flag+"="); ok {
			return after
		}
	}
	return defaultValue
}

// PromptString prompts the user for string input
func PromptString(prompt string) string {
	return PromptStringWithIO(prompt, os.Stdin, os.Stdout)
}

// PromptStringWithValidation prompts for string input with validation
func PromptStringWithValidation(prompt string, validator func(string) error) string {
	for {
		input := PromptString(prompt)
		if err := validator(input); err != nil {
			PrintError(fmt.Sprintf("Invalid input: %v", err))
			continue
		}
		return input
	}
}

// PromptStringWithIO prompts the user for string input using the provided
// io.Reader and io.Writer. This is useful for testing where stdin/stdout can
// be simulated.
func PromptStringWithIO(prompt string, in io.Reader, out io.Writer) string {
	_, _ = fmt.Fprint(out, prompt)
	scanner := bufio.NewScanner(in)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return ""
}

// PromptPassword prompts the user for a password (secure string) without echoing
// the input to the terminal. If stdin is not a terminal, it falls back to
// reading a normal line (echoed).
func PromptPassword(prompt string) string {
	// Print prompt without newline so password can be entered on same line
	return PromptPasswordWithIO(prompt, os.Stdin, os.Stdout)
}

// PromptPasswordWithValidation prompts for a secure string and validates it
// using the provided validator function. It will re-prompt until the validator
// returns nil.
func PromptPasswordWithValidation(prompt string, validator func(string) error) string {
	return PromptPasswordWithValidationIO(prompt, os.Stdin, os.Stdout, validator)
}

// PromptPasswordWithIO prompts for a secure string using the provided reader
// and writer. If the reader is a terminal (os.File), it will attempt to use
// term.ReadPassword to disable echo. Otherwise it falls back to a normal
// line-read (echoed).
func PromptPasswordWithIO(prompt string, in io.Reader, out io.Writer) string {
	// Print prompt without newline so password can be entered on same line
	_, _ = fmt.Fprint(out, prompt)

	// If in is an *os.File and a terminal, use ReadPassword to disable echo
	if f, ok := in.(*os.File); ok {
		if term.IsTerminal(int(f.Fd())) {
			b, err := term.ReadPassword(int(f.Fd()))
			_, _ = fmt.Fprintln(out)
			if err == nil {
				return strings.TrimSpace(string(b))
			}
			// Fall through to fallback on error
		}
	}

	// Fallback: read a normal line (echoed)
	reader := bufio.NewReader(in)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// PromptPasswordWithValidationIO prompts for a secure string and validates it
// using the provided validator function and I/O streams. It will re-prompt
// until the validator returns nil.
func PromptPasswordWithValidationIO(prompt string, in io.Reader, out io.Writer, validator func(string) error) string {
	for {
		pw := PromptPasswordWithIO(prompt, in, out)
		if validator == nil {
			return pw
		}
		if err := validator(pw); err != nil {
			PrintError(fmt.Sprintf("Invalid input: %v", err))
			continue
		}
		return pw
	}
}

// PromptBool prompts the user for a yes/no response
func PromptBool(prompt string) bool {
	for {
		response := strings.ToLower(PromptString(prompt))
		switch response {
		case "y", "yes", "true", "1":
			return true
		case "n", "no", "false", "0":
			return false
		default:
			PrintWarning("Please enter y/yes or n/no")
		}
	}
}

// PromptChoice prompts the user to select from a list of options
func PromptChoice(prompt string, options []string) int {
	fmt.Println(prompt)
	for i, option := range options {
		fmt.Printf("  %d) %s\n", i+1, option)
	}

	for {
		input := PromptString("Enter your choice (number): ")
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(options) {
			PrintError(fmt.Sprintf("Please enter a number between 1 and %d", len(options)))
			continue
		}
		return choice - 1 // Return 0-based index
	}
}

// ValidationFunc is a function type for input validation
type ValidationFunc func(string) error

// Color constants for colored output
type Color string

const (
	ColorReset   Color = "\033[0m"
	ColorRed     Color = "\033[31m"
	ColorGreen   Color = "\033[32m"
	ColorYellow  Color = "\033[33m"
	ColorBlue    Color = "\033[34m"
	ColorMagenta Color = "\033[35m"
	ColorCyan    Color = "\033[36m"
	ColorWhite   Color = "\033[37m"
	ColorBold    Color = "\033[1m"
)

// PrintColored prints text styled with the given Color using the default Console.
func PrintColored(text string, color Color) {
	_ = DefaultConsole().PrintColored(text, color)
}

// PrintSuccess prints a success message using the default Console.
func PrintSuccess(message string) {
	_ = DefaultConsole().Success(message)
}

// PrintError prints an error message using the default Console.
func PrintError(message string) {
	_ = DefaultConsole().Error(message)
}

// PrintWarning prints a warning message using the default Console.
func PrintWarning(message string) {
	_ = DefaultConsole().Warning(message)
}

// PrintInfo prints an info message using the default Console.
func PrintInfo(message string) {
	_ = DefaultConsole().Info(message)
}

// PrintTable prints data as a formatted table using the default Console.
func PrintTable(data [][]string) {
	_ = DefaultConsole().PrintTable(data)
}

// ProgressBar represents a console progress bar
type ProgressBar struct {
	console *Console
	total   int
	current int
	width   int
	prefix  string
}

// NewProgressBar creates a new progress bar that writes to the default Console.
func NewProgressBar(total int) *ProgressBar {
	return NewProgressBarWithConsole(DefaultConsole(), total, 50, "Progress")
}

// NewProgressBarWithOptions creates a progress bar with custom options that
// writes to the default Console.
func NewProgressBarWithOptions(total int, width int, prefix string) *ProgressBar {
	return NewProgressBarWithConsole(DefaultConsole(), total, width, prefix)
}

// NewProgressBarWithConsole creates a progress bar that writes through c.
func NewProgressBarWithConsole(c *Console, total int, width int, prefix string) *ProgressBar {
	return &ProgressBar{
		console: c,
		total:   total,
		width:   width,
		prefix:  prefix,
	}
}

// Update updates the progress bar
func (pb *ProgressBar) Update(current int) {
	pb.current = current
	pb.render()
}

// Finish completes the progress bar
func (pb *ProgressBar) Finish() {
	pb.current = pb.total
	pb.render()
	if pb.console.formatter.ProgressCharFilled() != "" {
		_ = pb.console.writeTo(pb.console.out, "\n")
	}
}

// render renders the progress bar
func (pb *ProgressBar) render() {
	if pb.total == 0 {
		return
	}

	current := pb.current
	if current < 0 {
		current = 0
	}
	if current > pb.total {
		current = pb.total
	}

	percentage := float64(current) / float64(pb.total)
	filled := int(percentage * float64(pb.width))
	if filled < 0 {
		filled = 0
	}
	if filled > pb.width {
		filled = pb.width
	}

	f := pb.console.formatter
	filledChar := f.ProgressCharFilled()
	emptyChar := f.ProgressCharEmpty()
	if filledChar == "" {
		// Formatter suppresses animations (e.g. JSONFormatter)
		return
	}

	bar := strings.Repeat(filledChar, filled) + strings.Repeat(emptyChar, pb.width-filled)
	line := fmt.Sprintf("\r%s: [%s] %.1f%% (%d/%d)",
		pb.prefix, bar, percentage*100, current, pb.total)
	_ = pb.console.writeTo(pb.console.out, line)
}

// Spinner represents a console spinner
type Spinner struct {
	console  *Console
	message  string
	frames   []string
	active   bool
	stopChan chan bool
	mu       sync.RWMutex
}

// NewSpinner creates a new spinner that writes to the default Console.
func NewSpinner(message string) *Spinner {
	return NewSpinnerWithConsole(DefaultConsole(), message)
}

// NewSpinnerWithConsole creates a spinner that writes through c.
func NewSpinnerWithConsole(c *Console, message string) *Spinner {
	frames := c.formatter.SpinnerFrames()
	if len(frames) == 0 {
		frames = []string{"|"} // Sentinel; Start() will treat empty frames as suppressed
	}
	return &Spinner{
		console:  c,
		message:  message,
		frames:   frames,
		stopChan: make(chan bool),
		mu:       sync.RWMutex{},
	}
}

// Start starts the spinner animation
func (s *Spinner) Start() {
	// If the formatter suppresses animations, do nothing.
	if len(s.console.formatter.SpinnerFrames()) == 0 {
		return
	}

	s.mu.Lock()
	if s.active {
		s.mu.Unlock()
		return
	}
	s.active = true
	s.mu.Unlock()

	go func() {
		i := 0
		for {
			select {
			case <-s.stopChan:
				return
			default:
				s.mu.RLock()
				message := s.message
				s.mu.RUnlock()
				_ = s.console.writeTo(s.console.out,
					fmt.Sprintf("\r%s %s", s.frames[i%len(s.frames)], message))
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

// Stop stops the spinner animation
func (s *Spinner) Stop() {
	if len(s.console.formatter.SpinnerFrames()) == 0 {
		return
	}

	s.mu.Lock()
	if !s.active {
		s.mu.Unlock()
		return
	}
	s.active = false
	s.mu.Unlock()

	select {
	case s.stopChan <- true:
	default:
		// spinner goroutine already exited
	}
	s.mu.RLock()
	msgLen := len(s.message)
	s.mu.RUnlock()
	_ = s.console.writeTo(s.console.out, "\r"+strings.Repeat(" ", msgLen+2)+"\r")
}

// UpdateMessage updates the spinner message
func (s *Spinner) UpdateMessage(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.message = message
}
