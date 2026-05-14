package sanstorage

import (
	"encoding/json"
	utils "terraform-provider-hitachi/hitachi/common/utils"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	"testing"
)

// go test -v -run ^TestGetNvmSubsystem$
func ^TestGetNvmSubsystem(t *testing.T) {
	/*
	   Expected Response Format:
	   {
	     "nvmSubsystemId" : 1,
	     "nvmSubsystemName" : "rest_subsystem",
	     "resourceGroupId" : 0,
	     "namespaceSecuritySetting" : "Enable",
	     "t10piMode" : "Disable",
	     "hostMode" : "LINUX/IRIX",
	     "nvmSubsystemNqn" : "nqn.2015-04.com.example:nvme:storage-subsystem-sn.5-10088-nvmssid.00001",
	     "namespaces" : [
	           {
	             "namespaceId" : 1,
	             "ldevId" : 2000
	           }
	       ],
	     "portIds" : [ "CL1-A", "CL1-B" ]
	   }
	*/
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating test manager: %v", err)
	}

	// Replace with a valid NVM Subsystem ID from your environment
	subsystemID := 1
	resp, err := psm.GetNvmSubsystem(subsystemID)
	if err != nil {
		t.Errorf("Unexpected error in GetNvmSubsystem: %v", err)
		return
	}

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Logf("Failed to marshal Response: %v", err)
		return
	}

	t.Logf("NVM Subsystem %d Response:\n%s", subsystemID, string(data))
}

// go test -v -run ^TestGetAllNvmSubsystems$
func ^TestGetAllNvmSubsystems(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	infoType := "basic"
	params := sanmodel.GetNvmSubsystemsParams{
		NvmSubsystemInfo: &infoType,
	}

	resp, err := psm.GetAllNvmSubsystems(params)
	if err != nil {
		t.Errorf("Unexpected error in GetAllNvmSubsystems %v", err)
		return
	}

	data, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("Response:\n%s", string(data))
}

// go test -v -run ^TestNvmSubsystemQueries$
func ^TestNvmSubsystemQueries(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	tests := []struct {
		name   string
		params sanmodel.GetNvmSubsystemsParams
	}{
		{
			name: "GetBasicInfo", // Explicitly requesting 'basic'
			params: sanmodel.GetNvmSubsystemsParams{
				NvmSubsystemInfo: utils.Ptr("basic"),
			},
		},
		{
			name: "GetNQNs",
			params: sanmodel.GetNvmSubsystemsParams{
				NvmSubsystemInfo: utils.Ptr("nqn"),
			},
		},
		{
			name: "GetNamespaces",
			params: sanmodel.GetNvmSubsystemsParams{
				NvmSubsystemInfo: utils.Ptr("namespace"),
			},
		},
		{
			name: "GetPorts",
			params: sanmodel.GetNvmSubsystemsParams{
				NvmSubsystemInfo: utils.Ptr("port"),
			},
		},
		// {
		//     name: "GetUnimplementedIDs",
		//     params: sanmodel.GetNvmSubsystemsParams{
		//         NvmSubsystemOption: utils.Ptr("undefined"),
		//     },
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := psm.GetAllNvmSubsystems(tt.params)
			if err != nil {
				t.Errorf("%s: API call failed: %v", tt.name, err)
				return
			}

			if resp == nil {
				t.Errorf("%s: Received nil response", tt.name)
				return
			}

			data, _ := json.MarshalIndent(resp, "", "  ")
			t.Logf("%s Response (Count: %d):\n%s", tt.name, len(resp.Data), string(data))
		})
	}
}

// go test -v -run ^TestCreateNvmSubsystem$
func ^TestCreateNvmSubsystem(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	targetId := 1

	req := sanmodel.CreateNvmSubsystemRequest{
		NvmSubsystemId:           targetId,
		NvmSubsystemName:         utils.Ptr("TF_TEST_SUBSYSTEM"),
		HostMode:                 utils.Ptr("LINUX/IRIX"),
		NamespaceSecuritySetting: utils.Ptr("Enable"),
		// HostModeOptions:          &[]int{68},
	}

	resIds, err := psm.CreateNvmSubsystem(req)
	if err != nil {
		t.Errorf("Failed to create NVM Subsystem: %v", err)
		return
	}

	t.Logf("Create Job Submitted Successfully: %v", resIds)
}

