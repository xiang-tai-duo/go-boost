//go:build windows

// Package defender
// File:        mpdefender.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/windows/defender/mpdefender.go
// Author:      TRAE.AI
// Created:     2026/06/24 00:00:00
// Description: Windows Defender exclusion management using Registry API
// --------------------------------------------------------------------------------

package defender

/*
#cgo CFLAGS: -DWIN32_LEAN_AND_MEAN
#cgo LDFLAGS: -ladvapi32
#include "mpdefender.h"
*/
import "C"
import (
	"fmt"
	"path/filepath"
	"syscall"
	"unsafe"

	"github.com/xiang-tai-duo/go-boost/logger"
)

const (
	ERROR_MPDEFENDER_ADD_FAILED    = "failed to add Windows Defender exclusion"
	ERROR_MPDEFENDER_REMOVE_FAILED = "failed to remove Windows Defender exclusion"
	ERROR_MPDEFENDER_INVALID_PATH  = "path cannot be empty"
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	MODULE_NAME_MPDEFENDER = "windows.defender.mpdefender"
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_MPDEFENDER, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_MPDEFENDER, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_MPDEFENDER, logger.SKIP_STACK_FRAMES_BASE)
}

func AddFileToWhiteList(filePath string) error {
	return addPathToWhiteList(filePath, 0)
}

func AddFolderToWhiteList(folderPath string) error {
	return addPathToWhiteList(folderPath, 0)
}

//goland:noinspection GoUnusedExportedFunction
func RemoveFileFromWhiteList(filePath string) error {
	return removePathFromWhiteList(filePath, 0)
}

//goland:noinspection GoUnusedExportedFunction
func RemoveFolderFromWhiteList(folderPath string) error {
	return removePathFromWhiteList(folderPath, 0)
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_MPDEFENDER, logger.SKIP_STACK_FRAMES_BASE)
}

func addPathToWhiteList(path string, exclusionType int) error {
	result := error(nil)
	err := error(nil)
	if path == "" {
		err = fmt.Errorf(ERROR_MPDEFENDER_INVALID_PATH)
	} else {
		absolutePath := ""
		if absolutePath, err = filepath.Abs(path); err == nil {
			pathPtr, _ := syscall.UTF16PtrFromString(absolutePath)
			pathW := (*C.wchar_t)(unsafe.Pointer(pathPtr))
			hr := C.MpDefenderAddExclusion(pathW)
			if hr != C.MPDEFENDER_S_OK {
				err = fmt.Errorf("%s: HRESULT=0x%08X", ERROR_MPDEFENDER_ADD_FAILED, uint32(hr))
			}
		}
		result = err
	}
	return result
}

func removePathFromWhiteList(path string, exclusionType int) error {
	result := error(nil)
	err := error(nil)
	if path == "" {
		err = fmt.Errorf(ERROR_MPDEFENDER_INVALID_PATH)
	} else {
		absolutePath := ""
		if absolutePath, err = filepath.Abs(path); err == nil {
			pathPtr, _ := syscall.UTF16PtrFromString(absolutePath)
			pathW := (*C.wchar_t)(unsafe.Pointer(pathPtr))
			hr := C.MpDefenderRemoveExclusion(pathW)
			if hr != C.MPDEFENDER_S_OK {
				err = fmt.Errorf("%s: HRESULT=0x%08X", ERROR_MPDEFENDER_REMOVE_FAILED, uint32(hr))
			}
		}
		result = err
	}
	return result
}
