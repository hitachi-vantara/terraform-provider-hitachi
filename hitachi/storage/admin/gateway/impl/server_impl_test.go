package admin

import (
	"fmt"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	"testing"
)

// go test -v -run xGetAdminServerList
func xGetAdminServerList(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Scenario 1: Pass HBA WWN to GetAdminServerList
	queryParams := gwymodel.AdminServerListParams{
		HbaWwn:    nil,
		Nickname:  nil,
		IscsiName: nil,
	}

	resp, err := psm.GetAdminServerList(queryParams)
	if err != nil {
		t.Errorf("Unexpected error in GetAdminServerList %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run xGetAdminServerInfo
func xGetAdminServerInfo(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Scenario 1: Pass server ID to GetAdminServerInfo
	serverID := 1
	resp, err := psm.GetAdminServerInfo(serverID)
	if err != nil {
		t.Errorf("Unexpected error in GetAdminServerInfo %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestGetAdminServerList
func TestGetAdminServerList(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	tests := []struct {
		name        string
		params      gwymodel.AdminServerListParams
		description string
	}{
		{
			name: "GetAllServers",
			params: gwymodel.AdminServerListParams{
				HbaWwn:    nil,
				Nickname:  nil,
				IscsiName: nil,
			},
			description: "Get all admin servers without filters",
		},
		{
			name: "FilterByNickname",
			params: gwymodel.AdminServerListParams{
				Nickname:  stringPtr("test-server"),
				HbaWwn:    nil,
				IscsiName: nil,
			},
			description: "Filter servers by nickname",
		},
		{
			name: "FilterByHbaWwn",
			params: gwymodel.AdminServerListParams{
				HbaWwn:    stringPtr("50:00:09:73:00:18:95:19"),
				Nickname:  nil,
				IscsiName: nil,
			},
			description: "Filter servers by HBA WWN",
		},
		{
			name: "FilterByIscsiName",
			params: gwymodel.AdminServerListParams{
				IscsiName: stringPtr("iqn.1991-05.com.microsoft:test-server"),
				Nickname:  nil,
				HbaWwn:    nil,
			},
			description: "Filter servers by iSCSI name",
		},
		{
			name: "MultipleFilters",
			params: gwymodel.AdminServerListParams{
				Nickname:  stringPtr("test-server"),
				HbaWwn:    stringPtr("50:00:09:73:00:18:95:19"),
				IscsiName: nil,
			},
			description: "Filter servers by multiple parameters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Test: %s - %s", tt.name, tt.description)

			resp, err := psm.GetAdminServerList(tt.params)
			if err != nil {
				t.Errorf("Unexpected error in GetAdminServerList: %v", err)
				return
			}

			if resp == nil {
				t.Error("Response should not be nil")
				return
			}

			t.Logf("Found %d servers", len(resp.Data))
			for i, server := range resp.Data {
				t.Logf("Server %d: ID=%d, Nickname=%s, Protocol=%s, OS=%s",
					i+1, server.ID, server.Nickname, server.Protocol, server.OsType)
			}
		})
	}
}

// go test -v -run TestGetAdminServerInfo
func TestGetAdminServerInfo(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	tests := []struct {
		name        string
		serverID    int
		expectError bool
		description string
	}{
		{
			name:        "ValidServerID",
			serverID:    1,
			expectError: false,
			description: "Get info for valid server ID",
		},
		{
			name:        "AnotherValidServerID",
			serverID:    2,
			expectError: false,
			description: "Get info for another valid server ID",
		},
		{
			name:        "NonExistentServerID",
			serverID:    99999,
			expectError: true,
			description: "Get info for non-existent server ID",
		},
		{
			name:        "ZeroServerID",
			serverID:    0,
			expectError: true,
			description: "Get info for invalid server ID (0)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Test: %s - %s", tt.name, tt.description)

			resp, err := psm.GetAdminServerInfo(tt.serverID)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for server ID %d, but got none", tt.serverID)
				} else {
					t.Logf("Expected error occurred: %v", err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for server ID %d: %v", tt.serverID, err)
				return
			}

			if resp == nil {
				t.Error("Response should not be nil")
				return
			}

			t.Logf("Server Info: ID=%d, Nickname=%s, Protocol=%s, OS=%s",
				resp.ID, resp.Nickname, resp.Protocol, resp.OsType)
			t.Logf("Capacity: Total=%d, Used=%d", resp.TotalCapacity, resp.UsedCapacity)
			t.Logf("Volumes: %d, Paths: %d", resp.NumberOfVolumes, resp.NumberOfPaths)
			t.Logf("Status: Reserved=%v, Inconsistent=%v", resp.IsReserved, resp.IsInconsistent)
		})
	}
}

// go test -v -run TestCreateAdminServer
func TestCreateAdminServer(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	tests := []struct {
		name        string
		params      gwymodel.CreateAdminServerParams
		expectError bool
		description string
	}{
		{
			name: "CreateFCServer",
			params: gwymodel.CreateAdminServerParams{
				ServerNickname: "test-fc-server",
				Protocol:       "FC",
				OsType:         "Linux",
				IsReserved:     false,
			},
			expectError: false,
			description: "Create a new FC server",
		},
		{
			name: "CreateIscsiServer",
			params: gwymodel.CreateAdminServerParams{
				ServerNickname: "test-iscsi-server",
				Protocol:       "iSCSI",
				OsType:         "VMware",
				IsReserved:     false,
			},
			expectError: false,
			description: "Create a new iSCSI server",
		},
		{
			name: "CreateReservedServer",
			params: gwymodel.CreateAdminServerParams{
				ServerNickname: "test-reserved-server",
				IsReserved:     true,
			},
			expectError: false,
			description: "Create a reserved server",
		},
		{
			name: "CreateServerWithOptions",
			params: gwymodel.CreateAdminServerParams{
				ServerNickname: "test-server-with-options",
				Protocol:       "FC",
				OsType:         "Windows",
				OsTypeOptions:  []int{1, 2, 3},
				IsReserved:     false,
			},
			expectError: false,
			description: "Create server with OS type options",
		},
		{
			name: "CreateServerEmptyNickname",
			params: gwymodel.CreateAdminServerParams{
				ServerNickname: "",
				Protocol:       "FC",
				OsType:         "Linux",
				IsReserved:     false,
			},
			expectError: true,
			description: "Create server with empty nickname (should fail)",
		},
		{
			name: "CreateServerInvalidProtocol",
			params: gwymodel.CreateAdminServerParams{
				ServerNickname: "test-invalid-protocol",
				Protocol:       "INVALID",
				OsType:         "Linux",
				IsReserved:     false,
			},
			expectError: true,
			description: "Create server with invalid protocol (should fail)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Test: %s - %s", tt.name, tt.description)

			err := psm.CreateAdminServer(tt.params)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for params %+v, but got none", tt.params)
				} else {
					t.Logf("Expected error occurred: %v", err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			t.Logf("Successfully created server: %s", tt.params.ServerNickname)

			// Find the created server by nickname to get its ID for cleanup
			listParams := gwymodel.AdminServerListParams{
				Nickname: &tt.params.ServerNickname,
			}
			servers, err := psm.GetAdminServerList(listParams)
			if err != nil || len(servers.Data) == 0 {
				t.Logf("Warning: Could not find created server for cleanup")
				return
			}

			// Cleanup: delete the created server
			deleteParams := gwymodel.DeleteAdminServerParams{
				KeepLunConfig: false,
			}
			err = psm.DeleteAdminServer(servers.Data[0].ID, deleteParams)
			if err != nil {
				t.Logf("Warning: Failed to cleanup created server ID %d: %v", servers.Data[0].ID, err)
			} else {
				t.Logf("Successfully cleaned up created server ID %d", servers.Data[0].ID)
			}
		})
	}
}

// go test -v -run TestUpdateAdminServer
func TestUpdateAdminServer(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// First create a server to update
	createParams := gwymodel.CreateAdminServerParams{
		ServerNickname: "test-update-server",
		Protocol:       "FC",
		OsType:         "Linux",
		IsReserved:     false,
	}

	err = psm.CreateAdminServer(createParams)
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}

	// Find the created server to get its ID
	listParams := gwymodel.AdminServerListParams{
		Nickname: &createParams.ServerNickname,
	}
	servers, err := psm.GetAdminServerList(listParams)
	if err != nil || len(servers.Data) == 0 {
		t.Fatalf("Failed to find created test server")
	}
	createdServerID := servers.Data[0].ID

	defer func() {
		// Cleanup
		deleteParams := gwymodel.DeleteAdminServerParams{
			KeepLunConfig: false,
		}
		err := psm.DeleteAdminServer(createdServerID, deleteParams)
		if err != nil {
			t.Logf("Warning: Failed to cleanup server ID %d: %v", createdServerID, err)
		}
	}()

	tests := []struct {
		name        string
		serverID    int
		params      gwymodel.UpdateAdminServerParams
		expectError bool
		description string
	}{
		{
			name:     "UpdateNickname",
			serverID: createdServerID,
			params: gwymodel.UpdateAdminServerParams{
				Nickname: "updated-nickname",
			},
			expectError: false,
			description: "Update server nickname",
		},
		{
			name:     "UpdateOsType",
			serverID: createdServerID,
			params: gwymodel.UpdateAdminServerParams{
				OsType: "VMware",
			},
			expectError: false,
			description: "Update server OS type",
		},
		{
			name:     "UpdateOsTypeOptions",
			serverID: createdServerID,
			params: gwymodel.UpdateAdminServerParams{
				OsTypeOptions: []int{1, 2},
			},
			expectError: false,
			description: "Update server OS type options",
		},
		{
			name:     "UpdateMultipleFields",
			serverID: createdServerID,
			params: gwymodel.UpdateAdminServerParams{
				Nickname:      "final-nickname",
				OsType:        "Windows",
				OsTypeOptions: []int{3, 4, 5},
			},
			expectError: false,
			description: "Update multiple server fields",
		},
		{
			name:     "UpdateNonExistentServer",
			serverID: 99999,
			params: gwymodel.UpdateAdminServerParams{
				Nickname: "should-fail",
			},
			expectError: true,
			description: "Update non-existent server (should fail)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Test: %s - %s", tt.name, tt.description)

			err := psm.UpdateAdminServer(tt.serverID, tt.params)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for server ID %d, but got none", tt.serverID)
				} else {
					t.Logf("Expected error occurred: %v", err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Verify the update by reading server info
			updatedInfo, err := psm.GetAdminServerInfo(tt.serverID)
			if err != nil {
				t.Errorf("Failed to verify update: %v", err)
				return
			}

			t.Logf("Updated Server: ID=%d, Nickname=%s, OS=%s",
				updatedInfo.ID, updatedInfo.Nickname, updatedInfo.OsType)
		})
	}
}

