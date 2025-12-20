// Package electron
// File:        electron.go
// Url:         https://github.com/xiang-tai-duo/go-bootstrap/blob/master/electron/electron.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Electron is a wrapper for Electron application operations, providing methods for application management.
// --------------------------------------------------------------------------------
package electron

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/google/uuid"
	__websocket "github.com/gorilla/websocket"
	"github.com/xiang-tai-duo/go-boost/debugger"
	"github.com/xiang-tai-duo/go-boost/embed"
	"github.com/xiang-tai-duo/go-boost/file"
	"github.com/xiang-tai-duo/go-boost/hash"
	"github.com/xiang-tai-duo/go-boost/http2"
	"github.com/xiang-tai-duo/go-boost/logger"
	"github.com/xiang-tai-duo/go-boost/process"
	"github.com/xiang-tai-duo/go-boost/serve"
	websocket "github.com/xiang-tai-duo/go-boost/websocket/server"
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
type (
	ELECTRON_EXIT_PROCEDURE func()
	STANDARD_CALLBACK       struct {
		callback func(string)
	}
	TOKEN_DATA struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	WEB_SOCKET_DATA struct {
		Data   interface{} `json:"data"`
		Silent bool        `json:"silent"`
		Type   string      `json:"type"`
	}
	WEB_SOCKET_PORT_DATA struct {
		Port int `json:"port"`
	}
	WEB_SOCKET_PORT_RESPONSE struct {
		WebSocket WEB_SOCKET_PORT_DATA `json:"websocket"`
	}
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName,SpellCheckingInspection
const (
	API_ELECTRON_READY                      = "/api/electron/ready"
	API_ELECTRON_WEB_SOCKET                 = "/api/electron/websocket"
	ARGUMENTS_PARAMETERS                    = "--params"
	CHROMIUM_ARG_ENABLE_FEATURES_DARWIN     = "--enable-features=Metal"
	CHROMIUM_ARG_ENABLE_FEATURES_LINUX      = "--enable-features=VaapiVideoDecoder"
	CHROMIUM_ARG_ENABLE_FEATURES_WINDOWS    = "--enable-features=Vulkan"
	CHROMIUM_ARG_FONT_ANTIALIASING          = "--enable-font-antialiasing"
	CHROMIUM_ARG_NO_SANDBOX                 = "--no-sandbox"
	CHROMIUM_ARG_FORCE_COLOR_PROFILE        = "--force-color-profile=srgb"
	CHROMIUM_ARG_HIGH_DPI_SUPPORT           = "--high-dpi-support=1"
	CONTENT_TYPE_JSON                       = "application/json"
	DANGEROUS_SCRIPT_CLOSE_TAG              = "</script"
	DANGEROUS_SCRIPT_OPEN_TAG               = "<script"
	DEFAULT_WINDOW_HEIGHT                   = -1
	DEFAULT_WINDOW_WIDTH                    = -1
	ELECTRON_EXECUTABLE_NOT_FOUND_ERROR     = "Executable not found: %s or %s"
	ELECTRON_READY_LOG_MESSAGE              = "console.log('[electron.go] BOTH FRONTEND AND BACKEND ARE READY')"
	ELECTRON_UNPACKED_DIRECTORY_PERMISSION  = 0755
	HOME_PAGE_FORMAT                        = "http://127.0.0.1:%d"
	MINIMUM_PRINTABLE_RUNE                  = 0x20
	MODULE_NAME_ELECTRON                    = "electron"
	OS_DARWIN                               = "darwin"
	OS_LINUX                                = "linux"
	OS_WINDOWS                              = "windows"
	PARAMETER_KEY_BACKEND_PID               = "backendPid"
	PARAMETER_KEY_DEBUG_MODE                = "isDebugMode"
	PARAMETER_KEY_HEIGHT                    = "height"
	PARAMETER_KEY_HOME_PAGE                 = "homePage"
	PARAMETER_KEY_MAXIMIZE                  = "isWindowMaximize"
	PARAMETER_KEY_PREVIEW                   = "preview"
	PARAMETER_KEY_RESIZEABLE                = "isWindowResizeable"
	PARAMETER_KEY_WEB_SOCKET_PORT           = "websocketPort"
	PARAMETER_KEY_WIDTH                     = "width"
	PARAMETER_KEY_WINDOW_SIZE               = "windowSize"
	TOKEN_KEY                               = "ELECTRON_TOKEN"
	URL_SCHEME_SEPARATOR                    = "://"
	WEB_API_PATH                            = "/api/"
	WEB_SOCKET_DATA_TYPE_EVAL               = "eval"
	WEB_SOCKET_DATA_TYPE_EXECUTE_JAVASCRIPT = "executeJavaScript"
	WEB_SOCKET_DATA_TYPE_TOKEN              = "token"
)

//goland:noinspection GoUnhandledErrorResult
var (
	broadcastMutex       sync.Mutex
	electronProcess      = (*exec.Cmd)(nil)
	electronProcessMutex sync.RWMutex
	isDebuggerPresent    = debugger.IsPresent()
	isInitialized        = false
	isInitializedMutex   sync.Mutex
	isWindowMaximize     = true
	isWindowResizeable   = true
	standardError        = func(s string) {}
	standardOutput       = func(s string) {}
	token                = ""
	webSocketPort        int
	webSocketServer      *websocket.WEB_SOCKET_SERVER
)

func init() {
	standardError = func(s string) {
		__debug(s)
	}
	standardOutput = func(s string) {
		__debug(s)
	}
}

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_ELECTRON, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_ELECTRON, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_ELECTRON, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_ELECTRON, logger.SKIP_STACK_FRAMES_BASE)
}

