package terraform

import (
	"os"
	"path/filepath"
	"testing"
)

type mockDiff map[string]interface{}

func (m mockDiff) Get(key string) interface{} {
	return m[key]
}

func (m mockDiff) GetOk(key string) (interface{}, bool) {
	val, ok := m[key]
	return val, ok
}


func TestValidateConfigurationFileInputsLogic(t *testing.T) {
	// Create a temporary directory to test existing path validation
	existingDir := t.TempDir()
	validPath := filepath.Join(existingDir, "config.txt")

	tests := []struct {
		name    string
		input   mockDiff
		wantErr bool
	}{
		{
			name: "valid create + download with path",
			input: mockDiff{
				"download_existconfig_only": false,
				"create_only":               false,
				"download_path":             validPath,
			},
			wantErr: false,
		},
		{
			name: "valid create + download with dirpath",
			input: mockDiff{
				"download_existconfig_only": false,
				"create_only":               false,
				"download_path":             existingDir,
			},
			wantErr: false,
		},
		{
			name: "valid download only with path",
			input: mockDiff{
				"download_existconfig_only": true,
				"create_only":               false,
				"download_path":             validPath,
			},
			wantErr: false,
		},
		{
			name: "valid create only without download",
			input: mockDiff{
				"download_existconfig_only": false,
				"create_only":               true,
				"download_path":             "",
			},
			wantErr: false,
		},
		{
			name: "missing path when download is triggered",
			input: mockDiff{
				"download_existconfig_only": true,
				"create_only":               false,
				"download_path":             "",
			},
			wantErr: true,
		},
		{
			name: "missing path when both create and download implied",
			input: mockDiff{
				"download_existconfig_only": false,
				"create_only":               false,
				"download_path":             "   ",
			},
			wantErr: true,
		},
		{
			name: "create_only ignored when download_only true but path is fine",
			input: mockDiff{
				"download_existconfig_only": true,
				"create_only":               true,
				"download_path":             validPath,
			},
			wantErr: false,
		},
		{
			name: "download path points to non-existent directory",
			input: mockDiff{
				"download_existconfig_only": true,
				"create_only":               false,
				"download_path":             "/nonexistent123/path/config.txt",
			},
			wantErr: true,
		},
		{
			name: "download path points to a file as dir (invalid)",
			input: func() mockDiff {
				tmpFile, err := os.CreateTemp("", "conf-file-*")
				if err != nil {
					t.Fatal(err)
				}
				defer tmpFile.Close()
				path := filepath.Join(tmpFile.Name(), "config.txt") // tmpFile.Name() is a file, not dir
				return mockDiff{
					"download_existconfig_only": true,
					"create_only":               false,
					"download_path":             path,
				}
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfigurationFileInputsLogic(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfigurationFileInputsLogic() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
