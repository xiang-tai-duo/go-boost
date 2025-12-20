//go:build windows

// Package logger
// File:        logger_windows.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/logger/logger_windows.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Logger service and disk space helpers for Windows.
// --------------------------------------------------------------------------------
package logger

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	FSCTL_SET_COMPRESSION      = 0x0009C040
	COMPRESSION_FORMAT_NONE    = 0x0000
	COMPRESSION_FORMAT_DEFAULT = 0x0001
)

//goland:noinspection GoUnusedExportedFunction
func IsWindowsService() bool {
	result, _ := svc.IsWindowsService()
	return result
}

//goland:noinspection SpellCheckingInspection
func getDiskFreeSpace(directory string) (uint64, error) {
	result := uint64(0)
	err := error(nil)
	var lpszDirectory *uint16
	if lpszDirectory, err = windows.UTF16PtrFromString(directory); err == nil {
		err = windows.GetDiskFreeSpaceEx(lpszDirectory, &result, nil, nil)
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func SetFileCompression(filePath string, enable bool) error {
	err := error(nil)
	if filePath == "" {
		err = fmt.Errorf("file path is empty")
	} else {
		var filePathW *uint16
		if filePathW, err = windows.UTF16PtrFromString(filePath); err == nil {
			hFile := windows.Handle(0)
			if hFile, err = windows.CreateFile(
				filePathW,
				windows.GENERIC_READ|windows.GENERIC_WRITE,
				windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE,
				nil,
				windows.OPEN_EXISTING,
				windows.FILE_FLAG_BACKUP_SEMANTICS,
				0,
			); err == nil {
				defer func(h windows.Handle) {
					_ = windows.CloseHandle(h)
				}(hFile)
				format := uint16(COMPRESSION_FORMAT_NONE)
				if enable {
					format = uint16(COMPRESSION_FORMAT_DEFAULT)
				}
				var bytesReturned uint32
				err = windows.DeviceIoControl(
					hFile,
					FSCTL_SET_COMPRESSION,
					(*byte)(unsafe.Pointer(&format)),
					uint32(unsafe.Sizeof(format)),
					nil,
					0,
					&bytesReturned,
					nil,
				)
			}
		}
	}
	return err
}
