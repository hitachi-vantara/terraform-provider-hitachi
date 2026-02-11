package admin

import (
	// "fmt"
	// adminmodel "terraform-provider-hitachi/hitachi/storage/admin/provisioner/model"
	"testing"
)

// go test -v -run TestGetStorageSystemInfo
func xTestGetStorageSystemInfo(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Scenario 1: Pass withEstimatedCapacities 'true' to GetStorageAdminInfo
	configurable_capacities := true
	resp, err := psm.GetStorageAdminInfo(configurable_capacities)
	if err != nil {
		t.Errorf("Unexpected error in Get (true) %v", err)
		return
	}
	t.Logf("Response (true): %v", resp)

	// Scenario 2: Pass withEstimatedCapacitie 'false' to GetStorageAdminInfo
	configurable_capacities = false
	resp, err = psm.GetStorageAdminInfo(configurable_capacities)
	if err != nil {
		t.Errorf("Unexpected error in Get (false) %v", err)
		return
	}
	t.Logf("Response (false): %v", resp)
}
