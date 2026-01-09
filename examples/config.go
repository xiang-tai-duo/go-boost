// Package boost
// File:        config.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/examples/config.go
// Author:      Vibe Coding
// Created:     12/30/2025 11:03:46
// Description: Example for CONFIG utility functions
// --------------------------------------------------------------------------------

package main

import (
	"fmt"

	. "github.com/xiang-tai-duo/go-boost"
)

func main() {

	// Create a new CONFIG instance
	config := NewConfig()

	// Example: Set various types of values
	fmt.Println("--- Setting Config Values ---")

	// Set string value
	err := config.SetString("app_name", "Go _ App")
	if err == nil {
		fmt.Println("Set app_name = 'Go _ App'")
	}

	// Set integer value
	err = config.SetInteger("version", 1)
	if err == nil {
		fmt.Println("Set version = 1")
	}

	// Set boolean value
	err = config.SetBoolean("debug", true)
	if err == nil {
		fmt.Println("Set debug = true")
	}

	// Set float value
	err = config.SetFloat("pi", 3.14159)
	if err == nil {
		fmt.Println("Set pi = 3.14159")
	}

	// Set string slice
	err = config.SetStrings("allowed_hosts", []string{"localhost", "example.com"})
	if err == nil {
		fmt.Println("Set allowed_hosts = [localhost, example.com]")
	}

	// Example: Get values
	fmt.Println("\n--- Getting Config Values ---")

	// Get string value
	appName, exists, err := config.GetString("app_name")
	if err == nil && exists {
		fmt.Printf("app_name: %s (exists: %v)\n", appName, exists)
	}

	// Get integer value
	version, exists, err := config.GetInteger("version")
	if err == nil && exists {
		fmt.Printf("version: %d (exists: %v)\n", version, exists)
	}

	// Get boolean value
	debug, exists, err := config.GetBoolean("debug")
	if err == nil && exists {
		fmt.Printf("debug: %v (exists: %v)\n", debug, exists)
	}

	// Get float value
	pi, exists, err := config.GetFloat("pi")
	if err == nil && exists {
		fmt.Printf("pi: %.5f (exists: %v)\n", pi, exists)
	}

	// Get string slice
	allowedHosts, exists, err := config.GetStringSlice("allowed_hosts")
	if err == nil && exists {
		fmt.Printf("allowed_hosts: %v (exists: %v)\n", allowedHosts, exists)
	}

	// Example: Check if key exists
	fmt.Println("\n--- Checking Key Existence ---")
	exists, _ = config.Exists("app_name")
	fmt.Printf("app_name exists: %v\n", exists)
	exists, _ = config.Exists("nonexistent_key")
	fmt.Printf("nonexistent_key exists: %v\n", exists)

	// Example: Get non-existent key
	fmt.Println("\n--- Getting Non-existent Key ---")
	value, exists, err := config.GetString("nonexistent_key")
	if err == nil {
		fmt.Printf("nonexistent_key: %s (exists: %v)\n", value, exists)
	}

	// Example: Get all values
	fmt.Println("\n--- Getting All Values ---")
	allValues, err := config.GetAll()
	if err == nil {
		fmt.Printf("All values: %v\n", allValues)
	}

	// Example: Save config to file
	fmt.Println("\n--- Saving Config to File ---")
	configPath := "./config.json"
	err = config.Save(configPath)
	if err == nil {
		fmt.Printf("Config saved to %s\n", configPath)
	} else {
		fmt.Printf("Error saving config: %v\n", err)
	}

	// Example: Create a new config instance and load from file
	fmt.Println("\n--- Loading Config from File ---")
	newConfig := NewConfig()
	err = newConfig.Load(configPath)
	if err == nil {
		fmt.Printf("Config loaded from %s\n", configPath)

		// Verify loaded values
		appName, _, _ := newConfig.GetString("app_name")
		version, _, _ := newConfig.GetInteger("version")
		fmt.Printf("Loaded app_name: %s, version: %d\n", appName, version)
	} else {
		fmt.Printf("Error loading config: %v\n", err)
	}

	// Example: Delete a key
	fmt.Println("\n--- Deleting a Key ---")
	err = config.Delete("debug")
	if err == nil {
		fmt.Println("Deleted key 'debug'")
		exists, _ = config.Exists("debug")
		fmt.Printf("debug exists after delete: %v\n", exists)
	}

	// Example: Clear all config
	fmt.Println("\n--- Clearing All Config ---")
	err = config.Clear()
	if err == nil {
		fmt.Println("Config cleared")
		allValues, _ := config.GetAll()
		fmt.Printf("All values after clear: %v\n", allValues)
	}

	// Example: Set a nested map (using Set directly)
	fmt.Println("\n--- Setting Nested Values ---")
	nestedConfig := NewConfig()
	appConfig := map[string]interface{}{
		"name":    "Test App",
		"version": 2,
		"database": map[string]interface{}{
			"host": "localhost",
			"port": 3306,
			"name": "test_db",
		},
	}
	nestedConfig.Set("app", appConfig)

	// Get the nested map
	appValue, exists, _ := nestedConfig.Get("app")
	if exists {
		fmt.Printf("app config: %v\n", appValue)
	}
	fmt.Println("\n--- CONFIG Examples Complete ---")
}
