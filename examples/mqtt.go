// Package boost
// File:        mqtt.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/examples/mqtt.go
// Author:      Vibe Coding
// Created:     12/30/2025 11:03:46
// Description: Example demonstrating the use of the go-boost MQTT client and server
// --------------------------------------------------------------------------------

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mochi-mqtt/server/v2/packets"
	. "github.com/xiang-tai-duo/go-boost"
)

func main() {
	fmt.Println("MQTT Example - Testing Client and Server Functionality")
	fmt.Println("=========================================================")

	// Test 1: MQTT Server Functionality
	testMQTTServer()

	// Test 2: MQTT Client Functionality
	testMQTTClient()

	// Test 3: Third-party Publish Callback Functionality
	testThirdPartyPublishCallback()

	// Test 4: Simple Publish Callback Test
	testSimplePublishCallback()

	// Test 5: All Callback Functionality
	testAllCallbacks()

	// Test 6: TLS Server Functionality (Demo)
	testMQTTServerTLS()

	fmt.Println("All MQTT tests completed!")
}

// testMQTTServer tests MQTT server functionality
func testMQTTServer() {
	fmt.Println("\n=== MQTT Server Test ===")
	fmt.Println("Testing MQTT server functionality...")

	// Create MQTT server on default port 1883
	server := NewMQTTServer()
	fmt.Printf("Created MQTT server: %s:%d\n", server.GetHost(), server.GetPort())

	// Start the server
	fmt.Println("Starting MQTT server...")
	err := server.Start()
	if err != nil {
		log.Printf("Failed to start MQTT server: %v\n", err)
		return
	}
	fmt.Println("MQTT server started successfully")

	// Check if server is running
	if server.IsRunning() {
		fmt.Println("Server status: Running")
	} else {
		fmt.Println("Server status: Not running")
	}

	// Give server time to initialize
	time.Sleep(1 * time.Second)

	// Test server publishing functionality with different parameter combinations
	topic := "server/test/topic"
	payload := "Hello from MQTT server"

	// Test with minimal parameters (only topic and payload)
	fmt.Printf("Publishing message with minimal parameters: Topic='%s', Payload='%s'\n", topic, payload)
	err = server.Publish(topic, payload)
	if err != nil {
		log.Printf("Failed to publish with minimal parameters: %v\n", err)
	} else {
		fmt.Println("Message published successfully with minimal parameters")
	}

	// Test with QoS parameter
	fmt.Printf("Publishing message with QoS: Topic='%s', Payload='%s', QoS=1\n", topic, payload)
	err = server.Publish(topic, payload, 1)
	if err != nil {
		log.Printf("Failed to publish with QoS: %v\n", err)
	} else {
		fmt.Println("Message published successfully with QoS")
	}

	// Test with all parameters
	fmt.Printf("Publishing message with all parameters: Topic='%s', Payload='%s', QoS=2, Retained=true\n", topic, payload)
	err = server.Publish(topic, payload, 2, true)
	if err != nil {
		log.Printf("Failed to publish with all parameters: %v\n", err)
	} else {
		fmt.Println("Message published successfully with all parameters")
	}

	// Test server subscription functionality
	fmt.Println("\n=== Testing Server Subscription Functionality ===")

	// Get the server's address
	serverAddr := fmt.Sprintf("tcp://%s:%d", server.GetHost(), server.GetPort())

	// Create a separate client to test server subscription
	testClient := NewMQTT(serverAddr)
	if err := testClient.Connect(); err == nil {
		fmt.Println("Test client connected successfully")

		// Set message handler for test client
		testClient.SetMessageHandler(func(msg *MQTT_MESSAGE) {
			log.Printf("Test client received message: Topic='%s', Payload='%s', QoS=%d\n",
				msg.Topic, msg.Payload, msg.QoS)
		})

		// Test client subscribes to a topic
		testTopic := "server/subscribe/test"
		fmt.Printf("Test client subscribing to topic: %s\n", testTopic)
		if err := testClient.Subscribe(testTopic, 0); err == nil {
			fmt.Printf("Test client subscribed to topic: %s\n", testTopic)

			// Server publishes to the topic
			payload := "Hello from server to subscribed client"
			fmt.Printf("Server publishing to topic: %s\n", testTopic)
			if err := server.Publish(testTopic, payload); err == nil {
				fmt.Println("Server published message successfully")
				// Wait for message to be received
				time.Sleep(2 * time.Second)
			} else {
				log.Printf("Failed to publish from server: %v\n", err)
			}

			// Test client unsubscribes
			fmt.Printf("Test client unsubscribing from topic: %s\n", testTopic)
			testClient.Unsubscribe(testTopic)
		}

		// Disconnect test client
		testClient.Disconnect(1 * time.Second)
		fmt.Println("Test client disconnected")
	} else {
		log.Printf("Failed to connect test client: %v\n", err)
	}
	// Stop the server after test
	fmt.Println("Stopping MQTT server...")
	err = server.Stop()
	if err != nil {
		log.Printf("Failed to stop MQTT server: %v\n", err)
	} else {
		fmt.Println("MQTT server stopped successfully")
	}

	// Verify server is stopped
	if !server.IsRunning() {
		fmt.Println("Server status: Stopped")
	} else {
		fmt.Println("Server status: Still running")
	}

	fmt.Println("MQTT server test completed!")
	fmt.Println("========================================\n")
}

