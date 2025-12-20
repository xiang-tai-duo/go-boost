// Package hash
// File:        hash.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/hash/hash.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Hash provides utility methods for hash operations, including MD5 hashing with configurable bit lengths
// --------------------------------------------------------------------------------

package hash

import (
	"crypto/md5"
	"crypto/sha3"
	"encoding/hex"
	"github.com/xiang-tai-duo/go-boost/logger"
	"io"
	"os"
	"strings"
)

//goland:noinspection GoUnusedConst,GoSnakeCaseUsage
const (
	MODULE_NAME_HASH        = "hash"
	RANDOM_IDENTIFIER_BYTES = 16
)

//goland:noinspection GoUnusedGlobalVariable,GoSnakeCaseUsage
var (
	CLIENT_PREFIX = "go-boost-"
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_HASH, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_HASH, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_HASH, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func GetFileMD5(filePath string) (string, error) {
	result := ""
	err := error(nil)
	var file *os.File
	if file, err = os.Open(filePath); err == nil {
		defer file.Close()
		hash := md5.New()
		if _, err = io.Copy(hash, file); err == nil {
			result = strings.ToLower(hex.EncodeToString(hash.Sum(nil)))
		}
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func GetFileSHA3(filePath string) (string, error) {
	result := ""
	err := error(nil)
	var file *os.File
	if file, err = os.Open(filePath); err == nil {
		defer file.Close()
		hash := sha3.New256()
		if _, err = io.Copy(hash, file); err == nil {
			result = strings.ToLower(hex.EncodeToString(hash.Sum(nil)))
		}
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func SHA3(input string) string {
	result := ""
	hash := sha3.Sum256([]byte(input))
	result = strings.ToLower(hex.EncodeToString(hash[:]))
	return result
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_HASH, logger.SKIP_STACK_FRAMES_BASE)
}
