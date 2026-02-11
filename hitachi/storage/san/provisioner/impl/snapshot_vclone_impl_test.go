package sanstorage

import (
	"testing"
)

// go test -v -run TestGetVirtualCloneParentVolumesProvisioner
func xTestGetVirtualCloneParentVolumesProvisioner(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetVirtualCloneParentVolumes()
	if err != nil {
		t.Errorf("Unexpected error in GetVirtualCloneParentVolumes: %v", err)
		return
	}
	t.Logf("Response: %+v", resp)
}

// go test -v -run TestGetSnapshotFamilyProvisioner
func xTestGetSnapshotFamilyProvisioner(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	ldevID := 307
	resp, err := psm.GetSnapshotFamily(ldevID)
	if err != nil {
		t.Errorf("Unexpected error in GetSnapshotFamily: %v", err)
		return
	}
	t.Logf("Response: %+v", resp)
}

// go test -v -run TestCreateSnapshotVCloneProvisioner
func xTestCreateSnapshotVCloneProvisioner(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	jobID, err := psm.CreateSnapshotVClone(307, 3)
	if err != nil {
		t.Errorf("Unexpected error in CreateSnapshotVClone: %v", err)
		return
	}
	t.Logf("JobID: %v", jobID)
}

// go test -v -run TestConvertSnapshotVCloneProvisioner
func xTestConvertSnapshotVCloneProvisioner(t *testing.T) {
	psm, err := newSnapshotTestManager()
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

// go test -v -run TestCreateSnapshotGroupVCloneProvisioner
func xTestCreateSnapshotGroupVCloneProvisioner(t *testing.T) {
	psm, err := newSnapshotTestManager()
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

// go test -v -run TestConvertSnapshotGroupVCloneProvisioner
func xTestConvertSnapshotGroupVCloneProvisioner(t *testing.T) {
	psm, err := newSnapshotTestManager()
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

// go test -v -run TestRestoreSnapshotFromVCloneProvisioner
func xTestRestoreSnapshotFromVCloneProvisioner(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	jobID, err := psm.RestoreSnapshotFromVClone(307, 3)
	if err != nil {
		t.Errorf("Unexpected error in RestoreSnapshotFromVClone: %v", err)
		return
	}
	t.Logf("JobID: %v", jobID)
}

// go test -v -run TestRestoreSnapshotGroupFromVCloneProvisioner
func xTestRestoreSnapshotGroupFromVCloneProvisioner(t *testing.T) {
	psm, err := newSnapshotTestManager()
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
