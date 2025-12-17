package admin

import (
	"fmt"
	"net/url"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/admin/gateway/http"
	model "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
)

func (psm *adminStorageManager) GetVolumes(params model.GetVolumeParams) (*model.VolumeInfoList, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var volumeInfoList model.VolumeInfoList

	// Build query string dynamically based on non-nil fields
	q := url.Values{}
	if params.PoolID != nil {
		q.Add("poolId", fmt.Sprintf("%d", *params.PoolID))
	}
	if params.PoolName != nil {
		q.Add("poolName", *params.PoolName)
	}
	if params.ServerID != nil {
		q.Add("serverId", fmt.Sprintf("%d", *params.ServerID))
	}
	if params.ServerNickname != nil {
		q.Add("serverNickname", *params.ServerNickname)
	}
	if params.Nickname != nil {
		q.Add("nickname", *params.Nickname)
	}
	if params.MinTotalCapacity != nil {
		q.Add("minTotalCapacity", fmt.Sprintf("%d", *params.MinTotalCapacity))
	}
	if params.MaxTotalCapacity != nil {
		q.Add("maxTotalCapacity", fmt.Sprintf("%d", *params.MaxTotalCapacity))
	}
	if params.MinUsedCapacity != nil {
		q.Add("minUsedCapacity", fmt.Sprintf("%d", *params.MinUsedCapacity))
	}
	if params.MaxUsedCapacity != nil {
		q.Add("maxUsedCapacity", fmt.Sprintf("%d", *params.MaxUsedCapacity))
	}
	if params.StartVolumeID != nil {
		q.Add("startVolumeId", fmt.Sprintf("%d", *params.StartVolumeID))
	}
	if params.Count != nil {
		q.Add("count", fmt.Sprintf("%d", *params.Count))
	}

	log.WriteDebug("TFDebug| QueryParams:%+v", q)

	apiSuf := fmt.Sprintf("objects/volumes?%s", q.Encode())
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &volumeInfoList)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	return &volumeInfoList, nil
}

func (psm *adminStorageManager) GetVolumeByID(id int) (*model.VolumeInfoByID, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var volumeInfo model.VolumeInfoByID

	apiSuf := fmt.Sprintf("objects/volumes/%d", id)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &volumeInfo)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to get volume with id %d: %w", id, err))
		return nil, err
	}

	return &volumeInfo, nil
}

func (psm *adminStorageManager) CreateVolume(params model.CreateVolumeParams) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := "objects/volumes"
	resIds, err := httpmethod.PostCall(psm.storageSetting, apiSuf, params)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to create volume: %w", err))
		return "", err
	}

	log.WriteDebug("TFDebug | volumeIDs = %v", *resIds)
	return *resIds, nil
}

func (psm *adminStorageManager) DeleteVolume(volumeID int) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/volumes/%d", volumeID)

	_, err := httpmethod.DeleteCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to delete volume %d: %w", volumeID, err))
		return err
	}

	log.WriteDebug("TFDebug | deleted volume ID=%d", volumeID)
	return nil
}

func (psm *adminStorageManager) ExpandVolume(volumeID int, params model.ExpandVolumeParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/volumes/%d/actions/expand/invoke", volumeID)

	// Perform POST call to expand the volume
	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, params)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to expand volume %d: %w", volumeID, err))
		return err
	}

	log.WriteInfo("Successfully triggered expansion for volume ID %d (increment: %d MiB)", volumeID, params.Capacity)
	return nil
}

func (psm *adminStorageManager) UpdateVolumeNickname(volumeID int, params model.UpdateVolumeNicknameParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/volumes/%d", volumeID)

	_, err := httpmethod.PatchCallSync(psm.storageSetting, apiSuf, params)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to update nickname for volume %d: %w", volumeID, err))
		return err
	}

	log.WriteInfo("Successfully updated nickname for volume ID %d to %q", volumeID, params.Nickname)
	return nil
}

func (psm *adminStorageManager) UpdateVolumeReductionSettings(volumeID int, params model.UpdateVolumeReductionParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/volumes/%d", volumeID)

	_, err := httpmethod.PatchCall(psm.storageSetting, apiSuf, params)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to update capacity reduction settings for volume %d: %w", volumeID, err))
		return err
	}

	log.WriteInfo("Successfully triggered capacity reduction setting update for volume ID %d", volumeID)
	return nil
}
