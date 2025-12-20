// Package mqttclient
// File:        mqtt.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/mqtt/client/client.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: MQTT client implementation
// --------------------------------------------------------------------------------
package mqttclient

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	mqttcommon "github.com/xiang-tai-duo/go-boost/mqtt/common"
)

//goland:noinspection GoUnusedConst,GoSnakeCaseUsage
const (
	DEFAULT_MQTT_BROKER             = "tcp://localhost:1883"
	DEFAULT_MQTT_KEEP_ALIVE         = 60 * time.Second
	DEFAULT_MQTT_CONNECTION_TIMEOUT = 30 * time.Second
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
type (
	MQTT struct {
		client            paho.Client
		broker            string
		clientID          string
		username          string
		password          string
		qualityOfService  int
		keepAlive         time.Duration
		connectionTimeout time.Duration
		messageHandler    func(*mqttcommon.MQTT_MESSAGE)
		connectHandler    func()
		disconnectHandler func()
		lock              sync.Mutex
		connected         bool
		subscriptions     map[string]int
	}
)

//goland:noinspection GoUnusedExportedFunction
func New(broker string, params ...interface{}) *MQTT {
	user := ""
	pass := ""
	qualityOfService := mqttcommon.DEFAULT_MQTT_QOS
	keepAlive := DEFAULT_MQTT_KEEP_ALIVE
	connectionTimeout := DEFAULT_MQTT_CONNECTION_TIMEOUT

	// Process params - only handle username and password
	switch len(params) {
	case 1:
		if u, ok := params[0].(string); ok {
			user = u
		}
	case 2:
		if u, ok := params[0].(string); ok {
			user = u
		}
		if p, ok := params[1].(string); ok {
			pass = p
		}
	}

	clientID := "go-boost-" + generateRandomID()
	return &MQTT{
		broker:            broker,
		clientID:          clientID,
		username:          user,
		password:          pass,
		qualityOfService:  qualityOfService,
		keepAlive:         keepAlive,
		connectionTimeout: connectionTimeout,
		subscriptions:     make(map[string]int),
	}
}

func (m *MQTT) Connect() error {
	err := error(nil)
	if m.client != nil && m.client.IsConnected() {
		err = nil
	} else {
		opts := paho.NewClientOptions()
		opts.AddBroker(m.broker)
		opts.SetClientID(m.clientID)
		opts.SetUsername(m.username)
		opts.SetPassword(m.password)
		opts.SetKeepAlive(m.keepAlive)
		opts.SetPingTimeout(m.connectionTimeout)
		opts.SetOnConnectHandler(func(client paho.Client) {
			m.lock.Lock()
			m.connected = true
			m.lock.Unlock()
			if m.connectHandler != nil {
				m.connectHandler()
			}
			m.lock.Lock()
			for topic, qualityOfService := range m.subscriptions {
				client.Subscribe(topic, byte(qualityOfService), m.internalMessageHandler)
			}
			m.lock.Unlock()
		})
		opts.SetConnectionLostHandler(func(_ paho.Client, reason error) {
			m.lock.Lock()
			m.connected = false
			m.lock.Unlock()
		})
		opts.SetDefaultPublishHandler(m.internalMessageHandler)
		m.client = paho.NewClient(opts)
		token := m.client.Connect()
		token.WaitTimeout(m.connectionTimeout)
		err = token.Error()
	}
	return err
}

func (m *MQTT) Disconnect(timeoutDuration time.Duration) error {
	if m.client != nil && m.client.IsConnected() {
		m.client.Disconnect(uint(timeoutDuration.Milliseconds()))
		m.lock.Lock()
		m.connected = false
		m.lock.Unlock()
		if m.disconnectHandler != nil {
			m.disconnectHandler()
		}
	}
	return nil
}

func (m *MQTT) GetBroker() string {
	return m.broker
}

func (m *MQTT) GetClientID() string {
	return m.clientID
}

func (m *MQTT) GetConnectionTimeout() time.Duration {
	return m.connectionTimeout
}

func (m *MQTT) GetKeepAlive() time.Duration {
	return m.keepAlive
}

func (m *MQTT) GetPassword() string {
	return m.password
}

func (m *MQTT) GetQualityOfService() int {
	return m.qualityOfService
}

func (m *MQTT) GetSubscriptions() map[string]int {
	m.lock.Lock()
	defer func() {
		m.lock.Unlock()
	}()
	subscriptionsCopy := make(map[string]int)
	for topic, qos := range m.subscriptions {
		subscriptionsCopy[topic] = qos
	}
	return subscriptionsCopy
}

func (m *MQTT) GetUsername() string {
	return m.username
}

func (m *MQTT) IsConnected() bool {
	m.lock.Lock()
	defer func() {
		m.lock.Unlock()
	}()
	return m.connected && m.client != nil && m.client.IsConnected()
}

func (m *MQTT) Publish(topic string, payload string, params ...interface{}) error {
	err := error(nil)
	qualityOfService := mqttcommon.DEFAULT_MQTT_QOS
	retained := false

	if len(params) > 0 {
		if q, ok := params[0].(int); ok {
			qualityOfService = q
		}
	}
	if len(params) > 1 {
		if r, ok := params[1].(bool); ok {
			retained = r
		}
	}

	if m.IsConnected() && topic != "" {
		token := m.client.Publish(topic, byte(qualityOfService), retained, payload)
		token.WaitTimeout(m.connectionTimeout)
		err = token.Error()
	} else if !m.IsConnected() {
		err = fmt.Errorf("client is not connected")
	} else if topic == "" {
		err = fmt.Errorf("topic cannot be empty")
	}
	return err
}

func (m *MQTT) SetBroker(brokerAddress string) error {
	err := error(nil)
	if brokerAddress != "" {
		m.broker = brokerAddress
		err = nil
	} else {
		err = fmt.Errorf("broker address cannot be empty")
	}
	return err
}

func (m *MQTT) SetClientID(clientIdentifier string) error {
	err := error(nil)
	if clientIdentifier != "" {
		m.clientID = clientIdentifier
		err = nil
	} else {
		err = fmt.Errorf("client ID cannot be empty")
	}
	return err
}

func (m *MQTT) SetConnectHandler(handler func()) error {
	m.connectHandler = handler
	return nil
}

func (m *MQTT) SetConnectionTimeout(timeoutDuration time.Duration) error {
	err := error(nil)
	if timeoutDuration > 0 {
		m.connectionTimeout = timeoutDuration
		err = nil
	} else {
		err = fmt.Errorf("connection timeout must be greater than 0")
	}
	return err
}

func (m *MQTT) SetDisconnectHandler(handler func()) error {
	m.disconnectHandler = handler
	return nil
}

func (m *MQTT) SetKeepAlive(keepAliveDuration time.Duration) error {
	err := error(nil)
	if keepAliveDuration > 0 {
		m.keepAlive = keepAliveDuration
		err = nil
	} else {
		err = fmt.Errorf("keep alive time must be greater than 0")
	}
	return err
}

func (m *MQTT) SetMessageHandler(handler func(*mqttcommon.MQTT_MESSAGE)) error {
	m.messageHandler = handler
	return nil
}

func (m *MQTT) SetPassword(passwordValue string) error {
	m.password = passwordValue
	return nil
}

func (m *MQTT) SetQualityOfService(qualityOfServiceLevel int) error {
	err := error(nil)
	if qualityOfServiceLevel >= 0 && qualityOfServiceLevel <= 2 {
		m.qualityOfService = qualityOfServiceLevel
		err = nil
	} else {
		err = fmt.Errorf("quality of Service level must be 0, 1, or 2")
	}
	return err
}

func (m *MQTT) SetUsername(usernameValue string) error {
	m.username = usernameValue
	return nil
}

func (m *MQTT) Subscribe(topic string, qualityOfService int) error {
	err := error(nil)
	if m.IsConnected() && topic != "" && qualityOfService >= 0 && qualityOfService <= 2 {
		token := m.client.Subscribe(topic, byte(qualityOfService), m.internalMessageHandler)
		token.WaitTimeout(m.connectionTimeout)
		if err = token.Error(); err == nil {
			m.lock.Lock()
			m.subscriptions[topic] = qualityOfService
			m.lock.Unlock()
		}
	} else if !m.IsConnected() {
		err = fmt.Errorf("client is not connected")
	} else if topic == "" {
		err = fmt.Errorf("topic cannot be empty")
	} else if qualityOfService < 0 || qualityOfService > 2 {
		err = fmt.Errorf("quality of Service level must be 0, 1, or 2")
	}
	return err
}

func (m *MQTT) Unsubscribe(topic string) error {
	err := error(nil)
	if m.IsConnected() && topic != "" {
		token := m.client.Unsubscribe(topic)
		token.WaitTimeout(m.connectionTimeout)
		if err = token.Error(); err == nil {
			m.lock.Lock()
			delete(m.subscriptions, topic)
			m.lock.Unlock()
		}
	} else if !m.IsConnected() {
		err = fmt.Errorf("client is not connected")
	} else if topic == "" {
		err = fmt.Errorf("topic cannot be empty")
	}
	return err
}

func (m *MQTT) internalMessageHandler(_ paho.Client, msg paho.Message) {
	if m.messageHandler != nil {
		timestamp := time.Now()
		mqttMessage := &mqttcommon.MQTT_MESSAGE{
			Topic:            msg.Topic(),
			Payload:          string(msg.Payload()),
			Timestamp:        timestamp,
			QualityOfService: int(msg.Qos()),
			Retained:         msg.Retained(),
			Duplicate:        msg.Duplicate(),
		}
		m.messageHandler(mqttMessage)
	}
}

func generateRandomID() string {
	id := ""
	b := make([]byte, 8)
	if _, err := rand.Read(b); err == nil {
		id = hex.EncodeToString(b)
	} else {
		id = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return id
}
