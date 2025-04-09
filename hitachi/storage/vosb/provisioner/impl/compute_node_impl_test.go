package vssbstorage

import (
	"fmt"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/provisioner/model"
	"testing"
)

func newTestManager() (*vssbStorageManager, error) {

	objStorage := vssbmodel.StorageDeviceSettings{
		Username:       "admin",    // admin
		Password:       "vssb-789", // vssb-789
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

	serverID := "697659d5-5856-4f94-bdf3-bbe47f3487bd"
	resp, err := psm.GetComputeNode(serverID)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestGetAllComputeNode
func xTestGetAllComputeNode(t *testing.T) {
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

// go test -v -run TestCreateComputeResource
func xTestCreateComputeResource(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	input := vssbmodel.ComputeResource{
		Name:   "ComputeNode-RESTAPI",
		OsType: "VMware",

		IscsiConnections: []vssbmodel.IscsiConnector{
			{IscsiInitiator: "iqn.1998-01.com.vmware:node-06-0723aa94", PortNames: []string{"002-iSCSI-001"}},
			//{IscsiInitiator: "iqn.1998-01.com.vmware:node-06-0723aa94"},
			//{IscsiInitiator: "iqn.1998-01.com.vmware:node-06-0723aa94", PortNames: []string{"002-iSCSI-001", "001-iSCSI-002"}},
			//{IscsiInitiator: "iqn.1998-01.com.vmware:node-06-0723aa95", PortNames: []string{"002-iSCSI-001", "001-iSCSI-002"}},
			//{IscsiInitiator: "iqn.1998-01.com.vmware:node-06-0723aa95", PortNames: []string{"002-iSCSI-001", "001-iSCSI-002xx"}},
		},

		FcConnections: []vssbmodel.FcConnector{
			{HostWWN: "50060E8107595000"},
		},
	}
	err = psm.CreateComputeResource(&input)
	if err != nil {
		t.Errorf("Unexpected error in CreateComputeResource %v", err)
		return
	}
	t.Logf("Successfully created compute reseource")
}

// go test -v -run TestGetComputeResourceInfo
func xTestGetComputeResourceInfo(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	//computeNodeName := "esxi-151"
	computeNodeName := "ComputeNode-RESTAPI2"
	resp, err := psm.GetComputeResourceInfo(computeNodeName)
	if err != nil {
		t.Errorf("Unexpected error in GetComputeResourceInfo %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestDeleteComputeNodeResource
func xTestDeleteComputeNodeResource(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	serverId := "9a6e6219-4397-443a-8427-bef1e7e19396"

	err = psm.DeleteComputeNodeResource(serverId)
	if err != nil {
		t.Errorf("Unexpected error in DeleteComputeNodeResource %v", err)
		return
	}
}