// go test -v -run ^TestUpdateNvmSubsystem$
func ^TestUpdateNvmSubsystem(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	targetId := 1
	req := sanmodel.UpdateNvmSubsystemRequest{
		NvmSubsystemName: utils.Ptr("TF_UPDATED_NAME"),
	}

	resIds, err := psm.UpdateNvmSubsystem(targetId, req)
	if err != nil {
		t.Errorf("Failed to update NVM Subsystem: %v", err)
		return
	}

	t.Logf("Update successful. Affected Resource ID: %v", resIds)
}

// go test -v -run ^TestGetNvmSubsystemPort$
func ^TestGetNvmSubsystemPort(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	subsystemID := 1
	portID := "CL1-A"
	resp, err := psm.GetNvmSubsystemPort(subsystemID, portID)
	if err != nil {
		t.Errorf("Unexpected error in GetNvmSubsystemPort: %v", err)
		return
	}

	data, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("Port Detail for Subsystem %d, Port %s:\n%s", subsystemID, portID, string(data))
}

// go test -v -run ^TestGetNvmSubsystemPorts$
func ^TestGetNvmSubsystemPorts(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	subsystemID := 1
	resp, err := psm.GetNvmSubsystemPorts(subsystemID)
	if err != nil {
		t.Errorf("Unexpected error in GetNvmSubsystemPorts: %v", err)
		return
	}

	data, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("Subsystem %d Detailed Port Info:\n%s", subsystemID, string(data))
}

// go test -v -run ^TestAddNvmSubsystemPort$
func ^TestAddNvmSubsystemPort(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	req := sanmodel.AddNvmSubsystemPortRequest{
		NvmSubsystemId: 1,
		PortId:         "CL1-A",
	}

	jobId, err := psm.AddNvmSubsystemPort(req)
	if err != nil {
		t.Errorf("Failed to add port: %v", err)
		return
	}

	t.Logf("Add Port Job Submitted: %v", jobId)
}

// go test -v -run ^TestDeleteNvmSubsystemPort$
func ^TestDeleteNvmSubsystemPort(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	subsystemID := 1
	portID := "CL1-A"

	jobId, err := psm.DeleteNvmSubsystemPort(subsystemID, portID)
	if err != nil {
		t.Errorf("Failed to delete port: %v", err)
		return
	}

	t.Logf("Delete Port Job Submitted: %v", jobId)
}

// go test -v -run ^TestGetAllHostNqns$
func ^TestGetAllHostNqns(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	subsystemID := 1
	resp, err := psm.GetAllHostNqns(subsystemID)
	if err != nil {
		t.Errorf("Failed to get host NQNs: %v", err)
		return
	}

	data, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("Host NQNs for Subsystem %d:\n%s", subsystemID, string(data))
}

// go test -v -run ^TestRegisterHostNqn$
func ^TestRegisterHostNqn(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	req := sanmodel.RegisterHostNqnRequest{
		NvmSubsystemId: 1,
		HostNqn:        "nqn.2014-08.org.example:uuid:4b73e622-ddc1-449a-99f7-412c0d3baa39",
	}

	jobId, err := psm.RegisterHostNqn(req)
	if err != nil {
		t.Errorf("Failed to register host NQN: %v", err)
		return
	}

	t.Logf("Register NQN Job Submitted: %v", jobId)
}

// go test -v -run ^TestSetHostNqnNickname$
func ^TestSetHostNqnNickname(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	subsystemID := 1
	hostNqn := "nqn.2014-08.org.example:uuid:4b73e622-ddc1-449a-99f7-412c0d3baa39"
	req := sanmodel.SetHostNqnNicknameRequest{
		HostNqnNickname: "TF_HOST_1",
	}

	jobId, err := psm.SetHostNqnNickname(subsystemID, hostNqn, req)
	if err != nil {
		t.Errorf("Failed to set nickname: %v", err)
		return
	}

	t.Logf("Set Nickname Job Submitted: %v", jobId)
}

// go test -v -run ^TestDeleteLoginHostNqn$
func ^TestDeleteLoginHostNqn(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	portID := "CL1-A"
	jobId, err := psm.DeleteLoginHostNqn(portID)
	if err != nil {
		t.Errorf("Failed to delete login info: %v", err)
		return
	}

	t.Logf("Delete Login Info Job Submitted: %v", jobId)
}

