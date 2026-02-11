package admin

import (
	"encoding/json"
	"testing"

	model "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
)

// go test -v -run TestGateway_GetIscsiTargets
func xTestGateway_GetIscsiTargets(t *testing.T) {
	gwy, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating gateway manager: %v", err)
	}

	serverId := 290
	resp, err := gwy.GetIscsiTargets(serverId)
	if err != nil {
		t.Errorf("Unexpected error in GetIscsiTargets: %v", err)
		return
	}

	if resp == nil {
		t.Errorf("Expected response, got nil")
		return
	}

	logIscsiTargetInfoList(t, resp)
}

// go test -v -run TestGateway_GetIscsiTargetByPort
func xTestGateway_GetIscsiTargetByPort(t *testing.T) {
	gwy, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating gateway manager: %v", err)
	}

	serverId := 290
	portId := "CL1-C"

	resp, err := gwy.GetIscsiTargetByPort(serverId, portId)
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

// go test -v -run TestGateway_ChangeIscsiTargetName
func xTestGateway_ChangeIscsiTargetName(t *testing.T) {
	gwy, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating gateway manager: %v", err)
	}

	serverId := 290
	portId := "CL1-C"
	targetIscsiName := "iqn.1994-04.jp.co.hitachi:rsd.has.t.10045.1c00e"

	err = gwy.ChangeIscsiTargetName(serverId, portId, targetIscsiName)
	if err != nil {
		t.Errorf("Unexpected error in ChangeIscsiTargetName: %v", err)
		return
	}

	t.Logf("Successfully changed iSCSI target name")
}

// -------------- Helper Functions --------------
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
