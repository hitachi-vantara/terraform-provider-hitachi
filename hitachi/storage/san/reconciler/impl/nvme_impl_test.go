package sanstorage

import (
	"encoding/json"
	"fmt"
	"testing"

	// "terraform-provider-hitachi/hitachi/common/utils"
	reconmodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
)

func newNvmReconcilerTestManager() (*sanStorageManager, error) {
	objStorage := reconmodel.StorageDeviceSettings{
		Serial:   12345,
		Username: "user1",
		Password: "mypswd",
		MgmtIP:   "10.10.11.12",
	}
	psm, err := newSanStorageManagerEx(objStorage)
	if err != nil {
		return nil, fmt.Errorf("unexpected error while creating NvmReconcilerTestManager %v", err)
	}
	return psm, nil
}

// go test -v -run ^TestReconcileGetNvmSubsystem$
func ^TestReconcileGetNvmSubsystem(t *testing.T) {
	psm, err := newNvmReconcilerTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	subsystemID := 1

	t.Run("Get One NVM Subsystem", func(t *testing.T) {
		subsystem, err := psm.ReconcileGetNvmSubsystem(subsystemID)
		if err != nil {
			t.Fatalf("Backend call failed: %v", err)
		}

		if subsystem == nil {
			t.Fatal("Expected subsystem data, got nil")
		}

		data, _ := json.MarshalIndent(subsystem, "", "  ")
		t.Logf("Response:\n%s", string(data))
	})
}

// go test -v -run ^TestReconcileGetMultipleNvmSubsystems$
func ^TestReconcileGetMultipleNvmSubsystems(t *testing.T) {
	psm, err := newNvmReconcilerTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	t.Run("Hydrate All Fields", func(t *testing.T) {
		input := reconmodel.NvmSubsystemGetMultipleInput{
			IncludeNamespaces: true,
			IncludePorts:      true,
			IncludeHostNqns:   true,
		}

		resp, err := psm.ReconcileGetMultipleNvmSubsystems(input)
		if err != nil {
			t.Fatalf("Reconcile failed: %v", err)
		}

		if resp == nil || len(resp.Data) == 0 {
			t.Log("No subsystems found to validate")
			return
		}

		// Validate that hydration fields are populated for the first item
		first := resp.Data[0]
		t.Logf("Validating Subsystem ID: %d", first.NvmSubsystemId)
		
		if first.NvmSubsystemNqn == "" {
			t.Errorf("Target NQN was not hydrated")
		}
		t.Logf("Target NQN: %s", first.NvmSubsystemNqn)

		if len(first.Namespaces) == 0 {
			t.Log("Warning: No namespaces found (could be valid, check storage state)")
		}
		
		if len(first.HostNqns) == 0 {
			t.Log("Warning: No Host NQNs found (could be valid, check storage state)")
		}

		data, _ := json.MarshalIndent(resp.Data, "", "  ")
		t.Logf("Full Hydrated Response:\n%s", string(data))
	})

	t.Run("Filter by ID and Name", func(t *testing.T) {
		subID := 1
		subName := "TF_TEST_SUBSYSTEM" // Change to an existing name in your array
		
		input := reconmodel.NvmSubsystemGetMultipleInput{
			NvmSubsystemId:   &subID,
			NvmSubsystemName: &subName,
		}

		resp, err := psm.ReconcileGetMultipleNvmSubsystems(input)
		if err != nil {
			t.Fatalf("Filtering failed: %v", err)
		}

		for _, sub := range resp.Data {
			if sub.NvmSubsystemId != subID {
				t.Errorf("Expected ID %d, got %d", subID, sub.NvmSubsystemId)
			}
		}
		t.Logf("Filtered results: %d", len(resp.Data))
	})
}

