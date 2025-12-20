// Package file
// File:        file.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/file/file.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: File provides utility methods for file operations
// --------------------------------------------------------------------------------
package file

import (
	"os"
	"path/filepath"
	"time"
)

type (
	FILE struct {
		path string
	}
)

//goland:noinspection GoUnusedGlobalVariable
var (
	File = FILE{}
)

//goland:noinspection GoUnusedExportedFunction
func New(path ...string) *FILE {
	var filePath string
	if len(path) > 0 {
		filePath = path[0]
	}
	return &FILE{path: filePath}
}

func (f FILE) Abs() (FILE, error) {
	var result FILE
	absPath, err := filepath.Abs(f.path)
	if err == nil {
		result = FILE{path: absPath}
	}
	return result, err
}

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

func (f FILE) AppendString(content string) error {
	return f.Append([]byte(content))
}

func (f FILE) Base() string {
	return filepath.Base(f.path)
}

func (f FILE) Chmod(mode os.FileMode) error {
	return os.Chmod(f.path, mode)
}

func (f FILE) Chown(uid int, gid int) error {
	return os.Chown(f.path, uid, gid)
}

//goland:noinspection SpellCheckingInspection
func (f FILE) Chtimes(atime time.Time, mtime time.Time) error {
	return os.Chtimes(f.path, atime, mtime)
}

func (f FILE) Clean() FILE {
	return FILE{path: filepath.Clean(f.path)}
}

func (f FILE) Copy(dest FILE) error {
	var err error
	if f.path == "" || dest.path == "" {
		return os.ErrInvalid
	}
	data, readErr := f.ReadAll()
	if readErr == nil {
		err = dest.Write(data)
	} else {
		err = readErr
	}
	return err
}

func (f FILE) Delete() error {
	if f.path == "" {
		return os.ErrInvalid
	}
	return os.Remove(f.path)
}

func (f FILE) IsExists() bool {
	var exists bool
	info, err := os.Stat(f.path)
	if err == nil && !info.IsDir() {
		exists = true
	} else {
		exists = false
	}
	return exists
}

func (f FILE) Ext() string {
	return filepath.Ext(f.path)
}

func (f FILE) FromPath(path string) FILE {
	return FILE{path: path}
}

func (f FILE) FromSlash() FILE {
	return FILE{path: filepath.FromSlash(f.path)}
}

func (f FILE) IsAbs() bool {
	return filepath.IsAbs(f.path)
}

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

func (f FILE) Mode() os.FileMode {
	var mode os.FileMode
	info, err := os.Stat(f.path)
	if err == nil {
		mode = info.Mode()
	}
	return mode
}

func (f FILE) ModTime() time.Time {
	var modTime time.Time
	info, err := os.Stat(f.path)
	if err == nil {
		modTime = info.ModTime()
	}
	return modTime
}

func (f FILE) Name() string {
	return filepath.Base(f.path)
}

func (f FILE) Path() string {
	return f.path
}

func (f FILE) ReadAll() ([]byte, error) {
	if f.path == "" {
		return nil, os.ErrInvalid
	}
	return os.ReadFile(f.path)
}

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

func (f FILE) Rel(basepath string) (string, error) {
	return filepath.Rel(basepath, f.path)
}

func (f FILE) Rename(newPath FILE) error {
	if f.path == "" || newPath.path == "" {
		return os.ErrInvalid
	}
	return os.Rename(f.path, newPath.path)
}

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

func (f FILE) ToSlash() FILE {
	return FILE{path: filepath.ToSlash(f.path)}
}

func (f FILE) Write(data []byte) error {
	if f.path == "" {
		return os.ErrInvalid
	}
	return os.WriteFile(f.path, data, 0644)
}

func (f FILE) WriteString(content string) error {
	return os.WriteFile(f.path, []byte(content), 0644)
}
