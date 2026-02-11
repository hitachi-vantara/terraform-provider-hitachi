package admin

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/admin/provisioner/message-catalog"
)

// GetIscsiTargets - thin pass-through to gateway
func (psm *adminStorageManager) GetIscsiTargets(serverId int) (*gwymodel.IscsiTargetInfoList, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ISCSI_TARGETS_BEGIN), serverId, psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return nil, err
	}

	iscsiTargetInfoList, err := gatewayObj.GetIscsiTargets(serverId)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_ISCSI_TARGETS_FAILED), serverId, psm.storageSetting.Serial)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ISCSI_TARGETS_END), serverId, psm.storageSetting.Serial)
	return iscsiTargetInfoList, nil
}

// GetIscsiTargetByPort - thin pass-through to gateway
func (psm *adminStorageManager) GetIscsiTargetByPort(serverId int, portId string) (*gwymodel.IscsiTargetInfoByPort, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ISCSI_TARGET_BY_PORT_BEGIN), serverId, portId, psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return nil, err
	}

	targetInfo, err := gatewayObj.GetIscsiTargetByPort(serverId, portId)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_ISCSI_TARGET_BY_PORT_FAILED), serverId, portId, psm.storageSetting.Serial)
		return nil, err
	}

	log.WriteDebug("TFDebug| TargetInfoByPort: %+v\n", targetInfo)
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ISCSI_TARGET_BY_PORT_END), serverId, portId, psm.storageSetting.Serial)

	return targetInfo, nil
}

// ChangeIscsiTargetName - thin pass-through to gateway
func (psm *adminStorageManager) ChangeIscsiTargetName(serverId int, portId string, targetIscsiName string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_CHANGE_ISCSI_TARGET_NAME_BEGIN), serverId, portId, targetIscsiName, psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return err
	}

	log.WriteDebug("TFDebug| ServerId:%v, PortId:%v, TargetIscsiName:%v\n", serverId, portId, targetIscsiName)

	err = gatewayObj.ChangeIscsiTargetName(serverId, portId, targetIscsiName)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_CHANGE_ISCSI_TARGET_NAME_FAILED), serverId, portId, psm.storageSetting.Serial)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CHANGE_ISCSI_TARGET_NAME_END), serverId, portId, psm.storageSetting.Serial)
	return nil
}