// go test -v -run ^TestReconcileGetMultipleNvmSubsystems_IndividualHydration$
func ^TestReconcileGetMultipleNvmSubsystems_IndividualHydration(t *testing.T) {
	psm, _ := newNvmReconcilerTestManager()

	tests := []struct {
		name  string
		input reconmodel.NvmSubsystemGetMultipleInput
	}{
		{
			name: "Only Namespaces",
			input: reconmodel.NvmSubsystemGetMultipleInput{IncludeNamespaces: true},
		},
		{
			name: "Only Ports",
			input: reconmodel.NvmSubsystemGetMultipleInput{IncludePorts: true},
		},
		{
			name: "Only Host NQNs",
			input: reconmodel.NvmSubsystemGetMultipleInput{IncludeHostNqns: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := psm.ReconcileGetMultipleNvmSubsystems(tt.input)
			if err != nil {
				t.Fatalf("%s failed: %v", tt.name, err)
			}
			t.Logf("%s returned %d items", tt.name, len(resp.Data))
		})
	}
}

// go test -v -run ^TestNvmSubsystemFullLifecycle$
func ^TestNvmSubsystemFullLifecycle(t *testing.T) {
    psm, err := newNvmReconcilerTestManager()
    if err != nil {
        t.Fatalf("Failed to initialize manager: %v", err)
    }

    // Generate a unique ID and Name to avoid collisions
    targetID, _ := psm.getFirstAvailableNvmSubsystemId()
    subName := fmt.Sprintf("TEST_SUB_%d", targetID)
    hostNqn := "nqn.2014-08.org.nvmexpress:uuid:550e8400-e29b-41d4-a716-446655440000"
    ldevID := 304 // Ensure this LDEV exists and is unassigned in your lab

    // 1. CREATE
    t.Run("Step 1: Create Subsystem", func(t *testing.T) {
        input := reconmodel.NvmSubsystemReconcilerInput{
            NvmSubsystemName: &subName,
            HostMode:         utils.Ptr("LINUX/IRIX"),
            Ports:            &[]string{"CL1-A"},
            HostNqns: &[]reconmodel.HostNqnInput{
                {Nqn: hostNqn, Nickname: "TEST_HOST"},
            },
            Namespaces: &[]reconmodel.NamespaceInput{
                {LdevId: ldevID, Nickname: "TEST_VOL", Path: []string{hostNqn}},
            },
        }

        sub, err := psm.reconcileNvmSubsystemCreate(targetID, input)
        if err != nil {
            t.Fatalf("Create failed: %v", err)
        }
        t.Logf("Created Subsystem ID: %d", sub.NvmSubsystemId)
    })

    // 2. UPDATE (Add Port and Change Nickname)
    t.Run("Step 2: Update Subsystem Resources", func(t *testing.T) {
        // Fetch current state
        current, err := psm.ReconcileGetNvmSubsystem(targetID)
        if err != nil {
            t.Fatalf("Fetch current failed: %v", err)
        }

        updatedName := subName + "_REV"
        input := reconmodel.NvmSubsystemReconcilerInput{
            NvmSubsystemName: &updatedName,
            Ports:            &[]string{"CL1-A", "CL2-A"}, // Adding CL2-A
            Namespaces: &[]reconmodel.NamespaceInput{
                {LdevId: ldevID, Nickname: "UPDATED_VOL", Path: []string{hostNqn}},
            },
        }

        sub, err := psm.reconcileNvmSubsystemUpdate(targetID, input, current)
        if err != nil {
            t.Fatalf("Update failed: %v", err)
        }

        if len(sub.PortIds) != 2 {
            t.Errorf("Expected 2 ports, got %d", len(sub.PortIds))
        }
    })

    // 3. DELETE
    t.Run("Step 3: Delete Subsystem", func(t *testing.T) {
        err := psm.ReconcileNvmSubsystemDelete(targetID)
        if err != nil {
            t.Fatalf("Delete failed: %v", err)
        }

        // Verify deletion
        _, err = psm.ReconcileGetNvmSubsystem(targetID)
        if err == nil {
            t.Error("Expected 404/error after deletion, but subsystem still exists")
        } else {
            t.Log("Subsystem successfully deleted and verified.")
        }
    })
}