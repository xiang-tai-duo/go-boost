// Package mqttclient
// File:        client.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/mqtt/client/client.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: MQTT client provides functionality for MQTT client communication
// --------------------------------------------------------------------------------
package mqttclient

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/xiang-tai-duo/go-boost/ca"
	"github.com/xiang-tai-duo/go-boost/logger"
	mqttcommon "github.com/xiang-tai-duo/go-boost/mqtt/common"
)

type (
	MQTT struct {
		broker            string
		client            paho.Client
		clientIdentifier  string
		connectHandler    func()
		connected         bool
		connectionTimeout time.Duration
		disconnectHandler func()
		keepAlive         time.Duration
		messageHandler    func(*mqttcommon.MQTT_MESSAGE)
		mutex             sync.Mutex
		password          string
		QoS               int
		subscriptions     map[string]int
		username          string
		state             ConnectionState
	}
)

//goland:noinspection GoSnakeCaseUsage,SpellCheckingInspection
const (
	CLIENT_PREFIX                     = "go-boost-"
	DEFAULT_CONNECTION_TIMEOUT        = 30 * time.Second
	DEFAULT_KEEP_ALIVE                = 60 * time.Second
	QUALITY_OF_SERVICE_PARAMETER      = 0
	QUALITY_OF_SERVICE_PARAMETER_SIZE = 1
	RETAINED_PARAMETER                = 1
	RETAINED_PARAMETER_SIZE           = 2
	MAXIMUM_QUALITY_OF_SERVICE        = 2
	MINIMUM_QUALITY_OF_SERVICE        = 0
	RANDOM_IDENTIFIER_BYTES           = 8
	TRANSMISSION_CONTROL_SCHEME       = "tcp://"
	SECURITY_SCHEME                   = "ssl://"
	SECURE_MQTT_SCHEME                = "mqtts://"
	SECURE_TRANSPORT_SCHEME           = "tls://"
	MODULE_NAME_CLIENT                = "mqtt.client"
)

type ConnectionState int

const (
	ConnectionStateDisconnected ConnectionState = iota
	ConnectionStateConnecting
	ConnectionStateConnected
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_CLIENT, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_CLIENT, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_CLIENT, logger.SKIP_STACK_FRAMES_BASE)
}

func New(broker string) *MQTT {
	result := (*MQTT)(nil)
	result = NewWithUserAndPassword(broker, "", "")
	return result
}

func (clientConnection *MQTT) Connect() error {
	err := error(nil)
	clientConnection.mutex.Lock()
	clientConnection.state = ConnectionStateConnecting
	clientConnection.mutex.Unlock()
	if clientConnection.client == nil || !clientConnection.client.IsConnected() {
		__debug(fmt.Sprintf("[MQTT-Client] Connect requested: clientID=%s, broker=%s, username=%s", clientConnection.clientIdentifier, clientConnection.broker, clientConnection.username))
		if err = clientConnection.connect(clientConnection.broker, nil); err != nil {
			__debug(fmt.Sprintf("[MQTT-Client] Plain MQTT connect failed: broker=%s, clientID=%s, error=%v", clientConnection.broker, clientConnection.clientIdentifier, err))
			transportLayerSecurityConfiguration, hasClientCertificate := buildClientTLSConfig()
			if hasClientCertificate && transportLayerSecurityConfiguration != nil {
				brokerAddress := upgradeBrokerToTLS(clientConnection.broker)
				if brokerAddress != clientConnection.broker {
					__debug(fmt.Sprintf("[MQTT-Client][TLS] Broker URL upgraded to TLS: %s -> %s", clientConnection.broker, brokerAddress))
				}
				err = clientConnection.connect(brokerAddress, transportLayerSecurityConfiguration)
			} else {
				__debug("[MQTT-Client][TLS] TLS retry skipped (no client certificate found)")
			}
		}
	}
	return err
}

