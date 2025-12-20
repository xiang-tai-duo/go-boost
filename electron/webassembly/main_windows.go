package main

/*
#include <windows.h>
static HANDLE openProcess(DWORD pid) { return OpenProcess(SYNCHRONIZE, FALSE, pid); }
static DWORD waitForSingleObject(HANDLE hHandle, DWORD dwMilliseconds) { return WaitForSingleObject(hHandle, dwMilliseconds); }
static BOOL closeHandle(HANDLE hObject) { return CloseHandle(hObject); }
*/
import "C"

import (
	"fmt"
)

const (
	ERROR_MESSAGE_PROCESS_NOT_FOUND = "process not found or access denied"
	ERROR_MESSAGE_WAIT_FAILED       = "wait failed with code: %d"
	WAIT_100_MILLISECONDS           = 100
	WAIT_TIMEOUT                    = 0x00000102
	WAIT_OBJECT_0                   = 0x00000000
	WAIT_FAILED                     = 0xFFFFFFFF
)

func isProcessAlive(pid int) (bool, error) {
	result := false
	err := error(nil)
	if handle := C.openProcess(C.DWORD(pid)); handle != nil {
		waitResult := C.waitForSingleObject(handle, C.DWORD(WAIT_100_MILLISECONDS))
		if waitResult == WAIT_TIMEOUT {
			result = true
		} else if waitResult == WAIT_OBJECT_0 {
			result = false
		} else {
			err = fmt.Errorf(ERROR_MESSAGE_WAIT_FAILED, waitResult)
		}
		C.closeHandle(handle)
	} else {
		err = fmt.Errorf(ERROR_MESSAGE_PROCESS_NOT_FOUND)
	}
	return result, err
}
