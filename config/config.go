// Package config
// File:        config.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/config/config.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: Config provides functionality for saving and loading configuration
// --------------------------------------------------------------------------------
package config

import (
	"encoding/json"
	"fmt"
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

func (c *Config) Exists(key string) (bool, error) {
	err := error(nil)
	exists := false
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, exists = c.data[key]
	return exists, err
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

func (c *Config) GetBoolean(key string) (bool, error) {
	result := false
	err := error(nil)
	raw := interface{}(nil)
	exists := false
	if raw, exists = c.data[key]; exists {
		if val, ok := raw.(bool); ok {
			result = val
			err = nil
		} else {
			result = false
			err = fmt.Errorf("key %q is not boolean", key)
		}
	} else {
		result = false
		err = fmt.Errorf("key %q not found", key)
	}
	return result, err
}

func (c *Config) GetFloat(key string) (float64, error) {
	result := 0.0
	err := error(nil)
	raw := interface{}(nil)
	exists := false
	if raw, exists = c.data[key]; exists {
		if val, ok := raw.(float64); ok {
			result = val
			err = nil
		} else if val, ok := raw.(int); ok {
			result = float64(val)
			err = nil
		} else if val, ok := raw.(int64); ok {
			result = float64(val)
			err = nil
		} else {
			result = 0
			err = fmt.Errorf("key %q is not number", key)
		}
	} else {
		result = 0
		err = fmt.Errorf("key %q not found", key)
	}
	return result, err
}

func (c *Config) GetInteger(key string) (int, error) {
	result := 0
	err := error(nil)
	raw := interface{}(nil)
	exists := false
	if raw, exists = c.data[key]; exists {
		if val, ok := raw.(int); ok {
			result = val
			err = nil
		} else if val, ok := raw.(float64); ok {
			result = int(val)
			err = nil
		} else if val, ok := raw.(int64); ok {
			result = int(val)
			err = nil
		} else {
			result = 0
			err = fmt.Errorf("key %q is not number", key)
		}
	} else {
		result = 0
		err = fmt.Errorf("key %q not found", key)
	}
	return result, err
}

func (c *Config) GetString(key string) (string, error) {
	result := ""
	err := error(nil)
	raw := interface{}(nil)
	exists := false
	if raw, exists = c.data[key]; exists {
		if val, ok := raw.(string); ok {
			result = val
			err = nil
		} else {
			result = ""
			err = fmt.Errorf("key %q is not string", key)
		}
	} else {
		result = ""
		err = fmt.Errorf("key %q not found", key)
	}
	return result, err
}

func (c *Config) GetStringSlice(key string) ([]string, error) {
	result := []string(nil)
	err := error(nil)
	raw := interface{}(nil)
	exists := false
	if raw, exists = c.data[key]; exists {
		if val, ok := raw.([]string); ok {
			result = val
			err = nil
		} else {
			result = nil
			err = fmt.Errorf("key %q is not string slice", key)
		}
	} else {
		result = nil
		err = fmt.Errorf("key %q not found", key)
	}
	return result, err
}

func (c *Config) Load(filePath ...interface{}) error {
	err := error(nil)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if len(filePath) == 0 {
		file := (*os.File)(nil)
		if file, err = os.Open(DEFAULT_CONFIG_FILE_NAME); err == nil {
			defer func(file *os.File) {
				_ = file.Close()
			}(file)
			err = json.NewDecoder(file).Decode(&c.data)
		}
	} else {
		firstArg := filePath[0]
		switch v := firstArg.(type) {
		case string:
			path := v
			file := (*os.File)(nil)
			if file, err = os.Open(path); err == nil {
				defer func(file *os.File) {
					_ = file.Close()
				}(file)
				err = json.NewDecoder(file).Decode(&c.data)
			}
		default:
			if jsonData, err := json.Marshal(v); err == nil {
				err = json.Unmarshal(jsonData, &c.data)
			}
		}
	}

	return err
}

func (c *Config) Save(filePath ...string) error {
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
				err = encoder.Encode(c.data)
			}
		}
	}
	return err
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
	return c.Set(key, value)
}

func (c *Config) SetString(key string, value string) error {
	return c.Set(key, value)
}

func (c *Config) SetStrings(key string, value []string) error {
	return c.Set(key, value)
}
