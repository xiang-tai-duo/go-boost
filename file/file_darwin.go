//go:build darwin

// Package file
// File:        file_darwin.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/file/file_darwin.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: File provides Darwin-specific file operations
// --------------------------------------------------------------------------------
package file

import (
	"errors"
	"os"

	"golang.org/x/sys/unix"
)

func IsExecutable(filePath string) bool {
	return false
}

func OpenLockedRead(path string) (*os.File, error) {
	return openExclusiveDarwin(path, os.O_RDONLY)
}

func OpenLockedReadOnly(path string) (*os.File, error) {
	return openExclusiveDarwin(path, os.O_RDONLY)
}

//goland:noinspection DuplicatedCode
func openExclusiveDarwin(path string, flag int) (*os.File, error) {
	result := (*os.File)(nil)
	err := error(nil)
	if path == "" {
		err = os.ErrInvalid
	} else {
		var file *os.File
		if file, err = os.OpenFile(path, flag, 0); err == nil {
			if err = lockFileExclusiveDarwin(file); err == nil {
				result = file
			} else {
				file.Close()
			}
		}
	}
	return result, err
}

func lockFileExclusiveDarwin(file *os.File) error {
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
