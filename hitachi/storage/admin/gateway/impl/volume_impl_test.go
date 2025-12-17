package admin

import (
	"encoding/json"
	model "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	"testing"
)

// go test -v -run TestGetVolumes
func xTestGetVolumes(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	// Optional params; use pointers for integers and strings
	poolID := 0
	startVolumeID := 6680
	count := 10
	// poolName := "poolName_example"
	// nickname := "vol_example"

	params := model.GetVolumeParams{
		PoolID:           &poolID,
		PoolName:         nil,
		ServerID:         nil,
		ServerNickname:   nil,
		Nickname:         nil,
		MinTotalCapacity: nil,
		MaxTotalCapacity: nil,
		MinUsedCapacity:  nil,
		MaxUsedCapacity:  nil,
		StartVolumeID:    &startVolumeID,
		Count:            &count,
	}

	resp, err := psm.GetVolumes(params)
	if err != nil {
		t.Errorf("Unexpected error calling GetVolumes: %v", err)
		return
	}

	if resp == nil {
		t.Errorf("Expected response, got nil")
		return
	}

	logVolumeInfoList(t, resp)
}

func logVolumeInfoList(t *testing.T, resp *model.VolumeInfoList) {
	if resp == nil {
		t.Log("Response: <nil>")
		return
	}

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Logf("Failed to marshal VolumeInfoList: %v", err)
		return
	}

	t.Logf("Response:\n%s", string(data))
}

// go test -v -run TestGetVolumeByID
func xTestGetVolumeByID(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	testVolumeID := 6678

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

	params := model.CreateVolumeParams{
		Capacity: 1024,
		Number:   ptr(1), // helper function for *int
		NicknameParam: model.VolumeNicknameParam{
			BaseName:       "test-volume",
			StartNumber:    ptr(1),
			NumberOfDigits: ptr(3),
		},
		SavingSetting:               ptr("COMPRESSION"), // or "DEDUPLICATION_AND_COMPRESSION"
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

	for id := 6625; id <= 6625; id++ {
		testVolumeID := id

		err = psm.DeleteVolume(testVolumeID)
		if err != nil {
			t.Errorf("Unexpected error in DeleteVolume for ID %d: %v", testVolumeID, err)
			continue
		}

		t.Logf("Successfully deleted volume ID %d", testVolumeID)
	}
}

// go test -v -run TestUpdateVolumeNickname
func xTestUpdateVolumeNickname(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	volumeID := 6678
	params := model.UpdateVolumeNicknameParams{
		Nickname: "new-nick",
	}

	err = psm.UpdateVolumeNickname(volumeID, params)
	if err != nil {
		t.Errorf("Unexpected error in UpdateVolumeNickname: %v", err)
		return
	}

	t.Logf("Successfully updated nickname for volume ID %d to %q", volumeID, params.Nickname)
}

// go test -v -run TestUpdateVolumeReductionSettings
func xestUpdateVolumeReductionSettings(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	volumeID := 6678
	params := model.UpdateVolumeReductionParams{
		SavingSetting:           ptr("COMPRESSION"),
		CompressionAcceleration: ptr(false),
		// SavingSetting:           ptr("DEDUPLICATION_AND_COMPRESSION"),
		// CompressionAcceleration: ptr(false),
	}

	err = psm.UpdateVolumeReductionSettings(volumeID, params)
	if err != nil {
		t.Errorf("Unexpected error in UpdateVolumeReductionSettings: %v", err)
		return
	}

	t.Logf("Successfully triggered capacity reduction setting update for volume ID %d", volumeID)
}

// go test -v -run TestExpandVolume
func xTestExpandVolume(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	volumeID := 6678
	params := model.ExpandVolumeParams{
		Capacity: 2048, // increment in MiB
	}

	err = psm.ExpandVolume(volumeID, params)
	if err != nil {
		t.Errorf("Unexpected error in ExpandVolume: %v", err)
		return
	}

	t.Logf("Successfully triggered expansion for volume ID %d by %d MiB", volumeID, params.Capacity)
}

// helper to convert a value to a pointer
func ptr[T any](v T) *T { return &v }
