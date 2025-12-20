// Package file
// File:        file.go
// Url:         https://github.com/xiang-tai-duo/go-bootstrap/blob/master/file/file.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: File provides utility methods for file operations
// --------------------------------------------------------------------------------
package file

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/xiang-tai-duo/go-bootstrap/strings2"
)

//goland:noinspection GoSnakeCaseUsage
const (
	LARGE_FILE_THRESHOLD int64 = 100 * 1024 * 1024
	CHUNK_SIZE           int64 = 10 * 1024 * 1024
)

//goland:noinspection GoUnusedExportedFunction
func Abs(path string) string {
	result := ""
	if absPath, err := filepath.Abs(path); err == nil {
		result = absPath
	}
	return result
}

//goland:noinspection GoUnhandledErrorResult
func Append(path string, data []byte) error {
	err := error(nil)
	var file *os.File
	if path == "" {
		err = os.ErrInvalid
	} else if file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644); err == nil {
		defer file.Close()
		_, err = file.Write(data)
	}
	return err
}

//goland:noinspection GoUnusedExportedFunction
func AppendString(path string, content string) error {
	return Append(path, []byte(content))
}

//goland:noinspection GoUnusedExportedFunction
func CreateTemporaryFile() string {
	result := filepath.Join(os.TempDir(), strings2.Random(strings2.DEFAULT_RANDOM_SIZE))
	if file, err := os.Create(result); err == nil {
		//goland:noinspection GoUnhandledErrorResult
		file.Close()
	} else {
		result = ""
	}
	return result
}

func Copy(from string, to string) error {
	result := os.ErrInvalid
	if from != "" && to != "" {
		var fileSize int64
		if stat, err := os.Stat(from); err == nil {
			fileSize = stat.Size()
			if fileSize > LARGE_FILE_THRESHOLD {
				result = copyLargeFile(from, to, fileSize)
			} else {
				var data []byte
				var readErr error
				if data, readErr = ReadAll(from); readErr == nil {
					result = os.WriteFile(to, data, 0644)
				} else {
					result = readErr
				}
			}
		} else {
			result = err
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func IsExists(path string) bool {
	result := false
	if path != "" {
		if stat, err := os.Stat(path); err == nil && !stat.IsDir() {
			result = true
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func IsDirectory(path string) bool {
	result := false
	if path != "" {
		if stat, err := os.Stat(path); err == nil && stat.IsDir() {
			result = true
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func IsFile(path string) bool {
	result := false
	if path != "" {
		if stat, err := os.Stat(path); err == nil && !stat.IsDir() {
			result = true
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func IsRegular(path string) bool {
	result := false
	if path != "" {
		if stat, err := os.Stat(path); err == nil && stat.Mode().IsRegular() {
			result = true
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func IsSymbolLink(path string) bool {
	result := false
	if path != "" {
		if stat, err := os.Lstat(path); err == nil && stat.Mode()&os.ModeSymlink != 0 {
			result = true
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func ModifyTime(path string) time.Time {
	result := time.Time{}
	if path != "" {
		if info, err := os.Stat(path); err == nil {
			result = info.ModTime()
		}
	}
	return result
}

func ReadAll(path string) ([]byte, error) {
	result := make([]byte, 0)
	err := os.ErrInvalid
	if path != "" {
		var data []byte
		if data, err = os.ReadFile(path); err == nil {
			result = data
		}
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func ReadString(path string) string {
	result := ""
	if path != "" {
		if data, err := os.ReadFile(path); err == nil {
			result = string(data)
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func Size(path string) int64 {
	result := int64(-1)
	if path != "" {
		if stat, err := os.Stat(path); err == nil {
			result = stat.Size()
		}
	}
	return result
}

//goland:noinspection GoUnhandledErrorResult
func copyLargeFile(from string, to string, fileSize int64) error {
	err := error(nil)
	var src *os.File
	if src, err = os.Open(from); err == nil {
		defer src.Close()
		var dst *os.File
		if dst, err = os.Create(to); err == nil {
			defer dst.Close()
			var bytesCopied int64
			for bytesCopied < fileSize {
				remaining := fileSize - bytesCopied
				chunkSize := CHUNK_SIZE
				if remaining < CHUNK_SIZE {
					chunkSize = remaining
				}
				var bytesWritten int64
				if bytesWritten, err = io.CopyN(dst, src, chunkSize); err == nil {
					bytesCopied += bytesWritten
				} else {
					break
				}
			}
			if bytesCopied != fileSize {
				err = io.ErrUnexpectedEOF
			}
		}
	}
	return err
}
