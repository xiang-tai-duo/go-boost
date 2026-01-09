// Package boost
// File:        json.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/examples/json.go
// Author:      Vibe Coding
// Created:     12/30/2025 11:03:46
// Description: Example for JSON utility functions
// --------------------------------------------------------------------------------

package main

import (
	"fmt"

	. "github.com/xiang-tai-duo/go-boost"
)

func main() {

	// Create a JSON instance from a map
	userMap := map[string]interface{}{
		"name":  "John Doe",
		"age":   30,
		"email": "john@example.com",
		"address": map[string]interface{}{
			"street": "123 Main St",
			"city":   "New York",
			"zip":    "10001",
		},
	}
	jsonObj := NewJSON(userMap)

	// Get JSON string
	jsonStr, err := jsonObj.Marshal()
	if err == nil {
		fmt.Println("Original JSON:", jsonStr)
	}

	// Format JSON with indentation
	formatted, err := jsonObj.Format("  ")
	if err == nil {
		fmt.Println("\nFormatted JSON:")
		fmt.Println(formatted)
	}

	// Minify JSON
	minified, err := jsonObj.Minify()
	if err == nil {
		fmt.Println("\nMinified JSON:", minified)
	}

	// Validate JSON
	isValid, err := jsonObj.Validate()
	if err == nil {
		fmt.Printf("\nJSON is valid: %v\n", isValid)
	}

	// Set value at path
	jsonObj.SetValue("address.country", "USA")
	jsonObj.SetInteger("age", 31)
	jsonObj.SetString("email", "john.doe@example.com")

	// Get updated JSON
	updatedJson, _ := jsonObj.Marshal()
	fmt.Println("\nUpdated JSON:", updatedJson)

	// Get specific values
	name, _ := jsonObj.GetString("name")
	age, _ := jsonObj.GetInteger("age")
	country, _ := jsonObj.GetString("address.country")
	fmt.Printf("\nRetrieved values:\n")
	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Age: %d\n", age)
	fmt.Printf("Country: %s\n", country)

	// Get all values at a path
	addressValues, _ := jsonObj.GetMap("address")
	fmt.Printf("\nAddress values: %v\n", addressValues)

	// Unmarshal JSON into a struct
	type Address struct {
		Street  string `json:"street"`
		City    string `json:"city"`
		Zip     string `json:"zip"`
		Country string `json:"country"`
	}
	type User struct {
		Name    string  `json:"name"`
		Age     int     `json:"age"`
		Email   string  `json:"email"`
		Address Address `json:"address"`
	}
	var user User
	err = jsonObj.Unmarshal(&user)
	if err == nil {
		fmt.Printf("\nUnmarshaled User:\n")
		fmt.Printf("Name: %s\n", user.Name)
		fmt.Printf("Age: %d\n", user.Age)
		fmt.Printf("Email: %s\n", user.Email)
		fmt.Printf("Address: %s, %s, %s, %s\n", user.Address.Street, user.Address.City, user.Address.Zip, user.Address.Country)
	}

	// Create JSON from string
	jsonStr2 := `{"product":"Laptop","price":999.99,"inStock":true}`
	jsonObj2 := NewJSON(jsonStr2)
	formatted2, _ := jsonObj2.Format("  ")
	fmt.Println("\nJSON from string:")
	fmt.Println(formatted2)

	// Write JSON to file
	err = jsonObj.WriteFile("./user.json", "  ")
	if err == nil {
		fmt.Println("\nJSON written to user.json")
	} else {
		fmt.Printf("Error writing to file: %v\n", err)
	}
}
