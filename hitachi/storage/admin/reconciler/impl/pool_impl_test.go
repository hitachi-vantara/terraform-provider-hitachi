package admin

import (
	"encoding/json"
	"testing"

	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
)

// go test -v -run TestReconcileCreateAdminPool
func xTestReconcileCreateAdminPool(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	// Scenario 1: Create pool with SSD drives
	params := gwymodel.CreateAdminPoolParams{
		Name:                "test-pool-reconciler-ssd",
		IsEncryptionEnabled: false,
		Drives: []gwymodel.CreateAdminPoolDrive{
			{
				DriveTypeCode:   "SSD",
				DataDriveCount:  4,
				RaidLevel:       "RAID5",
				ParityGroupType: "3D+1P",
			},
		},
	}

	poolID, err := psm.ReconcileCreateAdminPool(params)
	if err != nil {
		t.Errorf("Unexpected error in ReconcileCreateAdminPool: %v", err)
		return
	}

	if poolID <= 0 {
		t.Errorf("Expected valid pool ID but got %d", poolID)
		return
	}

	t.Logf("Successfully created pool '%s' with ID: %d", params.Name, poolID)

	// Scenario 2: Create pool with encryption enabled
	params = gwymodel.CreateAdminPoolParams{
		Name:                "test-pool-reconciler-encrypted",
		IsEncryptionEnabled: true,
		Drives: []gwymodel.CreateAdminPoolDrive{
			{
				DriveTypeCode:   "SAS",
				DataDriveCount:  6,
				RaidLevel:       "RAID6",
				ParityGroupType: "6D+2P",
			},
		},
	}

	poolID, err = psm.ReconcileCreateAdminPool(params)
	if err != nil {
		t.Errorf("Unexpected error in ReconcileCreateAdminPool (encrypted): %v", err)
		return
	}

	t.Logf("Successfully created encrypted pool '%s' with ID: %d", params.Name, poolID)
}

// go test -v -run TestReconcileReadAdminPool
func xTestReconcileReadAdminPool(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	// Get a pool list to find a valid pool ID for testing
	listParams := gwymodel.AdminPoolListParams{}
	poolsList, err := psm.GetAdminPoolList(listParams)
	if err != nil {
		t.Fatalf("Unexpected error getting pools list for setup: %v", err)
	}

	if len(poolsList.Data) == 0 {
		t.Skip("No pools available for ReconcileReadAdminPool test")
		return
	}

	testPoolID := poolsList.Data[0].ID

	// Test reading the pool
	poolInfo, timedOut, err := psm.ReconcileReadAdminPool(testPoolID)
	if err != nil {
		t.Errorf("Unexpected error in ReconcileReadAdminPool: %v", err)
		return
	}

	if timedOut {
		t.Logf("Note: Pool read operation timed out for pool %d", testPoolID)
	}

	if poolInfo == nil {
		t.Errorf("Expected pool information but got nil")
		return
	}

	if poolInfo.ID != testPoolID {
		t.Errorf("Expected pool ID %d but got %d", testPoolID, poolInfo.ID)
		return
	}

	b, _ := json.MarshalIndent(poolInfo, "", "  ")
	t.Logf("Pool ID: %d\nPool Info:\n%s", testPoolID, string(b))

	// Test with invalid pool ID
	invalidPoolID := 99999
	poolInfo, timedOut, err = psm.ReconcileReadAdminPool(invalidPoolID)
	if err == nil {
		t.Logf("Note: ReconcileReadAdminPool with invalid ID '%d' did not return an error (this may be expected behavior)", invalidPoolID)
	} else {
		t.Logf("Expected error for invalid pool ID %d: %v", invalidPoolID, err)
	}
	if timedOut {
		t.Logf("Note: Invalid pool read timed out")
	}
}

