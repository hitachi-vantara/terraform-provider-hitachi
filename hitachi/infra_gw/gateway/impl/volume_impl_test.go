package infra_gw

import (
	"testing"
)

func xTestGetMtPartner(t *testing.T) {
	psm, err := newMTTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	volumes, err := psm.GetMTVolumes("storage-12d27566fa9feb38f728801ae15997b3")
	if err != nil {
		t.Errorf("Unexpected error in GetDynamicPools %v", err)
		return
	}
	
	t.Logf("Response: %v", volumes)
}

func xTestGetVolumeByID(t *testing.T) {
	psm, err := newMTTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	volumes, err := psm.GetVolumeByID("storage-39f4eef0175c754bb90417358b0133c3", "storagevolume-beb340bb55d9dcc34182655d074444d3")
	if err != nil {
		t.Errorf("Unexpected error in GetVolumeByID %v", err)
		return
	}

	t.Logf("Response: %v", volumes)
}
