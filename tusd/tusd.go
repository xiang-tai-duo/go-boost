// Package tusd
// File:        tusd.go
// Url:         https://github.com/xiang-tai-duo/go-bootstrap/blob/master/tusd/tusd.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: TUSD is a wrapper for the tus resumable upload protocol, providing HTTP handler and file storage capabilities.
// --------------------------------------------------------------------------------
package tusd

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/tus/tusd/v2/pkg/filestore"
	"github.com/tus/tusd/v2/pkg/handler"
)

type (
	UPLOADED_HANDLER func(id string, filePath string, metaData map[string]string)
)

const (
	DEFAULT_ADDRESS       = ":1080"
	DEFAULT_BASE_PATH     = "/files/"
	DEFAULT_STORE_PATH    = "./data"
	DEFAULT_IDLE_TIMEOUT  = 60 * time.Second
	DEFAULT_READ_TIMEOUT  = 15 * time.Second
	DEFAULT_WRITE_TIMEOUT = 15 * time.Second
	WAIT_CLIENT_TIMEOUT   = 5 * time.Second
	WAIT_RETRY_INTERVAL   = 10 * time.Millisecond
)

var (
	address          string
	basePath         string
	storePath        string
	maxSize          int64
	behindProxy      bool
	disableDownload  bool
	enableMoveFile   bool
	targetDirectory  string
	uploadedHandlers []UPLOADED_HANDLER
	errorHandler     func(error)
	mutex            sync.Mutex
	server           *http.Server
	serverCancel     context.CancelFunc
	serverContext    context.Context
	serverListener   net.Listener
)

//goland:noinspection GoUnusedExportedFunction
func init() {
	serverContext, serverCancel = context.WithCancel(context.Background())
	server = &http.Server{
		Addr:         DEFAULT_ADDRESS,
		ReadTimeout:  DEFAULT_READ_TIMEOUT,
		WriteTimeout: DEFAULT_WRITE_TIMEOUT,
		IdleTimeout:  DEFAULT_IDLE_TIMEOUT,
	}
	address = DEFAULT_ADDRESS
	basePath = DEFAULT_BASE_PATH
	storePath = DEFAULT_STORE_PATH
	maxSize = 0
	behindProxy = false
	disableDownload = false
	enableMoveFile = true
	targetDirectory = ""
}

