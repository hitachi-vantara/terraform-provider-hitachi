package terraform

import (
	"strings"
	"testing"
)

// --- mock for schema.ResourceDiff ---
type mockDiff struct {
	id     string
	values map[string]interface{}
}

func (m *mockDiff) Get(key string) interface{} {
	return m.values[key]
}

func (m *mockDiff) GetOk(key string) (interface{}, bool) {
	v, ok := m.values[key]
	return v, ok
}

func (m *mockDiff) Id() string {
	return m.id
}

// go test -v -run TestValidateVolumeIDValues
func TestValidateVolumeIDValues(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		values    map[string]interface{}
		expectErr bool
	}{
		{
			name:      "Create with no volume_id — OK",
			id:        "",
			values:    map[string]interface{}{},
			expectErr: false,
		},
		{
			name:      "Create with volume_id — Error",
			id:        "",
			values:    map[string]interface{}{"volume_id": 5},
			expectErr: true,
		},
		{
			name:      "Update missing volume_id — Error",
			id:        "100,101",
			values:    map[string]interface{}{},
			expectErr: true,
		},
		{
			name:      "Update with volume_id not in state — Error",
			id:        "100,101",
			values:    map[string]interface{}{"volume_id": 200},
			expectErr: true,
		},
		{
			name:      "Update with valid volume_id — OK",
			id:        "100,101",
			values:    map[string]interface{}{"volume_id": 100},
			expectErr: false,
		},
		{
			name:      "Both volume_id and number_of_volumes set — Error",
			id:        "100",
			values:    map[string]interface{}{"volume_id": 100, "number_of_volumes": 2},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &mockDiff{id: tt.id, values: tt.values}
			err := ValidateVolumeIDValues(d)
			if (err != nil) != tt.expectErr {
				t.Errorf("[%s] expected error=%v, got %v", tt.name, tt.expectErr, err)
			}
		})
	}
}

// go test -v -run TestValidateNicknameParamValues
func TestValidateNicknameParamValues(t *testing.T) {
	tests := []struct {
		name      string
		values    map[string]interface{}
		expectErr bool
	}{
		{
			name: "Valid base_name only",
			values: map[string]interface{}{
				"nickname_param": []interface{}{
					map[string]interface{}{
						"base_name":       "DATA",
						"start_number":    -1,
						"number_of_digits": 0,
					},
				},
			},
			expectErr: false,
		},
		{
			name: "number_of_digits set but start_number missing",
			values: map[string]interface{}{
				"nickname_param": []interface{}{
					map[string]interface{}{
						"base_name":       "VOL",
						"start_number":    -1,
						"number_of_digits": 3,
					},
				},
			},
			expectErr: true,
		},
		{
			name: "Nickname too long",
			values: map[string]interface{}{
				"nickname_param": []interface{}{
					map[string]interface{}{
						"base_name":       strings.Repeat("A", 31),
						"start_number":    0,
						"number_of_digits": 3,
					},
				},
			},
			expectErr: true,
		},
		{
			name: "Valid nickname with suffix digits",
			values: map[string]interface{}{
				"nickname_param": []interface{}{
					map[string]interface{}{
						"base_name":       "DATA",
						"start_number":    0,
						"number_of_digits": 3,
					},
				},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &mockDiff{values: tt.values}
			err := ValidateNicknameParamValues(d)
			if (err != nil) != tt.expectErr {
				t.Errorf("[%s] expected error=%v, got %v", tt.name, tt.expectErr, err)
			}
		})
	}
}

// go test -v -run TestValidateDataReductionSettingsValues
func TestValidateDataReductionSettingsValues(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		values    map[string]interface{}
		expectErr bool
	}{
		{
			name: "Valid create — compression_acceleration not set",
			id:   "",
			values: map[string]interface{}{
				"capacity_saving":               "COMPRESSION",
				"is_data_reduction_share_enabled": false,
			},
			expectErr: false,
		},
		{
			name: "Create with compression_acceleration — Error",
			id:   "",
			values: map[string]interface{}{
				"compression_acceleration": true,
				"capacity_saving":          "COMPRESSION",
			},
			expectErr: true,
		},
		{
			name: "Update with DISABLE and compression_acceleration — Error",
			id:   "100",
			values: map[string]interface{}{
				"compression_acceleration": true,
				"capacity_saving":          "DISABLE",
			},
			expectErr: true,
		},
		{
			name: "Update with valid compression_acceleration — OK",
			id:   "100",
			values: map[string]interface{}{
				"compression_acceleration": true,
				"capacity_saving":          "COMPRESSION",
			},
			expectErr: false,
		},
		{
			name: "is_data_reduction_share_enabled true but DISABLE — Error",
			id:   "",
			values: map[string]interface{}{
				"capacity_saving":               "DISABLE",
				"is_data_reduction_share_enabled": true,
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &mockDiff{id: tt.id, values: tt.values}
			err := ValidateDataReductionSettingsValues(d)
			if (err != nil) != tt.expectErr {
				t.Errorf("[%s] expected error=%v, got %v", tt.name, tt.expectErr, err)
			}
		})
	}
}