// testMQTTClient tests MQTT client functionality
func testMQTTClient() {
	fmt.Println("=== MQTT Client Test ===")
	fmt.Println("Testing MQTT client functionality...")

	// Define test parameters
	topic := "client/test/topic"

	// Test with public broker
	broker := "tcp://test.mosquitto.org:1883"
	fmt.Printf("Using broker: %s\n", broker)

	// Create MQTT client
	client := NewMQTT(broker)
	fmt.Printf("Created MQTT client with ID: %s\n", client.GetClientID())

	// Configure client
	client.SetQoS(0)
	client.SetKeepAlive(30 * time.Second)
	client.SetConnectionTimeout(15 * time.Second)

	// Set message handler
	client.SetMessageHandler(func(msg *MQTT_MESSAGE) {
		log.Printf("Received message: Topic='%s', Payload='%s', QoS=%d\n",
			msg.Topic, msg.Payload, msg.QoS)
	})

	// Set connect handler
	client.SetConnectHandler(func() {
		fmt.Println("Connected to MQTT broker")
	})

	// Set disconnect handler
	client.SetDisconnectHandler(func() {
		fmt.Println("Disconnected from MQTT broker")
	})

	// Connect to broker
	fmt.Println("Connecting to MQTT broker...")
	err := client.Connect()
	if err != nil {
		log.Printf("Failed to connect: %v\n", err)
		return
	}

	// Give time for connection to establish
	time.Sleep(2 * time.Second)

	// Subscribe to topic
	fmt.Printf("Subscribing to topic: %s\n", topic)
	err = client.Subscribe(topic, 0)
	if err != nil {
		log.Printf("Failed to subscribe: %v\n", err)
	} else {
		fmt.Printf("Subscribed to topic: %s\n", topic)
	}

	// Publish a message with minimal parameters
	payload := "Hello from MQTT client"
	fmt.Printf("Publishing message with minimal parameters: Topic='%s', Payload='%s'\n", topic, payload)
	err = client.Publish(topic, payload)
	if err != nil {
		log.Printf("Failed to publish with minimal parameters: %v\n", err)
	} else {
		fmt.Println("Message published successfully with minimal parameters")
	}

	// Wait for message to be received
	time.Sleep(2 * time.Second)

	// Publish a message with QoS parameter
	payload = "Hello from MQTT client with QoS"
	fmt.Printf("Publishing message with QoS: Topic='%s', Payload='%s', QoS=1\n", topic, payload)
	err = client.Publish(topic, payload, 1)
	if err != nil {
		log.Printf("Failed to publish with QoS: %v\n", err)
	} else {
		fmt.Println("Message published successfully with QoS")
	}

	// Wait for message to be received
	time.Sleep(2 * time.Second)

	// Test multiple subscriptions
	fmt.Println("\n=== Testing Multiple Subscriptions ===")
	topics := []string{
		"client/test/topic1",
		"client/test/topic2",
		"client/test/topic3",
	}

	// Subscribe to multiple topics
	for _, t := range topics {
		fmt.Printf("Subscribing to topic: %s\n", t)
		err = client.Subscribe(t, 0)
		if err != nil {
			log.Printf("Failed to subscribe to %s: %v\n", t, err)
		} else {
			fmt.Printf("Subscribed to topic: %s\n", t)
		}

		// Publish to each topic
		payload := fmt.Sprintf("Hello to %s", t)
		err = client.Publish(t, payload)
		if err != nil {
			log.Printf("Failed to publish to %s: %v\n", t, err)
		} else {
			fmt.Printf("Published to topic: %s\n", t)
		}
	}

	// Wait for messages to be received
	time.Sleep(3 * time.Second)

	// Get all subscriptions
	subscriptions := client.GetSubscriptions()
	fmt.Printf("\nAll subscriptions: %v\n", subscriptions)

	// Unsubscribe from all topics
	fmt.Println("\nUnsubscribing from all topics...")
	for t := range subscriptions {
		err = client.Unsubscribe(t)
		if err != nil {
			log.Printf("Failed to unsubscribe from %s: %v\n", t, err)
		} else {
			fmt.Printf("Unsubscribed from topic: %s\n", t)
		}
	}

	// Disconnect client
	fmt.Println("Disconnecting from MQTT broker...")
	client.Disconnect(5 * time.Second)

	fmt.Println("MQTT client test completed!")
	fmt.Println("========================================\n")
}