// go test -v -run TestReconcileUpdateAdminPool
func xTestReconcileUpdateAdminPool(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	// Get a pool list to find a valid pool ID for testing
	listParams := gwymodel.AdminPoolListParams{}
	poolsList, err := psm.GetAdminPoolList(listParams)
	if err != nil {
		t.Fatalf("Unexpected error getting pools list for setup: %v", err)
	}

	if len(poolsList.Data) == 0 {
		t.Skip("No pools available for ReconcileUpdateAdminPool test")
		return
	}

	testPoolID := poolsList.Data[0].ID
	originalName := poolsList.Data[0].Name

	tests := []struct {
		name       string
		poolID     int
		updateFunc func() gwymodel.UpdateAdminPoolParams
	}{
		{
			name:   "Update Pool Name",
			poolID: testPoolID,
			updateFunc: func() gwymodel.UpdateAdminPoolParams {
				return gwymodel.UpdateAdminPoolParams{
					Name: "updated-pool-name-reconciler",
				}
			},
		},
		{
			name:   "Update Pool Thresholds",
			poolID: testPoolID,
			updateFunc: func() gwymodel.UpdateAdminPoolParams {
				return gwymodel.UpdateAdminPoolParams{
					Name:               originalName, // restore original name
					ThresholdWarning:   75,
					ThresholdDepletion: 88,
				}
			},
		},
		{
			name:   "Update Pool Name and Thresholds",
			poolID: testPoolID,
			updateFunc: func() gwymodel.UpdateAdminPoolParams {
				return gwymodel.UpdateAdminPoolParams{
					Name:               "comprehensive-update-reconciler",
					ThresholdWarning:   70,
					ThresholdDepletion: 85,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateParams := tt.updateFunc()

			err := psm.ReconcileUpdateAdminPool(tt.poolID, updateParams)
			if err != nil {
				t.Errorf("Unexpected error in ReconcileUpdateAdminPool (%s): %v", tt.name, err)
				return
			}

			t.Logf("Successfully updated pool %d: %s", tt.poolID, tt.name)

			// Verify the update by reading the pool
			updatedPool, _, err := psm.ReconcileReadAdminPool(tt.poolID)
			if err != nil {
				t.Errorf("Failed to read updated pool: %v", err)
				return
			}

			if updateParams.Name != "" && updatedPool.Name != updateParams.Name {
				t.Errorf("Expected updated name '%s', got '%s'", updateParams.Name, updatedPool.Name)
			}

			if updateParams.ThresholdWarning > 0 && updatedPool.CapacityManage.ThresholdWarning != updateParams.ThresholdWarning {
				t.Errorf("Expected threshold warning %d, got %d",
					updateParams.ThresholdWarning, updatedPool.CapacityManage.ThresholdWarning)
			}

			if updateParams.ThresholdDepletion > 0 && updatedPool.CapacityManage.ThresholdDepletion != updateParams.ThresholdDepletion {
				t.Errorf("Expected threshold depletion %d, got %d",
					updateParams.ThresholdDepletion, updatedPool.CapacityManage.ThresholdDepletion)
			}
		})
	}
}

// go test -v -run TestReconcileExpandAdminPool
func xTestReconcileExpandAdminPool(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	// Get a pool list to find a valid pool ID for testing
	listParams := gwymodel.AdminPoolListParams{}
	poolsList, err := psm.GetAdminPoolList(listParams)
	if err != nil {
		t.Fatalf("Unexpected error getting pools list for setup: %v", err)
	}

	if len(poolsList.Data) == 0 {
		t.Skip("No pools available for ReconcileExpandAdminPool test")
		return
	}

	testPoolID := poolsList.Data[0].ID

	// Get initial pool info for comparison
	initialPool, _, err := psm.ReconcileReadAdminPool(testPoolID)
	if err != nil {
		t.Fatalf("Failed to get initial pool info: %v", err)
	}

	tests := []struct {
		name        string
		poolID      int
		expandFunc  func() gwymodel.ExpandAdminPoolParams
		description string
	}{
		{
			name:   "Expand with SSD drives",
			poolID: testPoolID,
			expandFunc: func() gwymodel.ExpandAdminPoolParams {
				return gwymodel.ExpandAdminPoolParams{
					AdditionalDrives: []gwymodel.ExpandAdminPoolDrive{
						{
							DriveTypeCode:   "SSD",
							DataDriveCount:  2,
							RaidLevel:       "RAID5",
							ParityGroupType: "3D+1P",
						},
					},
				}
			},
			description: "Adding 2 SSD drives with RAID5",
		},
		{
			name:   "Expand with mixed drives",
			poolID: testPoolID,
			expandFunc: func() gwymodel.ExpandAdminPoolParams {
				return gwymodel.ExpandAdminPoolParams{
					AdditionalDrives: []gwymodel.ExpandAdminPoolDrive{
						{
							DriveTypeCode:   "SSD",
							DataDriveCount:  2,
							RaidLevel:       "RAID5",
							ParityGroupType: "3D+1P",
						},
						{
							DriveTypeCode:   "SAS",
							DataDriveCount:  4,
							RaidLevel:       "RAID6",
							ParityGroupType: "6D+2P",
						},
					},
				}
			},
			description: "Adding mixed drives: 2 SSD + 4 SAS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expandParams := tt.expandFunc()

			t.Logf("Expanding pool %d: %s", tt.poolID, tt.description)

			err := psm.ReconcileExpandAdminPool(tt.poolID, expandParams)
			if err != nil {
				t.Errorf("Unexpected error in ReconcileExpandAdminPool (%s): %v", tt.name, err)
				return
			}

			t.Logf("Successfully expanded pool %d: %s", tt.poolID, tt.name)

			// Verify the expansion by reading the pool
			expandedPool, _, err := psm.ReconcileReadAdminPool(tt.poolID)
			if err != nil {
				t.Errorf("Failed to read expanded pool: %v", err)
				return
			}

			// Check that capacity increased (this is a basic check)
			if expandedPool.TotalCapacity <= initialPool.TotalCapacity {
				t.Logf("Note: Total capacity did not increase after expansion. Initial: %d, Current: %d (this may be expected if the expansion is still in progress)",
					initialPool.TotalCapacity, expandedPool.TotalCapacity)
			} else {
				t.Logf("Pool capacity increased from %d to %d",
					initialPool.TotalCapacity, expandedPool.TotalCapacity)
			}

			// Log expanded pool details
			b, _ := json.MarshalIndent(expandedPool, "", "  ")
			t.Logf("Expanded Pool Info:\n%s", string(b))
		})
	}
}

