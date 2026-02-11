package vssbstorage

import (
	"testing"
)

// go test -v -run TestGetVolumeDetails
func xTestGetVolumeDetails(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	var volumename = "Mongonode3_vol3"

	resp, err := psm.GetVolumeDetails(volumename)
	if err != nil {
		t.Errorf("Unexpected error in GetVolumeDetails %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

