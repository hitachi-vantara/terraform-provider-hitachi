package sanstorage

import (
	"fmt"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	"testing"
)

// newDynamicPoolTestManager is for Testing and provide structure information for connection
func newDynamicPoolTestManager() (*sanStorageManager, error) {

	// Following storage has iscsi port
	objStorageIscsi := sanmodel.StorageDeviceSettings{
		Serial:   12345,
		Username: "user1",
		Password: "mypswd",
		MgmtIP:   "10.10.11.12",
	}
	psm, err := newSanStorageManagerEx(objStorageIscsi)
	if err != nil {
		return nil, fmt.Errorf("unexpected error while creating newDynamicPoolTestManager %v", err)
	}
	return psm, nil
}

// go test -v -run TestGetDynamicPools
func xTestGetDynamicPools(t *testing.T) {
	psm, err := newDynamicPoolTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetDynamicPools(nil, "")
	if err != nil {
		t.Errorf("Unexpected error in TestGetDynamicPools %v", err)
		return
	}
	t.Logf("Response: %+v", resp)
}

// go test -v -run TestGetDynamicPoolById
func xTestGetDynamicPoolById(t *testing.T) {
	psm, err := newDynamicPoolTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	poolId := 45
	resp, err := psm.GetDynamicPoolById(poolId)
	if err != nil {
		t.Errorf("Unexpected error in GetDynamicPoolById %v", err)
		return
	}
	t.Logf("Response: %+v", resp)
}
