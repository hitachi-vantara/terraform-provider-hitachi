package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Predefined list of sensitive fields to mask (case-insensitive).
var fieldsToMask = []string{
	"username",
	"password",
	"authorization",
	"token",
	"secret",
	"apiKey",
}

// MaskSensitiveData accepts a JSON string, parses it, and recursively masks sensitive fields.
// It returns the masked JSON string.
func MaskSensitiveData(jsonStr string) (string, error) {
	if jsonStr == "" {
		return "", nil
	}

	// Parse the JSON string into a map
	var data interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling JSON string: %v", err)
	}

	// Recursively mask sensitive data
	maskedData := maskDataRecursively(data)

	// Convert the masked data back to a JSON string
	maskedJSON, err := json.Marshal(maskedData)
	if err != nil {
		return "", fmt.Errorf("error marshalling masked data: %v", err)
	}

	return string(maskedJSON), nil
}

// maskDataRecursively checks each key and recursively masks sensitive data in maps and slices.
func maskDataRecursively(data interface{}) interface{} {
	// If the data is a map (dictionary), recursively check each key
	if dataMap, ok := data.(map[string]interface{}); ok {
		for key, value := range dataMap {
			// Mask the value if the key matches one of the sensitive keys
			if containsSubstring(fieldsToMask, key) {
				dataMap[key] = "*******" // Mask the sensitive key
			} else {
				// Recurse if not a sensitive key
				dataMap[key] = maskDataRecursively(value)
			}
		}
		return dataMap
	}

	// If the data is a slice (list), recurse through the slice elements
	if dataSlice, ok := data.([]interface{}); ok {
		for i, value := range dataSlice {
			dataSlice[i] = maskDataRecursively(value)
		}
		return dataSlice
	}

	// If the data is a string, mask it if it matches one of the sensitive fields
	if str, ok := data.(string); ok {
		if containsSubstring(fieldsToMask, str) {
			return "*******"
		}
	}

	return data
}

// containsSubstring checks if a string contains any of the substrings in the slice, case-insensitive.
func containsSubstring(slice []string, value string) bool {
	for _, item := range slice {
		// Use case-insensitive comparison
		if strings.Contains(strings.ToLower(value), strings.ToLower(item)) {
			return true
		}
	}
	return false
}
