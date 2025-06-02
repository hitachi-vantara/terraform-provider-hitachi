package vssbstorage

import (
	// vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"
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

// go test -v -run TestExpandStoragePool
func xTestExpandStoragePool(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	storagePoolName := "SP01"
	driveIds := []string{
		"0437c9f8-ec5a-4527-900b-300519321f1d",
		"cbf7b144-593e-451d-9a49-d62e6b7e1334",
	}

	err = psm.ExpandStoragePool(storagePoolName, driveIds)
	if err != nil {
		t.Errorf("Unexpected error in ExpandStoragePool %v", err)
		return
	}
}

// go test -v -run TestAddOfflineDrivesToStoragePool
func xTestAddOfflineDrivesToStoragePool(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	storagePoolName := "SP01"
	err = psm.AddOfflineDrivesToStoragePool(storagePoolName)
	if err != nil {
		t.Errorf("Unexpected error in AddOfflineDrivesToStoragePool %v", err)
		return
	}
}
