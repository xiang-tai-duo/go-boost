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
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
	__websocket "github.com/gorilla/websocket"
	"github.com/xiang-tai-duo/go-bootstrap/debugger"
	"github.com/xiang-tai-duo/go-bootstrap/embed"
	"github.com/xiang-tai-duo/go-bootstrap/file"
	"github.com/xiang-tai-duo/go-bootstrap/logger"
	"github.com/xiang-tai-duo/go-bootstrap/process"
	"github.com/xiang-tai-duo/go-bootstrap/serve"
	websocket "github.com/xiang-tai-duo/go-bootstrap/websocket/server"
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
type (
	ELECTRON struct {
		websocketServer *websocket.WEB_SOCKET_SERVER
		websocketPort   int
		websocketMutex  sync.Mutex
		token           string
		tokenKey        string
		stdout          func(string)
		stderr          func(string)
	}

	ELECTRON_EXIT_PROC func()

	TOKEN_DATA struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	WEBSOCKET_DATA struct {
		Type string      `json:"type"`
		Data interface{} `json:"data"`
	}

	STD_CALLBACK struct {
		f func(string)
	}
)

//goland:noinspection GoSnakeCaseUsage
const (
	ARGS_PARAMS               = "--params"
	PARAM_KEY_DEBUG_MODE      = "isDebugMode"
	PARAM_KEY_HEIGHT          = "height"
	PARAM_KEY_HOME_PAGE       = "homePage"
	PARAM_KEY_WEBSOCKET_PORT  = "websocketPort"
	PARAM_KEY_WIDTH           = "width"
	PARAM_KEY_WINDOW_SIZE     = "windowSize"
	TOKEN_KEY                 = "ELECTRON_TOKEN"
	WEBSOCKET_DATA_TYPE_EVAL  = "eval"
	WEBSOCKET_DATA_TYPE_TOKEN = "token"
)

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func New() *ELECTRON {
	return &ELECTRON{
		stdout: func(s string) { fmt.Print(s) },
		stderr: func(s string) { fmt.Fprint(os.Stderr, s) },
	}
}

func (electron *ELECTRON) SetStdout(stdout func(string)) {
	electron.stdout = stdout
}

func (electron *ELECTRON) SetStderr(stderr func(string)) {
	electron.stderr = stderr
}

//goland:noinspection GoUnhandledErrorResult
func (electron *ELECTRON) Launch(homePage string) (*exec.Cmd, error) {
	var result *exec.Cmd
	err := error(nil)
	process.KillByName(ELECTRON_NAME)
	electron.websocketMutex.Lock()
	defer electron.websocketMutex.Unlock()
	if electron.websocketServer == nil {
		electron.websocketServer = websocket.NewWebSocketServer()
		electron.token = uuid.New().String()
		electron.tokenKey = TOKEN_KEY
		serve.OnBeforeStaticServe(func(r *http.Request, w http.ResponseWriter) bool {
			result := false
			if debugger.IsPresent() {
				result = true
			} else {
				if cookie, err := r.Cookie(electron.tokenKey); err == nil {
					if cookie.Value == electron.token {
						result = true
					}
				}
			}
			if !result {
				http.Error(w, "", http.StatusForbidden)
				logger.Logger.Warning(fmt.Sprintf("Access Denied, %s -> %s", r.RemoteAddr, r.URL.String()))
			}
			return result
		})
		electron.websocketServer.SetConnectHandler(func(conn *__websocket.Conn, uuid string) error {
			data := WEBSOCKET_DATA{
				Type: WEBSOCKET_DATA_TYPE_TOKEN,
				Data: TOKEN_DATA{
					Key:   electron.tokenKey,
					Value: electron.token,
				},
			}
			if jsonData, err := json.Marshal(data); err == nil {
				electron.websocketServer.Broadcast(jsonData)
			}
			return nil
		})
		if electron.websocketPort, err = electron.websocketServer.LaunchAsync(websocket.ANY_PORT); err == nil {
			if isEmpty, _ := embed.IsEmpty(ELECTRON_DIST_FILES); !isEmpty {
				if err = os.MkdirAll(ELECTRON_DIST_PATH, 0755); err == nil {
					embed.RestoreAll(ELECTRON_DIST_FILES, true)
				}
			}
			if err == nil {
				executable := ""
				if executable, err = os.Executable(); err == nil {
					exeFileName := filepath.Base(executable)
					electronExecuteFilePath := ""
					if electronExecuteFilePath, err = electron.findExecutable(exeFileName); err != nil {
						if electronExecuteFilePath, err = electron.findExecutable(ELECTRON_NAME); err == nil {
							newElectronFilePath := filepath.Join(filepath.Dir(electronExecuteFilePath), filepath.Base(executable))
							if file.IsExists(newElectronFilePath) {
								os.Remove(newElectronFilePath)
							}
							if err = os.Rename(electronExecuteFilePath, newElectronFilePath); err == nil {
								electronExecuteFilePath = newElectronFilePath
							}
						}
					}
					if _, err = os.Stat(electronExecuteFilePath); err == nil {
						if result, err = startElectronCmd(electronExecuteFilePath, homePage, electron.websocketPort, electron.stdout, electron.stderr); err == nil {
							err = result.Wait()
						}
					} else {
						err = errors.New(fmt.Sprintf("cannot found %s(%s)", ELECTRON_NAME, exeFileName))
					}
				}
			}
		}
	}
	return result, err
}

