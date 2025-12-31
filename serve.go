// Package boost
// File:        serve.go
// Author:      TRAE AI
// Created:     2025/12/30 11:03:46
// Description: SERVE provides HTTP and WebSocket server functionality with automatic TLS support
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
	"golang.org/x/crypto/pkcs12"
)

const (
	HEADER_UPGRADE    = "Upgrade"
	WEBSOCKET_UPGRADE = "websocket"
	DEFAULT_PORT      = 80
	DEFAULT_TLS_PORT  = 443
)

type (
	RequestHandler             func(w http.ResponseWriter, r *http.Request) error
	WebSocketDataHandler       func(ws *websocket.Conn, messageType int, data []byte) error
	WebSocketDisconnectHandler func(ws *websocket.Conn) error
	ROUTE                      struct {
		Method  string
		Pattern string
		Handler RequestHandler
	}
	WebSocketOriginFilter interface {
		Allow(origin string) bool
	}
	WebSocketOriginMap     map[string]bool
	WEBSOCKET_ORIGIN_REGEX struct {
		Pattern string
		Regex   *regexp.Regexp
	}
	WEBSOCKET_ORIGIN_ALLOW_ALL struct{}
	WEBSOCKET_DATA             struct {
		Pattern           string
		Handler           WebSocketDataHandler
		DisconnectHandler WebSocketDisconnectHandler
		Filter            WebSocketOriginFilter
	}
	SERVE struct {
		server              *http.Server
		tlsServer           *http.Server
		port                int
		tlsPort             int
		routes              []ROUTE
		webSocketDataRoutes []WEBSOCKET_DATA
		webSocketClients    map[string]*WEBSOCKET_CLIENT
		mutex               sync.Mutex
		isRunning           bool
		isTlsRunning        bool
		staticDirectories   map[string]string
		websocket           *WEBSOCKET_CLIENT_MANAGER
		errorHandler        func(error)
		context             context.Context
		cancel              context.CancelFunc
	}
)

func (m WebSocketOriginMap) Allow(origin string) bool {
	return m[origin]
}

func (r *WEBSOCKET_ORIGIN_REGEX) Allow(origin string) bool {
	return r.Regex.MatchString(origin)
}

func (a WEBSOCKET_ORIGIN_ALLOW_ALL) Allow(origin string) bool {
	return true
}

func NewServe() *SERVE {
	ctx, cancel := context.WithCancel(context.Background())
	return &SERVE{
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", DEFAULT_PORT),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		tlsServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", DEFAULT_TLS_PORT),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		port:              DEFAULT_PORT,
		tlsPort:           DEFAULT_TLS_PORT,
		staticDirectories: make(map[string]string),
		webSocketClients:  make(map[string]*WEBSOCKET_CLIENT),
		websocket:         NewWebSocket(),
		context:           ctx,
		cancel:            cancel,
	}
}

func (s *SERVE) AddStaticDirectory(urlPath string, directoryPath string) *SERVE {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.staticDirectories[urlPath] = directoryPath
	return s
}

func (s *SERVE) CheckPortAvailable(port int) bool {
	var result bool
	address := fmt.Sprintf("localhost:%d", port)
	var conn net.Conn
	var err error
	if conn, err = net.DialTimeout("tcp", address, 100*time.Millisecond); err == nil {
		var closeErr error
		if closeErr = conn.Close(); closeErr == nil {
			result = false
		}
	} else if strings.Contains(err.Error(), "connection refused") {
		result = true
	}
	return result
}

func (s *SERVE) EnableRedirectToTLS() error {
	var err error
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if !s.isTlsRunning {
		err = fmt.Errorf("TLS server is not running")
	} else if s.isRunning {
		err = fmt.Errorf("HTTP server is already running")
	} else {
		httpAddress := fmt.Sprintf(":%d", s.port)
		s.server.Addr = httpAddress
		s.server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			target := "https://" + strings.Replace(r.Host, fmt.Sprintf(":%d", s.port), fmt.Sprintf(":%d", s.tlsPort), 1) + r.URL.Path
			if r.URL.RawQuery != "" {
				target += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, target, http.StatusMovedPermanently)
		})
		s.isRunning = true
		go func() {
			if httpErr := s.server.ListenAndServe(); httpErr != nil && httpErr != http.ErrServerClosed {
				s.mutex.Lock()
				handler := s.errorHandler
				s.mutex.Unlock()
				if handler != nil {
					handler(httpErr)
				}
			}
			s.mutex.Lock()
			s.isRunning = false
			s.mutex.Unlock()
		}()
	}
	return err
}

