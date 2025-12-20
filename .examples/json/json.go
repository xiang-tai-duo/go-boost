// Package main
// File:        https://github.com/xiang-tai-duo/go-boost/blob/master/.examples/json/json.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: JSON utility usage example
// --------------------------------------------------------------------------------
package main

import (
	"fmt"

	__json "github.com/xiang-tai-duo/go-boost/json"
)

func main() {
	// Example 1: Create JSON from string
	jsonStr := `{"name": "John", "age": 30, "city": "New York", "contact": {"email": "john@example.com", "phone": "123-456-7890"}}`
	json1 := __json.New(jsonStr)
	fmt.Println("Example 1: JSON from string")
	fmt.Println("Original JSON:", jsonStr)

	// Example 2: Create JSON from struct
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
		City string `json:"city"`
	}
	person := Person{Name: "Alice", Age: 25, City: "London"}
	json2 := __json.New(person)
	var json2Str string
	var err error
	if json2Str, err = json2.Marshal(); err == nil {
		fmt.Println("\nExample 2: JSON from struct")
		fmt.Println("Generated JSON:", json2Str)
	} else {
		fmt.Printf("Error marshaling JSON: %v\n", err)
	}

	// Example 3: Format JSON with indentation
	var formatted string
	if formatted, err = json1.Format("  "); err == nil {
		fmt.Println("\nExample 3: Formatted JSON")
		fmt.Println(formatted)
	} else {
		fmt.Printf("Error formatting JSON: %v\n", err)
	}

	// Example 4: Minify JSON
	var minified string
	if minified, err = json1.Minify(); err == nil {
		fmt.Println("\nExample 4: Minified JSON")
		fmt.Println(minified)
	} else {
		fmt.Printf("Error minifying JSON: %v\n", err)
	}

	// Example 5: Get values from JSON
	var name string
	if name, err = json1.GetString("name"); err == nil {
		fmt.Println("\nExample 5: Get values from JSON")
		fmt.Printf("Name: %s\n", name)
	}

	var age int
	if age, err = json1.GetInteger("age"); err == nil {
		fmt.Printf("Age: %d\n", age)
	}

	var email string
	if email, err = json1.GetString("contact.email"); err == nil {
		fmt.Printf("Email: %s\n", email)
	}

	// Example 6: Set values in JSON
	json1.SetString("name", "Bob")
	json1.SetInteger("age", 35)
	json1.SetString("contact.email", "bob@example.com")

	var updatedJSON string
	if updatedJSON, err = json1.Marshal(); err == nil {
		fmt.Println("\nExample 6: Updated JSON after setting values")
		fmt.Println(updatedJSON)
	} else {
		fmt.Printf("Error marshaling updated JSON: %v\n", err)
	}

	// Example 7: Validate JSON
	var valid bool
	if valid, err = json1.Validate(); err == nil {
		fmt.Println("\nExample 7: Validate JSON")
		fmt.Printf("Is valid: %v, Error: %v\n", valid, err)
	} else {
		fmt.Println("\nExample 7: Validate JSON")
		fmt.Printf("Is valid: %v, Error: %v\n", valid, err)
	}

	// Example 8: Write JSON to file
	if err = json1.WriteFile("person.json", "  "); err == nil {
		fmt.Println("\nExample 8: Write JSON to file")
		fmt.Println("JSON written to person.json successfully")
	} else {
		fmt.Printf("Error writing JSON to file: %v\n", err)
	}

	// Example 9: Unmarshal JSON to struct
	var personFromJSON Person
	if err = json1.Unmarshal(&personFromJSON); err == nil {
		fmt.Println("\nExample 9: Unmarshal JSON to struct")
		fmt.Printf("Unmarshaled person: %+v\n", personFromJSON)
	} else {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
	}
}
