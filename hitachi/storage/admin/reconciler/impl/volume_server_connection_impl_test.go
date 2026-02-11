package admin

import (
	"encoding/json"
	"reflect"
	"testing"

	recmodel "terraform-provider-hitachi/hitachi/storage/admin/reconciler/model"
)

// go test -v -run TestReconcileReadVolumeServerConnections
func xTestReconcileReadVolumeServerConnections(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	// 🔹 Replace with valid VolumeID–ServerID pairs in your backend
	pairs := []recmodel.VolumeServerPair{
		{VolumeID: 6079, ServerID: 9},
		{VolumeID: 6079, ServerID: 10},
	}

	results, existing, err := psm.ReconcileReadVolumeServerConnections(pairs)
	if err != nil {
		t.Errorf("Unexpected error in ReconcileReadVolumeServerConnections: %v", err)
		return
	}

	if len(existing) == 0 {
		t.Logf("No existing connections found for provided pairs")
	} else {
		b, _ := json.MarshalIndent(results, "", "  ")
		t.Logf("Existing connections (%d): %v\nDetails:\n%s", len(existing), existing, string(b))
	}
}

// go test -v -run TestReconcileDeleteVolumeServerConnections
func xTestReconcileDeleteVolumeServerConnections(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	// 🔹 Replace with real pairs that exist in your backend
	pairs := []recmodel.VolumeServerPair{
		{VolumeID: 6079, ServerID: 9},
		{VolumeID: 6079, ServerID: 10},
	}

	t.Logf("Attempting to delete volume-server connections: %+v", pairs)
	err = psm.ReconcileDeleteVolumeServerConnections(pairs)
	if err != nil {
		t.Errorf("Unexpected error in ReconcileDeleteVolumeServerConnections: %v", err)
		return
	}

	t.Logf("ReconcileDeleteVolumeServerConnections succeeded for: %+v", pairs)
}

// go test -v -run TestReconcileUpdateVolumeServerConnections
func xTestReconcileUpdateVolumeServerConnections(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	// 🔹 Replace with valid IDs in your system
	existingPairs := []recmodel.VolumeServerPair{
		{VolumeID: 6079, ServerID: 9},
		{VolumeID: 6079, ServerID: 10},
	}
	desiredPairs := []recmodel.VolumeServerPair{
		{VolumeID: 6079, ServerID: 9}, // kept
		// {VolumeID: 6079, ServerID: 10}, // removed
		{VolumeID: 6078, ServerID: 10}, // added
	}

	existingJSON, _ := json.MarshalIndent(existingPairs, "", "  ")
	desiredJSON, _ := json.MarshalIndent(desiredPairs, "", "  ")
	t.Logf("Existing pairs:\n%s\nDesired pairs:\n%s", string(existingJSON), string(desiredJSON))

	err = psm.ReconcileUpdateVolumeServerConnections(existingPairs, desiredPairs)
	if err != nil {
		t.Errorf("Unexpected error in ReconcileUpdateVolumeServerConnections: %v", err)
		return
	}

	t.Logf("ReconcileUpdateVolumeServerConnections completed successfully")
}

// go test -v -run TestDiffVolumeServerPairsRemoved
func xTestDiffVolumeServerPairsRemoved(t *testing.T) {
	a := []recmodel.VolumeServerPair{
		{VolumeID: 1, ServerID: 1}, // deleted
		{VolumeID: 2, ServerID: 2},
		{VolumeID: 3, ServerID: 3}, //deleted
	}
	b := []recmodel.VolumeServerPair{
		{VolumeID: 2, ServerID: 2},
		{VolumeID: 4, ServerID: 4}, // added
	}

	expected := []recmodel.VolumeServerPair{
		{VolumeID: 1, ServerID: 1},
		{VolumeID: 3, ServerID: 3},
	}
	got := diffVolumeServerPairs(a, b)

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}

// go test -v -run TestDiffVolumeServerPairsAdded
func xTestDiffVolumeServerPairsAdded(t *testing.T) {
	a := []recmodel.VolumeServerPair{
		{VolumeID: 1, ServerID: 1}, // deleted
		{VolumeID: 2, ServerID: 2},
		{VolumeID: 3, ServerID: 3}, //deleted
	}
	b := []recmodel.VolumeServerPair{
		{VolumeID: 2, ServerID: 2},
		{VolumeID: 4, ServerID: 4}, // added
	}

	expected := []recmodel.VolumeServerPair{
		{VolumeID: 4, ServerID: 4},
	}
	got := diffVolumeServerPairs(b, a)

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}
