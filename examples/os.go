// Package boost
// File:        os.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/examples/os.go
// Author:      Vibe Coding
// Created:     12/30/2025 11:03:46
// Description: Example for OS utility functions
// --------------------------------------------------------------------------------

package main

import (
	"fmt"

	. "github.com/xiang-tai-duo/go-boost"
)

func main() {

	os := OS{}

	fmt.Println("--- OS Information ---")
	fmt.Printf("OS Name: %s\n", os.Name())
	fmt.Printf("OS Version: %s\n", os.Version())
	fmt.Printf("Architecture: %s\n", os.Architecture())

	fmt.Println("\n--- OS Type Checks ---")
	fmt.Printf("Is Linux: %v\n", os.IsLinux())
	fmt.Printf("Is Mac: %v\n", os.IsMac())
	fmt.Printf("Is Windows: %v\n", os.IsWindows())
	fmt.Printf("Is Unix: %v\n", os.IsUnix())

	fmt.Println("\n--- Separators ---")
	fmt.Printf("Path Separator: %q\n", os.PathSeparator())
	fmt.Printf("Environment Separator: %q\n", os.Separator())

	fmt.Printf("Executable Path: %s\n", os.Executable())

	hostname, err := os.Hostname()
	if err == nil {
		fmt.Printf("Hostname: %s\n", hostname)
	}

	fmt.Println("\n--- Directories ---")
	fmt.Printf("Temporary Directory: %s\n", os.TemporaryDirectory())
	fmt.Printf("User Home Directory: %s\n", os.UserHomeDirectory())

	fmt.Println("\n--- Environment Variables ---")
	envVars := os.Environment()
	fmt.Printf("Environment Variable Count: %d\n", len(envVars))

	fmt.Printf("HOME: %s\n", os.GetEnvironment("HOME"))
	fmt.Printf("PATH: %s\n", os.GetEnvironment("PATH"))
	fmt.Printf("GOPATH: %s\n", os.GetEnvironment("GOPATH"))

	fmt.Printf("Has GOPATH: %v\n", os.HasEnvironment("GOPATH"))
	fmt.Printf("Has NONEXISTENT_VAR: %v\n", os.HasEnvironment("NONEXISTENT_VAR"))

	testStr := "Home: $HOME, Path: $PATH"
	fmt.Printf("ExpandEnv(%q): %q\n", testStr, os.ExpandEnvironment(testStr))

	fmt.Println("\n--- Environment Variable Operations ---")
	testVarKey := "BOOST_TEST_VAR"
	testVarValue := "test_value"

	err = os.SetEnvironment(testVarKey, testVarValue)
	if err == nil {
		fmt.Printf("Set %s=%s\n", testVarKey, testVarValue)
		fmt.Printf("Get %s: %s\n", testVarKey, os.GetEnvironment(testVarKey))
	}

	err = os.UnsetEnvironment(testVarKey)
	if err == nil {
		fmt.Printf("Unset %s\n", testVarKey)
		fmt.Printf("Get %s after unset: %s\n", testVarKey, os.GetEnvironment(testVarKey))
	}

	fmt.Println("\n--- Executable Lookup ---")
	lsPath, err := os.LookPath("ls")
	if err == nil {
		fmt.Printf("'ls' executable found at: %s\n", lsPath)
	} else {
		fmt.Printf("'ls' executable not found: %v\n", err)
	}

	nonexistentPath, err := os.LookPath("nonexistent_executable")
	if err == nil {
		fmt.Printf("'nonexistent_executable' found at: %s\n", nonexistentPath)
	} else {
		fmt.Printf("'nonexistent_executable' not found: %v\n", err)
	}

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

	fmt.Println("\n--- OS Examples Complete ---")
	fmt.Println("Note: os.Exit() is not demonstrated to avoid program termination")
	fmt.Println("Note: os.Kill() is not demonstrated to avoid killing processes")
}
