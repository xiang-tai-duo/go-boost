//go:build windows

// Package mesa3d
// File:        mesa3d.go
// Url:         https://github.com/xiang-tai-duo/go-boost-tools/blob/master/mesa3d/mesa3d.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Mesa3D utility functions for extracting and configuring OpenGL software rendering.
// --------------------------------------------------------------------------------
package mesa3d

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	MESA3D_UNPACKED_DIRECTORY_PERMISSION = 0755
	GOARCH_AMD64                         = "amd64"
	MODULE_NAME_MESA3D                   = "mesa3d"
)

// restoreAll extracts every file in embedFS into the executable directory,
// using only the base file name (matching the previous embed.RestoreAll with
// preserveEmbedRelativeDirectoryPath=false). Existing files are overwritten.
func restoreAll(embedFS fs.FS) (int, error) {
	result := 0
	cwd, err := os.Getwd()
	if err != nil {
		return 0, fmt.Errorf("failed to get working directory: %v", err)
	}
	err = fs.WalkDir(embedFS, ".", func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() || path == "." {
			return nil
		}
		data, err := fs.ReadFile(embedFS, path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %v", path, err)
		}
		target := filepath.Join(cwd, filepath.Base(path))
		if err = os.WriteFile(target, data, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %v", target, err)
		}
		result++
		return nil
	})
	if err != nil {
		return result, err
	}
	return result, nil
}

// isEmpty reports whether embedFS contains any file.
func isEmpty(embedFS fs.FS) (bool, error) {
	result := true
	err := fs.WalkDir(embedFS, ".", func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if walkErr == nil && path != "." && !d.IsDir() {
			result = false
			return fs.SkipAll
		}
		return nil
	})
	if err != nil && err != fs.SkipAll {
		return true, fmt.Errorf("failed to walk embedded filesystem: %v", err)
	}
	return result, nil
}

//goland:noinspection GoUnhandledErrorResult
func ExtractEnvironment() (int, error) {
	result := 0
	err := error(nil)
	executable := ""
	if executable, err = os.Executable(); err == nil {
		executableDirectory := filepath.Dir(executable)
		isEmptyEmbedFiles := false
		is64Bit := runtime.GOARCH == GOARCH_AMD64
		var embedFS fs.FS
		if is64Bit {
			embedFS = MESA3D_DISTRIBUTION_FILES_X64
		} else {
			embedFS = MESA3D_DISTRIBUTION_FILES_X86
		}
		if isEmptyEmbedFiles, err = isEmpty(embedFS); err == nil && !isEmptyEmbedFiles {
			if err = os.MkdirAll(executableDirectory, MESA3D_UNPACKED_DIRECTORY_PERMISSION); err == nil {
				if result, err = restoreAll(embedFS); err != nil {
					err = fmt.Errorf("failed to restore Mesa3D distribution files: %v", err)
				}
			} else {
				err = fmt.Errorf("failed to create directory %s: %v", executableDirectory, err)
			}
		}
		if err != nil {
			err = fmt.Errorf("failed to check embed files: %v", err)
		}
	}
	return result, err
}
