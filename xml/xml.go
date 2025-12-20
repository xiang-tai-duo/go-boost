// Package xml
// File:        xml.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/xml/xml.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: XML provides utility methods for XML operations, including marshal, unmarshal, validation, formatting, minification, and file operations.
// --------------------------------------------------------------------------------
package xml

import (
	__xml "encoding/xml"
	"errors"
	"os"
	"strings"
)

type (
	XML struct {
		value interface{}
		text  string
	}
)

//goland:noinspection GoUnusedExportedFunction,DuplicatedCode
func New(xml ...interface{}) *XML {
	x := &XML{}
	if len(xml) > 0 {
		x.value = xml[0]
		if s, ok := xml[0].(string); ok {
			var parsedValue interface{}
			if err := __xml.Unmarshal([]byte(s), &parsedValue); err == nil {
				x.text = s
			} else if xmlBytes, err := __xml.Marshal(xml[0]); err == nil {
				x.text = string(xmlBytes)
			}
		} else if xmlBytes, err := __xml.Marshal(xml[0]); err == nil {
			x.text = string(xmlBytes)
		}
	}
	return x
}

func (x *XML) Format(indent string) (string, error) {
	result := ""
	err := error(nil)
	var formatted []byte
	var parsedValue interface{}
	if err = __xml.Unmarshal([]byte(x.text), &parsedValue); err == nil {
		if formatted, err = __xml.MarshalIndent(parsedValue, "", indent); err == nil {
			result = string(formatted)
		}
	}
	return result, err
}

func (x *XML) GetValue(path string) (interface{}, error) {
	result := interface{}(nil)
	err := error(nil)
	var parsedValue interface{}
	if err = __xml.Unmarshal([]byte(x.text), &parsedValue); err == nil {
		result, err = getNestedValue(parsedValue, path)
	}
	return result, err
}

func (x *XML) GetInteger(path string) (int, error) {
	result := 0
	err := error(nil)
	var value interface{}
	if value, err = x.GetValue(path); err == nil {
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

func (x *XML) GetString(path string) (string, error) {
	result := ""
	err := error(nil)
	var value interface{}
	if value, err = x.GetValue(path); err == nil {
		if s, ok := value.(string); ok {
			result = s
		} else {
			err = errors.New("value is not a string")
		}
	}

	return result, err
}

//goland:noinspection DuplicatedCode
func (x *XML) GetMap(path string) (map[string]interface{}, error) {
	result := map[string]interface{}(nil)
	err := error(nil)
	var parsedValue interface{}
	if err = __xml.Unmarshal([]byte(x.text), &parsedValue); err == nil {
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

func (x *XML) Marshal() (string, error) {
	return x.text, nil
}

func (x *XML) Minify() (string, error) {
	result := ""
	err := error(nil)
	var minified []byte
	var parsedValue interface{}
	if err = __xml.Unmarshal([]byte(x.text), &parsedValue); err == nil {
		if minified, err = __xml.Marshal(parsedValue); err == nil {
			result = string(minified)
		}
	}
	return result, err
}

//goland:noinspection DuplicatedCode
func (x *XML) SetValue(path string, value interface{}) (*XML, error) {
	err := error(nil)
	var parsedValue interface{}
	if err = __xml.Unmarshal([]byte(x.text), &parsedValue); err == nil {
		if err = setNestedValue(parsedValue, path, value); err == nil {
			var xmlBytes []byte
			if xmlBytes, err = __xml.Marshal(parsedValue); err == nil {
				x.text = string(xmlBytes)
				x.value = parsedValue
			}
		}
	}
	return x, err
}

func (x *XML) SetInteger(path string, value int) (*XML, error) {
	return x.SetValue(path, value)
}

func (x *XML) SetString(path string, value string) (*XML, error) {
	return x.SetValue(path, value)
}

func (x *XML) Unmarshal(target interface{}) error {
	return __xml.Unmarshal([]byte(x.text), target)
}

func (x *XML) Validate() (bool, error) {
	result := false
	err := error(nil)
	var parsedValue interface{}
	if err = __xml.Unmarshal([]byte(x.text), &parsedValue); err == nil {
		result = true
	}
	return result, err
}

func (x *XML) WriteFile(filePath string, indent ...string) error {
	err := error(nil)
	var file *os.File
	var xmlBytes []byte
	var indentValue string
	if len(indent) > 0 {
		indentValue = indent[0]
	}
	if indentValue != "" {
		if xmlBytes, err = __xml.MarshalIndent(x.value, "", indentValue); err != nil {
			return err
		}
	} else {
		if xmlBytes, err = __xml.Marshal(x.value); err != nil {
			return err
		}
	}
	if file, err = os.Create(filePath); err == nil {
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
		_, err = file.Write(xmlBytes)
	}
	return err
}

//goland:noinspection DuplicatedCode
func getNestedValue(data interface{}, path string) (interface{}, error) {
	result := interface{}(nil)
	err := error(nil)
	parts := strings.Split(path, ".")
	current := data
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
	err := error(nil)
	parts := strings.Split(path, ".")
	current := data
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
