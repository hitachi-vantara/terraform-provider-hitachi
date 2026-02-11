package admin

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
)

// ReconcileCreateAdminServerHBAs creates server HBAs
func (reconcilerManager *adminStorageManager) ReconcileCreateAdminServerHBAs(serverID int, params gwymodel.CreateServerHBAParams) (*gwymodel.ServerHBAList, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("ReconcileCreateAdminServerHBAs called for server ID: %d", serverID)

	provObj, err := reconcilerManager.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	serverHBAList, err := provObj.CreateServerHBAs(serverID, params)
	if err != nil {
		log.WriteError("failed to create server HBAs for server ID %d: %v", serverID, err)
		return nil, err
	}

	log.WriteInfo("Successfully created server HBAs for server ID %d", serverID)
	return serverHBAList, nil
}

// ReconcileDeleteAdminServerHBA deletes a server HBA
func (reconcilerManager *adminStorageManager) ReconcileDeleteAdminServerHBA(serverID int, initiatorName string) (*gwymodel.ServerHBAList, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("ReconcileDeleteAdminServerHBA called for server ID: %d, initiator: %s", serverID, initiatorName)

	provObj, err := reconcilerManager.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	serverHBAList, err := provObj.DeleteServerHBA(serverID, initiatorName)
	if err != nil {
		log.WriteError("failed to delete server HBA for server ID %d, initiator %s: %v", serverID, initiatorName, err)
		return nil, err
	}

	log.WriteInfo("Successfully deleted server HBA for server ID %d, initiator: %s", serverID, initiatorName)
	return serverHBAList, nil
}
