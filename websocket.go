// --------------------------------------------------------------------------------
// File:        http_websocket.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: WebSocket client functionality for Go applications
// --------------------------------------------------------------------------------

package boost

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

//goland:noinspection SpellCheckingInspection
type (
	// WebSocketClientHandler defines the signature for WebSocket client message handlers
	WebSocketClientHandler func(conn *websocket.Conn, messageType int, data []byte) error

	// WebSocketClientErrorHandler defines the signature for WebSocket client error handlers
	WebSocketClientErrorHandler func(err error)

	// WEBSOCKET_CLIENT_CONFIG contains configuration for WebSocket client connections
	WEBSOCKET_CLIENT_CONFIG struct {
		Url                  string
		Origin               string
		Headers              http.Header
		Dialer               *websocket.Dialer
		ReconnectInterval    time.Duration
		MaxReconnectAttempts int
		SendChannelSize      int
	}

	// WEBSOCKET_CLIENT represents a WebSocket client connection
	WEBSOCKET_CLIENT struct {
		Config         WEBSOCKET_CLIENT_CONFIG
		Conn           *websocket.Conn
		MessageHandler WebSocketClientHandler
		ErrorHandler   WebSocketClientErrorHandler
		SendChan       chan []byte
		IsConnected    bool
		ReconnectCount int
		Mutex          sync.Mutex
	}

	// WEBSOCKET_CLIENT_MANAGER provides WebSocket client functionality
	WEBSOCKET_CLIENT_MANAGER struct {
		Clients       map[string]*WEBSOCKET_CLIENT
		DefaultConfig WEBSOCKET_CLIENT_CONFIG
		Mutex         sync.Mutex
	}
)

// NewWebSocket creates a new WEBSOCKET_CLIENT_MANAGER instance
// Returns: A new WEBSOCKET_CLIENT_MANAGER instance
func NewWebSocket() *WEBSOCKET_CLIENT_MANAGER {
	return &WEBSOCKET_CLIENT_MANAGER{
		Clients: make(map[string]*WEBSOCKET_CLIENT),
		DefaultConfig: WEBSOCKET_CLIENT_CONFIG{
			ReconnectInterval:    5 * time.Second,
			MaxReconnectAttempts: 5,
			SendChannelSize:      100,
		},
	}
}

// Connect creates a WebSocket client connection with default settings
// uri: WebSocket server URI, e.g., "ws://localhost:8080/ws"
// messageHandler: Function to handle received WebSocket messages
// Returns: The client ID and any error encountered
func (ws *WEBSOCKET_CLIENT_MANAGER) Connect(uri string, messageHandler WebSocketClientHandler) (string, error) {
	return ws.ConnectEx(uri, "", nil, nil, 0, 0, 0, messageHandler)
}

// ConnectEx creates a WebSocket client connection with full configuration options
// uri: WebSocket server URI, e.g., "ws://localhost:8080/ws"
// origin: Optional origin header to send
// headers: Optional custom headers to send
// dialer: Optional custom dialer for WebSocket connection
// reconnectInterval: Optional interval between reconnect attempts (default: 5s)
// maxReconnectAttempts: Optional maximum number of reconnect attempts (default: 5)
// sendChannelSize: Optional size of the send channel (default: 100)
// messageHandler: Function to handle received WebSocket messages
// Returns: The client ID and any error encountered
func (ws *WEBSOCKET_CLIENT_MANAGER) ConnectEx(uri string, origin string, headers http.Header, dialer *websocket.Dialer, reconnectInterval time.Duration, maxReconnectAttempts int, sendChannelSize int, messageHandler WebSocketClientHandler) (string, error) {
	var clientID string
	var err error

	// Validate required parameters
	if uri == "" {
		err = fmt.Errorf("websocket server URI is required")
	} else if messageHandler == nil {
		err = fmt.Errorf("websocket message handler is required")
	} else {
		ws.Mutex.Lock()
		defer ws.Mutex.Unlock()

		// Set default values if not provided
		if headers == nil {
			headers = make(http.Header)
		}
		if dialer == nil {
			dialer = &websocket.Dialer{}
		}
		if reconnectInterval == 0 {
			reconnectInterval = ws.DefaultConfig.ReconnectInterval
		}
		if maxReconnectAttempts == 0 {
			maxReconnectAttempts = ws.DefaultConfig.MaxReconnectAttempts
		}
		if sendChannelSize == 0 {
			sendChannelSize = ws.DefaultConfig.SendChannelSize
		}

		// Generate client ID
		clientID = fmt.Sprintf("websocket-client-%d", len(ws.Clients)+1)

		// Create WebSocket client configuration
		config := WEBSOCKET_CLIENT_CONFIG{
			Url:                  uri,
			Origin:               origin,
			Headers:              headers,
			Dialer:               dialer,
			ReconnectInterval:    reconnectInterval,
			MaxReconnectAttempts: maxReconnectAttempts,
			SendChannelSize:      sendChannelSize,
		}

		// Create new WebSocket client
		client := &WEBSOCKET_CLIENT{
			Config:         config,
			MessageHandler: messageHandler,
			SendChan:       make(chan []byte, sendChannelSize),
			IsConnected:    false,
			ReconnectCount: 0,
		}

		// Store client
		ws.Clients[clientID] = client

		// Start connection goroutine
		go client.connect()
	}

	return clientID, err
}

