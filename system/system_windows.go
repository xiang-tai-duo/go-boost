//go:build windows

// Package system
// File:        system_windows.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/system/system_windows.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: SYSTEM is a wrapper for operating system interactions, providing Windows specific implementations for architecture detection and long path name resolution via kernel32.
// --------------------------------------------------------------------------------

package system

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/xiang-tai-duo/go-boost/hash"
	"github.com/xiang-tai-duo/go-boost/logger"
	"github.com/xiang-tai-duo/go-boost/windows/kernel32"
	"golang.org/x/sys/windows"
)

//goland:noinspection GoSnakeCaseUsage
const (
	WINDOWS_SINGLETON_MUTEX_PREFIX = `Local\go-boost-singleton-`
)

//goland:noinspection GoUnusedGlobalVariable
var (
	singletonMutex windows.Handle
)

//goland:noinspection GoUnusedExportedFunction
func Architecture() string {
	return kernel32.GetNativeSystemInfo()
}

func GetLongPathName(path string) (string, error) {
	logger.Logger.Debug(fmt.Sprintf("GetLongPathName enter path=%s", path))
	result := ""
	err := error(nil)
	result, err = kernel32.GetLongPathNameW(path)
	logger.Logger.Debug(fmt.Sprintf("GetLongPathName result=%s err=%v", result, err))
	return result, err
}

//goland:noinspection GoUnhandledErrorResult
func acquireSingletonLock(executableFilePath string) error {
	logger.Logger.Debug(fmt.Sprintf("acquireSingletonLock enter executableFilePath=%s", executableFilePath))
	result := error(nil)
	err := error(nil)
	mutexName := WINDOWS_SINGLETON_MUTEX_PREFIX + strings.ToLower(hash.SHA3(filepath.Clean(executableFilePath)))
	logger.Logger.Debug(fmt.Sprintf("acquireSingletonLock mutexName=%s", mutexName))
	var name *uint16
	if name, err = windows.UTF16PtrFromString(mutexName); err == nil {
		var handle windows.Handle
		if handle, err = windows.CreateMutex(nil, false, name); err == nil {
			singletonMutex = handle
		} else {
			if errors.Is(err, windows.ERROR_ALREADY_EXISTS) {
				if handle != 0 {
					windows.CloseHandle(handle)
				}
				result = fmt.Errorf("singleton mutex already exists: %s", mutexName)
			} else {
				result = err
			}
		}
	} else {
		result = err
	}
	logger.Logger.Debug(fmt.Sprintf("acquireSingletonLock result err=%v", result))
	return result
}
