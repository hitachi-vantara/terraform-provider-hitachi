package infra_gw

import (
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
	"testing"
)

func TestOnbaordStorageDevice(t *testing.T) {
	psm, err := newMTTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	reqBody := model.CreateStorageDeviceParam{
		SerialNumber:      "40014",
		ManagementAddress: "172.25.47.115",
		// GatewayAddress:    "172.25.20.35",
		Username:  "ms_vmware",
		Password:  "Hitachi1",
		UcpSystem: "UCP-SYS1",
	}

	sid, err := psm.AddStorageDevice(reqBody)
	if err != nil {
		t.Errorf("Unexpected error in GetPartnerIdWithStatus %v", err)
		return
	}
	t.Logf("Response:  %v", sid)
}
