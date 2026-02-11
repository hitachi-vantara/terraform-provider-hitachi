package sanstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/storage/san/gateway/impl"
	sangatewaymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/san/provisioner/message-catalog"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"

	"github.com/jinzhu/copier"
)

// GetSupportedHostModes returns supported host modes and host mode options.
func (psm *sanStorageManager) GetSupportedHostModes() (*sanmodel.HostModeAndOptions, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := sangatewaymodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	log.WriteDebug("TFDebug| Storage Serial:%d, ManagementIP:%s\n", psm.storageSetting.Serial, psm.storageSetting.MgmtIP)
	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_HOSTGROUP_BEGIN), objStorage.Serial)
	hostModesAndOptions, err := gatewayObj.GetSupportedHostModes()
	if err != nil {
		log.WriteDebug("TFError| error in GetSupportedHostModes gateway call, err: %v", err)
		return nil, err
	}

	provModel := sanmodel.HostModeAndOptions{}
	if err := copier.Copy(&provModel, hostModesAndOptions); err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_HOSTGROUP_END), objStorage.Serial)
	return &provModel, nil
}