func (clientConnection *MQTT) Disconnect(timeoutDuration time.Duration) error {
	err := error(nil)
	if clientConnection.client != nil && clientConnection.client.IsConnected() {
		__debug(fmt.Sprintf("[MQTT-Client] Disconnect requested: clientID=%s, timeout=%s", clientConnection.clientIdentifier, timeoutDuration))
		clientConnection.client.Disconnect(uint(timeoutDuration.Milliseconds()))
		clientConnection.mutex.Lock()
		clientConnection.connected = false
		clientConnection.state = ConnectionStateDisconnected
		clientConnection.mutex.Unlock()
		if clientConnection.disconnectHandler != nil {
			clientConnection.disconnectHandler()
		}
		__debug(fmt.Sprintf("[MQTT-Client] Disconnected: clientID=%s", clientConnection.clientIdentifier))
	}
	return err
}

func (clientConnection *MQTT) GetBroker() string {
	result := ""
	result = clientConnection.broker
	return result
}

func (clientConnection *MQTT) GetClientID() string {
	result := ""
	result = clientConnection.clientIdentifier
	return result
}

func (clientConnection *MQTT) GetConnectionTimeout() time.Duration {
	result := time.Duration(0)
	result = clientConnection.connectionTimeout
	return result
}

func (clientConnection *MQTT) GetKeepAlive() time.Duration {
	result := time.Duration(0)
	result = clientConnection.keepAlive
	return result
}

func (clientConnection *MQTT) GetPassword() string {
	result := ""
	result = clientConnection.password
	return result
}

func (clientConnection *MQTT) GetQualityOfService() int {
	result := MINIMUM_QUALITY_OF_SERVICE
	result = clientConnection.QoS
	return result
}

func (clientConnection *MQTT) GetSubscriptions() map[string]int {
	result := make(map[string]int)
	clientConnection.mutex.Lock()
	defer clientConnection.mutex.Unlock()
	for topic, qualityOfService := range clientConnection.subscriptions {
		result[topic] = qualityOfService
	}
	return result
}

func (clientConnection *MQTT) GetUsername() string {
	result := ""
	result = clientConnection.username
	return result
}

func (clientConnection *MQTT) IsConnected() bool {
	result := false
	clientConnection.mutex.Lock()
	defer clientConnection.mutex.Unlock()
	if clientConnection.client == nil {
	} else if clientConnection.connected && clientConnection.client.IsConnected() {
		result = true
	}
	return result
}

func (clientConnection *MQTT) IsConnecting() bool {
	result := false
	clientConnection.mutex.Lock()
	defer clientConnection.mutex.Unlock()
	result = clientConnection.state == ConnectionStateConnecting
	return result
}

//goland:noinspection GoUnusedExportedFunction
func NewWithUser(broker string, user string) *MQTT {
	result := (*MQTT)(nil)
	result = NewWithUserAndPassword(broker, user, "")
	return result
}

func NewWithUserAndPassword(broker string, user string, pass string) *MQTT {
	result := (*MQTT)(nil)
	result = &MQTT{
		broker:            broker,
		clientIdentifier:  CLIENT_PREFIX + generateRandomID(),
		connectionTimeout: DEFAULT_CONNECTION_TIMEOUT,
		keepAlive:         DEFAULT_KEEP_ALIVE,
		password:          pass,
		QoS:               mqttcommon.DEFAULT_MQTT_QUALITY_OF_SERVICE,
		subscriptions:     make(map[string]int),
		username:          user,
	}
	return result
}

func (clientConnection *MQTT) Publish(topic string, payload string, parameters ...interface{}) error {
	err := error(nil)
	qualityOfService := mqttcommon.DEFAULT_MQTT_QUALITY_OF_SERVICE
	retained := false
	if len(parameters) >= QUALITY_OF_SERVICE_PARAMETER_SIZE {
		qualityOfServiceValue := MINIMUM_QUALITY_OF_SERVICE
		typeMatched := false
		qualityOfServiceValue, typeMatched = parameters[QUALITY_OF_SERVICE_PARAMETER].(int)
		if typeMatched {
			qualityOfService = qualityOfServiceValue
		}
	}
	if len(parameters) >= RETAINED_PARAMETER_SIZE {
		retainedValue := false
		typeMatched := false
		retainedValue, typeMatched = parameters[RETAINED_PARAMETER].(bool)
		if typeMatched {
			retained = retainedValue
		}
	}
	if !clientConnection.IsConnected() {
		err = fmt.Errorf("client is not connected")
		__debug(fmt.Sprintf("[MQTT-Client] Publish failed: topic=%s, error=%v", topic, err))
	} else if topic == "" {
		err = fmt.Errorf("topic cannot be empty")
		__debug(fmt.Sprintf("[MQTT-Client] Publish failed: %v", err))
	} else {
		__debug(fmt.Sprintf("[MQTT-Client] Publish topic=%s, qos=%d, retained=%t, payloadLen=%d", topic, qualityOfService, retained, len(payload)))
		publishToken := clientConnection.client.Publish(topic, byte(qualityOfService), retained, payload)
		publishToken.WaitTimeout(clientConnection.connectionTimeout)
		if err = publishToken.Error(); err == nil {
		} else {
			__debug(fmt.Sprintf("[MQTT-Client] Publish error on topic=%s: %v", topic, err))
		}
	}
	return err
}