// go test -v -run ^TestGetAllNamespaces$
func ^TestGetAllNamespaces(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	subsystemID := 1
	resp, err := psm.GetAllNamespaces(subsystemID)
	if err != nil {
		t.Errorf("Failed to get namespaces: %v", err)
		return
	}

	data, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("Namespaces for Subsystem %d:\n%s", subsystemID, string(data))
}

// go test -v -run ^TestGetNamespace$
func ^TestGetNamespace(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	subsystemID := 1
	namespaceID := 1
	resp, err := psm.GetNamespace(subsystemID, namespaceID)
	if err != nil {
		t.Errorf("Unexpected error in GetNamespace: %v", err)
		return
	}

	data, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("Namespace %d Detail:\n%s", namespaceID, string(data))
}

// go test -v -run ^TestCreateNamespace$
func ^TestCreateNamespace(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	req := sanmodel.CreateNamespaceRequest{
		NvmSubsystemId: 1,
		LdevId:         2000,
		NamespaceId:    utils.Ptr(1),
	}

	jobId, err := psm.CreateNamespace(req)
	if err != nil {
		t.Errorf("Failed to create namespace: %v", err)
		return
	}

	t.Logf("Create Namespace Job Submitted: %v", jobId)
}

// go test -v -run ^TestSetNamespaceNickname$
func ^TestSetNamespaceNickname(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	subsystemID := 1
	namespaceID := 1
	req := sanmodel.SetNamespaceNicknameRequest{
		NamespaceNickname: "TF_NS_NICKNAME",
	}

	resIds, err := psm.SetNamespaceNickname(subsystemID, namespaceID, req)
	if err != nil {
		t.Errorf("Failed to set nickname: %v", err)
		return
	}

	t.Logf("Nickname Update Job Submitted: %v", resIds)
}

// go test -v -run ^TestDeleteNamespace$
func ^TestDeleteNamespace(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	subsystemID := 1
	namespaceID := 1

	resIds, err := psm.DeleteNamespace(subsystemID, namespaceID)
	if err != nil {
		t.Errorf("Failed to delete namespace: %v", err)
		return
	}

	t.Logf("Namespace Delete Job Submitted: %v", resIds)
}

// go test -v -run ^TestGetNamespacePaths$
func ^TestGetNamespacePaths(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	params := sanmodel.GetNamespacePathsParams{
		NvmSubsystemId: 1,
		NamespaceId:    utils.Ptr(1),
	}

	resp, err := psm.GetNamespacePaths(params)
	if err != nil {
		t.Errorf("Unexpected error in GetNamespacePaths: %v", err)
		return
	}

	data, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("Namespace Paths Response:\n%s", string(data))
}

// go test -v -run ^TestGetNamespacePathDetail$
func ^TestGetNamespacePathDetail(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	subsystemId := 1
	hostNqn := "nqn.2014-08.org.example:uuid:ff533865-1d0b-4043-8345-a90afbb80d8b"
	namespaceId := 1

	resp, err := psm.GetNamespacePathDetail(subsystemId, hostNqn, namespaceId)
	if err != nil {
		t.Errorf("Unexpected error in GetNamespacePathDetail: %v", err)
		return
	}

	data, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("Namespace Path Detail:\n%s", string(data))
}

// go test -v -run ^TestRegisterNamespacePath$
func ^TestRegisterNamespacePath(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	req := sanmodel.RegisterNamespacePathRequest{
		NvmSubsystemId: 1,
		HostNqn:        "nqn.2014-08.org.example:uuid:4b73e622-ddc1-449a-99f7-412c0d3baa39",
		NamespaceId:    1,
	}

	resIds, err := psm.RegisterNamespacePath(req)
	if err != nil {
		t.Errorf("Failed to register namespace path: %v", err)
		return
	}

	t.Logf("Namespace Path Registration Job Submitted: %v", resIds)
}

// go test -v -run ^TestDeleteNamespacePath$
func ^TestDeleteNamespacePath(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	subsystemId := 1
	hostNqn := "nqn.2014-08.org.example:uuid:4b73e622-ddc1-449a-99f7-412c0d3baa39"
	namespaceId := 1

	resIds, err := psm.DeleteNamespacePath(subsystemId, hostNqn, namespaceId)
	if err != nil {
		t.Errorf("Failed to delete namespace path: %v", err)
		return
	}

	t.Logf("Namespace Path Deletion Job Submitted: %v", resIds)
}
