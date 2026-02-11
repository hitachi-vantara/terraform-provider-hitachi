package admin

import (
	"testing"

	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
)

// go test -v -run TestReconcileCreateAdminServer
func xTestReconcileCreateAdminServer(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test Case 1: Create a standard server
	params := gwymodel.CreateAdminServerParams{
		ServerNickname: "test-server-001",
		Protocol:       "FC",
		OsType:         "Linux",
		IsReserved:     false,
		OsTypeOptions:  []int{1, 2, 3},
	}

	// Empty add host groups parameters
	addHgParams := gwymodel.AddHostGroupsToServerParam{}

	serverID, err := psm.ReconcileCreateAdminServer(params, addHgParams)
	if err != nil {
		t.Errorf("Unexpected error in ReconcileCreateAdminServer: %v", err)
		return
	}

	if serverID == 0 {
		t.Error("Expected server ID, got 0")
		return
	}

	// Get server info to verify creation
	serverInfo, err := psm.ReconcileReadAdminServer(serverID)
	if err != nil {
		t.Errorf("Failed to read created server: %v", err)
		return
	}

	if serverInfo.Nickname != params.ServerNickname {
		t.Errorf("Expected nickname %s, got %s", params.ServerNickname, serverInfo.Nickname)
	}

	t.Logf("Created server with ID %d: %+v", serverID, serverInfo)
}

// go test -v -run TestReconcileCreateAdminServerReserved
func xTestReconcileCreateAdminServerReserved(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test Case 2: Create a reserved server
	params := gwymodel.CreateAdminServerParams{
		ServerNickname: "reserved-server-001",
		IsReserved:     true,
		// Protocol and OsType should not be required for reserved servers
	}

	// Empty add host groups parameters
	addHgParams := gwymodel.AddHostGroupsToServerParam{}

	serverID, err := psm.ReconcileCreateAdminServer(params, addHgParams)
	if err != nil {
		t.Errorf("Unexpected error in ReconcileCreateAdminServer for reserved server: %v", err)
		return
	}

	if serverID == 0 {
		t.Error("Expected server ID, got 0")
		return
	}

	// Get server info to verify creation
	serverInfo, err := psm.ReconcileReadAdminServer(serverID)
	if err != nil {
		t.Errorf("Failed to read created reserved server: %v", err)
		return
	}

	if !serverInfo.IsReserved {
		t.Error("Expected reserved server, got non-reserved")
	}

	t.Logf("Created reserved server with ID %d: %+v", serverID, serverInfo)
}

// go test -v -run TestReconcileReadAdminServer
func xTestReconcileReadAdminServer(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test Case 1: Read existing server
	serverID := 1
	serverInfo, err := psm.ReconcileReadAdminServer(serverID)
	if err != nil {
		t.Errorf("Unexpected error in ReconcileReadAdminServer: %v", err)
		return
	}

	if serverInfo == nil {
		t.Error("Expected server info, got nil")
		return
	}

	if serverInfo.ID != serverID {
		t.Errorf("Expected server ID %d, got %d", serverID, serverInfo.ID)
	}

	t.Logf("Read server: %+v", serverInfo)
}

// go test -v -run TestReconcileReadAdminServerNotFound
func xTestReconcileReadAdminServerNotFound(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test Case 2: Read non-existent server
	nonExistentID := 99999
	serverInfo, err := psm.ReconcileReadAdminServer(nonExistentID)

	// Should handle gracefully - either return error or nil
	if err != nil {
		t.Logf("Expected error for non-existent server: %v", err)
	} else if serverInfo != nil {
		t.Errorf("Expected nil for non-existent server, got: %+v", serverInfo)
	}
}

// go test -v -run TestReconcileUpdateAdminServer
func xTestReconcileUpdateAdminServer(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test Case 1: Update server nickname
	serverID := 1
	params := gwymodel.UpdateAdminServerParams{
		Nickname:      "updated-test-server",
		OsType:        "VMware",
		OsTypeOptions: []int{4, 5, 6},
	}

	// Empty add host groups parameters
	addHgParams := gwymodel.AddHostGroupsToServerParam{}

	updatedServerID, err := psm.ReconcileUpdateAdminServer(serverID, params, addHgParams)
	if err != nil {
		t.Errorf("Unexpected error in ReconcileUpdateAdminServer: %v", err)
		return
	}

	if updatedServerID != serverID {
		t.Errorf("Expected server ID %d, got %d", serverID, updatedServerID)
		return
	}

	// Get server info to verify update
	serverInfo, err := psm.ReconcileReadAdminServer(serverID)
	if err != nil {
		t.Errorf("Failed to read updated server: %v", err)
		return
	}

	if serverInfo.Nickname != params.Nickname {
		t.Errorf("Expected updated nickname %s, got %s", params.Nickname, serverInfo.Nickname)
	}

	t.Logf("Updated server: %+v", serverInfo)
}

// go test -v -run TestReconcileDeleteAdminServer
func xTestReconcileDeleteAdminServer(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test Case 1: Delete server with LUN config retention
	serverID := 1
	params := gwymodel.DeleteAdminServerParams{
		KeepLunConfig: true,
	}

	err = psm.ReconcileDeleteAdminServer(serverID, params)
	if err != nil {
		t.Errorf("Unexpected error in ReconcileDeleteAdminServer: %v", err)
		return
	}

	t.Logf("Successfully deleted server %d with LUN config retention", serverID)
}

