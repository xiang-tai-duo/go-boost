// --------------------------------------------------------------------------------
// File:        example_websocket.go
// Author:      TRAE AI
// Created:     2025/12/23 15:00:00
// Description: Example of using WebSocketData function to create a WebSocket server
// --------------------------------------------------------------------------------

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go-boost"

	"github.com/gorilla/websocket"
)

func main() {
	// 1. Create a new server instance
	serve := boost.NewServe()

	// 2. Register WebSocket endpoints

	// Example 1: Basic Echo Server
	// This endpoint simply echoes back any message it receives
	serve.WebSocketData("/echo", nil, func(conn *websocket.Conn, messageType int, data []byte) error {
		fmt.Printf("[Echo] Received: %s\n", string(data))
		// Just echo the message back to the client
		return conn.WriteMessage(messageType, data)
	})

	// Example 2: JSON API Server
	// This endpoint processes JSON messages with actions
	serve.WebSocketData("/api", map[string]bool{
		"http://localhost:3000": true,
		"https://example.com":   true,
	}, func(conn *websocket.Conn, messageType int, data []byte) error {
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

	// Example 3: Binary Data Server
	// This endpoint handles binary messages
	serve.WebSocketData("/binary", nil, func(conn *websocket.Conn, messageType int, data []byte) error {
		if messageType == websocket.BinaryMessage {
			fmt.Printf("[Binary] Received %d bytes\n", len(data))
			// Echo back binary data with a prefix
			response := append([]byte("ECHO: "), data...)
			return conn.WriteMessage(messageType, response)
		}
		return nil
	})

	// 3. Start the server
	fmt.Println("Starting WebSocket server on :8080")
	fmt.Println("WebSocket endpoints:")
	fmt.Println("  - ws://localhost:8080/echo     (Echo server)")
	fmt.Println("  - ws://localhost:8080/api      (JSON API server)")
	fmt.Println("  - ws://localhost:8080/binary   (Binary data server)")
	fmt.Println("\nPress Ctrl+C to stop the server")

	if err := serve.Listen(":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