func (clientConnection *MQTT) RemoveMessageHandler() error {
	err := error(nil)
	clientConnection.messageHandler = nil
	return err
}

func (clientConnection *MQTT) SetBroker(brokerAddress string) error {
	err := error(nil)
	if brokerAddress != "" {
		clientConnection.broker = brokerAddress
	} else {
		err = fmt.Errorf("broker address cannot be empty")
	}
	return err
}

func (clientConnection *MQTT) SetClientID(clientIdentifier string) error {
	err := error(nil)
	if clientIdentifier != "" {
		clientConnection.clientIdentifier = clientIdentifier
	} else {
		err = fmt.Errorf("client ID cannot be empty")
	}
	return err
}

func (clientConnection *MQTT) SetConnectHandler(handler func()) error {
	err := error(nil)
	clientConnection.connectHandler = handler
	return err
}

func (clientConnection *MQTT) SetConnectionTimeout(timeoutDuration time.Duration) error {
	err := error(nil)
	if timeoutDuration > 0 {
		clientConnection.connectionTimeout = timeoutDuration
	} else {
		err = fmt.Errorf("connection timeout must be greater than 0")
	}
	return err
}

func (clientConnection *MQTT) SetDisconnectHandler(handler func()) error {
	err := error(nil)
	clientConnection.disconnectHandler = handler
	return err
}

func (clientConnection *MQTT) SetKeepAlive(keepAliveDuration time.Duration) error {
	err := error(nil)
	if keepAliveDuration > 0 {
		clientConnection.keepAlive = keepAliveDuration
	} else {
		err = fmt.Errorf("keep alive time must be greater than 0")
	}
	return err
}

func (clientConnection *MQTT) SetMessageHandler(handler func(*mqttcommon.MQTT_MESSAGE)) error {
	err := error(nil)
	clientConnection.messageHandler = handler
	return err
}

func (clientConnection *MQTT) SetPassword(passwordValue string) error {
	err := error(nil)
	clientConnection.password = passwordValue
	return err
}

func (clientConnection *MQTT) SetQualityOfService(qualityOfServiceLevel int) error {
	err := error(nil)
	if qualityOfServiceLevel >= MINIMUM_QUALITY_OF_SERVICE && qualityOfServiceLevel <= MAXIMUM_QUALITY_OF_SERVICE {
		clientConnection.QoS = qualityOfServiceLevel
	} else {
		err = fmt.Errorf("quality of Service level must be 0, 1, or 2")
	}
	return err
}

func (clientConnection *MQTT) SetUsername(usernameValue string) error {
	err := error(nil)
	clientConnection.username = usernameValue
	return err
}

