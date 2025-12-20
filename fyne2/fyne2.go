// Package fyne2
// File:        fyne.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/fyne2/fyne.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Fyne utility functions for GUI applications, providing rendering configuration helpers.
// --------------------------------------------------------------------------------
package fyne2

import (
	"os"

	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection SpellCheckingInspection,GoSnakeCaseUsage,GoUnusedConst
const (
	GALLIUM_DRIVER    = "GALLIUM_DRIVER"
	LLVMPIPE_DRIVER   = "llvmpipe"
	MODULE_NAME_FYNE2 = "fyne2"
)

func init() {
	os.Setenv(GALLIUM_DRIVER, LLVMPIPE_DRIVER)
}

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_FYNE2, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_FYNE2, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_FYNE2, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_FYNE2, logger.SKIP_STACK_FRAMES_BASE)
}
