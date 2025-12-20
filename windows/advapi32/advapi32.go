//go:build windows

// Package advapi32
// File:        advapi32.go
// Url:         https://github.com/xiang-tai-duo/go-bootstrap/blob/master/windows/advapi32/advapi32.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: advapi32.dll wrapper for Windows registry and service API functions
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

	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst,SpellCheckingInspection
const (
	ERROR_FAILED_SERVICE_CONTROLLER_CONNECT = 1063
	ERROR_NO_MORE_ITEMS                     = 259
	ERROR_SERVICE_DISABLED                  = 1058
	ERROR_SERVICE_DOES_NOT_EXIST            = 1060
	ERROR_SUCCESS                           = 0
	HKEY_CLASSES_ROOT                       = 0x80000000
	HKEY_CURRENT_CONFIG                     = 0x80000005
	HKEY_CURRENT_USER                       = 0x80000001
	HKEY_CURRENT_USER_LOCAL_SETTINGS        = 0x80000007
	HKEY_DYN_DATA                           = 0x80000006
	HKEY_LOCAL_MACHINE                      = 0x80000002
	HKEY_PERFORMANCE_DATA                   = 0x80000004
	HKEY_PERFORMANCE_NLSTEXT                = 0x80000060
	HKEY_PERFORMANCE_TEXT                   = 0x80000050
	HKEY_USERS                              = 0x80000003
	KEY_ALL_ACCESS                          = 0xF003F
	KEY_CREATE_LINK                         = 0x0020
	KEY_CREATE_SUB_KEY                      = 0x0004
	KEY_ENUMERATE_SUB_KEYS                  = 0x0008
	KEY_EXECUTE                             = 0x20019
	KEY_NOTIFY                              = 0x0010
	KEY_QUERY_VALUE                         = 0x0001
	KEY_READ                                = 0x20019
	KEY_SET_VALUE                           = 0x0002
	KEY_WOW64_32KEY                         = 0x0200
	KEY_WOW64_64KEY                         = 0x0100
	KEY_WOW64_RES                           = 0x0300
	KEY_WRITE                               = 0x20006
	MAX_VALUE_DATA_LENGTH                   = 16384
	MAX_VALUE_NAME_LENGTH                   = 256
	MODULE_NAME_ADVAPI32                    = "windows.advapi32"
	REG_BINARY                              = 3
	REG_DWORD                               = 4
	REG_DWORD_BIG_ENDIAN                    = 5
	REG_EXPAND_SZ                           = 2
	REG_FULL_RESOURCE_DESCRIPTOR            = 9
	REG_LINK                                = 6
	REG_MULTI_SZ                            = 7
	REG_NONE                                = 0
	REG_QWORD                               = 11
	REG_RESOURCE_LIST                       = 8
	REG_RESOURCE_REQUIREMENTS_LIST          = 10
	REG_SZ                                  = 1
	REG_UNKNOWN                             = 0xFFFFFFFF
	SC_MANAGER_ALL_ACCESS                   = 0xF003F
	SC_MANAGER_CONNECT                      = 0x0001
	SC_MANAGER_CREATE_SERVICE               = 0x0002
	SC_MANAGER_ENUMERATE_SERVICE            = 0x0004
	SC_MANAGER_LOCK                         = 0x0008
	SC_MANAGER_MODIFY_BOOT_CONFIG           = 0x0020
	SC_MANAGER_QUERY_LOCK_STATUS            = 0x0010
	SERVICE_ALL_ACCESS                      = 0xF01FF
	SERVICE_AUTO_START                      = 0x00000002
	SERVICE_CHANGE_CONFIG                   = 0x0002
	SERVICE_CONFIG_DELAYED_AUTO_START_INFO  = 3
	SERVICE_CONFIG_DESCRIPTION              = 1
	SERVICE_CONFIG_FAILURE_ACTIONS          = 2
	SERVICE_CONFIG_FAILURE_ACTIONS_FLAG     = 4
	SERVICE_CONFIG_PREFERRED_NODE           = 9
	SERVICE_CONFIG_PRESHUTDOWN_INFO         = 7
	SERVICE_CONFIG_REQUIRED_PRIVILEGES_INFO = 6
	SERVICE_CONFIG_SERVICE_SID_INFO         = 5
	SERVICE_CONFIG_TRIGGER_INFO             = 8
	SERVICE_CONTINUE_PENDING                = 5
	SERVICE_DEMAND_START                    = 0x00000003
	SERVICE_DISABLED                        = 0x00000004
	SERVICE_ENUMERATE_DEPENDENTS            = 0x0008
	SERVICE_ERROR_NORMAL                    = 0x00000001
	SERVICE_INTERROGATE                     = 0x0080
	SERVICE_PAUSE_CONTINUE                  = 0x0040
	SERVICE_PAUSE_PENDING                   = 6
	SERVICE_PAUSED                          = 7
	SERVICE_QUERY_STATUS                    = 0x0004
	SERVICE_RUNNING                         = 4
	SERVICE_START                           = 0x0010
	SERVICE_START_PENDING                   = 2
	SERVICE_STOP                            = 0x0020
	SERVICE_STOP_PENDING                    = 3
	SERVICE_STOPPED                         = 1
	SERVICE_WIN32_OWN_PROCESS               = 0x00000010
	WCHAR_SIZE                              = 2
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_ADVAPI32, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_ADVAPI32, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_ADVAPI32, logger.SKIP_STACK_FRAMES_BASE)
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
		__error(fmt.Sprintf("GoString: unsupported type %T", lpwsz))
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction,GoSnakeCaseUsage,GoVetUnsafePointer
func CloseServiceHandle(hSCObject uintptr) bool {
	return C.CloseServiceHandle(C.SC_HANDLE(unsafe.Pointer(hSCObject))) != 0
}

