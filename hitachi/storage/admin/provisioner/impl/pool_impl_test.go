package admin

import (
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	"testing"
)

// go test -v -run TestGetAdminPoolList
func xTestGetAdminPoolList(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Scenario 1: Get all pools without filters
	params := gwymodel.AdminPoolListParams{}
	resp, err := psm.GetAdminPoolList(params)
	if err != nil {
		t.Errorf("Unexpected error in GetAdminPoolList (no filter) %v", err)
		return
	}
	t.Logf("Response (all pools): Count=%d, Pools=%+v", resp.Count, resp.Data)

	// Scenario 2: Get pools by name filter
	poolName := "pool1"
	params = gwymodel.AdminPoolListParams{
		Name: &poolName,
	}
	resp, err = psm.GetAdminPoolList(params)
	if err != nil {
		t.Errorf("Unexpected error in GetAdminPoolList (name filter) %v", err)
		return
	}
	t.Logf("Response (filtered by name '%s'): Count=%d, Pools=%+v", poolName, resp.Count, resp.Data)

	// Scenario 3: Get pools by status filter
	status := "Normal"
	params = gwymodel.AdminPoolListParams{
		Status: &status,
	}
	resp, err = psm.GetAdminPoolList(params)
	if err != nil {
		t.Errorf("Unexpected error in GetAdminPoolList (status filter) %v", err)
		return
	}
	t.Logf("Response (filtered by status '%s'): Count=%d, Pools=%+v", status, resp.Count, resp.Data)
}

// go test -v -run TestGetAdminPoolInfo
func xTestGetAdminPoolInfo(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// First, get all pools to find a valid pool ID for testing
	params := gwymodel.AdminPoolListParams{}
	poolsList, err := psm.GetAdminPoolList(params)
	if err != nil {
		t.Fatalf("Unexpected error getting pools list for setup: %v", err)
	}

	if len(poolsList.Data) == 0 {
		t.Skip("No pools available for GetAdminPoolInfo test")
		return
	}

	// Use the first pool ID for testing
	testPoolID := poolsList.Data[0].ID

	// Scenario 1: Get specific pool by valid ID
	resp, err := psm.GetAdminPoolInfo(testPoolID)
	if err != nil {
		t.Errorf("Unexpected error in GetAdminPoolInfo (valid ID) %v", err)
		return
	}
	t.Logf("Response (pool %d): %+v", testPoolID, resp)

	// Validate the returned pool ID matches what we requested
	if resp.ID != testPoolID {
		t.Errorf("Expected pool ID %d, got %d", testPoolID, resp.ID)
	}

	// Scenario 2: Test with a non-existent pool ID (this should fail)
	invalidPoolID := 99999
	_, err = psm.GetAdminPoolInfo(invalidPoolID)
	if err == nil {
		t.Logf("Note: GetAdminPoolInfo with invalid ID '%d' did not return an error (this may be expected behavior)", invalidPoolID)
	} else {
		t.Logf("Expected error for invalid pool ID %d: %v", invalidPoolID, err)
	}
}

