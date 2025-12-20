// --------------------------------------------------------------------------------
// File:        serve.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: HTTP server with WebSocket support for Go applications
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

// Constants for WebSocket upgrade
const (
	// HeaderUpgrade is the HTTP header name for WebSocket upgrade
	HeaderUpgrade = "Upgrade"
	// WebSocketUpgrade is the WebSocket upgrade value
	WebSocketUpgrade = "websocket"
)

//goland:noinspection SpellCheckingInspection
type (
	// RequestHandler defines the signature for HTTP request handlers
	RequestHandler func(w http.ResponseWriter, r *http.Request) error

	// WebSocketDataHandler defines the signature for WebSocket data processing handlers
	WebSocketDataHandler func(ws *websocket.Conn, messageType int, data []byte) error

	// Route represents an HTTP route with method, pattern, and handler
	Route struct {
		Method  string
		Pattern string
		Handler RequestHandler
	}

	// WebSocketOriginFilter defines the interface for WebSocket origin filtering
	WebSocketOriginFilter interface {
		Allow(origin string) bool
	}

	// WebSocketOriginMap implements WebSocketOriginFilter using a map of allowed origins
	WebSocketOriginMap map[string]bool

	// WebSocketOriginRegex implements WebSocketOriginFilter using a regular expression
	WebSocketOriginRegex struct {
		Pattern string
		Regex   *regexp.Regexp
	}

	// WebSocketData defines a WebSocket data route with pattern, data handler and origin filter
	WebSocketData struct {
		Pattern string
		Handler WebSocketDataHandler
		Filter  WebSocketOriginFilter
	}

	// Serve represents an HTTP server with WebSocket support
	Serve struct {
		Server              *http.Server
		Routes              []Route
		WebSocketDataRoutes []WebSocketData
		WebSocketClients    map[string]*WEBSOCKET_CLIENT
		Mutex               sync.Mutex
		isRunning           bool
		StaticDirectories   map[string]string
		websocket           *WEBSOCKET_CLIENT_MANAGER
	}
)

// Allow checks if the origin is in the allowed origins map
func (m WebSocketOriginMap) Allow(origin string) bool {
	return m[origin]
}

// Allow checks if the origin matches the regex pattern
func (r *WebSocketOriginRegex) Allow(origin string) bool {
	return r.Regex.MatchString(origin)
}

// AddStaticDirectory maps a URL path to a local directory for serving static files
// urlPath: The URL path to expose, e.g., "/static"
// directoryPath: The local directory path, e.g., "./public"
// Returns: The Serve instance for method chaining
func (s *Serve) AddStaticDirectory(urlPath string, directoryPath string) *Serve {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.StaticDirectories[urlPath] = directoryPath
	return s
}

// GetStaticDirectory retrieves the local directory path mapped to a given URL path
// urlPath: The URL path to query
// Returns: The local directory path, or empty string if not found
func (s *Serve) GetStaticDirectory(urlPath string) string {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	return s.StaticDirectories[urlPath]
}

// IsRunning checks if the HTTP server is currently running
// Returns: true if server is running, false otherwise
func (s *Serve) IsRunning() bool {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	return s.isRunning
}

// GetAvailablePort returns a random available TCP port greater than 1024
// Returns: A random available port number and any error encountered
func (s *Serve) GetAvailablePort() (int, error) {
	// Create a listener with random available port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()

	// Get the actual port assigned by the system
	port := listener.Addr().(*net.TCPAddr).Port

	// Ensure port is greater than 1024
	if port > 1024 {
		return port, nil
	}

	// If port is 1024 or less, try again (should be very rare)
	return s.GetAvailablePort()
}

// listen starts the HTTP server and begins listening for requests
// addr: Address to listen on, e.g., ":8080"
// Returns: Error encountered during server startup or shutdown
func (s *Serve) listen(addr string) error {
	s.Mutex.Lock()
	if s.isRunning {
		s.Mutex.Unlock()
		return fmt.Errorf("server is already running")
	}

	s.isRunning = true
	s.Server.Addr = addr
	s.Mutex.Unlock()

	s.Server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.handleRequest(w, r)
	})

	return s.Server.ListenAndServe()
}

// listenTLS starts the HTTPS server with TLS encryption
// addr: Address to listen on, e.g., ":443"
// certFile: Path to TLS certificate file, leave empty to generate self-signed certificate
// keyFile: Path to TLS private key file, leave empty to generate self-signed certificate
// Returns: Error encountered during server startup or shutdown
func (s *Serve) listenTLS(addr string, certFile string, keyFile string) error {
	s.Mutex.Lock()
	if s.isRunning {
		s.Mutex.Unlock()
		return fmt.Errorf("server is already running")
	}

	s.isRunning = true
	s.Server.Addr = addr
	s.Mutex.Unlock()

	s.Server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.handleRequest(w, r)
	})

	if certFile != "" && keyFile != "" {
		return s.Server.ListenAndServeTLS(certFile, keyFile)
	}

	certificate, err := generateSelfSignedCertificate()
	if err != nil {
		return err
	}

	if certificate == nil {
		return fmt.Errorf("failed to generate self-signed certificate: certificate is nil")
	}

	s.Server.TLSConfig = &tls.Config{
		Certificates: []tls.Certificate{*certificate},
	}

	return s.Server.ListenAndServeTLS("", "")
}

