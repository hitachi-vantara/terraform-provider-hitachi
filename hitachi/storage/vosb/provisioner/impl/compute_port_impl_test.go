package vssbstorage

import (
	"testing"
)

// go test -v -run TestDeleteAllChapUsersFromComputePort
func xTestDeleteAllChapUsersFromComputePort(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	portId := "5f07176a-e10d-47b7-99b4-57b93806048b" /// name of the port "001-iSCSI-002"

	err = psm.DeleteAllChapUsersFromComputePort(portId)
	if err != nil {
		t.Errorf("Unexpected error in DeleteAllChapUsersFromComputePort %v", err)
		return
	}
	t.Logf("Error: %v", err)
}

// go test -v -run TestAddChapUsersToComputePort
func xTestAddChapUsersToComputePort(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	portId := "5f07176a-e10d-47b7-99b4-57b93806048b" /// name of the port "001-iSCSI-002"
	chapUsers := []string{"chapuser1", "chapuser2"}

	err = psm.AddChapUsersToComputePort(portId, chapUsers)
	if err != nil {
		t.Errorf("Unexpected error in AddChapUsersToComputePort %v", err)
		return
	}
	t.Logf("Error: %v", err)
}
