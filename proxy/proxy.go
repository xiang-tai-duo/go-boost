// Package proxy
// File:        proxy.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/proxy/proxy.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: PROXY is an HTTP/HTTPS forward proxy with rule-based routing support.
// --------------------------------------------------------------------------------
package proxy

import (
	"embed"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/elazarl/goproxy"
	"github.com/xiang-tai-duo/go-boost/logger"
)

type (
	PROXY_CONFIG struct {
		DefaultProxy string       `json:"default_proxy"`
		Listen       string       `json:"listen"`
		Rules        []PROXY_RULE `json:"rules"`
	}

	PROXY_MATCHER struct {
		rules []PROXY_RULE
	}

	PROXY_RULE struct {
		Direct           bool   `json:"direct,omitempty"`
		Domain           string `json:"domain,omitempty"`
		InternetProtocol string `json:"ip,omitempty"`
		Proxy            string `json:"proxy,omitempty"`
	}

	PROXY_SERVER struct {
		config  *PROXY_CONFIG
		matcher *PROXY_MATCHER
		proxy   *goproxy.ProxyHttpServer
	}
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	CONNECT_REQUEST_FORMAT            = "CONNECT %s HTTP/1.1\r\nHost: %s\r\n\r\n"
	DEFAULT_CONFIGURATION_FILENAME    = "config.json"
	DOMAIN_WILDCARD_PREFIX            = "*."
	EMBEDDED_CONFIGURATION_PATH       = "config.json"
	HTTP_VERSION_OKAY_PREFIX          = "200"
	INTERNET_PROTOCOL_WILDCARD_SUFFIX = "*"
	MODULE_NAME_PROXY                 = "proxy"
	READ_BUFFER_SIZE                  = 1024
	RESPONSE_MINIMUM_LENGTH           = 12
	RESPONSE_STATUS_START             = 9
	TRANSMISSION_CONTROL_PROTOCOL     = "tcp"
)

//go:embed config.json
var CONFIG_JSON_FILE embed.FS

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_PROXY, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_PROXY, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_PROXY, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedExportedFunction
func New() (*PROXY_SERVER, error) {
	server := (*PROXY_SERVER)(nil)
	configuration := (*PROXY_CONFIG)(nil)
	err := error(nil)

	__debug("New: called")
	if configuration, err = loadConfigurationFromExecutableDirectory(); err == nil {
		server = newServer(configuration)
		__debug("New: server created successfully")
	}

	return server, err
}

func (matcher *PROXY_MATCHER) Match(host string) (string, bool) {
	proxyUniformResourceLocator := ""
	direct := false
	found := false

	__debug(fmt.Sprintf("PROXY_MATCHER.Match: host=%s, rules=%d", host, len(matcher.rules)))
	for _, rule := range matcher.rules {
		if rule.Domain != "" && matchDomain(host, rule.Domain) {
			proxyUniformResourceLocator = rule.Proxy
			direct = rule.Direct
			found = true
			__debug(fmt.Sprintf("PROXY_MATCHER.Match: domain matched, host=%s, pattern=%s, proxy=%s, direct=%v", host, rule.Domain, rule.Proxy, rule.Direct))
		} else {
			if rule.InternetProtocol != "" && matchInternetProtocol(host, rule.InternetProtocol) {
				proxyUniformResourceLocator = rule.Proxy
				direct = rule.Direct
				found = true
				__debug(fmt.Sprintf("PROXY_MATCHER.Match: internet protocol matched, host=%s, pattern=%s, proxy=%s, direct=%v", host, rule.InternetProtocol, rule.Proxy, rule.Direct))
			}
		}
		if found {
			break
		}
	}
	if !found {
		__debug(fmt.Sprintf("PROXY_MATCHER.Match: no rule matched for host=%s, returning default", host))
	}

	return proxyUniformResourceLocator, direct
}

func (server *PROXY_SERVER) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	__debug(fmt.Sprintf("PROXY_SERVER.ServeHTTP: method=%s, uniform_resource_locator=%s, remote=%s", request.Method, request.URL.String(), request.RemoteAddr))
	server.proxy.ServeHTTP(writer, request)
}

