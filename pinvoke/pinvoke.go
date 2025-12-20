// Package pinvoke
// File:        pinvoke.go
// Url:         https://github.com/xiang-tai-duo/go-bootstrap/blob/master/pinvoke/pinvoke.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Platform invoke functionality for calling native code
// --------------------------------------------------------------------------------
package pinvoke

/*
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/xiang-tai-duo/go-bootstrap/logger"
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	WCHAR_SIZE = 2
)

//goland:noinspection DuplicatedCode,GoUnusedExportedFunction
func CString(s string) *C.wchar_t {
	result := (*C.wchar_t)(nil)
	if utf16, err := syscall.UTF16FromString(s); err == nil {
		if ptr := C.malloc(C.size_t((len(utf16) + 1) * int(unsafe.Sizeof(uint16(0))))); ptr != nil {
			src := (*[1 << 30]uint16)(unsafe.Pointer(&utf16[0]))[:len(utf16):len(utf16)]
			dst := (*[1 << 30]uint16)(ptr)[: len(utf16)+1 : len(utf16)+1]
			copy(dst, src)
			dst[len(utf16)] = 0
			result = (*C.wchar_t)(ptr)
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func FreeCString(ptr *C.wchar_t) {
	if ptr != nil {
		C.free(unsafe.Pointer(ptr))
	}
}

//goland:noinspection SpellCheckingInspection,DuplicatedCode,GoUnusedExportedFunction
func GoString(lpwsz interface{}) string {
	result := ""
	if ptr, ok := lpwsz.(*C.wchar_t); ok {
		if pwsz := (*uint16)(unsafe.Pointer(ptr)); pwsz != nil {
			result = syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(pwsz))[:])
		}
	} else if ptr, ok := lpwsz.(*uint16); ok {
		if ptr != nil {
			result = syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(ptr))[:])
		}
	} else if ptr, ok := lpwsz.(C.LPWSTR); ok {
		if ptr != nil {
			result = syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(ptr))[:])
		}
	} else if ptr, ok := lpwsz.(C.LPSTR); ok {
		if ptr != nil {
			result = syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(ptr))[:])
		}
	} else if arr, ok := lpwsz.(*[32]uint16); ok {
		result = syscall.UTF16ToString(arr[:])
	} else if arr, ok := lpwsz.(*[64]uint16); ok {
		result = syscall.UTF16ToString(arr[:])
	} else if arr, ok := lpwsz.(*[256]uint16); ok {
		result = syscall.UTF16ToString(arr[:])
	} else if slice, ok := lpwsz.([]uint16); ok {
		result = syscall.UTF16ToString(slice)
	} else {
		logger.Logger.Error(fmt.Sprintf("GoString: unsupported type %T", lpwsz))
	}
	return result
}
