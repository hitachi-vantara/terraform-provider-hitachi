package admin

import (
	"fmt"
	"testing"

	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
)

// go test -v -run TestGetServerHBAs
func xTestGetServerHBAs(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test server ID - you may need to adjust this based on your test environment
	testServerID := 23

	// Scenario 1: Get all HBAs for a specific server
	resp, err := psm.GetServerHBAs(testServerID)
	if err != nil {
		t.Errorf("Unexpected error in GetServerHBAs for serverID %d: %v", testServerID, err)
		return
	}

	// Validate response structure
	if resp == nil {
		t.Errorf("Expected non-nil response, got nil")
		return
	}

	t.Logf("Response (server %d HBAs): Count=%d, HBAs=%+v", testServerID, resp.Count, resp.Data)

	// Log details about each HBA
	for i, hba := range resp.Data {
		t.Logf("HBA %d: WWN=%s, IscsiName=%s, PortIds=%v", i, hba.HbaWwn, hba.IscsiName, hba.PortIds)
	}
}

// go test -v -run TestGetServerHBAByWwn
func xTestGetServerHBAByWwn(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test server ID
	testServerID := 23

	// First, get all HBAs for the server to find a valid WWN for testing
	hbasList, err := psm.GetServerHBAs(testServerID)
	if err != nil {
		t.Fatalf("Unexpected error getting HBAs list for setup: %v", err)
	}

	if len(hbasList.Data) == 0 {
		t.Skip("No HBAs available for GetServerHBAByWwn test")
		return
	}

	// Use the first HBA IscsiName for testing
	testIscsiName := hbasList.Data[0].IscsiName

	// Scenario 1: Get specific HBA by valid server ID and IscsiName
	resp, err := psm.GetServerHBAByWwn(testServerID, testIscsiName)
	if err != nil {
		t.Errorf("Unexpected error in GetServerHBAByWwn (valid serverID and IscsiName) %v", err)
		return
	}
	t.Logf("Response (server %d, HBA %s): %+v", testServerID, testIscsiName, resp)

	// Validate the returned HBA IscsiName matches what we requested
	if resp.IscsiName != testIscsiName {
		t.Errorf("Expected HBA IscsiName %s, got %s", testIscsiName, resp.IscsiName)
	}

	// Scenario 2: Test with a non-existent HBA WWN (this should fail)
	invalidHBAWwn := "INVALID-WWN-12345"
	_, err = psm.GetServerHBAByWwn(testServerID, invalidHBAWwn)
	if err == nil {
		t.Logf("Note: GetServerHBAByWwn with invalid WWN '%s' did not return an error (this may be expected behavior)", invalidHBAWwn)
	} else {
		t.Logf("Expected behavior: GetServerHBAByWwn with invalid WWN '%s' returned error: %v", invalidHBAWwn, err)
	}

	// Scenario 3: Test with a non-existent server ID
	invalidServerID := 99999
	_, err = psm.GetServerHBAByWwn(invalidServerID, testIscsiName)
	if err == nil {
		t.Logf("Note: GetServerHBAByWwn with invalid serverID %d did not return an error (this may be expected behavior)", invalidServerID)
	} else {
		t.Logf("Expected behavior: GetServerHBAByWwn with invalid serverID %d returned error: %v", invalidServerID, err)
	}
}

// go test -v -run TestGetServerHBAsWithMultipleServers
func xTestGetServerHBAsWithMultipleServers(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test with multiple server IDs - adjust these based on your test environment
	testServerIDs := []int{1, 2, 3}

	for _, serverID := range testServerIDs {
		t.Run(fmt.Sprintf("ServerID_%d", serverID), func(t *testing.T) {
			resp, err := psm.GetServerHBAs(serverID)
			if err != nil {
				t.Logf("Server %d: Error getting HBAs (may not exist): %v", serverID, err)
				return
			}

			t.Logf("Server %d: HBA Count=%d", serverID, resp.Count)

			// Validate HBA data structure for each server
			for i, hba := range resp.Data {
				if hba.HbaWwn == "" {
					t.Errorf("Server %d, HBA %d: Expected non-empty HBA WWN", serverID, i)
				}
				if hba.IscsiName == "" {
					t.Logf("Server %d, HBA %d: IscsiName is empty (may be expected for FC HBAs)", serverID, i)
				}
				t.Logf("Server %d, HBA %d: HbaWwn=%s, IscsiName=%s, PortIds=%v",
					serverID, i, hba.HbaWwn, hba.IscsiName, hba.PortIds)
			}
		})
	}
}