func CreateDefaultHomePage(port int) string {
	result := ""
	result = fmt.Sprintf(HOME_PAGE_FORMAT, port)
	return result
}

//goland:noinspection GoUnusedExportedFunction
func DisableDebugMode() {
	isDebuggerPresent = false
}

//goland:noinspection GoUnusedExportedFunction
func EnableDebugMode() {
	isDebuggerPresent = true
}

//goland:noinspection GoUnusedExportedFunction
func Eval(code string) error {
	err := error(nil)
	if err = sanitizeJSCode(code); err == nil {
		err = sendElectronData(code, WEB_SOCKET_DATA_TYPE_EVAL, false)
	}
	return err
}

//goland:noinspection GoUnusedExportedFunction
func EvalSilent(code string) error {
	err := error(nil)
	if err = sanitizeJSCode(code); err == nil {
		err = sendElectronData(code, WEB_SOCKET_DATA_TYPE_EVAL, true)
	}
	return err
}

//goland:noinspection GoUnusedExportedFunction
func ExecuteJavaScript(code string) error {
	err := error(nil)
	if err = sanitizeJSCode(code); err == nil {
		err = sendElectronData(code, WEB_SOCKET_DATA_TYPE_EXECUTE_JAVASCRIPT, false)
	}
	return err
}

