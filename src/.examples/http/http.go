// Package main
// File:        https://github.com/xiang-tai-duo/go-boost/blob/master/.examples/http/http.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: HTTP client usage example
// --------------------------------------------------------------------------------
package main

import (
	"fmt"
	"time"

	__http "github.com/xiang-tai-duo/go-boost/http"
)

//goland:noinspection GoUnhandledErrorResult
func main() {
	// Example 1: Create HTTP client
	http := __http.New()
	fmt.Println("Example 1: Created HTTP client")
	fmt.Printf("Default timeout: %v\n", http.GetTimeout())

	// Example 2: Set custom timeout
	http.SetTimeout(10 * time.Second)
	fmt.Println("\nExample 2: Set custom timeout")
	fmt.Printf("Updated timeout: %v\n", http.GetTimeout())

	// Example 3: Enable self-signed certificates (if needed)
	__http.SetAllowSelfSignedCertificates(true)
	http.UpdateClientTransport()
	fmt.Println("\nExample 3: Enabled self-signed certificates")

	// Example 4: GET request
	fmt.Println("\nExample 4: GET request")
	getURL := "https://jsonplaceholder.typicode.com/posts/1"
	var getResponse string
	var getStatusCode int
	var getErr error
	if getResponse, getStatusCode, getErr = http.Get(getURL); getErr == nil {
		fmt.Printf("GET status code: %d\n", getStatusCode)
		fmt.Println("GET response:", getResponse)
	} else {
		fmt.Printf("GET request failed: %v\n", getErr)
	}

	// Example 5: POST request
	fmt.Println("\nExample 5: POST request")
	postURL := "https://jsonplaceholder.typicode.com/posts"
	postBody := `{"title": "foo", "body": "bar", "userId": 1}`
	var postResponse string
	var postStatusCode int
	var postErr error
	if postResponse, postStatusCode, postErr = http.Post(postURL, "application/json", postBody); postErr == nil {
		fmt.Printf("POST status code: %d\n", postStatusCode)
		fmt.Println("POST response:", postResponse)
	} else {
		fmt.Printf("POST request failed: %v\n", postErr)
	}

	// Example 6: PUT request
	fmt.Println("\nExample 6: PUT request")
	putURL := "https://jsonplaceholder.typicode.com/posts/1"
	putBody := `{"id": 1, "title": "updated title", "body": "updated body", "userId": 1}`
	var putResponse string
	var putStatusCode int
	var putErr error
	if putResponse, putStatusCode, putErr = http.Put(putURL, "application/json", putBody); putErr == nil {
		fmt.Printf("PUT status code: %d\n", putStatusCode)
		fmt.Println("PUT response:", putResponse)
	} else {
		fmt.Printf("PUT request failed: %v\n", putErr)
	}

	// Example 7: DELETE request
	fmt.Println("\nExample 7: DELETE request")
	deleteURL := "https://jsonplaceholder.typicode.com/posts/1"
	var deleteResponse string
	var deleteStatusCode int
	var deleteErr error
	if deleteResponse, deleteStatusCode, deleteErr = http.Delete(deleteURL); deleteErr == nil {
		fmt.Printf("DELETE status code: %d\n", deleteStatusCode)
		fmt.Println("DELETE response:", deleteResponse)
	} else {
		fmt.Printf("DELETE request failed: %v\n", deleteErr)
	}

	// Example 8: Custom HTTP method using Do()
	fmt.Println("\nExample 8: Custom HTTP method using Do()")
	customURL := "https://jsonplaceholder.typicode.com/posts"
	customBody := `{"title": "custom", "body": "custom content", "userId": 1}`
	var customResponse string
	var customStatusCode int
	var customErr error
	if customResponse, customStatusCode, customErr = http.Do("POST", customURL, "application/json", customBody); customErr == nil {
		fmt.Printf("Custom status code: %d\n", customStatusCode)
		fmt.Println("Custom response:", customResponse)
	} else {
		fmt.Printf("Custom request failed: %v\n", customErr)
	}
}
