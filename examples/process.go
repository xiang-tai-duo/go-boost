// --------------------------------------------------------------------------------
// File:        process.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: Example for Process and ProcessBuilder utility functions
// --------------------------------------------------------------------------------

package main

import (
	"fmt"
	"time"

	. "github.com/xiang-tai-duo/go-boost"
)

func main() {
	// Create a new Process instance
	process := NewProcess()

	// Get process information
	fmt.Println("--- Process Information ---")
	fmt.Printf("Process ID: %d\n", process.ProcessID())
	fmt.Printf("Parent Process ID: %d\n", process.ParentProcessID())
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

	// Create a ProcessBuilder for the current process
	fmt.Println("\n--- ProcessBuilder (Current Process) ---")
	currentProcess := NewProcessBuilder(process.ProcessID())
	fmt.Printf("ProcessBuilder PID: %d\n", currentProcess.ProcessID())
	fmt.Printf("Is Current Process: %v\n", currentProcess.IsCurrent())
	fmt.Printf("Process Exists: %v\n", currentProcess.Exists())
	// Get parent process info
	parentProcess := currentProcess.Parent()
	fmt.Printf("Parent Process PID: %d\n", parentProcess.ProcessID())

	// Example: Check if a command exists
	fmt.Println("\n--- Command Execution ---")
	processBuilder := NewProcessBuilder(0)
	fmt.Printf("ls command exists: %v\n", processBuilder.CommandExists("ls"))
	fmt.Printf("nonexistentcommand exists: %v\n", processBuilder.CommandExists("nonexistentcommand"))

	// Example: Execute command and wait for completion
	fmt.Println("\nExecuting 'ls -la' command:")
	output, err := processBuilder.ExecuteCommandWithOutput("ls", "-la")
	if err == nil {
		fmt.Println(output)
	} else {
		fmt.Printf("Error executing command: %v\n", err)
	}

	// Example: Execute command without waiting (async)
	fmt.Println("\nExecuting 'sleep 2' asynchronously:")
	cmd, err := processBuilder.ExecuteCommand("sleep", "2")
	if err == nil {
		fmt.Printf("Command started with PID: %d\n", cmd.Process.Pid)
		// Wait for command to complete
		fmt.Println("Waiting for command to complete...")
		time.Sleep(3 * time.Second)
	}

	// Example: Execute command and wait
	fmt.Println("\nExecuting 'echo Hello, Process!' and waiting:")
	err = processBuilder.ExecuteCommandAndWait("echo", "Hello, Process!")
	if err == nil {
		fmt.Println("Command executed successfully")
	} else {
		fmt.Printf("Error executing command: %v\n", err)
	}

	// Example: Handle signals (for demonstration purposes)
	fmt.Println("\n--- Signal Handling (Demo) ---")
	fmt.Println("Creating a ProcessBuilder for PID 1 (init/systemd)")
	systemProcess := NewProcessBuilder(1)
	fmt.Printf("System process exists: %v\n", systemProcess.Exists())

	// Example: Show how to use Signal method (commented out to avoid actual signals)
	// fmt.Println("Sending SIGUSR1 to current process (commented out)")
	// err = currentProcess.Signal(syscall.SIGUSR1)
	// if err == nil {
	//     fmt.Println("Signal sent successfully")
	// } else {
	//     fmt.Printf("Error sending signal: %v\n", err)
	// }

	fmt.Println("\n--- Process Examples Complete ---")
}
