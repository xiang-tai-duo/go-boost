// Package boost
// File:        network.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/network.go
// Author:      Vibe Coding
// Created:     01/07/2026 11:03:46
// Description: NETWORK provides functions to get network IP addresses with the smallest metric.
// --------------------------------------------------------------------------------
package boost

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

//goland:noinspection GoSnakeCaseUsage
type (
	IP_ADDRESSES struct {
		IPv4 string
		IPv6 string
	}
)

//goland:noinspection SpellCheckingInspection
func GetNetworkIpAddresses() (IP_ADDRESSES, error) {
	var result IP_ADDRESSES
	var err error
	os := runtime.GOOS
	switch os {
	case "windows":
		result, err = getWindowsIpAddressesWithSmallestMetric()
	case "linux":
		result, err = getLinuxIpAddressesWithSmallestMetric()
	case "darwin":
		result, err = getDarwinIpAddressesWithSmallestMetric()
	default:
		err = fmt.Errorf("unsupported OS: %s", os)
	}
	return result, err
}

//goland:noinspection SpellCheckingInspection
func getDarwinIpAddressesWithSmallestMetric() (IP_ADDRESSES, error) {
	var result IP_ADDRESSES
	var err error
	var bestInterfaceName string
	var found bool
	cmd := exec.Command("netstat", "-rn")
	output, err := cmd.Output()
	if err == nil {
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

// getIpAddressesFromAllInterfaces extracts IPv4 and IPv6 addresses from all active network interfaces
func getIpAddressesFromAllInterfaces() IP_ADDRESSES {
	var result IP_ADDRESSES
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
	var result IP_ADDRESSES
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
	var result IP_ADDRESSES
	var err error
	var bestInterfaceName string
	var bestMetric int
	var found bool
	cmd := exec.Command("ip", "route")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if strings.HasPrefix(line, "default") {
				fields := strings.Fields(line)
				var networkInterfaceName string
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

//goland:noinspection GoUnhandledErrorResult,SpellCheckingInspection
func getWindowsIpAddressesWithSmallestMetric() (IP_ADDRESSES, error) {
	var result IP_ADDRESSES
	var err error
	var bestInterfaceIndex int
	var bestMetric int
	var found bool
	cmd := exec.Command("netstat", "-rn")
	output, err := cmd.Output()
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
