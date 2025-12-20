// --------------------------------------------------------------------------------
// File:        filepath.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: FilePath provides utility methods for file path operations,
//              while FilePathBuilder is used for operations on specific file paths.
// --------------------------------------------------------------------------------

package boost

import (
	"io/fs"
	"path/filepath"
	"strings"
)

// FILEPATH provides utility methods for file path operations.
type FILEPATH struct {
	_filepath string
}

// NewFilePath creates a new FILEPATH instance with the given path.
// filePath: File path to use
// Returns: New FILEPATH instance
// Usage:
// filePath := NewFilePath("/path/to/file") // creates with specified path
func NewFilePath(filePath string) *FILEPATH {
	return &FILEPATH{_filepath: filePath}
}

// Path returns the current file path or sets a new file path.
// filePath: Optional file path to set
// Returns: Current file path as a string
// Usage:
// path := filePath.Path() // returns current file path
// filePath.Path("/new/path") // sets new file path
func (f *FILEPATH) Path(filePath ...string) string {
	if len(filePath) > 0 {
		f._filepath = filePath[0]
	}
	return f._filepath
}

// Abs returns a new FILEPATH instance with the absolute path.
// Returns: New FILEPATH instance with absolute path and any error encountered
// Usage:
// absPath, err := NewFilePath("relative/path").Abs()
func (f FILEPATH) Abs() (FILEPATH, error) {
	var result FILEPATH
	absPath, err := filepath.Abs(f._filepath)
	if err == nil {
		result = FILEPATH{_filepath: absPath}
	}
	return result, err
}

// Base returns the last element of the file path.
// Returns: The base name of the file path
// Usage:
// baseName := NewFilePath("/path/to/file.txt").Base() // returns "file.txt"
func (f FILEPATH) Base() string {
	return filepath.Base(f._filepath)
}

// Clean returns a new FILEPATH instance with the cleaned file path.
// Returns: New FILEPATH instance with cleaned path
// Usage:
// cleaned := NewFilePath("/path//to/../file.txt").Clean() // returns FILEPATH{_filepath: "/path/file.txt"}
func (f FILEPATH) Clean() FILEPATH {
	return FILEPATH{_filepath: filepath.Clean(f._filepath)}
}

// Dir returns the parent directory of the file path.
// Returns: The parent directory path as string
// Usage:
// dir := NewFilePath("/path/to/file.txt").Dir() // returns "/path/to"
func (f FILEPATH) Dir() string {
	return filepath.Dir(f._filepath)
}

// EvalSymlinks evaluates any symbolic links in the file path.
// Returns: New FILEPATH instance with evaluated path and any error encountered
// Usage:
// evalPath, err := NewFilePath("/link/to/file").EvalSymlinks()
func (f FILEPATH) EvalSymlinks() (FILEPATH, error) {
	var result FILEPATH
	evalPath, err := filepath.EvalSymlinks(f._filepath)
	if err == nil {
		result = FILEPATH{_filepath: evalPath}
	}
	return result, err
}

// Ext returns the file extension of the file path.
// Returns: The file extension including the dot
// Usage:
// ext := NewFilePath("file.txt").Ext() // returns ".txt"
func (f FILEPATH) Ext() string {
	return filepath.Ext(f._filepath)
}

// FromSlash converts forward slashes to the operating system's path separator.
// Returns: New FILEPATH instance with converted path
// Usage:
// path := NewFilePath("/path/to/file").FromSlash()
func (f FILEPATH) FromSlash() FILEPATH {
	return FILEPATH{_filepath: filepath.FromSlash(f._filepath)}
}

// GetFileExtension returns the file extension from the file path.
// Returns: The file extension including the dot
// Usage:
// extension := NewFilePath("/path/to/file.txt").GetFileExtension() // returns ".txt"
func (f FILEPATH) GetFileExtension() string {
	return filepath.Ext(f._filepath)
}

// GetFileName returns the file name from the file path.
// Returns: The file name
// Usage:
// filename := NewFilePath("/path/to/file.txt").GetFileName() // returns "file.txt"
func (f FILEPATH) GetFileName() string {
	return filepath.Base(f._filepath)
}

// GetFileNameWithoutExtension returns the file name without extension from the file path.
// Returns: The file name without extension
// Usage:
// filename := NewFilePath("/path/to/file.txt").GetFileNameWithoutExtension() // returns "file"
func (f FILEPATH) GetFileNameWithoutExtension() string {
	filename := filepath.Base(f._filepath)
	extension := filepath.Ext(filename)
	return strings.TrimSuffix(filename, extension)
}

// GetParentDirectory returns the parent directory from the file path.
// Returns: The parent directory path
// Usage:
// parent := NewFilePath("/path/to/file.txt").GetParentDirectory() // returns "/path/to"
func (f FILEPATH) GetParentDirectory() string {
	return filepath.Dir(f._filepath)
}

// Glob returns all file paths matching the pattern.
// Returns: Slice of matching file paths and any error encountered
// Usage:
// matches, err := NewFilePath("/path/*.txt").Glob()
func (f FILEPATH) Glob() ([]string, error) {
	return filepath.Glob(f._filepath)
}