//goland:noinspection SpellCheckingInspection,GoBoolExpressions
func startElectronCmd(electronExecuteFilePath string, homePage string, websocketPort int, stdout func(string), stderr func(string)) (*exec.Cmd, error) {
	var result *exec.Cmd
	err := error(nil)
	if file.IsExists(electronExecuteFilePath) {
		params := map[string]interface{}{
			PARAM_KEY_DEBUG_MODE:     debugger.IsPresent(),
			PARAM_KEY_HOME_PAGE:      homePage,
			PARAM_KEY_WEBSOCKET_PORT: websocketPort,
			PARAM_KEY_WINDOW_SIZE: map[string]int{
				PARAM_KEY_WIDTH:  -1,
				PARAM_KEY_HEIGHT: -1,
			},
		}
		jsonData, _ := json.Marshal(params)
		encoded := base64.StdEncoding.EncodeToString(jsonData)
		args := []string{ARGS_PARAMS, encoded}
		result = exec.Command(electronExecuteFilePath, args...)
		result.Dir = filepath.Dir(electronExecuteFilePath)
		result.Stdout = &STD_CALLBACK{f: stdout}
		result.Stderr = &STD_CALLBACK{f: stderr}
		err = result.Start()
	} else {
		err = fmt.Errorf("executable file not found: %s", electronExecuteFilePath)
	}
	return result, err
}

func (electron *ELECTRON) findExecutable(fileName string) (string, error) {
	result := ""
	err := error(nil)
	err = filepath.Walk(electron.getUnpackedDirectoryPath(), func(path string, info os.FileInfo, walkErr error) error {
		err := walkErr
		if walkErr == nil && !info.IsDir() {
			if strings.Contains(info.Name(), fileName) {
				result = path
			}
		}
		return err
	})
	if err == nil && result == "" {
		err = fmt.Errorf("executable file not found: %s", fileName)
	}
	return result, err
}

func (electron *ELECTRON) getUnpackedDirectoryPath() string {
	result := ELECTRON_DIST_PATH
	if wordingDirectoryPath, err := os.Getwd(); err == nil {
		result = filepath.Join(wordingDirectoryPath, ELECTRON_DIST_PATH)
	}
	return result
}

func (electron *ELECTRON) Eval(code string) error {
	err := error(nil)
	if electron.websocketServer == nil {
		err = fmt.Errorf("websocket server not started")
	} else {
		data := WEBSOCKET_DATA{
			Type: WEBSOCKET_DATA_TYPE_EVAL,
			Data: base64.StdEncoding.EncodeToString([]byte(code)),
		}
		if jsonData, err := json.Marshal(data); err == nil {
			electron.websocketServer.Broadcast(jsonData)
		}
	}
	return err
}

func (w *STD_CALLBACK) Write(p []byte) (n int, err error) {
	w.f(string(p))
	return len(p), nil
}
