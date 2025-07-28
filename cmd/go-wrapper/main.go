package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
		os.Exit(1)
	}

	// Look for main.go in the current directory
	mainGoPath := filepath.Join(wd, "main.go")
	if _, err := os.Stat(mainGoPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: main.go not found in current directory\n")
		os.Exit(1)
	}

	args := os.Args[1:] // Remove program name

	var cmd *exec.Cmd

	if len(args) == 0 {
		// Just "go" - pass through to regular go command
		cmd = exec.Command("go")
	} else if args[0] == "run" && len(args) == 1 {
		// "go run" with no additional args - run main.go
		cmd = exec.Command("go", "run", "main.go")
	} else if args[0] == "fresh" {
		// "go fresh" - run main.go with fresh argument
		cmd = exec.Command("go", "run", "main.go", "fresh")
	} else if args[0] == "run" && len(args) > 1 {
		// Check if the second argument looks like a .go file or path
		secondArg := args[1]
		if filepath.Ext(secondArg) == ".go" || secondArg == "." || secondArg[:1] == "./" {
			// Normal go run with file - pass through
			cmd = exec.Command("go", args...)
		} else {
			// "go run <arg>" where arg is not a file - treat as argument to main.go
			runArgs := append([]string{"run", "main.go"}, args[1:]...)
			cmd = exec.Command("go", runArgs...)
		}
	} else {
		// All other go commands - pass through
		cmd = exec.Command("go", args...)
	}

	// Set up command to inherit our stdin/stdout/stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Execute the command
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
		os.Exit(1)
	}
}
