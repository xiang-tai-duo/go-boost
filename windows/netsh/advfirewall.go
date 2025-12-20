//go:build windows

// Package netsh
// File:        advfirewall.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/windows/netsh/advfirewall.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Windows netsh advfirewall wrapper
// --------------------------------------------------------------------------------

package netsh

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	ACTION_BLOCK            = "Block"
	ACTION_PREFIX           = "Action:"
	LOCAL_PORT_ANY          = "Any"
	LOCAL_PORT_PREFIX       = "LocalPort:"
	MODULE_NAME_ADVFIREWALL = "windows.netsh.advfirewall"
	NETSH_ACTION_ADD        = "add"
	NETSH_ADVFIREWALL       = "advfirewall"
	NETSH_BIN               = "netsh"
	NETSH_DELETE            = "delete"
	NETSH_FIREWALL          = "firewall"
	NETSH_RULE              = "rule"
	NETSH_SHOW              = "show"
	NETSH_VERBOSE           = "verbose"
	PROGRAM_PREFIX          = "Program:"
	PROTOCOL_TCP            = "tcp"
	PROTOCOL_TCP_UDP        = "tcp,udp"
	PROTOCOL_UDP            = "udp"
	RULE_NAME_PREFIX        = "Rule Name:"
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_ADVFIREWALL, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_ADVFIREWALL, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_ADVFIREWALL, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func AddFirewallRule(ruleName string, exeFilePath string, port int) error {
	return AddFirewallRuleEx(ruleName, exeFilePath, port, PROTOCOL_TCP)
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func AddFirewallRuleEx(ruleName string, exeFilePath string, port int, protocol string) error {
	err := error(nil)
	cmd := exec.Command(NETSH_BIN, NETSH_ADVFIREWALL, NETSH_FIREWALL, NETSH_ACTION_ADD, NETSH_RULE,
		fmt.Sprintf("name=%s", ruleName),
		"dir=in",
		"action=allow",
		"enable=yes",
		"profile=any",
		fmt.Sprintf("description=%s", ruleName),
		fmt.Sprintf("program=%s", exeFilePath),
		fmt.Sprintf("protocol=%s", protocol),
		fmt.Sprintf("localport=%d", port),
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output := make([]byte, 0)
	if output, err = cmd.CombinedOutput(); err != nil {
		err = fmt.Errorf("failed to add firewall rule: %v\nOutput: %s", err, output)
	}
	return err
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func DeleteFirewallRule(ruleName string) error {
	err := error(nil)
	cmd := exec.Command(NETSH_BIN, NETSH_ADVFIREWALL, NETSH_FIREWALL, NETSH_DELETE, NETSH_RULE,
		fmt.Sprintf("name=%s", ruleName),
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output := make([]byte, 0)
	if output, err = cmd.CombinedOutput(); err != nil {
		err = fmt.Errorf("failed to delete firewall rule: %v\nOutput: %s", err, output)
	}
	return err
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func GetFirewallBlockedRules(exeFilePath string, port int) map[string]bool {
	result := make(map[string]bool)
	cmd := exec.Command(NETSH_BIN, NETSH_ADVFIREWALL, NETSH_FIREWALL, NETSH_SHOW, NETSH_RULE, "name=all", "dir=in", NETSH_VERBOSE)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output := make([]byte, 0)
	output, _ = cmd.CombinedOutput()
	rule := ""
	action := ""
	program := ""
	localPort := ""
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, RULE_NAME_PREFIX) {
			rule = strings.TrimSpace(strings.TrimPrefix(line, RULE_NAME_PREFIX))
		} else if strings.HasPrefix(line, ACTION_PREFIX) {
			action = strings.TrimSpace(strings.TrimPrefix(line, ACTION_PREFIX))
		} else if strings.HasPrefix(line, PROGRAM_PREFIX) {
			program = strings.TrimSpace(strings.TrimPrefix(line, PROGRAM_PREFIX))
		} else if strings.HasPrefix(line, LOCAL_PORT_PREFIX) {
			localPort = strings.TrimSpace(strings.TrimPrefix(line, LOCAL_PORT_PREFIX))
		} else if line == "" && rule != "" {
			blockedPort := false
			if localPort != LOCAL_PORT_ANY && localPort != "" {
				ports := strings.Split(localPort, ",")
				for _, p := range ports {
					if strings.TrimSpace(p) == strconv.Itoa(port) {
						blockedPort = true
						break
					}
				}
			}
			if program != "" {
				blockedExe := strings.EqualFold(program, exeFilePath)
				if (blockedPort || blockedExe) && action == ACTION_BLOCK {
					result[rule] = true
				}
			}
			rule, action, program, localPort = "", "", "", ""
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func GetFirewallRuleLocalPort(ruleName string) (string, error) {
	result := ""
	err := error(nil)
	cmd := exec.Command(NETSH_BIN, NETSH_ADVFIREWALL, NETSH_FIREWALL, NETSH_SHOW, NETSH_RULE, fmt.Sprintf("name=%s", ruleName), NETSH_VERBOSE)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output := make([]byte, 0)
	if output, err = cmd.CombinedOutput(); err == nil {
		for _, line := range strings.Split(string(output), "\n") {
			if strings.Contains(line, LOCAL_PORT_PREFIX) {
				parts := strings.Split(line, LOCAL_PORT_PREFIX)
				if len(parts) > 1 {
					result = strings.TrimSpace(parts[1])
					break
				}
			}
		}
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func IsFirewallRuleExists(exeFilePath string, port int) (bool, error) {
	result := false
	err := error(nil)
	cmd := exec.Command(NETSH_BIN, NETSH_ADVFIREWALL, NETSH_FIREWALL, NETSH_SHOW, NETSH_RULE, "name=all", NETSH_VERBOSE)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output := make([]byte, 0)
	if output, err = cmd.CombinedOutput(); err == nil {
		rule := ""
		program := ""
		localPort := ""
		for _, line := range strings.Split(string(output), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, RULE_NAME_PREFIX) {
				rule = strings.TrimSpace(strings.TrimPrefix(line, RULE_NAME_PREFIX))
			} else if strings.HasPrefix(line, PROGRAM_PREFIX) {
				program = strings.TrimSpace(strings.TrimPrefix(line, PROGRAM_PREFIX))
			} else if strings.HasPrefix(line, LOCAL_PORT_PREFIX) {
				localPort = strings.TrimSpace(strings.TrimPrefix(line, LOCAL_PORT_PREFIX))
			} else if line == "" && rule != "" {
				if strings.EqualFold(program, exeFilePath) && localPort == strconv.Itoa(port) {
					result = true
					break
				}
				rule, program, localPort = "", "", ""
			}
		}
	} else {
		err = fmt.Errorf("failed to list firewall rules: %v", err)
	}
	return result, err
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_ADVFIREWALL, logger.SKIP_STACK_FRAMES_BASE)
}
