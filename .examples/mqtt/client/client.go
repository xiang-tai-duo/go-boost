// Package main
// File:        https://github.com/xiang-tai-duo/go-boost/blob/master/.examples/mqtt/client/client.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: MQTT client usage example
// --------------------------------------------------------------------------------
package main

import (
	"fmt"
	"time"

	mqttclient "github.com/xiang-tai-duo/go-boost/mqtt/client"
	mqttcommon "github.com/xiang-tai-duo/go-boost/mqtt/common"
)

//goland:noinspection GoUnhandledErrorResult
func main() {
	// Example 1: Create MQTT client with broker address
	broker := "tcp://localhost:1883"
	mqtt := mqttclient.New(broker)

	// Example 2: Set connect and disconnect handlers
	mqtt.SetConnectHandler(func() {
		fmt.Println("Connected to MQTT broker")
	})

	mqtt.SetDisconnectHandler(func() {
		fmt.Println("Disconnected from MQTT broker")
	})

	// Example 3: Set message handler
	mqtt.SetMessageHandler(func(message *mqttcommon.MQTT_MESSAGE) {
		fmt.Printf("Received message: Topic=%s, Payload=%s, QualityOfService=%d\n",
			message.Topic, message.Payload, message.QualityOfService)
	})

	// Example 4: Connect to MQTT broker
	var err error
	if err = mqtt.Connect(); err == nil {
		defer func() {
			mqtt.Disconnect(5 * time.Second)
		}()

		// Example 5: Subscribe to a topic
		topic := "sensors/temperature"
		qos := 1 // Quality of Service level
		if err = mqtt.Subscribe(topic, qos); err == nil {
			fmt.Printf("Subscribed to topic: %s\n", topic)

			// Example 6: Publish a message
			payload := "{\"temperature\": 25.5, \"humidity\": 60}"
			if err = mqtt.Publish(topic, payload, qos, false); err == nil { // qos=1, retained=false
				fmt.Printf("Published message: %s\n", payload)

				// Example 7: Get client information
				fmt.Printf("Broker: %s\n", mqtt.GetBroker())
				fmt.Printf("Client ID: %s\n", mqtt.GetClientID())
				fmt.Printf("Is connected: %v\n", mqtt.IsConnected())

				// Wait for messages
				fmt.Println("Waiting for messages... (Press Ctrl+C to exit)")
				select {
				case <-time.After(10 * time.Second):
					fmt.Println("Timeout: Exiting after 10 seconds")
				}
			} else {
				fmt.Printf("Failed to publish: %v\n", err)
			}
		} else {
			fmt.Printf("Failed to subscribe: %v\n", err)
		}
	} else {
		fmt.Printf("Failed to connect: %v\n", err)
	}
}
