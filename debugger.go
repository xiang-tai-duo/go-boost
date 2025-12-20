// Package boost
// File:        debugger.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: DEBUGGER is a utility for detecting debuggers attached to the process.
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

func init() {
	NewDebugger().Check()
}

func (debugger DEBUGGER) Check() bool {
	var isPresent bool
	p, err := ps.FindProcess(os.Getpid())
	if p == nil || err != nil {
		debuggerState = DEBUGGER_DETACHED
		isPresent = false
		return isPresent
	}

	executable := strings.ToLower(p.Executable())
	pids := make(map[int]ps.Process)
	pid := os.Getppid()
	isRecursion := false

	for pid != 0 && !isRecursion {
		var proc ps.Process
		proc, err = ps.FindProcess(pid)

		switch {
		case proc == nil || err != nil:
			isPresent = false
			pid = 0
		case strings.ToLower(proc.Executable()) == executable:
			isPresent = false
			pid = 0
		case strings.ToLower(proc.Executable()) == "dlv" || strings.ToLower(proc.Executable()) == "goland64.exe":
			isPresent = true
			pid = 0
		default:
			pid = proc.PPid()
			if _, exists := pids[pid]; exists {
				isRecursion = true
				break
			}
			pids[proc.Pid()] = proc
		}
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
