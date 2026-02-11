package vssbstorage

import (
	"testing"
)

// go test -v -run TestGetStorageVersionInfo
func xTestGetStorageVersionInfo(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetStorageVersionInfo()
	if err != nil {
		t.Errorf("Unexpected error in GetStorageVersionInfo %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestGetHealthStatuses
func xTestGetHealthStatuses(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetHealthStatuses()
	if err != nil {
		t.Errorf("Unexpected error in GetHealthStatuses %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestGetStorageClusterInfo
func xTestGetStorageClusterInfo(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetStorageClusterInfo()
	if err != nil {
		t.Errorf("Unexpected error in GetStorageClusterInfo %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestGetDrivesInfo
func xTestGetDrivesInfo(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetDrivesInfo("Normal")
	if err != nil {
		t.Errorf("Unexpected error in GetDrivesInfo %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}
