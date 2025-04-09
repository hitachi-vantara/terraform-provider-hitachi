package vssbstorage

import (
	"fmt"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/provisioner/model"
	"testing"
)

func newChapUserTestManager() (*vssbStorageManager, error) {

	objStorage := vssbmodel.StorageDeviceSettings{
		Username:       "admin",    // admin
		Password:       "vssb-789", // vssb-789
		ClusterAddress: "10.76.47.55",
	}
	psm, err := newVssbStorageManagerEx(objStorage)
	if err != nil {
		return nil, fmt.Errorf("unexpected error while creating newVssbStorageManagerEx %v", err)
	}
	return psm, nil
}

// go test -v -run TestGetAllChapUser
func xTestGetAllChapUser(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetAllChapUsers()
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestCreateChapUser
func xTestCreateChapUser(t *testing.T) {
	psm, err := newChapUserTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	req := vssbmodel.ChapUserReq{
		TargetChapUserName: "rahul",
		TargetChapSecret:   "HitachiVantara",
	}

	err1 := psm.CreateChapUser(&req)
	if err1 != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Chap User Created successfully")
}

// go test -v -run TestDeleteChapUser
func xTestDeleteChapUser(t *testing.T) {
	psm, err := newChapUserTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	chapUserId := "eb95f496-973a-4831-9c0c-a3292bbe6ff6"

	err = psm.DeleteChapUser(chapUserId)
	if err != nil {
		t.Errorf("Unexpected error in DeleteChapUser %v", err)
		return
	}
	t.Logf("Successfully deleted chap user: %v", chapUserId)
}
