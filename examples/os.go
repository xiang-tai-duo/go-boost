// --------------------------------------------------------------------------------
// File:        os.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: Example for OS utility functions
// --------------------------------------------------------------------------------

package main

import (
	"fmt"

	. "github.com/xiang-tai-duo/go-boost"
)

func main() {
	// Create a new OS instance
	os := OS{}

	// Get OS information
	fmt.Println("--- OS Information ---")
	fmt.Printf("OS Name: %s\n", os.Name())
	fmt.Printf("OS Version: %s\n", os.Version())
	fmt.Printf("Architecture: %s\n", os.Architecture())

	// Check OS type
	fmt.Println("\n--- OS Type Checks ---")
	fmt.Printf("Is Linux: %v\n", os.IsLinux())
	fmt.Printf("Is Mac: %v\n", os.IsMac())
	fmt.Printf("Is Windows: %v\n", os.IsWindows())
	fmt.Printf("Is Unix: %v\n", os.IsUnix())

	// Get separators
	fmt.Println("\n--- Separators ---")
	fmt.Printf("Line Separator: %q\n", os.LineSeparator())
	fmt.Printf("Path Separator: %q\n", os.PathSeparator())
	fmt.Printf("Environment Separator: %q\n", os.Separator())

	// Get process information
	fmt.Println("\n--- Process Information ---")
	fmt.Printf("Process ID: %d\n", os.GetProcessID())
	fmt.Printf("Parent Process ID: %d\n", os.GetParentProcessID())
	// Get executable information
	fmt.Printf("Executable Path: %s\n", os.Executable())

	// Get hostname
	hostname, err := os.Hostname()
	if err == nil {
		fmt.Printf("Hostname: %s\n", hostname)
	}

	// Get directories
	fmt.Println("\n--- Directories ---")
	fmt.Printf("Temporary Directory: %s\n", os.TemporaryDirectory())
	fmt.Printf("User Home Directory: %s\n", os.UserHomeDirectory())

	// Environment variables
	fmt.Println("\n--- Environment Variables ---")
	envVars := os.Environment()
	fmt.Printf("Environment Variable Count: %d\n", len(envVars))
	// Get environment variable map
	envMap := os.GetEnvironmentMap()
	fmt.Printf("Environment Map Keys Count: %d\n", len(envMap))
	// Get specific environment variables
	fmt.Printf("HOME: %s\n", os.GetEnvironment("HOME"))
	fmt.Printf("PATH: %s\n", os.GetEnvironment("PATH"))
	fmt.Printf("GOPATH: %s\n", os.GetEnvironment("GOPATH"))
	// Check if environment variable exists
	fmt.Printf("Has GOPATH: %v\n", os.HasEnvironment("GOPATH"))
	fmt.Printf("Has NONEXISTENT_VAR: %v\n", os.HasEnvironment("NONEXISTENT_VAR"))
	// Expand environment variables
	testStr := "Home: $HOME, Path: $PATH"
	fmt.Printf("ExpandEnv(%q): %q\n", testStr, os.ExpandEnvironment(testStr))

	// Example: Set and unset environment variable
	fmt.Println("\n--- Environment Variable Operations ---")
	testVarKey := "BOOST_TEST_VAR"
	testVarValue := "test_value"
	// Set environment variable
	err = os.SetEnvironment(testVarKey, testVarValue)
	if err == nil {
		fmt.Printf("Set %s=%s\n", testVarKey, testVarValue)
		fmt.Printf("Get %s: %s\n", testVarKey, os.GetEnvironment(testVarKey))
	}
	// Unset environment variable
	err = os.UnsetEnvironment(testVarKey)
	if err == nil {
		fmt.Printf("Unset %s\n", testVarKey)
		fmt.Printf("Get %s after unset: %s\n", testVarKey, os.GetEnvironment(testVarKey))
	}

	// Example: Look for executable
	fmt.Println("\n--- Executable Lookup ---")
	lsPath, err := os.LookPath("ls")
	if err == nil {
		fmt.Printf("'ls' executable found at: %s\n", lsPath)
	} else {
		fmt.Printf("'ls' executable not found: %v\n", err)
	}
	// Example: Look for non-existent executable
	nonexistentPath, err := os.LookPath("nonexistent_executable")
	if err == nil {
		fmt.Printf("'nonexistent_executable' found at: %s\n", nonexistentPath)
	} else {
		fmt.Printf("'nonexistent_executable' not found: %v\n", err)
	}

	// Example: OS type specific behavior
	fmt.Println("\n--- OS Specific Behavior ---")
	if os.IsWindows() {
		fmt.Println("Running on Windows - using Windows-specific behavior")
	} else if os.IsMac() {
		fmt.Println("Running on macOS - using macOS-specific behavior")
	} else if os.IsLinux() {
		fmt.Println("Running on Linux - using Linux-specific behavior")
	} else {
		fmt.Println("Running on unknown OS - using generic behavior")
	}

	// Note: os.Exit() is not called in this example to avoid terminating the program prematurely
	fmt.Println("\n--- OS Examples Complete ---")
	fmt.Println("Note: os.Exit() is not demonstrated to avoid program termination")
	fmt.Println("Note: os.Kill() is not demonstrated to avoid killing processes")
}