// go test -v -run TestReconcileDeleteAdminServerComplete
func xTestReconcileDeleteAdminServerComplete(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test Case 2: Delete server completely (including LUN config)
	serverID := 2
	params := gwymodel.DeleteAdminServerParams{
		KeepLunConfig: false,
	}

	err = psm.ReconcileDeleteAdminServer(serverID, params)
	if err != nil {
		t.Errorf("Unexpected error in ReconcileDeleteAdminServer (complete): %v", err)
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

	// Test Case 1: Get all servers
	params := gwymodel.AdminServerListParams{}

	response, err := psm.GetAdminServerList(params)
	if err != nil {
		t.Errorf("Unexpected error in GetAdminServerList: %v", err)
		return
	}

	if response == nil {
		t.Error("Expected server list response, got nil")
		return
	}

	t.Logf("Server list count: %d", len(response.Data))
	for i, server := range response.Data {
		t.Logf("Server %d: ID=%d, Nickname=%s, Protocol=%s", i, server.ID, server.Nickname, server.Protocol)
	}
}

// go test -v -run TestGetAdminServerListFiltered
func xTestGetAdminServerListFiltered(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test Case 2: Get servers filtered by nickname
	nickname := "test-server"
	params := gwymodel.AdminServerListParams{
		Nickname: &nickname,
	}

	response, err := psm.GetAdminServerList(params)
	if err != nil {
		t.Errorf("Unexpected error in GetAdminServerList (filtered): %v", err)
		return
	}

	if response == nil {
		t.Error("Expected filtered server list response, got nil")
		return
	}

	t.Logf("Filtered server list count: %d", len(response.Data))
	for i, server := range response.Data {
		t.Logf("Filtered Server %d: ID=%d, Nickname=%s", i, server.ID, server.Nickname)
	}
}

// go test -v -run TestGetAdminServerInfo
func xTestGetAdminServerInfo(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test Case 1: Get specific server info
	serverID := 1
	serverInfo, err := psm.GetAdminServerInfo(serverID)
	if err != nil {
		t.Errorf("Unexpected error in GetAdminServerInfo: %v", err)
		return
	}

	if serverInfo == nil {
		t.Error("Expected server info, got nil")
		return
	}

	if serverInfo.ID != serverID {
		t.Errorf("Expected server ID %d, got %d", serverID, serverInfo.ID)
	}

	t.Logf("Server info: ID=%d, Nickname=%s, Protocol=%s, OS=%s",
		serverInfo.ID, serverInfo.Nickname, serverInfo.Protocol, serverInfo.OsType)
	t.Logf("Capacity: Total=%d, Used=%d", serverInfo.TotalCapacity, serverInfo.UsedCapacity)
	t.Logf("Paths: %d", serverInfo.NumberOfPaths)

	for i, path := range serverInfo.Paths {
		t.Logf("Path %d: HBA WWN=%s, iSCSI Name=%s, Port IDs=%v",
			i, path.HbaWwn, path.IscsiName, path.PortIds)
	}
}

// go test -v -run TestAdminServerCRUDWorkflow
func xTestAdminServerCRUDWorkflow(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Complete CRUD workflow test
	t.Log("=== Starting Admin Server CRUD Workflow Test ===")

	// Step 1: Create a server
	t.Log("Step 1: Creating server...")
	createParams := gwymodel.CreateAdminServerParams{
		ServerNickname: "workflow-test-server",
		Protocol:       "iSCSI",
		OsType:         "Linux",
		IsReserved:     false,
		OsTypeOptions:  []int{1, 2},
	}

	// Empty add host groups parameters
	addHgParams := gwymodel.AddHostGroupsToServerParam{}

	createdID, err := psm.ReconcileCreateAdminServer(createParams, addHgParams)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	t.Logf("✓ Created server with ID: %d", createdID)

	// Step 2: Read the created server
	t.Log("Step 2: Reading server...")
	readInfo, err := psm.ReconcileReadAdminServer(createdID)
	if err != nil {
		t.Errorf("Failed to read created server: %v", err)
	} else {
		t.Logf("✓ Read server: %s (ID: %d)", readInfo.Nickname, readInfo.ID)
	}

	// Step 3: Update the server
	t.Log("Step 3: Updating server...")
	updateParams := gwymodel.UpdateAdminServerParams{
		Nickname:      "workflow-test-server-updated",
		OsType:        "VMware",
		OsTypeOptions: []int{3, 4, 5},
	}

	// Empty add host groups parameters for update
	updateAddHgParams := gwymodel.AddHostGroupsToServerParam{}

	updatedServerID, err := psm.ReconcileUpdateAdminServer(createdID, updateParams, updateAddHgParams)
	if err != nil {
		t.Errorf("Failed to update server: %v", err)
	} else {
		// Get server info to verify update
		updatedInfo, err := psm.ReconcileReadAdminServer(updatedServerID)
		if err != nil {
			t.Errorf("Failed to read updated server: %v", err)
		} else {
			t.Logf("✓ Updated server nickname to: %s", updatedInfo.Nickname)
		}
	}

	// Step 4: Delete the server
	t.Log("Step 4: Deleting server...")
	deleteParams := gwymodel.DeleteAdminServerParams{
		KeepLunConfig: false,
	}

	err = psm.ReconcileDeleteAdminServer(createdID, deleteParams)
	if err != nil {
		t.Errorf("Failed to delete server: %v", err)
	} else {
		t.Logf("✓ Deleted server with ID: %d", createdID)
	}

	t.Log("=== Admin Server CRUD Workflow Test Completed ===")
}
