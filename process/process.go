// Package process
// File:        process.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/process/process.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: Process handles operations related to the current process
// --------------------------------------------------------------------------------
package process

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-ps"
)

type (
	PROCESS struct {
		pid int
	}
)

func New(pid ...int) *PROCESS {
	p := &PROCESS{}
	if len(pid) > 0 {
		p.pid = pid[0]
	} else {
		p.pid = os.Getpid()
	}
	return p
}

func (p PROCESS) ArgumentCount() int {
	if p.pid == os.Getpid() {
		return len(os.Args)
	}
	return 0
}

func (p PROCESS) Arguments() []string {
	if p.pid == os.Getpid() {
		return os.Args
	}
	return []string{}
}

func (p PROCESS) Environment() []string {
	if p.pid == os.Getpid() {
		return os.Environ()
	}
	return []string{}
}

//goland:noinspection GoUnhandledErrorResult
func (p PROCESS) Exit(code int) {
	if p.pid == os.Getpid() {
		os.Exit(code)
	} else {
		p.Kill()
	}
}

func (p PROCESS) GetArgument(index int) string {
	if p.pid == os.Getpid() {
		var argument string
		if index >= 0 && index < len(os.Args) {
			argument = os.Args[index]
		}
		return argument
	}
	return ""
}

func (p PROCESS) GetArgumentValue(flag string) string {
	if p.pid == os.Getpid() {
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
	return ""
}

func (p PROCESS) GetEnvironment(key string) string {
	if p.pid == os.Getpid() {
		return os.Getenv(key)
	}
	return ""
}

func (p PROCESS) GetParentProcessID() int {
	if p.pid == os.Getpid() {
		return os.Getppid()
	}

	// 使用go-ps库获取指定进程的父进程ID
	processes, err := ps.Processes()
	if err == nil {
		for _, proc := range processes {
			if proc.Pid() == p.pid {
				return proc.PPid()
			}
		}
	}
	return -1
}

func (p PROCESS) GetProcessID() int {
	return p.pid
}

func (p PROCESS) HasArgument(arg string) bool {
	if p.pid == os.Getpid() {
		var hasArg bool
		for _, a := range os.Args {
			if a == arg {
				hasArg = true
				break
			}
		}
		return hasArg
	}
	return false
}

func (p PROCESS) Kill() error {
	var err error
	process, findErr := os.FindProcess(p.pid)
	if findErr != nil {
		err = findErr
	} else {
		err = process.Kill()
	}
	return err
}

func (p PROCESS) Name() string {
	if p.pid == os.Getpid() {
		return filepath.Base(os.Args[0])
	}
	return ""
}

func (p PROCESS) ParentProcessID() int {
	if p.pid == os.Getpid() {
		return os.Getppid()
	}
	return -1
}

func (p PROCESS) Path() string {
	if p.pid == os.Getpid() {
		var path string
		if execPath, err := os.Executable(); err == nil {
			path = execPath
		}
		return path
	}
	return ""
}

func (p PROCESS) ProcessID() int {
	if p.pid == os.Getpid() {
		return os.Getpid()
	}
	return p.pid
}

func (p PROCESS) Wait() (*os.ProcessState, error) {
	var state *os.ProcessState
	var err error
	process, findErr := os.FindProcess(p.pid)
	if findErr != nil {
		err = findErr
	} else {
		state, err = process.Wait()
	}
	return state, err
}

func (p PROCESS) WorkingDirectory() string {
	if p.pid == os.Getpid() {
		var workingDirectory string
		if wd, err := os.Getwd(); err == nil {
			workingDirectory = wd
		}
		return workingDirectory
	}
	return ""
}

func (p PROCESS) CommandExists(cmd string) bool {
	var exists bool
	if _, err := exec.LookPath(cmd); err == nil {
		exists = true
	} else {
		exists = false
	}
	return exists
}

func (p PROCESS) ExecuteCommand(name string, args ...string) (*exec.Cmd, error) {
	var cmd *exec.Cmd
	var err error
	cmd = exec.Command(name, args...)
	err = cmd.Start()
	if err != nil {
		cmd = nil
	}
	return cmd, err
}

func (p PROCESS) ExecuteCommandAndWait(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

func (p PROCESS) ExecuteCommandWithOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func (p PROCESS) Exists() bool {
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

func (p PROCESS) IsCurrent() bool {
	return p.pid == os.Getpid()
}

func (p PROCESS) Parent() *PROCESS {
	var ppid int
	if p.IsCurrent() {
		ppid = os.Getppid()
	} else {
		ppid = -1
	}
	return New(ppid)
}

func (p PROCESS) Process() (*os.Process, error) {
	return os.FindProcess(p.pid)
}

func (p PROCESS) Signal(sig os.Signal) error {
	var err error
	process, procErr := p.Process()
	if procErr != nil {
		err = procErr
	} else {
		err = process.Signal(sig)
	}
	return err
}

func (p PROCESS) WaitProcess() (*os.ProcessState, error) {
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
