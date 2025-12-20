// Package messagebox
// File:        messagebox.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/messagebox/messagebox.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Cross-platform message box wrapper
// --------------------------------------------------------------------------------
package messagebox

import (
	"errors"

	"github.com/ncruces/zenity"
	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	DEFAULT_TITLE          = "Message"
	MODULE_NAME_MESSAGEBOX = "messagebox"
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_MESSAGEBOX, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_MESSAGEBOX, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_MESSAGEBOX, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedExportedFunction
func Confirm(title string, message string) (bool, error) {
	result, err := Question(title, message)
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func Error(title string, message string) error {
	result := zenity.Error(message, zenity.Title(getTitle(title)))
	return result
}

//goland:noinspection GoUnusedExportedFunction
func Info(title string, message string) error {
	result := zenity.Info(message, zenity.Title(getTitle(title)))
	return result
}

func Question(title string, message string) (bool, error) {
	result := false
	err := error(nil)
	if err = zenity.Question(message, zenity.Title(getTitle(title))); err == nil {
		result = true
	} else if errors.Is(err, zenity.ErrCanceled) {
		err = nil
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func Warning(title string, message string) error {
	result := zenity.Warning(message, zenity.Title(getTitle(title)))
	return result
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_MESSAGEBOX, logger.SKIP_STACK_FRAMES_BASE)
}

func getTitle(title string) string {
	result := title
	if result == "" {
		result = DEFAULT_TITLE
	}
	return result
}
