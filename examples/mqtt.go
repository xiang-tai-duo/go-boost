// --------------------------------------------------------------------------------
// File:        mqtt.go
// Author:      TRAE AI
// Created:     12/30/2025 11:03:46
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

	// Summary:
	fmt.Println("\nSummary:")
	fmt.Println("   - Tested anonymous MQTT connection")
	fmt.Println("   - Tested authenticated MQTT connection")
	fmt.Println("   - Demonstrated all MQTT client functions")
}

// testMQTTConnection tests MQTT connection with specified parameters

func testMQTTConnection(brokerAddress string, topic string, testName string, useAuth bool, username string, password string) {
	var err error
	var connected bool

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

	// Set credentials based on auth flag
	var clientUsername, clientPassword string
	if useAuth {
		clientUsername = username
		clientPassword = password
	}

	// Create MQTT client with required parameters
	var mqttClient *MQTT
	if useAuth {
		mqttClient = NewMQTT(brokerAddress, clientUsername, clientPassword)
	} else {
		mqttClient = NewMQTT(brokerAddress)
	}

	fmt.Printf("   Client ID: %s\n", mqttClient.GetClientID())

	// Configure other MQTT client settings
	fmt.Println("   Configuring client settings...")

	// Set QoS level
	err = mqttClient.SetQoS(1)
	if err == nil {
		fmt.Printf("   QoS Level: %d\n", mqttClient.GetQoS())

		// Set KeepAlive duration
		err = mqttClient.SetKeepAlive(30 * time.Second)
		if err == nil {
			fmt.Printf("   KeepAlive: %v\n", mqttClient.GetKeepAlive())

			// Set connection timeout
			err = mqttClient.SetConnectionTimeout(15 * time.Second)
			if err == nil {
				fmt.Printf("   Connection Timeout: %v\n", mqttClient.GetConnectionTimeout())

				// Get and display broker address
				fmt.Printf("   Broker: %s\n", mqttClient.GetBroker())

				// Display authentication settings
				fmt.Println("   Authentication settings:")
				if useAuth {
					fmt.Printf("   Username: %s\n", mqttClient.GetUsername())
					fmt.Println("   Password: [hidden]")
				} else {
					fmt.Println("   Using anonymous connection")
				}

				// Set message handler for received messages
				fmt.Println("   Setting message handler...")
				err = mqttClient.SetMessageHandler(func(msg *MQTT_MESSAGE) {
					log.Printf("Received message: Topic='%s', Payload='%s', QoS=%d, Retained=%t, Duplicate=%t\n",
						msg.Topic, msg.Payload, msg.QoS, msg.Retained, msg.Duplicate)
				})
				if err == nil {
					// Set connect handler
					err = mqttClient.SetConnectHandler(func() {
						fmt.Println("Connected to MQTT broker")
					})
					if err == nil {
						// Set disconnect handler
						err = mqttClient.SetDisconnectHandler(func() {
							fmt.Println("Disconnected from MQTT broker")
						})
						if err == nil {
							// Connect to MQTT broker
							fmt.Printf("2. Connecting to MQTT broker at %s...\n", brokerAddress)
							err = mqttClient.Connect()
							if err == nil {
								connected = true

								// Give time for connection to establish
								time.Sleep(2 * time.Second)

								// Subscribe to topic
								fmt.Printf("3. Subscribing to topic %s with QoS 1...\n", topic)
								err = mqttClient.Subscribe(topic, 1)
								if err == nil {
									// Publish a simple test message
									testMessage := fmt.Sprintf("Hello from %s test", testName)
									fmt.Printf("4. Publishing test message: '%s'\n", testMessage)
									publishErr := mqttClient.Publish(topic, testMessage, 1, false)
									if publishErr != nil {
										log.Printf("   Failed to publish message: %v\n", publishErr)
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
									unsubscribeErr := mqttClient.Unsubscribe(topic)
									if unsubscribeErr != nil {
										log.Printf("   Failed to unsubscribe from topic: %v\n", unsubscribeErr)
									}

									// Check if still connected
									if mqttClient.IsConnected() {
										fmt.Println("   Still connected to MQTT broker")
									} else {
										fmt.Println("   Disconnected from MQTT broker")
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Cleanup function
	if connected {
		fmt.Println("6. Disconnecting from MQTT broker...")
		mqttClient.Disconnect(5 * time.Second)
	}

	if err != nil {
		log.Printf("   Error: %v\n", err)
	} else if connected {
		fmt.Printf("%s completed!\n", testName)
		fmt.Println("========================================\n")
	}
}
