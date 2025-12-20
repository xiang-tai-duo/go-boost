package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Test case 2: Load from a custom file
	testFile := "test_config.json"
	testConfig := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
		"key3": true,
	}
	tempConfig := New()
	for k, v := range testConfig {
		tempConfig.Set(k, v)
	}
	tempConfig.Save(testFile)
	defer os.Remove(testFile)

	c2 := New()
	c2.Clear()
	err := c2.Load(testFile)
	if err != nil {
		t.Errorf("Failed to load from file: %v", err)
	}
	for k, expected := range testConfig {
		actual, exists, _ := c2.Get(k)
		if !exists {
			t.Errorf("For key %q, expected %v, got nil", k, expected)
			continue
		}
		// Special handling for key2 which is a number
		if k == "key2" {
			expectedInt, ok1 := expected.(int)
			actualFloat, ok2 := actual.(float64)
			if ok1 && ok2 {
				if int(actualFloat) != expectedInt {
					t.Errorf("For key %q, expected %v, got %v", k, expected, actual)
				}
			} else {
				t.Errorf("For key %q, expected type %T, got type %T", k, expected, actual)
			}
		} else {
			if actual != expected {
				t.Errorf("For key %q, expected %v, got %v", k, expected, actual)
			}
		}
	}

	type TestStruct struct {
		Key1 string `json:"key1"`
		Key2 int    `json:"key2"`
		Key3 bool   `json:"key3"`
	}
	testStruct := TestStruct{
		Key1: "value1",
		Key2: 123,
		Key3: true,
	}

	c3 := New()
	c3.Clear()
	err = c3.Load(testStruct)
	if err != nil {
		t.Errorf("Failed to load from struct: %v", err)
	}

	expectedKey1, exists, _ := c3.Get("key1")
	if !exists || expectedKey1 != "value1" {
		t.Errorf("Expected key1 to be 'value1', got %v", expectedKey1)
	}

	expectedKey2, exists, _ := c3.Get("key2")
	if !exists || expectedKey2 != float64(123) {
		t.Errorf("Expected key2 to be 123, got %v", expectedKey2)
	}

	expectedKey3, exists, _ := c3.Get("key3")
	if !exists || expectedKey3 != true {
		t.Errorf("Expected key3 to be true, got %v", expectedKey3)
	}
}
