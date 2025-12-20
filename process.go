// Package boost
// File:        process.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: Process handles operations related to the current process,
package boost

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type (
	PROCESS         struct{}
	PROCESS_BUILDER struct {
		pid int
	}
)

func NewProcess() *PROCESS {
	return &PROCESS{}
}

func NewProcessBuilder(pid int) PROCESS_BUILDER {
	return PROCESS_BUILDER{pid: pid}
}

func (p PROCESS) ArgumentCount() int {
	return len(os.Args)
}

func (p PROCESS) Arguments() []string {
	return os.Args
}

func (p PROCESS) Environment() []string {
	return os.Environ()
}

func (p PROCESS) GetArgument(index int) string {
	var argument string
	if index >= 0 && index < len(os.Args) {
		argument = os.Args[index]
	}
	return argument
}

func (p PROCESS) GetArgumentValue(flag string) string {
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

func (p PROCESS) GetEnvironment(key string) string {
	return os.Getenv(key)
}

func (p PROCESS) HasArgument(arg string) bool {
	var hasArg bool
	for _, a := range os.Args {
		if a == arg {
			hasArg = true
			break
		}
	}
	return hasArg
}

func (p PROCESS) Name() string {
	return filepath.Base(os.Args[0])
}

func (p PROCESS) ParentProcessID() int {
	return os.Getppid()
}

func (p PROCESS) Path() string {
	var path string
	if execPath, err := os.Executable(); err == nil {
		path = execPath
	}
	return path
}

func (p PROCESS) ProcessID() int {
	return os.Getpid()
}

func (p PROCESS) WorkingDirectory() string {
	var workingDirectory string
	if wd, err := os.Getwd(); err == nil {
		workingDirectory = wd
	}
	return workingDirectory
}

func (PROCESS_BUILDER) CommandExists(cmd string) bool {
	var exists bool
	if _, err := exec.LookPath(cmd); err == nil {
		exists = true
	} else {
		exists = false
	}
	return exists
}

func (PROCESS_BUILDER) ExecuteCommand(name string, args ...string) (*exec.Cmd, error) {
	var cmd *exec.Cmd
	var err error
	cmd = exec.Command(name, args...)
	err = cmd.Start()
	if err != nil {
		cmd = nil
	}
	return cmd, err
}

func (PROCESS_BUILDER) ExecuteCommandAndWait(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

func (PROCESS_BUILDER) ExecuteCommandWithOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func (p PROCESS_BUILDER) Exists() bool {
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

func (p PROCESS_BUILDER) IsCurrent() bool {
	return p.pid == os.Getpid()
}

func (p PROCESS_BUILDER) Kill() error {
	var err error
	process, procErr := p.Process()
	if procErr != nil {
		err = procErr
	} else {
		err = process.Kill()
	}
	return err
}

func (p PROCESS_BUILDER) Parent() PROCESS_BUILDER {
	var ppid int
	if p.IsCurrent() {
		ppid = os.Getppid()
	} else {
		ppid = -1
	}
	return NewProcessBuilder(ppid)
}

func (p PROCESS_BUILDER) Process() (*os.Process, error) {
	return os.FindProcess(p.pid)
}

func (p PROCESS_BUILDER) ProcessID() int {
	return p.pid
}

func (p PROCESS_BUILDER) Signal(sig os.Signal) error {
	var err error
	process, procErr := p.Process()
	if procErr != nil {
		err = procErr
	} else {
		err = process.Signal(sig)
	}
	return err
}

func (p PROCESS_BUILDER) Wait() (*os.ProcessState, error) {
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
