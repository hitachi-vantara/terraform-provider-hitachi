package admin

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/storage/admin/gateway/impl"
	admingatewaymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/admin/provisioner/message-catalog"
	adminmodel "terraform-provider-hitachi/hitachi/storage/admin/provisioner/model"

	"github.com/jinzhu/copier"
)

// GetStorageAdminInfo Obtains the storage admin information.
func (psm *adminStorageManager) GetStorageAdminInfo(configurable_capacities bool) (*adminmodel.StorageAdminInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := admingatewaymodel.StorageDeviceSettings{
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		Serial:   psm.storageSetting.Serial,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGEADMININFO_BEGIN))
	storageAdminInfo, err := gatewayObj.GetStorageAdminInfo(configurable_capacities)
	if err != nil {
		log.WriteDebug("TFError| failed to call GetStorageAdminInfo, err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_STORAGEADMININFO_FAILED))
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGEADMININFO_END))

	provStorageAdminInfo := adminmodel.StorageAdminInfo{}
	err = copier.Copy(&provStorageAdminInfo, storageAdminInfo)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}

	return &provStorageAdminInfo, nil
}
