// Package boost
// File:        soap.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/examples/soap.go
// Author:      Vibe Coding
// Created:     12/30/2025 11:03:46
// Description: Example for SOAP utility functions
// --------------------------------------------------------------------------------

package main

import (
	"fmt"

	. "github.com/xiang-tai-duo/go-boost"
)

func main() {
	soapClient := NewSOAP()

	fmt.Println("--- SOAP Examples ---")

	fmt.Println("\n1. Creating SOAP Envelope:")
	body := "<GetWeather xmlns=\"http://example.com/weather\"><City>London</City></GetWeather>"
	envelope, err := soapClient.CreateEnvelope(body)
	if err == nil {
		fmt.Println("SOAP Envelope created successfully:")
		fmt.Println(envelope)
	} else {
		fmt.Printf("Error creating envelope: %v\n", err)
	}

	fmt.Println("\n2. Extracting Body from SOAP Response:")
	soapResponse := `<?xml version="1.0" encoding="UTF-8"?><soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"><soapenv:Body><GetWeatherResponse xmlns="http://example.com/weather"><Temperature>22</Temperature></GetWeatherResponse></soapenv:Body></soapenv:Envelope>`
	bodyContent := soapClient.ExtractBody(soapResponse)
	fmt.Println("Extracted body:")
	fmt.Println(bodyContent)

	fmt.Println("\n--- SOAP Examples Complete ---")
}
