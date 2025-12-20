//go:build windows

// Package http2
// File:        http2_windows.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/http2/http2_windows.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Windows-specific helpers for HTTP package, including service process detection.
// --------------------------------------------------------------------------------

package http2

import (
	"net/url"
	"strings"

	"golang.org/x/sys/windows/registry"
	"golang.org/x/sys/windows/svc"
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	MODULE_NAME_HTTP2_WINDOWS = "http2_windows"
)

func isWindowsService() bool {
	result := false
	if isService, err := svc.IsWindowsService(); err == nil {
		result = isService
	}
	return result
}

func loadSystemProxy() (*url.URL, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.QUERY_VALUE)
	if err != nil {
		return nil, err
	}
	defer key.Close()

	proxyEnable, _, err := key.GetIntegerValue("ProxyEnable")
	if err != nil || proxyEnable == 0 {
		return nil, nil
	}

	proxyServer, _, err := key.GetStringValue("ProxyServer")
	if err != nil {
		return nil, nil
	}
	proxyServer = strings.TrimSpace(proxyServer)
	if proxyServer == "" {
		return nil, nil
	}

	return parseWindowsProxyString(proxyServer)
}

func parseWindowsProxyString(proxy string) (*url.URL, error) {
	result := proxy
	if strings.Contains(proxy, "=") {
		result = ""
		parts := strings.Split(proxy, ";")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			lower := strings.ToLower(part)
			if strings.HasPrefix(lower, "https=") {
				result = strings.TrimSpace(part[6:])
				break
			}
			if result == "" && strings.HasPrefix(lower, "http=") {
				result = strings.TrimSpace(part[5:])
			}
		}
	}
	if result == "" {
		return nil, nil
	}
	if !strings.Contains(result, "://") {
		result = "http://" + result
	}
	return url.Parse(result)
}
