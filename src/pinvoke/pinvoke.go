// Package pinvoke
// File:        pinvoke.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/src/pinvoke/pinvoke.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: Platform invoke functionality for calling native code
// --------------------------------------------------------------------------------
package pinvoke

/*

 */
import "C"
import (
	"fmt"
	"syscall"
	"unsafe"
)

type (
	PINVOKE struct {
	}
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	WCHAR_SIZE = 2
)

//goland:noinspection GoUnusedGlobalVariable
var (
	Pinvoke = PINVOKE{}
)

//goland:noinspection SpellCheckingInspection,GoUnusedFunction,GoUnusedExportedFunction,GoUnusedParameter
func GoStringW(lpwsz *C.wchar_t) string {
	fmt.Println(`
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
	`)
	sz := ""
	if pwsz := (*uint16)(unsafe.Pointer(lpwsz)); pwsz != nil {
		for ptr := pwsz; *ptr != 0; ptr = (*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + WCHAR_SIZE)) {
			sz += string(rune(*ptr))
		}
	}
	return sz
}

//goland:noinspection GoSnakeCaseUsage,GoUnusedFunction,GoUnusedExportedFunction,GoUnusedParameter
func CStringW(s string) *C.wchar_t {
	fmt.Println(`
		//goland:noinspection GoSnakeCaseUsage,GoUnusedFunction,GoUnusedExportedFunction
		func CStringW(s string) *C.wchar_t {
			var wsz *C.wchar_t = nil
			if ptr, err := syscall.UTF16FromString(s); err == nil {
				wsz = (*C.wchar_t)(unsafe.Pointer(&ptr[0]))
			}
			return wsz
		}
	`)
	var wsz *C.wchar_t = nil
	if ptr, err := syscall.UTF16FromString(s); err == nil {
		wsz = (*C.wchar_t)(unsafe.Pointer(&ptr[0]))
	}
	return wsz
}
