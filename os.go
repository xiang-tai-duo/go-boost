// --------------------------------------------------------------------------------
// File:        os.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: OS provides utility methods for operating system operations,
//              while OSBuilder is used for operations that require parameters.
// --------------------------------------------------------------------------------

package boost

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// OS provides utility methods for operating system operations related to the current operating system.
type OS struct{}

// Architecture returns the operating system architecture.
// Returns: The operating system architecture string
// Usage:
// arch := OS{}.Architecture() // returns "amd64", "arm64", etc.
func (OS) Architecture() string {
	return runtime.GOARCH
}

// Environment returns the environment variables as a slice of strings.
// Returns: Slice of environment variables
// Usage:
// envVars := OS{}.Environment()
func (OS) Environment() []string {
	return os.Environ()
}

// Executable returns the path to the current executable.
// Returns: Path to the current executable, or empty string if error occurs
// Usage:
// execPath := OS{}.Executable()
func (OS) Executable() string {
	var executablePath string
	execPath, err := os.Executable()
	if err == nil {
		executablePath = execPath
	} else {
		executablePath = ""
	}
	return executablePath
}

// Exit exits the current process with the given status code.
// code: Exit status code
// Usage:
// OS{}.Exit(0) // success exit
func (OS) Exit(code int) {
	os.Exit(code)
}

// ExpandEnvironment replaces ${var} or $var in the string according to the environment variables.
// s: String to expand
// Returns: Expanded string
// Usage:
// expanded := OS.ExpandEnvironment("$HOME/.config")
func (OS) ExpandEnvironment(s string) string {
	return os.ExpandEnv(s)
}

// GetEnvironment returns the value of the environment variable named by the key.
// key: Environment variable name
// Returns: Value of the environment variable
// Usage:
// homeDir := OS{}.GetEnvironment("HOME")
func (OS) GetEnvironment(key string) string {
	return os.Getenv(key)
}

// GetEnvironmentMap returns all environment variables as a map.
// Returns: Map of environment variables
// Usage:
// envMap := OS{}.GetEnvironmentMap()
func (OS) GetEnvironmentMap() map[string]string {
	envMap := make(map[string]string)
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}
	return envMap
}

// GetProcessID returns the process ID of the current process.
// Returns: Process ID
// Usage:
// pid := OS{}.GetProcessID()
func (OS) GetProcessID() int {
	return os.Getpid()
}

// GetParentProcessID returns the parent process ID of the current process.
// Returns: Parent process ID
// Usage:
// ppid := OS{}.GetParentProcessID()
func (OS) GetParentProcessID() int {
	return os.Getppid()
}

// HasEnvironment reports whether the environment variable named by the key exists.
// key: Environment variable name
// Returns: true if environment variable exists, false otherwise
// Usage:
//
//	if OS{}.HasEnvironment("PATH") {
//	    // PATH environment variable exists
//	}
func (OS) HasEnvironment(key string) bool {
	return os.Getenv(key) != ""
}

// Hostname returns the hostname reported by the kernel.
// Returns: Hostname string and any error encountered
// Usage:
// hostname, err := OS{}.Hostname()
func (OS) Hostname() (string, error) {
	return os.Hostname()
}

// IsLinux checks if the operating system is Linux.
// Returns: true if operating system is Linux, false otherwise
// Usage:
//
//	if OS{}.IsLinux() {
//	    // running on Linux
//	}
func (OS) IsLinux() bool {
	return runtime.GOOS == "linux"
}

// IsMac checks if the operating system is macOS.
// Returns: true if operating system is macOS, false otherwise
// Usage:
//
//	if OS{}.IsMac() {
//	    // running on macOS
//	}
func (OS) IsMac() bool {
	return runtime.GOOS == "darwin"
}

// IsUnix checks if the operating system is a Unix-like system (Linux or macOS).
// Returns: true if operating system is Unix-like, false otherwise
// Usage:
//
//	if OS{}.IsUnix() {
//	    // running on Unix-like system
//	}
func (OS) IsUnix() bool {
	return OS{}.IsMac() || OS{}.IsLinux()
}

