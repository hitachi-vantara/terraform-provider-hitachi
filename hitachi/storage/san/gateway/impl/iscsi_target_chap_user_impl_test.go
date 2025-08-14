package sanstorage

import (
	"fmt"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	"testing"
)

// newChapUserTestManager is for Testing and provide structure information for connection
func newChapUserTestManager() (*sanStorageManager, error) {

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

// go test -v -run TestGetChapUsers
func xTestGetChapUsers(t *testing.T) {
	psm, err := newChapUserTestManager()
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
	psm, err := newChapUserTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	/*
		// OUTPUT
	*/
	portId := "CL4-C"
	iscsiTargetNumber := 1
	wayOfChapUser := "INI"
	chapUserName := "muChap"

	resp, err := psm.GetChapUser(portId, iscsiTargetNumber, chapUserName, wayOfChapUser)
	if err != nil {
		t.Errorf("Unexpected error in GetChapUser %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestSetChapUserName
func xTestSetChapUserName(t *testing.T) {
	psm, err := newChapUserTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	portId := "CL4-C"
	iscsiTargetNumber := 1
	wayOfChapUser := "INI"
	chapUserName := "muChap"

	err = psm.SetChapUserName(portId, iscsiTargetNumber, wayOfChapUser, chapUserName)
	if err != nil {
		t.Errorf("Unexpected error in SetChapUserName %v", err)
		return
	}
	t.Logf("SetChapUserName Success")
}

// go test -v -run TestSetChapUserSecret
func xTestSetChapUserSecret(t *testing.T) {
	psm, err := newChapUserTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	portId := "CL4-C"
	iscsiTargetNumber := 1
	wayOfChapUser := "INI"
	chapUserName := "muChap"
	chapUserPassword := "xxx"

	err = psm.SetChapUserSecret(portId, iscsiTargetNumber, wayOfChapUser, chapUserName, chapUserPassword)
	if err != nil {
		t.Errorf("Unexpected error in SetChapUserSecret %v", err)
		return
	}
	t.Logf("SetChapUserSecret Success")
}

// go test -v -run TestDeleteChapUser
func xTestDeleteChapUser(t *testing.T) {
	psm, err := newChapUserTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	portId := "CL4-C"
	iscsiTargetNumber := 1
	wayOfChapUser := "INI"
	chapUserName := "muChap"

	err = psm.DeleteChapUser(portId, iscsiTargetNumber, wayOfChapUser, chapUserName)
	if err != nil {
		t.Errorf("Unexpected error in DeleteChapUser %v", err)
		return
	}
	t.Logf("Chap User deleted successfully...")
}
