// --------------------------------------------------------------------------------
// File:        serve.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: NewServe creates a wrapper for HTTP server operations, providing a
//              convenient way to create and manage HTTP servers in Go, including WebSocket support.
// --------------------------------------------------------------------------------

package boost

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// RequestHandler defines the signature for HTTP request handlers.
// w: Response writer for sending response
// r: HTTP request containing client data
// Returns: Error encountered during request handling
type RequestHandler func(w http.ResponseWriter, r *http.Request) error

// WebSocketDataHandler defines the signature for WebSocket data processing handlers.
// It's called when a message is received on a WebSocket connection.
// conn: WebSocket connection for sending responses back to the client
// messageType: Type of message received (websocket.TextMessage, websocket.BinaryMessage, etc.)
// data: The raw data received from the client
// Returns: Error if any processing error occurs, or nil if successful
type WebSocketDataHandler func(conn *websocket.Conn, messageType int, data []byte) error

// WebSocketHandler defines the signature for WebSocket connection handlers.
// w: Response writer for upgrading to WebSocket
// r: HTTP request containing WebSocket handshake data
// DEPRECATED: Use WebSocketDataHandler instead for easier data handling
type WebSocketHandler func(w http.ResponseWriter, r *http.Request)

// ROUTE represents an HTTP route with method, pattern, and handler.
type ROUTE struct {
	Method  string
	Pattern string
	Handler RequestHandler
}

// WebSocketOriginFilter defines the interface for WebSocket origin filtering.
type WebSocketOriginFilter interface {
	// Allow checks if the given origin is allowed.
	// origin: The Origin header value from the request
	// Returns: true if the origin is allowed, false otherwise
	Allow(origin string) bool
}

// WebSocketOriginMap implements WebSocketOriginFilter using a map of allowed origins.
type WebSocketOriginMap map[string]bool

// Allow checks if the origin is in the allowed origins map.
func (m WebSocketOriginMap) Allow(origin string) bool {
	return m[origin]
}

// WebSocketOriginRegex implements WebSocketOriginFilter using a regular expression.
type WebSocketOriginRegex struct {
	Pattern string         // The original regex pattern
	Regex   *regexp.Regexp // Compiled regex object
}

// NewWebSocketOriginRegex creates a new WebSocketOriginRegex from a regex pattern.
// pattern: The regex pattern to match allowed origins
// Returns: A new WebSocketOriginRegex or error if the pattern is invalid
func NewWebSocketOriginRegex(pattern string) (*WebSocketOriginRegex, error) {
	r, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &WebSocketOriginRegex{
		Pattern: pattern,
		Regex:   r,
	}, nil
}

// Allow checks if the origin matches the regex pattern.
func (r *WebSocketOriginRegex) Allow(origin string) bool {
	return r.Regex.MatchString(origin)
}

// WEBSOCKETROUTE represents a legacy WebSocket route with pattern, handler and origin filter.
type WEBSOCKETROUTE struct {
	Pattern string                // URL pattern to match
	Handler WebSocketHandler      // Legacy WebSocket connection handler
	Filter  WebSocketOriginFilter // Origin filter for this route
}

// WEBSOCKETDATA defines a WebSocket data route with pattern, data handler and origin filter.
// It's used by the new WebSocketData method which handles connection lifecycle internally.
type WEBSOCKETDATA struct {
	Pattern string                // URL pattern to match
	Handler WebSocketDataHandler  // Data handler for processing WebSocket messages
	Filter  WebSocketOriginFilter // Origin filter for this route
}

// NewServe creates a new SERVE instance with default configurations.
// Usage:
// serve := NewServe()
func NewServe() *SERVE {
	return &SERVE{
		server:            &http.Server{},
		staticDirectories: make(map[string]string),
	}
}

// SERVE represents an HTTP server with WebSocket support
// It manages routes, WebSocket connections, and static file serving
type SERVE struct {
	server              *http.Server
	routes              []ROUTE
	webSocketRoutes     []WEBSOCKETROUTE
	webSocketDataRoutes []WEBSOCKETDATA
	mutex               sync.Mutex
	isRunning           bool
	staticDirectories   map[string]string
}

