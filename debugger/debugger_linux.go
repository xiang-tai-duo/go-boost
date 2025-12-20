//go:build linux

// Package debugger
// File:        debugger_linux.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/debugger/debugger_linux.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: DEBUGGER native detection fallback for Linux
// --------------------------------------------------------------------------------

package debugger

func isDebuggerPresent() bool {
	return false
}
