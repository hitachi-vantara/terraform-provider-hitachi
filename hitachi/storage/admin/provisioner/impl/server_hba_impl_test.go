package admin

import (
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	"testing"
)

// go test -v -run TestGetServerHBAs
func xTestGetServerHBAs(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test with a valid server ID (replace with an actual server ID from your environment)
	testServerID := 1

	// Scenario 1: Get all HBAs for a server
	resp, err := psm.GetServerHBAs(testServerID)
	if err != nil {
		t.Errorf("Unexpected error in GetServerHBAs %v", err)
		return
	}
	t.Logf("Response (server %d HBAs): Count=%d, HBAs=%+v", testServerID, resp.Count, resp.Data)

	// Validate response structure
	if resp.Count != len(resp.Data) {
		t.Errorf("Count mismatch: expected %d, got %d", len(resp.Data), resp.Count)
	}

	// Log individual HBA details
	for i, hba := range resp.Data {
		t.Logf("HBA %d: ServerID=%d, HbaWwn=%s, IscsiName=%s, PortIds=%v",
			i, hba.ServerID, hba.HbaWwn, hba.IscsiName, hba.PortIds)
	}
}

// go test -v -run TestGetServerHBAByWwn
func xTestGetServerHBAByWwn(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	testServerID := 51

	// First, get all HBAs to find a valid WWN for testing
	hbasList, err := psm.GetServerHBAs(testServerID)
	if err != nil {
		t.Fatalf("Unexpected error getting HBAs list for setup: %v", err)
	}

	if len(hbasList.Data) == 0 {
		t.Skip("No HBAs available for GetServerHBAByWwn test")
		return
	}

	// Find an HBA with a valid WWN
	var testHbaWwn string
	for _, hba := range hbasList.Data {
		if hba.HbaWwn != "" {
			testHbaWwn = hba.HbaWwn
			break
		}
	}

	if testHbaWwn == "" {
		t.Skip("No HBAs with valid WWN found for GetServerHBAByWwn test")
		return
	}

	// Scenario 1: Get specific HBA by valid WWN
	resp, err := psm.GetServerHBAByWwn(testServerID, testHbaWwn)
	if err != nil {
		t.Errorf("Unexpected error in GetServerHBAByWwn (valid WWN) %v", err)
		return
	}
	t.Logf("Response (server %d, WWN %s): %+v", testServerID, testHbaWwn, resp)

	// Validate the returned HBA WWN matches what we requested
	if resp.HbaWwn != testHbaWwn {
		t.Errorf("Expected HBA WWN %s, got %s", testHbaWwn, resp.HbaWwn)
	}

	// Validate the server ID matches
	if resp.ServerID != testServerID {
		t.Errorf("Expected server ID %d, got %d", testServerID, resp.ServerID)
	}

	// Scenario 2: Test with a non-existent WWN (this should fail)
	invalidWwn := "20:00:00:00:00:00:00:99"
	_, err = psm.GetServerHBAByWwn(testServerID, invalidWwn)
	if err == nil {
		t.Logf("Note: GetServerHBAByWwn with invalid WWN '%s' did not return an error (this may be expected behavior)", invalidWwn)
	} else {
		t.Logf("Expected behavior: GetServerHBAByWwn with invalid WWN '%s' returned error: %v", invalidWwn, err)
	}
}

