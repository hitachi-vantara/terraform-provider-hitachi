package admin

import (
	model "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	"testing"
)

// go test -v -run TestAddHostGroupsToServer
func xTestAddHostGroupsToServer(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
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

	err = psm.AddHostGroupsToServer(serverId, params)
	if err != nil {
		t.Errorf("Unexpected error in AddHostGroupsToServer: %v", err)
		return
	}

	t.Logf("Successfully added hostgroups")
}

// go test -v -run TestSyncHostGroupsWithServer
func xTestSyncHostGroupsWithServer(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	serverId := 290

	err = psm.SyncHostGroupsWithServer(serverId)
	if err != nil {
		t.Errorf("Unexpected error in SyncHostGroupsWithServer: %v", err)
		return
	}

	t.Logf("Successfully added hostgroups")
}
