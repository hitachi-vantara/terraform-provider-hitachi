package terraform

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	terraform_schema "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// // mockResourceDiff used in the newer test setup
// type mockResourceDiff struct {
// 	data map[string]interface{}
// }

// func (m *mockResourceDiff) Get(key string) interface{} {
// 	return m.data[key]
// }

func xTestSchemaValidationFuncs(t *testing.T) {
	schemaParam := terraform_schema.ResourceVssbConfigurationFileSchema["create_configuration_file_param"].Elem.(*schema.Resource).Schema
	addressSetting := schemaParam["address_setting"].Elem.(*schema.Resource).Schema

	testGroups := map[string][]struct {
		name        string
		input       interface{}
		validateFn  schema.SchemaValidateFunc
		field       string
		expectError bool
	}{
		"CloudProvider": {
			{
				name:        "invalid cloud provider",
				input:       "invalid_provider",
				validateFn:  schemaParam["expected_cloud_provider"].ValidateFunc,
				field:       "expected_cloud_provider",
				expectError: true,
			},
			{
				name:        "valid baremetal",
				input:       "baremetal",
				validateFn:  schemaParam["expected_cloud_provider"].ValidateFunc,
				field:       "expected_cloud_provider",
				expectError: false,
			},
			{
				name:        "valid aws",
				input:       "aws",
				validateFn:  schemaParam["expected_cloud_provider"].ValidateFunc,
				field:       "expected_cloud_provider",
				expectError: false,
			},
			{
				name:        "google cloud provider",
				input:       "google",
				validateFn:  schemaParam["expected_cloud_provider"].ValidateFunc,
				field:       "expected_cloud_provider",
				expectError: false,
			},
			{
				name:        "azure cloud provider",
				input:       "azure",
				validateFn:  schemaParam["expected_cloud_provider"].ValidateFunc,
				field:       "expected_cloud_provider",
				expectError: false,
			},
		},
		"ExportFileType": {
			{
				name:        "invalid Type",
				input:       "InvalidType",
				validateFn:  schemaParam["export_file_type"].ValidateFunc,
				field:       "export_file_type",
				expectError: true,
			},
			{
				name:        "valid Normal",
				input:       "Normal",
				validateFn:  schemaParam["export_file_type"].ValidateFunc,
				field:       "export_file_type",
				expectError: false,
			},
			{
				name:        "valid AddStorageNodes",
				input:       "AddStorageNodes",
				validateFn:  schemaParam["export_file_type"].ValidateFunc,
				field:       "export_file_type",
				expectError: false,
			},
			{
				name:        "valid ReplaceStorageNode",
				input:       "ReplaceStorageNode",
				validateFn:  schemaParam["export_file_type"].ValidateFunc,
				field:       "export_file_type",
				expectError: false,
			},
			{
				name:        "valid AddDrives",
				input:       "AddDrives",
				validateFn:  schemaParam["export_file_type"].ValidateFunc,
				field:       "export_file_type",
				expectError: false,
			},
			{
				name:        "valid ReplaceDrive",
				input:       "ReplaceDrive",
				validateFn:  schemaParam["export_file_type"].ValidateFunc,
				field:       "export_file_type",
				expectError: false,
			},
		},
		"DriveID": {
			{
				name:        "valid UUID",
				input:       "6f1f57a3-8b5a-4aef-9c8f-0a617e2a73d9",
				validateFn:  schemaParam["drive_id"].ValidateFunc,
				field:       "drive_id",
				expectError: false,
			},
			{
				name:        "invalid UUID",
				input:       "invalid-uuid",
				validateFn:  schemaParam["drive_id"].ValidateFunc,
				field:       "drive_id",
				expectError: true,
			},
		},
		"NodeID": {
			{
				name:        "valid UUID",
				input:       "6f1f57a3-8b5a-4aef-9c8f-0a617e2a73d9",
				validateFn:  schemaParam["node_id"].ValidateFunc,
				field:       "node_id",
				expectError: false,
			},
			{
				name:        "invalid UUID",
				input:       "invalid-uuid",
				validateFn:  schemaParam["node_id"].ValidateFunc,
				field:       "node_id",
				expectError: true,
			},
		},
		"Index": {
			{
				name:        "valid index: 1",
				input:       1,
				validateFn:  addressSetting["index"].ValidateFunc,
				field:       "index",
				expectError: false,
			},
			{
				name:        "valid index: 6",
				input:       6,
				validateFn:  addressSetting["index"].ValidateFunc,
				field:       "index",
				expectError: false,
			},
			{
				name:        "invalid index: 0 (below min)",
				input:       0,
				validateFn:  addressSetting["index"].ValidateFunc,
				field:       "index",
				expectError: true,
			},
			{
				name:        "invalid index: 7 (above max)",
				input:       7,
				validateFn:  addressSetting["index"].ValidateFunc,
				field:       "index",
				expectError: true,
			},
		},

		"NumberOfDrives": {
			{
				name:        "valid number_of_drives: 6",
				input:       6,
				validateFn:  schemaParam["number_of_drives"].ValidateFunc,
				field:       "number_of_drives",
				expectError: false,
			},
			{
				name:        "valid number_of_drives: 24",
				input:       24,
				validateFn:  schemaParam["number_of_drives"].ValidateFunc,
				field:       "number_of_drives",
				expectError: false,
			},
			{
				name:        "invalid number_of_drives: 5 (below min)",
				input:       5,
				validateFn:  schemaParam["number_of_drives"].ValidateFunc,
				field:       "number_of_drives",
				expectError: true,
			},
			{
				name:        "invalid number_of_drives: 25 (above max)",
				input:       25,
				validateFn:  schemaParam["number_of_drives"].ValidateFunc,
				field:       "number_of_drives",
				expectError: true,
			},
		},
		"IPv4Address": {
			{
				name:        "valid IPv4",
				input:       "192.168.1.1",
				validateFn:  addressSetting["control_port_ipv4_address"].ValidateFunc,
				field:       "control_port_ipv4_address",
				expectError: false,
			},
			{
				name:        "invalid IPv4",
				input:       "bad.ip.addr",
				validateFn:  addressSetting["control_port_ipv4_address"].ValidateFunc,
				field:       "control_port_ipv4_address",
				expectError: true,
			},
		},
		"IPv6Address": {
			{
				name:        "valid IPv6",
				input:       "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
				validateFn:  addressSetting["compute_port_ipv6_address"].ValidateFunc,
				field:       "compute_port_ipv6_address",
				expectError: false,
			},
			{
				name:        "invalid IPv6",
				input:       "bad::ipv6",
				validateFn:  addressSetting["compute_port_ipv6_address"].ValidateFunc,
				field:       "compute_port_ipv6_address",
				expectError: true,
			},
		},
	}

	for groupName, cases := range testGroups {
		t.Run(groupName, func(t *testing.T) {
			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					_, errs := tc.validateFn(tc.input, tc.field)
					assert.Equal(t, tc.expectError, len(errs) > 0, "input: %v, errors: %v", tc.input, errs)
				})
			}
		})
	}
}