// NewServe creates a new Serve instance with default configurations
// Returns: A new Serve instance
func NewServe() *Serve {
	return &Serve{
		Server: &http.Server{
			Addr:         ":8080",          // Default listen port
			ReadTimeout:  15 * time.Second, // Default read timeout
			WriteTimeout: 15 * time.Second, // Default write timeout
			IdleTimeout:  60 * time.Second, // Default idle timeout
		},
		StaticDirectories: make(map[string]string),
		WebSocketClients:  make(map[string]*WEBSOCKET_CLIENT),
		websocket:         NewWebSocket(),
	}
}

// NewWebSocketOriginRegex creates a new WebSocketOriginRegex from a regex pattern
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

// On registers an HTTP request handler for a specific method and URL pattern
// method: HTTP method, e.g., http.MethodGet, http.MethodPost, or "*" for all methods
// pattern: URL pattern to match, e.g., "/api/users" or "/api/users/*"
// handler: Function to handle matching requests
// Returns: The Serve instance for method chaining
func (s *Serve) On(method string, pattern string, handler RequestHandler) *Serve {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Routes = append(s.Routes, Route{
		Method:  method,
		Pattern: pattern,
		Handler: handler,
	})

	// Auto-start server if not already running
	if !s.isRunning {
		go func() {
			if err := s.listen(s.Server.Addr); err != nil {
				fmt.Printf("Server error: %v\n", err)
			}
		}()
	}

	return s
}

// RemoveStaticDirectory removes the mapping between a URL path and local directory
// urlPath: The URL path to remove, e.g., "/static"
// Returns: The Serve instance for method chaining
func (s *Serve) RemoveStaticDirectory(urlPath string) *Serve {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	delete(s.StaticDirectories, urlPath)
	return s
}

// Shutdown gracefully shuts down all servers (HTTP and HTTPS) with a given context
// ctx: Context to control shutdown timeout
// Returns: Error encountered during shutdown
func (s *Serve) Shutdown(ctx context.Context) error {
	s.Mutex.Lock()
	if !s.isRunning {
		s.Mutex.Unlock()
		return fmt.Errorf("server is not running")
	}
	s.Mutex.Unlock()

	err := s.Server.Shutdown(ctx)
	if err == nil {
		s.Mutex.Lock()
		s.isRunning = false
		s.Mutex.Unlock()
	}
	return err
}

// shutdownAll gracefully shuts down all servers (HTTP and HTTPS) with a given context
// ctx: Context to control shutdown timeout
// Returns: Error encountered during shutdown
func (s *Serve) shutdownAll(ctx context.Context) error {
	s.Mutex.Lock()
	if !s.isRunning {
		s.Mutex.Unlock()
		return fmt.Errorf("server is not running")
	}
	s.Mutex.Unlock()

	err := s.Server.Shutdown(ctx)
	if err == nil {
		s.Mutex.Lock()
		s.isRunning = false
		s.Mutex.Unlock()
	}
	return err
}

// shutdownTLS gracefully shuts down the HTTPS server with a given context
// ctx: Context to control shutdown timeout
// Returns: Error encountered during shutdown
func (s *Serve) shutdownTLS(ctx context.Context) error {
	s.Mutex.Lock()
	if !s.isRunning {
		s.Mutex.Unlock()
		return fmt.Errorf("server is not running")
	}
	s.Mutex.Unlock()

	err := s.Server.Shutdown(ctx)
	if err == nil {
		s.Mutex.Lock()
		s.isRunning = false
		s.Mutex.Unlock()
	}
	return err
}

// OnWebSocket registers a WebSocket data handler for a specific URL pattern
// It handles the entire WebSocket connection lifecycle internally
// pattern: URL pattern to match, e.g., "/ws"
// filter: Origin filter - can be nil (allow all), string, or string array
// dataHandler: Function to process received WebSocket data
// Returns: The Serve instance for method chaining
func (s *Serve) OnWebSocket(pattern string, filter interface{}, dataHandler WebSocketDataHandler) *Serve {
	var originFilter WebSocketOriginFilter
	if filter != nil {
		switch f := filter.(type) {
		case string:
			// Single origin filter
			originMap := make(WebSocketOriginMap)
			originMap[f] = true
			originFilter = originMap
		case []string:
			// Multiple origin filters
			originMap := make(WebSocketOriginMap)
			for _, origin := range f {
				originMap[origin] = true
			}
			originFilter = originMap
		}
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.WebSocketDataRoutes = append(s.WebSocketDataRoutes, WebSocketData{
		Pattern: pattern,
		Handler: dataHandler,
		Filter:  originFilter,
	})

	// Auto-start server if not already running
	if !s.isRunning {
		go func() {
			if err := s.listen(s.Server.Addr); err != nil {
				fmt.Printf("Server error: %v\n", err)
			}
		}()
	}

	return s
}

// generateSelfSignedCertificate generates a self-signed TLS certificate
func generateSelfSignedCertificate() (*tls.Certificate, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	validFromTime := time.Now()
	validToTime := validFromTime.Add(365 * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}

	certificateTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Go Solution SDK"},
		},
		NotBefore:             validFromTime,
		NotAfter:              validToTime,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
	}

	certificateDER, err := x509.CreateCertificate(rand.Reader, &certificateTemplate, &certificateTemplate, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, err
	}

	privateKeyDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	certificatePEM := new(bytes.Buffer)
	_ = pem.Encode(certificatePEM, &pem.Block{Type: "CERTIFICATE", Bytes: certificateDER})

	privateKeyPEM := new(bytes.Buffer)
	_ = pem.Encode(privateKeyPEM, &pem.Block{Type: "PRIVATE KEY", Bytes: privateKeyDER})

	tlsCertificate, err := tls.X509KeyPair(certificatePEM.Bytes(), privateKeyPEM.Bytes())
	if err != nil {
		return nil, err
	}

	return &tlsCertificate, nil
}

