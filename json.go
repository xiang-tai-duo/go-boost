// --------------------------------------------------------------------------------
// File:        json.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: JSON provides utility methods for JSON operations, including
//              marshal, unmarshal, validation, formatting, minification, and file operations.
// --------------------------------------------------------------------------------

package boost

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

// JSON provides utility methods for JSON operations.
type JSON struct {
	_expr interface{}
	_json string
}

// NewJSON creates a new JSON instance with optional initial value.
// initialValue: Optional initial value to set
// Returns: JSON instance with the initial value
// Usage:
// jsonInstance := NewJSON()
// returns new JSON instance with empty values
// jsonInstance := NewJSON(`{"key":"value"}`)
// returns new JSON instance with parsed JSON value
func NewJSON(initialValue ...interface{}) *JSON {
	jsonInstance := &JSON{}
	if len(initialValue) > 0 {
		jsonInstance._expr = initialValue[0]
		// Initialize _json string
		if str, ok := initialValue[0].(string); ok {
			// Check if string is valid JSON
			var parsedValue interface{}
			if err := json.Unmarshal([]byte(str), &parsedValue); err == nil {
				// It's a valid JSON string, use it directly
				jsonInstance._json = str
			} else {
				// It's not a valid JSON string, marshal the value
				if jsonBytes, err := json.Marshal(initialValue[0]); err == nil {
					jsonInstance._json = string(jsonBytes)
				}
			}
		} else {
			// It's not a string, marshal the value
			if jsonBytes, err := json.Marshal(initialValue[0]); err == nil {
				jsonInstance._json = string(jsonBytes)
			}
		}
	}
	return jsonInstance
}

// Format formats the JSON value with indentation for better readability.
// indent: Indentation string to use (e.g., "  " for 2 spaces)
// Returns: formatted JSON string, error if any occurred
// Usage:
// jsonInstance := NewJSON(`{"key":"value"}`)
// formatted, err := jsonInstance.Format("  ")
// returns "{\n  \"key\": \"value\"\n}", nil on success
func (j *JSON) Format(indent string) (string, error) {
	var err error
	var formatted []byte
	var parsedValue interface{}
	if err = json.Unmarshal([]byte(j._json), &parsedValue); err == nil {
		if formatted, err = json.MarshalIndent(parsedValue, "", indent); err == nil {
			// Formatting successful
		}
	}
	return string(formatted), err
}

// Json returns the JSON string or reinitializes with new value.
// initialValue: Optional initial value to reinitialize
// Returns: JSON string, error if any occurred
// Usage:
// jsonInstance := NewJSON(`{"key":"value"}`)
// jsonStr, err := jsonInstance.Json()
// returns "{\"key\":\"value\"}", nil
// jsonStr, err := jsonInstance.Json(map[string]int{"num": 42})
// returns "{\"num\":42}", nil
func (j *JSON) Json(initialValue ...interface{}) (string, error) {
	var err error
	if len(initialValue) > 0 {
		// Reinitialize with new value
		j._expr = initialValue[0]
		if str, ok := initialValue[0].(string); ok {
			// Check if string is valid JSON
			var parsedValue interface{}
			if err = json.Unmarshal([]byte(str), &parsedValue); err == nil {
				// It's a valid JSON string, use it directly
				j._json = str
				err = nil
			} else {
				// It's not a valid JSON string, marshal the value
				if jsonBytes, marshalErr := json.Marshal(initialValue[0]); marshalErr == nil {
					j._json = string(jsonBytes)
					err = nil
				} else {
					err = marshalErr
				}
			}
		} else {
			// It's not a string, marshal the value
			if jsonBytes, marshalErr := json.Marshal(initialValue[0]); marshalErr == nil {
				j._json = string(jsonBytes)
				err = nil
			} else {
				err = marshalErr
			}
		}
	}
	return j._json, err
}

// Marshal converts the JSON value to a JSON string.
// Returns: JSON string, error if any occurred
// Usage:
// jsonInstance := NewJSON(map[string]string{"key": "value"})
// jsonString, err := jsonInstance.Marshal()
// returns "{\"key\":\"value\"}", nil on success
func (j *JSON) Marshal() (string, error) {
	return j._json, nil
}

// Minify removes all whitespace from the JSON value to reduce its size.
// Returns: minified JSON string, error if any occurred
// Usage:
// jsonInstance := NewJSON(`{\n  \"key\": \"value\"\n}`)
// minified, err := jsonInstance.Minify()
// returns "{\"key\":\"value\"}", nil on success
func (j *JSON) Minify() (string, error) {
	var err error
	var minified []byte
	var parsedValue interface{}
	if err = json.Unmarshal([]byte(j._json), &parsedValue); err == nil {
		if minified, err = json.Marshal(parsedValue); err == nil {
			// Minification successful
		}
	}
	return string(minified), err
}

