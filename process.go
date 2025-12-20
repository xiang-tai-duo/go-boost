// --------------------------------------------------------------------------------
// File:        process.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: Process handles operations related to the current process,
//              while ProcessBuilder is used for operations on other processes.
// --------------------------------------------------------------------------------

package boost

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Process provides utility methods specifically for the current process.
type Process struct{}

// ProcessBuilder provides utility methods for operations on other processes.
type ProcessBuilder struct {
	pid int
}

// NewProcessBuilder creates a new ProcessBuilder instance with the given PID.
// pid: Process ID to create the builder for
// Returns: New ProcessBuilder instance
// Usage:
// procBuilder := NewProcessBuilder(1234)
func NewProcessBuilder(pid int) ProcessBuilder {
	return ProcessBuilder{pid: pid}
}

// NewProcess creates a new Process instance.
// Usage:
// process := NewProcess()
func NewProcess() *Process {
	return &Process{}
}

// ArgumentCount returns the number of arguments passed to the current process.
// Returns: Number of arguments
// Usage:
// argCount := Process{}.ArgumentCount()
func (p Process) ArgumentCount() int {
	return len(os.Args)
}

// Arguments returns the command-line arguments passed to the current process.
// Returns: Slice of arguments
// Usage:
// args := Process{}.Arguments()
func (p Process) Arguments() []string {
	return os.Args
}

// Environment returns the environment variables for the current process.
// Returns: Slice of environment variables
// Usage:
// envVars := Process{}.Environment()
func (p Process) Environment() []string {
	return os.Environ()
}

// GetArgument returns the command-line argument at the specified index.
// index: Index of the argument to retrieve
// Returns: Argument string, or empty string if index is out of range
// Usage:
// arg := Process{}.GetArgument(0) // returns the command name
func (p Process) GetArgument(index int) string {
	var argument string
	if index >= 0 && index < len(os.Args) {
		argument = os.Args[index]
	}
	return argument
}

// GetArgumentValue returns the value of a command-line flag for the current process.
// flag: Flag name to look for (e.g., "--output" or "-o")
// Returns: Flag value, or empty string if flag not found
// Usage:
// outputPath := Process{}.GetArgumentValue("--output")
func (p Process) GetArgumentValue(flag string) string {
	var value string
	for i, arg := range os.Args {
		if arg == flag && i+1 < len(os.Args) {
			value = os.Args[i+1]
			break
		}
		if strings.HasPrefix(arg, flag+"=") {
			value = strings.TrimPrefix(arg, flag+"=")
			break
		}
	}
	return value
}

// GetEnvironment returns the value of the environment variable named by the key for the current process.
// key: Environment variable name
// Returns: Value of the environment variable
// Usage:
// homeDir := Process{}.GetEnvironment("HOME")
func (p Process) GetEnvironment(key string) string {
	return os.Getenv(key)
}

// HasArgument checks if the current process was started with the given command-line argument.
// arg: Argument to check for
// Returns: true if argument exists, false otherwise
// Usage:
//
//	if Process{}.HasArgument("--verbose") {
//	    // verbose mode enabled
//	}
func (p Process) HasArgument(arg string) bool {
	var hasArg bool
	for _, a := range os.Args {
		if a == arg {
			hasArg = true
			break
		}
	}
	return hasArg
}

// Name returns the name of the current process executable.
// Returns: Process name
// Usage:
// procName := Process{}.Name()
func (p Process) Name() string {
	return filepath.Base(os.Args[0])
}

// ParentProcessID returns the parent process ID of the current process.
// Returns: Parent process ID
// Usage:
// parentPID := Process{}.ParentProcessID()
func (p Process) ParentProcessID() int {
	return os.Getppid()
}

// Path returns the path to the current process executable.
// Returns: Executable path, or empty string if error occurs
// Usage:
// execPath := Process{}.Path()
func (p Process) Path() string {
	var path string
	execPath, err := os.Executable()
	if err == nil {
		path = execPath
	}
	return path
}

// ProcessID returns the process ID of the current process.
// Returns: Process ID
// Usage:
// processID := Process{}.ProcessID()
func (p Process) ProcessID() int {
	return os.Getpid()
}

// WorkingDirectory returns the working directory of the current process.
// Returns: Working directory path, or empty string if error occurs
// Usage:
// workingDirectory := Process{}.WorkingDirectory()
func (p Process) WorkingDirectory() string {
	var workingDirectory string
	wd, err := os.Getwd()
	if err == nil {
		workingDirectory = wd
	}
	return workingDirectory
}