// IsWindows checks if the operating system is Windows.
// Returns: true if operating system is Windows, false otherwise
// Usage:
//
//	if OS{}.IsWindows() {
//	    // running on Windows
//	}
func (OS) IsWindows() bool {
	return runtime.GOOS == "windows"
}

// Kill kills the process with the given process ID.
// pid: Process ID to kill
// Returns: Error encountered during process termination
// Usage:
// err := OS.Kill(1234)
func (OS) Kill(pid int) error {
	var err error
	process, findErr := os.FindProcess(pid)
	if findErr != nil {
		err = findErr
	} else {
		err = process.Kill()
	}
	return err
}

// LineSeparator returns the line separator for the current operating system.
// Returns: Line separator string (\r\n for Windows, \n for Unix-like)
// Usage:
// lineSep := OS.LineSeparator()
func (o OS) LineSeparator() string {
	var separator string
	if o.IsWindows() {
		separator = "\r\n"
	} else {
		separator = "\n"
	}
	return separator
}

// LookPath searches for an executable file named file in the directories named by the PATH environment variable.
// file: Name of the executable to search for
// Returns: Path to the executable and any error encountered
// Usage:
// cmdPath, err := OS.LookPath("go")
func (OS) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

// Name returns the operating system name.
// Returns: Operating system name ("windows", "linux", "darwin", etc.)
// Usage:
// osName := OS{}.Name()
func (OS) Name() string {
	return runtime.GOOS
}

// PathSeparator returns the path separator for the current operating system.
// Returns: Path separator rune (\ for Windows, / for Unix-like)
// Usage:
// pathSep := OS{}.PathSeparator()
func (o OS) PathSeparator() rune {
	var separator rune
	if o.IsWindows() {
		separator = '\\'
	} else {
		separator = '/'
	}
	return separator
}

// Wait waits for a process to exit and returns its ProcessState.
// If pid is 0, it waits for the current process.
// If pid is non-zero, it waits for the specified process ID.
// Returns: ProcessState and any error encountered
// Usage:
// state, err := OS{}.Wait(0) // Wait for current process
// state, err := OS{}.Wait(1234) // Wait for specified PID
func (OS) Wait(pid int) (*os.ProcessState, error) {
	var state *os.ProcessState
	var err error
	var targetPid int

	if pid == 0 {
		targetPid = os.Getpid()
	} else {
		targetPid = pid
	}

	process, findErr := os.FindProcess(targetPid)
	if findErr != nil {
		err = findErr
	} else {
		state, err = process.Wait()
	}
	return state, err
}

// Separator returns the path list separator for the current operating system.
// Returns: Path list separator string (; for Windows, : for Unix-like)
// Usage:
// pathListSep := OS{}.Separator()
func (o OS) Separator() string {
	var separator string
	if o.IsWindows() {
		separator = ";"
	} else {
		separator = ":"
	}
	return separator
}

// SetEnvironment sets the value of the environment variable named by the key.
// key: Environment variable name
// value: Environment variable value
// Returns: Error encountered during environment variable set
// Usage:
// err := OS.SetEnvironment("MY_VAR", "value")
func (OS) SetEnvironment(key string, value string) error {
	return os.Setenv(key, value)
}

// TemporaryDirectory returns the default directory for temporary files.
// Returns: Path to temporary directory
// Usage:
// tempDir := OS{}.TemporaryDirectory()
func (OS) TemporaryDirectory() string {
	return os.TempDir()
}

// UnsetEnvironment unsets the environment variable named by the key.
// key: Environment variable name
// Returns: Error encountered during environment variable unset
// Usage:
// err := OS.UnsetEnvironment("MY_VAR")
func (OS) UnsetEnvironment(key string) error {
	return os.Unsetenv(key)
}

// UserHomeDirectory returns the current user's home directory.
// Returns: Path to user's home directory, or empty string if error occurs
// Usage:
// homeDir := OS{}.UserHomeDirectory()
func (OS) UserHomeDirectory() string {
	var homeDirectory string
	home, err := os.UserHomeDir()
	if err == nil {
		homeDirectory = home
	} else {
		homeDirectory = ""
	}
	return homeDirectory
}

// Version returns the Go runtime version.
// Returns: Go runtime version string
// Usage:
// goVersion := OS{}.Version()
func (OS) Version() string {
	return runtime.Version()
}