//goland:noinspection GoUnusedExportedFunction
func (server *PROXY_SERVER) Start() error {
	err := error(nil)
	__debug(fmt.Sprintf("Proxy server listening on %s", server.config.Listen))
	__debug("Powered by: github.com/elazarl/goproxy")
	err = http.ListenAndServe(server.config.Listen, server)
	if err != nil {
		__debug(fmt.Sprintf("PROXY_SERVER.Start: http.ListenAndServe failed: %v", err))
	}
	return err
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_PROXY, logger.SKIP_STACK_FRAMES_BASE)
}

func (server *PROXY_SERVER) createConnectDial(proxyUniformResourceLocator *url.URL) func(network, address string) (net.Conn, error) {
	return func(network, address string) (net.Conn, error) {
		err := error(nil)
		connection := (net.Conn)(nil)
		connectRequest := ""
		buffer := make([]byte, READ_BUFFER_SIZE)
		count := 0
		response := ""
		__debug(fmt.Sprintf("createConnectDial: dialing proxy %s for target %s", proxyUniformResourceLocator.Host, address))
		if connection, err = net.Dial(TRANSMISSION_CONTROL_PROTOCOL, proxyUniformResourceLocator.Host); err == nil {
			__debug(fmt.Sprintf("createConnectDial: connected to proxy %s", proxyUniformResourceLocator.Host))
			connectRequest = fmt.Sprintf(CONNECT_REQUEST_FORMAT, address, address)
			if _, err = connection.Write([]byte(connectRequest)); err == nil {
				__debug(fmt.Sprintf("createConnectDial: CONNECT request sent for %s", address))
				if count, err = connection.Read(buffer); err == nil {
					response = string(buffer[:count])
					__debug(fmt.Sprintf("createConnectDial: proxy response (%d bytes): %s", count, strings.TrimSpace(response)))
					if len(response) < RESPONSE_MINIMUM_LENGTH || response[RESPONSE_STATUS_START:RESPONSE_MINIMUM_LENGTH] != HTTP_VERSION_OKAY_PREFIX {
						connection.Close()
						connection = nil
						err = fmt.Errorf("proxy connect failed: %s", response)
						__debug(fmt.Sprintf("createConnectDial: proxy returned non-200: %s", strings.TrimSpace(response)))
					} else {
						__debug(fmt.Sprintf("createConnectDial: tunnel established to %s via %s", address, proxyUniformResourceLocator.Host))
					}
				} else {
					__debug(fmt.Sprintf("createConnectDial: connection.Read failed: %v", err))
					connection.Close()
					connection = nil
				}
			} else {
				__debug(fmt.Sprintf("createConnectDial: connection.Write failed: %v", err))
				connection.Close()
				connection = nil
			}
		} else {
			__debug(fmt.Sprintf("createConnectDial: net.Dial failed to %s: %v", proxyUniformResourceLocator.Host, err))
		}

		return connection, err
	}
}

func extractInternetProtocol(host string) net.IP {
	result := (net.IP)(nil)
	if strings.Contains(host, ":") {
		host, _, _ = net.SplitHostPort(host)
		__debug(fmt.Sprintf("extractInternetProtocol: split host:port, host=%s", host))
	}
	result = net.ParseIP(host)
	__debug(fmt.Sprintf("extractInternetProtocol: host=%s, result=%v", host, result))
	return result
}

func getExecutableDirectory() string {
	result := ""
	executablePath := ""
	executablePath, _ = os.Executable()
	if executablePath != "" {
		result = filepath.Dir(executablePath)
	} else {
		result, _ = os.Getwd()
	}
	__debug(fmt.Sprintf("getExecutableDirectory: result=%s", result))
	return result
}