// Shutdown closes a WebSocket client connection
// clientID: Client ID used for reference
// Returns: Error if the client is not found
func (ws *WEBSOCKET_CLIENT_MANAGER) Shutdown(clientID string) error {
	var err error

	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()

	client, exists := ws.Clients[clientID]
	if !exists {
		err = fmt.Errorf("websocket client not found: %s", clientID)
	} else {
		// Close send channel to signal connection goroutine to exit
		close(client.SendChan)

		// Remove client from map
		delete(ws.Clients, clientID)
	}

	return err
}

// Send sends a message to a WebSocket client connection
// clientID: Client ID used for reference
// message: Message to send
// Returns: Error if the client is not found or disconnected
func (ws *WEBSOCKET_CLIENT_MANAGER) Send(clientID string, message []byte) error {
	var err error
	var client *WEBSOCKET_CLIENT
	var exists bool

	ws.Mutex.Lock()
	client, exists = ws.Clients[clientID]
	ws.Mutex.Unlock()

	if !exists {
		err = fmt.Errorf("websocket client not found: %s", clientID)
	} else {
		select {
		case client.SendChan <- message:
			// Send successful
		default:
			err = fmt.Errorf("websocket client send channel is full or client is disconnected: %s", clientID)
		}
	}

	return err
}

// connect establishes a WebSocket connection and handles message exchange
func (c *WEBSOCKET_CLIENT) connect() {
	for {
		// Establish connection
		conn, _, err := c.Config.Dialer.Dial(c.Config.Url, c.Config.Headers)
		if err != nil {
			// Check if we should reconnect
			c.Mutex.Lock()
			if c.ReconnectCount < c.Config.MaxReconnectAttempts {
				c.ReconnectCount++
				reconnectInterval := c.Config.ReconnectInterval
				c.Mutex.Unlock()
				time.Sleep(reconnectInterval)
				continue
			}
			c.Mutex.Unlock()
			break
		}

		// Update client state
		c.Mutex.Lock()
		c.Conn = conn
		c.IsConnected = true
		c.ReconnectCount = 0
		c.Mutex.Unlock()

		// Start read and write goroutines
		readDone := make(chan struct{})
		writeDone := make(chan struct{})

		// Read goroutine
		go c.readLoop(readDone)

		// Write goroutine
		go c.writeLoop(writeDone)

		// Wait for either read or write loop to exit
		select {
		case <-readDone:
			close(writeDone)
		case <-writeDone:
			close(readDone)
		}

		// Cleanup connection
		conn.Close()

		c.Mutex.Lock()
		c.IsConnected = false
		c.Mutex.Unlock()

		// Check if we should reconnect
		select {
		case <-c.SendChan:
			// If there's still data to send, reconnect
			c.Mutex.Lock()
			c.ReconnectCount++
			reconnectInterval := c.Config.ReconnectInterval
			c.Mutex.Unlock()
			time.Sleep(reconnectInterval)
		default:
			// No more data to send, exit
			return
		}
	}
}

// readLoop reads messages from the WebSocket connection
func (c *WEBSOCKET_CLIENT) readLoop(done chan struct{}) {
	defer close(done)

	var shouldExit bool
	for !shouldExit {
		messageType, message, err := c.Conn.ReadMessage()
		if err != nil {
			shouldExit = true
		} else {
			// Call message handler
			if c.MessageHandler != nil {
				if err := c.MessageHandler(c.Conn, messageType, message); err != nil {
					// Handle error internally
				}
			}
		}
	}
}

// writeLoop writes messages to the WebSocket connection
func (c *WEBSOCKET_CLIENT) writeLoop(done chan struct{}) {
	defer close(done)

	var shouldExit bool
	for !shouldExit {
		select {
		case message, ok := <-c.SendChan:
			if !ok {
				// Send channel closed, exit
				shouldExit = true
			} else if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				shouldExit = true
			}
		case <-done:
			shouldExit = true
		}
	}
}