// HasPrefix checks if the file path has the given prefix.
// prefix: Prefix to check
// Returns: true if path has the prefix, false otherwise
// Usage:
//
//	if NewFilePath("/path/to/file").HasPrefix("/path") {
//	    // path has prefix "/path"
//	}
func (f FILEPATH) HasPrefix(prefix string) bool {
	return filepath.HasPrefix(f._filepath, prefix)
}

// IsAbs checks if the file path is absolute.
// Returns: true if path is absolute, false otherwise
// Usage:
//
//	if NewFilePath("/absolute/path").IsAbs() {
//	    // path is absolute
//	}
func (f FILEPATH) IsAbs() bool {
	return filepath.IsAbs(f._filepath)
}

// IsLocal checks if the file path is local.
// Returns: true if path is local, false otherwise
// Usage:
//
//	if NewFilePath("local/path").IsLocal() {
//	    // path is local
//	}
func (f FILEPATH) IsLocal() bool {
	return filepath.IsLocal(f._filepath)
}

// Join joins the file path with the given elements.
// elem: Elements to join with the path
// Returns: New FILEPATH instance with joined path
// Usage:
// joined := NewFilePath("/path").Join("to", "file.txt") // returns FILEPATH{_filepath: "/path/to/file.txt"}
func (f FILEPATH) Join(elem ...string) FILEPATH {
	args := make([]string, len(elem)+1)
	args[0] = f._filepath
	copy(args[1:], elem)
	return FILEPATH{_filepath: filepath.Join(args...)}
}

// Localize localizes the file path for the current operating system.
// Returns: New FILEPATH instance with localized path and any error encountered
// Usage:
// localPath, err := NewFilePath("/path/to/file").Localize()
func (f FILEPATH) Localize() (FILEPATH, error) {
	var result FILEPATH
	localPath, err := filepath.Localize(f._filepath)
	if err == nil {
		result = FILEPATH{_filepath: localPath}
	}
	return result, err
}

// Match checks if the file path matches the given pattern.
// pattern: Pattern to match against
// Returns: true if path matches pattern, false otherwise, and any error encountered
// Usage:
// matched, err := NewFilePath("file.txt").Match("*.txt") // returns true, nil
func (f FILEPATH) Match(pattern string) (bool, error) {
	return filepath.Match(pattern, f._filepath)
}

// Rel returns the relative path from basepath to the file path.
// basepath: Base path to compute relative path from
// Returns: New FILEPATH instance with relative path and any error encountered
// Usage:
// relPath, err := NewFilePath("/path/to/file").Rel("/path") // returns FILEPATH{_filepath: "to/file"}, nil
func (f FILEPATH) Rel(basepath string) (FILEPATH, error) {
	var result FILEPATH
	relPath, err := filepath.Rel(basepath, f._filepath)
	if err == nil {
		result = FILEPATH{_filepath: relPath}
	}
	return result, err
}

// Split splits the file path into directory and file name.
// Returns: Directory and file name as separate strings
// Usage:
// dir, file := NewFilePath("/path/to/file.txt").Split() // returns "/path/to", "file.txt"
func (f FILEPATH) Split() (string, string) {
	return filepath.Split(f._filepath)
}

// SplitList splits the file path list into separate paths.
// Returns: Slice of split paths
// Usage:
// paths := NewFilePath("/path1:/path2:/path3").SplitList() // returns ["/path1", "/path2", "/path3"]
func (f FILEPATH) SplitList() []string {
	return filepath.SplitList(f._filepath)
}

// ToSlash converts the operating system's path separator to forward slashes.
// Returns: New FILEPATH instance with forward slashes
// Usage:
// path := NewFilePath("\\windows\\path").ToSlash() // returns FILEPATH{_filepath: "/windows/path"}
func (f FILEPATH) ToSlash() FILEPATH {
	return FILEPATH{_filepath: filepath.ToSlash(f._filepath)}
}

// VolumeName returns the volume name of the file path.
// Returns: The volume name
// Usage:
// volume := NewFilePath("C:\\path").VolumeName() // returns "C:" on Windows
func (f FILEPATH) VolumeName() string {
	return filepath.VolumeName(f._filepath)
}

// Walk walks the file path tree, calling fn for each file and directory.
// fn: Function to call for each file and directory
// Returns: Error encountered during walk
// Usage:
//
//	err := NewFilePath("/path").Walk(func(path string, info os.FileInfo, err error) error {
//	    // handle file or directory
//	    return nil
//	})
func (f FILEPATH) Walk(fn filepath.WalkFunc) error {
	return filepath.Walk(f._filepath, fn)
}

// WalkDir walks the file path tree, calling fn for each file and directory using fs.DirEntry.
// fn: Function to call for each file and directory
// Returns: Error encountered during walk
// Usage:
//
//	err := NewFilePath("/path").WalkDir(func(path string, d fs.DirEntry, err error) error {
//	    // handle file or directory
//	    return nil
//	})
func (f FILEPATH) WalkDir(fn fs.WalkDirFunc) error {
	return filepath.WalkDir(f._filepath, fn)
}
