// Package os
// File:        os.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/os/os.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: OS provides utility methods for operating system operations, including environment variables and platform detection
// --------------------------------------------------------------------------------
package os

import (
	__os "os"
	"os/exec"
	"runtime"
	"strings"
)

//goland:noinspection GoUnusedExportedFunction
func Architecture() string {
	return runtime.GOARCH
}

//goland:noinspection GoUnusedExportedFunction
func Environments() []string {
	return __os.Environ()
}

//goland:noinspection GoUnusedExportedFunction
func ExecutableFilePath() string {
	var path string
	if execPath, err := __os.Executable(); err == nil {
		path = execPath
	} else {
		path = ""
	}
	return path
}

//goland:noinspection GoUnusedExportedFunction
func ExpandEnvironment(s string) string {
	return __os.ExpandEnv(s)
}

//goland:noinspection GoUnusedExportedFunction
func GetEnvironment(key string) string {
	return __os.Getenv(key)
}

//goland:noinspection GoUnusedExportedFunction
func HasEnvironment(key string) bool {
	return __os.Getenv(key) != ""
}

//goland:noinspection GoUnusedExportedFunction
func Hostname() (string, error) {
	return __os.Hostname()
}

//goland:noinspection SpellCheckingInspection
func IsLinux() bool {
	return strings.ToLower(runtime.GOOS) == "linux"
}

//goland:noinspection SpellCheckingInspection
func IsMac() bool {
	return strings.ToLower(runtime.GOOS) == "darwin"
}

//goland:noinspection GoUnusedExportedFunction
func IsUnix() bool {
	return IsMac() || IsLinux()
}

//goland:noinspection GoUnusedExportedFunction
func IsWindows() bool {
	return strings.ToLower(runtime.GOOS) == "windows"
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
func PathSeparator() rune {
	var separator rune
	if IsWindows() {
		separator = '\\'
	} else {
		separator = '/'
	}
	return separator
}

//goland:noinspection GoUnusedExportedFunction
func Separator() string {
	var separator string
	if IsWindows() {
		separator = ";"
	} else {
		separator = ":"
	}
	return separator
}

//goland:noinspection GoUnusedExportedFunction
func SetEnvironment(key string, value string) error {
	return __os.Setenv(key, value)
}

//goland:noinspection GoUnusedExportedFunction
func TemporaryDirectory() string {
	return __os.TempDir()
}

//goland:noinspection GoUnusedExportedFunction
func UnsetEnvironment(key string) error {
	return __os.Unsetenv(key)
}

//goland:noinspection GoUnusedExportedFunction
func UserHomeDirectory() string {
	var homeDirectory string
	if home, err := __os.UserHomeDir(); err == nil {
		homeDirectory = home
	} else {
		homeDirectory = ""
	}
	return homeDirectory
}

//goland:noinspection GoUnusedExportedFunction
func Version() string {
	return runtime.Version()
}
