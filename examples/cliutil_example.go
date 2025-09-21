package main

import (
	"fmt"
	"time"

	"github.com/julianstephens/go-utils/cliutil"
	"github.com/julianstephens/go-utils/logger"
)

func main() {
	// Example 1: Parse command-line arguments
	fmt.Println("=== CLI Argument Parsing ===")
	testArgs := []string{"--config", "app.yaml", "--verbose", "--output", "result.txt", "input1.txt", "input2.txt"}
	args := cliutil.ParseArgs(testArgs)
	
	fmt.Printf("Parsed arguments:\n")
	fmt.Printf("  Config file: %s\n", args.GetFlagWithDefault("config", "default.yaml"))
	fmt.Printf("  Verbose mode: %t\n", args.HasFlag("verbose"))
	fmt.Printf("  Output file: %s\n", args.GetFlag("output"))
	fmt.Printf("  Input files: %v\n", args.Positional)
	
	// Configure logger based on verbose flag
	if args.HasFlag("verbose") {
		logger.SetLogLevel("debug")
		logger.Info("Verbose mode enabled")
	}

	// Example 2: Colored output
	fmt.Println("\n=== Colored Output Examples ===")
	cliutil.PrintSuccess("Operation completed successfully!")
	cliutil.PrintError("Something went wrong")
	cliutil.PrintWarning("This is a warning message")
	cliutil.PrintInfo("Here's some information")
	cliutil.PrintColored("Custom blue message", cliutil.ColorBlue)

	// Example 3: Table output
	fmt.Println("\n=== Table Output ===")
	tableData := [][]string{
		{"Name", "Age", "Role", "Department"},
		{"Alice", "30", "Engineer", "Development"},
		{"Bob", "28", "Designer", "UX/UI"},
		{"Charlie", "35", "Manager", "Operations"},
	}
	cliutil.PrintTable(tableData)

	// Example 4: Progress bar
	fmt.Println("\n=== Progress Bar Demo ===")
	total := 50
	pb := cliutil.NewProgressBarWithOptions(total, 40, "Processing files")
	
	for i := 0; i <= total; i++ {
		pb.Update(i)
		time.Sleep(50 * time.Millisecond)
	}
	pb.Finish()
	fmt.Println("Processing complete!")

	// Example 5: Spinner demo
	fmt.Println("\n=== Spinner Demo ===")
	spinner := cliutil.NewSpinner("Loading configuration...")
	spinner.Start()
	time.Sleep(2 * time.Second)
	
	spinner.UpdateMessage("Connecting to database...")
	time.Sleep(1 * time.Second)
	
	spinner.UpdateMessage("Initializing modules...")
	time.Sleep(1 * time.Second)
	
	spinner.Stop()
	cliutil.PrintSuccess("System initialized successfully!")

	// Example 6: Validation functions
	fmt.Println("\n=== Validation Examples ===")
	testEmails := []string{"valid@example.com", "invalid-email", "@example.com", "user@domain.org"}
	
	for _, email := range testEmails {
		if err := cliutil.ValidateEmail(email); err != nil {
			cliutil.PrintError(fmt.Sprintf("Invalid email '%s': %v", email, err))
		} else {
			cliutil.PrintSuccess(fmt.Sprintf("Valid email: %s", email))
		}
	}

	// Example 7: Flag handling utilities
	fmt.Println("\n=== Flag Handling Utilities ===")
	sampleArgs := []string{"myapp", "--debug", "--port=8080", "--config", "prod.yaml", "serve"}
	
	if cliutil.HasFlag(sampleArgs, "--debug") {
		cliutil.PrintInfo("Debug mode is enabled")
	}
	
	port := cliutil.GetFlagValue(sampleArgs, "--port", "3000")
	cliutil.PrintInfo(fmt.Sprintf("Server will run on port: %s", port))
	
	configFile := cliutil.GetFlagValue(sampleArgs, "--config", "config.yaml")
	cliutil.PrintInfo(fmt.Sprintf("Using configuration file: %s", configFile))

	// Example 8: Interactive prompts (commented out since they require user input)
	/*
	fmt.Println("\n=== Interactive Prompts ===")
	name := cliutil.PromptString("Enter your name: ")
	cliutil.PrintInfo(fmt.Sprintf("Hello, %s!", name))
	
	email := cliutil.PromptStringWithValidation("Enter your email: ", cliutil.ValidateEmail)
	cliutil.PrintSuccess(fmt.Sprintf("Email %s is valid!", email))
	
	proceed := cliutil.PromptBool("Do you want to continue? (y/n): ")
	if proceed {
		cliutil.PrintSuccess("Continuing with the operation...")
	} else {
		cliutil.PrintWarning("Operation cancelled by user")
	}
	
	options := []string{"Create new project", "Open existing project", "Exit"}
	choice := cliutil.PromptChoice("What would you like to do?", options)
	cliutil.PrintInfo(fmt.Sprintf("You selected: %s", options[choice]))
	*/

	fmt.Println("\n=== CLI Utilities Demo Complete ===")
	cliutil.PrintInfo("All cliutil features demonstrated successfully!")
}