package sanstorage

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	// "terraform-provider-hitachi/hitachi/common/utils"
	sangatewaymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	provisonerimpl "terraform-provider-hitachi/hitachi/storage/san/provisioner/impl"
	provisonermodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"
	mc "terraform-provider-hitachi/hitachi/storage/san/reconciler/message-catalog"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"

	// "github.com/jinzhu/copier"
)

// GetLun get Storage Lun information
func (psm *sanStorageManager) GetLun(ldevID int) (*sangatewaymodel.LogicalUnit, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_LUN_BEGIN), ldevID)
	objStorage := provisonermodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}
	reconcilerLogicalUnit, err := provObj.GetLun(ldevID)
	if err != nil {
		log.WriteDebug("TFError| error in GetLun provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_LUN_FAILED), ldevID)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_LUN_END), ldevID)
	return reconcilerLogicalUnit, nil
}

// GetRangeOfLuns gets the desired luns based on range specified
func (psm *sanStorageManager) GetRangeOfLuns(startLdevID int, endLdevID int, IsUndefinedLdev bool) (*[]sangatewaymodel.LogicalUnit, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_LUN_RANGE_BEGIN), startLdevID, endLdevID)
	objStorage := provisonermodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}
	reconcilerLogicalUnits, err := provObj.GetRangeOfLuns(startLdevID, endLdevID, IsUndefinedLdev)
	if err != nil {
		log.WriteDebug("TFError| error in GetRangeOfLuns provisioner call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_LUN_RANGE_END), startLdevID, endLdevID)
	return reconcilerLogicalUnits, nil

}

// SetLun will create or update the lun
func (psm *sanStorageManager) SetLun(lunRequest *reconcilermodel.LunRequest) (*int, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisonermodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	var provLogicalUnit *sangatewaymodel.LogicalUnit

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	isCreateLunOperation := false
	if lunRequest.LdevID == nil {
		isCreateLunOperation = true
	} else {
		// is it existing?
		provLogicalUnit, err = provObj.GetLun(*lunRequest.LdevID)
		if err != nil {
			log.WriteDebug("TFError| error in GetLun reconciler call, err: %v", err)
			log.WriteError(mc.GetMessage(mc.ERR_GET_LUN_FAILED), *lunRequest.LdevID)
			return nil, err
		}
		if provLogicalUnit.EmulationType == "NOT DEFINED" {
			isCreateLunOperation = true
		}
	}

	var ldev *int
	var defaultdataReductionMode = "disabled"
	if isCreateLunOperation {
		if lunRequest.DataReductionMode == nil {
			lunRequest.DataReductionMode = &defaultdataReductionMode
		}
		// create lun
		log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_LUN_BEGIN), objStorage.Serial)

		createReq := sangatewaymodel.CreateLunRequestGwy{
			LdevID:                             lunRequest.LdevID,
			PoolID:                             lunRequest.PoolID,
			ParityGroupID:                      lunRequest.ParityGroupID,
			ExternalParityGroupID:              lunRequest.ExternalParityGroupID,
			ByteFormatCapacity:                 lunRequest.ByteFormatCapacity,
			DataReductionMode:                  lunRequest.DataReductionMode,
			IsDataReductionSharedVolumeEnabled: lunRequest.IsDataReductionSharedVolumeEnabled,
			IsCompressionAccelerationEnabled:   lunRequest.IsCompressionAccelerationEnabled,
		}

		ldev, err = provObj.CreateLun(createReq)
		if err != nil {
			log.WriteDebug("TFError| error in CreateLun reconciler call, err: %v", err)
			log.WriteError(mc.GetMessage(mc.ERR_CREATE_LUN_FAILED), objStorage.Serial)
			return nil, err
		}

		log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_LUN_END), *ldev, objStorage.Serial)

		if lunRequest.Name != nil && *lunRequest.Name != "" {
			log.WriteDebug("Updating lun name/label to %s", *lunRequest.Name)
			lunUpdateRequest := reconcilermodel.UpdateLunRequest{
				LdevID: ldev,
				Name:   lunRequest.Name,
			}

			_, err := psm.UpdateLun(&lunUpdateRequest)
			if err != nil {
				return nil, err
			}
		}

		return ldev, nil

	} else {
		// ldevid is already used, error out
		log.WriteDebug("The specified LDEV_ID has already been allocated and is no longer available")
		return nil, fmt.Errorf("the specified LDEV_ID %d has already been allocated and is no longer available", *lunRequest.LdevID)
	}
}

