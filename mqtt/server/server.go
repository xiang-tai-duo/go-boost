// Package mqttserver
// File:        server.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/mqtt/server/server.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: MQTT server provides functionality for MQTT server operations
// --------------------------------------------------------------------------------
package mqttserver

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	server "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/xiang-tai-duo/go-boost/ca"
	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection GoSnakeCaseUsage
type (
	AUTH_HOOK struct {
		server.HookBase
		ledger *auth.Ledger
	}
	HOOK struct {
		server.HookBase
		messageHandler    MESSAGE_HANDLER
		subscriberHandler SUBSCRIBE_HANDLER
	}
	MESSAGE_HANDLER   func(client_identification, topic string, payload []byte)
	SUBSCRIBE_HANDLER func(client_identification, topic string, quality_of_service byte)
)

//goland:noinspection GoSnakeCaseUsage,SpellCheckingInspection
const (
	ADDRESS_FORMAT            = "%s:%d"
	AUTHENTICATION_HOOK_ID    = "go-boost-auth-hook"
	DEFAULT_MESSAGE_QUALITY   = byte(0)
	DEFAULT_MQTT_SERVER_HOST  = "0.0.0.0"
	DEFAULT_MQTT_SERVER_PORT  = 1883
	DEFAULT_MQTTS_SERVER_PORT = 8883
	DEFAULT_PUBLISH_RETAINED  = false
	HOOK_ID                   = "go-boost-hook"
	INLINE_ID_FORMAT          = "inline-%d"
	MODULE_NAME_SERVER        = "mqtt.server"
	TCP_LISTENER_ID           = "tcp"
	TLS_TRANSPORT_CONTROL_ID  = "tls"
)

//goland:noinspection GoSnakeCaseUsage,SpellCheckingInspection
var (
	allowInlineMqttClient  bool
	anonymousMode          bool
	authenticationPassword string
	authenticationUsername string
	hostAddress            string
	isRunning              bool
	messageHandler         MESSAGE_HANDLER
	mqttServer             *server.Server
	mutexProtection        sync.Mutex
	mqttsPortNumber        int
	portNumber             int
	subscriberHandler      SUBSCRIBE_HANDLER
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_SERVER, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_SERVER, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_SERVER, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_SERVER, logger.SKIP_STACK_FRAMES_BASE)
}

func New() {
	NewWithHostAndPort(DEFAULT_MQTT_SERVER_HOST, DEFAULT_MQTT_SERVER_PORT)
}

//goland:noinspection GoUnhandledErrorResult
func Close() error {
	result := error(nil)
	mutexProtection.Lock()
	defer mutexProtection.Unlock()
	if isRunning && mqttServer != nil {
		__debug("Closing MQTT server")
		if close_err := mqttServer.Close(); close_err != nil {
			__debug(fmt.Sprintf("Error closing MQTT server: %v", close_err))
			result = close_err
		}
		isRunning = false
		__debug("MQTT server closed")
	}
	return result
}

func DisconnectClientByID(clientID string) error {
	result := error(nil)
	if clientID != "" && mqttServer != nil && mqttServer.Clients != nil {
		if client, ok := mqttServer.Clients.Get(clientID); ok && client != nil {
			result = mqttServer.DisconnectClient(client, packets.ErrSessionTakenOver)
		}
	}
	return result
}

func GetHost() string {
	result := hostAddress
	return result
}

