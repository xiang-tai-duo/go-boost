//go:build windows

// Package service
// File:        service_windows.go
// Url:         https://github.com/xiang-tai-duo/go-bootstrap/blob/master/service/service_windows.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Windows service management functions
// --------------------------------------------------------------------------------
package service

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/xiang-tai-duo/go-boost/hash"
	"github.com/xiang-tai-duo/go-boost/logger"
	"github.com/xiang-tai-duo/go-boost/system"
	__svc "golang.org/x/sys/windows/svc"
	__mgr "golang.org/x/sys/windows/svc/mgr"
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
type (
	SERVICE_HANDLER struct {
		OnStart SERVICE_START_HANDLER
		OnStop  SERVICE_STOP_HANDLER
	}
	SERVICE_START_HANDLER func()
	SERVICE_STOP_HANDLER  func()
)

//goland:noinspection GoSnakeCaseUsage,SpellCheckingInspection,GoNameStartsWithPackageName,GoUnusedConst
const (
	MODULE_NAME_SERVICE_WINDOWS = "service_windows"
	SERVICE_COMMAND_ACCEPTED    = __svc.AcceptStop | __svc.AcceptShutdown
	SERVICE_STOP_POLL_INTERVAL  = 500 * time.Millisecond
	SERVICE_STOP_TIMEOUT        = 1 * time.Minute
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_SERVICE_WINDOWS, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_SERVICE_WINDOWS, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_SERVICE_WINDOWS, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func EnableService(serviceName string) error {
	err := error(nil)
	var mgr *__mgr.Mgr
	if mgr, err = __mgr.Connect(); err == nil {
		defer mgr.Disconnect()
		var svc *__mgr.Service
		if svc, err = mgr.OpenService(serviceName); err == nil {
			var config __mgr.Config
			if config, err = svc.Config(); err == nil {
				if config.StartType != __mgr.StartAutomatic {
					config.StartType = __mgr.StartAutomatic
					err = svc.UpdateConfig(config)
				}
			}
			svc.Close()
		}
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult,GoUnusedParameter
func (s *SERVICE_HANDLER) Execute(args []string, r <-chan __svc.ChangeRequest, status chan<- __svc.Status) (bool, uint32) {
	result := uint32(0)
	status <- __svc.Status{State: __svc.StartPending}
	if s.OnStart != nil {
		s.OnStart()
	}
	status <- __svc.Status{State: __svc.Running, Accepts: SERVICE_COMMAND_ACCEPTED}
loop:
	for {
		c := <-r
		switch c.Cmd {
		case __svc.Interrogate:
			status <- c.CurrentStatus
		case __svc.Stop, __svc.Shutdown:
			break loop
		default:
			__error(fmt.Sprintf("SERVICE_HANDLER: 收到未知的服务控制命令: %v", c.Cmd))
		}
	}
	status <- __svc.Status{State: __svc.StopPending}
	if s.OnStop != nil {
		s.OnStop()
	}
	status <- __svc.Status{State: __svc.Stopped}
	return false, result
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func GetServiceBinaryPath(serviceName string) (string, error) {
	result := ""
	err := error(nil)
	var mgr *__mgr.Mgr
	if mgr, err = __mgr.Connect(); err == nil {
		defer mgr.Disconnect()
		var svc *__mgr.Service
		if svc, err = mgr.OpenService(serviceName); err == nil {
			var config __mgr.Config
			if config, err = svc.Config(); err == nil {
				result, _ = filepath.Abs(strings.Trim(config.BinaryPathName, `"`))
			}
			svc.Close()
		}
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func GetServiceStatus(serviceName string) (uint32, error) {
	result := uint32(0)
	err := error(nil)
	var mgr *__mgr.Mgr
	if mgr, err = __mgr.Connect(); err == nil {
		defer mgr.Disconnect()
		var svc *__mgr.Service
		if svc, err = mgr.OpenService(serviceName); err == nil {
			var status __svc.Status
			if status, err = svc.Query(); err == nil {
				result = uint32(status.State)
			}
			svc.Close()
		}
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func IsServiceAutomatic(serviceName string) (bool, error) {
	result := false
	err := error(nil)
	var mgr *__mgr.Mgr
	if mgr, err = __mgr.Connect(); err == nil {
		defer mgr.Disconnect()
		var svc *__mgr.Service
		if svc, err = mgr.OpenService(serviceName); err == nil {
			var config __mgr.Config
			if config, err = svc.Config(); err == nil {
				result = config.StartType == __mgr.StartAutomatic
			}
			svc.Close()
		}
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func IsServiceExists(serviceName string) bool {
	result := false
	err := error(nil)
	var mgr *__mgr.Mgr
	if mgr, err = __mgr.Connect(); err == nil {
		defer mgr.Disconnect()
		var svc *__mgr.Service
		if svc, err = mgr.OpenService(serviceName); err == nil {
			result = true
			svc.Close()
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func IsWindowsService() (bool, error) {
	result := false
	err := error(nil)
	result, err = __svc.IsWindowsService()
	return result, err
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func RegisterService(serviceName string, displayName string, description string) error {
	err := error(nil)
	exeFilePath := ""
	if exeFilePath, err = system.GetExecutableFilePath(); err == nil {
		var mgr *__mgr.Mgr
		if mgr, err = __mgr.Connect(); err == nil {
			defer mgr.Disconnect()
			var svc *__mgr.Service
			if svc, err = mgr.OpenService(serviceName); err == nil {
				var config __mgr.Config
				if config, err = svc.Config(); err == nil {
					binaryPathName := ""
					if binaryPathName, err = filepath.Abs(strings.Trim(config.BinaryPathName, `"`)); err == nil {
						isExpired := false
						if strings.EqualFold(binaryPathName, exeFilePath) {
							isExpired = true
							var serviceFileMd5 string
							if serviceFileMd5, err = hash.GetFileMD5(binaryPathName); err == nil {
								var exeFileMd5 string
								if exeFileMd5, err = hash.GetFileMD5(exeFilePath); err == nil {
									if strings.EqualFold(exeFileMd5, serviceFileMd5) {
										isExpired = false
									}
								}
							}
						} else {
							isExpired = true
						}
						if isExpired {
							if err = StopService(serviceName); err == nil {
								if err = svc.Delete(); err == nil {
									err = registerServiceInternal(mgr, serviceName, displayName, description, exeFilePath)
								}
							}
						}
					}
				}
				svc.Close()
			} else {
				err = registerServiceInternal(mgr, serviceName, displayName, description, exeFilePath)
			}
		}
	}
	return err
}

func Run(serviceName string, start SERVICE_START_HANDLER, stop SERVICE_STOP_HANDLER) error {
	err := error(nil)
	isService := false
	if isService, err = IsWindowsService(); err == nil && isService {
		err = __svc.Run(serviceName, &SERVICE_HANDLER{
			OnStart: start,
			OnStop:  stop,
		})
	}
	return err
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func StartService(serviceName string) error {
	err := error(nil)
	var mgr *__mgr.Mgr
	if mgr, err = __mgr.Connect(); err == nil {
		defer mgr.Disconnect()
		var svc *__mgr.Service
		if svc, err = mgr.OpenService(serviceName); err == nil {
			var config __mgr.Config
			if config, err = svc.Config(); err == nil {
				if config.StartType == __mgr.StartDisabled {
					err = fmt.Errorf("service is disabled")
				} else {
					err = svc.Start()
				}
			}
			svc.Close()
		}
	}
	return err
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func StopService(serviceName string) error {
	err := error(nil)
	var mgr *__mgr.Mgr
	if mgr, err = __mgr.Connect(); err == nil {
		defer mgr.Disconnect()
		var svc *__mgr.Service
		if svc, err = mgr.OpenService(serviceName); err == nil {
			var status __svc.Status
			if status, err = svc.Query(); err == nil {
				if status.State == __svc.Running {
					if _, err = svc.Control(__svc.Stop); err == nil {
						waitTimeout := time.Now().Add(SERVICE_STOP_TIMEOUT)
						stopped := false
						for time.Now().Before(waitTimeout) {
							if status, err = svc.Query(); err != nil {
								break
							}
							if status.State == __svc.Stopped {
								stopped = true
								break
							}
							time.Sleep(SERVICE_STOP_POLL_INTERVAL)
						}
						if err == nil && !stopped {
							err = fmt.Errorf("service did not stop within %v", SERVICE_STOP_TIMEOUT)
						}
					}
				}
			}
			svc.Close()
		}
	}
	return err
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func UnregisterService(serviceName string) bool {
	result := false
	err := error(nil)
	var mgr *__mgr.Mgr
	if mgr, err = __mgr.Connect(); err == nil {
		defer mgr.Disconnect()
		var svc *__mgr.Service
		if svc, err = mgr.OpenService(serviceName); err == nil {
			if err = svc.Delete(); err == nil {
				result = true
			}
			svc.Close()
		} else {
			result = true
		}
	}
	return result
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_SERVICE_WINDOWS, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnhandledErrorResult
func registerServiceInternal(mgr *__mgr.Mgr, serviceName string, displayName string, description string, exeFilePath string) error {
	err := error(nil)
	if mgr == nil {
		err = fmt.Errorf("mgr is nil")
	} else {
		var svc *__mgr.Service
		if svc, err = mgr.CreateService(
			serviceName,
			exeFilePath,
			__mgr.Config{
				DisplayName: displayName,
				Description: description,
				StartType:   __mgr.StartAutomatic,
			}); err == nil {
			svc.Close()
		}
	}
	return err
}