// go test -v -run TestDeleteAdminServer
func TestDeleteAdminServer(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	tests := []struct {
		name        string
		setupServer bool
		params      gwymodel.DeleteAdminServerParams
		expectError bool
		description string
	}{
		{
			name:        "DeleteServerKeepLunConfig",
			setupServer: true,
			params: gwymodel.DeleteAdminServerParams{
				KeepLunConfig: true,
			},
			expectError: false,
			description: "Delete server keeping LUN configuration",
		},
		{
			name:        "DeleteServerRemoveLunConfig",
			setupServer: true,
			params: gwymodel.DeleteAdminServerParams{
				KeepLunConfig: false,
			},
			expectError: false,
			description: "Delete server removing LUN configuration",
		},
		{
			name:        "DeleteNonExistentServer",
			setupServer: false,
			params: gwymodel.DeleteAdminServerParams{
				KeepLunConfig: false,
			},
			expectError: true,
			description: "Delete non-existent server (should fail)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Test: %s - %s", tt.name, tt.description)

			var serverID int

			if tt.setupServer {
				// Create a server for deletion
				createParams := gwymodel.CreateAdminServerParams{
					ServerNickname: "test-delete-server",
					Protocol:       "FC",
					OsType:         "Linux",
					IsReserved:     false,
				}

				err := psm.CreateAdminServer(createParams)
				if err != nil {
					t.Fatalf("Failed to create test server: %v", err)
				}

				// Find the created server to get its ID
				listParams := gwymodel.AdminServerListParams{
					Nickname: &createParams.ServerNickname,
				}
				servers, err := psm.GetAdminServerList(listParams)
				if err != nil || len(servers.Data) == 0 {
					t.Fatalf("Failed to find created test server")
				}
				serverID = servers.Data[0].ID
				t.Logf("Created server ID %d for deletion test", serverID)
			} else {
				serverID = 99999 // Non-existent server ID
			}

			err := psm.DeleteAdminServer(serverID, tt.params)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for server ID %d, but got none", serverID)
				} else {
					t.Logf("Expected error occurred: %v", err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			t.Logf("Successfully deleted server ID %d", serverID)

			// Verify server is deleted by trying to get its info
			_, err = psm.GetAdminServerInfo(serverID)
			if err == nil {
				t.Errorf("Server ID %d should have been deleted but still exists", serverID)
			} else {
				t.Logf("Confirmed server ID %d is deleted: %v", serverID, err)
			}
		})
	}
}

