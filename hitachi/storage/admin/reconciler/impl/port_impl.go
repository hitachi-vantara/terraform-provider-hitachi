package admin

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
)

// ReconcileReadAdminPort reads port information
func (reconcilerManager *adminStorageManager) ReconcileReadAdminPort(portID string) (*gwymodel.PortInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("ReconcileReadAdminPort called for port ID: %s", portID)

	provObj, err := reconcilerManager.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	portInfo, err := provObj.GetPortByID(portID)
	if err != nil {
		log.WriteError("failed to get port %s: %v", portID, err)
		return nil, err
	}

	log.WriteInfo("Successfully read port %s information", portID)
	return portInfo, nil
}

// ReconcileUpdateAdminPort updates port configuration
func (reconcilerManager *adminStorageManager) ReconcileUpdateAdminPort(portID string, params gwymodel.UpdatePortParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("ReconcileUpdateAdminPort called for port ID: %s", portID)

	provObj, err := reconcilerManager.getProvisionerManager()
	if err != nil {
		return err
	}

	err = provObj.UpdatePort(portID, params)
	if err != nil {
		log.WriteError("failed to update port %s: %v", portID, err)
		return err
	}

	log.WriteInfo("Successfully updated port %s", portID)
	return nil
}
