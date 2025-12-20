//go:build darwin

// Package hardware
// File:        hardware_darwin.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/hardware/hardware_darwin.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: HARDWARE is a wrapper for hardware information collection, providing macOS (darwin) specific implementations via cgo/IOKit.
// --------------------------------------------------------------------------------

package hardware

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -framework CoreFoundation -framework IOKit
#include "hardware_darwin.h"
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"net"
	"strconv"
	"unsafe"

	"github.com/xiang-tai-duo/go-boost/logger"
)

const (
	DARWIN_SINGLETON_PORT_LOWER_BOUND  = 65500
	DARWIN_SINGLETON_PORT_RANGE_LENGTH = 31
	LOCAL_HOST_IPV4                    = "127.0.0.1"
)

var singletonListener net.Listener

func AcquireSingletonSocket() error {
	logger.Logger.Debug("AcquireSingletonSocket enter")
	result := error(nil)
	err := error(nil)
	seed := uint64(C.hardware_fingerprint_seed())
	port := DARWIN_SINGLETON_PORT_LOWER_BOUND + int(seed%DARWIN_SINGLETON_PORT_RANGE_LENGTH)
	address := net.JoinHostPort(LOCAL_HOST_IPV4, strconv.Itoa(port))
	logger.Logger.Debug(fmt.Sprintf("AcquireSingletonSocket address=%s seed=%d", address, seed))
	var listener net.Listener
	if listener, err = net.Listen("tcp", address); err == nil {
		singletonListener = listener
	} else {
		result = err
	}
	logger.Logger.Debug(fmt.Sprintf("AcquireSingletonSocket result err=%v", result))
	return result
}

func getBIOSInfo() string {
	logger.Logger.Debug("getBIOSInfo enter")
	result := ""
	value := C.hardware_platform_uuid()
	if value != nil {
		result = C.GoString(value)
		C.free(unsafe.Pointer(value))
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
	value := C.hardware_platform_serial()
	if value != nil {
		result = C.GoString(value)
		C.free(unsafe.Pointer(value))
	}
	logger.Logger.Debug(fmt.Sprintf("getDiskSerial resultLength=%d", len(result)))
	return result
}

func getMemoryInfo() string {
	logger.Logger.Debug("getMemoryInfo enter")
	result := fmt.Sprintf(FORMAT_NUMBER, uint64(C.hardware_physical_memory()))
	logger.Logger.Debug(fmt.Sprintf("getMemoryInfo result=%s", result))
	return result
}
