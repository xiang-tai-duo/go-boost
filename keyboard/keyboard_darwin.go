//go:build darwin
// +build darwin

package keyboard

import (
	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	MODULE_NAME_KEYBOARD_DARWIN = "keyboard_darwin"
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_KEYBOARD_DARWIN, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_KEYBOARD_DARWIN, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_KEYBOARD_DARWIN, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_KEYBOARD_DARWIN, logger.SKIP_STACK_FRAMES_BASE)
}

func isKeyPressed(vk int) bool {
	return false
}
