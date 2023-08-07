package sanstorage

import (
	"fmt"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	"testing"
)

// newResourceTestManager is for Testing and provide structure information for connection
func newResourceTestManager() (*sanStorageManager, error) {

	objStorageIscsi := sanmodel.StorageDeviceSettings{
		Serial:   30078,
		Username: "ms_vmware",
		Password: "Hitachi1",
		MgmtIP:   "172.25.47.120",
	}
	psm, err := newSanStorageManagerEx(objStorageIscsi)
	if err != nil {
		return nil, fmt.Errorf("unexpected error while creating newResourceTestManager %v", err)
	}
	return psm, nil
}

// go test -v -run TestLockResources
func TestLockResources(t *testing.T) {
	psm, err := newResourceTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	reqBody := sanmodel.LockResourcesReq{}
	reqBody.Parameters.WaitTime = 30
	err = psm.LockResources(reqBody)
	if err != nil {
		t.Errorf("Unexpected error in LockResources %v", err)
		return
	}
}

// go test -v -run TestUnlockResources
func TestUnlockResources(t *testing.T) {
	psm, err := newResourceTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	err = psm.UnlockResources()
	if err != nil {
		t.Errorf("Unexpected error in UnlockResources %v", err)
		return
	}
}
