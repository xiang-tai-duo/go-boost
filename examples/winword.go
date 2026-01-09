// Package main
// File:        winword.go
// Url:         `https://github.com/xiang-tai-duo/go-boost/blob/master/examples/winword.go`
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: Example for Word document operations, demonstrating how to read text from a DOCX file.
// --------------------------------------------------------------------------------

package main

import (
	"fmt"

	"github.com/xiang-tai-duo/go-boost"
)

func main() {
	// Create a new WORD_DOCUMENT instance
	wordDocument := boost.NewWordDocument()

	// Path to the DOCX file
	docxPath := "bin/word.docx"

	// Read text from the DOCX file
	if err := wordDocument.Load(docxPath); err == nil {
		// Get the extracted text
		text := wordDocument.Word.Document.Text()

		// Print the extracted text
		fmt.Println("Extracted text from DOCX file:")
		fmt.Println("-----------------------------------")
		fmt.Println(text)
		fmt.Println("-----------------------------------")
		fmt.Printf("Total characters: %d\n", len(text))
	} else {
		fmt.Printf("Error reading DOCX file: %v\n", err)
		return
	}
}
