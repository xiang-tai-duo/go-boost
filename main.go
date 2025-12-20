// --------------------------------------------------------------------------------
// File:        main.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: Main entry point for go-boost application, demonstrating WebSocket functionality
// --------------------------------------------------------------------------------

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"go-boost"
)

// Config holds the application configuration
// Contains server addresses and WebSocket settings
var Config = struct {
	HttpAddress  string
	HttpsAddress string
	WebSocketPath string
	WebsocketSslPath string
}{
	HttpAddress:  ":8080",
	HttpsAddress: ":8443",
	WebSocketPath: "/ws",
	WebsocketSslPath: "/wss",
}

// SetHttpAddress sets the HTTP server address
// address: New HTTP server address, e.g., ":8080"
// Usage:
// SetHttpAddress(":8080")
func SetHttpAddress(address string) {
	Config.HttpAddress = address
}

// GetHttpAddress returns the current HTTP server address
// Returns: Current HTTP server address
// Usage:
// address := GetHttpAddress()
// returns ":8080"
func GetHttpAddress() string {
	return Config.HttpAddress
}

// SetHttpsAddress sets the HTTPS server address
// address: New HTTPS server address, e.g., ":8443"
// Usage:
// SetHttpsAddress(":8443")
func SetHttpsAddress(address string) {
	Config.HttpsAddress = address
}

// GetHttpsAddress returns the current HTTPS server address
// Returns: Current HTTPS server address
// Usage:
// address := GetHttpsAddress()
// returns ":8443"
func GetHttpsAddress() string {
	return Config.HttpsAddress
}

// SetWebSocketPath sets the WebSocket path
// path: New WebSocket path, e.g., "/ws"
// Usage:
// SetWebSocketPath("/ws")
func SetWebSocketPath(path string) {
	Config.WebSocketPath = path
}

// GetWebSocketPath returns the current WebSocket path
// Returns: Current WebSocket path
// Usage:
// path := GetWebSocketPath()
// returns "/ws"
func GetWebSocketPath() string {
	return Config.WebSocketPath
}

// SetWebsocketSslPath sets the WebSocket SSL path
// path: New WebSocket SSL path, e.g., "/wss"
// Usage:
// SetWebsocketSslPath("/wss")
func SetWebsocketSslPath(path string) {
	Config.WebsocketSslPath = path
}

// GetWebsocketSslPath returns the current WebSocket SSL path
// Returns: Current WebSocket SSL path
// Usage:
// path := GetWebsocketSslPath()
// returns "/wss"
func GetWebsocketSslPath() string {
	return Config.WebsocketSslPath
}

// main is the entry point of the application
// Starts HTTP and HTTPS servers with WebSocket support
func main() {
	// Create WebSocket upgrader
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// Create HTTP server builder
	httpBuilder := boost.NewServeBuilder()

	// Register WebSocket handler for HTTP
	httpBuilder.WebSocket(Config.WebSocketPath, func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}
		defer func(conn *websocket.Conn) {
			_ = conn.Close()
		}(conn)
		
		log.Printf("New WebSocket connection established: %s", r.RemoteAddr)
		handleWebSocketConnection(conn)
	})

	// Create HTTPS server builder
	httpsBuilder := boost.NewServeBuilder()

	// Register WebSocket handler for HTTPS
	httpsBuilder.WebSocket(Config.WebsocketSslPath, func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket SSL upgrade error: %v", err)
			return
		}
		defer func(conn *websocket.Conn) {
			_ = conn.Close()
		}(conn)
		
		log.Printf("New WebSocket SSL connection established: %s", r.RemoteAddr)
		handleWebSocketConnection(conn)
	})

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("Starting HTTP server on %s", Config.HttpAddress)
		if err := httpBuilder.Listen(Config.HttpAddress); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Start HTTPS server in a goroutine
	go func() {
		log.Printf("Starting HTTPS server on %s", Config.HttpsAddress)
		if err := httpsBuilder.ListenTLS(Config.HttpsAddress, "", ""); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTPS server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the servers
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down servers...")

	// Create a deadline to wait for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := httpBuilder.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server forced to shutdown: %v", err)
	}

	// Shutdown HTTPS server
	if err := httpsBuilder.Shutdown(ctx); err != nil {
		log.Fatalf("HTTPS server forced to shutdown: %v", err)
	}

	log.Println("Servers exited properly")
}

// handleWebSocketConnection handles WebSocket messages from clients
// conn: WebSocket connection to handle
func handleWebSocketConnection(conn *websocket.Conn) {
	for {
		// Read message from client
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		log.Printf("Received message: %s", message)

		// Echo message back to client
		if err = conn.WriteMessage(messageType, message); err != nil {
			log.Printf("WebSocket write error: %v", err)
			break
		}
	}
}
