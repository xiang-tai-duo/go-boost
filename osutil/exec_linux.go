//go:build linux

package osutil

import "os/exec"

func SetHideWindow(cmd *exec.Cmd) {
	// No-op on Linux
}
