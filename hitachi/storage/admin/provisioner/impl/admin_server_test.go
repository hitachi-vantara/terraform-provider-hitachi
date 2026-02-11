package admin

import (
	model "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"

	"testing"
)

// go test -v -run TestCreateAdminServer
func xTestCreateAdminServer(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test Case 1: Create a standard server
	params := model.CreateAdminServerParams{
		ServerNickname: "test-provisioner-server-001",
		Protocol:       "FC",
		OsType:         "Linux",
		IsReserved:     false,
		OsTypeOptions:  []int{1, 2, 3},
	}

	err = psm.CreateAdminServer(params)
	if err != nil {
		t.Errorf("Unexpected error in CreateAdminServer: %v", err)
		return
	}

	// Retrieve the created server to verify creation
	listParams := model.AdminServerListParams{
		Nickname: &params.ServerNickname,
	}
	servers, err := psm.GetAdminServerList(listParams)
	if err != nil {
		t.Errorf("Failed to get server list: %v", err)
		return
	}

	if len(servers.Data) == 0 {
		t.Error("Expected to find created server, got empty list")
		return
	}

	serverInfo := servers.Data[0]
	if serverInfo.Nickname != params.ServerNickname {
		t.Errorf("Expected nickname %s, got %s", params.ServerNickname, serverInfo.Nickname)
	}

	t.Logf("Created server: ID=%d, Nickname=%s, Protocol=%s",
		serverInfo.ID, serverInfo.Nickname, serverInfo.Protocol)
}

// go test -v -run TestCreateAdminServerReserved
func xTestCreateAdminServerReserved(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test Case 2: Create a reserved server
	params := model.CreateAdminServerParams{
		ServerNickname: "reserved-provisioner-server-001",
		IsReserved:     true,
	}

	err = psm.CreateAdminServer(params)
	if err != nil {
		t.Errorf("Unexpected error in CreateAdminServer for reserved server: %v", err)
		return
	}

	// Retrieve the created server to verify creation
	listParams := model.AdminServerListParams{
		Nickname: &params.ServerNickname,
	}
	servers, err := psm.GetAdminServerList(listParams)
	if err != nil {
		t.Errorf("Failed to get server list: %v", err)
		return
	}

	if len(servers.Data) == 0 {
		t.Error("Expected to find created server, got empty list")
		return
	}

	serverInfo := servers.Data[0]
	if !serverInfo.IsReserved {
		t.Error("Expected reserved server, got non-reserved")
	}

	t.Logf("Created reserved server: ID=%d, Nickname=%s, IsReserved=%t",
		serverInfo.ID, serverInfo.Nickname, serverInfo.IsReserved)
}

// go test -v -run TestCreateAdminServerWithOptions
func xTestCreateAdminServerWithOptions(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test Case 3: Create server with OS type options
	params := model.CreateAdminServerParams{
		ServerNickname: "multi-option-server-001",
		Protocol:       "iSCSI",
		OsType:         "VMware",
		IsReserved:     false,
		OsTypeOptions:  []int{4, 5, 6, 7},
	}

	err = psm.CreateAdminServer(params)
	if err != nil {
		t.Errorf("Unexpected error in CreateAdminServer with options: %v", err)
		return
	}

	// Retrieve the created server to verify creation
	listParams := model.AdminServerListParams{
		Nickname: &params.ServerNickname,
	}
	servers, err := psm.GetAdminServerList(listParams)
	if err != nil {
		t.Errorf("Failed to get server list: %v", err)
		return
	}

	if len(servers.Data) == 0 {
		t.Error("Expected to find created server, got empty list")
		return
	}

	serverInfo := servers.Data[0]
	t.Logf("Created server with options: ID=%d, Protocol=%s, OS=%s",
		serverInfo.ID, serverInfo.Protocol, serverInfo.OsType)

	// Server created successfully with options
	t.Logf("Server created successfully with OS type options")
}

