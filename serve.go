// --------------------------------------------------------------------------------
// File:        serve.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: ServeBuilder is a wrapper for HTTP server operations, providing a
//              convenient way to create and manage HTTP servers in Go.
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
	"strings"
	"sync"
	"time"
)

// RequestHandler defines the signature for HTTP request handlers.
// w: Response writer for sending response
// r: HTTP request containing client data
// Returns: Error encountered during request handling
type RequestHandler func(w http.ResponseWriter, r *http.Request) error

// WebSocketHandler defines the signature for WebSocket connection handlers.
// w: Response writer for upgrading to WebSocket
// r: HTTP request containing WebSocket handshake data
type WebSocketHandler func(w http.ResponseWriter, r *http.Request)

// ROUTE represents an HTTP route with method, pattern, and handler.
type ROUTE struct {
	Method  string
	Pattern string
	Handler RequestHandler
}

// WEBSOCKETROUTE represents a WebSocket route with pattern and handler.
type WEBSOCKETROUTE struct {
	Pattern string
	Handler WebSocketHandler
}

// NewServeBuilder creates a new SERVE instance with default configurations.
// Usage:
// builder := NewServeBuilder()
func NewServeBuilder() *SERVE {
	return &SERVE{
		server:            &http.Server{},
		staticDirectories: make(map[string]string),
	}
}

type SERVE struct {
	server            *http.Server
	routes            []ROUTE
	webSocketRoutes   []WEBSOCKETROUTE
	mutex             sync.Mutex
	isRunning         bool
	staticDirectories map[string]string
}

// AddStaticDirectory maps a URL path to a local directory for serving static files.
// urlPath: The URL path to expose, e.g., "/static"
// directoryPath: The local directory path, e.g., "./public"
// Usage:
// builder.AddStaticDirectory("/static", "./public")
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
// dir := builder.GetStaticDirectory("/static")
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
//	builder.Handler(http.MethodGet, "/api/users", func(w http.ResponseWriter, r *http.Request) error {
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
//	if builder.IsRunning() {
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
// err := builder.Listen(":8080")
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
// err := builder.ListenTLS(":443", "cert.pem", "key.pem")
// // Use self-signed certificate
// err := builder.ListenTLS(":443", "", "")
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
// builder.RemoveStaticDirectory("/static")
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
// err := builder.Shutdown(ctx)
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
// err := builder.ShutdownTLS(ctx)
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
// err := builder.ShutdownAll(ctx)
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

// WebSocket registers a WebSocket handler for a specific URL pattern.
// pattern: URL pattern to match, e.g., "/ws"
// handler: Function to handle WebSocket connections
// Usage:
//
//	builder.WebSocket("/ws", func(w http.ResponseWriter, r *http.Request) {
//	    conn, err := upgrader.Upgrade(w, r, nil)
//	    if err != nil {
//	        return
//	    }
//	    defer conn.Close()
//	    // WebSocket handling logic
//	})
func (s *SERVE) WebSocket(pattern string, handler WebSocketHandler) *SERVE {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.webSocketRoutes = append(s.webSocketRoutes, WEBSOCKETROUTE{
		Pattern: pattern,
		Handler: handler,
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

func (s *SERVE) handleRequest(w http.ResponseWriter, r *http.Request) {
	s.mutex.Lock()
	routes := make([]ROUTE, len(s.routes))
	copy(routes, s.routes)
	webSocketRoutes := make([]WEBSOCKETROUTE, len(s.webSocketRoutes))
	copy(webSocketRoutes, s.webSocketRoutes)
	s.mutex.Unlock()

	var handled bool
	if strings.ToLower(r.Header.Get("Upgrade")) == "websocket" {
		for _, wsRoute := range webSocketRoutes {
			if s.matchPath(wsRoute.Pattern, r.URL.Path) {
				wsRoute.Handler(w, r)
				handled = true
				break
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
