package sanstorage

import (
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	"testing"
)

// go test -v -run TestGetLun
func xTestGetLun(t *testing.T) {
	/*
		OUTPUT:
		{
		  "ldevId" : 44,
		  "clprId" : 0,
		  "emulationType" : "OPEN-V-CVS",
		  "byteFormatCapacity" : "2.00 G",
		  "blockCapacity" : 4194304,
		  "attributes" : [ "CVS" ],
		  "raidLevel" : "RAID5",
		  "raidType" : "3D+1P",
		  "numOfParityGroups" : 1,
		  "parityGroupIds" : [ "1-1" ],
		  "driveType" : "SLB5G-M7R6SS",
		  "driveByteFormatCapacity" : "6.98 T",
		  "driveBlockCapacity" : 15000000128,
		  "status" : "BLK",
		  "mpBladeId" : 0,
		  "ssid" : "0004",
		  "resourceGroupId" : 0,
		  "isAluaEnabled" : false
		}
	*/
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	ldevId := 44
	resp, err := psm.GetLun(ldevId)
	if err != nil {
		t.Errorf("Unexpected error in Get Lun %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestGetAllLun
func xTestGetAllLun(t *testing.T) {
	/*
				OUTPUT:
		[{
		  "data" : [ {
		    "ldevId" : 0,
		    "clprId" : 0,
		    "emulationType" : "OPEN-V-CVS-CM",
		    "byteFormatCapacity" : "50.00 M",
		    "blockCapacity" : 102400,
		    "numOfPorts" : 1,
		    "ports" : [ {
		      "portId" : "CL1-A",
		      "hostGroupNumber" : 2,
		      "hostGroupName" : "agalica-58.180",
		      "lun" : 0
		    } ],
		    "attributes" : [ "CMD", "CVS" ],
		    "raidLevel" : "RAID5",
		    "raidType" : "3D+1P",
		    "numOfParityGroups" : 1,
		    "parityGroupIds" : [ "1-1" ],
		    "driveType" : "SLB5G-M7R6SS",
		    "driveByteFormatCapacity" : "6.98 T",
		    "driveBlockCapacity" : 15000000128,
		    "label" : "agalica-58.180-cmd",
		    "status" : "NML",
		    "mpBladeId" : 0,
		    "ssid" : "0004",
		    "resourceGroupId" : 0,
		    "isAluaEnabled" : false
		  }, {
		    "ldevId" : 1,
		    "clprId" : 0,
		    "emulationType" : "OPEN-V-CVS",
		    "byteFormatCapacity" : "20.00 G",
		    "blockCapacity" : 41943040,
		    "composingPoolId" : 0,
		    "attributes" : [ "CVS", "POOL" ],
		    "raidLevel" : "RAID5",
		    "raidType" : "3D+1P",
		    "numOfParityGroups" : 1,
		    "parityGroupIds" : [ "1-2" ],
		    "driveType" : "SLB5G-M7R6SS",
		    "driveByteFormatCapacity" : "6.98 T",
		    "driveBlockCapacity" : 15000000128,
		    "label" : "newLabel",
		    "status" : "NML",
		    "mpBladeId" : 3,
		    "ssid" : "0004",
		    "resourceGroupId" : 2,
		    "isAluaEnabled" : false
		  },
		  {
		    "ldevId" : 99,
		    "emulationType" : "NOT DEFINED",
		    "ssid" : "0004",
		    "resourceGroupId" : 8
		  } ]
		}]
	*/
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetAllLun()
	if err != nil {
		t.Errorf("Unexpected error in Get Lun %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestCreateLun
func xTestCreateLun(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	//ldevId := 43
	poolID := 0
	byteFormatCapacity := "2G"
	//dataReductionMode := "disabled" // disabled

	reqBody := sanmodel.CreateLunRequestGwy{
		//LdevID: &ldevId,
		PoolID: &poolID,
		ByteFormatCapacity: byteFormatCapacity,
		//DataReductionMode:  &dataReductionMode,
	}

	resp, err := psm.CreateLun(reqBody)
	if err != nil {
		t.Errorf("Unexpected error in Create Lun %v", err)
		return
	}
	t.Logf("Response: %v", *resp)
}

// go test -v -run TestCreateBasicVolume
func xTestCreateBasicVolume(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	//ldevId := 43
	parityGroupID := "1-1"
	byteFormatCapacity := "2G"

	reqBody := sanmodel.CreateLunRequestGwy{
		//LdevID: &ldevId,
		ParityGroupID:      &parityGroupID,
		ByteFormatCapacity: byteFormatCapacity,
	}

	resp, err := psm.CreateLun(reqBody)
	if err != nil {
		t.Errorf("Unexpected error in Create Lun %v", err)
		return
	}
	t.Logf("Response: %v", *resp)
}

// go test -v -run TestCreateExternalVolume
func xTestCreateExternalVolume(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	//ldevId := 43
	externalParityGroupID := "E1-1"
	byteFormatCapacity := "2G"

	reqBody := sanmodel.CreateLunRequestGwy{
		//LdevID: &ldevId,
		ExternalParityGroupID:      &externalParityGroupID,
		ByteFormatCapacity: byteFormatCapacity,
	}

	resp, err := psm.CreateLun(reqBody)
	if err != nil {
		t.Errorf("Unexpected error in Create Lun %v", err)
		return
	}
	t.Logf("Response: %v", *resp)
}

// go test -v -run TestUpdateLun
func xTestUpdateLun(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	dataReductionMode := "compression"
	newLabel := "newLabel"
	// dataReductionProcessMode := "inline"
	// isCompressionAccelerationEnabled := true
	// isAluaEnabled := false

	reqBody := sanmodel.UpdateLunRequestGwy{
		Label:             &newLabel,
		DataReductionMode: &dataReductionMode,
		// DataReductionProcessMode: &dataReductionProcessMode,
		// IsCompressionAccelerationEnabled: &isCompressionAccelerationEnabled,
		// IsAluaEnabled: &isAluaEnabled,
	}
	ldevId := 281

	resp, err := psm.UpdateLun(reqBody, ldevId)
	if err != nil {
		t.Errorf("Unexpected error in Update Lun %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestExpandLun
func xTestExpandLun(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	reqBody := sanmodel.ExpandLunRequestGwy{
		Parameters: sanmodel.ExpandLunParameters{
			AdditionalByteFormatCapacity: "25G",
		},
	}

	ldevId := 100
	resp, err := psm.ExpandLun(reqBody, ldevId)
	if err != nil {
		t.Errorf("Unexpected error in Expand Lun %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestDeleteLun
func xTestDeleteLun(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	ldevId := 44
	err = psm.DeleteLun(ldevId, false)
	if err != nil {
		t.Errorf("Unexpected error in Delete Lun %v", err)
		return
	}
	t.Logf("Lun deleted successfully...")
}
