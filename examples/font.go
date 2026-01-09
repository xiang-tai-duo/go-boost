// Package main
// File:        font.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: Example usage of font-related functions from the boost package.
// --------------------------------------------------------------------------------

package main

import (
	"fmt"
	"os"

	"boost"
)

// CustomFont is a custom implementation of FONT
// that simulates font checking for demonstration purposes
type CustomFont struct {
	boost.FONT
}

// CheckCharacter checks if a character can be displayed in the given font
// This is a custom implementation for demonstration
func (c *CustomFont) CheckCharacter(fontName string, char rune) (bool, error) {
	var err error
	// Simulate font checking
	// For demonstration, we'll say that only ASCII characters can be displayed
	return char < 128, err
}

func main() {
	// Read input from char.txt in bin directory
	charFile := "bin/char.txt"
	content, err := os.ReadFile(charFile)
	if err != nil {
		fmt.Printf("Error reading char.txt: %v\n", err)
		return
	}
	inputString := string(content)
	fontName := "Arial"

	// Example 1: Using FONT struct instance
	fmt.Println("Example 1: Using FONT struct instance")
	font := boost.NewFont()

	undisplayableChars, displayableString, err := font.CheckFontCharacters(inputString, fontName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Input string: %s\n", inputString)
		fmt.Printf("Font name: %s\n", fontName)
		fmt.Printf("Undisplayable characters: %v\n", undisplayableChars)
		fmt.Printf("Displayable string: %s\n", displayableString)
	}

	// Example 2: Get system fonts using FONT instance
	fmt.Println()
	fmt.Println("Example 2: Getting system fonts")
	systemFonts, err := font.GetSystemFonts()
	if err != nil {
		fmt.Printf("Error getting system fonts: %v\n", err)
	} else {
		fmt.Printf("System fonts found: %d\n", len(systemFonts))
		if len(systemFonts) > 0 {
			fmt.Printf("First few fonts: %v\n", systemFonts[:min(len(systemFonts), 5)])
		}
	}

	fmt.Println()

	// Example 5: Using custom font implementation
	fmt.Println("Example 5: Using custom font implementation")
	customFont := &CustomFont{}

	undisplayableChars, displayableString, err = customFont.CheckFontCharacters(inputString, fontName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Input string: %s\n", inputString)
		fmt.Printf("Font name: %s\n", fontName)
		fmt.Printf("Undisplayable characters: %v\n", undisplayableChars)
		fmt.Printf("Displayable string: %s\n", displayableString)
	}
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
