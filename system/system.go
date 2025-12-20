// Package system
// File:        system.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/system/system.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: SYSTEM is a wrapper for operating system interactions, providing cross-platform helpers for environment variables, process/executable info, path operations, command execution, GUID generation, panic recovery and singleton detection.
// --------------------------------------------------------------------------------
package system

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/google/uuid"
	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection GoSnakeCaseUsage
type (
	PANIC_LOGGER func(message string)
)

//goland:noinspection GoSnakeCaseUsage,SpellCheckingInspection
const (
	COMMAND_LINUX_OPEN                 = "xdg-open"
	COMMAND_MACOS_OPEN                 = "open"
	COMMAND_WINDOWS_SHELL              = "cmd"
	COMMAND_WINDOWS_SHELL_ARGUMENT     = "/c"
	COMMAND_WINDOWS_SHELL_START        = "start"
	DEFAULT_TEMP_FILE_PATTERN          = "go-boost-*"
	DEFAULT_TEMP_FILE_PREFIX           = "go-boost-"
	DOT                                = "."
	ERROR_FILE_EXISTS                  = "%s(%s) already exists"
	ERROR_FILE_PATH_EMPTY              = "filePath1 or filePath2 is empty"
	ERROR_PATH_EMPTY                   = "oldFilePath or newFilePath is empty"
	ERROR_UNSUPPORTED_OPERATING_SYSTEM = "unsupported operating system"
	EXIT_CODE_FAILURE                  = 1
	LINUX                              = "linux"
	MACOS                              = "darwin"
	MODULE_NAME_SYSTEM                 = "system"
	PATH_SEPARATOR_UNIX                = '/'
	PATH_SEPARATOR_WINDOWS             = '\\'
	SEPARATOR_UNIX                     = ":"
	SEPARATOR_WINDOWS                  = ";"
	WINDOWS                            = "windows"
)

var isPreviousInstanceRunning = true

func init() {
	if executableFilePath, err := GetExecutableFilePath(); err == nil {
		if err = acquireSingletonLock(executableFilePath); err == nil {
			isPreviousInstanceRunning = false
		} else {
			__debug(fmt.Sprintf("cannot acquire singleton lock: %v", err))
		}
	}
}

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_SYSTEM, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __error(message any) {
	logger.Logger.ErrorEx(message, MODULE_NAME_SYSTEM, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_SYSTEM, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_SYSTEM, logger.SKIP_STACK_FRAMES_BASE)
}

