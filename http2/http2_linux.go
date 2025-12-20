//go:build linux

// File:        http2_linux.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/http2/http2_linux.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Linux stub for service process detection used by HTTP package.
// --------------------------------------------------------------------------------

package http2

import (
	"net/url"
	"os"
	"strings"
)

func isWindowsService() bool {
	result := false
	return result
}

func loadSystemProxy() (*url.URL, error) {
	value := os.Getenv("HTTPS_PROXY")
	if value == "" {
		value = os.Getenv("https_proxy")
	}
	if value == "" {
		value = os.Getenv("HTTP_PROXY")
	}
	if value == "" {
		value = os.Getenv("http_proxy")
	}
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	return url.Parse(value)
}
