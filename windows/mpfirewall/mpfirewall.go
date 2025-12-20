//go:build windows

// Package mpfirewall
// File:        mpfirewall.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/windows/mpfirewall/mpfirewall.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Windows Firewall COM API wrapper using INetFwPolicy2 interface
// --------------------------------------------------------------------------------

package mpfirewall

/*
#cgo CFLAGS: -DWIN32_LEAN_AND_MEAN
#cgo LDFLAGS: -lole32 -loleaut32 -ladvapi32
#include "mpfirewall.h"
*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/xiang-tai-duo/go-boost/logger"
	"github.com/xiang-tai-duo/go-boost/strings2"
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	ACTION_ALLOW  = int(C.MPFIREWALL_ACTION_ALLOW)
	ACTION_BLOCK  = int(C.MPFIREWALL_ACTION_BLOCK)
	PROTOCOL_TCP  = int(C.MPFIREWALL_PROTOCOL_TCP)
	PROTOCOL_UDP  = int(C.MPFIREWALL_PROTOCOL_UDP)
	PROTOCOL_ANY  = int(C.MPFIREWALL_PROTOCOL_ANY)
	DIRECTION_IN  = int(C.MPFIREWALL_RULE_DIR_INBOUND)
	DIRECTION_OUT = int(C.MPFIREWALL_RULE_DIR_OUTBOUND)
	PROFILE_ALL   = int(C.MPFIREWALL_PROFILE_ALL)

	ERROR_MPFIREWALL_INIT_FAILED   = "failed to initialize firewall COM object"
	ERROR_MPFIREWALL_ADD_FAILED    = "failed to add firewall rule"
	ERROR_MPFIREWALL_DELETE_FAILED = "failed to delete firewall rule"
	ERROR_MPFIREWALL_QUERY_FAILED  = "failed to query firewall rules"

	MODULE_NAME_MPFIREWALL = "windows.mpfirewall.mpfirewall"
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_MPFIREWALL, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_MPFIREWALL, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_MPFIREWALL, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedExportedFunction
func AddFirewallRule(ruleName string, exeFilePath string, port int) error {
	return AddFirewallRuleEx(ruleName, exeFilePath, port, PROTOCOL_TCP)
}

//goland:noinspection GoUnusedExportedFunction
func AddFirewallRuleEx(ruleName string, exeFilePath string, port int, protocol int) error {
	result := error(nil)
	err := error(nil)
	holder, initErr := newPolicy2Holder()
	if initErr != nil {
		err = initErr
	} else {
		defer holder.release()
		portStr := C.CString(strings2.Itoa(port))
		defer C.free(unsafe.Pointer(portStr))
		nameStr := C.CString(ruleName)
		defer C.free(unsafe.Pointer(nameStr))
		descStr := C.CString(ruleName)
		defer C.free(unsafe.Pointer(descStr))
		appStr := C.CString(exeFilePath)
		defer C.free(unsafe.Pointer(appStr))
		hr := C.mpfirewall_add_rule(
			holder.ptr,
			nameStr,
			descStr,
			appStr,
			portStr,
			C.int(DIRECTION_IN),
			C.int(ACTION_ALLOW),
			C.int(protocol),
			C.int(PROFILE_ALL),
			1,
		)
		if hr != C.MPFIREWALL_S_OK {
			err = fmt.Errorf("%s: HRESULT=0x%08X", ERROR_MPFIREWALL_ADD_FAILED, uint32(hr))
		}
	}
	result = err
	return result
}

//goland:noinspection GoUnusedExportedFunction
func DeleteFirewallRule(ruleName string) error {
	result := error(nil)
	err := error(nil)
	holder, initErr := newPolicy2Holder()
	if initErr != nil {
		err = initErr
	} else {
		defer holder.release()
		nameStr := C.CString(ruleName)
		defer C.free(unsafe.Pointer(nameStr))
		hr := C.mpfirewall_delete_rule(holder.ptr, nameStr)
		if hr != C.MPFIREWALL_S_OK {
			err = fmt.Errorf("%s: HRESULT=0x%08X", ERROR_MPFIREWALL_DELETE_FAILED, uint32(hr))
		}
	}
	result = err
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetFirewallBlockedRules(exeFilePath string, port int) map[string]bool {
	result := make(map[string]bool)
	holder, initErr := newPolicy2Holder()
	if initErr != nil {
	} else {
		defer holder.release()
		exeStr := C.CString(exeFilePath)
		defer C.free(unsafe.Pointer(exeStr))
		listPtr := C.mpfirewall_get_blocked_rules(holder.ptr, exeStr, C.int(port))
		if listPtr != nil {
			if listPtr.count > 0 && listPtr.rule_names != nil {
				names := unsafe.Slice((**C.char)(unsafe.Pointer(listPtr.rule_names)), int(listPtr.count))
				for i := 0; i < int(listPtr.count); i++ {
					name := C.GoString(names[i])
					result[name] = true
				}
			}
			C.mpfirewall_free_rule_name_list(listPtr)
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetFirewallRuleLocalPort(ruleName string) (string, error) {
	result := ""
	err := error(nil)
	holder, initErr := newPolicy2Holder()
	if initErr != nil {
		err = initErr
	} else {
		defer holder.release()
		nameStr := C.CString(ruleName)
		defer C.free(unsafe.Pointer(nameStr))
		var localPortStr *C.char
		hr := C.mpfirewall_get_rule_local_port(holder.ptr, nameStr, &localPortStr)
		if hr == C.MPFIREWALL_S_OK && localPortStr != nil {
			result = C.GoString(localPortStr)
			C.free(unsafe.Pointer(localPortStr))
		} else if hr != C.MPFIREWALL_S_OK {
			err = fmt.Errorf("%s: HRESULT=0x%08X", ERROR_MPFIREWALL_QUERY_FAILED, uint32(hr))
		}
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func IsFirewallRuleExists(exeFilePath string, port int) (bool, error) {
	result := false
	err := error(nil)
	holder, initErr := newPolicy2Holder()
	if initErr != nil {
		err = initErr
	} else {
		defer holder.release()
		exeStr := C.CString(exeFilePath)
		defer C.free(unsafe.Pointer(exeStr))
		infoPtr := C.mpfirewall_is_rule_exists(holder.ptr, exeStr, C.int(port))
		if infoPtr != nil {
			if infoPtr.found != 0 {
				result = true
			}
			C.mpfirewall_free_rule_info(infoPtr)
		}
	}
	return result, err
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_MPFIREWALL, logger.SKIP_STACK_FRAMES_BASE)
}

type policy2Holder struct {
	ptr unsafe.Pointer
}

func newPolicy2Holder() (*policy2Holder, error) {
	result := &policy2Holder{}
	err := error(nil)
	var ptr unsafe.Pointer
	hr := C.mpfirewall_init(&ptr)
	if hr != C.MPFIREWALL_S_OK {
		err = fmt.Errorf("%s: HRESULT=0x%08X", ERROR_MPFIREWALL_INIT_FAILED, uint32(hr))
	} else {
		result.ptr = ptr
	}
	return result, err
}

func (h *policy2Holder) release() {
	if h.ptr != nil {
		C.mpfirewall_release(h.ptr)
		h.ptr = nil
	}
}
