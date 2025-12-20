// --------------------------------------------------------------------------------
// File:        directory.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: Directory handles operations related to the current directory,
//              while DirectoryBuilder is used for operations on specific directories.
// --------------------------------------------------------------------------------

package boost

import (
	"os"
	"path/filepath"
	"sort"
)

// DIRECTORY provides utility methods for directory operations.
type DIRECTORY struct {
	path string
}

// NewDirectory creates a new DIRECTORY instance with the given path.
// path: Optional directory path to use
// Returns: New DIRECTORY instance
// Usage:
// dir := NewDirectory() // creates with empty path
// dir := NewDirectory("/path/to/dir") // creates with specified path
func NewDirectory(path ...string) *DIRECTORY {
	var dirPath string
	if len(path) > 0 {
		dirPath = path[0]
	}
	return &DIRECTORY{path: dirPath}
}

// FromPath creates a new DIRECTORY instance with the given path.
// path: Directory path to use
// Returns: New DIRECTORY instance
// Usage:
// dir := DIRECTORY{}.FromPath("/path/to/dir")
func (d DIRECTORY) FromPath(path string) DIRECTORY {
	return DIRECTORY{path: path}
}

// Path returns the current directory path.
// Returns: Current directory path as a string
// Usage:
// path := dir.Path()
func (d DIRECTORY) Path() string {
	return d.path
}

// Base returns the last element of the directory path.
// Returns: The base name of the directory path
// Usage:
// baseName := DIRECTORY{}.FromPath("/path/to/dir").Base() // returns "dir"
func (d DIRECTORY) Base() string {
	return filepath.Base(d.path)
}

// Clean returns a new DIRECTORY instance with the cleaned directory path.
// Returns: New DIRECTORY instance with cleaned path
// Usage:
// cleaned := DIRECTORY{}.FromPath("/path//to/../dir").Clean() // returns DIRECTORY{path: "/path/dir"}
func (d DIRECTORY) Clean() DIRECTORY {
	return DIRECTORY{path: filepath.Clean(d.path)}
}

// Create creates a new directory with mode 0755.
// Returns: Error encountered during directory creation
// Usage:
// err := DIRECTORY{}.FromPath("/new/dir").Create()
func (d DIRECTORY) Create() error {
	if d.path == "" {
		return os.ErrInvalid
	}
	return os.Mkdir(d.path, 0755)
}

// Delete removes the directory. The directory must be empty.
// Returns: Error encountered during directory deletion
// Usage:
// err := DIRECTORY{}.FromPath("/empty/dir").Delete()
func (d DIRECTORY) Delete() error {
	if d.path == "" {
		return os.ErrInvalid
	}
	return os.Remove(d.path)
}

// DeleteAll removes the directory and all its contents recursively.
// Returns: Error encountered during directory deletion
// Usage:
// err := DIRECTORY{}.FromPath("/dir/with/contents").DeleteAll()
func (d DIRECTORY) DeleteAll() error {
	if d.path == "" {
		return os.ErrInvalid
	}
	return os.RemoveAll(d.path)
}

// Dir returns a new DIRECTORY instance representing the parent directory.
// Returns: New DIRECTORY instance for parent directory
// Usage:
// parent := DIRECTORY{}.FromPath("/path/to/dir").Dir() // returns DIRECTORY{path: "/path/to"}
func (d DIRECTORY) Dir() DIRECTORY {
	return DIRECTORY{path: filepath.Dir(d.path)}
}

// Exists checks if the directory exists and is a directory.
// Returns: true if directory exists, false otherwise
// Usage:
//
//	if DIRECTORY{}.FromPath("/existing/dir").Exists() {
//	    // directory exists
//	}
func (d DIRECTORY) Exists() bool {
	var exists bool
	info, err := os.Stat(d.path)
	if err == nil && info.IsDir() {
		exists = true
	} else {
		exists = false
	}
	return exists
}

// FromSlash converts forward slashes to the operating system's path separator.
// Returns: New DIRECTORY instance with converted path
// Usage:
// path := DIRECTORY{}.FromPath("/path/to/dir").FromSlash()
func (d DIRECTORY) FromSlash() DIRECTORY {
	return DIRECTORY{path: filepath.FromSlash(d.path)}
}

