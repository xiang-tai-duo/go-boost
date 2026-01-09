// Package boost
// File:        directory.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/examples/directory.go
// Author:      Vibe Coding
// Created:     12/30/2025 11:03:46
// Description: Example for DIRECTORY utility functions
// --------------------------------------------------------------------------------

package main

import (
	"fmt"

	. "github.com/xiang-tai-duo/go-boost"
)

func main() {

	// Create a new DIRECTORY instance
	dir := NewDirectory("./test-dir")

	// Check if directory exists
	fmt.Printf("Directory '%s' exists: %v\n", dir.Path(), dir.Exists())

	// Create directory
	err := dir.Create()
	if err == nil {
		fmt.Printf("Created directory '%s'\n", dir.Path())
	} else {
		fmt.Printf("Error creating directory: %v\n", err)
	}

	// Check if directory exists again
	fmt.Printf("Directory '%s' exists after create: %v\n", dir.Path(), dir.Exists())

	// Create subdirectory using Join
	subdir := dir.Join("subdir")
	fmt.Printf("Subdirectory path: %s\n", subdir.Path())
	err = subdir.MakeAll() // MakeAll creates parent directories if needed
	if err == nil {
		fmt.Printf("Created subdirectory '%s'\n", subdir.Path())
	} else {
		fmt.Printf("Error creating subdirectory: %v\n", err)
	}

	// Get directory contents
	contents, err := dir.GetAll()
	if err == nil {
		fmt.Printf("Directory contents: %v\n", contents)
	} else {
		fmt.Printf("Error getting directory contents: %v\n", err)
	}

	// Get only directories
	dirs, err := dir.GetDirectories()
	if err == nil {
		fmt.Printf("Directories only: %v\n", dirs)
	} else {
		fmt.Printf("Error getting directories: %v\n", err)
	}

	// Get only files (should be empty for our test dir)
	files, err := dir.GetFiles()
	if err == nil {
		fmt.Printf("Files only: %v\n", files)
	} else {
		fmt.Printf("Error getting files: %v\n", err)
	}

	// Get working directory
	wd, err := dir.GetWorkingDirectory()
	if err == nil {
		fmt.Printf("Current working directory: %s\n", wd.Path())
	} else {
		fmt.Printf("Error getting working directory: %v\n", err)
	}

	// Check if path is absolute
	fmt.Printf("Is '%s' absolute: %v\n", dir.Path(), dir.IsAbs())
	fmt.Printf("Is '%s' absolute: %v\n", wd.Path(), wd.IsAbs())

	// Get base name
	fmt.Printf("Base name of '%s': %s\n", dir.Path(), dir.Base())

	// Get parent directory
	parentDir := dir.Dir()
	fmt.Printf("Parent directory of '%s': %s\n", dir.Path(), parentDir.Path())

	// Convert to slash format
	fmt.Printf("ToSlash: %s\n", dir.ToSlash().Path())

	// Clean path
	dirtyDir := NewDirectory("./test-dir/../test-dir")
	fmt.Printf("Dirty path: %s\n", dirtyDir.Path())
	fmt.Printf("Clean path: %s\n", dirtyDir.Clean().Path())

	// Delete subdirectory
	err = subdir.Delete()
	if err == nil {
		fmt.Printf("Deleted subdirectory '%s'\n", subdir.Path())
	} else {
		fmt.Printf("Error deleting subdirectory: %v\n", err)
	}

	// Delete main directory
	err = dir.Delete()
	if err == nil {
		fmt.Printf("Deleted directory '%s'\n", dir.Path())
	} else {
		fmt.Printf("Error deleting directory: %v\n", err)
	}

	// Check if directory exists after deletion
	fmt.Printf("Directory '%s' exists after delete: %v\n", dir.Path(), dir.Exists())

	// Create a nested directory structure and delete all
	nestedDir := NewDirectory("./nested/dir/structure")
	err = nestedDir.MakeAll()
	if err == nil {
		fmt.Printf("Created nested directory structure '%s'\n", nestedDir.Path())
	}

	// Delete all nested directories
	rootNestedDir := NewDirectory("./nested")
	err = rootNestedDir.DeleteAll()
	if err == nil {
		fmt.Printf("Deleted all nested directories\n")
	} else {
		fmt.Printf("Error deleting nested directories: %v\n", err)
	}
}
