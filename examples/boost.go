// --------------------------------------------------------------------------------
// File:        boost.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: Example for BOOST utility functions
// --------------------------------------------------------------------------------

package main

import (
	"fmt"

	"github.com/xiang-tai-duo/go-boost"
)

func main() {
	// Example: Using the Boost() function to wrap values
	fmt.Println("--- Boost() Function Examples ---")
	// Wrap a string path and convert to DIRECTORY
	pathStr := "./test-dir"
	dir := boost.Boost(pathStr).AsDirectory()
	fmt.Printf("Boost('%s').AsDirectory(): %s\n", pathStr, dir.Path())
	fmt.Printf("Directory exists: %v\n", dir.Exists())
	// Wrap a string path and convert to FILE
	filePath := "./test.txt"
	file := boost.Boost(filePath).AsFile()
	fmt.Printf("Boost('%s').AsFile(): %s\n", filePath, file.Path())
	fmt.Printf("File exists: %v\n", file.Exists())
	// Wrap a string path and convert to FILEPATH
	filePathObj := boost.Boost(filePath).AsFilePath()
	fmt.Printf("Boost('%s').AsFilePath(): %s\n", filePath, filePathObj.Path())
	fmt.Printf("Is absolute path: %v\n", filePathObj.IsAbs())
	fmt.Printf("File name: %s\n", filePathObj.GetFileName())
	fmt.Printf("File extension: %s\n", filePathObj.GetFileExtension())
	// Wrap a map and convert to JSON
	userMap := map[string]interface{}{
		"name":  "John Doe",
		"age":   30,
		"email": "john@example.com",
	}
	jsonObj := boost.Boost(userMap).AsJson()
	jsonStr, _ := jsonObj.Json()
	fmt.Printf("Boost(map).AsJson(): %s\n", jsonStr)
	// Wrap a JSON string and convert to JSON
	jsonStrRaw := `{"product":"Laptop","price":999.99}`
	jsonFromStr := boost.Boost(jsonStrRaw).AsJson()
	formattedJson, _ := jsonFromStr.Format("  ")
	fmt.Printf("Boost(jsonStr).AsJson() formatted:\n%s\n", formattedJson)

	// Example: Using global Debugger instance
	fmt.Println("\n--- Global Debugger Instance ---")
	isDebuggerPresent := boost.Debugger.Check()
	fmt.Printf("Global Debugger check: %v\n", isDebuggerPresent)
	fmt.Printf("Is debugger present: %v\n", boost.Debugger.IsPresent())

	// Example: Chaining operations with Boost()
	fmt.Println("\n--- Chaining Operations ---")
	// Wrap a path, convert to FILEPATH, then get file name without extension
	fileName := boost.Boost("/path/to/document.pdf").AsFilePath().GetFileNameWithoutExtension()
	fmt.Printf("Boost('/path/to/document.pdf').AsFilePath().GetFileNameWithoutExtension(): %s\n", fileName)
	// Wrap a path, convert to DIRECTORY, then check if it exists
	isDirExists := boost.Boost("./nonexistent-dir").AsDirectory().Exists()
	fmt.Printf("Boost('./nonexistent-dir').AsDirectory().Exists(): %v\n", isDirExists)

	// Example: Using Boost() with different types
	fmt.Println("\n--- Boost() with Different Types ---")
	// Wrap an integer (will be converted to string internally)
	intWrapped := boost.Boost(123)
	intAsDir := intWrapped.AsDirectory()
	fmt.Printf("Boost(123).AsDirectory().Path(): '%s'\n", intAsDir.Path())
	// Wrap a boolean (will be converted to string internally)
	boolWrapped := boost.Boost(true)
	boolAsFile := boolWrapped.AsFile()
	fmt.Printf("Boost(true).AsFile().Path(): '%s'\n", boolAsFile.Path())

	// Example: Comparing direct creation vs Boost() conversion
	fmt.Println("\n--- Direct vs Boost() Conversion ---")
	// Direct creation
	directFile := boost.NewFile("./test.txt")
	// Using Boost() conversion
	wrappedFile := boost.Boost("./test.txt").AsFile()

	fmt.Printf("Direct File Path: %s\n", directFile.Path())
	fmt.Printf("Boost() File Path: %s\n", wrappedFile.Path())
	fmt.Printf("Both paths are the same: %v\n", directFile.Path() == wrappedFile.Path())

	fmt.Println("\n--- BOOST Examples Complete ---")
}
