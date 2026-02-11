package sanstorage

import (
	"fmt"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/san/gateway/http"
	model "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
)

// GetVirtualCloneParentVolumes retrieves a list of LDEV IDs that are valid virtual clone parents.
func (psm *sanStorageManager) GetVirtualCloneParentVolumes() (*model.VirtualCloneParentVolumeList, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var response model.VirtualCloneParentVolumeList

	apiSuf := "objects/virtual-clone-parent-volumes"

	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &response)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	log.WriteDebug("TFDebug| VirtualCloneParentVolumes: %+v", response)

	return &response, nil
}

// GetSnapshotFamily retrieves the snapshot family information for a specific LDEV.
func (psm *sanStorageManager) GetSnapshotFamily(ldevID int) (*model.SnapshotFamilyListResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var response model.SnapshotFamilyListResponse

	apiSuf := fmt.Sprintf("objects/snapshot-family?ldevId=%d", ldevID)

	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &response)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	log.WriteDebug("TFDebug| SnapshotFamily for LDEV %d: %+v", ldevID, response)
	return &response, nil
}

// CreateSnapshotVClone initiates the creation of a virtual clone for a specific snapshot.
func (psm *sanStorageManager) CreateSnapshotVClone(pvolLdevID int, muNumber int) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/snapshots/%d,%d/actions/virtual-clone/invoke", pvolLdevID, muNumber)

	request := model.VirtualCloneRequest{
		Parameters: model.VirtualCloneParams{
			OperationType: "create",
		},
	}

	resIds, err := httpmethod.PostCall(psm.storageSetting, apiSuf, request)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to create virtual clone for snapshot %d,%d: %w", pvolLdevID, muNumber, err))
		return "", err
	}

	log.WriteDebug("TFDebug | createSnapshotVCloneJobID = %v", *resIds)
	return *resIds, nil
}

// ConvertSnapshotVClone converts an existing snapshot to a virtual clone.
func (psm *sanStorageManager) ConvertSnapshotVClone(pvolLdevID int, muNumber int) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/snapshots/%d,%d/actions/virtual-clone/invoke", pvolLdevID, muNumber)

	request := model.VirtualCloneRequest{
		Parameters: model.VirtualCloneParams{
			OperationType: "convert",
		},
	}

	resIds, err := httpmethod.PostCall(psm.storageSetting, apiSuf, request)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to convert snapshot %d,%d to virtual clone: %w", pvolLdevID, muNumber, err))
		return "", err
	}

	log.WriteDebug("TFDebug | convertSnapshotVCloneJobID = %v", *resIds)
	return *resIds, nil
}

// RestoreSnapshotFromVClone restores snapshot data from a virtual clone.
func (psm *sanStorageManager) RestoreSnapshotFromVClone(pvolLdevID int, muNumber int) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/snapshots/%d,%d/actions/virtual-clone/invoke", pvolLdevID, muNumber)

	request := model.VirtualCloneRequest{
		Parameters: model.VirtualCloneParams{
			OperationType: "restore",
		},
	}

	resIds, err := httpmethod.PostCall(psm.storageSetting, apiSuf, request)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to restore snapshot %d,%d from virtual clone: %w", pvolLdevID, muNumber, err))
		return "", err
	}

	log.WriteDebug("TFDebug | restoreSnapshotFromVCloneJobID = %v", *resIds)
	return *resIds, nil
}

// CreateSnapshotGroupVClone initiates virtual clone creation for a snapshot group.
func (psm *sanStorageManager) CreateSnapshotGroupVClone(snapshotGroupName string) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/snapshot-groups/%s/actions/virtual-clone/invoke", snapshotGroupName)

	request := model.VirtualCloneRequest{
		Parameters: model.VirtualCloneParams{
			OperationType: "create",
		},
	}

	resIds, err := httpmethod.PostCall(psm.storageSetting, apiSuf, request)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to create virtual clone for snapshot group %s: %w", snapshotGroupName, err))
		return "", err
	}

	log.WriteDebug("TFDebug | createSnapshotGroupVCloneJobID = %v", *resIds)
	return *resIds, nil
}

// ConvertSnapshotGroupVClone converts a snapshot group to virtual clones.
func (psm *sanStorageManager) ConvertSnapshotGroupVClone(snapshotGroupName string) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/snapshot-groups/%s/actions/virtual-clone/invoke", snapshotGroupName)

	request := model.VirtualCloneRequest{
		Parameters: model.VirtualCloneParams{
			OperationType: "convert",
		},
	}

	resIds, err := httpmethod.PostCall(psm.storageSetting, apiSuf, request)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to convert snapshot group %s to virtual clone: %w", snapshotGroupName, err))
		return "", err
	}

	log.WriteDebug("TFDebug | convertSnapshotGroupVCloneJobID = %v", *resIds)
	return *resIds, nil
}

// RestoreSnapshotGroupFromVClone restores an entire snapshot group from virtual clones.
func (psm *sanStorageManager) RestoreSnapshotGroupFromVClone(snapshotGroupName string) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/snapshot-groups/%s/actions/virtual-clone/invoke", snapshotGroupName)

	request := model.VirtualCloneRequest{
		Parameters: model.VirtualCloneParams{
			OperationType: "restore",
		},
	}

	resIds, err := httpmethod.PostCall(psm.storageSetting, apiSuf, request)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to restore snapshot group %s from virtual clone: %w", snapshotGroupName, err))
		return "", err
	}

	log.WriteDebug("TFDebug | restoreSnapshotGroupFromVCloneJobID = %v", *resIds)
	return *resIds, nil
}
