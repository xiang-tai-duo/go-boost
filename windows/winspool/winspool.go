//go:build windows

// Package winspool
// File:        winspool.go
// Url:         https://github.com/xiang-tai-duo/go-bootstrap/blob/master/windows/winspool/winspool.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: winspool is a wrapper for Windows print spooler operations, providing methods for printer management and paper information retrieval.
// --------------------------------------------------------------------------------
package winspool

/*
#include <windows.h>
#include <winspool.h>
#include <stdlib.h>
#cgo LDFLAGS: -lwinspool

#define NAME_SIZE 64

*/
import "C"

import (
	"fmt"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/xiang-tai-duo/go-bootstrap/logger"
)

//goland:noinspection GoSnakeCaseUsage
type (
	PAPER_SIZE struct {
		Width  int32
		Length int32
	}
	PAPER_INFO struct {
		Id        uint16
		PaperName string
		Size      PAPER_SIZE
	}
	PRINTER_INFO struct {
		Name      string
		PortName  string
		IsDefault bool
	}
	PRINTER_INFO_2 struct {
		pServerName         *uint16
		pPrinterName        *uint16
		pShareName          *uint16
		pPortName           *uint16
		pDriverName         *uint16
		pComment            *uint16
		pLocation           *uint16
		pDevMode            *byte
		pSepFile            *uint16
		pPrintProcessor     *uint16
		pDatatype           *uint16
		pParameters         *uint16
		pSecurityDescriptor *byte
		Attributes          uint32
		Priority            uint32
		DefaultPriority     uint32
		StartTime           uint32
		UntilTime           uint32
		Status              uint32
		cJobs               uint32
		AveragePPM          uint32
	}
	DRIVER_INFO_3 struct {
		cVersion         uint32
		pName            *uint16
		pEnvironment     *uint16
		pDriverPath      *uint16
		pDataFile        *uint16
		pConfigFile      *uint16
		pHelpFile        *uint16
		pDependentFiles  *uint16
		pMonitorName     *uint16
		pDefaultDataType *uint16
	}
	PORT_INFO_2 struct {
		pPortName    *uint16
		pMonitorName *uint16
		pDescription *uint16
		fPortType    uint32
		Reserved     uint32
	}
	JOB_INFO struct {
		JobID        uint32
		Document     string
		Status       uint32
		StatusText   string
		PagesPrinted uint32
		TotalPages   uint32
	}
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst,SpellCheckingInspection
const (
	WCHAR_SIZE                   = 2
	PRINTER_ENUM_LOCAL           = 0x00000002
	PRINTER_ENUM_CONNECTIONS     = 0x00000004
	PRINTER_ATTRIBUTE_DEFAULT    = 0x00000004
	PORT_TYPE_TCPIP_MONITOR      = 0x00000001
	MONITOR_STANDARD_TCPIP       = "Standard TCP/IP Port"
	MONITOR_TCPMON_DLL           = "TCPMON.DLL"
	MONITOR_TCPIP_KEYWORD        = "tcpip"
	MONITOR_TCP_KEYWORD          = "tcp"
	JOB_STATUS_PAUSED            = 0x00000001
	JOB_STATUS_ERROR             = 0x00000002
	JOB_STATUS_DELETING          = 0x00000004
	JOB_STATUS_SPOOLING          = 0x00000008
	JOB_STATUS_PRINTING          = 0x00000010
	JOB_STATUS_OFFLINE           = 0x00000020
	JOB_STATUS_PAPEROUT          = 0x00000040
	JOB_STATUS_PRINTED           = 0x00000080
	JOB_STATUS_DELETED           = 0x00000100
	JOB_STATUS_BLOCKED_DEVQ      = 0x00000200
	JOB_STATUS_USER_INTERVENTION = 0x00000400
	JOB_STATUS_RESTART           = 0x00000800
	JOB_STATUS_COMPLETE          = 0x00001000
	JOB_STATUS_RETAINED          = 0x00002000
	JOB_STATUS_RENDERING_LOCALLY = 0x00004000
	JOB_WAIT_MILLISECONDS        = 1
	JOB_WAIT_TIMEOUT             = 600
	JOB_NOT_FOUND_TIMEOUT        = 10
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

//goland:noinspection GoUnusedExportedFunction
func GetPapersInfoW(printerName string) []PAPER_INFO {
	result := make([]PAPER_INFO, 0)
	if IsPrinterExist(printerName) {
		if printerNameW := CString(printerName); printerNameW != nil {
			defer FreeCString(printerNameW)
			if size := int(C.DeviceCapabilitiesW(printerNameW, nil, C.DC_PAPERS, nil, nil)); size > 0 {
				if paperIDs := make([]uint16, size); len(paperIDs) > 0 {
					C.DeviceCapabilitiesW(printerNameW, nil, C.DC_PAPERS, (*C.wchar_t)(unsafe.Pointer(&paperIDs[0])), nil)
					if paperSizes := make([]PAPER_SIZE, size); len(paperSizes) > 0 {
						C.DeviceCapabilitiesW(printerNameW, nil, C.DC_PAPERSIZE, (*C.wchar_t)(unsafe.Pointer(&paperSizes[0])), nil)
						if paperNames := make([]uint16, size*C.NAME_SIZE); len(paperNames) > 0 {
							C.DeviceCapabilitiesW(printerNameW, nil, C.DC_PAPERNAMES, (*C.wchar_t)(unsafe.Pointer(&paperNames[0])), nil)
							for i := 0; i < size; i++ {
								namePtr := unsafe.Pointer(&paperNames[i*C.NAME_SIZE])
								name := syscall.UTF16ToString((*[C.NAME_SIZE]uint16)(namePtr)[:])
								result = append(result, PAPER_INFO{
									Id:        paperIDs[i],
									PaperName: name,
									Size:      paperSizes[i],
								})
							}
						}
					}
				}
			}
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func GetPrinterDriverFiles(printerName string) []string {
	result := make([]string, 0)
	if printerName != "" {
		var hPrinter C.HANDLE
		if printerNameW := CString(printerName); printerNameW != nil {
			defer FreeCString(printerNameW)
			if errCode := C.OpenPrinterW(printerNameW, &hPrinter, nil); errCode == C.TRUE {
				defer C.ClosePrinter(hPrinter)
				var cbNeeded C.DWORD
				if C.GetPrinterDriverW(hPrinter, nil, 3, nil, 0, &cbNeeded); cbNeeded > 0 {
					buffer := make([]byte, cbNeeded)
					pBuffer := unsafe.Pointer(&buffer[0])
					if err := C.GetPrinterDriverW(hPrinter, nil, 3, (*C.BYTE)(pBuffer), cbNeeded, &cbNeeded); err == C.TRUE {
						driverInfo := (*DRIVER_INFO_3)(pBuffer)
						for _, pDependentFiles := range []*uint16{
							driverInfo.pDataFile,
							driverInfo.pConfigFile,
							driverInfo.pHelpFile,
						} {
							result = append(result, GoString(pDependentFiles))
						}
						if driverInfo.pDependentFiles != nil {
							p := driverInfo.pDependentFiles
							for {
								if pDependentFile := GoString(p); pDependentFile == "" {
									break
								} else {
									result = append(result, pDependentFile)
									p = (*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + uintptr(len(pDependentFile)+1)*WCHAR_SIZE))
								}
							}
						}
					}
				}
			}
		}

	}
	return result
}

func IsPrinterExist(printerName string) bool {
	result := false
	if printerNameW := CString(printerName); printerNameW != nil {
		defer FreeCString(printerNameW)
		var hPrinter C.HANDLE = nil
		if ret := C.OpenPrinterW(printerNameW, &hPrinter, nil); ret != 0 {
			result = true
			C.ClosePrinter(hPrinter)
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetPrinters() []PRINTER_INFO {
	result := make([]PRINTER_INFO, 0)
	var cbNeeded C.DWORD
	var cReturned C.DWORD
	flags := C.DWORD(PRINTER_ENUM_LOCAL | PRINTER_ENUM_CONNECTIONS)
	if C.EnumPrintersW(flags, nil, 2, nil, 0, &cbNeeded, &cReturned); cbNeeded > 0 {
		buffer := make([]byte, cbNeeded)
		pBuffer := unsafe.Pointer(&buffer[0])
		if C.EnumPrintersW(flags, nil, 2, (*C.BYTE)(pBuffer), cbNeeded, &cbNeeded, &cReturned) != C.FALSE {
			printerInfoArray := (*[1 << 20]PRINTER_INFO_2)(pBuffer)[:cReturned:cReturned]
			for _, info := range printerInfoArray {
				name := GoString(info.pPrinterName)
				portName := GoString(info.pPortName)
				if name != "" {
					isDefault := (info.Attributes & PRINTER_ATTRIBUTE_DEFAULT) != 0
					result = append(result, PRINTER_INFO{
						Name:      name,
						PortName:  portName,
						IsDefault: isDefault,
					})
				}
			}
		}
	}
	return result
}

func GetPorts() map[string]PORT_INFO_2 {
	result := make(map[string]PORT_INFO_2)
	var cbNeeded C.DWORD
	var cReturned C.DWORD
	if C.EnumPortsW(nil, 2, nil, 0, &cbNeeded, &cReturned); cbNeeded > 0 {
		buffer := make([]byte, cbNeeded)
		pBuffer := unsafe.Pointer(&buffer[0])
		if C.EnumPortsW(nil, 2, (*C.BYTE)(pBuffer), cbNeeded, &cbNeeded, &cReturned) != C.FALSE {
			portInfoArray := (*[1 << 20]PORT_INFO_2)(pBuffer)[:cReturned:cReturned]
			for _, info := range portInfoArray {
				portName := GoString(info.pPortName)
				if portName != "" {
					result[portName] = info
				}
			}
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetTcpIpPrinters() []PRINTER_INFO {
	result := make([]PRINTER_INFO, 0)
	ports := GetPorts()
	allPrinters := GetPrinters()
	for _, printer := range allPrinters {
		if port, ok := ports[printer.PortName]; ok {
			monitorName := GoString(port.pMonitorName)
			if strings.EqualFold(monitorName, MONITOR_STANDARD_TCPIP) ||
				strings.EqualFold(monitorName, MONITOR_TCPMON_DLL) ||
				strings.Contains(strings.ToLower(monitorName), MONITOR_TCPIP_KEYWORD) ||
				strings.Contains(strings.ToLower(monitorName), MONITOR_TCP_KEYWORD) {
				result = append(result, printer)
			}
		}
	}
	return result
}

func GetJobs(printerName string) []JOB_INFO {
	result := make([]JOB_INFO, 0)
	if printerNameW := CString(printerName); printerNameW != nil {
		defer FreeCString(printerNameW)
		var hPrinter C.HANDLE
		if C.OpenPrinterW(printerNameW, &hPrinter, nil) != 0 {
			defer C.ClosePrinter(hPrinter)
			var bytesNeeded C.DWORD
			var jobsReturned C.DWORD
			C.EnumJobsW(hPrinter, 0, 0xFFFF, 1, nil, 0, &bytesNeeded, &jobsReturned)
			if bytesNeeded > 0 {
				buffer := make([]byte, bytesNeeded)
				pBuffer := unsafe.Pointer(&buffer[0])
				if C.EnumJobsW(hPrinter, 0, 0xFFFF, 1, (*C.BYTE)(pBuffer), bytesNeeded, &bytesNeeded, &jobsReturned) != 0 {
					jobInfoArray := (*[1 << 20]C.JOB_INFO_1W)(pBuffer)[:jobsReturned:jobsReturned]
					for _, info := range jobInfoArray {
						result = append(result, JOB_INFO{
							JobID:        uint32(info.JobId),
							Document:     GoString(info.pDocument),
							Status:       uint32(info.Status),
							StatusText:   GoString(info.pStatus),
							PagesPrinted: uint32(info.PagesPrinted),
							TotalPages:   uint32(info.TotalPages),
						})
					}
				}
			}
		}
	}
	return result
}

func GetJob(printerName string, jobName ...interface{}) []JOB_INFO {
	result := make([]JOB_INFO, 0)
	if len(jobName) > 0 {
		jobs := GetJobs(printerName)
		if len(jobs) > 0 {
			switch arg := jobName[0].(type) {
			case string:
				for _, j := range jobs {
					if j.Document == arg {
						result = append(result, j)
					}
				}
			case uint32:
				for _, j := range jobs {
					if j.JobID == arg {
						result = append(result, j)
					}
				}
			}
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func IsPrinterIdle(printerName string) bool {
	result := true
	jobs := GetJobs(printerName)
	for _, job := range jobs {
		if (job.Status&JOB_STATUS_PRINTING) != 0 ||
			(job.Status&JOB_STATUS_SPOOLING) != 0 {
			result = false
			break
		}
	}
	return result
}

func IsJobCompleted(jobStatus uint32) bool {
	result := false
	if (jobStatus&JOB_STATUS_PRINTING) == 0 &&
		(jobStatus&JOB_STATUS_SPOOLING) == 0 {
		unfinishedMask := uint32(JOB_STATUS_PAUSED | JOB_STATUS_ERROR | JOB_STATUS_DELETING |
			JOB_STATUS_OFFLINE | JOB_STATUS_PAPEROUT | JOB_STATUS_BLOCKED_DEVQ |
			JOB_STATUS_USER_INTERVENTION | JOB_STATUS_RESTART)
		if (jobStatus & unfinishedMask) == 0 {
			completedMask := uint32(JOB_STATUS_PRINTED | JOB_STATUS_COMPLETE | JOB_STATUS_RETAINED)
			result = (jobStatus & completedMask) != 0
		}
	}
	return result
}

func WaitJobCompleted(printerName string, jobName ...interface{}) error {
	var result error
	var jobInfo *JOB_INFO
	startTime := time.Now()
	found := false
	for !found {
		if jobInfos := GetJob(printerName, jobName...); len(jobInfos) > 0 {
			jobInfo = &jobInfos[0]
			found = true
		} else {
			if time.Since(startTime) > time.Duration(JOB_NOT_FOUND_TIMEOUT)*time.Second {
				result = fmt.Errorf("job not found within %d seconds", JOB_NOT_FOUND_TIMEOUT)
				found = true
			} else {
				time.Sleep(JOB_WAIT_MILLISECONDS * time.Millisecond)
			}
		}
	}
	if result == nil && jobInfo != nil {
		completed := false
		for !completed {
			if IsJobCompleted(jobInfo.Status) {
				completed = true
			} else {
				if time.Since(startTime) > time.Duration(JOB_WAIT_TIMEOUT)*time.Second {
					result = fmt.Errorf("job wait timeout after %d seconds", JOB_WAIT_TIMEOUT)
					completed = true
				} else {
					time.Sleep(time.Second)
					if jobInfos := GetJob(printerName, jobName...); len(jobInfos) > 0 {
						jobInfo = &jobInfos[0]
					} else {
						completed = true
					}
				}
			}
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func SetDefaultPrinter(printerName string) error {
	var result error
	if !IsPrinterExist(printerName) {
		result = fmt.Errorf("printer not found: %s", printerName)
	} else {
		if printerNameW := CString(printerName); printerNameW != nil {
			defer FreeCString(printerNameW)
			if C.SetDefaultPrinterW(printerNameW) == 0 {
				result = fmt.Errorf("failed to set default printer: %s", printerName)
			}
		}
	}
	return result
}
