package sanstorage

import (
	"fmt"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	"testing"
)

// newIscsiTargetTestManager is for Testing and provide structure information for connection
func newIscsiTargetTestManager() (*sanStorageManager, error) {

	// Following storage has iscsi port
	objStorageIscsi := sanmodel.StorageDeviceSettings{
		Serial:   12345,
		Username: "user1",
		Password: "mypswd",
		MgmtIP:   "10.10.11.12",
	}
	psm, err := newSanStorageManagerEx(objStorageIscsi)
	if err != nil {
		return nil, fmt.Errorf("unexpected error while creating newIscsiTargetTestManager %v", err)
	}
	return psm, nil
}

// go test -v -run TestGetIscsiTarget
func xTestGetIscsiTarget(t *testing.T) {
	psm, err := newIscsiTargetTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	/*
		// OUTPUT
		{
			"hostGroupId" : "CL4-C,1",
			"portId" : "CL4-C",
			"hostGroupNumber" : 1,
			"hostGroupName" : "snewar-tgt",
			"iscsiName" : "iqn.1994-04.jp.co.hitachi:rsd.r90.t.30078.host112",
			"authenticationMode" : "NONE",
			"iscsiTargetDirection" : "S",
			"hostMode" : "VMWARE_EX"
		}
	*/
	portId := "CL4-C"
	hgNum := 7
	resp, err := psm.GetIscsiTarget(portId, hgNum)
	if err != nil {
		t.Errorf("Unexpected error in GetIscsiTarget %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestGetAllIscsiTargets
func xTestGetAllIscsiTargets(t *testing.T) {
	psm, err := newIscsiTargetTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetAllIscsiTargets()
	if err != nil {
		t.Errorf("Unexpected error in TestGetAllIscsiTargets %v", err)
		return
	}
	t.Logf("Response: %+v", resp)
}

// go test -v -run TestCreateIscsiTarget
func xTestCreateIscsiTarget(t *testing.T) {
	psm, err := newIscsiTargetTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	portID := "CL4-C"
	iscsiTargetName := "TESTING_KSH_REST_API_31"
	//myHostGroupNumber := 0
	myHostMode := "AIX"
	hostModeOptions := []int{12, 13}
	crReq := sanmodel.CreateIscsiTargetReq{
		PortID:          portID,
		IscsiTargetName: iscsiTargetName,
		//HostGroupNumber: &myHostGroupNumber,
		HostMode:        &myHostMode,
		HostModeOptions: &hostModeOptions,
	}
	portId, iscsiTargetNumber, err := psm.CreateIscsiTarget(crReq)
	if err != nil {
		t.Errorf("Unexpected error in CreateIscsiTarget %v", err)
		return
	}
	t.Logf("portId: %v, iscsiTargetNumber: %v", *portId, *iscsiTargetNumber)
}

// go test -v -run TestSetIscsiNameForIscsiTarget
func xTestSetIscsiNameForIscsiTarget(t *testing.T) {
	psm, err := newIscsiTargetTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	portID := "CL4-C"
	myHostGroupNumber := 1
	iscsiNameIqn := "iqn.1994-05.com.redhat:496799ba83"

	req := sanmodel.SetIscsiNameReq{
		PortID:             portID,
		IscsiTargetNumber:  myHostGroupNumber,
		IscsiTargetNameIqn: iscsiNameIqn,
	}

	err = psm.SetIscsiNameForIscsiTarget(req)
	if err != nil {
		t.Errorf("Unexpected error in CreateIscsiTarget %v", err)
		return
	}

}

// go test -v -run TestSetNicknameForIscsiName
func xTestSetNicknameForIscsiName(t *testing.T) {
	psm, err := newIscsiTargetTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	portID := "CL4-C"
	hostGroupNumber := 1
	iscsiNameIqn := "iqn.1994-05.com.redhat:496799ba83"
	iscsiNickName := "ksh-test-name"

	req := sanmodel.SetNicknameIscsiReq{
		IscsiNickname: iscsiNickName,
	}

	err = psm.SetNicknameForIscsiName(portID, hostGroupNumber, iscsiNameIqn, req)
	if err != nil {
		t.Errorf("Unexpected error in SetNicknameForIscsiName %v", err)
		return
	}

}

// go test -v -run TestDeleteIscsiNameFromIscsiTarget
func xTestDeleteIscsiNameFromIscsiTarget(t *testing.T) {
	psm, err := newIscsiTargetTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	portID := "CL4-C"
	hostGroupNumber := 1
	iscsiNameIqn := "iqn.1994-05.com.redhat:496799ba83"

	err = psm.DeleteIscsiNameFromIscsiTarget(portID, hostGroupNumber, iscsiNameIqn)
	if err != nil {
		t.Errorf("Unexpected error in DeleteIscsiNameFromIscsiTarget %v", err)
		return
	}

}

// go test -v -run TestAddLdevToIscsiTarget
func xTestAddLdevToIscsiTarget(t *testing.T) {
	psm, err := newIscsiTargetTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	portID := "CL4-C"
	//PortIds := []string{"CL1-A"}
	iscsiTargetNumber := 7
	ldevID := 3
	lunID := 4

	crReq := sanmodel.AddLdevToHgReqGwy{
		//PortIds:         PortIds,
		PortID:          &portID,
		HostGroupNumber: &iscsiTargetNumber,
		LdevID:          &ldevID,
		Lun:             &lunID,
	}

	err = psm.AddLdevToHG(crReq)
	if err != nil {
		t.Errorf("Unexpected error in AddLdevToHG %v", err)
		return
	}
	t.Logf("Sucessfully Added Ldev: %+v, for PortId: %s, IscsiTargetNumber: %d", crReq, portID, iscsiTargetNumber)
}

// go test -v -run TestRemoveLdevFromIscsiTarget
func xTestRemoveLdevFromIscsiTarget(t *testing.T) {
	psm, err := newIscsiTargetTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	portID := "CL4-C"
	iscsiTargetNumber := 7
	lunID := 2

	err = psm.RemoveLdevFromHG(portID, iscsiTargetNumber, lunID)
	if err != nil {
		t.Errorf("Unexpected error in RemoveLdevFromHG %v", err)
		return
	}
	t.Logf("Sucessfully Removed Ldev from Iscsi Target")
}

// go test -v -run TestGetIscsiNameInformation
func xTestGetIscsiNameInformation(t *testing.T) {
	psm, err := newIscsiTargetTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	portId := "CL4-C"
	iscsiTargetNumber := 1
	//iscsiNameIqn := "iqn.1991-05.com.microsoft:win-e4sar20hjrv"

	resp, err := psm.GetIscsiNameInformation(portId, iscsiTargetNumber)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestGetIscsiTargetLuPaths
func xTestGetIscsiTargetLuPaths(t *testing.T) {
	psm, err := newIscsiTargetTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetIscsiTargetGroupLuPaths("CL4-C", 7)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestSetIScsiTargetHostModeAndHostModeOptions
func xTestSetIScsiTargetHostModeAndHostModeOptions(t *testing.T) {
	psm, err := newIscsiTargetTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	PortId := "CL4-C"
	HostgroupNumber := 5
	//HostgroupNumber := 24
	myHostModeOptions := []int{11, 12}
	//myHostModeOptions := []int{12, 13}
	//myHostModeOptions := []int{-1}
	myHostMode := "AIX"
	//myHostMode := "LINUX/IRIX"
	// Both Set: { "hostMode": "WIN", "hostModeOptions": [12,33]}
	// RestHostMode: { "hostMode": "WIN", "hostModeOptions": [-1]}
	// Only HostMode: { "hostMode": "AIX"}

	crReq := sanmodel.SetIscsiHostModeAndOptions{
		HostMode:        myHostMode,
		HostModeOptions: &myHostModeOptions,
	}

	err = psm.SetIScsiTargetHostModeAndHostModeOptions(PortId, HostgroupNumber, crReq)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Sucessfully set HostModeAndOption: %+v, for PortId: %s, HostgroupNumber: %d", crReq, PortId, HostgroupNumber)
	/**
		BEFORE:
			"hostGroupId" : "CL4-C,5",
			"portId" : "CL4-C",
			"hostGroupNumber" : 5,
			"hostGroupName" : "TESTING_KSH_REST_API_11",
			"iscsiName" : "iqn.1994-05.com.redhat:496799ba94",
			"authenticationMode" : "BOTH",
			"iscsiTargetDirection" : "S",
			"hostMode" : "LINUX/IRIX"
		-----------
		AFTER:
			"hostGroupId" : "CL4-C,5",
			"portId" : "CL4-C",
			"hostGroupNumber" : 5,
			"hostGroupName" : "TESTING_KSH_REST_API_11",
			"iscsiName" : "iqn.1994-05.com.redhat:496799ba94",
			"authenticationMode" : "BOTH",
			"iscsiTargetDirection" : "S",
			"hostMode" : "WIN",   <=====
			"hostModeOptions" : [ 11, 12 ] <=====
	**/
}

// go test -v -run TestDeleteIscsiTarget
func xTestDeleteIscsiTarget(t *testing.T) {
	psm, err := newIscsiTargetTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	portId := "CL4-C"
	iscsiTargetNumber := 7

	err = psm.DeleteIscsiTarget(portId, iscsiTargetNumber)
	if err != nil {
		t.Errorf("Unexpected error in DeleteIscsiTarget %v", err)
		return
	}
	t.Logf("Sucessfully deleted PortId: %s, IscsiTargetNumber: %d", portId, iscsiTargetNumber)
}
