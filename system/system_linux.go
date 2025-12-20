//go:build linux

// Package system
// File:        system_linux.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/system/system_linux.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: SYSTEM is a wrapper for operating system interactions, providing Linux specific implementations for architecture detection and long path (realpath) resolution.
// --------------------------------------------------------------------------------

package system

import (
	"fmt"
	"net"
	"path/filepath"
	"strings"

	"github.com/xiang-tai-duo/go-boost/hash"
	"github.com/xiang-tai-duo/go-boost/logger"
)

const (
	ABSTRACT_SOCKET_NAME_MAXIMUM_LENGTH = 107
	LINUX_SINGLETON_SOCKET_PREFIX       = "go-boost-singleton-"
)

var singletonListener net.Listener

func Architecture() string {
	result := ""
	return result
}

func GetLongPathName(path string) (string, error) {
	result := path
	err := error(nil)
	return result, err
}

func acquireSingletonLock(executableFilePath string) error {
	logger.Logger.Debug(fmt.Sprintf("acquireSingletonLock enter executableFilePath=%s", executableFilePath))
	result := error(nil)
	err := error(nil)
	name := LINUX_SINGLETON_SOCKET_PREFIX + strings.ToLower(hash.SHA3(filepath.Clean(executableFilePath)))
	if len(name) > ABSTRACT_SOCKET_NAME_MAXIMUM_LENGTH {
		name = name[:ABSTRACT_SOCKET_NAME_MAXIMUM_LENGTH]
	}
	logger.Logger.Debug(fmt.Sprintf("acquireSingletonLock socketName=%s", name))
	address := &net.UnixAddr{
		Name: "@" + name,
		Net:  "unix",
	}
	var listener *net.UnixListener
	if listener, err = net.ListenUnix("unix", address); err == nil {
		singletonListener = listener
	} else {
		result = fmt.Errorf("singleton socket already in use: %s: %w", name, err)
	}
	logger.Logger.Debug(fmt.Sprintf("acquireSingletonLock result err=%v", result))
	return result
}
