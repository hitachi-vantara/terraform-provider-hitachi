package sanstorage

import (
	"fmt"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
	"testing"
)

// newiScsiChapUserTestManager is for Testing and provide structure information for connection
func newiScsiChapUserTestManager() (*sanStorageManager, error) {

	objStorage := sanmodel.StorageDeviceSettings{
		Serial:   12345,
		Username: "user1",
		Password: "mypswd",
		MgmtIP:   "10.10.11.12",
	}
	psm, err := newSanStorageManagerEx(objStorage)
	if err != nil {
		return nil, fmt.Errorf("unexpected error while creating newiScsiChapUserTestManager %v", err)
	}
	return psm, nil
}

// go test -v -run TestGetChapUsers
func xTestGetChapUsers(t *testing.T) {
	psm, err := newiScsiChapUserTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	/*
		// OUTPUT
	*/
	portId := "CL4-C"
	iscsiTargetNumber := 1
	resp, err := psm.GetChapUsers(portId, iscsiTargetNumber)
	if err != nil {
		t.Errorf("Unexpected error in GetChapUsers %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestGetChapUser
func xTestGetChapUser(t *testing.T) {
	psm, err := newiScsiChapUserTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	portId := "CL4-C"
	iscsiTargetNumber := 1
	wayOfChapUser := "INI"
	chapUserName := "myChap"

	resp, err := psm.GetChapUser(portId, iscsiTargetNumber, chapUserName, wayOfChapUser)
	if err != nil {
		t.Errorf("Unexpected error in GetChapUser %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run xTestCreateChapUser
func xTestReconsileChapUser(t *testing.T) {
	psm, err := newiScsiChapUserTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	input := sanmodel.ChapUserRequest{
		PortID:            "CL4-C",
		IscsiTargetNumber: 1,
		WayOfChapUser:     "INI",
		ChapUserName:      "myChap",
		ChapUserSecret:    "abcdef123456",
	}

	chapUserInfo, err := psm.ReconcileChapUser(&input)
	if err != nil {
		t.Errorf("Unexpected error in ReconcileChapUser  %v", err)
		return
	}
	t.Logf("Response: %v", chapUserInfo)
	t.Logf("Chap User Reconciled successfully...")
}

// go test -v -run TestDeleteChapUser
func xTestDeleteChapUser(t *testing.T) {
	psm, err := newiScsiChapUserTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	portId := "CL4-C"
	iscsiTargetNumber := 1
	wayOfChapUser := "INI"
	chapUserName := "muChap"

	err = psm.DeleteChapUser(portId, iscsiTargetNumber, chapUserName, wayOfChapUser)
	if err != nil {
		t.Errorf("Unexpected error in Delete chapuser %v", err)
		return
	}
	t.Logf("Chap User deleted successfully...")
}