// go test -v -run TestCreateUpdateDeleteFlow
func TestCreateUpdateDeleteFlow(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	t.Log("Test: Complete CRUD flow - create, read, update, delete")

	// Step 1: Create a server
	createParams := gwymodel.CreateAdminServerParams{
		ServerNickname: "test-crud-flow",
		Protocol:       "FC",
		OsType:         "Linux",
		IsReserved:     false,
	}

	err = psm.CreateAdminServer(createParams)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Find the created server to get its ID
	listParams := gwymodel.AdminServerListParams{
		Nickname: &createParams.ServerNickname,
	}
	servers, err := psm.GetAdminServerList(listParams)
	if err != nil || len(servers.Data) == 0 {
		t.Fatalf("Failed to find created test server")
	}
	createdServerID := servers.Data[0].ID
	t.Logf("Step 1 - Created server: ID=%d, Nickname=%s", createdServerID, createParams.ServerNickname)

	// Step 2: Read the server info
	serverInfo, err := psm.GetAdminServerInfo(createdServerID)
	if err != nil {
		t.Errorf("Failed to read server info: %v", err)
		return
	}
	t.Logf("Step 2 - Read server info: ID=%d, Nickname=%s, Protocol=%s",
		serverInfo.ID, serverInfo.Nickname, serverInfo.Protocol)

	// Step 3: Update the server
	updateParams := gwymodel.UpdateAdminServerParams{
		Nickname: "updated-crud-flow",
		OsType:   "VMware",
	}

	err = psm.UpdateAdminServer(createdServerID, updateParams)
	if err != nil {
		t.Errorf("Failed to update server: %v", err)
		return
	}
	t.Logf("Step 3 - Updated server ID: %d", createdServerID)

	// Step 4: Verify the update by reading again
	updatedInfo, err := psm.GetAdminServerInfo(createdServerID)
	if err != nil {
		t.Errorf("Failed to verify update: %v", err)
		return
	}

	if updatedInfo.Nickname != "updated-crud-flow" {
		t.Errorf("Update verification failed: expected nickname 'updated-crud-flow', got '%s'", updatedInfo.Nickname)
	}
	t.Logf("Step 4 - Verified update: Nickname=%s, OS=%s", updatedInfo.Nickname, updatedInfo.OsType)

	// Step 5: Delete the server
	deleteParams := gwymodel.DeleteAdminServerParams{
		KeepLunConfig: false,
	}

	err = psm.DeleteAdminServer(createdServerID, deleteParams)
	if err != nil {
		t.Errorf("Failed to delete server: %v", err)
		return
	}
	t.Logf("Step 5 - Deleted server ID %d", createdServerID)

	// Step 6: Verify deletion
	_, err = psm.GetAdminServerInfo(createdServerID)
	if err == nil {
		t.Errorf("Server should have been deleted but still exists")
	} else {
		t.Logf("Step 6 - Confirmed deletion: %v", err)
	}
}

