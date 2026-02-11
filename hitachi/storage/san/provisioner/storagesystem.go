package sanstorage

import (
	// "encoding/json"
	// "errors"
	// "context"
	// "fmt"
	// "io/ioutil"
	// "strconv"
	// "time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	"terraform-provider-hitachi/hitachi/common/utils"
	// mc "terraform-provider-hitachi/hitachi/messagecatalog"
	sangateway "terraform-provider-hitachi/hitachi/storage/san/gateway"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/model"
)

func GetStorageSystemInfo(storageSetting sanmodel.StorageDeviceSettings) (*sanmodel.StorageSystem, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	ssInfo, err := sangateway.GetStorageSystemInfo(storageSetting)
	if err != nil {
		return nil, err
	}

	log.WriteDebug("SS: %+v\n", ssInfo)

	mgmtIP := ssInfo.SvpIP
	if mgmtIP == "" {
		mgmtIP = ssInfo.Ctl1IP
	}

	ss := sanmodel.StorageSystem{
		StorageDeviceID: ssInfo.StorageDeviceID,
		Model:           ssInfo.Model,
		SerialNumber:    ssInfo.SerialNumber,
		MgmtIP:          mgmtIP,
		SvpIP:           ssInfo.SvpIP,
		ControllerIP1:   ssInfo.Ctl1IP,
		ControllerIP2:   ssInfo.Ctl2IP,
		MicroVersion:    ssInfo.DkcMicroVersion,
	}

	return &ss, nil
}

func GetStorageSystem(storageSetting sanmodel.StorageDeviceSettings) (*sanmodel.StorageSystem, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	ssInfo, err := sangateway.GetStorageSystemInfo(storageSetting)
	if err != nil {
		return nil, err
	}

	ssCapacity, err := sangateway.GetStorageCapacity(storageSetting)
	if err != nil {
		return nil, err
	}

	log.WriteDebug("SS: %+v %+v\n", ssInfo, ssCapacity)

	mgmtIP := ssInfo.SvpIP
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
		SvpIP:             ssInfo.SvpIP,
		ControllerIP1:     ssInfo.Ctl1IP,
		ControllerIP2:     ssInfo.Ctl2IP,
		MicroVersion:      ssInfo.DkcMicroVersion,
		TotalCapacityInMB: totalCapInMB,
		FreeCapacityInMB:  freeCapInMB,
		UsedCapacityInMB:  usedCapInMB,
	}

	return &ss, nil
}
