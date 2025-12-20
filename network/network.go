// Package network
// File:        network.go
// Url:         https://github.com/xiang-tai-duo/go-bootstrap/blob/master/network/network.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: NETWORK provides functions to get network IP addresses with the smallest metric.
// --------------------------------------------------------------------------------
package network

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

//goland:noinspection GoSnakeCaseUsage
type (
	IP_ADDRESSES struct {
		IPv4 string
		IPv6 string
	}
)

//goland:noinspection SpellCheckingInspection,GoUnusedExportedFunction
func GetNetworkIpAddresses() (IP_ADDRESSES, error) {
	result := IP_ADDRESSES{}
	err := error(nil)
	operatingSystem := runtime.GOOS
	switch operatingSystem {
	case "windows":
		result, err = getWindowsIpAddressesWithSmallestMetric()
	case "linux":
		result, err = getLinuxIpAddressesWithSmallestMetric()
	case "darwin":
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
	case "windows":
		result, err = getWindowsProcessIdByPort(port)
	case "linux":
		result, err = getLinuxProcessIdByPort(port)
	case "darwin":
		result, err = getDarwinProcessIdByPort(port)
	default:
		err = fmt.Errorf("unsupported OS: %s", operatingSystem)
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func GetRandomPort() int {
	port := 0
	for {
		if listener, err := net.Listen("tcp", ":0"); err == nil {
			port = listener.Addr().(*net.TCPAddr).Port
			if err := listener.Close(); err == nil {
				if port > 1024 {
					break
				}
			}
		} else {
			port = 0
			break
		}
	}
	return port
}

//goland:noinspection GoUnusedExportedFunction
func IsPortAvailable(port int) bool {
	var result bool
	address := fmt.Sprintf("localhost:%d", port)
	var conn net.Conn
	var err error
	if conn, err = net.DialTimeout("tcp", address, 100*time.Millisecond); err == nil {
		var closeErr error
		if closeErr = conn.Close(); closeErr == nil {
			result = false
		}
	} else if strings.Contains(err.Error(), "refused") {
		result = true
	}
	return result
}

//goland:noinspection SpellCheckingInspection
func getDarwinIpAddressesWithSmallestMetric() (IP_ADDRESSES, error) {
	result := IP_ADDRESSES{}
	err := error(nil)
	bestInterfaceName := ""
	found := false
	command := exec.Command("netstat", "-rn")
	output := make([]byte, 0)
	if output, err = command.Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if strings.HasPrefix(line, "Routing tables") || strings.HasPrefix(line, "Internet:") || strings.HasPrefix(line, "Destination") {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) >= 5 && fields[0] == "default" {
				bestInterfaceName = fields[3]
				found = true
			}
		}
	}
	if found && bestInterfaceName != "" {
		networkInterface, interfaceError := net.InterfaceByName(bestInterfaceName)
		if interfaceError == nil {
			result = getIpAddressesFromInterface(networkInterface)
		}
	} else {
		result = getIpAddressesFromAllInterfaces()
	}
	return result, err
}

func getDarwinProcessIdByPort(port int) (*int, error) {
	result := (*int)(nil)
	err := error(nil)
	command := exec.Command("lsof", "-i", fmt.Sprintf(":%d", port))
	output, err := command.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			result = nil
			err = nil
		}
	} else {
		lines := strings.Split(string(output), "\n")
		for i, line := range lines {
			if i == 0 {
				continue
			}
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				pidStr := fields[1]
				pid, parseError := strconv.Atoi(pidStr)
				if parseError == nil {
					result = &pid
					break
				}
			}
		}
	}
	return result, err
}

func getIpAddressesFromAllInterfaces() IP_ADDRESSES {
	result := IP_ADDRESSES{}
	interfaces, interfaceError := net.Interfaces()
	if interfaceError == nil {
		for _, networkInterface := range interfaces {
			if networkInterface.Flags&net.FlagUp != 0 {
				addresses, addressError := networkInterface.Addrs()
				if addressError == nil {
					for _, address := range addresses {
						ipNet, ok := address.(*net.IPNet)
						if ok {
							ip := ipNet.IP
							if ip.To4() != nil && !ip.IsLoopback() && result.IPv4 == "" {
								result.IPv4 = ip.String()
							} else if ip.To16() != nil && !ip.IsLoopback() && result.IPv6 == "" {
								result.IPv6 = ip.String()
							}
						}
					}
				}
			}
		}
	}
	return result
}

