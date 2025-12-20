//go:build windows

// Package electron
// File:        electron_windows.go
// Url:         https://github.com/xiang-tai-duo/go-bootstrap/blob/master/electron/electron_windows.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Electron Windows specific implementation
// --------------------------------------------------------------------------------
package electron

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName,GoUnusedConst
const (
	ELECTRON_NAME       = "go-boost-electron.exe"
	ELECTRON_DIST_PATH  = "dist/win-unpacked"
	ELECTRON_NO_SANDBOX = ""
)
