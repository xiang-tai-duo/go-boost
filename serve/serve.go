// Package serve
// File:        serve.go
// Url:         https://github.com/xiang-tai-duo/go-bootstrap/blob/master/serve/serve.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: SERVE provides HTTP server functionality with automatic TLS support
// --------------------------------------------------------------------------------
package serve

import (
	"bytes"
	_context "context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"embed"
	"encoding/pem"
	"errors"
	"fmt"
	"io/fs"
	"math/big"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tus/tusd/v2/pkg/filestore"
	tushandler "github.com/tus/tusd/v2/pkg/handler"
	"github.com/xiang-tai-duo/go-bootstrap/logger"
	x509util "github.com/xiang-tai-duo/go-bootstrap/x509"
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst,GoNameStartsWithPackageName,SpellCheckingInspection
const (
	ANY_METHOD                        = "*"
	BYTE_UNITS                        = "KMGTPE"
	CERTIFICATE_BLOCK_TYPE            = "CERTIFICATE"
	CERTIFICATE_FILE_NOT_FOUND_FORMAT = "certificate file not found: %s"
	CERTIFICATE_NIL_ERROR             = "failed to generate self-signed certificate: certificate is nil"
	CONTENT_TYPE                      = "Content-Type"
	CONTENT_TYPE_HTML_UTF8            = "text/html; charset=utf-8"
	CURRENT_DIRECTORY                 = "."
	DEFAULT_ADDRESS                   = ":80"
	DEFAULT_DIRECTORY_PERMISSION      = 0755
	DEFAULT_HTTP_PORT                 = 80
	DEFAULT_IDLE_TIMEOUT              = 60 * time.Second
	DEFAULT_INDEX_HTML                = "index.html"
	DEFAULT_READ_TIMEOUT              = 15 * time.Second
	DEFAULT_RSA_KEY_BITS              = 2048
	DEFAULT_SERIAL_NUMBER_BITS        = 128
	DEFAULT_TLS_ADDRESS               = ":443"
	DEFAULT_VALID_DAYS                = 365
	DEFAULT_WRITE_TIMEOUT             = 15 * time.Second
	DELETE                            = "DELETE"
	DIRECTORY_LIST_ENTRY_FORMAT       = "                <td><a href=\"%s\" class=\"%s\">%s</a></td>\n"
	DIRECTORY_LIST_HEADING_FORMAT     = "    <h1>%s - /</h1>\n"
	DIRECTORY_LIST_PARENT_FORMAT      = "                <td><a href=\"%s\" class=\"parent-dir\">Parent Directory</a></td>\n"
	DIRECTORY_LIST_TD_FORMAT          = "                <td>%s</td>\n"
	DIRECTORY_LIST_TITLE_FORMAT       = "    <title>%s - /</title>\n"
	EMPTY_PLACEHOLDER                 = "-"
	EXECUTE                           = "EXECUTE"
	FILENAME_KEY                      = "filename"
	FILE_RENAME_FORMAT                = "%s_%d%s"
	FILE_SIZE_BYTE_FORMAT             = "%d B"
	FILE_SIZE_FORMAT                  = "%.1f %cB"
	FILE_SIZE_UNIT                    = 1024
	FORBIDDEN_LOG_FORMAT              = "%s forbidden"
	GET                               = "GET"
	HTML_DOCTYPE                      = "<!DOCTYPE html>\n"
	HTML_HTML_OPEN                    = "<html lang=\"en\">\n"
	HTML_HEAD_OPEN                    = "<head>\n"
	HTML_HEAD_CLOSE                   = "</head>\n"
	HTML_META_CHARSET                 = "    <meta charset=\"UTF-8\">\n"
	HTML_META_VIEWPORT                = "    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n"
	HTML_STYLE_OPEN                   = "    <style>\n"
	HTML_STYLE_BODY                   = "        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; margin: 0; padding: 20px; background-color: #fff; }\n"
	HTML_STYLE_H1                     = "        h1 { color: #003399; border-bottom: 1px solid #ccc; padding-bottom: 10px; }\n"
	HTML_STYLE_TABLE                  = "        table { width: 100%; border-collapse: collapse; margin-top: 20px; }\n"
	HTML_STYLE_TH_TD                  = "        th, td { text-align: left; padding: 8px; border-bottom: 1px solid #ddd; }\n"
	HTML_STYLE_TH                     = "        th { background-color: #f2f2f2; font-weight: bold; color: #333; }\n"
	HTML_STYLE_TR_HOVER               = "        tr:hover { background-color: #f5f5f5; }\n"
	HTML_STYLE_A                      = "        a { color: #0066cc; text-decoration: none; }\n"
	HTML_STYLE_A_HOVER                = "        a:hover { text-decoration: underline; }\n"
	HTML_STYLE_DIR_ICON               = "        .dir-icon::before { content: '📁'; margin-right: 8px; }\n"
	HTML_STYLE_FILE_ICON              = "        .file-icon::before { content: '📄'; margin-right: 8px; }\n"
	HTML_STYLE_PARENT_DIR             = "        .parent-dir::before { content: '⬆️'; margin-right: 8px; }\n"
	HTML_STYLE_CLOSE                  = "    </style>\n"
	HTML_BODY_OPEN                    = "<body>\n"
	HTML_BODY_CLOSE                   = "</body>\n"
	HTML_HTML_CLOSE                   = "</html>\n"
	HTML_TABLE_OPEN                   = "    <table>\n"
	HTML_TABLE_CLOSE                  = "    </table>\n"
	HTML_THEAD_OPEN                   = "        <thead>\n"
	HTML_THEAD_CLOSE                  = "        </thead>\n"
	HTML_TBODY_OPEN                   = "        <tbody>\n"
	HTML_TBODY_CLOSE                  = "        </tbody>\n"
	HTML_TR_OPEN                      = "            <tr>\n"
	HTML_TR_CLOSE                     = "            </tr>\n"
	HTML_TH_NAME                      = "                <th>Name</th>\n"
	HTML_TH_SIZE                      = "                <th>Size</th>\n"
	HTML_TH_LAST_MODIFIED             = "                <th>Last Modified</th>\n"
	HTML_TD_DASH                      = "                <td>-</td>\n"
	HTTPS_PREFIX                      = "https://"
	HTTP_PORT_FORMAT                  = ":%d"
	HTTP_SCHEME                       = "http"
	HTTP_SERVER_RUNNING_ERROR         = "HTTP server is already running"
	HTTPS_SCHEME                      = "https"
	ICON_CLASS_DIR                    = "dir-icon"
	ICON_CLASS_FILE                   = "file-icon"
	INFO_FILE_SUFFIX                  = ".info"
	INVALID_PARAMETERS_ERROR          = "invalid number of parameters"
	KEY_FILE_NOT_FOUND_FORMAT         = "key file not found: %s"
	MAX_PORT                          = 65535
	MIME_DEFAULT                      = "application/octet-stream"
	MIN_PORT                          = 1
	MODIFIED_TIME_FORMAT              = "2006-01-02 15:04:05"
	NOT_FOUND_LOG_FORMAT              = "%s not found"
	ORGANIZATION_NAME                 = "Go Solution SDK"
	PATH_SEPARATOR                    = "/"
	PFX_EXTENSION                     = ".pfx"
	PFX_NO_VALID_CERT_BYTES_ERROR     = "no valid certificates found in PFX byte array"
	PFX_NO_VALID_CERT_ERROR           = "no valid certificates found in PFX file"
	POST                              = "POST"
	PRIVATE_KEY_BLOCK_TYPE            = "PRIVATE KEY"
	PUT                               = "PUT"
	ROUTE_TUSD_CONFLICT_FORMAT        = "pattern %s conflicts with TUSD mount point %s"
	SERVER_NOT_RUNNING_ERROR          = "server is not running"
	SET                               = "SET"
	SHOW                              = "SHOW"
	TCP_NETWORK                       = "tcp"
	TLS_SERVER_NOT_RUNNING_ERROR      = "TLS server is not running"
	TLS_SERVER_RUNNING_ERROR          = "TLS server is already running"
	TUSD_DEFAULT_BASE_PATH            = "/tusd/"
	TUSD_DEFAULT_STORE_PATH           = "./tusd"
	TUSD_MOUNT_CONFLICT_FORMAT        = "TUSD mount %s conflicts with existing route %s"
	TUSD_MOUNT_NOT_FOUND_FORMAT       = "TUSD mount %s does not exist"
	URL_FORMAT                        = "%s://%s/"
	WAIT_CLIENT_TIMEOUT               = 5 * time.Second
	WAIT_RETRY_INTERVAL               = 10 * time.Millisecond
	WILDCARD_SUFFIX                   = "/*"
	WRITE_ERROR_FORMAT                = "%s write error: %s"
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName,SpellCheckingInspection
type (
	AFTER_SERVE_STATIC_HANDLER  func(request *http.Request, response http.ResponseWriter, filePath string)
	BEFORE_REQUEST_HANDLER      func(request *http.Request, response http.ResponseWriter) bool
	BEFORE_ROUTE_HANDLER        func(request *http.Request, response http.ResponseWriter) bool
	BEFORE_SERVE_STATIC_HANDLER func(request *http.Request, response http.ResponseWriter) bool
	BEFORE_TUSD_HANDLER         func(request *http.Request, response http.ResponseWriter) bool
	REQUEST_HANDLER             func(request *http.Request, response http.ResponseWriter) error
	ROUTE                       struct {
		Method  string
		Pattern string
		Handler REQUEST_HANDLER
	}
	TUSD_MOUNT struct {
		Uri             string
		DirectoryPath   string
		MaxSize         int64
		EnableMoveFile  bool
		TargetDirectory string
		Handler         *tushandler.Handler
		UploadedHandler TUSD_UPLOADED_HANDLER
	}
	TUSD_UPLOADED_HANDLER func(basePath string, id string, filePath string, metaData map[string]string)
)

//goland:noinspection SpellCheckingInspection
var (
	afterServeStaticHandlers  []AFTER_SERVE_STATIC_HANDLER
	beforeRequestHandlers     []BEFORE_REQUEST_HANDLER
	beforeRouteHandlers       []BEFORE_ROUTE_HANDLER
	beforeServeStaticHandlers []BEFORE_SERVE_STATIC_HANDLER
	beforeTusdHandlers        []BEFORE_TUSD_HANDLER
	embedDirectories          map[string]embed.FS
	enableListDirectory       bool
	errorHandler              func(error)
	mutex                     sync.Mutex
	routes                    []ROUTE
	server                    *http.Server
	serverCancel              _context.CancelFunc
	serverContext             _context.Context
	serverListener            net.Listener
	staticDirectories         map[string]string
	tlsPort                   int
	tlsServer                 *http.Server
	tlsServerCancel           _context.CancelFunc
	tlsServerContext          _context.Context
	tlsServerListener         net.Listener
	tusdMounts                map[string]*TUSD_MOUNT
)

//goland:noinspection GoUnusedExportedFunction
func init() {
	serverContext, serverCancel = _context.WithCancel(_context.Background())
	tlsServerContext, tlsServerCancel = _context.WithCancel(_context.Background())
	server = &http.Server{
		Addr:         DEFAULT_ADDRESS,
		ReadTimeout:  DEFAULT_READ_TIMEOUT,
		WriteTimeout: DEFAULT_WRITE_TIMEOUT,
		IdleTimeout:  DEFAULT_IDLE_TIMEOUT,
	}
	tlsServer = &http.Server{
		Addr:         DEFAULT_TLS_ADDRESS,
		ReadTimeout:  DEFAULT_READ_TIMEOUT,
		WriteTimeout: DEFAULT_WRITE_TIMEOUT,
		IdleTimeout:  DEFAULT_IDLE_TIMEOUT,
	}
	staticDirectories = make(map[string]string)
	embedDirectories = make(map[string]embed.FS)
	tusdMounts = make(map[string]*TUSD_MOUNT)
	tusdMounts[TUSD_DEFAULT_BASE_PATH] = &TUSD_MOUNT{
		Uri:             TUSD_DEFAULT_BASE_PATH,
		DirectoryPath:   TUSD_DEFAULT_STORE_PATH,
		MaxSize:         0,
		EnableMoveFile:  true,
		TargetDirectory: "",
	}
}

//goland:noinspection GoUnusedExportedFunction
func AddEmbedDirectory(urlPath string, embedFs embed.FS) {
	mutex.Lock()
	defer mutex.Unlock()
	embedDirectories[urlPath] = embedFs
}

//goland:noinspection GoUnusedExportedFunction
func AddStaticDirectory(urlPath string, directoryPath string) {
	mutex.Lock()
	defer mutex.Unlock()
	staticDirectories[urlPath] = directoryPath
}

//goland:noinspection GoUnusedExportedFunction
func EnableDirectoryListing(enabled bool) {
	mutex.Lock()
	defer mutex.Unlock()
	enableListDirectory = enabled
}

//goland:noinspection DuplicatedCode,GoUnusedExportedFunction
func EnableRedirectToTLS() error {
	err := error(nil)
	mutex.Lock()
	defer mutex.Unlock()
	if tlsServerListener == nil {
		err = fmt.Errorf(TLS_SERVER_NOT_RUNNING_ERROR)
	} else if serverListener != nil {
		err = fmt.Errorf(HTTP_SERVER_RUNNING_ERROR)
	} else {
		httpAddress := fmt.Sprintf(HTTP_PORT_FORMAT, DEFAULT_HTTP_PORT)
		server.Addr = httpAddress
		server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			target := HTTPS_PREFIX + strings.Replace(r.Host, fmt.Sprintf(HTTP_PORT_FORMAT, DEFAULT_HTTP_PORT), fmt.Sprintf(HTTP_PORT_FORMAT, tlsPort), 1) + r.URL.Path
			if r.URL.RawQuery != "" {
				target += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, target, http.StatusMovedPermanently)
		})
		go func() {
			if ln, err := net.Listen(TCP_NETWORK, server.Addr); err == nil {
				mutex.Lock()
				serverListener = ln
				mutex.Unlock()
				if err := server.Serve(serverListener); !errors.Is(err, http.ErrServerClosed) {
					if errorHandler != nil {
						errorHandler(err)
					}
				}
			}
			mutex.Lock()
			serverListener = nil
			mutex.Unlock()
		}()
	}
	return err
}

