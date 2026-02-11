package sanstorage

import (
	"encoding/json"
	// "fmt"
	model "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	"testing"
)

// go test -v -run ^TestGetSnapshotGroups$
func xTestGetSnapshotGroups(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating sanStorageManager: %v", err)
	}

	params := model.GetSnapshotGroupsParams{}

	resp, err := psm.GetSnapshotGroups(params)
	if err != nil {
		t.Errorf("Unexpected error calling GetSnapshotGroups: %v", err)
		return
	}

	if resp == nil {
		t.Errorf("Expected response, got nil")
		return
	}

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Logf("Failed to marshal SnapshotGroupListResponse: %v", err)
		return
	}

	t.Logf("Response:\n%s", string(data))
}

// go test -v -run ^TestGetSnapshotGroupsByName$
func xTestGetSnapshotGroupsByName(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating sanStorageManager: %v", err)
	}

	snapgroupName := "snewar-tia-grp-02"
	params := model.GetSnapshotGroupsParams{
		SnapshotGroupName: &snapgroupName,
	}

	resp, err := psm.GetSnapshotGroups(params)
	if err != nil {
		t.Errorf("Unexpected error calling GetSnapshotGroups: %v", err)
		return
	}

	if resp == nil {
		t.Errorf("Expected response, got nil")
		return
	}

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Logf("Failed to marshal SnapshotGroupListResponse: %v", err)
		return
	}

	t.Logf("Response:\n%s", string(data))
}

// go test -v -run ^TestGetSnapshotGroupsDetailedPair$
func xTestGetSnapshotGroupsDetailedPair(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating sanStorageManager: %v", err)
	}

	detailType := "pair"
	params := model.GetSnapshotGroupsParams{
		DetailInfoType: &detailType,
	}

	resp, err := psm.GetSnapshotGroups(params)
	if err != nil {
		t.Errorf("Unexpected error calling GetSnapshotGroups: %v", err)
		return
	}

	if resp == nil {
		t.Errorf("Expected response, got nil")
		return
	}

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Logf("Failed to marshal SnapshotGroupListResponse: %v", err)
		return
	}

	t.Logf("Response:\n%s", string(data))
}

// go test -v -run ^TestGetSnapshotGroupsRetention$
func xTestGetSnapshotGroupsRetention(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating sanStorageManager: %v", err)
	}

	snapgroupName := "snewar-tia-grp-02"
	detailType := "retention"
	params := model.GetSnapshotGroupsParams{
		SnapshotGroupName: &snapgroupName,
		DetailInfoType: &detailType,
	}

	resp, err := psm.GetSnapshotGroups(params)
	if err != nil {
		t.Errorf("Unexpected error calling GetSnapshotGroups: %v", err)
		return
	}

	if resp == nil {
		t.Errorf("Expected response, got nil")
		return
	}

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Logf("Failed to marshal SnapshotGroupListResponse: %v", err)
		return
	}

	t.Logf("Response:\n%s", string(data))
}

// go test -v -run ^TestGetSnapshotGroup$
func xTestGetSnapshotGroup(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating sanStorageManager: %v", err)
	}

	// snapshotGroupID := "halbaSnapshotGroup1"
	snapshotGroupID := "snewar-tia-grp-02"

	params := model.GetSnapshotGroupsParams{}

	resp, err := psm.GetSnapshotGroup(snapshotGroupID, params)
	if err != nil {
		t.Errorf("Unexpected error calling GetSnapshotGroup: %v", err)
		return
	}

	if resp == nil {
		t.Errorf("Expected response, got nil")
		return
	}

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Logf("Failed to marshal SnapshotGroup: %v", err)
		return
	}

	t.Logf("Response:\n%s", string(data))
}

// go test -v -run ^TestGetSnapshotGroupRetention$
func xTestGetSnapshotGroupRetention(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating sanStorageManager: %v", err)
	}

	// snapshotGroupID := "halbaSnapshotGroup1"
	snapshotGroupID := "snewar-tia-grp-02"

	detailType := "retention"
	params := model.GetSnapshotGroupsParams{
		DetailInfoType: &detailType,
	}

	resp, err := psm.GetSnapshotGroup(snapshotGroupID, params)
	if err != nil {
		t.Errorf("Unexpected error calling GetSnapshotGroup: %v", err)
		return
	}

	if resp == nil {
		t.Errorf("Expected response, got nil")
		return
	}

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Logf("Failed to marshal SnapshotGroup: %v", err)
		return
	}

	t.Logf("Response:\n%s", string(data))
}

