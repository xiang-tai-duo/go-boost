//go:build darwin

package osutil

import "os/exec"

func SetHideWindow(cmd *exec.Cmd) {
	// No-op on macOS
}
