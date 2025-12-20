// --------------------------------------------------------------------------------
// File:        file.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: File provides utility methods for file operations,
//              while FileBuilder is used for operations on specific files.
// --------------------------------------------------------------------------------

package boost

import (
	"os"
	"path/filepath"
	"time"
)

// FILE provides utility methods for file operations.
type FILE struct {
	path string
}

// NewFile creates a new FILE instance with the given path.
// path: Optional file path to use
// Returns: New FILE instance
// Usage:
// file := NewFile() // creates with empty path
// file := NewFile("/path/to/file") // creates with specified path
func NewFile(path ...string) *FILE {
	var filePath string
	if len(path) > 0 {
		filePath = path[0]
	}
	return &FILE{path: filePath}
}

// FromPath creates a new FILE instance with the given path.
// path: File path to use
// Returns: New FILE instance
// Usage:
// file := FILE{}.FromPath("/path/to/file")
func (f FILE) FromPath(path string) FILE {
	return FILE{path: path}
}

// Path returns the current file path.
// Returns: Current file path as a string
// Usage:
// path := file.Path()
func (f FILE) Path() string {
	return f.path
}

// Abs returns a new FILE instance with the absolute path.
// Returns: New FILE instance with absolute path and any error encountered
// Usage:
// absPath, err := FILE{}.FromPath("relative/path").Abs()
func (f FILE) Abs() (FILE, error) {
	var result FILE
	absPath, err := filepath.Abs(f.path)
	if err == nil {
		result = FILE{path: absPath}
	}
	return result, err
}

// Append appends data to the file.
// data: Data to append to the file
// Returns: Error encountered during file append operation
// Usage:
// err := FILE{}.FromPath("file.txt").Append([]byte("appended data"))
func (f FILE) Append(data []byte) error {
	var err error
	var file *os.File
	if f.path == "" {
		err = os.ErrInvalid
	} else if file, err = os.OpenFile(f.path, os.O_APPEND|os.O_WRONLY, 0644); err == nil {
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
		_, err = file.Write(data)
	}
	return err
}

// AppendString appends a string to the file.
// content: String to append to the file
// Returns: Error encountered during file append operation
// Usage:
// err := FILE{}.FromPath("file.txt").AppendString("appended string")
func (f FILE) AppendString(content string) error {
	return f.Append([]byte(content))
}

// Base returns the last element of the file path.
// Returns: The base name of the file path
// Usage:
// baseName := FILE{}.FromPath("/path/to/file.txt").Base() // returns "file.txt"
func (f FILE) Base() string {
	return filepath.Base(f.path)
}

// Chmod changes the file mode.
// mode: New file mode to set
// Returns: Error encountered during file mode change
// Usage:
// err := FILE{}.FromPath("file.txt").Chmod(0644)
func (f FILE) Chmod(mode os.FileMode) error {
	return os.Chmod(f.path, mode)
}

// Chown changes the file owner and group.
// uid: User ID to set
// gid: Group ID to set
// Returns: Error encountered during file ownership change
// Usage:
// err := FILE{}.FromPath("file.txt").Chown(1000, 1000)
func (f FILE) Chown(uid, gid int) error {
	return os.Chown(f.path, uid, gid)
}

// Chtimes changes the file access and modification times.
// atime: New access time to set
// mtime: New modification time to set
// Returns: Error encountered during file times change
// Usage:
// err := FILE{}.FromPath("file.txt").Chtimes(time.Now(), time.Now())
func (f FILE) Chtimes(atime, mtime time.Time) error {
	return os.Chtimes(f.path, atime, mtime)
}

// Clean returns a new FILE instance with the cleaned file path.
// Returns: New FILE instance with cleaned path
// Usage:
// cleaned := FILE{}.FromPath("/path//to/../file.txt").Clean() // returns FILE{path: "/path/file.txt"}
func (f FILE) Clean() FILE {
	return FILE{path: filepath.Clean(f.path)}
}