// DeleteLun delete a lun by ldevId
func (psm *sanStorageManager) DeleteLun(ldevID int) error {
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
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_LUN_BEGIN), ldevID, objStorage.Serial)

	err = provObj.DeleteLun(ldevID)
	if err != nil {
		log.WriteDebug("TFError| error in DeleteLun call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_LUN_FAILED), ldevID, objStorage.Serial)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_LUN_END), ldevID, objStorage.Serial)
	return nil
}

// UpdateLun updates a lun
func (psm *sanStorageManager) UpdateLun(lunUpdateRequest *reconcilermodel.UpdateLunRequest) (*int, error) {
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
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	if lunUpdateRequest.Name != nil || lunUpdateRequest.DataReductionMode != nil || lunUpdateRequest.IsAluaEnabled != nil {
		log.WriteInfo("Updating Name/Data Reduction Mode for LUN %v", lunUpdateRequest.LdevID)
		log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_LUN_BEGIN), lunUpdateRequest.LdevID, objStorage.Serial)
		updReq := sangatewaymodel.UpdateLunRequestGwy{
			Label:             lunUpdateRequest.Name,
			DataReductionMode: lunUpdateRequest.DataReductionMode,
			IsAluaEnabled:     lunUpdateRequest.IsAluaEnabled,
		}
		_, err := provObj.UpdateLun(*lunUpdateRequest.LdevID, updReq)
		if err != nil {
			log.WriteDebug("TFError| error in UpdateLun reconciler call, err: %v", err)
			log.WriteError(mc.GetMessage(mc.ERR_UPDATE_LUN_FAILED), lunUpdateRequest.LdevID, objStorage.Serial)
			return nil, err
		}
		log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_LUN_END), lunUpdateRequest.LdevID, objStorage.Serial)
	}

	if lunUpdateRequest.DataReductionProcessMode != nil {
		log.WriteInfo("Updating Data Reduction Process Mode to %v for LUN %v", *lunUpdateRequest.DataReductionProcessMode, lunUpdateRequest.LdevID)
		log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_LUN_BEGIN), lunUpdateRequest.LdevID, objStorage.Serial)
		updReq := sangatewaymodel.UpdateLunRequestGwy{
			DataReductionProcessMode: lunUpdateRequest.DataReductionProcessMode,
		}
		_, err := provObj.UpdateLun(*lunUpdateRequest.LdevID, updReq)
		if err != nil {
			log.WriteDebug("TFError| error in UpdateLun reconciler call, err: %v", err)
			log.WriteError(mc.GetMessage(mc.ERR_UPDATE_LUN_FAILED), lunUpdateRequest.LdevID, objStorage.Serial)
			return nil, err
		}
		log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_LUN_END), lunUpdateRequest.LdevID, objStorage.Serial)
	}

	if lunUpdateRequest.IsCompressionAccelerationEnabled != nil {
		log.WriteInfo("Updating Compression Acceleration Enabled to %v for LUN %v", *lunUpdateRequest.IsCompressionAccelerationEnabled, lunUpdateRequest.LdevID)
		compAccEnabled := *lunUpdateRequest.IsCompressionAccelerationEnabled
		log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_LUN_BEGIN), lunUpdateRequest.LdevID, objStorage.Serial)
		updReq := sangatewaymodel.UpdateLunRequestGwy{
			IsCompressionAccelerationEnabled: &compAccEnabled,
		}
		_, err := provObj.UpdateLun(*lunUpdateRequest.LdevID, updReq)
		if err != nil {
			log.WriteDebug("TFError| error in UpdateLun reconciler call, err: %v", err)
			log.WriteError(mc.GetMessage(mc.ERR_UPDATE_LUN_FAILED), lunUpdateRequest.LdevID, objStorage.Serial)
			return nil, err
		}
		log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_LUN_END), lunUpdateRequest.LdevID, objStorage.Serial)
	}

	if lunUpdateRequest.ByteFormatCapacity != nil {
		log.WriteInfo(mc.GetMessage(mc.INFO_EXPAND_LUN_BEGIN), lunUpdateRequest.LdevID, objStorage.Serial)
		_, err := provObj.ExpandLun(*lunUpdateRequest.LdevID, *lunUpdateRequest.ByteFormatCapacity)
		if err != nil {
			log.WriteDebug("TFError| error in ExpandLun call, err: %v", err)
			log.WriteError(mc.GetMessage(mc.ERR_EXPAND_LUN_FAILED), lunUpdateRequest.LdevID, objStorage.Serial)
			return nil, err
		}
		log.WriteInfo(mc.GetMessage(mc.INFO_EXPAND_LUN_END), lunUpdateRequest.LdevID, objStorage.Serial)
	}

	lunID := lunUpdateRequest.LdevID
	return lunID, nil
}
