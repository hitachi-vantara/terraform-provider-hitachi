package sanstorage

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	sanmodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
	// gwymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	reconmodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
)

func newSnapshotTestManager() (*sanStorageManager, error) {
	objStorage := sanmodel.StorageDeviceSettings{
		Serial:   12345,
		Username: "user1",
		Password: "mypswd",
		MgmtIP:   "10.10.11.12",
	}
	psm, err := newSanStorageManagerEx(objStorage)
	if err != nil {
		return nil, fmt.Errorf("unexpected error while creating TestManager %v", err)
	}
	return psm, nil
}

// go test -v -run ^TestReconcileGetSnapshot$
func xTestReconcileGetSnapshot(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	pvolID := 313
	mu := 3

	t.Run("Get One Snapshot", func(t *testing.T) {
		snapshot, err := psm.ReconcileGetSnapshot(&pvolID, &mu)
		if err != nil {
			t.Fatalf("Backend call failed: %v", err)
		}

		if snapshot == nil {
			t.Fatal("Expected snapshot data, got nil")
		}

		data, err := json.MarshalIndent(snapshot, "", "  ")
		if err != nil {
			t.Logf("Failed to marshal Response: %v", err)
			return
		}

		t.Logf("Response:\n%s", string(data))
	})
}

// go test -v -run ^TestReconcileReadSnapshot$
func xTestReconcileReadSnapshot(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	pvolID := 8145
	mu := 3
	snapshotGroupName := "TEST_SNAPSHOT_GROUP"

	input := reconmodel.SnapshotReconcilerInput{
		PvolLdevID:        &pvolID,
		MuNumber:          &mu,
		SnapshotGroupName: &snapshotGroupName,
	}

	t.Run("Reading Existing Snapshot", func(t *testing.T) {
		snapshot, err := psm.reconcileReadSnapshot(input)
		if err != nil {
			t.Fatalf("Backend call failed: %v", err)
		}

		if snapshot == nil {
			t.Fatal("Expected snapshot data, got nil")
		}

		data, err := json.MarshalIndent(snapshot, "", "  ")
		if err != nil {
			t.Logf("Failed to marshal Response: %v", err)
			return
		}

		t.Logf("Response:\n%s", string(data))
	})
}

// go test -v -run ^TestReconcileGetMultipleSnapshots_Combinations$
func xTestReconcileGetMultipleSnapshots_Combinations(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// --- Test Data (Replace with existing values from your array) ---
	pvolID := 8145
	svolID := 8146
	muNum := 3
	groupName := "spc-snapshot"

	// 1. Pattern: P-VOL LDEV + Snapshot Group Name
	t.Run("Pattern_Pvol_And_Group", func(t *testing.T) {
		input := reconmodel.SnapshotGetMultipleInput{
			PvolLdevID:        &pvolID,
			SnapshotGroupName: &groupName,
		}
		executeReconcileGetMultipleSnapshots(t, psm, input)
	})

	// 2. Pattern: P-VOL LDEV + MU Number
	t.Run("Pattern_Pvol_And_MU", func(t *testing.T) {
		input := reconmodel.SnapshotGetMultipleInput{
			PvolLdevID: &pvolID,
			MuNumber:   &muNum,
		}
		executeReconcileGetMultipleSnapshots(t, psm, input)
	})

	// 3. Pattern: Only P-VOL LDEV
	t.Run("Pattern_Only_Pvol", func(t *testing.T) {
		input := reconmodel.SnapshotGetMultipleInput{
			PvolLdevID: &pvolID,
		}
		executeReconcileGetMultipleSnapshots(t, psm, input)
	})

	// 4. Pattern: Only S-VOL LDEV
	t.Run("Pattern_Only_Svol", func(t *testing.T) {
		input := reconmodel.SnapshotGetMultipleInput{
			SvolLdevID: &svolID,
		}
		executeReconcileGetMultipleSnapshots(t, psm, input)
	})

	// 5. Pattern: Invalid Combination (Should fail validation)
	t.Run("Pattern_Invalid_Mixed", func(t *testing.T) {
		// Mixing Svol and Group is NOT in the allowed list
		input := reconmodel.SnapshotGetMultipleInput{
			SvolLdevID:        &svolID,
			SnapshotGroupName: &groupName,
		}
		_, err := psm.ReconcileGetMultipleSnapshots(input)
		if err == nil {
			t.Error("Expected validation error for mixed Svol/Group parameters, but got nil")
		} else {
			t.Logf("Correctly caught invalid combination: %v", err)
		}
	})
}