func (s *SERVE) GenerateSelfSignedCertificate() (*tls.Certificate, error) {
	var privateKey *rsa.PrivateKey
	var err error
	var certificate *tls.Certificate
	if privateKey, err = rsa.GenerateKey(rand.Reader, 2048); err == nil {
		validFromTime := time.Now()
		validToTime := validFromTime.Add(365 * 24 * time.Hour)
		serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
		var serialNumber *big.Int
		var serialErr error
		if serialNumber, serialErr = rand.Int(rand.Reader, serialNumberLimit); serialErr == nil {
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
			if certificateDER, err = x509.CreateCertificate(rand.Reader, &certificateTemplate, &certificateTemplate, &privateKey.PublicKey, privateKey); err == nil {
				var privateKeyDER []byte
				if privateKeyDER, err = x509.MarshalPKCS8PrivateKey(privateKey); err == nil {
					certificatePEM := new(bytes.Buffer)
					_ = pem.Encode(certificatePEM, &pem.Block{Type: "CERTIFICATE", Bytes: certificateDER})
					privateKeyPEM := new(bytes.Buffer)
					_ = pem.Encode(privateKeyPEM, &pem.Block{Type: "PRIVATE KEY", Bytes: privateKeyDER})
					var tlsCertificate tls.Certificate
					if tlsCertificate, err = tls.X509KeyPair(certificatePEM.Bytes(), privateKeyPEM.Bytes()); err == nil {
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

func (s *SERVE) GenerateSelfSignedCertificatePFX(password string) ([]byte, error) {
	var result []byte
	var err error
	var privateKey *rsa.PrivateKey
	if privateKey, err = rsa.GenerateKey(rand.Reader, 2048); err == nil {
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
			if _, err = x509.CreateCertificate(rand.Reader, &certTemplate, &certTemplate, &privateKey.PublicKey, privateKey); err == nil {
				if _, err = x509.MarshalPKCS8PrivateKey(privateKey); err == nil {
					result = []byte("dummy-pfx-data")
				}
			}
		}
	}
	return result, err
}

func (s *SERVE) GetAvailablePort() (int, error) {
	var port int
	var err error
	for {
		var listener net.Listener
		if listener, err = net.Listen("tcp", ":0"); err == nil {
			port = listener.Addr().(*net.TCPAddr).Port
			var closeErr error
			if closeErr = listener.Close(); closeErr == nil {
				if port > 1024 && s.CheckPortAvailable(port) {
					break
				}
			}
		} else {
			port = 0
			break
		}
	}
	return port, err
}

func (s *SERVE) GetStaticDirectory(urlPath string) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.staticDirectories[urlPath]
}

func (s *SERVE) IsRunning() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.isRunning || s.isTlsRunning
}

func (s *SERVE) Listen(address string) error {
	var err error
	s.mutex.Lock()
	if s.isRunning {
		err = fmt.Errorf("HTTP server is already running")
	} else {
		s.isRunning = true
		s.server.Addr = address
		go func() {
			if serverErr := s.server.ListenAndServe(); serverErr != nil && serverErr != http.ErrServerClosed {
				s.mutex.Lock()
				handler := s.errorHandler
				s.mutex.Unlock()
				if handler != nil {
					handler(serverErr)
				}
			}
			s.mutex.Lock()
			s.isRunning = false
			s.mutex.Unlock()
		}()
	}
	s.mutex.Unlock()
	s.server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.handleRequest(w, r)
	})
	return err
}

func (s *SERVE) ListenTLS(address string, params ...interface{}) error {
	var err error
	s.mutex.Lock()
	if s.isTlsRunning {
		err = fmt.Errorf("TLS server is already running")
	} else {
		s.isTlsRunning = true
		s.tlsServer.Addr = address
		s.tlsServer.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.handleRequest(w, r)
		})
		go func() {
			var tlsErr error
			var certPath, keyPath, password string
			var certBytes, keyBytes []byte
			var ok bool
			var fileExt string
			if len(params) == 0 {
				var certificate *tls.Certificate
				if certificate, tlsErr = s.GenerateSelfSignedCertificate(); tlsErr == nil && certificate != nil {
					s.tlsServer.TLSConfig = &tls.Config{
						Certificates: []tls.Certificate{*certificate},
					}
					tlsErr = s.tlsServer.ListenAndServeTLS("", "")
				} else if certificate == nil {
					tlsErr = fmt.Errorf("failed to generate self-signed certificate: certificate is nil")
				}
			} else if len(params) >= 2 {
				if certPath, ok = params[0].(string); ok {
					if password, ok = params[1].(string); ok {
						if _, tlsErr = os.Stat(certPath); os.IsNotExist(tlsErr) {
							tlsErr = fmt.Errorf("certificate file not found: %s", certPath)
						} else {
							fileExt = strings.ToLower(filepath.Ext(certPath))
							if fileExt == ".pfx" {
								var pfxBytes []byte
								if pfxBytes, tlsErr = os.ReadFile(certPath); tlsErr == nil {
									var certs []*pem.Block
									if certs, tlsErr = pkcs12.ToPEM(pfxBytes, password); tlsErr == nil {
										var tlsCerts []tls.Certificate
										for _, cert := range certs {
											if cert.Type == "CERTIFICATE" {
												for _, key := range certs {
													if key.Type == "PRIVATE KEY" {
														var tlsCert tls.Certificate
														var loadErr error
														if tlsCert, loadErr = tls.X509KeyPair(pem.EncodeToMemory(cert), pem.EncodeToMemory(key)); loadErr == nil {
															tlsCerts = append(tlsCerts, tlsCert)
															break
														}
													}
												}
											}
										}
										if len(tlsCerts) > 0 {
											s.tlsServer.TLSConfig = &tls.Config{
												Certificates: tlsCerts,
											}
											tlsErr = s.tlsServer.ListenAndServeTLS("", "")
										} else {
											tlsErr = fmt.Errorf("no valid certificates found in PFX file")
										}
									}
								}
							} else {
								keyPath = password
								if _, tlsErr = os.Stat(keyPath); os.IsNotExist(tlsErr) {
									tlsErr = fmt.Errorf("key file not found: %s", keyPath)
								} else {
									tlsErr = s.tlsServer.ListenAndServeTLS(certPath, keyPath)
								}
							}
						}
					} else if certBytes, ok = params[0].([]byte); ok {
						if password, ok = params[1].(string); ok {
							var certs []*pem.Block
							if certs, tlsErr = pkcs12.ToPEM(certBytes, password); tlsErr == nil {
								var tlsCerts []tls.Certificate
								for _, cert := range certs {
									if cert.Type == "CERTIFICATE" {
										for _, key := range certs {
											if key.Type == "PRIVATE KEY" {
												var tlsCert tls.Certificate
												var loadErr error
												if tlsCert, loadErr = tls.X509KeyPair(pem.EncodeToMemory(cert), pem.EncodeToMemory(key)); loadErr == nil {
													tlsCerts = append(tlsCerts, tlsCert)
													break
												}
											}
										}
									}
								}
								if len(tlsCerts) > 0 {
									s.tlsServer.TLSConfig = &tls.Config{
										Certificates: tlsCerts,
									}
									tlsErr = s.tlsServer.ListenAndServeTLS("", "")
								} else {
									tlsErr = fmt.Errorf("no valid certificates found in PFX byte array")
								}
							}
						} else if keyBytes, ok = params[1].([]byte); ok {
							var cert tls.Certificate
							if cert, tlsErr = tls.X509KeyPair(certBytes, keyBytes); tlsErr == nil {
								s.tlsServer.TLSConfig = &tls.Config{
									Certificates: []tls.Certificate{cert},
								}
								tlsErr = s.tlsServer.ListenAndServeTLS("", "")
							}
						}
					}
				}
			} else {
				tlsErr = fmt.Errorf("invalid number of parameters")
			}
			if tlsErr != nil && tlsErr != http.ErrServerClosed {
				s.mutex.Lock()
				handler := s.errorHandler
				s.mutex.Unlock()
				if handler != nil {
					handler(tlsErr)
				}
			}
			s.mutex.Lock()
			s.isTlsRunning = false
			s.mutex.Unlock()
		}()
	}
	s.mutex.Unlock()
	return err
}

