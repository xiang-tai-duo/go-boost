// Package boost
// File:        filepath.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/examples/filepath.go
// Author:      Vibe Coding
// Created:     12/30/2025 11:03:46
// Description: Example for FILEPATH utility functions
// --------------------------------------------------------------------------------

package main

import (
	"fmt"
	"os"

	. "github.com/xiang-tai-duo/go-boost"
)

func main() {

	fp := NewFilePath("/path/to/file.txt")

	fmt.Println("Path:", fp.Path())

	fp.Path("/new/path/to/file.txt")
	fmt.Println("Updated Path:", fp.Path())

	absPath, err := fp.Abs()
	if err == nil {
		fmt.Println("Absolute Path:", absPath.Path())
	}

	fmt.Println("Base Name:", fp.Base())

	fpDirty := NewFilePath("/path/../to/./file.txt")
	cleanPath := fpDirty.Clean()
	fmt.Println("Dirty Path:", fpDirty.Path())
	fmt.Println("Clean Path:", cleanPath.Path())

	fmt.Println("Directory:", fp.Dir())

	fmt.Println("Extension:", fp.Ext())
	fmt.Println("GetFileExtension:", fp.GetFileExtension())

	fmt.Println("GetFileName:", fp.GetFileName())

	fmt.Println("GetFileNameWithoutExtension:", fp.GetFileNameWithoutExtension())

	fmt.Println("GetParentDirectory:", fp.GetParentDirectory())

	fmt.Println("IsAbs:", fp.IsAbs())

	joinedPath := fp.Join("subdir", "anotherfile.txt")
	fmt.Println("Joined Path:", joinedPath.Path())

	dir, file := fp.Split()
	fmt.Println("Split - Dir:", dir, "File:", file)

	fmt.Println("VolumeName:", fp.VolumeName())

	toSlashResult := fp.ToSlash()

	fmt.Println("ToSlash:", (&toSlashResult).Path())

	fpSlash := NewFilePath("/path/to/file.txt")
	fromSlashResult := fpSlash.FromSlash()
	fmt.Println("FromSlash:", (&fromSlashResult).Path())

	fmt.Println("HasPrefix '/new':", fp.HasPrefix("/new"))

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
