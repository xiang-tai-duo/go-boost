//go:build windows

// Package file
// File:        file_windows.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/file/file_windows.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: File provides Windows-specific file locking implementation that denies other processes both read and write access
// --------------------------------------------------------------------------------
package file

import (
	"encoding/binary"
	"os"

	"golang.org/x/sys/windows"
)

const (
	DOS_HEADER_SIZE   = 64
	PE_HEADER_OFFSET  = 0x3C
	PE_SIGNATURE      = "PE\x00\x00"
	PE_SIGNATURE_SIZE = 4
)

func IsExecutable(filePath string) bool {
	result := false
	fileInfo, err := os.Stat(filePath)
	if err == nil {
		if !fileInfo.IsDir() {
			var file *os.File
			if file, err = os.Open(filePath); err == nil {
				defer file.Close()
				dosBuffer := make([]byte, DOS_HEADER_SIZE)
				if n, err := file.Read(dosBuffer); err == nil && n >= DOS_HEADER_SIZE {
					if string(dosBuffer[:2]) == "MZ" {
						peOffset := int64(binary.LittleEndian.Uint32(dosBuffer[PE_HEADER_OFFSET:]))
						if _, err = file.Seek(peOffset, 0); err == nil {
							peBuffer := make([]byte, PE_SIGNATURE_SIZE)
							if n, err = file.Read(peBuffer); err == nil && n >= PE_SIGNATURE_SIZE {
								if string(peBuffer) == PE_SIGNATURE {
									result = true
								}
							}
						}
					}
				}
			}
		}
	}
	return result
}

func OpenLockedRead(path string) (*os.File, error) {
	return openExclusiveWindows(path, windows.GENERIC_READ)
}

func OpenLockedReadOnly(path string) (*os.File, error) {
	return openExclusiveWindows(path, windows.GENERIC_READ)
}

func openExclusiveWindows(path string, access uint32) (*os.File, error) {
	result := (*os.File)(nil)
	err := error(nil)
	if path == "" {
		err = os.ErrInvalid
	} else {
		var pathPointer *uint16
		if pathPointer, err = windows.UTF16PtrFromString(path); err == nil {
			var handle windows.Handle
			if handle, err = windows.CreateFile(pathPointer, access, 0, nil, windows.OPEN_EXISTING, windows.FILE_ATTRIBUTE_NORMAL, 0); err == nil {
				result = os.NewFile(uintptr(handle), path)
			}
		}
	}
	return result, err
}
