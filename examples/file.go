// Package boost
// File:        file.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/examples/file.go
// Author:      Vibe Coding
// Created:     12/30/2025 11:03:46
// Description: Example for FILE utility functions
// --------------------------------------------------------------------------------

package main

import (
	"fmt"
	"time"

	. "github.com/xiang-tai-duo/go-boost"
)

func main() {

	file := NewFile("./test.txt")

	fmt.Printf("File '%s' exists: %v\n", file.Path(), file.IsExists())

	err := file.WriteString("Hello, Go _!\n")
	if err == nil {
		fmt.Printf("Wrote content to file '%s'\n", file.Path())
	} else {
		fmt.Printf("Error writing to file: %v\n", err)
	}

	err = file.AppendString("Appended line 1\n")
	if err == nil {
		fmt.Printf("Appended content to file '%s'\n", file.Path())
	}

	content, err := file.ReadString()
	if err == nil {
		fmt.Printf("File content:\n%s\n", content)
	}

	data, err := file.ReadAll()
	if err == nil {
		fmt.Printf("File size: %d bytes\n", len(data))
	}

	fmt.Printf("File name: %s\n", file.Name())
	fmt.Printf("File extension: %s\n", file.Ext())
	fmt.Printf("File size: %d bytes\n", file.Size())
	fmt.Printf("File mode: %v\n", file.Mode())
	fmt.Printf("Last modified: %v\n", file.ModTime())
	fmt.Printf("Is absolute path: %v\n", file.IsAbs())
	fmt.Printf("Is directory: %v\n", file.IsDir())
	fmt.Printf("Is regular file: %v\n", file.IsRegular())
	fmt.Printf("Is symlink: %v\n", file.IsSymlink())

	dir := file.Dir()
	fmt.Printf("File directory: %s\n", dir.Path())

	copyFile := NewFile("./test_copy.txt")
	err = file.Copy(*copyFile)
	if err == nil {
		fmt.Printf("Copied file to '%s'\n", copyFile.Path())
	}

	fmt.Printf("Copy file exists: %v\n", copyFile.IsExists())

	renameFile := NewFile("./test_renamed.txt")
	err = file.Rename(*renameFile)
	if err == nil {
		fmt.Printf("Renamed file to '%s'\n", renameFile.Path())
	}

	fmt.Printf("Original file exists after rename: %v\n", file.IsExists())

	fmt.Printf("Renamed file exists: %v\n", renameFile.IsExists())

	err = renameFile.Chmod(0666)
	if err == nil {
		fmt.Printf("Changed file permissions to 0666\n")
	}

	now := time.Now()
	err = renameFile.Chtimes(now, now)
	if err == nil {
		fmt.Printf("Updated file times to current time\n")
	}

	err = renameFile.Delete()
	if err == nil {
		fmt.Printf("Deleted file '%s'\n", renameFile.Path())
	}
	err = copyFile.Delete()
	if err == nil {
		fmt.Printf("Deleted file '%s'\n", copyFile.Path())
	}

	fmt.Printf("Renamed file exists after delete: %v\n", renameFile.IsExists())
	fmt.Printf("Copy file exists after delete: %v\n", copyFile.IsExists())

	nestedFile := NewFile("./nested/test.txt")

	err = nestedFile.WriteString("This won't work without parent directory\n")
	if err != nil {
		fmt.Printf("Expected error writing to nested file: %v\n", err)
	}
}
