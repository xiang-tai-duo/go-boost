// Package boost
// File:        config.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/config.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: Config provides functionality for saving and loading configuration
// --------------------------------------------------------------------------------
package boost

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type (
	CONFIG struct {
		_mutex sync.Mutex
		_data  map[string]interface{}
	}
)

//goland:noinspection GoUnhandledErrorResult
func NewConfig() *CONFIG {
	config := &CONFIG{
		_data: make(map[string]interface{}),
	}
	config.Load()
	return config
}

func (c *CONFIG) Clear() error {
	var err error
	c._mutex.Lock()
	defer c._mutex.Unlock()
	c._data = make(map[string]interface{})
	return err
}

func (c *CONFIG) Delete(key string) error {
	var err error
	c._mutex.Lock()
	defer c._mutex.Unlock()
	delete(c._data, key)
	return err
}

func (c *CONFIG) Exists(key string) (bool, error) {
	var err error
	var exists bool
	c._mutex.Lock()
	defer c._mutex.Unlock()
	_, exists = c._data[key]
	return exists, err
}

func (c *CONFIG) Get(key string) (interface{}, bool, error) {
	var err error
	var value interface{}
	var exists bool
	c._mutex.Lock()
	defer c._mutex.Unlock()
	value, exists = c._data[key]
	return value, exists, err
}

func (c *CONFIG) GetBoolean(key string) (bool, error) {
	var value bool
	var raw interface{}
	if raw, value = c._data[key]; !value {
		return false, fmt.Errorf("key %q not found", key)
	}
	if val, ok := raw.(bool); ok {
		return val, nil
	}
	return false, fmt.Errorf("key %q is not boolean", key)
}

func (c *CONFIG) GetFloat(key string) (float64, error) {
	var raw interface{}
	if raw, _ = c._data[key]; raw == nil {
		return 0, fmt.Errorf("key %q not found", key)
	}
	if val, ok := raw.(float64); ok {
		return val, nil
	} else if val, ok := raw.(int); ok {
		return float64(val), nil
	} else if val, ok := raw.(int64); ok {
		return float64(val), nil
	}
	return 0, fmt.Errorf("key %q is not number", key)
}

func (c *CONFIG) GetInteger(key string) (int, error) {
	var raw interface{}
	if raw, _ = c._data[key]; raw == nil {
		return 0, fmt.Errorf("key %q not found", key)
	}
	if val, ok := raw.(int); ok {
		return val, nil
	} else if val, ok := raw.(float64); ok {
		return int(val), nil
	} else if val, ok := raw.(int64); ok {
		return int(val), nil
	}
	return 0, fmt.Errorf("key %q is not number", key)
}

func (c *CONFIG) GetString(key string) (string, error) {
	var raw interface{}
	if raw, _ = c._data[key]; raw == nil {
		return "", fmt.Errorf("key %q not found", key)
	}
	if val, ok := raw.(string); ok {
		return val, nil
	}
	return "", fmt.Errorf("key %q is not string", key)
}

func (c *CONFIG) GetStringSlice(key string) ([]string, error) {
	var raw interface{}
	if raw, _ = c._data[key]; raw == nil {
		return nil, fmt.Errorf("key %q not found", key)
	}
	if val, ok := raw.([]string); ok {
		return val, nil
	}
	return nil, fmt.Errorf("key %q is not string slice", key)
}

func (c *CONFIG) Load(filePath ...string) error {
	var err error
	var path string
	if len(filePath) == 0 {
		path = "config.json"
	} else {
		path = filePath[0]
	}
	var file *os.File
	if file, err = os.Open(path); err == nil {
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
		err = json.NewDecoder(file).Decode(&c._data)
	}
	return err
}

func (c *CONFIG) Save(filePath ...string) error {
	var err error
	var path string
	if len(filePath) == 0 {
		path = "config.json"
	} else {
		path = filePath[0]
	}
	var dir string
	var file *os.File
	if dir, err = filepath.Abs(filepath.Dir(path)); err == nil {
		if err = os.MkdirAll(dir, 0755); err == nil {
			if file, err = os.Create(path); err == nil {
				defer func(file *os.File) {
					_ = file.Close()
				}(file)
				encoder := json.NewEncoder(file)
				encoder.SetIndent("", "  ")
				err = encoder.Encode(c._data)
			}
		}
	}
	return err
}

func (c *CONFIG) Set(key string, value interface{}) error {
	var err error
	c._mutex.Lock()
	defer c._mutex.Unlock()
	c._data[key] = value
	return err
}

func (c *CONFIG) SetBoolean(key string, value bool) error {
	return c.Set(key, value)
}

func (c *CONFIG) SetFloat(key string, value float64) error {
	return c.Set(key, value)
}

func (c *CONFIG) SetInteger(key string, value int) error {
	return c.Set(key, value)
}

func (c *CONFIG) SetString(key string, value string) error {
	return c.Set(key, value)
}

func (c *CONFIG) SetStrings(key string, value []string) error {
	return c.Set(key, value)
}