// AddStaticDirectory maps a URL path to a local directory for serving static files.
// urlPath: The URL path to expose, e.g., "/static"
// directoryPath: The local directory path, e.g., "./public"
// Usage:
// serve.AddStaticDirectory("/static", "./public")
func (s *SERVE) AddStaticDirectory(urlPath string, directoryPath string) *SERVE {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.staticDirectories[urlPath] = directoryPath
	return s
}

// GetStaticDirectory retrieves the local directory path mapped to a given URL path.
// urlPath: The URL path to query
// Returns: The local directory path, or empty string if not found
// Usage:
// dir := serve.GetStaticDirectory("/static")
func (s *SERVE) GetStaticDirectory(urlPath string) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.staticDirectories[urlPath]
}

// Handler registers an HTTP request handler for a specific method and URL pattern.
// method: HTTP method, e.g., http.MethodGet, http.MethodPost, or "*" for all methods
// pattern: URL pattern to match, e.g., "/api/users" or "/api/users/*"
// handler: Function to handle matching requests
// Usage:
//
//	serve.Handler(http.MethodGet, "/api/users", func(w http.ResponseWriter, r *http.Request) error {
//	    w.Write([]byte("Hello, World!"))
//	    return nil
//	})
func (s *SERVE) Handler(method string, pattern string, handler RequestHandler) *SERVE {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.routes = append(s.routes, ROUTE{
		Method:  method,
		Pattern: pattern,
		Handler: handler,
	})
	return s
}

// IsRunning checks if the HTTP server is currently running.
// Returns: true if server is running, false otherwise
// Usage:
//
//	if serve.IsRunning() {
//	    fmt.Println("Server is running")
//	}
func (s *SERVE) IsRunning() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.isRunning
}

// Listen starts the HTTP server and begins listening for requests.
// addr: Address to listen on, e.g., ":8080"
// Returns: Error encountered during server startup or shutdown
// Usage:
// err := serve.Listen(":8080")
//
//	if err != nil && err != http.ErrServerClosed {
//	    log.Fatal(err)
//	}
func (s *SERVE) Listen(addr string) error {
	var err error
	s.mutex.Lock()
	if s.isRunning {
		s.mutex.Unlock()
		err = fmt.Errorf("server is already running")
	} else {
		s.isRunning = true
		s.server.Addr = addr
		s.mutex.Unlock()

		s.server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.handleRequest(w, r)
		})

		err = s.server.ListenAndServe()
	}
	return err
}

// ListenTLS starts the HTTPS server with TLS encryption.
// addr: Address to listen on, e.g., ":443"
// certFile: Path to TLS certificate file, leave empty to generate self-signed certificate
// keyFile: Path to TLS private key file, leave empty to generate self-signed certificate
// Returns: Error encountered during server startup or shutdown
// Usage:
// err := serve.ListenTLS(":443", "cert.pem", "key.pem")
// // Use self-signed certificate
// err := serve.ListenTLS(":443", "", "")
func (s *SERVE) ListenTLS(addr string, certFile string, keyFile string) error {
	var err error
	s.mutex.Lock()
	if s.isRunning {
		s.mutex.Unlock()
		err = fmt.Errorf("server is already running")
	} else {
		s.isRunning = true
		s.server.Addr = addr
		s.mutex.Unlock()

		s.server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.handleRequest(w, r)
		})

		if certFile != "" && keyFile != "" {
			err = s.server.ListenAndServeTLS(certFile, keyFile)
		} else {
			// Generate self-signed certificate
			var cert *tls.Certificate
			if cert, err = generateSelfSignedCertificate(); err == nil {
				s.server.TLSConfig = &tls.Config{
					Certificates: []tls.Certificate{*cert},
				}
				err = s.server.ListenAndServeTLS("", "")
			}
		}
	}
	return err
}

// RemoveStaticDirectory removes the mapping between a URL path and local directory.
// urlPath: The URL path to remove, e.g., "/static"
// Usage:
// serve.RemoveStaticDirectory("/static")
func (s *SERVE) RemoveStaticDirectory(urlPath string) *SERVE {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.staticDirectories, urlPath)
	return s
}