// Validate checks if the JSON value is a valid JSON format.
// Returns: boolean indicating if the value is valid JSON, error if any occurred
// Usage:
// jsonInstance := NewJSON(`{"key":"value"}`)
// isValid, err := jsonInstance.Validate()
// returns true, nil on success
func (j *JSON) Validate() (bool, error) {
	var err error
	var isValid bool
	var parsedValue interface{}
	// Try to unmarshal the JSON string
	if err = json.Unmarshal([]byte(j._json), &parsedValue); err == nil {
		isValid = true
		err = nil
	} else {
		isValid = false
		err = nil
	}
	return isValid, err
}

// Unmarshal converts the JSON value to a Go value.
// target: Pointer to the Go value to store the result
// Returns: error if any occurred
// Usage:
// jsonInstance := NewJSON(`{"key":"value"}`)
// var result map[string]string
// err := jsonInstance.Unmarshal(&result)
// returns nil on success, result will be map[string]string{"key": "value"}
func (j *JSON) Unmarshal(target interface{}) error {
	var err error
	if err = json.Unmarshal([]byte(j._json), target); err == nil {
		// Unmarshaling successful
	}
	return err
}

// SetValue sets a value in the JSON using dot notation.
// path: Path to the value using dot notation (e.g., "Root.Id")
// value: Value to set
// Returns: JSON instance for method chaining, error if any occurred
// Usage:
// jsonInstance := NewJSON(`{"Root": {"Id": 1}}`)
// jsonInstance.SetValue("Root.Id", 3).SetValue("Root.Name", "test")
// jsonStr, _ := jsonInstance.Json()
// returns "{\"Root\":{\"Id\":3,\"Name\":\"test\"}}"
func (j *JSON) SetValue(path string, value interface{}) (*JSON, error) {
	var err error
	var parsedValue interface{}
	if err = json.Unmarshal([]byte(j._json), &parsedValue); err == nil {
		if err = setNestedValue(parsedValue, path, value); err == nil {
			if jsonBytes, marshalErr := json.Marshal(parsedValue); marshalErr == nil {
				j._json = string(jsonBytes)
				j._expr = parsedValue
				err = nil
			} else {
				err = marshalErr
			}
		}
	}
	return j, err
}

// GetValue gets a value from the JSON using dot notation.
// path: Path to the value using dot notation (e.g., "Root.Id")
// Returns: Value as interface{}, error if any occurred
// Usage:
// jsonInstance := NewJSON(`{"Root": {"Id": 3}}`)
// value, err := jsonInstance.GetValue("Root.Id")
// returns 3, nil
func (j *JSON) GetValue(path string) (interface{}, error) {
	var err error
	var parsedValue interface{}
	if err = json.Unmarshal([]byte(j._json), &parsedValue); err == nil {
		return getNestedValue(parsedValue, path)
	}
	return nil, err
}

// SetValueInt sets an integer value in the JSON using dot notation.
// path: Path to the value using dot notation (e.g., "Root.Id")
// value: Integer value to set
// Returns: JSON instance for method chaining, error if any occurred
// Usage:
// jsonInstance := NewJSON(`{"Root": {"Id": 1}}`)
// jsonInstance.SetValueInt("Root.Id", 5)
func (j *JSON) SetValueInt(path string, value int) (*JSON, error) {
	return j.SetValue(path, value)
}

// SetValueString sets a string value in the JSON using dot notation.
// path: Path to the value using dot notation (e.g., "Root.Name")
// value: String value to set
// Returns: JSON instance for method chaining, error if any occurred
// Usage:
// jsonInstance := NewJSON(`{"Root": {"Name": "old"}}`)
// jsonInstance.SetValueString("Root.Name", "new")
func (j *JSON) SetValueString(path string, value string) (*JSON, error) {
	return j.SetValue(path, value)
}

// GetValueInt gets an integer value from the JSON using dot notation.
// path: Path to the value using dot notation (e.g., "Root.Id")
// Returns: Integer value, error if any occurred
// Usage:
// jsonInstance := NewJSON(`{"Root": {"Id": 5}}`)
// value, err := jsonInstance.GetValueInt("Root.Id")
// returns 5, nil
func (j *JSON) GetValueInt(path string) (int, error) {
	var err error
	var value interface{}
	if value, err = j.GetValue(path); err == nil {
		if intValue, ok := value.(float64); ok {
			return int(intValue), nil
		} else if intValue, ok := value.(int); ok {
			return intValue, nil
		}
		err = errors.New("value is not an integer")
	}
	return 0, err
}

