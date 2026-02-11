package sanstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	"terraform-provider-hitachi/hitachi/common/utils"
	gatewayimpl "terraform-provider-hitachi/hitachi/storage/san/gateway/impl"
	sangatewaymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/san/provisioner/message-catalog"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"
)

// GetStorageSystemInfo get storage system info
func (psm *sanStorageManager) GetStorageSystemInfo(detailInfoType ...string) (*sanmodel.StorageSystem, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_SYSTEM_BEGIN), psm.storageSetting.MgmtIP)
	objStorage := sangatewaymodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	log.WriteDebug("TFDebug| Storage Serial:%d, ManagementIP:%s\n", psm.storageSetting.Serial, psm.storageSetting.MgmtIP)
	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}
	var ssInfo *sangatewaymodel.StorageSystemInfo
	if len(detailInfoType) > 0 {
		ssInfo, err = gatewayObj.GetStorageSystemInfo(detailInfoType[0])
	} else {
		ssInfo, err = gatewayObj.GetStorageSystemInfo()
	}
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_STORAGE_SYSTEM_FAILED), objStorage.MgmtIP)
		return nil, err
	}

	log.WriteDebug("TFDebug| Storage setting: %+v\n", ssInfo)

	mgmtIP := ssInfo.IP
	if mgmtIP == "" {
		mgmtIP = ssInfo.Ctl1IP
	}

	ss := sanmodel.StorageSystem{
		StorageDeviceID:                    ssInfo.StorageDeviceID,
		Model:                              ssInfo.Model,
		SerialNumber:                       ssInfo.SerialNumber,
		MgmtIP:                             mgmtIP,
		IP:                                 ssInfo.IP,
		ControllerIP1:                      ssInfo.Ctl1IP,
		ControllerIP2:                      ssInfo.Ctl2IP,
		MicroVersion:                       ssInfo.DkcMicroVersion,
		DetailDkcMicroVersion:              ssInfo.DetailDkcMicroVersion,
		IsCompressionAccelerationAvailable: ssInfo.IsCompressionAccelerationAvailable,
		IsSecure:                           ssInfo.IsSecure,
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_SYSTEM_END), objStorage.MgmtIP)
	return &ss, nil
}

// GetStorageSystem to get Storage system information (storage + capacity info)
func (psm *sanStorageManager) GetStorageSystem(detailInfoType ...string) (*sanmodel.StorageSystem, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_SYSTEM_BEGIN), psm.storageSetting.MgmtIP)
	objStorage := sangatewaymodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	log.WriteDebug("TFDebug| Storage Serial:%d, ManagementIP:%s\n", psm.storageSetting.Serial, psm.storageSetting.MgmtIP)
	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	var ssInfo *sangatewaymodel.StorageSystemInfo
	if len(detailInfoType) > 0 {
		ssInfo, err = gatewayObj.GetStorageSystemInfo(detailInfoType[0])
	} else {
		ssInfo, err = gatewayObj.GetStorageSystemInfo()
	}
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_STORAGE_SYSTEM_FAILED), objStorage.MgmtIP)
		return nil, err
	}

	ssCapacity, err := gatewayObj.GetStorageCapacity()
	if err != nil {
		log.WriteDebug("TFError| error in GetStorageCapacity gateway call, err: %v", err)
		return nil, err
	}

	log.WriteDebug("TFDebug|SS: %+v , %+v\n", ssInfo, ssCapacity)

	mgmtIP := ssInfo.IP
	if mgmtIP == "" {
		mgmtIP = ssInfo.Ctl1IP
	}

	totalCapInMB := utils.ConvertSizeFromKbToMb(ssCapacity.Total.TotalCapacity)
	freeCapInMB := utils.ConvertSizeFromKbToMb(ssCapacity.Total.FreeSpace)
	usedCapInMB := totalCapInMB - freeCapInMB

	ss := sanmodel.StorageSystem{
		StorageDeviceID:   ssInfo.StorageDeviceID,
		Model:             ssInfo.Model,
		SerialNumber:      ssInfo.SerialNumber,
		MgmtIP:            mgmtIP,
		IP:                ssInfo.IP,
		ControllerIP1:     ssInfo.Ctl1IP,
		ControllerIP2:     ssInfo.Ctl2IP,
		MicroVersion:      ssInfo.DkcMicroVersion,
		TotalCapacityInMB: totalCapInMB,
		FreeCapacityInMB:  freeCapInMB,
		UsedCapacityInMB:  usedCapInMB,
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_SYSTEM_END), objStorage.MgmtIP)
	return &ss, nil
}
