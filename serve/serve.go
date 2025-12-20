// Package serve
// File:        serve.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/serve/serve.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: SERVE provides HTTP server functionality with automatic TLS support
// --------------------------------------------------------------------------------
package serve

import (
	"bytes"
	__context "context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"embed"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math/big"
	"mime"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tus/tusd/v2/pkg/filelocker"
	"github.com/tus/tusd/v2/pkg/filestore"
	tushandler "github.com/tus/tusd/v2/pkg/handler"
	"github.com/xiang-tai-duo/go-boost/ca"
	"github.com/xiang-tai-duo/go-boost/hash"
	"github.com/xiang-tai-duo/go-boost/logger"
	"github.com/xiang-tai-duo/go-boost/network"
	"github.com/xiang-tai-duo/go-boost/sqlite"
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName,SpellCheckingInspection
type (
	AFTER_SERVE_STATIC_HANDLER  func(request *http.Request, response http.ResponseWriter, filePath string)
	BEFORE_REQUEST_HANDLER      func(request *http.Request, response http.ResponseWriter) bool
	BEFORE_ROUTE_HANDLER        func(request *http.Request, response http.ResponseWriter) bool
	BEFORE_SERVE_STATIC_HANDLER func(request *http.Request, response http.ResponseWriter) bool
	BEFORE_TUSD_HANDLER         func(request *http.Request, response http.ResponseWriter) bool
	HTTP_RESPONSE               struct {
		Message string `json:"message,omitempty"`
		Success bool   `json:"success"`
	}
	REQUEST_HANDLER func(request *http.Request, response http.ResponseWriter) error
	ROUTE           struct {
		Handler REQUEST_HANDLER
		Method  string
		Pattern string
	}
	SERVE_CONFIG struct {
		RevokedCertificates []string `json:"revokedCertificates"`
	}
	STATE_SNAPSHOT struct {
		AfterServeStaticHandlers  []AFTER_SERVE_STATIC_HANDLER
		BeforeRequestHandlers     []BEFORE_REQUEST_HANDLER
		BeforeRouteHandlers       []BEFORE_ROUTE_HANDLER
		BeforeServeStaticHandlers []BEFORE_SERVE_STATIC_HANDLER
		BeforeTusdHandlers        []BEFORE_TUSD_HANDLER
		EmbedDirectories          map[string]embed.FS
		ProtectedPatterns         []string
		IgnoredProtectionPatterns []string
		Routes                    []ROUTE
		StaticDirectories         map[string]string
		TusdMounts                map[string]*TUSD_MOUNT
	}
	TUSD_MOUNT struct {
		DirectoryPath   string
		Handler         *tushandler.Handler
		MaxSize         int64
		TargetDirectory string
		UploadedHandler TUSD_UPLOADED_HANDLER
		Uri             string
	}
	TUSD_UPLOADED_HANDLER func(uri string, uploadedFilePath string, fileName string, fileInfo tushandler.FileInfo)
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst,GoNameStartsWithPackageName,SpellCheckingInspection,SqlResolve,SqlDialectInspection
const (
	ANY_METHOD                        = "*"
	AUTOINCREMENT                     = "AUTOINCREMENT"
	BYTE_UNITS                        = "KMGTPE"
	CERTIFICATE_BLOCK_TYPE            = "CERTIFICATE"
	CONTENT_TYPE                      = "Content-Type"
	CONTENT_TYPE_HTML_UTF8            = "text/html; charset=utf-8"
	CONTENT_TYPE_JSON                 = "application/json"
	DEFAULT_ADDRESS                   = "127.0.0.1:80"
	DEFAULT_DIRECTORY_PERMISSION      = 0755
	DEFAULT_HTTP_PORT                 = 80
	DEFAULT_IDLE_TIMEOUT              = 60 * time.Second
	DEFAULT_INDEX_HTML                = "index.html"
	DEFAULT_READ_TIMEOUT              = 15 * time.Second
	DEFAULT_RSA_KEY_BITS              = 2048
	DEFAULT_SERIAL_NUMBER_BITS        = 128
	DEFAULT_STATIC_ROUTE              = "/*"
	DEFAULT_PROTECTED_PATTERN         = "/api/*"
	DEFAULT_TLS_ADDRESS               = ":443"
	DEFAULT_VALID_DAYS                = 365
	DEFAULT_WRITE_TIMEOUT             = 15 * time.Second
	DELETE                            = "DELETE"
	DIRECTORY_LIST_ENTRY_FORMAT       = "                <td><a href=\"%s\" class=\"%s\">%s</a></td>\n"
	DIRECTORY_LIST_PARENT_FORMAT      = "                <td><a href=\"%s\" class=\"parent-dir\">Parent Directory</a></td>\n"
	DIRECTORY_LIST_TD_FORMAT          = "                <td>%s</td>\n"
	DOT                               = "."
	EMPTY_PLACEHOLDER                 = "-"
	EPHEMERAL_PORT                    = 0
	EXECUTE                           = "EXECUTE"
	FILE_PERMISSION_PRIVATE           = 0600
	FILE_PERMISSION_PUBLIC            = 0644
	FILE_RENAME_FORMAT                = "%s_%d%s"
	FILE_SIZE_BYTE_FORMAT             = "%d B"
	FILE_SIZE_FORMAT                  = "%.1f %cB"
	FILE_SIZE_UNIT                    = 1024
	FILENAME_KEY                      = "filename"
	FORBIDDEN_LOG_FORMAT              = "%s forbidden"
	FORWARDED_PROTOCOL_HEADER_NAME    = "X-Forwarded-Proto"
	GET                               = "GET"
	HTTP_PORT_FORMAT                  = ":%d"
	HTTP_PROTOCOL                     = "http"
	HTTPS_PREFIX                      = "https://"
	HTTPS_PROTOCOL                    = "https"
	ICON_CLASS_DIR                    = "dir-icon"
	ICON_CLASS_FILE                   = "file-icon"
	INFO_FILE_SUFFIX                  = ".info"
	LOOPBACK_ADDRESS_IPV4             = "127.0.0.1"
	LOOPBACK_ADDRESS_IPV6             = "::1"
	LOOPBACK_ADDRESS_LOCALHOST        = "localhost"
	MAX_PORT                          = 65535
	MIME_TYPE_CSS                     = "text/css; charset=utf-8"
	MIME_TYPE_GIF                     = "image/gif"
	MIME_TYPE_HTML                    = "text/html; charset=utf-8"
	MIME_TYPE_ICON                    = "image/x-icon"
	MIME_TYPE_JAVASCRIPT              = "application/javascript; charset=utf-8"
	MIME_TYPE_JPEG                    = "image/jpeg"
	MIME_TYPE_JSON                    = "application/json; charset=utf-8"
	MIME_TYPE_PNG                     = "image/png"
	MIME_TYPE_SVG                     = "image/svg+xml"
	MIME_TYPE_TEXT                    = "text/plain; charset=utf-8"
	MIME_TYPE_WOFF                    = "font/woff"
	MIME_TYPE_WOFF2                   = "font/woff2"
	MODIFIED_TIME_FORMAT              = "2006-01-02 15:04:05"
	MIN_PORT                          = 1
	MODULE_NAME_SERVE                 = "serve"
	ORGANIZATION_NAME                 = "Go Solution SDK"
	PATH_PARAMETER_PREFIX             = "{"
	PATH_PARAMETER_SUFFIX             = "}"
	PFX_EXTENSION                     = ".pfx"
	POST                              = "POST"
	PRIVATE_KEY_BLOCK_TYPE            = "PRIVATE KEY"
	PUT                               = "PUT"
	QUERY_STRING_SEPARATOR            = "?"
	REQUEST_HASH_FORBIDDEN_RESPONSE   = "access denied"
	REQUEST_HASH_HEADER_NAME          = "X-Request-Hash"
	REQUEST_TIMESTAMP_HEADER_NAME     = "X-Request-Timestamp"
	REQUEST_UUID_HEADER_NAME          = "X-Request-UUID"
	RESET                             = "RESET"
	ROOT_ROUTE                        = "/"
	SALT_BYTES_LENGTH                 = 32
	SALT_RESPONSE_HASH_KEY            = "hash"
	SALT_RESPONSE_NO_CERTIFICATE      = "access denied"
	SAVE                              = "SAVE"
	SERVE_CONFIG_FILE_NAME            = "serve.json"
	SERVE_DATABASE_FILE_NAME          = "serve.db"
	SERVE_RECORD_TYPE_REQUEST         = "request"
	SERVE_REPLAY_RECORD_RETENTION     = 30 * 24 * time.Hour
	SERVE_SQL_CREATE_INDEX            = "CREATE UNIQUE INDEX IF NOT EXISTS idx_history_uuid ON history (uuid)"
	SERVE_SQL_CREATE_TABLE            = "CREATE TABLE IF NOT EXISTS history (id INTEGER PRIMARY KEY AUTOINCREMENT, uuid TEXT NOT NULL, timestamp TEXT NOT NULL, type TEXT NOT NULL, created_at INTEGER NOT NULL)"
	SERVE_SQL_DELETE_EXPIRED_REQUESTS = "DELETE FROM history WHERE created_at < ?"
	SERVE_SQL_INSERT_REQUEST          = "INSERT INTO history (uuid, timestamp, type, created_at) VALUES (?, ?, ?, ?)"
	SET                               = "SET"
	SHOW                              = "SHOW"
	SQLITE_DATABASE_DEFAULT_PASSWORD  = "sa"
	TCP                               = "tcp"
	TLS_HANDSHAKE_ERROR_WAIT_DURATION = 50 * time.Millisecond
	WAIT_CLIENT_TIMEOUT               = 5 * time.Second
	WAIT_RETRY_INTERVAL               = 10 * time.Millisecond
	WEBAPI_PATH_SALT                  = "/api/salt"
	WILDCARD_SUFFIX                   = "/*"
)

//goland:noinspection SpellCheckingInspection,GoSnakeCaseUsage
var (
	afterServeStaticHandlers      []AFTER_SERVE_STATIC_HANDLER
	beforeRequestHandlers         []BEFORE_REQUEST_HANDLER
	beforeRouteHandlers           []BEFORE_ROUTE_HANDLER
	beforeServeStaticHandlers     []BEFORE_SERVE_STATIC_HANDLER
	beforeTusdHandlers            []BEFORE_TUSD_HANDLER
	embedDirectories              map[string]embed.FS
	errorHandler                  func(error)
	htmlListDirectoryTemplate     string
	ignoredProtectionPatterns     []string
	isListDirectoryEnabled        bool
	jsonBodyHashMutex             sync.RWMutex
	jsonBodyHashValidation        bool
	loggedClientCertificates      map[string]bool
	loggedClientCertificatesMutex sync.Mutex
	protectedPatterns             []string
	requestHandlerRegistered      bool
	revokedCertificates           map[string]bool
	revokedCertificatesMutex      sync.RWMutex
	routes                        []ROUTE
	salt                          string
	saltMutex                     sync.RWMutex
	serveConfig                   *SERVE_CONFIG
	serveMutex                    sync.Mutex
	server                        *http.Server
	serverCancel                  __context.CancelFunc
	serverContext                 __context.Context
	serverListener                net.Listener
	sqliteDatabase                *sqlite.SQLITE
	sqliteDatabaseFilePath        string
	sqliteDatabaseInitOnce        sync.Once
	staticDirectories             map[string]string
	tlsClientCAPool               *x509.CertPool
	tlsPort                       int
	tlsServer                     *http.Server
	tlsServerCancel               __context.CancelFunc
	tlsServerContext              __context.Context
	tlsServerListener             net.Listener
	tusdMounts                    map[string]*TUSD_MOUNT
	tusdEnabled                   bool
	tusdMutex                     sync.Mutex
)

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func init() {
	serverContext, serverCancel = __context.WithCancel(__context.Background())
	tlsServerContext, tlsServerCancel = __context.WithCancel(__context.Background())
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
	htmlListDirectoryTemplate = "<!DOCTYPE html>\n" +
		"<html lang=\"en\">\n" +
		"<head>\n" +
		"    <meta charset=\"UTF-8\">\n" +
		"    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n" +
		"    <title>%s - /</title>\n" +
		"    <style>\n" +
		"        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; margin: 0; padding: 20px; background-color: #fff; }\n" +
		"        h1 { color: #003399; border-bottom: 1px solid #ccc; padding-bottom: 10px; }\n" +
		"        table { width: 100%; border-collapse: collapse; margin-top: 20px; }\n" +
		"        th, td { text-align: left; padding: 8px; border-bottom: 1px solid #ddd; }\n" +
		"        th { background-color: #f2f2f2; font-weight: bold; color: #333; }\n" +
		"        tr:hover { background-color: #f5f5f5; }\n" +
		"        a { color: #0066cc; text-decoration: none; }\n" +
		"        a:hover { text-decoration: underline; }\n" +
		"        .dir-icon::before { content: '馃搧'; margin-right: 8px; }\n" +
		"        .file-icon::before { content: '馃搫'; margin-right: 8px; }\n" +
		"        .parent-dir::before { content: '猬嗭笍'; margin-right: 8px; }\n" +
		"    </style>\n" +
		"</head>\n" +
		"<body>\n" +
		"    <h1>%s - /</h1>\n" +
		"    <table>\n" +
		"        <thead>\n" +
		"            <tr>\n" +
		"                <th>Name</th>\n" +
		"                <th>Size</th>\n" +
		"                <th>Last Modified</th>\n" +
		"            </tr>\n" +
		"        </thead>\n" +
		"        </tbody>\n"
	cleanupExpiredReplayRecords()
	protectedPatterns = []string{DEFAULT_PROTECTED_PATTERN}
	ignoredProtectionPatterns = []string{WEBAPI_PATH_SALT}
	tlsClientCAPool = loadDefaultCertificatePool()
	tusdMounts = make(map[string]*TUSD_MOUNT)
	revokedCertificates = make(map[string]bool)
	loadServeConfig()
	salt = generateRandomSalt()
	routes = append(routes, ROUTE{
		Method:  POST,
		Pattern: WEBAPI_PATH_SALT,
		Handler: handleGetSalt,
	})
	registerMimeTypes()
}

func registerMimeTypes() {
	err := error(nil)
	mappings := []struct {
		extension string
		mimeType  string
	}{
		{".css", MIME_TYPE_CSS},
		{".gif", MIME_TYPE_GIF},
		{".htm", MIME_TYPE_HTML},
		{".html", MIME_TYPE_HTML},
		{".ico", MIME_TYPE_ICON},
		{".jpeg", MIME_TYPE_JPEG},
		{".jpg", MIME_TYPE_JPEG},
		{".js", MIME_TYPE_JAVASCRIPT},
		{".json", MIME_TYPE_JSON},
		{".mjs", MIME_TYPE_JAVASCRIPT},
		{".png", MIME_TYPE_PNG},
		{".svg", MIME_TYPE_SVG},
		{".txt", MIME_TYPE_TEXT},
		{".woff", MIME_TYPE_WOFF},
		{".woff2", MIME_TYPE_WOFF2},
	}
	for _, mapping := range mappings {
		if err = mime.AddExtensionType(mapping.extension, mapping.mimeType); err != nil {
			__warning(fmt.Sprintf("register mime type %s failed: %v", mapping.extension, err))
		}
	}
}

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_SERVE, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_SERVE, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_SERVE, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_SERVE, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedExportedFunction
func AddEmbedDirectoryMapping(urlPath string, embedFileSystem embed.FS) {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	embedDirectories[urlPath] = embedFileSystem
}

//goland:noinspection GoUnusedExportedFunction
func AddIgnoredProtectionPattern(pattern string) {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	duplicateFound := false
	for _, existing := range ignoredProtectionPatterns {
		if existing == pattern {
			duplicateFound = true
			break
		}
	}
	if !duplicateFound {
		ignoredProtectionPatterns = append(ignoredProtectionPatterns, pattern)
	}
}

//goland:noinspection GoUnusedExportedFunction
func AddIgnoredProtectionPatterns(patterns []string) {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	for _, pattern := range patterns {
		duplicateFound := false
		for _, existing := range ignoredProtectionPatterns {
			if existing == pattern {
				duplicateFound = true
				break
			}
		}
		if !duplicateFound {
			ignoredProtectionPatterns = append(ignoredProtectionPatterns, pattern)
		}
	}
}

//goland:noinspection GoUnusedExportedFunction
func AddProtectedPattern(pattern string) {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	duplicateFound := false
	for _, existing := range protectedPatterns {
		if existing == pattern {
			duplicateFound = true
			break
		}
	}
	if !duplicateFound {
		protectedPatterns = append(protectedPatterns, pattern)
	}
}

//goland:noinspection GoUnusedExportedFunction
func AddStaticDirectoryMapping(urlPath string, directoryPath string) {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	staticDirectories[urlPath] = directoryPath
}

//goland:noinspection GoUnusedExportedFunction
func ClearRevokedCertificates() {
	revokedCertificatesMutex.Lock()
	defer revokedCertificatesMutex.Unlock()
	revokedCertificates = make(map[string]bool)
	saveServeConfig()
}

//goland:noinspection GoUnusedExportedFunction
func DisableJSONBodyValidation() {
	jsonBodyHashMutex.Lock()
	defer jsonBodyHashMutex.Unlock()
	jsonBodyHashValidation = false
}

//goland:noinspection GoUnusedExportedFunction
func EnableDirectoryListing(enabled bool) {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	isListDirectoryEnabled = enabled
}

//goland:noinspection GoUnusedExportedFunction
func EnableJSONBodyValidation() {
	jsonBodyHashMutex.Lock()
	defer jsonBodyHashMutex.Unlock()
	jsonBodyHashValidation = true
}

//goland:noinspection GoUnusedExportedFunction
func EnableRedirectToTLS() error {
	result := error(nil)
	serveMutex.Lock()
	if tlsServerListener == nil {
		result = fmt.Errorf("TLS server is not running")
	} else if serverListener != nil {
		result = fmt.Errorf("HTTP server is already running")
	}
	serveMutex.Unlock()
	if result == nil {
		httpAddress := fmt.Sprintf(HTTP_PORT_FORMAT, DEFAULT_HTTP_PORT)
		serveMutex.Lock()
		server.Addr = httpAddress
		serveMutex.Unlock()
		http.HandleFunc(ROOT_ROUTE, func(responseWriter http.ResponseWriter, request *http.Request) {
			target := HTTPS_PREFIX + strings.Replace(request.Host, fmt.Sprintf(HTTP_PORT_FORMAT, DEFAULT_HTTP_PORT), fmt.Sprintf(HTTP_PORT_FORMAT, tlsPort), 1) + request.URL.Path
			if request.URL.RawQuery != "" {
				target += QUERY_STRING_SEPARATOR + request.URL.RawQuery
			}
			http.Redirect(responseWriter, request, target, http.StatusMovedPermanently)
		})
		go func() {
			if listener, err := net.Listen(TCP, httpAddress); err == nil {
				serveMutex.Lock()
				serverListener = listener
				serveMutex.Unlock()
				if err = http.Serve(serverListener, nil); !errors.Is(err, http.ErrServerClosed) {
					if errorHandler != nil {
						errorHandler(err)
					}
				}
			}
			serveMutex.Lock()
			serverListener = nil
			serveMutex.Unlock()
		}()
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GenerateSelfSignedCertificate() (*tls.Certificate, error) {
	var result *tls.Certificate
	err := error(nil)
	var certificatePEM []byte
	var privateKeyPEM []byte
	if certificatePEM, privateKeyPEM, err = generateSelfSignedCertificatePEM(); err == nil {
		var tlsCertificate tls.Certificate
		if tlsCertificate, err = tls.X509KeyPair(certificatePEM, privateKeyPEM); err == nil {
			result = &tlsCertificate
		}
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func GenerateSelfSignedCertificateFiles(certificatePath string, keyPath string) error {
	result := error(nil)
	var certificatePEM []byte
	var privateKeyPEM []byte
	if certificatePEM, privateKeyPEM, result = generateSelfSignedCertificatePEM(); result == nil {
		if result = os.MkdirAll(filepath.Dir(certificatePath), DEFAULT_DIRECTORY_PERMISSION); result == nil {
			if result = os.MkdirAll(filepath.Dir(keyPath), DEFAULT_DIRECTORY_PERMISSION); result == nil {
				if result = os.WriteFile(certificatePath, certificatePEM, FILE_PERMISSION_PUBLIC); result == nil {
					result = os.WriteFile(keyPath, privateKeyPEM, FILE_PERMISSION_PRIVATE)
				}
			}
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetEnableJsonBodyHashValidation() bool {
	result := false
	jsonBodyHashMutex.RLock()
	if jsonBodyHashValidation {
		result = true
	}
	jsonBodyHashMutex.RUnlock()
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetIdleTimeout() time.Duration {
	result := time.Duration(0)
	serveMutex.Lock()
	result = server.IdleTimeout
	serveMutex.Unlock()
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetIgnoredProtectionPatterns() []string {
	result := make([]string, 0)
	serveMutex.Lock()
	result = append(result, ignoredProtectionPatterns...)
	serveMutex.Unlock()
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetPort() int {
	result := 0
	var err error
	serveMutex.Lock()
	var portString string
	if _, portString, err = net.SplitHostPort(server.Addr); err == nil {
		var port int
		if port, err = strconv.Atoi(portString); err == nil {
			result = port
		}
	}
	serveMutex.Unlock()
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetProtectedPatterns() []string {
	result := make([]string, 0)
	serveMutex.Lock()
	result = append(result, protectedPatterns...)
	serveMutex.Unlock()
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetReadTimeout() time.Duration {
	result := time.Duration(0)
	serveMutex.Lock()
	result = server.ReadTimeout
	serveMutex.Unlock()
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetRevokedCertificates() []string {
	revokedCertificatesMutex.RLock()
	defer revokedCertificatesMutex.RUnlock()
	result := make([]string, 0, len(revokedCertificates))
	for certificate := range revokedCertificates {
		result = append(result, certificate)
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetSalt() string {
	result := ""
	saltMutex.RLock()
	result = salt
	saltMutex.RUnlock()
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetStaticDirectory(urlPath string) string {
	result := ""
	serveMutex.Lock()
	result = staticDirectories[urlPath]
	serveMutex.Unlock()
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetTlsIdleTimeout() time.Duration {
	result := time.Duration(0)
	serveMutex.Lock()
	result = tlsServer.IdleTimeout
	serveMutex.Unlock()
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetTlsPort() int {
	result := 0
	var err error
	serveMutex.Lock()
	var portString string
	if _, portString, err = net.SplitHostPort(tlsServer.Addr); err == nil {
		var port int
		if port, err = strconv.Atoi(portString); err == nil {
			result = port
		}
	}
	serveMutex.Unlock()
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetTlsReadTimeout() time.Duration {
	result := time.Duration(0)
	serveMutex.Lock()
	result = tlsServer.ReadTimeout
	serveMutex.Unlock()
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetTlsWriteTimeout() time.Duration {
	result := time.Duration(0)
	serveMutex.Lock()
	result = tlsServer.WriteTimeout
	serveMutex.Unlock()
	return result
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func GetTusdMount(uri string) *TUSD_MOUNT {
	var result *TUSD_MOUNT
	tusdMutex.Lock()
	result = tusdMounts[uri]
	tusdMutex.Unlock()
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetTusdEnabled() bool {
	result := false
	serveMutex.Lock()
	result = tusdEnabled
	serveMutex.Unlock()
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetWriteTimeout() time.Duration {
	result := time.Duration(0)
	serveMutex.Lock()
	result = server.WriteTimeout
	serveMutex.Unlock()
	return result
}

//goland:noinspection GoUnusedExportedFunction
func IsHttpRequest(request *http.Request) bool {
	result := false
	if request != nil {
		forwardedProtocol := request.Header.Get(FORWARDED_PROTOCOL_HEADER_NAME)
		if forwardedProtocol != "" {
			if strings.EqualFold(forwardedProtocol, HTTP_PROTOCOL) {
				result = true
			}
		} else if request.URL != nil && request.URL.Scheme != "" {
			if strings.EqualFold(request.URL.Scheme, HTTP_PROTOCOL) {
				result = true
			}
		} else {
			if request.TLS == nil {
				result = true
			}
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func IsHttpsRequest(request *http.Request) bool {
	result := false
	if request != nil {
		forwardedProtocol := request.Header.Get(FORWARDED_PROTOCOL_HEADER_NAME)
		if forwardedProtocol != "" {
			if strings.EqualFold(forwardedProtocol, HTTPS_PROTOCOL) {
				result = true
			}
		} else if request.URL != nil && request.URL.Scheme != "" {
			if strings.EqualFold(request.URL.Scheme, HTTPS_PROTOCOL) {
				result = true
			}
		} else {
			if request.TLS != nil {
				result = true
			}
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func IsListening() bool {
	result := false
	serveMutex.Lock()
	if serverListener != nil {
		result = true
	}
	serveMutex.Unlock()
	return result
}

//goland:noinspection GoUnusedExportedFunction
func IsLoopbackAddress(request *http.Request) bool {
	result := false
	if request != nil {
		host, _, err := net.SplitHostPort(request.RemoteAddr)
		if err != nil {
			host = request.RemoteAddr
		}
		if host == LOOPBACK_ADDRESS_IPV4 || host == LOOPBACK_ADDRESS_LOCALHOST || host == LOOPBACK_ADDRESS_IPV6 {
			result = true
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func IsTlsListening() bool {
	result := false
	serveMutex.Lock()
	if tlsServerListener != nil {
		result = true
	}
	serveMutex.Unlock()
	return result
}

//goland:noinspection GoUnusedExportedFunction
func resolveAddressPort(address string) string {
	err := error(nil)
	result := address
	hostname := ""
	portString := ""
	port := 0
	randomPort := 0
	if hostname, portString, err = net.SplitHostPort(address); err == nil {
		if port, err = strconv.Atoi(portString); err == nil && port == EPHEMERAL_PORT {
			randomPort = network.GetRandomPort()
			if randomPort > 0 {
				result = net.JoinHostPort(hostname, strconv.Itoa(randomPort))
			}
		}
	}
	return result
}

func ListenAsync() error {
	result := error(nil)
	__debug("Starting ListenAsync")
	var serverAddr string
	serveMutex.Lock()
	if serverListener != nil {
		__debug("HTTP server is already running")
		result = fmt.Errorf("HTTP server is already running")
		serveMutex.Unlock()
	} else {
		server.Addr = resolveAddressPort(server.Addr)
		serverAddr = server.Addr
		serveMutex.Unlock()
		__debug("Initializing TUSd mounts")
		if tusdEnabled {
			if result = initializeTusdMounts(); result == nil {
				__debug("TUSd mounts initialized successfully")
			} else {
				__debug("Failed to initialize TUSd mounts: " + result.Error())
			}
		} else {
			__debug("TUSd disabled, skipping initialization")
		}
		if result == nil {
			registerRequestHandler()
			go func() {
				__debug(fmt.Sprintf("Attempting to listen on TCP network: %s", serverAddr))
				var listener net.Listener
				var err error
				if listener, err = net.Listen(TCP, serverAddr); err == nil {
					__debug("TCP listener created successfully")
					serveMutex.Lock()
					serverListener = listener
					serveMutex.Unlock()
					__debug("Starting HTTP server")
					if err = http.Serve(serverListener, nil); !errors.Is(err, http.ErrServerClosed) {
						__debug("HTTP server error occurred: " + err.Error())
						if errorHandler != nil {
							__debug("Calling error handler")
							errorHandler(err)
						}
					} else {
						__debug("HTTP server closed normally")
					}
				} else {
					__debug("Failed to create TCP listener: " + err.Error())
				}
				serveMutex.Lock()
				serverListener = nil
				serveMutex.Unlock()
				__debug("Server listener cleared")
			}()
		}
	}
	if result == nil {
		__debug("Waiting for server to be ready")
		waitServe(serverAddr)
		__debug("Server is ready")
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func ListenTlsAsync(address string) error {
	result := error(nil)
	resolvedAddress := resolveAddressPort(address)
	serveMutex.Lock()
	if tlsServerListener == nil {
		serveMutex.Unlock()
		if tusdEnabled {
			if result = initializeTusdMounts(); result == nil {
				__debug("TUSd mounts initialized successfully")
			} else {
				__debug("Failed to initialize TUSd mounts: " + result.Error())
			}
		} else {
			__debug("TUSd disabled, skipping initialization")
		}
		if result == nil {
			registerRequestHandler()
			go func() {
				err := error(nil)
				var tlsConfig *tls.Config
				if tlsConfig, err = createTLSConfig(); err == nil {
					var listener net.Listener
					if listener, err = tls.Listen(TCP, resolvedAddress, tlsConfig); err == nil {
						serveMutex.Lock()
						tlsServerListener = listener
						serveMutex.Unlock()
						err = http.Serve(tlsServerListener, nil)
					}
				}
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					serveMutex.Lock()
					handler := errorHandler
					serveMutex.Unlock()
					if handler != nil {
						handler(err)
					}
				}
				serveMutex.Lock()
				tlsServerListener = nil
				serveMutex.Unlock()
			}()
		}
	} else {
		result = fmt.Errorf("TLS server is already running")
		serveMutex.Unlock()
	}
	if result == nil {
		__debug("Waiting for tls server to be ready")
		waitServeTls(resolvedAddress)
		__debug("Tls server is ready, setting up request handler")
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func On(method string, pattern string, handler REQUEST_HANDLER) error {
	result := error(nil)
	tusdMutex.Lock()
	for basePath := range tusdMounts {
		if isPathInUseForTusd(pattern, basePath) {
			result = fmt.Errorf("pattern %s conflicts with TUSD mount point %s", pattern, basePath)
			break
		}
	}
	tusdMutex.Unlock()
	if result == nil {
		serveMutex.Lock()
		routes = append(routes, ROUTE{
			Method:  method,
			Pattern: pattern,
			Handler: handler,
		})
		serveMutex.Unlock()
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func OnAfterStaticServe(handler AFTER_SERVE_STATIC_HANDLER) {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	afterServeStaticHandlers = append(afterServeStaticHandlers, handler)
}

//goland:noinspection GoUnusedExportedFunction
func OnBeforeRequest(callback BEFORE_REQUEST_HANDLER) {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	beforeRequestHandlers = append(beforeRequestHandlers, callback)
}

//goland:noinspection GoUnusedExportedFunction
func OnBeforeRoute(callback BEFORE_ROUTE_HANDLER) {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	beforeRouteHandlers = append(beforeRouteHandlers, callback)
}

//goland:noinspection GoUnusedExportedFunction
func OnBeforeStaticServe(handler BEFORE_SERVE_STATIC_HANDLER) {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	beforeServeStaticHandlers = append(beforeServeStaticHandlers, handler)
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func OnBeforeTusd(handler BEFORE_TUSD_HANDLER) {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	beforeTusdHandlers = append(beforeTusdHandlers, handler)
}

//goland:noinspection GoUnusedExportedFunction
func OnError(handler func(error)) {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	errorHandler = handler
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func OnTusd(uri string, directoryPath string, handler TUSD_UPLOADED_HANDLER) error {
	result := error(nil)
	if !strings.HasPrefix(uri, ROOT_ROUTE) {
		uri = ROOT_ROUTE + uri
	}
	if !strings.HasSuffix(uri, ROOT_ROUTE) {
		uri = uri + ROOT_ROUTE
	}
	serveMutex.Lock()
	for _, route := range routes {
		if isPathInUseForTusd(route.Pattern, uri) {
			result = fmt.Errorf("TUSD mount %s conflicts with existing route %s", uri, route.Pattern)
			break
		}
	}
	serveMutex.Unlock()
	if result == nil {
		tusdMutex.Lock()
		tusdMounts[uri] = &TUSD_MOUNT{
			Uri:             uri,
			DirectoryPath:   directoryPath,
			MaxSize:         0,
			TargetDirectory: "",
			UploadedHandler: handler,
		}
		tusdMutex.Unlock()
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func RemoveIgnoredProtectionPattern(pattern string) {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	updated := make([]string, 0, len(ignoredProtectionPatterns))
	for _, existing := range ignoredProtectionPatterns {
		if existing != pattern {
			updated = append(updated, existing)
		}
	}
	ignoredProtectionPatterns = updated
}

//goland:noinspection GoUnusedExportedFunction
func RemoveIgnoredProtectionPatterns() {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	ignoredProtectionPatterns = make([]string, 0)
}

//goland:noinspection GoUnusedExportedFunction
func RemoveProtectedPattern(pattern string) {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	updated := make([]string, 0, len(protectedPatterns))
	for _, existing := range protectedPatterns {
		if existing != pattern {
			updated = append(updated, existing)
		}
	}
	protectedPatterns = updated
}

//goland:noinspection GoUnusedExportedFunction
func RemoveProtectedPatterns() {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	protectedPatterns = make([]string, 0)
}

//goland:noinspection GoUnusedExportedFunction
func RemoveStaticDirectoryMapping(urlPath string) {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	delete(staticDirectories, urlPath)
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func RemoveTusdMount(uri string) {
	tusdMutex.Lock()
	defer tusdMutex.Unlock()
	delete(tusdMounts, uri)
}

//goland:noinspection GoUnusedExportedFunction
func SetAddress(address string) {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	server.Addr = address
}

//goland:noinspection GoUnusedExportedFunction
func SetContext(context __context.Context) {
	serveMutex.Lock()
	if serverCancel != nil {
		serverCancel()
	}
	if tlsServerCancel != nil {
		tlsServerCancel()
	}
	newServerContext, serverCancelFunction := __context.WithCancel(context)
	newTlsServerContext, tlsServerCancelFunction := __context.WithCancel(context)
	serverContext = newServerContext
	tlsServerContext = newTlsServerContext
	serverCancel = serverCancelFunction
	tlsServerCancel = tlsServerCancelFunction
	serveMutex.Unlock()
}

//goland:noinspection GoUnusedExportedFunction
func SetIdleTimeout(timeout time.Duration) {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	server.IdleTimeout = timeout
}

//goland:noinspection GoUnusedExportedFunction
func SetIdleTimeouts(timeout time.Duration) {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	server.IdleTimeout = timeout
	tlsServer.IdleTimeout = timeout
}

//goland:noinspection GoUnusedExportedFunction
func SetProtectedPatterns(patterns []string) {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	protectedPatterns = append([]string(nil), patterns...)
}

//goland:noinspection GoUnusedExportedFunction
func SetReadTimeout(timeout time.Duration) {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	server.ReadTimeout = timeout
}

//goland:noinspection GoUnusedExportedFunction
func SetReadTimeouts(timeout time.Duration) {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	server.ReadTimeout = timeout
	tlsServer.ReadTimeout = timeout
}

//goland:noinspection GoUnusedExportedFunction
func SetRevokedCertificates(certificates []string) {
	revokedCertificatesMutex.Lock()
	defer revokedCertificatesMutex.Unlock()
	revokedCertificates = make(map[string]bool)
	for _, certificate := range certificates {
		revokedCertificates[certificate] = true
	}
	saveServeConfig()
}

//goland:noinspection GoUnusedExportedFunction
func SetTimeouts(readTimeout, writeTimeout, idleTimeout time.Duration) {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	server.ReadTimeout = readTimeout
	server.WriteTimeout = writeTimeout
	server.IdleTimeout = idleTimeout
	tlsServer.ReadTimeout = readTimeout
	tlsServer.WriteTimeout = writeTimeout
	tlsServer.IdleTimeout = idleTimeout
}

//goland:noinspection GoUnusedExportedFunction
func SetTlsAddress(address string) {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	tlsServer.Addr = address
}

//goland:noinspection GoUnusedExportedFunction
func SetTlsIdleTimeout(timeout time.Duration) {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	tlsServer.IdleTimeout = timeout
}

//goland:noinspection GoUnusedExportedFunction
func SetTlsReadTimeout(timeout time.Duration) {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	tlsServer.ReadTimeout = timeout
}

//goland:noinspection GoUnusedExportedFunction
func SetTlsWriteTimeout(timeout time.Duration) {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	tlsServer.WriteTimeout = timeout
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func SetTusdMount(basePath string, maxSize int64, targetDirectory string) error {
	result := error(nil)
	tusdMutex.Lock()
	if mount, ok := tusdMounts[basePath]; ok {
		mount.MaxSize = maxSize
		mount.TargetDirectory = targetDirectory
	} else {
		result = fmt.Errorf("TUSD mount %s does not exist", basePath)
	}
	tusdMutex.Unlock()
	return result
}

//goland:noinspection GoUnusedExportedFunction
func SetTusdEnabled(enabled bool) {
	serveMutex.Lock()
	tusdEnabled = enabled
	serveMutex.Unlock()
}

//goland:noinspection GoUnusedExportedFunction
func SetWriteTimeout(timeout time.Duration) {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	server.WriteTimeout = timeout
}

//goland:noinspection GoUnusedExportedFunction
func SetWriteTimeouts(timeout time.Duration) {
	serveMutex.Lock()
	defer serveMutex.Unlock()
	server.WriteTimeout = timeout
	tlsServer.WriteTimeout = timeout
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func Shutdown() {
	shutdown()
	shutdownTls()
}

func createConfigurationSnapshot() STATE_SNAPSHOT {
	result := STATE_SNAPSHOT{}
	serveMutex.Lock()
	result.Routes = make([]ROUTE, len(routes))
	result.BeforeRequestHandlers = make([]BEFORE_REQUEST_HANDLER, len(beforeRequestHandlers))
	result.BeforeRouteHandlers = make([]BEFORE_ROUTE_HANDLER, len(beforeRouteHandlers))
	result.BeforeServeStaticHandlers = make([]BEFORE_SERVE_STATIC_HANDLER, len(beforeServeStaticHandlers))
	result.BeforeTusdHandlers = make([]BEFORE_TUSD_HANDLER, len(beforeTusdHandlers))
	result.AfterServeStaticHandlers = make([]AFTER_SERVE_STATIC_HANDLER, len(afterServeStaticHandlers))
	result.ProtectedPatterns = append([]string(nil), protectedPatterns...)
	result.IgnoredProtectionPatterns = append([]string(nil), ignoredProtectionPatterns...)
	result.StaticDirectories = make(map[string]string)
	for key, value := range staticDirectories {
		result.StaticDirectories[key] = value
	}
	result.EmbedDirectories = make(map[string]embed.FS)
	for key, value := range embedDirectories {
		result.EmbedDirectories[key] = value
	}
	copy(result.Routes, routes)
	copy(result.BeforeRequestHandlers, beforeRequestHandlers)
	copy(result.BeforeRouteHandlers, beforeRouteHandlers)
	copy(result.BeforeServeStaticHandlers, beforeServeStaticHandlers)
	copy(result.BeforeTusdHandlers, beforeTusdHandlers)
	copy(result.AfterServeStaticHandlers, afterServeStaticHandlers)
	serveMutex.Unlock()
	tusdMutex.Lock()
	result.TusdMounts = make(map[string]*TUSD_MOUNT)
	for key, value := range tusdMounts {
		result.TusdMounts[key] = value
	}
	tusdMutex.Unlock()
	return result
}

//goland:noinspection SpellCheckingInspection
func createDefaultServeConfig() {
	defaultRevokedCertificates := []string{
		"0000000000000000000000000000000000000000000000000000000000000000",
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
	}
	serveConfig = &SERVE_CONFIG{
		RevokedCertificates: defaultRevokedCertificates,
	}
	revokedCertificatesMutex.Lock()
	revokedCertificates = make(map[string]bool)
	for _, certificate := range defaultRevokedCertificates {
		revokedCertificates[certificate] = true
	}
	revokedCertificatesMutex.Unlock()
	saveServeConfig()
	__debug("Created default serve config file")
}

func createTLSConfig() (*tls.Config, error) {
	var result *tls.Config
	err := error(nil)
	certificate := tls.Certificate{}
	certificatePath := findCertificateFile(ca.SERVER_CERTIFICATE_FILE_NAME)
	privateKeyPath := findCertificateFile(ca.SERVER_PRIVATE_KEY_FILE_NAME)
	if certificatePath != "" && privateKeyPath != "" {
		certificate, err = tls.LoadX509KeyPair(certificatePath, privateKeyPath)
		if err == nil {
			__debug("[TLS] Loaded certificate from files: " + certificatePath)
		} else {
			__debug("[TLS] Failed to load certificate files, falling back to self-signed certificate: " + err.Error())
		}
	}
	if err != nil {
		var generatedCertificate *tls.Certificate
		if generatedCertificate, err = GenerateSelfSignedCertificate(); err == nil && generatedCertificate != nil {
			certificate = *generatedCertificate
			__debug("[TLS] Generated self-signed certificate")
		} else if generatedCertificate == nil {
			err = fmt.Errorf("failed to generate self-signed certificate: certificate is nil")
		}
	}
	if err == nil {
		if tlsClientCAPool == nil {
			tlsClientCAPool = loadDefaultCertificatePool()
		}
		result = &tls.Config{
			Certificates: []tls.Certificate{certificate},
			MinVersion:   tls.VersionTLS12,
		}
		if tlsClientCAPool != nil {
			result.ClientAuth = tls.RequireAndVerifyClientCert
			result.ClientCAs = tlsClientCAPool
		} else {
			result.ClientAuth = tls.RequestClientCert
		}
	}
	return result, err
}

//goland:noinspection DuplicatedCode
func findCertificateFile(fileName string) string {
	result := ""
	__debug(fmt.Sprintf("[Certificate] Looking for certificate file: %s", fileName))
	var executablePath string
	var err error
	if executablePath, err = os.Executable(); err == nil {
		executableDirectory := filepath.Dir(executablePath)
		candidate := filepath.Join(executableDirectory, fileName)
		__debug(fmt.Sprintf("[Certificate] Checking executable directory: %s", candidate))
		if _, err = os.Stat(candidate); err == nil {
			result = candidate
			__debug(fmt.Sprintf("[Certificate] Found in executable directory: %s", candidate))
		} else {
			__debug(fmt.Sprintf("[Certificate] Not found in executable directory: %s, error: %v", candidate, err))
		}
	} else {
		__debug(fmt.Sprintf("[Certificate] Failed to get executable path: %v", err))
	}
	if result == "" {
		var workingDirectory string
		if workingDirectory, err = os.Getwd(); err == nil {
			candidate := filepath.Join(workingDirectory, fileName)
			__debug(fmt.Sprintf("[Certificate] Checking working directory: %s", candidate))
			if _, err = os.Stat(candidate); err == nil {
				result = candidate
				__debug(fmt.Sprintf("[Certificate] Found in working directory: %s", candidate))
			} else {
				__debug(fmt.Sprintf("[Certificate] Not found in working directory: %s, error: %v", candidate, err))
			}
		} else {
			__debug(fmt.Sprintf("[Certificate] Failed to get working directory: %v", err))
		}
	}
	if result == "" {
		__debug(fmt.Sprintf("[Certificate] File not found: %s", fileName))
	}
	return result
}

func findServeConfigFile() string {
	result := ""
	var executablePath string
	var err error
	if executablePath, err = os.Executable(); err == nil {
		executableDirectory := filepath.Dir(executablePath)
		candidate := filepath.Join(executableDirectory, SERVE_CONFIG_FILE_NAME)
		if _, err = os.Stat(candidate); err == nil {
			result = candidate
		}
	}
	if result == "" {
		var workingDirectory string
		if workingDirectory, err = os.Getwd(); err == nil {
			candidate := filepath.Join(workingDirectory, SERVE_CONFIG_FILE_NAME)
			if _, err = os.Stat(candidate); err == nil {
				result = candidate
			}
		}
	}
	return result
}

func formatFileSize(size int64) string {
	result := ""
	if size < FILE_SIZE_UNIT {
		result = fmt.Sprintf(FILE_SIZE_BYTE_FORMAT, size)
	} else {
		divisor, exponent := int64(FILE_SIZE_UNIT), 0
		for n := size / FILE_SIZE_UNIT; n >= FILE_SIZE_UNIT; n /= FILE_SIZE_UNIT {
			divisor *= FILE_SIZE_UNIT
			exponent++
		}
		result = fmt.Sprintf(FILE_SIZE_FORMAT, float64(size)/float64(divisor), BYTE_UNITS[exponent])
	}
	return result
}

func generateRandomSalt() string {
	result := ""
	bytesValue := make([]byte, SALT_BYTES_LENGTH)
	var err error
	if _, err = rand.Read(bytesValue); err == nil {
		result = hex.EncodeToString(bytesValue)
	} else {
		result = strconv.FormatInt(time.Now().UnixNano(), 16)
	}
	return result
}

func generateSelfSignedCertificatePEM() ([]byte, []byte, error) {
	resultCertificate := make([]byte, 0)
	resultPrivateKey := make([]byte, 0)
	err := error(nil)
	var privateKey *rsa.PrivateKey
	if privateKey, err = rsa.GenerateKey(rand.Reader, DEFAULT_RSA_KEY_BITS); err == nil {
		validFromTime := time.Now()
		validToTime := validFromTime.Add(DEFAULT_VALID_DAYS * 24 * time.Hour)
		serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), DEFAULT_SERIAL_NUMBER_BITS)
		var serialNumber *big.Int
		if serialNumber, err = rand.Int(rand.Reader, serialNumberLimit); err == nil {
			certificateTemplate := x509.Certificate{
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
			var certificateDER []byte
			if certificateDER, err = x509.CreateCertificate(rand.Reader, &certificateTemplate, &certificateTemplate, &privateKey.PublicKey, privateKey); err == nil {
				var privateKeyDER []byte
				if privateKeyDER, err = x509.MarshalPKCS8PrivateKey(privateKey); err == nil {
					certificatePEM := new(bytes.Buffer)
					_ = pem.Encode(certificatePEM, &pem.Block{Type: CERTIFICATE_BLOCK_TYPE, Bytes: certificateDER})
					privateKeyPEM := new(bytes.Buffer)
					_ = pem.Encode(privateKeyPEM, &pem.Block{Type: PRIVATE_KEY_BLOCK_TYPE, Bytes: privateKeyDER})
					resultCertificate = certificatePEM.Bytes()
					resultPrivateKey = privateKeyPEM.Bytes()
				}
			}
		}
	}
	return resultCertificate, resultPrivateKey, err
}

func cleanupExpiredReplayRecords() {
	databaseFilePath := getSqliteDatabaseFilePath()
	if _, err := os.Stat(databaseFilePath); err == nil {
		database := sqlite.New()
		database.PragmaKey = SQLITE_DATABASE_DEFAULT_PASSWORD
		if err = database.Open(databaseFilePath); err == nil {
			expiresAt := time.Now().Add(-SERVE_REPLAY_RECORD_RETENTION).Unix()
			if err = database.Exec(SERVE_SQL_DELETE_EXPIRED_REQUESTS, expiresAt); err != nil {
				__debug(fmt.Sprintf("Failed to delete expired replay records: %v", err))
			}
		} else {
			__debug(fmt.Sprintf("Failed to open serve database for replay record cleanup: %v", err))
		}
		database.Close()
	} else if !os.IsNotExist(err) {
		__debug(fmt.Sprintf("Failed to check serve database file %s: %v", databaseFilePath, err))
	}
}

func getCertificateFingerprint(certificate *x509.Certificate) string {
	result := ""
	if certificate != nil {
		sum := sha256.Sum256(certificate.Raw)
		result = hex.EncodeToString(sum[:])
	}
	return result
}

func getDirectoryListHTML(path string, urlPath string, embedFileSystemPointer *embed.FS) string {
	result := ""
	var entries []fs.DirEntry
	var err error
	if embedFileSystemPointer != nil {
		entries, err = fs.ReadDir(*embedFileSystemPointer, path)
	} else {
		entries, err = os.ReadDir(path)
	}
	if err == nil {
		html := strings.Builder{}
		html.WriteString(fmt.Sprintf(htmlListDirectoryTemplate, urlPath, urlPath))
		if urlPath != ROOT_ROUTE {
			parentPath := strings.TrimSuffix(urlPath, ROOT_ROUTE)
			lastSlash := strings.LastIndex(parentPath, ROOT_ROUTE)
			if lastSlash >= 0 {
				parentPath = parentPath[:lastSlash]
				if parentPath == "" {
					parentPath = ROOT_ROUTE
				}
			}
			html.WriteString("            <tr>\n")
			html.WriteString(fmt.Sprintf(DIRECTORY_LIST_PARENT_FORMAT, parentPath))
			html.WriteString("                <td>-</td>\n")
			html.WriteString("                <td>-</td>\n")
			html.WriteString("            </tr>\n")
		}
		for _, entry := range entries {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			name := info.Name()
			size := ""
			if !info.IsDir() {
				size = formatFileSize(info.Size())
			}
			lastModified := info.ModTime().Format(MODIFIED_TIME_FORMAT)
			linkPath := strings.TrimSuffix(urlPath, ROOT_ROUTE) + ROOT_ROUTE + name
			className := ICON_CLASS_FILE
			if info.IsDir() {
				className = ICON_CLASS_DIR
			}
			html.WriteString("            <tr>\n")
			html.WriteString(fmt.Sprintf(DIRECTORY_LIST_ENTRY_FORMAT, linkPath, className, name))
			html.WriteString(fmt.Sprintf(DIRECTORY_LIST_TD_FORMAT, size))
			html.WriteString(fmt.Sprintf(DIRECTORY_LIST_TD_FORMAT, lastModified))
			html.WriteString("            </tr>\n")
		}
		html.WriteString("        </tbody>\n")
		html.WriteString("    </table>\n")
		html.WriteString("</body>\n")
		html.WriteString("</html>\n")
		result = html.String()
	}
	return result
}

func getEmbedFsRootDirectory(embedFileSystem embed.FS) string {
	result := ""
	embedFilesPath := make([]string, 0)
	err := error(nil)
	if err = fs.WalkDir(embedFileSystem, DOT, func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			embedFilesPath = append(embedFilesPath, path)
		}
		return nil
	}); err == nil && len(embedFilesPath) > 0 {
		if len(embedFilesPath) == 1 {
			parts := strings.Split(embedFilesPath[0], ROOT_ROUTE)
			if len(parts) > 1 {
				result = parts[0]
			}
		} else {
			filePathParts := strings.Split(embedFilesPath[0], ROOT_ROUTE)
			for index := 1; index < len(embedFilesPath); index++ {
				parts := strings.Split(embedFilesPath[index], ROOT_ROUTE)
				minLength := len(filePathParts)
				if len(parts) < minLength {
					minLength = len(parts)
				}
				matchLength := 0
				for j := 0; j < minLength; j++ {
					if filePathParts[j] == parts[j] {
						matchLength++
					} else {
						break
					}
				}
				if matchLength == 0 {
					filePathParts = make([]string, 0)
					break
				}
				filePathParts = filePathParts[:matchLength]
			}
			if len(filePathParts) > 0 {
				result = strings.Join(filePathParts, ROOT_ROUTE)
			}
		}
	}
	return result
}

func getPathAfterPattern(pattern string, path string) string {
	result := ""
	if strings.HasSuffix(pattern, WILDCARD_SUFFIX) {
		prefix := strings.TrimSuffix(pattern, WILDCARD_SUFFIX)
		if prefix == "" {
			result = strings.TrimPrefix(path, ROOT_ROUTE)
		} else {
			result = strings.TrimPrefix(strings.TrimPrefix(path, prefix), ROOT_ROUTE)
		}
	} else if pattern == path {
		result = ""
	} else {
		normalizedPattern := strings.TrimSuffix(pattern, ROOT_ROUTE)
		result = strings.TrimPrefix(path, normalizedPattern+ROOT_ROUTE)
	}
	return result
}

func getServeConfigFilePath() string {
	result := SERVE_CONFIG_FILE_NAME
	var err error
	var executablePath string
	if executablePath, err = os.Executable(); err == nil {
		result = filepath.Join(filepath.Dir(executablePath), SERVE_CONFIG_FILE_NAME)
	} else {
		var workingDirectory string
		if workingDirectory, err = os.Getwd(); err == nil {
			result = filepath.Join(workingDirectory, SERVE_CONFIG_FILE_NAME)
		}
	}
	return result
}

func getSqliteDatabaseFilePath() string {
	result := SERVE_DATABASE_FILE_NAME
	var err error
	var executablePath string
	if executablePath, err = os.Executable(); err == nil {
		result = filepath.Join(filepath.Dir(executablePath), SERVE_DATABASE_FILE_NAME)
	}
	return result
}

func handleGetSalt(request *http.Request, response http.ResponseWriter) error {
	result := error(nil)
	if request.TLS == nil {
		__debug(fmt.Sprintf("[Salt] Request rejected: not a TLS/HTTPS connection (remote=%s, proto=%s)", request.RemoteAddr, request.Proto))
		http.Error(response, SALT_RESPONSE_NO_CERTIFICATE, http.StatusForbidden)
	} else if len(request.TLS.PeerCertificates) == 0 {
		__debug(fmt.Sprintf("[Salt] Request rejected: TLS connection without client certificate (remote=%s)", request.RemoteAddr))
		http.Error(response, SALT_RESPONSE_NO_CERTIFICATE, http.StatusForbidden)
	} else {
		clientCertificate := request.TLS.PeerCertificates[0]
		if !verifyClientCertificate(clientCertificate) {
			__debug(fmt.Sprintf("[Salt] Request rejected: client certificate not issued by trusted CA (remote=%s, subject=%s, issuer=%s)", request.RemoteAddr, clientCertificate.Subject.String(), clientCertificate.Issuer.String()))
			http.Error(response, SALT_RESPONSE_NO_CERTIFICATE, http.StatusForbidden)
		} else if isCertificateRevoked(clientCertificate) {
			__debug(fmt.Sprintf("[Salt] Request rejected: client certificate revoked (remote=%s, subject=%s, issuer=%s)", request.RemoteAddr, clientCertificate.Subject.String(), clientCertificate.Issuer.String()))
			http.Error(response, SALT_RESPONSE_NO_CERTIFICATE, http.StatusForbidden)
		} else {
			fingerprint := getCertificateFingerprint(clientCertificate)
			saltMutex.RLock()
			currentSalt := salt
			saltMutex.RUnlock()
			hashValue := hash.SHA3(fingerprint + currentSalt)
			response.Header().Set(CONTENT_TYPE, CONTENT_TYPE_JSON)
			payload := map[string]string{SALT_RESPONSE_HASH_KEY: hashValue}
			var data []byte
			var err error
			if data, err = json.Marshal(payload); err == nil {
				var writeResult int
				if writeResult, err = response.Write(data); err != nil {
					result = err
				} else if writeResult == 0 {
					result = fmt.Errorf("failed to write response")
				} else {
					__debug(fmt.Sprintf("[Salt] Response sent successfully: remote=%s", request.RemoteAddr))
				}
			} else {
				result = err
			}
		}
	}
	return result
}

type responseWriterWrapper struct {
	http.ResponseWriter
	wroteHeader bool
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.wroteHeader = true
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriterWrapper) Write(data []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(data)
}

//goland:noinspection SpellCheckingInspection
func handleRequest(responseWriter http.ResponseWriter, request *http.Request) {
	snapshot := createConfigurationSnapshot()
	isRequestForbidden := false
	isTusdHandled := false
	isRouteHandled := false
	runBeforeRequestHandlers := func() bool {
		result := false
		for _, handler := range snapshot.BeforeRequestHandlers {
			if !handler(request, responseWriter) {
				result = true
				break
			}
		}
		return result
	}
	isProtectedRequest := func() bool {
		result := false
		if !IsLoopbackAddress(request) && request.TLS != nil && isProtectedURI(request.URL.Path, snapshot) {
			result = true
		}
		return result
	}
	validateClientCertificate := func() bool {
		result := false
		writeClientCertificateToLogger(request)
		if len(request.TLS.PeerCertificates) == 0 {
			__debug("[mTLS] No client certificate provided")
			http.Error(responseWriter, "Client certificate required", http.StatusForbidden)
			result = true
		} else {
			clientCertificate := request.TLS.PeerCertificates[0]
			fingerprint := getCertificateFingerprint(clientCertificate)
			if !verifyClientCertificate(clientCertificate) {
				__debug(fmt.Sprintf("[mTLS] Certificate verification failed: CN=%s, fingerprint=%s", clientCertificate.Subject.CommonName, fingerprint))
				http.Error(responseWriter, "Certificate verification failed", http.StatusForbidden)
				result = true
			} else if isCertificateRevoked(clientCertificate) {
				__debug(fmt.Sprintf("[mTLS] Certificate revoked: CN=%s, fingerprint=%s", clientCertificate.Subject.CommonName, fingerprint))
				http.Error(responseWriter, "Certificate revoked", http.StatusForbidden)
				result = true
			}
		}
		return result
	}
	validateReplayProtection := func() bool {
		result := false
		requestUUID := request.Header.Get(REQUEST_UUID_HEADER_NAME)
		requestTimestamp := request.Header.Get(REQUEST_TIMESTAMP_HEADER_NAME)
		if isReplayAttack(requestUUID, requestTimestamp) {
			__debug(fmt.Sprintf("[mTLS] Replay or forged request detected: uuid=%s, timestamp=%s, path=%s", requestUUID, requestTimestamp, request.URL.Path))
			http.Error(responseWriter, "replay or forged request forbidden", http.StatusForbidden)
			result = true
		}
		return result
	}
	isJSONBodyHashValidationRequired := func() bool {
		result := false
		if GetEnableJsonBodyHashValidation() && request.URL.Path != WEBAPI_PATH_SALT && strings.Contains(strings.ToLower(request.Header.Get(CONTENT_TYPE)), CONTENT_TYPE_JSON) {
			result = true
		}
		return result
	}
	validateJSONBodyHash := func() bool {
		result := false
		if len(request.TLS.PeerCertificates) > 0 {
			clientCertificate := request.TLS.PeerCertificates[0]
			fingerprint := getCertificateFingerprint(clientCertificate)
			saltMutex.RLock()
			currentSalt := salt
			saltMutex.RUnlock()
			certificateSaltHash := hash.SHA3(fingerprint + currentSalt)
			requestHashHeader := request.Header.Get(REQUEST_HASH_HEADER_NAME)
			if requestHashHeader == "" {
				__debug(fmt.Sprintf("[mTLS] Missing %s header: path=%s", REQUEST_HASH_HEADER_NAME, request.URL.Path))
				http.Error(responseWriter, REQUEST_HASH_FORBIDDEN_RESPONSE, http.StatusForbidden)
				result = true
			} else {
				var bodyBytes []byte
				var err error
				if bodyBytes, err = io.ReadAll(request.Body); err != nil {
					__debug(fmt.Sprintf("[mTLS] Failed to read request body for hash validation: %v", err))
					http.Error(responseWriter, REQUEST_HASH_FORBIDDEN_RESPONSE, http.StatusForbidden)
					result = true
				} else {
					_ = request.Body.Close()
					request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
					expectedHash := hash.SHA3(string(bodyBytes) + certificateSaltHash)
					if !strings.EqualFold(expectedHash, requestHashHeader) {
						__debug(fmt.Sprintf("[mTLS] Request hash mismatch: path=%s, expected=%s, actual=%s", request.URL.Path, expectedHash, requestHashHeader))
						http.Error(responseWriter, REQUEST_HASH_FORBIDDEN_RESPONSE, http.StatusForbidden)
						result = true
					}
				}
			}
		} else {
			__debug(fmt.Sprintf("[mTLS] No client certificate for JSON body hash validation: path=%s", request.URL.Path))
			http.Error(responseWriter, REQUEST_HASH_FORBIDDEN_RESPONSE, http.StatusForbidden)
			result = true
		}
		return result
	}
	validateProtectedRequest := func() bool {
		result := false
		if isProtectedRequest() {
			if validateClientCertificate() {
				result = true
			} else if validateReplayProtection() {
				result = true
			} else if isJSONBodyHashValidationRequired() && validateJSONBodyHash() {
				result = true
			}
		}
		return result
	}
	runBeforeTusdHandlers := func() bool {
		result := false
		for _, handler := range snapshot.BeforeTusdHandlers {
			if !handler(request, responseWriter) {
				result = true
				break
			}
		}
		return result
	}
	handleTusdRequest := func() bool {
		result := false
		if tusdMount := searchTusdMount(snapshot.TusdMounts, request.URL.Path); tusdMount != nil && tusdMount.Handler != nil {
			if runBeforeTusdHandlers() {
				result = true
			} else {
				strippedPath := strings.TrimPrefix(request.URL.Path, tusdMount.Uri)
				if strippedPath == request.URL.Path {
					strippedPath = strings.TrimPrefix(request.URL.Path, strings.TrimSuffix(tusdMount.Uri, ROOT_ROUTE))
				}
				if strippedPath == "" {
					strippedPath = ROOT_ROUTE
				}
				request.URL.Path = strippedPath
				tusdMount.Handler.ServeHTTP(responseWriter, request)
				result = true
			}
		}
		return result
	}
	runBeforeRouteHandlers := func() bool {
		result := false
		for _, handler := range snapshot.BeforeRouteHandlers {
			if !handler(request, responseWriter) {
				result = true
				break
			}
		}
		return result
	}
	handleRouteRequest := func() bool {
		result := false
		if runBeforeRouteHandlers() {
			result = true
		} else {
			for _, route := range snapshot.Routes {
				if (route.Method == ANY_METHOD || route.Method == request.Method) && MatchPath(route.Pattern, request.URL.Path) {
					var err error
					wrappedWriter := &responseWriterWrapper{ResponseWriter: responseWriter}
					if err = route.Handler(request, wrappedWriter); err != nil {
						if !wrappedWriter.wroteHeader {
							http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
						}
						__debug(err.Error())
					}
					result = true
					break
				}
			}
		}
		return result
	}
	runBeforeServeStaticHandlers := func() bool {
		result := false
		if snapshot.BeforeServeStaticHandlers != nil {
			for _, handler := range snapshot.BeforeServeStaticHandlers {
				if !handler(request, responseWriter) {
					result = true
					break
				}
			}
		}
		return result
	}
	runAfterServeStaticHandlers := func() {
		if snapshot.AfterServeStaticHandlers != nil {
			for _, handler := range snapshot.AfterServeStaticHandlers {
				handler(request, responseWriter, request.URL.Path)
			}
		}
	}
	handleStaticRequest := func() {
		if !runBeforeServeStaticHandlers() {
			serveStatic(responseWriter, request, snapshot)
			runAfterServeStaticHandlers()
		}
	}

	isRequestForbidden = runBeforeRequestHandlers()
	if !isRequestForbidden {
		isRequestForbidden = validateProtectedRequest()
	}
	if !isRequestForbidden {
		isTusdHandled = handleTusdRequest()
	}
	if !isRequestForbidden && !isTusdHandled {
		isRouteHandled = handleRouteRequest()
	}
	if !isRequestForbidden && !isTusdHandled && !isRouteHandled {
		handleStaticRequest()
	}
}

//goland:noinspection SpellCheckingInspection
func handleTusdCompleted(mount *TUSD_MOUNT, event tushandler.HookEvent) {
	__debug(fmt.Sprintf("TUSd upload completed - ID: %s, MetaData: %v", event.Upload.ID, event.Upload.MetaData))
	uploadedFilePath := filepath.Join(mount.DirectoryPath, event.Upload.ID)
	targetFileName := event.Upload.ID
	targetDirectory := mount.DirectoryPath
	if mount.TargetDirectory != "" {
		targetDirectory = mount.TargetDirectory
	}
	fileName := event.Upload.ID
	if name, ok := event.Upload.MetaData[FILENAME_KEY]; ok && name != "" {
		fileName = name
	}
	var err error
	if err = os.MkdirAll(targetDirectory, DEFAULT_DIRECTORY_PERMISSION); err != nil {
		__debug(fmt.Sprintf("Failed to create target directory: %s, error: %v", targetDirectory, err))
	}
	targetFilePath := filepath.Join(targetDirectory, fileName)
	if _, err = os.Stat(targetFilePath); err == nil {
		extension := filepath.Ext(fileName)
		baseName := fileName[:len(fileName)-len(extension)]
		counter := 1
		for {
			newFileName := fmt.Sprintf(FILE_RENAME_FORMAT, baseName, counter, extension)
			targetFilePath = filepath.Join(targetDirectory, newFileName)
			if _, err = os.Stat(targetFilePath); os.IsNotExist(err) {
				break
			}
			counter++
		}
	}
	targetFileName = filepath.Base(targetFilePath)
	if mount.UploadedHandler != nil {
		mount.UploadedHandler(mount.Uri, uploadedFilePath, targetFileName, event.Upload)
	}
}

func initializeSqliteDatabase() {
	sqliteDatabaseFilePath = getSqliteDatabaseFilePath()
	sqliteDatabase = sqlite.New()
	sqliteDatabase.PragmaKey = SQLITE_DATABASE_DEFAULT_PASSWORD
	if err := sqliteDatabase.Create(sqliteDatabaseFilePath); err != nil {
		__debug(fmt.Sprintf("Failed to create serve database at %s: %v", sqliteDatabaseFilePath, err))
		sqliteDatabase.Close()
		if err = os.Remove(sqliteDatabaseFilePath); err != nil && !os.IsNotExist(err) {
			__debug(fmt.Sprintf("Failed to remove corrupted serve database file %s: %v", sqliteDatabaseFilePath, err))
		}
		sqliteDatabase = sqlite.New()
		sqliteDatabase.PragmaKey = SQLITE_DATABASE_DEFAULT_PASSWORD
		if err = sqliteDatabase.Create(sqliteDatabaseFilePath); err != nil {
			__debug(fmt.Sprintf("Failed to recreate serve database at %s: %v", sqliteDatabaseFilePath, err))
			sqliteDatabase.Close()
			sqliteDatabase = nil
		}
	}
	if sqliteDatabase != nil {
		if err := sqliteDatabase.ExecNonQuery(SERVE_SQL_CREATE_TABLE); err != nil {
			__debug(fmt.Sprintf("Failed to create serve database table: %v", err))
			sqliteDatabase.Close()
			sqliteDatabase = nil
		} else if err = sqliteDatabase.ExecNonQuery(SERVE_SQL_CREATE_INDEX); err != nil {
			__debug(fmt.Sprintf("Failed to create serve database index: %v", err))
			sqliteDatabase.Close()
			sqliteDatabase = nil
		}
	}
}

//goland:noinspection SpellCheckingInspection
func initializeTusdMounts() error {
	result := error(nil)
	tusdMutex.Lock()
	mounts := make(map[string]*TUSD_MOUNT)
	for key, value := range tusdMounts {
		mounts[key] = value
	}
	tusdMutex.Unlock()
	for uri, mount := range mounts {
		if mount.Handler == nil {
			if result = os.MkdirAll(mount.DirectoryPath, DEFAULT_DIRECTORY_PERMISSION); result == nil {
				store := filestore.New(mount.DirectoryPath)
				locker := filelocker.New(mount.DirectoryPath)
				composer := tushandler.NewStoreComposer()
				store.UseIn(composer)
				locker.UseIn(composer)
				handlerConfig := tushandler.Config{
					StoreComposer:           composer,
					BasePath:                uri,
					MaxSize:                 mount.MaxSize,
					NotifyCompleteUploads:   true,
					NotifyTerminatedUploads: true,
					NotifyUploadProgress:    true,
					NotifyCreatedUploads:    true,
				}
				var handler *tushandler.Handler
				if handler, result = tushandler.NewHandler(handlerConfig); result == nil {
					tusdMutex.Lock()
					mount.Handler = handler
					tusdMutex.Unlock()
					currentMount := mount
					go func() {
						for event := range currentMount.Handler.CompleteUploads {
							__debug(fmt.Sprintf("Tusd CompleteUploads event - URI: %s, ID: %s, Size: %d, Offset: %d, MetaData: %v", currentMount.Uri, event.Upload.ID, event.Upload.Size, event.Upload.Offset, event.Upload.MetaData))
							handleTusdCompleted(currentMount, event)
						}
					}()
					go func() {
						for event := range currentMount.Handler.TerminatedUploads {
							__debug(fmt.Sprintf("Tusd TerminatedUploads event - URI: %s, ID: %s, Size: %d, Offset: %d, MetaData: %v", currentMount.Uri, event.Upload.ID, event.Upload.Size, event.Upload.Offset, event.Upload.MetaData))
						}
					}()
					go func() {
						for event := range currentMount.Handler.UploadProgress {
							__debug(fmt.Sprintf("Tusd UploadProgress event - URI: %s, ID: %s, Size: %d, Offset: %d, MetaData: %v", currentMount.Uri, event.Upload.ID, event.Upload.Size, event.Upload.Offset, event.Upload.MetaData))
						}
					}()
					go func() {
						for event := range currentMount.Handler.CreatedUploads {
							__debug(fmt.Sprintf("Tusd CreatedUploads event - URI: %s, ID: %s, Size: %d, Offset: %d, MetaData: %v", currentMount.Uri, event.Upload.ID, event.Upload.Size, event.Upload.Offset, event.Upload.MetaData))
						}
					}()
				}
			} else {
				break
			}
		}
	}
	return result
}

func isCertificateRevoked(certificate *x509.Certificate) bool {
	result := false
	if certificate != nil {
		fingerprint := getCertificateFingerprint(certificate)
		revokedCertificatesMutex.RLock()
		if revokedCertificates[fingerprint] {
			result = true
		}
		revokedCertificatesMutex.RUnlock()
	}
	return result
}

//goland:noinspection SpellCheckingInspection
func isPathInUseForTusd(pattern string, tusdBasePath string) bool {
	result := false
	if pattern == tusdBasePath || pattern == strings.TrimSuffix(tusdBasePath, ROOT_ROUTE) {
		result = true
	}
	if !result && strings.HasPrefix(pattern, tusdBasePath) {
		result = true
	}
	if !result && strings.HasSuffix(pattern, WILDCARD_SUFFIX) {
		prefix := strings.TrimSuffix(pattern, WILDCARD_SUFFIX)
		if strings.HasPrefix(tusdBasePath, prefix+ROOT_ROUTE) || tusdBasePath == prefix+ROOT_ROUTE {
			result = true
		}
	}
	return result
}

func isProtectedURI(path string, snapshot STATE_SNAPSHOT) bool {
	result := false
	isIgnoredPattern := false
	for _, pattern := range snapshot.IgnoredProtectionPatterns {
		if MatchPath(pattern, path) {
			isIgnoredPattern = true
			break
		}
	}
	if !isIgnoredPattern {
		for _, pattern := range snapshot.ProtectedPatterns {
			if MatchPath(pattern, path) {
				result = true
				break
			}
		}
	}
	return result
}

func isReplayAttack(uuid string, timestamp string) bool {
	sqliteDatabaseInitOnce.Do(initializeSqliteDatabase)
	result := false
	if uuid == "" || timestamp == "" {
		result = true
	} else if sqliteDatabase == nil {
		result = true
	} else {
		if sqliteDatabase.Exec(SERVE_SQL_INSERT_REQUEST, uuid, timestamp, SERVE_RECORD_TYPE_REQUEST, time.Now().Unix()) != nil {
			result = true
		}
	}
	return result
}

func loadDefaultCertificatePool() *x509.CertPool {
	var result *x509.CertPool
	__debug("[Certificate] Loading default CA certificate pool...")
	var certificatePath string
	if certificatePath = findCertificateFile(ca.CERTIFICATE_AUTHORITY_CERTIFICATE_FILE_NAME); certificatePath != "" {
		__debug(fmt.Sprintf("[Certificate] Loading CA certificate from: %s", certificatePath))
		var certificateBytes []byte
		var err error
		if certificateBytes, err = os.ReadFile(certificatePath); err == nil {
			pool := x509.NewCertPool()
			if pool.AppendCertsFromPEM(certificateBytes) {
				result = pool
				__debug(fmt.Sprintf("[Certificate] CA certificate pool loaded successfully from: %s", certificatePath))
			} else {
				__debug(fmt.Sprintf("[Certificate] Failed to append CA certificate to pool from: %s", certificatePath))
			}
		} else {
			__debug(fmt.Sprintf("[Certificate] Failed to read CA certificate file: %v", err))
		}
	} else {
		__debug("[Certificate] CA certificate file not found")
	}
	return result
}

//goland:noinspection DuplicatedCode
func loadServeConfig() {
	configPath := findServeConfigFile()
	if configPath == "" {
		createDefaultServeConfig()
	} else {
		var data []byte
		var err error
		if data, err = os.ReadFile(configPath); err == nil {
			var config SERVE_CONFIG
			if err = json.Unmarshal(data, &config); err == nil {
				serveConfig = &config
				revokedCertificatesMutex.Lock()
				revokedCertificates = make(map[string]bool)
				for _, certificate := range config.RevokedCertificates {
					revokedCertificates[certificate] = true
				}
				revokedCertificatesMutex.Unlock()
				__debug("Loaded serve config with " + strconv.Itoa(len(config.RevokedCertificates)) + " revoked certificates")
			} else {
				__debug("Failed to parse serve config file: " + err.Error())
				createDefaultServeConfig()
			}
		} else {
			__debug("Failed to read serve config file: " + err.Error())
			createDefaultServeConfig()
		}
	}
}

func MatchPath(pattern string, path string) bool {
	result := false
	if pattern == ROOT_ROUTE {
		result = true
	}
	if !result && pattern == path {
		result = true
	}
	if !result && strings.HasSuffix(pattern, WILDCARD_SUFFIX) {
		prefix := strings.TrimSuffix(pattern, WILDCARD_SUFFIX)
		if prefix == "" || strings.HasPrefix(path, prefix+ROOT_ROUTE) || path == prefix {
			result = true
		}
	}
	if !result && pattern != ROOT_ROUTE {
		normalizedPattern := strings.TrimSuffix(pattern, ROOT_ROUTE)
		if strings.HasPrefix(path, normalizedPattern+ROOT_ROUTE) {
			result = true
		}
	}
	if !result && pattern != ROOT_ROUTE && path != ROOT_ROUTE {
		patternParts := strings.Split(pattern, ROOT_ROUTE)
		pathParts := strings.Split(path, ROOT_ROUTE)
		if len(patternParts) == len(pathParts) {
			partsMatch := true
			for index, patternPart := range patternParts {
				if patternPart == "" {
					continue
				}
				if strings.HasPrefix(patternPart, PATH_PARAMETER_PREFIX) && strings.HasSuffix(patternPart, PATH_PARAMETER_SUFFIX) {
					continue
				}
				if patternPart != pathParts[index] {
					partsMatch = false
					break
				}
			}
			if partsMatch {
				result = true
			}
		}
	}
	return result
}

func registerRequestHandler() {
	serveMutex.Lock()
	if !requestHandlerRegistered {
		http.HandleFunc(ROOT_ROUTE, handleRequest)
		requestHandlerRegistered = true
	}
	serveMutex.Unlock()
}

func saveServeConfig() {
	if serveConfig != nil {
		revokedCertificatesMutex.RLock()
		serveConfig.RevokedCertificates = make([]string, 0, len(revokedCertificates))
		for certificate := range revokedCertificates {
			serveConfig.RevokedCertificates = append(serveConfig.RevokedCertificates, certificate)
		}
		revokedCertificatesMutex.RUnlock()
		var data []byte
		var err error
		if data, err = json.MarshalIndent(serveConfig, "", "  "); err == nil {
			configPath := getServeConfigFilePath()
			if err = os.WriteFile(configPath, data, FILE_PERMISSION_PUBLIC); err != nil {
				__debug("Failed to write serve config file: " + err.Error())
			}
		} else {
			__debug("Failed to marshal serve config: " + err.Error())
		}
	}
}

func searchEmbedFilePath(embedFileSystem embed.FS, filePath string) (fs.FileInfo, string, error) {
	var result os.FileInfo
	err := error(nil)
	embedFilePath := ""
	if result, err = fs.Stat(embedFileSystem, filePath); err == nil {
		embedFilePath = filePath
	} else {
		var rootDirectoryName string
		if rootDirectoryName = getEmbedFsRootDirectory(embedFileSystem); rootDirectoryName != "" {
			targetFilePath := rootDirectoryName
			if filePath != DOT {
				targetFilePath = rootDirectoryName + ROOT_ROUTE + filePath
			}
			if result, err = fs.Stat(embedFileSystem, targetFilePath); err == nil {
				embedFilePath = targetFilePath
			}
		}
	}
	return result, embedFilePath, err
}

//goland:noinspection SpellCheckingInspection
func searchTusdMount(mounts map[string]*TUSD_MOUNT, path string) *TUSD_MOUNT {
	var result *TUSD_MOUNT
	bestLength := 0
	for basePath, mount := range mounts {
		if strings.HasPrefix(path, basePath) || path == strings.TrimSuffix(basePath, ROOT_ROUTE) {
			if len(basePath) > bestLength {
				result = mount
				bestLength = len(basePath)
			}
		}
	}
	return result
}

//goland:noinspection SpellCheckingInspection
func serveStatic(responseWriter http.ResponseWriter, request *http.Request, snapshot STATE_SNAPSHOT) {
	__debug("serveStatic called for path: " + request.URL.Path)
	httpStatusCode := http.StatusNotFound
	writeErrorMessage := ""
	isEmbedFileSystem := false
	physicalFilePath := ""
	embedFileSystem := embed.FS{}
	isDirectory := false
	__debug(fmt.Sprintf("Checking %d static directories", len(snapshot.StaticDirectories)))
	for urlPath, directoryPath := range snapshot.StaticDirectories {
		__debug("Checking static directory mapping: " + urlPath + " -> " + directoryPath)
		if MatchPath(urlPath, request.URL.Path) {
			__debug("URL path matches static directory: " + urlPath)
			var absoluteDirectoryPath string
			var err error
			if absoluteDirectoryPath, err = filepath.Abs(directoryPath); err != nil {
				absoluteDirectoryPath = directoryPath
			}
			requestFilePath := getPathAfterPattern(urlPath, request.URL.Path)
			absoluteFilePath := filepath.Join(absoluteDirectoryPath, requestFilePath)
			__debug("Resolved absolute file path: " + absoluteFilePath)
			var stat os.FileInfo
			if stat, err = os.Stat(absoluteFilePath); !os.IsNotExist(err) {
				__debug("File exists: " + absoluteFilePath + ", isDir: " + strconv.FormatBool(stat.IsDir()))
				if stat.IsDir() {
					defaultIndexHTML := filepath.Join(absoluteFilePath, DEFAULT_INDEX_HTML)
					if _, err = os.Stat(defaultIndexHTML); !os.IsNotExist(err) {
						__debug("Found index.html in directory: " + defaultIndexHTML)
						physicalFilePath = defaultIndexHTML
						isEmbedFileSystem = false
						httpStatusCode = http.StatusOK
					} else if isListDirectoryEnabled {
						__debug("Directory listing enabled for: " + absoluteFilePath)
						isDirectory = true
						isEmbedFileSystem = false
						physicalFilePath = absoluteFilePath
						httpStatusCode = http.StatusOK
					} else {
						__debug("Directory listing disabled, returning Forbidden")
						httpStatusCode = http.StatusForbidden
					}
				} else {
					__debug("Serving static file: " + absoluteFilePath)
					physicalFilePath = absoluteFilePath
					isEmbedFileSystem = false
					httpStatusCode = http.StatusOK
				}
				break
			} else {
				__debug("File not found in static directory: " + absoluteFilePath)
			}
		}
	}
	if physicalFilePath == "" {
		__debug(fmt.Sprintf("Checking %d embed directories", len(snapshot.EmbedDirectories)))
		for embedDirectoryNameMapToUriPath, embeddedFS := range snapshot.EmbedDirectories {
			__debug("Checking embed directory mapping: " + embedDirectoryNameMapToUriPath)
			if MatchPath(embedDirectoryNameMapToUriPath, request.URL.Path) {
				__debug("URL path matches embed directory: " + embedDirectoryNameMapToUriPath)
				requestEmbedFilePath := getPathAfterPattern(embedDirectoryNameMapToUriPath, request.URL.Path)
				if requestEmbedFilePath == "" {
					requestEmbedFilePath = DOT
				}
				__debug("Request file path in embed: " + requestEmbedFilePath)
				var stat fs.FileInfo
				var embedFilePath string
				var err error
				if stat, embedFilePath, err = searchEmbedFilePath(embeddedFS, requestEmbedFilePath); err == nil {
					__debug("Found in embed: " + embedFilePath + ", is directory: " + strconv.FormatBool(stat.IsDir()))
					if stat.IsDir() || strings.HasSuffix(requestEmbedFilePath, ROOT_ROUTE) {
						relativeIndexHtmlFilePath := ""
						if requestEmbedFilePath == DOT {
							relativeIndexHtmlFilePath = DEFAULT_INDEX_HTML
						} else {
							relativeIndexHtmlFilePath = strings.TrimSuffix(requestEmbedFilePath, ROOT_ROUTE) + ROOT_ROUTE + DEFAULT_INDEX_HTML
						}
						__debug("Checking embed index.html: " + relativeIndexHtmlFilePath)
						var embedIndexHtmlFilePath string
						if _, embedIndexHtmlFilePath, err = searchEmbedFilePath(embeddedFS, relativeIndexHtmlFilePath); err == nil {
							__debug("Found index.html in embed directory: " + embedIndexHtmlFilePath)
							physicalFilePath = embedIndexHtmlFilePath
							isEmbedFileSystem = true
							embedFileSystem = embeddedFS
							httpStatusCode = http.StatusOK
						} else if isListDirectoryEnabled {
							__debug("Directory listing enabled for embed directory: " + embedFilePath)
							isDirectory = true
							isEmbedFileSystem = true
							embedFileSystem = embeddedFS
							physicalFilePath = embedFilePath
							httpStatusCode = http.StatusOK
						} else {
							__debug("Directory listing disabled for embed, returning Forbidden")
							httpStatusCode = http.StatusForbidden
						}
					} else {
						__debug("Serving embed file: " + embedFilePath)
						physicalFilePath = embedFilePath
						isEmbedFileSystem = true
						embedFileSystem = embeddedFS
						httpStatusCode = http.StatusOK
					}
					break
				} else {
					__debug("File not found in embed: " + requestEmbedFilePath)
				}
			}
		}
	}
	if httpStatusCode == http.StatusOK {
		__debug(fmt.Sprintf("Serving content with status OK, isDirectory: %v, isEmbedFs: %v, path: %s", isDirectory, isEmbedFileSystem, physicalFilePath))
		if mimeType := mime.TypeByExtension(filepath.Ext(physicalFilePath)); mimeType != "" {
			responseWriter.Header().Set(CONTENT_TYPE, mimeType)
		}
		if isDirectory {
			var embedFileSystemPointer *embed.FS
			if isEmbedFileSystem {
				embedFileSystemPointer = &embedFileSystem
			}
			var html string
			if html = getDirectoryListHTML(physicalFilePath, request.URL.Path, embedFileSystemPointer); html == "" {
				__debug("Failed to generate directory listing HTML")
				httpStatusCode = http.StatusInternalServerError
			} else {
				responseWriter.Header().Set(CONTENT_TYPE, CONTENT_TYPE_HTML_UTF8)
				var err error
				if _, err = responseWriter.Write([]byte(html)); err != nil {
					__debug("Failed to write directory listing: " + err.Error())
					httpStatusCode = http.StatusInternalServerError
					writeErrorMessage = fmt.Sprintf("%s write error: %s", request.URL.Path, err.Error())
				} else {
					__debug("Successfully served directory listing")
				}
			}
		} else {
			if isEmbedFileSystem {
				__debug("Calling http.ServeFileFS for embed file")
				http.ServeFileFS(responseWriter, request, embedFileSystem, physicalFilePath)
			} else {
				__debug("Calling http.ServeFile for static file")
				http.ServeFile(responseWriter, request, physicalFilePath)
			}
		}
	} else {
		__debug("Not serving content, status code: " + strconv.Itoa(httpStatusCode))
	}
	switch httpStatusCode {
	case http.StatusNotFound:
		http.Error(responseWriter, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		__debug(fmt.Sprintf("%s not found", request.URL.Path))
	case http.StatusForbidden:
		http.Error(responseWriter, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		__debug(fmt.Sprintf(FORBIDDEN_LOG_FORMAT, request.URL.Path))
	case http.StatusInternalServerError:
		http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		__debug(writeErrorMessage)
	}
	__debug("serveStatic completed")
}

func shutdown() error {
	result := error(nil)
	serveMutex.Lock()
	if serverListener == nil {
		result = fmt.Errorf("server is not running")
	}
	serveMutex.Unlock()
	if result == nil {
		result = server.Shutdown(serverContext)
		serveMutex.Lock()
		serverListener = nil
		serveMutex.Unlock()
	}
	return result
}

func shutdownTls() error {
	result := error(nil)
	serveMutex.Lock()
	if tlsServerListener == nil {
		result = fmt.Errorf("TLS server is not running")
	}
	serveMutex.Unlock()
	if result == nil {
		result = tlsServer.Shutdown(tlsServerContext)
		serveMutex.Lock()
		tlsServerListener = nil
		serveMutex.Unlock()
	}
	return result
}

func verifyClientCertificate(certificate *x509.Certificate) bool {
	result := false
	if certificate != nil {
		serveMutex.Lock()
		pool := tlsClientCAPool
		serveMutex.Unlock()
		if pool != nil {
			verifyOptions := x509.VerifyOptions{
				Roots:     pool,
				KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			}
			var err error
			if _, err = certificate.Verify(verifyOptions); err == nil {
				result = true
			} else {
				__debug(fmt.Sprintf("[mTLS] Certificate verification failed: %v", err))
			}
		} else {
			__debug("[mTLS] CA certificate pool is not loaded, cannot verify client certificates")
		}
	} else {
		__debug("[mTLS] Cannot verify nil certificate")
	}
	return result
}

func waitServe(address string) {
	start := time.Now()
	for {
		var connection net.Conn
		var err error
		if connection, err = net.DialTimeout(TCP, address, WAIT_CLIENT_TIMEOUT); err == nil {
			_ = connection.Close()
			break
		}
		if time.Since(start) >= WAIT_CLIENT_TIMEOUT {
			break
		}
		time.Sleep(WAIT_RETRY_INTERVAL)
	}
}

//goland:noinspection DuplicatedCode
func waitServeTls(address string) {
	start := time.Now()
	for {
		var connection *tls.Conn
		err := error(nil)
		if connection, err = tls.Dial(TCP, address, &tls.Config{InsecureSkipVerify: true}); err == nil {
			go func() {
				time.Sleep(TLS_HANDSHAKE_ERROR_WAIT_DURATION)
				__warning("The error can be safely ignored: 'http: TLS handshake error from 127.0.0.1:56402: tls: client didn't provide a certificate'")
			}()
			_ = connection.Close()
			break
		}
		if time.Since(start) >= WAIT_CLIENT_TIMEOUT {
			break
		}
		time.Sleep(WAIT_RETRY_INTERVAL)
	}
}

func writeClientCertificateToLogger(request *http.Request) {
	if request.TLS != nil && len(request.TLS.PeerCertificates) > 0 {
		clientCertificate := request.TLS.PeerCertificates[0]
		fingerprint := getCertificateFingerprint(clientCertificate)
		b := false
		loggedClientCertificatesMutex.Lock()
		if loggedClientCertificates == nil {
			loggedClientCertificates = make(map[string]bool)
		}
		if !loggedClientCertificates[fingerprint] {
			loggedClientCertificates[fingerprint] = true
			b = true
		}
		loggedClientCertificatesMutex.Unlock()
		if b {
			verified := verifyClientCertificate(clientCertificate)
			__debug(fmt.Sprintf("[mTLS] Client authenticated: CN=%s, fingerprint=%s, verifiedByCA=%t", clientCertificate.Subject.CommonName, fingerprint, verified))
		}
	}
}