//goland:noinspection GoUnusedExportedFunction
func ListenAsync() error {
	err := error(nil)
	mutex.Lock()
	if serverListener == nil {
		store := filestore.New(storePath)
		composer := handler.NewStoreComposer()
		store.UseIn(composer)
		handlerConfig := handler.Config{
			StoreComposer:           composer,
			BasePath:                basePath,
			MaxSize:                 maxSize,
			RespectForwardedHeaders: behindProxy,
			DisableDownload:         disableDownload,
			NotifyCompleteUploads:   true,
		}
		if tusdHandler, handlerErr := handler.NewHandler(handlerConfig); handlerErr == nil {
			go func() {
				for event := range tusdHandler.CompleteUploads {
					if enableMoveFile {
						sourcePath := filepath.Join(storePath, event.Upload.ID)
						targetDir := storePath
						if targetDirectory != "" {
							targetDir = targetDirectory
						}
						fileName := event.Upload.ID
						if name, ok := event.Upload.MetaData["filename"]; ok {
							fileName = name
						}
						targetPath := filepath.Join(targetDir, fileName)
						if _, err := os.Stat(targetPath); err == nil {
							ext := filepath.Ext(fileName)
							base := fileName[:len(fileName)-len(ext)]
							counter := 1
							for {
								newFileName := fmt.Sprintf("%s_%d%s", base, counter, ext)
								targetPath = filepath.Join(targetDir, newFileName)
								if _, err := os.Stat(targetPath); os.IsNotExist(err) {
									break
								}
								counter++
							}
						}
						if err := os.Rename(sourcePath, targetPath); err == nil {
							infoPath := filepath.Join(storePath, event.Upload.ID+".info")
							_ = os.Remove(infoPath)
							for _, handler := range uploadedHandlers {
								handler(event.Upload.ID, targetPath, event.Upload.MetaData)
							}
						} else {
							for _, handler := range uploadedHandlers {
								handler(event.Upload.ID, sourcePath, event.Upload.MetaData)
							}
						}
					} else {
						sourcePath := filepath.Join(storePath, event.Upload.ID)
						for _, handler := range uploadedHandlers {
							handler(event.Upload.ID, sourcePath, event.Upload.MetaData)
						}
					}
				}
			}()
			server.Handler = tusdHandler
			go func() {
				if listener, err := net.Listen("tcp", server.Addr); err == nil {
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
		} else {
			err = handlerErr
		}
	} else {
		err = fmt.Errorf("TUSD server is already running")
	}
	mutex.Unlock()
	waitServe()
	return err
}

//goland:noinspection GoUnusedExportedFunction
func OnUploaded(handler UPLOADED_HANDLER) {
	mutex.Lock()
	defer mutex.Unlock()
	uploadedHandlers = append(uploadedHandlers, handler)
}

//goland:noinspection GoUnusedExportedFunction
func OnError(handler func(error)) {
	mutex.Lock()
	defer mutex.Unlock()
	errorHandler = handler
}

//goland:noinspection GoUnusedExportedFunction
func SetAddress(addr string) {
	mutex.Lock()
	defer mutex.Unlock()
	address = addr
	server.Addr = addr
}

//goland:noinspection GoUnusedExportedFunction
func GetAddress() string {
	mutex.Lock()
	defer mutex.Unlock()
	return address
}

//goland:noinspection GoUnusedExportedFunction
func SetBasePath(path string) {
	mutex.Lock()
	defer mutex.Unlock()
	basePath = path
}

//goland:noinspection GoUnusedExportedFunction
func GetBasePath() string {
	mutex.Lock()
	defer mutex.Unlock()
	return basePath
}

//goland:noinspection GoUnusedExportedFunction
func SetStorePath(path string) {
	mutex.Lock()
	defer mutex.Unlock()
	storePath = path
}

//goland:noinspection GoUnusedExportedFunction
func GetStorePath() string {
	mutex.Lock()
	defer mutex.Unlock()
	return storePath
}

//goland:noinspection GoUnusedExportedFunction
func SetMaxSize(size int64) {
	mutex.Lock()
	defer mutex.Unlock()
	maxSize = size
}

//goland:noinspection GoUnusedExportedFunction
func GetMaxSize() int64 {
	mutex.Lock()
	defer mutex.Unlock()
	return maxSize
}

//goland:noinspection GoUnusedExportedFunction
func SetBehindProxy(enabled bool) {
	mutex.Lock()
	defer mutex.Unlock()
	behindProxy = enabled
}

//goland:noinspection GoUnusedExportedFunction
func GetBehindProxy() bool {
	mutex.Lock()
	defer mutex.Unlock()
	return behindProxy
}

//goland:noinspection GoUnusedExportedFunction
func SetDisableDownload(enabled bool) {
	mutex.Lock()
	defer mutex.Unlock()
	disableDownload = enabled
}

//goland:noinspection GoUnusedExportedFunction
func GetDisableDownload() bool {
	mutex.Lock()
	defer mutex.Unlock()
	return disableDownload
}

//goland:noinspection GoUnusedExportedFunction
func SetEnableMoveFile(enabled bool) {
	mutex.Lock()
	defer mutex.Unlock()
	enableMoveFile = enabled
}

//goland:noinspection GoUnusedExportedFunction
func GetEnableMoveFile() bool {
	mutex.Lock()
	defer mutex.Unlock()
	return enableMoveFile
}

//goland:noinspection GoUnusedExportedFunction
func SetTargetDirectory(directory string) {
	mutex.Lock()
	defer mutex.Unlock()
	targetDirectory = directory
}

//goland:noinspection GoUnusedExportedFunction
func GetTargetDirectory() string {
	mutex.Lock()
	defer mutex.Unlock()
	return targetDirectory
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
func IsListening() bool {
	mutex.Lock()
	defer mutex.Unlock()
	return serverListener != nil
}

//goland:noinspection GoUnusedExportedFunction
func Shutdown() {
	mutex.Lock()
	if serverListener != nil {
		_ = server.Shutdown(serverContext)
		serverListener = nil
	}
	mutex.Unlock()
}

func waitServe() {
	client := &http.Client{
		Timeout: WAIT_CLIENT_TIMEOUT,
	}
	url := fmt.Sprintf("http://%s%s", server.Addr, basePath)
	for {
		if resp, err := client.Get(url); err == nil {
			_ = resp.Body.Close()
			break
		}
		time.Sleep(WAIT_RETRY_INTERVAL)
	}
}
