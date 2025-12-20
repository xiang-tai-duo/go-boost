// Package debugger
// File:        debugger.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/debugger/debugger.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: DEBUGGER is a utility for detecting debuggers attached to the process
// --------------------------------------------------------------------------------
package debugger

import (
	"os"
	"strings"

	"github.com/mitchellh/go-ps"
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
const (
	DEBUGGER_NONE     = 0
	DEBUGGER_ATTACHED = 1
	DEBUGGER_DETACHED = 2
)

var (
	debuggerState = DEBUGGER_NONE
)

func init() {
	InitializeDebuggerStatus()
}

func InitializeDebuggerStatus() bool {
	result := false
	err := error(nil)
	var p ps.Process
	if p, err = ps.FindProcess(os.Getpid()); err == nil && p != nil {
		executable := strings.ToLower(p.Executable())
		pids := make(map[int]ps.Process)
		pid := os.Getppid()
		isRecursion := false
		for pid != 0 && !isRecursion {
			var proc ps.Process
			if proc, err = ps.FindProcess(pid); err == nil && proc != nil {
				procExecutable := strings.ToLower(proc.Executable())
				if procExecutable == executable {
					result = false
					pid = 0
				} else if procExecutable == "dlv" || procExecutable == "goland64.exe" {
					result = true
					pid = 0
				} else {
					pid = proc.PPid()
					if _, exists := pids[pid]; exists {
						isRecursion = true
					} else {
						pids[proc.Pid()] = proc
					}
				}
			} else {
				result = false
				pid = 0
			}
		}
	} else {
		result = false
	}

	if result {
		debuggerState = DEBUGGER_ATTACHED
	} else {
		debuggerState = DEBUGGER_DETACHED
	}
	return result
}

func IsPresent() bool {
	result := false
	if debuggerState == DEBUGGER_ATTACHED {
		result = true
	}
	return result
}
