package main

import (
	"fmt"
	"os"
	"syscall"
)

const (
	ERROR_MESSAGE_SIGNAL = "signal error: %s"
)

func isProcessAlive(pid int) (bool, error) {
	result := false
	err := error(nil)
	process := (*os.Process)(nil)
	if process, err = os.FindProcess(pid); err == nil {
		err = process.Signal(syscall.Signal(0))
		if err == nil {
			result = true
		} else {
			if err == syscall.ESRCH {
				result = false
				err = nil
			} else if err == syscall.EPERM {
				result = true
				err = nil
			} else {
				err = fmt.Errorf(ERROR_MESSAGE_SIGNAL, err.Error())
			}
		}
	} else {
		err = nil
	}
	return result, err
}
