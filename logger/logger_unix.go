//go:build !windows

// Package logger
// File:        logger_unix.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/logger/logger_unix.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Logger disk space helpers for Unix systems.
// --------------------------------------------------------------------------------
package logger

import "syscall"

func getDiskFreeSpace(directory string) (uint64, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(directory, &stat)
	if err != nil {
		return 0, err
	}
	return stat.Bavail * uint64(stat.Bsize), nil
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func SetFileCompression(filePath string, enable bool) error {
	_ = filePath
	_ = enable
	return nil
}