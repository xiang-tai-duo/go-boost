// Package network
// File:        network.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/network/network.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: NETWORK provides functions to get network IP addresses with the smallest metric.
// --------------------------------------------------------------------------------
package network

import (
	"errors"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/xiang-tai-duo/go-boost/logger"
	"github.com/xiang-tai-duo/go-boost/osutil"
)

//goland:noinspection GoSnakeCaseUsage
type (
	INTERNET_PROTOCOL_ADDRESSES struct {
		InternetProtocolVersion4 string
		InternetProtocolVersion6 string
	}
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst,GoNameStartsWithPackageName,SpellCheckingInspection
const (
	ALL_NUMERIC_OWNER_FLAG                                  = "-ano"
	DARWIN_OPERATING_SYSTEM                                 = "darwin"
	DARWIN_ROUTE_DESTINATION_FIELD_INDEX                    = 0
	DARWIN_ROUTE_INTERFACE_FIELD_INDEX                      = 3
	DARWIN_ROUTE_MINIMUM_FIELD_COUNT                        = 5
	DECIMAL_FORMAT                                          = "%d"
	DEFAULT_ROUTE_ADDRESS                                   = "0.0.0.0"
	DEFAULT_ROUTE_KEYWORD                                   = "default"
	DESTINATION_HEADER                                      = "Destination"
	DEVICE_FIELD_KEYWORD                                    = "dev"
	EPHEMERAL_PORT_ADDRESS                                  = ":0"
	INITIAL_CONNECTION_TIMEOUT                              = 100 * time.Millisecond
	INTERNET_ADDRESSES_FLAG                                 = "-i"
	INTERNET_HEADER                                         = "Internet:"
	INTERNET_PROTOCOL_COMMAND                               = "ip"
	LINUX_OPERATING_SYSTEM                                  = "linux"
	LIST_OPEN_FILES_COMMAND                                 = "lsof"
	LSOF_MINIMUM_FIELD_COUNT                                = 2
	LSOF_PROCESS_IDENTIFIER_FIELD_INDEX                     = 1
	LOCALHOST_ADDRESS_FORMAT                                = "127.0.0.1:%d"
	MAXIMUM_CONNECTION_TIMEOUT                              = 3 * time.Second
	MAX_PORT                                                = 65535
	METRIC_FIELD_KEYWORD                                    = "metric"
	MINIMUM_PORT_NUMBER                                     = 1024
	MIN_PORT                                                = 1
	MODULE_NAME_NETWORK                                     = "network"
	NETSTAT_MINIMUM_FIELD_COUNT                             = 5
	NETWORK_DESTINATION_HEADER                              = "Network Destination"
	NETWORK_STATISTICS_COMMAND                              = "netstat"
	NEWLINE_CHARACTER                                       = "\n"
	NEXT_FIELD_OFFSET                                       = 1
	PORT_SUFFIX_FORMAT                                      = ":%d"
	PROTOCOL_HEADER                                         = "Proto"
	REFUSED_ERROR_KEYWORD                                   = "refused"
	ROUTE_NUMERIC_FLAG                                      = "-rn"
	ROUTE_SUBCOMMAND                                        = "route"
	ROUTING_TABLES_HEADER                                   = "Routing tables"
	SEPARATOR_LINE                                          = "==========="
	TIMEOUT_ERROR_KEYWORD                                   = "timeout"
	TRANSMISSION_CONTROL_PROTOCOL                           = "tcp"
	WINDOWS_NETSTAT_INTERFACE_ADDRESS_FIELD_INDEX           = 1
	WINDOWS_NETSTAT_INTERFACE_INTERNET_PROTOCOL_FIELD_INDEX = 3
	WINDOWS_NETSTAT_METRIC_FIELD_INDEX                      = 4
	WINDOWS_NETSTAT_PROCESS_IDENTIFIER_FIELD_INDEX          = 4
	WINDOWS_OPERATING_SYSTEM                                = "windows"
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_NETWORK, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_NETWORK, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_NETWORK, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection SpellCheckingInspection,GoUnusedExportedFunction
func GetNetworkIpAddresses() (INTERNET_PROTOCOL_ADDRESSES, error) {
	result := INTERNET_PROTOCOL_ADDRESSES{}
	err := error(nil)
	operatingSystem := runtime.GOOS
	switch operatingSystem {
	case WINDOWS_OPERATING_SYSTEM:
		result, err = getWindowsIpAddressesWithSmallestMetric()
	case LINUX_OPERATING_SYSTEM:
		result, err = getLinuxIpAddressesWithSmallestMetric()
	case DARWIN_OPERATING_SYSTEM:
		result, err = getDarwinIpAddressesWithSmallestMetric()
	default:
		err = fmt.Errorf("unsupported OS: %s", operatingSystem)
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func GetProcessIdByPort(port int) (*int, error) {
	result := (*int)(nil)
	err := error(nil)
	operatingSystem := runtime.GOOS
	switch operatingSystem {
	case WINDOWS_OPERATING_SYSTEM:
		result, err = getWindowsProcessIdByPort(port)
	case LINUX_OPERATING_SYSTEM:
		result, err = getLinuxProcessIdByPort(port)
	case DARWIN_OPERATING_SYSTEM:
		result, err = getDarwinProcessIdByPort(port)
	default:
		err = fmt.Errorf("unsupported OS: %s", operatingSystem)
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func GetRandomPort() int {
	result := 0
	err := error(nil)
	for {
		var listener net.Listener
		if listener, err = net.Listen(TRANSMISSION_CONTROL_PROTOCOL, EPHEMERAL_PORT_ADDRESS); err == nil {
			result = listener.Addr().(*net.TCPAddr).Port
			if err = listener.Close(); err == nil {
				if result > MINIMUM_PORT_NUMBER {
					break
				}
			}
		} else {
			result = 0
			break
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func IsInvalidPortNumber(port int) bool {
	return port < MIN_PORT || port > MAX_PORT
}

//goland:noinspection GoUnusedExportedFunction
func IsPortAvailable(port int) bool {
	result := false
	err := error(nil)
	address := fmt.Sprintf(LOCALHOST_ADDRESS_FORMAT, port)
	var conn net.Conn
	timeout := INITIAL_CONNECTION_TIMEOUT
	for {
		if conn, err = net.DialTimeout(TRANSMISSION_CONTROL_PROTOCOL, address, timeout); err == nil {
			if err = conn.Close(); err == nil {
				result = false
			}
			__debug(fmt.Sprintf("%s: %v (dial succeeded)", address, result))
			break
		} else if strings.Contains(err.Error(), REFUSED_ERROR_KEYWORD) {
			result = true
			__debug(fmt.Sprintf("%s: %v (connection refused)", address, result))
			break
		} else if strings.Contains(err.Error(), TIMEOUT_ERROR_KEYWORD) {
			if timeout >= MAXIMUM_CONNECTION_TIMEOUT {
				result = true
				__debug(fmt.Sprintf("Detection %s: %v (timeout exceeded)", address, result))
				break
			}
			timeout += INITIAL_CONNECTION_TIMEOUT
			continue
		} else {
			__debug(fmt.Sprintf("%s: %v (unknown error: %v)", address, result, err))
			break
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func IsValidPortNumber(port int) bool {
	return port >= MIN_PORT && port <= MAX_PORT
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_NETWORK, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection SpellCheckingInspection
func getDarwinIpAddressesWithSmallestMetric() (INTERNET_PROTOCOL_ADDRESSES, error) {
	result := INTERNET_PROTOCOL_ADDRESSES{}
	err := error(nil)
	bestInterfaceName := ""
	found := false
	command := exec.Command(NETWORK_STATISTICS_COMMAND, ROUTE_NUMERIC_FLAG)
	var output []byte
	if output, err = command.Output(); err == nil {
		lines := strings.Split(string(output), NEWLINE_CHARACTER)
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if strings.HasPrefix(line, ROUTING_TABLES_HEADER) || strings.HasPrefix(line, INTERNET_HEADER) || strings.HasPrefix(line, DESTINATION_HEADER) {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) >= DARWIN_ROUTE_MINIMUM_FIELD_COUNT && fields[DARWIN_ROUTE_DESTINATION_FIELD_INDEX] == DEFAULT_ROUTE_KEYWORD {
				bestInterfaceName = fields[DARWIN_ROUTE_INTERFACE_FIELD_INDEX]
				found = true
			}
		}
	}
	if found && bestInterfaceName != "" {
		var networkInterface *net.Interface
		if networkInterface, err = net.InterfaceByName(bestInterfaceName); err == nil {
			result = getIpAddressesFromInterface(networkInterface)
		}
	} else {
		result = getIpAddressesFromAllInterfaces()
	}
	return result, err
}

//goland:noinspection DuplicatedCode
func getDarwinProcessIdByPort(port int) (*int, error) {
	result := (*int)(nil)
	err := error(nil)
	command := exec.Command(LIST_OPEN_FILES_COMMAND, INTERNET_ADDRESSES_FLAG, fmt.Sprintf(PORT_SUFFIX_FORMAT, port))
	var output []byte
	if output, err = command.Output(); err == nil {
		lines := strings.Split(string(output), NEWLINE_CHARACTER)
		for index, line := range lines {
			if index == 0 {
				continue
			}
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) >= LSOF_MINIMUM_FIELD_COUNT {
				processIdentifierString := fields[LSOF_PROCESS_IDENTIFIER_FIELD_INDEX]
				var processIdentifier int
				var parsingError error
				if processIdentifier, parsingError = strconv.Atoi(processIdentifierString); parsingError == nil {
					result = &processIdentifier
					break
				}
			}
		}
	} else {
		if exitError, ok := errors.AsType[*exec.ExitError](err); ok && exitError.ExitCode() == 1 {
			result = nil
			err = nil
		}
	}
	return result, err
}

func getIpAddressesFromAllInterfaces() INTERNET_PROTOCOL_ADDRESSES {
	result := INTERNET_PROTOCOL_ADDRESSES{}
	interfaces, interfaceError := net.Interfaces()
	if interfaceError == nil {
		for _, networkInterface := range interfaces {
			if networkInterface.Flags&net.FlagUp != 0 {
				addresses, addressError := networkInterface.Addrs()
				if addressError == nil {
					for _, address := range addresses {
						internetProtocolNetwork, ok := address.(*net.IPNet)
						if ok {
							internetProtocol := internetProtocolNetwork.IP
							if internetProtocol.To4() != nil && !internetProtocol.IsLoopback() && result.InternetProtocolVersion4 == "" {
								result.InternetProtocolVersion4 = internetProtocol.String()
							} else if internetProtocol.To16() != nil && !internetProtocol.IsLoopback() && result.InternetProtocolVersion6 == "" {
								result.InternetProtocolVersion6 = internetProtocol.String()
							}
						}
					}
				}
			}
		}
	}
	return result
}

func getIpAddressesFromInterface(networkInterface *net.Interface) INTERNET_PROTOCOL_ADDRESSES {
	result := INTERNET_PROTOCOL_ADDRESSES{}
	addresses, addressError := networkInterface.Addrs()
	if addressError == nil {
		for _, address := range addresses {
			internetProtocolNetwork, ok := address.(*net.IPNet)
			if ok {
				internetProtocol := internetProtocolNetwork.IP
				if internetProtocol.To4() != nil && !internetProtocol.IsLoopback() {
					result.InternetProtocolVersion4 = internetProtocol.String()
				} else if internetProtocol.To16() != nil && !internetProtocol.IsLoopback() {
					result.InternetProtocolVersion6 = internetProtocol.String()
				}
			}
		}
	}
	return result
}

//goland:noinspection SpellCheckingInspection,GoUnhandledErrorResult
func getLinuxIpAddressesWithSmallestMetric() (INTERNET_PROTOCOL_ADDRESSES, error) {
	result := INTERNET_PROTOCOL_ADDRESSES{}
	err := error(nil)
	bestInterfaceName := ""
	bestMetric := 0
	found := false
	command := exec.Command(INTERNET_PROTOCOL_COMMAND, ROUTE_SUBCOMMAND)
	var output []byte
	if output, err = command.Output(); err == nil {
		lines := strings.Split(string(output), NEWLINE_CHARACTER)
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if strings.HasPrefix(line, DEFAULT_ROUTE_KEYWORD) {
				fields := strings.Fields(line)
				networkInterfaceName := ""
				metric := 0
				for index, field := range fields {
					if field == DEVICE_FIELD_KEYWORD && index+NEXT_FIELD_OFFSET < len(fields) {
						networkInterfaceName = fields[index+NEXT_FIELD_OFFSET]
					} else if field == METRIC_FIELD_KEYWORD && index+NEXT_FIELD_OFFSET < len(fields) {
						fmt.Sscanf(fields[index+NEXT_FIELD_OFFSET], DECIMAL_FORMAT, &metric)
					}
				}
				if networkInterfaceName != "" {
					if !found || metric < bestMetric {
						bestMetric = metric
						bestInterfaceName = networkInterfaceName
						found = true
					}
				}
			}
		}
	}
	if found && bestInterfaceName != "" {
		var networkInterface *net.Interface
		if networkInterface, err = net.InterfaceByName(bestInterfaceName); err == nil {
			result = getIpAddressesFromInterface(networkInterface)
		}
	} else {
		result = getIpAddressesFromAllInterfaces()
	}
	return result, err
}

//goland:noinspection DuplicatedCode
func getLinuxProcessIdByPort(port int) (*int, error) {
	result := (*int)(nil)
	err := error(nil)
	command := exec.Command(LIST_OPEN_FILES_COMMAND, INTERNET_ADDRESSES_FLAG, fmt.Sprintf(PORT_SUFFIX_FORMAT, port))
	var output []byte
	if output, err = command.Output(); err == nil {
		lines := strings.Split(string(output), NEWLINE_CHARACTER)
		for index, line := range lines {
			if index == 0 {
				continue
			}
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) >= LSOF_MINIMUM_FIELD_COUNT {
				processIdentifierString := fields[LSOF_PROCESS_IDENTIFIER_FIELD_INDEX]
				var processIdentifier int
				var parsingError error
				if processIdentifier, parsingError = strconv.Atoi(processIdentifierString); parsingError == nil {
					result = &processIdentifier
					break
				}
			}
		}
	} else {
		if exitError, ok := errors.AsType[*exec.ExitError](err); ok && exitError.ExitCode() == 1 {
			result = nil
			err = nil
		}
	}
	return result, err
}

//goland:noinspection GoUnhandledErrorResult,SpellCheckingInspection
func getWindowsIpAddressesWithSmallestMetric() (INTERNET_PROTOCOL_ADDRESSES, error) {
	result := INTERNET_PROTOCOL_ADDRESSES{}
	err := error(nil)
	bestInterfaceIndex := 0
	bestMetric := 0
	found := false
	command := exec.Command(NETWORK_STATISTICS_COMMAND, ROUTE_NUMERIC_FLAG)
	osutil.SetHideWindow(command)
	var output []byte
	if output, err = command.Output(); err == nil {
		lines := strings.Split(string(output), NEWLINE_CHARACTER)
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if strings.HasPrefix(line, NETWORK_DESTINATION_HEADER) || strings.HasPrefix(line, SEPARATOR_LINE) {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) >= NETSTAT_MINIMUM_FIELD_COUNT && fields[0] == DEFAULT_ROUTE_ADDRESS {
				interfaceInternetProtocol := fields[WINDOWS_NETSTAT_INTERFACE_INTERNET_PROTOCOL_FIELD_INDEX]
				metric := 0
				if len(fields) > WINDOWS_NETSTAT_METRIC_FIELD_INDEX {
					fmt.Sscanf(fields[WINDOWS_NETSTAT_METRIC_FIELD_INDEX], DECIMAL_FORMAT, &metric)
				}
				if !found || metric < bestMetric {
					bestMetric = metric
					found = true
					interfaces, interfaceError := net.Interfaces()
					if interfaceError == nil {
						for _, networkInterface := range interfaces {
							addresses, addressError := networkInterface.Addrs()
							if addressError == nil {
								for _, address := range addresses {
									internetProtocolNetwork, ok := address.(*net.IPNet)
									if ok {
										internetProtocol := internetProtocolNetwork.IP
										if internetProtocol.String() == interfaceInternetProtocol {
											bestInterfaceIndex = networkInterface.Index
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	if found && bestInterfaceIndex != 0 {
		var networkInterface *net.Interface
		if networkInterface, err = net.InterfaceByIndex(bestInterfaceIndex); err == nil {
			result = getIpAddressesFromInterface(networkInterface)
		}
	} else {
		result = getIpAddressesFromAllInterfaces()
	}
	return result, err
}

func getWindowsProcessIdByPort(port int) (*int, error) {
	result := (*int)(nil)
	err := error(nil)
	command := exec.Command(NETWORK_STATISTICS_COMMAND, ALL_NUMERIC_OWNER_FLAG)
	osutil.SetHideWindow(command)
	var output []byte
	if output, err = command.Output(); err == nil {
		lines := strings.Split(string(output), NEWLINE_CHARACTER)
		portString := fmt.Sprintf(PORT_SUFFIX_FORMAT, port)
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if strings.HasPrefix(line, PROTOCOL_HEADER) {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) >= NETSTAT_MINIMUM_FIELD_COUNT {
				localAddress := fields[WINDOWS_NETSTAT_INTERFACE_ADDRESS_FIELD_INDEX]
				if strings.HasSuffix(localAddress, portString) {
					processIdentifierString := fields[WINDOWS_NETSTAT_PROCESS_IDENTIFIER_FIELD_INDEX]
					var processIdentifier int
					var parsingError error
					if processIdentifier, parsingError = strconv.Atoi(processIdentifierString); parsingError == nil {
						result = &processIdentifier
						break
					}
				}
			}
		}
	}
	return result, err
}
