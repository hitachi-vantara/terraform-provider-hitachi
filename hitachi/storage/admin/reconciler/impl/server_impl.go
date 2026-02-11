package admin

import (
	"encoding/json"
	"fmt"
	"strings"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
)

func (psm *adminStorageManager) ReconcileCreateAdminServer(params gwymodel.CreateAdminServerParams, addHgParams gwymodel.AddHostGroupsToServerParam) (int, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Starting admin server creation reconciliation for '%s'", params.ServerNickname)

	// Debug log the parameters
	b, _ := json.MarshalIndent(params, "", "  ")
	log.WriteDebug("Create Params: %v", string(b))

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		log.WriteError("Failed to get provisioner manager: %v", err)
		return 0, err
	}

	err = provObj.CreateAdminServer(params)
	if err != nil {
		log.WriteError("Failed to create admin server '%s': %v", params.ServerNickname, err)
		return 0, fmt.Errorf("admin server creation failed: %w", err)
	}

	// Find the created server by nickname since we don't get ID back
	listParams := gwymodel.AdminServerListParams{
		Nickname: &params.ServerNickname,
	}

	serverList, err := provObj.GetAdminServerList(listParams)
	if err != nil {
		log.WriteError("Failed to get created server info for '%s': %v", params.ServerNickname, err)
		return 0, fmt.Errorf("failed to retrieve created server info: %w", err)
	}

	if serverList == nil || len(serverList.Data) == 0 {
		return 0, fmt.Errorf("created server '%s' not found", params.ServerNickname)
	}

	serverID := serverList.Data[0].ID

	// add/sync hostgroups
	err = psm.AddHostGroupsToServer(serverID, addHgParams)
	if err != nil {
		// clean-up
		derr := provObj.DeleteAdminServer(serverID, gwymodel.DeleteAdminServerParams{})
		if derr != nil {
			// Graceful handling if deletion is asynchronous or object already gone
			if isNotFoundError(derr) {
				log.WriteWarn("Admin server %d not found or already deleted: %v", serverID, derr)
				return 0, err // delete successful if already gone
			}
			log.WriteError("Clean-up after failure. Failed to delete admin server ID %d: %v", serverID, derr)
			return 0, fmt.Errorf("clean-up after failure, admin server deletion failed: %w", derr)
		}
		return 0, err // add hg failed
	}

	log.WriteInfo("Admin server '%s' created successfully with ID %d", params.ServerNickname, serverID)
	return serverID, nil
}

func (psm *adminStorageManager) ReconcileReadAdminServer(serverID int) (*gwymodel.AdminServerInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Starting admin server read reconciliation for ID %d", serverID)

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		log.WriteError("Failed to get provisioner manager: %v", err)
		return nil, err
	}

	result, err := provObj.GetAdminServerInfo(serverID)
	if err != nil {
		if IsNotFoundError(err) {
			log.WriteWarn("Admin server %d not found or already deleted", serverID)
			return nil, err
		}
		log.WriteError("Failed to get admin server info for ID %d: %v", serverID, err)
		return nil, fmt.Errorf("failed to get admin server info: %w", err)
	}

	if result == nil {
		return nil, fmt.Errorf("no server info returned for ID %d", serverID)
	}

	log.WriteInfo("Admin server ID %d read successfully", serverID)
	return result, nil
}

func (psm *adminStorageManager) ReconcileUpdateAdminServer(serverID int, params gwymodel.UpdateAdminServerParams, addHgParams gwymodel.AddHostGroupsToServerParam) (int, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Starting admin server update reconciliation for ID %d", serverID)

	// Debug log the parameters
	b, _ := json.MarshalIndent(params, "", "  ")
	log.WriteDebug("Update Params: %v", string(b))

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		log.WriteError("Failed to get provisioner manager: %v", err)
		return 0, err
	}

	if params.Nickname != "" || params.OsType != "" || len(params.OsTypeOptions) != 0 {
		err = provObj.UpdateAdminServer(serverID, params)
		if err != nil {
			if IsNotFoundError(err) {
				log.WriteWarn("Admin server %d not found during update", serverID)
				return 0, err
			}
			log.WriteError("Failed to update admin server ID %d: %v", serverID, err)
			return 0, fmt.Errorf("admin server update failed: %w", err)
		}
	}

	// add/sync hostgroups
	err = psm.AddHostGroupsToServer(serverID, addHgParams)
	if err != nil {
		return 0, err
	}

	log.WriteInfo("Admin server ID %d updated successfully", serverID)
	return serverID, nil
}

