// --------------------------------------------------------------------------------
// File:        string.go
// Author:      TRAE AI
// Created:     12/30/2025 11:03:46
// Description: Example for STRING utility functions
// --------------------------------------------------------------------------------

package main

import (
	"fmt"
	"unicode"

	. "github.com/xiang-tai-duo/go-boost"
)

func main() {

	// Create a new STRING instance
	var str STRING
	str.String("Hello, Go _! This is a test string.")

	// Basic string operations
	fmt.Println("Original string:", str.String())
	fmt.Println("Length:", len(str.String()))

	// Contains operations
	fmt.Println("Contains 'Go':", str.Contains("Go"))
	fmt.Println("Contains 'Java':", str.Contains("Java"))
	fmt.Println("ContainsAny 'xyz':", str.ContainsAny("xyz"))
	fmt.Println("ContainsRune '!':", str.ContainsRune('!'))

	// Count occurrences
	fmt.Println("Count 'is':", str.Count("is"))

	// Cut operations
	before, after, found := str.Cut(",")
	fmt.Printf("Cut by ',': found=%v, before='%s', after='%s'\n", found, before.String(), after.String())

	// Cut prefix/suffix
	var withPrefix STRING
	withPrefix.String("Prefix_Text")
	result, found := withPrefix.CutPrefix("Prefix_")
	fmt.Printf("CutPrefix 'Prefix_': found=%v, result='%s'\n", found, result.String())
	var withSuffix STRING
	withSuffix.String("Text_Suffix")
	result, found = withSuffix.CutSuffix("_Suffix")
	fmt.Printf("CutSuffix '_Suffix': found=%v, result='%s'\n", found, result.String())

	// Case operations
	toLowerResult := str.ToLower()
	fmt.Println("ToLower:", toLowerResult.String())
	toUpperResult := str.ToUpper()
	fmt.Println("ToUpper:", toUpperResult.String())
	titleResult := str.Title()
	fmt.Println("Title:", titleResult.String())
	fmt.Println("EqualFold 'hello, go boost!':", str.EqualFold("hello, go boost!"))

	// Fields and splitting
	fmt.Println("Fields:", str.Fields())
	fmt.Println("Split by ' ':", str.Split(" "))
	fmt.Println("SplitN by ' ' (3 parts):", str.SplitN(" ", 3))

	// Index operations
	fmt.Println("Index '_':", str.Index("_"))
	fmt.Println("LastIndex 'is':", str.LastIndex("is"))
	fmt.Println("IndexByte 'o':", str.IndexByte('o'))

	// Prefix/suffix checks
	fmt.Println("HasPrefix 'Hello':", str.HasPrefix("Hello"))
	fmt.Println("HasSuffix '.':", str.HasSuffix("."))

	// Substring operations
	leftResult := str.Left(5)
	fmt.Println("Left 5:", leftResult.String())
	rightResult := str.Right(10)
	fmt.Println("Right 10:", rightResult.String())
	midResult := str.Mid(7, 5)
	fmt.Println("Mid 7, 5:", midResult.String())
	subStrResult := str.SubStr(7)
	fmt.Println("SubStr 7:", subStrResult.String())

	// Trim operations
	var whitespaceStr STRING
	whitespaceStr.String("   Trim me   ")
	trimSpaceResult := whitespaceStr.TrimSpace()
	fmt.Printf("Trim whitespace: '%s' -> '%s'\n", whitespaceStr.String(), trimSpaceResult.String())
	var xxStr STRING
	xxStr.String("xxTrim mexx")
	trimmedX := xxStr.Trim("x")
	fmt.Printf("Trim 'x': '%s' -> '%s'\n", xxStr.String(), trimmedX.String())

	// Replace operations
	replaceResult := str.Replace("a", "A", -1)
	fmt.Println("Replace 'a' with 'A':", replaceResult.String())
	replaceAllResult := str.ReplaceAll("is", "was")
	fmt.Println("ReplaceAll 'is' with 'was':", replaceAllResult.String())

	// Join operation
	var joinStr STRING
	joinStr.String(", ")
	joined := joinStr.Join([]string{"apple", "banana", "cherry"})
	fmt.Println("Join with ', ': ", joined.String())

	// Repeat operation
	var repeatStr STRING
	repeatStr.String("ab")
	repeated := repeatStr.Repeat(3)
	fmt.Println("Repeat 3 times:", repeated.String())

	// Map operation
	var mapStr STRING
	mapStr.String("Hello 123")
	mapped := mapStr.Map(func(r rune) rune {
		if unicode.IsDigit(r) {
			return '*'
		}
		return r
	})
	fmt.Printf("Map digits to '*': '%s' -> '%s'\n", mapStr.String(), mapped.String())

	// FieldsFunc
	var fieldsFuncStr STRING
	fieldsFuncStr.String("Hello,123-World")
	fields := fieldsFuncStr.FieldsFunc(func(r rune) bool {
		return !unicode.IsLetter(r)
	})
	fmt.Println("FieldsFunc (letters only):", fields)

	// Trim functions
	var trimFuncStr STRING
	trimFuncStr.String("!!Hello!!")
	trimmed := trimFuncStr.TrimFunc(func(r rune) bool {
		return r == '!'
	})
	fmt.Printf("TrimFunc '!': '%s' -> '%s'\n", trimFuncStr.String(), trimmed.String())
}
