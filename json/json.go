// Package json
// File:        json.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/json/json.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: JSON provides utility methods for JSON operations, including marshal, unmarshal, validation, formatting, minification, and file operations.
// --------------------------------------------------------------------------------
package json

import (
	_json "encoding/json"
	"errors"
	"os"
	"strings"
)

type (
	JSON struct {
		value interface{}
		text  string
	}
)

//goland:noinspection DuplicatedCode
func New(json ...interface{}) *JSON {
	j := &JSON{}
	if len(json) > 0 {
		j.value = json[0]
		if s, ok := json[0].(string); ok {
			var parsedValue interface{}
			if err := _json.Unmarshal([]byte(s), &parsedValue); err == nil {
				j.text = s
			} else if jsonBytes, err := _json.Marshal(json[0]); err == nil {
				j.text = string(jsonBytes)
			}
		} else if jsonBytes, err := _json.Marshal(json[0]); err == nil {
			j.text = string(jsonBytes)
		}
	}
	return j
}

func (j *JSON) Format(indent string) (string, error) {
	var err error
	var formatted []byte
	var parsedValue interface{}
	var result string
	if err = _json.Unmarshal([]byte(j.text), &parsedValue); err == nil {
		if formatted, err = _json.MarshalIndent(parsedValue, "", indent); err == nil {
			result = string(formatted)
		}
	}
	return result, err
}

func (j *JSON) GetInteger(path string) (int, error) {
	var err error
	var value interface{}
	var result int
	if value, err = j.GetValue(path); err == nil {
		if f, ok := value.(float64); ok {
			result = int(f)
		} else if i, ok := value.(int); ok {
			result = i
		} else {
			err = errors.New("value is not an integer")
		}
	}
	return result, err
}

//goland:noinspection DuplicatedCode
func (j *JSON) GetMap(path string) (map[string]interface{}, error) {
	var parsedValue interface{}
	var err error
	var result map[string]interface{}
	if err = _json.Unmarshal([]byte(j.text), &parsedValue); err == nil {
		var nestedValue interface{}
		if nestedValue, err = getNestedValue(parsedValue, path); err == nil {
			if m, ok := nestedValue.(map[string]interface{}); ok {
				result = make(map[string]interface{})
				for key, value := range m {
					result[key] = value
				}
			} else {
				err = errors.New("specified path does not point to a map")
			}
		}
	}
	return result, err
}

func (j *JSON) GetString(path string) (string, error) {
	var err error
	var value interface{}
	var result string
	if value, err = j.GetValue(path); err == nil {
		if s, ok := value.(string); ok {
			result = s
		} else {
			err = errors.New("value is not a string")
		}
	}

	return result, err
}

func (j *JSON) GetValue(path string) (interface{}, error) {
	var err error
	var parsedValue interface{}
	var result interface{}
	if err = _json.Unmarshal([]byte(j.text), &parsedValue); err == nil {
		result, err = getNestedValue(parsedValue, path)
	}
	return result, err
}

func (j *JSON) Marshal() (string, error) {
	return j.text, nil
}

func (j *JSON) Minify() (string, error) {
	var err error
	var minified []byte
	var parsedValue interface{}
	var result string
	if err = _json.Unmarshal([]byte(j.text), &parsedValue); err == nil {
		if minified, err = _json.Marshal(parsedValue); err == nil {
			result = string(minified)
		}
	}
	return result, err
}

func (j *JSON) SetInteger(path string, value int) (*JSON, error) {
	return j.SetValue(path, value)
}

func (j *JSON) SetString(path string, value string) (*JSON, error) {
	return j.SetValue(path, value)
}

//goland:noinspection DuplicatedCode
func (j *JSON) SetValue(path string, value interface{}) (*JSON, error) {
	var err error
	var parsedValue interface{}
	if err = _json.Unmarshal([]byte(j.text), &parsedValue); err == nil {
		if err = setNestedValue(parsedValue, path, value); err == nil {
			var jsonBytes []byte
			if jsonBytes, err = _json.Marshal(parsedValue); err == nil {
				j.text = string(jsonBytes)
				j.value = parsedValue
			}
		}
	}
	return j, err
}

func (j *JSON) Unmarshal(target interface{}) error {
	return _json.Unmarshal([]byte(j.text), target)
}

func (j *JSON) Validate() (bool, error) {
	var parsedValue interface{}
	var result bool
	var err error
	if err = _json.Unmarshal([]byte(j.text), &parsedValue); err == nil {
		result = true
	}
	return result, err
}

func (j *JSON) WriteFile(filePath string, indent ...string) error {
	var err error
	var file *os.File
	var jsonBytes []byte
	var indentValue string
	if len(indent) > 0 {
		indentValue = indent[0]
	}
	if indentValue != "" {
		if jsonBytes, err = _json.MarshalIndent(j.value, "", indentValue); err == nil {
			// No return here, continue
		} else {
			return err
		}
	} else {
		if jsonBytes, err = _json.Marshal(j.value); err == nil {
			// No return here, continue
		} else {
			return err
		}
	}
	if file, err = os.Create(filePath); err == nil {
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
		_, err = file.Write(jsonBytes)
	}
	return err
}

//goland:noinspection DuplicatedCode
func getNestedValue(data interface{}, path string) (interface{}, error) {
	parts := strings.Split(path, ".")
	current := data
	var result interface{}
	var err error
	for _, part := range parts {
		switch m := current.(type) {
		case map[string]interface{}:
			if value, ok := m[part]; ok {
				current = value
			} else {
				err = errors.New("path not found")
				break
			}
		default:
			err = errors.New("cannot navigate path: parent is not a map")
			break
		}
		if err != nil {
			break
		}
	}
	if err == nil {
		result = current
	}
	return result, err
}

//goland:noinspection DuplicatedCode
func setNestedValue(data interface{}, path string, value interface{}) error {
	parts := strings.Split(path, ".")
	current := data
	var err error
	for i, part := range parts {
		if i == len(parts)-1 {
			switch m := current.(type) {
			case map[string]interface{}:
				m[part] = value
			default:
				err = errors.New("cannot set value at path: parent is not a map")
				break
			}
		} else {
			switch m := current.(type) {
			case map[string]interface{}:
				if _, ok := m[part]; !ok {
					m[part] = make(map[string]interface{})
				}
				current = m[part]
			default:
				err = errors.New("cannot navigate path: parent is not a map")
				break
			}
		}
		if err != nil {
			break
		}
	}
	return err
}