func (server *PROXY_SERVER) handleConnect(host string, context *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	proxyUniformResourceLocator, direct := server.matcher.Match(host)

	__debug(fmt.Sprintf("handleConnect: host=%s, matched_proxy=%s, direct=%v", host, proxyUniformResourceLocator, direct))
	if direct {
		context.Proxy.ConnectDial = nil
		__debug(fmt.Sprintf("handleConnect: DIRECT connect to %s", host))
	} else {
		if proxyUniformResourceLocator == "" {
			proxyUniformResourceLocator = server.config.DefaultProxy
			__debug(fmt.Sprintf("handleConnect: using default proxy=%s for %s", proxyUniformResourceLocator, host))
		}

		if proxyUniformResourceLocator != "" {
			parsedUniformResourceLocator := (*url.URL)(nil)
			err := error(nil)
			if parsedUniformResourceLocator, err = url.Parse(proxyUniformResourceLocator); err == nil {
				context.Proxy.ConnectDial = server.createConnectDial(parsedUniformResourceLocator)
				__debug(fmt.Sprintf("handleConnect: PROXY connect to %s via %s", host, proxyUniformResourceLocator))
			} else {
				__debug(fmt.Sprintf("handleConnect: failed to parse proxy uniform resource locator %s: %v", proxyUniformResourceLocator, err))
			}
		} else {
			__debug(fmt.Sprintf("handleConnect: no proxy configured for %s and no default proxy, using direct", host))
		}
	}

	return goproxy.OkConnect, host
}

func (server *PROXY_SERVER) handleRequest(request *http.Request, context *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	host := request.URL.Host
	if host == "" {
		host = request.Host
	}

	proxyUniformResourceLocator, direct := server.matcher.Match(host)

	__debug(fmt.Sprintf("handleRequest: method=%s, host=%s, uniform_resource_locator=%s, matched_proxy=%s, direct=%v", request.Method, host, request.URL.String(), proxyUniformResourceLocator, direct))
	if direct {
		context.Proxy.Tr = &http.Transport{}
		__debug(fmt.Sprintf("handleRequest: DIRECT request to %s %s", request.Method, host))
	} else {
		if proxyUniformResourceLocator == "" {
			proxyUniformResourceLocator = server.config.DefaultProxy
			__debug(fmt.Sprintf("handleRequest: using default proxy=%s for %s", proxyUniformResourceLocator, host))
		}

		if proxyUniformResourceLocator != "" {
			parsedUniformResourceLocator := (*url.URL)(nil)
			err := error(nil)
			if parsedUniformResourceLocator, err = url.Parse(proxyUniformResourceLocator); err == nil {
				context.Proxy.Tr = &http.Transport{
					Proxy: http.ProxyURL(parsedUniformResourceLocator),
				}
				__debug(fmt.Sprintf("handleRequest: PROXY request %s %s via %s", request.Method, host, proxyUniformResourceLocator))
			} else {
				__debug(fmt.Sprintf("handleRequest: failed to parse proxy uniform resource locator %s: %v", proxyUniformResourceLocator, err))
			}
		} else {
			__debug(fmt.Sprintf("handleRequest: no proxy configured for %s and no default proxy, using direct", host))
		}
	}

	return request, nil
}

func loadConfiguration(path string) (configuration *PROXY_CONFIG, err error) {
	err = error(nil)
	data := make([]byte, 0)
	configuration = (*PROXY_CONFIG)(nil)

	__debug(fmt.Sprintf("loadConfiguration: path=%s", path))
	if data, err = os.ReadFile(path); err == nil {
		__debug(fmt.Sprintf("loadConfiguration: read %d bytes", len(data)))
		configuration = &PROXY_CONFIG{}
		if err = json.Unmarshal(data, configuration); err != nil {
			__debug(fmt.Sprintf("loadConfiguration: json.Unmarshal failed: %v", err))
			configuration = nil
		} else {
			__debug(fmt.Sprintf("loadConfiguration: success, listen=%s, default_proxy=%s, rules=%d", configuration.Listen, configuration.DefaultProxy, len(configuration.Rules)))
		}
	} else {
		__debug(fmt.Sprintf("loadConfiguration: os.ReadFile failed: %v", err))
	}

	return configuration, err
}

