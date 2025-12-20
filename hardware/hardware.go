// Package hardware
// File:        hardware.go
// Url:         https://github.com/xiang-tai-duo/go-bootstrap/blob/master/hardware/hardware.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Hardware provides utility methods for hardware information retrieval, including hardware fingerprint generation for license registration
// --------------------------------------------------------------------------------
package hardware

import (
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/xiang-tai-duo/go-bootstrap/hash"
)

//goland:noinspection GoSnakeCaseUsage,SpellCheckingInspection
const (
	ARGUMENT_C                         = "-Command"
	ARGUMENT_CPUINFO_FILE              = "/proc/cpuinfo"
	ARGUMENT_DISKUTIL_INFO             = "info"
	ARGUMENT_DISKUTIL_ROOT             = "/"
	ARGUMENT_FORMAT_LIST               = "/format:list"
	ARGUMENT_FSUTIL_C_DRIVE            = "C:\\"
	ARGUMENT_FSUTIL_FSINFO             = "fsinfo"
	ARGUMENT_FSUTIL_VOLUMEINFO         = "volumeinfo"
	ARGUMENT_GET_CIMINSTANCE_BIOS      = "Get-CimInstance Win32_BIOS | Select-Object -ExpandProperty SerialNumber -First 1"
	ARGUMENT_GET_CIMINSTANCE_PROCESSOR = "Get-CimInstance Win32_Processor | Select-Object -ExpandProperty Name -First 1"
	ARGUMENT_GET_VOLUME_C              = "Get-Volume | Where-Object { $_.DriveLetter -eq 'C' } | Select-Object -ExpandProperty FileSystemLabel"
	ARGUMENT_GET_WMIOBJECT_BASEBOARD   = "Get-WmiObject Win32_BaseBoard | Select-Object -ExpandProperty SerialNumber -First 1"
	ARGUMENT_GET_WMIOBJECT_DISKDRIVE   = "Get-WmiObject Win32_DiskDrive | Select-Object -ExpandProperty SerialNumber -First 1"
	ARGUMENT_GET_WMIOBJECT_PROCESSOR   = "Get-WmiObject Win32_Processor | Select-Object -ExpandProperty ProcessorId -First 1"
	ARGUMENT_IOREG_L                   = "-l"
	ARGUMENT_LSBLK_NO_SERIAL           = "-no"
	ARGUMENT_MACHDEP_CPU_BRAND_STRING  = "machdep.cpu.brand_string"
	ARGUMENT_MACHDEP_CPU_SIGNATURE     = "machdep.cpu.signature"
	ARGUMENT_MODEL_NAME                = "model name"
	ARGUMENT_NAME                      = "Name"
	ARGUMENT_BLKID_O_VALUE             = "-o"
	ARGUMENT_BLKID_S_UUID              = "-s"
	ARGUMENT_BLKID_UUID                = "UUID"
	ARGUMENT_BLKID_VALUE               = "value"
	ARGUMENT_PROCESSOR_ID              = "ProcessorId"
	ARGUMENT_PRODUCT_SERIAL_FILE       = "/sys/devices/virtual/dmi/id/product_serial"
	ARGUMENT_SERIAL                    = "serial"
	ARGUMENT_SERIAL_NUMBER             = "SerialNumber"
	ARGUMENT_SERIAL_NUMBER_FILE        = "/sys/devices/virtual/dmi/id/board_serial"
	ARGUMENT_SPHARDWARE_DATATYPE       = "SPHardwareDataType"
	ARGUMENT_SYSCTL_N                  = "-n"
	ARGUMENT_VOLUME_SERIAL_NUMBER      = "VolumeSerialNumber"
	COMMAND_BLKID                      = "blkid"
	COMMAND_DISKUTIL                   = "diskutil"
	COMMAND_FSUTIL                     = "fsutil"
	COMMAND_IOREG                      = "ioreg"
	COMMAND_LSBLK                      = "lsblk"
	COMMAND_POWERSHELL                 = "powershell"
	COMMAND_SYSTEM_PROFILER            = "system_profiler"
	COMMAND_SYSCTL                     = "sysctl"
	COMMAND_WMIC_BASEBOARD             = "wmic"
	COMMAND_WMIC_BIOS                  = "wmic"
	COMMAND_WMIC_CPU                   = "wmic"
	COMMAND_WMIC_DISKDRIVE             = "wmic"
	COMMAND_WMIC_LOGICALDISK           = "wmic"
	DARWIN                             = "darwin"
	DEVICE_MEDIA_NAME                  = "Device / Media Name:"
	ENVIRONMENT_USER                   = "USER"
	ENVIRONMENT_USERNAME               = "USERNAME"
	IOPLATFORM_SERIAL_NUMBER           = "IOPlatformSerialNumber"
	LINUX                              = "linux"
	NAME_EQUAL                         = "Name="
	PROCESSOR_ID_EQUAL                 = "ProcessorId="
	SERIAL_NUMBER_EQUAL                = "SerialNumber="
	SERIAL_NUMBER_SYSTEM               = "Serial Number (system):"
	VOLUME_SERIAL_NUMBER               = "Volume Serial Number"
	VOLUME_SERIAL_NUMBER_EQUAL         = "VolumeSerialNumber="
	VOLUME_UUID_COLON                  = "Volume UUID:"
	WINDOWS                            = "windows"
)

