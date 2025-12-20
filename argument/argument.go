// Package argument
// File:        argument.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/argument/argument.go
// Author:      TRAE.AI
// Created:     2026/05/21 12:00:00
// Description: ARGUMENT is a wrapper for command-line argument operations, providing a set of methods for argument parsing.
// --------------------------------------------------------------------------------

package argument

import (
	"github.com/xiang-tai-duo/go-boost/logger"
	"github.com/xiang-tai-duo/go-boost/strings2"
	"os"
)

const (
	EMPTY_STRING         = ""
	MODULE_NAME_ARGUMENT = "argument"
	PREFIX_DASH          = "-"
	PREFIX_LENGTH        = 1
	STEP_INCREMENT       = 1
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_ARGUMENT, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_ARGUMENT, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_ARGUMENT, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedExportedFunction
func Get(key string) string {
	result := EMPTY_STRING
	args := os.Args
	found := false
	for i := 0; i < len(args) && !found; i += STEP_INCREMENT {
		if args[i] == key && i+STEP_INCREMENT < len(args) {
			next := args[i+STEP_INCREMENT]
			if strings2.Left(next, PREFIX_LENGTH) != PREFIX_DASH {
				result = next
			}
			found = true
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func IsExist(key string) bool {
	result := false
	args := os.Args
	for i := 0; i < len(args) && !result; i += STEP_INCREMENT {
		if args[i] == key {
			result = true
		}
	}
	return result
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_ARGUMENT, logger.SKIP_STACK_FRAMES_BASE)
}
