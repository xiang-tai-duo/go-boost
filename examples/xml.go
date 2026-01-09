// Package boost
// File:        xml.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/examples/xml.go
// Author:      Vibe Coding
// Created:     12/20/2025 12:31:58
// Description: Example for XML utility functions
// --------------------------------------------------------------------------------

package main

import (
	"encoding/xml"
	"fmt"

	. "github.com/xiang-tai-duo/go-boost"
)

func main() {

	type Address struct {
		XMLName xml.Name `xml:"address"`
		Street  string   `xml:"street"`
		City    string   `xml:"city"`
		Zip     string   `xml:"zip"`
	}
	type User struct {
		XMLName xml.Name `xml:"user"`
		Name    string   `xml:"name"`
		Age     int      `xml:"age"`
		Email   string   `xml:"email"`
		Address Address  `xml:"address"`
	}

	user := User{
		Name:  "John Doe",
		Age:   30,
		Email: "john@example.com",
		Address: Address{
			Street: "123 Main St",
			City:   "New York",
			Zip:    "10001",
		},
	}

	xmlObj := NewXML(user)

	xmlStr, err := xmlObj.Marshal()
	if err == nil {
		fmt.Println("Original XML:", xmlStr)
	}

	formatted, err := xmlObj.Format("  ")
	if err == nil {
		fmt.Println("\nFormatted XML:")
		fmt.Println(formatted)
	}

	minified, err := xmlObj.Minify()
	if err == nil {
		fmt.Println("\nMinified XML:", minified)
	}

	isValid, err := xmlObj.Validate()
	if err == nil {
		fmt.Printf("\nXML is valid: %v\n", isValid)
	}

	var parsedUser User
	err = xmlObj.Unmarshal(&parsedUser)
	if err == nil {
		fmt.Printf("\nUnmarshaled User:\n")
		fmt.Printf("Name: %s\n", parsedUser.Name)
		fmt.Printf("Age: %d\n", parsedUser.Age)
		fmt.Printf("Email: %s\n", parsedUser.Email)
		fmt.Printf("Address: %s, %s, %s\n", parsedUser.Address.Street, parsedUser.Address.City, parsedUser.Address.Zip)
	}

	err = xmlObj.WriteFile("./user.xml", "  ")
	if err == nil {
		fmt.Println("\nXML written to user.xml")
	} else {
		fmt.Printf("Error writing to file: %v\n", err)
	}
}
