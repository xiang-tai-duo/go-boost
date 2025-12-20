// --------------------------------------------------------------------------------
// File:        serve.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: Examples of using the go-boost HTTP server with WebSocket support
// --------------------------------------------------------------------------------

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xiang-tai-duo/go-boost"
)

// sample_serve_main is the main function for the serve examples
func sample_serve_main() {
	// Example 1: Basic HTTP Server
	fmt.Println("=== Example 1: Basic HTTP Server ===")
	sample_serve_basic_http_server()

	// Example 2: HTTP Server with Routes
	fmt.Println("\n=== Example 2: HTTP Server with Routes ===")
	sample_serve_http_server_with_routes()

	// Example 3: HTTPS Server with TLS
	fmt.Println("\n=== Example 3: HTTPS Server with TLS ===")
	sample_serve_https_server_with_tls()

	// Example 4: WebSocket Echo Server
	fmt.Println("\n=== Example 4: WebSocket Echo Server ===")
	sample_serve_websocket_echo_server()

	// Example 5: WebSocket JSON API Server
	fmt.Println("\n=== Example 5: WebSocket JSON API Server ===")
	sample_serve_websocket_json_api_server()

	// Example 6: Static File Server
	fmt.Println("\n=== Example 6: Static File Server ===")
	sample_serve_static_file_server()

	// Example 7: WebSocket Complete Demo (Server + Client + Data Exchange)
	fmt.Println("\n=== Example 7: WebSocket Complete Demo ===")
	sample_serve_websocket_complete_demo()
}

// sample_serve_basic_http_server demonstrates a basic HTTP server setup
func sample_serve_basic_http_server() {
	// Create a new server instance
	serve := boost.NewServe()

	// Register a simple GET route
	serve.On(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) error {
		_, err := fmt.Fprintf(w, "Hello, World!\n")
		return err
	})

	fmt.Println("Basic HTTP server started on http://localhost:8080")
	fmt.Println("Route: GET / -> Hello, World!")

	// Server automatically starts when On is called
	// No need to explicitly call Listen()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)
}

// sample_serve_http_server_with_routes demonstrates a server with multiple routes
func sample_serve_http_server_with_routes() {
	// Create a new server instance
	serve := boost.NewServe()

	// Register multiple routes
	serve.On(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) error {
		_, err := fmt.Fprintf(w, "Home Page\n")
		return err
	})

	serve.On(http.MethodGet, "/about", func(w http.ResponseWriter, r *http.Request) error {
		_, err := fmt.Fprintf(w, "About Page\n")
		return err
	})

	serve.On(http.MethodPost, "/api/data", func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")
		return json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Data received"})
	})

	fmt.Println("HTTP server with routes started on http://localhost:8081")
	fmt.Println("Routes:")
	fmt.Println("  GET / -> Home Page")
	fmt.Println("  GET /about -> About Page")
	fmt.Println("  POST /api/data -> JSON response")

	// Server automatically starts when On is called
	// No need to explicitly call Listen()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)
}

// sample_serve_https_server_with_tls demonstrates an HTTPS server with TLS encryption
func sample_serve_https_server_with_tls() {
	// Create a new server instance
	serve := boost.NewServe()

	// Register a route
	serve.On(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) error {
		_, err := fmt.Fprintf(w, "Secure HTTPS Server\n")
		return err
	})

	fmt.Println("HTTPS server with TLS started on https://localhost:8443")
	fmt.Println("Note: Using self-signed certificate (may show security warning in browser)")

	// Server automatically starts when On is called
	// No need to explicitly call ListenTLS()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)
}

// sample_serve_websocket_echo_server demonstrates a WebSocket echo server
func sample_serve_websocket_echo_server() {
	// Create a new server instance
	serve := boost.NewServe()

	// Register WebSocket data handler
	serve.OnWebSocket("/ws", nil, func(conn *websocket.Conn, messageType int, data []byte) error {
		fmt.Printf("[Echo] Received: %s\n", string(data))
		// Echo message back to client
		return conn.WriteMessage(messageType, data)
	})

	fmt.Println("WebSocket echo server started on ws://localhost:8084")
	fmt.Println("WebSocket endpoint: ws://localhost:8084/ws")

	// Server automatically starts when OnWebSocket is called
	// No need to explicitly call Listen()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)
}

