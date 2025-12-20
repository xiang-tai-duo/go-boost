// Package boost
// File:        json.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: JSON provides utility methods for JSON operations, including
//
//	marshal, unmarshal, validation, formatting, minification, and file operations.
package boost

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

type (
	JSON struct {
		_expr interface{}
		_json string
	}
)

func NewJSON(initialValue ...interface{}) *JSON {
	jsonInstance := &JSON{}
	if len(initialValue) > 0 {
		jsonInstance._expr = initialValue[0]
		if str, ok := initialValue[0].(string); ok {
			var parsedValue interface{}
			if err := json.Unmarshal([]byte(str), &parsedValue); err == nil {
				jsonInstance._json = str
			} else if jsonBytes, err := json.Marshal(initialValue[0]); err == nil {
				jsonInstance._json = string(jsonBytes)
			}
		} else if jsonBytes, err := json.Marshal(initialValue[0]); err == nil {
			jsonInstance._json = string(jsonBytes)
		}
	}
	return jsonInstance
}

func (j *JSON) Format(indent string) (string, error) {
	var err error
	var formatted []byte
	var parsedValue interface{}
	if err = json.Unmarshal([]byte(j._json), &parsedValue); err == nil {
		if formatted, err = json.MarshalIndent(parsedValue, "", indent); err != nil {
			return "", err
		}
	}
	return string(formatted), err
}

func (j *JSON) GetValue(path string) (interface{}, error) {
	var err error
	var parsedValue interface{}
	if err = json.Unmarshal([]byte(j._json), &parsedValue); err == nil {
		return getNestedValue(parsedValue, path)
	}
	return nil, err
}

func (j *JSON) GetValueInt(path string) (int, error) {
	var err error
	var value interface{}
	if value, err = j.GetValue(path); err != nil {
		return 0, err
	}
	if intValue, ok := value.(float64); ok {
		return int(intValue), nil
	} else if intValue, ok := value.(int); ok {
		return intValue, nil
	}
	return 0, errors.New("value is not an integer")
}

func (j *JSON) GetValueString(path string) (string, error) {
	var err error
	var value interface{}
	if value, err = j.GetValue(path); err != nil {
		return "", err
	}
	if strValue, ok := value.(string); ok {
		return strValue, nil
	}
	return "", errors.New("value is not a string")
}

func (j *JSON) GetValues(path string) (map[string]interface{}, error) {
	var parsedValue interface{}
	var err error
	if err = json.Unmarshal([]byte(j._json), &parsedValue); err != nil {
		return nil, err
	}
	parentValue, getErr := getNestedValue(parsedValue, path)
	if getErr != nil {
		return nil, getErr
	}
	if parentMap, ok := parentValue.(map[string]interface{}); ok {
		result := make(map[string]interface{})
		for key, value := range parentMap {
			result[key] = value
		}
		return result, nil
	}
	return nil, errors.New("specified path does not point to a map")
}

func (j *JSON) Json(initialValue ...interface{}) (string, error) {
	var err error
	if len(initialValue) > 0 {
		j._expr = initialValue[0]
		if str, ok := initialValue[0].(string); ok {
			var parsedValue interface{}
			if err = json.Unmarshal([]byte(str), &parsedValue); err == nil {
				j._json = str
				err = nil
			} else if jsonBytes, marshalErr := json.Marshal(initialValue[0]); marshalErr == nil {
				j._json = string(jsonBytes)
				err = nil
			} else {
				err = marshalErr
			}
		} else if jsonBytes, marshalErr := json.Marshal(initialValue[0]); marshalErr == nil {
			j._json = string(jsonBytes)
			err = nil
		} else {
			err = marshalErr
		}
	}
	return j._json, err
}

func (j *JSON) Marshal() (string, error) {
	return j._json, nil
}

func (j *JSON) Minify() (string, error) {
	var err error
	var minified []byte
	var parsedValue interface{}
	if err = json.Unmarshal([]byte(j._json), &parsedValue); err == nil {
		if minified, err = json.Marshal(parsedValue); err != nil {
			return "", err
		}
	}
	return string(minified), err
}

func (j *JSON) SetValue(path string, value interface{}) (*JSON, error) {
	var err error
	var parsedValue interface{}
	if err = json.Unmarshal([]byte(j._json), &parsedValue); err != nil {
		return j, err
	}
	if err = setNestedValue(parsedValue, path, value); err != nil {
		return j, err
	}
	if jsonBytes, marshalErr := json.Marshal(parsedValue); marshalErr != nil {
		return j, marshalErr
	} else {
		j._json = string(jsonBytes)
		j._expr = parsedValue
		return j, nil
	}
}

func (j *JSON) SetValueInt(path string, value int) (*JSON, error) {
	return j.SetValue(path, value)
}

func (j *JSON) SetValueString(path string, value string) (*JSON, error) {
	return j.SetValue(path, value)
}

func (j *JSON) Unmarshal(target interface{}) error {
	var err error
	if err = json.Unmarshal([]byte(j._json), target); err == nil {
		return nil
	}
	return err
}

func (j *JSON) Validate() (bool, error) {
	var parsedValue interface{}
	if err := json.Unmarshal([]byte(j._json), &parsedValue); err == nil {
		return true, nil
	}
	return false, nil
}

func (j *JSON) WriteFile(filePath string, indent string) error {
	var err error
	var file *os.File
	var jsonBytes []byte
	if indent != "" {
		if jsonBytes, err = json.MarshalIndent(j._expr, "", indent); err != nil {
			return err
		}
	} else if jsonBytes, err = json.Marshal(j._expr); err != nil {
		return err
	}
	if file, err = os.Create(filePath); err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	if _, err = file.Write(jsonBytes); err != nil {
		return err
	}
	return nil
}

func getNestedValue(data interface{}, path string) (interface{}, error) {
	parts := strings.Split(path, ".")
	current := data

	for _, part := range parts {
		switch m := current.(type) {
		case map[string]interface{}:
			if val, ok := m[part]; ok {
				current = val
			} else {
				return nil, errors.New("path not found")
			}
		default:
			return nil, errors.New("cannot navigate path: parent is not a map")
		}
	}
	return current, nil
}

func setNestedValue(data interface{}, path string, value interface{}) error {
	parts := strings.Split(path, ".")
	current := data

	for i, part := range parts {
		if i == len(parts)-1 {
			switch m := current.(type) {
			case map[string]interface{}:
				m[part] = value
			default:
				return errors.New("cannot set value at path: parent is not a map")
			}
		} else {
			switch m := current.(type) {
			case map[string]interface{}:
				if _, ok := m[part]; !ok {
					m[part] = make(map[string]interface{})
				}
				current = m[part]
			default:
				return errors.New("cannot navigate path: parent is not a map")
			}
		}
	}
	return nil
}
