package sanstorage

import (
	"fmt"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"
	"testing"
)

// newTestManager is for Testing and provide structure information for connection
func newVolumeTestManager() (*sanStorageManager, error) {

	objStorage := sanmodel.StorageDeviceSettings{
		Serial:   40014,
		Username: "bXNfdm13YXJl",
		Password: "SGl0YWNoaTE=",
		MgmtIP:   "172.25.47.115",
	}
	psm, err := newSanStorageManagerEx(objStorage)
	if err != nil {
		return nil, fmt.Errorf("unexpected error while creating newSanStorageManagerEx %v", err)
	}
	return psm, nil
}

// go test -v -run TestGetUndefinedLun
func xTestGetUndefinedLun(t *testing.T) {
	psm, err := newVolumeTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetUndefinedLun(3)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestGetLun
func xTestGetLun(t *testing.T) {
	psm, err := newVolumeTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	ldevId := 106
	resp, err := psm.GetLun(ldevId)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestGetRangeOfLuns
func xTestGetRangeOfLuns(t *testing.T) {
	psm, err := newVolumeTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	startLdevID := 280
	endLdevID := 285
	IsUndefinedLdev := true

	resp, err := psm.GetRangeOfLuns(startLdevID, endLdevID, IsUndefinedLdev)
	if err != nil {
		t.Errorf("Unexpected error in GetRangeOfLuns %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestCreateLunInDynamicPoolWithLDevId
func xTestCreateLunInDynamicPoolWithLDevId(t *testing.T) {
	psm, err := newVolumeTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	ldevId := 107
	var sizeInGB float64 = 2
	var dynamicPool uint = 0
	dataReductionMode := "compression"

	newLdevId, err := psm.CreateLunInDynamicPoolWithLDevId(ldevId, sizeInGB, dynamicPool, dataReductionMode)
	if err != nil {
		t.Errorf("Unexpected error in CreateLunInDynamicPoolWithLDevId %v", err)
		return
	}
	t.Logf("newLdevId: %v", &newLdevId)

	lun, err := psm.GetLun(*newLdevId)
	if err != nil {
		t.Errorf("Unexpected error in CreateLunInDynamicPoolWithLDevId %v", err)
		return
	}
	t.Logf("Response: %v", lun)
}

// go test -v -run TestCreateLunInDynamicPool
func xTestCreateLunInDynamicPool(t *testing.T) {
	psm, err := newVolumeTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	var sizeInGB float64 = 2
	var dynamicPool uint = 0
	dataReductionMode := "compression"

	newLdevId, err := psm.CreateLunInDynamicPool(sizeInGB, dynamicPool, dataReductionMode)
	if err != nil {
		t.Errorf("Unexpected error in CreateLunInDynamicPool %v", err)
		return
	}
	t.Logf("newLdevId: %v", &newLdevId)

	lun, err := psm.GetLun(*newLdevId)
	if err != nil {
		t.Errorf("Unexpected error in CreateLunInDynamicPool %v", err)
		return
	}
	t.Logf("Response: %v", lun)
}

// go test -v -run TestCreateLunInParityGroupWithLDevId
func xTestCreateLunInParityGroupWithLDevId(t *testing.T) {
	psm, err := newVolumeTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	ldevId := 108
	var sizeInGB float64 = 2
	parityGroupId := "1-1"
	dataReductionMode := "compression"

	newLdevId, err := psm.CreateLunInParityGroupWithLDevId(ldevId, sizeInGB, parityGroupId, dataReductionMode)
	if err != nil {
		t.Errorf("Unexpected error in CreateLunInParityGroupWithLDevId %v", err)
		return
	}
	t.Logf("newLdevId: %v", &newLdevId)

	lun, err := psm.GetLun(*newLdevId)
	if err != nil {
		t.Errorf("Unexpected error in CreateLunInParityGroupWithLDevId %v", err)
		return
	}
	t.Logf("Response: %v", lun)
}

// go test -v -run TestCreateLunInParityGroup
func xTestCreateLunInParityGroup(t *testing.T) {
	psm, err := newVolumeTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	var sizeInGB float64 = 2
	parityGroupId := "1-1"
	dataReductionMode := "compression"

	newLdevId, err := psm.CreateLunInParityGroup(sizeInGB, parityGroupId, dataReductionMode)
	if err != nil {
		t.Errorf("Unexpected error in CreateLunInParityGroup %v", err)
		return
	}
	t.Logf("newLdevId: %v", &newLdevId)

	lun, err := psm.GetLun(*newLdevId)
	if err != nil {
		t.Errorf("Unexpected error in CreateLunInParityGroup %v", err)
		return
	}
	t.Logf("Response: %v", lun)
}

// go test -v -run TestExpandLun
func xTestExpandLun(t *testing.T) {
	psm, err := newVolumeTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	var sizeInGB float64 = 2
	ldevId := 221

	newLdevId, err := psm.ExpandLun(ldevId, sizeInGB)
	if err != nil {
		t.Errorf("Unexpected error in ExpandLun %v", err)
		return
	}
	t.Logf("newLdevId: %v", &newLdevId)

	lun, err := psm.GetLun(*newLdevId)
	if err != nil {
		t.Errorf("Unexpected error in ExpandLun %v", err)
		return
	}
	t.Logf("Response: %v", lun)
}

// go test -v -run TestDeleteLun
func xTestDeleteLun(t *testing.T) {
	psm, err := newVolumeTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	ldevId := 282

	err = psm.DeleteLun(ldevId)
	if err != nil {
		t.Errorf("Unexpected error in DeleteLun %v", err)
		return
	}
}

// go test -v -run TestUpdateLun
func xTestUpdateLun(t *testing.T) {
	psm, err := newVolumeTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	ldevId := 281
	var label string = "label51"
	var dataReductionMode string = "disabled"

	logicalUnit, err := psm.UpdateLun(ldevId, &label, &dataReductionMode)
	if err != nil {
		t.Errorf("Unexpected error in UpdateLun %v", err)
		return
	}
	t.Logf("logicalUnit: %v", &logicalUnit)
}
