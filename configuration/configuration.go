// Package configuration
// File:        configuration.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/configuration/configuration.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Configuration provides functionality for saving and loading configuration
// --------------------------------------------------------------------------------

package configuration

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/xiang-tai-duo/go-boost/aagon2"
	"github.com/xiang-tai-duo/go-boost/logger"
)

type (
	CONFIGURATION struct {
		data  map[string]interface{}
		mutex sync.Mutex
	}
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	DEFAULT_CONFIG_FILE_NAME              = "config.json"
	DEFAULT_DIR_PERM_MODE     os.FileMode = 0755
	JSON_INDENT_PREFIX                    = ""
	JSON_INDENT_STRING                    = "  "
	MODULE_NAME_CONFIGURATION             = "configuration"
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedGlobalVariable
var (
	errSecretCipherTextTooShort = errors.New("configuration: secret cipher text too short")
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_CONFIGURATION, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_CONFIGURATION, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_CONFIGURATION, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnhandledErrorResult,GoUnusedExportedFunction
func New() *CONFIGURATION {
	result := &CONFIGURATION{
		data: make(map[string]interface{}),
	}
	result.Load()
	return result
}

func (c *CONFIGURATION) Clear() error {
	err := error(nil)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data = make(map[string]interface{})
	return err
}

func (c *CONFIGURATION) Delete(key string) error {
	err := error(nil)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.data, key)
	return err
}

func (c *CONFIGURATION) Exists(key string) bool {
	result := false
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, result = c.data[key]
	return result
}

func (c *CONFIGURATION) Get(key string) (interface{}, bool, error) {
	result := interface{}(nil)
	exists := false
	err := error(nil)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	result, exists = c.data[key]
	return result, exists, err
}

func (c *CONFIGURATION) GetBoolean(key string, defaultValue ...bool) bool {
	result := false
	success := false
	if raw, exists := c.data[key]; exists {
		if val, ok := raw.(bool); ok {
			result = val
			success = true
		}
	}
	if !success && len(defaultValue) > 0 {
		result = defaultValue[0]
	}
	return result
}

//goland:noinspection DuplicatedCode
func (c *CONFIGURATION) GetFloat(key string, defaultValue ...float64) float64 {
	result := 0.0
	success := false
	if raw, exists := c.data[key]; exists {
		if val, ok := raw.(float64); ok {
			result = val
			success = true
		} else if val, ok := raw.(int); ok {
			result = float64(val)
			success = true
		} else if val, ok := raw.(int64); ok {
			result = float64(val)
			success = true
		}
	}
	if !success && len(defaultValue) > 0 {
		result = defaultValue[0]
	}
	return result
}

//goland:noinspection DuplicatedCode
func (c *CONFIGURATION) GetInteger(key string, defaultValue ...int) int {
	result := 0
	success := false
	if raw, exists := c.data[key]; exists {
		if val, ok := raw.(int); ok {
			result = val
			success = true
		} else if val, ok := raw.(float64); ok {
			result = int(val)
			success = true
		} else if val, ok := raw.(int64); ok {
			result = int(val)
			success = true
		}
	}
	if !success && len(defaultValue) > 0 {
		result = defaultValue[0]
	}
	return result
}

func (c *CONFIGURATION) GetSecretString(key string, defaultValue ...string) string {
	result := ""
	success := false
	if raw, exists := c.data[key]; exists {
		if val, ok := raw.(string); ok {
			if plainText, err := c.decryptSecret(val); err == nil {
				result = plainText
				success = true
			}
		}
	}
	if !success && len(defaultValue) > 0 {
		result = defaultValue[0]
	}
	return result
}

func (c *CONFIGURATION) GetString(key string, defaultValue ...string) string {
	result := ""
	success := false
	if raw, exists := c.data[key]; exists {
		if val, ok := raw.(string); ok {
			result = val
			success = true
		}
	}
	if !success && len(defaultValue) > 0 {
		result = defaultValue[0]
	}
	return result
}

func (c *CONFIGURATION) GetStringSlice(key string, defaultValue ...[]string) []string {
	result := make([]string, 0)
	success := false
	if raw, exists := c.data[key]; exists {
		if val, ok := raw.([]string); ok {
			result = val
			success = true
		}
	}
	if !success && len(defaultValue) > 0 {
		result = defaultValue[0]
	}
	return result
}

func (c *CONFIGURATION) IsEmpty() bool {
	result := false
	result = len(c.data) == 0
	return result
}

func (c *CONFIGURATION) Load(filePath ...*string) error {
	err := error(nil)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	configFilePath := ""
	if len(filePath) == 0 || filePath[0] == nil {
		configFilePath = DEFAULT_CONFIG_FILE_NAME
	} else {
		configFilePath = *filePath[0]
	}
	var file *os.File
	if file, err = os.Open(configFilePath); err == nil {
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
		err = json.NewDecoder(file).Decode(&c.data)
	} else {
		err = c.saveAs(configFilePath)
	}
	return err
}

func (c *CONFIGURATION) Save() error {
	return c.SaveAs(DEFAULT_CONFIG_FILE_NAME)
}

func (c *CONFIGURATION) SaveAs(filePath string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.saveAs(filePath)
}

func (c *CONFIGURATION) Set(key string, value interface{}) error {
	err := error(nil)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = value
	return err
}

func (c *CONFIGURATION) SetBoolean(key string, value bool) error {
	return c.Set(key, value)
}

func (c *CONFIGURATION) SetFloat(key string, value float64) error {
	return c.Set(key, value)
}

func (c *CONFIGURATION) SetInteger(key string, value int) error {
	return c.Set(key, value)
}

func (c *CONFIGURATION) SetSecretString(key string, value string) error {
	err := error(nil)
	var cipherText string
	if cipherText, err = c.encryptSecret(value); err == nil {
		err = c.Set(key, cipherText)
	}
	return err
}

func (c *CONFIGURATION) SetString(key string, value string) error {
	return c.Set(key, value)
}

func (c *CONFIGURATION) SetStrings(key string, value []string) error {
	return c.Set(key, value)
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_CONFIGURATION, logger.SKIP_STACK_FRAMES_BASE)
}