// go test -v -run TestCreateServerHBAWithWwn
func xTestCreateServerHBAWithWwn(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test server ID and HBA WWN
	testServerID := 51
	testHBAWwn := "10009440c9d053ac"

	params := gwymodel.CreateServerHBAParams{
		HBAs: []gwymodel.ServerHBARequest{
			{
				HbaWwn:    testHBAWwn,
				IscsiName: "",
			},
		},
	}

	t.Logf("Creating HBA with WWN: %s on server %d", testHBAWwn, testServerID)

	resp, err := psm.CreateServerHBAs(testServerID, params)
	if err != nil {
		t.Errorf("Unexpected error in CreateServerHBAs for WWN %s: %v", testHBAWwn, err)
		return
	}

	// Validate response structure
	if resp == nil {
		t.Errorf("Expected non-nil response, got nil")
		return
	}

	t.Logf("Response (server %d, HBA WWN %s): Count=%d, HBAs=%+v", testServerID, testHBAWwn, resp.Count, resp.Data)

	// Validate that at least one HBA was created
	if resp.Count == 0 {
		t.Errorf("Expected at least 1 HBA to be created, got count=%d", resp.Count)
	}
}

// go test -v -run TestCreateServerHBAWithISCSI
func xTestCreateServerHBAWithISCSI(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test server ID and iSCSI name
	testServerID := 20
	testISCSIName := "iqn.server1041"

	params := gwymodel.CreateServerHBAParams{
		HBAs: []gwymodel.ServerHBARequest{
			{
				HbaWwn:    "",
				IscsiName: testISCSIName,
			},
		},
	}

	t.Logf("Creating HBA with iSCSI name: %s on server %d", testISCSIName, testServerID)

	resp, err := psm.CreateServerHBAs(testServerID, params)
	if err != nil {
		t.Errorf("Unexpected error in CreateServerHBAs for iSCSI name %s: %v", testISCSIName, err)
		return
	}

	// Validate response structure
	if resp == nil {
		t.Errorf("Expected non-nil response, got nil")
		return
	}

	t.Logf("Response (server %d, iSCSI name %s): Count=%d, HBAs=%+v", testServerID, testISCSIName, resp.Count, resp.Data)

	// Validate that at least one HBA was created
	if resp.Count == 0 {
		t.Errorf("Expected at least 1 HBA to be created, got count=%d", resp.Count)
	}
}

// go test -v -run TestDeleteServerHBAByWwn
func xTestDeleteServerHBAByWwn(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test server ID - you may need to adjust this based on your test environment
	testServerID := 51
	testHBAWwn := "10009440c9d053ac"

	t.Logf("Attempting to delete HBA with WWN: %s from server %d", testHBAWwn, testServerID)

	// Delete the HBA
	resp, err := psm.DeleteServerHBA(testServerID, testHBAWwn)
	if err != nil {
		t.Errorf("Unexpected error in DeleteServerHBA for WWN %s: %v", testHBAWwn, err)
		return
	}

	// Validate response structure
	if resp == nil {
		t.Errorf("Expected non-nil response, got nil")
		return
	}

	t.Logf("Response after deleting HBA (server %d): Count=%d, HBAs=%+v", testServerID, resp.Count, resp.Data)

	// Validate that the HBA was actually deleted
	for _, hba := range resp.Data {
		if hba.HbaWwn == testHBAWwn {
			t.Errorf("Expected HBA with WWN %s to be deleted, but it still exists", testHBAWwn)
		}
	}

	t.Logf("Successfully deleted HBA with WWN: %s", testHBAWwn)
}

func xTestDeleteServerHBAByISCSI(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test server ID - you may need to adjust this based on your test environment
	testServerID := 20
	testISCSIName := "iqn.server1041"

	t.Logf("Attempting to delete HBA with iSCSI name: %s from server %d", testISCSIName, testServerID)

	// Delete the HBA
	resp, err := psm.DeleteServerHBA(testServerID, testISCSIName)
	if err != nil {
		t.Errorf("Unexpected error in DeleteServerHBA for iSCSI name %s: %v", testISCSIName, err)
		return
	}

	// Validate response structure
	if resp == nil {
		t.Errorf("Expected non-nil response, got nil")
		return
	}

	t.Logf("Response after deleting HBA (server %d): Count=%d, HBAs=%+v", testServerID, resp.Count, resp.Data)

	// Validate that the HBA was actually deleted
	for _, hba := range resp.Data {
		if hba.IscsiName == testISCSIName {
			t.Errorf("Expected HBA with iSCSI name %s to be deleted, but it still exists", testISCSIName)
		}
	}

	t.Logf("Successfully deleted HBA with iSCSI name: %s", testISCSIName)
}

