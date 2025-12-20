//go:build windows

// Package advapi32
// File:        advapi32.go
// Url:         https://github.com/xiang-tai-duo/go-bootstrap/blob/master/windows/advapi32/advapi32.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: advapi32.dll wrapper for Windows registry API functions
// --------------------------------------------------------------------------------
package advapi32

/*
#include <windows.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/xiang-tai-duo/go-bootstrap/logger"
)

//goland:noinspection GoSnakeCaseUsage
const (
	WCHAR_SIZE = 2
)

//goland:noinspection GoUnusedConst,GoSnakeCaseUsage,SpellCheckingInspection
const (
	HKEY_CLASSES_ROOT                = 0x80000000
	HKEY_CURRENT_USER                = 0x80000001
	HKEY_LOCAL_MACHINE               = 0x80000002
	HKEY_USERS                       = 0x80000003
	HKEY_PERFORMANCE_DATA            = 0x80000004
	HKEY_CURRENT_CONFIG              = 0x80000005
	HKEY_DYN_DATA                    = 0x80000006
	HKEY_CURRENT_USER_LOCAL_SETTINGS = 0x80000007
	HKEY_PERFORMANCE_TEXT            = 0x80000050
	HKEY_PERFORMANCE_NLSTEXT         = 0x80000060
	KEY_QUERY_VALUE                  = 0x0001
	KEY_SET_VALUE                    = 0x0002
	KEY_CREATE_SUB_KEY               = 0x0004
	KEY_ENUMERATE_SUB_KEYS           = 0x0008
	KEY_NOTIFY                       = 0x0010
	KEY_CREATE_LINK                  = 0x0020
	KEY_WOW64_32KEY                  = 0x0200
	KEY_WOW64_64KEY                  = 0x0100
	KEY_WOW64_RES                    = 0x0300
	KEY_READ                         = 0x20019
	KEY_WRITE                        = 0x20006
	KEY_EXECUTE                      = 0x20019
	KEY_ALL_ACCESS                   = 0xF003F
	REG_NONE                         = 0
	REG_SZ                           = 1
	REG_EXPAND_SZ                    = 2
	REG_BINARY                       = 3
	REG_DWORD                        = 4
	REG_DWORD_BIG_ENDIAN             = 5
	REG_LINK                         = 6
	REG_MULTI_SZ                     = 7
	REG_RESOURCE_LIST                = 8
	REG_FULL_RESOURCE_DESCRIPTOR     = 9
	REG_RESOURCE_REQUIREMENTS_LIST   = 10
	REG_QWORD                        = 11
	REG_UNKNOWN                      = 0xFFFFFFFF
	ERROR_SUCCESS                    = 0
	ERROR_NO_MORE_ITEMS              = 259
	MAX_VALUE_NAME_LENGTH            = 256
	MAX_VALUE_DATA_LENGTH            = 16384
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

//goland:noinspection GoUnusedExportedFunction,GoSnakeCaseUsage,GoVetUnsafePointer
func RegOpenKeyEx(hKey uint, subKey string, options uint32, samDesired uint32) uint32 {
	var result C.HKEY
	subKeyUtf16, _ := syscall.UTF16FromString(subKey)
	C.RegOpenKeyExW(C.HKEY(unsafe.Pointer(uintptr(hKey))), (*C.WCHAR)(unsafe.Pointer(&subKeyUtf16[0])), C.DWORD(options), C.REGSAM(samDesired), (*C.HKEY)(unsafe.Pointer(&result)))
	return uint32(uintptr(unsafe.Pointer(result)))
}

//goland:noinspection GoUnusedExportedFunction,GoSnakeCaseUsage,GoVetUnsafePointer
func RegEnumValue(hKey uint32, index uint32) (valueName string, value []byte, valueType uint32) {
	var nameLen C.DWORD = C.DWORD(MAX_VALUE_NAME_LENGTH)
	var dataLen C.DWORD = C.DWORD(MAX_VALUE_DATA_LENGTH)
	var regType C.DWORD
	nameBuffer := make([]uint16, MAX_VALUE_NAME_LENGTH)
	valueBuffer := make([]byte, MAX_VALUE_DATA_LENGTH)
	if err := int32(C.RegEnumValueW(C.HKEY(unsafe.Pointer(uintptr(hKey))), C.DWORD(index), (*C.WCHAR)(unsafe.Pointer(&nameBuffer[0])), &nameLen, nil, &regType, (*C.BYTE)(unsafe.Pointer(&valueBuffer[0])), &dataLen)); err == ERROR_SUCCESS {
		valueName = syscall.UTF16ToString(nameBuffer[:nameLen])
		value = valueBuffer[:dataLen]
		valueType = uint32(regType)
	} else {
		valueType = REG_UNKNOWN
	}
	return
}

//goland:noinspection GoUnusedExportedFunction,GoSnakeCaseUsage,GoVetUnsafePointer
func RegCloseKey(hKey uint32) int32 {
	return int32(C.RegCloseKey(C.HKEY(unsafe.Pointer(uintptr(hKey)))))
}
