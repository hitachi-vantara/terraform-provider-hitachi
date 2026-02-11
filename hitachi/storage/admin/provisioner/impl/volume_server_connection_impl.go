package admin

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/admin/provisioner/message-catalog"
)

// thin pass-through function to call gateway
func (psm *adminStorageManager) GetVolumeServerConnections(params gwymodel.GetVolumeServerConnectionsParams) (*gwymodel.VolumeServerConnectionsResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_VOLUME_SERVER_CONNECTIONS_BEGIN), psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return nil, err
	}

	log.WriteDebug("TFDebug| QueryParams:%+v\n", params)

	connList, err := gatewayObj.GetVolumeServerConnections(params)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_VOLUME_SERVER_CONNECTIONS_FAILED), psm.storageSetting.Serial)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_VOLUME_SERVER_CONNECTIONS_END), psm.storageSetting.Serial)
	return connList, nil
}

func (psm *adminStorageManager) GetOneVolumeServerConnection(volumeId, serverId int) (*gwymodel.VolumeServerConnectionDetail, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_VOLUME_SERVER_CONNECTION_BEGIN), volumeId, serverId, psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return nil, err
	}

	connDetail, err := gatewayObj.GetOneVolumeServerConnection(volumeId, serverId)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_VOLUME_SERVER_CONNECTION_FAILED), volumeId, serverId, psm.storageSetting.Serial)
		return nil, err
	}

	log.WriteDebug("TFDebug| VolumeServerConnectionDetail: %+v\n", connDetail)
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_VOLUME_SERVER_CONNECTION_END), volumeId, serverId, psm.storageSetting.Serial)

	return connDetail, nil
}

func (psm *adminStorageManager) AttachVolumeToServers(params gwymodel.AttachVolumeServerConnectionParam) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_ATTACH_VOLUME_SERVER_CONNECTION_BEGIN), psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return "", err
	}

	log.WriteDebug("TFDebug| Params:%+v\n", params)

	connIDs, err := gatewayObj.AttachVolumeToServers(params)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_ATTACH_VOLUME_SERVER_CONNECTION_FAILED), psm.storageSetting.Serial)
		return "", err
	}

	log.WriteDebug("TFDebug| Created connection IDs: %v\n", connIDs)
	log.WriteInfo(mc.GetMessage(mc.INFO_ATTACH_VOLUME_SERVER_CONNECTION_END), psm.storageSetting.Serial)

	return connIDs, nil
}

func (psm *adminStorageManager) DetachVolumeFromServer(volumeId, serverId int) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_DETACH_VOLUME_SERVER_CONNECTION_BEGIN), volumeId, serverId, psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return err
	}

	err = gatewayObj.DetachVolumeToServers(volumeId, serverId)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_DETACH_VOLUME_SERVER_CONNECTION_FAILED), volumeId, serverId, psm.storageSetting.Serial)
		return err
	}

	log.WriteDebug("TFDebug| Detached volume-server connection %d,%d\n", volumeId, serverId)
	log.WriteInfo(mc.GetMessage(mc.INFO_DETACH_VOLUME_SERVER_CONNECTION_END), volumeId, serverId, psm.storageSetting.Serial)

	return nil
}