// go test -v -run TestGetServerHBAsTypes
func xTestGetServerHBAsTypes(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	testServerID := 1

	resp, err := psm.GetServerHBAs(testServerID)
	if err != nil {
		t.Fatalf("Unexpected error getting HBAs for type analysis: %v", err)
	}

	if len(resp.Data) == 0 {
		t.Skip("No HBAs available for type analysis")
		return
	}

	// Count HBAs by type (FC vs iSCSI based on whether IscsiName is populated)
	fcCount := 0
	iscsiCount := 0

	for _, hba := range resp.Data {
		if hba.IscsiName != "" {
			iscsiCount++
		} else {
			fcCount++
		}
	}

	t.Logf("HBA type distribution for server %d:", testServerID)
	t.Logf("  FC HBAs (no iSCSI name): %d", fcCount)
	t.Logf("  iSCSI HBAs (with iSCSI name): %d", iscsiCount)

	// Test individual HBA retrieval for each type
	if fcCount > 0 {
		t.Run("Type_FC", func(t *testing.T) {
			// Find first FC HBA (no iSCSI name)
			for _, hba := range resp.Data {
				if hba.IscsiName == "" {
					// Test individual HBA retrieval
					retrievedHBA, err := psm.GetServerHBAByWwn(testServerID, hba.HbaWwn)
					if err != nil {
						t.Errorf("Error retrieving FC HBA %s: %v", hba.HbaWwn, err)
						return
					}

					// Validate consistency
					if retrievedHBA.HbaWwn != hba.HbaWwn {
						t.Errorf("HBA WWN mismatch: expected %s, got %s", hba.HbaWwn, retrievedHBA.HbaWwn)
					}
					if retrievedHBA.IscsiName != hba.IscsiName {
						t.Errorf("IscsiName mismatch for WWN %s: expected %s, got %s",
							hba.HbaWwn, hba.IscsiName, retrievedHBA.IscsiName)
					}

					t.Logf("Successfully validated FC HBA %s", hba.HbaWwn)
					break // Test only first FC HBA
				}
			}
		})
	}

	if iscsiCount > 0 {
		t.Run("Type_iSCSI", func(t *testing.T) {
			// Find first iSCSI HBA (with iSCSI name)
			for _, hba := range resp.Data {
				if hba.IscsiName != "" {
					// Test individual HBA retrieval
					retrievedHBA, err := psm.GetServerHBAByWwn(testServerID, hba.HbaWwn)
					if err != nil {
						t.Errorf("Error retrieving iSCSI HBA %s: %v", hba.HbaWwn, err)
						return
					}

					// Validate consistency
					if retrievedHBA.HbaWwn != hba.HbaWwn {
						t.Errorf("HBA WWN mismatch: expected %s, got %s", hba.HbaWwn, retrievedHBA.HbaWwn)
					}
					if retrievedHBA.IscsiName != hba.IscsiName {
						t.Errorf("IscsiName mismatch for WWN %s: expected %s, got %s",
							hba.HbaWwn, hba.IscsiName, retrievedHBA.IscsiName)
					}

					t.Logf("Successfully validated iSCSI HBA %s with name %s", hba.HbaWwn, hba.IscsiName)
					break // Test only first iSCSI HBA
				}
			}
		})
	}
}

// go test -v -run TestGetServerHBAErrorHandling
func xTestGetServerHBAErrorHandling(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test various error scenarios
	errorTestCases := []struct {
		name     string
		serverID int
		hbaWwn   string
		testType string
	}{
		{"InvalidServerID_Negative", -1, "", "GetServerHBAs"},
		{"InvalidServerID_Zero", 0, "", "GetServerHBAs"},
		{"InvalidServerID_Large", 999999, "", "GetServerHBAs"},
		{"InvalidWWN_Empty", 1, "", "GetServerHBAByWwn"},
		{"InvalidWWN_Invalid", 1, "INVALID-WWN-FORMAT", "GetServerHBAByWwn"},
		{"InvalidWWN_NonExistent", 1, "10:00:00:00:00:00:00:FF", "GetServerHBAByWwn"},
	}

	for _, tc := range errorTestCases {
		t.Run(tc.name, func(t *testing.T) {
			switch tc.testType {
			case "GetServerHBAs":
				_, err := psm.GetServerHBAs(tc.serverID)
				if err != nil {
					t.Logf("Expected error for %s with serverID %d: %v", tc.name, tc.serverID, err)
				} else {
					t.Logf("No error returned for %s with serverID %d (may be valid behavior)", tc.name, tc.serverID)
				}
			case "GetServerHBAByWwn":
				_, err := psm.GetServerHBAByWwn(tc.serverID, tc.hbaWwn)
				if err != nil {
					t.Logf("Expected error for %s with serverID %d, WWN %s: %v", tc.name, tc.serverID, tc.hbaWwn, err)
				} else {
					t.Logf("No error returned for %s with serverID %d, WWN %s (may be valid behavior)", tc.name, tc.serverID, tc.hbaWwn)
				}
			}
		})
	}
}
