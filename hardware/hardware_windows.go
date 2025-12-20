//go:build windows

// Package hardware
// File:        hardware_windows.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/hardware/hardware_windows.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: HARDWARE is a wrapper for hardware information collection, providing Windows specific implementations via kernel32 syscalls.
// --------------------------------------------------------------------------------

package hardware

import (
	"encoding/hex"
	"fmt"
	"syscall"
	"unsafe"

	"github.com/xiang-tai-duo/go-boost/logger"
)

const (
	BIOS_BUFFER_MAX              = 64
	BIOS_BUFFER_OFFSET           = 8
	BUFFER_SIZE_256              = 256
	FIRMWARE_TABLE_PROVIDER_RSMB = 0x52534D42
	WINDOWS_DEFAULT_ROOT_PATH    = "C:\\"
)

type MEMORY_STATUS_EX struct {
	Length               uint32
	MemoryLoad           uint32
	TotalPhys            uint64
	AvailPhys            uint64
	TotalPageFile        uint64
	AvailPageFile        uint64
	TotalVirtual         uint64
	AvailVirtual         uint64
	AvailExtendedVirtual uint64
}

var (
	kernel32               = syscall.NewLazyDLL("kernel32.dll")
	getComputerNameW       = kernel32.NewProc("GetComputerNameW")
	getSystemFirmwareTable = kernel32.NewProc("GetSystemFirmwareTable")
	getVolumeInformationW  = kernel32.NewProc("GetVolumeInformationW")
	globalMemoryStatusEx   = kernel32.NewProc("GlobalMemoryStatusEx")
)

func getBIOSInfo() string {
	logger.Logger.Debug("getBIOSInfo enter")
	result := ""
	size, _, _ := getSystemFirmwareTable.Call(
		uintptr(FIRMWARE_TABLE_PROVIDER_RSMB),
		uintptr(0),
		uintptr(0),
		uintptr(0),
	)
	logger.Logger.Debug(fmt.Sprintf("getBIOSInfo firmware table size=%d", size))
	if size > 0 {
		buffer := make([]byte, size)
		_, _, _ = getSystemFirmwareTable.Call(
			uintptr(FIRMWARE_TABLE_PROVIDER_RSMB),
			uintptr(0),
			uintptr(unsafe.Pointer(&buffer[0])),
			size,
		)
		if len(buffer) > BIOS_BUFFER_OFFSET {
			end := BIOS_BUFFER_OFFSET + BIOS_BUFFER_MAX
			if end > len(buffer) {
				end = len(buffer)
			}
			result = hex.EncodeToString(buffer[BIOS_BUFFER_OFFSET:end])
		}
	}
	logger.Logger.Debug(fmt.Sprintf("getBIOSInfo resultLength=%d", len(result)))
	return result
}

func getComputerName() string {
	logger.Logger.Debug("getComputerName enter")
	result := ""
	buffer := make([]uint16, BUFFER_SIZE_256)
	size := uint32(len(buffer))
	getComputerNameW.Call(
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&size)),
	)
	result = syscall.UTF16ToString(buffer)
	logger.Logger.Debug(fmt.Sprintf("getComputerName result=%s", result))
	return result
}

func getDiskSerial() string {
	logger.Logger.Debug("getDiskSerial enter")
	result := ""
	volumeName := make([]uint16, BUFFER_SIZE_256)
	serialNumber := uint32(0)
	fileSystemName := make([]uint16, BUFFER_SIZE_256)
	rootPath, _ := syscall.UTF16PtrFromString(WINDOWS_DEFAULT_ROOT_PATH)
	maximumComponentLength := uint32(0)
	fileSystemFlags := uint32(0)
	getVolumeInformationW.Call(
		uintptr(unsafe.Pointer(rootPath)),
		uintptr(unsafe.Pointer(&volumeName[0])),
		uintptr(len(volumeName)),
		uintptr(unsafe.Pointer(&serialNumber)),
		uintptr(unsafe.Pointer(&maximumComponentLength)),
		uintptr(unsafe.Pointer(&fileSystemFlags)),
		uintptr(unsafe.Pointer(&fileSystemName[0])),
		uintptr(len(fileSystemName)),
	)
	result = fmt.Sprintf(FORMAT_NUMBER, serialNumber)
	logger.Logger.Debug(fmt.Sprintf("getDiskSerial result=%s", result))
	return result
}

func getMemoryInfo() string {
	logger.Logger.Debug("getMemoryInfo enter")
	result := ""
	memoryStatus := MEMORY_STATUS_EX{}
	memoryStatus.Length = uint32(unsafe.Sizeof(memoryStatus))
	globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memoryStatus)))
	result = fmt.Sprintf(FORMAT_NUMBER, memoryStatus.TotalPhys)
	logger.Logger.Debug(fmt.Sprintf("getMemoryInfo result=%s", result))
	return result
}
