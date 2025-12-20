// Package websocket_client
// File:        websocket.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/websocket/websocket.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: WebSocket client functionality for Go applications
// --------------------------------------------------------------------------------
package websocket

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection SpellCheckingInspection,GoSnakeCaseUsage,GoNameStartsWithPackageName
type (
	WEBSOCKET_CLIENT struct {
		Config         WEBSOCKET_CLIENT_CONFIG
		Conn           *websocket.Conn
		ErrorHandler   WEBSOCKET_CLIENT_ERROR_HANDLER
		IsConnected    bool
		MessageHandler WEBSOCKET_CLIENT_HANDLER
		Mutex          sync.Mutex
		ReconnectCount int
		SendChan       chan []byte
	}

	WEBSOCKET_CLIENT_CONFIG struct {
		Dialer               *websocket.Dialer
		HandshakeTimeout     time.Duration
		Headers              http.Header
		MaxReconnectAttempts int
		Origin               string
		ReconnectInterval    time.Duration
		SendChannelSize      int
		Url                  string
	}

	WEBSOCKET_CLIENT_ERROR_HANDLER func(err error)

	WEBSOCKET_CLIENT_HANDLER func(conn *websocket.Conn, messageType int, data []byte) error

	WEBSOCKET_CLIENT_MANAGER struct {
		Clients       map[string]*WEBSOCKET_CLIENT
		DefaultConfig WEBSOCKET_CLIENT_CONFIG
		Mutex         sync.Mutex
	}
)

const (
	CLIENT_ID_FORMAT               = "websocket-client-%d"
	DEFAULT_HANDSHAKE_TIMEOUT      = 10 * time.Second
	DEFAULT_MAX_RECONNECT_ATTEMPTS = 5
	DEFAULT_RECONNECT_INTERVAL     = 5 * time.Second
	DEFAULT_SEND_CHANNEL_SIZE      = 100
	MODULE_NAME_WEBSOCKET          = "websocket"
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_WEBSOCKET, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_WEBSOCKET, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_WEBSOCKET, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedExportedFunction
func New() *WEBSOCKET_CLIENT_MANAGER {
	return &WEBSOCKET_CLIENT_MANAGER{
		Clients: make(map[string]*WEBSOCKET_CLIENT),
		DefaultConfig: WEBSOCKET_CLIENT_CONFIG{
			HandshakeTimeout:     DEFAULT_HANDSHAKE_TIMEOUT,
			ReconnectInterval:    DEFAULT_RECONNECT_INTERVAL,
			MaxReconnectAttempts: DEFAULT_MAX_RECONNECT_ATTEMPTS,
			SendChannelSize:      DEFAULT_SEND_CHANNEL_SIZE,
		},
	}
}

func (ws *WEBSOCKET_CLIENT_MANAGER) Connect(uri string, messageHandler WEBSOCKET_CLIENT_HANDLER) (string, error) {
	return ws.ConnectEx(uri, "", nil, nil, 0, 0, 0, 0, messageHandler)
}

func (ws *WEBSOCKET_CLIENT_MANAGER) ConnectEx(uri string, origin string, headers http.Header, dialer *websocket.Dialer, handshakeTimeout time.Duration, reconnectInterval time.Duration, maxReconnectAttempts int, sendChannelSize int, messageHandler WEBSOCKET_CLIENT_HANDLER) (string, error) {
	clientID := ""
	err := error(nil)
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
		clientID = fmt.Sprintf(CLIENT_ID_FORMAT, len(ws.Clients)+1)
		config := WEBSOCKET_CLIENT_CONFIG{
			Dialer:               dialer,
			Headers:              headers,
			MaxReconnectAttempts: maxReconnectAttempts,
			Origin:               origin,
			ReconnectInterval:    reconnectInterval,
			SendChannelSize:      sendChannelSize,
			Url:                  uri,
		}
		client := &WEBSOCKET_CLIENT{
			Config:         config,
			IsConnected:    false,
			MessageHandler: messageHandler,
			ReconnectCount: 0,
			SendChan:       make(chan []byte, sendChannelSize),
		}
		ws.Clients[clientID] = client
		go client.connect()
	}
	return clientID, err
}

func (ws *WEBSOCKET_CLIENT_MANAGER) Send(clientID string, message []byte) error {
	err := error(nil)
	ws.Mutex.Lock()
	client, exists := ws.Clients[clientID]
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

func (ws *WEBSOCKET_CLIENT_MANAGER) Shutdown(clientID string) error {
	err := error(nil)
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()
	if client, exists := ws.Clients[clientID]; exists {
		close(client.SendChan)
		delete(ws.Clients, clientID)
	} else {
		err = fmt.Errorf("websocket client not found: %s", clientID)
	}
	return err
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_WEBSOCKET, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnhandledErrorResult
func (c *WEBSOCKET_CLIENT) connect() {
	shouldExit := false
	for !shouldExit {
		conn, _, err := c.Config.Dialer.Dial(c.Config.Url, c.Config.Headers)
		if err == nil {
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
				shouldExit = true
			}
		} else {
			c.Mutex.Lock()
			if c.ReconnectCount < c.Config.MaxReconnectAttempts {
				c.ReconnectCount++
				reconnectInterval := c.Config.ReconnectInterval
				c.Mutex.Unlock()
				time.Sleep(reconnectInterval)
			} else {
				c.Mutex.Unlock()
				shouldExit = true
			}
		}
	}
}

func (c *WEBSOCKET_CLIENT) readLoop(done chan struct{}) {
	defer close(done)
	shouldExit := false
	for !shouldExit {
		messageType, message, err := c.Conn.ReadMessage()
		if err == nil {
			if c.MessageHandler != nil {
				_ = c.MessageHandler(c.Conn, messageType, message)
			}
		} else {
			shouldExit = true
		}
	}
}

func (c *WEBSOCKET_CLIENT) writeLoop(done chan struct{}) {
	defer close(done)
	shouldExit := false
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