// Shutdown gracefully shuts down the HTTP server with a given context.
// ctx: Context to control shutdown timeout
// Returns: Error encountered during shutdown
// Usage:
// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// defer cancel()
// err := serve.Shutdown(ctx)
func (s *SERVE) Shutdown(ctx context.Context) error {
	var err error
	s.mutex.Lock()
	if !s.isRunning {
		s.mutex.Unlock()
		err = fmt.Errorf("server is not running")
	} else {
		s.mutex.Unlock()

		if err = s.server.Shutdown(ctx); err == nil {
			s.mutex.Lock()
			s.isRunning = false
			s.mutex.Unlock()
		}
	}
	return err
}

// ShutdownTLS gracefully shuts down the HTTPS server with a given context.
// ctx: Context to control shutdown timeout
// Returns: Error encountered during shutdown
// Usage:
// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// defer cancel()
// err := serve.ShutdownTLS(ctx)
func (s *SERVE) ShutdownTLS(ctx context.Context) error {
	var err error
	s.mutex.Lock()
	if !s.isRunning {
		s.mutex.Unlock()
		err = fmt.Errorf("server is not running")
	} else {
		s.mutex.Unlock()

		if err = s.server.Shutdown(ctx); err == nil {
			s.mutex.Lock()
			s.isRunning = false
			s.mutex.Unlock()
		}
	}
	return err
}

// ShutdownAll gracefully shuts down all servers (HTTP and HTTPS) with a given context.
// ctx: Context to control shutdown timeout
// Returns: Error encountered during shutdown
// Usage:
// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// defer cancel()
// err := serve.ShutdownAll(ctx)
func (s *SERVE) ShutdownAll(ctx context.Context) error {
	var err error
	s.mutex.Lock()
	if !s.isRunning {
		s.mutex.Unlock()
		err = fmt.Errorf("server is not running")
	} else {
		s.mutex.Unlock()

		if err = s.server.Shutdown(ctx); err == nil {
			s.mutex.Lock()
			s.isRunning = false
			s.mutex.Unlock()
		}
	}
	return err
}

func (s *SERVE) WebSocket(pattern string, handler WebSocketHandler) *SERVE {
	// Backward compatibility: call new WebSocket method with nil filter
	return s.WebSocket(pattern, nil, handler)
}