// testThirdPartyPublishCallback tests the publish callback functionality when third-party clients publish messages
func testThirdPartyPublishCallback() {
	fmt.Println("=== Third-party Publish Callback Test ===")
	fmt.Println("Testing publish callback functionality...")

	// Create MQTT server on port 1884
	server := NewMQTTServer("0.0.0.0", 1884)
	fmt.Printf("Created MQTT server: %s:%d\n", server.GetHost(), server.GetPort())

	// Set up publish callback for the server
	messageReceived := make(chan *MQTT_MESSAGE, 10)
	server.SetPublishCallback(func(client *server.Client, packet packets.Packet) {
		msg := server.GetMQTTMessage(packet)
		fmt.Printf("Server received published message: Topic='%s', Payload='%s', QoS=%d\n",
			msg.Topic, msg.Payload, msg.QoS)
		messageReceived <- msg
	})

	// Start the server
	fmt.Println("Starting MQTT server...")
	err := server.Start()
	if err != nil {
		log.Printf("Failed to start MQTT server: %v\n", err)
		return
	}
	fmt.Println("MQTT server started successfully")

	// Give server time to initialize
	time.Sleep(1 * time.Second)

	// Create a third-party client
	serverAddr := fmt.Sprintf("tcp://%s:%d", server.GetHost(), server.GetPort())
	thirdPartyClient := NewMQTT(serverAddr)
	fmt.Printf("Created third-party client with ID: %s\n", thirdPartyClient.GetClientID())

	// Connect third-party client to the server
	if err := thirdPartyClient.Connect(); err == nil {
		fmt.Println("Third-party client connected successfully")

		// Test multiple publish scenarios
		testCases := []struct {
			topic    string
			payload  string
			qos      int
			retained bool
		}{
			{"thirdparty/test/topic1", "Hello from third party client - minimal", 0, false},
			{"thirdparty/test/topic2", "Hello from third party client - QoS 1", 1, false},
			{"thirdparty/test/topic3", "Hello from third party client - retained", 0, true},
			{"thirdparty/sensor/temp", "23.5", 0, false},
			{"thirdparty/sensor/humidity", "65%", 0, false},
		}

		// Execute all test cases
		for i, test := range testCases {
			fmt.Printf("\nTest case %d:\n", i+1)
			fmt.Printf("Third-party client publishing: Topic='%s', Payload='%s', QoS=%d, Retained=%v\n",
				test.topic, test.payload, test.qos, test.retained)

			// Publish with appropriate parameters
			var publishErr error
			if test.retained {
				publishErr = thirdPartyClient.Publish(test.topic, test.payload, test.qos, test.retained)
			} else if test.qos > 0 {
				publishErr = thirdPartyClient.Publish(test.topic, test.payload, test.qos)
			} else {
				publishErr = thirdPartyClient.Publish(test.topic, test.payload)
			}

			if publishErr != nil {
				log.Printf("Third-party client failed to publish: %v\n", publishErr)
				continue
			}

			fmt.Println("Third-party client published message successfully")

			// Wait for callback to be triggered
			time.Sleep(1 * time.Second)
		}

		// Wait for all messages to be processed
		time.Sleep(2 * time.Second)

		// Disconnect third-party client
		thirdPartyClient.Disconnect(1 * time.Second)
		fmt.Println("Third-party client disconnected")
	} else {
		log.Printf("Failed to connect third-party client: %v\n", err)
	}

	// Stop the server
	fmt.Println("\nStopping MQTT server...")
	err = server.Stop()
	if err != nil {
		log.Printf("Failed to stop MQTT server: %v\n", err)
	} else {
		fmt.Println("MQTT server stopped successfully")
	}

	fmt.Println("Third-party publish callback test completed!")
	fmt.Println("========================================\n")
}