// CommandExists checks if the given command exists in the PATH.
// cmd: Command name to check
// Returns: true if command exists, false otherwise
// Usage:
//
//	if ProcessBuilder(0).CommandExists("go") {
//	    // go command exists
//	}
func (ProcessBuilder) CommandExists(cmd string) bool {
	var exists bool
	_, err := exec.LookPath(cmd)
	if err == nil {
		exists = true
	} else {
		exists = false
	}
	return exists
}

// ExecuteCommand starts a new command process but does not wait for it to complete.
// name: Command name to execute
// args: Command arguments
// Returns: Executed command and any error encountered
// Usage:
// cmd, err := ProcessBuilder(0).ExecuteCommand("sleep", "5")
func (ProcessBuilder) ExecuteCommand(name string, args ...string) (*exec.Cmd, error) {
	var cmd *exec.Cmd
	var err error
	cmd = exec.Command(name, args...)
	err = cmd.Start()
	if err != nil {
		cmd = nil
	}
	return cmd, err
}

// ExecuteCommandAndWait starts a new command process and waits for it to complete.
// name: Command name to execute
// args: Command arguments
// Returns: Error encountered during command execution
// Usage:
// err := ProcessBuilder(0).ExecuteCommandAndWait("echo", "Hello, World!")
func (ProcessBuilder) ExecuteCommandAndWait(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

// ExecuteCommandWithOutput starts a new command process, waits for it to complete, and returns its output.
// name: Command name to execute
// args: Command arguments
// Returns: Combined stdout and stderr output, and any error encountered
// Usage:
// output, err := ProcessBuilder(0).ExecuteCommandWithOutput("ls", "-la")
func (ProcessBuilder) ExecuteCommandWithOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// Exists checks if the process with the given PID exists.
// Returns: true if process exists, false otherwise
// Usage:
//
//	if ProcessBuilder(1234).Exists() {
//	    // process exists
//	}
func (p ProcessBuilder) Exists() bool {
	var exists bool
	process, err := os.FindProcess(p.pid)
	if err == nil {
		err = process.Signal(os.Signal(nil))
		exists = (err == nil)
	} else {
		exists = false
	}
	return exists
}

// IsCurrent checks if the ProcessBuilder represents the current process.
// Returns: true if it's the current process, false otherwise
// Usage:
//
//	if proc.IsCurrent() {
//	    // this is the current process
//	}
func (p ProcessBuilder) IsCurrent() bool {
	return p.pid == os.Getpid()
}

// Kill terminates the process.
// Returns: Error encountered during process termination
// Usage:
// err := ProcessBuilder(1234).Kill()
func (p ProcessBuilder) Kill() error {
	var err error
	process, procErr := p.Process()
	if procErr != nil {
		err = procErr
	} else {
		err = process.Kill()
	}
	return err
}

// Parent returns a ProcessBuilder representing the parent process.
// Returns: ProcessBuilder for parent process
// Usage:
// parentProc := ProcessBuilder(1234).Parent()
func (p ProcessBuilder) Parent() ProcessBuilder {
	var ppid int
	if p.IsCurrent() {
		ppid = os.Getppid()
	} else {
		ppid = -1
	}
	return NewProcessBuilder(ppid)
}

// ProcessID returns the process ID.
// Returns: Process ID
// Usage:
// processID := ProcessBuilder(1234).ProcessID() // returns 1234
func (p ProcessBuilder) ProcessID() int {
	return p.pid
}

// Process returns the underlying os.Process.
// Returns: os.Process pointer and any error encountered
// Usage:
// process, err := ProcessBuilder(1234).Process()
func (p ProcessBuilder) Process() (*os.Process, error) {
	return os.FindProcess(p.pid)
}

// Signal sends a signal to the process.
// sig: Signal to send
// Returns: Error encountered during signal delivery
// Usage:
// err := ProcessBuilder(1234).Signal(os.Interrupt)
func (p ProcessBuilder) Signal(sig os.Signal) error {
	var err error
	process, procErr := p.Process()
	if procErr != nil {
		err = procErr
	} else {
		err = process.Signal(sig)
	}
	return err
}

// Wait waits for the process to exit and returns its ProcessState.
// Returns: ProcessState and any error encountered
// Usage:
// state, err := ProcessBuilder(1234).Wait()
func (p ProcessBuilder) Wait() (*os.ProcessState, error) {
	var state *os.ProcessState
	var err error
	process, procErr := p.Process()
	if procErr != nil {
		err = procErr
	} else {
		state, err = process.Wait()
	}
	return state, err
}
