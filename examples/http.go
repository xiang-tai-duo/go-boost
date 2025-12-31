// --------------------------------------------------------------------------------
// File:        http.go
// Author:      TRAE AI
// Created:     12/30/2025 11:03:46
// Description: Example for HTTP utility functions
// --------------------------------------------------------------------------------

package main

import (
	"fmt"
	"time"

	. "github.com/xiang-tai-duo/go-boost"
)

func main() {

	// Create a new HTTP client
	httpClient := NewHTTP()

	// Set custom timeout
	httpClient.SetTimeout(10 * time.Second)
	fmt.Printf("HTTP client timeout: %v\n", httpClient.GetTimeout())

	// Example: GET request to a public API
	fmt.Println("\n--- GET Request Example ---")
	getURL := "https://jsonplaceholder.typicode.com/posts/1"
	responseBody, statusCode, err := httpClient.Get(getURL)
	if err == nil {
		fmt.Printf("GET %s\n", getURL)
		fmt.Printf("Status Code: %d\n", statusCode)
		fmt.Printf("Response Body: %s\n", responseBody)
	} else {
		fmt.Printf("GET request failed: %v\n", err)
	}

	// Example: POST request
	fmt.Println("\n--- POST Request Example ---")
	postURL := "https://jsonplaceholder.typicode.com/posts"
	postBody := `{
		"title": "foo",
		"body": "bar",
		"userId": 1
	}`
	responseBody, statusCode, err = httpClient.Post(postURL, "application/json", postBody)
	if err == nil {
		fmt.Printf("POST %s\n", postURL)
		fmt.Printf("Status Code: %d\n", statusCode)
		fmt.Printf("Response Body: %s\n", responseBody)
	} else {
		fmt.Printf("POST request failed: %v\n", err)
	}

	// Example: PUT request
	fmt.Println("\n--- PUT Request Example ---")
	putURL := "https://jsonplaceholder.typicode.com/posts/1"
	putBody := `{
		"id": 1,
		"title": "updated title",
		"body": "updated body",
		"userId": 1
	}`
	responseBody, statusCode, err = httpClient.Put(putURL, "application/json", putBody)
	if err == nil {
		fmt.Printf("PUT %s\n", putURL)
		fmt.Printf("Status Code: %d\n", statusCode)
		fmt.Printf("Response Body: %s\n", responseBody)
	} else {
		fmt.Printf("PUT request failed: %v\n", err)
	}

	// Example: DELETE request
	fmt.Println("\n--- DELETE Request Example ---")
	deleteURL := "https://jsonplaceholder.typicode.com/posts/1"
	responseBody, statusCode, err = httpClient.Delete(deleteURL)
	if err == nil {
		fmt.Printf("DELETE %s\n", deleteURL)
		fmt.Printf("Status Code: %d\n", statusCode)
		fmt.Printf("Response Body: %s\n", responseBody)
	} else {
		fmt.Printf("DELETE request failed: %v\n", err)
	}

	// Example: Custom Do request
	fmt.Println("\n--- Custom Do Request Example ---")
	customURL := "https://jsonplaceholder.typicode.com/posts"
	customBody := `{
		"title": "custom request",
		"body": "custom body",
		"userId": 1
	}`
	responseBody, statusCode, err = httpClient.Do("POST", customURL, "application/json", customBody)
	if err == nil {
		fmt.Printf("CUSTOM POST %s\n", customURL)
		fmt.Printf("Status Code: %d\n", statusCode)
		fmt.Printf("Response Body: %s\n", responseBody)
	} else {
		fmt.Printf("Custom request failed: %v\n", err)
	}

	// Example: Self-signed certificate handling
	fmt.Println("\n--- Self-signed Certificate Example ---")

	// Enable self-signed certificates
	SetAllowSelfSignedCertificates(true)

	// Update client transport to use the new setting
	httpClient.UpdateClientTransport()
	fmt.Println("Self-signed certificates enabled for future requests")
}
