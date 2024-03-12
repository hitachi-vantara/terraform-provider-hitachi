package infra_gw

import (
	"fmt"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
	"testing"
)

// newDynamicPoolTestManager is for Testing and provide structure information for connection
func newReconcilerestManager() (*infraGwManager, error) {

	// Following storage has iscsi port
	// subscrierId := "ucpadmin"
	// partnerId := "ucpadmin"

	setting := model.InfraGwSettings{
		Username: "ucpadmin",
		Password: "Passw0rd!",
		Address:  "172.25.22.81",
		V3API:    false,
		// PartnerId:    &partnerId,
		// SubscriberId: &subscrierId,
	}

	psm, err := newInfraGwManagerEx(setting)
	if err != nil {
		return nil, fmt.Errorf("unexpected error while creating newDynamicPoolTestManager %v", err)
	}
	return psm, nil
}

func TestCreateUpdateVolume(t *testing.T) {
	psm, err := newReconcilerestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	pooldId := 0
	// lunid := 636
	// {"name":"VolumeTest1111121","poolId":4,"parityGroupId":"1-3","capacity":"1GB","ucpSystem":"UCP-SYS1"}
	storageId := "storage-349a72cc2d6b6b131ac5f2c4d557c6d6"
	createInput := model.CreateVolumeParams{Capacity: "100MB",
		System: "Logical-UCP-95054", PoolID: &pooldId}
	sid, err := psm.ReconcileVolume(storageId, &createInput, nil)
	if err != nil {
		t.Errorf("Unexpected error in GetPartnerIdWithStatus %v", err)
		return
	}
	t.Logf("Response: %v", sid)
}

func xTestGetVolumeByLdevId(t *testing.T) {
	psm, err := newReconcilerestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// {"name":"VolumeTest1111121","poolId":4,"parityGroupId":"1-3","capacity":"1GB","ucpSystem":"UCP-SYS1"}
	storageId := "storage-39f4eef0175c754bb90417358b0133c3"

	pvol, mvol, err := psm.GetVolumeByLDevId(storageId, 562)
	if err != nil {
		t.Errorf("Unexpected error in GetPartnerIdWithStatus %v", err)
		return
	}
	t.Logf("Response: %v %v", pvol, mvol)
}

func xTestGetVolumeByID(t *testing.T) {
	psm, err := newReconcilerestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// {"name":"VolumeTest1111121","poolId":4,"parityGroupId":"1-3","capacity":"1GB","ucpSystem":"UCP-SYS1"}
	storageId := "storage-349a72cc2d6b6b131ac5f2c4d557c6d6"
	volId := "storagevolume-91a060e705ade94cedf4c149c20d7d6a"
	pvol, err := psm.GetVolumeByID(storageId, volId)
	if err != nil {
		t.Errorf("Unexpected error in GetPartnerIdWithStatus %v", err)
		return
	}
	t.Logf("Response: %v", pvol)
}
