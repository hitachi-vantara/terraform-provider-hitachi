package sanstorage

import (
	"fmt"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"
	"testing"
)

// newHgTestManager is for Testing and provide structure information for connection
func newHgTestManager() (*sanStorageManager, error) {

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

// go test -v -run TestCreateHostGroup
func xTestCreateHostGroup(t *testing.T) {
	psm, err := newHgTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	myPortID := "CL1-A"
	myHostGroupName := "TESTING_REST_API"
	myHostGroupNumber := 0
	myHostModeOptions := []int{12, 13}
	myHostMode := "AIX"
	//myLdevIds := []int{12, 13}
	myWwn := []sanmodel.Wwn{
		{Wwn: "www", Name: "name"},
	}
	//==============
	crReq := sanmodel.CreateHostGroupRequest{
		PortID:          &myPortID,
		HostGroupName:   &myHostGroupName,
		HostGroupNumber: &myHostGroupNumber,
		HostModeOptions: myHostModeOptions,
		HostMode:        &myHostMode,
		//Ldevs:           myLdevIds,
		Wwns: myWwn,
	}

	resp, err := psm.CreateHostGroup(crReq)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestDeleteHostGroup
func xTestDeleteHostGroup(t *testing.T) {
	psm, err := newHgTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	portID := "CL1-A"
	hostGroupNumber := 24
	err = psm.DeleteHostGroup(portID, hostGroupNumber)
	if err != nil {
		t.Errorf("Unexpected error in DeleteHostGroup %v", err)
		return
	}

}

// go test -v -run TestGetAllHostGroups
func xTestGetAllHostGroups(t *testing.T) {
	psm, err := newHgTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetAllHostGroups()
	if err != nil {
		t.Errorf("Unexpected error in GetAllHostGroups %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}