func loadConfigurationFromExecutableDirectory() (*PROXY_CONFIG, error) {
	configuration := (*PROXY_CONFIG)(nil)
	err := error(nil)
	configurationPath := ""

	__debug("loadConfigurationFromExecutableDirectory: called")
	if configurationPath, err = newConfiguration(); err == nil {
		__debug(fmt.Sprintf("loadConfigurationFromExecutableDirectory: loading configuration from %s", configurationPath))
		configuration, err = loadConfiguration(configurationPath)
	} else {
		__debug(fmt.Sprintf("loadConfigurationFromExecutableDirectory: newConfiguration failed: %v", err))
	}

	return configuration, err
}

func matchDomain(host, pattern string) bool {
	result := false
	host = strings.ToLower(host)
	pattern = strings.ToLower(pattern)

	if strings.HasPrefix(pattern, DOMAIN_WILDCARD_PREFIX) {
		suffix := pattern[1:]
		result = strings.HasSuffix(host, suffix) || host == pattern[2:]
	} else {
		result = host == pattern
	}

	__debug(fmt.Sprintf("matchDomain: host=%s, pattern=%s, result=%v", host, pattern, result))
	return result
}

func matchInternetProtocol(host, pattern string) bool {
	result := false
	hostInternetProtocol := extractInternetProtocol(host)

	if hostInternetProtocol != nil {
		if strings.HasSuffix(pattern, INTERNET_PROTOCOL_WILDCARD_SUFFIX) {
			prefix := pattern[:len(pattern)-1]
			result = strings.HasPrefix(hostInternetProtocol.String(), prefix)
		} else {
			patternInternetProtocol := net.ParseIP(pattern)
			if patternInternetProtocol != nil {
				result = hostInternetProtocol.Equal(patternInternetProtocol)
			}
		}
	}

	__debug(fmt.Sprintf("matchInternetProtocol: host=%s, pattern=%s, host_internet_protocol=%v, result=%v", host, pattern, hostInternetProtocol, result))
	return result
}

func newConfiguration() (string, error) {
	err := error(nil)
	targetPath := ""
	data := make([]byte, 0)
	directory := ""

	directory = getExecutableDirectory()
	targetPath = filepath.Join(directory, DEFAULT_CONFIGURATION_FILENAME)

	__debug(fmt.Sprintf("newConfiguration: directory=%s, target=%s", directory, targetPath))
	if _, err = os.Stat(targetPath); os.IsNotExist(err) {
		__debug(fmt.Sprintf("newConfiguration: configuration not found, releasing embedded configuration to %s", targetPath))
		if data, err = CONFIG_JSON_FILE.ReadFile(EMBEDDED_CONFIGURATION_PATH); err == nil {
			__debug(fmt.Sprintf("newConfiguration: embedded configuration read %d bytes", len(data)))
			err = os.WriteFile(targetPath, data, os.FileMode(0644))
			if err == nil {
				__debug(fmt.Sprintf("newConfiguration: configuration written to %s", targetPath))
			} else {
				__debug(fmt.Sprintf("newConfiguration: os.WriteFile failed: %v", err))
			}
		} else {
			__debug(fmt.Sprintf("newConfiguration: CONFIG_JSON_FILE.ReadFile failed: %v", err))
		}
	} else {
		err = nil
		__debug(fmt.Sprintf("newConfiguration: configuration already exists at %s", targetPath))
	}

	return targetPath, err
}

func newMatcher(rules []PROXY_RULE) *PROXY_MATCHER {
	__debug(fmt.Sprintf("newMatcher: creating matcher with %d rules", len(rules)))
	return &PROXY_MATCHER{rules: rules}
}

func newServer(configuration *PROXY_CONFIG) *PROXY_SERVER {
	__debug(fmt.Sprintf("newServer: creating server, listen=%s, default_proxy=%s, rules=%d", configuration.Listen, configuration.DefaultProxy, len(configuration.Rules)))
	matcher := newMatcher(configuration.Rules)
	proxyServer := goproxy.NewProxyHttpServer()
	proxyServer.Verbose = false

	server := &PROXY_SERVER{
		config:  configuration,
		matcher: matcher,
		proxy:   proxyServer,
	}

	proxyServer.OnRequest().DoFunc(server.handleRequest)
	proxyServer.OnRequest().HandleConnectFunc(server.handleConnect)

	__debug("newServer: server created, handlers registered")
	return server
}
