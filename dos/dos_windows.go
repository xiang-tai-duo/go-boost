// Package dos
// File:        dos_windows.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/dos/dos_windows.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: DOS is a wrapper for Windows DOS commands, providing Go implementations for common file operations.
// --------------------------------------------------------------------------------

//go:build windows

package dos

import (
	"io"
	"os"
	"path/filepath"

	"github.com/xiang-tai-duo/go-boost/logger"
	"golang.org/x/sys/windows"
)

//goland:noinspection SpellCheckingInspection,GoSnakeCaseUsage,GoUnusedConst
const (
	ACTION_OPEN             = "open"
	ACTION_RUNAS            = "runas"
	MODE_DIR                = 0755
	WINDOW_SHOW             = windows.SW_SHOWNORMAL
	WINDOW_HIDE             = windows.SW_HIDE
	MODULE_NAME_DOS_WINDOWS = "dos_windows"
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_DOS_WINDOWS, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_DOS_WINDOWS, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_DOS_WINDOWS, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedExportedFunction
func Call(command string) error {
	err := error(nil)
	cmdPtr := (*uint16)(nil)
	if cmdPtr, err = windows.UTF16PtrFromString(command); err == nil {
		err = callCreateProcess(cmdPtr)
	}
	return err
}

func Cd(dir string) error {
	return os.Chdir(dir)
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func Copy(src string, dst string) error {
	err := error(nil)
	sourceFile := (*os.File)(nil)
	destFile := (*os.File)(nil)
	if sourceFile, err = os.Open(src); err == nil {
		defer sourceFile.Close()
		destDir := filepath.Dir(dst)
		if err = os.MkdirAll(destDir, MODE_DIR); err == nil {
			if destFile, err = os.Create(dst); err == nil && destFile != nil {
				defer destFile.Close()
				if _, err = io.Copy(destFile, sourceFile); err == nil {
					err = destFile.Sync()
				}
			}
		}
	}
	return err
}

//goland:noinspection GoUnusedExportedFunction
func Del(path string) error {
	return os.Remove(path)
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func Deltree(path string) error {
	return os.RemoveAll(path)
}

//goland:noinspection GoUnusedExportedFunction
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

//goland:noinspection GoUnusedExportedFunction
func Mkdir(dir string) error {
	return os.MkdirAll(dir, MODE_DIR)
}

//goland:noinspection GoUnusedExportedFunction
func Move(src, dst string) error {
	if src == dst {
		return nil
	}
	err := error(nil)
	destDir := filepath.Dir(dst)
	if err = os.MkdirAll(destDir, MODE_DIR); err == nil {
		err = os.Rename(src, dst)
	}
	return err
}

func Pwd() (string, error) {
	return os.Getwd()
}

func Rem(message string) {
	__info(message)
}

//goland:noinspection GoUnusedExportedFunction
func Start(path string) error {
	return StartWithAction(path, ACTION_OPEN)
}

//goland:noinspection GoUnusedExportedFunction
func StartWithAction(path, action string) error {
	err := error(nil)
	pathPtr := (*uint16)(nil)
	operationPtr := (*uint16)(nil)
	if pathPtr, err = windows.UTF16PtrFromString(path); err == nil {
		if operationPtr, err = windows.UTF16PtrFromString(action); err == nil {
			err = windows.ShellExecute(0, operationPtr, pathPtr, nil, nil, WINDOW_SHOW)
		}
	}
	return err
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func Xcopy(src, dst string) error {
	err := error(nil)
	srcInfo := (os.FileInfo)(nil)
	if srcInfo, err = os.Stat(src); err == nil {
		if srcInfo.IsDir() {
			err = xcopyDir(src, dst)
		} else {
			err = copyFile(src, dst)
		}
	}
	return err
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_DOS_WINDOWS, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoSnakeCaseUsage,GoUnusedFunction
func callCreateProcess(cmdPtr *uint16) error {
	si := windows.StartupInfo{}
	pi := windows.ProcessInformation{}
	si.ShowWindow = WINDOW_HIDE
	si.Flags = windows.STARTF_USESHOWWINDOW
	return windows.CreateProcess(nil, cmdPtr, nil, nil, false, 0, nil, nil, &si, &pi)
}

//goland:noinspection GoSnakeCaseUsage,GoUnusedFunction,GoUnhandledErrorResult
func copyFile(src, dst string) error {
	err := error(nil)
	sourceFile := (*os.File)(nil)
	destFile := (*os.File)(nil)
	destDir := filepath.Dir(dst)
	if err = os.MkdirAll(destDir, MODE_DIR); err == nil {
		if sourceFile, err = os.Open(src); err == nil {
			defer sourceFile.Close()
			if destFile, err = os.Create(dst); err == nil {
				defer destFile.Close()
				if _, err = io.Copy(destFile, sourceFile); err == nil {
					err = destFile.Sync()
				}
			}
		}
	}
	return err
}

//goland:noinspection GoSnakeCaseUsage,GoUnusedFunction,SpellCheckingInspection
func xcopyDir(src, dst string) error {
	err := error(nil)
	srcInfo := (os.FileInfo)(nil)
	entries := make([]os.DirEntry, 0)
	if srcInfo, err = os.Stat(src); err == nil {
		if err = os.MkdirAll(dst, srcInfo.Mode()); err == nil {
			if entries, err = os.ReadDir(src); err == nil {
				for _, entry := range entries {
					srcPath := filepath.Join(src, entry.Name())
					dstPath := filepath.Join(dst, entry.Name())
					if entry.IsDir() {
						err = xcopyDir(srcPath, dstPath)
					} else {
						err = copyFile(srcPath, dstPath)
					}
					if err != nil {
						break
					}
				}
			}
		}
	}
	return err
}