//goland:noinspection GoUnusedExportedFunction,GoSnakeCaseUsage,GoVetUnsafePointer
func CreateServiceW(hSCManager uintptr, serviceName string, displayName string, dwDesiredAccess uint32, dwServiceType uint32, dwStartType uint32, dwErrorControl uint32, binaryPathName string, loadOrderGroup string, lpdwTagId *uint32, dependencies string, serviceStartName string, password string) uintptr {
	result := uintptr(0)
	serviceNameUtf16, _ := syscall.UTF16FromString(serviceName)
	displayNameUtf16, _ := syscall.UTF16FromString(displayName)
	binaryPathNameUtf16, _ := syscall.UTF16FromString(binaryPathName)
	loadOrderGroupUtf16 := (*uint16)(nil)
	if loadOrderGroup != "" {
		loadOrderGroupArr, _ := syscall.UTF16FromString(loadOrderGroup)
		loadOrderGroupUtf16 = &loadOrderGroupArr[0]
	}
	dependenciesUtf16 := (*uint16)(nil)
	if dependencies != "" {
		dependenciesArr, _ := syscall.UTF16FromString(dependencies)
		dependenciesUtf16 = &dependenciesArr[0]
	}
	serviceStartNameUtf16 := (*uint16)(nil)
	if serviceStartName != "" {
		serviceStartNameArr, _ := syscall.UTF16FromString(serviceStartName)
		serviceStartNameUtf16 = &serviceStartNameArr[0]
	}
	passwordUtf16 := (*uint16)(nil)
	if password != "" {
		passwordArr, _ := syscall.UTF16FromString(password)
		passwordUtf16 = &passwordArr[0]
	}
	result = uintptr(unsafe.Pointer(C.CreateServiceW(
		C.SC_HANDLE(unsafe.Pointer(hSCManager)),
		(*C.WCHAR)(unsafe.Pointer(&serviceNameUtf16[0])),
		(*C.WCHAR)(unsafe.Pointer(&displayNameUtf16[0])),
		C.DWORD(dwDesiredAccess),
		C.DWORD(dwServiceType),
		C.DWORD(dwStartType),
		C.DWORD(dwErrorControl),
		(*C.WCHAR)(unsafe.Pointer(&binaryPathNameUtf16[0])),
		(*C.WCHAR)(unsafe.Pointer(loadOrderGroupUtf16)),
		(*C.DWORD)(unsafe.Pointer(lpdwTagId)),
		(*C.WCHAR)(unsafe.Pointer(dependenciesUtf16)),
		(*C.WCHAR)(unsafe.Pointer(serviceStartNameUtf16)),
		(*C.WCHAR)(unsafe.Pointer(passwordUtf16)),
	)))
	return result
}

//goland:noinspection GoUnusedExportedFunction,GoSnakeCaseUsage,GoVetUnsafePointer
func DeleteService(hService uintptr) bool {
	return C.DeleteService(C.SC_HANDLE(unsafe.Pointer(hService))) != 0
}

//goland:noinspection GoSnakeCaseUsage
func GetLastError() uint32 {
	return uint32(C.GetLastError())
}