func GetPort() int {
	result := portNumber
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetSubscribedTopics() []string {
	result := make([]string, 0)
	if mqttServer != nil && mqttServer.Clients != nil {
		topic_set := make(map[string]struct{})
		for _, client := range mqttServer.Clients.GetAll() {
			if client != nil && client.State.Subscriptions != nil {
				for filter := range client.State.Subscriptions.GetAll() {
					topic_set[filter] = struct{}{}
				}
			}
		}
		for filter := range topic_set {
			result = append(result, filter)
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetSubscribers(topic string) []string {
	result := make([]string, 0)
	if mqttServer != nil && mqttServer.Topics != nil {
		if subscribers := mqttServer.Topics.Subscribers(topic); subscribers != nil {
			for client_identification := range subscribers.Subscriptions {
				result = append(result, client_identification)
			}
			for _, group := range subscribers.Shared {
				for client_identification := range group {
					result = append(result, client_identification)
				}
			}
			for id := range subscribers.InlineSubscriptions {
				result = append(result, fmt.Sprintf(INLINE_ID_FORMAT, id))
			}
		}
	}
	return result
}

func GetTLSPort() int {
	result := mqttsPortNumber
	return result
}

//goland:noinspection GoUnusedExportedFunction
func HasSubscribers(topic string) bool {
	result := false
	subscribers := GetSubscribers(topic)
	if len(subscribers) > 0 {
		result = true
	}
	return result
}

func (hook *AUTH_HOOK) ID() string {
	result := AUTHENTICATION_HOOK_ID
	return result
}

func (h *HOOK) ID() string {
	result := HOOK_ID
	return result
}

//goland:noinspection GoUnusedExportedFunction
func IsAnonymous() bool {
	result := anonymousMode
	return result
}

func IsClientConnected(clientID string) bool {
	result := false
	if clientID != "" && mqttServer != nil && mqttServer.Clients != nil {
		if client, ok := mqttServer.Clients.Get(clientID); ok && client != nil && !client.Closed() {
			result = true
		}
	}
	return result
}

func IsRunning() bool {
	result := false
	mutexProtection.Lock()
	defer mutexProtection.Unlock()
	if isRunning {
		result = true
	}
	return result
}

//goland:noinspection GoUnhandledErrorResult
func ListenAsync() error {
	result := error(nil)
	mutexProtection.Lock()
	defer mutexProtection.Unlock()
	if !isRunning && mqttServer != nil {
		__debug(fmt.Sprintf("Starting MQTT server, host=%s, port=%d, anonymous=%t", hostAddress, portNumber, anonymousMode))
		if anonymousMode {
			__debug("Adding AllowHook (anonymous mode)")
			if result = mqttServer.AddHook(new(auth.AllowHook), nil); result != nil {
				__debug(fmt.Sprintf("Failed to add AllowHook: %v", result))
			}
		} else if result == nil {
			__debug(fmt.Sprintf("Adding AuthHook with username=%s", authenticationUsername))
			ledger := &auth.Ledger{
				Users: auth.Users{
					authenticationUsername: auth.UserRule{
						Username: auth.RString(authenticationUsername),
						Password: auth.RString(authenticationPassword),
						ACL: auth.Filters{
							auth.RString("#"): auth.ReadWrite,
						},
					},
				},
			}
			if result = mqttServer.AddHook(&AUTH_HOOK{ledger: ledger}, nil); result != nil {
				__debug(fmt.Sprintf("Failed to add AuthHook: %v", result))
			}
		}
		if result == nil {
			hook := &HOOK{
				messageHandler:    messageHandler,
				subscriberHandler: subscriberHandler,
			}
			if hook_err := mqttServer.AddHook(hook, nil); hook_err != nil {
				__debug(fmt.Sprintf("Failed to add message/subscribe hook: %v", hook_err))
			}
			address := fmt.Sprintf(ADDRESS_FORMAT, hostAddress, portNumber)
			listener_config := listeners.Config{
				ID:      TCP_LISTENER_ID,
				Address: address,
			}
			__debug(fmt.Sprintf("TLS disabled, listening plain TCP on %s", address))
			tcp := listeners.NewTCP(listener_config)
			if result = mqttServer.AddListener(tcp); result == nil {
				__debug(fmt.Sprintf("Listener %s added at %s", TCP_LISTENER_ID, address))
				tls_config := buildServerTLSConfig()
				if tls_config != nil {
					tls_address := fmt.Sprintf(ADDRESS_FORMAT, hostAddress, mqttsPortNumber)
					tls_listener_config := listeners.Config{
						ID:        TLS_TRANSPORT_CONTROL_ID,
						Address:   tls_address,
						TLSConfig: tls_config,
					}
					tls_tcp := listeners.NewTCP(tls_listener_config)
					if result = mqttServer.AddListener(tls_tcp); result == nil {
						__debug(fmt.Sprintf("Listener %s added at %s (mTLS=%t)", TLS_TRANSPORT_CONTROL_ID, tls_address, tls_config.ClientAuth == tls.RequireAndVerifyClientCert))
					} else {
						__debug(fmt.Sprintf("Failed to add TLS listener on %s: %v", tls_address, result))
					}
				}
				if result == nil {
					go func() {
						if serve_err := mqttServer.Serve(); serve_err != nil {
							__debug(fmt.Sprintf("mqtt server error: %v", serve_err))
							slog.Error("mqtt server error", "error", serve_err)
						}
					}()
					isRunning = true
					__debug("Server started successfully")
				}
			} else {
				__debug(fmt.Sprintf("Failed to add listener on %s: %v", address, result))
			}
		}
	}
	return result
}

func NewWithHost(h string) {
	NewWithHostAndPort(h, DEFAULT_MQTT_SERVER_PORT)
}

func NewWithHostAndPort(h string, p int) {
	NewWithHostAndPorts(h, p, DEFAULT_MQTTS_SERVER_PORT)
}

func NewWithHostAndPorts(h string, p int, tls_port int) {
	hostAddress = h
	portNumber = p
	mqttsPortNumber = tls_port
	mqttServer = server.New(&server.Options{InlineClient: allowInlineMqttClient})
}

func NewWithPort(p int) {
	NewWithHostAndPort(DEFAULT_MQTT_SERVER_HOST, p)
}

func NewWithPorts(p int, tls_port int) {
	NewWithHostAndPorts(DEFAULT_MQTT_SERVER_HOST, p, tls_port)
}

func (hook *AUTH_HOOK) OnACLCheck(client *server.Client, topic string, write bool) bool {
	result := true
	if hook.ledger != nil {
		_, result = hook.ledger.ACLOk(client, topic, write)
	}
	return result
}

func (hook *AUTH_HOOK) OnConnectAuthenticate(client *server.Client, packet packets.Packet) bool {
	result := false
	if hook.ledger != nil {
		_, result = hook.ledger.AuthOk(client, packet)
	}
	if !result {
		__debug(fmt.Sprintf("[MQTT-Server][Auth] Authentication failed: clientID=%s, remote=%s, username=%s", client.ID, client.Net.Remote, string(packet.Connect.Username)))
		slog.Warn("mqtt client authentication failed", "clientID", client.ID, "remote", client.Net.Remote, "username", string(packet.Connect.Username), "password", string(packet.Connect.Password))
	} else {
		__debug(fmt.Sprintf("[MQTT-Server][Auth] Authentication succeeded: clientID=%s, remote=%s, username=%s", client.ID, client.Net.Remote, string(packet.Connect.Username)))
	}
	return result
}

func (h *HOOK) OnPublished(cl *server.Client, pk packets.Packet) {
	if h.messageHandler != nil && cl != nil {
		h.messageHandler(cl.ID, pk.TopicName, pk.Payload)
	}
}

func (h *HOOK) OnSubscribed(cl *server.Client, pk packets.Packet, reason_codes []byte) {
	if h.subscriberHandler != nil {
		for i, sub := range pk.Filters {
			if i < len(reason_codes) {
				h.subscriberHandler(cl.ID, sub.Filter, reason_codes[i])
			}
		}
	}
}

func (hook *AUTH_HOOK) Provides(method byte) bool {
	result := false
	if method == server.OnConnectAuthenticate || method == server.OnACLCheck {
		result = true
	}
	return result
}

func (h *HOOK) Provides(b byte) bool {
	result := false
	if b == server.OnSubscribed || b == server.OnPublished {
		result = true
	}
	return result
}

func Publish(topic string, payload string) error {
	result := error(nil)
	result = PublishEx(topic, []byte(payload), DEFAULT_MESSAGE_QUALITY, DEFAULT_PUBLISH_RETAINED)
	return result
}

func PublishEx(topic string, payload []byte, qualityOfService byte, retained bool) error {
	result := error(nil)
	if mqttServer == nil {
		result = fmt.Errorf("mqtt server not initialized")
		__debug(fmt.Sprintf("PublishEx failed: %v", result))
	} else {
		__debug(fmt.Sprintf("PublishEx topic=%s, qos=%d, retained=%t, payloadLen=%d", topic, qualityOfService, retained, len(payload)))
		if result = mqttServer.Publish(topic, payload, retained, qualityOfService); result != nil {
			__debug(fmt.Sprintf("PublishEx error on topic=%s: %v", topic, result))
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func SetAllowInlineClient(enabled bool) {
	allowInlineMqttClient = enabled
}

//goland:noinspection GoUnusedExportedFunction
func SetAnonymous(enabled bool) {
	anonymousMode = enabled
}

//goland:noinspection GoUnusedExportedFunction
func SetAuthentication(username string, password string) {
	authenticationUsername = username
	authenticationPassword = password
}

func SetOnMessageHandler(handler func(client_identification, topic string, payload []byte)) {
	messageHandler = handler
}

func SetOnSubscribeHandler(handler func(client_identification, topic string, quality_of_service byte)) {
	subscriberHandler = handler
}

func buildServerTLSConfig() *tls.Config {
	var result *tls.Config
	__debug("Attempting to build TLS config")
	certificate_path := findMqttCertificateFile(ca.SERVER_CERTIFICATE_FILE_NAME)
	private_key_path := findMqttCertificateFile(ca.SERVER_PRIVATE_KEY_FILE_NAME)
	if certificate_path == "" || private_key_path == "" {
		__debug(fmt.Sprintf("Server certificate or key not found, TLS not enabled (cert='%s', key='%s')", certificate_path, private_key_path))
	} else {
		var certificate tls.Certificate
		var err error
		if certificate, err = tls.LoadX509KeyPair(certificate_path, private_key_path); err == nil {
			__debug(fmt.Sprintf("Server certificate loaded from %s", certificate_path))
			result = &tls.Config{
				Certificates: []tls.Certificate{certificate},
				MinVersion:   tls.VersionTLS12,
			}
			if ca_pool := loadMqttCAPool(); ca_pool != nil {
				result.ClientCAs = ca_pool
				result.ClientAuth = tls.RequireAndVerifyClientCert
				__debug("mTLS enabled: clients must present a certificate signed by the loaded CA")
			} else {
				result.ClientAuth = tls.NoClientCert
				__debug("CA certificate not found, mTLS disabled (server-only TLS)")
			}
		} else {
			__debug(fmt.Sprintf("Failed to load server key pair (cert='%s', key='%s'): %v", certificate_path, private_key_path, err))
		}
	}
	return result
}

func findMqttCertificateFile(file_name string) string {
	result := ""
	__debug(fmt.Sprintf("Looking for certificate file: %s", file_name))
	var executable_path string
	var err error
	if executable_path, err = os.Executable(); err == nil {
		executable_directory := filepath.Dir(executable_path)
		candidate := filepath.Join(executable_directory, file_name)
		__debug(fmt.Sprintf("Checking executable directory: %s", candidate))
		if _, stat_err := os.Stat(candidate); stat_err == nil {
			result = candidate
			__debug(fmt.Sprintf("Found in executable directory: %s", candidate))
		} else {
			__debug(fmt.Sprintf("Not found in executable directory: %s, error: %v", candidate, stat_err))
		}
	} else {
		__debug(fmt.Sprintf("Failed to get executable path: %v", err))
	}
	if result == "" {
		var working_directory string
		if working_directory, err = os.Getwd(); err == nil {
			candidate := filepath.Join(working_directory, file_name)
			__debug(fmt.Sprintf("Checking working directory: %s", candidate))
			if _, stat_err := os.Stat(candidate); stat_err == nil {
				result = candidate
				__debug(fmt.Sprintf("Found in working directory: %s", candidate))
			} else {
				__debug(fmt.Sprintf("Not found in working directory: %s, error: %v", candidate, stat_err))
			}
		} else {
			__debug(fmt.Sprintf("Failed to get working directory: %v", err))
		}
	}
	if result == "" {
		__debug(fmt.Sprintf("File not found: %s", file_name))
	}
	return result
}

func loadMqttCAPool() *x509.CertPool {
	var result *x509.CertPool
	__debug("Loading CA certificate pool")
	certificate_path := findMqttCertificateFile(ca.CERTIFICATE_AUTHORITY_CERTIFICATE_FILE_NAME)
	if certificate_path != "" {
		var bytes []byte
		var err error
		if bytes, err = os.ReadFile(certificate_path); err == nil {
			pool := x509.NewCertPool()
			if pool.AppendCertsFromPEM(bytes) {
				__debug(fmt.Sprintf("CA certificate pool loaded from %s", certificate_path))
				result = pool
			} else {
				__debug(fmt.Sprintf("Failed to parse CA certificate from %s", certificate_path))
			}
		} else {
			__debug(fmt.Sprintf("Failed to read CA certificate file %s: %v", certificate_path, err))
		}
	} else {
		__debug("CA certificate file not found")
	}
	return result
}
