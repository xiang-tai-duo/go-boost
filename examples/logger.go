// --------------------------------------------------------------------------------
// File:        logger.go
// Author:      TRAE AI
// Created:     12/30/2025 11:03:46
// Description: Example for LOGGER utility functions
// --------------------------------------------------------------------------------

package main

import (
	"fmt"
	"time"

	. "github.com/xiang-tai-duo/go-boost"
)

func main() {

	// Example: Using the global Logger instance
	fmt.Println("--- Logger Example ---")

	// Log different levels of messages
	fmt.Println("\n1. Logging different levels:")

	// Debug level (will not show if current level is Info)
	Logger.Debug("This is a debug message")

	// Info level
	Logger.Info("This is an info message")

	// Warning level
	Logger.Warning("This is a warning message")

	// Error level with string
	Logger.Error("This is an error message")

	// Error level with error object
	Logger.Error(fmt.Errorf("This is an error object"))

	// Secret level (will be masked in console)
	Logger.Secret("This is a secret message that should not be exposed in plain text")

	// Example: Changing log level
	fmt.Println("\n2. Changing log level:")

	// Set log level to Debug to see debug messages
	fmt.Println("Setting log level to Debug...")
	Logger.CurrentLevel = LOG_LEVEL_DEBUG

	// Now debug messages will be shown
	Logger.Debug("This debug message should now be visible")

	// Example: Logging with context
	fmt.Println("\n3. Logging with context:")

	// Log messages from different functions
	logFromFunction()

	// Example: Logging in a goroutine
	fmt.Println("\n4. Logging from goroutines:")
	for i := 0; i < 3; i++ {
		i := i
		go func() {
			Logger.Info(fmt.Sprintf("This is a log from goroutine %d", i))
		}()
	}

	// Give goroutines time to complete
	time.Sleep(100 * time.Millisecond)

	// Example: Demonstrate logger functionality
	fmt.Println("\n5. Logger functionality:")

	// Log structured data
	Logger.Info(fmt.Sprintf("User login attempt: username='testuser', ip='192.168.1.1', timestamp='%v'", time.Now()))

	// Log with formatted strings
	Logger.Warning(fmt.Sprintf("Disk usage at %.2f%%", 85.5))

	// Log errors with stack trace information
	Logger.Error("Critical error: database connection failed")

	// Example: Secret logging
	fmt.Println("\n6. Secret logging:")

	// Log multiple secret messages
	Logger.Secret("API Key: sk_1234567890abcdef")

	// Log multiple secret messages
	Logger.Secret("Database Password: mysecretpassword")

	// Log multiple secret messages
	Logger.Secret("JWT Secret: myjwtsecretkey")

	// Note: The decrypt functionality would require a log file path,
	// but for this example we'll just demonstrate the API usage
	fmt.Println("\n7. Logger API usage:")
	fmt.Println("- Logger.Debug(message string)")
	fmt.Println("- Logger.Info(message string)")
	fmt.Println("- Logger.Warning(message string)")
	fmt.Println("- Logger.Error(message interface{}, skipStackFrames ...int)")
	fmt.Println("- Logger.Secret(message string)")
	fmt.Println("- Logger.DecryptSecretLogs(logFilePath string) ([]string, error)")

	// Example: Reset log level
	Logger.CurrentLevel = LOG_LEVEL_INFO
	fmt.Println("\n8. Reset log level to Info")
	Logger.Debug("This debug message should not be visible anymore")
	Logger.Info("This info message should still be visible")
	fmt.Println("\n--- Logger Examples Complete ---")
	fmt.Println("Check the Logs directory for generated log files")
}

// Helper function to demonstrate logging from different functions

func logFromFunction() {
	Logger.Info("This log is from the logFromFunction() helper function")
	Logger.Debug("This debug log is also from logFromFunction()")
}
