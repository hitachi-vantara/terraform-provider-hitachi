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

// go test -v -run TestGetDashboardInfo
func xTestGetDashboardInfo(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetDashboardInfo()
	if err != nil {
		t.Errorf("Unexpected error in GetDashboardInfo %v", err)
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

	resp, err := psm.GetDrivesInfo("")
	if err != nil {
		t.Errorf("Unexpected error in GetDrivesInfo %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}
