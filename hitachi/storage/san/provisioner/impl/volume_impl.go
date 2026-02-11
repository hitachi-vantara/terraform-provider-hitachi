package sanstorage

import (
	// "fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	"terraform-provider-hitachi/hitachi/common/utils"

	mc "terraform-provider-hitachi/hitachi/storage/san/provisioner/message-catalog"

	gatewayimpl "terraform-provider-hitachi/hitachi/storage/san/gateway/impl"
	sangatewaymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
)

// GetLun get a lun by ldevId
func (psm *sanStorageManager) GetLun(ldevID int) (*sangatewaymodel.LogicalUnit, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_LUN_BEGIN), ldevID)
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
	lun, err := gatewayObj.GetLun(ldevID)
	if err != nil {
		log.WriteDebug("TFError| error in GetLun call, ldevID:%d, err: %v", ldevID, err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_LUN_FAILED), ldevID)
		return nil, err
	}

	if lun.ByteFormatCapacity == "" {
		// does not exist, or in the process of being deleted
		log.WriteDebug("TFDebug| ByteFormatCapacity is blank\n")
		return lun, nil
	}

	// calculate capacity
	totalInBytes, err := utils.ConvertSizeToBytes(lun.ByteFormatCapacity)
	if err != nil {
		log.WriteDebug("TFError| error in ConvertSizeToBytes, err: %v", err)
		return nil, err
	}
	totalInMB := utils.ConvertSizeFromBytesToMb(totalInBytes)

	blockSize := totalInBytes / lun.BlockCapacity

	usedInBytes := lun.NumOfUsedBlock * blockSize
	usedInMB := utils.ConvertSizeFromBytesToMb(usedInBytes)

	lun.TotalCapacityInMB = totalInMB
	lun.UsedCapacityInMB = usedInMB
	lun.FreeCapacityInMB = totalInMB - usedInMB

	log.WriteDebug("TFDebug| lun=%+v", lun)
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_LUN_END), ldevID)
	return lun, nil
}

// GetRangeOfLuns gets the desired luns based on range specified
func (psm *sanStorageManager) GetRangeOfLuns(startLdevID int, endLdevID int, IsUndefinedLdev bool, filterOption string, detailInfoType string) (*[]sangatewaymodel.LogicalUnit, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	desiredLogicalUnits := []sangatewaymodel.LogicalUnit{}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_LUN_RANGE_BEGIN), startLdevID, endLdevID)

	// If filterOption is specified, use bulk API call through gateway
	if filterOption != "" || detailInfoType != "" {
		objStorage := sangatewaymodel.StorageDeviceSettings{
			Serial:   psm.storageSetting.Serial,
			Username: psm.storageSetting.Username,
			Password: psm.storageSetting.Password,
			MgmtIP:   psm.storageSetting.MgmtIP,
		}

		gatewayObj, err := gatewayimpl.NewEx(objStorage)
		if err != nil {
			log.WriteDebug("TFError| error in gateway NewEx call, err: %v", err)
			return nil, err
		}

		// Use bulk API with new parameters
		bulkLuns, err := gatewayObj.GetRangeOfLunsWithOptions(startLdevID, endLdevID, IsUndefinedLdev, filterOption, detailInfoType)
		if err != nil {
			log.WriteDebug("TFError| error in GetRangeOfLunsWithOptions call, err: %v", err)
			return nil, err
		}

		// The gateway already returns the provider's LogicalUnit model.
		// Preserve all fields (status/cylinder/attributes/ports/virtual* etc.).
		for _, gatewayLun := range *bulkLuns {
			desiredLogicalUnits = append(desiredLogicalUnits, gatewayLun)
		}
	} else {
		// Fallback to individual calls for backward compatibility
		for i := startLdevID; i <= endLdevID; i++ {
			logicalUnit, err := psm.GetLun(i)
			if err != nil {
				log.WriteDebug("TFError| error in GetRangeOfLuns call, err: %v", err)
				log.WriteError(mc.GetMessage(mc.ERR_GET_LUN_FAILED), i)
				return nil, err
			}
			if IsUndefinedLdev {
				if logicalUnit.EmulationType == "NOT DEFINED" {
					desiredLogicalUnits = append(desiredLogicalUnits, *logicalUnit)
				}
			} else {
				if logicalUnit.EmulationType != "NOT DEFINED" {
					desiredLogicalUnits = append(desiredLogicalUnits, *logicalUnit)
				}
			}
		}
	}

	if len(desiredLogicalUnits) == 0 {
		log.WriteDebug("TFDebug| GetRangeOfLuns - No luns found based on given criteria.")
	} else {
		log.WriteDebug("TFDebug| GetRangeOfLuns - Found luns:  %+v", desiredLogicalUnits)
		log.WriteInfo(mc.GetMessage(mc.INFO_GET_LUN_RANGE_END), startLdevID, endLdevID)
	}

	return &desiredLogicalUnits, nil
}

