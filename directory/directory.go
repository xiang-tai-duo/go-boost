// Package directory
// File:        directory.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/directory/directory.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Directory handles operations related to the current directory, including creation, deletion, and listing
// --------------------------------------------------------------------------------
package directory

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/xiang-tai-duo/go-boost/file"
	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	MODULE_NAME_DIRECTORY = "directory"
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_DIRECTORY, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_DIRECTORY, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_DIRECTORY, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedExportedFunction
func Copy(source string, destination string) error {
	err := error(nil)
	sourceDirectory := ""
	destinationDirectory := ""
	if source != "" && destination != "" {
		sourceDirectory = filepath.Clean(source)
		destinationDirectory = filepath.Clean(destination)
		if sourceDirectory != destinationDirectory {
			err = filepath.WalkDir(sourceDirectory, func(path string, d os.DirEntry, walkErr error) error {
				result := error(nil)
				if walkErr != nil {
					result = walkErr
				} else {
					relativePath := ""
					if relativePath, result = filepath.Rel(sourceDirectory, path); result == nil {
						destinationPath := filepath.Join(destinationDirectory, relativePath)
						if d.IsDir() {
							result = os.MkdirAll(destinationPath, os.ModePerm)
						} else {
							if result = os.MkdirAll(filepath.Dir(destinationPath), os.ModePerm); result == nil {
								result = file.Copy(path, destinationPath)
							}
						}
					}
				}
				return result
			})
		}
	} else {
		err = os.ErrInvalid
	}
	return err
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
	err := error(nil)
	var info os.FileInfo
	if path != "" {
		if info, err = os.Stat(path); err == nil {
			if info.IsDir() {
				result = true
			}
		}
	}
	return result
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_DIRECTORY, logger.SKIP_STACK_FRAMES_BASE)
}
