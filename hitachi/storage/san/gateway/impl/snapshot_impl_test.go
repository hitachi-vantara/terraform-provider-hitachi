package sanstorage

import (
	"encoding/json"
	// "fmt"
	model "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	"testing"
)

// go test -v -run TestGetSnapshots
func xTestGetSnapshots(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating sanStorageManager: %v", err)
	}

	// The following shows how to specify parameters:
	// The LDEV number of the P-VOL and the snapshot group name
	// The LDEV number and the MU number of the P-VOL
	// Only the LDEV number of the P-VOL
	// If the secondary volume exists, only the LDEV number of the S-VOL
	// If no parameters are specified, an error occurs.

	// snapshotGroupName := "example_group"
	// pvolLdevID := 1251
	pvolLdevID := 359
	// muNumber := 0
	detailInfoType := "retention"

	params := model.GetSnapshotsParams{
		// SnapshotGroupName: &snapshotGroupName,
		PvolLdevID: &pvolLdevID,
		// SvolLdevID:        nil,
		// MuNumber:          &muNumber,
		DetailInfoType: &detailInfoType,
	}

	resp, err := psm.GetSnapshots(params)
	if err != nil {
		t.Errorf("Unexpected error calling GetSnapshots: %v", err)
		return
	}

	if resp == nil {
		t.Errorf("Expected response, got nil")
		return
	}

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Logf("Failed to marshal SnapshotListResponse: %v", err)
		return
	}

	t.Logf("Response:\n%s", string(data))
}

// go test -v -run ^TestGetSnapshot$
func xTestGetSnapshot(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating sanStorageManager: %v", err)
	}

	pvolLdevID := 307
	muNumber := 3

	resp, err := psm.GetSnapshot(pvolLdevID, muNumber)
	if err != nil {
		t.Errorf("Unexpected error calling GetSnapshot: %v", err)
		return
	}

	if resp == nil {
		t.Errorf("Expected response, got nil")
		return
	}

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Logf("Failed to marshal Snapshot: %v", err)
		return
	}

	t.Logf("Response:\n%s", string(data))
}

// go test -v -run TestGetSnapshotReplicationsRange
// only for VSP 5000 series
func xTestGetSnapshotReplicationsRange(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating sanStorageManager: %v", err)
	}

	startID := 0
	// endID := 100

	params := model.GetSnapshotReplicationsRangeParams{
		StartPvolLdevID: &startID,
		// EndPvolLdevID:   &endID,
	}

	resp, err := psm.GetSnapshotReplicationsRange(params)
	if err != nil {
		t.Errorf("Unexpected error calling GetSnapshotReplicationsRange: %v", err)
		return
	}

	if resp == nil {
		t.Errorf("Expected response, got nil")
		return
	}

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Logf("Failed to marshal Response: %v", err)
		return
	}

	t.Logf("Response:\n%s", string(data))
}

// go test -v -run TestCreateSnapshot
func xTestCreateSnapshot(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating sanStorageManager: %v", err)
	}

	snapshotGroupName := "halbaSnapshotGroup1"
	autoSplit := false
	isConsistencyGroup := false

	params := model.CreateSnapshotParams{
		SnapshotGroupName:  snapshotGroupName,
		SnapshotPoolID:     103,
		PvolLdevID:         307,
		AutoSplit:          &autoSplit,
		IsConsistencyGroup: &isConsistencyGroup,
	}

	resIds, err := psm.CreateSnapshot(params)
	if err != nil {
		t.Errorf("Unexpected error calling CreateSnapshot: %v", err)
		return
	}

	t.Logf("Successfully initiated CreateSnapshot. %v", resIds)
}

// go test -v -run TestSplitSnapshot
func xTestSplitSnapshot(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating sanStorageManager: %v", err)
	}

	pvolLdevID := 307
	muNumber := 3
	retentionPeriod := 120 // hours

	request := model.SplitSnapshotRequest{
		Parameters: model.SplitSnapshotParams{
			RetentionPeriod: &retentionPeriod,
		},
	}

	resIds, err := psm.SplitSnapshot(pvolLdevID, muNumber, request)
	if err != nil {
		t.Errorf("Unexpected error calling SplitSnapshot: %v", err)
		return
	}

	t.Logf("Successfully initiated SplitSnapshot. %v", resIds)
}