// go test -v -run TestReconcileDeleteAdminPool
func xTestReconcileDeleteAdminPool(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	// First create a pool to delete
	createParams := gwymodel.CreateAdminPoolParams{
		Name:                "test-pool-to-delete-reconciler",
		IsEncryptionEnabled: false,
		Drives: []gwymodel.CreateAdminPoolDrive{
			{
				DriveTypeCode:   "SSD",
				DataDriveCount:  2,
				RaidLevel:       "RAID1",
				ParityGroupType: "1D+1D",
			},
		},
	}

	poolID, err := psm.ReconcileCreateAdminPool(createParams)
	if err != nil {
		t.Fatalf("Failed to create test pool for deletion: %v", err)
	}

	t.Logf("Created test pool '%s' with ID: %d for deletion test", createParams.Name, poolID)

	// Verify the pool exists
	poolInfo, _, err := psm.ReconcileReadAdminPool(poolID)
	if err != nil {
		t.Fatalf("Failed to read created test pool: %v", err)
	}

	if poolInfo.ID != poolID {
		t.Fatalf("Created pool ID mismatch: expected %d, got %d", poolID, poolInfo.ID)
	}

	// Now delete the pool
	err = psm.ReconcileDeleteAdminPool(poolID)
	if err != nil {
		t.Errorf("Unexpected error in ReconcileDeleteAdminPool: %v", err)
		return
	}

	t.Logf("Successfully deleted pool %d", poolID)

	// Verify pool is deleted by trying to read it
	_, _, err = psm.ReconcileReadAdminPool(poolID)
	if err == nil {
		t.Errorf("Expected error when reading deleted pool %d, but got none", poolID)
	} else {
		t.Logf("Confirmed pool %d is deleted: %v", poolID, err)
	}
}

// go test -v -run TestGetAdminPoolList
func xTestGetAdminPoolList(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	// Test basic pool list functionality
	params := gwymodel.AdminPoolListParams{}
	resp, err := psm.GetAdminPoolList(params)
	if err != nil {
		t.Errorf("Unexpected error in GetAdminPoolList: %v", err)
		return
	}

	t.Logf("Total pools found: %d", resp.Count)

	if len(resp.Data) > 0 {
		// Log details of first few pools
		maxPools := 3
		if len(resp.Data) < maxPools {
			maxPools = len(resp.Data)
		}

		for i := 0; i < maxPools; i++ {
			pool := resp.Data[i]
			t.Logf("Pool %d: ID=%d, Name=%s, Status=%s, TotalCapacity=%d",
				i+1, pool.ID, pool.Name, pool.Status, pool.TotalCapacity)
		}
	}
}

// go test -v -run TestGetAdminPoolInfo
func xTestGetAdminPoolInfo(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	// Get a pool list to find a valid pool ID for testing
	listParams := gwymodel.AdminPoolListParams{}
	poolsList, err := psm.GetAdminPoolList(listParams)
	if err != nil {
		t.Fatalf("Unexpected error getting pools list for setup: %v", err)
	}

	if len(poolsList.Data) == 0 {
		t.Skip("No pools available for GetAdminPoolInfo test")
		return
	}

	testPoolID := poolsList.Data[0].ID

	// Test getting pool info
	poolInfo, err := psm.GetAdminPoolInfo(testPoolID)
	if err != nil {
		t.Errorf("Unexpected error in GetAdminPoolInfo: %v", err)
		return
	}

	if poolInfo == nil {
		t.Errorf("Expected pool information but got nil")
		return
	}

	if poolInfo.ID != testPoolID {
		t.Errorf("Expected pool ID %d but got %d", testPoolID, poolInfo.ID)
		return
	}

	b, _ := json.MarshalIndent(poolInfo, "", "  ")
	t.Logf("Pool ID: %d\nPool Info:\n%s", testPoolID, string(b))
}