// WebSocket registers a WebSocket handler for a specific URL pattern with origin filtering.
// pattern: URL pattern to match, e.g., "/ws"
// filter: Origin filter - can be nil (allow all), map[string]bool, or regex string
// handler: Legacy WebSocket connection handler function
// DEPRECATED: Use WebSocketData instead for easier data handling
//
// Usage:
//
//	// Basic usage with nil filter (allow all origins)
//	serve.WebSocket("/ws", nil, func(w http.ResponseWriter, r *http.Request) {
//	    // Handle WebSocket connection here
//	})
//
//	// With origin filtering using map
//	serve.WebSocket("/ws", map[string]bool{
//	    "https://example.com": true,
//	    "http://localhost:3000": true,
//	}, func(w http.ResponseWriter, r *http.Request) {
//	    // Handle WebSocket connection here
//	})
//
//	// With origin filtering using regex
//	serve.WebSocket("/ws", `^https?://.*\.example\.com$`, func(w http.ResponseWriter, r *http.Request) {
//	    // Handle WebSocket connection here
//	})
func (s *SERVE) WebSocket(pattern string, filter interface{}, handler WebSocketHandler) *SERVE {
	var originFilter WebSocketOriginFilter
	if filter != nil {
		switch f := filter.(type) {
		case map[string]bool:
			originFilter = WebSocketOriginMap(f)
		case string:
			if r, err := NewWebSocketOriginRegex(f); err == nil {
				originFilter = r
			}
		}
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.webSocketRoutes = append(s.webSocketRoutes, WEBSOCKETROUTE{
		Pattern: pattern,
		Handler: handler,
		Filter:  originFilter,
	})
	return s
}

// WebSocketData registers a WebSocket data handler for a specific URL pattern.
// It handles the entire WebSocket connection lifecycle internally, so you DON'T need to:
// - Upgrade HTTP connections to WebSocket
// - Manage connection state or lifecycle
// - Read messages from clients
// - Handle connection closing and cleanup
//
// You only need to provide a dataHandler function that processes the ready-to-use raw data
// received from clients. The server handles all connection management automatically.
//
// pattern: URL pattern to match, e.g., "/ws"
// filter: Origin filter - can be:
//   - nil: Allow all origins (default, for development)
//   - map[string]bool: Allow specific origins listed in the map
//   - string: Regular expression pattern to match allowed origins
//
// dataHandler: Function called automatically when a message is received
//
//	conn: WebSocket connection for sending responses back to the client
//	messageType: Type of message received (websocket.TextMessage, websocket.BinaryMessage, etc.)
//	data: Raw, ready-to-use message data received from the client
//	Returns: Error if any, or nil if successful
//
// Complete Usage Examples:
//
// Example 1: Basic Echo Server
// This endpoint simply echoes back any message it receives
//
//	// Create a simple echo server
//	serve.WebSocketData("/echo", nil, func(conn *websocket.Conn, messageType int, data []byte) error {
//	    fmt.Printf("[Echo] Received: %s\n", string(data))
//	    // Just echo the message back to the client
//	    return conn.WriteMessage(messageType, data)
//	})
//
// Example 2: JSON API Server with Origin Filtering
// This endpoint processes JSON messages with different actions (ping, time, greet)
//
//	// Create a JSON API server
//	serve.WebSocketData("/api", map[string]bool{
//	    "http://localhost:3000": true,
//	    "https://example.com":   true,
//	}, func(conn *websocket.Conn, messageType int, data []byte) error {
//	    // Define request and response structs
//	    type Request struct {
//	        Action string      `json:"action"`
//	        Data   interface{} `json:"data"`
//	    }
//
//	    type Response struct {
//	        Status  string      `json:"status"`
//	        Message string      `json:"message"`
//	        Data    interface{} `json:"data,omitempty"`
//	    }
//
//	    // Parse incoming JSON
//	    var req Request
//	    if err := json.Unmarshal(data, &req); err != nil {
//	        resp := Response{
//	            Status:  "error",
//	            Message: "Invalid JSON format",
//	        }
//	        respData, _ := json.Marshal(resp)
//	        return conn.WriteMessage(messageType, respData)
//	    }
//
//	    // Process based on action
//	    var resp Response
//	    switch req.Action {
//	    case "ping":
//	        resp = Response{
//	            Status:  "success",
//	            Message: "Pong response",
//	            Data: map[string]string{
//	                "timestamp": time.Now().Format(time.RFC3339),
//	            },
//	        }
//	    case "time":
//	        resp = Response{
//	            Status:  "success",
//	            Message: "Current time",
//	            Data: map[string]string{
//	                "time": time.Now().Format(time.RFC3339),
//	            },
//	        }
//	    case "greet":
//	        name, ok := req.Data.(string)
//	        if !ok {
//	            name = "Guest"
//	        }
//	        resp = Response{
//	            Status:  "success",
//	            Message: "Hello " + name + "!",
//	        }
//	    default:
//	        resp = Response{
//	            Status:  "error",
//	            Message: "Unknown action: " + req.Action,
//	        }
//	    }
//
//	    // Send response
//	    respData, _ := json.Marshal(resp)
//	    return conn.WriteMessage(messageType, respData)
//	})
//
// Example 3: Binary Data Server
// This endpoint handles binary messages with prefix echo
//
//	// Create a binary data server
//	serve.WebSocketData("/binary", nil, func(conn *websocket.Conn, messageType int, data []byte) error {
//	    if messageType == websocket.BinaryMessage {
//	        fmt.Printf("[Binary] Received %d bytes\n", len(data))
//	        // Echo back binary data with a prefix
//	        response := append([]byte("ECHO: "), data...)
//	        return conn.WriteMessage(messageType, response)
//	    }
//	    return nil
//	})
//
// Example 4: Simple Chat Server
// This example shows a chat server that broadcasts messages to all clients
//
//	// Create a map to store all connected clients
//	clients := make(map[*websocket.Conn]bool)
//	clientsMutex := sync.Mutex{}
//
//	// Create a chat server
//	serve.WebSocketData("/chat", nil, func(conn *websocket.Conn, messageType int, data []byte) error {
//	    // Add client to map when connected
//	    clientsMutex.Lock()
//	    clients[conn] = true
//	    clientsMutex.Unlock()
//
//	    // Broadcast message to all clients
//	    clientsMutex.Lock()
//	    for client := range clients {
//	        if client != conn {
//	            client.WriteMessage(messageType, data)
//	        }
//	    }
//	    clientsMutex.Unlock()
//
//	    return nil
//	})
//
// Example 5: Authentication Required WebSocket
// This example shows how to add authentication check before processing messages
//
//	// Create a secure WebSocket endpoint
//	serve.WebSocketData("/secure", nil, func(conn *websocket.Conn, messageType int, data []byte) error {
//	    // Parse authentication request
//	    type AuthRequest struct {
//	        Token string `json:"token"`
//	    }
//
//	    var authReq AuthRequest
//	    if err := json.Unmarshal(data, &authReq); err != nil {
//	        conn.WriteMessage(messageType, []byte(`{"status":"error","message":"Invalid auth format"}`))
//	        return err
//	    }
//
//	    // Validate token (replace with your own auth logic)
//	    if authReq.Token != "valid-token" {
//	        conn.WriteMessage(messageType, []byte(`{"status":"error","message":"Invalid token"}`))
//	        return fmt.Errorf("invalid token")
//	    }
//
//	    // Authentication successful
//	    conn.WriteMessage(messageType, []byte(`{"status":"ok","message":"Authenticated"}`))
//	    return nil
//	})
//
// Example 6: File Transfer Server
// This example shows how to handle file transfer via WebSocket
//
//	// Create a file transfer server
//	serve.WebSocketData("/file", nil, func(conn *websocket.Conn, messageType int, data []byte) error {
//	    if messageType == websocket.BinaryMessage {
//	        // Save binary data to a file
//	        fileName := fmt.Sprintf("file_%d.dat", time.Now().Unix())
//	        if err := os.WriteFile(fileName, data, 0644); err != nil {
//	            conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("{\"status\":\"error\",\"message\":\"Failed to save file: %v\"}", err)))
//	            return err
//	        }
//	        conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("{\"status\":\"success\",\"message\":\"File saved as %s\"}", fileName)))
//	    }
//	    return nil
//	})
func (s *SERVE) WebSocketData(pattern string, filter interface{}, dataHandler WebSocketDataHandler) *SERVE {
	var originFilter WebSocketOriginFilter
	if filter != nil {
		switch f := filter.(type) {
		case map[string]bool:
			originFilter = WebSocketOriginMap(f)
		case string:
			if r, err := NewWebSocketOriginRegex(f); err == nil {
				originFilter = r
			}
		}
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.webSocketDataRoutes = append(s.webSocketDataRoutes, WEBSOCKETDATA{
		Pattern: pattern,
		Handler: dataHandler,
		Filter:  originFilter,
	})
	return s
}

func generateSelfSignedCertificate() (*tls.Certificate, error) {
	var cert *tls.Certificate
	var err error

	var priv *rsa.PrivateKey
	if priv, err = rsa.GenerateKey(rand.Reader, 2048); err == nil {
		validFrom := time.Now()
		validTo := validFrom.Add(365 * 24 * time.Hour)

		serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
		var serialNumber *big.Int
		if serialNumber, err = rand.Int(rand.Reader, serialNumberLimit); err == nil {
			certTemplate := x509.Certificate{
				SerialNumber: serialNumber,
				Subject: pkix.Name{
					Organization: []string{"Go Solution SDK"},
				},
				NotBefore:             validFrom,
				NotAfter:              validTo,
				KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
				ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
				BasicConstraintsValid: true,
				IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
			}

			var certDER []byte
			if certDER, err = x509.CreateCertificate(rand.Reader, &certTemplate, &certTemplate, &priv.PublicKey, priv); err == nil {
				var privDER []byte
				if privDER, err = x509.MarshalPKCS8PrivateKey(priv); err == nil {
					certPEM := new(bytes.Buffer)
					pem.Encode(certPEM, &pem.Block{Type: "CERTIFICATE", Bytes: certDER})

					privPEM := new(bytes.Buffer)
					pem.Encode(privPEM, &pem.Block{Type: "PRIVATE KEY", Bytes: privDER})

					var tlsCert tls.Certificate
					if tlsCert, err = tls.X509KeyPair(certPEM.Bytes(), privPEM.Bytes()); err == nil {
						cert = &tlsCert
					}
				}
			}
		}
	}

	return cert, err
}

// handleRequest processes incoming HTTP requests
// It handles:
// - WebSocket upgrade requests (legacy and new WebSocketData routes)
// - HTTP route requests
// - Static file serving
func (s *SERVE) handleRequest(w http.ResponseWriter, r *http.Request) {
	s.mutex.Lock()
	routes := make([]ROUTE, len(s.routes))
	copy(routes, s.routes)
	webSocketRoutes := make([]WEBSOCKETROUTE, len(s.webSocketRoutes))
	copy(webSocketRoutes, s.webSocketRoutes)
	webSocketDataRoutes := make([]WEBSOCKETDATA, len(s.webSocketDataRoutes))
	copy(webSocketDataRoutes, s.webSocketDataRoutes)
	s.mutex.Unlock()

	var handled bool
	if strings.ToLower(r.Header.Get("Upgrade")) == "websocket" {
		// Handle legacy WebSocket routes
		for _, wsRoute := range webSocketRoutes {
			if s.matchPath(wsRoute.Pattern, r.URL.Path) {
				wsRoute.Handler(w, r)
				handled = true
				break
			}
		}

		// Handle new WebSocketData routes
		if !handled {
			for _, wsDataRoute := range webSocketDataRoutes {
				if s.matchPath(wsDataRoute.Pattern, r.URL.Path) {
					// Create WebSocket upgrader
					upgrader := websocket.Upgrader{
						CheckOrigin: func(r *http.Request) bool {
							origin := r.Header.Get("Origin")
							if wsDataRoute.Filter == nil {
								return true // Allow all origins
							}
							return wsDataRoute.Filter.Allow(origin)
						},
					}

					// Upgrade connection
					conn, err := upgrader.Upgrade(w, r, nil)
					if err == nil {
						// Handle connection lifecycle
						defer func() {
							_ = conn.Close()
						}()

						// Continuously read messages
						for {
							messageType, message, err := conn.ReadMessage()
							if err != nil {
								break // Connection closed or error occurred
							}

							// Call user data handler
							if err := wsDataRoute.Handler(conn, messageType, message); err != nil {
								// If handler returns error, send error message and close connection
								_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: %v", err)))
								break
							}
						}
					}
					handled = true
					break
				}
			}
		}
	}

	if !handled {
		for _, route := range routes {
			if (route.Method == "*" || route.Method == r.Method) && s.matchPath(route.Pattern, r.URL.Path) {
				if err := route.Handler(w, r); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				handled = true
				break
			}
		}
	}

	if !handled {
		if !s.serveStatic(w, r) {
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	}
}

func (s *SERVE) matchPath(pattern string, path string) bool {
	var matched bool
	if pattern == path {
		matched = true
	} else if strings.HasSuffix(pattern, "/*") {
		prefix := strings.TrimSuffix(pattern, "/*")
		if prefix == "" || strings.HasPrefix(path, prefix+"/") || path == prefix {
			matched = true
		}
	} else {
		// Special handling for root path to ensure only matching root path
		if pattern == "/" || path == "/" {
			return false
		}

		patternParts := strings.Split(pattern, "/")
		pathParts := strings.Split(path, "/")

		if len(patternParts) != len(pathParts) {
			return false
		}

		matched = true
		for i, patternPart := range patternParts {
			if patternPart == "" {
				continue
			}

			if strings.HasPrefix(patternPart, "{") && strings.HasSuffix(patternPart, "}") {
				// This is a path parameter, skip matching
				continue
			}

			if patternPart != pathParts[i] {
				matched = false
				break
			}
		}
	}
	return matched
}

func (s *SERVE) serveStatic(w http.ResponseWriter, r *http.Request) bool {
	var served bool
	s.mutex.Lock()
	staticDirectories := make(map[string]string)
	for k, v := range s.staticDirectories {
		staticDirectories[k] = v
	}
	s.mutex.Unlock()
	for urlPath, directoryPath := range staticDirectories {
		if strings.HasPrefix(r.URL.Path, urlPath) {
			filePath := filepath.Join(directoryPath, strings.TrimPrefix(r.URL.Path, urlPath))

			if _, err := os.Stat(filePath); !os.IsNotExist(err) {
				http.ServeFile(w, r, filePath)
				served = true
				break
			}
		}
	}
	return served
}