func GenerateSelfSignedCertificate() (*tls.Certificate, error) {
	var privateKey *rsa.PrivateKey
	err := error(nil)
	var certificate *tls.Certificate
	var serialNumber *big.Int
	var serialErr error
	var certificateTemplate x509.Certificate
	var certificateDER []byte
	var privateKeyDER []byte
	var tlsCertificate tls.Certificate
	if privateKey, err = rsa.GenerateKey(rand.Reader, DEFAULT_RSA_KEY_BITS); err == nil {
		validFromTime := time.Now()
		validToTime := validFromTime.Add(DEFAULT_VALID_DAYS * 24 * time.Hour)
		serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), DEFAULT_SERIAL_NUMBER_BITS)
		if serialNumber, serialErr = rand.Int(rand.Reader, serialNumberLimit); serialErr == nil {
			certificateTemplate = x509.Certificate{
				SerialNumber: serialNumber,
				Subject: pkix.Name{
					Organization: []string{ORGANIZATION_NAME},
				},
				NotBefore:             validFromTime,
				NotAfter:              validToTime,
				KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
				ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
				BasicConstraintsValid: true,
				IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
			}
			if certificateDER, err = x509.CreateCertificate(rand.Reader, &certificateTemplate, &certificateTemplate, &privateKey.PublicKey, privateKey); err == nil {
				if privateKeyDER, err = x509.MarshalPKCS8PrivateKey(privateKey); err == nil {
					certificatePEM := new(bytes.Buffer)
					_ = pem.Encode(certificatePEM, &pem.Block{Type: CERTIFICATE_BLOCK_TYPE, Bytes: certificateDER})
					privateKeyPEM := new(bytes.Buffer)
					_ = pem.Encode(privateKeyPEM, &pem.Block{Type: PRIVATE_KEY_BLOCK_TYPE, Bytes: privateKeyDER})
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

//goland:noinspection GoUnusedExportedFunction
func GetFileExtensionName(mimeType string) string {
	var result string
	if fileExtensionName, ok := EXTENSION_MAP[strings.ToLower(mimeType)]; ok {
		result = fileExtensionName
	} else {
		parts := strings.Split(mimeType, "/")
		if len(parts) == 2 {
			subtype := parts[1]
			if strings.HasPrefix(subtype, "x-") {
				subtype = strings.TrimPrefix(subtype, "x-")
			}
			result = "." + subtype
		} else {
			result = ""
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetMimeType(fileExtensionName string) string {
	var result string
	if mimeType, ok := MIME_TYPES[strings.ToLower(fileExtensionName)]; ok {
		result = mimeType
	} else {
		result = MIME_DEFAULT
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetPort() int {
	var port int
	var portStr string
	err := error(nil)
	mutex.Lock()
	defer mutex.Unlock()
	_, portStr, err = net.SplitHostPort(server.Addr)
	if err == nil {
		port, err = strconv.Atoi(portStr)
		if err != nil {
			port = 0
		}
	} else {
		port = 0
	}
	return port
}

//goland:noinspection GoUnusedExportedFunction
func GetStaticDirectory(urlPath string) string {
	mutex.Lock()
	defer mutex.Unlock()
	return staticDirectories[urlPath]
}

//goland:noinspection GoUnusedExportedFunction
func GetTlsPort() int {
	var port int
	var portStr string
	err := error(nil)
	mutex.Lock()
	defer mutex.Unlock()
	_, portStr, err = net.SplitHostPort(tlsServer.Addr)
	if err == nil {
		port, err = strconv.Atoi(portStr)
		if err != nil {
			port = 0
		}
	} else {
		port = 0
	}
	return port
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func GetTusdMount(uri string) *TUSD_MOUNT {
	mutex.Lock()
	defer mutex.Unlock()
	result := tusdMounts[uri]
	return result
}

//goland:noinspection GoUnusedExportedFunction
func IsListening() bool {
	mutex.Lock()
	defer mutex.Unlock()
	return serverListener != nil
}

//goland:noinspection GoUnusedExportedFunction
func IsTlsListening() bool {
	mutex.Lock()
	defer mutex.Unlock()
	return tlsServerListener != nil
}

//goland:noinspection DuplicatedCode,GoUnusedExportedFunction
func ListenAsync() error {
	err := error(nil)
	mutex.Lock()
	if serverListener == nil {
		err = initTusdMounts()
		if err == nil {
			go func() {
				if listener, err := net.Listen(TCP_NETWORK, server.Addr); err == nil {
					mutex.Lock()
					serverListener = listener
					mutex.Unlock()
					if err := server.Serve(serverListener); !errors.Is(err, http.ErrServerClosed) {
						if errorHandler != nil {
							errorHandler(err)
						}
					}
				}
				mutex.Lock()
				serverListener = nil
				mutex.Unlock()
			}()
		}
	} else {
		err = fmt.Errorf(HTTP_SERVER_RUNNING_ERROR)
	}
	mutex.Unlock()
	if err == nil {
		waitServe(server.Addr, false)
		server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handleRequest(w, r)
		})
	}
	return err
}

//goland:noinspection GoUnusedExportedFunction
func ListenTlsAsync(address string, params ...interface{}) error {
	err := error(nil)
	mutex.Lock()
	if tlsServerListener != nil {
		err = fmt.Errorf(TLS_SERVER_RUNNING_ERROR)
	} else {
		err = initTusdMounts()
		if err == nil {
			tlsServer.Addr = address
			tlsServer.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handleRequest(w, r)
			})
			go func() {
				var err error
				var certPath, keyPath, password string
				var certBytes, keyBytes []byte
				var ok bool
				var fileExt string
				if len(params) == 0 {
					var certificate *tls.Certificate
					if certificate, err = GenerateSelfSignedCertificate(); err == nil && certificate != nil {
						tlsServer.TLSConfig = &tls.Config{
							Certificates: []tls.Certificate{*certificate},
						}
						var ln net.Listener
						if ln, err = tls.Listen(TCP_NETWORK, tlsServer.Addr, tlsServer.TLSConfig); err == nil {
							mutex.Lock()
							tlsServerListener = ln
							mutex.Unlock()
							err = tlsServer.Serve(tlsServerListener)
						}
					} else if certificate == nil {
						err = fmt.Errorf(CERTIFICATE_NIL_ERROR)
					}
				} else if len(params) >= 2 {
					if certPath, ok = params[0].(string); ok {
						if password, ok = params[1].(string); ok {
							if _, err = os.Stat(certPath); os.IsNotExist(err) {
								err = fmt.Errorf(CERTIFICATE_FILE_NOT_FOUND_FORMAT, certPath)
							} else {
								fileExt = strings.ToLower(filepath.Ext(certPath))
								if fileExt == PFX_EXTENSION {
									var pfxBytes []byte
									var tlsCerts []tls.Certificate
									if pfxBytes, err = os.ReadFile(certPath); err == nil {
										if tlsCerts, err = x509util.Load(pfxBytes, password); err == nil {
											if len(tlsCerts) > 0 {
												tlsServer.TLSConfig = &tls.Config{
													Certificates: tlsCerts,
												}
												var ln net.Listener
												if ln, err = tls.Listen(TCP_NETWORK, tlsServer.Addr, tlsServer.TLSConfig); err == nil {
													mutex.Lock()
													tlsServerListener = ln
													mutex.Unlock()
													go waitServe(tlsServer.Addr, true)
													err = tlsServer.Serve(tlsServerListener)
												}
											} else {
												err = fmt.Errorf(PFX_NO_VALID_CERT_ERROR)
											}
										}
									}
								} else {
									keyPath = password
									if _, err = os.Stat(keyPath); os.IsNotExist(err) {
										err = fmt.Errorf(KEY_FILE_NOT_FOUND_FORMAT, keyPath)
									} else {
										var ln net.Listener
										if ln, err = tls.Listen(TCP_NETWORK, tlsServer.Addr, tlsServer.TLSConfig); err == nil {
											mutex.Lock()
											tlsServerListener = ln
											mutex.Unlock()
											err = tlsServer.Serve(tlsServerListener)
										}
									}
								}
							}
						} else if certBytes, ok = params[0].([]byte); ok {
							if password, ok = params[1].(string); ok {
								var tlsCerts []tls.Certificate
								if tlsCerts, err = x509util.Load(certBytes, password); err == nil {
									if len(tlsCerts) > 0 {
										tlsServer.TLSConfig = &tls.Config{
											Certificates: tlsCerts,
										}
										var ln net.Listener
										if ln, err = tls.Listen(TCP_NETWORK, tlsServer.Addr, tlsServer.TLSConfig); err == nil {
											mutex.Lock()
											tlsServerListener = ln
											mutex.Unlock()
											err = tlsServer.Serve(tlsServerListener)
										}
									} else {
										err = fmt.Errorf(PFX_NO_VALID_CERT_BYTES_ERROR)
									}
								}
							} else if keyBytes, ok = params[1].([]byte); ok {
								var cert tls.Certificate
								if cert, err = tls.X509KeyPair(certBytes, keyBytes); err == nil {
									tlsServer.TLSConfig = &tls.Config{
										Certificates: []tls.Certificate{cert},
									}
									var ln net.Listener
									if ln, err = tls.Listen(TCP_NETWORK, tlsServer.Addr, tlsServer.TLSConfig); err == nil {
										mutex.Lock()
										tlsServerListener = ln
										mutex.Unlock()
										err = tlsServer.Serve(tlsServerListener)
									}
								}
							}
						}
					}
				} else {
					err = fmt.Errorf(INVALID_PARAMETERS_ERROR)
				}
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					mutex.Lock()
					handler := errorHandler
					mutex.Unlock()
					if handler != nil {
						handler(err)
					}
				}
				mutex.Lock()
				tlsServerListener = nil
				mutex.Unlock()
			}()
		}
	}
	mutex.Unlock()
	if err == nil {
		waitServe(tlsServer.Addr, true)
	}
	return err
}

//goland:noinspection GoUnusedExportedFunction
func On(method string, pattern string, handler REQUEST_HANDLER) error {
	err := error(nil)
	mutex.Lock()
	defer mutex.Unlock()
	for basePath := range tusdMounts {
		if isPathConflictWithTusd(pattern, basePath) {
			err = fmt.Errorf(ROUTE_TUSD_CONFLICT_FORMAT, pattern, basePath)
			break
		}
	}
	if err == nil {
		routes = append(routes, ROUTE{
			Method:  method,
			Pattern: pattern,
			Handler: handler,
		})
	}
	return err
}

//goland:noinspection GoUnusedExportedFunction
func OnAfterStaticServe(handler AFTER_SERVE_STATIC_HANDLER) {
	mutex.Lock()
	defer mutex.Unlock()
	afterServeStaticHandlers = append(afterServeStaticHandlers, handler)
}

//goland:noinspection GoUnusedExportedFunction
func OnBeforeRequest(callback BEFORE_REQUEST_HANDLER) {
	mutex.Lock()
	defer mutex.Unlock()
	beforeRequestHandlers = append(beforeRequestHandlers, callback)
}

//goland:noinspection GoUnusedExportedFunction
func OnBeforeRoute(callback BEFORE_ROUTE_HANDLER) {
	mutex.Lock()
	defer mutex.Unlock()
	beforeRouteHandlers = append(beforeRouteHandlers, callback)
}

//goland:noinspection GoUnusedExportedFunction
func OnBeforeStaticServe(handler BEFORE_SERVE_STATIC_HANDLER) {
	mutex.Lock()
	defer mutex.Unlock()
	beforeServeStaticHandlers = append(beforeServeStaticHandlers, handler)
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func OnBeforeTusd(handler BEFORE_TUSD_HANDLER) {
	mutex.Lock()
	defer mutex.Unlock()
	beforeTusdHandlers = append(beforeTusdHandlers, handler)
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func OnTusd(uri string, directoryPath string, handler TUSD_UPLOADED_HANDLER) error {
	err := error(nil)
	mutex.Lock()
	defer mutex.Unlock()
	if !strings.HasPrefix(uri, PATH_SEPARATOR) {
		uri = PATH_SEPARATOR + uri
	}
	if !strings.HasSuffix(uri, PATH_SEPARATOR) {
		uri = uri + PATH_SEPARATOR
	}
	for _, route := range routes {
		if isPathConflictWithTusd(route.Pattern, uri) {
			err = fmt.Errorf(TUSD_MOUNT_CONFLICT_FORMAT, uri, route.Pattern)
			break
		}
	}
	if err == nil {
		tusdMounts[uri] = &TUSD_MOUNT{
			Uri:             uri,
			DirectoryPath:   directoryPath,
			MaxSize:         0,
			EnableMoveFile:  true,
			TargetDirectory: "",
			UploadedHandler: handler,
		}
	}
	return err
}

//goland:noinspection GoUnusedExportedFunction
func OnError(handler func(error)) {
	mutex.Lock()
	defer mutex.Unlock()
	errorHandler = handler
}

//goland:noinspection GoUnusedExportedFunction
func RemoveStaticDirectory(urlPath string) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(staticDirectories, urlPath)
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func RemoveTusdMount(uri string) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(tusdMounts, uri)
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func SetTusdMount(basePath string, maxSize int64, enableMoveFile bool, targetDirectory string) error {
	err := error(nil)
	mutex.Lock()
	defer mutex.Unlock()
	if mount, ok := tusdMounts[basePath]; ok {
		mount.MaxSize = maxSize
		mount.EnableMoveFile = enableMoveFile
		mount.TargetDirectory = targetDirectory
	} else {
		err = fmt.Errorf(TUSD_MOUNT_NOT_FOUND_FORMAT, basePath)
	}
	return err
}

//goland:noinspection GoUnusedExportedFunction
func SetAddress(address string) {
	mutex.Lock()
	defer mutex.Unlock()
	server.Addr = address
}

//goland:noinspection GoUnusedExportedFunction
func SetContext(ctx _context.Context) {
	mutex.Lock()
	if serverCancel != nil {
		serverCancel()
	}
	if tlsServerCancel != nil {
		tlsServerCancel()
	}
	newServerCtx, serverCancelFunc := _context.WithCancel(ctx)
	newTlsServerCtx, tlsServerCancelFunc := _context.WithCancel(ctx)
	serverContext = newServerCtx
	tlsServerContext = newTlsServerCtx
	serverCancel = serverCancelFunc
	tlsServerCancel = tlsServerCancelFunc
	mutex.Unlock()
}

//goland:noinspection GoUnusedExportedFunction
func SetTlsAddress(address string) {
	mutex.Lock()
	defer mutex.Unlock()
	tlsServer.Addr = address
}

//goland:noinspection GoUnusedExportedFunction
func SetReadTimeout(timeout time.Duration) {
	mutex.Lock()
	defer mutex.Unlock()
	server.ReadTimeout = timeout
}

//goland:noinspection GoUnusedExportedFunction
func GetReadTimeout() time.Duration {
	mutex.Lock()
	defer mutex.Unlock()
	return server.ReadTimeout
}

//goland:noinspection GoUnusedExportedFunction
func SetWriteTimeout(timeout time.Duration) {
	mutex.Lock()
	defer mutex.Unlock()
	server.WriteTimeout = timeout
}

//goland:noinspection GoUnusedExportedFunction
func GetWriteTimeout() time.Duration {
	mutex.Lock()
	defer mutex.Unlock()
	return server.WriteTimeout
}

//goland:noinspection GoUnusedExportedFunction
func SetIdleTimeout(timeout time.Duration) {
	mutex.Lock()
	defer mutex.Unlock()
	server.IdleTimeout = timeout
}

//goland:noinspection GoUnusedExportedFunction
func GetIdleTimeout() time.Duration {
	mutex.Lock()
	defer mutex.Unlock()
	return server.IdleTimeout
}

//goland:noinspection GoUnusedExportedFunction
func SetTlsReadTimeout(timeout time.Duration) {
	mutex.Lock()
	defer mutex.Unlock()
	tlsServer.ReadTimeout = timeout
}

//goland:noinspection GoUnusedExportedFunction
func SetTlsWriteTimeout(timeout time.Duration) {
	mutex.Lock()
	defer mutex.Unlock()
	tlsServer.WriteTimeout = timeout
}

//goland:noinspection GoUnusedExportedFunction
func SetTlsIdleTimeout(timeout time.Duration) {
	mutex.Lock()
	defer mutex.Unlock()
	tlsServer.IdleTimeout = timeout
}

//goland:noinspection GoUnusedExportedFunction
func SetReadTimeouts(timeout time.Duration) {
	mutex.Lock()
	defer mutex.Unlock()
	server.ReadTimeout = timeout
	tlsServer.ReadTimeout = timeout
}

//goland:noinspection GoUnusedExportedFunction
func SetWriteTimeouts(timeout time.Duration) {
	mutex.Lock()
	defer mutex.Unlock()
	server.WriteTimeout = timeout
	tlsServer.WriteTimeout = timeout
}

//goland:noinspection GoUnusedExportedFunction
func SetIdleTimeouts(timeout time.Duration) {
	mutex.Lock()
	defer mutex.Unlock()
	server.IdleTimeout = timeout
	tlsServer.IdleTimeout = timeout
}

//goland:noinspection GoUnusedExportedFunction
func SetTimeouts(readTimeout, writeTimeout, idleTimeout time.Duration) {
	mutex.Lock()
	defer mutex.Unlock()
	server.ReadTimeout = readTimeout
	server.WriteTimeout = writeTimeout
	server.IdleTimeout = idleTimeout
	tlsServer.ReadTimeout = readTimeout
	tlsServer.WriteTimeout = writeTimeout
	tlsServer.IdleTimeout = idleTimeout
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func Shutdown() {
	shutdown()
	shutdownTls()
}

//goland:noinspection SpellCheckingInspection
func formatFileSize(size int64) string {
	const unit = FILE_SIZE_UNIT
	if size < unit {
		return fmt.Sprintf(FILE_SIZE_BYTE_FORMAT, size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf(FILE_SIZE_FORMAT, float64(size)/float64(div), BYTE_UNITS[exp])
}

func findEmbedFilePath(embedFs embed.FS, path string) (fs.FileInfo, string, error) {
	var stat fs.FileInfo
	resultPath := ""
	err := error(nil)
	if stat, err = fs.Stat(embedFs, path); err == nil {
		resultPath = path
	} else {
		var entries []fs.DirEntry
		entries, err = fs.ReadDir(embedFs, CURRENT_DIRECTORY)
		if err == nil && len(entries) == 1 && entries[0].IsDir() {
			rootDir := entries[0].Name()
			tryPath := rootDir
			if path != CURRENT_DIRECTORY {
				tryPath = rootDir + PATH_SEPARATOR + path
			}
			stat, err = fs.Stat(embedFs, tryPath)
			if err == nil {
				resultPath = tryPath
			}
		}
	}
	if err != nil {
		err = os.ErrNotExist
	}
	return stat, resultPath, err
}

//goland:noinspection SpellCheckingInspection,DuplicatedCode
func getDirectoryListHTML(path string, urlPath string, embedFS *embed.FS) string {
	result := ""
	var entries []fs.DirEntry
	var err error
	if embedFS != nil {
		entries, err = fs.ReadDir(*embedFS, path)
	} else {
		entries, err = os.ReadDir(path)
	}
	if err == nil {
		var html strings.Builder
		html.WriteString(HTML_DOCTYPE)
		html.WriteString(HTML_HTML_OPEN)
		html.WriteString(HTML_HEAD_OPEN)
		html.WriteString(HTML_META_CHARSET)
		html.WriteString(HTML_META_VIEWPORT)
		html.WriteString(fmt.Sprintf(DIRECTORY_LIST_TITLE_FORMAT, urlPath))
		html.WriteString(HTML_STYLE_OPEN)
		html.WriteString(HTML_STYLE_BODY)
		html.WriteString(HTML_STYLE_H1)
		html.WriteString(HTML_STYLE_TABLE)
		html.WriteString(HTML_STYLE_TH_TD)
		html.WriteString(HTML_STYLE_TH)
		html.WriteString(HTML_STYLE_TR_HOVER)
		html.WriteString(HTML_STYLE_A)
		html.WriteString(HTML_STYLE_A_HOVER)
		html.WriteString(HTML_STYLE_DIR_ICON)
		html.WriteString(HTML_STYLE_FILE_ICON)
		html.WriteString(HTML_STYLE_PARENT_DIR)
		html.WriteString(HTML_STYLE_CLOSE)
		html.WriteString(HTML_HEAD_CLOSE)
		html.WriteString(HTML_BODY_OPEN)
		html.WriteString(fmt.Sprintf(DIRECTORY_LIST_HEADING_FORMAT, urlPath))
		html.WriteString(HTML_TABLE_OPEN)
		html.WriteString(HTML_THEAD_OPEN)
		html.WriteString(HTML_TR_OPEN)
		html.WriteString(HTML_TH_NAME)
		html.WriteString(HTML_TH_SIZE)
		html.WriteString(HTML_TH_LAST_MODIFIED)
		html.WriteString(HTML_TR_CLOSE)
		html.WriteString(HTML_THEAD_CLOSE)
		html.WriteString(HTML_TBODY_OPEN)
		if urlPath != PATH_SEPARATOR {
			parentPath := strings.TrimSuffix(urlPath, PATH_SEPARATOR)
			lastSlash := strings.LastIndex(parentPath, PATH_SEPARATOR)
			if lastSlash >= 0 {
				parentPath = parentPath[:lastSlash]
				if parentPath == "" {
					parentPath = PATH_SEPARATOR
				}
			}
			html.WriteString(HTML_TR_OPEN)
			html.WriteString(fmt.Sprintf(DIRECTORY_LIST_PARENT_FORMAT, parentPath))
			html.WriteString(HTML_TD_DASH)
			html.WriteString(HTML_TD_DASH)
			html.WriteString(HTML_TR_CLOSE)
		}
		for _, entry := range entries {
			info, entryErr := entry.Info()
			if entryErr != nil {
				continue
			}
			name := info.Name()
			size := ""
			if !info.IsDir() {
				size = formatFileSize(info.Size())
			}
			lastModified := info.ModTime().Format(MODIFIED_TIME_FORMAT)
			linkPath := strings.TrimSuffix(urlPath, PATH_SEPARATOR) + PATH_SEPARATOR + name
			className := ICON_CLASS_FILE
			if info.IsDir() {
				className = ICON_CLASS_DIR
			}
			html.WriteString(HTML_TR_OPEN)
			html.WriteString(fmt.Sprintf(DIRECTORY_LIST_ENTRY_FORMAT, linkPath, className, name))
			html.WriteString(fmt.Sprintf(DIRECTORY_LIST_TD_FORMAT, size))
			html.WriteString(fmt.Sprintf(DIRECTORY_LIST_TD_FORMAT, lastModified))
			html.WriteString(HTML_TR_CLOSE)
		}
		html.WriteString(HTML_TBODY_CLOSE)
		html.WriteString(HTML_TABLE_CLOSE)
		html.WriteString(HTML_BODY_CLOSE)
		html.WriteString(HTML_HTML_CLOSE)
		result = html.String()
	}
	return result
}

//goland:noinspection GoBoolExpressions,SpellCheckingInspection
func handleRequest(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	_routes := make([]ROUTE, len(routes))
	_beforeRequestHandlers := make([]BEFORE_REQUEST_HANDLER, len(beforeRequestHandlers))
	_beforeRouteHandlers := make([]BEFORE_ROUTE_HANDLER, len(beforeRouteHandlers))
	_beforeServeStaticHandlers := make([]BEFORE_SERVE_STATIC_HANDLER, len(beforeServeStaticHandlers))
	_beforeTusdHandlers := make([]BEFORE_TUSD_HANDLER, len(beforeTusdHandlers))
	_afterServeStaticHandlers := make([]AFTER_SERVE_STATIC_HANDLER, len(afterServeStaticHandlers))
	_tusdMounts := make(map[string]*TUSD_MOUNT)
	for k, v := range tusdMounts {
		_tusdMounts[k] = v
	}
	copy(_routes, routes)
	copy(_beforeRequestHandlers, beforeRequestHandlers)
	copy(_beforeRouteHandlers, beforeRouteHandlers)
	copy(_beforeServeStaticHandlers, beforeServeStaticHandlers)
	copy(_beforeTusdHandlers, beforeTusdHandlers)
	copy(_afterServeStaticHandlers, afterServeStaticHandlers)
	mutex.Unlock()
	isForbidden := false
	if !isForbidden {
		for _, handler := range _beforeRequestHandlers {
			if !handler(r, w) {
				isForbidden = true
				break
			}
		}
	}
	isTusdHandled := false
	if !isForbidden {
		if mount := matchTusdMount(_tusdMounts, r.URL.Path); mount != nil && mount.Handler != nil {
			isTusdAllow := true
			for _, handler := range _beforeTusdHandlers {
				if !handler(r, w) {
					isTusdAllow = false
					break
				}
			}
			if isTusdAllow {
				mount.Handler.ServeHTTP(w, r)
			}
			isTusdHandled = true
		}
	}
	if !isForbidden && !isTusdHandled {
		isRouteHandled := false
		for _, handler := range _beforeRouteHandlers {
			if !handler(r, w) {
				isForbidden = true
				break
			}
		}
		if !isForbidden {
			for _, route := range _routes {
				if (route.Method == ANY_METHOD || route.Method == r.Method) && matchPath(route.Pattern, r.URL.Path) {
					if err := route.Handler(r, w); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						logger.Logger.Error(err.Error(), 1)
					}
					isRouteHandled = true
					break
				}
			}
		}
		if !isRouteHandled {
			ok := true
			if _beforeServeStaticHandlers != nil {
				for _, handler := range _beforeServeStaticHandlers {
					ok = handler(r, w)
					if !ok {
						break
					}
				}
			}
			if ok {
				serveStatic(w, r)
			}
			if _afterServeStaticHandlers != nil {
				for _, handler := range _afterServeStaticHandlers {
					handler(r, w, r.URL.Path)
				}
			}
		}
	}
}

//goland:noinspection SpellCheckingInspection,DuplicatedCode
func handleTusdComplete(mount *TUSD_MOUNT, event tushandler.HookEvent) {
	sourcePath := filepath.Join(mount.DirectoryPath, event.Upload.ID)
	finalPath := sourcePath
	if mount.EnableMoveFile {
		targetDir := mount.DirectoryPath
		if mount.TargetDirectory != "" {
			targetDir = mount.TargetDirectory
		}
		fileName := event.Upload.ID
		if name, ok := event.Upload.MetaData[FILENAME_KEY]; ok && name != "" {
			fileName = name
		}
		_ = os.MkdirAll(targetDir, DEFAULT_DIRECTORY_PERMISSION)
		targetPath := filepath.Join(targetDir, fileName)
		if _, err := os.Stat(targetPath); err == nil {
			ext := filepath.Ext(fileName)
			base := fileName[:len(fileName)-len(ext)]
			counter := 1
			for {
				newFileName := fmt.Sprintf(FILE_RENAME_FORMAT, base, counter, ext)
				targetPath = filepath.Join(targetDir, newFileName)
				if _, statErr := os.Stat(targetPath); os.IsNotExist(statErr) {
					break
				}
				counter++
			}
		}
		if err := os.Rename(sourcePath, targetPath); err == nil {
			infoPath := filepath.Join(mount.DirectoryPath, event.Upload.ID+INFO_FILE_SUFFIX)
			_ = os.Remove(infoPath)
			finalPath = targetPath
		}
	}
	if mount.UploadedHandler != nil {
		mount.UploadedHandler(mount.Uri, event.Upload.ID, finalPath, event.Upload.MetaData)
	}
}

//goland:noinspection SpellCheckingInspection
func initTusdMounts() error {
	err := error(nil)
	for basePath, mount := range tusdMounts {
		if mount.Handler == nil {
			if mkdirErr := os.MkdirAll(mount.DirectoryPath, DEFAULT_DIRECTORY_PERMISSION); mkdirErr == nil {
				store := filestore.New(mount.DirectoryPath)
				composer := tushandler.NewStoreComposer()
				store.UseIn(composer)
				handlerConfig := tushandler.Config{
					StoreComposer:         composer,
					BasePath:              basePath,
					MaxSize:               mount.MaxSize,
					NotifyCompleteUploads: true,
				}
				if tusdHandlerInstance, handlerErr := tushandler.NewHandler(handlerConfig); handlerErr == nil {
					mount.Handler = tusdHandlerInstance
					currentMount := mount
					go func() {
						for event := range currentMount.Handler.CompleteUploads {
							handleTusdComplete(currentMount, event)
						}
					}()
				} else {
					err = handlerErr
					break
				}
			} else {
				err = mkdirErr
				break
			}
		}
	}
	return err
}

//goland:noinspection SpellCheckingInspection
func isPathConflictWithTusd(pattern string, tusdBasePath string) bool {
	result := false
	if pattern == tusdBasePath || pattern == strings.TrimSuffix(tusdBasePath, PATH_SEPARATOR) {
		result = true
	} else if strings.HasPrefix(pattern, tusdBasePath) {
		result = true
	} else if strings.HasSuffix(pattern, WILDCARD_SUFFIX) {
		prefix := strings.TrimSuffix(pattern, WILDCARD_SUFFIX)
		if strings.HasPrefix(tusdBasePath, prefix+PATH_SEPARATOR) || tusdBasePath == prefix+PATH_SEPARATOR {
			result = true
		}
	}
	return result
}

func matchPath(pattern string, path string) bool {
	var matched bool
	if pattern == path {
		matched = true
	} else if strings.HasSuffix(pattern, WILDCARD_SUFFIX) {
		prefix := strings.TrimSuffix(pattern, WILDCARD_SUFFIX)
		if prefix == "" || strings.HasPrefix(path, prefix+PATH_SEPARATOR) || path == prefix {
			matched = true
		}
	} else {
		if pattern != PATH_SEPARATOR && path != PATH_SEPARATOR {
			patternParts := strings.Split(pattern, PATH_SEPARATOR)
			pathParts := strings.Split(path, PATH_SEPARATOR)
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

//goland:noinspection SpellCheckingInspection
func matchTusdMount(mounts map[string]*TUSD_MOUNT, path string) *TUSD_MOUNT {
	var result *TUSD_MOUNT
	bestLength := 0
	for basePath, mount := range mounts {
		if strings.HasPrefix(path, basePath) || path == strings.TrimSuffix(basePath, PATH_SEPARATOR) {
			if len(basePath) > bestLength {
				result = mount
				bestLength = len(basePath)
			}
		}
	}
	return result
}

//goland:noinspection
func serveStatic(w http.ResponseWriter, r *http.Request) {
	httpStatusCode := http.StatusNotFound
	writeErrorMsg := ""
	isEmbedFs := false
	physicalFilePath := ""
	embedFs := embed.FS{}
	isDirectory := false
	mutex.Lock()
	_staticDirectories := make(map[string]string)
	for k, v := range staticDirectories {
		_staticDirectories[k] = v
	}
	_embedDirectories := make(map[string]embed.FS)
	for k, v := range embedDirectories {
		_embedDirectories[k] = v
	}
	mutex.Unlock()
	for urlPath, dirPath := range _staticDirectories {
		if strings.HasPrefix(r.URL.Path, urlPath) {
			absoluteDirectoryPath, _ := filepath.Abs(dirPath)
			if absoluteDirectoryPath == "" {
				absoluteDirectoryPath = dirPath
			}
			absoluteFilePath := filepath.Join(absoluteDirectoryPath, strings.TrimPrefix(r.URL.Path, urlPath))
			if stat, err := os.Stat(absoluteFilePath); !os.IsNotExist(err) {
				if stat.IsDir() {
					defaultIndexHTML := filepath.Join(absoluteFilePath, DEFAULT_INDEX_HTML)
					if _, err := os.Stat(defaultIndexHTML); !os.IsNotExist(err) {
						physicalFilePath = defaultIndexHTML
						isEmbedFs = false
						httpStatusCode = http.StatusOK
					} else if enableListDirectory {
						isDirectory = true
						isEmbedFs = false
						physicalFilePath = absoluteFilePath
						httpStatusCode = http.StatusOK
					} else {
						httpStatusCode = http.StatusForbidden
					}
				} else {
					physicalFilePath = absoluteFilePath
					isEmbedFs = false
					httpStatusCode = http.StatusOK
				}
				break
			}
		}
	}
	if physicalFilePath == "" {
		for urlPath, e := range _embedDirectories {
			if strings.HasPrefix(r.URL.Path, urlPath) {
				requestFilePath := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, urlPath), PATH_SEPARATOR)
				if requestFilePath == "" {
					requestFilePath = CURRENT_DIRECTORY
				}
				if stat, embedFilePath, err := findEmbedFilePath(e, requestFilePath); err == nil {
					if stat.IsDir() || strings.HasSuffix(requestFilePath, PATH_SEPARATOR) {
						defaultIndexHTML := ""
						if requestFilePath == CURRENT_DIRECTORY {
							defaultIndexHTML = DEFAULT_INDEX_HTML
						} else {
							defaultIndexHTML = strings.TrimSuffix(requestFilePath, PATH_SEPARATOR) + PATH_SEPARATOR + DEFAULT_INDEX_HTML
						}
						if _, actualIndexPath, err := findEmbedFilePath(e, defaultIndexHTML); err == nil {
							physicalFilePath = actualIndexPath
							isEmbedFs = true
							embedFs = e
							httpStatusCode = http.StatusOK
						} else if enableListDirectory {
							isDirectory = true
							isEmbedFs = true
							embedFs = e
							physicalFilePath = embedFilePath
							httpStatusCode = http.StatusOK
						} else {
							httpStatusCode = http.StatusForbidden
						}
					} else {
						physicalFilePath = embedFilePath
						isEmbedFs = true
						embedFs = e
						httpStatusCode = http.StatusOK
					}
					break
				}
			}
		}
	}
	if httpStatusCode == http.StatusOK {
		if isDirectory {
			var pEmbedFs *embed.FS
			if isEmbedFs {
				pEmbedFs = &embedFs
			}
			if html := getDirectoryListHTML(physicalFilePath, r.URL.Path, pEmbedFs); html == "" {
				httpStatusCode = http.StatusInternalServerError
			} else {
				w.Header().Set(CONTENT_TYPE, CONTENT_TYPE_HTML_UTF8)
				if _, err := w.Write([]byte(html)); err != nil {
					httpStatusCode = http.StatusInternalServerError
					writeErrorMsg = fmt.Sprintf(WRITE_ERROR_FORMAT, r.URL.Path, err.Error())
				}
			}
		} else {
			if isEmbedFs {
				http.ServeFileFS(w, r, embedFs, physicalFilePath)
			} else {
				http.ServeFile(w, r, physicalFilePath)
			}
		}
	}
	switch httpStatusCode {
	case http.StatusNotFound:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		logger.Logger.Error(fmt.Sprintf(NOT_FOUND_LOG_FORMAT, r.URL.Path))
	case http.StatusForbidden:
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		logger.Logger.Error(fmt.Sprintf(FORBIDDEN_LOG_FORMAT, r.URL.Path))
	case http.StatusInternalServerError:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		logger.Logger.Error(writeErrorMsg)
	}
}

//goland:noinspection GoUnusedFunction
func shutdown() error {
	err := error(nil)
	mutex.Lock()
	if serverListener == nil {
		err = fmt.Errorf(SERVER_NOT_RUNNING_ERROR)
	}
	mutex.Unlock()
	if err == nil {
		err = server.Shutdown(serverContext)
		mutex.Lock()
		serverListener = nil
		mutex.Unlock()
	}
	return err
}

//goland:noinspection GoUnusedFunction
func shutdownTls() error {
	err := error(nil)
	mutex.Lock()
	if tlsServerListener == nil {
		err = fmt.Errorf(TLS_SERVER_NOT_RUNNING_ERROR)
	}
	mutex.Unlock()
	if err == nil {
		err = tlsServer.Shutdown(tlsServerContext)
		mutex.Lock()
		tlsServerListener = nil
		mutex.Unlock()
	}
	return err
}

func waitServe(address string, isTLS bool) {
	client := &http.Client{
		Timeout: WAIT_CLIENT_TIMEOUT,
	}
	if isTLS {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	scheme := HTTP_SCHEME
	if isTLS {
		scheme = HTTPS_SCHEME
	}
	url := fmt.Sprintf(URL_FORMAT, scheme, address)
	for {
		if resp, err := client.Get(url); err == nil {
			_ = resp.Body.Close()
			break
		}
		time.Sleep(WAIT_RETRY_INTERVAL)
	}
}