// Copy copies the file to the destination path.
// dest: Destination file path
// Returns: Error encountered during file copy operation
// Usage:
// err := FILE{}.FromPath("source.txt").Copy(FILE{}.FromPath("dest.txt"))
func (f FILE) Copy(dest FILE) error {
	var err error
	if f.path == "" || dest.path == "" {
		return os.ErrInvalid
	}
	data, readErr := f.ReadAll()
	if readErr != nil {
		err = readErr
	} else {
		err = dest.Write(data)
	}
	return err
}

// Delete removes the file.
// Returns: Error encountered during file deletion
// Usage:
// err := FILE{}.FromPath("file.txt").Delete()
func (f FILE) Delete() error {
	if f.path == "" {
		return os.ErrInvalid
	}
	return os.Remove(f.path)
}

// Dir returns a new DIRECTORY instance representing the parent directory.
// Returns: New DIRECTORY instance for parent directory
// Usage:
// parent := FILE{}.FromPath("/path/to/file.txt").Dir() // returns DIRECTORY{path: "/path/to"}
func (f FILE) Dir() DIRECTORY {
	return DIRECTORY{path: filepath.Dir(f.path)}
}

// Exists checks if the file exists and is a regular file.
// Returns: true if file exists, false otherwise
// Usage:
//
//	if FILE{}.FromPath("existing.txt").Exists() {
//	    // file exists
//	}
func (f FILE) Exists() bool {
	var exists bool
	info, err := os.Stat(f.path)
	if err == nil && !info.IsDir() {
		exists = true
	} else {
		exists = false
	}
	return exists
}

// Ext returns the file extension.
// Returns: The file extension including the dot
// Usage:
// ext := FILE{}.FromPath("file.txt").Ext() // returns ".txt"
func (f FILE) Ext() string {
	return filepath.Ext(f.path)
}

// FromSlash converts forward slashes to the operating system's path separator.
// Returns: New FILE instance with converted path
// Usage:
// path := FILE{}.FromPath("/path/to/file.txt").FromSlash()
func (f FILE) FromSlash() FILE {
	return FILE{path: filepath.FromSlash(f.path)}
}

// IsAbs checks if the file path is absolute.
// Returns: true if path is absolute, false otherwise
// Usage:
//
//	if FILE{}.FromPath("/absolute/path.txt").IsAbs() {
//	    // path is absolute
//	}
func (f FILE) IsAbs() bool {
	return filepath.IsAbs(f.path)
}

// IsDir checks if the path is a directory.
// Returns: true if path is a directory, false otherwise
// Usage:
//
//	if FILE{}.FromPath("/path/to/dir").IsDir() {
//	    // path is a directory
//	}
func (f FILE) IsDir() bool {
	var isDirectory bool
	info, err := os.Stat(f.path)
	if err == nil && info.IsDir() {
		isDirectory = true
	} else {
		isDirectory = false
	}
	return isDirectory
}

// IsRegular checks if the file is a regular file.
// Returns: true if file is regular, false otherwise
// Usage:
//
//	if FILE{}.FromPath("file.txt").IsRegular() {
//	    // file is regular
//	}
func (f FILE) IsRegular() bool {
	var isRegular bool
	info, err := os.Stat(f.path)
	if err == nil && info.Mode().IsRegular() {
		isRegular = true
	} else {
		isRegular = false
	}
	return isRegular
}

// IsSymlink checks if the file is a symbolic link.
// Returns: true if file is a symlink, false otherwise
// Usage:
//
//	if FILE{}.FromPath("symlink.txt").IsSymlink() {
//	    // file is a symlink
//	}
func (f FILE) IsSymlink() bool {
	var isSymlink bool
	info, err := os.Lstat(f.path)
	if err == nil && info.Mode()&os.ModeSymlink != 0 {
		isSymlink = true
	} else {
		isSymlink = false
	}
	return isSymlink
}

