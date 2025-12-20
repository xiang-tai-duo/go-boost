//go:build !distribution_all_in_one

// Package electron
// File:        electron_distribution.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/electron/electron_distribution.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Electron distribution file holder.
// --------------------------------------------------------------------------------
package electron

import "embed"

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	MODULE_NAME_ELECTRON_DISTRIBUTION = "electron_distribution"
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
var (
	ELECTRON_DISTRIBUTION_FILES embed.FS
)