// testSimplePublishCallback tests the basic publish callback functionality
func testSimplePublishCallback() {
	fmt.Println("=== Simple Publish Callback Test ===")
	fmt.Println("Testing basic publish callback functionality...")

	// Create MQTT server on port 1885
	server := NewMQTTServer("0.0.0.0", 1885)
	fmt.Printf("Created MQTT server: %s:%d\n", server.GetHost(), server.GetPort())

	// Set up publish callback
	messageCount := 0
	server.SetPublishCallback(func(client *server.Client, packet packets.Packet) {
		msg := server.GetMQTTMessage(packet)
		messageCount++
		fmt.Printf("\n=== Received Published Message %d ===\n", messageCount)
		fmt.Printf("Topic:     %s\n", msg.Topic)
		fmt.Printf("Payload:   %s\n", msg.Payload)
		fmt.Printf("Timestamp: %s\n", msg.Timestamp.Format("2006-01-02 15:04:05.000"))
		fmt.Printf("QoS:       %d\n", msg.QoS)
		fmt.Printf("Retained:  %v\n", msg.Retained)
		fmt.Printf("Duplicate: %v\n", msg.Duplicate)
		fmt.Println("====================================")
	})

	// Start the server
	fmt.Println("\nStarting MQTT server...")
	err := server.Start()
	if err != nil {
		log.Printf("Failed to start MQTT server: %v\n", err)
		return
	}
	fmt.Println("MQTT server started successfully")

	// Test server publishing functionality
	testCases := []struct {
		topic    string
		payload  string
		qos      int
		retained bool
	}{{
		topic:    "test/simple",
		payload:  "Hello from server!",
		qos:      0,
		retained: false,
	}, {
		topic:    "test/with_qos",
		payload:  "Message with QoS 1",
		qos:      1,
		retained: false,
	}, {
		topic:    "test/retained",
		payload:  "This is a retained message",
		qos:      0,
		retained: true,
	}, {
		topic:    "test/complex/topic/structure",
		payload:  "Message with complex topic",
		qos:      2,
		retained: false,
	}}

	// Publish messages with different parameters
	for i, tc := range testCases {
		time.Sleep(500 * time.Millisecond)
		fmt.Printf("\nPublishing test %d: Topic='%s', Payload='%s', QoS=%d, Retained=%v\n",
			i+1, tc.topic, tc.payload, tc.qos, tc.retained)

		// Test with full parameter set
		err := server.Publish(tc.topic, tc.payload, tc.qos, tc.retained)
		if err != nil {
			log.Printf("Failed to publish test %d: %v\n", i+1, err)
		} else {
			fmt.Printf("Published test %d successfully\n", i+1)
		}
	}

	// Test with minimal parameters (only topic and payload)
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("\nPublishing with minimal parameters: Topic='test/minimal', Payload='Minimal params test'\n")
	err = server.Publish("test/minimal", "Minimal params test")
	if err != nil {
		log.Printf("Failed to publish minimal test: %v\n", err)
	} else {
		fmt.Printf("Published minimal test successfully\n")
	}

	// Test with only topic, payload, and QoS
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("\nPublishing with topic, payload, QoS: Topic='test/partial', Payload='Partial params test', QoS=1\n")
	err = server.Publish("test/partial", "Partial params test", 1)
	if err != nil {
		log.Printf("Failed to publish partial test: %v\n", err)
	} else {
		fmt.Printf("Published partial test successfully\n")
	}

	// Wait for all messages to be processed
	time.Sleep(2 * time.Second)

	// Summary
	fmt.Printf("\n\n=== Test Summary ===\n")
	fmt.Printf("Total messages received: %d\n", messageCount)
	fmt.Printf("Expected messages: %d\n", len(testCases)+2)
	fmt.Printf("Test %s\n", map[bool]string{true: "PASSED", false: "FAILED"}[messageCount == len(testCases)+2])

	// Stop the server
	fmt.Println("\nStopping MQTT server...")
	server.Stop()
	fmt.Println("MQTT server stopped successfully")
	fmt.Println("\nSimple publish callback test completed!")
	fmt.Println("========================================\n")
}

