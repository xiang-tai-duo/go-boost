// Package boost
// File:        serve.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: HTTP server with WebSocket support for Go applications
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

const (
	HEADER_UPGRADE    = "Upgrade"
	WEBSOCKET_UPGRADE = "websocket"
	DEFAULT_PORT      = 80
)

//goland:noinspection SpellCheckingInspection
type (
	RequestHandler func(w http.ResponseWriter, r *http.Request) error
	WebSocketDataHandler func(ws *websocket.Conn, messageType int, data []byte) error
	WebSocketDisconnectHandler func(ws *websocket.Conn) error
	ROUTE struct {
		Method  string
		Pattern string
		Handler RequestHandler
	}
	WebSocketOriginFilter interface {
		Allow(origin string) bool
	}
	WebSocketOriginMap map[string]bool
	WEBSOCKET_ORIGIN_REGEX struct {
		Pattern string
		Regex   *regexp.Regexp
	}
	WEBSOCKET_ORIGIN_ALLOW_ALL struct{}
	WEBSOCKET_DATA struct {
		Pattern           string
		Handler           WebSocketDataHandler
		DisconnectHandler WebSocketDisconnectHandler
		Filter            WebSocketOriginFilter
	}
	SERVE struct {
		Server              *http.Server
		Port                int
		Routes              []ROUTE
		WebSocketDataRoutes []WEBSOCKET_DATA
		WebSocketClients    map[string]*WEBSOCKET_CLIENT
		Mutex               sync.Mutex
		isRunning           bool
		StaticDirectories   map[string]string
		websocket           *WEBSOCKET_CLIENT_MANAGER
	}
)

var (
	WebSocketAllowAll WEBSOCKET_ORIGIN_ALLOW_ALL
)

func NewServe() *SERVE {
	port := DEFAULT_PORT
	return &SERVE{
		Server: &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		Port:              port,
		StaticDirectories: make(map[string]string),
		WebSocketClients:  make(map[string]*WEBSOCKET_CLIENT),
		websocket:         NewWebSocket(),
	}
}

func (s *SERVE) AddStaticDirectory(urlPath string, directoryPath string) *SERVE {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.StaticDirectories[urlPath] = directoryPath
	return s
}

func (s *SERVE) GetAvailablePort() (int, error) {
	var port int
	var err error

	for {
		var listener net.Listener
		if listener, err = net.Listen("tcp", ":0"); err == nil {
			port = listener.Addr().(*net.TCPAddr).Port
			_ = listener.Close()

			if port > 1024 && s.CheckPortAvailable(port) {
				break
			}
		} else {
			port = 0
			break
		}
	}

	return port, err
}

func (s *SERVE) CheckPortAvailable(port int) bool {
	address := fmt.Sprintf("localhost:%d", port)
	conn, err := net.DialTimeout("tcp", address, 100*time.Millisecond)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			return true
		}
		return false
	}
	_ = conn.Close()
	return false
}

func (s *SERVE) GetStaticDirectory(urlPath string) string {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	return s.StaticDirectories[urlPath]
}

func (s *SERVE) IsRunning() bool {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	return s.isRunning
}

func (s *SERVE) On(method string, pattern string, handler RequestHandler) *SERVE {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Routes = append(s.Routes, ROUTE{
		Method:  method,
		Pattern: pattern,
		Handler: handler,
	})
	if !s.isRunning {
		go func() {
			if err := s.listen(s.Server.Addr); err != nil {
				if err.Error() != "server is already running" {
					fmt.Printf("Server error: %v\n", err)
				}
			}
		}()
	}
	return s
}

func (s *SERVE) OnWebSocket(pattern string, dataHandler WebSocketDataHandler) *SERVE {
	return s.OnWebSocketEx(pattern, map[string]bool{"*": true}, dataHandler, nil)
}

func (s *SERVE) OnWebSocketEx(pattern string, filter interface{}, dataHandler WebSocketDataHandler, disconnectHandler WebSocketDisconnectHandler) *SERVE {
	var originFilter WebSocketOriginFilter
	if filter != nil {
		switch f := filter.(type) {
		case string:
			originMap := make(WebSocketOriginMap)
			originMap[f] = true
			originFilter = originMap
		case []string:
			originMap := make(WebSocketOriginMap)
			for _, origin := range f {
				originMap[origin] = true
			}
			originFilter = originMap
		}
	}
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.WebSocketDataRoutes = append(s.WebSocketDataRoutes, WEBSOCKET_DATA{
		Pattern:           pattern,
		Handler:           dataHandler,
		DisconnectHandler: disconnectHandler,
		Filter:            originFilter,
	})
	if !s.isRunning {
		go func() {
			if err := s.listen(s.Server.Addr); err != nil {
				if err.Error() != "server is already running" {
					fmt.Printf("Server error: %v\n", err)
				}
			}
		}()
	}
	return s
}