// sample_serve_websocket_json_api_server demonstrates a WebSocket server that handles JSON messages
func sample_serve_websocket_json_api_server() {
	// Create a new server instance
	serve := boost.NewServe()

	// Define request and response structs
	type Request struct {
		Action string      `json:"action"`
		Data   interface{} `json:"data"`
	}

	type Response struct {
		Status  string      `json:"status"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}

	// Register WebSocket data handler with origin filter (allow localhost only)
	serve.OnWebSocket("/api", map[string]bool{
		"http://localhost:3000": true,
		"http://127.0.0.1:3000": true,
	}, func(conn *websocket.Conn, messageType int, data []byte) error {
		// Parse incoming JSON
		var req Request
		if err := json.Unmarshal(data, &req); err != nil {
			resp := Response{
				Status:  "error",
				Message: "Invalid JSON format",
			}
			respData, _ := json.Marshal(resp)
			return conn.WriteMessage(messageType, respData)
		}

		// Process based on action
		var resp Response
		switch req.Action {
		case "ping":
			resp = Response{
				Status:  "success",
				Message: "Pong response",
				Data: map[string]string{
					"timestamp": time.Now().Format(time.RFC3339),
				},
			}
		case "time":
			resp = Response{
				Status:  "success",
				Message: "Current time",
				Data: map[string]string{
					"time": time.Now().Format(time.RFC3339),
				},
			}
		case "greet":
			name, ok := req.Data.(string)
			if !ok {
				name = "Guest"
			}
			resp = Response{
				Status:  "success",
				Message: "Hello " + name + "!",
			}
		default:
			resp = Response{
				Status:  "error",
				Message: "Unknown action: " + req.Action,
			}
		}

		// Send response
		respData, _ := json.Marshal(resp)
		return conn.WriteMessage(messageType, respData)
	})

	fmt.Println("WebSocket JSON API server started on ws://localhost:8085")
	fmt.Println("WebSocket endpoint: ws://localhost:8085/api")
	fmt.Println("Allowed origins: http://localhost:3000, http://127.0.0.1:3000")
	fmt.Println("Supported actions: ping, time, greet")

	// Server automatically starts when OnWebSocket is called
	// No need to explicitly call Listen()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)
}

// sample_serve_static_file_server demonstrates a static file server
func sample_serve_static_file_server() {
	// Create a new server instance
	serve := boost.NewServe()

	// Add static directory mapping
	// Note: You need to create a 'public' directory with some files to test this
	serve.AddStaticDirectory("/static", "./public")

	// Register a simple route to display static file info
	serve.On(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) error {
		_, err := fmt.Fprintf(w, "Static File Server\n")
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, "Static files available at: /static/*\n")
		return err
	})

	fmt.Println("Static file server started on http://localhost:8086")
	fmt.Println("Static files served from: ./public -> http://localhost:8086/static/")
	fmt.Println("Note: Create a 'public' directory with files to test this example")

	// Server automatically starts when AddStaticDirectory/On is called
	// No need to explicitly call Listen()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)
}

// sample_serve_websocket_complete_demo demonstrates a complete WebSocket flow:
// 1. Start a WebSocket server
// 2. Connect to it as a client
// 3. Send messages from client to server
// 4. Server processes and echoes messages back
// 5. Log all activities
func sample_serve_websocket_complete_demo() {
	// Step 1: Start WebSocket server
	serve := boost.NewServe()

	// Get a random available port greater than 1024
	port, err := serve.GetAvailablePort()
	if err != nil {
		log.Fatalf("Failed to get random available port: %v", err)
	}
	url := fmt.Sprintf("ws://localhost:%d/ws", port)

	fmt.Printf("Starting WebSocket complete demo on port %d\n", port)

	// Set server to listen on the random port
	serve.Server.Addr = fmt.Sprintf(":%d", port)

	// Register WebSocket data handler that logs messages and echoes them back
	serve.OnWebSocket("/ws", nil, func(conn *websocket.Conn, messageType int, data []byte) error {
		log.Printf("[Server] Received: %s", string(data))

		// Echo message back with timestamp
		response := fmt.Sprintf("[Echo] %s - %s", string(data), time.Now().Format(time.RFC3339))
		log.Printf("[Server] Sending: %s", response)

		return conn.WriteMessage(messageType, []byte(response))
	})

	// Wait for server to start
	time.Sleep(500 * time.Millisecond)

	// Step 2: Connect to WebSocket server as client
	log.Printf("[Client] Connecting to %s", url)

	// Create WebSocket client instance directly
	wsClient := boost.NewWebSocket()

	// Messages to send
	testMessages := []string{
		"Hello, WebSocket!",
		"How are you?",
		"This is a test message",
		"WebSocket is awesome!",
	}

	// Channel to signal when all messages are processed
	done := make(chan struct{})
	messageCount := 0
	expectedMessages := len(testMessages)

	// Test 1: Connect with explicit name
	clientName, err := wsClient.Connect(
		url,
		// Message handler for client
		func(conn *websocket.Conn, messageType int, data []byte) error {
			log.Printf("[Client] Received: %s", string(data))
			messageCount++
			if messageCount == expectedMessages {
				close(done)
			}
			return nil
		},
	)
	if err != nil {
		log.Printf("[Client] Failed to connect: %v", err)
		return
	}
	log.Printf("[Client] Connected with explicit name: %s", clientName)

	// Test 2: Connect with auto-generated name
	autoGenName, err := wsClient.Connect(
		url,
		// Message handler for auto-generated client
		func(conn *websocket.Conn, messageType int, data []byte) error {
			log.Printf("[Auto-client] Received: %s", string(data))
			return nil
		},
	)
	if err != nil {
		log.Printf("[Client] Failed to connect with auto-generated name: %v", err)
		return
	}
	log.Printf("[Client] Connected with auto-generated name: %s", autoGenName)

	// Step 3: Send test messages from client to server
	for _, msg := range testMessages {
		log.Printf("[Client] Sending: %s", msg)
		err := wsClient.Send(clientName, []byte(msg))
		if err != nil {
			log.Printf("[Client] Send error: %v", err)
		}
		// Wait a bit between messages
		time.Sleep(300 * time.Millisecond)
	}

	// Wait for all messages to be processed
	select {
	case <-done:
		log.Println("[Demo] All messages processed successfully!")
	case <-time.After(5 * time.Second):
		log.Println("[Demo] Timeout waiting for messages")
	}

	log.Println("[Demo] WebSocket complete demo finished!")
}

// Keep a main function for testing if needed
func main() {
	// Directly run only the WebSocket complete demo to see the effect
	sample_serve_websocket_complete_demo()

	// Keep main function running to demonstrate all examples
	fmt.Println("\nWebSocket demo completed. Press Ctrl+C to exit.")
	select {}
}