func (clientConnection *MQTT) Subscribe(topic string, qualityOfService int) error {
	err := error(nil)
	if !clientConnection.IsConnected() {
		err = fmt.Errorf("client is not connected")
		__debug(fmt.Sprintf("[MQTT-Client] Subscribe failed: topic=%s, error=%v", topic, err))
	} else if topic == "" {
		err = fmt.Errorf("topic cannot be empty")
		__debug(fmt.Sprintf("[MQTT-Client] Subscribe failed: %v", err))
	} else if qualityOfService < MINIMUM_QUALITY_OF_SERVICE || qualityOfService > MAXIMUM_QUALITY_OF_SERVICE {
		err = fmt.Errorf("quality of Service level must be 0, 1, or 2")
		__debug(fmt.Sprintf("[MQTT-Client] Subscribe failed for topic=%s: %v", topic, err))
	} else {
		__debug(fmt.Sprintf("[MQTT-Client] Subscribe topic=%s, qos=%d", topic, qualityOfService))
		subscriptionToken := clientConnection.client.Subscribe(topic, byte(qualityOfService), clientConnection.internalMessageHandler)
		subscriptionToken.WaitTimeout(clientConnection.connectionTimeout)
		if err = subscriptionToken.Error(); err == nil {
			clientConnection.mutex.Lock()
			clientConnection.subscriptions[topic] = qualityOfService
			clientConnection.mutex.Unlock()
			__debug(fmt.Sprintf("[MQTT-Client] Subscribed topic=%s, qos=%d", topic, qualityOfService))
		} else {
			__debug(fmt.Sprintf("[MQTT-Client] Subscribe error on topic=%s: %v", topic, err))
		}
	}
	return err
}

func (clientConnection *MQTT) Unsubscribe(topic string) error {
	err := error(nil)
	if !clientConnection.IsConnected() {
		err = fmt.Errorf("client is not connected")
		__debug(fmt.Sprintf("[MQTT-Client] Unsubscribe failed: topic=%s, error=%v", topic, err))
	} else if topic == "" {
		err = fmt.Errorf("topic cannot be empty")
		__debug(fmt.Sprintf("[MQTT-Client] Unsubscribe failed: %v", err))
	} else {
		__debug(fmt.Sprintf("[MQTT-Client] Unsubscribe topic=%s", topic))
		unsubscribeToken := clientConnection.client.Unsubscribe(topic)
		unsubscribeToken.WaitTimeout(clientConnection.connectionTimeout)
		if err = unsubscribeToken.Error(); err == nil {
			clientConnection.mutex.Lock()
			delete(clientConnection.subscriptions, topic)
			clientConnection.mutex.Unlock()
			__debug(fmt.Sprintf("[MQTT-Client] Unsubscribed topic=%s", topic))
		} else {
			__debug(fmt.Sprintf("[MQTT-Client] Unsubscribe error on topic=%s: %v", topic, err))
		}
	}
	return err
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_CLIENT, logger.SKIP_STACK_FRAMES_BASE)
}

func buildClientTLSConfig() (*tls.Config, bool) {
	result := (*tls.Config)(nil)
	hasClientCertificate := false
	err := error(nil)
	__debug("[MQTT-Client][TLS] Attempting to build TLS config")
	certificatePath := findMqttClientCertificateFile(ca.CLIENT_CERTIFICATE_FILE_NAME)
	privateKeyPath := findMqttClientCertificateFile(ca.CLIENT_PRIVATE_KEY_FILE_NAME)
	if certificatePath == "" || privateKeyPath == "" {
		__debug(fmt.Sprintf("[MQTT-Client][TLS] Client certificate or key not found, TLS not enabled (cert='%s', key='%s')", certificatePath, privateKeyPath))
	} else {
		certificate := tls.Certificate{}
		if certificate, err = tls.LoadX509KeyPair(certificatePath, privateKeyPath); err == nil {
			__debug(fmt.Sprintf("[MQTT-Client][TLS] Client certificate loaded from %s", certificatePath))
			hasClientCertificate = true
			result = &tls.Config{
				Certificates: []tls.Certificate{certificate},
				MinVersion:   tls.VersionTLS12,
			}
			certificateAuthorityPool := loadMqttClientCAPool()
			if certificateAuthorityPool == nil {
				result.InsecureSkipVerify = true
				__debug("[MQTT-Client][TLS] CA certificate not found, server certificate verification disabled (insecure)")
			} else {
				result.RootCAs = certificateAuthorityPool
				__debug("[MQTT-Client][TLS] mTLS enabled: will verify server certificate against loaded CA")
			}
		} else {
			__debug(fmt.Sprintf("[MQTT-Client][TLS] Failed to load client key pair (cert='%s', key='%s'): %v", certificatePath, privateKeyPath, err))
		}
	}
	return result, hasClientCertificate
}

