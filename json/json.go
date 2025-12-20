// Package json
// File:        json.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/json/json.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: JSON provides utility methods for JSON operations, including marshal, unmarshal, validation, formatting, minification, and file operations.
// --------------------------------------------------------------------------------
package json

import (
	jsonlib "encoding/json"
	"errors"
	"github.com/xiang-tai-duo/go-boost/logger"
	"os"
	"strings"
)

type (
	JSON struct {
		text  string
		value interface{}
	}
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	EMPTY_INDENT_PREFIX = ""
	MODULE_NAME_JSON    = "json"
	PATH_SEPARATOR      = "."
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_JSON, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_JSON, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_JSON, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection DuplicatedCode
func New(json ...interface{}) *JSON {
	j := &JSON{}
	if len(json) > 0 {
		j.value = json[0]
		err := error(nil)
		if s, ok := json[0].(string); ok {
			var parsedValue interface{}
			if err = jsonlib.Unmarshal([]byte(s), &parsedValue); err == nil {
				j.text = s
			} else {
				var jsonBytes []byte
				if jsonBytes, err = jsonlib.Marshal(json[0]); err == nil {
					j.text = string(jsonBytes)
				}
			}
		} else {
			var jsonBytes []byte
			if jsonBytes, err = jsonlib.Marshal(json[0]); err == nil {
				j.text = string(jsonBytes)
			}
		}
	}
	return j
}

func (j *JSON) Format(indent string) (string, error) {
	result := ""
	err := error(nil)
	var parsedValue interface{}
	if err = jsonlib.Unmarshal([]byte(j.text), &parsedValue); err == nil {
		var formatted []byte
		if formatted, err = jsonlib.MarshalIndent(parsedValue, EMPTY_INDENT_PREFIX, indent); err == nil {
			result = string(formatted)
		}
	}
	return result, err
}

func (j *JSON) GetInteger(path string) (int, error) {
	result := 0
	err := error(nil)
	var value interface{}
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
	result := make(map[string]interface{})
	err := error(nil)
	var parsedValue interface{}
	if err = jsonlib.Unmarshal([]byte(j.text), &parsedValue); err == nil {
		var nestedValue interface{}
		if nestedValue, err = getNestedValue(parsedValue, path); err == nil {
			if m, ok := nestedValue.(map[string]interface{}); ok {
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
	result := ""
	err := error(nil)
	var value interface{}
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
	err := error(nil)
	var result interface{}
	var parsedValue interface{}
	if err = jsonlib.Unmarshal([]byte(j.text), &parsedValue); err == nil {
		result, err = getNestedValue(parsedValue, path)
	}
	return result, err
}

func (j *JSON) Marshal() (string, error) {
	result := j.text
	err := error(nil)
	return result, err
}

func (j *JSON) Minify() (string, error) {
	result := ""
	err := error(nil)
	var parsedValue interface{}
	if err = jsonlib.Unmarshal([]byte(j.text), &parsedValue); err == nil {
		var minified []byte
		if minified, err = jsonlib.Marshal(parsedValue); err == nil {
			result = string(minified)
		}
	}
	return result, err
}

func (j *JSON) SetInteger(path string, value int) (*JSON, error) {
	result, err := j.SetValue(path, value)
	return result, err
}

func (j *JSON) SetString(path string, value string) (*JSON, error) {
	result, err := j.SetValue(path, value)
	return result, err
}

//goland:noinspection DuplicatedCode
func (j *JSON) SetValue(path string, value interface{}) (*JSON, error) {
	result := j
	err := error(nil)
	var parsedValue interface{}
	if err = jsonlib.Unmarshal([]byte(j.text), &parsedValue); err == nil {
		if err = setNestedValue(parsedValue, path, value); err == nil {
			var jsonBytes []byte
			if jsonBytes, err = jsonlib.Marshal(parsedValue); err == nil {
				j.text = string(jsonBytes)
				j.value = parsedValue
			}
		}
	}
	return result, err
}

func (j *JSON) Unmarshal(target interface{}) error {
	err := jsonlib.Unmarshal([]byte(j.text), target)
	return err
}

func (j *JSON) Validate() (bool, error) {
	result := false
	err := error(nil)
	var parsedValue interface{}
	if err = jsonlib.Unmarshal([]byte(j.text), &parsedValue); err == nil {
		result = true
	}
	return result, err
}

func (j *JSON) WriteFile(filePath string, indent ...string) error {
	err := error(nil)
	indentValue := EMPTY_INDENT_PREFIX
	if len(indent) > 0 {
		indentValue = indent[0]
	}
	var jsonBytes []byte
	if indentValue != EMPTY_INDENT_PREFIX {
		jsonBytes, err = jsonlib.MarshalIndent(j.value, EMPTY_INDENT_PREFIX, indentValue)
	} else {
		jsonBytes, err = jsonlib.Marshal(j.value)
	}
	if err == nil {
		var file *os.File
		if file, err = os.Create(filePath); err == nil {
			defer func(file *os.File) {
				_ = file.Close()
			}(file)
			_, err = file.Write(jsonBytes)
		}
	}
	return err
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_JSON, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection DuplicatedCode
func getNestedValue(data interface{}, path string) (interface{}, error) {
	err := error(nil)
	var result interface{}
	parts := strings.Split(path, PATH_SEPARATOR)
	current := data
	for _, part := range parts {
		switch m := current.(type) {
		case map[string]interface{}:
			if value, ok := m[part]; ok {
				current = value
			} else {
				err = errors.New("path not found")
			}
		default:
			err = errors.New("cannot navigate path: parent is not a map")
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
	err := error(nil)
	parts := strings.Split(path, PATH_SEPARATOR)
	current := data
	for i, part := range parts {
		if i == len(parts)-1 {
			switch m := current.(type) {
			case map[string]interface{}:
				m[part] = value
			default:
				err = errors.New("cannot set value at path: parent is not a map")
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
			}
		}
		if err != nil {
			break
		}
	}
	return err
}
