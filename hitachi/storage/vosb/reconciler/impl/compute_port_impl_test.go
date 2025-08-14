package vssbstorage

import (
	"testing"
)

// go test -v -run TestAllowChapUsersToAccessComputePort
func xTestAllowChapUsersToAccessComputePort(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	portId := "5f07176a-e10d-47b7-99b4-57b93806048b" /// name of the port "001-iSCSI-002"
	authMode := "CHAP"
	chapUsers := []string{"chapuser1", "chapuser2"}

	err = psm.AllowChapUsersToAccessComputePort(portId, authMode, chapUsers)
	if err != nil {
		t.Errorf("Unexpected error in AllowChapUsersToAccessComputePort %v", err)
		return
	}
	t.Logf("Error: %v", err)
}
