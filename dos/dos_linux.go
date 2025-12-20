// Package dos
// File:        dos_linux.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/dos/dos_linux.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: DOS is a wrapper for Windows DOS commands, providing Go implementations for common file operations.
// --------------------------------------------------------------------------------

//go:build linux
// +build linux

package dos

import (
	"errors"

	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection SpellCheckingInspection,GoSnakeCaseUsage,GoUnusedConst
const (
	ACTION_OPEN           = "open"
	ACTION_RUNAS          = "runas"
	MODE_DIR              = 0755
	MODULE_NAME_DOS_LINUX = "dos_linux"
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_DOS_LINUX, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_DOS_LINUX, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_DOS_LINUX, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedExportedFunction
func Call(command string) error {
	err := error(nil)
	err = errors.New("unsupported on this platform")
	return err
}

func Cd(dir string) error {
	err := error(nil)
	err = errors.New("unsupported on this platform")
	return err
}

//goland:noinspection GoUnusedExportedFunction
func Copy(src, dst string) error {
	err := error(nil)
	err = errors.New("unsupported on this platform")
	return err
}

//goland:noinspection GoUnusedExportedFunction
func Del(path string) error {
	err := error(nil)
	err = errors.New("unsupported on this platform")
	return err
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func Deltree(path string) error {
	err := error(nil)
	err = errors.New("unsupported on this platform")
	return err
}

//goland:noinspection GoUnusedExportedFunction
func Exists(path string) bool {
	return false
}

//goland:noinspection GoUnusedExportedFunction
func Mkdir(dir string) error {
	err := error(nil)
	err = errors.New("unsupported on this platform")
	return err
}

//goland:noinspection GoUnusedExportedFunction
func Move(src, dst string) error {
	err := error(nil)
	err = errors.New("unsupported on this platform")
	return err
}

func Pwd() (string, error) {
	err := errors.New("unsupported on this platform")
	return "", err
}

func Rem(message string) {
}

//goland:noinspection GoUnusedExportedFunction
func Start(path string) error {
	err := error(nil)
	err = errors.New("unsupported on this platform")
	return err
}

//goland:noinspection GoUnusedExportedFunction
func StartWithAction(path, action string) error {
	err := error(nil)
	err = errors.New("unsupported on this platform")
	return err
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func Xcopy(src, dst string) error {
	err := error(nil)
	err = errors.New("unsupported on this platform")
	return err
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_DOS_LINUX, logger.SKIP_STACK_FRAMES_BASE)
}
