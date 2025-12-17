package admin

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/admin/provisioner/message-catalog"
)

// thin pass-through function to call gateway
// no separate model (uses gateway), but has message catalog logging
func (psm *adminStorageManager) GetPorts(queryParams gwymodel.GetPortParams) (*gwymodel.PortInfoList, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PORTS_BEGIN), psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return nil, err
	}

	log.WriteDebug("TFDebug| QueryParams:%+v\n", queryParams)

	portInfoList, err := gatewayObj.GetPorts(queryParams)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_PORTS_FAILED), psm.storageSetting.Serial)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PORTS_END), psm.storageSetting.Serial)
	return portInfoList, nil
}

// thin pass-through function to call gateway
func (psm *adminStorageManager) GetPortByID(portID string) (*gwymodel.PortInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PORT_BY_ID_BEGIN), portID, psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return nil, err
	}

	portInfo, err := gatewayObj.GetPortByID(portID)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_PORT_BY_ID_FAILED), portID, psm.storageSetting.Serial)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PORT_BY_ID_END), portID, psm.storageSetting.Serial)
	return portInfo, nil
}

// thin pass-through function to call gateway
func (psm *adminStorageManager) UpdatePort(portID string, params gwymodel.UpdatePortParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_PORT_BEGIN), portID, psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return err
	}

	log.WriteDebug("TFDebug| UpdateParams:%+v\n", params)

	err = gatewayObj.UpdatePort(portID, params)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_UPDATE_PORT_FAILED), portID, psm.storageSetting.Serial)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_PORT_END), portID, psm.storageSetting.Serial)
	return nil
}
