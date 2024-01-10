package sanstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	provisonerimpl "terraform-provider-hitachi/hitachi/storage/san/provisioner/impl"
	provisonermodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"
	mc "terraform-provider-hitachi/hitachi/storage/san/reconciler/message-catalog"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"

	"github.com/jinzhu/copier"
)

// GetParityGroups is used to fetch all parity group details and also we can filter by ids
func (psm *sanStorageManager) GetParityGroups(parityGroupIds ...[]string) (*[]sanmodel.ParityGroup, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PARITY_GROUP_BEGIN), objStorage.Serial)
	provParityGroups, err := provObj.GetParityGroups(parityGroupIds...)
	if err != nil {
		log.WriteDebug("TFError| error in GetParityGroups provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_PARITY_GROUP_FAILED), objStorage.Serial)
		return nil, err
	}

	// Converting Prov to Reconciler
	reconcileParityGroups := []sanmodel.ParityGroup{}
	err = copier.Copy(&reconcileParityGroups, provParityGroups)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PARITY_GROUP_END), objStorage.Serial)
	return &reconcileParityGroups, nil
}