// go test -v -run ^TestSplitSnapshotGroup$
func xTestSplitSnapshotGroup(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating sanStorageManager: %v", err)
	}

	snapshotGroupID := "snapshotGroup"
	// retentionPeriod := 1 // hours

	request := model.SplitSnapshotRequest{
		Parameters: model.SplitSnapshotParams{
			// RetentionPeriod: &retentionPeriod,
		},
	}

	resIds, err := psm.SplitSnapshotGroup(snapshotGroupID, request)
	if err != nil {
		t.Errorf("Unexpected error calling SplitSnapshotGroup: %v", err)
		return
	}

	t.Logf("Successfully initiated SplitSnapshotGroup. %v", resIds)
}

// go test -v -run ^TestResyncSnapshotGroup$
func xTestResyncSnapshotGroup(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating sanStorageManager: %v", err)
	}

	snapshotGroupID := "snapshotGroup"
	autoSplit := false

	request := model.ResyncSnapshotRequest{
		Parameters: model.ResyncSnapshotParams{
			AutoSplit: &autoSplit,
			// RetentionPeriod: &retentionPeriod,
		},
	}

	resIds, err := psm.ResyncSnapshotGroup(snapshotGroupID, request)
	if err != nil {
		t.Errorf("Unexpected error calling ResyncSnapshotGroup: %v", err)
		return
	}

	t.Logf("Successfully initiated ResyncSnapshotGroup. %v", resIds)
}

// go test -v -run ^TestRestoreSnapshotGroup$
func xTestRestoreSnapshotGroup(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating sanStorageManager: %v", err)
	}

	snapshotGroupID := "snapshotGroup"
	autoSplit := false

	request := model.RestoreSnapshotRequest{
		Parameters: model.RestoreSnapshotParams{
			AutoSplit: &autoSplit,
		},
	}

	resIds, err := psm.RestoreSnapshotGroup(snapshotGroupID, request)
	if err != nil {
		t.Errorf("Unexpected error calling RestoreSnapshotGroup: %v", err)
		return
	}

	t.Logf("Successfully initiated RestoreSnapshotGroup. %v", resIds)
}

// go test -v -run ^TestDeleteSnapshotGroup$
func xTestDeleteSnapshotGroup(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating sanStorageManager: %v", err)
	}

	snapshotGroupID := "snapshotGroup"

	resIds, err := psm.DeleteSnapshotGroup(snapshotGroupID)
	if err != nil {
		t.Errorf("Unexpected error calling DeleteSnapshotGroup: %v", err)
		return
	}

	t.Logf("Successfully initiated DeleteSnapshotGroup. %v", resIds)
}

// go test -v -run ^TestCloneSnapshotGroup$
func xTestCloneSnapshotGroup(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating sanStorageManager: %v", err)
	}

	snapshotGroupID := "snapshotGroup"
	copySpeed := "faster"

	request := model.CloneSnapshotRequest{
		Parameters: model.CloneSnapshotParams{
			CopySpeed: &copySpeed,
		},
	}

	resIds, err := psm.CloneSnapshotGroup(snapshotGroupID, request)
	if err != nil {
		t.Errorf("Unexpected error calling CloneSnapshotGroup: %v", err)
		return
	}

	t.Logf("Successfully initiated CloneSnapshotGroup. %v", resIds)
}

// go test -v -run ^TestDeleteSnapshotTree$
func xTestDeleteSnapshotTree(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating sanStorageManager: %v", err)
	}

	request := model.DeleteSnapshotTreeRequest{
		Parameters: model.DeleteSnapshotTreeParams{
			LdevID: 1,
		},
	}

	resIds, err := psm.DeleteSnapshotTree(request)
	if err != nil {
		t.Errorf("Unexpected error calling DeleteSnapshotTree: %v", err)
		return
	}

	t.Logf("Successfully initiated DeleteSnapshotTree. %v", resIds)
}

// go test -v -run ^TestDeleteGarbageData$
func xTestDeleteGarbageData(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating sanStorageManager: %v", err)
	}

	request := model.DeleteGarbageDataRequest{
		Parameters: model.DeleteGarbageDataParams{
			LdevID:        66,
			OperationType: "start",
		},
	}

	resIds, err := psm.DeleteGarbageData(request)
	if err != nil {
		t.Errorf("Unexpected error calling DeleteGarbageData: %v", err)
		return
	}

	t.Logf("Successfully initiated DeleteGarbageData. %v", resIds)
}
