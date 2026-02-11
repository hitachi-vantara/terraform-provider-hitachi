package admin

import (
	"fmt"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/admin/gateway/http"
	model "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
)

func (psm *adminStorageManager) AddHostGroupsToServer(serverId int, params model.AddHostGroupsToServerParam) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// hostgroup must already be existing
	apiSuf := fmt.Sprintf("objects/servers/%d/actions/add-host-groups/invoke", serverId)
	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, params)
	if err != nil {
		log.WriteError(err)
		return err
	}

	return nil
}

// Renames the host group names so they match the server nickname.
func (psm *adminStorageManager) SyncHostGroupsWithServer(serverId int) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/servers/%d/actions/sync-host-group-names/invoke", serverId)
	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError(err)
		return err
	}

	return nil
}