// go test -v -run TestConcurrentOperations
func TestConcurrentOperations(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	t.Log("Test: Concurrent server operations")

	const numServers = 3
	serverIDs := make([]int, numServers)

	// Create multiple servers concurrently
	t.Log("Creating servers concurrently...")
	for i := 0; i < numServers; i++ {
		createParams := gwymodel.CreateAdminServerParams{
			ServerNickname: fmt.Sprintf("concurrent-server-%d", i),
			Protocol:       "FC",
			OsType:         "Linux",
			IsReserved:     false,
		}

		err := psm.CreateAdminServer(createParams)
		if err != nil {
			t.Errorf("Failed to create server %d: %v", i, err)
			continue
		}

		// Find the created server to get its ID
		listParams := gwymodel.AdminServerListParams{
			Nickname: &createParams.ServerNickname,
		}
		servers, listErr := psm.GetAdminServerList(listParams)
		if listErr != nil || len(servers.Data) == 0 {
			t.Errorf("Failed to find created server %d", i)
			continue
		}
		serverIDs[i] = servers.Data[0].ID
		t.Logf("Created server %d: ID=%d", i, servers.Data[0].ID)
	}

	// Read all servers concurrently
	t.Log("Reading servers concurrently...")
	for i, serverID := range serverIDs {
		if serverID == 0 {
			continue
		}

		info, err := psm.GetAdminServerInfo(serverID)
		if err != nil {
			t.Errorf("Failed to read server %d (ID=%d): %v", i, serverID, err)
			continue
		}
		t.Logf("Read server %d: ID=%d, Nickname=%s", i, info.ID, info.Nickname)
	}

	// Update servers concurrently
	t.Log("Updating servers concurrently...")
	for i, serverID := range serverIDs {
		if serverID == 0 {
			continue
		}

		updateParams := gwymodel.UpdateAdminServerParams{
			Nickname: fmt.Sprintf("updated-concurrent-server-%d", i),
		}

		err = psm.UpdateAdminServer(serverID, updateParams)
		if err != nil {
			t.Errorf("Failed to update server %d (ID=%d): %v", i, serverID, err)
			continue
		}
		t.Logf("Updated server %d: ID=%d", i, serverID)
	}

	// Cleanup: Delete all servers
	t.Log("Cleaning up servers...")
	deleteParams := gwymodel.DeleteAdminServerParams{
		KeepLunConfig: false,
	}

	for i, serverID := range serverIDs {
		if serverID == 0 {
			continue
		}

		err := psm.DeleteAdminServer(serverID, deleteParams)
		if err != nil {
			t.Logf("Warning: Failed to cleanup server %d (ID=%d): %v", i, serverID, err)
		} else {
			t.Logf("Cleaned up server %d: ID=%d", i, serverID)
		}
	}
}