// go test -v -run ^TestReconcileGetMultipleSnapshotsRange$
func xTestReconcileGetMultipleSnapshotsRange(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	start := 100
	end := 200
	input := reconmodel.SnapshotGetMultipleRangeInput{
		StartPvolLdevID: &start,
		EndPvolLdevID:   &end,
	}

	t.Run("Get Multiple Snapshots Range", func(t *testing.T) {
		resp, err := psm.ReconcileGetMultipleSnapshotsRange(input)
		if err != nil {
			t.Fatalf("Backend range call failed: %v", err)
		}

		data, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			t.Logf("Failed to marshal Response: %v", err)
			return
		}

		t.Logf("Response:\n%s", string(data))
		t.Logf("Range query returned %d snapshots", len(resp.Data))
	})
}

// go test -v -run TestReconcileGetMultipleSnapshotsRangeAll
func xTestReconcileGetMultipleSnapshotsRangeAll(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	input := reconmodel.SnapshotGetMultipleRangeInput{
		StartPvolLdevID: nil,
		EndPvolLdevID:   nil,
	}

	t.Run("Get All Snapshots", func(t *testing.T) {
		resp, err := psm.ReconcileGetMultipleSnapshotsRange(input)
		if err != nil {
			t.Fatalf("Backend range call failed: %v", err)
		}

		data, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			t.Logf("Failed to marshal Response: %v", err)
			return
		}

		t.Logf("Response:\n%s", string(data))
		t.Logf("Range query returned %d snapshots", len(resp.Data))
	})
}

// go test -v -run ^TestReconcileSnapshotCreate_ParseCheck$
func xTestReconcileSnapshotCreate_ParseCheck(t *testing.T) {
	psm, _ := newSnapshotTestManager()

	// Testing the helper used inside reconcileSnapshotCreate
	testID := "123,1"
	pvol, mu, err := psm.parseSnapshotResID(testID)

	if err != nil {
		t.Fatalf("Failed to parse valid ID: %v", err)
	}

	if *pvol != 123 || *mu != 1 {
		t.Errorf("Expected 123,1 but got %d,%d", *pvol, *mu)
	}
}

// go test -v -run ^TestReconcileSnapshotCreate_No_Svol$
func xTestReconcileSnapshotCreate_No_Svol(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Failed to initialize manager: %v", err)
	}

	pvolID := 313
	poolID := 11
	groupName := "dummySnapG1"
	muNumber := 3

	input := reconmodel.SnapshotReconcilerInput{
		PvolLdevID:        &pvolID,
		SnapshotPoolID:    &poolID,
		SnapshotGroupName: &groupName,
		MuNumber:          &muNumber,
	}

	t.Run("Create_and_Wait_for_Pair", func(t *testing.T) {
		snapshot, err := psm.reconcileSnapshotCreate(input)

		// Assertions
		if err != nil {
			t.Fatalf("ReconcileSnapshotCreate failed: %v", err)
		}

		data, err := json.MarshalIndent(snapshot, "", "  ")
		if err != nil {
			t.Logf("Failed to marshal Response: %v", err)
			return
		}

		t.Logf("Response:\n%s", string(data))
	})
}

// go test -v -run ^TestReconcileSnapshotCreate_With_Svol$
func xTestReconcileSnapshotCreate_With_Svol(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Failed to initialize manager: %v", err)
	}

	pvolID := 313
	poolID := 11
	groupName := "dummySnapG1"
	muNumber := 3
	svolID := 314
	// isClone := true
	// canCascade := true

	input := reconmodel.SnapshotReconcilerInput{
		PvolLdevID:        &pvolID,
		SnapshotPoolID:    &poolID,
		SnapshotGroupName: &groupName,
		MuNumber:          &muNumber,
		SvolLdevID:        &svolID,
		// IsClone:           &isClone,
		// CanCascade:        &canCascade,
	}

	t.Run("Create_and_Wait_for_Pair", func(t *testing.T) {
		snapshot, err := psm.reconcileSnapshotCreate(input)

		// Assertions
		if err != nil {
			t.Fatalf("ReconcileSnapshotCreate failed: %v", err)
		}

		data, err := json.MarshalIndent(snapshot, "", "  ")
		if err != nil {
			t.Logf("Failed to marshal Response: %v", err)
			return
		}

		t.Logf("Response:\n%s", string(data))
	})
}

