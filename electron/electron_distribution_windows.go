//go:build windows && distribution_all_in_one

// Package electron
// File:        electron_distribution_windows.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/electron/electron_distribution_windows.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Embedded Electron distribution files for Windows.
// --------------------------------------------------------------------------------
package electron

import "embed"

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	MODULE_NAME_ELECTRON_DISTRIBUTION_WINDOWS = "electron_distribution_windows"
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
var (
	//go:embed dist/win-unpacked
	ELECTRON_DISTRIBUTION_FILES embed.FS
)
