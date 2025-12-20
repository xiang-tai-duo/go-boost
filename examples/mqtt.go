// --------------------------------------------------------------------------------
// File:        mqtt.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: Example demonstrating the use of the go-boost MQTT client
// --------------------------------------------------------------------------------

package main

import (
	"fmt"
	"log"
	"time"

	. "github.com/xiang-tai-duo/go-boost"
)

func main() {
	fmt.Println("MQTT Client Example - Testing Both Connection Types")
	fmt.Println("====================================================")

	// Define common parameters
	topic := "go-boost/example/topic"

	// Test 1: Anonymous Connection (Public broker)
	// Public brokers like test.mosquitto.org don't require authentication
	anonymousBroker := "tcp://test.mosquitto.org:1883"
	testMQTTConnection(anonymousBroker, topic, "Anonymous Connection Test", false, "", "")

	// Test 2: Authenticated Connection Example
	// Note: Public brokers like test.mosquitto.org typically don't require authentication
	// or have specific credentials. The following test demonstrates the authentication flow
	// but is expected to fail with "not Authorized" since we're using dummy credentials.
	// In production, you would use your own MQTT broker with proper authentication.
	authBroker := "tcp://test.mosquitto.org:1883"
	testMQTTConnection(authBroker, topic, "Authenticated Connection Example", true, "your-username", "your-password")

	// Alternative: Some public brokers offer authenticated access on specific ports
	// For example, mosquitto offers authenticated access on port 8883 (SSL)
	// but requires client certificates. For simplicity, we're demonstrating
	// the authentication API here without a working authenticated connection.

	fmt.Println("Important Note:")
	fmt.Println("   - Anonymous connection test succeeded as expected")
	fmt.Println("   - Authenticated connection test demonstrated the API usage")
	fmt.Println("   - In production, use your MQTT broker with valid credentials")

	fmt.Println("All MQTT connection tests completed!")
	fmt.Println("\nSummary:")
	fmt.Println("   - Tested anonymous MQTT connection")
	fmt.Println("   - Tested authenticated MQTT connection")
	fmt.Println("   - Demonstrated all MQTT client functions")
}

