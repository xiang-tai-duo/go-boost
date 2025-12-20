// Package debugger
// File:        debugger_windows.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/debugger/debugger_windows.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: DEBUGGER native detection for Windows platform
// --------------------------------------------------------------------------------

//go:build windows

package debugger

/*
#include <windows.h>
*/
import "C"

func isDebuggerPresent() bool {
	result := false
	if isLocalDebuggerPresent() || isRemoteDebuggerPresent() {
		result = true
	}
	return result
}

func isLocalDebuggerPresent() bool {
	return C.IsDebuggerPresent() != 0
}

func isRemoteDebuggerPresent() bool {
	result := false
	var present C.BOOL
	if C.CheckRemoteDebuggerPresent(C.GetCurrentProcess(), &present) == 0 {
		result = false
	} else {
		result = present != 0
	}
	return result
}
