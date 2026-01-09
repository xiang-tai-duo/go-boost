// Package boost
// File:        xml.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/xml.go
// Author:      Vibe Coding
// Created:     12/30/2025 11:03:46
// Description: XML provides utility methods for XML operations, including marshal, unmarshal, validation, formatting, minification, and file operations.
// --------------------------------------------------------------------------------
package boost

import (
	"encoding/xml"
	"errors"
	"os"
)

type (
	XML struct {
		value interface{}
		text  string
	}
)

func NewXML(initialValue ...interface{}) *XML {
	j := &XML{}
	if len(initialValue) > 0 {
		j.value = initialValue[0]
		if s, ok := initialValue[0].(string); ok {
			var parsedValue interface{}
			if err := xml.Unmarshal([]byte(s), &parsedValue); err == nil {
				j.text = s
			} else if xmlBytes, err := xml.Marshal(initialValue[0]); err == nil {
				j.text = string(xmlBytes)
			}
		} else if xmlBytes, err := xml.Marshal(initialValue[0]); err == nil {
			j.text = string(xmlBytes)
		}
	}
	return j
}

func (j *XML) Format(indent string) (string, error) {
	var err error
	var formatted []byte
	var parsedValue interface{}
	var result string
	if err = xml.Unmarshal([]byte(j.text), &parsedValue); err == nil {
		if formatted, err = xml.MarshalIndent(parsedValue, "", indent); err == nil {
			result = string(formatted)
		}
	}
	return result, err
}

func (j *XML) GetValue(path string) (interface{}, error) {
	var err error
	var parsedValue interface{}
	var result interface{}
	if err = xml.Unmarshal([]byte(j.text), &parsedValue); err == nil {
		result, err = getNestedValue(parsedValue, path)
	}
	return result, err
}

func (j *XML) GetInteger(path string) (int, error) {
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

func (j *XML) GetString(path string) (string, error) {
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

func (j *XML) GetMap(path string) (map[string]interface{}, error) {
	var parsedValue interface{}
	var err error
	var result map[string]interface{}
	if err = xml.Unmarshal([]byte(j.text), &parsedValue); err == nil {
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

func (j *XML) Marshal() (string, error) {
	return j.text, nil
}

func (j *XML) Minify() (string, error) {
	var err error
	var minified []byte
	var parsedValue interface{}
	var result string
	if err = xml.Unmarshal([]byte(j.text), &parsedValue); err == nil {
		if minified, err = xml.Marshal(parsedValue); err == nil {
			result = string(minified)
		}
	}
	return result, err
}

func (j *XML) SetValue(path string, value interface{}) (*XML, error) {
	var err error
	var parsedValue interface{}
	if err = xml.Unmarshal([]byte(j.text), &parsedValue); err == nil {
		if err = setNestedValue(parsedValue, path, value); err == nil {
			var xmlBytes []byte
			if xmlBytes, err = xml.Marshal(parsedValue); err == nil {
				j.text = string(xmlBytes)
				j.value = parsedValue
			}
		}
	}
	return j, err
}

func (j *XML) SetInteger(path string, value int) (*XML, error) {
	return j.SetValue(path, value)
}

func (j *XML) SetString(path string, value string) (*XML, error) {
	return j.SetValue(path, value)
}

func (j *XML) Unmarshal(target interface{}) error {
	return xml.Unmarshal([]byte(j.text), target)
}

func (j *XML) Validate() (bool, error) {
	var parsedValue interface{}
	var result bool
	var err error
	if err = xml.Unmarshal([]byte(j.text), &parsedValue); err == nil {
		result = true
	}
	return result, err
}

func (j *XML) WriteFile(filePath string, indent ...string) error {
	var err error
	var file *os.File
	var xmlBytes []byte
	var indentValue string
	if len(indent) > 0 {
		indentValue = indent[0]
	}
	if indentValue != "" {
		if xmlBytes, err = xml.MarshalIndent(j.value, "", indentValue); err != nil {
			return err
		}
	} else {
		if xmlBytes, err = xml.Marshal(j.value); err != nil {
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
