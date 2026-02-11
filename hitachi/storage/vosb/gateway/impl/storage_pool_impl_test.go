package vssbstorage

import (
	"testing"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/gateway/model"
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
		t.Errorf("Unexpected error in GetStoragePoolsByPoolName %v", err)
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

	poolId := "b40f471b-110c-4a5f-9970-2e4315c72b07"
	driveIds := []string{
		"0437c9f8-ec5a-4527-900b-300519321f1d",
		"cbf7b144-593e-451d-9a49-d62e6b7e1334",
	}
	req := vssbmodel.ExpandStoragePoolReq{
		DriveIds: driveIds,
	}

	err = psm.ExpandStoragePool(poolId, &req)
	if err != nil {
		t.Errorf("Unexpected error in ExpandStoragePool %v", err)
		return
	}
}
