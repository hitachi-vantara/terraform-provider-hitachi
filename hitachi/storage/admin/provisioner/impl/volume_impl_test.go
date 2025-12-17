package admin

import (
	// "fmt"
	"encoding/json"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	"testing"
)

// go test -v -run TestGetVolumes
func xTestGetVolumes(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	// Example query with no filters (all optional params omitted)
	params := gwymodel.GetVolumeParams{}
	resp, err := psm.GetVolumes(params)
	if err != nil {
		t.Errorf("Unexpected error in GetVolumes: %v", err)
		return
	}

	if resp == nil {
		t.Errorf("Expected non-nil response from GetVolumes")
		return
	}

	t.Logf("GetVolumes Response: %+v", resp)
}

// go test -v -run TestGetVolumeByID
func xTestGetVolumeByID(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	// Replace with a valid ID for your test system
	testVolumeID := 6694

	resp, err := psm.GetVolumeByID(testVolumeID)
	if err != nil {
		t.Errorf("Unexpected error in GetVolumeByID: %v", err)
		return
	}

	if resp == nil {
		t.Errorf("Expected non-nil response from GetVolumeByID for ID %d", testVolumeID)
		return
	}

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Logf("Failed to marshal VolumeInfoList: %v", err)
		return
	}

	t.Logf("Response:\n%s", string(data))
}

// go test -v -run TestCreateVolume
func xTestCreateVolume(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	params := gwymodel.CreateVolumeParams{
		Capacity: 1024,
		Number:   ptr(1), // helper function for *int
		NicknameParam: gwymodel.VolumeNicknameParam{
			BaseName: "test-volume",
			// StartNumber:    ptr(1),
			// NumberOfDigits: ptr(3),
		},
		// SavingSetting:               ptr("DISABLE"), // or "DEDUPLICATION_AND_COMPRESSION"
		IsDataReductionShareEnabled: ptr(false),
		PoolID:                      0,
	}

	resp, err := psm.CreateVolume(params)
	if err != nil {
		t.Errorf("Unexpected error in CreateVolume: %v", err)
		return
	}

	if resp == "" {
		t.Errorf("Expected non-empty response from CreateVolume")
		return
	}

	t.Logf("Response:\n%s", resp)
}

// go test -v -run TestDeleteVolume
func xTestDeleteVolume(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	for id := 6708; id <= 6708; id++ {
		testVolumeID := id

		err = psm.DeleteVolume(testVolumeID)
		if err != nil {
			t.Errorf("Unexpected error in DeleteVolume for ID %d: %v", testVolumeID, err)
			continue
		}

		t.Logf("Successfully deleted volume ID %d", testVolumeID)
	}
}

// go test -v -run TestExpandVolume
func xTestExpandVolume(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	params := gwymodel.ExpandVolumeParams{
		Capacity: 512, // incremental MiB
	}

	testVolumeID := 6694
	err = psm.ExpandVolume(testVolumeID, params)
	if err != nil {
		t.Errorf("Unexpected error in ExpandVolume for ID %d: %v", testVolumeID, err)
		return
	}

	t.Logf("Successfully triggered expansion for volume ID %d with increment %d MiB", testVolumeID, params.Capacity)
}

// go test -v -run TestUpdateVolumeNickname
func xTestUpdateVolumeNickname(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	params := gwymodel.UpdateVolumeNicknameParams{
		Nickname: "new-nickname",
	}

	testVolumeID := 6694
	err = psm.UpdateVolumeNickname(testVolumeID, params)
	if err != nil {
		t.Errorf("Unexpected error in UpdateVolumeNickname for ID %d: %v", testVolumeID, err)
		return
	}

	t.Logf("Successfully updated nickname for volume ID %d to %q", testVolumeID, params.Nickname)
}

// go test -v -run TestUpdateVolumeReductionSettings
func xTestUpdateVolumeReductionSettings(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	params := gwymodel.UpdateVolumeReductionParams{
		SavingSetting:           ptr("DEDUPLICATION_AND_COMPRESSION"),
		CompressionAcceleration: ptr(true),
	}

	testVolumeID := 6694
	err = psm.UpdateVolumeReductionSettings(testVolumeID, params)
	if err != nil {
		t.Errorf("Unexpected error in UpdateVolumeReductionSettings for ID %d: %v", testVolumeID, err)
		return
	}

	t.Logf("Successfully triggered capacity reduction settings update for volume ID %d", testVolumeID)
}

// helper to convert a value to a pointer
func ptr[T any](v T) *T { return &v }
