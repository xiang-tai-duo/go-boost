// Package os
// File:        os.go
// Url:         https://github.com/xiang-tai-duo/go-bootstrap/blob/master/os/os.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: OS provides utility methods for operating system operations, including environment variables and platform detection
// --------------------------------------------------------------------------------
package system

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/xiang-tai-duo/go-bootstrap/directory"
)

const (
	WINDOWS = "windows"
	LINUX   = "linux"
	MACOS   = "darwin"
)

//goland:noinspection GoUnusedExportedFunction
func Architecture() string {
	return runtime.GOARCH
}

//goland:noinspection GoUnusedExportedFunction
func CreateTemporaryDirectory() string {
	return directory.CreateTemporaryDirectory()
}

//goland:noinspection GoUnusedExportedFunction
func Environments() []string {
	return os.Environ()
}

//goland:noinspection GoUnusedExportedFunction
func ExecutableFilePath() string {
	path := ""
	if exeFilePath, err := os.Executable(); err == nil {
		path = exeFilePath
	}
	return path
}

//goland:noinspection GoUnusedExportedFunction
func ExecutableFileName() string {
	path := ""
	if exeFilePath, err := os.Executable(); err == nil {
		path = exeFilePath
	}
	return filepath.Base(path)
}

//goland:noinspection GoUnusedExportedFunction
func ExpandEnvironment(s string) string {
	return os.ExpandEnv(s)
}

//goland:noinspection GoUnusedExportedFunction
func GetEnvironment(key string) string {
	return os.Getenv(key)
}

//goland:noinspection GoUnusedExportedFunction
func GetTemporaryDirectory() string {
	return os.TempDir()
}

//goland:noinspection GoUnusedExportedFunction
func HasEnvironment(key string) bool {
	return os.Getenv(key) != ""
}

//goland:noinspection GoUnusedExportedFunction
func Hostname() (string, error) {
	return os.Hostname()
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
	return os.Setenv(key, value)
}

//goland:noinspection GoUnusedExportedFunction
func UnsetEnvironment(key string) error {
	return os.Unsetenv(key)
}

//goland:noinspection GoUnusedExportedFunction
func UserHomeDirectory() string {
	var homeDirectory string
	if home, err := os.UserHomeDir(); err == nil {
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
