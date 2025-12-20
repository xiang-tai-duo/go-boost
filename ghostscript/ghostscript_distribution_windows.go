//go:build distribution_all_in_one && windows

// Package ghostscript
// File:        ghostscript_distribution_windows.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/ghostscript/ghostscript_distribution_windows.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Embedded Ghostscript distribution files for Windows.
// --------------------------------------------------------------------------------
package ghostscript

import __embed "embed"

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
var (
	//go:embed ghostpcl-10.07.1-win32/*
	GHOSTSCRIPT_DISTRIBUTION_FILES __embed.FS
)
