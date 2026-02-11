package vssbstorage

import (
	"os"
	"testing"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"
)

// go test -v -run TestRestoreConfigurationDefinitionFile
func xTestRestoreConfigurationDefinitionFile(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	err = psm.RestoreConfigurationDefinitionFile(&reconcilermodel.CreateConfigurationFileParam{})
	if err != nil {
		t.Errorf("Unexpected error in RestoreConfigurationDefinitionFile %v", err)
		return
	}
	t.Logf("Successfully restored configuration file")
}

// go test -v -run TestDownloadConfigurationFile
func xTestDownloadConfigurationFile(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error from newTestManager: %v", err)
	}

	filePath, err := psm.DownloadConfigurationFile("/tmp")
	if err != nil {
		t.Fatalf("DownloadConfigurationFile failed: %v", err)
	}

	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		t.Errorf("path does not exist: %s", filePath)
	}
	if info.IsDir() {
		t.Errorf("expected a file, but got a directory: %s", filePath)
	}

	t.Logf("Successfully found downloaded file: %s", filePath)
}
