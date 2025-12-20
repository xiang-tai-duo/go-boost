//go:build darwin

// File:        http2_darwin.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/http2/http2_darwin.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Darwin (macOS) stub for service process detection used by HTTP package.
// --------------------------------------------------------------------------------

package http2

import (
	"bufio"
	"fmt"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
)

func isWindowsService() bool {
	result := false
	return result
}

func loadSystemProxy() (*url.URL, error) {
	output, err := exec.Command("scutil", "--proxy").Output()
	if err != nil {
		return nil, err
	}
	httpHost, httpPort := parseScutilProxy(string(output), "HTTPEnable", "HTTPProxy", "HTTPPort")
	httpsHost, httpsPort := parseScutilProxy(string(output), "HTTPSEnable", "HTTPSProxy", "HTTPSPort")
	if httpsHost != "" {
		return buildScutilProxyURL("https", httpsHost, httpsPort), nil
	}
	if httpHost != "" {
		return buildScutilProxyURL("http", httpHost, httpPort), nil
	}
	return nil, nil
}

func parseScutilProxy(output, enableKey, hostKey, portKey string) (string, int) {
	enabled := false
	host := ""
	port := 0
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case enableKey:
			if n, parseErr := strconv.Atoi(value); parseErr == nil && n != 0 {
				enabled = true
			}
		case hostKey:
			host = value
		case portKey:
			if n, parseErr := strconv.Atoi(value); parseErr == nil {
				port = n
			}
		}
	}
	if !enabled || host == "" {
		return "", 0
	}
	return host, port
}

func buildScutilProxyURL(scheme, host string, port int) *url.URL {
	if port > 0 {
		return &url.URL{
			Scheme: scheme,
			Host:   fmt.Sprintf("%s:%d", host, port),
		}
	}
	return &url.URL{
		Scheme: scheme,
		Host:   host,
	}
}
