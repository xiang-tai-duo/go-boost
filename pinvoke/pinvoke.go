// Package pinvoke
// File:        pinvoke.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/pinvoke/pinvoke.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Platform invoke functionality for calling native code
// --------------------------------------------------------------------------------
package pinvoke

/*
#include <stdlib.h>
#include <windows.h>
*/
import "C"
import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	BUFFER_SIZE_256     = 256
	BUFFER_SIZE_32      = 32
	BUFFER_SIZE_64      = 64
	MAX_ARRAY_INDEX_20  = 1 << 20
	MODULE_NAME_PINVOKE = "pinvoke"
	NULL_TERMINATOR     = 0
	WCHAR_NULL_PADDING  = 1
	WCHAR_SIZE          = 2
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_PINVOKE, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_PINVOKE, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_PINVOKE, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection DuplicatedCode,GoUnusedExportedFunction
func CString(s string) *C.wchar_t {
	result := (*C.wchar_t)(nil)
	if utf16, err := syscall.UTF16FromString(s); err == nil {
		if ptr := C.malloc(C.size_t((len(utf16) + WCHAR_NULL_PADDING) * int(unsafe.Sizeof(uint16(0))))); ptr != nil {
			src := utf16
			dst := unsafe.Slice((*uint16)(ptr), len(utf16)+WCHAR_NULL_PADDING)
			copy(dst, src)
			dst[len(utf16)] = NULL_TERMINATOR
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
			result = syscall.UTF16ToString((*[MAX_ARRAY_INDEX_20]uint16)(unsafe.Pointer(pwsz))[:])
		}
	} else if ptr, ok := lpwsz.(*uint16); ok {
		if ptr != nil {
			result = syscall.UTF16ToString((*[MAX_ARRAY_INDEX_20]uint16)(unsafe.Pointer(ptr))[:])
		}
	} else if ptr, ok := lpwsz.(C.LPWSTR); ok {
		if ptr != nil {
			result = syscall.UTF16ToString((*[MAX_ARRAY_INDEX_20]uint16)(unsafe.Pointer(ptr))[:])
		}
	} else if ptr, ok := lpwsz.(C.LPSTR); ok {
		if ptr != nil {
			result = syscall.UTF16ToString((*[MAX_ARRAY_INDEX_20]uint16)(unsafe.Pointer(ptr))[:])
		}
	} else if arr, ok := lpwsz.(*[BUFFER_SIZE_32]uint16); ok {
		result = syscall.UTF16ToString(arr[:])
	} else if arr, ok := lpwsz.(*[BUFFER_SIZE_64]uint16); ok {
		result = syscall.UTF16ToString(arr[:])
	} else if arr, ok := lpwsz.(*[BUFFER_SIZE_256]uint16); ok {
		result = syscall.UTF16ToString(arr[:])
	} else if slice, ok := lpwsz.([]uint16); ok {
		result = syscall.UTF16ToString(slice)
	} else {
		__error(fmt.Sprintf("GoString: unsupported type %T", lpwsz))
	}
	return result
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_PINVOKE, logger.SKIP_STACK_FRAMES_BASE)
}