// GetUndefinedLun this will be internal function to get UndefinedLun
func (psm *sanStorageManager) GetUndefinedLun(numberOfUndefineLun int) ([]int, error) {
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
	// Get all lun
	luns, err := gatewayObj.GetAllLun()
	if err != nil {
		log.WriteDebug("TFError| error in GetAllLun call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_LUN_FAILED))
		return nil, err
	}

	var arrUndefiedLdev []int
	count := 0
	for _, value := range luns.ListOfLun {
		//"emulationType" : "NOT DEFINED"
		if value.EmulationType == "NOT DEFINED" {
			log.WriteDebug("TFDebug| value.LdevID %d", value.LdevID)
			arrUndefiedLdev = append(arrUndefiedLdev, value.LdevID)
			count++
			// If required undefined lun fetch then return
			if numberOfUndefineLun == count {
				break
			}
		}
	}
	return arrUndefiedLdev, nil
}

// CreateLun
func (psm *sanStorageManager) CreateLun(reqBody sangatewaymodel.CreateLunRequestGwy) (*int, error) {

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

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_LUN_BEGIN), objStorage.Serial)
	lunId, err := gatewayObj.CreateLun(reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in CreateLun call, ldevID: %d, err: %v", reqBody.LdevID, err)
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_LUN_FAILED), objStorage.Serial)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_LUN_END), objStorage.Serial)

	return lunId, nil
}

// ExpandLun expands a lun by newSize
func (psm *sanStorageManager) ExpandLun(ldevId int, newSize string) (*int, error) {

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

	// size := fmt.Sprintf("%dG", newSize)
	reqBody := sangatewaymodel.ExpandLunRequestGwy{
		Parameters: sangatewaymodel.ExpandLunParameters{
			AdditionalByteFormatCapacity: newSize,
		},
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_EXPAND_LUN_BEGIN), ldevId, objStorage.Serial)
	lunId, err := gatewayObj.ExpandLun(reqBody, ldevId)
	if err != nil {
		log.WriteDebug("TFError| error in ExpandLun call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_EXPAND_LUN_FAILED), ldevId, objStorage.Serial)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_EXPAND_LUN_END), ldevId, objStorage.Serial)

	return lunId, nil
}

// FormatLdev invokes the LDEV format action via gateway
func (psm *sanStorageManager) FormatLdev(ldevId int, req sangatewaymodel.FormatLdevRequestGwy) (*int, error) {
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

	log.WriteInfo("Invoking format for LDEV %d", ldevId)
	lunId, err := gatewayObj.FormatLdev(req, ldevId)
	if err != nil {
		log.WriteDebug("TFError| error in FormatLdev call, err: %v", err)
		return nil, err
	}

	return lunId, nil
}

// BlockLun requests the gateway to change the LDEV status to blocked ('blk')
func (psm *sanStorageManager) BlockLun(ldevId int) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	gatewayObj, err := gatewayimpl.NewEx(sangatewaymodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	})
	if err != nil {
		log.WriteDebug("TFError| error in gateway NewEx call, err: %v", err)
		return err
	}

	if err := gatewayObj.BlockLun(ldevId); err != nil {
		log.WriteDebug("TFError| error in BlockLun gateway call, err: %v", err)
		return err
	}

	return nil
}

// UnblockLun requests the gateway to change the LDEV status back to normal ('nml')
func (psm *sanStorageManager) UnblockLun(ldevId int) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	gatewayObj, err := gatewayimpl.NewEx(sangatewaymodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	})
	if err != nil {
		log.WriteDebug("TFError| error in gateway NewEx call, err: %v", err)
		return err
	}

	if err := gatewayObj.UnblockLun(ldevId); err != nil {
		log.WriteDebug("TFError| error in UnblockLun gateway call, err: %v", err)
		return err
	}

	return nil
}

// StopAllVolumeFormat invokes the gateway to stop all ongoing/normal volume formats.
func (psm *sanStorageManager) StopAllVolumeFormat() error {
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
		log.WriteDebug("TFError| error in gateway NewEx call, err: %v", err)
		return err
	}

	if err := gatewayObj.StopAllVolumeFormat(); err != nil {
		log.WriteDebug("TFError| error in StopAllVolumeFormat gateway call, err: %v", err)
		return err
	}

	return nil
}

// DeleteLun delete a lun by ldevId
func (psm *sanStorageManager) DeleteLun(ldevId int) error {

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

	lun, err := gatewayObj.GetLun(ldevId)
	if err != nil {
		log.WriteDebug("TFError| error in GetLun call, ldevID:%d, err: %v", ldevId, err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_LUN_FAILED), ldevId)
		return err
	}

	capacitySaving := true
	if lun.DataReductionMode == "" || lun.DataReductionMode == "disabled" {
		capacitySaving = false
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_LUN_BEGIN), ldevId, objStorage.Serial)
	err = gatewayObj.DeleteLun(ldevId, capacitySaving)
	if err != nil {
		log.WriteDebug("TFError| error in DeleteLun call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_LUN_FAILED), ldevId, objStorage.Serial)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_LUN_END), ldevId, objStorage.Serial)

	return nil
}

// UpdateLun updates a lun
func (psm *sanStorageManager) UpdateLun(ldevId int, updReq sangatewaymodel.UpdateLunRequestGwy) (*int, error) {
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

	log.WriteDebug("updReq: %+v", updReq)

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_LUN_BEGIN), ldevId, objStorage.Serial)
	_, err = gatewayObj.UpdateLun(updReq, ldevId)
	if err != nil {
		log.WriteDebug("TFError| error in UpdateLun call: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_UPDATE_LUN_FAILED), ldevId, objStorage.Serial)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_LUN_END), ldevId, objStorage.Serial)

	return &ldevId, nil
}

// SetEseVolume toggles ESE enablement for a mainframe volume.
func (psm *sanStorageManager) SetEseVolume(ldevId int, isEseVolume bool) error {
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

	return gatewayObj.SetEseVolume(ldevId, isEseVolume)
}