func (clientConnection *MQTT) connect(brokerAddress string, transportLayerSecurityConfiguration *tls.Config) error {
	err := error(nil)
	clientOptions := paho.NewClientOptions()
	if transportLayerSecurityConfiguration == nil {
		__info(fmt.Sprintf("[MQTT-Client] Using plain MQTT connection: broker=%s", brokerAddress))
	} else {
		clientOptions.SetTLSConfig(transportLayerSecurityConfiguration)
		__info(fmt.Sprintf("[MQTT-Client] Using TLS encrypted MQTT connection: broker=%s, mTLS=%t", brokerAddress, len(transportLayerSecurityConfiguration.Certificates) > 0 && transportLayerSecurityConfiguration.RootCAs != nil))
		__debug(fmt.Sprintf("[MQTT-Client][TLS] TLS config applied (mTLS-ready=%t, rootCAsLoaded=%t)", len(transportLayerSecurityConfiguration.Certificates) > 0, transportLayerSecurityConfiguration.RootCAs != nil))
	}
	clientOptions.AddBroker(brokerAddress)
	clientOptions.SetClientID(clientConnection.clientIdentifier)
	clientOptions.SetUsername(clientConnection.username)
	clientOptions.SetPassword(clientConnection.password)
	clientOptions.SetKeepAlive(clientConnection.keepAlive)
	clientOptions.SetPingTimeout(clientConnection.connectionTimeout)
	clientOptions.SetOnConnectHandler(func(client paho.Client) {
		err := error(nil)
		__debug(fmt.Sprintf("[MQTT-Client] Connected to broker %s as %s", brokerAddress, clientConnection.clientIdentifier))
		clientConnection.mutex.Lock()
		clientConnection.connected = true
		clientConnection.state = ConnectionStateConnected
		clientConnection.mutex.Unlock()
		if clientConnection.connectHandler != nil {
			clientConnection.connectHandler()
		}
		clientConnection.mutex.Lock()
		for topic, qualityOfService := range clientConnection.subscriptions {
			__debug(fmt.Sprintf("[MQTT-Client] Re-subscribing to topic=%s, qos=%d", topic, qualityOfService))
			subscriptionToken := client.Subscribe(topic, byte(qualityOfService), clientConnection.internalMessageHandler)
			subscriptionToken.WaitTimeout(clientConnection.connectionTimeout)
			if err = subscriptionToken.Error(); err == nil {
			} else {
				__debug(fmt.Sprintf("[MQTT-Client] Re-subscribe failed for topic=%s: %v", topic, err))
			}
		}
		clientConnection.mutex.Unlock()
	})
	clientOptions.SetConnectionLostHandler(func(_ paho.Client, reason error) {
		__debug(fmt.Sprintf("[MQTT-Client] Connection lost from broker %s: %v", brokerAddress, reason))
		clientConnection.mutex.Lock()
		clientConnection.connected = false
		clientConnection.state = ConnectionStateDisconnected
		clientConnection.mutex.Unlock()
		if clientConnection.disconnectHandler != nil {
			clientConnection.disconnectHandler()
		}
	})
	clientOptions.SetDefaultPublishHandler(clientConnection.internalMessageHandler)
	clientConnection.client = paho.NewClient(clientOptions)
	connectToken := clientConnection.client.Connect()
	connectToken.WaitTimeout(clientConnection.connectionTimeout)
	clientConnection.mutex.Lock()
	if err = connectToken.Error(); err == nil {
		clientConnection.state = ConnectionStateConnected
		__debug(fmt.Sprintf("[MQTT-Client] Connect token completed successfully: broker=%s", brokerAddress))
	} else {
		clientConnection.state = ConnectionStateDisconnected
		__debug(fmt.Sprintf("[MQTT-Client] Connect failed: broker=%s, clientID=%s, error=%v", brokerAddress, clientConnection.clientIdentifier, err))
		clientConnection.client = nil
	}
	clientConnection.mutex.Unlock()
	return err
}

