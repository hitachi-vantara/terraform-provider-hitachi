package sanstorage

import (
	"fmt"
	"net/url"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/san/gateway/http"
	model "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
)

// GetSnapshots retrieves a list of Thin Image pair information based on filtering parameters.
func (psm *sanStorageManager) GetSnapshots(params model.GetSnapshotsParams) (*model.SnapshotListResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var pairResponse model.SnapshotListResponse

	q := url.Values{}

	if params.SnapshotGroupName != nil {
		q.Add("snapshotGroupName", *params.SnapshotGroupName)
	}
	if params.PvolLdevID != nil {
		q.Add("pvolLdevId", fmt.Sprintf("%d", *params.PvolLdevID))
	}
	if params.SvolLdevID != nil {
		q.Add("svolLdevId", fmt.Sprintf("%d", *params.SvolLdevID))
	}
	if params.MuNumber != nil {
		q.Add("muNumber", fmt.Sprintf("%d", *params.MuNumber))
	}
	if params.DetailInfoType != nil {
		q.Add("detailInfoType", *params.DetailInfoType)
	}

	log.WriteDebug("TFDebug| QueryParams:%+v", q)

	apiSuf := fmt.Sprintf("objects/snapshots?%s", q.Encode())
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &pairResponse)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	return &pairResponse, nil
}

// GetSnapshot retrieves detailed information for a specific Thin Image pair
func (psm *sanStorageManager) GetSnapshot(pvolLdevID int, muNumber int) (*model.Snapshot, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var pairInfo model.Snapshot

	apiSuf := fmt.Sprintf("objects/snapshots/%d,%d", pvolLdevID, muNumber)

	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &pairInfo)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	return &pairInfo, nil
}

// GetSnapshotReplicationsRange retrieves information about all Thin Image pairs
// This API can be used when the storage system is the VSP 5000 series.
func (psm *sanStorageManager) GetSnapshotReplicationsRange(params model.GetSnapshotReplicationsRangeParams) (*model.SnapshotAllListResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var pairResponse model.SnapshotAllListResponse

	q := url.Values{}

	if params.StartPvolLdevID != nil {
		q.Add("startPvolLdevId", fmt.Sprintf("%d", *params.StartPvolLdevID))
	}
	if params.EndPvolLdevID != nil {
		q.Add("endPvolLdevId", fmt.Sprintf("%d", *params.EndPvolLdevID))
	}
	log.WriteDebug("TFDebug| QueryParams:%+v", q)

	apiSuf := fmt.Sprintf("objects/snapshot-replications?%s", q.Encode())

	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &pairResponse)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	return &pairResponse, nil
}

// CreateSnapshot initiates the creation of one Thin Image pairs.
func (psm *sanStorageManager) CreateSnapshot(request model.CreateSnapshotParams) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := "objects/snapshots"

	resIds, err := httpmethod.PostCall(psm.storageSetting, apiSuf, request)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to create thin image pair: %w", err))
		return "", err
	}

	log.WriteDebug("TFDebug | Resource IDs = %v", *resIds)
	return *resIds, nil
}

// SplitSnapshot initiates the split of a specific Thin Image pair to store snapshot data.
func (psm *sanStorageManager) SplitSnapshot(pvolLdevID int, muNumber int, request model.SplitSnapshotRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/snapshots/%d,%d/actions/split/invoke", pvolLdevID, muNumber)

	resIds, err := httpmethod.PostCall(psm.storageSetting, apiSuf, request)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to split thin image pair %d,%d: %w", pvolLdevID, muNumber, err))
		return "", err
	}

	log.WriteDebug("TFDebug | Resource IDs = %v", *resIds)
	return *resIds, nil
}

// ResyncSnapshot resynchronizes a specific Thin Image pair and deletes snapshot data.
func (psm *sanStorageManager) ResyncSnapshot(pvolLdevID int, muNumber int, request model.ResyncSnapshotRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/snapshots/%d,%d/actions/resync/invoke", pvolLdevID, muNumber)

	resIds, err := httpmethod.PostCall(psm.storageSetting, apiSuf, request)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to resync thin image pair %d,%d: %w", pvolLdevID, muNumber, err))
		return "", err
	}

	log.WriteDebug("TFDebug | Resource IDs = %v", *resIds)
	return *resIds, nil
}

