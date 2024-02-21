package infra_gw

import (
	"fmt"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
	"testing"
)

// newDynamicPoolTestManager is for Testing and provide structure information for connection
func newReconcilerestManager() (*infraGwManager, error) {

	// Following storage has iscsi port
	subscrierId := "46519299-c43c-4c6e-a680-81dce45a3fca"
	partnerId := "a8d1f065-a9e7-42cf-b565-a67466fec549"

	setting := model.InfraGwSettings{
		Username:     "ucpadmin",
		Password:     "Passw0rd!",
		Address:      "172.25.22.81",
		V3API:        false,
		PartnerId:    &partnerId,
		SubscriberId: &subscrierId,
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
	// {"name":"VolumeTest1111121","poolId":4,"parityGroupId":"1-3","capacity":"1GB","ucpSystem":"UCP-SYS1"}
	storageId := "storage-39f4eef0175c754bb90417358b0133c3"
	createInput := model.CreateVolumeParams{Name: "testVol222Name", Capacity: "100MB", DeduplicationCompressionMode: "ENABLED",
		System: "Logical-UCP-30595", PoolID: &pooldId, ParityGroupId: "1-1"}
	sid, err := psm.ReconcileVolume(storageId, &createInput, nil)
	if err != nil {
		t.Errorf("Unexpected error in GetPartnerIdWithStatus %v", err)
		return
	}
	t.Logf("Response: %v", sid)
}

func xTestGetPartnetSubscribervolumeVolume(t *testing.T) {
	psm, err := newReconcilerestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// {"name":"VolumeTest1111121","poolId":4,"parityGroupId":"1-3","capacity":"1GB","ucpSystem":"UCP-SYS1"}
	storageId := "storage-39f4eef0175c754bb90417358b0133c3"

	sid, err := psm.GetVolumesByPartnerSubscriberID(storageId, 0, 10)
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
