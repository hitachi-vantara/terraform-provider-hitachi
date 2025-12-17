package admin

import (
	"encoding/json"
	"testing"

	model "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
)

// go test -v -run TestGetVolumeServerConnections
func xTestGetVolumeServerConnections(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	serverId := 10
	startVolumeId := 6075
	count := 5
	// serverNickname := "serverNickname_example"

	params := model.GetVolumeServerConnectionsParams{
		ServerId:       &serverId,
		ServerNickname: nil, // or &serverNickname
		StartVolumeId:  &startVolumeId,
		Count:          &count,
	}

	resp, err := psm.GetVolumeServerConnections(params)
	if err != nil {
		t.Errorf("Unexpected error calling GetVolumeServerConnections: %v", err)
		return
	}

	if resp == nil {
		t.Errorf("Expected non-nil response from GetVolumeServerConnections")
		return
	}

	logVolumeServerConnections(t, resp)
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
		t.Errorf("Unexpected error calling GetOneVolumeServerConnection: %v", err)
		return
	}

	if resp == nil {
		t.Errorf("Expected non-nil response from GetOneVolumeServerConnection")
		return
	}

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Logf("Failed to marshal VolumeServerConnectionDetail: %v", err)
		return
	}

	t.Logf("Response:\n%s", string(data))
}

func logVolumeServerConnections(t *testing.T, resp *model.VolumeServerConnectionsResponse) {
	if resp == nil {
		t.Log("Response: <nil>")
		return
	}

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Logf("Failed to marshal VolumeServerConnectionsResponse: %v", err)
		return
	}

	t.Logf("Response:\n%s", string(data))
}

// go test -v -run TestAttachVolumeToServers
func xTestAttachVolumeToServers(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	params := model.AttachVolumeServerConnectionParam{
		ServerIds: []int{9, 10},
		VolumeIds: []int{6079},
	}

	resp, err := psm.AttachVolumeToServers(params)
	if err != nil {
		t.Errorf("Unexpected error calling AttachVolumeToServers: %v", err)
		return
	}

	if resp == "" {
		t.Errorf("Expected non-empty connection ID string from AttachVolumeToServers")
		return
	}

	t.Logf("Successfully attached volumes to servers, connection IDs: %s", resp)
}

// go test -v -run TestDetachVolumeToServers
func xTestDetachVolumeToServers(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	volumeId := 6079
	serverId := 9

	err = psm.DetachVolumeToServers(volumeId, serverId)
	if err != nil {
		t.Errorf("Unexpected error calling DetachVolumeToServers: %v", err)
		return
	}

	t.Logf("Successfully detached volume %d from server %d", volumeId, serverId)
}
