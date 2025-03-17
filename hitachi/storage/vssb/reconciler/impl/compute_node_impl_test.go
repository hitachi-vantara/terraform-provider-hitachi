package vssbstorage

import (
	"fmt"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vssb/reconciler/model"
	"testing"
)

func newTestManager() (*vssbStorageManager, error) {

	objStorage := vssbmodel.StorageDeviceSettings{
		Username:       "admin",
		Password:       "vssb-789",
		ClusterAddress: "10.76.47.55",
		// ClusterAddress: "172.25.58.151",
	}
	psm, err := newVssbStorageManagerEx(objStorage)
	if err != nil {
		return nil, fmt.Errorf("unexpected error while creating newVssbStorageManagerEx %v", err)
	}
	return psm, nil
}

// go test -v -run TestGetComputeNodes
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

// go test -v -run TestReconcileComputeNode
func xTestReconcileComputeNode(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	input := vssbmodel.ComputeResource{
		Name:   "ComputeNode-RESTAPI123",
		OsType: "VMware",

		// IscsiConnections: []vssbmodel.IscsiConnector{
		// 	{IscsiInitiator: "iqn.1998-01.com.vmware:node-06-0723aa94", PortNames: []string{"002-iSCSI-001"}},
		// 	//{IscsiInitiator: "iqn.1998-01.com.vmware:node-06-0723aa94"},
		// 	//{IscsiInitiator: "iqn.1998-01.com.vmware:node-06-0723aa94", PortNames: []string{"002-iSCSI-001", "001-iSCSI-002"}},
		// 	//{IscsiInitiator: "iqn.1998-01.com.vmware:node-06-0723aa95", PortNames: []string{"002-iSCSI-001", "001-iSCSI-002"}},
		// 	//{IscsiInitiator: "iqn.1998-01.com.vmware:node-06-0723aa95", PortNames: []string{"002-iSCSI-001", "001-iSCSI-002xx"}},
		// },
		FcConnections: []vssbmodel.FcConnector{
			{HostWWN: "56360e8107595362"},
		},
	}
	resp, err := psm.ReconcileComputeNode(&input)
	if err != nil {
		t.Errorf("Unexpected error in ReconcileComputeNode %v", err)
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

// go test -v -run TestReconcileUpdateComputeNode
func xTestReconcileUpdateComputeNode(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// #1 Update Name
	/*
		input := vssbmodel.ComputeResource{
			Name: "ComputeNode-RESTAPI",
			//Name: "ComputeNode-RESTAPI123",
			ID: "47b2bcdf-9fc7-434f-9bcf-cd23fd3791d1",
		}
	*/

	// #2 Update Name and OS Type
	/*
		input := vssbmodel.ComputeResource{
			Name:   "ComputeNode-RESTAPI",
			OsType: "VMware",
			//OsType: "Windows",
			ID: "47b2bcdf-9fc7-434f-9bcf-cd23fd3791d1",
		}
	*/

	// #3 Update Only IscsiInitiator
	/*
		input := vssbmodel.ComputeResource{
			Name:   "ComputeNode-RESTAPI",
			OsType: "VMware",
			ID:     "47b2bcdf-9fc7-434f-9bcf-cd23fd3791d1",

			IscsiConnections: []vssbmodel.IscsiConnector{
				//{IscsiInitiator: "iqn.1998-01.com.vmware:node-06-0723aa95"},
				{IscsiInitiator: "iqn.1998-01.com.vmware:node-06-0723aa94"},
			},
		}
	*/
	// #4 Update Iscsi Initiator + PortNames

	input := vssbmodel.ComputeResource{
		Name:   "ComputeNode-RESTAPI",
		OsType: "VMware",
		ID:     "47b2bcdf-9fc7-434f-9bcf-cd23fd3791d1",

		IscsiConnections: []vssbmodel.IscsiConnector{
			//{IscsiInitiator: "iqn.1998-01.com.vmware:node-06-0723aa94"},
			//{IscsiInitiator: "iqn.1998-01.com.vmware:node-06-0723aa94", PortNames: []string{"001-iSCSI-002"}},
			{IscsiInitiator: "iqn.1998-01.com.vmware:node-06-0723aa94", PortNames: []string{"002-iSCSI-001", "001-iSCSI-002"}},
			//{IscsiInitiator: "iqn.1998-01.com.vmware:node-06-0723aa95", PortNames: []string{"002-iSCSI-001", "001-iSCSI-002"}},
			//{IscsiInitiator: "iqn.1998-01.com.vmware:node-06-0723aa95", PortNames: []string{"002-iSCSI-001", "001-iSCSI-002xx"}},
		},
	}

	resp, err := psm.ReconcileComputeNode(&input)
	if err != nil {
		t.Errorf("Unexpected error in ReconcileComputeNode %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}