// Get returns all entries in the directory.
// Returns: Slice of entry names and any error encountered
// Usage:
// entries, err := DIRECTORY{}.FromPath("/path/to/dir").Get()
func (d DIRECTORY) Get() ([]string, error) {
	var result []string
	if d.path == "" {
		return nil, os.ErrInvalid
	}
	entries, err := os.ReadDir(d.path)
	if err != nil {
		result = nil
	} else {
		result = make([]string, 0, len(entries))
		for _, entry := range entries {
			result = append(result, entry.Name())
		}
		sort.Strings(result)
	}
	return result, err
}

// GetDirectories returns all subdirectories in the directory.
// Returns: Slice of subdirectory names and any error encountered
// Usage:
// dirs, err := DIRECTORY{}.FromPath("/path/to/dir").GetDirectories()
func (d DIRECTORY) GetDirectories() ([]string, error) {
	var result []string
	if d.path == "" {
		return nil, os.ErrInvalid
	}
	entries, err := os.ReadDir(d.path)
	if err != nil {
		result = nil
	} else {
		result = make([]string, 0)
		for _, entry := range entries {
			if entry.IsDir() {
				result = append(result, entry.Name())
			}
		}
		sort.Strings(result)
	}
	return result, err
}

// GetFiles returns all files in the directory.
// Returns: Slice of file names and any error encountered
// Usage:
// files, err := DIRECTORY{}.FromPath("/path/to/dir").GetFiles()
func (d DIRECTORY) GetFiles() ([]string, error) {
	var result []string
	if d.path == "" {
		return nil, os.ErrInvalid
	}
	entries, err := os.ReadDir(d.path)
	if err != nil {
		result = nil
	} else {
		result = make([]string, 0)
		for _, entry := range entries {
			if !entry.IsDir() {
				result = append(result, entry.Name())
			}
		}
		sort.Strings(result)
	}
	return result, err
}

// GetWorkingDirectory returns a DIRECTORY instance representing the current working directory.
// Returns: DIRECTORY instance for current working directory and any error encountered
// Usage:
// cwd, err := DIRECTORY{}.GetWorkingDirectory()
func (d DIRECTORY) GetWorkingDirectory() (DIRECTORY, error) {
	var result DIRECTORY
	path, err := os.Getwd()
	if err == nil {
		result = DIRECTORY{path: path}
	}
	return result, err
}

// IsAbs checks if the directory path is absolute.
// Returns: true if path is absolute, false otherwise
// Usage:
//
//	if DIRECTORY{}.FromPath("/absolute/path").IsAbs() {
//	    // path is absolute
//	}
func (d DIRECTORY) IsAbs() bool {
	return filepath.IsAbs(d.path)
}

// Join returns a new DIRECTORY instance with the joined path.
// Returns: New DIRECTORY instance with joined path
// Usage:
// joined := DIRECTORY{}.FromPath("/path").Join("to", "dir") // returns DIRECTORY{path: "/path/to/dir"}
func (d DIRECTORY) Join(elem ...string) DIRECTORY {
	return DIRECTORY{path: filepath.Join(d.path, filepath.Join(elem...))}
}

// MakeAll creates the directory and all parent directories with mode 0755.
// Returns: Error encountered during directory creation
// Usage:
// err := DIRECTORY{}.FromPath("/deep/nested/dir").MakeAll()
func (d DIRECTORY) MakeAll() error {
	return os.MkdirAll(d.path, 0755)
}

// SetWorkingDirectory changes the current working directory to the directory path.
// Returns: Error encountered during directory change
// Usage:
// err := DIRECTORY{}.FromPath("/new/working/dir").SetWorkingDirectory()
func (d DIRECTORY) SetWorkingDirectory() error {
	if d.path == "" {
		return os.ErrInvalid
	}
	return os.Chdir(d.path)
}

// ToSlash converts the operating system's path separator to forward slashes.
// Returns: New DIRECTORY instance with forward slashes
// Usage:
// path := DIRECTORY{}.FromPath("\\windows\\path").ToSlash() // returns DIRECTORY{path: "/windows/path"}
func (d DIRECTORY) ToSlash() DIRECTORY {
	return DIRECTORY{path: filepath.ToSlash(d.path)}
}