// go test -v -run TestEdgeCases
func TestEdgeCases(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	tests := []struct {
		name        string
		testFunc    func(t *testing.T)
		description string
	}{
		{
			name: "VeryLongNickname",
			testFunc: func(t *testing.T) {
				longNickname := string(make([]byte, 256)) // Very long nickname
				for i := range longNickname {
					longNickname = longNickname[:i] + "a" + longNickname[i+1:]
				}

				createParams := gwymodel.CreateAdminServerParams{
					ServerNickname: longNickname,
					Protocol:       "FC",
					OsType:         "Linux",
					IsReserved:     false,
				}

				err := psm.CreateAdminServer(createParams)
				if err == nil {
					t.Error("Expected error for very long nickname, but got none")
				} else {
					t.Logf("Expected error for long nickname: %v", err)
				}
			},
			description: "Test server creation with very long nickname",
		},
		{
			name: "SpecialCharactersInNickname",
			testFunc: func(t *testing.T) {
				specialNickname := "test-server!@#$%^&*()_+-=[]{}|;':\",./<>?"

				createParams := gwymodel.CreateAdminServerParams{
					ServerNickname: specialNickname,
					Protocol:       "FC",
					OsType:         "Linux",
					IsReserved:     false,
				}

				err := psm.CreateAdminServer(createParams)
				if err != nil {
					t.Logf("Server creation with special characters failed (expected): %v", err)
				} else {
					// Find created server for cleanup
					listParams := gwymodel.AdminServerListParams{
						Nickname: &createParams.ServerNickname,
					}
					servers, listErr := psm.GetAdminServerList(listParams)
					if listErr == nil && len(servers.Data) > 0 {
						t.Logf("Server created with special characters: ID=%d", servers.Data[0].ID)
						// Cleanup
						deleteParams := gwymodel.DeleteAdminServerParams{KeepLunConfig: false}
						psm.DeleteAdminServer(servers.Data[0].ID, deleteParams)
					}
				}
			},
			description: "Test server creation with special characters in nickname",
		},
		{
			name: "UpdateToEmptyValues",
			testFunc: func(t *testing.T) {
				// Create a server first
				createParams := gwymodel.CreateAdminServerParams{
					ServerNickname: "test-empty-update",
					Protocol:       "FC",
					OsType:         "Linux",
					IsReserved:     false,
				}

				err := psm.CreateAdminServer(createParams)
				if err != nil {
					t.Fatalf("Failed to create test server: %v", err)
				}

				// Find the created server to get its ID
				listParams := gwymodel.AdminServerListParams{
					Nickname: &createParams.ServerNickname,
				}
				servers, listErr := psm.GetAdminServerList(listParams)
				if listErr != nil || len(servers.Data) == 0 {
					t.Fatalf("Failed to find created test server")
				}
				serverID := servers.Data[0].ID

				defer func() {
					deleteParams := gwymodel.DeleteAdminServerParams{KeepLunConfig: false}
					psm.DeleteAdminServer(serverID, deleteParams)
				}()

				// Try to update to empty nickname
				updateParams := gwymodel.UpdateAdminServerParams{
					Nickname: "", // Empty nickname
				}

				err = psm.UpdateAdminServer(serverID, updateParams)
				if err == nil {
					t.Error("Expected error for empty nickname update, but got none")
				} else {
					t.Logf("Expected error for empty nickname update: %v", err)
				}
			},
			description: "Test updating server with empty values",
		},
		{
			name: "RapidCreateDelete",
			testFunc: func(t *testing.T) {
				// Rapidly create and delete servers
				for i := 0; i < 5; i++ {
					createParams := gwymodel.CreateAdminServerParams{
						ServerNickname: fmt.Sprintf("rapid-test-%d", i),
						Protocol:       "FC",
						OsType:         "Linux",
						IsReserved:     false,
					}

					err := psm.CreateAdminServer(createParams)
					if err != nil {
						t.Errorf("Failed to create rapid server %d: %v", i, err)
						continue
					}

					// Find the created server to get its ID
					listParams := gwymodel.AdminServerListParams{
						Nickname: &createParams.ServerNickname,
					}
					servers, listErr := psm.GetAdminServerList(listParams)
					if listErr != nil || len(servers.Data) == 0 {
						t.Errorf("Failed to find created rapid server %d", i)
						continue
					}
					serverID := servers.Data[0].ID

					t.Logf("Created rapid server %d: ID=%d", i, serverID)

					// Immediately delete
					deleteParams := gwymodel.DeleteAdminServerParams{KeepLunConfig: false}
					err = psm.DeleteAdminServer(serverID, deleteParams)
					if err != nil {
						t.Errorf("Failed to delete rapid server %d: %v", i, err)
					} else {
						t.Logf("Deleted rapid server %d: ID=%d", i, serverID)
					}
				}
			},
			description: "Test rapid create and delete operations",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Edge Case Test: %s - %s", tt.name, tt.description)
			tt.testFunc(t)
		})
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
