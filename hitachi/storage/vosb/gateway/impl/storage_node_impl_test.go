package vssbstorage

import (
	"os"
	"testing"
	// vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/gateway/model"
)

// go test -v -run TestGetAllStorageNodes
func xTestGetAll(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetStorageNodes()
	if err != nil {
		t.Errorf("Unexpected error in GetStorageNodes %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestGetAllStorageNodes
func xTestGetOne(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetStorageNode("5d991526-c48f-4c31-9a1d-40ab914915fb")
	if err != nil {
		t.Errorf("Unexpected error in GetStorageNodes %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestAddStorageNode
func xTestAddStorageNode(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if true {
		err := psm.doAddStorageNode(
			"/root/huitest2.csv",
			"",
			"vssb-789",
			"baremetal",
		)
		if err != nil {
			t.Errorf("Unexpected error in GetStorageNodes %v", err)
			return
		}
		return
	}

	resp, err := psm.GetStorageNodes()
	if err != nil {
		t.Errorf("Unexpected error in GetStorageNodes %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// Helper function to remove test files
func removeTestFile(t *testing.T, filepath string) {
	t.Helper()
	if err := os.Remove(filepath); err != nil {
		t.Logf("Warning: Failed to remove test file %s: %v", filepath, err)
	}
}

// Integration test for doAddStorageNodeAzure function (binary file support)
// go test -v -run TestDoAddStorageNodeAzure -count=1
func xTestDoAddStorageNodeAzure(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	tests := []struct {
		name                      string
		configurationFile         string
		exportedConfigurationFile string
		setupUserPassword         string
		expectedCloudProvider     string
		wantError                 bool
		description               string
	}{
		{
			name:                      "azure_with_large_binary_file",
			configurationFile:         "",
			exportedConfigurationFile: "/tmp/test_large_binary.bin",
			setupUserPassword:         "",
			expectedCloudProvider:     "azure",
			wantError:                 false,
			description:               "Test Azure deployment with large binary file (>32KB)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Running test: %s - %s", tt.name, tt.description)

			if tt.exportedConfigurationFile != "" {
				var binaryData = generateBinaryData(64 * 1024) // 64KB for large file test
				createTestBinaryFile(t, tt.exportedConfigurationFile, binaryData)
				defer removeTestFile(t, tt.exportedConfigurationFile)
			}

			// Execute the function
			err := psm.doAddStorageNodeAzure(
				tt.configurationFile,
				tt.exportedConfigurationFile,
				tt.setupUserPassword,
				tt.expectedCloudProvider,
			)

			// Validate results
			if tt.wantError && err == nil {
				t.Errorf("Expected error but got none for test case: %s", tt.name)
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error for test case %s: %v", tt.name, err)
			}

			if err == nil {
				t.Logf("Test case %s completed successfully", tt.name)
			} else {
				t.Logf("Test case %s completed with expected error: %v", tt.name, err)
			}
		})
	}
}

// Helper function to create binary test files
func createTestBinaryFile(t *testing.T, filepath string, data []byte) {
	t.Helper()
	file, err := os.Create(filepath)
	if err != nil {
		t.Fatalf("Failed to create binary test file %s: %v", filepath, err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		t.Fatalf("Failed to write to binary test file %s: %v", filepath, err)
	}
	t.Logf("Created binary test file %s with %d bytes", filepath, len(data))
}

// Helper function to generate binary test data
func generateBinaryData(size int) []byte {
	data := make([]byte, size)
	// Create a pattern with various byte values including null bytes
	for i := 0; i < size; i++ {
		data[i] = byte(i % 256)
	}
	// Add some specific binary patterns
	if size > 10 {
		// Add PNG header-like bytes
		data[0] = 0x89
		data[1] = 0x50 // 'P'
		data[2] = 0x4E // 'N'
		data[3] = 0x47 // 'G'
		data[4] = 0x0D
		data[5] = 0x0A
		data[6] = 0x1A
		data[7] = 0x0A
	}
	return data
}
