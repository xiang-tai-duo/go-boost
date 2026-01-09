// Package boost
// File:        serve.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/serve.go
// Author:      TRAE AI
// Created:     12/30/2025 11:03:46
// Description: SERVE provides HTTP and WebSocket server functionality with automatic TLS support
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
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	server "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"
	"golang.org/x/crypto/pkcs12"
)

//goland:noinspection GoSnakeCaseUsage
const (
	HEADER_UPGRADE           = "Upgrade"
	WEBSOCKET_UPGRADE        = "websocket"
	DEFAULT_PORT             = 80
	DEFAULT_TLS_PORT         = 443
	DEFAULT_MQTT_SERVER_PORT = 1883
	DEFAULT_MQTT_SERVER_HOST = "0.0.0.0"
)

//goland:noinspection GoSnakeCaseUsage
type (
	RequestHandler             func(w http.ResponseWriter, r *http.Request) error
	ServeStaticCallbackHandler func(served bool, uri string, filePath string)
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
		serveStaticCallback ServeStaticCallbackHandler
	}

	PublishCallback      func(client *server.Client, packet packets.Packet)
	SubscribeCallback    func(client *server.Client, packet packets.Packet)
	SubscribedCallback   func(client *server.Client, packet packets.Packet)
	UnsubscribeCallback  func(client *server.Client, packet packets.Packet)
	UnsubscribedCallback func(client *server.Client, packet packets.Packet)
	ConnectCallback      func(client *server.Client, packet packets.Packet)
	DisconnectCallback   func(client *server.Client, err error)
	AuthCallback         func(client *server.Client, packet packets.Packet) bool
	ACLCheckCallback     func(client *server.Client, topic string, write bool) bool

	MQTT_SERVER struct {
		Server               *server.Server
		Host                 string
		Port                 int
		CertFile             string
		KeyFile              string
		TLSConfig            *tls.Config
		isRunning            bool
		lock                 sync.Mutex
		publishCallback      PublishCallback
		subscribeCallback    SubscribeCallback
		subscribedCallback   SubscribedCallback
		unsubscribeCallback  UnsubscribeCallback
		unsubscribedCallback UnsubscribedCallback
		connectCallback      ConnectCallback
		disconnectCallback   DisconnectCallback
		authCallback         AuthCallback
		aclCheckCallback     ACLCheckCallback
	}

	MQTT_SERVE_HOOK struct {
		server *MQTT_SERVER
		server.HookBase
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
	} else if strings.Contains(err.Error(), "refused") {
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
			if tlsErr != nil && !errors.Is(tlsErr, http.ErrServerClosed) {
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

// OnStaticServe sets the static serve callback
func (s *SERVE) OnStaticServe(callback ServeStaticCallbackHandler) *SERVE {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.serveStaticCallback = callback
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
	var foundFilePath string
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
				foundFilePath = filePath
				break
			}
		}
	}
	if s.serveStaticCallback != nil {
		s.serveStaticCallback(served, r.URL.Path, foundFilePath)
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

// MQTT_SERVER related methods

func NewMQTTServer(params ...interface{}) *MQTT_SERVER {
	host := DEFAULT_MQTT_SERVER_HOST
	port := DEFAULT_MQTT_SERVER_PORT

	if len(params) > 0 {
		if h, ok := params[0].(string); ok && h != "" {
			host = h
		}
	}

	if len(params) > 1 {
		if p, ok := params[1].(int); ok && p > 0 {
			port = p
		}
	}

	return &MQTT_SERVER{
		Host:     host,
		Port:     port,
		CertFile: "",
		KeyFile:  "",
		Server:   server.New(nil),
	}
}

func NewMQTTServerTLS(host string, port int) *MQTT_SERVER {
	if host == "" {
		host = DEFAULT_MQTT_SERVER_HOST
	}

	if port <= 0 {
		port = DEFAULT_MQTT_SERVER_PORT
	}

	return &MQTT_SERVER{
		Host:     host,
		Port:     port,
		CertFile: "",
		KeyFile:  "",
		Server:   server.New(nil),
	}
}

func (ms *MQTT_SERVER) GetHost() string {
	return ms.Host
}

func (ms *MQTT_SERVER) GetPort() int {
	return ms.Port
}

func (ms *MQTT_SERVER) IsRunning() bool {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	return ms.isRunning
}

func (h *MQTT_SERVE_HOOK) ID() string {
	return "go-boost-hook"
}

func (h *MQTT_SERVE_HOOK) OnACLCheck(client *server.Client, topic string, write bool) bool {
	var result bool = true
	if h.server.aclCheckCallback != nil {
		result = h.server.aclCheckCallback(client, topic, write)
	}
	return result
}

func (h *MQTT_SERVE_HOOK) OnConnect(client *server.Client, packet packets.Packet) error {
	if h.server.connectCallback != nil {
		h.server.connectCallback(client, packet)
	}
	return h.HookBase.OnConnect(client, packet)
}

func (h *MQTT_SERVE_HOOK) OnConnectAuthenticate(client *server.Client, packet packets.Packet) bool {
	var result bool = true
	if h.server.authCallback != nil {
		result = h.server.authCallback(client, packet)
	}
	return result
}

func (h *MQTT_SERVE_HOOK) OnDisconnect(client *server.Client, err error, expire bool) {
	if h.server.disconnectCallback != nil {
		h.server.disconnectCallback(client, err)
	}
	h.HookBase.OnDisconnect(client, err, expire)
}

func (h *MQTT_SERVE_HOOK) OnPublish(client *server.Client, packet packets.Packet) (packets.Packet, error) {
	if h.server.publishCallback != nil {
		h.server.publishCallback(client, packet)
	}
	return h.HookBase.OnPublish(client, packet)
}

func (h *MQTT_SERVE_HOOK) OnSubscribe(client *server.Client, packet packets.Packet) packets.Packet {
	if h.server.subscribeCallback != nil {
		h.server.subscribeCallback(client, packet)
	}
	return h.HookBase.OnSubscribe(client, packet)
}

func (h *MQTT_SERVE_HOOK) OnSubscribed(client *server.Client, packet packets.Packet, reasonCodes []byte) {
	if h.server.subscribedCallback != nil {
		h.server.subscribedCallback(client, packet)
	}
	h.HookBase.OnSubscribed(client, packet, reasonCodes)
}

func (h *MQTT_SERVE_HOOK) OnUnsubscribe(client *server.Client, packet packets.Packet) packets.Packet {
	if h.server.unsubscribeCallback != nil {
		h.server.unsubscribeCallback(client, packet)
	}
	return h.HookBase.OnUnsubscribe(client, packet)
}

func (h *MQTT_SERVE_HOOK) OnUnsubscribed(client *server.Client, packet packets.Packet) {
	if h.server.unsubscribedCallback != nil {
		h.server.unsubscribedCallback(client, packet)
	}
	h.HookBase.OnUnsubscribed(client, packet)
}

func (h *MQTT_SERVE_HOOK) Provides(b byte) bool {
	return true
}

func (ms *MQTT_SERVER) Publish(topic string, payload string, params ...interface{}) error {
	var err error
	var qos int = DEFAULT_MQTT_QOS
	var retained bool = false

	if len(params) > 0 {
		if q, ok := params[0].(int); ok {
			qos = q
		}
	}
	if len(params) > 1 {
		if r, ok := params[1].(bool); ok {
			retained = r
		}
	}

	// 使用mochi-mqtt服务器的Publish方法发布消息
	if ms.Server != nil {
		// 调用正确的Publish方法，参数为：topic, payload, retain, qos
		err = ms.Server.Publish(topic, []byte(payload), retained, byte(qos))
	}

	return err
}

func (ms *MQTT_SERVER) SetACLCheckCallback(callback ACLCheckCallback) error {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	ms.aclCheckCallback = callback
	return nil
}

func (ms *MQTT_SERVER) SetAuthCallback(callback AuthCallback) error {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	ms.authCallback = callback
	return nil
}

func (ms *MQTT_SERVER) SetConnectCallback(callback ConnectCallback) error {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	ms.connectCallback = callback
	return nil
}

func (ms *MQTT_SERVER) SetDisconnectCallback(callback DisconnectCallback) error {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	ms.disconnectCallback = callback
	return nil
}

func (ms *MQTT_SERVER) SetPublishCallback(callback PublishCallback) error {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	ms.publishCallback = callback
	return nil
}

func (ms *MQTT_SERVER) SetSubscribeCallback(callback SubscribeCallback) error {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	ms.subscribeCallback = callback
	return nil
}

func (ms *MQTT_SERVER) SetSubscribedCallback(callback SubscribedCallback) error {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	ms.subscribedCallback = callback
	return nil
}

func (ms *MQTT_SERVER) SetUnsubscribeCallback(callback UnsubscribeCallback) error {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	ms.unsubscribeCallback = callback
	return nil
}

func (ms *MQTT_SERVER) SetUnsubscribedCallback(callback UnsubscribedCallback) error {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	ms.unsubscribedCallback = callback
	return nil
}

func (ms *MQTT_SERVER) SetCertFile(certFilePath string) error {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	ms.CertFile = certFilePath
	return nil
}

func (ms *MQTT_SERVER) SetKeyFile(keyFilePath string) error {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	ms.KeyFile = keyFilePath
	return nil
}

func (ms *MQTT_SERVER) Start(params ...interface{}) error {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	var err error
	if ms.isRunning {
		err = nil
	} else {
		hook := &MQTT_SERVE_HOOK{server: ms}
		err = ms.Server.AddHook(hook, nil)
		if err == nil {
			var listener listeners.Listener
			address := fmt.Sprintf("%s:%d", ms.Host, ms.Port)
			var useTLS bool = len(params) > 0
			var tlsConfig *tls.Config
			var tlsErr error

			if useTLS {
				if len(params) == 0 {
					var cert *tls.Certificate
					if cert, tlsErr = ms.GenerateSelfSignedCertificate(); tlsErr == nil && cert != nil {
						tlsConfig = &tls.Config{
							Certificates: []tls.Certificate{*cert},
							MinVersion:   tls.VersionTLS12,
						}
					} else {
						tlsErr = fmt.Errorf("failed to generate self-signed certificate")
					}
				} else if len(params) >= 2 {
					var certPath, passwordOrKey string
					var ok bool
					var fileExt string

					if certPath, ok = params[0].(string); ok {
						if passwordOrKey, ok = params[1].(string); ok {
							if _, tlsErr = os.Stat(certPath); os.IsNotExist(tlsErr) {
								tlsErr = fmt.Errorf("certificate file not found: %s", certPath)
							} else {
								fileExt = strings.ToLower(filepath.Ext(certPath))
								if fileExt == ".pfx" {
									var pfxBytes []byte
									var certs []*pem.Block
									var tlsCerts []tls.Certificate

									if pfxBytes, tlsErr = os.ReadFile(certPath); tlsErr == nil {
										if certs, tlsErr = pkcs12.ToPEM(pfxBytes, passwordOrKey); tlsErr == nil {
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
										}
									}

									if len(tlsCerts) > 0 {
										tlsConfig = &tls.Config{
											Certificates: tlsCerts,
											MinVersion:   tls.VersionTLS12,
										}
									} else {
										tlsErr = fmt.Errorf("no valid certificates found in PFX file")
									}
								} else {
									var tlsCert tls.Certificate
									if tlsCert, tlsErr = tls.LoadX509KeyPair(certPath, passwordOrKey); tlsErr != nil {
										tlsErr = fmt.Errorf("failed to load certificate: %v", tlsErr)
									} else {
										tlsConfig = &tls.Config{
											Certificates: []tls.Certificate{tlsCert},
											MinVersion:   tls.VersionTLS12,
										}
									}
								}
							}
						}
					}
				}

				if tlsErr != nil {
					err = tlsErr
				} else if tlsConfig != nil {
					var tlsListener net.Listener
					if tlsListener, tlsErr = tls.Listen("tcp", address, tlsConfig); tlsErr != nil {
						err = tlsErr
					} else {
						listener = listeners.NewNet("tls", tlsListener)
					}
				}
			} else {
				listener = listeners.NewTCP(listeners.Config{
					ID:      "tcp",
					Address: address,
				})
			}

			if err == nil && listener != nil {
				if err = ms.Server.AddListener(listener); err == nil {
					go func() {
						if serveErr := ms.Server.Serve(); serveErr != nil {
							fmt.Printf("MQTT server error: %v\n", serveErr)
						}
					}()

					ms.isRunning = true
				}
			}
		}
	}
	return err
}

func (ms *MQTT_SERVER) Stop() error {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	var err error
	if !ms.isRunning {
		err = nil
	} else {
		if ms.Server != nil {
			ms.Server.Close()
		}

		ms.isRunning = false
		err = nil
	}
	return err
}

// GetMQTTMessage converts a packets.Packet to MQTT_MESSAGE
func (ms *MQTT_SERVER) GetMQTTMessage(packet packets.Packet) *MQTT_MESSAGE {
	return &MQTT_MESSAGE{
		Topic:     packet.TopicName,
		Payload:   string(packet.Payload),
		Timestamp: time.Now(),
		QoS:       int(packet.FixedHeader.Qos),
		Retained:  packet.FixedHeader.Retain,
		Duplicate: packet.FixedHeader.Dup,
	}
}

// GetClientID returns the client ID from the client object
func (ms *MQTT_SERVER) GetClientID(client *server.Client) string {
	return client.ID
}

// GetClientIP returns the client IP address and port from the client object
func (ms *MQTT_SERVER) GetClientIP(client *server.Client) (string, int) {
	var ip string
	var port int
	remoteAddr := client.Net.Remote
	if remoteAddr != "" {
		// 处理IPv6地址，IPv6地址格式为"[IP]:端口"
		if strings.Contains(remoteAddr, "[") {
			// 分割IPv6地址：[IP]:端口
			ipv6Parts := strings.Split(remoteAddr, "]:")
			if len(ipv6Parts) == 2 {
				ip = ipv6Parts[0][1:]                 // 移除开头的[获取IP
				fmt.Sscanf(ipv6Parts[1], "%d", &port) // 解析端口
			}
		} else {
			// 处理IPv4地址：IP:端口
			ipv4Parts := strings.Split(remoteAddr, ":")
			if len(ipv4Parts) >= 2 {
				// 对于IPv4，最后一部分是端口，前面的都是IP
				ip = strings.Join(ipv4Parts[:len(ipv4Parts)-1], ":")
				fmt.Sscanf(ipv4Parts[len(ipv4Parts)-1], "%d", &port) // 解析端口
			}
		}
	}
	// 仅在函数末尾返回结果
	return ip, port
}

// GetClientPort returns the client port from the client object
func (ms *MQTT_SERVER) GetClientPort(client *server.Client) int {
	// 使用反射获取客户端端口
	value := reflect.ValueOf(client)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// 尝试获取Conn字段
	connField := value.FieldByName("Conn")
	if connField.IsValid() && connField.CanInterface() {
		if conn, ok := connField.Interface().(net.Conn); ok {
			if addr, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
				return addr.Port
			}
		}
	}

	// 尝试获取其他可能的连接字段
	connField = value.FieldByName("Connection")
	if connField.IsValid() && connField.CanInterface() {
		if conn, ok := connField.Interface().(net.Conn); ok {
			if addr, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
				return addr.Port
			}
		}
	}

	return 0
}

