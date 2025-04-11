package vssbstorage

import (
	"testing"
)

// go test -v -run TestGetAllStoragePools
func xTestGetAllStoragePools(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetAllStoragePools()
	if err != nil {
		t.Errorf("Unexpected error in GetAllStoragePools %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestGetStoragePoolsByPoolNames
func xTestGetStoragePoolsByPoolNames(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	poolNames := []string{"SP01", "SP02"}
	resp, err := psm.GetStoragePoolsByPoolNames(poolNames)
	if err != nil {
		t.Errorf("Unexpected error in GetStoragePoolsByPoolNames %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}
