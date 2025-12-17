package admin

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/admin/provisioner/message-catalog"
)

// AddHostGroupsToServer - thin pass-through to gateway
func (psm *adminStorageManager) AddHostGroupsToServer(serverId int, params gwymodel.AddHostGroupsToServerParam) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_ADD_HOSTGROUPS_TO_SERVER_BEGIN), serverId, psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return err
	}

	log.WriteDebug("TFDebug| Params:%+v\n", params)

	err = gatewayObj.AddHostGroupsToServer(serverId, params)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_ADD_HOSTGROUPS_TO_SERVER_FAILED), serverId, psm.storageSetting.Serial)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_ADD_HOSTGROUPS_TO_SERVER_END), serverId, psm.storageSetting.Serial)
	return nil
}

// SyncHostGroupsWithServer - thin pass-through to gateway
func (psm *adminStorageManager) SyncHostGroupsWithServer(serverId int) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_SYNC_HOSTGROUPS_WITH_SERVER_BEGIN), serverId, psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return err
	}

	err = gatewayObj.SyncHostGroupsWithServer(serverId)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_SYNC_HOSTGROUPS_WITH_SERVER_FAILED), serverId, psm.storageSetting.Serial)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_SYNC_HOSTGROUPS_WITH_SERVER_END), serverId, psm.storageSetting.Serial)
	return nil
}
