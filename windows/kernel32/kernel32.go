//go:build windows

// Package kernel32
// File:        kernel32.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/windows/kernel32/kernel32.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: kernel32.dll wrapper for Windows API functions
// --------------------------------------------------------------------------------
package kernel32

/*
#include <windows.h>
#include <winioctl.h>

static WORD goGetNativeProcessorArchitecture() {
    SYSTEM_INFO si;
    GetNativeSystemInfo(&si);
    return si.wProcessorArchitecture;
}
*/
import "C"
import (
	"errors"
	"fmt"
	"log"
	"syscall"
	"unsafe"

	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	ERROR_CSTRING_FAILED         = "CString failed"
	ERROR_GET_LONG_PATH_FAILED   = "GetLongPathNameW failed"
	ERROR_GET_WINDOWS_DIR_FAILED = "GetWindowsDirectoryW failed"
	ERROR_PATH_EMPTY             = "path is empty"
	MAX_PATH                     = 260
	WCHAR_SIZE                   = 2
	COMPRESSION_FORMAT_NONE      = 0x0000
	COMPRESSION_FORMAT_DEFAULT   = 0x0001
	COMPRESSION_FORMAT_LZNT1     = 0x0002
	MODULE_NAME_KERNEL32         = "windows.kernel32"
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_KERNEL32, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_KERNEL32, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_KERNEL32, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection DuplicatedCode,GoUnusedExportedFunction
func CString(s string) *C.wchar_t {
	result := (*C.wchar_t)(nil)
	if utf16, err := syscall.UTF16FromString(s); err == nil {
		if ptr := C.malloc(C.size_t((len(utf16) + 1) * int(unsafe.Sizeof(uint16(0))))); ptr != nil {
			src := utf16
			dst := unsafe.Slice((*uint16)(ptr), len(utf16)+1)
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
		log.Printf("GoString: unsupported type %T", lpwsz)
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction,GoSnakeCaseUsage
func CheckRemoteDebuggerPresent() bool {
	result := false
	var present C.BOOL
	if C.CheckRemoteDebuggerPresent(C.GetCurrentProcess(), &present) != 0 {
		result = present != 0
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction,GoSnakeCaseUsage
func GetLongPathNameW(path string) (string, error) {
	result := ""
	err := error(nil)
	if path == "" {
		err = errors.New(ERROR_PATH_EMPTY)
	} else {
		if p := CString(path); p == nil {
			err = errors.New(ERROR_CSTRING_FAILED)
		} else {
			defer FreeCString(p)
			buffer := make([]uint16, MAX_PATH)
			n := uint32(C.GetLongPathNameW(p, (*C.WCHAR)(unsafe.Pointer(&buffer[0])), C.DWORD(len(buffer))))
			if n == 0 {
				err = errors.New(ERROR_GET_LONG_PATH_FAILED)
			} else if int(n) > len(buffer) {
				buffer = make([]uint16, n)
				n = uint32(C.GetLongPathNameW(p, (*C.WCHAR)(unsafe.Pointer(&buffer[0])), C.DWORD(len(buffer))))
				if n == 0 {
					err = errors.New(ERROR_GET_LONG_PATH_FAILED)
				} else {
					result = GoString(buffer)
				}
			} else {
				result = GoString(buffer)
			}
		}
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction,GoSnakeCaseUsage
func GetNativeSystemInfo() string {
	result := ""
	switch C.goGetNativeProcessorArchitecture() {
	case C.PROCESSOR_ARCHITECTURE_AMD64:
		result = "amd64"
	case C.PROCESSOR_ARCHITECTURE_ARM:
		result = "arm"
	case C.PROCESSOR_ARCHITECTURE_ARM64:
		result = "arm64"
	case C.PROCESSOR_ARCHITECTURE_IA64:
		result = "ia64"
	case C.PROCESSOR_ARCHITECTURE_INTEL:
		result = "386"
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction,GoSnakeCaseUsage
func GetWindowsDirectoryW() (string, error) {
	result := ""
	err := error(nil)
	buf := make([]uint16, MAX_PATH)
	n := uint32(C.GetWindowsDirectoryW((*C.WCHAR)(unsafe.Pointer(&buf[0])), C.UINT(len(buf))))
	if n == 0 {
		err = errors.New(ERROR_GET_WINDOWS_DIR_FAILED)
	} else if int(n) > len(buf) {
		buf = make([]uint16, n)
		n = uint32(C.GetWindowsDirectoryW((*C.WCHAR)(unsafe.Pointer(&buf[0])), C.UINT(len(buf))))
		if n == 0 {
			err = errors.New(ERROR_GET_WINDOWS_DIR_FAILED)
		} else {
			result = GoString(buf)
		}
	} else {
		result = GoString(buf)
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction,GoSnakeCaseUsage
func IsDebuggerPresent() bool {
	return C.IsDebuggerPresent() != 0
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func SetFileCompression(filePath string, enable bool) error {
	err := error(nil)
	if filePath == "" {
		err = fmt.Errorf("file path is empty")
	} else {
		if filePathW := CString(filePath); filePathW == nil {
			err = fmt.Errorf("failed to convert file path to wide string: %s", filePath)
		} else {
			defer FreeCString(filePathW)
			hFile := C.CreateFileW(
				filePathW,
				C.GENERIC_READ|C.GENERIC_WRITE,
				C.FILE_SHARE_READ|C.FILE_SHARE_WRITE,
				nil,
				C.OPEN_EXISTING,
				C.FILE_FLAG_BACKUP_SEMANTICS,
				nil,
			)
			if hFile == C.HANDLE(C.INVALID_HANDLE_VALUE) {
				err = fmt.Errorf("failed to open file: %s, error code: %d", filePath, uint32(C.GetLastError()))
			} else {
				defer C.CloseHandle(hFile)
				format := C.USHORT(COMPRESSION_FORMAT_NONE)
				if enable {
					format = C.USHORT(COMPRESSION_FORMAT_DEFAULT)
				}
				var bytesReturned C.DWORD
				if C.DeviceIoControl(
					hFile,
					C.FSCTL_SET_COMPRESSION,
					C.LPVOID(unsafe.Pointer(&format)),
					C.DWORD(unsafe.Sizeof(format)),
					nil,
					0,
					&bytesReturned,
					nil,
				) == C.FALSE {
					err = fmt.Errorf("failed to set compression on file: %s, error code: %d", filePath, uint32(C.GetLastError()))
				}
			}
		}
	}
	return err
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_KERNEL32, logger.SKIP_STACK_FRAMES_BASE)
}
