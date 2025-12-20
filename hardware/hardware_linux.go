//go:build linux

// Package hardware
// File:        hardware_linux.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/hardware/hardware_linux.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: HARDWARE is a wrapper for hardware information collection, providing Linux specific implementations via sysfs and syscall.
// --------------------------------------------------------------------------------

package hardware

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/xiang-tai-duo/go-boost/logger"
)

const (
	LINUX_MACHINE_ID_PATH             = "/etc/machine-id"
	LINUX_SYSFS_BLOCK_PATH            = "/sys/block"
	LINUX_SYSFS_DEVICE_SERIAL_PATH    = "device/serial"
	LINUX_SYSFS_DMI_PRODUCT_UUID_PATH = "/sys/class/dmi/id/product_uuid"
	LINUX_VIRTUAL_BLOCK_PREFIX_LOOP   = "loop"
	LINUX_VIRTUAL_BLOCK_PREFIX_RAM    = "ram"
)

func getBIOSInfo() string {
	logger.Logger.Debug("getBIOSInfo enter")
	result := ""
	if data, err := os.ReadFile(LINUX_SYSFS_DMI_PRODUCT_UUID_PATH); err == nil {
		result = strings.TrimSpace(string(data))
	}
	if result == "" {
		if data, err := os.ReadFile(LINUX_MACHINE_ID_PATH); err == nil {
			result = strings.TrimSpace(string(data))
		}
	}
	logger.Logger.Debug(fmt.Sprintf("getBIOSInfo resultLength=%d", len(result)))
	return result
}

func getComputerName() string {
	logger.Logger.Debug("getComputerName enter")
	result := getHostname()
	logger.Logger.Debug(fmt.Sprintf("getComputerName result=%s", result))
	return result
}

func getDiskSerial() string {
	logger.Logger.Debug("getDiskSerial enter")
	result := ""
	if entries, err := os.ReadDir(LINUX_SYSFS_BLOCK_PATH); err == nil {
		for _, entry := range entries {
			name := entry.Name()
			if !strings.HasPrefix(name, LINUX_VIRTUAL_BLOCK_PREFIX_LOOP) && !strings.HasPrefix(name, LINUX_VIRTUAL_BLOCK_PREFIX_RAM) {
				if data, readError := os.ReadFile(filepath.Join(LINUX_SYSFS_BLOCK_PATH, name, LINUX_SYSFS_DEVICE_SERIAL_PATH)); readError == nil {
					if serial := strings.TrimSpace(string(data)); serial != "" {
						result = serial
					}
				}
			}
		}
	}
	logger.Logger.Debug(fmt.Sprintf("getDiskSerial result=%s", result))
	return result
}

func getMemoryInfo() string {
	logger.Logger.Debug("getMemoryInfo enter")
	result := fmt.Sprintf(FORMAT_NUMBER, 0)
	info := syscall.Sysinfo_t{}
	if err := syscall.Sysinfo(&info); err == nil {
		result = fmt.Sprintf(FORMAT_NUMBER, uint64(info.Totalram)*uint64(info.Unit))
	}
	logger.Logger.Debug(fmt.Sprintf("getMemoryInfo result=%s", result))
	return result
}
