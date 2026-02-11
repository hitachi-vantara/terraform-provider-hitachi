package sanstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewaymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	provisonerimpl "terraform-provider-hitachi/hitachi/storage/san/provisioner/impl"
	provisonermodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"
)

// GetPavAliases returns PAV alias entries filtered by CU number when provided.
func (psm *sanStorageManager) GetPavAliases(cuNumber *int) (*[]gatewaymodel.PavAlias, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisonermodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx, err: %v", err)
		return nil, err
	}

	provList, err := provObj.GetPavAliases(cuNumber)
	if err != nil {
		return nil, err
	}

	return provList, nil
}

func (psm *sanStorageManager) AssignPavAlias(baseLdevID int, aliasLdevIDs []int) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisonermodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx, err: %v", err)
		return err
	}

	if err := provObj.AssignPavAlias(baseLdevID, aliasLdevIDs); err != nil {
		log.WriteDebug("TFError| error in AssignPavAlias provisioner call, err: %v", err)
		return err
	}
	return nil
}

func (psm *sanStorageManager) UnassignPavAlias(aliasLdevIDs []int) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisonermodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx, err: %v", err)
		return err
	}

	if err := provObj.UnassignPavAlias(aliasLdevIDs); err != nil {
		log.WriteDebug("TFError| error in UnassignPavAlias provisioner call, err: %v", err)
		return err
	}
	return nil
}
