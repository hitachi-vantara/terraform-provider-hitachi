package sanstorage

import (
	"fmt"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"
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

// go test -v -run TestGetAllIscsitarget
func xTestGetAllIscsitarget(t *testing.T) {
	psm, err := newIscsiTargetTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetIscsiTarget("CL4-C", 6)
	if err != nil {
		t.Errorf("Unexpected error in GetIscsiTarget %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestGetAllIscsitargets
func xTestGetAllIscsitargets(t *testing.T) {
	psm, err := newIscsiTargetTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetAllIscsiTargets()
	if err != nil {
		t.Errorf("Unexpected error in GetAllIscsiTargets %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestCreateIscsiTarget
func xTestCreateIscsiTarget(t *testing.T) {
	psm, err := newIscsiTargetTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	portID := "CL4-C"
	iscsiTargetName := "TESTING_KSH_REST_API_61"
	//myHostGroupNumber := 0
	myHostMode := "AIX"
	hostModeOptions := []int{17, 18}

	one := 1
	two := 2
	three := 3
	four := 4
	lDev := []sanmodel.IscsiLuns{
		{Lun: &one, LdevId: &two},
		{Lun: &three, LdevId: &four},
	}
	initiators := []sanmodel.Initiator{
		{IscsiTargetNameIqn: "iqn.1994-05.com.redhat:496799ba87", IscsiNickname: "ksh-test-name-1"},
		{IscsiTargetNameIqn: "iqn.1994-05.com.redhat:496799ba88", IscsiNickname: "ksh-test-name-2"},
	}
	crReq := sanmodel.CreateIscsiTargetReq{
		PortID:          portID,
		IscsiTargetName: iscsiTargetName,
		//HostGroupNumber: &myHostGroupNumber,
		HostMode:        &myHostMode,
		HostModeOptions: &hostModeOptions,
		Ldevs:           &lDev,
		Initiators:      &initiators,
	}

	resp, err := psm.CreateIscsiTarget(crReq)
	if err != nil {
		t.Errorf("Unexpected error in CreateIscsiTarget %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestDeleteIscsiTarget
func xTestDeleteIscsiTarget(t *testing.T) {
	psm, err := newIscsiTargetTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	portID := "CL4-C"
	iscscTargetNumber := 9
	err = psm.DeleteIscsiTarget(portID, iscscTargetNumber)
	if err != nil {
		t.Errorf("Unexpected error in DeleteIscsiTarget %v", err)
		return
	}
}
