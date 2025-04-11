package vssbstorage

import (
	"fmt"
	"strconv"
	"strings"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/vosb/gateway/http"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/gateway/model"
)

// GetAllVolumes gets all available volume details
func (psm *vssbStorageManager) GetAllVolumes() (*vssbmodel.Volumes, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var volumes vssbmodel.Volumes
	apiSuf := "objects/volumes"
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &volumes)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &volumes, nil
}

// Create Volume
func (psm *vssbStorageManager) CreateVolume(reqBody *vssbmodel.CreateVolumeRequestGwy) (*int, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := "objects/volumes"
	affRes, err := httpmethod.PostCall(psm.storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	arr := strings.Split(*affRes, ",")
	sldevID := arr[0]
	ldevID, _ := strconv.Atoi(sldevID)
	log.WriteDebug("TFDebug | ldevID= %d", ldevID)
	return &ldevID, nil
}

// Create AddVolumeToComputeNode
func (psm *vssbStorageManager) AddVolumeToComputeNode(reqBody *vssbmodel.AddVolumeToComputeNodeReq) (*int, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := "objects/volume-server-connections"
	affRes, err := httpmethod.PostCall(psm.storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	arr := strings.Split(*affRes, ",")
	sldevID := arr[0]
	ldevID, _ := strconv.Atoi(sldevID)
	log.WriteDebug("TFDebug | ldevID= %d", ldevID)
	return &ldevID, nil
}

// RemoveVolumeFromComputeNode
func (psm *vssbStorageManager) RemoveVolumeFromComputeNode(volumeId *string, serverId *string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	if volumeId == nil || serverId == nil {
		return fmt.Errorf("volumeId or serverId is nil")
	}

	apiSuf := fmt.Sprintf("objects/volume-server-connections/%s,%s", *volumeId, *serverId)
	_, err := httpmethod.DeleteCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	return nil
}

// UpdateVolumeNickName
func (psm *vssbStorageManager) UpdateVolumeNickName(volumeId *string, nickName *vssbmodel.UpdateVolumeNickNameReq) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if volumeId == nil {
		return fmt.Errorf("volumeId is nil")
	}
	apiSuf := fmt.Sprintf("objects/volumes/%s", *volumeId)
	_, err := httpmethod.PatchCall(psm.storageSetting, apiSuf, nickName)
	if err != nil {
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}

	return nil

}

// ExtendVolumeSize
func (psm *vssbStorageManager) ExtendVolumeSize(volumeId *string, capacity *vssbmodel.UpdateVolumeSizeReq) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	if volumeId == nil {
		return fmt.Errorf("volumeId is nil")
	}
	apiSuf := fmt.Sprintf("objects/volumes/%s/actions/expand/invoke", *volumeId)
	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, capacity)
	if err != nil {
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}

	return nil
}

// DeleteVolume
func (psm *vssbStorageManager) DeleteVolume(volumeId *string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if volumeId == nil {
		return fmt.Errorf("volumeId is nil")
	}
	apiSuf := fmt.Sprintf("objects/volumes/%s", *volumeId)
	_, err := httpmethod.DeleteCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}

	return nil
}