// handleHTTPRequest processes incoming HTTP requests
func (s *Serve) handleHTTPRequest(w http.ResponseWriter, r *http.Request, routes []Route) bool {
	var handled bool
	// Try to match HTTP routes first
	for _, route := range routes {
		if (route.Method == "*" || route.Method == r.Method) && s.matchPath(route.Pattern, r.URL.Path) {
			if err := route.Handler(w, r); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			handled = true
			break
		}
	}

	// Try to serve static files if no route matched
	if !handled {
		handled = s.serveStatic(w, r)
	}

	return handled
}

// handleRequest processes incoming HTTP requests
func (s *Serve) handleRequest(w http.ResponseWriter, r *http.Request) {
	s.Mutex.Lock()
	routes := make([]Route, len(s.Routes))
	copy(routes, s.Routes)
	webSocketDataRoutes := make([]WebSocketData, len(s.WebSocketDataRoutes))
	copy(webSocketDataRoutes, s.WebSocketDataRoutes)
	s.Mutex.Unlock()

	var handled bool
	if strings.ToLower(r.Header.Get(HeaderUpgrade)) == WebSocketUpgrade {
		// Handle WebSocket requests
		handled = s.handleWebSocketRequest(w, r, webSocketDataRoutes)
	} else {
		// Handle HTTP requests
		handled = s.handleHTTPRequest(w, r, routes)
	}

	if !handled {
		// If not handled by any route, return 404
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

// handleWebSocketRequest processes incoming WebSocket upgrade requests
func (s *Serve) handleWebSocketRequest(w http.ResponseWriter, r *http.Request, webSocketDataRoutes []WebSocketData) bool {
	var handled bool
	// Try to match WebSocket data routes
	for _, wsDataRoute := range webSocketDataRoutes {
		if s.matchPath(wsDataRoute.Pattern, r.URL.Path) {
			upgrader := websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					origin := r.Header.Get("Origin")
					if wsDataRoute.Filter == nil {
						return true
					}
					return wsDataRoute.Filter.Allow(origin)
				},
			}

			var conn *websocket.Conn
			var err error
			if conn, err = upgrader.Upgrade(w, r, nil); err == nil {
				// Handle WebSocket connection in a separate goroutine
				go func(conn *websocket.Conn, route WebSocketData) {
					defer func() {
						_ = conn.Close()
					}()

					for {
						messageType, message, err := conn.ReadMessage()
						if err != nil {
							break
						}

						if err := route.Handler(conn, messageType, message); err != nil {
							_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: %v", err)))
							break
						}
					}
				}(conn, wsDataRoute)
				handled = true
			}
			break
		}
	}

	return handled
}

// matchPath checks if a URL path matches a given pattern
func (s *Serve) matchPath(pattern string, path string) bool {
	var matched bool
	if pattern == path {
		matched = true
	} else if strings.HasSuffix(pattern, "/*") {
		prefix := strings.TrimSuffix(pattern, "/*")
		if prefix == "" || strings.HasPrefix(path, prefix+"/") || path == prefix {
			matched = true
		}
	} else {
		if pattern != "/" && path != "/" {
			patternParts := strings.Split(pattern, "/")
			pathParts := strings.Split(path, "/")

			if len(patternParts) == len(pathParts) {
				matched = true
				for i, patternPart := range patternParts {
					if patternPart == "" {
						continue
					}

					if strings.HasPrefix(patternPart, "{") && strings.HasSuffix(patternPart, "}") {
						continue
					}

					if patternPart != pathParts[i] {
						matched = false
						break
					}
				}
			}
		}
	}
	return matched
}

// serveStatic serves static files from configured directories
func (s *Serve) serveStatic(w http.ResponseWriter, r *http.Request) bool {
	var served bool
	s.Mutex.Lock()
	staticDirectories := make(map[string]string)
	for k, v := range s.StaticDirectories {
		staticDirectories[k] = v
	}
	s.Mutex.Unlock()
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