func (c *CONFIGURATION) decryptSecret(cipherTextBase64 string) (string, error) {
	result := ""
	err := error(nil)
	cipherData := make([]byte, 0)
	if cipherData, err = base64.StdEncoding.DecodeString(cipherTextBase64); err == nil {
		var block cipher.Block
		if block, err = aes.NewCipher(aagon2.DeriveKey()); err == nil {
			var gcm cipher.AEAD
			if gcm, err = cipher.NewGCM(block); err == nil {
				nonceSize := gcm.NonceSize()
				if len(cipherData) < nonceSize {
					err = errSecretCipherTextTooShort
				} else {
					nonce, cipherText := cipherData[:nonceSize], cipherData[nonceSize:]
					var plainText []byte
					if plainText, err = gcm.Open(nil, nonce, cipherText, nil); err == nil {
						result = string(plainText)
					}
				}
			}
		}
	}
	return result, err
}

func (c *CONFIGURATION) encryptSecret(plainText string) (string, error) {
	result := ""
	err := error(nil)
	var block cipher.Block
	if block, err = aes.NewCipher(aagon2.DeriveKey()); err == nil {
		var gcm cipher.AEAD
		if gcm, err = cipher.NewGCM(block); err == nil {
			nonce := make([]byte, gcm.NonceSize())
			if _, err = io.ReadFull(rand.Reader, nonce); err == nil {
				cipherData := gcm.Seal(nonce, nonce, []byte(plainText), nil)
				result = base64.StdEncoding.EncodeToString(cipherData)
			}
		}
	}
	return result, err
}

func (c *CONFIGURATION) saveAs(filePath string) error {
	result := error(nil)
	path := filePath
	var dir string
	if dir, result = filepath.Abs(filepath.Dir(path)); result == nil {
		if result = os.MkdirAll(dir, DEFAULT_DIR_PERM_MODE); result == nil {
			var file *os.File
			if file, result = os.Create(path); result == nil {
				defer func(file *os.File) {
					_ = file.Close()
				}(file)
				encoder := json.NewEncoder(file)
				encoder.SetIndent(JSON_INDENT_PREFIX, JSON_INDENT_STRING)
				result = encoder.Encode(c.data)
			}
		}
	}
	return result
}
