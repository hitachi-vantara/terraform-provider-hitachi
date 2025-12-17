package admin

import (
	"fmt"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/admin/gateway/http"
	model "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
)

func (psm *adminStorageManager) GetIscsiTargets(serverId int) (*model.IscsiTargetInfoList, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var iscsiTargetInfoList model.IscsiTargetInfoList

	apiSuf := fmt.Sprintf("objects/servers/%d/target-iscsi-ports", serverId)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &iscsiTargetInfoList)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	return &iscsiTargetInfoList, nil
}

func (psm *adminStorageManager) GetIscsiTargetByPort(serverId int, portId string) (*model.IscsiTargetInfoByPort, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var iscsiTargetInfo model.IscsiTargetInfoByPort

	apiSuf := fmt.Sprintf("objects/servers/%d/target-iscsi-ports/%s", serverId, portId)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &iscsiTargetInfo)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	return &iscsiTargetInfo, nil
}

func (psm *adminStorageManager) ChangeIscsiTargetName(serverId int, portId string, targetIscsiName string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	params := model.RenameIscsiTargetNameParam{
		TargetIscsiName: targetIscsiName,
	}
	apiSuf := fmt.Sprintf("objects/servers/%d/target-iscsi-ports/%s", serverId, portId)
	_, err := httpmethod.PatchCall(psm.storageSetting, apiSuf, params)
	if err != nil {
		log.WriteError(err)
		return err
	}

	return nil
}