func xTestValidateConfigurationFileInputsLogic(t *testing.T) {
	existingDir := t.TempDir()
	validPath := filepath.Join(existingDir, "config.txt")

	// Create a temp file to simulate invalid directory case
	tmpFile, err := os.CreateTemp("", "conf-file-*")
	if err != nil {
		t.Fatal(err)
	}
	defer tmpFile.Close()
	invalidDirPath := filepath.Join(tmpFile.Name(), "config.txt") // Not actually a directory

	testGroups := map[string][]struct {
		testNum     int
		name        string
		diff        *mockResourceDiff
		expectError bool
	}{
		"DownloadOnly": {
			{
				testNum: 1,
				name:    "valid create + download with path",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               false,
					"download_path":             validPath,
				}},
				expectError: false,
			},
			{
				testNum: 2,
				name:    "valid create + download with dirpath",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               false,
					"download_path":             existingDir,
				}},
				expectError: false,
			},
			{
				testNum: 3,
				name:    "valid download only with path",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": true,
					"create_only":               false,
					"download_path":             validPath,
				}},
				expectError: false,
			},
			{
				testNum: 4,
				name:    "valid create only without download",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
				}},
				expectError: false,
			},
			{
				testNum: 5,
				name:    "missing path when download is triggered",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": true,
					"create_only":               false,
					"download_path":             "",
				}},
				expectError: true,
			},
			{
				testNum: 6,
				name:    "missing path when both create and download implied",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               false,
					"download_path":             "   ",
				}},
				expectError: true,
			},
			{
				testNum: 7,
				name:    "create_only ignored when download_only true but path is fine",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": true,
					"create_only":               true,
					"download_path":             validPath,
				}},
				expectError: false,
			},
			{
				testNum: 8,
				name:    "download path points to non-existent directory",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": true,
					"create_only":               false,
					"download_path":             "/nonexistent123/path/config.txt",
				}},
				expectError: true,
			},
			{
				testNum: 9,
				name:    "download path points to a file as dir (invalid)",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": true,
					"create_only":               false,
					"download_path":             invalidDirPath,
				}},
				expectError: true,
			},
		},
		"ExpectedCloudProvider": {
			{
				testNum: 1,
				name:    "missing required cloud provider",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{},
					},
				}},
				expectError: true,
			},
			{
				testNum: 2,
				name:    "unsupported cloud provider",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "unknowncloud",
						},
					},
				}},
				expectError: true,
			},
		},
		"ExportFileType": {
			{
				testNum: 1,
				name:    "missing export_file_type",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "google",
						},
					},
				}},
				expectError: true,
			},
			{
				testNum: 2,
				name:    "unsupported export_file_type",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "azure",
							"export_file_type":        "ReplaceDrive",
						},
					},
				}},
				expectError: true,
			},
		},
		"ReplaceDrive": {
			{
				testNum: 1,
				name:    "valid ReplaceDrive with recover_single_drive=true (no drive_id)",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "google",
							"export_file_type":        "ReplaceDrive",
							"recover_single_drive":    true,
							"machine_image_id":        "ami-123",
						},
					},
				}},
				expectError: false,
			},
			{
				testNum: 2,
				name:    "invalid ReplaceDrive with recover_single_drive=false and missing drive_id",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "google",
							"export_file_type":        "ReplaceDrive",
							// "recover_single_drive": false,
						},
					},
				}},
				expectError: true,
			},
			{
				testNum: 3,
				name:    "valid ReplaceDrive with recover_single_drive=false and with drive_id",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "google",
							"export_file_type":        "ReplaceDrive",
							"recover_single_drive":    false,
							"drive_id":                "6f1f57a3-8b5a-4aef-9c8f-0a617e2a73d9",
						},
					},
				}},
				expectError: false,
			},
			{
				testNum: 4,
				name:    "valid ReplaceDrive in azure. Ignored.",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "azure",
							"export_file_type":        "ReplaceDrive",
							"recover_single_drive":    true,
							"machine_image_id":        "ami-123",
						},
					},
				}},
				expectError: true,
			},
		},
		"AddDrives": {
			{
				testNum: 1,
				name:    "valid AddDrives with 6 number_of_drives",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "azure",
							"export_file_type":        "AddDrives",
							"number_of_drives":        6,
						},
					},
				}},
				expectError: false,
			},
			{
				testNum: 2,
				name:    "valid AddDrives with 24 number_of_drives",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "google",
							"export_file_type":        "AddDrives",
							"number_of_drives":        24,
						},
					},
				}},
				expectError: false,
			},
			{
				testNum: 5,
				name:    "invalid AddDrives with no number_of_drives",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "azure",
							"export_file_type":        "AddDrives",
						},
					},
				}},
				expectError: true,
			},
		},
		"ReplaceStorageNode_Google": {
			{
				testNum: 1,
				name:    "invalid ReplaceStorageNode with missing machine_image_id",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "google",
							"export_file_type":        "ReplaceStorageNode",
						},
					},
				}},
				expectError: true,
			},
			{
				testNum: 2,
				name:    "invalid ReplaceStorageNode with missing node_id",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "google",
							"export_file_type":        "ReplaceStorageNode",
							"machine_image_id":        "ami-123",
							"recover_single_node":     false,
						},
					},
				}},
				expectError: true,
			},
			{
				testNum: 3,
				name:    "valid ReplaceStorageNode with recover_single_node and node_id",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "google",
							"export_file_type":        "ReplaceStorageNode",
							"machine_image_id":        "ami-123",
							"recover_single_node":     true,
							"node_id":                 "uuid123456",
						},
					},
				}},
				expectError: false,
			},
			{
				testNum: 6,
				name:    "valid ReplaceStorageNode with just node_id in google",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "google",
							"export_file_type":        "ReplaceStorageNode",
							"machine_image_id":        "ami-123",
							"node_id":                 "uuid123456",
						},
					},
				}},
				expectError: false,
			},
		},
		"ReplaceStorageNode_Azure": {
			{
				testNum: 1,
				name:    "invalid ReplaceStorageNode with missing machine_image_id",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "azure",
							"export_file_type":        "ReplaceStorageNode",
						},
					},
				}},
				expectError: true,
			},
			{
				testNum: 2,
				name:    "valid ReplaceStorageNode with only machine_image_id",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "azure",
							"export_file_type":        "ReplaceStorageNode",
							"machine_image_id":        "ami-123",
						},
					},
				}},
				expectError: false,
			},
			{
				testNum: 3,
				name:    "valid ReplaceStorageNode with ignored recover_single_node",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "azure",
							"export_file_type":        "ReplaceStorageNode",
							"machine_image_id":        "ami-123",
							"recover_single_node":     true, //ignored
						},
					},
				}},
				expectError: false,
			},
			{
				testNum: 4,
				name:    "valid ReplaceStorageNode with ignored node_id",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "azure",
							"export_file_type":        "ReplaceStorageNode",
							"machine_image_id":        "ami-123",
							"node_id":                 "uuid123456", //ignored
						},
					},
				}},
				expectError: false,
			},
		},
		"AddStorageNodes": {
			{
				testNum: 1,
				name:    "invalid AddStorageNodes with missing machine_image_id",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "google",
							"export_file_type":        "AddStorageNodes",
						},
					},
				}},
				expectError: true,
			},
			{
				testNum: 2,
				name:    "invalid AddStorageNodes with missing address_setting",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "google",
							"export_file_type":        "AddStorageNodes",
							"machine_image_id":        "ami-123",
						},
					},
				}},
				expectError: true,
			},
			{
				testNum: 3,
				name:    "valid AddStorageNodes with address_setting in google",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "google",
							"export_file_type":        "AddStorageNodes",
							"machine_image_id":        "ami-123",
							"address_setting": []interface{}{
								map[string]interface{}{
									"index":                       1,
									"control_port_ipv4_address":   "10.0.0.1",
									"internode_port_ipv4_address": "10.0.0.2",
									"compute_port_ipv4_address":   "10.0.0.3",
								},
							},
						},
					},
				}},
				expectError: false,
			},
			{
				testNum: 4,
				name:    "valid AddStorageNodes with address_setting in azure",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "azure",
							"export_file_type":        "AddStorageNodes",
							"machine_image_id":        "ami-123",
							"address_setting": []interface{}{
								map[string]interface{}{
									"index":                       1,
									"control_port_ipv4_address":   "10.0.0.1",
									"internode_port_ipv4_address": "10.0.0.2",
									"compute_port_ipv4_address":   "10.0.0.3",
									"compute_port_ipv6_address":   "2001:db8:85a3::8a2e:370:7334",
								},
							},
						},
					},
				}},
				expectError: false,
			},
		},
		"AddressSetting": {
			{
				testNum: 1,
				name:    "invalid AddStorageNodes with empty address_setting",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "google",
							"export_file_type":        "AddStorageNodes",
							"machine_image_id":        "ami-123",
							"address_setting":         []interface{}{},
						},
					},
				}},
				expectError: true,
			},
			{
				testNum: 2,
				name:    "invalid AddStorageNodes address_setting with missing index",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "google",
							"export_file_type":        "AddStorageNodes",
							"machine_image_id":        "ami-123",
							"address_setting": []interface{}{
								map[string]interface{}{},
							},
						},
					},
				}},
				expectError: true,
			},
			{
				testNum: 3,
				name:    "invalid AddStorageNodes address_setting with missing control_port_ipv4_address",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "google",
							"export_file_type":        "AddStorageNodes",
							"machine_image_id":        "ami-123",
							"address_setting": []interface{}{
								map[string]interface{}{
									"index": 7,
								},
							},
						},
					},
				}},
				expectError: true,
			},
			{
				testNum: 4,
				name:    "invalid AddStorageNodes address_setting with missing internode_port_ipv4_address",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "azure",
							"export_file_type":        "AddStorageNodes",
							"machine_image_id":        "ami-123",
							"address_setting": []interface{}{
								map[string]interface{}{
									"index":                     1,
									"control_port_ipv4_address": "10.0.0.1",
									// "internode_port_ipv4_address": "10.0.0.2",
									// "compute_port_ipv4_address":   "10.0.0.3",
								},
							},
						},
					},
				}},
				expectError: true,
			},
			{
				testNum: 5,
				name:    "invalid AddStorageNodes address_setting with missing compute_port_ipv4_address",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "azure",
							"export_file_type":        "AddStorageNodes",
							"machine_image_id":        "ami-123",
							"address_setting": []interface{}{
								map[string]interface{}{
									"index":                       1,
									"control_port_ipv4_address":   "10.0.0.1",
									"internode_port_ipv4_address": "10.0.0.2",
									// "compute_port_ipv4_address":   "10.0.0.3",
								},
							},
						},
					},
				}},
				expectError: true,
			},
			{
				testNum: 6,
				name:    "valid AddStorageNodes address_setting with computePortIpv6Address",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "azure",
							"export_file_type":        "AddStorageNodes",
							"machine_image_id":        "ami-123",
							"address_setting": []interface{}{
								map[string]interface{}{
									"index":                       1,
									"control_port_ipv4_address":   "10.0.0.1",
									"internode_port_ipv4_address": "10.0.0.2",
									"compute_port_ipv4_address":   "10.0.0.3",
									"compute_port_ipv6_address":   "2001:db8:85a3::8a2e:370:7334", // for Azure
								},
							},
						},
					},
				}},
				expectError: false,
			},
			{
				testNum: 7,
				name:    "valid AddStorageNodes with 6 address_settings",
				diff: &mockResourceDiff{data: map[string]interface{}{
					"download_existconfig_only": false,
					"create_only":               true,
					"download_path":             "",
					"create_configuration_file_param": []interface{}{
						map[string]interface{}{
							"expected_cloud_provider": "azure",
							"export_file_type":        "AddStorageNodes",
							"machine_image_id":        "ami-123",
							"address_setting": []interface{}{
								map[string]interface{}{
									"index":                       1,
									"control_port_ipv4_address":   "10.0.1.1",
									"internode_port_ipv4_address": "10.0.1.2",
									"compute_port_ipv4_address":   "10.0.1.3",
								},
								map[string]interface{}{
									"index":                       2,
									"control_port_ipv4_address":   "10.0.2.1",
									"internode_port_ipv4_address": "10.0.2.2",
									"compute_port_ipv4_address":   "10.0.2.3",
								},
								map[string]interface{}{
									"index":                       3,
									"control_port_ipv4_address":   "10.0.3.1",
									"internode_port_ipv4_address": "10.0.3.2",
									"compute_port_ipv4_address":   "10.0.3.3",
								},
								map[string]interface{}{
									"index":                       4,
									"control_port_ipv4_address":   "10.0.4.1",
									"internode_port_ipv4_address": "10.0.4.2",
									"compute_port_ipv4_address":   "10.0.4.3",
								},
								map[string]interface{}{
									"index":                       5,
									"control_port_ipv4_address":   "10.0.5.1",
									"internode_port_ipv4_address": "10.0.5.2",
									"compute_port_ipv4_address":   "10.0.5.3",
								},
								map[string]interface{}{
									"index":                       6,
									"control_port_ipv4_address":   "10.0.6.1",
									"internode_port_ipv4_address": "10.0.6.2",
									"compute_port_ipv4_address":   "10.0.6.3",
								},
							},
						},
					},
				}},
				expectError: false,
			},
		},
	}

	for groupName, tests := range testGroups {
		t.Run(groupName, func(t *testing.T) {
			for _, tt := range tests {
				t.Run(fmt.Sprintf("SubTest%d: %s", tt.testNum, tt.name), func(t *testing.T) {
					err := ValidateConfigurationFileInputsLogic(tt.diff)
					if tt.expectError {
						if err == nil {
							t.Errorf("FAILED: SubTest%d: expected error but got none", tt.testNum)
						} else {
							t.Logf("PASS: SubTest%d: expected error occurred: %v", tt.testNum, err)
						}
					} else {
						if err != nil {
							t.Errorf("FAILED: SubTest%d: expected no error but got: %v", tt.testNum, err)
						} else {
							t.Logf("PASS: SubTest%d: no error occurred", tt.testNum)
						}
					}
				})
			}
		})
	}
}
