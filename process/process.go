// Package process
// File:        process.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/process/process.go
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
	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	MODULE_NAME_PROCESS = "process"
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_PROCESS, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_PROCESS, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_PROCESS, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedExportedFunction
func ArgumentCount() int {
	result := len(os.Args)
	return result
}

//goland:noinspection GoUnusedExportedFunction
func Arguments() []string {
	result := os.Args
	return result
}

//goland:noinspection GoUnusedExportedFunction
func Environment() []string {
	result := os.Environ()
	return result
}

//goland:noinspection GoUnusedExportedFunction
func Execute(name string, args ...string) (*exec.Cmd, error) {
	err := error(nil)
	result := exec.Command(name, args...)
	if err = result.Start(); err != nil {
		result = nil
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func ExecuteAndWait(name string, args ...string) error {
	err := error(nil)
	cmd := exec.Command(name, args...)
	err = cmd.Run()
	return err
}

//goland:noinspection GoUnusedExportedFunction
func ExecuteWithOutput(name string, args ...string) (string, error) {
	err := error(nil)
	result := ""
	cmd := exec.Command(name, args...)
	var output []byte
	output, err = cmd.CombinedOutput()
	result = string(output)
	return result, err
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func Exit(processId int, code int) {
	if processId == os.Getpid() {
		os.Exit(code)
	} else {
		Kill(processId)
	}
}

//goland:noinspection GoUnusedExportedFunction
func GetArgument(index int) string {
	result := ""
	if index >= 0 && index < len(os.Args) {
		result = os.Args[index]
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetArgumentValue(flag string) string {
	result := ""
	for i, arg := range os.Args {
		if arg == flag && i+1 < len(os.Args) {
			result = os.Args[i+1]
			break
		}
		if strings.HasPrefix(arg, flag+"=") {
			result = strings.TrimPrefix(arg, flag+"=")
			break
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetEnvironment(key string) string {
	result := os.Getenv(key)
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetParentProcessID(processId int) int {
	result := -1
	processes := make([]ps.Process, 0)
	err := error(nil)
	if processes, err = ps.Processes(); err == nil {
		for _, proc := range processes {
			if proc.Pid() == processId {
				result = proc.PPid()
				break
			}
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetProcessID(processId int) int {
	result := processId
	return result
}

//goland:noinspection GoUnusedExportedFunction
func HasArgument(arg string) bool {
	result := false
	for _, a := range os.Args {
		if a == arg {
			result = true
			break
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func IsCommandExists(cmd string) bool {
	result := false
	if _, err := exec.LookPath(cmd); err == nil {
		result = true
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func IsCurrent(processId int) bool {
	result := processId == os.Getpid()
	return result
}

//goland:noinspection GoUnusedExportedFunction
func IsExists(processId int) bool {
	result := false
	if process, err := os.FindProcess(processId); err == nil {
		if err = process.Signal(os.Signal(nil)); err == nil {
			result = true
		}
	}
	return result
}

func Kill(processId int) error {
	err := error(nil)
	var process *os.Process
	if process, err = os.FindProcess(processId); err == nil {
		err = process.Kill()
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult
func KillByFilePath(filePath string) error {
	err := error(nil)
	processes := make([]ps.Process, 0)
	if processes, err = ps.Processes(); err == nil {
		for _, proc := range processes {
			if strings.EqualFold(proc.Executable(), filePath) && proc.Pid() != os.Getpid() {
				pid := proc.Pid()
				__info(fmt.Sprintf("Kill %d", pid))
				Kill(pid)
			}
		}
		if processes, err = ps.Processes(); err == nil {
			for _, proc := range processes {
				if strings.EqualFold(proc.Executable(), filePath) && proc.Pid() != os.Getpid() {
					pid := proc.Pid()
					err = fmt.Errorf("process %d still exists after kill: %s", pid, filePath)
					__debug(err.Error())
					break
				}
			}
		}
	}
	return err
}

//goland:noinspection GoUnusedExportedFunction
func Name() string {
	result := filepath.Base(os.Args[0])
	return result
}

//goland:noinspection GoUnusedExportedFunction
func Parent(processId int) int {
	result := -1
	if processId == os.Getpid() {
		result = os.Getppid()
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func ParentProcessID() int {
	result := os.Getppid()
	return result
}

//goland:noinspection GoUnusedExportedFunction
func Path() string {
	result := ""
	if execPath, err := os.Executable(); err == nil {
		result = execPath
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func Process(processId int) (*os.Process, error) {
	err := error(nil)
	var result *os.Process
	result, err = os.FindProcess(processId)
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func Signal(processId int, sig os.Signal) error {
	err := error(nil)
	var process *os.Process
	if process, err = Process(processId); err == nil {
		err = process.Signal(sig)
	}
	return err
}

//goland:noinspection GoUnusedExportedFunction
func Wait(processId int) (*os.ProcessState, error) {
	err := error(nil)
	var result *os.ProcessState
	var process *os.Process
	if process, err = os.FindProcess(processId); err == nil {
		result, err = process.Wait()
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func WaitProcess(processId int) (*os.ProcessState, error) {
	err := error(nil)
	var result *os.ProcessState
	var process *os.Process
	if process, err = Process(processId); err == nil {
		result, err = process.Wait()
	}
	return result, err
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_PROCESS, logger.SKIP_STACK_FRAMES_BASE)
}