func CreateTemporaryDirectory() string {
	logger.Logger.Debug("CreateTemporaryDirectory enter")
	result := ""
	if dir, err := os.MkdirTemp("", DEFAULT_TEMP_FILE_PREFIX); err == nil {
		result = dir
	}
	logger.Logger.Debug(fmt.Sprintf("CreateTemporaryDirectory result=%s", result))
	return result
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func CreateTemporaryFile(extension string) (string, error) {
	logger.Logger.Debug(fmt.Sprintf("CreateTemporaryFile enter extension=%s", extension))
	result := ""
	err := error(nil)
	pattern := DEFAULT_TEMP_FILE_PATTERN
	if extension != "" {
		if !strings.HasPrefix(extension, DOT) {
			extension = DOT + extension
		}
		pattern = DEFAULT_TEMP_FILE_PREFIX + "*" + extension
	}
	file := (*os.File)(nil)
	if file, err = os.CreateTemp("", pattern); err == nil {
		result = file.Name()
		file.Close()
	}
	logger.Logger.Debug(fmt.Sprintf("CreateTemporaryFile result=%s err=%v", result, err))
	return result, err
}

func ExecuteCommand(command string) error {
	logger.Logger.Debug(fmt.Sprintf("ExecuteCommand enter command=%s", command))
	result := error(nil)
	var cmd *exec.Cmd
	switch {
	case IsWindows():
		cmd = exec.Command(COMMAND_WINDOWS_SHELL, COMMAND_WINDOWS_SHELL_ARGUMENT, COMMAND_WINDOWS_SHELL_START, command)
	case IsMacOS():
		cmd = exec.Command(COMMAND_MACOS_OPEN, command)
	case IsLinux():
		cmd = exec.Command(COMMAND_LINUX_OPEN, command)
	default:
		result = fmt.Errorf(ERROR_UNSUPPORTED_OPERATING_SYSTEM)
	}
	if result == nil {
		result = cmd.Start()
	}
	logger.Logger.Debug(fmt.Sprintf("ExecuteCommand result err=%v", result))
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetEnvironment(key string) string {
	return os.Getenv(key)
}

//goland:noinspection GoUnusedExportedFunction
func GetEnvironments() []string {
	return os.Environ()
}

//goland:noinspection GoUnusedExportedFunction
func GetExecutableFileName() string {
	result := ""
	if executableFilePath, err := os.Executable(); err == nil {
		result = filepath.Base(executableFilePath)
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetExecutableFilePath() (string, error) {
	logger.Logger.Debug("GetExecutableFilePath enter")
	result := ""
	err := error(nil)
	executableFilePath := ""
	if executableFilePath, err = os.Executable(); err == nil {
		result = executableFilePath
		if absoluteExecutableFilePath, absoluteError := filepath.Abs(executableFilePath); absoluteError == nil {
			result = absoluteExecutableFilePath
		}
	}
	logger.Logger.Debug(fmt.Sprintf("GetExecutableFilePath result=%s err=%v", result, err))
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func GetExecuteDirectory() string {
	logger.Logger.Debug("GetExecuteDirectory enter")
	result := ""
	if executableFilePath, err := os.Executable(); err == nil {
		if absoluteExecutableFilePath, absoluteError := filepath.Abs(executableFilePath); absoluteError == nil {
			executableFilePath = absoluteExecutableFilePath
		}
		result = filepath.Dir(executableFilePath)
	}
	logger.Logger.Debug(fmt.Sprintf("GetExecuteDirectory result=%s", result))
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetExpandEnvironment(value string) string {
	return os.ExpandEnv(value)
}

//goland:noinspection GoUnusedExportedFunction
func GetHostName() (string, error) {
	return os.Hostname()
}

//goland:noinspection GoUnusedExportedFunction
func GetTemporaryDirectory() string {
	return os.TempDir()
}

//goland:noinspection GoUnusedExportedFunction
func IsEnvironmentExists(key string) bool {
	return os.Getenv(key) != ""
}

//goland:noinspection GoUnusedExportedFunction
func IsLinux() bool {
	return strings.EqualFold(runtime.GOOS, LINUX)
}

//goland:noinspection GoUnusedExportedFunction
func IsMacOS() bool {
	return strings.EqualFold(runtime.GOOS, MACOS)
}

//goland:noinspection GoUnusedExportedFunction
func IsPreviousInstanceRunning() bool {
	return isPreviousInstanceRunning
}

//goland:noinspection GoUnusedExportedFunction
func IsSameFilePath(filePath1 string, filePath2 string) (bool, error) {
	logger.Logger.Debug(fmt.Sprintf("IsSameFilePath enter filePath1=%s filePath2=%s", filePath1, filePath2))
	result := false
	err := error(nil)
	if filePath1 != "" && filePath2 != "" {
		var absolutePath1 string
		if absolutePath1, err = filepath.Abs(filePath1); err == nil {
			if longPathName, longPathError := GetLongPathName(absolutePath1); longPathError == nil {
				absolutePath1 = longPathName
			}
			var absolutePath2 string
			if absolutePath2, err = filepath.Abs(filePath2); err == nil {
				if longPathName, longPathError := GetLongPathName(absolutePath2); longPathError == nil {
					absolutePath2 = longPathName
				}
				result = strings.EqualFold(absolutePath1, absolutePath2)
			}
		}
	} else {
		err = errors.New(ERROR_FILE_PATH_EMPTY)
	}
	logger.Logger.Debug(fmt.Sprintf("IsSameFilePath result=%v err=%v", result, err))
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func IsUnix() bool {
	return IsMacOS() || IsLinux()
}

//goland:noinspection GoUnusedExportedFunction
func IsWindows() bool {
	return strings.EqualFold(runtime.GOOS, WINDOWS)
}

//goland:noinspection GoUnusedExportedFunction
func LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

//goland:noinspection GoUnusedExportedFunction
func Name() string {
	return runtime.GOOS
}

//goland:noinspection GoUnusedExportedFunction
func NewGUID() string {
	return uuid.NewString()
}

//goland:noinspection GoUnusedExportedFunction
func PathSeparator() rune {
	result := PATH_SEPARATOR_UNIX
	if IsWindows() {
		result = PATH_SEPARATOR_WINDOWS
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func RecoverPanicAndExit(source string, messagePrefix string, log PANIC_LOGGER) {
	if recovered := recover(); recovered != nil {
		if log != nil {
			log(fmt.Sprintf("%s - source=%s, panic=%v\n%s", messagePrefix, source, recovered, string(debug.Stack())))
		}
		os.Exit(EXIT_CODE_FAILURE)
	}
}

//goland:noinspection GoUnusedExportedFunction
func Rename(oldFilePath string, newFilePath string) error {
	return RenameEx(oldFilePath, newFilePath, false)
}

//goland:noinspection GoUnusedExportedFunction
func RenameEx(oldFilePath string, newFilePath string, overwrite bool) error {
	logger.Logger.Debug(fmt.Sprintf("RenameEx enter oldFilePath=%s newFilePath=%s overwrite=%v", oldFilePath, newFilePath, overwrite))
	result := error(nil)
	err := error(nil)
	if oldFilePath != "" && newFilePath != "" {
		var absoluteOldPath string
		if absoluteOldPath, err = filepath.Abs(oldFilePath); err == nil {
			if longPathName, longPathError := GetLongPathName(absoluteOldPath); longPathError == nil {
				absoluteOldPath = longPathName
			}
			var absoluteNewPath string
			if absoluteNewPath, err = filepath.Abs(newFilePath); err == nil {
				var same bool
				if same, err = IsSameFilePath(absoluteOldPath, absoluteNewPath); err == nil && !same {
					if _, err = os.Stat(absoluteOldPath); err == nil {
						if overwrite {
							if _, statError := os.Stat(absoluteNewPath); statError == nil {
								err = os.RemoveAll(absoluteNewPath)
							}
						}
						if err == nil {
							if _, statError := os.Stat(absoluteNewPath); statError != nil {
								if newParentDirectory := filepath.Dir(absoluteNewPath); newParentDirectory != "" {
									err = os.MkdirAll(newParentDirectory, os.ModePerm)
								}
								if err == nil {
									err = os.Rename(absoluteOldPath, absoluteNewPath)
								}
							} else {
								err = fmt.Errorf(ERROR_FILE_EXISTS, newFilePath, absoluteNewPath)
							}
						}
					}
				}
			}
		}
	} else {
		err = errors.New(ERROR_PATH_EMPTY)
	}
	result = err
	logger.Logger.Debug(fmt.Sprintf("RenameEx result err=%v", result))
	return result
}

//goland:noinspection GoUnusedExportedFunction
func SafeGo(source string, messagePrefix string, log PANIC_LOGGER, handler func()) {
	go func() {
		defer RecoverPanicAndExit(source, messagePrefix, log)
		handler()
	}()
}

//goland:noinspection GoUnusedExportedFunction
func Separator() string {
	result := SEPARATOR_UNIX
	if IsWindows() {
		result = SEPARATOR_WINDOWS
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func SetEnvironment(key string, value string) error {
	return os.Setenv(key, value)
}

//goland:noinspection GoUnusedExportedFunction
func UnsetEnvironment(key string) error {
	return os.Unsetenv(key)
}

//goland:noinspection GoUnusedExportedFunction
func UserHomeDirectory() string {
	result := ""
	if home, err := os.UserHomeDir(); err == nil {
		result = home
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func Version() string {
	return runtime.Version()
}