func GetBoardInfo() string {
	result := ""
	switch runtime.GOOS {
	case WINDOWS:
		result = GetBoardInfoWindows()
	case LINUX:
		result = GetBoardInfoLinux()
	case DARWIN:
		result = GetBoardInfoDarwin()
	default:
		result = ""
	}
	return result
}

func GetBoardInfoDarwin() string {
	result := ""
	cmd := exec.Command(COMMAND_SYSTEM_PROFILER, ARGUMENT_SPHARDWARE_DATATYPE)
	if output, err := cmd.Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, SERIAL_NUMBER_SYSTEM) {
				if parts := strings.Split(line, ":"); len(parts) > 1 {
					result = strings.TrimSpace(parts[1])
					break
				}
			}
		}
	}
	if result == "" {
		cmd := exec.Command(COMMAND_IOREG, ARGUMENT_IOREG_L)
		if output, err := cmd.Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, IOPLATFORM_SERIAL_NUMBER) {
					if parts := strings.Split(line, "="); len(parts) > 1 {
						result = strings.TrimSpace(strings.ReplaceAll(parts[1], "\"", ""))
						break
					}
				}
			}
		}
	}
	return result
}

func GetBoardInfoLinux() string {
	result := ""
	if data, err := os.ReadFile(ARGUMENT_SERIAL_NUMBER_FILE); err == nil {
		result = strings.TrimSpace(string(data))
	}
	if result == "" {
		if data, err := os.ReadFile(ARGUMENT_PRODUCT_SERIAL_FILE); err == nil {
			result = strings.TrimSpace(string(data))
		}
	}
	return result
}

//goland:noinspection DuplicatedCode
func GetBoardInfoWindows() string {
	result := ""
	cmd := exec.Command(COMMAND_WMIC_BASEBOARD, "baseboard", "get", ARGUMENT_SERIAL_NUMBER, ARGUMENT_FORMAT_LIST)
	if output, err := cmd.Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, SERIAL_NUMBER_EQUAL) {
				result = strings.TrimSpace(strings.TrimPrefix(line, SERIAL_NUMBER_EQUAL))
				break
			}
		}
	}
	if result == "" {
		cmd := exec.Command(COMMAND_POWERSHELL, ARGUMENT_C, ARGUMENT_GET_WMIOBJECT_BASEBOARD)
		if output, err := cmd.Output(); err == nil {
			result = strings.TrimSpace(string(output))
		}
	}
	if result == "" {
		cmd := exec.Command(COMMAND_POWERSHELL, ARGUMENT_C, ARGUMENT_GET_CIMINSTANCE_BIOS)
		if output, err := cmd.Output(); err == nil {
			result = strings.TrimSpace(string(output))
		}
	}
	if result == "" {
		cmd := exec.Command(COMMAND_WMIC_BIOS, "bios", "get", ARGUMENT_SERIAL_NUMBER, ARGUMENT_FORMAT_LIST)
		if output, err := cmd.Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, SERIAL_NUMBER_EQUAL) {
					result = strings.TrimSpace(strings.TrimPrefix(line, SERIAL_NUMBER_EQUAL))
					break
				}
			}
		}
	}
	return result
}

func GetCPUInfo() string {
	result := ""
	switch runtime.GOOS {
	case WINDOWS:
		result = GetCPUInfoWindows()
	case LINUX:
		result = GetCPUInfoLinux()
	case DARWIN:
		result = GetCPUInfoDarwin()
	default:
		result = ""
	}
	return result
}

