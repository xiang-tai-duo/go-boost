// Package hash
// File:        hash.go
// Url:         https://github.com/xiang-tai-duo/go-bootstrap/blob/master/hash/hash.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Hash provides utility methods for hash operations, including MD5 hashing with configurable bit lengths
// --------------------------------------------------------------------------------
package hash

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

//goland:noinspection GoUnusedExportedFunction
func MD5(input string) string {
	hash := md5.Sum([]byte(input))
	return strings.ToLower(hex.EncodeToString(hash[:]))
}
