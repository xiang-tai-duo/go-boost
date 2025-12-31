// Package boost
// File:        os.go
// Author:      TRAE AI
// Created:     12/30/2025 11:03:46
// Description: OS provides utility methods for operating system operations,
package boost

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type (
	OS struct{}
)

func (o OS) Architecture() string {
	return runtime.GOARCH
}

func (o OS) Environment() []string {
	return os.Environ()
}

func (o OS) Executable() string {
	var executablePath string
	if execPath, err := os.Executable(); err == nil {
		executablePath = execPath
	} else {
		executablePath = ""
	}
	return executablePath
}

func (o OS) ExpandEnvironment(s string) string {
	return os.ExpandEnv(s)
}

func (o OS) GetEnvironment(key string) string {
	return os.Getenv(key)
}

func (o OS) GetEnvironmentMap() map[string]string {
	envMap := make(map[string]string)
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}
	return envMap
}

func (o OS) HasEnvironment(key string) bool {
	return os.Getenv(key) != ""
}

func (o OS) Hostname() (string, error) {
	return os.Hostname()
}

func (o OS) IsLinux() bool {
	return runtime.GOOS == "linux"
}

func (o OS) IsMac() bool {
	return runtime.GOOS == "darwin"
}

func (o OS) IsUnix() bool {
	return o.IsMac() || o.IsLinux()
}

func (o OS) IsWindows() bool {
	return runtime.GOOS == "windows"
}

func (o OS) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

func (o OS) Name() string {
	return runtime.GOOS
}

func (o OS) PathSeparator() rune {
	var separator rune
	if o.IsWindows() {
		separator = '\\'
	} else {
		separator = '/'
	}
	return separator
}

func (o OS) Separator() string {
	var separator string
	if o.IsWindows() {
		separator = ";"
	} else {
		separator = ":"
	}
	return separator
}

func (o OS) SetEnvironment(key string, value string) error {
	return os.Setenv(key, value)
}

func (o OS) TemporaryDirectory() string {
	return os.TempDir()
}

func (o OS) UnsetEnvironment(key string) error {
	return os.Unsetenv(key)
}

func (o OS) UserHomeDirectory() string {
	var homeDirectory string
	if home, err := os.UserHomeDir(); err == nil {
		homeDirectory = home
	} else {
		homeDirectory = ""
	}
	return homeDirectory
}

func (o OS) Version() string {
	return runtime.Version()
}
