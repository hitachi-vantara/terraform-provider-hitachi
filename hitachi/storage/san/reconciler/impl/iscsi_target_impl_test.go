package sanstorage

import (
	"fmt"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
	"testing"
)

// newIscsiTestManager is for Testing and provide structure information for connection
func newIscsiTestManager() (*sanStorageManager, error) {

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

// go test -v -run TestUpdateIscsiLun
func xTestUpdateIscsiLun(t *testing.T) {
	psm, err := newIscsiTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	myPortID := "CL4-C"
	myHostGroupName := "TESTING_REST_API_CL4-C5"
	myHostGroupNumber := 5
	myHostModeOptions := []int{12, 32}
	myHostMode := "AIX"

	//  16, 12 ==> TF:  16, 13 : Delete 16, 12 and Add 16, 13 (if 13 already there then error)
	//  16, 12 ==> TF:  16 ==> Delete 16, 12 and Create new 16, XX
	//  16, 12 ==> 16, 12 | 20, 13 ==> Only Add 20, 13 and skip 16, 12
	//  16, 12 => TF: 20, 12    ( Another LDEV is already mapped to LUN) Give Error

	//Ldev0 := 16
	Ldev1 := 16
	Lun1 := 12
	//Lun1 := 13//

	//Ldev2 := 20
	//Lun2 := 13

	// For Single HG Lun ID will be unique
	Ldev2 := 20
	Lun2 := 13

	//Ldev0 := 21

	//Ldev3 := 20
	//Lun3 := 14

	myLdevIds := []sanmodel.IscsiLuns{
		//{LdevId: &Ldev0}, // Automatic Lun ID Assing
		//{LdevId: &Ldev3, Lun: &Lun3},
		{LdevId: &Ldev1, Lun: &Lun1},
		{LdevId: &Ldev2, Lun: &Lun2},
	}

	//==============
	crReq := sanmodel.CreateIscsiTargetReq{
		PortID:            myPortID,
		IscsiTargetName:   myHostGroupName,
		IscsiTargetNumber: &myHostGroupNumber,
		HostModeOptions:   &myHostModeOptions,
		HostMode:          &myHostMode,
		Ldevs:             myLdevIds,
	}

	resp, err := psm.ReconcileIscsiTarget(&crReq)
	if err != nil {
		t.Errorf("Unexpected error in ReconcileIscsiTarget %v", err)
		return
	}
	t.Logf("Response: %v", resp)

}

// go test -v -run TestUpdateIscsiInitiator
func xTestUpdateIscsiInitiator(t *testing.T) {
	psm, err := newIscsiTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	myPortID := "CL4-C"
	myHostGroupName := "TESTING_REST_API_CL4-C5"
	myHostGroupNumber := 5
	myHostModeOptions := []int{12, 32}
	myHostMode := "AIX"

	myInitiator := []sanmodel.Initiator{
		{IscsiTargetNameIqn: "iqn.1994-05.com.redhat:496799baaa", IscsiNickname: "test-iscsiName1"},
		//{IscsiTargetNameIqn: "iqn.1994-05.com.redhat:496799bbbb", IscsiNickname: "test-iscsiName2"},
		//{IscsiTargetNameIqn: "iqn.1994-05.com.redhat:496799bbbb"},
	}

	//==============
	crReq := sanmodel.CreateIscsiTargetReq{
		PortID:            myPortID,
		IscsiTargetName:   myHostGroupName,
		IscsiTargetNumber: &myHostGroupNumber,
		HostModeOptions:   &myHostModeOptions,
		HostMode:          &myHostMode,
		Initiators:        myInitiator,
	}

	resp, err := psm.ReconcileIscsiTarget(&crReq)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)

	/*
			"data" : [ {
		    "hostIscsiId" : "CL4-C,5,iqn.1994-05.com.redhat%3A496799baaa",
		    "portId" : "CL4-C",
		    "hostGroupNumber" : 5,
		    "hostGroupName" : "TESTING_KSH_REST_API_11",
		    "iscsiName" : "iqn.1994-05.com.redhat:496799baaa",
		    "iscsiNickname" : "test-iscsiName1"
		  }
	*/
}

// go test -v -run TestUpdateIscsiHostGroupMode
func xTestUpdateIscsiHostGroupMode(t *testing.T) {
	psm, err := newIscsiTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	/*
		1) Hostgroup Mode and Option changes
		2) Hostgroup Mode removed from TF but available in existing HG
		3) Hostgroup Option removed from TF but available in existing HG
		4) Both removed
	*/

	myPortID := "CL4-C"
	myHostGroupName := "TESTING_REST_API_CL4-C5"
	myHostGroupNumber := 5
	myHostModeOptions := []int{12, 13}
	myHostMode := "LINUX/IRIX"

	//==============
	crReq := sanmodel.CreateIscsiTargetReq{
		PortID:            myPortID,
		IscsiTargetName:   myHostGroupName,
		IscsiTargetNumber: &myHostGroupNumber,
		HostModeOptions:   &myHostModeOptions,
		HostMode:          &myHostMode,
	}

	resp, err := psm.ReconcileIscsiTarget(&crReq)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestGetIscsiTargetsByPortIds
func xTestGetIscsiTargetsByPortIds(t *testing.T) {
	psm, err := newIscsiTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	/*
		1) Hostgroup Mode and Option changes
		2) Hostgroup Mode removed from TF but available in existing HG
		3) Hostgroup Option removed from TF but available in existing HG
		4) Both removed
	*/

	myPortID := []string{"CL4-C", "CL5-A"}

	resp, err := psm.GetIscsiTargetsByPortIds(myPortID)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}
