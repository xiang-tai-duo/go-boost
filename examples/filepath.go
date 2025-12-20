// --------------------------------------------------------------------------------
// File:        filepath.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: Example for FILEPATH utility functions
// --------------------------------------------------------------------------------

package main

import (
	"fmt"
	"os"

	. "github.com/xiang-tai-duo/go-boost"
)

func main() {
	// Create a new FILEPATH instance
	fp := NewFilePath("/path/to/file.txt")

	// Get the path
	fmt.Println("Path:", fp.Path())
	// Change the path
	fp.Path("/new/path/to/file.txt")
	fmt.Println("Updated Path:", fp.Path())
	// Get absolute path
	absPath, err := fp.Abs()
	if err == nil {
		fmt.Println("Absolute Path:", absPath.Path())
	}
	// Get base name
	fmt.Println("Base Name:", fp.Base())
	// Clean the path
	fpDirty := NewFilePath("/path/../to/./file.txt")
	cleanPath := fpDirty.Clean()
	fmt.Println("Dirty Path:", fpDirty.Path())
	fmt.Println("Clean Path:", cleanPath.Path())
	// Get directory
	fmt.Println("Directory:", fp.Dir())
	// Get file extension
	fmt.Println("Extension:", fp.Ext())
	fmt.Println("GetFileExtension:", fp.GetFileExtension())
	// Get file name
	fmt.Println("GetFileName:", fp.GetFileName())
	// Get file name without extension
	fmt.Println("GetFileNameWithoutExtension:", fp.GetFileNameWithoutExtension())
	// Get parent directory
	fmt.Println("GetParentDirectory:", fp.GetParentDirectory())
	// Check if absolute path
	fmt.Println("IsAbs:", fp.IsAbs())
	// Join paths
	joinedPath := fp.Join("subdir", "anotherfile.txt")
	fmt.Println("Joined Path:", joinedPath.Path())
	// Split path
	dir, file := fp.Split()
	fmt.Println("Split - Dir:", dir, "File:", file)
	// Get volume name (Windows only, returns empty on Unix)
	fmt.Println("VolumeName:", fp.VolumeName())
	// Convert to slash format
	toSlashResult := fp.ToSlash()
	// Create a pointer to the FILEPATH value
	fmt.Println("ToSlash:", (&toSlashResult).Path())
	// Convert from slash format
	fpSlash := NewFilePath("/path/to/file.txt")
	fromSlashResult := fpSlash.FromSlash()
	fmt.Println("FromSlash:", (&fromSlashResult).Path())
	// Check has prefix
	fmt.Println("HasPrefix '/new':", fp.HasPrefix("/new"))

	// Walk example
	fmt.Println("\n--- Walking Directory ---")
	walkPath := NewFilePath(".")
	err = walkPath.Walk(func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", path)
		return nil
	})
	if err != nil {
		fmt.Println("Walk error:", err)
	}
}
