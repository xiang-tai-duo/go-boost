// Package boost
// File:        process.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/examples/process.go
// Author:      Vibe Coding
// Created:     12/30/2025 11:03:46
// Description: Example for Process utility functions
// --------------------------------------------------------------------------------

package main

import (
	"fmt"
	"time"

	. "github.com/xiang-tai-duo/go-boost"
)

func main() {

	// Create a new Process instance for current process
	process := NewProcess()

	// Get process information
	fmt.Println("--- Process Information ---")
	fmt.Printf("Process ID: %d\n", process.GetProcessID())
	fmt.Printf("Parent Process ID: %d\n", process.GetParentProcessID())
	fmt.Printf("Process Name: %s\n", process.Name())
	fmt.Printf("Process Path: %s\n", process.Path())
	fmt.Printf("Working Directory: %s\n", process.WorkingDirectory())

	// Get argument information
	fmt.Println("\n--- Arguments ---")
	fmt.Printf("Argument Count: %d\n", process.ArgumentCount())
	fmt.Printf("All Arguments: %v\n", process.Arguments())
	fmt.Printf("First Argument (executable): %s\n", process.GetArgument(0))

	// Example: Check for specific arguments
	fmt.Printf("Has '--verbose' argument: %v\n", process.HasArgument("--verbose"))
	fmt.Printf("Argument at index 1: %s\n", process.GetArgument(1))

	// Example: Get argument value
	fmt.Printf("Value for '--test': %s\n", process.GetArgumentValue("--test"))
	fmt.Printf("Value for '--env=': %s\n", process.GetArgumentValue("--env"))

	// Get environment information
	fmt.Println("\n--- Environment Variables ---")
	envVars := process.Environment()
	fmt.Printf("Environment Variable Count: %d\n", len(envVars))

	// Get specific environment variables
	fmt.Printf("HOME: %s\n", process.GetEnvironment("HOME"))
	fmt.Printf("PATH: %s\n", process.GetEnvironment("PATH"))
	fmt.Printf("GOPATH: %s\n", process.GetEnvironment("GOPATH"))

	// Create a Process instance for the current process using explicit PID
	fmt.Println("\n--- Process (Current Process with Explicit PID) ---")
	currentProcess := NewProcess(process.GetProcessID())
	fmt.Printf("Process PID: %d\n", currentProcess.GetProcessID())
	fmt.Printf("Is Current Process: %v\n", currentProcess.IsCurrent())
	fmt.Printf("Process IsExists: %v\n", currentProcess.Exists())

	// Get parent process info
	fmt.Printf("Parent Process PID: %d\n", currentProcess.GetParentProcessID())

	// Example: Check if a command exists
	fmt.Println("\n--- Command Execution ---")
	fmt.Printf("ls command exists: %v\n", currentProcess.CommandExists("ls"))
	fmt.Printf("nonexistentcommand exists: %v\n", currentProcess.CommandExists("nonexistentcommand"))

	// Example: Execute command and wait for completion
	fmt.Println("\nExecuting 'ls -la' command:")
	output, err := currentProcess.ExecuteCommandWithOutput("ls", "-la")
	if err == nil {
		fmt.Println(output)
	} else {
		fmt.Printf("Error executing command: %v\n", err)
	}

	// Example: Execute command without waiting (async)
	fmt.Println("\nExecuting 'sleep 2' asynchronously:")
	cmd, err := currentProcess.ExecuteCommand("sleep", "2")
	if err == nil {
		fmt.Printf("Command started with PID: %d\n", cmd.Process.Pid)

		// Create a Process instance for the child process
		childProcess := NewProcess(cmd.Process.Pid)
		fmt.Printf("Child Process PID: %d\n", childProcess.GetProcessID())
		fmt.Printf("Child Parent PID: %d\n", childProcess.GetParentProcessID())

		// Wait for command to complete
		fmt.Println("Waiting for command to complete...")
		time.Sleep(3 * time.Second)
	}

	// Example: Execute command and wait
	fmt.Println("\nExecuting 'echo Hello, Process!' and waiting:")
	err = currentProcess.ExecuteCommandAndWait("echo", "Hello, Process!")
	if err == nil {
		fmt.Println("Command executed successfully")
	} else {
		fmt.Printf("Error executing command: %v\n", err)
	}

	// Example: Test GetParentProcessID for a specific PID
	fmt.Println("\n--- Testing GetParentProcessID ---")

	// Test with current process
	fmt.Printf("Current process PID: %d, Parent PID: %d\n", currentProcess.GetProcessID(), currentProcess.GetParentProcessID())

	// Test with init/systemd process
	initProcess := NewProcess(1)
	fmt.Printf("Init process PID: %d, Parent PID: %d\n", initProcess.GetProcessID(), initProcess.GetParentProcessID())
	fmt.Println("\n--- Process Examples Complete ---")
}