func GetCPUInfoDarwin() string {
	result := ""
	cmd := exec.Command(COMMAND_SYSCTL, ARGUMENT_SYSCTL_N, ARGUMENT_MACHDEP_CPU_SIGNATURE)
	if output, err := cmd.Output(); err == nil {
		result = strings.TrimSpace(string(output))
	}
	if result == "" {
		cmd := exec.Command(COMMAND_SYSCTL, ARGUMENT_SYSCTL_N, ARGUMENT_MACHDEP_CPU_BRAND_STRING)
		if output, err := cmd.Output(); err == nil {
			result = strings.TrimSpace(string(output))
		}
	}
	return result
}

func GetCPUInfoFromEnv() string {
	result := ""
	processorCount := runtime.NumCPU()
	arch := runtime.GOARCH
	result = arch + "_" + strconv.Itoa(processorCount)
	return result
}

func GetCPUInfoLinux() string {
	result := ""
	if data, err := os.ReadFile(ARGUMENT_CPUINFO_FILE); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.Contains(line, ARGUMENT_SERIAL) {
				if parts := strings.Split(line, ":"); len(parts) > 1 {
					result = strings.TrimSpace(parts[1])
					break
				}
			}
		}
		if result == "" {
			for _, line := range lines {
				if strings.Contains(line, ARGUMENT_MODEL_NAME) {
					if parts := strings.Split(line, ":"); len(parts) > 1 {
						result = strings.TrimSpace(parts[1])
						break
					}
				}
			}
		}
	}
	return result
}

//goland:noinspection DuplicatedCode
func GetCPUInfoWindows() string {
	result := ""
	cmd := exec.Command(COMMAND_WMIC_CPU, "cpu", "get", ARGUMENT_PROCESSOR_ID, ARGUMENT_FORMAT_LIST)
	if output, err := cmd.Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, PROCESSOR_ID_EQUAL) {
				result = strings.TrimSpace(strings.TrimPrefix(line, PROCESSOR_ID_EQUAL))
				break
			}
		}
	}
	if result == "" {
		cmd := exec.Command(COMMAND_POWERSHELL, ARGUMENT_C, ARGUMENT_GET_WMIOBJECT_PROCESSOR)
		if output, err := cmd.Output(); err == nil {
			result = strings.TrimSpace(string(output))
		}
	}
	if result == "" {
		cmd := exec.Command(COMMAND_POWERSHELL, ARGUMENT_C, ARGUMENT_GET_CIMINSTANCE_PROCESSOR)
		if output, err := cmd.Output(); err == nil {
			result = strings.TrimSpace(string(output))
		}
	}
	if result == "" {
		cmd := exec.Command(COMMAND_WMIC_CPU, "cpu", "get", ARGUMENT_NAME, ARGUMENT_FORMAT_LIST)
		if output, err := cmd.Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, NAME_EQUAL) {
					result = strings.TrimSpace(strings.TrimPrefix(line, NAME_EQUAL))
					break
				}
			}
		}
	}
	if result == "" {
		result = GetCPUInfoFromEnv()
	}
	return result
}

func GetDiskSerial() string {
	result := ""
	switch runtime.GOOS {
	case WINDOWS:
		result = GetDiskSerialWindows()
	case LINUX:
		result = GetDiskSerialLinux()
	case DARWIN:
		result = GetDiskSerialDarwin()
	default:
		result = ""
	}
	return result
}

//goland:noinspection DuplicatedCode
func GetDiskSerialDarwin() string {
	result := ""
	cmd := exec.Command(COMMAND_DISKUTIL, ARGUMENT_DISKUTIL_INFO, ARGUMENT_DISKUTIL_ROOT)
	if output, err := cmd.Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, DEVICE_MEDIA_NAME) {
				if parts := strings.Split(line, ":"); len(parts) > 1 {
					result = strings.TrimSpace(parts[1])
					break
				}
			}
		}
	}
	if result == "" {
		cmd := exec.Command(COMMAND_DISKUTIL, ARGUMENT_DISKUTIL_INFO, ARGUMENT_DISKUTIL_ROOT)
		if output, err := cmd.Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, VOLUME_UUID_COLON) {
					if parts := strings.Split(line, ":"); len(parts) > 1 {
						result = strings.TrimSpace(parts[1])
						break
					}
				}
			}
		}
	}
	return result
}

func GetDiskSerialLinux() string {
	result := ""
	cmd := exec.Command(COMMAND_LSBLK, ARGUMENT_LSBLK_NO_SERIAL, "serial")
	if output, err := cmd.Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if line := strings.TrimSpace(line); line != "" {
				result = line
				break
			}
		}
	}
	if result == "" {
		cmd := exec.Command(COMMAND_BLKID, ARGUMENT_BLKID_S_UUID, ARGUMENT_BLKID_UUID, ARGUMENT_BLKID_O_VALUE, ARGUMENT_BLKID_VALUE)
		if output, err := cmd.Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if line := strings.TrimSpace(line); line != "" {
					result = line
					break
				}
			}
		}
	}
	return result
}

