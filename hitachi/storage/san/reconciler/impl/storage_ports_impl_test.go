package sanstorage

import (
	"testing"
)

// go test -v -run TestGetStoragePorts
func xTestGetStoragePorts(t *testing.T) {
	psm, err := newIscsiTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetStoragePorts()
	if err != nil {
		t.Errorf("Unexpected error in GetStoragePorts %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestGetStoragePortByPortId
func xTestGetStoragePortByPortId(t *testing.T) {
	psm, err := newIscsiTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	portId := "CL4-C"
	resp, err := psm.GetStoragePortByPortId(portId)
	if err != nil {
		t.Errorf("Unexpected error in GetStoragePortByPortId %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}
