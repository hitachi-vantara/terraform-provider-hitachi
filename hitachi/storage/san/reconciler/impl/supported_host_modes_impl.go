package sanstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	provisonerimpl "terraform-provider-hitachi/hitachi/storage/san/provisioner/impl"
	provisonermodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"
	mc "terraform-provider-hitachi/hitachi/storage/san/reconciler/message-catalog"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"

	"github.com/jinzhu/copier"
)

// GetSupportedHostModes returns supported host modes and host mode options.
func (psm *sanStorageManager) GetSupportedHostModes() (*sanmodel.HostModeAndOptions, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisonermodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	log.WriteDebug("TFDebug| Storage Serial:%d, ManagementIP:%s\n", psm.storageSetting.Serial, psm.storageSetting.MgmtIP)
	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_HOSTGROUP_BEGIN), objStorage.Serial)
	provModel, err := provObj.GetSupportedHostModes()
	if err != nil {
		log.WriteDebug("TFError| error in GetSupportedHostModes provisioner call, err: %v", err)
		return nil, err
	}

	reconModel := sanmodel.HostModeAndOptions{}
	if err := copier.Copy(&reconModel, provModel); err != nil {
		log.WriteDebug("TFError| error in Copy from provisioner to reconciler structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_HOSTGROUP_END), objStorage.Serial)
	return &reconModel, nil
}