func findMqttClientCertificateFile(fileName string) string {
	result := ""
	err := error(nil)
	__debug(fmt.Sprintf("[MQTT-Client][Certificate] Looking for certificate file: %s", fileName))
	executablePath := ""
	if executablePath, err = os.Executable(); err == nil {
		executableDirectory := filepath.Dir(executablePath)
		candidatePath := filepath.Join(executableDirectory, fileName)
		__debug(fmt.Sprintf("[MQTT-Client][Certificate] Checking executable directory: %s", candidatePath))
		if _, err = os.Stat(candidatePath); err == nil {
			result = candidatePath
			__debug(fmt.Sprintf("[MQTT-Client][Certificate] Found in executable directory: %s", candidatePath))
		} else {
			__debug(fmt.Sprintf("[MQTT-Client][Certificate] Not found in executable directory: %s, error: %v", candidatePath, err))
		}
	} else {
		__debug(fmt.Sprintf("[MQTT-Client][Certificate] Failed to get executable path: %v", err))
	}
	if result == "" {
		err = nil
		workingDirectory := ""
		if workingDirectory, err = os.Getwd(); err == nil {
			candidatePath := filepath.Join(workingDirectory, fileName)
			__debug(fmt.Sprintf("[MQTT-Client][Certificate] Checking working directory: %s", candidatePath))
			if _, err = os.Stat(candidatePath); err == nil {
				result = candidatePath
				__debug(fmt.Sprintf("[MQTT-Client][Certificate] Found in working directory: %s", candidatePath))
			} else {
				__debug(fmt.Sprintf("[MQTT-Client][Certificate] Not found in working directory: %s, error: %v", candidatePath, err))
			}
		} else {
			__debug(fmt.Sprintf("[MQTT-Client][Certificate] Failed to get working directory: %v", err))
		}
	}
	if result == "" {
		__debug(fmt.Sprintf("[MQTT-Client][Certificate] File not found: %s", fileName))
	}
	return result
}

func generateRandomID() string {
	result := ""
	err := error(nil)
	randomBytes := make([]byte, RANDOM_IDENTIFIER_BYTES)
	if _, err = rand.Read(randomBytes); err == nil {
		result = hex.EncodeToString(randomBytes)
	} else {
		result = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return result
}

func (clientConnection *MQTT) internalMessageHandler(_ paho.Client, message paho.Message) {
	if clientConnection.messageHandler != nil {
		receivedMessage := &mqttcommon.MQTT_MESSAGE{
			Topic:            message.Topic(),
			Payload:          string(message.Payload()),
			Timestamp:        time.Now(),
			QualityOfService: int(message.Qos()),
			Retained:         message.Retained(),
			Duplicate:        message.Duplicate(),
		}
		clientConnection.messageHandler(receivedMessage)
	}
}

func loadMqttClientCAPool() *x509.CertPool {
	result := (*x509.CertPool)(nil)
	err := error(nil)
	__debug("[MQTT-Client][TLS] Loading CA certificate pool")
	certificatePath := findMqttClientCertificateFile(ca.CERTIFICATE_AUTHORITY_CERTIFICATE_FILE_NAME)
	if certificatePath == "" {
		__debug("[MQTT-Client][TLS] CA certificate file not found")
	} else {
		certificateBytes := make([]byte, 0)
		if certificateBytes, err = os.ReadFile(certificatePath); err == nil {
			certificateAuthorityPool := x509.NewCertPool()
			if certificateAuthorityPool.AppendCertsFromPEM(certificateBytes) {
				__debug(fmt.Sprintf("[MQTT-Client][TLS] CA certificate pool loaded from %s", certificatePath))
				result = certificateAuthorityPool
			} else {
				__debug(fmt.Sprintf("[MQTT-Client][TLS] Failed to parse CA certificate from %s", certificatePath))
			}
		} else {
			__debug(fmt.Sprintf("[MQTT-Client][TLS] Failed to read CA certificate file %s: %v", certificatePath, err))
		}
	}
	return result
}

func upgradeBrokerToTLS(brokerAddress string) string {
	result := ""
	result = brokerAddress
	lowerBrokerAddress := strings.ToLower(brokerAddress)
	if strings.HasPrefix(lowerBrokerAddress, SECURITY_SCHEME) ||
		strings.HasPrefix(lowerBrokerAddress, SECURE_TRANSPORT_SCHEME) ||
		strings.HasPrefix(lowerBrokerAddress, SECURE_MQTT_SCHEME) {
	} else if strings.HasPrefix(lowerBrokerAddress, TRANSMISSION_CONTROL_SCHEME) {
		result = SECURITY_SCHEME + brokerAddress[len(TRANSMISSION_CONTROL_SCHEME):]
	}
	return result
}