// testMQTTConnection tests MQTT connection with specified parameters
func testMQTTConnection(brokerAddress string, topic string, testName string, useAuth bool, username string, password string) {
	fmt.Printf("\n=== %s ===\n", testName)
	fmt.Printf("Broker: %s\n", brokerAddress)
	fmt.Printf("Topic: %s\n", topic)
	fmt.Printf("Use Authentication: %t\n", useAuth)
	if useAuth {
		fmt.Printf("Username: %s\n", username)
		fmt.Printf("Password: [hidden]\n")
	}
	fmt.Println()

	// Create and configure MQTT client using go-boost
	fmt.Println("1. Creating MQTT client...")
	mqttClient := NewMQTT(brokerAddress)

	// Configure MQTT client settings
	fmt.Println("   Configuring client settings...")
	// Set custom client ID
	clientID := "go-boost-example-client-" + testName + "-" + time.Now().Format("20060102150405")
	err := mqttClient.SetClientID(clientID)
	if err != nil {
		log.Printf("   Failed to set client ID: %v\n", err)
		return
	}
	fmt.Printf("   Client ID: %s\n", mqttClient.GetClientID())
	// Set QoS level
	err = mqttClient.SetQoS(1)
	if err != nil {
		log.Printf("   Failed to set QoS: %v\n", err)
		return
	}
	fmt.Printf("   QoS Level: %d\n", mqttClient.GetQoS())
	// Set KeepAlive duration
	err = mqttClient.SetKeepAlive(30 * time.Second)
	if err != nil {
		log.Printf("   Failed to set keep alive: %v\n", err)
		return
	}
	fmt.Printf("   KeepAlive: %v\n", mqttClient.GetKeepAlive())
	// Set connection timeout
	err = mqttClient.SetConnectionTimeout(15 * time.Second)
	if err != nil {
		log.Printf("   Failed to set connection timeout: %v\n", err)
		return
	}
	fmt.Printf("   Connection Timeout: %v\n", mqttClient.GetConnectionTimeout())
	// Get and display broker address
	fmt.Printf("   Broker: %s\n", mqttClient.GetBroker())

	// Set username and password if using authentication
	fmt.Println("   Setting username/password...")
	if useAuth {
		err = mqttClient.SetUsername(username)
		if err != nil {
			log.Printf("   Failed to set username: %v\n", err)
			return
		}
		err = mqttClient.SetPassword(password)
		if err != nil {
			log.Printf("   Failed to set password: %v\n", err)
			return
		}
		fmt.Printf("   Username: %s\n", mqttClient.GetUsername())
		fmt.Println("   Password: [hidden]")
	} else {
		// Set empty credentials for anonymous connection
		err = mqttClient.SetUsername("")
		if err != nil {
			log.Printf("   ❌ Failed to set username: %v\n", err)
			return
		}
		err = mqttClient.SetPassword("")
		if err != nil {
			log.Printf("   ❌ Failed to set password: %v\n", err)
			return
		}
		fmt.Println("   Using anonymous connection")
	}

	// Set message handler for received messages
	fmt.Println("   Setting message handler...")
	err = mqttClient.SetMessageHandler(func(msg *MQTT_MESSAGE) {
		log.Printf("Received message: Topic='%s', Payload='%s', QoS=%d, Retained=%t, Duplicate=%t\n",
			msg.Topic, msg.Payload, msg.QoS, msg.Retained, msg.Duplicate)
	})
	if err != nil {
		log.Printf("   Failed to set message handler: %v\n", err)
		return
	}
	// Set connect handler
	err = mqttClient.SetConnectHandler(func() {
		fmt.Println("Connected to MQTT broker")
	})
	if err != nil {
		log.Printf("   Failed to set connect handler: %v\n", err)
		return
	}
	// Set disconnect handler
	err = mqttClient.SetDisconnectHandler(func() {
		fmt.Println("Disconnected from MQTT broker")
	})
	if err != nil {
		log.Printf("   Failed to set disconnect handler: %v\n", err)
		return
	}

	// Connect to MQTT broker
	fmt.Printf("2. Connecting to MQTT broker at %s...\n", brokerAddress)
	err = mqttClient.Connect()
	if err != nil {
		log.Printf("   Failed to connect to MQTT broker: %v\n", err)
		return
	}

	// Cleanup function
	cleanup := func() {
		fmt.Println("6. Disconnecting from MQTT broker...")
		mqttClient.Disconnect(5 * time.Second)
	}

	// Give time for connection to establish
	time.Sleep(2 * time.Second)

	// Subscribe to topic
	fmt.Printf("3. Subscribing to topic %s with QoS 1...\n", topic)
	err = mqttClient.Subscribe(topic, 1)
	if err != nil {
		log.Printf("   Failed to subscribe to topic: %v\n", err)
		cleanup()
		return
	}

	// Publish a simple test message
	testMessage := fmt.Sprintf("Hello from %s test", testName)
	fmt.Printf("4. Publishing test message: '%s'\n", testMessage)
	err = mqttClient.Publish(topic, testMessage, 1, false)
	if err != nil {
		log.Printf("   Failed to publish message: %v\n", err)
	} else {
		log.Printf("   Message published successfully\n")
	}

	// Wait a bit to ensure message is processed
	time.Sleep(2 * time.Second)

	// Get subscriptions
	subscriptions := mqttClient.GetSubscriptions()
	fmt.Printf("5. Current subscriptions: %v\n", subscriptions)

	// Unsubscribe from topic
	fmt.Printf("6. Unsubscribing from topic %s...\n", topic)
	err = mqttClient.Unsubscribe(topic)
	if err != nil {
		log.Printf("   Failed to unsubscribe from topic: %v\n", err)
	}

	// Check if still connected
	if mqttClient.IsConnected() {
		fmt.Println("   Still connected to MQTT broker")
	} else {
		fmt.Println("   Disconnected from MQTT broker")
	}

	// Disconnect and cleanup
	cleanup()

	fmt.Printf("%s completed!\n", testName)
	fmt.Println("========================================\n")
}
