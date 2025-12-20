// --------------------------------------------------------------------------------
// File:        file.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: Example for FILE utility functions
// --------------------------------------------------------------------------------

package main

import (
	"fmt"
	"time"

	. "github.com/xiang-tai-duo/go-boost"
)

func main() {
	// Create a new FILE instance
	file := NewFile("./test.txt")

	// Check if file exists
	fmt.Printf("File '%s' exists: %v\n", file.Path(), file.Exists())
	// Write content to file
	err := file.WriteString("Hello, Go Boost!\n")
	if err == nil {
		fmt.Printf("Wrote content to file '%s'\n", file.Path())
	} else {
		fmt.Printf("Error writing to file: %v\n", err)
	}
	// Append content to file
	err = file.AppendString("Appended line 1\n")
	if err == nil {
		fmt.Printf("Appended content to file '%s'\n", file.Path())
	}
	// Read file content as string
	content, err := file.ReadString()
	if err == nil {
		fmt.Printf("File content:\n%s\n", content)
	}
	// Read file content as bytes
	data, err := file.ReadAll()
	if err == nil {
		fmt.Printf("File size: %d bytes\n", len(data))
	}
	// Get file information
	fmt.Printf("File name: %s\n", file.Name())
	fmt.Printf("File extension: %s\n", file.Ext())
	fmt.Printf("File size: %d bytes\n", file.Size())
	fmt.Printf("File mode: %v\n", file.Mode())
	fmt.Printf("Last modified: %v\n", file.ModTime())
	fmt.Printf("Is absolute path: %v\n", file.IsAbs())
	fmt.Printf("Is directory: %v\n", file.IsDir())
	fmt.Printf("Is regular file: %v\n", file.IsRegular())
	fmt.Printf("Is symlink: %v\n", file.IsSymlink())
	// Get file directory
	dir := file.Dir()
	fmt.Printf("File directory: %s\n", dir.Path())
	// Copy file
	copyFile := NewFile("./test_copy.txt")
	err = file.Copy(*copyFile)
	if err == nil {
		fmt.Printf("Copied file to '%s'\n", copyFile.Path())
	}
	// Check copy exists
	fmt.Printf("Copy file exists: %v\n", copyFile.Exists())
	// Rename file
	renameFile := NewFile("./test_renamed.txt")
	err = file.Rename(*renameFile)
	if err == nil {
		fmt.Printf("Renamed file to '%s'\n", renameFile.Path())
	}
	// Check original file no longer exists
	fmt.Printf("Original file exists after rename: %v\n", file.Exists())
	// Check renamed file exists
	fmt.Printf("Renamed file exists: %v\n", renameFile.Exists())
	// Change file permissions
	err = renameFile.Chmod(0666)
	if err == nil {
		fmt.Printf("Changed file permissions to 0666\n")
	}
	// Change file times
	now := time.Now()
	err = renameFile.Chtimes(now, now)
	if err == nil {
		fmt.Printf("Updated file times to current time\n")
	}
	// Delete files
	err = renameFile.Delete()
	if err == nil {
		fmt.Printf("Deleted file '%s'\n", renameFile.Path())
	}
	err = copyFile.Delete()
	if err == nil {
		fmt.Printf("Deleted file '%s'\n", copyFile.Path())
	}
	// Check files don't exist
	fmt.Printf("Renamed file exists after delete: %v\n", renameFile.Exists())
	fmt.Printf("Copy file exists after delete: %v\n", copyFile.Exists())
	// Create nested file
	nestedFile := NewFile("./nested/test.txt")
	// Note: This will fail because the directory doesn't exist
	err = nestedFile.WriteString("This won't work without parent directory\n")
	if err != nil {
		fmt.Printf("Expected error writing to nested file: %v\n", err)
	}
}
