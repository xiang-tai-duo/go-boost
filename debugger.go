// --------------------------------------------------------------------------------
// File:        debugger.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: DEBUGGER is a utility for detecting debuggers attached to the process.
// --------------------------------------------------------------------------------

package boost

import (
	"os"
	"strings"

	"github.com/mitchellh/go-ps"
)

// DEBUGGER represents a utility for detecting debuggers attached to the process.
type DEBUGGER struct{}

var (
	// Current debugger state
	debuggerState = DEBUGGER_NONE
)

// NewDebugger creates a new DEBUGGER instance
// Returns: New DEBUGGER instance
// Usage:
// debugger := NewDebugger()
func NewDebugger() DEBUGGER {
	return DEBUGGER{}
}

const (
	// DEBUGGER_NONE indicates no debugger is present
	DEBUGGER_NONE = 0
	// DEBUGGER_ATTACHED indicates a debugger is attached
	DEBUGGER_ATTACHED = 1
	// DEBUGGER_DETACHED indicates a debugger was previously attached but is now detached
	DEBUGGER_DETACHED = 2
)

func init() {
	NewDebugger().Check()
}

// Check detects if a debugger is attached to the process.
// Returns: true if debugger is attached, false otherwise
// Usage:
// isDebuggerAttached := Debugger.Check()
func (debugger DEBUGGER) Check() bool {
	var isPresent bool
	p, err := ps.FindProcess(os.Getpid())
	if p == nil || err != nil {
		debuggerState = DEBUGGER_DETACHED
		isPresent = false
	} else {
		executable := strings.ToLower(p.Executable())
		pids := make(map[int]ps.Process)
		pid := os.Getppid()
		isRecursion := false
		for pid != 0 && !isRecursion {
			switch p, err := ps.FindProcess(pid); {
			case p == nil || err != nil:
				isPresent = false
				pid = 0
			case strings.ToLower(p.Executable()) == executable:
				isPresent = false
				pid = 0
			case strings.ToLower(p.Executable()) == "dlv" ||
				strings.ToLower(p.Executable()) == "goland64.exe":
				isPresent = true
				pid = 0
			default:
				pid = p.PPid()
				if _, exists := pids[pid]; exists {
					isRecursion = true
					break
				}
				pids[p.Pid()] = p
			}
		}
		if isPresent {
			debuggerState = DEBUGGER_ATTACHED
		} else {
			debuggerState = DEBUGGER_DETACHED
		}
	}
	return isPresent
}

// IsPresent returns whether a debugger is currently attached.
// Returns: true if debugger is attached, false otherwise
// Usage:
//
//	if Debugger.IsPresent() {
//	    fmt.Println("DEBUGGER detected!")
//	}
func (debugger DEBUGGER) IsPresent() bool {
	var result bool
	if debuggerState == DEBUGGER_ATTACHED {
		result = true
	} else {
		result = false
	}
	return result
}
