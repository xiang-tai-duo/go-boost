// Package main
// File:        network.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/examples/network.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: Example for network functionality
// --------------------------------------------------------------------------------
package main

import (
	"fmt"

	"github.com/xiang-tai-duo/go-boost"
)

func main() {
	// Get network IPs with smallest metric
	ips, err := boost.GetNetworkIPWithSmallestMetric()
	if err != nil {
		fmt.Printf("Error getting network IPs: %v\n", err)
		return
	}
	fmt.Printf("Network IPs with smallest metric:\n")
	fmt.Printf("  IPv4: %s\n", ips.IPv4)
	fmt.Printf("  IPv6: %s\n", ips.IPv6)
}
