package admin

import (
	"encoding/json"
	"fmt"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
)

func (psm *adminStorageManager) AddHostGroupsToServer(serverId int, params gwymodel.AddHostGroupsToServerParam) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting adding hostgroups to server")

	b, _ := json.MarshalIndent(params, "", "  ")
	log.WriteDebug("Params: %v", string(b))

	// If no host groups, nothing to reconcile
	if len(params.HostGroups) == 0 {
		log.WriteDebug("No new host_groups provided, skipping add hostgroups to server.")
		return nil
	}

	log.WriteInfo(fmt.Sprintf("Adding new host groups to server: ID %d", serverId))

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return err
	}

	// Step 1: Add host groups
	err = provObj.AddHostGroupsToServer(serverId,  params)
	if err != nil {
		log.WriteError("failed to add host groups to server %d: %w", serverId, err)
		return err
	}

	log.WriteInfo("Successfully added host groups to server: ID %d", serverId)

	// Step 2: Sync host groups with server nickname
	err = provObj.SyncHostGroupsWithServer(serverId)
	if err != nil {
		log.WriteError("failed to synchronize host groups with server: %d: %w", serverId, err)
		return err
	}

	log.WriteInfo("Successfully synchronized host groups with server: ID %d", serverId)

	return nil
}
