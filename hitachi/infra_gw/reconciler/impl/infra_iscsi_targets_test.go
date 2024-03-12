package infra_gw

import (
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
	"testing"
)

func TestCreateUpdateISCSITarget(t *testing.T) {
	psm, err := newReconcilerestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	// lunid := 636
	// {"name":"VolumeTest1111121","poolId":4,"parityGroupId":"1-3","capacity":"1GB","ucpSystem":"UCP-SYS1"}
	storageId := "storage-349a72cc2d6b6b131ac5f2c4d557c6d6"
	createInput := model.CreateIscsiTargetParam{
		IscsiName: "NewTargetName1",
		Port:      "CL2-C",
		HostMode:  "VMWARE",
		UcpSystem: "Logical-UCP-95054",
	}
	sid, err := psm.ReconcileIscsiTarget(storageId, &createInput)
	if err != nil {
		t.Errorf("Unexpected error in GetPartnerIdWithStatus %v", err)
		return
	}
	t.Logf("Response: %v", sid)
}
