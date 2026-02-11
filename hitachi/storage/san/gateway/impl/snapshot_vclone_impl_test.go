package sanstorage

import (
	"testing"
	"encoding/json"
)

// go test -v -run TestGetVirtualCloneParentVolumes
func xTestGetVirtualCloneParentVolumes(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetVirtualCloneParentVolumes()
	if err != nil {
		t.Errorf("Unexpected error in GetVirtualCloneParentVolumes: %v", err)
		return
	}

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Logf("Failed to marshal Snapshot: %v", err)
		return
	}

	t.Logf("Response:\n%s", string(data))
}

// go test -v -run ^TestGetSnapshotFamily$
func xTestGetSnapshotFamily(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	ldevID := 681
	resp, err := psm.GetSnapshotFamily(ldevID)
	if err != nil {
		t.Errorf("Unexpected error in GetSnapshotFamily: %v", err)
		return
	}

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Logf("Failed to marshal Snapshot: %v", err)
		return
	}

	t.Logf("Response:\n%s", string(data))
}

// snapshot vclone tests

// go test -v -run ^TestCreateSnapshotVClone$
func xTestCreateSnapshotVClone(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	jobID, err := psm.CreateSnapshotVClone(681, 3)
	if err != nil {
		t.Errorf("Unexpected error in CreateSnapshotVClone: %v", err)
		return
	}
	t.Logf("JobID: %v", jobID)
}

// go test -v -run TestConvertSnapshotVClone
func xTestConvertSnapshotVClone(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	jobID, err := psm.ConvertSnapshotVClone(307, 3)
	if err != nil {
		t.Errorf("Unexpected error in ConvertSnapshotVClone: %v", err)
		return
	}
	t.Logf("JobID: %v", jobID)
}

// go test -v -run ^TestRestoreSnapshotFromVClone$
func xTestRestoreSnapshotFromVClone(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	jobID, err := psm.RestoreSnapshotFromVClone(1938, 3)
	if err != nil {
		t.Errorf("Unexpected error in RestoreSnapshotFromVClone: %v", err)
		return
	}
	t.Logf("JobID: %v", jobID)
}

// snapshot group virtual clone tests

// go test -v -run TestCreateSnapshotGroupVClone
func xTestCreateSnapshotGroupVClone(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	jobID, err := psm.CreateSnapshotGroupVClone("halbaSnapshotGroup1")
	if err != nil {
		t.Errorf("Unexpected error in CreateSnapshotGroupVClone: %v", err)
		return
	}
	t.Logf("JobID: %v", jobID)
}

// go test -v -run TestConvertSnapshotGroupVClone
func xTestConvertSnapshotGroupVClone(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	jobID, err := psm.ConvertSnapshotGroupVClone("halbaSnapshotGroup1")
	if err != nil {
		t.Errorf("Unexpected error in ConvertSnapshotGroupVClone: %v", err)
		return
	}
	t.Logf("JobID: %v", jobID)
}

// go test -v -run TestRestoreSnapshotGroupFromVClone
func xTestRestoreSnapshotGroupFromVClone(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	jobID, err := psm.RestoreSnapshotGroupFromVClone("halbaSnapshotGroup1")
	if err != nil {
		t.Errorf("Unexpected error in RestoreSnapshotGroupFromVClone: %v", err)
		return
	}
	t.Logf("JobID: %v", jobID)
}
