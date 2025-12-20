//go:build windows

// Package kernel32
// File:        kernel32.go
// Url:         https://github.com/xiang-tai-duo/go-bootstrap/blob/master/windows/kernel32/kernel32.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: kernel32.dll wrapper for Windows API functions
// --------------------------------------------------------------------------------
package kernel32

/*
#include <windows.h>
*/
import "C"
import (
	"unsafe"
)

const (
	MAX_PATH = 260
)

//goland:noinspection GoUnusedExportedFunction,GoSnakeCaseUsage
func GetWindowsDirectoryW(lpBuffer *uint16, uSize uint32) uint32 {
	return uint32(C.GetWindowsDirectoryW((*C.WCHAR)(unsafe.Pointer(lpBuffer)), C.UINT(uSize)))
}
