// Package mesa3d
// File:        mesa3d_distribution.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/mesa3d/mesa3d_distribution.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Mesa3D distribution file holder (empty; Mesa3D DLLs are deployed externally).
// --------------------------------------------------------------------------------
package mesa3d

import "embed"

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	MODULE_NAME_MESA3D_DISTRIBUTION = "mesa3d_distribution"
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
var (
	MESA3D_DISTRIBUTION_FILES_X64 embed.FS
	MESA3D_DISTRIBUTION_FILES_X86 embed.FS
)
