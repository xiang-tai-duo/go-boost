// Package config
// File:        config.go
// Url:         https://github.com/xiang-tai-duo/go-bootstrap/blob/master/config/config.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Config provides functionality for saving and loading configuration
// --------------------------------------------------------------------------------
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type (
	Config struct {
		mutex sync.Mutex
		data  map[string]interface{}
	}
)

//goland:noinspection GoSnakeCaseUsage
const (
	DEFAULT_CONFIG_FILE_NAME = "config.json"
)

//goland:noinspection GoUnhandledErrorResult,GoUnusedExportedFunction
func New() *Config {
	config := &Config{
		data: make(map[string]interface{}),
	}
	config.Load()
	return config
}

func (c *Config) Clear() error {
	err := error(nil)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data = make(map[string]interface{})
	return err
}

func (c *Config) Delete(key string) error {
	err := error(nil)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.data, key)
	return err
}

func (c *Config) Exists(key string) bool {
	exists := false
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, exists = c.data[key]
	return exists
}

func (c *Config) Get(key string) (interface{}, bool, error) {
	err := error(nil)
	value := interface{}(nil)
	exists := false
	c.mutex.Lock()
	defer c.mutex.Unlock()
	value, exists = c.data[key]
	return value, exists, err
}

func (c *Config) GetBoolean(key string, defaultValue ...bool) bool {
	result := false
	raw := interface{}(nil)
	exists := false
	success := false
	if raw, exists = c.data[key]; exists {
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
func (c *Config) GetFloat(key string, defaultValue ...float64) float64 {
	result := 0.0
	raw := interface{}(nil)
	exists := false
	success := false
	if raw, exists = c.data[key]; exists {
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
func (c *Config) GetInteger(key string, defaultValue ...int) int {
	result := 0
	raw := interface{}(nil)
	exists := false
	success := false
	if raw, exists = c.data[key]; exists {
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

func (c *Config) GetString(key string, defaultValue ...string) string {
	result := ""
	raw := interface{}(nil)
	exists := false
	success := false
	if raw, exists = c.data[key]; exists {
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

func (c *Config) GetStringSlice(key string, defaultValue ...[]string) []string {
	result := []string(nil)
	raw := interface{}(nil)
	exists := false
	success := false
	if raw, exists = c.data[key]; exists {
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

func (c *Config) IsEmpty() bool {
	return len(c.data) == 0
}

func (c *Config) Load(filePath ...interface{}) error {
	err := error(nil)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	configFilePath := ""
	if len(filePath) == 0 {
		configFilePath = DEFAULT_CONFIG_FILE_NAME
	} else {
		firstArg := filePath[0]
		switch v := firstArg.(type) {
		case string:
			configFilePath = v
		}
	}
	file := (*os.File)(nil)
	if file, err = os.Open(configFilePath); err == nil {
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
		err = json.NewDecoder(file).Decode(&c.data)
	} else {
		c.Save(configFilePath)
		err = nil
	}
	return err
}

func (c *Config) Save(filePath ...string) {
	go func() {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		err := error(nil)
		path := ""
		if len(filePath) == 0 {
			path = "config.json"
		} else {
			path = filePath[0]
		}
		dir := ""
		file := (*os.File)(nil)
		if dir, err = filepath.Abs(filepath.Dir(path)); err == nil {
			if err = os.MkdirAll(dir, 0755); err == nil {
				if file, err = os.Create(path); err == nil {
					defer func(file *os.File) {
						_ = file.Close()
					}(file)
					encoder := json.NewEncoder(file)
					encoder.SetIndent("", "  ")
					_ = encoder.Encode(c.data)
				}
			}
		}
	}()
}

func (c *Config) Set(key string, value interface{}) error {
	err := error(nil)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = value
	return err
}

func (c *Config) SetBoolean(key string, value bool) error {
	return c.Set(key, value)
}

func (c *Config) SetFloat(key string, value float64) error {
	return c.Set(key, value)
}

func (c *Config) SetInteger(key string, value int) error {
	return  c.Set(key, value)
}

func (c *Config) SetString(key string, value string) error {
	return c.Set(key, value)
}

func (c *Config) SetStrings(key string, value []string) error {
	return c.Set(key, value)
}