// testAllCallbacks tests all available callback functions
func testAllCallbacks() {
	fmt.Println("=== All Callbacks Test ===")
	fmt.Println("Testing all MQTT callback functions...")

	// Create MQTT server on port 1886
	server := NewMQTTServer("0.0.0.0", 1886)
	fmt.Printf("Created MQTT server: %s:%d\n", server.GetHost(), server.GetPort())

	// Track which callbacks were triggered
	callbacksTriggered := make(map[string]bool)

	// Set up all callbacks

	// Set Auth callback
	server.SetAuthCallback(func(client *server.Client, packet packets.Packet) bool {
		clientID := client.ID
		username := string(packet.Connect.Username)
		password := string(packet.Connect.Password)
		fmt.Printf("AuthCallback triggered: ClientID='%s', Username='%s', Password='%s'\n", clientID, username, password)
		callbacksTriggered["AuthCallback"] = true
		return true // Allow all connections for testing
	})

	// Set Connect callback
	server.SetConnectCallback(func(client *server.Client, packet packets.Packet) {
		clientID := client.ID
		fmt.Printf("ConnectCallback triggered: ClientID='%s'\n", clientID)
		callbacksTriggered["ConnectCallback"] = true
	})

	// Set Disconnect callback
	server.SetDisconnectCallback(func(client *server.Client, err error) {
		clientID := client.ID
		fmt.Printf("DisconnectCallback triggered: ClientID='%s', Error='%v'\n", clientID, err)
		callbacksTriggered["DisconnectCallback"] = true
	})

	// Set Subscribe callback
	server.SetSubscribeCallback(func(client *server.Client, packet packets.Packet) {
		clientID := client.ID
		for _, sub := range packet.Filters {
			fmt.Printf("SubscribeCallback triggered: ClientID='%s', Topic='%s', QoS=%d\n", clientID, sub.Filter, sub.Qos)
		}
		callbacksTriggered["SubscribeCallback"] = true
	})

	// Set Subscribed callback
	server.SetSubscribedCallback(func(client *server.Client, packet packets.Packet) {
		clientID := client.ID
		for _, sub := range packet.Filters {
			fmt.Printf("SubscribedCallback triggered: ClientID='%s', Topic='%s', QoS=%d\n", clientID, sub.Filter, sub.Qos)
		}
		callbacksTriggered["SubscribedCallback"] = true
	})

	// Set Unsubscribe callback
	server.SetUnsubscribeCallback(func(client *server.Client, packet packets.Packet) {
		clientID := client.ID
		for _, sub := range packet.Filters {
			fmt.Printf("UnsubscribeCallback triggered: ClientID='%s', Topic='%s'\n", clientID, sub.Filter)
		}
		callbacksTriggered["UnsubscribeCallback"] = true
	})

	// Set Unsubscribed callback
	server.SetUnsubscribedCallback(func(client *server.Client, packet packets.Packet) {
		clientID := client.ID
		for _, sub := range packet.Filters {
			fmt.Printf("UnsubscribedCallback triggered: ClientID='%s', Topic='%s'\n", clientID, sub.Filter)
		}
		callbacksTriggered["UnsubscribedCallback"] = true
	})

	// Set Publish callback
	server.SetPublishCallback(func(client *server.Client, packet packets.Packet) {
		mqttMessage := server.GetMQTTMessage(packet)
		fmt.Printf("PublishCallback triggered: Topic='%s', Payload='%s', QoS=%d\n", mqttMessage.Topic, mqttMessage.Payload, mqttMessage.QoS)
		callbacksTriggered["PublishCallback"] = true
	})

	// Set ACL Check callback
	server.SetACLCheckCallback(func(client *server.Client, topic string, write bool) bool {
		clientID := client.ID
		action := "read"
		if write {
			action = "write"
		}
		fmt.Printf("ACLCheckCallback triggered: ClientID='%s', Topic='%s', Action='%s'\n", clientID, topic, action)
		callbacksTriggered["ACLCheckCallback"] = true
		return true // Allow all operations for testing
	})

	// Start the server
	fmt.Println("\nStarting MQTT server...")
	err := server.Start()
	if err != nil {
		log.Printf("Failed to start MQTT server: %v\n", err)
		return
	}
	fmt.Println("MQTT server started successfully")

	// Give server time to initialize
	time.Sleep(1 * time.Second)

	// Create a client
	serverAddr := fmt.Sprintf("tcp://%s:%d", server.GetHost(), server.GetPort())
	client := NewMQTT(serverAddr)
	fmt.Printf("Created client with ID: %s\n", client.GetClientID())

	// Connect client
	fmt.Println("\nConnecting client to server...")
	if err := client.Connect(); err != nil {
		log.Printf("Failed to connect client: %v\n", err)
		server.Stop()
		return
	}
	fmt.Println("Client connected successfully")

	// Give time for connect callbacks to be triggered
	time.Sleep(1 * time.Second)

	// Subscribe to a topic
	topic := "callback/test/topic"
	fmt.Printf("\nSubscribing client to topic: %s\n", topic)
	if err := client.Subscribe(topic, 0); err != nil {
		log.Printf("Failed to subscribe: %v\n", err)
	} else {
		fmt.Printf("Client subscribed to topic: %s\n", topic)
	}

	// Give time for subscribe callbacks to be triggered
	time.Sleep(1 * time.Second)

	// Publish a message
	payload := "Hello from client to test publish callback"
	fmt.Printf("\nPublishing message: Topic='%s', Payload='%s'\n", topic, payload)
	if err := client.Publish(topic, payload); err != nil {
		log.Printf("Failed to publish: %v\n", err)
	} else {
		fmt.Println("Client published message successfully")
	}

	// Give time for publish callbacks to be triggered
	time.Sleep(2 * time.Second)

	// Unsubscribe from topic
	fmt.Printf("\nUnsubscribing client from topic: %s\n", topic)
	if err := client.Unsubscribe(topic); err != nil {
		log.Printf("Failed to unsubscribe: %v\n", err)
	} else {
		fmt.Printf("Client unsubscribed from topic: %s\n", topic)
	}

	// Give time for unsubscribe callbacks to be triggered
	time.Sleep(1 * time.Second)

	// Disconnect client
	fmt.Println("\nDisconnecting client...")
	client.Disconnect(1 * time.Second)
	fmt.Println("Client disconnected")

	// Give time for disconnect callbacks to be triggered
	time.Sleep(1 * time.Second)

	// Stop the server
	fmt.Println("\nStopping MQTT server...")
	server.Stop()
	fmt.Println("MQTT server stopped successfully")

	// Summary
	fmt.Printf("\n\n=== Test Summary ===\n")
	fmt.Println("Callbacks expected to be triggered:")
	expectedCallbacks := []string{
		"AuthCallback",
		"ConnectCallback",
		"SubscribeCallback",
		"SubscribedCallback",
		"ACLCheckCallback",
		"PublishCallback",
		"UnsubscribeCallback",
		"UnsubscribedCallback",
		"DisconnectCallback",
	}

	allPassed := true
	for _, callback := range expectedCallbacks {
		status := "PASSED"
		if !callbacksTriggered[callback] {
			status = "FAILED"
			allPassed = false
		}
		fmt.Printf("%-20s: %s\n", callback, status)
	}

	fmt.Printf("\nOverall Test: %s\n", map[bool]string{true: "PASSED", false: "FAILED"}[allPassed])
	fmt.Println("All callbacks test completed!")
	fmt.Println("========================================\n")
}