// go test -v -run ^TestReconcileSnapshot_Split$
func xTestReconcileSnapshot_Split(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	psm, _ := newSnapshotTestManager()

	pvolID, mu := 313, 3
	input := reconmodel.SnapshotReconcilerInput{
		PvolLdevID: &pvolID,
		MuNumber:   &mu,
	}

	t.Run("Split_Snapshot", func(t *testing.T) {
		snapshot, err := psm.reconcileSnapshotSplit(input)
		if err != nil {
			t.Fatalf("Split failed: %v", err)
		}
		data, err := json.MarshalIndent(snapshot, "", "  ")
		if err != nil {
			t.Logf("Failed to marshal Response: %v", err)
			return
		}

		t.Logf("Response:\n%s", string(data))
	})
}

// go test -v -run ^TestReconcileSnapshot_Resync$
func xTestReconcileSnapshot_Resync(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	psm, _ := newSnapshotTestManager()

	pvolID, mu := 313, 3
	input := reconmodel.SnapshotReconcilerInput{
		PvolLdevID: &pvolID,
		MuNumber:   &mu,
	}

	t.Run("Resync_Snapshot", func(t *testing.T) {
		snapshot, err := psm.reconcileSnapshotResync(input)
		if err != nil {
			t.Fatalf("Resync failed: %v", err)
		}
		data, err := json.MarshalIndent(snapshot, "", "  ")
		if err != nil {
			t.Logf("Failed to marshal Response: %v", err)
			return
		}

		t.Logf("Response:\n%s", string(data))
	})
}

// go test -v -run ^TestReconcileSnapshot_Restore$
func xTestReconcileSnapshot_Restore(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	psm, _ := newSnapshotTestManager()
	pvolID, mu := 8145, 1
	input := reconmodel.SnapshotReconcilerInput{
		PvolLdevID: &pvolID,
		MuNumber:   &mu,
	}

	t.Run("Restore_From_Svol", func(t *testing.T) {
		// Step 1: Ensure it's split
		_, _ = psm.reconcileSnapshotSplit(input)

		// Step 2: Restore
		snapshot, err := psm.reconcileSnapshotRestore(input)
		if err != nil {
			t.Fatalf("Restore failed: %v", err)
		}
		data, err := json.MarshalIndent(snapshot, "", "  ")
		if err != nil {
			t.Logf("Failed to marshal Response: %v", err)
			return
		}

		t.Logf("Response:\n%s", string(data))
	})
}

// go test -v -run ^TestReconcileSnapshot_Unassign_Volume$
func xTestReconcileSnapshot_Unassign_Volume(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	psm, _ := newSnapshotTestManager()

	pvolID, mu := 313, 3
	input := reconmodel.SnapshotReconcilerInput{
		PvolLdevID: &pvolID,
		MuNumber:   &mu,
	}

	t.Run("Unassign_Volume", func(t *testing.T) {
		snapshot, err := psm.reconcileSnapshotUnassign(input)
		if err != nil {
			t.Logf("error: %v", err)
		}

		data, err := json.MarshalIndent(snapshot, "", "  ")
		if err != nil {
			t.Logf("Failed to marshal Response: %v", err)
			return
		}

		t.Logf("Response:\n%s", string(data))

	})
}

// go test -v -run ^TestReconcileSnapshot_Assign_Volume$
func xTestReconcileSnapshot_Assign_Volume(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	psm, _ := newSnapshotTestManager()

	pvolID, mu, svolID := 313, 3, 314
	input := reconmodel.SnapshotReconcilerInput{
		PvolLdevID: &pvolID,
		MuNumber:   &mu,
		SvolLdevID: &svolID,
	}
	t.Run("Assign_Volume", func(t *testing.T) {
		snapshot, err := psm.reconcileSnapshotAssign(input)
		if err != nil {
			t.Fatalf("Assign failed: %v", err)
		}
		data, err := json.MarshalIndent(snapshot, "", "  ")
		if err != nil {
			t.Logf("Failed to marshal Response: %v", err)
			return
		}

		t.Logf("Response:\n%s", string(data))
	})
}