// go test -v -run TestUpdateAdminServer
func xTestUpdateAdminServer(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test Case 1: Update server nickname and OS type
	serverID := 1
	params := model.UpdateAdminServerParams{
		Nickname:      "updated-provisioner-server",
		OsType:        "Windows",
		OsTypeOptions: []int{6, 7, 8},
	}

	err = psm.UpdateAdminServer(serverID, params)
	if err != nil {
		t.Errorf("Unexpected error in UpdateAdminServer: %v", err)
		return
	}

	// Retrieve the updated server to verify changes
	serverInfo, err := psm.GetAdminServerInfo(serverID)
	if err != nil {
		t.Errorf("Failed to get updated server info: %v", err)
		return
	}

	if serverInfo.Nickname != params.Nickname {
		t.Errorf("Expected updated nickname %s, got %s", params.Nickname, serverInfo.Nickname)
	}

	if serverInfo.OsType != params.OsType {
		t.Errorf("Expected updated OS type %s, got %s", params.OsType, serverInfo.OsType)
	}

	t.Logf("Updated server: ID=%d, Nickname=%s, OS=%s",
		serverInfo.ID, serverInfo.Nickname, serverInfo.OsType)
}

// go test -v -run TestUpdateAdminServerOsOptions
func xTestUpdateAdminServerOsOptions(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test Case 2: Update server with different OS type options
	serverID := 2
	params := model.UpdateAdminServerParams{
		OsType:        "Solaris",
		OsTypeOptions: []int{8, 9, 10},
	}

	err = psm.UpdateAdminServer(serverID, params)
	if err != nil {
		t.Errorf("Unexpected error in UpdateAdminServer with OS options: %v", err)
		return
	}

	// Retrieve the updated server to verify changes
	serverInfo, err := psm.GetAdminServerInfo(serverID)
	if err != nil {
		t.Errorf("Failed to get updated server info: %v", err)
		return
	}

	t.Logf("Updated server OS options: ID=%d, OS=%s", serverInfo.ID, serverInfo.OsType)
}

// go test -v -run TestDeleteAdminServer
func xTestDeleteAdminServer(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test Case 1: Delete server with LUN config retention
	serverID := 3
	params := model.DeleteAdminServerParams{
		KeepLunConfig: true,
	}

	err = psm.DeleteAdminServer(serverID, params)
	if err != nil {
		t.Errorf("Unexpected error in DeleteAdminServer: %v", err)
		return
	}

	t.Logf("Successfully deleted server %d with LUN config retention", serverID)
}

// go test -v -run TestDeleteAdminServerComplete
func xTestDeleteAdminServerComplete(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test Case 2: Delete server completely (including LUN config)
	serverID := 4
	params := model.DeleteAdminServerParams{
		KeepLunConfig: false,
	}

	err = psm.DeleteAdminServer(serverID, params)
	if err != nil {
		t.Errorf("Unexpected error in DeleteAdminServer (complete): %v", err)
		return
	}

	t.Logf("Successfully deleted server %d completely", serverID)
}