func (psm *adminStorageManager) ReconcileDeleteAdminServer(serverID int, params gwymodel.DeleteAdminServerParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Starting admin server deletion reconciliation for ID %d", serverID)

	// Debug log the parameters
	b, _ := json.MarshalIndent(params, "", "  ")
	log.WriteDebug("Delete Params: %v", string(b))

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		log.WriteError("Failed to get provisioner manager: %v", err)
		return err
	}

	err = provObj.DeleteAdminServer(serverID, params)
	if err != nil {
		// Graceful handling if deletion is asynchronous or object already gone
		if isNotFoundError(err) {
			log.WriteWarn("Admin server %d not found or already deleted: %v", serverID, err)
			return nil // Consider successful if already gone
		}
		log.WriteError("Failed to delete admin server ID %d: %v", serverID, err)
		return fmt.Errorf("admin server deletion failed: %w", err)
	}

	log.WriteInfo("Admin server ID %d deleted successfully", serverID)
	return nil
}

func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "not found") ||
		strings.Contains(errStr, "does not exist") ||
		strings.Contains(errStr, "404")
}

func (psm *adminStorageManager) GetAdminServerList(params gwymodel.AdminServerListParams) (*gwymodel.AdminServerListResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Starting admin server list retrieval")

	// Debug log the parameters
	b, _ := json.MarshalIndent(params, "", "  ")
	log.WriteDebug("List Params: %v", string(b))

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		log.WriteError("Failed to get provisioner manager: %v", err)
		return nil, err
	}

	result, err := provObj.GetAdminServerList(params)
	if err != nil {
		log.WriteError("Failed to get admin server list: %v", err)
		return nil, fmt.Errorf("admin server list retrieval failed: %w", err)
	}

	if result == nil {
		log.WriteWarn("No admin server list returned")
		return nil, fmt.Errorf("no admin server list returned")
	}

	log.WriteInfo("Admin server list retrieved successfully with %d servers", len(result.Data))
	return result, nil
}

func (psm *adminStorageManager) GetAdminServerInfo(serverID int) (*gwymodel.AdminServerInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Starting admin server info retrieval for ID %d", serverID)

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		log.WriteError("Failed to get provisioner manager: %v", err)
		return nil, err
	}

	result, err := provObj.GetAdminServerInfo(serverID)
	if err != nil {
		if isNotFoundError(err) {
			log.WriteWarn("Admin server %d not found: %v", serverID, err)
			return nil, err // Propagate not found error
		}
		log.WriteError("Failed to get admin server info for ID %d: %v", serverID, err)
		return nil, fmt.Errorf("admin server info retrieval failed: %w", err)
	}

	if result == nil {
		log.WriteWarn("No admin server info returned for ID %d", serverID)
		return nil, fmt.Errorf("no admin server info returned for ID %d", serverID)
	}

	log.WriteInfo("Admin server info retrieved successfully for ID %d (nickname: %s)", serverID, result.Nickname)
	return result, nil
}

func (psm *adminStorageManager) ReconcileSetAdminServerPath(serverID int, params gwymodel.SetAdminServerPathParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Starting admin server path reconciliation for server ID %d", serverID)

	// Debug log the parameters
	b, _ := json.MarshalIndent(params, "", "  ")
	log.WriteDebug("Set Path Params: %v", string(b))

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		log.WriteError("Failed to get provisioner manager: %v", err)
		return err
	}

	err = provObj.SetAdminServerPath(serverID, params)
	if err != nil {
		log.WriteError("Failed to set admin server path for ID %d: %v", serverID, err)
		return fmt.Errorf("admin server path setting failed: %w", err)
	}

	log.WriteInfo("Admin server path set successfully for server ID %d", serverID)
	return nil
}

func (psm *adminStorageManager) ReconcileDeleteAdminServerPath(serverID int, params gwymodel.DeleteAdminServerPathParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Starting admin server path deletion reconciliation for server ID %d", serverID)

	// Debug log the parameters
	b, _ := json.MarshalIndent(params, "", "  ")
	log.WriteDebug("Delete Path Params: %v", string(b))

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		log.WriteError("Failed to get provisioner manager: %v", err)
		return err
	}

	err = provObj.DeleteAdminServerPath(serverID, params)
	if err != nil {
		log.WriteError("Failed to delete admin server path for ID %d: %v", serverID, err)
		return fmt.Errorf("admin server path deletion failed: %w", err)
	}

	log.WriteInfo("Admin server path deleted successfully for server ID %d", serverID)
	return nil
}

func (psm *adminStorageManager) GetAdminServerPath(params gwymodel.AdminServerPathParams) (*gwymodel.AdminServerPathInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Starting admin server path retrieval for server ID %d", params.ServerID)

	// Debug log the parameters
	b, _ := json.MarshalIndent(params, "", "  ")
	log.WriteDebug("Get Path Params: %v", string(b))

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		log.WriteError("Failed to get provisioner manager: %v", err)
		return nil, err
	}

	result, err := provObj.GetAdminServerPath(params)
	if err != nil {
		log.WriteError("Failed to get admin server path for ID %d: %v", params.ServerID, err)
		return nil, fmt.Errorf("admin server path retrieval failed: %w", err)
	}

	log.WriteInfo("Admin server path retrieved successfully for server ID %d", params.ServerID)
	return result, nil
}
