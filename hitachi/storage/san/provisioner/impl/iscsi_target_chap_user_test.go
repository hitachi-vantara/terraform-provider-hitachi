package sanstorage

import (
	"fmt"

	sanmodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"
	"testing"
)

// newIscsiTargetChapUserTestManager is for Testing and provide structure information for connection
func newIscsiTargetChapUserTestManager() (*sanStorageManager, error) {

	// Following storage has iscsi port
	objStorageIscsi := sanmodel.StorageDeviceSettings{
		Serial:   30078,
		Username: "bXNfdm13YXJl",
		Password: "SGl0YWNoaTE=",
		MgmtIP:   "172.25.47.120",
	}
	psm, err := newSanStorageManagerEx(objStorageIscsi)
	if err != nil {
		return nil, fmt.Errorf("unexpected error while creating newIscsiTargetChapUserTestManager %v", err)
	}
	return psm, nil
}

// go test -v -run TestGetChapUsers
func xTestGetChapUsers(t *testing.T) {
	psm, err := newIscsiTargetChapUserTestManager()
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
	psm, err := newIscsiTargetChapUserTestManager()
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

// go test -v -run xTestCreateChapUser
func xTestCreateChapUser(t *testing.T) {
	psm, err := newIscsiTargetChapUserTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	portId := "CL4-C"
	iscsiTargetNumber := 1
	wayOfChapUser := "INI"
	chapUserName := "muChap"
	chapUserSecret := "abcd"

	err = psm.CreateChapUser(portId, iscsiTargetNumber, wayOfChapUser, chapUserName, chapUserSecret)
	if err != nil {
		t.Errorf("Unexpected error in Create chapuser %v", err)
		return
	}
	t.Logf("Chap User created successfully...")
}

// go test -v -run TestDeleteChapUser
func xTestDeleteChapUser(t *testing.T) {
	psm, err := newIscsiTargetChapUserTestManager()
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

// go test -v -run TestChangeChapUserSecret
func xTestChangeChapUserSecret(t *testing.T) {
	psm, err := newIscsiTargetChapUserTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	portId := "CL4-C"
	iscsiTargetNumber := 1
	wayOfChapUser := "INI"
	chapUserName := "muChap"
	chapUserSecret := "abcd"

	err = psm.ChangeChapUserSecret(portId, iscsiTargetNumber, wayOfChapUser, chapUserName, chapUserSecret)
	if err != nil {
		t.Errorf("Unexpected error in Change chapuser %v", err)
		return
	}
	t.Logf("Chap User changed successfully...")
}
