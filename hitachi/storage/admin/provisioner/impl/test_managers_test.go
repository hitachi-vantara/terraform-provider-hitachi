package admin

import (
	"fmt"

	model "terraform-provider-hitachi/hitachi/storage/admin/provisioner/model"
)

// NOTE: This is a placeholder for unit tests.
// Do NOT commit real values to the repository.
// Create a local copy of this file with your actual test storage details.

func newTestManager() (*adminStorageManager, error) {
	// Example placeholder storage settings
	objStorage := model.StorageDeviceSettings{
		Serial:   12345,
		Username: "username",
		Password: "password",
		MgmtIP:   "127.0.0.1",
	}

	psm, err := newAdminStorageManagerEx(objStorage)
	if err != nil {
		return nil, fmt.Errorf("unexpected error while creating newAdminStorageManagerEx: %v", err)
	}
	return psm, nil
}
