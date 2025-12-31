// Package boost
// File:        mqtt.go
// Author:      TRAE AI
// Created:     2025/12/30 11:03:46
// Description: MQTT client wrapper for go-boost library
package boost

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	DEFAULT_MQTT_BROKER             = "tcp://localhost:1883"
	DEFAULT_MQTT_QOS                = 0
	DEFAULT_MQTT_KEEP_ALIVE         = 60 * time.Second
	DEFAULT_MQTT_CONNECTION_TIMEOUT = 30 * time.Second
)

type (
	MQTT_MESSAGE struct {
		Topic     string    `json:"topic"`
		Payload   string    `json:"payload"`
		Timestamp time.Time `json:"timestamp"`
		QoS       int       `json:"qos"`
		Retained  bool      `json:"retained"`
		Duplicate bool      `json:"duplicate"`
	}

	MQTT struct {
		client            mqtt.Client
		broker            string
		clientID          string
		username          string
		password          string
		qos               int
		keepAlive         time.Duration
		connectionTimeout time.Duration
		messageHandler    func(*MQTT_MESSAGE)
		connectHandler    func()
		disconnectHandler func()
		lock              sync.Mutex
		connected         bool
		subscriptions     map[string]int
	}
)

func NewMQTT(brokerAddress string, params ...string) *MQTT {
	clientID := "go-boost-" + generateRandomID()
	user := ""
	pass := ""
	if len(params) > 0 {
		user = params[0]
		if len(params) > 1 {
			pass = params[1]
		}
	}
	return &MQTT{
		broker:            brokerAddress,
		clientID:          clientID,
		username:          user,
		password:          pass,
		qos:               DEFAULT_MQTT_QOS,
		keepAlive:         DEFAULT_MQTT_KEEP_ALIVE,
		connectionTimeout: DEFAULT_MQTT_CONNECTION_TIMEOUT,
		subscriptions:     make(map[string]int),
	}
}

func (m *MQTT) Connect() error {
	var err error
	if m.client != nil && m.client.IsConnected() {
		err = nil
	} else {
		opts := mqtt.NewClientOptions()
		opts.AddBroker(m.broker)
		opts.SetClientID(m.clientID)
		opts.SetUsername(m.username)
		opts.SetPassword(m.password)
		opts.SetKeepAlive(m.keepAlive)
		opts.SetPingTimeout(m.connectionTimeout)
		opts.SetOnConnectHandler(func(client mqtt.Client) {
			m.lock.Lock()
			m.connected = true
			m.lock.Unlock()
			if m.connectHandler != nil {
				m.connectHandler()
			}
			m.lock.Lock()
			for topic, qos := range m.subscriptions {
				client.Subscribe(topic, byte(qos), m.internalMessageHandler)
			}
			m.lock.Unlock()
		})
		opts.SetConnectionLostHandler(func(client mqtt.Client, reason error) {
			m.lock.Lock()
			m.connected = false
			m.lock.Unlock()
		})
		opts.SetDefaultPublishHandler(m.internalMessageHandler)
		m.client = mqtt.NewClient(opts)
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

func (m *MQTT) GetQoS() int {
	return m.qos
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

func (m *MQTT) Publish(topic string, payload string, qos int, retained bool) error {
	var err error
	if !m.IsConnected() {
		err = fmt.Errorf("client is not connected")
	} else if topic == "" {
		err = fmt.Errorf("topic cannot be empty")
	} else {
		token := m.client.Publish(topic, byte(qos), retained, payload)
		token.WaitTimeout(m.connectionTimeout)
		err = token.Error()
	}
	return err
}

func (m *MQTT) SetBroker(brokerAddress string) error {
	var err error
	if brokerAddress == "" {
		err = fmt.Errorf("broker address cannot be empty")
	} else {
		m.broker = brokerAddress
		err = nil
	}
	return err
}

func (m *MQTT) SetClientID(clientIdentifier string) error {
	var err error
	if clientIdentifier == "" {
		err = fmt.Errorf("client ID cannot be empty")
	} else {
		m.clientID = clientIdentifier
		err = nil
	}
	return err
}

func (m *MQTT) SetConnectHandler(handler func()) error {
	m.connectHandler = handler
	return nil
}

func (m *MQTT) SetConnectionTimeout(timeoutDuration time.Duration) error {
	var err error
	if timeoutDuration <= 0 {
		err = fmt.Errorf("connection timeout must be greater than 0")
	} else {
		m.connectionTimeout = timeoutDuration
		err = nil
	}
	return err
}

func (m *MQTT) SetDisconnectHandler(handler func()) error {
	m.disconnectHandler = handler
	return nil
}

func (m *MQTT) SetKeepAlive(keepAliveDuration time.Duration) error {
	var err error
	if keepAliveDuration <= 0 {
		err = fmt.Errorf("keep alive time must be greater than 0")
	} else {
		m.keepAlive = keepAliveDuration
		err = nil
	}
	return err
}

func (m *MQTT) SetMessageHandler(handler func(*MQTT_MESSAGE)) error {
	m.messageHandler = handler
	return nil
}

func (m *MQTT) SetPassword(passwordValue string) error {
	m.password = passwordValue
	return nil
}

func (m *MQTT) SetQoS(qosLevel int) error {
	var err error
	if qosLevel < 0 || qosLevel > 2 {
		err = fmt.Errorf("QoS level must be 0, 1, or 2")
	} else {
		m.qos = qosLevel
		err = nil
	}
	return err
}

func (m *MQTT) SetUsername(usernameValue string) error {
	m.username = usernameValue
	return nil
}

func (m *MQTT) Subscribe(topic string, qos int) error {
	var err error
	if !m.IsConnected() {
		err = fmt.Errorf("client is not connected")
	} else if topic == "" {
		err = fmt.Errorf("topic cannot be empty")
	} else if qos < 0 || qos > 2 {
		err = fmt.Errorf("QoS level must be 0, 1, or 2")
	} else {
		token := m.client.Subscribe(topic, byte(qos), m.internalMessageHandler)
		token.WaitTimeout(m.connectionTimeout)
		if token.Error() != nil {
			err = token.Error()
		} else {
			m.lock.Lock()
			m.subscriptions[topic] = qos
			m.lock.Unlock()
			err = nil
		}
	}
	return err
}

func (m *MQTT) Unsubscribe(topic string) error {
	var err error
	if !m.IsConnected() {
		err = fmt.Errorf("client is not connected")
	} else if topic == "" {
		err = fmt.Errorf("topic cannot be empty")
	} else {
		token := m.client.Unsubscribe(topic)
		token.WaitTimeout(m.connectionTimeout)
		if token.Error() != nil {
			err = token.Error()
		} else {
			m.lock.Lock()
			delete(m.subscriptions, topic)
			m.lock.Unlock()
			err = nil
		}
	}
	return err
}

func generateRandomID() string {
	var id string
	b := make([]byte, 8)
	if _, err := rand.Read(b); err == nil {
		id = hex.EncodeToString(b)
	} else {
		id = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return id
}

func (m *MQTT) internalMessageHandler(client mqtt.Client, msg mqtt.Message) {
	if m.messageHandler != nil {
		timestamp := time.Now()
		mqttMessage := &MQTT_MESSAGE{
			Topic:     string(msg.Topic()),
			Payload:   string(msg.Payload()),
			Timestamp: timestamp,
			QoS:       int(msg.Qos()),
			Retained:  msg.Retained(),
			Duplicate: msg.Duplicate(),
		}
		m.messageHandler(mqttMessage)
	}
}
