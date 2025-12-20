//go:build darwin

// Package system
// File:        system_darwin.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/system/system_darwin.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: SYSTEM is a wrapper for operating system interactions, providing macOS (darwin) specific implementations for architecture detection and long path resolution.
// --------------------------------------------------------------------------------

package system

import (
	"fmt"

	"github.com/xiang-tai-duo/go-boost/hardware"
	"github.com/xiang-tai-duo/go-boost/logger"
)

func Architecture() string {
	result := ""
	return result
}

func GetLongPathName(path string) (string, error) {
	result := path
	err := error(nil)
	return result, err
}

func acquireSingletonLock(executableFilePath string) error {
	logger.Logger.Debug("acquireSingletonLock enter executableFilePath=" + executableFilePath)
	result := hardware.AcquireSingletonSocket()
	logger.Logger.Debug("acquireSingletonLock result err=" + fmt.Sprintf("%v", result))
	return result
}
