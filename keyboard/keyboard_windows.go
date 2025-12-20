//go:build windows
// +build windows

package keyboard

/*
#include <windows.h>
*/
import "C"

import (
	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection GoSnakeCaseUsage,SpellCheckingInspection,GoUnusedConst
const (
	VK_LCONTROL                  = 0xA2
	KEY_PRESSED_MASK             = C.SHORT(-32768)
	MODULE_NAME_KEYBOARD_WINDOWS = "keyboard_windows"
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_KEYBOARD_WINDOWS, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_KEYBOARD_WINDOWS, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_KEYBOARD_WINDOWS, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_KEYBOARD_WINDOWS, logger.SKIP_STACK_FRAMES_BASE)
}

func isKeyPressed(vk int) bool {
	result := false
	ret := C.GetAsyncKeyState(C.int(vk))
	if (ret & KEY_PRESSED_MASK) != 0 {
		result = true
	}
	return result
}