// go test -v -run TestGetAdminServerList
func xTestGetAdminServerList(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Scenario 1: Pass serial to GetAdminServerList
	queryParams := model.AdminServerListParams{
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

// go test -v -run TestGetAdminServerListFiltered
func xTestGetAdminServerListFiltered(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test Case 2: Get servers filtered by HBA WWN
	hbaWwn := "50:01:04:f0:00:12:34:56"
	queryParams := model.AdminServerListParams{
		HbaWwn: &hbaWwn,
	}

	resp, err := psm.GetAdminServerList(queryParams)
	if err != nil {
		t.Errorf("Unexpected error in GetAdminServerList (filtered by HBA WWN): %v", err)
		return
	}

	if resp == nil {
		t.Error("Expected response, got nil")
		return
	}

	t.Logf("Filtered by HBA WWN response: Count=%d", len(resp.Data))
	for i, server := range resp.Data {
		t.Logf("Server %d: ID=%d, Nickname=%s", i, server.ID, server.Nickname)
	}
}

// go test -v -run TestGetAdminServerListByNickname
func xTestGetAdminServerListByNickname(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test Case 3: Get servers filtered by nickname
	nickname := "test-server"
	queryParams := model.AdminServerListParams{
		Nickname: &nickname,
	}

	resp, err := psm.GetAdminServerList(queryParams)
	if err != nil {
		t.Errorf("Unexpected error in GetAdminServerList (filtered by nickname): %v", err)
		return
	}

	if resp == nil {
		t.Error("Expected response, got nil")
		return
	}

	t.Logf("Filtered by nickname response: Count=%d", len(resp.Data))
	for i, server := range resp.Data {
		t.Logf("Server %d: ID=%d, Nickname=%s", i, server.ID, server.Nickname)
	}
}

// go test -v -run TestGetAdminServerInfo
func xTestGetAdminServerInfo(t *testing.T) {
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

// go test -v -run TestGetAdminServerInfoDetailed
func xTestGetAdminServerInfoDetailed(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test Case 2: Get detailed server info with capacity and path information
	serverID := 2
	serverInfo, err := psm.GetAdminServerInfo(serverID)
	if err != nil {
		t.Errorf("Unexpected error in GetAdminServerInfo (detailed): %v", err)
		return
	}

	if serverInfo == nil {
		t.Error("Expected server info, got nil")
		return
	}

	t.Logf("Detailed Server Info:")
	t.Logf("  ID: %d", serverInfo.ID)
	t.Logf("  Nickname: %s", serverInfo.Nickname)
	t.Logf("  Protocol: %s", serverInfo.Protocol)
	t.Logf("  OS Type: %s", serverInfo.OsType)
	t.Logf("  Is Reserved: %t", serverInfo.IsReserved)
	t.Logf("  Total Capacity: %d", serverInfo.TotalCapacity)
	t.Logf("  Used Capacity: %d", serverInfo.UsedCapacity)
	t.Logf("  Number of Paths: %d", serverInfo.NumberOfPaths)

	for i, path := range serverInfo.Paths {
		t.Logf("  Path %d:", i)
		t.Logf("    HBA WWN: %s", path.HbaWwn)
		t.Logf("    iSCSI Name: %s", path.IscsiName)
		t.Logf("    Port IDs: %v", path.PortIds)
	}
}

// go test -v -run TestProvisionerCRUDWorkflow
func xTestProvisionerCRUDWorkflow(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	t.Log("=== Starting Provisioner CRUD Workflow Test ===")

	// Step 1: Create a server
	t.Log("Step 1: Creating server...")
	createParams := model.CreateAdminServerParams{
		ServerNickname: "provisioner-workflow-server",
		Protocol:       "FC",
		OsType:         "AIX",
		IsReserved:     false,
		OsTypeOptions:  []int{9, 10},
	}

	err = psm.CreateAdminServer(createParams)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Find the created server to get its ID
	listParams := model.AdminServerListParams{
		Nickname: &createParams.ServerNickname,
	}
	servers, err := psm.GetAdminServerList(listParams)
	if err != nil || len(servers.Data) == 0 {
		t.Fatalf("Failed to find created server")
	}
	createdID := servers.Data[0].ID
	t.Logf("✓ Created server with ID: %d", createdID)

	// Step 2: Read the created server
	t.Log("Step 2: Reading server...")
	readInfo, err := psm.GetAdminServerInfo(createdID)
	if err != nil {
		t.Errorf("Failed to read created server: %v", err)
	} else {
		t.Logf("✓ Read server: %s (ID: %d)", readInfo.Nickname, readInfo.ID)
	}

	// Step 3: Update the server
	t.Log("Step 3: Updating server...")
	updateParams := model.UpdateAdminServerParams{
		Nickname:      "provisioner-workflow-server-updated",
		OsType:        "Solaris",
		OsTypeOptions: []int{11, 12, 13},
	}

	err = psm.UpdateAdminServer(createdID, updateParams)
	if err != nil {
		t.Errorf("Failed to update server: %v", err)
	} else {
		// Verify the update by getting server info
		updatedInfo, getErr := psm.GetAdminServerInfo(createdID)
		if getErr != nil {
			t.Errorf("Failed to get updated server info: %v", getErr)
		} else {
			t.Logf("✓ Updated server nickname to: %s", updatedInfo.Nickname)
		}
	}

	// Step 4: List servers to verify update
	t.Log("Step 4: Listing servers...")
	listParams2 := model.AdminServerListParams{}
	listResp, err := psm.GetAdminServerList(listParams2)
	if err != nil {
		t.Errorf("Failed to list servers: %v", err)
	} else {
		t.Logf("✓ Listed %d servers", len(listResp.Data))

		// Find our updated server
		for _, server := range listResp.Data {
			if server.ID == createdID {
				t.Logf("  Found our server: ID=%d, Nickname=%s", server.ID, server.Nickname)
				break
			}
		}
	}

	// Step 5: Delete the server
	t.Log("Step 5: Deleting server...")
	deleteParams := model.DeleteAdminServerParams{
		KeepLunConfig: false,
	}

	err = psm.DeleteAdminServer(createdID, deleteParams)
	if err != nil {
		t.Errorf("Failed to delete server: %v", err)
	} else {
		t.Logf("✓ Deleted server with ID: %d", createdID)
	}

	t.Log("=== Provisioner CRUD Workflow Test Completed ===")
}

// go test -v -run TestProvisionerErrorHandling
func xTestProvisionerErrorHandling(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	t.Log("=== Testing Provisioner Error Handling ===")

	// Test 1: Create server with invalid nickname
	t.Log("Test 1: Create server with empty nickname...")
	invalidParams := model.CreateAdminServerParams{
		ServerNickname: "", // Invalid: empty nickname
		Protocol:       "FC",
		OsType:         "Linux",
		IsReserved:     false,
	}

	err = psm.CreateAdminServer(invalidParams)
	if err == nil {
		t.Error("Expected error for empty nickname, got nil")
	} else {
		t.Logf("✓ Correctly caught error for empty nickname: %v", err)
	}

	// Test 2: Update non-existent server
	t.Log("Test 2: Update non-existent server...")
	nonExistentID := 99999
	updateParams := model.UpdateAdminServerParams{
		Nickname: "should-not-work",
	}

	err = psm.UpdateAdminServer(nonExistentID, updateParams)
	if err == nil {
		t.Error("Expected error for updating non-existent server, got nil")
	} else {
		t.Logf("✓ Correctly caught error for updating non-existent server: %v", err)
	}

	// Test 3: Delete non-existent server
	t.Log("Test 3: Delete non-existent server...")
	deleteParams := model.DeleteAdminServerParams{
		KeepLunConfig: false,
	}

	err = psm.DeleteAdminServer(nonExistentID, deleteParams)
	if err == nil {
		t.Logf("No error for deleting non-existent server (may be valid behavior)")
	} else {
		t.Logf("✓ Got error for deleting non-existent server: %v", err)
	}

	// Test 4: Get info for non-existent server
	t.Log("Test 4: Get info for non-existent server...")
	_, err = psm.GetAdminServerInfo(nonExistentID)
	if err == nil {
		t.Logf("No error for getting non-existent server info (may be valid behavior)")
	} else {
		t.Logf("✓ Got error for getting non-existent server info: %v", err)
	}

	t.Log("=== Provisioner Error Handling Tests Completed ===")
}

// go test -v -run TestProvisionerConcurrency
func xTestProvisionerConcurrency(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	t.Log("=== Testing Provisioner Concurrent Operations ===")

	// Test concurrent server creation
	t.Log("Testing concurrent server creation...")

	createCount := 3
	results := make(chan error, createCount)

	for i := 0; i < createCount; i++ {
		go func(index int) {
			params := model.CreateAdminServerParams{
				ServerNickname: "concurrent-server-" + string(rune('A'+index)),
				Protocol:       "iSCSI",
				OsType:         "Linux",
				IsReserved:     false,
				OsTypeOptions:  []int{index + 1, index + 2},
			}

			err := psm.CreateAdminServer(params)
			results <- err
		}(i)
	}

	// Wait for all goroutines to complete
	successCount := 0
	for i := 0; i < createCount; i++ {
		err := <-results
		if err == nil {
			successCount++
		} else {
			t.Logf("Concurrent creation %d failed: %v", i, err)
		}
	}

	t.Logf("✓ Concurrent creation results: %d/%d successful", successCount, createCount)

	// Test concurrent server listing
	t.Log("Testing concurrent server listing...")

	listCount := 5
	listResults := make(chan error, listCount)

	for i := 0; i < listCount; i++ {
		go func() {
			params := model.AdminServerListParams{}
			_, err := psm.GetAdminServerList(params)
			listResults <- err
		}()
	}

	// Wait for all list operations to complete
	listSuccessCount := 0
	for i := 0; i < listCount; i++ {
		err := <-listResults
		if err == nil {
			listSuccessCount++
		} else {
			t.Logf("Concurrent listing %d failed: %v", i, err)
		}
	}

	t.Logf("✓ Concurrent listing results: %d/%d successful", listSuccessCount, listCount)

	t.Log("=== Provisioner Concurrent Operations Tests Completed ===")
}
