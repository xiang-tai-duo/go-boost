//go:build linux && distribution_all_in_one

// Package electron
// File:        electron_distribution_linux.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/electron/electron_distribution_linux.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Embedded Electron distribution files for Linux.
// --------------------------------------------------------------------------------
package electron

import "embed"

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	MODULE_NAME_ELECTRON_DISTRIBUTION_LINUX = "electron_distribution_linux"
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
var (
	//go:embed dist/linux-unpacked
	ELECTRON_DISTRIBUTION_FILES embed.FS
)
