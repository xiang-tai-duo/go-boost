//go:build windows

// Package osutil - Windows specific exec utilities
package osutil

import (
	"os/exec"
	"syscall"
)

func SetHideWindow(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
