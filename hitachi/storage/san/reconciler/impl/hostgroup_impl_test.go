package sanstorage

import (
	"fmt"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
	"testing"
)

// newHgTestManager is for Testing and provide structure information for connection
func newHgTestManager() (*sanStorageManager, error) {

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

// go test -v -run TestGetHostGroup
func xTestGetHostGroup(t *testing.T) {
	psm, err := newHgTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetHostGroup("CL1-A", 24)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestUpdateLun
func xTestUpdateLun(t *testing.T) {
	psm, err := newHgTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	myPortID := "CL1-A"
	myHostGroupName := "TESTING-HOSTGROUP"
	myHostGroupNumber := 23
	myHostModeOptions := []int{12, 32}
	myHostMode := "AIX"

	//  25, 12 ==> TF:  25, 13 : Delete 25, 12 and Add 25, 13 (if 13 already there then error)
	//  25, 12 ==> TF:  25 ==> Delete 25, 12 and Create new 25, XX
	//  25, 12 ==> 25, 12 | 21, 13 ==> Only Add 21, 13 and skip 23, 12
	//  25, 12 => TF: 21, 12    ( Another LDEV is already mapped to LUN) Give Error

	//Ldev0 := 25
	Ldev1 := 25
	Lun1 := 12
	//Lun1 := 13
	// For Single HG Lun ID will be unique
	Ldev2 := 21
	Lun2 := 13
	//Ldev0 := 21
	Ldev3 := 22
	Lun3 := 14
	myLdevIds := []sanmodel.Luns{
		//{LdevId: &Ldev0}, // Automatic Lun ID Assing
		{LdevId: &Ldev3, Lun: &Lun3},
		{LdevId: &Ldev1, Lun: &Lun1},
		{LdevId: &Ldev2, Lun: &Lun2},
	}

	/*
		myWwn := []sanmodel.Wwn{
			{Wwn: "100000109b3dfbbb", Name: "test-wwn1b"},
			{Wwn: "100000109b3dfbbc", Name: "test-wwn1c"},
		}
	*/

	//==============
	crReq := sanmodel.CreateHostGroupRequest{
		PortID:          &myPortID,
		HostGroupName:   &myHostGroupName,
		HostGroupNumber: &myHostGroupNumber,
		HostModeOptions: myHostModeOptions,
		HostMode:        &myHostMode,
		Ldevs:           myLdevIds,
		//Wwns:            myWwn,
	}

	resp, err := psm.ReconcileHostGroup(&crReq)
	if err != nil {
		t.Errorf("Unexpected error in ReconcileHostGroup %v", err)
		return
	}
	t.Logf("Response: %v", resp)

}

// go test -v -run TestUpdateWwn
func xTestUpdateWwn(t *testing.T) {
	psm, err := newHgTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	myPortID := "CL1-A"
	myHostGroupName := "TESTING-HOSTGROUP"
	myHostGroupNumber := 23
	myHostModeOptions := []int{12, 32}
	myHostMode := "AIX"

	myWwn := []sanmodel.Wwn{
		{Wwn: "100000109b3dfbbb", Name: "test-wwn1b"},
		{Wwn: "100000109b3dfbbc", Name: "test-wwn1c"},
	}

	//==============
	crReq := sanmodel.CreateHostGroupRequest{
		PortID:          &myPortID,
		HostGroupName:   &myHostGroupName,
		HostGroupNumber: &myHostGroupNumber,
		HostModeOptions: myHostModeOptions,
		HostMode:        &myHostMode,
		Wwns:            myWwn,
	}

	resp, err := psm.ReconcileHostGroup(&crReq)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestUpdateHostGroupMode
func xTestUpdateHostGroupMode(t *testing.T) {
	psm, err := newHgTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	/*
		1) Hostgroup Mode and Option changes
		2) Hostgroup Mode removed from TF but available in existing HG
		3) Hostgroup Option removed from TF but available in existing HG
		4) Both removed
	*/

	myPortID := "CL1-A"
	myHostGroupName := "TESTING_REST_API"
	myHostGroupNumber := 23
	myHostModeOptions := []int{12, 13}
	myHostMode := "LINUX/IRIX"

	//==============
	crReq := sanmodel.CreateHostGroupRequest{
		PortID:          &myPortID,
		HostGroupName:   &myHostGroupName,
		HostGroupNumber: &myHostGroupNumber,
		HostModeOptions: myHostModeOptions,
		HostMode:        &myHostMode,
	}

	resp, err := psm.ReconcileHostGroup(&crReq)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestCreateHostGroup
func xTestCreateHostGroup(t *testing.T) {
	psm, err := newHgTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	myPortID := "CL1-A"
	myHostGroupName := "TESTING_REST_API"
	myHostGroupNumber := 23
	myHostModeOptions := []int{12, 13}
	myHostMode := "AIX"

	//Ldev0 := 25
	Ldev1 := 25
	Lun1 := 12
	// For Single HG Lun ID will be unique
	//Ldev2 := 21
	//Lun2 := 13
	myLdevIds := []sanmodel.Luns{
		//{LdevId: &Ldev0},  // Automatic Lun ID Assing
		{LdevId: &Ldev1, Lun: &Lun1},
		//{LdevId: &Ldev2, Lun: &Lun2},
	}

	myWwn := []sanmodel.Wwn{
		{Wwn: "100000109b3dfbbb", Name: "test-name"},
	}

	//==============
	crReq := sanmodel.CreateHostGroupRequest{
		PortID:          &myPortID,
		HostGroupName:   &myHostGroupName,
		HostGroupNumber: &myHostGroupNumber,
		HostModeOptions: myHostModeOptions,
		HostMode:        &myHostMode,
		Ldevs:           myLdevIds,
		Wwns:            myWwn,
	}

	resp, err := psm.ReconcileHostGroup(&crReq)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

