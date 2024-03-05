package sanstorage

import (
	// "encoding/json"
	// "errors"
	// "context"
	// "fmt"
	// "io/ioutil"

	// "time"

	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	"terraform-provider-hitachi/hitachi/common/utils"

	mc "terraform-provider-hitachi/hitachi/storage/san/provisioner/message-catalog"

	gatewayimpl "terraform-provider-hitachi/hitachi/storage/san/gateway/impl"
	sangatewaymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"

	"github.com/jinzhu/copier"
)

// GetLun get a lun by ldevId
func (psm *sanStorageManager) GetLun(ldevID int) (*sanmodel.LogicalUnit, error) {
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
	lunGateway, err := gatewayObj.GetLun(ldevID)
	if err != nil {
		log.WriteDebug("TFError| error in GetLun call, ldevID:%d, err: %v", ldevID, err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_LUN_FAILED), ldevID)
		return nil, err
	}

	// Converting Gatway to Provisionar
	lun := sanmodel.LogicalUnit{}
	err = copier.Copy(&lun, lunGateway)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to prov structure, err: %v", err)
		return nil, err
	}

	if lun.ByteFormatCapacity == "" {
		// does not exist, or in the process of being deleted
		log.WriteDebug("TFDebug| ByteFormatCapacity is blank\n")
		return &lun, nil
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
	return &lun, nil
}

// GetRangeOfLuns gets the desired luns based on range specified
func (psm *sanStorageManager) GetRangeOfLuns(startLdevID int, endLdevID int, IsUndefinedLdev bool) (*[]sanmodel.LogicalUnit, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	desiredLogicalUnits := []sanmodel.LogicalUnit{}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_LUN_RANGE_BEGIN), startLdevID, endLdevID)

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

// CreateLunInDynamicPoolWithLDevId creates lun in dynamic pool with ldevId
func (psm *sanStorageManager) CreateLunInDynamicPoolWithLDevId(ldevId int, sizeInGB float64, dynamicPool uint, dataReductionMode string) (*int, error) {

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

	dpool := int(dynamicPool)
	size := fmt.Sprintf("%dM", uint(sizeInGB*1024))
	reqBody := sangatewaymodel.CreateLunRequestGwy{
		LdevID:             &ldevId,
		PoolID:             &dpool,
		ByteFormatCapacity: size,
		DataReductionMode:  &dataReductionMode,
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_LUN_BEGIN), objStorage.Serial)
	lunId, err := gatewayObj.CreateLun(reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in CreateLun call, ldevID: %d, err: %v", ldevId, err)
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_LUN_FAILED), objStorage.Serial)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_LUN_END), objStorage.Serial)

	return lunId, nil
}

// CreateLunInParityGroupWithLDevId creates lun in parity group with ldevId
func (psm *sanStorageManager) CreateLunInParityGroupWithLDevId(ldevId int, sizeInGB float64, parityGroup string, dataReductionMode string) (*int, error) {

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

	size := fmt.Sprintf("%dM", uint(sizeInGB*1024))
	reqBody := sangatewaymodel.CreateLunRequestGwy{
		LdevID:             &ldevId,
		ParityGroupID:      &parityGroup,
		ByteFormatCapacity: size,
		DataReductionMode:  &dataReductionMode,
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_LUN_BEGIN), objStorage.Serial)
	lunId, err := gatewayObj.CreateLun(reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in CreateLun call, ldevID: %d, err: %v", ldevId, err)
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_LUN_FAILED), objStorage.Serial)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_LUN_END), objStorage.Serial)

	return lunId, nil
}

// CreateLunInDynamicPool creates lun in dynamic
func (psm *sanStorageManager) CreateLunInDynamicPool(sizeInGB float64, dynamicPool uint, dataReductionMode string) (*int, error) {

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

	dpool := int(dynamicPool)
	size := fmt.Sprintf("%dM", uint(sizeInGB*1024))
	reqBody := sangatewaymodel.CreateLunRequestGwy{
		PoolID:             &dpool,
		ByteFormatCapacity: size,
		DataReductionMode:  &dataReductionMode,
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_LUN_BEGIN), objStorage.Serial)
	lunId, err := gatewayObj.CreateLun(reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in CreateLun call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_LUN_FAILED), objStorage.Serial)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_LUN_END), objStorage.Serial)

	return lunId, nil
}

// CreateLunInParityGroup creates lun in parity group
func (psm *sanStorageManager) CreateLunInParityGroup(sizeInGB float64, parityGroup string, dataReductionMode string) (*int, error) {

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

	size := fmt.Sprintf("%dM", uint(sizeInGB*1024))
	reqBody := sangatewaymodel.CreateLunRequestGwy{
		ParityGroupID:      &parityGroup,
		ByteFormatCapacity: size,
		DataReductionMode:  &dataReductionMode,
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_LUN_BEGIN), objStorage.Serial)
	lunId, err := gatewayObj.CreateLun(reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in CreateLun call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_LUN_FAILED), objStorage.Serial)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_LUN_END), objStorage.Serial)

	return lunId, nil
}

// ExpandLun expands a lun by newSize
func (psm *sanStorageManager) ExpandLun(ldevId int, newSize float64) (*int, error) {

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

	size := fmt.Sprintf("%dM", uint(newSize*1024))
	reqBody := sangatewaymodel.ExpandLunRequestGwy{
		Parameters: sangatewaymodel.ExpandLunParameters{
			AdditionalByteFormatCapacity: size,
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
func (psm *sanStorageManager) UpdateLun(ldevId int, label *string, dataReductionMode *string) (*sanmodel.LogicalUnit, error) {
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

	if label != nil || dataReductionMode != nil {
		updReq := sangatewaymodel.UpdateLunRequestGwy{
			Label:             label,
			DataReductionMode: dataReductionMode,
		}
		log.WriteDebug("updReq: %+v", updReq)

		log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_LUN_BEGIN), ldevId, objStorage.Serial)
		_, err := gatewayObj.UpdateLun(updReq, ldevId)
		if err != nil {
			log.WriteDebug("TFError| error in UpdateLun call: %+v", err)
			log.WriteError(mc.GetMessage(mc.ERR_UPDATE_LUN_FAILED), ldevId, objStorage.Serial)
			return nil, err
		}
		log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_LUN_END), ldevId, objStorage.Serial)
	}

	lun, err := psm.GetLun(ldevId)
	if err != nil {
		log.WriteDebug("TFError| error in GetLun call, ldevID:%d, err: %v", ldevId, err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_LUN_FAILED), ldevId)
		return nil, err
	}

	return lun, nil
}
