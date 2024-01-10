package vssbstorage

import (
	"fmt"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vssb/gateway/model"
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

// go test -v -run TestGetAllChapUsers
func xTestGetAllChapUsers(t *testing.T) {
	psm, err := newChapUserTestManager()
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

// go test -v -run TestGetChapUserInfo
func xTestGetChapUserInfo(t *testing.T) {
	psm, err := newChapUserTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	chapUserId := "eb95f496-973a-4831-9c0c-a3292bbe6ff6"

	chapUser, err := psm.GetChapUserInfo(chapUserId)
	if err != nil {
		t.Errorf("Unexpected error in DeleteChapUser %v", err)
		return
	}
	t.Logf("Successfully got chap user info: %v, chap user info %v", chapUserId, chapUser)
}

// go test -v -run TestUpdateChapUser
func xTestUpdateChapUser(t *testing.T) {
	psm, err := newChapUserTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	req := vssbmodel.ChapUserReq{
		TargetChapUserName: "rahul",
		TargetChapSecret:   "HitachiVantara",
	}
	chapUserId := "eb95f496-973a-4831-9c0c-a3292bbe6ff6"

	err1 := psm.UpdateChapUser(chapUserId, &req)
	if err1 != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Chap User Updated successfully")
}

// go test -v -run TestGetPortAuthSettings
func xTestGetPortAuthSettings(t *testing.T) {
	psm, err := newChapUserTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	portId := "5f07176a-e10d-47b7-99b4-57b93806048b"

	pas, err1 := psm.GetPortAuthSettings(portId)
	if err1 != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Successfully got port auth settings for port id %v, port auth settings %v", portId, pas)
}

// go test -v -run TestUpdatePortAuthSettings
func xTestUpdatePortAuthSettings(t *testing.T) {
	psm, err := newChapUserTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	req := vssbmodel.PortAuthSettings{
		AuthMode:         "None",
		IsMutualChapAuth: true,
	}

	portId := "5f07176a-e10d-47b7-99b4-57b93806048b"

	err1 := psm.UpdatePortAuthSettings(portId, &req)
	if err1 != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Successfully updated port auth settings for port id %v", portId)
}

// go test -v -run TestGetChapUsersAllowedToAccessPort
func xTestGetChapUsersAllowedToAccessPort(t *testing.T) {
	psm, err := newChapUserTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	portId := "5f07176a-e10d-47b7-99b4-57b93806048b"

	cuList, err1 := psm.GetChapUsersAllowedToAccessPort(portId)
	if err1 != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Successfully got chap user list allowed to access port id %v, chap users %v", portId, cuList)
}

// go test -v -run TestAllowChapUserToAccessPort
func xTestAllowChapUserToAccessPort(t *testing.T) {
	psm, err := newChapUserTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	portId := "5f07176a-e10d-47b7-99b4-57b93806048b"
	req := vssbmodel.ChapUserIdReq{
		ChapUserId: "eb95f496-973a-4831-9c0c-a3292bbe6ff6",
	}

	err1 := psm.AllowChapUserToAccessPort(portId, &req)
	if err1 != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Successfully alowwed chap userto access port id %v, chap users %v", portId, req.ChapUserId)
}

// go test -v -run TestDeletePortAccessForChapUser
func xTestDeletePortAccessForChapUser(t *testing.T) {
	psm, err := newChapUserTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	portId := "5f07176a-e10d-47b7-99b4-57b93806048b"
	chapUserId := "eb95f496-973a-4831-9c0c-a3292bbe6ff6"

	err1 := psm.DeletePortAccessForChapUser(portId, chapUserId)
	if err1 != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Successfully deleted  chap user to access port id %v, chap users %v", portId, chapUserId)
}

// go test -v -run TestGetChapUserInfoAllowedToAccessPort
func xTestGetChapUserInfoAllowedToAccessPort(t *testing.T) {
	psm, err := newChapUserTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	portId := "5f07176a-e10d-47b7-99b4-57b93806048b"
	chapUserId := "eb95f496-973a-4831-9c0c-a3292bbe6ff6"

	cuList, err1 := psm.GetChapUserInfoAllowedToAccessPort(portId, chapUserId)
	if err1 != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Successfully got chap user list allowed to access port id %v, chap users %v", portId, cuList)
}
