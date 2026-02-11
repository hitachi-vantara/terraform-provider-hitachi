package admin

import (
	"testing"

	model "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
)

// go test -v -run TestGateway_AddHostGroupsToServer
func xTestGateway_AddHostGroupsToServer(t *testing.T) {
	gwy, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating gateway manager: %v", err)
	}

	serverId := 290
	portId := "CL1-C"
	hostGroupName := "iscsiserver1"
	params := model.AddHostGroupsToServerParam{
		HostGroups: []model.HostGroupForAddToServerParam{
			{
				PortID:        portId,
				HostGroupID:   nil,
				HostGroupName: &hostGroupName,
			},
		},
	}

	err = gwy.AddHostGroupsToServer(serverId, params)
	if err != nil {
		t.Errorf("Unexpected error in AddHostGroupsToServer: %v", err)
		return
	}

	t.Logf("Successfully added hostgroups to server")
}

// go test -v -run TestGateway_SyncHostGroupsWithServer
func xTestGateway_SyncHostGroupsWithServer(t *testing.T) {
	gwy, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating gateway manager: %v", err)
	}

	serverId := 290

	err = gwy.SyncHostGroupsWithServer(serverId)
	if err != nil {
		t.Errorf("Unexpected error in SyncHostGroupsWithServer: %v", err)
		return
	}

	t.Logf("Successfully synchronized hostgroups names with server nickname")
}
