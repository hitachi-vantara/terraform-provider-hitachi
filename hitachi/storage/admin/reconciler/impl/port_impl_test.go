package admin

import (
	"encoding/json"
	"testing"

	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
)

// go test -v -run TestReconcileReadAdminPort
func xTestReconcileReadAdminPort(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	// Test with a valid port ID (adjust to actual test data)
	portID := "CL4-C" // replace with valid port ID in your system

	portInfo, err := psm.ReconcileReadAdminPort(portID)
	if err != nil {
		t.Errorf("Unexpected error in ReconcileReadAdminPort: %v", err)
		return
	}

	if portInfo == nil {
		t.Errorf("Expected port information but got nil")
		return
	}

	b, _ := json.MarshalIndent(portInfo, "", "  ")
	t.Logf("Port ID: %s\nPort Info:\n%s", portID, string(b))
}

// go test -v -run TestReconcileUpdateAdminPort
func xTestReconcileUpdateAdminPort(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	tests := []struct {
		name       string
		portID     string
		updateFunc func() gwymodel.UpdatePortParams
	}{
		{
			name:   "Update Port Security",
			portID: "CL4-C",
			updateFunc: func() gwymodel.UpdatePortParams {
				enabled := true
				return gwymodel.UpdatePortParams{
					PortSecurity: &enabled,
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

			err := psm.ReconcileUpdateAdminPort(tt.portID, params)
			if err != nil {
				t.Errorf("Unexpected error in %s: %v", tt.name, err)
				return
			}

			t.Logf("%s succeeded for Port ID %s", tt.name, tt.portID)
		})
	}
}

// go test -v -run TestReconcileReadAdminPortWithDifferentPorts
func xTestReconcileReadAdminPortWithDifferentPorts(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	// Test with different port IDs (adjust to actual test data)
	portIDs := []string{"CL4-C", "CL1-A"} // replace with valid port IDs in your system

	for _, portID := range portIDs {
		t.Run("ReadPort_"+portID, func(t *testing.T) {
			portInfo, err := psm.ReconcileReadAdminPort(portID)
			if err != nil {
				t.Logf("Port %s not available or error: %v", portID, err)
				return
			}

			if portInfo == nil {
				t.Logf("Port %s returned nil information", portID)
				return
			}

			t.Logf("Port %s - Protocol: %s, Speed: %s, Security: %t",
				portID, portInfo.Protocol, portInfo.PortSpeed, portInfo.PortSecurity)
		})
	}
}

// go test -v -run TestReconcileUpdateAdminPortWithValidation
func xTestReconcileUpdateAdminPortWithValidation(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	portID := "CL4-C" // replace with valid port ID in your system

	// First, read the current port information
	currentPort, err := psm.ReconcileReadAdminPort(portID)
	if err != nil {
		t.Fatalf("Failed to read port %s for test setup: %v", portID, err)
	}
	t.Logf("Current port state before update: Security=%t, Speed=%s",
		currentPort.PortSecurity, currentPort.PortSpeed)

	// Update port security
	enableSecurity := true
	updateParams := gwymodel.UpdatePortParams{
		PortSecurity: &enableSecurity,
	}

	err = psm.ReconcileUpdateAdminPort(portID, updateParams)
	if err != nil {
		t.Errorf("Unexpected error in ReconcileUpdateAdminPort: %v", err)
		return
	}

	// Verify the update by reading the port again
	updatedPort, err := psm.ReconcileReadAdminPort(portID)
	if err != nil {
		t.Errorf("Failed to read updated port %s: %v", portID, err)
		return
	}

	t.Logf("Port state after update: Security=%t, Speed=%s",
		updatedPort.PortSecurity, updatedPort.PortSpeed)

	// Validate that port security is now enabled
	if !updatedPort.PortSecurity {
		t.Errorf("Expected port security to be enabled (true), but got %v", updatedPort.PortSecurity)
	} else {
		t.Logf("SUCCESS: Port security successfully enabled for port %s", portID)
	}
}
