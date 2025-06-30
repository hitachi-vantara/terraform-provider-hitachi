package vssbstorage

import (
	"fmt"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"
	"testing"
)

func newUserPasswordTestManager() (*vssbStorageManager, error) {

	objStorage := vssbmodel.StorageDeviceSettings{
		Username:       "user1",
		Password:       "mypswd",
		ClusterAddress: "10.10.12.13",
	}
	psm, err := newVssbStorageManagerEx(objStorage)
	if err != nil {
		return nil, fmt.Errorf("unexpected error while creating newVssbStorageManagerEx %v", err)
	}
	return psm, nil
}

// go test -v -run TestChangeUserPassword
func xTestChangeUserPassword(t *testing.T) {
	psm, err := newUserPasswordTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	req := vssbmodel.ChangeUserPasswordReq{
		CurrentPassword: "mypswd",
		NewPassword:     "mypswd",
	}
	userId := "testUser"

	resp, err1 := psm.ChangeUserPassword(userId, &req)
	if err1 != nil {
		t.Errorf("Unexpected error %v", err1)
		return
	}
	t.Logf("User Password Changed successfully")
	t.Logf("Response: %v", resp)
}
