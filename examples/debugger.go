// Package boost
// File:        debugger.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/examples/debugger.go
// Author:      Vibe Coding
// Created:     12/30/2025 11:03:46
// Description: Example for DEBUGGER utility functions
// --------------------------------------------------------------------------------

package main

import (
	"fmt"
	"os"

	. "github.com/xiang-tai-duo/go-boost"
)

func main() {

	// Create a new DEBUGGER instance
	debugger := NewDebugger()

	// Check for debugger presence
	isDebuggerPresent := debugger.Check()
	fmt.Printf("Debugger check result: %v\n", isDebuggerPresent)

	// Get debugger presence status
	isPresent := debugger.IsPresent()
	fmt.Printf("Is debugger present: %v\n", isPresent)

	// Show current debugger state
	if isPresent {
		fmt.Println("A debugger is attached to this process!")
	} else {
		fmt.Println("No debugger detected.")
	}

	// Example: Using debugger detection to modify behavior
	if isPresent {
		// In debug mode: enable verbose logging
		fmt.Println("[DEBUG MODE] Verbose logging enabled")
		fmt.Println("[DEBUG MODE] Process ID:", os.Getpid())
		fmt.Println("[DEBUG MODE] Parent Process ID:", os.Getppid())
	} else {

		// In release mode: normal operation
		fmt.Println("Running in release mode")
	}
}
