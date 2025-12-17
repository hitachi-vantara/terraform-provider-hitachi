package admin

import (
	"encoding/json"
	model "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	"testing"
)

// go test -v -run TestGetIscsiTargets
func xTestGetIscsiTargets(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	serverId := 290
	resp, err := psm.GetIscsiTargets(serverId)
	if err != nil {
		t.Errorf("Unexpected error calling GetIscsiTargets: %v", err)
		return
	}

	if resp == nil {
		t.Errorf("Expected response, got nil")
		return
	}

	logIscsiTargetInfoList(t, resp)
}

// go test -v -run TestGetIscsiTargetByPort
func xTestGetIscsiTargetByPort(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	serverId := 290
	portId := "CL1-C"

	resp, err := psm.GetIscsiTargetByPort(serverId, portId)
	if err != nil {
		t.Errorf("Unexpected error in GetIscsiTargetByPort: %v", err)
		return
	}

	if resp == nil {
		t.Errorf("Expected non-nil response from GetIscsiTargetByPort")
		return
	}

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Logf("Failed to marshal IscsiTargetInfoByPort: %v", err)
		return
	}

	t.Logf("Response:\n%s", string(data))
}

// go test -v -run TestChangeIscsiTargetName
func xTestChangeIscsiTargetName(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	serverId := 290
	portId := "CL1-C"
	targetIscsiName := "iqn.1994-04.jp.co.hitachi:rsd.has.t.10045.1c00e"
	// targetIscsiName := "iqn.aaa"

	err = psm.ChangeIscsiTargetName(serverId, portId, targetIscsiName)
	if err != nil {
		t.Errorf("Unexpected error in ChangeIscsiTargetName: %v", err)
		return
	}

	t.Logf("Successfully updated iscsitarget name")
}

func logIscsiTargetInfoList(t *testing.T, resp *model.IscsiTargetInfoList) {
	if resp == nil {
		t.Log("Response: <nil>")
		return
	}

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Logf("Failed to marshal IscsiTargetInfoList: %v", err)
		return
	}

	t.Logf("Response:\n%s", string(data))
}
