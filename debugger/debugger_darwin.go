//go:build darwin

// Package debugger
// File:        debugger_darwin.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/debugger/debugger_darwin.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: DEBUGGER native detection fallback for Darwin
// --------------------------------------------------------------------------------

package debugger

func isDebuggerPresent() bool {
	return false
}
