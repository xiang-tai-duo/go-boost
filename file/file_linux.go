//go:build !windows && !darwin

// Package file
// File:        file_linux.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/file/file_linux.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: File provides Linux-specific file operations
// --------------------------------------------------------------------------------
package file

import (
	"errors"
	"os"

	"golang.org/x/sys/unix"
)

//goland:noinspection GoSnakeCaseUsage
const (
	ELF_MAGIC_NUMBER  = "\x7fELF"
	MAGIC_BUFFER_SIZE = 4
	SHEBANG_PREFIX    = "#!"
)

//goland:noinspection GoUnhandledErrorResult
func IsExecutable(filePath string) bool {
	result := false
	if fileInfo, err := os.Stat(filePath); err == nil {
		if !fileInfo.IsDir() {
			var file *os.File
			if file, err = os.Open(filePath); err == nil {
				defer file.Close()
				buffer := make([]byte, MAGIC_BUFFER_SIZE)
				if n, err := file.Read(buffer); err == nil {
					if n >= MAGIC_BUFFER_SIZE && string(buffer[:MAGIC_BUFFER_SIZE]) == ELF_MAGIC_NUMBER {
						result = true
					} else if n >= 2 && string(buffer[:2]) == SHEBANG_PREFIX {
						result = true
					}
				}
			}
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func OpenLockedRead(path string) (*os.File, error) {
	return openExclusiveUnix(path, os.O_RDONLY)
}

//goland:noinspection GoUnusedExportedFunction
func OpenLockedReadOnly(path string) (*os.File, error) {
	return openExclusiveUnix(path, os.O_RDONLY)
}

//goland:noinspection GoUnhandledErrorResult,DuplicatedCode
func openExclusiveUnix(path string, flag int) (*os.File, error) {
	result := (*os.File)(nil)
	err := error(nil)
	if path == "" {
		err = os.ErrInvalid
	} else {
		var file *os.File
		if file, err = os.OpenFile(path, flag, 0); err == nil {
			if err = lockFileExclusive(file); err == nil {
				result = file
			} else {
				file.Close()
			}
		}
	}
	return result, err
}

func lockFileExclusive(file *os.File) error {
	err := error(nil)
	if file == nil {
		err = errors.New("file is nil")
	} else {
		if lockErr := unix.Flock(int(file.Fd()), unix.LOCK_EX|unix.LOCK_NB); lockErr != nil {
			err = errors.New("file is locked by another process")
		}
	}
	return err
}
