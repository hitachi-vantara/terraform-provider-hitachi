package infra_gw

import (
	"fmt"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
	"testing"

	"github.com/google/uuid"
)

// newDynamicPoolTestManager is for Testing and provide structure information for connection
func newMTTestManager() (*infraGwManager, error) {

	// Following storage has iscsi port
	PartnerId := "a8d1f065-a9e7-42cf-b565-a67466fec549"
	SubscriberId := "f260ca80-ec80-423c-9af8-dcac117bd068"
	setting := model.InfraGwSettings{
		Username:     "ucpadmin",
		Password:     "Passw0rd!",
		PartnerId:    &PartnerId,
		SubscriberId: &SubscriberId,
		Address:      "172.25.22.81",
	}

	psm, err := newInfraGwManagerEx(setting)
	if err != nil {
		return nil, fmt.Errorf("unexpected error while creating newDynamicPoolTestManager %v", err)
	}
	return psm, nil
}

// go test -v -run TestGetDynamicPools
func xTestGetPartners(t *testing.T) {
	psm, err := newMTTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	partners, err := psm.GetAllPartners()
	if err != nil {
		t.Errorf("Unexpected error in GetDynamicPools %v", err)
		return
	}
	t.Logf("Response: %v", partners)
}

func xTestCreatePartner(t *testing.T) {
	psm, err := newMTTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	newUUID := uuid.New().String()
	reqBody := &model.RegisterPartnerReq{
		Name:        "TestPartner",
		PartnerID:   newUUID,
		Description: "This is a test Part",
	}
	partners, err := psm.RegisterPartner(reqBody)
	if err != nil {
		t.Errorf("Unexpected error in GetDynamicPools %v", err)
		return
	}
	t.Logf("Response: %v", partners)
}
