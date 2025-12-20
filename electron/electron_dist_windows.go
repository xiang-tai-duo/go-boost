//go:build windows && electron_dist

package electron

import "embed"

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
var (
	//go:embed dist/win-unpacked
	ELECTRON_DIST_FILES embed.FS
)
