package terraform

import (
	"errors"
	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	common "terraform-provider-hitachi/hitachi/terraform/common"

	// mc "terraform-provider-hitachi/hitachi/messagecatalog"

	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"

	model "terraform-provider-hitachi/hitachi/infra_gw/model"
	reconimpl "terraform-provider-hitachi/hitachi/infra_gw/reconciler/impl"

	//terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetInfraChapUsers(d *schema.ResourceData) ([]string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := common.GetSerialString(d)
	storageId := d.Get("storage_id").(string)

	if serial == "" && storageId == "" {
		err := errors.New("both serial and storage_id can't be empty. Please specify one")
		return nil, err
	}

	if serial != "" && storageId != "" {
		err := errors.New("both serial and storage_id are not allowed. Either serial or storage_id can be specified")
		return nil, err
	}

	address, err := cache.GetCurrentAddress()
	if err != nil {
		return nil, err
	}

	if storageId == "" {
		storageId, err = common.GetStorageIdFromSerial(address, serial)
		if err != nil {
			return nil, err
		}
		d.Set("storage_id", storageId)
	}

	port := d.Get("port_id").(string)

	iscsi_id := -1
	hid, okId := d.GetOk("iscsi_target_number")
	if okId {
		iscsi_id = hid.(int)
	}

	iscsiTargetId, err := common.GetIscsiTargetId(address, storageId, port, iscsi_id)
	if err != nil {
		return nil, err
	}

	storageSetting, err := cache.GetInfraSettingsFromCache(address)
	if err != nil {
		return nil, err
	}

	setting := model.InfraGwSettings{
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		Address:  storageSetting.Address,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GW_GET_ISCSI_TARGETS_BEGIN), setting.Address)
	reconResponse, err := reconObj.GetIscsiTarget(storageId, iscsiTargetId)
	if err != nil {
		log.WriteDebug("TFError| error getting GetIscsiTarget, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_INFRA_GW_GET_ISCSI_TARGETS_FAILED), setting.Address)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GW_GET_ISCSI_TARGETS_END), setting.Address)

	return reconResponse.Data.ChapUsers, nil
}