//goland:noinspection GoUnusedExportedFunction,GoSnakeCaseUsage,GoVetUnsafePointer
func OpenSCManagerW(machineName string, databaseName string, dwDesiredAccess uint32) uintptr {
	result := uintptr(0)
	machineNameUtf16 := (*uint16)(nil)
	if machineName != "" {
		machineNameArr, _ := syscall.UTF16FromString(machineName)
		machineNameUtf16 = &machineNameArr[0]
	}
	databaseNameUtf16 := (*uint16)(nil)
	if databaseName != "" {
		databaseNameArr, _ := syscall.UTF16FromString(databaseName)
		databaseNameUtf16 = &databaseNameArr[0]
	}
	result = uintptr(unsafe.Pointer(C.OpenSCManagerW(
		(*C.WCHAR)(unsafe.Pointer(machineNameUtf16)),
		(*C.WCHAR)(unsafe.Pointer(databaseNameUtf16)),
		C.DWORD(dwDesiredAccess),
	)))
	return result
}

//goland:noinspection GoUnusedExportedFunction,GoSnakeCaseUsage,GoVetUnsafePointer
func OpenServiceW(hSCManager uintptr, serviceName string, dwDesiredAccess uint32) uintptr {
	serviceNameUtf16, _ := syscall.UTF16FromString(serviceName)
	return uintptr(unsafe.Pointer(C.OpenServiceW(
		C.SC_HANDLE(unsafe.Pointer(hSCManager)),
		(*C.WCHAR)(unsafe.Pointer(&serviceNameUtf16[0])),
		C.DWORD(dwDesiredAccess),
	)))
}

//goland:noinspection GoSnakeCaseUsage,GoVetUnsafePointer
func QueryServiceConfigW(hService uintptr) (binaryPathName string, startType uint32, err uint32) {
	var bytesNeeded C.DWORD
	C.QueryServiceConfigW(C.SC_HANDLE(unsafe.Pointer(hService)), nil, 0, &bytesNeeded)
	if bytesNeeded > 0 {
		buffer := make([]byte, bytesNeeded)
		if C.QueryServiceConfigW(C.SC_HANDLE(unsafe.Pointer(hService)), (*C.QUERY_SERVICE_CONFIGW)(unsafe.Pointer(&buffer[0])), bytesNeeded, &bytesNeeded) != 0 {
			config := (*C.QUERY_SERVICE_CONFIGW)(unsafe.Pointer(&buffer[0]))
			binaryPathName = GoString(config.lpBinaryPathName)
			startType = uint32(config.dwStartType)
			err = ERROR_SUCCESS
		} else {
			err = GetLastError()
		}
	} else {
		err = GetLastError()
	}
	return
}

//goland:noinspection GoSnakeCaseUsage,GoVetUnsafePointer
func QueryServiceStatus(hService uintptr) (currentState uint32, err uint32) {
	var status C.SERVICE_STATUS
	if C.QueryServiceStatus(C.SC_HANDLE(unsafe.Pointer(hService)), &status) != 0 {
		currentState = uint32(status.dwCurrentState)
		err = ERROR_SUCCESS
	} else {
		err = GetLastError()
	}
	return
}

//goland:noinspection GoUnusedExportedFunction,GoSnakeCaseUsage,GoVetUnsafePointer
func RegCloseKey(hKey uint32) int32 {
	return int32(C.RegCloseKey(C.HKEY(unsafe.Pointer(uintptr(hKey)))))
}

//goland:noinspection GoUnusedExportedFunction,GoSnakeCaseUsage,GoVetUnsafePointer
func RegEnumValue(hKey uint32, index uint32) (valueName string, value []byte, valueType uint32) {
	nameLen := C.DWORD(MAX_VALUE_NAME_LENGTH)
	dataLen := C.DWORD(MAX_VALUE_DATA_LENGTH)
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
func RegOpenKeyEx(hKey uint, subKey string, options uint32, samDesired uint32) uint32 {
	var result C.HKEY
	subKeyUtf16, _ := syscall.UTF16FromString(subKey)
	C.RegOpenKeyExW(C.HKEY(unsafe.Pointer(uintptr(hKey))), (*C.WCHAR)(unsafe.Pointer(&subKeyUtf16[0])), C.DWORD(options), C.REGSAM(samDesired), (*C.HKEY)(unsafe.Pointer(&result)))
	return uint32(uintptr(unsafe.Pointer(result)))
}

//goland:noinspection GoSnakeCaseUsage,GoVetUnsafePointer
func StartServiceCtrlDispatcherW(lpServiceStartTable *C.SERVICE_TABLE_ENTRYW) bool {
	return C.StartServiceCtrlDispatcherW(lpServiceStartTable) != 0
}

//goland:noinspection GoSnakeCaseUsage,GoVetUnsafePointer
func StartServiceW(hService uintptr) bool {
	return C.StartServiceW(C.SC_HANDLE(unsafe.Pointer(hService)), 0, nil) != 0
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_ADVAPI32, logger.SKIP_STACK_FRAMES_BASE)
}
