// Package filepath
// File:        filepath.go
// Url:         https://github.com/xiang-tai-duo/go-bootstrap/blob/master/filepath/filepath.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: FilePath provides utility methods for file path operations, including absolute paths, cleaning, and joining
// --------------------------------------------------------------------------------
package filepath

import (
	__filepath "path/filepath"
	"strings"
)

func GetFileNameWithoutExtension(path string) string {
	var result string
	filename := __filepath.Base(path)
	extension := __filepath.Ext(filename)
	result = strings.TrimSuffix(filename, extension)
	return result
}

func GetDirectoryName(path string) string {
	var result string
	result = __filepath.Dir(path)
	return result
}
