// Package main
// File:        https://github.com/xiang-tai-duo/go-boost/blob/master/.examples/serve/serve.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: HTTP and WebSocket server usage example
// --------------------------------------------------------------------------------
package main

import (
	"fmt"
	"net/http"
	"time"

	__websocket "github.com/gorilla/websocket"
	__serve "github.com/xiang-tai-duo/go-boost/serve"
)

func main() {
	// Example 1: Create HTTP server
	server := __serve.New()
	fmt.Println("Example 1: Created HTTP server")

	// Example 2: Add routes
	fmt.Println("\nExample 2: Adding routes")

	// Add GET route
	server.On("GET", "/", func(request *http.Request, response http.ResponseWriter) error {
		response.Header().Set("Content-Type", "text/plain")
		response.WriteHeader(http.StatusOK)
		_, err := response.Write([]byte("Hello, World! This is the root route.\n"))
		return err
	})

	// Add GET route with path parameter
	server.On("GET", "/user/{id}", func(request *http.Request, response http.ResponseWriter) error {
		response.Header().Set("Content-Type", "text/plain")
		response.WriteHeader(http.StatusOK)
		_, err := response.Write([]byte("User profile route\n"))
		return err
	})

	// Add POST route
	server.On("POST", "/api/data", func(request *http.Request, response http.ResponseWriter) error {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusCreated)
		_, err := response.Write([]byte(`{"status": "success", "message": "Data received"}`))
		return err
	})

	// Example 3: Add WebSocket route
	fmt.Println("\nExample 3: Adding WebSocket route")
	server.OnWebSocket("/ws", func(ws *__websocket.Conn, messageType int, data []byte) error {
		fmt.Printf("Received WebSocket message: %s\n", data)
		// Echo the message back
		return ws.WriteMessage(messageType, data)
	})

	// Example 4: Add static file directory
	fmt.Println("\nExample 4: Adding static file directory")
	// Uncomment the line below and replace "path/to/static/files" with actual directory
	// server.AddStaticDirectory("/static", "path/to/static/files")

	// Example 5: Check if port is available
	fmt.Println("\nExample 5: Checking port availability")
	port := 8080
	if server.CheckPortAvailable(port) {
		fmt.Printf("Port %d is available\n", port)
	} else {
		fmt.Printf("Port %d is not available\n", port)
		// Try to get an available port
		if availablePort, err := server.GetAvailablePort(); err == nil {
			fmt.Printf("Found available port: %d\n", availablePort)
			port = availablePort
		}
	}

	// Example 6: Start HTTP server
	address := fmt.Sprintf(":%d", port)
	fmt.Printf("\nExample 6: Starting HTTP server on %s\n", address)
	err := server.Listen(address)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		return
	}

	// Check if server is running
	fmt.Printf("Server is running: %v\n", server.IsRunning())

	// Example 7: Keep server running
	fmt.Println("\nExample 7: Server is running")
	fmt.Println("You can access the following endpoints:")
	fmt.Printf("- http://localhost:%d/ (Root route)\n", port)
	fmt.Printf("- http://localhost:%d/user/123 (User route with parameter)\n", port)
	fmt.Printf("- http://localhost:%d/api/data (POST route)\n", port)
	fmt.Printf("- ws://localhost:%d/ws (WebSocket endpoint)\n", port)
	fmt.Println("\nServer will run for 60 seconds...")
	fmt.Println("Press Ctrl+C to exit early")

	// Wait for 60 seconds
	time.Sleep(60 * time.Second)

	// Example 8: Shutdown server
	fmt.Println("\nExample 8: Shutting down server")
	err = server.Shutdown()
	if err == nil {
		fmt.Println("Server shutdown successfully")

		// Check if server is running after shutdown
		fmt.Printf("Server is running: %v\n", server.IsRunning())
	} else {
		fmt.Printf("Failed to shutdown server: %v\n", err)
	}
}
