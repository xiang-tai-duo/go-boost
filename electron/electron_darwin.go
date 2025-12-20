//go:build darwin

// Package electron
// File:        electron_darwin.go
// Description: Electron macOS specific implementation
// --------------------------------------------------------------------------------
package electron

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
const (
	ELECTRON_NAME               = "go-boost-electron"
	ELECTRON_DIST_PATH          = "dist/mac-unpacked"
	ELECTRON_DISTRIBUTION_PATH  = ELECTRON_DIST_PATH
	ELECTRON_NO_SANDBOX         = "--no-sandbox"
	MODULE_NAME_ELECTRON_DARWIN = "electron_darwin"
)
