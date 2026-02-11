package vssbstorage

import (
	// vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"
	"testing"
)

// go test -v -run TestGetAllStorageNodes
// go test -v -run TestGetAll
func xTestGetAll(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetStorageNodes()
	if err != nil {
		t.Errorf("Unexpected error in GetAllStorageNodes %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}
