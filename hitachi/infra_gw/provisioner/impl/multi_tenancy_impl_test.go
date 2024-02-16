package infra_gw

import (
	"fmt"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
	"testing"
)

// newDynamicPoolTestManager is for Testing and provide structure information for connection
func newMTTestManager() (*infraGwManager, error) {

	// Following storage has iscsi port

	setting := model.InfraGwSettings{
		Username: "ucpadmin",
		Password: "Passw0rd!",
		Address:  "172.25.22.81",
	}

	psm, err := newInfraGwManagerEx(setting)
	if err != nil {
		return nil, fmt.Errorf("unexpected error while creating newDynamicPoolTestManager %v", err)
	}
	return psm, nil
}

// go test -v -run TestGetPartnerIDwithStatus
func TestGetPartnerIDwithStatus(t *testing.T) {
	psm, err := newMTTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	found, pid, sid, err := psm.GetPartnerAndSubscriberId("ucpadmin")
	if err != nil {
		t.Errorf("Unexpected error in GetPartnerIdWithStatus %v", err)
		return
	}
	t.Logf("Response: %b %v , %v", &found, pid, sid)
}
