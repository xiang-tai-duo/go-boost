// Package process
// File:        process.go
// Url:         https://github.com/xiang-tai-duo/go-bootstrap/blob/master/process/process.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Process handles operations related to the current process
// --------------------------------------------------------------------------------
package process

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-ps"
	"github.com/xiang-tai-duo/go-bootstrap/logger"
)

//goland:noinspection GoUnusedExportedFunction
func ArgumentCount() int {
	return len(os.Args)
}

//goland:noinspection GoUnusedExportedFunction
func Arguments() []string {
	return os.Args
}

//goland:noinspection GoUnusedExportedFunction
func Environment() []string {
	return os.Environ()
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func Exit(processId int, code int) {
	if processId == os.Getpid() {
		os.Exit(code)
	}
	Kill(processId)
}

//goland:noinspection GoUnusedExportedFunction
func GetArgument(index int) string {
	var argument string
	if index >= 0 && index < len(os.Args) {
		argument = os.Args[index]
	}
	return argument
}

//goland:noinspection GoUnusedExportedFunction
func GetArgumentValue(flag string) string {
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

//goland:noinspection GoUnusedExportedFunction
func GetEnvironment(key string) string {
	return os.Getenv(key)
}

//goland:noinspection GoUnusedExportedFunction
func GetParentProcessID(processId int) int {
	processes, err := ps.Processes()
	if err == nil {
		for _, proc := range processes {
			if proc.Pid() == processId {
				return proc.PPid()
			}
		}
	}
	return -1
}

//goland:noinspection GoUnusedExportedFunction
func GetProcessID(processId int) int {
	return processId
}

//goland:noinspection GoUnusedExportedFunction
func HasArgument(arg string) bool {
	var hasArg bool
	for _, a := range os.Args {
		if a == arg {
			hasArg = true
			break
		}
	}
	return hasArg
}

func Kill(processId int) error {
	err := error(nil)
	process := (*os.Process)(nil)
	if process, err = os.FindProcess(processId); err == nil {
		err = process.Kill()
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult
func KillByName(name string) error {
	err := error(nil)
	processes := make([]ps.Process, 0)
	if processes, err = ps.Processes(); err == nil {
		for _, proc := range processes {
			if proc.Executable() == name && proc.Pid() != os.Getpid() {
				pid := proc.Pid()
				logger.Logger.Info(fmt.Sprintf("Kill %d", pid))
				Kill(pid)
			}
		}
	}
	return err
}

//goland:noinspection GoUnusedExportedFunction
func Name() string {
	return filepath.Base(os.Args[0])
}

//goland:noinspection GoUnusedExportedFunction
func ParentProcessID() int {
	return os.Getppid()
}

//goland:noinspection GoUnusedExportedFunction
func Path() string {
	var path string
	if execPath, err := os.Executable(); err == nil {
		path = execPath
	}
	return path
}

//goland:noinspection GoUnusedExportedFunction
func Wait(processId int) (*os.ProcessState, error) {
	var state *os.ProcessState
	var err error
	process, findErr := os.FindProcess(processId)
	if findErr != nil {
		err = findErr
	} else {
		state, err = process.Wait()
	}
	return state, err
}

//goland:noinspection GoUnusedExportedFunction
func IsCommandExists(cmd string) bool {
	var exists bool
	if _, err := exec.LookPath(cmd); err == nil {
		exists = true
	} else {
		exists = false
	}
	return exists
}

//goland:noinspection GoUnusedExportedFunction
func Execute(name string, args ...string) (*exec.Cmd, error) {
	var cmd *exec.Cmd
	var err error
	cmd = exec.Command(name, args...)
	err = cmd.Start()
	if err != nil {
		cmd = nil
	}
	return cmd, err
}

//goland:noinspection GoUnusedExportedFunction
func ExecuteAndWait(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

//goland:noinspection GoUnusedExportedFunction
func ExecuteWithOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

//goland:noinspection GoUnusedExportedFunction
func IsExists(processId int) bool {
	var exists bool
	process, err := os.FindProcess(processId)
	if err == nil {
		err = process.Signal(os.Signal(nil))
		exists = err == nil
	} else {
		exists = false
	}
	return exists
}

//goland:noinspection GoUnusedExportedFunction
func IsCurrent(processId int) bool {
	return processId == os.Getpid()
}

//goland:noinspection GoUnusedExportedFunction
func Parent(processId int) int {
	var ppid int
	if processId == os.Getpid() {
		ppid = os.Getppid()
	} else {
		ppid = -1
	}
	return ppid
}

//goland:noinspection GoUnusedExportedFunction
func Process(processId int) (*os.Process, error) {
	return os.FindProcess(processId)
}

//goland:noinspection GoUnusedExportedFunction
func Signal(processId int, sig os.Signal) error {
	var err error
	process, procErr := Process(processId)
	if procErr != nil {
		err = procErr
	} else {
		err = process.Signal(sig)
	}
	return err
}

//goland:noinspection GoUnusedExportedFunction
func WaitProcess(processId int) (*os.ProcessState, error) {
	var state *os.ProcessState
	var err error
	process, procErr := Process(processId)
	if procErr != nil {
		err = procErr
	} else {
		state, err = process.Wait()
	}
	return state, err
}
