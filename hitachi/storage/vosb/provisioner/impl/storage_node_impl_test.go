package vssbstorage

import (
	"testing"
)

// go test -v -run TestGetAllStoragePools
func xTestGetAll(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetStorageNodes()
	if err != nil {
		t.Errorf("Unexpected error in GetStorageNodes %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}