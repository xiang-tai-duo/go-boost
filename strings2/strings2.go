// Package strings2
// File:        strings2.go
// Url:         https://github.com/xiang-tai-duo/go-bootstrap/blob/master/strings/strings2.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: STRINGS is a wrapper for string operations, providing a set of methods for string manipulation.
// --------------------------------------------------------------------------------
package strings2

import (
	"math/rand"
	"time"
	"unicode"
)

//goland:noinspection SpellCheckingInspection,GoSnakeCaseUsage,GoUnusedConst
const (
	DEFAULT_RANDOM_SIZE = 8
	EMPTY               = ""
	LETTERS_LITE        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LETTERS             = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

var (
	RND = rand.New(rand.NewSource(time.Now().UnixNano()))
)

//goland:noinspection SpellCheckingInspection,GoUnusedExportedFunction
func Atoi(value string) int {
	result := 0
	sign := 1
	index := 0
	for index < len(value) && unicode.IsSpace(rune(value[index])) {
		index++
	}
	if index < len(value) {
		switch value[index] {
		case '+':
			index++
		case '-':
			sign = -1
			index++
		}
	}
	for index < len(value) && unicode.IsDigit(rune(value[index])) {
		result = result*10 + int(value[index]-'0')
		index++
	}
	if index < len(value) && !unicode.IsSpace(rune(value[index])) {
		result = 0
	} else {
		result = sign * result
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func Left(value string, length int) string {
	result := ""
	if length >= len(value) {
		result = value
	} else {
		result = value[:length]
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func Mid(value string, start int, length int) string {
	result := ""
	if start < 0 {
		start = 0
	}
	if start >= len(value) {
		result = ""
	} else {
		end := start + length
		if end > len(value) {
			end = len(value)
		}
		result = value[start:end]
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func Random(args ...interface{}) string {
	result := make([]byte, 0)
	randomLength := DEFAULT_RANDOM_SIZE
	if len(args) > 0 {
		if length, ok := args[0].(int); ok {
			randomLength = length
			if randomLength <= 0 {
				randomLength = DEFAULT_RANDOM_SIZE
			}
		}
	}
	result = make([]byte, randomLength)
	result[0] = LETTERS_LITE[RND.Intn(len(LETTERS_LITE))]
	for i := 1; i < randomLength; i++ {
		result[i] = LETTERS[RND.Intn(len(LETTERS))]
	}
	return string(result)
}

//goland:noinspection GoUnusedExportedFunction
func Right(value string, length int) string {
	result := ""
	if length >= len(value) {
		result = value
	} else {
		result = value[len(value)-length:]
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func SubString(value string, start int) string {
	result := ""
	if start < 0 {
		start = 0
	}
	if start >= len(value) {
		result = ""
	} else {
		result = value[start:]
	}
	return result
}
