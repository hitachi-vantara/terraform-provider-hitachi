package sanstorage

import (
	"fmt"
	sangatewaymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"
	"testing"
)

// newSnapshotTestManager is for Testing and provides structure information for connection
func newSnapshotTestManager() (*sanStorageManager, error) {
	objStorage := sanmodel.StorageDeviceSettings{
		Serial:   12345,
		Username: "user1",
		Password: "mypswd",
		MgmtIP:   "10.10.11.12",
	}
	psm, err := newSanStorageManagerEx(objStorage)
	if err != nil {
		return nil, fmt.Errorf("unexpected error while creating newSanStorageManagerEx %v", err)
	}
	return psm, nil
}

// go test -v -run TestGetSnapshots
func xTestGetSnapshots(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	params := sangatewaymodel.GetSnapshotsParams{}
	resp, err := psm.GetSnapshots(params)
	if err != nil {
		t.Errorf("Unexpected error in GetSnapshots: %v", err)
		return
	}
	t.Logf("Response: %+v", resp)
}

// go test -v -run TestGetSnapshot
func xTestGetSnapshot(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetSnapshot(100, 1)
	if err != nil {
		t.Errorf("Unexpected error in GetSnapshot: %v", err)
		return
	}
	t.Logf("Response: %+v", resp)
}

// go test -v -run TestGetSnapshotReplicationsRange
func xTestGetSnapshotReplicationsRange(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	params := sangatewaymodel.GetSnapshotReplicationsRangeParams{}
	resp, err := psm.GetSnapshotReplicationsRange(params)
	if err != nil {
		t.Errorf("Unexpected error in GetSnapshotReplicationsRange: %v", err)
		return
	}
	t.Logf("Response: %+v", resp)
}

// go test -v -run TestCreateSnapshot
func xTestCreateSnapshot(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	request := sangatewaymodel.CreateSnapshotParams{}
	jobID, err := psm.CreateSnapshot(request)
	if err != nil {
		t.Errorf("Unexpected error in CreateSnapshot: %v", err)
		return
	}
	t.Logf("JobID: %v", jobID)
}

// go test -v -run TestSplitSnapshot
func xTestSplitSnapshot(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	request := sangatewaymodel.SplitSnapshotRequest{}
	jobID, err := psm.SplitSnapshot(100, 1, request)
	if err != nil {
		t.Errorf("Unexpected error in SplitSnapshot: %v", err)
		return
	}
	t.Logf("JobID: %v", jobID)
}

// go test -v -run TestResyncSnapshot
func xTestResyncSnapshot(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	request := sangatewaymodel.ResyncSnapshotRequest{}
	jobID, err := psm.ResyncSnapshot(100, 1, request)
	if err != nil {
		t.Errorf("Unexpected error in ResyncSnapshot: %v", err)
		return
	}
	t.Logf("JobID: %v", jobID)
}

// go test -v -run TestRestoreSnapshot
func xTestRestoreSnapshot(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	request := sangatewaymodel.RestoreSnapshotRequest{}
	jobID, err := psm.RestoreSnapshot(100, 1, request)
	if err != nil {
		t.Errorf("Unexpected error in RestoreSnapshot: %v", err)
		return
	}
	t.Logf("JobID: %v", jobID)
}

// go test -v -run TestDeleteSnapshot
func xTestDeleteSnapshot(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	jobID, err := psm.DeleteSnapshot(100, 1)
	if err != nil {
		t.Errorf("Unexpected error in DeleteSnapshot: %v", err)
		return
	}
	t.Logf("JobID: %v", jobID)
}

// go test -v -run TestCloneSnapshot
func xTestCloneSnapshot(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	request := sangatewaymodel.CloneSnapshotRequest{}
	jobID, err := psm.CloneSnapshot(100, 1, request)
	if err != nil {
		t.Errorf("Unexpected error in CloneSnapshot: %v", err)
		return
	}
	t.Logf("JobID: %v", jobID)
}

// go test -v -run TestAssignSnapshotVolume
func xTestAssignSnapshotVolume(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	request := sangatewaymodel.AssignSnapshotVolumeRequest{}
	jobID, err := psm.AssignSnapshotVolume(100, 1, request)
	if err != nil {
		t.Errorf("Unexpected error in AssignSnapshotVolume: %v", err)
		return
	}
	t.Logf("JobID: %v", jobID)
}

// go test -v -run TestUnassignSnapshotVolume
func xTestUnassignSnapshotVolume(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	jobID, err := psm.UnassignSnapshotVolume(100, 1)
	if err != nil {
		t.Errorf("Unexpected error in UnassignSnapshotVolume: %v", err)
		return
	}
	t.Logf("JobID: %v", jobID)
}

// go test -v -run TestSetSnapshotRetentionPeriod
func xTestSetSnapshotRetentionPeriod(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	request := sangatewaymodel.SetSnapshotRetentionPeriodRequest{}
	jobID, err := psm.SetSnapshotRetentionPeriod(100, 1, request)
	if err != nil {
		t.Errorf("Unexpected error in SetSnapshotRetentionPeriod: %v", err)
		return
	}
	t.Logf("JobID: %v", jobID)
}