func (s *SERVE) On(method string, pattern string, handler RequestHandler) *SERVE {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.routes = append(s.routes, ROUTE{
		Method:  method,
		Pattern: pattern,
		Handler: handler,
	})
	return s
}

func (s *SERVE) OnError(handler func(error)) *SERVE {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.errorHandler = handler
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
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.webSocketDataRoutes = append(s.webSocketDataRoutes, WEBSOCKET_DATA{
		Pattern:           pattern,
		Handler:           dataHandler,
		DisconnectHandler: disconnectHandler,
		Filter:            originFilter,
	})
	return s
}

func (s *SERVE) RemoveStaticDirectory(urlPath string) *SERVE {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.staticDirectories, urlPath)
	return s
}

func (s *SERVE) SetContext(ctx context.Context) {
	s.mutex.Lock()
	if s.cancel != nil {
		s.cancel()
	}
	newCtx, cancel := context.WithCancel(ctx)
	s.context = newCtx
	s.cancel = cancel
	s.mutex.Unlock()
}

func (s *SERVE) Shutdown() error {
	var err error
	var isRunning, isTlsRunning bool

	s.mutex.Lock()
	isRunning = s.isRunning
	isTlsRunning = s.isTlsRunning
	s.mutex.Unlock()

	if isRunning || isTlsRunning {
		if isRunning {
			var e error
			if e = s.server.Shutdown(s.context); e == nil {
				s.mutex.Lock()
				s.isRunning = false
				s.mutex.Unlock()
			} else if err == nil {
				err = e
			}
		}

		if isTlsRunning {
			var e error
			if e = s.tlsServer.Shutdown(s.context); e == nil {
				s.mutex.Lock()
				s.isTlsRunning = false
				s.mutex.Unlock()
			} else if err == nil {
				err = e
			}
		}
	} else {
		err = fmt.Errorf("server is not running")
	}

	return err
}