// RestoreSnapshot initiates the restoration of data from a snapshot to the P-VOL.
func (psm *sanStorageManager) RestoreSnapshot(pvolLdevID int, muNumber int, request model.RestoreSnapshotRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/snapshots/%d,%d/actions/restore/invoke", pvolLdevID, muNumber)

	resIds, err := httpmethod.PostCall(psm.storageSetting, apiSuf, request)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to restore thin image pair %d,%d: %w", pvolLdevID, muNumber, err))
		return "", err
	}

	log.WriteDebug("TFDebug | Resource IDs = %v", *resIds)
	return *resIds, nil
}

// DeleteSnapshot initiates the deletion of a specific Thin Image pair.
func (psm *sanStorageManager) DeleteSnapshot(pvolLdevID int, muNumber int) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/snapshots/%d,%d", pvolLdevID, muNumber)

	resIds, err := httpmethod.DeleteCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to delete thin image pair %d,%d: %w", pvolLdevID, muNumber, err))
		return "", err
	}

	log.WriteDebug("TFDebug | Resource IDs = %v", *resIds)
	return *resIds, nil
}

// CloneSnapshot initiates the cloning of a specific Thin Image pair.
func (psm *sanStorageManager) CloneSnapshot(pvolLdevID int, muNumber int, request model.CloneSnapshotRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/snapshots/%d,%d/actions/clone/invoke", pvolLdevID, muNumber)

	resIds, err := httpmethod.PostCall(psm.storageSetting, apiSuf, request)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to clone thin image pair %d,%d: %w", pvolLdevID, muNumber, err))
		return "", err
	}

	log.WriteDebug("TFDebug | cloneSnapshotJobID = %v", *resIds)
	return *resIds, nil
}

// AssignSnapshotVolume assigns a secondary volume to snapshot data of a Thin Image pair.
func (psm *sanStorageManager) AssignSnapshotVolume(pvolLdevID int, muNumber int, request model.AssignSnapshotVolumeRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/snapshots/%d,%d/actions/assign-volume/invoke", pvolLdevID, muNumber)

	// Execute the POST call
	resIds, err := httpmethod.PostCall(psm.storageSetting, apiSuf, request)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to assign volume to thin image pair %d,%d: %w", pvolLdevID, muNumber, err))
		return "", err
	}

	log.WriteDebug("TFDebug | assignSnapshotVolumeJobID = %v", *resIds)
	return *resIds, nil
}

// UnassignSnapshotVolume unassigns the secondary volume for the snapshot data of a Thin Image pair.
func (psm *sanStorageManager) UnassignSnapshotVolume(pvolLdevID int, muNumber int) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/snapshots/%d,%d/actions/unassign-volume/invoke", pvolLdevID, muNumber)

	resIds, err := httpmethod.PostCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to unassign volume for thin image pair %d,%d: %w", pvolLdevID, muNumber, err))
		return "", err
	}

	log.WriteDebug("TFDebug | unassignSnapshotVolumeJobID = %v", *resIds)
	return *resIds, nil
}

// SetSnapshotRetentionPeriod sets the snapshot data retention period for the snapshot data of the Thin Image Advanced pair.
// You can also update (extend) a retention period that is already set.
func (psm *sanStorageManager) SetSnapshotRetentionPeriod(pvolLdevID int, muNumber int, request model.SetSnapshotRetentionPeriodRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/snapshots/%d,%d/actions/set-retention/invoke", pvolLdevID, muNumber)

	resIds, err := httpmethod.PostCall(psm.storageSetting, apiSuf, request)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to set retention period for thin image advanced pair %d,%d: %w", pvolLdevID, muNumber, err))
		return "", err
	}

	log.WriteDebug("TFDebug | setSnapshotRetentionPeriodJobID = %v", *resIds)
	return *resIds, nil
}
