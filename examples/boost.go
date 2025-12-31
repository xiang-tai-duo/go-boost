// --------------------------------------------------------------------------------
// File:        boost.go
// Author:      TRAE AI
// Created:     12/30/2025 11:03:46
// Description: Example for BOOST utility functions
// --------------------------------------------------------------------------------

package main

import (
	"fmt"

	. "github.com/xiang-tai-duo/go-boost"
)

func main() {

	// Example: Using the M() function to wrap values
	fmt.Println("--- M() Function Examples ---")

	// Wrap a string path and convert to DIRECTORY
	pathStr := "./test-dir"
	dir := M(pathStr).AsDirectory()
	fmt.Printf("M('%s').AsDirectory(): %s\n", pathStr, dir.Path())
	fmt.Printf("Directory exists: %v\n", dir.Exists())

	// Wrap a string path and convert to FILE
	filePath := "./test.txt"
	file := M(filePath).AsFile()
	fmt.Printf("M('%s').AsFile(): %s\n", filePath, file.Path())
	fmt.Printf("File exists: %v\n", file.Exists())

	// Wrap a string path and convert to FILEPATH
	filePathObj := M(filePath).AsFilePath()
	fmt.Printf("M('%s').AsFilePath(): %s\n", filePath, filePathObj.Path())
	fmt.Printf("Is absolute path: %v\n", filePathObj.IsAbs())
	fmt.Printf("File name: %s\n", filePathObj.GetFileName())
	fmt.Printf("File extension: %s\n", filePathObj.GetFileExtension())

	// Wrap a map and convert to JSON
	userMap := map[string]interface{}{
		"name":  "John Doe",
		"age":   30,
		"email": "john@example.com",
	}
	jsonObj := M(userMap).AsJson()
	jsonStr, _ := jsonObj.Marshal()
	fmt.Printf("M(map).AsJson(): %s\n", jsonStr)

	// Wrap a JSON string and convert to JSON
	jsonStrRaw := `{"product":"Laptop","price":999.99}`
	jsonFromStr := M(jsonStrRaw).AsJson()
	formattedJson, _ := jsonFromStr.Format("  ")
	fmt.Printf("M(jsonStr).AsJson() formatted:\n%s\n", formattedJson)

	// Example: Using global Debugger instance
	fmt.Println("\n--- Global Debugger Instance ---")
	isDebuggerPresent := Debugger.Check()
	fmt.Printf("Global Debugger check: %v\n", isDebuggerPresent)
	fmt.Printf("Is debugger present: %v\n", Debugger.IsPresent())

	// Example: Chaining operations with M()
	fmt.Println("\n--- Chaining Operations ---")

	// Wrap a path, convert to FILEPATH, then get file name without extension
	fileName := M("/path/to/document.pdf").AsFilePath().GetFileNameWithoutExtension()
	fmt.Printf("M('/path/to/document.pdf').AsFilePath().GetFileNameWithoutExtension(): %s\n", fileName)

	// Wrap a path, convert to DIRECTORY, then check if it exists
	isDirExists := M("./nonexistent-dir").AsDirectory().Exists()
	fmt.Printf("M('./nonexistent-dir').AsDirectory().Exists(): %v\n", isDirExists)

	// Example: Using M() with different types
	fmt.Println("\n--- M() with Different Types ---")

	// Wrap an integer (will be converted to string internally)
	intWrapped := M(123)
	intAsDir := intWrapped.AsDirectory()
	fmt.Printf("M(123).AsDirectory().Path(): '%s'\n", intAsDir.Path())

	// Wrap a boolean (will be converted to string internally)
	boolWrapped := M(true)
	boolAsFile := boolWrapped.AsFile()
	fmt.Printf("M(true).AsFile().Path(): '%s'\n", boolAsFile.Path())

	// Example: Comparing direct creation vs M() conversion
	fmt.Println("\n--- Direct vs M() Conversion ---")

	// Direct creation
	directFile := NewFile("./test.txt")

	// Using M() conversion
	wrappedFile := M("./test.txt").AsFile()
	fmt.Printf("Direct File Path: %s\n", directFile.Path())
	fmt.Printf("M() File Path: %s\n", wrappedFile.Path())
	fmt.Printf("Both paths are the same: %v\n", directFile.Path() == wrappedFile.Path())
	fmt.Println("\n--- BOOST Examples Complete ---")
}
