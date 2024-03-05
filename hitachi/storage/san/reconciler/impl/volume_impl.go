package sanstorage

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	"terraform-provider-hitachi/hitachi/common/utils"
	provisonerimpl "terraform-provider-hitachi/hitachi/storage/san/provisioner/impl"
	provisonermodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"
	mc "terraform-provider-hitachi/hitachi/storage/san/reconciler/message-catalog"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"

	"github.com/jinzhu/copier"
)

// GetLun get Storage Lun information
func (psm *sanStorageManager) GetLun(ldevID int) (*sanmodel.LogicalUnit, error) {
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
	provLogicalUnit, err := provObj.GetLun(ldevID)
	if err != nil {
		log.WriteDebug("TFError| error in GetLun provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_LUN_FAILED), ldevID)
		return nil, err
	}
	// Converting Prov to Reconciler
	reconcilerLogicalUnit := sanmodel.LogicalUnit{}
	err = copier.Copy(&reconcilerLogicalUnit, provLogicalUnit)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_LUN_END), ldevID)
	return &reconcilerLogicalUnit, nil
}

// GetRangeOfLuns gets the desired luns based on range specified
func (psm *sanStorageManager) GetRangeOfLuns(startLdevID int, endLdevID int, IsUndefinedLdev bool) (*[]sanmodel.LogicalUnit, error) {
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
	provLogicalUnits, err := provObj.GetRangeOfLuns(startLdevID, endLdevID, IsUndefinedLdev)
	if err != nil {
		log.WriteDebug("TFError| error in GetRangeOfLuns provisioner call, err: %v", err)
		return nil, err
	}

	// Converting Prov to Reconciler
	reconcilerLogicalUnits := []sanmodel.LogicalUnit{}
	err = copier.Copy(&reconcilerLogicalUnits, provLogicalUnits)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_LUN_RANGE_END), startLdevID, endLdevID)
	return &reconcilerLogicalUnits, nil

}

