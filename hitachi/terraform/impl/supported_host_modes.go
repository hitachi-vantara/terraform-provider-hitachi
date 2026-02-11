package terraform

import (
	"strconv"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	reconimpl "terraform-provider-hitachi/hitachi/storage/san/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

func GetSupportedHostModes(d *schema.ResourceData) (*terraformmodel.SupportedHostModes, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	storageSetting, err := cache.GetSanSettingsFromCache(strconv.Itoa(serial))
	if err != nil {
		return nil, err
	}

	setting := reconcilermodel.StorageDeviceSettings{
		Serial:   storageSetting.Serial,
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		MgmtIP:   storageSetting.MgmtIP,
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_HOSTGROUP_BEGIN), setting.Serial)
	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx, err: %v", err)
		return nil, err
	}

	reconModel, err := reconObj.GetSupportedHostModes()
	if err != nil {
		return nil, err
	}

	terraformModel := terraformmodel.SupportedHostModes{}
	if err := copier.Copy(&terraformModel, reconModel); err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_HOSTGROUP_END), setting.Serial)
	return &terraformModel, nil
}
