//go:build linux

package advapi32

func CloseServiceHandle(hSCObject uintptr) bool {
	return false
}

func CreateServiceW(hSCManager uintptr, serviceName string, displayName string, dwDesiredAccess uint32, dwServiceType uint32, dwStartType uint32, dwErrorControl uint32, binaryPathName string, loadOrderGroup string, lpdwTagId *uint32, dependencies string, serviceStartName string, password string) uintptr {
	return 0
}

func DeleteService(hService uintptr) bool {
	return false
}

func GetLastError() uint32 {
	return 0
}

func OpenSCManagerW(machineName string, databaseName string, dwDesiredAccess uint32) uintptr {
	return 0
}

func OpenServiceW(hSCManager uintptr, serviceName string, dwDesiredAccess uint32) uintptr {
	return 0
}

func QueryServiceConfigW(hService uintptr) (binaryPathName string, startType uint32, err uint32) {
	return "", 0, 0
}

func QueryServiceStatus(hService uintptr) (currentState uint32, err uint32) {
	return 0, 0
}

func RegCloseKey(hKey uint32) int32 {
	return 0
}

func RegEnumValue(hKey uint32, index uint32) (valueName string, value []byte, valueType uint32) {
	return "", nil, 0
}

func RegOpenKeyEx(hKey uint, subKey string, options uint32, samDesired uint32) uint32 {
	return 0
}

func StartServiceCtrlDispatcherW(lpServiceStartTable *uintptr) bool {
	return false
}

func StartServiceW(hService uintptr) bool {
	return false
}
