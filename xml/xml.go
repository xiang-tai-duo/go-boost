// Package xml
// File:        xml.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/xml/xml.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: XML provides utility methods for XML operations, including marshal, unmarshal, validation, formatting, minification, and file operations.
// --------------------------------------------------------------------------------
package xml

import (
	__xml "encoding/xml"
	"errors"
	"os"
	"strings"

	"github.com/beevik/etree"
	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection GoUnusedConst,GoSnakeCaseUsage
const (
	DEFAULT_INDENT  = "    "
	EMPTY_PREFIX    = ""
	MODULE_NAME_XML = "xml"
	PATH_SEPARATOR  = "."
	TAB_CHARACTER   = "\t"
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_XML, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_XML, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_XML, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedExportedFunction
func Format(xml string) (string, error) {
	result := ""
	err := error(nil)
	result, err = FormatIndent(xml, DEFAULT_INDENT)
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func FormatIndent(xml string, indent string) (string, error) {
	result := xml
	err := error(nil)
	doc := etree.NewDocument()
	doc.ReadSettings.Permissive = true
	doc.ReadSettings.PreserveCData = true
	doc.WriteSettings.CanonicalEndTags = false
	doc.WriteSettings.CanonicalText = false
	doc.WriteSettings.CanonicalAttrVal = false
	if err = doc.ReadFromString(xml); err == nil {
		if indent == "" {
			doc.Indent(etree.NoIndent)
		} else if strings.Contains(indent, TAB_CHARACTER) {
			doc.IndentTabs()
		} else {
			doc.Indent(len(indent))
		}
		var formattedXml string
		if formattedXml, err = doc.WriteToString(); err == nil {
			result = formattedXml
		}
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func GetInteger(xml string, path string) (int, error) {
	result := 0
	err := error(nil)
	var value interface{}
	if value, err = GetValue(xml, path); err == nil {
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

//goland:noinspection DuplicatedCode,GoUnusedExportedFunction
func GetMap(xml string, path string) (map[string]interface{}, error) {
	result := map[string]interface{}(nil)
	err := error(nil)
	var parsedValue interface{}
	if err = __xml.Unmarshal([]byte(xml), &parsedValue); err == nil {
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

//goland:noinspection GoUnusedExportedFunction
func GetString(xml string, path string) (string, error) {
	result := ""
	err := error(nil)
	var value interface{}
	if value, err = GetValue(xml, path); err == nil {
		if s, ok := value.(string); ok {
			result = s
		} else {
			err = errors.New("value is not a string")
		}
	}
	return result, err
}

func GetValue(xml string, path string) (interface{}, error) {
	result := interface{}(nil)
	err := error(nil)
	var parsedValue interface{}
	if err = __xml.Unmarshal([]byte(xml), &parsedValue); err == nil {
		result, err = getNestedValue(parsedValue, path)
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func Minify(xml string) (string, error) {
	result := ""
	err := error(nil)
	var parsedValue interface{}
	if err = __xml.Unmarshal([]byte(xml), &parsedValue); err == nil {
		var minified []byte
		if minified, err = __xml.Marshal(parsedValue); err == nil {
			result = string(minified)
		}
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func SetInteger(xml string, path string, value int) (string, error) {
	result := ""
	err := error(nil)
	result, err = SetValue(xml, path, value)
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func SetString(xml string, path string, value string) (string, error) {
	result := ""
	err := error(nil)
	result, err = SetValue(xml, path, value)
	return result, err
}

//goland:noinspection DuplicatedCode
func SetValue(xml string, path string, value interface{}) (string, error) {
	result := xml
	err := error(nil)
	var parsedValue interface{}
	if err = __xml.Unmarshal([]byte(xml), &parsedValue); err == nil {
		if err = setNestedValue(parsedValue, path, value); err == nil {
			var xmlBytes []byte
			if xmlBytes, err = __xml.Marshal(parsedValue); err == nil {
				result = string(xmlBytes)
			}
		}
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func Unmarshal(xml string, target interface{}) error {
	err := error(nil)
	err = __xml.Unmarshal([]byte(xml), target)
	return err
}

//goland:noinspection GoUnusedExportedFunction
func Validate(xml string) (bool, error) {
	result := false
	err := error(nil)
	var parsedValue interface{}
	if err = __xml.Unmarshal([]byte(xml), &parsedValue); err == nil {
		result = true
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func WriteFile(xml string, filePath string, indent ...string) error {
	err := error(nil)
	var xmlBytes []byte
	indentValue := ""
	if len(indent) > 0 {
		indentValue = indent[0]
	}
	if indentValue == "" {
		xmlBytes = []byte(xml)
	} else {
		var parsedValue interface{}
		if err = __xml.Unmarshal([]byte(xml), &parsedValue); err == nil {
			xmlBytes, err = __xml.MarshalIndent(parsedValue, EMPTY_PREFIX, indentValue)
		}
	}
	if err == nil {
		var file *os.File
		if file, err = os.Create(filePath); err == nil {
			defer func(file *os.File) {
				_ = file.Close()
			}(file)
			_, err = file.Write(xmlBytes)
		}
	}
	return err
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_XML, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection DuplicatedCode
func getNestedValue(data interface{}, path string) (interface{}, error) {
	result := interface{}(nil)
	err := error(nil)
	parts := strings.Split(path, PATH_SEPARATOR)
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
	parts := strings.Split(path, PATH_SEPARATOR)
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
