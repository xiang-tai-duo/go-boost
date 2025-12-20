// Package websocket
// File:        websocket.go
// Url:         https://github.com/xiang-tai-duo/go-bootstrap/blob/master/websocket/server/websocket.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: WebSocket server functionality for Go applications
// --------------------------------------------------------------------------------
package websocket

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

//goland:noinspection GoNameStartsWithPackageName,GoSnakeCaseUsage
const (
	MAX_MESSAGE_SIZE    = 512
	WRITE_WAIT          = 10 * time.Second
	PONG_WAIT           = 60 * time.Second
	PING_PERIOD         = (PONG_WAIT * 9) / 10
	MIN_PORT            = 1024
	MAX_PORT            = 65535
	MAX_ATTEMPTS        = 100
	BUFFER_SIZE         = 1024
	CHANNEL_BUFFER_SIZE = 256
	ANY_PORT            = 0
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
type (
	WEBSOCKET_CLIENT struct {
		conn *websocket.Conn
		send chan []byte
		uuid string
	}
	WEBSOCKET_DATA_HANDLER       func(websocket *websocket.Conn, messageType int, data []byte) error
	WEBSOCKET_DISCONNECT_HANDLER func(websocket *websocket.Conn) error
	WEBSOCKET_CONNECT_HANDLER    func(websocket *websocket.Conn, uuid string) error
	WEBSOCKET_ORIGIN_FILTER      interface {
		Allow(origin string) bool
	}
	WEBSOCKET_ORIGIN_MAP   map[string]bool
	WEBSOCKET_ORIGIN_REGEX struct {
		Pattern string
		Regex   *regexp.Regexp
	}
	WEBSOCKET_DATA struct {
		Pattern           string
		Handler           WEBSOCKET_DATA_HANDLER
		DisconnectHandler WEBSOCKET_DISCONNECT_HANDLER
		Filter            WEBSOCKET_ORIGIN_FILTER
	}
	WEB_SOCKET_SERVER struct {
		upgrader          websocket.Upgrader
		clients           map[string]*WEBSOCKET_CLIENT
		broadcast         chan []byte
		register          chan *WEBSOCKET_CLIENT
		unregister        chan *WEBSOCKET_CLIENT
		dataHandler       WEBSOCKET_DATA_HANDLER
		disconnectHandler WEBSOCKET_DISCONNECT_HANDLER
		connectHandler    WEBSOCKET_CONNECT_HANDLER
	}
)

func (m WEBSOCKET_ORIGIN_MAP) Allow(origin string) bool {
	return m[origin]
}

func (r *WEBSOCKET_ORIGIN_REGEX) Allow(origin string) bool {
	return r.Regex.MatchString(origin)
}

func NewWebSocketServer() *WEB_SOCKET_SERVER {
	return &WEB_SOCKET_SERVER{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  BUFFER_SIZE,
			WriteBufferSize: BUFFER_SIZE,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		broadcast:  make(chan []byte, CHANNEL_BUFFER_SIZE),
		register:   make(chan *WEBSOCKET_CLIENT),
		unregister: make(chan *WEBSOCKET_CLIENT),
		clients:    make(map[string]*WEBSOCKET_CLIENT),
	}
}

func (ws *WEB_SOCKET_SERVER) SetCheckOrigin(checkOrigin func(r *http.Request) bool) {
	ws.upgrader.CheckOrigin = checkOrigin
}

func (ws *WEB_SOCKET_SERVER) SetDataHandler(handler WEBSOCKET_DATA_HANDLER) {
	ws.dataHandler = handler
}

func (ws *WEB_SOCKET_SERVER) SetDisconnectHandler(handler WEBSOCKET_DISCONNECT_HANDLER) {
	ws.disconnectHandler = handler
}

func (ws *WEB_SOCKET_SERVER) SetConnectHandler(handler WEBSOCKET_CONNECT_HANDLER) {
	ws.connectHandler = handler
}

func (ws *WEB_SOCKET_SERVER) init() {
	for {
		select {
		case client := <-ws.register:
			ws.clients[client.uuid] = client
			if ws.connectHandler != nil {
				_ = ws.connectHandler(client.conn, client.uuid)
			}
		case client := <-ws.unregister:
			if _, ok := ws.clients[client.uuid]; ok {
				delete(ws.clients, client.uuid)
				close(client.send)
				if ws.disconnectHandler != nil {
					_ = ws.disconnectHandler(client.conn)
				}
			}
		case message := <-ws.broadcast:
			for _, client := range ws.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(ws.clients, client.uuid)
					if ws.disconnectHandler != nil {
						_ = ws.disconnectHandler(client.conn)
					}
				}
			}
		}
	}
}