// testMQTTServerTLS tests MQTT server TLS functionality
func testMQTTServerTLS() {
	fmt.Println("=== MQTT TLS Server Test ===")
	fmt.Println("Testing MQTT TLS server functionality...")

	// Test 1: Create TLS server with auto-generated self-signed certificate
	host := "0.0.0.0"
	port := 8883

	fmt.Printf("\n1. Testing TLS server with auto-generated self-signed certificate\n")
	tlsServer := NewMQTTServerTLS(host, port)
	fmt.Printf("   Created TLS server: %s:%d\n", tlsServer.GetHost(), tlsServer.GetPort())

	// Set up publish callback
	tlsServer.SetPublishCallback(func(client *server.Client, packet packets.Packet) {
		msg := tlsServer.GetMQTTMessage(packet)
		fmt.Printf("   TLS Server received message: Topic='%s', Payload='%s'\n", msg.Topic, msg.Payload)
	})

	// Start the TLS server with auto-generated cert
	fmt.Println("   Starting TLS server with auto-generated self-signed certificate...")
	err := tlsServer.Start()
	if err != nil {
		log.Printf("   Failed to start TLS server: %v\n", err)
	} else {
		fmt.Println("   TLS server started successfully with auto-generated certificate")

		// Wait a bit for server to initialize
		time.Sleep(2 * time.Second)

		// Stop the server
		fmt.Println("   Stopping TLS server...")
		tlsServer.Stop()
		fmt.Println("   TLS server stopped successfully")
	}

	// Test 2: Show how to use TLS with certificate files
	fmt.Printf("\n2. Example: Using TLS with certificate files\n")
	fmt.Println("   In a real scenario, you would use:")
	fmt.Println("   certFile := \"/path/to/cert.pem\"")
	fmt.Println("   keyFile := \"/path/to/key.pem\"")
	fmt.Println("   tlsServer := NewMQTTServerTLS(host, port)")
	fmt.Println("   tlsServer.Start(certFile, keyFile)")

	// Test 3: Show how to use TLS with PFX file
	fmt.Printf("\n3. Example: Using TLS with PFX file\n")
	fmt.Println("   In a real scenario, you would use:")
	fmt.Println("   pfxFile := \"/path/to/cert.pfx\"")
	fmt.Println("   pfxPassword := \"your_password\"")
	fmt.Println("   tlsServer := NewMQTTServerTLS(host, port)")
	fmt.Println("   tlsServer.Start(pfxFile, pfxPassword)")

	// Test 4: Create TLS server and test with different parameters
	fmt.Printf("\n4. Testing TLS server creation with different parameters\n")

	// Test with invalid cert path (should fail gracefully)
	invalidCertServer := NewMQTTServerTLS(host, port+1)
	fmt.Printf("   Testing with invalid certificate path...\n")
	err = invalidCertServer.Start("/invalid/path/to/cert.pem", "/invalid/path/to/key.pem")
	if err != nil {
		fmt.Printf("   Expected error: %v\n", err)
		fmt.Println("   ✓ Correctly handled invalid certificate path")
	}

	fmt.Println("\nMQTT TLS server tests completed!")
	fmt.Println("========================================\n")
}