// GetValueString gets a string value from the JSON using dot notation.
// path: Path to the value using dot notation (e.g., "Root.Name")
// Returns: String value, error if any occurred
// Usage:
// jsonInstance := NewJSON(`{"Root": {"Name": "test"}}`)
// value, err := jsonInstance.GetValueString("Root.Name")
// returns "test", nil
func (j *JSON) GetValueString(path string) (string, error) {
	var err error
	var value interface{}
	if value, err = j.GetValue(path); err == nil {
		if strValue, ok := value.(string); ok {
			return strValue, nil
		}
		err = errors.New("value is not a string")
	}
	return "", err
}

// GetValues gets all direct child nodes and their values from the JSON using dot notation.
// path: Path to the parent node using dot notation (e.g., "Root")
// Returns: Map of child node names to their values, error if any occurred
// Usage:
// jsonInstance := NewJSON(`{"Root": {"Id": 3, "Name": "test", "Nested": {"Value": 42}}}`)
// values, err := jsonInstance.GetValues("Root")
// returns map[string]interface{}{"Id": 3, "Name": "test", "Nested": map[string]interface{}{"Value": 42}}, nil
// Example of iterating the result:
//
//	for key, value := range values {
//	    fmt.Printf("Key: %s, Value: %v\n", key, value)
//	}
//
// Outputs:
// Key: Id, Value: 3
// Key: Name, Value: test
// Key: Nested, Value: map[Value:42]
//
// Example of checking value types:
//
//	for key, value := range values {
//	    switch v := value.(type) {
//	    case float64: // JSON numbers are parsed as float64 by default
//	        fmt.Printf("Key: %s, Type: number, Value: %v\n", key, v)
//	    case string:
//	        fmt.Printf("Key: %s, Type: string, Value: %v\n", key, v)
//	    case bool:
//	        fmt.Printf("Key: %s, Type: bool, Value: %v\n", key, v)
//	    case map[string]interface{}:
//	        fmt.Printf("Key: %s, Type: map, Value: %v\n", key, v)
//	    case []interface{}:
//	        fmt.Printf("Key: %s, Type: array, Value: %v\n", key, v)
//	    case nil:
//	        fmt.Printf("Key: %s, Type: null, Value: nil\n", key)
//	    default:
//	        fmt.Printf("Key: %s, Type: unknown, Value: %v\n", key, v)
//	    }
//	}
//
// Outputs:
// Key: Id, Type: number, Value: 3
// Key: Name, Type: string, Value: test
// Key: Nested, Type: map, Value: map[Value:42]
func (j *JSON) GetValues(path string) (map[string]interface{}, error) {
	var err error
	var parsedValue interface{}
	if err = json.Unmarshal([]byte(j._json), &parsedValue); err == nil {
		parentValue, getErr := getNestedValue(parsedValue, path)
		if getErr == nil {
			if parentMap, ok := parentValue.(map[string]interface{}); ok {
				// Return all direct child nodes, do not recursively process
				result := make(map[string]interface{})
				for key, value := range parentMap {
					result[key] = value
				}
				return result, nil
			}
			err = errors.New("specified path does not point to a map")
		} else {
			err = getErr
		}
	}
	return nil, err
}

// Helper function to set nested value using dot notation
func setNestedValue(data interface{}, path string, value interface{}) error {
	parts := strings.Split(path, ".")
	current := data

	for i, part := range parts {
		if i == len(parts)-1 {
			// Last part, set the value
			switch m := current.(type) {
			case map[string]interface{}:
				m[part] = value
			default:
				return errors.New("cannot set value at path: parent is not a map")
			}
		} else {
			// Navigate to next level
			switch m := current.(type) {
			case map[string]interface{}:
				if _, ok := m[part]; !ok {
					// Create missing map if not exists
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

// Helper function to get nested value using dot notation
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

// WriteFile writes the current JSON instance value to a file.
// filePath: Path to the file to write
// indent: Indentation string to use (e.g., "  " for 2 spaces), leave empty for minified JSON
// Returns: error if any occurred
// Usage:
// jsonInstance := NewJSON(map[string]string{"key": "value"})
// err := jsonInstance.WriteFile("output.json", "  ")
// returns nil on success, file will contain formatted JSON
func (j *JSON) WriteFile(filePath string, indent string) error {
	var err error
	var file *os.File
	var jsonBytes []byte
	if indent != "" {
		if jsonBytes, err = json.MarshalIndent(j._expr, "", indent); err == nil {
			// Marshaling with indentation successful
		}
	} else {
		if jsonBytes, err = json.Marshal(j._expr); err == nil {
			// Marshaling without indentation successful
		}
	}
	if err == nil {
		if file, err = os.Create(filePath); err == nil {
			defer func(file *os.File) {
				_ = file.Close()
			}(file)
			if _, err = file.Write(jsonBytes); err == nil {
				// Writing to file successful
			}
		}
	}
	return err
}
