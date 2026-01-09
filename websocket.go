// Package boost
// File:        websocket.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/websocket.go
// Author:      Vibe Coding
// Created:     12/30/2025 11:03:46
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

//goland:noinspection SpellCheckingInspection,GoSnakeCaseUsage
type (
	WebSocketClientHandler func(conn *websocket.Conn, messageType int, data []byte) error

	WebSocketClientErrorHandler func(err error)

	WEBSOCKET_CLIENT_CONFIG struct {
		Url                  string
		Origin               string
		Headers              http.Header
		Dialer               *websocket.Dialer
		HandshakeTimeout     time.Duration
		ReconnectInterval    time.Duration
		MaxReconnectAttempts int
		SendChannelSize      int
	}

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

	WEBSOCKET_CLIENT_MANAGER struct {
		Clients       map[string]*WEBSOCKET_CLIENT
		DefaultConfig WEBSOCKET_CLIENT_CONFIG
		Mutex         sync.Mutex
	}
)

func NewWebSocket() *WEBSOCKET_CLIENT_MANAGER {
	return &WEBSOCKET_CLIENT_MANAGER{
		Clients: make(map[string]*WEBSOCKET_CLIENT),
		DefaultConfig: WEBSOCKET_CLIENT_CONFIG{
			HandshakeTimeout:     10 * time.Second,
			ReconnectInterval:    5 * time.Second,
			MaxReconnectAttempts: 5,
			SendChannelSize:      100,
		},
	}
}

func (ws *WEBSOCKET_CLIENT_MANAGER) Connect(uri string, messageHandler WebSocketClientHandler) (string, error) {
	return ws.ConnectEx(uri, "", nil, nil, 0, 0, 0, 0, messageHandler)
}

func (ws *WEBSOCKET_CLIENT_MANAGER) ConnectEx(uri string, origin string, headers http.Header, dialer *websocket.Dialer, handshakeTimeout time.Duration, reconnectInterval time.Duration, maxReconnectAttempts int, sendChannelSize int, messageHandler WebSocketClientHandler) (string, error) {
	var clientID string
	var err error
	if uri == "" {
		err = fmt.Errorf("websocket server URI is required")
	} else if messageHandler == nil {
		err = fmt.Errorf("websocket message handler is required")
	} else {
		ws.Mutex.Lock()
		defer ws.Mutex.Unlock()
		if headers == nil {
			headers = make(http.Header)
		}
		if dialer == nil {
			dialer = &websocket.Dialer{}
		}
		if handshakeTimeout == 0 {
			handshakeTimeout = ws.DefaultConfig.HandshakeTimeout
		}
		dialer.HandshakeTimeout = handshakeTimeout
		if reconnectInterval == 0 {
			reconnectInterval = ws.DefaultConfig.ReconnectInterval
		}
		if maxReconnectAttempts == 0 {
			maxReconnectAttempts = ws.DefaultConfig.MaxReconnectAttempts
		}
		if sendChannelSize == 0 {
			sendChannelSize = ws.DefaultConfig.SendChannelSize
		}
		clientID = fmt.Sprintf("websocket-client-%d", len(ws.Clients)+1)
		config := WEBSOCKET_CLIENT_CONFIG{
			Url:                  uri,
			Origin:               origin,
			Headers:              headers,
			Dialer:               dialer,
			ReconnectInterval:    reconnectInterval,
			MaxReconnectAttempts: maxReconnectAttempts,
			SendChannelSize:      sendChannelSize,
		}
		client := &WEBSOCKET_CLIENT{
			Config:         config,
			MessageHandler: messageHandler,
			SendChan:       make(chan []byte, sendChannelSize),
			IsConnected:    false,
			ReconnectCount: 0,
		}
		ws.Clients[clientID] = client
		go client.connect()
	}
	return clientID, err
}

func (ws *WEBSOCKET_CLIENT_MANAGER) Shutdown(clientID string) error {
	var err error
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()
	client, exists := ws.Clients[clientID]
	if !exists {
		err = fmt.Errorf("websocket client not found: %s", clientID)
	} else {
		close(client.SendChan)
		delete(ws.Clients, clientID)
	}
	return err
}

func (ws *WEBSOCKET_CLIENT_MANAGER) Send(clientID string, message []byte) error {
	var err error
	var client *WEBSOCKET_CLIENT
	var exists bool
	ws.Mutex.Lock()
	client, exists = ws.Clients[clientID]
	ws.Mutex.Unlock()
	if exists {
		select {
		case client.SendChan <- message:
		default:
			err = fmt.Errorf("websocket client send channel is full or client is disconnected: %s", clientID)
		}
	} else {
		err = fmt.Errorf("websocket client not found: %s", clientID)
	}
	return err
}

func (c *WEBSOCKET_CLIENT) connect() {
	for {
		conn, _, err := c.Config.Dialer.Dial(c.Config.Url, c.Config.Headers)
		if err != nil {
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
		c.Mutex.Lock()
		c.Conn = conn
		c.IsConnected = true
		c.ReconnectCount = 0
		c.Mutex.Unlock()
		readDone := make(chan struct{})
		writeDone := make(chan struct{})
		go c.readLoop(readDone)
		go c.writeLoop(writeDone)
		select {
		case <-readDone:
			close(writeDone)
		case <-writeDone:
			close(readDone)
		}
		conn.Close()
		c.Mutex.Lock()
		c.IsConnected = false
		c.Mutex.Unlock()
		select {
		case <-c.SendChan:
			c.Mutex.Lock()
			c.ReconnectCount++
			reconnectInterval := c.Config.ReconnectInterval
			c.Mutex.Unlock()
			time.Sleep(reconnectInterval)
		default:
			return
		}
	}
}

func (c *WEBSOCKET_CLIENT) readLoop(done chan struct{}) {
	defer close(done)
	var shouldExit bool
	for !shouldExit {
		messageType, message, err := c.Conn.ReadMessage()
		if err != nil {
			shouldExit = true
		} else {
			if c.MessageHandler != nil {
				if err := c.MessageHandler(c.Conn, messageType, message); err != nil {
				}
			}
		}
	}
}

func (c *WEBSOCKET_CLIENT) writeLoop(done chan struct{}) {
	defer close(done)
	var shouldExit bool
	for !shouldExit {
		select {
		case message, ok := <-c.SendChan:
			if !ok {
				shouldExit = true
			} else if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				shouldExit = true
			}
		case <-done:
			shouldExit = true
		}
	}
}
