//go:build !electron_dist

package electron

import "embed"

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
var (
	ELECTRON_DIST_FILES embed.FS
)
