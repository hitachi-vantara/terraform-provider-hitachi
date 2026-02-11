package admin

import (
	"encoding/json"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	"testing"
)

// go test -v -run TestGetVolumeServerConnections
func xTestGetVolumeServerConnections(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	serverId := 10
	startVolumeId := 19
	count := 10
	// serverNickname := "serverNickname_example"

	params := gwymodel.GetVolumeServerConnectionsParams{
		ServerId:       &serverId,
		ServerNickname: nil, // or &serverNickname
		StartVolumeId:  &startVolumeId,
		Count:          &count,
	}

	resp, err := psm.GetVolumeServerConnections(params)
	if err != nil {
		t.Errorf("Unexpected error in GetVolumeServerConnections: %v", err)
		return
	}

	if resp == nil {
		t.Errorf("Expected non-nil response from GetVolumeServerConnections")
		return
	}

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Logf("Failed to marshal response: %v", err)
		return
	}

	t.Logf("GetVolumeServerConnections Response:\n%s", string(data))
}

// go test -v -run TestGetOneVolumeServerConnection
func xTestGetOneVolumeServerConnection(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	serverId := 9
	volumeId := 6079

	resp, err := psm.GetOneVolumeServerConnection(volumeId, serverId)
	if err != nil {
		t.Errorf("Unexpected error in GetOneVolumeServerConnection: %v", err)
		return
	}

	if resp == nil {
		t.Errorf("Expected non-nil response from GetOneVolumeServerConnection for VolumeID=%d ServerID=%d", volumeId, serverId)
		return
	}

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Logf("Failed to marshal VolumeServerConnectionDetail: %v", err)
		return
	}

	t.Logf("GetOneVolumeServerConnection Response:\n%s", string(data))
}

// go test -v -run TestAttachVolumeToServers
func xTestAttachVolumeToServers(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	params := gwymodel.AttachVolumeServerConnectionParam{
		ServerIds: []int{9, 10},
		VolumeIds: []int{6079},
	}

	resp, err := psm.AttachVolumeToServers(params)
	if err != nil {
		t.Errorf("Unexpected error in AttachVolumeToServers: %v", err)
		return
	}

	if resp == "" {
		t.Errorf("Expected non-empty connection ID(s) from AttachVolumeToServers")
		return
	}

	t.Logf("AttachVolumeToServers Response: %s", resp)
}

// go test -v -run TestDetachVolumeFromServer
func xTestDetachVolumeFromServer(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	volumeId := 6079
	serverId := 9

	err = psm.DetachVolumeFromServer(volumeId, serverId)
	if err != nil {
		t.Errorf("Unexpected error in DetachVolumeFromServer for VolumeID=%d ServerID=%d: %v", volumeId, serverId, err)
		return
	}

	t.Logf("Successfully detached volume ID %d from server ID %d", volumeId, serverId)
}