// go test -v -run TestResyncSnapshot
func xTestResyncSnapshot(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating sanStorageManager: %v", err)
	}

	pvolLdevID := 307
	muNumber := 3

	autoSplit := false
	request := model.ResyncSnapshotRequest{
		Parameters: model.ResyncSnapshotParams{
			AutoSplit: &autoSplit,
		},
	}

	resIds, err := psm.ResyncSnapshot(pvolLdevID, muNumber, request)
	if err != nil {
		t.Errorf("Unexpected error calling ResyncSnapshot: %v", err)
		return
	}

	t.Logf("Successfully initiated ResyncSnapshot. %v", resIds)
}

// go test -v -run TestRestoreSnapshot
func xTestRestoreSnapshot(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating sanStorageManager: %v", err)
	}

	pvolLdevID := 307
	muNumber := 3

	autoSplit := false
	request := model.RestoreSnapshotRequest{
		Parameters: model.RestoreSnapshotParams{
			AutoSplit: &autoSplit,
		},
	}

	resIds, err := psm.RestoreSnapshot(pvolLdevID, muNumber, request)
	if err != nil {
		t.Errorf("Unexpected error calling RestoreSnapshot: %v", err)
		return
	}

	t.Logf("Successfully initiated RestoreSnapshot. %v", resIds)
}

// go test -v -run TestDeleteSnapshot
func xTestDeleteSnapshot(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating sanStorageManager: %v", err)
	}

	pvolLdevID := 307
	muNumber := 3

	resIds, err := psm.DeleteSnapshot(pvolLdevID, muNumber)
	if err != nil {
		t.Errorf("Unexpected error calling DeleteSnapshot: %v", err)
		return
	}

	t.Logf("Successfully initiated DeleteSnapshot. %v", resIds)
}

// go test -v -run TestCloneSnapshot
func xTestCloneSnapshot(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating sanStorageManager: %v", err)
	}

	pvolLdevID := 307
	muNumber := 3

	copySpeed := "faster"
	request := model.CloneSnapshotRequest{
		Parameters: model.CloneSnapshotParams{
			CopySpeed: &copySpeed,
		},
	}

	resIds, err := psm.CloneSnapshot(pvolLdevID, muNumber, request)
	if err != nil {
		t.Errorf("Unexpected error calling CloneSnapshot: %v", err)
		return
	}

	t.Logf("Successfully initiated CloneSnapshot. %v", resIds)
}

// go test -v -run TestAssignSnapshotVolume
func xTestAssignSnapshotVolume(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating sanStorageManager: %v", err)
	}

	pvolLdevID := 307
	muNumber := 3
	svolLdevID := 101

	request := model.AssignSnapshotVolumeRequest{
		Parameters: model.AssignSnapshotVolumeParams{
			SvolLdevID: svolLdevID,
		},
	}

	resIds, err := psm.AssignSnapshotVolume(pvolLdevID, muNumber, request)
	if err != nil {
		t.Errorf("Unexpected error calling AssignSnapshotVolume: %v", err)
		return
	}

	t.Logf("Successfully initiated AssignSnapshotVolume. %v", resIds)
}

// go test -v -run TestUnassignSnapshotVolume
func xTestUnassignSnapshotVolume(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating sanStorageManager: %v", err)
	}

	pvolLdevID := 307
	muNumber := 3

	resIds, err := psm.UnassignSnapshotVolume(pvolLdevID, muNumber)
	if err != nil {
		t.Errorf("Unexpected error calling UnassignSnapshotVolume: %v", err)
		return
	}

	t.Logf("Successfully initiated UnassignSnapshotVolume. %v", resIds)
}

// go test -v -run TestSetSnapshotRetentionPeriod
func xTestSetSnapshotRetentionPeriod(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating sanStorageManager: %v", err)
	}

	pvolLdevID := 307
	muNumber := 3
	retentionPeriod := 120 // hours

	request := model.SetSnapshotRetentionPeriodRequest{
		Parameters: model.SetSnapshotRetentionPeriodParams{
			RetentionPeriod: &retentionPeriod,
		},
	}

	resIds, err := psm.SetSnapshotRetentionPeriod(pvolLdevID, muNumber, request)
	if err != nil {
		t.Errorf("Unexpected error calling SetSnapshotRetentionPeriod: %v", err)
		return
	}

	t.Logf("Successfully initiated SetSnapshotRetentionPeriod. %v", resIds)
}
