// Package boost
// File:        debugger.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/debugger.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: DEBUGGER is a utility for detecting debuggers attached to the process
// --------------------------------------------------------------------------------
package boost

import (
	"os"
	"strings"

	"github.com/mitchellh/go-ps"
)

const (
	DEBUGGER_NONE     = 0
	DEBUGGER_ATTACHED = 1
	DEBUGGER_DETACHED = 2
)

type (
	DEBUGGER struct{}
)

var (
	debuggerState = DEBUGGER_NONE
)

func NewDebugger() DEBUGGER {
	return DEBUGGER{}
}

func (debugger DEBUGGER) Check() bool {
	var err error
	var isPresent bool
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
					isPresent = false
					pid = 0
				} else if procExecutable == "dlv" || procExecutable == "goland64.exe" {
					isPresent = true
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
				isPresent = false
				pid = 0
			}
		}
	} else {
		isPresent = false
	}

	if isPresent {
		debuggerState = DEBUGGER_ATTACHED
	} else {
		debuggerState = DEBUGGER_DETACHED
	}
	return isPresent
}

func (debugger DEBUGGER) IsPresent() bool {
	var result bool
	if debuggerState == DEBUGGER_ATTACHED {
		result = true
	} else {
		result = false
	}
	return result
}

func init() {
	NewDebugger().Check()
}
