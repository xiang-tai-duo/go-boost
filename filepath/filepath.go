// Package filepath
// File:        filepath2.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/filepath/filepath2.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: FilePath provides utility methods for file path operations, including absolute paths, cleaning, and joining
// --------------------------------------------------------------------------------

package filepath2

import (
	__filepath "path/filepath"
	"strings"

	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	MODULE_NAME_FILEPATH = "filepath"
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_FILEPATH, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_FILEPATH, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_FILEPATH, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_FILEPATH, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedExportedFunction
func GetDirectoryName(path string) string {
	return __filepath.Dir(path)
}

//goland:noinspection GoUnusedExportedFunction
func GetFileNameWithoutExtension(path string) string {
	filename := __filepath.Base(path)
	extension := __filepath.Ext(filename)
	return strings.TrimSuffix(filename, extension)
}
