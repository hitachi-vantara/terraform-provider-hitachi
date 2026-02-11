package sanstorage

import (
	"fmt"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"
	sangatewaymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	"testing"
)

// newTestManager is for Testing and provide structure information for connection
func newVolumeTestManager() (*sanStorageManager, error) {

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
	var sizeInGB uint = 2
	var dpool int = 0
	// dataReductionMode := "compression"
	// sharedVol := true
	// compAccEnabled := true

	size := fmt.Sprintf("%dG", sizeInGB)

	reqBody := sangatewaymodel.CreateLunRequestGwy{
		LdevID:             &ldevId,
		PoolID:             &dpool,
		ByteFormatCapacity: size,
		// DataReductionMode:  &dataReductionMode,
		// IsDataReductionSharedVolumeEnabled: &sharedVol,
		// IsCompressionAccelerationEnabled: &compAccEnabled,
	}

	newLdevId, err := psm.CreateLun(reqBody)
	if err != nil {
		t.Errorf("Unexpected error in CreateLunInDynamicPoolWithLDevId %v", err)
		return
	}
	t.Logf("newLdevId: %v", *newLdevId)

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

	var sizeInGB uint = 2
	var dpool int = 0
	dataReductionMode := "compression"
	sharedVol := true
	compAccEnabled := true

	size := fmt.Sprintf("%dG", sizeInGB)

	reqBody := sangatewaymodel.CreateLunRequestGwy{
		// LdevID:             &ldevId,
		PoolID:             &dpool,
		ByteFormatCapacity: size,
		DataReductionMode:  &dataReductionMode,
		IsDataReductionSharedVolumeEnabled: &sharedVol,
		IsCompressionAccelerationEnabled: &compAccEnabled,
	}

	newLdevId, err := psm.CreateLun(reqBody)
	if err != nil {
		t.Errorf("Unexpected error in CreateLunInDynamicPoolWithLDevId %v", err)
		return
	}
	t.Logf("newLdevId: %v", *newLdevId)

	lun, err := psm.GetLun(*newLdevId)
	if err != nil {
		t.Errorf("Unexpected error in CreateLunInDynamicPoolWithLDevId %v", err)
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
	var sizeInGB uint = 2
	parityGroupId := "1-1"
	dataReductionMode := "compression"
	// sharedVol := true
	// compAccEnabled := true

	size := fmt.Sprintf("%dG", sizeInGB)

	reqBody := sangatewaymodel.CreateLunRequestGwy{
		LdevID:             &ldevId,
		ParityGroupID:      &parityGroupId,
		ByteFormatCapacity: size,
		DataReductionMode:  &dataReductionMode,
		// IsDataReductionSharedVolumeEnabled: &sharedVol,
		// IsCompressionAccelerationEnabled: &compAccEnabled,
	}

	newLdevId, err := psm.CreateLun(reqBody)
	if err != nil {
		t.Errorf("Unexpected error in CreateLunInDynamicPoolWithLDevId %v", err)
		return
	}
	t.Logf("newLdevId: %v", *newLdevId)

	lun, err := psm.GetLun(*newLdevId)
	if err != nil {
		t.Errorf("Unexpected error in CreateLunInDynamicPoolWithLDevId %v", err)
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

	// ldevId := 108
	var sizeInGB uint = 2
	parityGroupId := "1-1"
	dataReductionMode := "compression"
	// sharedVol := true
	// compAccEnabled := true

	size := fmt.Sprintf("%dG", sizeInGB)

	reqBody := sangatewaymodel.CreateLunRequestGwy{
		// LdevID:             &ldevId,
		ParityGroupID:      &parityGroupId,
		ByteFormatCapacity: size,
		DataReductionMode:  &dataReductionMode,
		// IsDataReductionSharedVolumeEnabled: &sharedVol,
		// IsCompressionAccelerationEnabled: &compAccEnabled,
	}

	newLdevId, err := psm.CreateLun(reqBody)
	if err != nil {
		t.Errorf("Unexpected error in CreateLunInDynamicPoolWithLDevId %v", err)
		return
	}
	t.Logf("newLdevId: %v", *newLdevId)

	lun, err := psm.GetLun(*newLdevId)
	if err != nil {
		t.Errorf("Unexpected error in CreateLunInDynamicPoolWithLDevId %v", err)
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

	var sizeInGB uint64 = 4
	ldevId := 107
	size := fmt.Sprintf("%dG", sizeInGB)

	newLdevId, err := psm.ExpandLun(ldevId, size)
	if err != nil {
		t.Errorf("Unexpected error in ExpandLun %v", err)
		return
	}
	t.Logf("newLdevId: %v", *newLdevId)

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

	ldevId := 125

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

	ldevId := 107
	// label := "Updated_Lun_Label_107"
	// dataReductionMode := "compression"
	// reductionMode := "post_process"
	compAccEnabled := true

	reqBody := sangatewaymodel.UpdateLunRequestGwy{
		// Label:             &label,
		// DataReductionMode:  &dataReductionMode,
		// DataReductionProcessMode: &reductionMode,
		IsCompressionAccelerationEnabled: &compAccEnabled,
	}

	_, err = psm.UpdateLun(ldevId, reqBody)
	if err != nil {
		t.Errorf("Unexpected error in UpdateLun %v", err)
		return
	}
	t.Logf("ldevId: %v", ldevId)
}
