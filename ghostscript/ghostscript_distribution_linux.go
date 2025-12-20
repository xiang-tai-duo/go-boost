//go:build distribution_all_in_one && linux

// Package ghostscript
// File:        ghostscript_distribution_linux.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/ghostscript/ghostscript_distribution_linux.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Embedded Ghostscript distribution files for Linux.
// --------------------------------------------------------------------------------
package ghostscript

import __embed "embed"

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
var (
	//go:embed ghostpdl-10.07.1/bin/*
	GHOSTSCRIPT_DISTRIBUTION_FILES __embed.FS
)