// SetLun will create or expand the lun
func (psm *sanStorageManager) SetLun(lunRequest *sanmodel.LunRequest) (*sanmodel.LogicalUnit, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisonermodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	var provLogicalUnit *provisonermodel.LogicalUnit

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	isCreateLunOperation := false
	if lunRequest.LdevID == nil {
		isCreateLunOperation = true
	} else {
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
		if lunRequest.LdevID != nil && lunRequest.PoolID != nil {
			ldev, err = provObj.CreateLunInDynamicPoolWithLDevId(*lunRequest.LdevID, lunRequest.CapacityInGB, uint(*lunRequest.PoolID), *lunRequest.DataReductionMode)
			if err != nil {
				log.WriteDebug("TFError| error in CreateLunInDynamicPoolWithLDevId reconciler call, err: %v", err)
				log.WriteError(mc.GetMessage(mc.ERR_CREATE_LUN_FAILED), objStorage.Serial)
				return nil, err
			}
		} else if lunRequest.LdevID != nil && lunRequest.ParityGroupID != nil {
			ldev, err = provObj.CreateLunInParityGroupWithLDevId(*lunRequest.LdevID, lunRequest.CapacityInGB, *lunRequest.ParityGroupID, *lunRequest.DataReductionMode)
			if err != nil {
				log.WriteDebug("TFError| error in CreateLunInParityGroupWithLDevId reconciler call, err: %v", err)
				log.WriteError(mc.GetMessage(mc.ERR_CREATE_LUN_FAILED), objStorage.Serial)
				return nil, err
			}
		} else if lunRequest.LdevID == nil && lunRequest.PoolID != nil {
			ldev, err = provObj.CreateLunInDynamicPool(lunRequest.CapacityInGB, uint(*lunRequest.PoolID), *lunRequest.DataReductionMode)
			if err != nil {
				log.WriteDebug("TFError| error in CreateLunInDynamicPool reconciler call, err: %v", err)
				log.WriteError(mc.GetMessage(mc.ERR_CREATE_LUN_FAILED), objStorage.Serial)
				return nil, err
			}
		} else if lunRequest.LdevID == nil && lunRequest.ParityGroupID != nil {
			ldev, err = provObj.CreateLunInParityGroup(lunRequest.CapacityInGB, *lunRequest.ParityGroupID, *lunRequest.DataReductionMode)
			if err != nil {
				log.WriteDebug("TFError| error in CreateLunInParityGroup reconciler call, err: %v", err)
				log.WriteError(mc.GetMessage(mc.ERR_CREATE_LUN_FAILED), objStorage.Serial)
				return nil, err
			}
		}
		log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_LUN_END), *ldev, objStorage.Serial)
	} else {
		// expand lun
		ldev = lunRequest.LdevID

		if uint64(lunRequest.CapacityInGB*1024) < provLogicalUnit.TotalCapacityInMB {
			msg := "cannot shrink the capacity of lun. Please enter more capacity than current one."
			return nil, fmt.Errorf(msg)
		}

		if uint64(lunRequest.CapacityInGB*1024) == provLogicalUnit.TotalCapacityInMB {
			msg := "cannot expand the lun as current and desired capacity is same."
			return nil, fmt.Errorf(msg)
		}

		if lunRequest.PoolID != nil && provLogicalUnit.PoolID != *lunRequest.PoolID {
			msg := "pool_id cannot be changed."
			return nil, fmt.Errorf(msg)
		}

		if lunRequest.ParityGroupID != nil {
			if !utils.IsParityGroupPresent(*lunRequest.ParityGroupID, provLogicalUnit.ParityGroupId) {
				msg := "paritygroup_id cannot be changed."
				return nil, fmt.Errorf(msg)
			}
		}

		additionalCapacityInMB := uint64(lunRequest.CapacityInGB*1024) - (provLogicalUnit.TotalCapacityInMB)
		log.WriteDebug("TFDebug| additionalCapacityInMB %d: ", additionalCapacityInMB)

		log.WriteInfo(mc.GetMessage(mc.INFO_EXPAND_LUN_BEGIN), *ldev, objStorage.Serial)
		_, err := provObj.ExpandLun(*ldev, float64(additionalCapacityInMB/1024))
		if err != nil {
			log.WriteDebug("TFError| error in ExpandLun reconciler call, err: %v", err)
			log.WriteError(mc.GetMessage(mc.ERR_EXPAND_LUN_FAILED), *ldev, objStorage.Serial)
			return nil, err
		}
		log.WriteDebug("TFDebug| Lun %d expanded successfully.", *ldev)
		log.WriteInfo(mc.GetMessage(mc.INFO_EXPAND_LUN_END), *ldev, objStorage.Serial)
	}

	if lunRequest.Name != nil || lunRequest.DataReductionMode != nil {
		log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_LUN_BEGIN), *ldev, objStorage.Serial)
		_, err := provObj.UpdateLun(*ldev, lunRequest.Name, lunRequest.DataReductionMode)
		if err != nil {
			log.WriteDebug("TFError| error in UpdateLun reconciler call, err: %v", err)
			log.WriteError(mc.GetMessage(mc.ERR_UPDATE_LUN_FAILED), *ldev, objStorage.Serial)
		}
		log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_LUN_END), *ldev, objStorage.Serial)
	}

	provLogicalUnit, err = provObj.GetLun(*ldev)
	if err != nil {
		log.WriteDebug("TFError| error in GetLun reconciler call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_LUN_FAILED), *lunRequest.LdevID)
		return nil, err
	}

	// Converting Prov to Reconciler
	reconcilerLogicalUnit := sanmodel.LogicalUnit{}
	err = copier.Copy(&reconcilerLogicalUnit, provLogicalUnit)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}

	return &reconcilerLogicalUnit, nil
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
func (psm *sanStorageManager) UpdateLun(lunUpdateRequest *sanmodel.UpdateLunRequest) (*sanmodel.LogicalUnit, error) {
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

	if lunUpdateRequest.Name != nil || lunUpdateRequest.DataReductionMode != nil {
		log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_LUN_BEGIN), lunUpdateRequest.LdevID, objStorage.Serial)
		_, err := provObj.UpdateLun(*lunUpdateRequest.LdevID, lunUpdateRequest.Name, lunUpdateRequest.DataReductionMode)
		if err != nil {
			log.WriteDebug("TFError| error in UpdateLun reconciler call, err: %v", err)
			log.WriteError(mc.GetMessage(mc.ERR_UPDATE_LUN_FAILED), lunUpdateRequest.LdevID, objStorage.Serial)
		}
		log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_LUN_END), lunUpdateRequest.LdevID, objStorage.Serial)
	}

	if lunUpdateRequest.CapacityInGB != nil {
		log.WriteInfo(mc.GetMessage(mc.INFO_EXPAND_LUN_BEGIN), lunUpdateRequest.LdevID, objStorage.Serial)

		_, err := provObj.ExpandLun(*lunUpdateRequest.LdevID, *lunUpdateRequest.CapacityInGB)
		if err != nil {
			log.WriteDebug("TFError| error in ExpandLun call, err: %v", err)
			log.WriteError(mc.GetMessage(mc.ERR_EXPAND_LUN_FAILED), lunUpdateRequest.LdevID, objStorage.Serial)
			return nil, err
		}
		log.WriteInfo(mc.GetMessage(mc.INFO_EXPAND_LUN_END), lunUpdateRequest.LdevID, objStorage.Serial)
	}

	lun, err := psm.GetLun(*lunUpdateRequest.LdevID)
	if err != nil {
		log.WriteDebug("TFError| error in GetLun reconciler call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_LUN_FAILED), lunUpdateRequest.LdevID)
		return nil, err
	}

	return lun, nil
}
