package infra_gw

import (
	"fmt"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
	"testing"
)

// newDynamicPoolTestManager is for Testing and provide structure information for connection
func newReconcilerestManagerForPartner() (*infraGwManager, error) {

	// Following storage has iscsi port
	// subscrierId := "46519299-c43c-4c6e-a680-81dce45a3fca"
	// partnerId := "a8d1f065-a9e7-42cf-b565-a67466fec549"

	setting := model.InfraGwSettings{
		Username: "ucpadmin",
		Password: "Passw0rd!",
		Address:  "172.25.22.81",
		V3API:    false,
		// 	PartnerId:    &partnerId,
		// 	SubscriberId: &subscrierId,
	}

	psm, err := newInfraGwManagerEx(setting)
	if err != nil {
		return nil, fmt.Errorf("unexpected error while creating newDynamicPoolTestManager %v", err)
	}
	return psm, nil
}

func TestGetPartnerIDandSubscriberID(t *testing.T) {
	psm, err := newReconcilerestManagerForPartner()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	sid, err := psm.GetPartnerAndSubscriberId("ucpadmin")
	if err != nil {
		t.Errorf("Unexpected error in GetPartnerIdWithStatus %v", err)
		return
	}
	t.Logf("Response: %v", sid)
}