// go test -v -run ^TestReconcileSnapshot_Delete$
func xTestReconcileSnapshot_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	psm, _ := newSnapshotTestManager()
	pvolID, mu := 313, 3
	input := reconmodel.SnapshotReconcilerInput{
		PvolLdevID: &pvolID,
		MuNumber:   &mu,
	}

	t.Run("Delete_Snapshot_Pair", func(t *testing.T) {
		_, err := psm.ReconcileSnapshotDelete(input)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}
	})
}

// Helper to dry up the test code
func executeReconcileGetMultipleSnapshots(t *testing.T, psm *sanStorageManager, input reconmodel.SnapshotGetMultipleInput) {
	t.Helper()
	resp, err := psm.ReconcileGetMultipleSnapshots(input)
	if err != nil {
		t.Fatalf("API Call failed: %v", err)
	}

	if resp == nil {
		t.Fatal("Response is nil")
	}

	t.Logf("Found %d snapshots", len(resp.Data))
	if len(resp.Data) > 0 {
		data, _ := json.MarshalIndent(resp.Data, "", "  ") // Log the first item for brevity
		t.Logf("First Result Sample:\n%s", string(data))
	}
}

// Helpers for pointers
func ptrString(s string) *string { return &s }
func ptrInt(i int) *int          { return &i }
func ptrBool(b bool) *bool       { return &b }