func (s *SERVE) RemoveStaticDirectory(urlPath string) *SERVE {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	delete(s.StaticDirectories, urlPath)
	return s
}

func (s *SERVE) Shutdown(ctx context.Context) error {
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

// Allow cecks if the origin is in the allowed origins map
func (m WebSocketOriginMap) Allow(origin string) bool {
	return m[origin]
}

func (r *WEBSOCKET_ORIGIN_REGEX) Allow(origin string) bool {
	return r.Regex.MatchString(origin)
}

func (a WEBSOCKET_ORIGIN_ALLOW_ALL) Allow(origin string) bool {
	return true
}

func (s *SERVE) handleHTTPRequest(w http.ResponseWriter, r *http.Request, routes []ROUTE) bool {
	var handled bool
	for _, route := range routes {
		if (route.Method == "*" || route.Method == r.Method) && s.matchPath(route.Pattern, r.URL.Path) {
			if err := route.Handler(w, r); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			handled = true
			break
		}
	}
	if !handled {
		handled = s.serveStatic(w, r)
	}
	return handled
}

func (s *SERVE) handleRequest(w http.ResponseWriter, r *http.Request) {
	s.Mutex.Lock()
	routes := make([]ROUTE, len(s.Routes))
	copy(routes, s.Routes)
	webSocketDataRoutes := make([]WEBSOCKET_DATA, len(s.WebSocketDataRoutes))
	copy(webSocketDataRoutes, s.WebSocketDataRoutes)
	s.Mutex.Unlock()
	var handled bool
	if strings.ToLower(r.Header.Get(HEADER_UPGRADE)) == WEBSOCKET_UPGRADE {
		handled = s.handleWebSocketRequest(w, r, webSocketDataRoutes)
	} else {
		handled = s.handleHTTPRequest(w, r, routes)
	}
	if !handled {
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func (s *SERVE) handleWebSocketRequest(w http.ResponseWriter, r *http.Request, webSocketDataRoutes []WEBSOCKET_DATA) bool {
	var handled bool
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
				go func(conn *websocket.Conn, route WEBSOCKET_DATA) {
					defer func() {
						if route.DisconnectHandler != nil {
							_ = route.DisconnectHandler(conn)
						}
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

func (s *SERVE) listen(addr string) error {
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

func (s *SERVE) listenTLS(addr string, certFile string, keyFile string) error {
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

func (s *SERVE) serveStatic(w http.ResponseWriter, r *http.Request) bool {
	var served bool
	s.Mutex.Lock()
	staticDirectories := make(map[string]string)
	for k, v := range s.StaticDirectories {
		staticDirectories[k] = v
	}
	s.Mutex.Unlock()
	for urlPath, directoryPath := range staticDirectories {
		if strings.HasPrefix(r.URL.Path, urlPath) {
			abs, err := filepath.Abs(directoryPath)
			if err != nil {
				abs = directoryPath
			}
			relFilePath := strings.TrimPrefix(r.URL.Path, urlPath)
			filePath := filepath.Join(abs, relFilePath)
			if _, err := os.Stat(filePath); !os.IsNotExist(err) {
				http.ServeFile(w, r, filePath)
				served = true
				break
			}
		}
	}
	return served
}

func (s *SERVE) shutdownAll(ctx context.Context) error {
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

func (s *SERVE) shutdownTLS(ctx context.Context) error {
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

func generateSelfSignedCertificate() (*tls.Certificate, error) {
	var privateKey *rsa.PrivateKey
	var err error
	var certificate *tls.Certificate
	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err == nil {
		validFromTime := time.Now()
		validToTime := validFromTime.Add(365 * 24 * time.Hour)
		serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
		serialNumber, serialErr := rand.Int(rand.Reader, serialNumberLimit)
		if serialErr == nil {
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
			var certificateDER []byte
			certificateDER, err = x509.CreateCertificate(rand.Reader, &certificateTemplate, &certificateTemplate, &privateKey.PublicKey, privateKey)
			if err == nil {
				var privateKeyDER []byte
				privateKeyDER, err = x509.MarshalPKCS8PrivateKey(privateKey)
				if err == nil {
					certificatePEM := new(bytes.Buffer)
					_ = pem.Encode(certificatePEM, &pem.Block{Type: "CERTIFICATE", Bytes: certificateDER})

					privateKeyPEM := new(bytes.Buffer)
					_ = pem.Encode(privateKeyPEM, &pem.Block{Type: "PRIVATE KEY", Bytes: privateKeyDER})

					var tlsCertificate tls.Certificate
					tlsCertificate, err = tls.X509KeyPair(certificatePEM.Bytes(), privateKeyPEM.Bytes())
					if err == nil {
						certificate = &tlsCertificate
					}
				}
			}
		} else {
			err = serialErr
		}
	}
	return certificate, err
}
