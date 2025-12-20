//go:build darwin

// File:        mesa3d_darwin.go
// Description: Darwin stub for mesa3d package (no-op)
// --------------------------------------------------------------------------------

package mesa3d

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	MESA3D_UNPACKED_DIRECTORY_PERMISSION = 0755
	GOARCH_AMD64                         = "amd64"
	MODULE_NAME_MESA3D                   = "mesa3d"
)

func ExtractEnvironment() (int, error) {
	return 0, nil
}
