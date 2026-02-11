package sanstorage

import (
	"encoding/json"
	"testing"

	model "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
)

// go test -v -run ^TestGetSnapshotGroups$
func xTestGetSnapshotGroups(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	params := model.GetSnapshotGroupsParams{}
	resp, err := psm.GetSnapshotGroups(params)
	if err != nil {
		t.Errorf("Unexpected error in GetSnapshotGroups: %v", err)
		return
	}
	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Logf("Failed to marshal SnapshotGroup: %v", err)
		return
	}

	t.Logf("Response:\n%s", string(data))
}

// go test -v -run ^TestGetSnapshotGroup$
func xTestGetSnapshotGroup(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	groupID := "SnapshotGroupTest"
	detailType := "retention"
	params := model.GetSnapshotGroupsParams{
		DetailInfoType: &detailType,
	}
	resp, err := psm.GetSnapshotGroup(groupID, params)
	if err != nil {
		t.Errorf("Unexpected error in GetSnapshotGroup: %v", err)
		return
	}
	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Logf("Failed to marshal SnapshotGroup: %v", err)
		return
	}

	t.Logf("Response:\n%s", string(data))
}

// go test -v -run TestSplitSnapshotGroup
func xTestSplitSnapshotGroup(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	groupID := "testGroup"
	request := model.SplitSnapshotRequest{}
	resp, err := psm.SplitSnapshotGroup(groupID, request)
	if err != nil {
		t.Errorf("Unexpected error in SplitSnapshotGroup: %v", err)
		return
	}
	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Logf("Failed to marshal SnapshotGroup: %v", err)
		return
	}

	t.Logf("Response:\n%s", string(data))
}

// go test -v -run TestResyncSnapshotGroup
func xTestResyncSnapshotGroup(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	groupID := "testGroup"
	request := model.ResyncSnapshotRequest{}
	resp, err := psm.ResyncSnapshotGroup(groupID, request)
	if err != nil {
		t.Errorf("Unexpected error in ResyncSnapshotGroup: %v", err)
		return
	}
	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Logf("Failed to marshal SnapshotGroup: %v", err)
		return
	}

	t.Logf("Response:\n%s", string(data))
}

// go test -v -run TestRestoreSnapshotGroup
func xTestRestoreSnapshotGroup(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	groupID := "testGroup"
	request := model.RestoreSnapshotRequest{}
	resp, err := psm.RestoreSnapshotGroup(groupID, request)
	if err != nil {
		t.Errorf("Unexpected error in RestoreSnapshotGroup: %v", err)
		return
	}
	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Logf("Failed to marshal SnapshotGroup: %v", err)
		return
	}

	t.Logf("Response:\n%s", string(data))
}

// go test -v -run TestDeleteSnapshotGroup
func xTestDeleteSnapshotGroup(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	groupID := "testGroup"
	resp, err := psm.DeleteSnapshotGroup(groupID)
	if err != nil {
		t.Errorf("Unexpected error in DeleteSnapshotGroup: %v", err)
		return
	}
	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Logf("Failed to marshal SnapshotGroup: %v", err)
		return
	}

	t.Logf("Response:\n%s", string(data))
}

// go test -v -run TestCloneSnapshotGroup
func xTestCloneSnapshotGroup(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	groupID := "testGroup"
	request := model.CloneSnapshotRequest{}
	resp, err := psm.CloneSnapshotGroup(groupID, request)
	if err != nil {
		t.Errorf("Unexpected error in CloneSnapshotGroup: %v", err)
		return
	}
	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Logf("Failed to marshal SnapshotGroup: %v", err)
		return
	}

	t.Logf("Response:\n%s", string(data))
}

// go test -v -run TestDeleteSnapshotTree
func xTestDeleteSnapshotTree(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	request := model.DeleteSnapshotTreeRequest{
		Parameters: model.DeleteSnapshotTreeParams{
			LdevID: 100,
		},
	}
	resIds, err := psm.DeleteSnapshotTree(request)
	if err != nil {
		t.Errorf("Unexpected error in DeleteSnapshotTree: %v", err)
		return
	}
	t.Logf("Successfully initiated DeleteSnapshotGroup. %v", resIds)
}

// go test -v -run TestDeleteGarbageData
func xTestDeleteGarbageData(t *testing.T) {
	psm, err := newSnapshotTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	request := model.DeleteGarbageDataRequest{
		Parameters: model.DeleteGarbageDataParams{
			LdevID:        100,
			OperationType: "all",
		},
	}
	resIds, err := psm.DeleteGarbageData(request)
	if err != nil {
		t.Errorf("Unexpected error in DeleteGarbageData: %v", err)
		return
	}
	t.Logf("Successfully initiated DeleteGarbageData. %v", resIds)
}