// go test -v -run TestCreateAdminPool
func xTestCreateAdminPool(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Scenario 1: Create a pool with SSD drives
	params := gwymodel.CreateAdminPoolParams{
		Name:                "test-pool-provisioner-ssd",
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

	err = psm.CreateAdminPool(params)
	if err != nil {
		t.Errorf("Unexpected error in CreateAdminPool (SSD): %v", err)
		return
	}
	t.Logf("Successfully created pool: %s", params.Name)

	// Scenario 2: Create a pool with multiple drive types
	params = gwymodel.CreateAdminPoolParams{
		Name:                "test-pool-provisioner-mixed",
		IsEncryptionEnabled: true,
		Drives: []gwymodel.CreateAdminPoolDrive{
			{
				DriveTypeCode:   "SSD",
				DataDriveCount:  2,
				RaidLevel:       "RAID1",
				ParityGroupType: "1D+1D",
			},
			{
				DriveTypeCode:   "SAS",
				DataDriveCount:  6,
				RaidLevel:       "RAID6",
				ParityGroupType: "6D+2P",
			},
		},
	}

	err = psm.CreateAdminPool(params)
	if err != nil {
		t.Errorf("Unexpected error in CreateAdminPool (mixed drives): %v", err)
		return
	}
	t.Logf("Successfully created pool: %s", params.Name)
}

// go test -v -run TestUpdateAdminPool
func xTestUpdateAdminPool(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Get a pool to update
	params := gwymodel.AdminPoolListParams{}
	poolsList, err := psm.GetAdminPoolList(params)
	if err != nil {
		t.Fatalf("Unexpected error getting pools list for setup: %v", err)
	}

	if len(poolsList.Data) == 0 {
		t.Skip("No pools available for UpdateAdminPool test")
		return
	}

	testPoolID := poolsList.Data[0].ID
	originalName := poolsList.Data[0].Name

	// Scenario 1: Update pool name
	updateParams := gwymodel.UpdateAdminPoolParams{
		Name: "updated-pool-name-provisioner",
	}

	err = psm.UpdateAdminPool(testPoolID, updateParams)
	if err != nil {
		t.Errorf("Unexpected error in UpdateAdminPool (name): %v", err)
		return
	}
	t.Logf("Successfully updated pool %d name to: %s", testPoolID, updateParams.Name)

	// Scenario 2: Update thresholds
	updateParams = gwymodel.UpdateAdminPoolParams{
		Name:               originalName, // restore original name
		ThresholdWarning:   75,
		ThresholdDepletion: 88,
	}

	err = psm.UpdateAdminPool(testPoolID, updateParams)
	if err != nil {
		t.Errorf("Unexpected error in UpdateAdminPool (thresholds): %v", err)
		return
	}
	t.Logf("Successfully updated pool %d thresholds: warning=%d, depletion=%d",
		testPoolID, updateParams.ThresholdWarning, updateParams.ThresholdDepletion)
}

// go test -v -run TestExpandAdminPool
func xTestExpandAdminPool(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Get a pool to expand
	params := gwymodel.AdminPoolListParams{}
	poolsList, err := psm.GetAdminPoolList(params)
	if err != nil {
		t.Fatalf("Unexpected error getting pools list for setup: %v", err)
	}

	if len(poolsList.Data) == 0 {
		t.Skip("No pools available for ExpandAdminPool test")
		return
	}

	testPoolID := poolsList.Data[0].ID

	// Scenario 1: Expand pool with additional SSD drives
	expandParams := gwymodel.ExpandAdminPoolParams{
		AdditionalDrives: []gwymodel.ExpandAdminPoolDrive{
			{
				DriveTypeCode:   "SSD",
				DataDriveCount:  2,
				RaidLevel:       "RAID5",
				ParityGroupType: "3D+1P",
			},
		},
	}

	err = psm.ExpandAdminPool(testPoolID, expandParams)
	if err != nil {
		t.Errorf("Unexpected error in ExpandAdminPool (SSD): %v", err)
		return
	}
	t.Logf("Successfully expanded pool %d with %d SSD drives",
		testPoolID, expandParams.AdditionalDrives[0].DataDriveCount)

	// Scenario 2: Expand pool with multiple drive types
	expandParams = gwymodel.ExpandAdminPoolParams{
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

	err = psm.ExpandAdminPool(testPoolID, expandParams)
	if err != nil {
		t.Errorf("Unexpected error in ExpandAdminPool (mixed drives): %v", err)
		return
	}
	t.Logf("Successfully expanded pool %d with mixed drives: SSD=%d, SAS=%d",
		testPoolID, expandParams.AdditionalDrives[0].DataDriveCount, expandParams.AdditionalDrives[1].DataDriveCount)
}

// go test -v -run TestDeleteAdminPool
func xTestDeleteAdminPool(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// First create a pool to delete
	createParams := gwymodel.CreateAdminPoolParams{
		Name:                "test-pool-to-delete-provisioner",
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

	err = psm.CreateAdminPool(createParams)
	if err != nil {
		t.Fatalf("Failed to create test pool for deletion: %v", err)
	}

	// Get the created pool ID
	listParams := gwymodel.AdminPoolListParams{
		Name: &createParams.Name,
	}
	poolsList, err := psm.GetAdminPoolList(listParams)
	if err != nil || len(poolsList.Data) == 0 {
		t.Fatalf("Failed to find created test pool: %v", err)
	}

	testPoolID := poolsList.Data[0].ID

	// Now delete the pool
	err = psm.DeleteAdminPool(testPoolID)
	if err != nil {
		t.Errorf("Unexpected error in DeleteAdminPool: %v", err)
		return
	}
	t.Logf("Successfully deleted pool %d", testPoolID)

	// Verify pool is deleted
	_, err = psm.GetAdminPoolInfo(testPoolID)
	if err == nil {
		t.Errorf("Expected error when getting deleted pool %d, but got none", testPoolID)
	} else {
		t.Logf("Confirmed pool %d is deleted: %v", testPoolID, err)
	}
}
