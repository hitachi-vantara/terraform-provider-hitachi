package vssbstorage

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/storage/vssb/gateway/impl"
	vssbgatewaymodel "terraform-provider-hitachi/hitachi/storage/vssb/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/vssb/provisioner/message-catelog"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vssb/provisioner/model"

	"github.com/jinzhu/copier"
)

// GetAllVolumes gets all available volume details
func (psm *vssbStorageManager) GetAllVolumes(computeNodeName string) (*vssbmodel.Volumes, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := vssbgatewaymodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_VOLUME_INFO_BEGIN))
	var volumeInfo *vssbgatewaymodel.Volumes = &vssbgatewaymodel.Volumes{}
	volumeData, err := gatewayObj.GetAllVolumes()
	if err != nil {
		log.WriteDebug("TFError| failed to call GetAllVolumes err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_VOLUME_INFO_FAILED))
		return nil, err
	}

	// Assign Volume based on condition
	if computeNodeName == "" {
		volumeInfo = volumeData
	} else {
		serverID, err := psm.GetServerIdByComputeNode(computeNodeName)
		if err != nil {
			log.WriteDebug("TFError| failed to call GetServerIdByComputeNode err: %+v", err)
			return nil, err
		}
		if serverID == "" {
			log.WriteDebug("TFError| ServerID not found in - GetServerIdByComputeNode")
			return nil, fmt.Errorf("unable to find compute node name, please add valid compute node name")
		}
		finalVolums, err := psm.GetVolumesByServerId(serverID, volumeData)
		if err != nil {
			log.WriteDebug("TFError| failed to call GetVolumesByServerId err: %+v", err)
			return nil, err
		}
		volumeInfo = finalVolums
	}

	provVolumes := vssbmodel.Volumes{}
	err = copier.Copy(&provVolumes, volumeInfo)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_VOLUME_INFO_END))
	return &provVolumes, nil
}

// GetServerIdByComputeNode .
func (psm *vssbStorageManager) GetServerIdByComputeNode(computeNodeName string) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := vssbgatewaymodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return "", err
	}
	computeNodes, err := gatewayObj.GetAllComputeNodes()
	if err != nil {
		log.WriteDebug("TFError| failed to call GetAllVolumes err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_SERVERS_FAILED))
		return "", err
	}

	for _, node := range computeNodes.Data {
		if node.Nickname == computeNodeName {
			return node.ID, nil
		}
	}

	return "", nil
}

// GetVolumesByServerId .
func (psm *vssbStorageManager) GetVolumesByServerId(serverId string, allVolumes *vssbgatewaymodel.Volumes) (*vssbgatewaymodel.Volumes, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := vssbgatewaymodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}
	volumeServerInfo, err := gatewayObj.GetConnectionInfoBtwnVolumeAndServerByServerID(serverId)
	if err != nil {
		log.WriteDebug("TFError| failed to call GetConnectionInfoBtwnVolumeAndServerByServerID err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_CONNECTION_BY_SERVER_FAILED), serverId)
		return nil, err
	}
	finalVolumes := []vssbgatewaymodel.Volume{}
	if volumeServerInfo.Data != nil {
		for _, volume := range allVolumes.Data {
			for _, relationVolume := range volumeServerInfo.Data {
				if volume.ID == relationVolume.VolumeId {
					finalVolumes = append(finalVolumes, volume)
					break
				}
			}
		}
	}
	resultVolume := vssbgatewaymodel.Volumes{
		Data: finalVolumes,
	}
	return &resultVolume, nil
}