// GetSubscriptions returns all subscription filters and QoS from the packet
func (ms *MQTT_SERVER) GetSubscriptions(packet packets.Packet) map[string]int {
	result := make(map[string]int)
	for _, sub := range packet.Filters {
		result[sub.Filter] = int(sub.Qos)
	}
	return result
}

// GetSubscriptionFilters returns all subscription filters from the packet
func (ms *MQTT_SERVER) GetSubscriptionFilters(packet packets.Packet) []string {
	var result []string
	for _, sub := range packet.Filters {
		result = append(result, sub.Filter)
	}
	return result
}

// GetPublishTopic returns the publish topic from the packet
func (ms *MQTT_SERVER) GetPublishTopic(packet packets.Packet) string {
	return packet.TopicName
}

// GetPublishPayload returns the publish payload from the packet
func (ms *MQTT_SERVER) GetPublishPayload(packet packets.Packet) string {
	return string(packet.Payload)
}

// GetConnectUsername returns the username from the connect packet
func (ms *MQTT_SERVER) GetConnectUsername(packet packets.Packet) string {
	return string(packet.Connect.Username)
}

// GetConnectPassword returns the password from the connect packet
func (ms *MQTT_SERVER) GetConnectPassword(packet packets.Packet) string {
	return string(packet.Connect.Password)
}

// GetOperationType returns "读取" for read operations and "写入" for write operations
func (ms *MQTT_SERVER) GetOperationType(write bool) string {
	if write {
		return "写入"
	}
	return "读取"
}

// FormatError returns the error message if error is not nil, otherwise returns "无"
func (ms *MQTT_SERVER) FormatError(err error) string {
	if err != nil {
		return err.Error()
	}
	return "无"
}

func (ms *MQTT_SERVER) GenerateSelfSignedCertificate() (*tls.Certificate, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Go-Boost"},
			CommonName:   "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, err
	}

	privDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privDER})

	cert, err := tls.X509KeyPair(certPEM, privPEM)
	if err != nil {
		return nil, err
	}

	return &cert, nil
}
