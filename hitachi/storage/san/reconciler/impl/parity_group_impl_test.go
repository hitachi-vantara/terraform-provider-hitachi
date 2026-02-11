package sanstorage

import (
	"fmt"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
	"testing"
)

// newParityGroupTestManager is for Testing and provide structure information for connection
func newParityGroupTestManager() (*sanStorageManager, error) {

	objStorage := sanmodel.StorageDeviceSettings{
		Serial:   12345,
		Username: "user1",
		Password: "mypswd",
		MgmtIP:   "10.10.11.12",
	}
	psm, err := newSanStorageManagerEx(objStorage)
	if err != nil {
		return nil, fmt.Errorf("unexpected error while creating newSanStorageManagerEx %v", err)
	}
	return psm, nil
}

// go test -v -run TestGetParityGroups
func xTestGetParityGroups(t *testing.T) {
	psm, err := newParityGroupTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetParityGroups()
	if err != nil {
		t.Errorf("Unexpected error in GetParityGroups %v", err)
		return
	}
	t.Logf("Response: %v", resp)

	// With filters
	//filters := []string{"1-2"}
	filters := []string{"1-2", "1-3"}
	resp, err = psm.GetParityGroups(filters)
	if err != nil {
		t.Errorf("Unexpected error in GetParityGroups %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}
