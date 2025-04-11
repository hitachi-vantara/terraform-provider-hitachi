package vssbstorage

import (
	"fmt"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/gateway/model"
	"testing"
)

func newTestManager() (*vssbStorageManager, error) {

	objStorage := vssbmodel.StorageDeviceSettings{
		Username: "admin",    // admin
		Password: "vssb-789", // vssb-789
		ClusterAddress: "10.76.47.55",
		// ClusterAddress: "172.25.58.151",
	}
	psm, err := newVssbStorageManagerEx(objStorage)
	if err != nil {
		return nil, fmt.Errorf("unexpected error while creating newVssbStorageManagerEx %v", err)
	}
	return psm, nil
}

// go test -v -run TestGetComputeNode
func xTestGetComputeNode(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	serverID := "6673a1c4-7152-44ac-8145-deeee473aba5"
	resp, err := psm.GetComputeNode(serverID)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestGetAllComputeNodes
func xTestGetAllComputeNodes(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetAllComputeNodes()
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestRegisterComputeNode
func xTestRegisterComputeNode(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	serverNickName := "ComputeNode11"
	osType := "VMware"

	computeNode := &vssbmodel.ComputeNodeCreateReq{
		ServerNickname: serverNickName,
		OsType:         osType,
	}

	err = psm.RegisterComputeNode(computeNode)
	if err != nil {
		t.Errorf("Unexpected error in RegisterComputeNode %v", err)
		return
	}
	t.Logf("Successfully registered server")
}

// go test -v -run TestDeleteComputeNode
func xTestDeleteComputeNode(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	serverID := "53c96ca6-fdd4-4e03-ab1e-4725e8e4e6d9"
	err = psm.DeleteComputeNode(serverID)
	if err != nil {
		t.Errorf("Unexpected error in DeleteComputeNode %v", err)
		return
	}
	t.Logf("Successfully deleted server: %v", serverID)
}

// go test -v -run TestEditComputeNode
func xTestEditComputeNode(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	serverID := "53c96ca6-fdd4-4e03-ab1e-4725e8e4e6d9"
	nickName := "ComputeNode2"
	osType := "VMware"

	computeNode := &vssbmodel.ComputeNodeInformation{
		Nickname: nickName,
		OsType:   osType,
	}

	err = psm.EditComputeNode(serverID, computeNode)
	if err != nil {
		t.Errorf("Unexpected error in EditComputeNode %v", err)
		return
	}
	t.Logf("Successfully edited server: %v", serverID)
}

// go test -v -run TestRegisterInitiatorInfoForComputeNode
func xTestRegisterInitiatorInfoForComputeNode(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	serverID := "bbad9cb1-7a1b-4c1d-b902-ead1f889c447"
	protocol := "iSCSI"
	//iscsiName := "iqn.1998-01.com.vmware:node-06-0723aa93"
	iscsiName := "iqn.1998-01.com.vmware:node-06-0723aa94"
	// Note: Unique iScsi Name require
	reqInfo := vssbmodel.RegisterInitiator{Protocol: protocol, IscsiName: iscsiName}
	err = psm.RegisterInitiatorInfoForComputeNode(serverID, &reqInfo)
	if err != nil {
		t.Errorf("Unexpected error in RegisterInitiatorInfoForComputeNode %v", err)
		return
	}
	t.Logf("Successfully edited server: %v", serverID)
}

// go test -v -run TestRegisterHbaInfoForComputeNode
func xTestRegisterHbaInfoForComputeNode(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	serverID := "79c00bf7-067f-4ace-a878-0cd64368eefa"
	protocol := "FC"
	hbaWwn := "50060e8107595005"
	isTargetAny := false

	reqInfo := vssbmodel.RegisterHba{Protocol: protocol, HbaWwn: hbaWwn, IsTargetAny: isTargetAny}
	err = psm.RegisterHbaInfoForComputeNode(serverID, &reqInfo)
	if err != nil {
		t.Errorf("Unexpected error in RegisterHbaInfoForComputeNode %v", err)
		return
	}
	t.Logf("Successfully edited server: %v", serverID)
}

// go test -v -run TestConfigureHbaPortsForComputeNode
func xTestConfigureHbaPortsForComputeNode(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	serverID := "6da76e88-dca8-416b-89b8-04e0f5062cca"

	err = psm.ConfigureHbaPortsForComputeNode(serverID)
	if err != nil {
		t.Errorf("Unexpected error in ConfigureHbaPortsForComputeNode %v", err)
		return
	}
	t.Logf("Successfully edited server: %v", serverID)
}

// go test -v -run TestDeleteInitiatorInfoForComputeNode
func xTestDeleteInitiatorInfoForComputeNode(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	serverID := "53c96ca6-fdd4-4e03-ab1e-4725e8e4e6d9"
	hbaID := "667fb6e2-1cf7-448e-8e5d-4309756f7863"
	err = psm.DeleteInitiatorInfoForComputeNode(serverID, hbaID)
	if err != nil {
		t.Errorf("Unexpected error in DeleteInitiatorInfoForComputeNode %v", err)
		return
	}
	t.Logf("Successfully deleted server: %v", serverID)
}

// go test -v -run TestGetInitiatorInformationForComputeNode
func xTestGetInitiatorInformationForComputeNode(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	/*
		{
			"data": [
			  {
				"id": "8408ad7a-2fa2-4acd-8c50-56d908a34a21",
				"serverId": "bbad9cb1-7a1b-4c1d-b902-ead1f889c447",
				"name": "iqn.1998-01.com.vmware:node-06-0723aa94",
				"protocol": "iSCSI",
				"portIds": []
			  }
			]
		  }
	*/

	serverID := "bbad9cb1-7a1b-4c1d-b902-ead1f889c447"
	hbaID := "8408ad7a-2fa2-4acd-8c50-56d908a34a21"
	resp, err := psm.GetInitiatorInformationForComputeNode(serverID, hbaID)
	if err != nil {
		t.Errorf("Unexpected error in GetInitiatorInformationForComputeNode %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestGetInitiatorsInformationForComputeNode
func xTestGetInitiatorsInformationForComputeNode(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	/*
		{
			"data": [
			  {
				"id": "8408ad7a-2fa2-4acd-8c50-56d908a34a21",
				"serverId": "bbad9cb1-7a1b-4c1d-b902-ead1f889c447",
				"name": "iqn.1998-01.com.vmware:node-06-0723aa94",
				"protocol": "iSCSI",
				"portIds": []
			  }
			]
		  }
	*/

	serverID := "bbad9cb1-7a1b-4c1d-b902-ead1f889c447"
	resp, err := psm.GetInitiatorsInformationForComputeNode(serverID)
	if err != nil {
		t.Errorf("Unexpected error in GetInitiatorsInformationForComputeNode %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestGetPathInfoForComputeNode
func xTestGetPathInfoForComputeNode(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	serverID := "53c96ca6-fdd4-4e03-ab1e-4725e8e4e6d9"
	hbaID := "667fb6e2-1cf7-448e-8e5d-4309756f7863"
	portID := "343db11a-21d0-45eb-a66c-c31ab1ac6863"

	reqBody := &vssbmodel.ComputeNodePathReq{
		HbaId:  hbaID,
		PortId: portID,
	}

	resp, err := psm.GetPathInfoForComputeNode(serverID, reqBody)
	if err != nil {
		t.Errorf("Unexpected error in GetPathInfoForComputeNode %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestGetPathsInfoForComputeNode
func xTestGetPathsInfoForComputeNode(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	serverID := "53c96ca6-fdd4-4e03-ab1e-4725e8e4e6d9"
	resp, err := psm.GetPathsInfoForComputeNode(serverID)
	if err != nil {
		t.Errorf("Unexpected error in GetPathsInfoForComputeNode %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestAddPathInfoToComputeNode
func xTestAddPathInfoToComputeNode(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	serverID := "53c96ca6-fdd4-4e03-ab1e-4725e8e4e6d9"
	hbaID := "667fb6e2-1cf7-448e-8e5d-4309756f7863"
	portID := "343db11a-21d0-45eb-a66c-c31ab1ac6863"

	reqBody := &vssbmodel.ComputeNodePathReq{
		HbaId:  hbaID,
		PortId: portID,
	}

	err = psm.AddPathInfoToComputeNode(serverID, reqBody)
	if err != nil {
		t.Errorf("Unexpected error in AddPathInfoToComputeNode %v", err)
		return
	}

	t.Logf("Successfully added path info to compute node %v", serverID)
}

// go test -v -run TestDeleteComputeNodePath
func xTestDeleteComputeNodePath(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	serverID := "53c96ca6-fdd4-4e03-ab1e-4725e8e4e6d9"
	hbaID := "667fb6e2-1cf7-448e-8e5d-4309756f7863"
	portID := "343db11a-21d0-45eb-a66c-c31ab1ac6863"

	reqBody := &vssbmodel.ComputeNodePathReq{
		HbaId:  hbaID,
		PortId: portID,
	}

	err = psm.DeleteComputeNodePath(serverID, reqBody)
	if err != nil {
		t.Errorf("Unexpected error in DeleteComputeNodePath %v", err)
		return
	}
	t.Logf("Successfully deleted path for compute node %v", serverID)
}

// go test -v -run TestGetConnectionInfoBtwnVolumeAndServerByVolumeID
func xTestGetConnectionInfoBtwnVolumeAndServerByVolumeID(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	volumeID := "8ec8062d-d1b5-414c-9557-f0a0279ca687"
	resp, err := psm.GetConnectionInfoBtwnVolumeAndServerByVolumeID(volumeID)
	if err != nil {
		t.Errorf("Unexpected error in GetConnectionInfoBtwnVolumeAndServerByVolumeID %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestGetConnectionInfoBtwnVolumeAndServerByServerID
func xTestGetConnectionInfoBtwnVolumeAndServerByServerID(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	serverID := "6673a1c4-7152-44ac-8145-deeee473aba5"
	resp, err := psm.GetConnectionInfoBtwnVolumeAndServerByServerID(serverID)
	if err != nil {
		t.Errorf("Unexpected error in GetConnectionInfoBtwnVolumeAndServerByServerID %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestGetConnectionInfoBtwnVolumeAndServerBoth
func xTestGetConnectionInfoBtwnVolumeAndServerBoth(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	volumeID := "8ec8062d-d1b5-414c-9557-f0a0279ca687"
	serverID := "6673a1c4-7152-44ac-8145-deeee473aba5"
	resp, err := psm.GetConnectionInfoBtwnVolumeAndServerBoth(volumeID, serverID)
	if err != nil {
		t.Errorf("Unexpected error in GetConnectionInfoBtwnVolumeAndServerBoth %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestSetPathBtwnVolumeAndServer
func xTestSetPathBtwnVolumeAndServer(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	volumeID := "8ec8062d-d1b5-414c-9557-f0a0279ca687"
	serverID := "6673a1c4-7152-44ac-8145-deeee473aba5"

	reqBody := &vssbmodel.SetPathVolumeServerReq{
		ServerId: serverID,
		VolumeId: volumeID,
	}

	err = psm.SetPathBtwnVolumeAndServer(reqBody)
	if err != nil {
		t.Errorf("Unexpected error in SetPathBtwnVolumeAndServer %v", err)
		return
	}

	t.Logf("Successfully set path between volume %s and server %s", volumeID, serverID)
}

// go test -v -run TestReleaseMultipleConnectionsBtwnVolumeAndServer
func xTestReleaseMultipleConnectionsBtwnVolumeAndServer(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	reqBody := &vssbmodel.ReleaseMultiConVolumeServerReq{
		ServerIds: []string{"6673a1c4-7152-44ac-8145-deeee473aba5"},
		VolumeIds: []string{"8ec8062d-d1b5-414c-9557-f0a0279ca687"},
	}

	err = psm.ReleaseMultipleConnectionsBtwnVolumeAndServer(reqBody)
	if err != nil {
		t.Errorf("Unexpected error in ReleaseMultipleConnectionsBtwnVolumeAndServer %v", err)
		return
	}

	t.Logf("Successfully released multiple connections between volume and server")
}

// go test -v -run TestReleaseConnectionBtwnVolumeAndServer
func xTestReleaseConnectionBtwnVolumeAndServer(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	volumeID := "8ec8062d-d1b5-414c-9557-f0a0279ca687"
	serverID := "6673a1c4-7152-44ac-8145-deeee473aba5"

	reqBody := &vssbmodel.SetPathVolumeServerReq{
		ServerId: serverID,
		VolumeId: volumeID,
	}

	err = psm.ReleaseConnectionBtwnVolumeAndServer(reqBody)
	if err != nil {
		t.Errorf("Unexpected error in ReleaseConnectionBtwnVolumeAndServer %v", err)
		return
	}

	t.Logf("Successfully set path between volume %s and server %s", volumeID, serverID)
}
