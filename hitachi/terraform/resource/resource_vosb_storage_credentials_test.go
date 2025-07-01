package terraform

import (
	// "fmt"
	"strings"
	"testing"
)

// mockResourceDiff implements only the methods required by minimalDiff
type mockResourceDiff struct {
	data map[string]interface{}
}

func (m *mockResourceDiff) Get(key string) interface{} {
	return m.data[key]
}

func (m *mockResourceDiff) GetOk(key string) (interface{}, bool) {
	val, ok := m.data[key]
	return val, ok
}

func xTestValidatePasswordChangeInputs(t *testing.T) {
	tests := []struct {
		name        string
		data        map[string]interface{}
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid input",
			data: map[string]interface{}{
				"current_password": "abc123!",
				"new_password":     "xyz456!",
				"user_id":          "user_abc",
			},
		},
		{
			name: "missing user_id",
			data: map[string]interface{}{
				"current_password": "abc123!",
				"new_password":     "xyz456!",
				"user_id":          "",
			},
			expectError: true,
			errorMsg:    "missing user_id",
		},
		{
			name: "passwords are the same",
			data: map[string]interface{}{
				"current_password": "abc123!",
				"new_password":     "abc123!",
				"user_id":          "user_xyz",
			},
			expectError: true,
			errorMsg:    "new_password must be different from current_password",
		},
		{
			name: "user_id too short",
			data: map[string]interface{}{
				"current_password": "abc123!",
				"new_password":     "xyz456!",
				"user_id":          "usr",
			},
			expectError: true,
			errorMsg:    "user_id must be 5 to 255 valid characters",
		},
		{
			name: "user_id too long",
			data: map[string]interface{}{
				"current_password": "abc123!",
				"new_password":     "xyz456!",
				"user_id":          strings.Repeat("a", 256),
			},
			expectError: true,
			errorMsg:    "user_id must be 5 to 255 valid characters",
		},
		{
			name: "user_id invalid characters",
			data: map[string]interface{}{
				"current_password": "abc123!",
				"new_password":     "xyz456!",
				"user_id":          "user$%^&*()", // contains disallowed characters
			},
			expectError: true,
			errorMsg:    "user_id must be 5 to 255 valid characters",
		},
		{
			name: "current_password too long",
			data: map[string]interface{}{
				"current_password": strings.Repeat("a", 257),
				"new_password":     "xyz456!",
				"user_id":          "validUser123",
			},
			expectError: true,
			errorMsg:    "current_password must be 1 to 256 valid characters",
		},
		{
			name: "new_password invalid characters",
			data: map[string]interface{}{
				"current_password": "abc123!",
				"new_password":     "pass with spaces", // space is not in allowed set
				"user_id":          "validUser123",
			},
			expectError: true,
			errorMsg:    "new_password must be 1 to 256 valid characters",
		},
		{
			name: "valid passwords with full symbol set",
			data: map[string]interface{}{
				"current_password": `!#$%&"'()*+,./:;<=>?@[\]^_{|}~`,
				"new_password":     `Aa1!@#_+=~`,
				"user_id":          "user@ok",
			},
		},
		{
			name: "empty current_password",
			data: map[string]interface{}{
				"current_password": "",
				"new_password":     "valid123!",
				"user_id":          "user123",
			},
			expectError: true,
			errorMsg:    "current_password must be 1 to 256 valid characters",
		},
		{
			name: "empty new_password",
			data: map[string]interface{}{
				"current_password": "valid123!",
				"new_password":     "",
				"user_id":          "user123",
			},
			expectError: true,
			errorMsg:    "new_password must be 1 to 256 valid characters",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			diff := &mockResourceDiff{data: tc.data}
			err := validatePasswordChangeInputsLogic(diff)

			if tc.expectError {
				if err == nil || err.Error() != tc.errorMsg {
					t.Errorf("expected error %q, got %v", tc.errorMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
