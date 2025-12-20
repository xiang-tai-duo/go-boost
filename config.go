// --------------------------------------------------------------------------------
// File:        config.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: Config provides functionality for saving and loading configuration
//              data to and from files.
// --------------------------------------------------------------------------------

package boost

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// CONFIG represents a configuration store that can be loaded from and saved to files.
type CONFIG struct {
	_mutex sync.Mutex
	_data  map[string]interface{}
}

// NewConfig creates a new CONFIG instance with default configurations.
// Usage:
// config := NewConfig()
func NewConfig() *CONFIG {
	return &CONFIG{
		_data: make(map[string]interface{}),
	}
}

// Clear removes all configuration entries.
// Returns: error if any occurred
// Usage:
// err := config.Clear()
func (c *CONFIG) Clear() error {
	var err error
	c._mutex.Lock()
	defer c._mutex.Unlock()
	c._data = make(map[string]interface{})
	return err
}

// Delete removes a configuration entry by key.
// key: Configuration key to delete
// Returns: error if any occurred
// Usage:
// err := config.Delete("database.host")
func (c *CONFIG) Delete(key string) error {
	var err error
	c._mutex.Lock()
	defer c._mutex.Unlock()
	delete(c._data, key)
	return err
}

// Exists checks if a configuration key exists.
// key: Configuration key to check
// Returns: bool indicating if the key exists, error if any occurred
// Usage:
// exists, err := config.Exists("database.host")
func (c *CONFIG) Exists(key string) (bool, error) {
	var err error
	var exists bool
	c._mutex.Lock()
	defer c._mutex.Unlock()
	_, exists = c._data[key]
	return exists, err
}

// Get retrieves a configuration value by key.
// key: Configuration key to retrieve
// Returns: interface{} value, bool indicating if the key exists, error if any occurred
// Usage:
// value, exists, err := config.Get("database.host")
func (c *CONFIG) Get(key string) (interface{}, bool, error) {
	var err error
	var value interface{}
	var exists bool
	c._mutex.Lock()
	defer c._mutex.Unlock()
	value, exists = c._data[key]
	return value, exists, err
}

// GetAll retrieves all configuration entries.
// Returns: map of all configuration entries, error if any occurred
// Usage:
// allConfig, err := config.GetAll()
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

// GetBoolean retrieves a configuration value as a boolean.
// key: Configuration key to retrieve
// Returns: bool value, bool indicating if the key exists, error if any occurred
// Usage:
// value, exists, err := config.GetBoolean("feature.enabled")
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

// GetFloat retrieves a configuration value as a float64.
// key: Configuration key to retrieve
// Returns: float64 value, bool indicating if the key exists, error if any occurred
// Usage:
// value, exists, err := config.GetFloat("server.timeout")
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

// GetInteger retrieves a configuration value as an int.
// key: Configuration key to retrieve
// Returns: int value, bool indicating if the key exists, error if any occurred
// Usage:
// value, exists, err := config.GetInteger("server.port")
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

// GetString retrieves a configuration value as a string.
// key: Configuration key to retrieve
// Returns: string value, bool indicating if the key exists, error if any occurred
// Usage:
// value, exists, err := config.GetString("database.host")
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

// GetStringSlice retrieves a configuration value as a slice of strings.
// key: Configuration key to retrieve
// Returns: []string value, bool indicating if the key exists, error if any occurred
// Usage:
// value, exists, err := config.GetStringSlice("server.allowedOrigins")
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

// Load loads configuration from a JSON file.
// filePath: Path to the JSON configuration file, can be just a filename to use current directory
// Returns: error if any occurred
// Usage:
// err := config.Load("./config.json")
// err := config.Load("config.json") // Uses current directory
func (c *CONFIG) Load(filePath string) error {
	var err error
	var file *os.File
	if file, err = os.Open(filePath); err == nil {
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
		if err = json.NewDecoder(file).Decode(&c._data); err == nil {
			// Load successful
		}
	}
	return err
}

// Save saves configuration to a JSON file.
// filePath: Path to the JSON configuration file, can be just a filename to use current directory
// Returns: error if any occurred
// Usage:
// err := config.Save("./config.json")
// err := config.Save("config.json") // Uses current directory
func (c *CONFIG) Save(filePath string) error {
	var err error
	var dir string
	if dir, err = filepath.Abs(filepath.Dir(filePath)); err == nil {
		if err = os.MkdirAll(dir, 0755); err == nil {
			var file *os.File
			if file, err = os.Create(filePath); err == nil {
				defer func(file *os.File) {
					_ = file.Close()
				}(file)
				encoder := json.NewEncoder(file)
				encoder.SetIndent("", "  ")
				if err = encoder.Encode(c._data); err == nil {
					// Save successful
				}
			}
		}
	}
	return err
}

// Set adds or updates a configuration entry.
// key: Configuration key to set
// value: Configuration value to set
// Returns: error if any occurred
// Usage:
// err := config.Set("database.host", "localhost")
func (c *CONFIG) Set(key string, value interface{}) error {
	var err error
	c._mutex.Lock()
	defer c._mutex.Unlock()
	c._data[key] = value
	return err
}

// SetBoolean adds or updates a boolean configuration entry.
// key: Configuration key to set
// value: Boolean value to set
// Returns: error if any occurred
// Usage:
// err := config.SetBoolean("feature.enabled", true)
func (c *CONFIG) SetBoolean(key string, value bool) error {
	return c.Set(key, value)
}

// SetFloat adds or updates a float64 configuration entry.
// key: Configuration key to set
// value: Float64 value to set
// Returns: error if any occurred
// Usage:
// err := config.SetFloat("server.timeout", 30.5)
func (c *CONFIG) SetFloat(key string, value float64) error {
	return c.Set(key, value)
}

// SetInteger adds or updates an integer configuration entry.
// key: Configuration key to set
// value: Integer value to set
// Returns: error if any occurred
// Usage:
// err := config.SetInteger("server.port", 8080)
func (c *CONFIG) SetInteger(key string, value int) error {
	return c.Set(key, value)
}

// SetString adds or updates a string configuration entry.
// key: Configuration key to set
// value: String value to set
// Returns: error if any occurred
// Usage:
// err := config.SetString("database.host", "localhost")
func (c *CONFIG) SetString(key string, value string) error {
	return c.Set(key, value)
}

// SetStrings adds or updates a slice of strings configuration entry.
// key: Configuration key to set
// value: Slice of strings to set
// Returns: error if any occurred
// Usage:
// err := config.SetStrings("server.allowedOrigins", []string{"http://localhost", "http://example.com"})
func (c *CONFIG) SetStrings(key string, value []string) error {
	return c.Set(key, value)
}