func (ws *WEB_SOCKET_SERVER) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if conn, err := ws.upgrader.Upgrade(w, r, nil); err == nil {
		client := &WEBSOCKET_CLIENT{
			conn: conn,
			send: make(chan []byte, CHANNEL_BUFFER_SIZE),
			uuid: uuid.New().String(),
		}
		ws.register <- client
		go ws.readPump(client)
		go ws.writePump(client)
	}
}

func (ws *WEB_SOCKET_SERVER) Serve(w http.ResponseWriter, r *http.Request, handler WEBSOCKET_DATA_HANDLER, disconnectHandler WEBSOCKET_DISCONNECT_HANDLER) {
	ws.dataHandler = handler
	ws.disconnectHandler = disconnectHandler
	ws.ServeHTTP(w, r)
}

//goland:noinspection GoUnhandledErrorResult
func (ws *WEB_SOCKET_SERVER) readPump(client *WEBSOCKET_CLIENT) {
	defer func() {
		ws.unregister <- client
		client.conn.Close()
	}()
	client.conn.SetReadLimit(MAX_MESSAGE_SIZE)
	client.conn.SetReadDeadline(time.Now().Add(PONG_WAIT))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(PONG_WAIT))
		return nil
	})
	for {
		if messageType, message, err := client.conn.ReadMessage(); err == nil {
			if ws.dataHandler != nil {
				if err := ws.dataHandler(client.conn, messageType, message); err != nil {
					_ = client.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: %v", err)))
					break
				}
			} else {
				ws.broadcast <- message
			}
		} else {
			break
		}
	}
}

//goland:noinspection GoUnhandledErrorResult
func (ws *WEB_SOCKET_SERVER) writePump(client *WEBSOCKET_CLIENT) {
	ticker := time.NewTicker(PING_PERIOD)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(WRITE_WAIT))
			if ok {
				length := uint32(len(message))
				prefix := []byte{
					byte(length >> 24),
					byte(length >> 16),
					byte(length >> 8),
					byte(length),
				}
				fullMessage := append(prefix, message...)
				client.conn.WriteMessage(websocket.BinaryMessage, fullMessage)
			} else {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(WRITE_WAIT))
			client.conn.WriteMessage(websocket.PingMessage, nil)
		}
	}
}

func (ws *WEB_SOCKET_SERVER) Broadcast(message interface{}) {
	var data []byte
	switch v := message.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		data = []byte(fmt.Sprintf("%v", v))
	}
	select {
	case ws.broadcast <- data:
	default:
	}
}

func (ws *WEB_SOCKET_SERVER) SendToClient(clientUUID string, message interface{}) bool {
	result := false
	if client, ok := ws.clients[clientUUID]; ok {
		var data []byte
		switch v := message.(type) {
		case string:
			data = []byte(v)
		case []byte:
			data = v
		default:
			data = []byte(fmt.Sprintf("%v", v))
		}
		select {
		case client.send <- data:
			result = true
		default:
		}
	}
	return result
}

//goland:noinspection GoUnhandledErrorResult
func (ws *WEB_SOCKET_SERVER) LaunchAsync(port int) (int, error) {
	websocketPort := port
	err := error(nil)
	if websocketPort == ANY_PORT {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		i := 0
		for i = 0; i < MAX_ATTEMPTS; i++ {
			websocketPort = r.Intn(MAX_PORT-MIN_PORT+1) + MIN_PORT
			if listener, err := net.Listen("tcp", fmt.Sprintf(":%d", websocketPort)); err == nil {
				listener.Close()
				break
			}
		}
		if websocketPort == 0 {
			err = fmt.Errorf("failed to find available port after %d attempts", i)
		}
	}
	if websocketPort > 0 {
		go func(port int) {
			go ws.init()
			if listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port)); err == nil {
				if err = http.Serve(listener, ws); err != nil {
					log.Printf("WebSocket server error: %v", err)
				}
			}
		}(websocketPort)
	}
	return websocketPort, err
}

func (ws *WEB_SOCKET_SERVER) ClientCount() int {
	return len(ws.clients)
}