// go test -v -run ^TestValidateSnapshotCreate_TI_Std$
func xTestValidateSnapshotCreate_TI_Std(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Failed to initialize manager: %v", err)
	}

	// 40014 data
	htiPoolID := 5
	hdpPoolID := 2
	pvolStandard := 320
	pvolWithCompression := 3240
	svolStandard := 314

	tests := []struct {
		name        string
		input       reconmodel.SnapshotReconcilerInput
		wantErr     bool
		errContains string
	}{
		// --- POSITIVE TESTS ---
		{
			name: "Positive: Standard TI (HTI Pool) No Svol",
			input: reconmodel.SnapshotReconcilerInput{
				SnapshotGroupName: ptrString("std_ti_test"),
				SnapshotPoolID:    &htiPoolID,
				PvolLdevID:        &pvolStandard,
				IsClone:           ptrBool(false),
			},
			wantErr: false,
		},
		{
			name: "Positive: Standard TI (HTI Pool) With Svol",
			input: reconmodel.SnapshotReconcilerInput{
				SnapshotGroupName: ptrString("std_ti_test"),
				SnapshotPoolID:    &htiPoolID,
				PvolLdevID:        &pvolStandard,
				SvolLdevID:        &svolStandard,
				IsClone:           ptrBool(false),
			},
			wantErr: false,
		},
		{
			name: "Positive: Standard TI (HDP Pool Allowed)",
			input: reconmodel.SnapshotReconcilerInput{
				SnapshotGroupName: ptrString("std_ti_test"),
				SnapshotPoolID:    &hdpPoolID,
				PvolLdevID:        &pvolStandard,
				SvolLdevID:        &svolStandard,
				IsClone:           ptrBool(false),
			},
			wantErr: false,
		},
		{
			name: "Positive: Standard TI Clone Case-Insensitive Speed",
			input: reconmodel.SnapshotReconcilerInput{
				SnapshotGroupName: ptrString("std_clone_ok"),
				SnapshotPoolID:    &htiPoolID,
				PvolLdevID:        &pvolStandard,
				SvolLdevID:        &svolStandard,
				IsClone:           ptrBool(true),
				CanCascade:        ptrBool(true),
				ClonesAutomation:  ptrBool(true),
				CopySpeed:         ptrString("FASTER"), // Case insensitive test
			},
			wantErr: false,
		},
		// --- NEGATIVE TESTS (GLOBAL & STD) ---
		{
			name: "Negative: Missing S-VOL for Clone",
			input: reconmodel.SnapshotReconcilerInput{
				SnapshotGroupName: ptrString("neg_clone_svol"),
				SnapshotPoolID:    &htiPoolID,
				PvolLdevID:        &pvolStandard,
				IsClone:           ptrBool(true),
				CanCascade:        ptrBool(true),
			},
			wantErr:     true,
			errContains: "svolLdevId is required when isClone is true",
		},
		{
			name: "Negative: ForceCopy required for compressed P-VOL",
			input: reconmodel.SnapshotReconcilerInput{
				SnapshotGroupName:        ptrString("neg_forcecopy"),
				SnapshotPoolID:           &htiPoolID,
				PvolLdevID:               &pvolWithCompression,
				IsDataReductionForceCopy: ptrBool(false),
			},
			wantErr:     true,
			errContains: "isDataReductionForceCopy must be true when P-VOL capacity saving is enabled",
		},
		{
			name: "Negative: isClone and autoSplit are mutually exclusive",
			input: reconmodel.SnapshotReconcilerInput{
				SnapshotGroupName: ptrString("neg_excl"),
				SnapshotPoolID:    &htiPoolID,
				PvolLdevID:        &pvolStandard,
				IsClone:           ptrBool(true),
				AutoSplit:         ptrBool(true),
				SvolLdevID:        ptrInt(svolStandard),
			},
			wantErr:     true,
			errContains: "cannot specify both isClone and autoSplit",
		},
		{
			name: "Negative Global: clonesAutomation without isClone",
			input: reconmodel.SnapshotReconcilerInput{
				SnapshotGroupName: ptrString("neg_auto_no_clone"),
				SnapshotPoolID:    &htiPoolID,
				PvolLdevID:        &pvolStandard,
				CanCascade:        ptrBool(true),
				IsClone:           ptrBool(false),
				ClonesAutomation:  ptrBool(true), // Conflict
			},
			wantErr:     true,
			errContains: "clonesAutomation can only be specified when isClone is true",
		},
		{
			name: "Negative Global: copySpeed without clonesAutomation",
			input: reconmodel.SnapshotReconcilerInput{
				SnapshotGroupName: ptrString("neg_speed_no_auto"),
				SnapshotPoolID:    &htiPoolID,
				PvolLdevID:        &pvolStandard,
				CanCascade:        ptrBool(true),
				IsClone:           ptrBool(true),
				ClonesAutomation:  ptrBool(false),
				CopySpeed:         ptrString("faster"), // Conflict
				SvolLdevID:        ptrInt(314),
			},
			wantErr:     true,
			errContains: "copySpeed requires both isClone and clonesAutomation",
		},
		{
			name: "Negative: copySpeed invalid value",
			input: reconmodel.SnapshotReconcilerInput{
				SnapshotGroupName: ptrString("neg_speed_val"),
				SnapshotPoolID:    &htiPoolID,
				PvolLdevID:        &pvolStandard,
				CanCascade:        ptrBool(true),
				IsClone:           ptrBool(true),
				ClonesAutomation:  ptrBool(true),
				CopySpeed:         ptrString("ultra-fast"),
				SvolLdevID:        ptrInt(svolStandard),
			},
			wantErr:     true,
			errContains: "invalid copySpeed: ultra-fast",
		},
		{
			name: "Negative: negative muNumber",
			input: reconmodel.SnapshotReconcilerInput{
				SnapshotGroupName: ptrString("neg_mu_val"),
				SnapshotPoolID:    &htiPoolID,
				PvolLdevID:        &pvolStandard,
				MuNumber:          ptrInt(-5),
			},
			wantErr:     true,
			errContains: "muNumber must be between 0 and 1023",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := psm.validateSnapshotCreateInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: got err %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if err != nil && tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("%s: error '%v' does not contain '%s'", tt.name, err, tt.errContains)
			}
		})
	}
}