//goland:noinspection DuplicatedCode,SpellCheckingInspection
func GetDiskSerialWindows() string {
	result := ""
	cmd := exec.Command(COMMAND_WMIC_DISKDRIVE, "diskdrive", "get", ARGUMENT_SERIAL_NUMBER, ARGUMENT_FORMAT_LIST)
	if output, err := cmd.Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, SERIAL_NUMBER_EQUAL) {
				result = strings.TrimSpace(strings.TrimPrefix(line, SERIAL_NUMBER_EQUAL))
				break
			}
		}
	}
	if result == "" {
		cmd := exec.Command(COMMAND_POWERSHELL, ARGUMENT_C, ARGUMENT_GET_WMIOBJECT_DISKDRIVE)
		if output, err := cmd.Output(); err == nil {
			result = strings.TrimSpace(string(output))
		}
	}
	if result == "" {
		cmd := exec.Command(COMMAND_POWERSHELL, ARGUMENT_C, ARGUMENT_GET_VOLUME_C)
		if output, err := cmd.Output(); err == nil {
			result = strings.TrimSpace(string(output))
		}
	}
	if result == "" {
		cmd := exec.Command(COMMAND_WMIC_LOGICALDISK, "logicaldisk", "get", ARGUMENT_VOLUME_SERIAL_NUMBER, ARGUMENT_FORMAT_LIST)
		if output, err := cmd.Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, VOLUME_SERIAL_NUMBER_EQUAL) {
					result = strings.TrimSpace(strings.TrimPrefix(line, VOLUME_SERIAL_NUMBER_EQUAL))
					break
				}
			}
		}
	}
	if result == "" {
		cmd := exec.Command(COMMAND_FSUTIL, ARGUMENT_FSUTIL_FSINFO, ARGUMENT_FSUTIL_VOLUMEINFO, ARGUMENT_FSUTIL_C_DRIVE)
		if output, err := cmd.Output(); err == nil {
			outputStr := string(output)
			if strings.Contains(outputStr, VOLUME_SERIAL_NUMBER) {
				lines := strings.Split(outputStr, "\n")
				for _, line := range lines {
					if strings.Contains(line, VOLUME_SERIAL_NUMBER) {
						if parts := strings.Split(line, ":"); len(parts) > 1 {
							result = strings.TrimSpace(parts[1])
							break
						}
					}
				}
			}
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetHardwareFingerprint() string {
	fingerprint := ""
	components := make([]string, 0)
	if cpuInfo := GetCPUInfo(); cpuInfo != "" {
		components = append(components, cpuInfo)
	}
	if macAddress := GetMACAddress(); macAddress != "" {
		components = append(components, macAddress)
	}
	if diskSerial := GetDiskSerial(); diskSerial != "" {
		components = append(components, diskSerial)
	}
	if boardInfo := GetBoardInfo(); boardInfo != "" {
		components = append(components, boardInfo)
	}
	if systemInfo := GetSystemInfo(); systemInfo != "" {
		components = append(components, systemInfo)
	}
	raw := strings.Join(components, "|")
	fingerprint = hash.MD5(raw)
	return fingerprint
}

//goland:noinspection SpellCheckingInspection
func GetMACAddress() string {
	result := ""

	interfaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range interfaces {
			if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
				if mac := iface.HardwareAddr.String(); mac != "" {
					result = mac
					break
				}
			}
		}
	}
	if result == "" && err == nil {
		for _, iface := range interfaces {
			if iface.Flags&net.FlagLoopback == 0 {
				if mac := iface.HardwareAddr.String(); mac != "" {
					result = mac
					break
				}
			}
		}
	}
	return result
}

func GetSystemInfo() string {
	var components []string

	if hostname, err := os.Hostname(); err == nil && hostname != "" {
		components = append(components, hostname)
	}

	username := os.Getenv(ENVIRONMENT_USERNAME)
	if username == "" {
		username = os.Getenv(ENVIRONMENT_USER)
	}
	if username != "" {
		components = append(components, username)
	}

	osName := runtime.GOOS
	components = append(components, osName)

	osArch := runtime.GOARCH
	components = append(components, osArch)

	cpuCount := strconv.Itoa(runtime.NumCPU())
	components = append(components, cpuCount)

	return strings.Join(components, "|")
}
