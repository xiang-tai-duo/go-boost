//go:build !distribution_all_in_one

// Package ghostscript
// File:        ghostscript_distribution.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/ghostscript/ghostscript_distribution.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Ghostscript distribution file holder.
// --------------------------------------------------------------------------------
package ghostscript

import __embed "embed"

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
var (
	GHOSTSCRIPT_DISTRIBUTION_FILES __embed.FS
)
