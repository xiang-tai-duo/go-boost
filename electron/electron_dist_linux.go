//go:build linux && electron_dist

package electron

import "embed"

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
var (
	//go:embed dist/linux-unpacked
	ELECTRON_DIST_FILES embed.FS
)
