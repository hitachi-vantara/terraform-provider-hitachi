package sanstorage

import (
	"fmt"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	"testing"
)

// newResourceTestManager is for Testing and provide structure information for connection
func newResourceTestManager() (*sanStorageManager, error) {

	objStorageIscsi := sanmodel.StorageDeviceSettings{
		Serial:   12345,
		Username: "user1",
		Password: "mypswd",
		MgmtIP:   "10.10.11.12",
	}
	psm, err := newSanStorageManagerEx(objStorageIscsi)
	if err != nil {
		return nil, fmt.Errorf("unexpected error while creating newResourceTestManager %v", err)
	}
	return psm, nil
}

// go test -v -run TestLockResources
func xTestLockResources(t *testing.T) {
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
func xTestUnlockResources(t *testing.T) {
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
