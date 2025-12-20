// Package directory
// File:        directory.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/directory/directory.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: Directory handles operations related to the current directory, including creation, deletion, and listing
// --------------------------------------------------------------------------------
package directory

import (
	"os"
	"path/filepath"
	"sort"
)

type (
	DIRECTORY struct {
		path string
	}
)

//goland:noinspection GoUnusedExportedFunction
func New(path ...string) *DIRECTORY {
	var dirPath string
	if len(path) > 0 {
		dirPath = path[0]
	}
	return &DIRECTORY{path: dirPath}
}

func (d *DIRECTORY) Base() string {
	return filepath.Base(d.path)
}

func (d *DIRECTORY) Clean() *DIRECTORY {
	return &DIRECTORY{path: filepath.Clean(d.path)}
}

func (d *DIRECTORY) Create() error {
	if d.path == "" {
		return os.ErrInvalid
	}
	return os.Mkdir(d.path, 0755)
}

func (d *DIRECTORY) Delete() error {
	if d.path == "" {
		return os.ErrInvalid
	}
	return os.Remove(d.path)
}

func (d *DIRECTORY) DeleteAll() error {
	if d.path == "" {
		return os.ErrInvalid
	}
	return os.RemoveAll(d.path)
}

func (d *DIRECTORY) FromPath(path string) *DIRECTORY {
	return &DIRECTORY{path: path}
}

func (d *DIRECTORY) FromSlash() *DIRECTORY {
	return &DIRECTORY{path: filepath.FromSlash(d.path)}
}

func (d *DIRECTORY) GetAll() ([]string, error) {
	result := []string{}
	err := error(nil)
	if d.path == "" {
		return nil, os.ErrInvalid
	}
	entries, err := os.ReadDir(d.path)
	if err == nil {
		result = make([]string, 0, len(entries))
		for _, entry := range entries {
			result = append(result, entry.Name())
		}
		sort.Strings(result)
	} else {
		result = nil
	}
	return result, err
}

func (d *DIRECTORY) GetDirectories() ([]string, error) {
	result := []string{}
	err := error(nil)
	if d.path == "" {
		return nil, os.ErrInvalid
	}
	entries, err := os.ReadDir(d.path)
	if err == nil {
		result = make([]string, 0)
		for _, entry := range entries {
			if entry.IsDir() {
				result = append(result, entry.Name())
			}
		}
		sort.Strings(result)
	} else {
		result = nil
	}
	return result, err
}

func (d *DIRECTORY) GetFiles() ([]string, error) {
	result := []string{}
	err := error(nil)
	if d.path == "" {
		return nil, os.ErrInvalid
	}
	entries, err := os.ReadDir(d.path)
	if err == nil {
		result = make([]string, 0)
		for _, entry := range entries {
			if !entry.IsDir() {
				result = append(result, entry.Name())
			}
		}
		sort.Strings(result)
	} else {
		result = nil
	}
	return result, err
}

func (d *DIRECTORY) GetWorkingDirectory() (*DIRECTORY, error) {
	result := (*DIRECTORY)(nil)
	err := error(nil)
	path, err := os.Getwd()
	if err == nil {
		result = &DIRECTORY{path: path}
	}
	return result, err
}

func (d *DIRECTORY) IsAbs() bool {
	return filepath.IsAbs(d.path)
}

func (d *DIRECTORY) IsExists() bool {
	result := false
	info, err := os.Stat(d.path)
	if err == nil && info.IsDir() {
		result = true
	}
	return result
}

func (d *DIRECTORY) Join(elem ...string) *DIRECTORY {
	return &DIRECTORY{path: filepath.Join(d.path, filepath.Join(elem...))}
}

func (d *DIRECTORY) MakeAll() error {
	return os.MkdirAll(d.path, 0755)
}

func (d *DIRECTORY) Path() string {
	return d.path
}

func (d *DIRECTORY) SetWorkingDirectory() error {
	result := error(nil)
	if d.path == "" {
		result = os.ErrInvalid
	} else {
		result = os.Chdir(d.path)
	}
	return result
}

func (d *DIRECTORY) ToSlash() *DIRECTORY {
	return &DIRECTORY{path: filepath.ToSlash(d.path)}
}
