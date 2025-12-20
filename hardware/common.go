// Package hardware
// File:        common.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/hardware/common.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: HARDWARE is a wrapper for hardware information collection, providing cross-platform methods to retrieve BIOS, CPU, disk, memory, MAC and fingerprint info.
// --------------------------------------------------------------------------------
package hardware

import (
	"crypto/sha3"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/xiang-tai-duo/go-boost/logger"
)

const (
	FORMAT_NUMBER                = "%d"
	HARDWARE_FINGERPRINT_SALT    = "qGkNvR7xLpZ2wYtB8sHmC3dJf01eKuXo"
	HARDWARE_INFO_FORMAT         = "%s%s: %s%s"
	HARDWARE_INFO_LABEL_BIOS     = "BIOS Info"
	HARDWARE_INFO_LABEL_COMPUTER = "Computer Name"
	HARDWARE_INFO_LABEL_CPU      = "CPU Count"
	HARDWARE_INFO_LABEL_DISK     = "Disk Serial"
	HARDWARE_INFO_LABEL_HOSTNAME = "Hostname"
	HARDWARE_INFO_LABEL_MAC      = "MAC Addresses"
	HARDWARE_INFO_LABEL_MEMORY   = "Physical Memory"
	LINE_SEPARATOR               = "\n"
	MAC_ADDRESS_SEPARATOR        = ", "
	MODULE_NAME_HARDWARE         = "hardware"
)

func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_HARDWARE, logger.SKIP_STACK_FRAMES_BASE)
}

func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_HARDWARE, logger.SKIP_STACK_FRAMES_BASE)
}

func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_HARDWARE, logger.SKIP_STACK_FRAMES_BASE)
}

func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_HARDWARE, logger.SKIP_STACK_FRAMES_BASE)
}

func GetHardwareFingerprint() string {
	logger.Logger.Debug("GetHardwareFingerprint enter")
	result := GetHardwareFingerprintEx(HARDWARE_FINGERPRINT_SALT)
	logger.Logger.Debug(fmt.Sprintf("GetHardwareFingerprint result=%s", result))
	return result
}

func GetHardwareFingerprintEx(salt string) string {
	logger.Logger.Debug(fmt.Sprintf("GetHardwareFingerprintEx enter salt=%s", salt))
	result := getSHA3(GetHardwareInfo() + salt)
	logger.Logger.Debug(fmt.Sprintf("GetHardwareFingerprintEx result=%s", result))
	return result
}

func GetHardwareInfo() string {
	logger.Logger.Debug("GetHardwareInfo enter")
	result := ""
	result = fmt.Sprintf(HARDWARE_INFO_FORMAT, result, HARDWARE_INFO_LABEL_BIOS, getBIOSInfo(), LINE_SEPARATOR)
	result = fmt.Sprintf(HARDWARE_INFO_FORMAT, result, HARDWARE_INFO_LABEL_COMPUTER, getComputerName(), LINE_SEPARATOR)
	result = fmt.Sprintf(HARDWARE_INFO_FORMAT, result, HARDWARE_INFO_LABEL_CPU, getCPUInfo(), LINE_SEPARATOR)
	result = fmt.Sprintf(HARDWARE_INFO_FORMAT, result, HARDWARE_INFO_LABEL_DISK, getDiskSerial(), LINE_SEPARATOR)
	result = fmt.Sprintf(HARDWARE_INFO_FORMAT, result, HARDWARE_INFO_LABEL_HOSTNAME, getHostname(), LINE_SEPARATOR)
	result = fmt.Sprintf(HARDWARE_INFO_FORMAT, result, HARDWARE_INFO_LABEL_MAC, getMACAddresses(), LINE_SEPARATOR)
	result = fmt.Sprintf(HARDWARE_INFO_FORMAT, result, HARDWARE_INFO_LABEL_MEMORY, getMemoryInfo(), LINE_SEPARATOR)
	logger.Logger.Debug(fmt.Sprintf("GetHardwareInfo result=%s", result))
	return result
}

func getCPUInfo() string {
	logger.Logger.Debug("getCPUInfo enter")
	result := fmt.Sprintf(FORMAT_NUMBER, runtime.NumCPU())
	logger.Logger.Debug(fmt.Sprintf("getCPUInfo result=%s", result))
	return result
}

func getHostname() string {
	logger.Logger.Debug("getHostname enter")
	result := ""
	if hostname, err := os.Hostname(); err == nil {
		result = hostname
	}
	logger.Logger.Debug(fmt.Sprintf("getHostname result=%s", result))
	return result
}

func getMACAddresses() string {
	logger.Logger.Debug("getMACAddresses enter")
	result := ""
	macAddresses := make([]string, 0)
	if interfaces, err := net.Interfaces(); err == nil {
		for _, i := range interfaces {
			if i.Flags&net.FlagUp != 0 && i.Flags&net.FlagLoopback == 0 {
				if mac := i.HardwareAddr.String(); mac != "" {
					macAddresses = append(macAddresses, mac)
				}
			}
		}
	}
	sort.Strings(macAddresses)
	result = strings.Join(macAddresses, MAC_ADDRESS_SEPARATOR)
	logger.Logger.Debug(fmt.Sprintf("getMACAddresses result=%s", result))
	return result
}

func getSHA3(input string) string {
	logger.Logger.Debug(fmt.Sprintf("getSHA3 enter inputLength=%d", len(input)))
	hash := sha3.Sum256([]byte(input))
	result := strings.ToLower(hex.EncodeToString(hash[:]))
	logger.Logger.Debug(fmt.Sprintf("getSHA3 result=%s", result))
	return result
}
