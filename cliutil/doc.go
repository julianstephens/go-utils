/*
Package cliutil provides helpers and utilities for building command-line interfaces (CLI) in Go projects.

This package offers reusable functions for common CLI patterns such as argument parsing,
flag handling, interactive prompts, progress indicators, and colored output. It is designed
to be simple to import, idiomatic, and consistent across julianstephens Go repositories.

Basic Usage:

The cliutil package provides functions for various CLI operations:

	package main

	import (
		"fmt"
		"github.com/julianstephens/go-utils/cliutil"
	)

	func main() {
		// Parse command-line arguments
		args := cliutil.ParseArgs(os.Args[1:])
		fmt.Printf("Parsed arguments: %+v\n", args)

		// Interactive prompts
		name := cliutil.PromptString("Enter your name: ")
		confirm := cliutil.PromptBool("Continue? (y/n): ")

		// Colored output
		cliutil.PrintSuccess("Operation completed successfully!")
		cliutil.PrintError("An error occurred")
		cliutil.PrintWarning("This is a warning")

		// Progress indicators
		progress := cliutil.NewProgressBar(100)
		for i := 0; i <= 100; i++ {
			progress.Update(i)
			time.Sleep(10 * time.Millisecond)
		}
		progress.Finish()
	}

Argument Parsing:

The package provides simple argument parsing utilities:

	// Parse arguments into a structured format
	args := cliutil.ParseArgs([]string{"--verbose", "--output", "file.txt", "input.txt"})

	// Access flags and positional arguments
	verbose := args.HasFlag("verbose")
	output := args.GetFlag("output")
	files := args.Positional

Flag Handling:

Support for common flag patterns:

	// Check for boolean flags
	if cliutil.HasFlag(os.Args, "--debug") {
		// Enable debug mode
	}

	// Get flag values
	logLevel := cliutil.GetFlagValue(os.Args, "--log-level", "info")

Interactive Prompts:

Functions for user input with validation:

	// String input with optional validation
	email := cliutil.PromptStringWithValidation("Email: ", cliutil.ValidateEmail)

	// Boolean prompts
	proceed := cliutil.PromptBool("Do you want to continue? (y/n): ")

	// Choice selection
	option := cliutil.PromptChoice("Select an option:", []string{"Option 1", "Option 2", "Option 3"})

Output Formatting:

Colored output and formatting utilities:

	// Colored text output
	cliutil.PrintSuccess("✓ Task completed")
	cliutil.PrintError("✗ Task failed")
	cliutil.PrintWarning("! Warning message")
	cliutil.PrintInfo("ℹ Information")

	// Custom colored output
	cliutil.PrintColored("Custom message", cliutil.ColorBlue)

	// Formatted output
	cliutil.PrintTable([][]string{
		{"Name", "Age", "City"},
		{"John", "30", "New York"},
		{"Jane", "25", "Boston"},
	})

Progress Indicators:

Progress bars and spinners for long-running operations:

	// Progress bar
	bar := cliutil.NewProgressBar(total)
	for i := 0; i < total; i++ {
		// Do work
		bar.Update(i + 1)
	}
	bar.Finish()

	// Spinner for indeterminate progress
	spinner := cliutil.NewSpinner("Processing...")
	spinner.Start()
	defer spinner.Stop()

	// Do work
	time.Sleep(5 * time.Second)

Integration:

This package is designed to be used across julianstephens Go repositories for consistent
CLI development. It integrates well with other go-utils packages like logger and config:

	// Use with logger for structured CLI logging
	if cliutil.HasFlag(os.Args, "--verbose") {
		logger.SetLogLevel("debug")
	}

	// Use with config for CLI configuration
	var cfg AppConfig
	configFile := cliutil.GetFlagValue(os.Args, "--config", "config.yaml")
	config.LoadFromFile(configFile, &cfg)

Thread Safety:

Progress indicators and interactive prompts are designed to be thread-safe for
concurrent CLI operations. However, colored output functions should be used
from a single goroutine to prevent mixed output.
*/
package cliutil