// Get volume infomration based on the volume name
func (psm *vssbStorageManager) GetVolumeDetails(volumeName string) (*vssbmodel.Volume, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := vssbgatewaymodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_VOLUME_INFO_BEGIN))
	volumesData, err := gatewayObj.GetAllVolumes()
	if err != nil {
		log.WriteDebug("TFError| failed to call GetAllVolumes err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_VOLUME_INFO_FAILED))
		return nil, err
	}

	var volumeId string = ""
	provVolume := vssbmodel.Volume{}

	// Assign Volume based on condition
	for _, volume := range volumesData.Data {
		if volume.Name == volumeName {
			volumeId = volume.ID
			err = copier.Copy(&provVolume, volume)
			if err != nil {
				log.WriteDebug("Copy Error: %v", err)
				return nil, err
			}
			break
		}
	}
	if volumeId == "" {
		return nil, fmt.Errorf("no volume found on the given volume name")
	}
	computeServerIDs, err := psm.GetComputeNodesIDsByVolumeId(&volumeId)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}
	computeNodes, err := psm.GetComputeNodesbyComputeNodeId(computeServerIDs)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_VOLUME_INFO_END))

	provcomputeNodes := vssbmodel.Servers{}
	if computeNodes == nil {
		return &provVolume, nil
	}

	err = copier.Copy(&provcomputeNodes, computeNodes)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}
	provVolume.ComputeNodes = append(provVolume.ComputeNodes, provcomputeNodes.Data...)

	return &provVolume, nil

}

// Get volume infomration based on the volume name
func (psm *vssbStorageManager) GetVolumeDetailsByName(volumeName string) (*vssbmodel.Volume, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := vssbgatewaymodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_VOLUME_INFO_BEGIN))
	volumesData, err := gatewayObj.GetAllVolumes()
	if err != nil {
		log.WriteDebug("TFError| failed to call GetAllVolumes err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_VOLUME_INFO_FAILED))
		return nil, err
	}

	var volumeId string = ""
	provVolume := vssbmodel.Volume{}

	// Assign Volume based on condition
	for _, volume := range volumesData.Data {
		if volume.Name == volumeName {
			volumeId = volume.ID
			err = copier.Copy(&provVolume, volume)
			if err != nil {
				log.WriteDebug("Copy Error: %v", err)
				return nil, err
			}
			break
		}
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_VOLUME_INFO_END))

	if volumeId == "" {
		return nil, fmt.Errorf("volume not found")
	} else {
		return &provVolume, nil
	}

}

// get the array of compute nodes by filtering the volume id
func (psm *vssbStorageManager) GetComputeNodesIDsByVolumeId(volumeID *string) (*[]map[string]interface{}, error) {

	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := vssbgatewaymodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	computeNodeServerInfo, err := gatewayObj.GetConnectionInfoBtwnVolumeAndServerByVolumeID(*volumeID)

	if err != nil {
		log.WriteDebug("TFError| failed to call GetConnectionInfoBtwnVolumeAndServerByServerID err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_CONNECTION_BY_SERVER_FAILED), volumeID)
		return nil, err
	}
	computeServerIDs := []map[string]interface{}{}

	if computeNodeServerInfo.Data != nil {
		for _, ServerInfo := range computeNodeServerInfo.Data {
			mapData := map[string]interface{}{
				"ID":  ServerInfo.ServerId,
				"LUN": ServerInfo.Lun,
			}
			computeServerIDs = append(computeServerIDs, mapData)

		}
	}
	return &computeServerIDs, nil

}

// Get the array of compute nodes depath details by using compute node id
func (psm *vssbStorageManager) GetComputeNodesbyComputeNodeId(computeNodeIds *[]map[string]interface{}) (*vssbgatewaymodel.Servers, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := vssbgatewaymodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	finalComputeNodes := []vssbgatewaymodel.Server{}
	allComputeNodes, err := gatewayObj.GetAllComputeNodes()

	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}
	// Retrive all the compute nodes and filter based on inputinterface
	if allComputeNodes.Data != nil {
		for _, nid := range *computeNodeIds {
			for _, node := range allComputeNodes.Data {
				if node.ID == nid["ID"] {
					node.Lun = nid["LUN"].(int)
					finalComputeNodes = append(finalComputeNodes, node)
				}
			}
		}
	} else {

		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}
	computeNodes := vssbgatewaymodel.Servers{
		Data: finalComputeNodes,
	}
	return &computeNodes, nil

}

