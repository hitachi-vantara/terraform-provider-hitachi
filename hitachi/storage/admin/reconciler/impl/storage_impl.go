package admin

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	provisonerimpl "terraform-provider-hitachi/hitachi/storage/admin/provisioner/impl"
	provisonermodel "terraform-provider-hitachi/hitachi/storage/admin/provisioner/model"
	mc "terraform-provider-hitachi/hitachi/storage/admin/reconciler/message-catalog"
	adminmodel "terraform-provider-hitachi/hitachi/storage/admin/reconciler/model"

	"github.com/jinzhu/copier"
)

// GetStorageAdminInfo get storage admin info
func (psm *adminStorageManager) GetStorageAdminInfo(configurable_capacities bool) (*adminmodel.StorageAdminInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisonermodel.StorageDeviceSettings{
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		Serial:   psm.storageSetting.Serial,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_ADMIN_INFO_BEGIN))
	provStorageAdminInfo, err := provObj.GetStorageAdminInfo(configurable_capacities)
	log.WriteDebug("TFError| Prov StorageAdminInfo: %v", provStorageAdminInfo)
	if err != nil {
		log.WriteDebug("TFError| error in GetStorageAdminInfo provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_STORAGE_ADMIN_INFO_FAILED))
		return nil, err
	}
	// Converting Prov to Reconciler
	reconcilerStorageAdminInfo := adminmodel.StorageAdminInfo{}
	err = copier.Copy(&reconcilerStorageAdminInfo, provStorageAdminInfo)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_ADMIN_INFO_END))
	log.WriteDebug("TFError| Recon StorageAdminInfo: %v", reconcilerStorageAdminInfo)
	return &reconcilerStorageAdminInfo, nil
}
