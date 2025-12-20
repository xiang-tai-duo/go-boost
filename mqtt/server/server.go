// Package mqttserver
// File:        server.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/mqtt/server/server.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: MQTT server functionality
// --------------------------------------------------------------------------------
package mqttserver

import (
	"crypto/tls"
	"fmt"
	"sync"

	server "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/listeners"
)

//goland:noinspection GoSnakeCaseUsage
const (
	DEFAULT_MQTT_SERVER_PORT = 1883
	DEFAULT_MQTT_SERVER_HOST = "0.0.0.0"
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
type (
	MqttServer struct {
		Server    *server.Server
		Host      string
		Port      int
		CertFile  string
		KeyFile   string
		TLSConfig *tls.Config
		isRunning bool
		lock      sync.Mutex
	}
)

//goland:noinspection GoUnusedExportedFunction
func New(params ...interface{}) *MqttServer {
	host := DEFAULT_MQTT_SERVER_HOST
	port := DEFAULT_MQTT_SERVER_PORT
	switch len(params) {
	case 1:
		if h, ok := params[0].(string); ok {
			host = h
		} else if p, ok := params[0].(int); ok {
			port = p
		}
	case 2:
		if h, ok := params[0].(string); ok {
			host = h
		}
		if p, ok := params[1].(int); ok {
			port = p
		}
	}

	return &MqttServer{
		Host:     host,
		Port:     port,
		CertFile: "",
		KeyFile:  "",
		Server:   server.New(nil),
	}
}

func (ms *MqttServer) GetHost() string {
	return ms.Host
}

func (ms *MqttServer) GetPort() int {
	return ms.Port
}

func (ms *MqttServer) IsRunning() bool {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	return ms.isRunning
}

func (ms *MqttServer) Start() error {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	err := error(nil)
	if !ms.isRunning && ms.Server != nil {
		address := fmt.Sprintf("%s:%d", ms.Host, ms.Port)
		tcp := listeners.NewTCP(listeners.Config{
			ID:      "tcp",
			Address: address,
		})
		if err = ms.Server.AddListener(tcp); err == nil {
			go func() {
				if serveErr := ms.Server.Serve(); serveErr != nil {
					// Handle error if needed
				}
			}()
			ms.isRunning = true
		}
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult
func (ms *MqttServer) Stop() error {
	result := error(nil)
	ms.lock.Lock()
	defer ms.lock.Unlock()
	if ms.isRunning && ms.Server != nil {
		ms.Server.Close()
		ms.isRunning = false
	}
	return result
}