func (psm *vssbStorageManager) CreateVolume(name string, nickName string, poolId string, capacity float32) (*int, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := vssbgatewaymodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_VOLUME_BEGIN))

	var capacityMB = (capacity * 1024)
	var capacityGBInt int32 = int32(capacityMB)

	Name := vssbgatewaymodel.NameParam{
		BaseName: &name,
	}
	nickname := vssbgatewaymodel.NickNameParam{
		BaseName: &nickName,
	}
	reqBody := &vssbgatewaymodel.CreateVolumeRequestGwy{
		PoolID:        &poolId,
		NameParam:     Name,
		NickNameParam: nickname,
		Capacity:      &capacityGBInt,
	}
	volId, err := gatewayObj.CreateVolume(reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in CreateVolume call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_VOLUME_FAILED))
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_VOLUME_END))

	return volId, nil
}

func (psm *vssbStorageManager) AddVolumeToComputeNode(volumeId string, serverId string) (*int, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := vssbgatewaymodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}
	gatwayModel := vssbgatewaymodel.AddVolumeToComputeNodeReq{
		VolumeID: volumeId,
		ServerID: serverId,
	}
	nodeAdditions, err := gatewayObj.AddVolumeToComputeNode(&gatwayModel)
	log.WriteInfo(mc.GetMessage(mc.INFO_ADD_VOLUME_TO_COMPUTE_NODES_BEGIN), volumeId)
	if err != nil {
		log.WriteDebug("TFError| error in AddVolumeToComputeNode call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_ADD_VOLUME_TO_COMPUTE_NODES_FAILED), volumeId)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_ADD_VOLUME_TO_COMPUTE_NODES_END), volumeId)
	return nodeAdditions, nil
}

func (psm *vssbStorageManager) UpdateVolumeNickName(serverId string, nickName string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := vssbgatewaymodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}
	NickNameodel := vssbgatewaymodel.UpdateVolumeNickNameReq{
		NickName: nickName,
	}
	err = gatewayObj.UpdateVolumeNickName(&serverId, &NickNameodel)
	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_VOLUME_NICKNAME_BEGIN), serverId)
	if err != nil {
		log.WriteDebug("TFError| error in UpdateVolumeNickName call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_UPDATE_VOLUME_NICKNAME_FAILED), serverId)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_VOLUME_NICKNAME_END), serverId)
	return nil
}

func (psm *vssbStorageManager) ExpandVolume(serverId string, additionalCapacity *int32) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := vssbgatewaymodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}
	additionalVolModel := vssbgatewaymodel.UpdateVolumeSizeReq{
		AdditionalCapacity: additionalCapacity,
	}
	err = gatewayObj.ExtendVolumeSize(&serverId, &additionalVolModel)
	log.WriteInfo(mc.GetMessage(mc.INFO_EXPAND_VOLUME_SIZE_BEGIN), serverId)
	if err != nil {
		log.WriteDebug("TFError| error in ExtendVolumeSize call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_EXPAND_VOLUME_SIZE_FAILED), serverId)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_EXPAND_VOLUME_SIZE_END), serverId)
	return nil
}

func (psm *vssbStorageManager) RemoveVolumeFromComputeNode(volumeId string, serverId string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := vssbgatewaymodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}
	err = gatewayObj.RemoveVolumeFromComputeNode(&volumeId, &serverId)
	log.WriteInfo(mc.GetMessage(mc.INFO_REMOVE_VOLUME_TO_COMPUTE_NODES_BEGIN), volumeId)
	if err != nil {
		log.WriteDebug("TFError| error in RemoveVolumeFromComputeNode call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_REMOVE_VOLUME_TO_COMPUTE_NODES_FAILED), volumeId)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_REMOVE_VOLUME_TO_COMPUTE_NODES_END), volumeId)
	return nil
}

func (psm *vssbStorageManager) DeleteVolume(volumeId string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := vssbgatewaymodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_VOLUME_BEGIN))
	err = gatewayObj.DeleteVolume(&volumeId)
	if err != nil {
		log.WriteDebug("TFError| error in Delete volume call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_VOLUME_FAILED))
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_VOLUME_END))
	return nil
}
