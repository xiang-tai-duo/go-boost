// --------------------------------------------------------------------------------
// File:        websocket.go
// Author:      TRAE AI
// Created:     12/30/2025 11:03:46
// Description: Examples of using the go-boost WebSocket support with test web page
// --------------------------------------------------------------------------------
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	. "github.com/xiang-tai-duo/go-boost"
)

var websocketEndpoint string

func main() {
	testWebSocketHandshake()
	if initWebSocketServer() {
		testPageURL, _ := initStaticFilesServer()
		if websocketEndpoint != "" {
			testPageURL += fmt.Sprintf("?ws=%s", websocketEndpoint)
		}
		fmt.Printf("Test page available at: %s\n", testPageURL)
	}
	fmt.Println("Running for up to 5 seconds...")

	// Block for at most 5 seconds
	select {
	case <-time.After(5 * time.Second):
		fmt.Println("Program has been running for 5 seconds, exiting automatically")
	}
}

func initStaticFilesServer() (string, string) {
	serve := NewServe()
	var port int
	if p, err := serve.GetAvailablePort(); err != nil {
		log.Fatalf("Failed to get available port: %v", err)
	} else {
		port = p
	}

	// Set the Port field in SERVE struct
	serve.Port = port
	serve.AddStaticDirectory("/", ".")
	serve.On(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) error {
		redirectURL := "/websocket.html"
		if r.URL.RawQuery != "" {
			redirectURL += "?" + r.URL.RawQuery
		}
		http.Redirect(w, r, redirectURL, http.StatusFound)
		return nil
	})
	serve.On(http.MethodGet, "/api/endpoints", func(w http.ResponseWriter, r *http.Request) error {
		response := struct {
			WebSocket string `json:"websocket"`
		}{
			WebSocket: websocketEndpoint,
		}
		w.Header().Set("Content-Type", "application/json")
		return json.NewEncoder(w).Encode(response)
	})
	
	// Start the server in a goroutine
	go func() {
		if err := serve.Listen(); err != nil {
			log.Printf("Static files server error: %v", err)
		}
	}()
	
	time.Sleep(100 * time.Millisecond)
	testPageURL := fmt.Sprintf("http://localhost:%d", port)
	return testPageURL, fmt.Sprintf("%s/api/endpoints", testPageURL)
}

func initWebSocketServer() bool {
	serve := NewServe()
	var port int
	if p, err := serve.GetAvailablePort(); err != nil {
		log.Fatalf("Failed to get available port: %v", err)
	} else {
		port = p
	}

	// Set the Port field in SERVE struct
	serve.Port = port
	serve.OnWebSocketEx("/api", WebSocketAllowAll, func(conn *websocket.Conn, messageType int, data []byte) error {
		var req struct {
			Data string `json:"data"`
		}
		if err := json.Unmarshal(data, &req); err != nil {
			respData, _ := json.Marshal(struct {
				Status  string `json:"status"`
				Message string `json:"message"`
			}{
				Status:  "error",
				Message: "Invalid JSON format",
			})
			return conn.WriteMessage(messageType, respData)
		}
		log.Printf("[SERVER] Received: %s\n", string(data))
		respData, _ := json.Marshal(struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}{
			Status:  "success",
			Message: "Hello " + req.Data + "!",
		})
		log.Printf("[SERVER] Send: %s\n", string(respData))
		return conn.WriteMessage(messageType, respData)
	}, func(conn *websocket.Conn) error {
		log.Printf("[SERVER] Client disconnected\n")
		return nil
	})
	endpoint := fmt.Sprintf("ws://localhost:%d/api", port)
	websocketEndpoint = endpoint
	fmt.Printf("WebSocket JSON API Server started on %s\n", endpoint)
	
	// Start the server in a goroutine
	go func() {
		if err := serve.Listen(); err != nil {
			log.Printf("WebSocket server error: %v", err)
		}
	}()
	
	time.Sleep(100 * time.Millisecond)
	
	// Initialize SSL WebSocket server
	return initSSLWebSocketServer()
}

