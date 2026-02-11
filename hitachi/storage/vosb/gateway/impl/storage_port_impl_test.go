package vssbstorage

import (
	"testing"
)

// go test -v -run TestGetStoragePorts
func xTestGetStoragePorts(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetStoragePorts()
	if err != nil {
		t.Errorf("Unexpected error in GetStoragePorts %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}
