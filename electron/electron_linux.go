//go:build linux

// Package electron
// File:        electron_linux.go
// Url:         https://github.com/xiang-tai-duo/go-bootstrap/blob/master/electron/electron_linux.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Electron Linux specific implementation
// --------------------------------------------------------------------------------
package electron

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
const (
	ELECTRON_NAME       = "go-boost-electron"
	ELECTRON_DIST_PATH  = "dist/linux-unpacked"
	ELECTRON_NO_SANDBOX = "--no-sandbox"
)