func initSSLWebSocketServer() bool {
	serve := NewServe()
	var port int
	if p, err := serve.GetAvailablePort(); err != nil {
		log.Printf("Failed to get available port for SSL WebSocket: %v", err)

		// Continue even if SSL server fails
		return true
	} else {
		port = p
	}

	// Set the TlsPort field in SERVE struct for SSL server
	serve.TlsPort = port
	serve.OnWebSocketEx("/api", WebSocketAllowAll, func(conn *websocket.Conn, messageType int, data []byte) error {
		var req struct {
			Data string `json:"data"`
		}
		if err := json.Unmarshal(data, &req); err != nil {
			respData, _ := json.Marshal(struct {
				Status  string `json:"status"`
				Message string `json:"message"`
			}{
				Status:  "error",
				Message: "Invalid JSON format",
			})
			return conn.WriteMessage(messageType, respData)
		}
		log.Printf("[SSL SERVER] Received: %s\n", string(data))
		respData, _ := json.Marshal(struct {
			Status  string `json:"status"`
			Message string `json:"message"`
			SSL     bool   `json:"ssl"`
		}{
			Status:  "success",
			Message: "Hello " + req.Data + "!",
			SSL:     true,
		})
		log.Printf("[SSL SERVER] Send: %s\n", string(respData))
		return conn.WriteMessage(messageType, respData)
	}, func(conn *websocket.Conn) error {
		log.Printf("[SSL SERVER] Client disconnected\n")
		return nil
	})
	endpoint := fmt.Sprintf("wss://localhost:%d/api", port)
	fmt.Printf("SSL WebSocket JSON API Server started on %s\n", endpoint)
	
	// Start the SSL server in a goroutine using serve.Listen()
	go func() {
		if err := serve.Listen(); err != nil {
			log.Printf("SSL WebSocket server error: %v", err)
		}
	}()
	
	time.Sleep(100 * time.Millisecond)
	return true
}

func testWebSocketHandshake() {
	serve := NewServe()
	var port int
	if p, err := serve.GetAvailablePort(); err != nil {
		log.Fatalf("Failed to get available port for handshake test: %v", err)
	} else {
		port = p
	}
	
	// Set the Port field in SERVE struct
	serve.Port = port
	serve.OnWebSocket("/handshake", func(conn *websocket.Conn, messageType int, data []byte) error {
		var clientMsg struct {
			Content string `json:"content"`
			From    string `json:"from"`
		}
		if err := json.Unmarshal(data, &clientMsg); err != nil {
			return err
		}
		log.Printf("[SERVER] Received: %s\n", string(data))
		serverMsg := struct {
			Content string `json:"content"`
			From    string `json:"from"`
			Server  string `json:"server"`
		}{
			Content: clientMsg.Content,
			From:    clientMsg.From,
			Server:  "Server added this content",
		}
		serverData, _ := json.Marshal(serverMsg)
		log.Printf("[SERVER] Sending: %s\n", string(serverData))
		return conn.WriteMessage(messageType, serverData)
	})
	
	// Start the server in a goroutine
	go func() {
		if err := serve.Listen(); err != nil {
			log.Printf("Handshake test server error: %v", err)
		}
	}()
	
	time.Sleep(200 * time.Millisecond)
	wsManager := NewWebSocket()
	uri := fmt.Sprintf("ws://localhost:%d/handshake", port)
	clientID, err := wsManager.Connect(uri, func(conn *websocket.Conn, messageType int, data []byte) error {
		log.Printf("[CLIENT] Received server response: %s\n", string(data))
		return nil
	})
	if err != nil {
		log.Fatalf("[CLIENT] Failed to connect: %v", err)
	}
	initialClientMsg := struct {
		Content string `json:"content"`
		From    string `json:"from"`
	}{
		Content: "Hello from client",
		From:    "client",
	}
	initialClientData, _ := json.Marshal(initialClientMsg)
	log.Printf("[CLIENT] Sending initial message: %s\n", string(initialClientData))
	if err := wsManager.Send(clientID, initialClientData); err != nil {
		log.Fatalf("[CLIENT] Failed to send initial client message: %v", err)
	}
	time.Sleep(500 * time.Millisecond)
	fmt.Println("WebSocket message exchange completed")
}
