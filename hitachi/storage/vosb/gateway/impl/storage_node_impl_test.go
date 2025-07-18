package vssbstorage

import (
	"testing"
	// vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/gateway/model"
)

// go test -v -run TestGetAllStorageNodes
func xTestGetAll(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetStorageNodes()
	if err != nil {
		t.Errorf("Unexpected error in GetStorageNodes %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestGetAllStorageNodes
func xTestGetOne(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetStorageNode("5d991526-c48f-4c31-9a1d-40ab914915fb")
	if err != nil {
		t.Errorf("Unexpected error in GetStorageNodes %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestAddStorageNode
func xTestAddStorageNode(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if true {
		err := psm.doAddStorageNode(
			"/root/huitest2.csv",
			"",
			"vssb-789",
			"baremetal",
		)
		if err != nil {
			t.Errorf("Unexpected error in GetStorageNodes %v", err)
			return
		}
		return
	}

	resp, err := psm.GetStorageNodes()
	if err != nil {
		t.Errorf("Unexpected error in GetStorageNodes %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}