// go test -v -run TestCreateServerHBAs
func xTestCreateServerHBAs(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	testServerID := 51

	// Scenario 1: Create HBAs with FC WWN
	fcHBAs := []gwymodel.ServerHBARequest{
		{
			HbaWwn: "210000e08b050504",
		},
		{
			HbaWwn: "210000e08b050505",
		},
	}

	createParams := gwymodel.CreateServerHBAParams{
		HBAs: fcHBAs,
	}

	resp, err := psm.CreateServerHBAs(testServerID, createParams)
	if err != nil {
		t.Errorf("Unexpected error in CreateServerHBAs (FC WWNs) %v", err)
		return
	}
	t.Logf("Response (create FC HBAs): Count=%d, HBAs=%+v", resp.Count, resp.Data)

	// Validate that HBAs were created
	if resp.Count == 0 {
		t.Error("Expected at least one HBA to be created, but count is 0")
	}

	// Scenario 2: Create HBAs with iSCSI names
	iscsiHBAs := []gwymodel.ServerHBARequest{
		{
			IscsiName: "iqn.1991-05.com.microsoft:test-server-01",
		},
		{
			IscsiName: "iqn.1991-05.com.microsoft:test-server-02",
		},
	}

	createParamsIscsi := gwymodel.CreateServerHBAParams{
		HBAs: iscsiHBAs,
	}

	respIscsi, err := psm.CreateServerHBAs(testServerID, createParamsIscsi)
	if err != nil {
		t.Errorf("Unexpected error in CreateServerHBAs (iSCSI names) %v", err)
		return
	}
	t.Logf("Response (create iSCSI HBAs): Count=%d, HBAs=%+v", respIscsi.Count, respIscsi.Data)

	// Validate that iSCSI HBAs were created
	if respIscsi.Count == 0 {
		t.Error("Expected at least one iSCSI HBA to be created, but count is 0")
	}

	// Scenario 3: Test with empty HBAs list (should fail or return empty)
	emptyParams := gwymodel.CreateServerHBAParams{
		HBAs: []gwymodel.ServerHBARequest{},
	}

	_, err = psm.CreateServerHBAs(testServerID, emptyParams)
	if err == nil {
		t.Logf("Note: CreateServerHBAs with empty HBAs list did not return an error")
	} else {
		t.Logf("Expected behavior: CreateServerHBAs with empty HBAs list returned error: %v", err)
	}
}

// go test -v -run TestDeleteServerHBA
func xTestDeleteServerHBA(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	testServerID := 1

	// First, get current HBAs to find one to delete
	hbasList, err := psm.GetServerHBAs(testServerID)
	if err != nil {
		t.Fatalf("Unexpected error getting HBAs list for setup: %v", err)
	}

	if len(hbasList.Data) == 0 {
		t.Skip("No HBAs available for DeleteServerHBA test")
		return
	}

	// Find an HBA with a valid initiator name (WWN or iSCSI name)
	var testInitiatorName string
	for _, hba := range hbasList.Data {
		if hba.HbaWwn != "" {
			testInitiatorName = hba.HbaWwn
			break
		}
		if hba.IscsiName != "" {
			testInitiatorName = hba.IscsiName
			break
		}
	}

	if testInitiatorName == "" {
		t.Skip("No HBAs with valid initiator name found for DeleteServerHBA test")
		return
	}

	// Get initial count
	initialCount := hbasList.Count
	t.Logf("Initial HBA count for server %d: %d", testServerID, initialCount)

	// Scenario 1: Delete HBA by initiator name
	resp, err := psm.DeleteServerHBA(testServerID, testInitiatorName)
	if err != nil {
		t.Errorf("Unexpected error in DeleteServerHBA %v", err)
		return
	}
	t.Logf("Response (delete HBA %s): Count=%d, remaining HBAs=%+v", testInitiatorName, resp.Count, resp.Data)

	// Validate that HBA count decreased
	if resp.Count >= initialCount {
		t.Errorf("Expected HBA count to decrease from %d, but got %d", initialCount, resp.Count)
	} else {
		t.Logf("SUCCESS: HBA count decreased from %d to %d", initialCount, resp.Count)
	}

	// Scenario 2: Try to delete a non-existent HBA (should fail)
	invalidInitiator := "20:00:00:00:00:00:00:99"
	_, err = psm.DeleteServerHBA(testServerID, invalidInitiator)
	if err == nil {
		t.Logf("Note: DeleteServerHBA with invalid initiator '%s' did not return an error (this may be expected behavior)", invalidInitiator)
	} else {
		t.Logf("Expected behavior: DeleteServerHBA with invalid initiator '%s' returned error: %v", invalidInitiator, err)
	}
}

