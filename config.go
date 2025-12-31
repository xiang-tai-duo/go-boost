// Package boost
// File:        config.go
// Author:      TRAE AI
// Created:     2025/12/30 11:03:46
// Description: Config provides functionality for saving and loading configuration
package boost

import (
	"encoding/json"
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

func NewConfig() *CONFIG {
	return &CONFIG{
		_data: make(map[string]interface{}),
	}
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

func (c *CONFIG) GetAll() (map[string]interface{}, error) {
	var err error
	result := make(map[string]interface{})
	c._mutex.Lock()
	defer c._mutex.Unlock()
	for k, v := range c._data {
		result[k] = v
	}
	return result, err
}

func (c *CONFIG) GetBoolean(key string) (bool, bool, error) {
	var err error
	var value bool
	var exists bool
	var raw interface{}
	if raw, exists, err = c.Get(key); err == nil && exists {
		if val, ok := raw.(bool); ok {
			value = val
		}
	}
	return value, exists, err
}

func (c *CONFIG) GetFloat(key string) (float64, bool, error) {
	var err error
	var value float64
	var exists bool
	var raw interface{}
	if raw, exists, err = c.Get(key); err == nil && exists {
		if val, ok := raw.(float64); ok {
			value = val
		} else if val, ok := raw.(int); ok {
			value = float64(val)
		} else if val, ok := raw.(int64); ok {
			value = float64(val)
		}
	}
	return value, exists, err
}

func (c *CONFIG) GetInteger(key string) (int, bool, error) {
	var err error
	var value int
	var exists bool
	var raw interface{}
	if raw, exists, err = c.Get(key); err == nil && exists {
		if val, ok := raw.(int); ok {
			value = val
		} else if val, ok := raw.(float64); ok {
			value = int(val)
		} else if val, ok := raw.(int64); ok {
			value = int(val)
		}
	}
	return value, exists, err
}

func (c *CONFIG) GetString(key string) (string, bool, error) {
	var err error
	var value string
	var exists bool
	var raw interface{}
	if raw, exists, err = c.Get(key); err == nil && exists {
		if val, ok := raw.(string); ok {
			value = val
		}
	}
	return value, exists, err
}

func (c *CONFIG) GetStringSlice(key string) ([]string, bool, error) {
	var err error
	var value []string
	var exists bool
	var raw interface{}
	if raw, exists, err = c.Get(key); err == nil && exists {
		if val, ok := raw.([]string); ok {
			value = val
		}
	}
	return value, exists, err
}

func (c *CONFIG) Load(filePath string) error {
	var err error
	var file *os.File
	if file, err = os.Open(filePath); err == nil {
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
		err = json.NewDecoder(file).Decode(&c._data)
	}
	return err
}

func (c *CONFIG) Save(filePath string) error {
	var err error
	var dir string
	var file *os.File
	if dir, err = filepath.Abs(filepath.Dir(filePath)); err == nil {
		if err = os.MkdirAll(dir, 0755); err == nil {
			if file, err = os.Create(filePath); err == nil {
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
