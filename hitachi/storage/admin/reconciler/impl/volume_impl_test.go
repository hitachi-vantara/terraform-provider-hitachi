package admin

import (
	"encoding/json"
	"testing"

	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
)

// go test -v -run TestReconcileReadAdminVolumes
func xTestReconcileReadAdminVolumes(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	// Provide existing volume IDs (adjust to actual test data)
	volumeIDs := []int{6739, 6756, 6757, 6758, 6783} // replace with valid IDs in your system

	volInfos, existing, err := psm.ReconcileReadAdminVolumes(volumeIDs)
	if err != nil {
		t.Errorf("Unexpected error in ReconcileReadAdminVolumes: %v", err)
		return
	}

	if len(existing) == 0 {
		t.Errorf("Expected some existing volumes but got none")
		return
	}

	b, _ := json.MarshalIndent(volInfos, "", "  ")
	t.Logf("Existing IDs: %v\nVolume Info:\n%s", existing, string(b))
}

// go test -v -run TestReconcileCreateAdminVolumes
func xTestReconcileCreateAdminVolumes(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	params := gwymodel.CreateVolumeParams{
		Capacity: 1024,
		Number:   ptr(4),
		NicknameParam: gwymodel.VolumeNicknameParam{
			BaseName:       "test",
			StartNumber:    ptr(1),
			NumberOfDigits: ptr(2),
		},
		IsDataReductionShareEnabled: ptr(false),
		PoolID:                      0,
	}

	resp, err := psm.ReconcileCreateAdminVolumes(params)
	if err != nil {
		t.Errorf("Unexpected error in ReconcileCreateAdminVolumes: %v", err)
		return
	}

	if len(resp) == 0 {
		t.Errorf("Expected non-empty response from ReconcileCreateAdminVolumes")
		return
	}

	t.Logf("ReconcileCreateAdminVolumes Response: %v", resp)
}

// go test -v -run TestReconcileDeleteAdminVolumes
func xTestReconcileDeleteAdminVolumes(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	// Replace with actual IDs that exist or just-created ones
	volumeIDs := []int{6783}

	err = psm.ReconcileDeleteAdminVolumes(volumeIDs)
	if err != nil {
		t.Errorf("Unexpected error in ReconcileDeleteAdminVolumes: %v", err)
		return
	}

	t.Logf("ReconcileDeleteAdminVolumes completed successfully for: %v", volumeIDs)
}

// go test -v -run TestReconcileUpdateAdminVolume
func xTestReconcileUpdateAdminVolume(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	tests := []struct {
		name       string
		volumeID   int
		updateFunc func() gwymodel.CreateVolumeParams
	}{
		{
			name:     "Update Capacity",
			volumeID: 6783,
			updateFunc: func() gwymodel.CreateVolumeParams {
				return gwymodel.CreateVolumeParams{
					Capacity: 2048, // MiB
				}
			},
		},
		{
			name:     "Update SavingSetting and CompressionAcceleration",
			volumeID: 6783,
			updateFunc: func() gwymodel.CreateVolumeParams {
				setting := "DEDUPLICATION_AND_COMPRESSION"
				accel := true
				return gwymodel.CreateVolumeParams{
					SavingSetting:           &setting,
					CompressionAcceleration: &accel,
				}
			},
		},
		{
			name:     "Update Nickname",
			volumeID: 6783,
			updateFunc: func() gwymodel.CreateVolumeParams {
				return gwymodel.CreateVolumeParams{
					Capacity: 2048, // MiB
					NicknameParam: gwymodel.VolumeNicknameParam{
						BaseName: "renamedVol",
					},
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := tt.updateFunc()

			// Pretty-print params for debugging
			paramsJSON, _ := json.MarshalIndent(params, "", "  ")
			t.Logf("Params for %s: %s", tt.name, string(paramsJSON))

			err := psm.ReconcileUpdateAdminVolume(tt.volumeID, params)
			if err != nil {
				t.Errorf("Unexpected error in %s: %v", tt.name, err)
				return
			}

			t.Logf("%s succeeded for VolumeID %d", tt.name, tt.volumeID)
		})
	}
}

// helper to convert a value to a pointer
func ptr[T any](v T) *T { return &v }
