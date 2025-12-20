// Package debugger
// File:        debugger.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/debugger/debugger.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: DEBUGGER is a utility for detecting debuggers attached to the process
// --------------------------------------------------------------------------------
package debugger

import (
	"os"
	"strings"

	"github.com/mitchellh/go-ps"
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName,SpellCheckingInspection
const (
	DEBUGGER_ATTACHED                     = 1
	DEBUGGER_DELVE                        = "dlv"
	DEBUGGER_DELVE_DEBUG_ADAPTER_PROTOCOL = "dlv-dap"
	DEBUGGER_DETACHED                     = 2
	DEBUGGER_GOGLAND                      = "gogland"
	DEBUGGER_GOGLAND_64                   = "gogland64"
	DEBUGGER_GOLAND                       = "goland"
	DEBUGGER_GOLAND_64                    = "goland64"
	DEBUGGER_NONE                         = 0
	DEBUG_BINARY_PREFIX                   = "__debug_bin"
	PROCESS_ID_NONE                       = 0
	WINDOWS_EXECUTABLE_SUFFIX             = ".exe"
)

var (
	debuggerState                = DEBUGGER_NONE
	isGoland64Attached           = false
	isThirdPartyDebuggerAttached = false
)

func init() {
	InitializeDebuggerStatus()
}

func InitializeDebuggerStatus() bool {
	result := false
	err := error(nil)
	var process ps.Process
	if process, err = ps.FindProcess(os.Getpid()); err == nil && process != nil {
		executable := strings.ToLower(process.Executable())
		processIds := make(map[int]ps.Process)
		parentProcessId := os.Getppid()
		isRecursion := false
		for parentProcessId != PROCESS_ID_NONE && !isRecursion {
			var parentProcess ps.Process
			if parentProcess, err = ps.FindProcess(parentProcessId); err == nil && parentProcess != nil {
				parentExecutable := strings.ToLower(parentProcess.Executable())
				if parentExecutable == executable {
					result = false
					parentProcessId = PROCESS_ID_NONE
				} else if isGoland64(parentExecutable) {
					result = true
					parentProcessId = PROCESS_ID_NONE
				} else {
					parentProcessId = parentProcess.PPid()
					if _, exists := processIds[parentProcessId]; exists {
						isRecursion = true
					} else {
						processIds[parentProcess.Pid()] = parentProcess
					}
				}
			} else {
				result = false
				parentProcessId = PROCESS_ID_NONE
			}
		}
	}
	isGoland64Attached = result
	if isDebuggerPresent() {
		isThirdPartyDebuggerAttached = true
		result = true
	}
	if result {
		debuggerState = DEBUGGER_ATTACHED
	} else {
		debuggerState = DEBUGGER_DETACHED
	}
	return result
}

func IsGoland64Attached() bool {
	return isGoland64Attached
}

func IsPresent() bool {
	result := false
	if debuggerState == DEBUGGER_ATTACHED {
		result = true
	}
	return result
}

func IsThirdPartyDebuggerAttached() bool {
	return isThirdPartyDebuggerAttached
}

func isGoland64(name string) bool {
	result := false
	if name != "" {
		lowerName := strings.TrimSuffix(strings.ToLower(name), WINDOWS_EXECUTABLE_SUFFIX)
		switch lowerName {
		case DEBUGGER_DELVE, DEBUGGER_DELVE_DEBUG_ADAPTER_PROTOCOL, DEBUGGER_GOLAND_64, DEBUGGER_GOLAND, DEBUGGER_GOGLAND_64, DEBUGGER_GOGLAND:
			result = true
		default:
			if strings.HasPrefix(lowerName, DEBUG_BINARY_PREFIX) || strings.Contains(lowerName, DEBUGGER_DELVE) {
				result = true
			}
		}
	}
	return result
}
