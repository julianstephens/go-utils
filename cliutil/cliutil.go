package cliutil

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
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
		
		if strings.HasPrefix(arg, "--") {
			// Long flag
			flagName := strings.TrimPrefix(arg, "--")
			
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
func isBooleanFlag(flagName string) bool {
	commonBoolFlags := map[string]bool{
		"verbose": true, "v": true,
		"debug": true, "d": true,
		"help": true, "h": true,
		"version": true,
		"quiet": true, "q": true,
		"force": true, "f": true,
		"dry-run": true,
		"interactive": true, "i": true,
		"recursive": true, "r": true,
		"all": true, "a": true,
	}
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
		if strings.HasPrefix(arg, flag+"=") {
			return strings.TrimPrefix(arg, flag+"=")
		}
	}
	return defaultValue
}

// PromptString prompts the user for string input
func PromptString(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return ""
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

// ValidateNonEmpty validates that input is not empty
func ValidateNonEmpty(input string) error {
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("input cannot be empty")
	}
	return nil
}

// ValidateEmail validates basic email format
func ValidateEmail(input string) error {
	if err := ValidateNonEmpty(input); err != nil {
		return err
	}
	parts := strings.Split(input, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("invalid email format")
	}
	if !strings.Contains(parts[1], ".") {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

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

// PrintColored prints text in the specified color
func PrintColored(text string, color Color) {
	fmt.Printf("%s%s%s\n", color, text, ColorReset)
}

// PrintSuccess prints a success message in green
func PrintSuccess(message string) {
	PrintColored("✓ "+message, ColorGreen)
}

// PrintError prints an error message in red
func PrintError(message string) {
	PrintColored("✗ "+message, ColorRed)
}

// PrintWarning prints a warning message in yellow
func PrintWarning(message string) {
	PrintColored("! "+message, ColorYellow)
}

// PrintInfo prints an info message in blue
func PrintInfo(message string) {
	PrintColored("ℹ "+message, ColorBlue)
}

// PrintTable prints data in a simple table format
func PrintTable(data [][]string) {
	if len(data) == 0 {
		return
	}

	// Calculate column widths
	colWidths := make([]int, len(data[0]))
	for _, row := range data {
		for i, cell := range row {
			if len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// Print rows
	for i, row := range data {
		for j, cell := range row {
			fmt.Printf("%-*s", colWidths[j]+2, cell)
		}
		fmt.Println()
		
		// Print separator after header
		if i == 0 {
			for j := range row {
				fmt.Print(strings.Repeat("-", colWidths[j]+2))
			}
			fmt.Println()
		}
	}
}

// ProgressBar represents a console progress bar
type ProgressBar struct {
	total   int
	current int
	width   int
	prefix  string
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total int) *ProgressBar {
	return &ProgressBar{
		total:  total,
		width:  50,
		prefix: "Progress",
	}
}

// NewProgressBarWithOptions creates a progress bar with custom options
func NewProgressBarWithOptions(total int, width int, prefix string) *ProgressBar {
	return &ProgressBar{
		total:  total,
		width:  width,
		prefix: prefix,
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
	fmt.Println()
}

// render renders the progress bar
func (pb *ProgressBar) render() {
	if pb.total == 0 {
		return
	}

	// Ensure current is within bounds
	current := pb.current
	if current < 0 {
		current = 0
	}
	if current > pb.total {
		current = pb.total
	}

	percentage := float64(current) / float64(pb.total)
	filled := int(percentage * float64(pb.width))
	
	// Ensure filled is within bounds
	if filled < 0 {
		filled = 0
	}
	if filled > pb.width {
		filled = pb.width
	}
	
	bar := strings.Repeat("█", filled) + strings.Repeat("░", pb.width-filled)
	
	fmt.Printf("\r%s: [%s] %.1f%% (%d/%d)", 
		pb.prefix, bar, percentage*100, current, pb.total)
}

// Spinner represents a console spinner
type Spinner struct {
	message   string
	frames    []string
	active    bool
	stopChan  chan bool
}

// NewSpinner creates a new spinner
func NewSpinner(message string) *Spinner {
	return &Spinner{
		message:  message,
		frames:   []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		stopChan: make(chan bool),
	}
}

// Start starts the spinner animation
func (s *Spinner) Start() {
	if s.active {
		return
	}
	
	s.active = true
	go func() {
		i := 0
		for {
			select {
			case <-s.stopChan:
				return
			default:
				fmt.Printf("\r%s %s", s.frames[i%len(s.frames)], s.message)
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

// Stop stops the spinner animation
func (s *Spinner) Stop() {
	if !s.active {
		return
	}
	
	s.active = false
	s.stopChan <- true
	fmt.Print("\r" + strings.Repeat(" ", len(s.message)+2) + "\r")
}

// UpdateMessage updates the spinner message
func (s *Spinner) UpdateMessage(message string) {
	s.message = message
}