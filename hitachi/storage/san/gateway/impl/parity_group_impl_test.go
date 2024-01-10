package sanstorage

import (
	"fmt"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	"testing"
)

// newParityGroupTestManager is for Testing and provide structure information for connection
func newParityGroupTestManager() (*sanStorageManager, error) {

	// Following storage has iscsi port
	objStorageIscsi := sanmodel.StorageDeviceSettings{
		Serial:   30078,
		Username: "ms_vmware",
		Password: "Hitachi1",
		MgmtIP:   "172.25.47.120",
	}
	psm, err := newSanStorageManagerEx(objStorageIscsi)
	if err != nil {
		return nil, fmt.Errorf("unexpected error while creating newParityGroupTestManager %v", err)
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
	t.Logf("Response: %+v", resp)
}
