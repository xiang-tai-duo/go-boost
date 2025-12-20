// Package ca
// File:        ca.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/ca/ca.go
// Author:      TRAE.AI
// Created:     2026/06/03 20:00:00
// Description: CA centralizes certificate-related file name constants shared across packages to avoid typos.
// --------------------------------------------------------------------------------
package ca

import (
	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection SpellCheckingInspection,GoSnakeCaseUsage,GoUnusedConst
const (
	CERTIFICATE_AUTHORITY_CERTIFICATE_FILE_NAME = "ca.crt"
	CLIENT_CERTIFICATE_FILE_NAME                = "client.crt"
	CLIENT_PRIVATE_KEY_FILE_NAME                = "client.key"
	MODULE_NAME_CA                              = "ca"
	SERVER_CERTIFICATE_FILE_NAME                = "serve.crt"
	SERVER_PRIVATE_KEY_FILE_NAME                = "serve.key"
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_CA, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_CA, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_CA, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_CA, logger.SKIP_STACK_FRAMES_BASE)
}
