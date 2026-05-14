package sanstorage

import (
	"fmt"
	"terraform-provider-hitachi/hitachi/common/utils"
	sangatewaymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"
	"testing"
)

// newNvmSubsystemTestManager is for Testing and provides structure information for connection
func newNvmSubsystemTestManager() (*sanStorageManager, error) {
	objStorage := sanmodel.StorageDeviceSettings{
		Serial:   12345,
		Username: "user1",
		Password: "mypswd",
		MgmtIP:   "10.10.11.12",
	}
	psm, err := newSanStorageManagerEx(objStorage)
	if err != nil {
		return nil, fmt.Errorf("unexpected error while creating newSanStorageManagerEx %v", err)
	}
	return psm, nil
}

// go test -v -run ^TestGetAllNvmSubsystems$
func ^TestGetAllNvmSubsystems(t *testing.T) {
	psm, err := newNvmSubsystemTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	params := sangatewaymodel.GetNvmSubsystemsParams{
		NvmSubsystemInfo: utils.Ptr("basic"),
	}
	resp, err := psm.GetAllNvmSubsystems(params)
	if err != nil {
		t.Errorf("Unexpected error in GetAllNvmSubsystems: %v", err)
		return
	}
	t.Logf("Response: %+v", resp)
}

// go test -v -run ^TestGetNvmSubsystem$
func ^TestGetNvmSubsystem(t *testing.T) {
	psm, err := newNvmSubsystemTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetNvmSubsystem(1)
	if err != nil {
		t.Errorf("Unexpected error in GetNvmSubsystem: %v", err)
		return
	}
	t.Logf("Response: %+v", resp)
}

// go test -v -run ^TestCreateNvmSubsystem$
func ^TestCreateNvmSubsystem(t *testing.T) {
	psm, err := newNvmSubsystemTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	request := sangatewaymodel.CreateNvmSubsystemRequest{
		NvmSubsystemId:           1,
		NvmSubsystemName:         utils.Ptr("TF_NVME_TEST"),
		HostMode:                 utils.Ptr("LINUX/IRIX"),
		NamespaceSecuritySetting: utils.Ptr("Enable"),
	}
	resId, err := psm.CreateNvmSubsystem(request)
	if err != nil {
		t.Errorf("Unexpected error in CreateNvmSubsystem: %v", err)
		return
	}
	t.Logf("ResId: %v", resId)
}

// go test -v -run ^TestUpdateNvmSubsystem$
func ^TestUpdateNvmSubsystem(t *testing.T) {
	psm, err := newNvmSubsystemTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Note: API allows only one attribute change per call
	request := sangatewaymodel.UpdateNvmSubsystemRequest{
		NvmSubsystemName: utils.Ptr("TF_NVME_TEST_UPDATED"),
	}
	resId, err := psm.UpdateNvmSubsystem(1, request)
	if err != nil {
		t.Errorf("Unexpected error in UpdateNvmSubsystem: %v", err)
		return
	}
	t.Logf("ResId: %v", resId)
}

// go test -v -run ^TestDeleteNvmSubsystem$
func ^TestDeleteNvmSubsystem(t *testing.T) {
	psm, err := newNvmSubsystemTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resId, err := psm.DeleteNvmSubsystem(1)
	if err != nil {
		t.Errorf("Unexpected error in DeleteNvmSubsystem: %v", err)
		return
	}
	t.Logf("ResId: %v", resId)
}

/* --- Port Management Tests --- */

// go test -v -run ^TestGetNvmSubsystemPorts$
func ^TestGetNvmSubsystemPorts(t *testing.T) {
	psm, err := newNvmSubsystemTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetNvmSubsystemPorts(1)
	if err != nil {
		t.Errorf("Unexpected error in GetNvmSubsystemPorts: %v", err)
		return
	}
	t.Logf("Response: %+v", resp)
}

// go test -v -run ^TestAddNvmSubsystemPort$
func ^TestAddNvmSubsystemPort(t *testing.T) {
	psm, err := newNvmSubsystemTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	request := sangatewaymodel.AddNvmSubsystemPortRequest{
		NvmSubsystemId: 1,
		PortId:         "CL1-A",
	}
	resId, err := psm.AddNvmSubsystemPort(request)
	if err != nil {
		t.Errorf("Unexpected error in AddNvmSubsystemPort: %v", err)
		return
	}
	t.Logf("ResId: %v", resId)
}

