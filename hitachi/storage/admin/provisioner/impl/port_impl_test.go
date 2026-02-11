package admin

import (
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	"testing"
)

// go test -v -run TestGetPorts
func xTestGetPorts(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Scenario 1: Get all ports without protocol filter
	params := gwymodel.GetPortParams{}
	resp, err := psm.GetPorts(params)
	if err != nil {
		t.Errorf("Unexpected error in GetPorts (no filter) %v", err)
		return
	}
	t.Logf("Response (all ports): Count=%d, Ports=%+v", resp.Count, resp.Data)

	// Scenario 2: Get FC ports only
	protocol := "FC"
	params = gwymodel.GetPortParams{
		Protocol: &protocol,
	}
	resp, err = psm.GetPorts(params)
	if err != nil {
		t.Errorf("Unexpected error in GetPorts (FC filter) %v", err)
		return
	}
	t.Logf("Response (FC ports): Count=%d, Ports=%+v", resp.Count, resp.Data)

	// Scenario 3: Get iSCSI ports only
	protocol = "iSCSI"
	params = gwymodel.GetPortParams{
		Protocol: &protocol,
	}
	resp, err = psm.GetPorts(params)
	if err != nil {
		t.Errorf("Unexpected error in GetPorts (ISCSI filter) %v", err)
		return
	}
	t.Logf("Response (iSCSI ports): Count=%d, Ports=%+v", resp.Count, resp.Data)

	// Scenario 4: Get NVMe-TCP ports only
	protocol = "NVME_TCP"
	params = gwymodel.GetPortParams{
		Protocol: &protocol,
	}
	resp, err = psm.GetPorts(params)
	if err != nil {
		t.Errorf("Unexpected error in GetPorts (NVME_TCP filter) %v", err)
		return
	}
	t.Logf("Response (NVMe-TCP ports): Count=%d, Ports=%+v", resp.Count, resp.Data)
}

// go test -v -run TestGetPortByID
func xTestGetPortByID(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// First, get all ports to find a valid port ID for testing
	params := gwymodel.GetPortParams{}
	portsList, err := psm.GetPorts(params)
	if err != nil {
		t.Fatalf("Unexpected error getting ports list for setup: %v", err)
	}

	if len(portsList.Data) == 0 {
		t.Skip("No ports available for GetPortByID test")
		return
	}

	// Use the first port ID for testing
	testPortID := portsList.Data[0].ID

	// Scenario 1: Get specific port by valid ID
	resp, err := psm.GetPortByID(testPortID)
	if err != nil {
		t.Errorf("Unexpected error in GetPortByID (valid ID) %v", err)
		return
	}
	t.Logf("Response (port %s): %+v", testPortID, resp)

	// Validate the returned port ID matches what we requested
	if resp.ID != testPortID {
		t.Errorf("Expected port ID %s, got %s", testPortID, resp.ID)
	}

	// Scenario 2: Test with a non-existent port ID (this should fail)
	invalidPortID := "INVALID-PORT-ID-12345"
	_, err = psm.GetPortByID(invalidPortID)
	if err == nil {
		t.Logf("Note: GetPortByID with invalid ID '%s' did not return an error (this may be expected behavior)", invalidPortID)
	} else {
		t.Logf("Expected behavior: GetPortByID with invalid ID '%s' returned error: %v", invalidPortID, err)
	}
}

// go test -v -run TestGetPortsWithSpecificProtocolTypes
func xTestGetPortsWithSpecificProtocolTypes(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	protocolTypes := []string{"FC", "iSCSI", "NVME_TCP"}

	for _, protocolType := range protocolTypes {
		t.Run("Protocol_"+protocolType, func(t *testing.T) {
			protocol := protocolType
			params := gwymodel.GetPortParams{
				Protocol: &protocol,
			}

			resp, err := psm.GetPorts(params)
			if err != nil {
				t.Errorf("Unexpected error in GetPorts (%s filter) %v", protocolType, err)
				return
			}

			t.Logf("Protocol %s: Count=%d", protocolType, resp.Count)

			// Validate that all returned ports match the requested protocol
			for i, port := range resp.Data {
				if port.Protocol != protocolType {
					t.Errorf("Port %d: Expected protocol %s, got %s", i, protocolType, port.Protocol)
				}
			}
		})
	}
}

func xTestUpdatePort(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test port ID
	testPortID := "CL4-C"

	// First, get the current state of the port to verify it exists
	currentPort, err := psm.GetPortByID(testPortID)
	if err != nil {
		t.Fatalf("Failed to get port %s for test setup: %v", testPortID, err)
	}
	t.Logf("Current port state before update: %+v", currentPort)

	// Prepare update parameters to enable port security
	enablePortSecurity := true
	updateParams := gwymodel.UpdatePortParams{
		PortSecurity: &enablePortSecurity,
	}

	// Perform the update
	err = psm.UpdatePort(testPortID, updateParams)
	if err != nil {
		t.Errorf("Unexpected error in UpdatePort for port %s: %v", testPortID, err)
		return
	}
	t.Logf("Port %s update completed successfully", testPortID)

	// Verify the update by getting the port again
	updatedPort, err := psm.GetPortByID(testPortID)
	if err != nil {
		t.Errorf("Failed to get updated port %s: %v", testPortID, err)
		return
	}
	t.Logf("Port state after update: %+v", updatedPort)

	// Validate that port security is now enabled
	if !updatedPort.PortSecurity {
		t.Errorf("Expected port security to be enabled (true), but got %v", updatedPort.PortSecurity)
	} else {
		t.Logf("SUCCESS: Port security successfully enabled for port %s", testPortID)
	}
}
