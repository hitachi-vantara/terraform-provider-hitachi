package sanstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/storage/san/gateway/impl"
	sangatewaymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
)

// GetPavAliases returns PAV alias entries filtered by CU number when provided.
func (psm *sanStorageManager) GetPavAliases(cuNumber *int) (*[]sangatewaymodel.PavAlias, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := sangatewaymodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	gwList, err := gatewayObj.GetPavAliases(cuNumber)
	if err != nil {
		log.WriteDebug("TFError| error in GetPavAliases gateway call, err: %v", err)
		return nil, err
	}

	return gwList, nil
}

func (psm *sanStorageManager) AssignPavAlias(baseLdevID int, aliasLdevIDs []int) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := sangatewaymodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	if err := gatewayObj.AssignPavAlias(baseLdevID, aliasLdevIDs); err != nil {
		log.WriteDebug("TFError| error in AssignPavAlias gateway call, err: %v", err)
		return err
	}
	return nil
}

func (psm *sanStorageManager) UnassignPavAlias(aliasLdevIDs []int) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := sangatewaymodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	if err := gatewayObj.UnassignPavAlias(aliasLdevIDs); err != nil {
		log.WriteDebug("TFError| error in UnassignPavAlias gateway call, err: %v", err)
		return err
	}
	return nil
}
