package infra_gw

import (
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
	"testing"
)

func xTestOnboardDeviceImpl(t *testing.T) {
	
	psm, err := newReconcilerestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// {"name":"VolumeTest1111121","poolId":4,"parityGroupId":"1-3","capacity":"1GB","ucpSystem":"UCP-SYS1"}
	// storageId := "storage-39f4eef0175c754bb90417358b0133c3"
	reqBody := model.CreateStorageDeviceParam{
		SerialNumber:      "30595",
		ManagementAddress: "172.25.47.112",
		GatewayAddress:    "172.25.20.35",
		Username:          "ms_vmware",
		Password:          "Hitachi1",
		// UcpSystem: "UCP-SYS12",
	}
	sid, err := psm.ReconcileStorageDevice("", &reqBody)
	if err != nil {
		t.Errorf("Unexpected error in GetPartnerIdWithStatus %v", err)
		return
	}
	t.Logf("Response: %v", sid)
}
