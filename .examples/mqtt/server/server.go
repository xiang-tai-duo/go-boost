// Package main
// File:        https://github.com/xiang-tai-duo/go-boost/blob/master/.examples/mqtt/server/server.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: MQTT server usage example
// --------------------------------------------------------------------------------
package main

import (
	"fmt"
	"time"

	mqttserver "github.com/xiang-tai-duo/go-boost/mqtt/server"
)

func main() {
	// Example 1: Create MQTT server with default settings
	// This will create a server with default host "0.0.0.0" and port 1883
	server1 := mqttserver.New()
	fmt.Printf("Server 1 - Host: %s, Port: %d\n", server1.GetHost(), server1.GetPort())

	// Example 2: Create MQTT server with custom host
	server2 := mqttserver.New("127.0.0.1")
	fmt.Printf("Server 2 - Host: %s, Port: %d\n", server2.GetHost(), server2.GetPort())

	// Example 3: Create MQTT server with custom host and port
	server3 := mqttserver.New("0.0.0.0", 1884)
	fmt.Printf("Server 3 - Host: %s, Port: %d\n", server3.GetHost(), server3.GetPort())

	// Example 4: Start the MQTT server
	fmt.Println("Starting MQTT server...")
	var err error
	if err = server1.Start(); err == nil {
		fmt.Println("MQTT server started successfully")

		// Check if server is running
		fmt.Printf("Server is running: %v\n", server1.IsRunning())

		// Example 5: Keep server running for a while
		fmt.Println("Server will run for 30 seconds...")
		fmt.Println("You can connect to it using an MQTT client (e.g., mosquitto_sub or mosquitto_pub)")
		fmt.Println("Press Ctrl+C to exit early")

		// Wait for 30 seconds
		time.Sleep(30 * time.Second)

		// Example 6: Stop the MQTT server
		fmt.Println("Stopping MQTT server...")
		if err = server1.Stop(); err == nil {
			fmt.Println("MQTT server stopped successfully")

			// Check if server is running after stop
			fmt.Printf("Server is running: %v\n", server1.IsRunning())
		} else {
			fmt.Printf("Failed to stop server: %v\n", err)
		}
	} else {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