// Mode returns the file mode.
// Returns: The file mode
// Usage:
// mode := FILE{}.FromPath("file.txt").Mode()
func (f FILE) Mode() os.FileMode {
	var mode os.FileMode
	info, err := os.Stat(f.path)
	if err == nil {
		mode = info.Mode()
	}
	return mode
}

// ModTime returns the file modification time.
// Returns: The file modification time
// Usage:
// modTime := FILE{}.FromPath("file.txt").ModTime()
func (f FILE) ModTime() time.Time {
	var modTime time.Time
	info, err := os.Stat(f.path)
	if err == nil {
		modTime = info.ModTime()
	}
	return modTime
}

// Name returns the last element of the file path.
// Returns: The base name of the file path
// Usage:
// name := FILE{}.FromPath("/path/to/file.txt").Name() // returns "file.txt"
func (f FILE) Name() string {
	return filepath.Base(f.path)
}

// ReadAll reads all content from the file.
// Returns: File content as byte slice and any error encountered
// Usage:
// content, err := FILE{}.FromPath("file.txt").ReadAll()
func (f FILE) ReadAll() ([]byte, error) {
	if f.path == "" {
		return nil, os.ErrInvalid
	}
	return os.ReadFile(f.path)
}

// ReadString reads all content from the file as a string.
// Returns: File content as string and any error encountered
// Usage:
// content, err := FILE{}.FromPath("file.txt").ReadString()
func (f FILE) ReadString() (string, error) {
	var content string
	if f.path == "" {
		return "", os.ErrInvalid
	}
	data, err := os.ReadFile(f.path)
	if err == nil {
		content = string(data)
	}
	return content, err
}

// Rel returns the relative path from basepath to the file.
// basepath: Base path to compute relative path from
// Returns: Relative path and any error encountered
// Usage:
// relPath, err := FILE{}.FromPath("/path/to/file.txt").Rel("/path") // returns "to/file.txt", nil
func (f FILE) Rel(basepath string) (string, error) {
	return filepath.Rel(basepath, f.path)
}

// Rename renames the file to the new path.
// newPath: New file path
// Returns: Error encountered during file rename
// Usage:
// err := FILE{}.FromPath("old.txt").Rename(FILE{}.FromPath("new.txt"))
func (f FILE) Rename(newPath FILE) error {
	if f.path == "" || newPath.path == "" {
		return os.ErrInvalid
	}
	return os.Rename(f.path, newPath.path)
}

// Size returns the file size in bytes.
// Returns: File size in bytes, or -1 if error occurs
// Usage:
// size := FILE{}.FromPath("file.txt").Size()
func (f FILE) Size() int64 {
	var size int64
	info, err := os.Stat(f.path)
	if err == nil {
		size = info.Size()
	} else {
		size = -1
	}
	return size
}

// ToSlash converts the operating system's path separator to forward slashes.
// Returns: New FILE instance with forward slashes
// Usage:
// path := FILE{}.FromPath("\\windows\\path.txt").ToSlash() // returns FILE{path: "/windows/path.txt"}
func (f FILE) ToSlash() FILE {
	return FILE{path: filepath.ToSlash(f.path)}
}

// Write writes data to the file, truncating it if it exists.
// data: _data to write to the file
// Returns: Error encountered during file write
// Usage:
// err := FILE{}.FromPath("file.txt").Write([]byte("content"))
func (f FILE) Write(data []byte) error {
	if f.path == "" {
		return os.ErrInvalid
	}
	return os.WriteFile(f.path, data, 0644)
}

// WriteString writes a string to the file, truncating it if it exists.
// content: String to write to the file
// Returns: Error encountered during file write
// Usage:
// err := FILE{}.FromPath("file.txt").WriteString("content")
func (f FILE) WriteString(content string) error {
	return os.WriteFile(f.path, []byte(content), 0644)
}