/* --- Host NQN Management Tests --- */

// go test -v -run ^TestGetAllHostNqns$
func ^TestGetAllHostNqns(t *testing.T) {
	psm, err := newNvmSubsystemTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetAllHostNqns(1)
	if err != nil {
		t.Errorf("Unexpected error in GetAllHostNqns: %v", err)
		return
	}
	t.Logf("Response: %+v", resp)
}

// go test -v -run ^TestRegisterHostNqn$
func ^TestRegisterHostNqn(t *testing.T) {
	psm, err := newNvmSubsystemTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	request := sangatewaymodel.RegisterHostNqnRequest{
		NvmSubsystemId:  1,
		HostNqn:         "nqn.2014-08.org.example:uuid:6ed48464-670d",
		HostNqnNickname: utils.Ptr("TF_HOST_01"),
	}
	resId, err := psm.RegisterHostNqn(request)
	if err != nil {
		t.Errorf("Unexpected error in RegisterHostNqn: %v", err)
		return
	}
	t.Logf("ResId: %v", resId)
}

// go test -v -run ^TestDeleteLoginHostNqn$
func ^TestDeleteLoginHostNqn(t *testing.T) {
	psm, err := newNvmSubsystemTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resId, err := psm.DeleteLoginHostNqn("CL1-A")
	if err != nil {
		t.Errorf("Unexpected error in DeleteLoginHostNqn: %v", err)
		return
	}
	t.Logf("ResId: %v", resId)
}

/* --- Namespace Management Tests --- */

// go test -v -run ^TestGetAllNamespaces$
func ^TestGetAllNamespaces(t *testing.T) {
	psm, err := newNvmSubsystemTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetAllNamespaces(1)
	if err != nil {
		t.Errorf("Unexpected error in GetAllNamespaces: %v", err)
		return
	}
	t.Logf("Response: %+v", resp)
}

// go test -v -run ^TestCreateNamespace$
func ^TestCreateNamespace(t *testing.T) {
	psm, err := newNvmSubsystemTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	request := sangatewaymodel.CreateNamespaceRequest{
		NvmSubsystemId:    1,
		LdevId:            100,
		NamespaceNickname: utils.Ptr("TF_NS_100"),
	}
	resId, err := psm.CreateNamespace(request)
	if err != nil {
		t.Errorf("Unexpected error in CreateNamespace: %v", err)
		return
	}
	t.Logf("ResId: %v", resId)
}

/* --- Namespace Path Tests --- */

// go test -v -run ^TestGetNamespacePaths$
func ^TestGetNamespacePaths(t *testing.T) {
	psm, err := newNvmSubsystemTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	params := sangatewaymodel.GetNamespacePathsParams{
		NvmSubsystemId: 1,
		NamespaceId:    utils.Ptr(1),
	}
	resp, err := psm.GetNamespacePaths(params)
	if err != nil {
		t.Errorf("Unexpected error in GetNamespacePaths: %v", err)
		return
	}
	t.Logf("Response: %+v", resp)
}

// go test -v -run ^TestRegisterNamespacePath$
func ^TestRegisterNamespacePath(t *testing.T) {
	psm, err := newNvmSubsystemTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	request := sangatewaymodel.RegisterNamespacePathRequest{
		NvmSubsystemId: 1,
		NamespaceId:    1,
		HostNqn:        "nqn.2014-08.org.example:uuid:6ed48464-670d",
	}
	resId, err := psm.RegisterNamespacePath(request)
	if err != nil {
		t.Errorf("Unexpected error in RegisterNamespacePath: %v", err)
		return
	}
	t.Logf("ResId: %v", resId)
}

// go test -v -run ^TestDeleteNamespacePath$
func ^TestDeleteNamespacePath(t *testing.T) {
	psm, err := newNvmSubsystemTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resId, err := psm.DeleteNamespacePath(1, "nqn.2014-08.org.example:uuid:6ed48464-670d", 1)
	if err != nil {
		t.Errorf("Unexpected error in DeleteNamespacePath: %v", err)
		return
	}
	t.Logf("ResId: %v", resId)
}