// go test -v -run ^TestValidateSnapshotCreate_TIA$
func xTestValidateSnapshotCreate_TIA(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Failed to initialize manager: %v", err)
	}

	// 40014 data
	hdpPoolID := 2
	hdpPoolIDMismatch := 6
	// htiPoolID := 5 //  This won't flag it as TIA
	pvolTIA := 3240 // DRS + Compression
	// pvolNoCompression := 313 // Standard HDP Vol, No DRS, No Compression. This won't flag it as TIA
	svolTIA := 3246         // Same Pool DRS
	svolMismatchPool := 314 // Different Pool

	tests := []struct {
		name        string
		input       reconmodel.SnapshotReconcilerInput
		wantErr     bool
		errContains string
	}{
		{
			name: "Positive: TIA Valid Request",
			input: reconmodel.SnapshotReconcilerInput{
				SnapshotGroupName:        ptrString("tia_pos"),
				SnapshotPoolID:           &hdpPoolID,
				PvolLdevID:               &pvolTIA,
				IsClone:                  ptrBool(false),
				CanCascade:               ptrBool(true),
				IsDataReductionForceCopy: ptrBool(true),
			},
			wantErr: false,
		},
		{
			name: "Positive: TIA (HDP Pool, Aligned, ForcedCopy)",
			input: reconmodel.SnapshotReconcilerInput{
				SnapshotGroupName:        ptrString("tia_pos_test"),
				SnapshotPoolID:           &hdpPoolID,
				PvolLdevID:               &pvolTIA,
				SvolLdevID:               &svolTIA,
				IsClone:                  ptrBool(false),
				CanCascade:               ptrBool(true),
				IsDataReductionForceCopy: ptrBool(true),
			},
			wantErr: false,
		},
		// it won't hit this because it thinks it is TI Std
		// {
		// 	name: "Negative TIA: (HTI Pool Not Allowed)",
		// 	input: reconmodel.SnapshotReconcilerInput{
		// 		SnapshotGroupName:        ptrString("tia_neg_cascade"),
		// 		SnapshotPoolID:           &htiPoolID,
		// 		PvolLdevID:               &pvolTIA,
		// 		CanCascade:               ptrBool(true), // Always True rule
		// 		IsDataReductionForceCopy: ptrBool(true),
		// 	},
		// 	wantErr:     true,
		// 	errContains: "Thin Image Advanced pairs require an HDP pool",
		// 	// go test -v -run "TestValidateSnapshotCreate_TIA/Negative_TIA:_\(HTI_Pool_Not_Allowed\)"
		// },
		{
			name: "Negative TIA: canCascade Must Be True",
			input: reconmodel.SnapshotReconcilerInput{
				SnapshotGroupName: ptrString("tia_neg_cascade"),
				SnapshotPoolID:    &hdpPoolID,
				PvolLdevID:        &pvolTIA,
				CanCascade:        ptrBool(false), // Always True rule

			},
			wantErr:     true,
			errContains: "canCascade must be true for Thin Image Advanced",
		},
		{
			name: "Negative TIA: P-VOL Pool Mismatch",
			input: reconmodel.SnapshotReconcilerInput{
				SnapshotGroupName: ptrString("tia_neg_pool_mismatch"),
				SnapshotPoolID:    &hdpPoolIDMismatch,
				PvolLdevID:        &pvolTIA,
			},
			wantErr:     true,
			errContains: "must be the same as snapshot pool ID",
		},
		// it won't hit this because it thinks it is TI Std (no DRS)
		// {
		// 	name: "Negative TIA: P-VOL No Compression",
		// 	input: reconmodel.SnapshotReconcilerInput{
		// 		SnapshotGroupName: ptrString("tia_neg_no_comp"),
		// 		SnapshotPoolID:    &hdpPoolID,
		// 		PvolLdevID:        &pvolNoCompression,
		// 	},
		// 	wantErr:     true,
		// 	errContains: "must have capacity saving enabled for TIA",
		// 	// go test -v -run "TestValidateSnapshotCreate_TIA/Negative_TIA:_P-VOL_No_Compression"
		// },
		{
			name: "Negative TIA: isClone Forbidden",
			input: reconmodel.SnapshotReconcilerInput{
				SnapshotGroupName: ptrString("tia_neg_clone"),
				SnapshotPoolID:    &hdpPoolID,
				PvolLdevID:        &pvolTIA,
				IsClone:           ptrBool(true),
			},
			wantErr:     true,
			errContains: "Thin Image Advanced pairs do not support isClone: true",
		},
		{
			name: "Negative TIA: S-VOL Pool Mismatch",
			input: reconmodel.SnapshotReconcilerInput{
				SnapshotGroupName: ptrString("tia_neg_svol_pool"),
				SnapshotPoolID:    &hdpPoolID,
				PvolLdevID:        &pvolTIA,
				SvolLdevID:        &svolMismatchPool,
				CanCascade:        ptrBool(true),
			},
			wantErr:     true,
			errContains: "S-VOL pool ID",
		},
		{
			name: "Negative TIA: ForceCopy Required",
			input: reconmodel.SnapshotReconcilerInput{
				SnapshotGroupName:        ptrString("tia_neg_force"),
				SnapshotPoolID:           &hdpPoolID,
				PvolLdevID:               &pvolTIA,
				CanCascade:               ptrBool(true),
				IsDataReductionForceCopy: ptrBool(false),
			},
			wantErr:     true,
			errContains: "isDataReductionForceCopy must be true for Thin Image Advanced",
		},
		{
			name: "Negative Global: isClone and autoSplit are mutually exclusive",
			input: reconmodel.SnapshotReconcilerInput{
				SnapshotGroupName: ptrString("neg_excl"),
				SnapshotPoolID:    &hdpPoolID,
				PvolLdevID:        &pvolTIA,
				IsClone:           ptrBool(false),
				AutoSplit:         ptrBool(true), // Conflict
				SvolLdevID:        ptrInt(314),
			},
			wantErr:     true,
			errContains: "cannot specify both isClone and autoSplit",
		},
		{
			name: "Negative Global: clonesAutomation without isClone",
			input: reconmodel.SnapshotReconcilerInput{
				SnapshotGroupName: ptrString("neg_auto_no_clone"),
				SnapshotPoolID:    &hdpPoolID,
				PvolLdevID:        &pvolTIA,
				IsClone:           ptrBool(false),
				ClonesAutomation:  ptrBool(true), // Conflict
			},
			wantErr:     true,
			errContains: "clonesAutomation can only be specified when isClone is true",
		},
		{
			name: "Negative Global: copySpeed without clonesAutomation",
			input: reconmodel.SnapshotReconcilerInput{
				SnapshotGroupName: ptrString("neg_speed_no_auto"),
				SnapshotPoolID:    &hdpPoolID,
				PvolLdevID:        &pvolTIA,
				IsClone:           ptrBool(false),
				ClonesAutomation:  ptrBool(false),
				CopySpeed:         ptrString("faster"), // Conflict
				SvolLdevID:        ptrInt(314),
			},
			wantErr:     true,
			errContains: "copySpeed requires both isClone and clonesAutomation",
		},
		{
			name: "Negative Global: invalid copySpeed value",
			input: reconmodel.SnapshotReconcilerInput{
				SnapshotGroupName: ptrString("neg_speed_val"),
				SnapshotPoolID:    &hdpPoolID,
				PvolLdevID:        &pvolTIA,
				IsClone:           ptrBool(false),
				ClonesAutomation:  ptrBool(true),
				CopySpeed:         ptrString("ultra-fast"), // Invalid string
				SvolLdevID:        ptrInt(314),
			},
			wantErr:     true,
			errContains: "invalid copySpeed: ultra-fast",
		},
		{
			name: "Negative Global: negative muNumber",
			input: reconmodel.SnapshotReconcilerInput{
				SnapshotGroupName: ptrString("neg_mu_val"),
				SnapshotPoolID:    &hdpPoolID,
				PvolLdevID:        &pvolTIA,
				MuNumber:          ptrInt(-5),
			},
			wantErr:     true,
			errContains: "muNumber must be 0 or greater",
		},
		{
			name: "Negative TIA: Retention Period requires AutoSplit",
			input: reconmodel.SnapshotReconcilerInput{
				SnapshotGroupName: ptrString("tia_neg_retention"),
				SnapshotPoolID:    &hdpPoolID,
				PvolLdevID:        &pvolTIA,
				AutoSplit:         ptrBool(false),
				RetentionPeriod:   ptrInt(24),
			},
			wantErr:     true,
			errContains: "retentionPeriod requires autoSplit to be true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := psm.validateSnapshotCreateInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: got err %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}
