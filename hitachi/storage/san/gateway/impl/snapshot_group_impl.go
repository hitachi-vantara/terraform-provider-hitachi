package sanstorage

import (
	"fmt"
	"net/url"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/san/gateway/http"
	model "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
)

// GetSnapshotGroups retrieves information about snapshot groups and their associated Thin Image pairs.
func (psm *sanStorageManager) GetSnapshotGroups(params model.GetSnapshotGroupsParams) (*model.SnapshotGroupListResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var response model.SnapshotGroupListResponse
	q := url.Values{}

	if params.SnapshotGroupName != nil {
		q.Add("snapshotGroupName", *params.SnapshotGroupName)
	}
	if params.DetailInfoType != nil { // pair (vsp5000) or retention (TIA)
		q.Add("detailInfoType", *params.DetailInfoType)
	}

	log.WriteDebug("TFDebug| QueryParams:%+v", q)

	apiSuf := "objects/snapshot-groups"
	if len(q) > 0 {
		apiSuf = fmt.Sprintf("%s?%s", apiSuf, q.Encode())
	}

	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &response)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	return &response, nil
}

// GetSnapshotGroup retrieves detailed information about a specific snapshot group and its pairs.
func (psm *sanStorageManager) GetSnapshotGroup(snapshotGroupID string, params model.GetSnapshotGroupsParams) (*model.SnapshotGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var response model.SnapshotGroup
	q := url.Values{}

	if params.DetailInfoType != nil { // retention
		q.Add("detailInfoType", *params.DetailInfoType)
	}

	log.WriteDebug("TFDebug| QueryParams:%+v", q)

	apiSuf := fmt.Sprintf("objects/snapshot-groups/%s", snapshotGroupID)
	if len(q) > 0 {
		apiSuf = fmt.Sprintf("%s?%s", apiSuf, q.Encode())
	}

	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &response)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	return &response, nil
}

func (psm *sanStorageManager) SplitSnapshotGroup(snapshotGroupID string, request model.SplitSnapshotRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/snapshot-groups/%s/actions/split/invoke", snapshotGroupID)

	resIds, err := httpmethod.PostCall(psm.storageSetting, apiSuf, request)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to split snapshot group %s: %w", snapshotGroupID, err))
		return "", err
	}

	if resIds == nil {
		return "", fmt.Errorf("failed to split snapshot group: no job ID returned")
	}

	log.WriteDebug("TFDebug | splitSnapshotGroupJobID = %v", *resIds)
	return *resIds, nil
}

func (psm *sanStorageManager) ResyncSnapshotGroup(snapshotGroupID string, request model.ResyncSnapshotRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/snapshot-groups/%s/actions/resync/invoke", snapshotGroupID)

	resIds, err := httpmethod.PostCall(psm.storageSetting, apiSuf, request)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to resync snapshot group %s: %w", snapshotGroupID, err))
		return "", err
	}

	if resIds == nil {
		return "", fmt.Errorf("failed to resync snapshot group: no job ID returned")
	}

	log.WriteDebug("TFDebug | resyncSnapshotGroupJobID = %v", *resIds)
	return *resIds, nil
}

func (psm *sanStorageManager) RestoreSnapshotGroup(snapshotGroupID string, request model.RestoreSnapshotRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/snapshot-groups/%s/actions/restore/invoke", snapshotGroupID)

	resIds, err := httpmethod.PostCall(psm.storageSetting, apiSuf, request)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to restore snapshot group %s: %w", snapshotGroupID, err))
		return "", err
	}

	if resIds == nil {
		return "", fmt.Errorf("failed to restore snapshot group: no job ID returned")
	}

	log.WriteDebug("TFDebug | restoreSnapshotGroupJobID = %v", *resIds)
	return *resIds, nil
}

func (psm *sanStorageManager) DeleteSnapshotGroup(snapshotGroupID string) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/snapshot-groups/%s", snapshotGroupID)

	resIds, err := httpmethod.DeleteCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to delete snapshot group %s: %w", snapshotGroupID, err))
		return "", err
	}

	if resIds == nil {
		return "", fmt.Errorf("failed to delete snapshot group: no job ID returned")
	}

	log.WriteDebug("TFDebug | deleteSnapshotGroupJobID = %v", *resIds)
	return *resIds, nil
}

func (psm *sanStorageManager) CloneSnapshotGroup(snapshotGroupID string, request model.CloneSnapshotRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/snapshot-groups/%s/actions/clone/invoke", snapshotGroupID)

	resIds, err := httpmethod.PostCall(psm.storageSetting, apiSuf, request)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to clone snapshot group %s: %w", snapshotGroupID, err))
		return "", err
	}

	if resIds == nil {
		return "", fmt.Errorf("failed to clone snapshot group: no job ID returned")
	}

	log.WriteDebug("TFDebug | cloneSnapshotGroupJobID = %v", *resIds)
	return *resIds, nil
}

// SetSnapshotGroupRetentionPeriod sets the snapshot data retention period for the snapshot data of the Thin Image Advanced pair.
// You can also update (extend) a retention period that is already set.
func (psm *sanStorageManager) SetSnapshotGroupRetentionPeriod(snapshotGroupID string, request model.SetSnapshotRetentionPeriodRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/snapshot-groups/%s/actions/set-retention/invoke", snapshotGroupID)

	resIds, err := httpmethod.PostCall(psm.storageSetting, apiSuf, request)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to set retention period for thin image advanced pair group%s: %w", snapshotGroupID, err))
		return "", err
	}

	log.WriteDebug("TFDebug | setSnapshotGroupRetetionPeriodJobID = %v", *resIds)
	return *resIds, nil
}

func (psm *sanStorageManager) DeleteSnapshotTree(request model.DeleteSnapshotTreeRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := "services/snapshot-tree/actions/delete/invoke"

	resIds, err := httpmethod.PostCall(psm.storageSetting, apiSuf, request)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to delete snapshot tree for root LDEV %d: %w", request.Parameters.LdevID, err))
		return "", err
	}

	if resIds == nil {
		return "", fmt.Errorf("failed to delete snapshot tree: no job ID returned")
	}

	log.WriteDebug("TFDebug | deleteSnapshotTreeJobID = %v", *resIds)
	return *resIds, nil
}

func (psm *sanStorageManager) DeleteGarbageData(request model.DeleteGarbageDataRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := "services/snapshot-tree/actions/delete-garbage-data/invoke"

	resIds, err := httpmethod.PostCall(psm.storageSetting, apiSuf, request)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to %s garbage data deletion for root LDEV %d: %w", request.Parameters.OperationType, request.Parameters.LdevID, err))
		return "", err
	}

	if resIds == nil {
		return "", fmt.Errorf("failed to delete garbage data: no job ID returned")
	}

	log.WriteDebug("TFDebug | deleteGarbageDataJobID = %v", *resIds)
	return *resIds, nil
}
