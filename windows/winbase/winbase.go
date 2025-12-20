//go:build windows

// Package winbase
// File:        winbase.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/windows/winbase.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: winbase is a wrapper for Windows base operations, providing methods for memory and string operations.
// --------------------------------------------------------------------------------
package winbase

/*
#include <windows.h>
#include <strsafe.h>

#cgo LDFLAGS: -lkernel32
*/
import "C"

import (
	"syscall"
	"unsafe"

	"github.com/xiang-tai-duo/go-boost/pinvoke"
)

//goland:noinspection SpellCheckingInspection,GoUnusedFunction,GoUnusedExportedFunction
func GoStringW(lpwsz *C.wchar_t) string {
	sz := ""
	if pwsz := (*uint16)(unsafe.Pointer(lpwsz)); pwsz != nil {
		for ptr := pwsz; *ptr != 0; ptr = (*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + pinvoke.WCHAR_SIZE)) {
			sz += string(rune(*ptr))
		}
	}
	return sz
}

//goland:noinspection GoSnakeCaseUsage,GoUnusedFunction,GoUnusedExportedFunction
func CStringW(s string) *C.wchar_t {
	var wsz *C.wchar_t = nil
	if ptr, err := syscall.UTF16FromString(s); err == nil {
		wsz = (*C.wchar_t)(unsafe.Pointer(&ptr[0]))
	}
	return wsz
}

//goland:noinspection GoUnusedExportedFunction
func IsBadStringPtr(p *uint16) bool {
	result := true
	if p != nil {
		result = C.IsBadStringPtrW((*C.WCHAR)(unsafe.Pointer(p)), C.STRSAFE_MAX_CCH) == C.TRUE
	}
	return result
}
