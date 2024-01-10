package sanstorage

import (
	"fmt"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"
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

	// Without Filters
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

// Sample:
/*"data" : [ {
  "parityGroupId" : "1-2",
  "numOfLdevs" : 252,
  "usedCapacityRate" : 99,
  "availableVolumeCapacity" : 0,
  "raidLevel" : "RAID5",
  "raidType" : "3D+1P",
  "clprId" : 0,
  "driveType" : "SLB5F-M960SS",
  "driveTypeName" : "SSD(MLC)",
  "totalCapacity" : 2640,
  "physicalCapacity" : 2640,
  "availablePhysicalCapacity" : 0,
  "isAcceleratedCompressionEnabled" : false,
  "availableVolumeCapacityInKB" : 200448
}] */
