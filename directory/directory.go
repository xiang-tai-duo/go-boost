// Package directory
// File:        directory.go
// Url:         https://github.com/xiang-tai-duo/go-bootstrap/blob/master/directory/directory.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Directory handles operations related to the current directory, including creation, deletion, and listing
// --------------------------------------------------------------------------------
package directory

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"

	"github.com/xiang-tai-duo/go-bootstrap/strings2"
)

//goland:noinspection GoUnusedExportedFunction
func CreateTemporaryDirectory() string {
	result := filepath.Join(os.TempDir(), strings2.Random(strings2.DEFAULT_RANDOM_SIZE))
	if os.Mkdir(result, GetDefaultPermission()) != nil {
		result = ""
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetDirectoriesAndFiles(path string) []string {
	result := make([]string, 0)
	if path != "" {
		if entries, err := os.ReadDir(path); err == nil {
			result = make([]string, 0, len(entries))
			for _, entry := range entries {
				result = append(result, entry.Name())
			}
			sort.Strings(result)
		}
	}
	return result
}

func GetDefaultPermission() os.FileMode {
	fileMode := os.FileMode(0755)
	if runtime.GOOS == "windows" {
		fileMode = os.FileMode(0755)
	}
	return fileMode
}

//goland:noinspection GoUnusedExportedFunction
func GetDirectories(path string) []string {
	result := make([]string, 0)
	if path != "" {
		if entries, err := os.ReadDir(path); err == nil {
			result = make([]string, 0)
			for _, entry := range entries {
				if entry.IsDir() {
					result = append(result, entry.Name())
				}
			}
			sort.Strings(result)
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetFiles(path string) []string {
	result := make([]string, 0)
	if path != "" {
		if entries, err := os.ReadDir(path); err == nil {
			result = make([]string, 0)
			for _, entry := range entries {
				if !entry.IsDir() {
					result = append(result, entry.Name())
				}
			}
			sort.Strings(result)
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func IsExists(path string) bool {
	result := false
	if path != "" {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			result = true
		}
	}
	return result
}
