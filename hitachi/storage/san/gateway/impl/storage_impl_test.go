package sanstorage

import (
	"fmt"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	"testing"
)

// newTestManager is for Testing and provide structure information for connection
func newTestManager() (*sanStorageManager, error) {

	objStorage := sanmodel.StorageDeviceSettings{
		Serial:   40014,
		Username: "ms_vmware",
		Password: "Hitachi1",
		MgmtIP:   "172.25.47.115",
	}
	psm, err := newSanStorageManagerEx(objStorage)
	if err != nil {
		return nil, fmt.Errorf("unexpected error while creating newSanStorageManagerEx %v", err)
	}
	return psm, nil
}

// go test -v -run TestGetStorageSystemInfo
func xTestGetStorageSystemInfo(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetStorageSystemInfo()
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestGetStorageCapacity
func xTestGetStorageCapacity(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	/*
		OUTPUT:
		======
		[{
			"internal" : {
			  "freeSpace" : 83975831040,
			  "totalCapacity" : 88615368192
			},
			"external" : {
			  "freeSpace" : 0,
			  "totalCapacity" : 41943040
			},
			"total" : {
			  "freeSpace" : 83975831040,
			  "totalCapacity" : 88657311232
			}
		  }]
	*/

	resp, err := psm.GetStorageCapacity()
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}
