//go:build windows

// Package winspool
// File:        winspool.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/windows/winspool/winspool.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: winspool is a wrapper for Windows print spooler operations, providing methods for printer management and paper information retrieval.
// --------------------------------------------------------------------------------
package winspool

/*
#include <windows.h>
#include <winspool.h>
#cgo LDFLAGS: -lwinspool

#define NAME_SIZE 64
*/
import "C"

import (
	"syscall"
	"unsafe"

	"github.com/xiang-tai-duo/go-boost/pinvoke"
	"github.com/xiang-tai-duo/go-boost/windows/winbase"
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
)

//goland:noinspection SpellCheckingInspection,GoUnusedFunction,GoUnusedExportedFunction
func GoStringW(lpwsz *C.wchar_t) string {
	sz := ""
	if pwsz := (*uint16)(unsafe.Pointer(lpwsz)); pwsz != nil {
		for ptr := pwsz; *ptr != 0; ptr = (*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + pinvoke.WCHAR_SIZE)) {
			sz += string(rune(*ptr))
		}
	}
	return sz
}

//goland:noinspection GoSnakeCaseUsage,GoUnusedFunction,GoUnusedExportedFunction
func CStringW(s string) *C.wchar_t {
	var wsz *C.wchar_t = nil
	if ptr, err := syscall.UTF16FromString(s); err == nil {
		wsz = (*C.wchar_t)(unsafe.Pointer(&ptr[0]))
	}
	return wsz
}

//goland:noinspection GoUnusedExportedFunction
func GetPapersInfoW(printerName string) []PAPER_INFO {
	result := make([]PAPER_INFO, 0)
	if IsPrinterExist(printerName) {
		if printerNameW := CStringW(printerName); printerNameW != nil {
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
		if errCode := C.OpenPrinterW(CStringW(printerName), &hPrinter, nil); errCode == C.TRUE {
			defer C.ClosePrinter(hPrinter)
			var cbNeeded C.DWORD
			if C.GetPrinterDriverW(hPrinter, nil, 3, nil, 0, &cbNeeded); cbNeeded > 0 {
				buffer := make([]byte, cbNeeded)
				pBuffer := unsafe.Pointer(&buffer[0])
				if errCode := C.GetPrinterDriverW(hPrinter, nil, 3, (*C.BYTE)(pBuffer), cbNeeded, &cbNeeded); errCode == C.TRUE {
					driverInfo := (*DRIVER_INFO_3)(pBuffer)
					for _, pDependentFiles := range []*uint16{
						driverInfo.pDataFile,
						driverInfo.pConfigFile,
						driverInfo.pHelpFile,
					} {
						if !winbase.IsBadStringPtr(pDependentFiles) {
							result = append(result, syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(pDependentFiles))[:]))
						}
					}
					if !winbase.IsBadStringPtr(driverInfo.pDependentFiles) {
						p := driverInfo.pDependentFiles
						for {
							if pDependentFile := syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(p))[:]); pDependentFile == "" {
								break
							} else {
								result = append(result, pDependentFile)
								p = (*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + uintptr(len(pDependentFile)+1)*pinvoke.WCHAR_SIZE))
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
	if printerNameW := CStringW(printerName); printerNameW != nil {
		var hPrinter C.HANDLE = nil
		if ret := C.OpenPrinterW(printerNameW, &hPrinter, nil); ret != 0 {
			result = true
			C.ClosePrinter(hPrinter)
		}
	}
	return result
}