// go test -v -run TestServerHBAWorkflow
func xTestServerHBAWorkflow(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	testServerID := 1

	// Step 1: Get initial state
	initialHBAs, err := psm.GetServerHBAs(testServerID)
	if err != nil {
		t.Fatalf("Failed to get initial HBAs: %v", err)
	}
	t.Logf("Initial HBA count: %d", initialHBAs.Count)

	// Step 2: Create new HBAs
	testWwn := "21:00:00:e0:8b:99:99:99"
	createParams := gwymodel.CreateServerHBAParams{
		HBAs: []gwymodel.ServerHBARequest{
			{HbaWwn: testWwn},
		},
	}

	createdHBAs, err := psm.CreateServerHBAs(testServerID, createParams)
	if err != nil {
		t.Errorf("Failed to create HBAs: %v", err)
		return
	}
	t.Logf("Created HBAs count: %d", createdHBAs.Count)

	// Step 3: Verify the HBA was created by getting it specifically
	retrievedHBA, err := psm.GetServerHBAByWwn(testServerID, testWwn)
	if err != nil {
		t.Errorf("Failed to retrieve created HBA: %v", err)
	} else {
		t.Logf("Successfully retrieved created HBA: %+v", retrievedHBA)
		if retrievedHBA.HbaWwn != testWwn {
			t.Errorf("Retrieved HBA WWN mismatch: expected %s, got %s", testWwn, retrievedHBA.HbaWwn)
		}
	}

	// Step 4: Delete the created HBA
	deletedHBAs, err := psm.DeleteServerHBA(testServerID, testWwn)
	if err != nil {
		t.Errorf("Failed to delete HBA: %v", err)
		return
	}
	t.Logf("HBAs after deletion count: %d", deletedHBAs.Count)

	// Step 5: Verify final state matches initial state
	finalHBAs, err := psm.GetServerHBAs(testServerID)
	if err != nil {
		t.Errorf("Failed to get final HBAs: %v", err)
		return
	}
	t.Logf("Final HBA count: %d", finalHBAs.Count)

	if finalHBAs.Count != initialHBAs.Count {
		t.Errorf("Final HBA count (%d) does not match initial count (%d)", finalHBAs.Count, initialHBAs.Count)
	} else {
		t.Logf("SUCCESS: Workflow completed successfully, HBA count restored to initial state")
	}
}

// go test -v -run TestServerHBAErrorCases
func xTestServerHBAErrorCases(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Test with invalid server ID
	invalidServerID := -1

	t.Run("GetServerHBAs_InvalidServerID", func(t *testing.T) {
		_, err := psm.GetServerHBAs(invalidServerID)
		if err == nil {
			t.Logf("Note: GetServerHBAs with invalid server ID %d did not return an error", invalidServerID)
		} else {
			t.Logf("Expected behavior: GetServerHBAs with invalid server ID %d returned error: %v", invalidServerID, err)
		}
	})

	t.Run("GetServerHBAByWwn_InvalidServerID", func(t *testing.T) {
		_, err := psm.GetServerHBAByWwn(invalidServerID, "21:00:00:e0:8b:05:05:04")
		if err == nil {
			t.Logf("Note: GetServerHBAByWwn with invalid server ID %d did not return an error", invalidServerID)
		} else {
			t.Logf("Expected behavior: GetServerHBAByWwn with invalid server ID %d returned error: %v", invalidServerID, err)
		}
	})

	t.Run("CreateServerHBAs_InvalidServerID", func(t *testing.T) {
		params := gwymodel.CreateServerHBAParams{
			HBAs: []gwymodel.ServerHBARequest{
				{HbaWwn: "21:00:00:e0:8b:05:05:04"},
			},
		}
		_, err := psm.CreateServerHBAs(invalidServerID, params)
		if err == nil {
			t.Logf("Note: CreateServerHBAs with invalid server ID %d did not return an error", invalidServerID)
		} else {
			t.Logf("Expected behavior: CreateServerHBAs with invalid server ID %d returned error: %v", invalidServerID, err)
		}
	})

	t.Run("DeleteServerHBA_InvalidServerID", func(t *testing.T) {
		_, err := psm.DeleteServerHBA(invalidServerID, "21:00:00:e0:8b:05:05:04")
		if err == nil {
			t.Logf("Note: DeleteServerHBA with invalid server ID %d did not return an error", invalidServerID)
		} else {
			t.Logf("Expected behavior: DeleteServerHBA with invalid server ID %d returned error: %v", invalidServerID, err)
		}
	})
}