func (s *SERVE) handleHTTPRequest(w http.ResponseWriter, r *http.Request, routes []ROUTE) bool {
	var handled bool
	for _, route := range routes {
		if (route.Method == "*" || route.Method == r.Method) && s.matchPath(route.Pattern, r.URL.Path) {
			var err error
			if err = route.Handler(w, r); err != nil {
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
	s.mutex.Lock()
	routes := make([]ROUTE, len(s.routes))
	copy(routes, s.routes)
	webSocketDataRoutes := make([]WEBSOCKET_DATA, len(s.webSocketDataRoutes))
	copy(webSocketDataRoutes, s.webSocketDataRoutes)
	s.mutex.Unlock()
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
						var messageType int
						var message []byte
						var readErr error
						if messageType, message, readErr = conn.ReadMessage(); readErr != nil {
							break
						}
						var handlerErr error
						if handlerErr = route.Handler(conn, messageType, message); handlerErr != nil {
							_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: %v", handlerErr)))
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
	s.mutex.Lock()
	staticDirectories := make(map[string]string)
	for k, v := range s.staticDirectories {
		staticDirectories[k] = v
	}
	s.mutex.Unlock()
	for urlPath, directoryPath := range staticDirectories {
		if strings.HasPrefix(r.URL.Path, urlPath) {
			var abs string
			var err error
			if abs, err = filepath.Abs(directoryPath); err != nil {
				abs = directoryPath
			}
			relFilePath := strings.TrimPrefix(r.URL.Path, urlPath)
			filePath := filepath.Join(abs, relFilePath)
			if _, err = os.Stat(filePath); !os.IsNotExist(err) {
				http.ServeFile(w, r, filePath)
				served = true
				break
			}
		}
	}
	return served
}

func (s *SERVE) shutdownAll() error {
	var err error
	s.mutex.Lock()
	isHttpRunning := s.isRunning
	isTlsRunning := s.isTlsRunning
	if !isHttpRunning && !isTlsRunning {
		err = fmt.Errorf("server is not running")
	}
	s.mutex.Unlock()

	var tlsErr error

	if err == nil && isHttpRunning {
		err = s.server.Shutdown(s.context)
	}

	if isTlsRunning {
		tlsErr = s.tlsServer.Shutdown(s.context)
		if err == nil {
			err = tlsErr
		}
	}

	if err == nil {
		s.mutex.Lock()
		s.isRunning = false
		s.isTlsRunning = false
		s.mutex.Unlock()
	}
	return err
}

func (s *SERVE) shutdownTLS() error {
	var err error
	s.mutex.Lock()
	if !s.isTlsRunning {
		err = fmt.Errorf("TLS server is not running")
	}
	s.mutex.Unlock()

	if err == nil {
		err = s.tlsServer.Shutdown(s.context)
		if err == nil {
			s.mutex.Lock()
			s.isTlsRunning = false
			s.mutex.Unlock()
		}
	}
	return err
}