func getIpAddressesFromInterface(networkInterface *net.Interface) IP_ADDRESSES {
	result := IP_ADDRESSES{}
	addresses, addressError := networkInterface.Addrs()
	if addressError == nil {
		for _, address := range addresses {
			ipNet, ok := address.(*net.IPNet)
			if ok {
				ip := ipNet.IP
				if ip.To4() != nil && !ip.IsLoopback() {
					result.IPv4 = ip.String()
				} else if ip.To16() != nil && !ip.IsLoopback() {
					result.IPv6 = ip.String()
				}
			}
		}
	}
	return result
}

//goland:noinspection SpellCheckingInspection,GoUnhandledErrorResult
func getLinuxIpAddressesWithSmallestMetric() (IP_ADDRESSES, error) {
	result := IP_ADDRESSES{}
	err := error(nil)
	bestInterfaceName := ""
	bestMetric := 0
	found := false
	command := exec.Command("ip", "route")
	output, err := command.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if strings.HasPrefix(line, "default") {
				fields := strings.Fields(line)
				networkInterfaceName := ""
				metric := 0
				for i, field := range fields {
					if field == "dev" && i+1 < len(fields) {
						networkInterfaceName = fields[i+1]
					} else if field == "metric" && i+1 < len(fields) {
						fmt.Sscanf(fields[i+1], "%d", &metric)
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
		networkInterface, interfaceError := net.InterfaceByName(bestInterfaceName)
		if interfaceError == nil {
			result = getIpAddressesFromInterface(networkInterface)
		}
	} else {
		result = getIpAddressesFromAllInterfaces()
	}
	return result, err
}

func getLinuxProcessIdByPort(port int) (*int, error) {
	result := (*int)(nil)
	err := error(nil)
	command := exec.Command("lsof", "-i", fmt.Sprintf(":%d", port))
	output, err := command.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			result = nil
			err = nil
		}
	} else {
		lines := strings.Split(string(output), "\n")
		for i, line := range lines {
			if i == 0 {
				continue
			}
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				pidStr := fields[1]
				pid, parseError := strconv.Atoi(pidStr)
				if parseError == nil {
					result = &pid
					break
				}
			}
		}
	}
	return result, err
}

//goland:noinspection GoUnhandledErrorResult,SpellCheckingInspection
func getWindowsIpAddressesWithSmallestMetric() (IP_ADDRESSES, error) {
	result := IP_ADDRESSES{}
	err := error(nil)
	bestInterfaceIndex := 0
	bestMetric := 0
	found := false
	command := exec.Command("netstat", "-rn")
	output, err := command.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if strings.HasPrefix(line, "Network Destination") || strings.HasPrefix(line, "===========") {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) >= 5 && fields[0] == "0.0.0.0" {
				interfaceIP := fields[3]
				metric := 0
				if len(fields) > 4 {
					fmt.Sscanf(fields[4], "%d", &metric)
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
									ipNet, ok := address.(*net.IPNet)
									if ok {
										ip := ipNet.IP
										if ip.String() == interfaceIP {
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
		networkInterface, interfaceError := net.InterfaceByIndex(bestInterfaceIndex)
		if interfaceError == nil {
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
	command := exec.Command("netstat", "-ano")
	output, err := command.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		portStr := fmt.Sprintf(":%d", port)
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if strings.HasPrefix(line, "Proto") {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				localAddr := fields[1]
				if strings.HasSuffix(localAddr, portStr) {
					pidStr := fields[4]
					pid, parseError := strconv.Atoi(pidStr)
					if parseError == nil {
						result = &pid
						break
					}
				}
			}
		}
	}
	return result, err
}
