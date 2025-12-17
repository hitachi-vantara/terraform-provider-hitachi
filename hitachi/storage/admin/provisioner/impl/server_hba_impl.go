package admin

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/admin/provisioner/message-catalog"
)

// thin pass-through function to call gateway
// no separate model (uses gateway), but has message catalog logging
func (psm *adminStorageManager) GetServerHBAs(serverID int) (*gwymodel.ServerHBAList, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_SERVER_HBAS_BEGIN), serverID)

	gwyObj, err := psm.getGatewayManager()
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_SERVER_HBAS_FAILED))
		return nil, err
	}

	serverHBAs, err := gwyObj.GetServerHBAs(serverID)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_SERVER_HBAS_FAILED))
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_SERVER_HBAS_END))
	return serverHBAs, nil
}

// thin pass-through function to call gateway
func (psm *adminStorageManager) GetServerHBAByWwn(serverID int, hbaWwn string) (*gwymodel.ServerHBA, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_SERVER_HBA_BEGIN), serverID, hbaWwn, psm.storageSetting.Serial)

	gwyObj, err := psm.getGatewayManager()
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_SERVER_HBA_FAILED))
		return nil, err
	}

	serverHBA, err := gwyObj.GetServerHBAByWwn(serverID, hbaWwn)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_SERVER_HBA_FAILED))
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_SERVER_HBA_END), serverHBA)
	return serverHBA, nil
}

// thin pass-through function to call gateway
func (psm *adminStorageManager) CreateServerHBAs(serverID int, params gwymodel.CreateServerHBAParams) (*gwymodel.ServerHBAList, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Creating server HBAs for server ID: %d", serverID)

	gwyObj, err := psm.getGatewayManager()
	if err != nil {
		return nil, err
	}

	log.WriteDebug("TFDebug| CreateParams:%+v\n", params)

	serverHBAList, err := gwyObj.CreateServerHBAs(serverID, params)
	if err != nil {
		log.WriteError("Failed to create server HBAs for server ID %d: %v", serverID, err)
		return nil, err
	}

	log.WriteInfo("Successfully created server HBAs for server ID: %d", serverID)
	return serverHBAList, nil
}

// thin pass-through function to call gateway
func (psm *adminStorageManager) DeleteServerHBA(serverID int, initiatorName string) (*gwymodel.ServerHBAList, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Deleting server HBA for server ID: %d, initiator: %s", serverID, initiatorName)

	gwyObj, err := psm.getGatewayManager()
	if err != nil {
		return nil, err
	}

	serverHBAList, err := gwyObj.DeleteServerHBA(serverID, initiatorName)
	if err != nil {
		log.WriteError("Failed to delete server HBA for server ID %d, initiator %s: %v", serverID, initiatorName, err)
		return nil, err
	}

	log.WriteInfo("Successfully deleted server HBA for server ID: %d, initiator: %s", serverID, initiatorName)
	return serverHBAList, nil
}