//goland:noinspection GoUnusedExportedFunction
func ExecuteJavaScriptSilent(code string) error {
	err := error(nil)
	if err = sanitizeJSCode(code); err == nil {
		err = sendElectronData(code, WEB_SOCKET_DATA_TYPE_EXECUTE_JAVASCRIPT, true)
	}
	return err
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func InitializeBackendServices() {
	isInitializedMutex.Lock()
	defer isInitializedMutex.Unlock()
	if !isInitialized {
		isInitialized = true
		initializeWebSocketServer()
		initializeServe()
	}
}

//goland:noinspection GoUnusedExportedFunction,GoImportUsedAsName
func IsElectronRunning() bool {
	result := false
	electronProcessMutex.RLock()
	process := electronProcess
	electronProcessMutex.RUnlock()
	if process != nil && process.Process != nil {
		if process.ProcessState == nil {
			result = true
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func IsWindowMaximize() bool {
	return isWindowMaximize
}

//goland:noinspection GoUnusedExportedFunction
func IsWindowResizeable() bool {
	return isWindowResizeable
}

//goland:noinspection GoUnhandledErrorResult,GoUnusedExportedFunction
func Launch(homePage string) (*exec.Cmd, error) {
	var result *exec.Cmd
	err := error(nil)
	webSocketPortValue := 0
	if webSocketPortValue, err = getWebSocketPort(homePage); err == nil {
		__debug(fmt.Sprintf("Launching Electron application with home page: %s, port: %d, isMaximize: %t", homePage, webSocketPortValue, isWindowMaximize))
		var electronExecuteFilePath string
		if electronExecuteFilePath, err = prepareElectronEnvironment(); err == nil {
			params := map[string]interface{}{
				PARAMETER_KEY_BACKEND_PID:     os.Getpid(),
				PARAMETER_KEY_DEBUG_MODE:      isDebuggerPresent,
				PARAMETER_KEY_HOME_PAGE:       homePage,
				PARAMETER_KEY_MAXIMIZE:        IsWindowMaximize(),
				PARAMETER_KEY_RESIZEABLE:      IsWindowResizeable(),
				PARAMETER_KEY_WEB_SOCKET_PORT: webSocketPortValue,
				PARAMETER_KEY_WINDOW_SIZE: map[string]int{
					PARAMETER_KEY_WIDTH:  DEFAULT_WINDOW_WIDTH,
					PARAMETER_KEY_HEIGHT: DEFAULT_WINDOW_HEIGHT,
				},
			}
			__debug(fmt.Sprintf("Starting Electron process: %s", electronExecuteFilePath))
			if result, err = startElectronProcess(electronExecuteFilePath, params); err == nil && result != nil {
				setElectronProcess(result)
				__debug(fmt.Sprintf("Electron process started with PID: %d", result.Process.Pid))
				err = result.Wait()
				clearElectronProcess(result)
				__debug("Electron process exited")
			} else {
				__debug(fmt.Sprintf("Start Electron failed: %v", err))
			}
		}
	} else {
		__debug(fmt.Sprintf("Failed to get WebSocket port: %v", err))
	}
	return result, err
}

//goland:noinspection GoUnhandledErrorResult,GoUnusedExportedFunction
func Preview(uri string) (*exec.Cmd, error) {
	var result *exec.Cmd
	err := error(nil)
	__debug(fmt.Sprintf("Launching Electron in preview mode with URI: %s (lock bypassed)", uri))
	var electronExecuteFilePath string
	if electronExecuteFilePath, err = prepareElectronEnvironment(); err == nil {
		params := map[string]interface{}{
			PARAMETER_KEY_DEBUG_MODE: debugger.IsPresent(),
			PARAMETER_KEY_PREVIEW:    uri,
			PARAMETER_KEY_WINDOW_SIZE: map[string]int{
				PARAMETER_KEY_WIDTH:  DEFAULT_WINDOW_WIDTH,
				PARAMETER_KEY_HEIGHT: DEFAULT_WINDOW_HEIGHT,
			},
		}
		__debug(fmt.Sprintf("Starting Electron preview process: %s", electronExecuteFilePath))
		if result, err = startElectronProcess(electronExecuteFilePath, params); err == nil && result != nil {
			__debug(fmt.Sprintf("Electron preview process started with PID: %d", result.Process.Pid))
			err = result.Wait()
			__debug("Electron preview process exited")
		} else {
			__debug(fmt.Sprintf("Start Electron preview failed: %v", err))
		}
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func SetStderr(f func(string)) {
	standardError = f
}

//goland:noinspection GoUnusedExportedFunction
func SetStdout(f func(string)) {
	standardOutput = f
}

//goland:noinspection GoUnusedExportedFunction
func SetWindowMaximize(value bool) {
	isWindowMaximize = value
}

//goland:noinspection GoUnusedExportedFunction
func SetWindowResizable(value bool) {
	isWindowResizeable = value
}

func (w *STANDARD_CALLBACK) Write(p []byte) (n int, err error) {
	w.callback(string(p))
	return len(p), nil
}

func buildChromiumArgs() []string {
	result := make([]string, 0)
	result = append(result, CHROMIUM_ARG_FORCE_COLOR_PROFILE)
	switch runtime.GOOS {
	case OS_WINDOWS:
		result = append(result,
			CHROMIUM_ARG_HIGH_DPI_SUPPORT,
			CHROMIUM_ARG_FONT_ANTIALIASING,
			CHROMIUM_ARG_ENABLE_FEATURES_WINDOWS,
		)
	case OS_DARWIN:
		result = append(result, CHROMIUM_ARG_ENABLE_FEATURES_DARWIN)
	case OS_LINUX:
		result = append(result,
			CHROMIUM_ARG_ENABLE_FEATURES_LINUX,
			CHROMIUM_ARG_NO_SANDBOX,
		)
	}
	return result
}

//goland:noinspection SpellCheckingInspection
func buildIsolatedWindowsEnvironment() []string {
	keys := []string{
		"PATH",
		"SystemRoot",
		"WINDIR",
		"TEMP",
		"TMP",
		"APPDATA",
		"LOCALAPPDATA",
		"USERPROFILE",
		"COMSPEC",
		"PATHEXT",
	}
	result := make([]string, 0, len(keys))
	for _, key := range keys {
		if value := os.Getenv(key); value != "" {
			result = append(result, key+"="+value)
		}
	}
	return result
}

func clearElectronProcess(process *exec.Cmd) {
	electronProcessMutex.Lock()
	defer electronProcessMutex.Unlock()
	if electronProcess == process {
		electronProcess = nil
	}
}

func getUnpackedDirectoryPath() string {
	result := ELECTRON_DISTRIBUTION_PATH
	err := error(nil)
	var workingDirectoryPath string
	if workingDirectoryPath, err = os.Getwd(); err == nil {
		result = filepath.Join(workingDirectoryPath, ELECTRON_DISTRIBUTION_PATH)
	}
	return result
}

func getWebSocketPort(homePage string) (int, error) {
	result := 0
	err := error(nil)
	var parsedAddress *url.URL
	if parsedAddress, err = url.Parse(homePage); err == nil {
		baseAddress := parsedAddress.Scheme + URL_SCHEME_SEPARATOR + parsedAddress.Host
		requestAddress := baseAddress + API_ELECTRON_WEB_SOCKET
		httpClient := http2.New()
		var body string
		var statusCode int
		if body, statusCode, err = httpClient.Get(requestAddress); err == nil {
			if statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices {
				var response WEB_SOCKET_PORT_RESPONSE
				if err = json.Unmarshal([]byte(body), &response); err == nil {
					result = response.WebSocket.Port
				}
			} else {
				err = fmt.Errorf("API request failed with status: %d", statusCode)
			}
		}
		if result == 0 && err == nil {
			err = fmt.Errorf("failed to get websocket port")
		}
	}
	return result, err
}

//goland:noinspection GoUnhandledErrorResult
func initializeServe() {
	serve.OnBeforeStaticServe(func(r *http.Request, w http.ResponseWriter) bool {
		result := false
		if serve.MatchPath(WEB_API_PATH+"*", r.URL.Path) {
			result = true
		} else if debugger.IsPresent() {
			result = true
		} else {
			err := error(nil)
			var cookie *http.Cookie
			if cookie, err = r.Cookie(TOKEN_KEY); err == nil {
				if cookie.Value == token {
					result = true
				}
			}
		}
		if !result {
			http.Error(w, "", http.StatusForbidden)
			__debug(fmt.Sprintf("Access Denied, %s -> %s", r.RemoteAddr, r.URL.String()))
		}
		return result
	})
	serve.On(serve.GET, API_ELECTRON_WEB_SOCKET, func(r *http.Request, w http.ResponseWriter) error {
		err := error(nil)
		w.Header().Set(serve.CONTENT_TYPE, CONTENT_TYPE_JSON)
		w.WriteHeader(http.StatusOK)
		response := WEB_SOCKET_PORT_RESPONSE{
			WebSocket: WEB_SOCKET_PORT_DATA{
				Port: webSocketPort,
			},
		}
		var jsonData []byte
		if jsonData, err = json.Marshal(response); err == nil {
			_, err = w.Write(jsonData)
		}
		return err
	})
	serve.On(serve.POST, API_ELECTRON_READY, func(r *http.Request, w http.ResponseWriter) error {
		err := error(nil)
		w.Header().Set(serve.CONTENT_TYPE, CONTENT_TYPE_JSON)
		w.WriteHeader(http.StatusOK)
		ExecuteJavaScript(ELECTRON_READY_LOG_MESSAGE)
		return err
	})
}

//goland:noinspection GoUnhandledErrorResult
func initializeWebSocketServer() error {
	err := error(nil)
	__debug("Creating new WebSocket server for Electron communication")
	webSocketServer = websocket.NewWebSocketServer()
	token = uuid.New().String()
	webSocketServer.SetConnectHandler(func(connection *__websocket.Conn, uuid string) error {
		err := error(nil)
		__debug(fmt.Sprintf("WebSocket client connected: %s", uuid))
		tokenData := TOKEN_DATA{
			Key:   TOKEN_KEY,
			Value: token,
		}
		tokenJson := make([]byte, 0)
		if tokenJson, err = json.Marshal(tokenData); err == nil {
			sendElectronData(string(tokenJson), WEB_SOCKET_DATA_TYPE_TOKEN, false)
			__debug("Token data broadcasted to WebSocket clients")
		}
		return err
	})
	if webSocketPort, err = webSocketServer.LaunchAsync(websocket.ANY_PORT); err == nil {
		__debug(fmt.Sprintf("WebSocket server started on port: %d", webSocketPort))
	} else {
		__debug(fmt.Sprintf("Failed to start WebSocket server: %s", err.Error()))
	}
	return err
}

func isVSCodeDebugEnvironment() bool {
	result := false
	for _, environment := range os.Environ() {
		key := strings.ToUpper(strings.SplitN(environment, "=", 2)[0])
		if strings.HasPrefix(key, "VSCODE_") || key == "VSCODE_PID" || key == "VSCODE_CWD" || key == "TERM_PROGRAM" && strings.EqualFold(os.Getenv("TERM_PROGRAM"), "vscode") {
			result = true
			break
		}
	}
	return result
}

//goland:noinspection GoUnhandledErrorResult,DuplicatedCode
func prepareElectronEnvironment() (string, error) {
	result := ""
	err := error(nil)
	executable := ""
	if executable, err = os.Executable(); err == nil {
		executableFileName := filepath.Base(executable)
		if debugger.IsPresent() && strings.HasPrefix(strings.ToLower(executableFileName), debugger.DEBUG_BINARY_PREFIX) {
			executableFileName = ELECTRON_NAME
		}
		isEmptyEmbedFiles := false
		if isEmptyEmbedFiles, err = embed.IsEmpty(ELECTRON_DISTRIBUTION_FILES); err == nil && !isEmptyEmbedFiles {
			__debug("Extracting Electron distribution files from embed")
			if err = os.MkdirAll(ELECTRON_DISTRIBUTION_PATH, ELECTRON_UNPACKED_DIRECTORY_PERMISSION); err == nil {
				if _, err = embed.RestoreAll(ELECTRON_DISTRIBUTION_FILES, true); err == nil {
					__debug("Electron distribution files extracted successfully")
				} else {
					__debug(fmt.Sprintf("Failed to restore Electron distribution files from embed: %s", err.Error()))
				}
			} else {
				__debug(fmt.Sprintf("Failed to create directory %s: %s", ELECTRON_DISTRIBUTION_PATH, err.Error()))
			}
		} else if err != nil {
			__debug(fmt.Sprintf("Failed to check embed files: %s", err.Error()))
		}
		appNameElectronExecuteFilePath, _ := searchElectronExecutableFilePath(executableFileName)
		goBoostElectronExecuteFilePath, _ := searchElectronExecutableFilePath(ELECTRON_NAME)
		if appNameElectronExecuteFilePath != "" && goBoostElectronExecuteFilePath != "" {
			if _, isAppNameElectronExists := os.Stat(appNameElectronExecuteFilePath); isAppNameElectronExists == nil {
				if _, isGoBoostElectronExists := os.Stat(goBoostElectronExecuteFilePath); isGoBoostElectronExists == nil {
					appNameElectronHash := ""
					if appNameElectronHash, err = hash.GetFileSHA3(appNameElectronExecuteFilePath); err != nil {
						__debug(fmt.Sprintf("Failed to get hash for %s: %v", appNameElectronExecuteFilePath, err))
					}
					goBoostElectronHash := ""
					if goBoostElectronHash, err = hash.GetFileSHA3(goBoostElectronExecuteFilePath); err != nil {
						__debug(fmt.Sprintf("Failed to get hash for %s: %s", goBoostElectronExecuteFilePath, err))
					}
					if appNameElectronHash == "" || goBoostElectronHash == "" || !strings.EqualFold(appNameElectronHash, goBoostElectronHash) {
						__debug(fmt.Sprintf("Hash mismatch: app=%s, hash=%s", appNameElectronExecuteFilePath, appNameElectronHash))
						__debug(fmt.Sprintf("               app=%s, hash=%s", goBoostElectronExecuteFilePath, goBoostElectronHash))
						process.KillByFilePath(appNameElectronExecuteFilePath)
						fileExtension := filepath.Ext(appNameElectronExecuteFilePath)
						baseName := strings.TrimSuffix(filepath.Base(appNameElectronExecuteFilePath), fileExtension)
						newBaseName := filepath.Join(filepath.Dir(appNameElectronExecuteFilePath), baseName+"_"+uuid.New().String()+fileExtension)
						if err = os.Rename(appNameElectronExecuteFilePath, newBaseName); err == nil {
							__debug(fmt.Sprintf("Renamed %s to %s", appNameElectronExecuteFilePath, newBaseName))
						} else {
							__debug(fmt.Sprintf("Failed to rename %s to %s: %v", appNameElectronExecuteFilePath, newBaseName, err))
						}
					}
				}
			}
		}
		if err == nil {
			if appNameElectronExecuteFilePath, err = searchElectronExecutableFilePath(executableFileName); err == nil {
				if _, err = os.Stat(appNameElectronExecuteFilePath); err == nil {
					result = appNameElectronExecuteFilePath
				} else {
					__debug(fmt.Sprintf("Electron executable not found: %s", err.Error()))
				}
			} else {
				__debug(fmt.Sprintf("Executable not found by name: %s, trying: %s", executableFileName, ELECTRON_NAME))
				if goBoostElectronExecuteFilePath, err = searchElectronExecutableFilePath(ELECTRON_NAME); err == nil {
					appNameElectronExecuteFilePath = filepath.Join(filepath.Dir(goBoostElectronExecuteFilePath), filepath.Base(executable))
					if file.IsExists(appNameElectronExecuteFilePath) {
						__debug(fmt.Sprintf("Renaming existing executable: %s", appNameElectronExecuteFilePath))
						process.KillByFilePath(appNameElectronExecuteFilePath)
						ext := filepath.Ext(appNameElectronExecuteFilePath)
						baseName := strings.TrimSuffix(filepath.Base(appNameElectronExecuteFilePath), ext)
						randomName := filepath.Join(filepath.Dir(appNameElectronExecuteFilePath), baseName+"_"+uuid.New().String()+ext)
						if err = os.Rename(appNameElectronExecuteFilePath, randomName); err == nil {
							__debug(fmt.Sprintf("Renamed %s to %s", appNameElectronExecuteFilePath, randomName))
						} else {
							__debug(fmt.Sprintf("%v", err))
						}
					}
					__debug(fmt.Sprintf("Renaming executable from %s to %s", goBoostElectronExecuteFilePath, appNameElectronExecuteFilePath))
					if err = os.Rename(goBoostElectronExecuteFilePath, appNameElectronExecuteFilePath); err == nil {
						goBoostElectronExecuteFilePath = appNameElectronExecuteFilePath
						if _, err = os.Stat(goBoostElectronExecuteFilePath); err == nil {
							result = goBoostElectronExecuteFilePath
						} else {
							err = fmt.Errorf("cannot found %s(%s)", ELECTRON_NAME, executableFileName)
							__debug(fmt.Sprintf("Electron executable not found: %s", err.Error()))
						}
					}
				} else {
					__debug(fmt.Sprintf(ELECTRON_EXECUTABLE_NOT_FOUND_ERROR, executableFileName, ELECTRON_NAME))
				}
			}
		}
	}
	return result, err
}

func sanitizeJSCode(code string) error {
	err := error(nil)
	lowered := strings.ToLower(code)
	if strings.Contains(lowered, DANGEROUS_SCRIPT_CLOSE_TAG) || strings.Contains(lowered, DANGEROUS_SCRIPT_OPEN_TAG) {
		err = fmt.Errorf("javascript code contains dangerous script tags")
	}
	if err == nil && (strings.ContainsRune(code, '\u2028') || strings.ContainsRune(code, '\u2029')) {
		err = fmt.Errorf("javascript code contains line separator characters")
	}
	if err == nil {
		for _, runeValue := range code {
			if runeValue < MINIMUM_PRINTABLE_RUNE && runeValue != '\t' && runeValue != '\n' && runeValue != '\r' {
				err = fmt.Errorf("javascript code contains illegal control character: %U", runeValue)
				break
			}
		}
	}
	return err
}

func searchElectronExecutableFilePath(fileName string) (string, error) {
	result := ""
	err := error(nil)
	searchDirectoryPath := getUnpackedDirectoryPath()
	__debug(fmt.Sprintf("Searching for executable: %s in: %s", fileName, searchDirectoryPath))
	if err = filepath.Walk(searchDirectoryPath, func(path string, info os.FileInfo, walkError error) error {
		if walkError == nil && !info.IsDir() {
			if strings.EqualFold(info.Name(), fileName) {
				result = path
				__debug(fmt.Sprintf("Executable found: %s", path))
			}
		}
		return walkError
	}); err == nil {
		if result == "" {
			err = fmt.Errorf("executable file not found: %s", fileName)
			__debug(err.Error())
		}
	} else {
		__debug(fmt.Sprintf("Error while walking directory: %s", err.Error()))
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func sendElectronData(data string, dataType string, silent bool) error {
	err := error(nil)
	if webSocketServer == nil {
		err = fmt.Errorf("websocket server not started")
	} else {
		webSocketData := WEB_SOCKET_DATA{
			Type: dataType,
			Data: base64.StdEncoding.EncodeToString([]byte(data)),
		}
		webSocketData.Silent = silent
		jsonData := make([]byte, 0)
		if jsonData, err = json.Marshal(webSocketData); err == nil {
			broadcastMutex.Lock()
			defer broadcastMutex.Unlock()
			webSocketServer.Broadcast(jsonData)
		}
	}
	return err
}

func setElectronProcess(process *exec.Cmd) {
	electronProcessMutex.Lock()
	defer electronProcessMutex.Unlock()
	electronProcess = process
}

func startElectronProcess(electronExecuteFilePath string, params map[string]interface{}) (*exec.Cmd, error) {
	var result *exec.Cmd
	err := error(nil)
	if file.IsExists(electronExecuteFilePath) {
		__debug(fmt.Sprintf("Electron executable found: %s", electronExecuteFilePath))
		__debug(fmt.Sprintf("Electron startup parameters - debugMode: %v", params[PARAMETER_KEY_DEBUG_MODE]))
		var jsonData []byte
		if jsonData, err = json.Marshal(params); err == nil {
			encoded := base64.StdEncoding.EncodeToString(jsonData)
			args := buildChromiumArgs()
			args = append(args, ARGUMENTS_PARAMETERS, encoded)
			result = exec.Command(electronExecuteFilePath, args...)
			if runtime.GOOS == OS_WINDOWS && debugger.IsPresent() && isVSCodeDebugEnvironment() {
				result.Env = buildIsolatedWindowsEnvironment()
			}
			result.Dir = filepath.Dir(electronExecuteFilePath)
			result.Stdout = &STANDARD_CALLBACK{callback: standardOutput}
			result.Stderr = &STANDARD_CALLBACK{callback: standardError}
			if runtime.GOOS == OS_LINUX {
				if fileInfo, statErr := os.Stat(electronExecuteFilePath); statErr == nil {
					if fileInfo.Mode()&0111 == 0 {
						__debug(fmt.Sprintf("Adding execute permission to: %s", electronExecuteFilePath))
						if chmodErr := os.Chmod(electronExecuteFilePath, fileInfo.Mode()|0111); chmodErr != nil {
							__debug(fmt.Sprintf("Failed to add execute permission: %s", chmodErr.Error()))
						}
					}
				}
			}
			if err = result.Start(); err != nil {
				__debug(fmt.Sprintf("Failed to start Electron process: %s", err.Error()))
			}
		}
	} else {
		err = fmt.Errorf("executable file not found: %s", electronExecuteFilePath)
		__debug(err.Error())
	}
	return result, err
}
