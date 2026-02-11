package sanstorage

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"

	// "terraform-provider-hitachi/hitachi/common/utils"
	gatewayimpl "terraform-provider-hitachi/hitachi/storage/san/gateway/impl"
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
func (psm *sanStorageManager) GetRangeOfLuns(startLdevID int, endLdevID int, IsUndefinedLdev bool, filterOption string, detailInfoType string) (*[]sangatewaymodel.LogicalUnit, error) {
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
	reconcilerLogicalUnits, err := provObj.GetRangeOfLuns(startLdevID, endLdevID, IsUndefinedLdev, filterOption, detailInfoType)
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
	if isCreateLunOperation {
		// Mainframe vs block is driven ONLY by cylinder.
		isMainframe := lunRequest.Cylinder != nil
		var defaultdataReductionMode = "disabled"
		if !isMainframe {
			if lunRequest.DataReductionMode == nil {
				lunRequest.DataReductionMode = &defaultdataReductionMode
			}
		}

		// create lun
		log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_LUN_BEGIN), objStorage.Serial)

		createReq := sangatewaymodel.CreateLunRequestGwy{
			LdevID:                lunRequest.LdevID,
			Cylinder:              lunRequest.Cylinder,
			EmulationType:         lunRequest.EmulationType,
			IsTseVolume:           lunRequest.IsTseVolume,
			IsEseVolume:           lunRequest.IsEseVolume,
			ClprID:                lunRequest.ClprID,
			MpBladeID:             lunRequest.MpBladeID,
			Ssid:                  lunRequest.Ssid,
			BlockCapacity:         lunRequest.BlockCapacity,
			PoolID:                lunRequest.PoolID,
			ParityGroupID:         lunRequest.ParityGroupID,
			ExternalParityGroupID: lunRequest.ExternalParityGroupID,
			ByteFormatCapacity:    lunRequest.ByteFormatCapacity,
		}
		if !isMainframe {
			createReq.DataReductionMode = lunRequest.DataReductionMode
			createReq.IsDataReductionSharedVolumeEnabled = lunRequest.IsDataReductionSharedVolumeEnabled
			createReq.IsCompressionAccelerationEnabled = lunRequest.IsCompressionAccelerationEnabled
		}

		// Mainframe create does not support specifying ldevId.
		if isMainframe {
			createReq.LdevID = nil
			// Mainframe uses label rather than name.
			createReq.Label = lunRequest.Name
			// For parity-group mainframe, capacity is expressed via blockCapacity.
			// Do not send cylinder in that case.
			if createReq.ParityGroupID != nil && createReq.BlockCapacity != nil {
				createReq.Cylinder = nil
			}
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

// SetEseVolume toggles ESE enablement for a mainframe volume.
func (psm *sanStorageManager) SetEseVolume(ldevId int, isEseVolume bool) error {
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

	return provObj.SetEseVolume(ldevId, isEseVolume)
}

// FormatLdev invokes format operation through provisioner
func (psm *sanStorageManager) FormatLdev(ldevID int, req reconcilermodel.FormatLdevRequest) (*int, error) {
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

	gwReq := sangatewaymodel.FormatLdevRequestGwy{}
	if req.OperationType != nil || req.IsDataReductionForceFormat != nil {
		gwReq.Parameters = sangatewaymodel.FormatLdevParameters{}
		if req.OperationType != nil {
			gwReq.Parameters.OperationType = *req.OperationType
		}
		if req.IsDataReductionForceFormat != nil {
			gwReq.Parameters.IsDataReductionForceFormat = *req.IsDataReductionForceFormat
		}
	}

	lunId, err := provObj.FormatLdev(ldevID, gwReq)
	if err != nil {
		log.WriteDebug("TFError| error in FormatLdev provisioner call, err: %v", err)
		return nil, err
	}

	return lunId, nil
}

// StopAllVolumeFormat performs an appliance-wide stop-format via provisioner.
func (psm *sanStorageManager) StopAllVolumeFormat() error {
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
		log.WriteDebug("TFError| error in provisioner NewEx call, err: %v", err)
		return err
	}

	// Acquire resource-group lock via gateway before invoking stop-format.
	gwObj, gwErr := gatewayimpl.NewEx(sangatewaymodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	})
	if gwErr != nil {
		log.WriteDebug("TFError| error in gateway NewEx call, err: %v", gwErr)
		return gwErr
	}

	lockReq := sangatewaymodel.LockResourcesReq{}
	lockReq.Parameters.WaitTime = 10
	if err := gwObj.LockResources(lockReq); err != nil {
		log.WriteDebug("TFError| error in LockResources gateway call, err: %v", err)
		return err
	}

	// Ensure unlock attempt when finished.
	defer func() {
		if err := gwObj.UnlockResources(); err != nil {
			log.WriteDebug("TFError| error in UnlockResources gateway call, err: %v", err)
		}
	}()

	if err := provObj.StopAllVolumeFormat(); err != nil {
		log.WriteDebug("TFError| error in StopAllVolumeFormat provisioner call, err: %v", err)
		return err
	}

	log.WriteInfo("Successfully invoked stop-format on storage serial %d", psm.storageSetting.Serial)

	return nil
}

// BlockLun requests the provisioner to change the given LDEV status to blocked.
func (psm *sanStorageManager) BlockLun(ldevID int) error {
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
		log.WriteDebug("TFError| error in provisioner NewEx call, err: %v", err)
		return err
	}

	if err := provObj.BlockLun(ldevID); err != nil {
		log.WriteDebug("TFError| error in BlockLun provisioner call, err: %v", err)
		return err
	}

	return nil
}

// UnblockLun requests the provisioner to change the given LDEV status back to normal.
func (psm *sanStorageManager) UnblockLun(ldevID int) error {
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
		log.WriteDebug("TFError| error in provisioner NewEx call, err: %v", err)
		return err
	}

	if err := provObj.UnblockLun(ldevID); err != nil {
		log.WriteDebug("TFError| error in UnblockLun provisioner call, err: %v", err)
		return err
	}

	return nil
}
