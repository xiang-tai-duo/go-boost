// Package strings2
// File:        strings2.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/strings/strings2.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: STRINGS is a wrapper for string operations, providing a set of methods for string manipulation.
// --------------------------------------------------------------------------------
package strings2

import (
	cryptoRand "crypto/rand"
	"math/rand"
	"time"
	"unicode"

	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection SpellCheckingInspection,GoSnakeCaseUsage,GoUnusedConst
const (
	DEFAULT_RANDOM_SIZE  = 8
	LETTERS_LITE         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LETTERS              = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	MODULE_NAME_STRINGS2 = "strings2"
)

var (
	RND = rand.New(rand.NewSource(time.Now().UnixNano()))
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_STRINGS2, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_STRINGS2, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_STRINGS2, logger.SKIP_STACK_FRAMES_BASE)
}

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
func Itoa(value int) string {
	result := ""
	negative := false
	num := value
	if num < 0 {
		negative = true
		num = -num
	}
	if num == 0 {
		result = "0"
	} else {
		buf := [32]byte{}
		pos := 31
		for num > 0 {
			buf[pos] = byte('0' + num%10)
			pos--
			num /= 10
		}
		if negative {
			pos--
			buf[pos] = '-'
		}
		result = string(buf[pos:32])
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
	buffer := make([]byte, randomLength)
	if _, err := cryptoRand.Read(buffer); err == nil {
		result[0] = LETTERS_LITE[int(buffer[0])%len(LETTERS_LITE)]
		for i := 1; i < randomLength; i++ {
			result[i] = LETTERS[int(buffer[i])%len(LETTERS)]
		}
	} else {
		result[0] = LETTERS_LITE[RND.Intn(len(LETTERS_LITE))]
		for i := 1; i < randomLength; i++ {
			result[i] = LETTERS[RND.Intn(len(LETTERS))]
		}
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

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_STRINGS2, logger.SKIP_STACK_FRAMES_BASE)
}
