package sanstorage

import (
	"fmt"
	"strconv"
	"strings"

	// "time"
	// "encoding/json"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	// "terraform-provider-hitachi/hitachi/common/utils"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/model"
)

func GetHostGroup(storageSetting sanmodel.StorageDeviceSettings, portID string, hostGroupNumber int) (*sanmodel.HostGroupGwy, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var hostGroup sanmodel.HostGroupGwy
	apiSuf := fmt.Sprintf("objects/host-groups/%v,%v", portID, hostGroupNumber)
	err := GetCall(storageSetting, apiSuf, &hostGroup)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}
	return &hostGroup, nil
}

func GetHostGroupWwns(storageSetting sanmodel.StorageDeviceSettings, portID string, hostGroupNumber int) (*[]sanmodel.HostWwnDetail, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var hgwwns sanmodel.HostWwnDetails
	apiSuf := fmt.Sprintf("objects/host-wwns?portId=%v&hostGroupNumber=%v", portID, hostGroupNumber)
	err := GetCall(storageSetting, apiSuf, &hgwwns)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}
	return &hgwwns.Data, nil
}

func GetHostGroupLuPaths(storageSetting sanmodel.StorageDeviceSettings, portID string, hostGroupNumber int) (*[]sanmodel.HostLuPath, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var hglupaths sanmodel.HostLuPaths
	apiSuf := fmt.Sprintf("objects/luns?portId=%v&hostGroupNumber=%v", portID, hostGroupNumber)
	err := GetCall(storageSetting, apiSuf, &hglupaths)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}
	return &hglupaths.Data, nil
}

func GetHostGroupModeAndOptions(storageSetting sanmodel.StorageDeviceSettings) (*sanmodel.HostModeAndOptions, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var hgModeAndOptions sanmodel.HostModeAndOptions
	apiSuf := fmt.Sprintf("objects/supported-host-modes/instance")
	err := GetCall(storageSetting, apiSuf, &hgModeAndOptions)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}
	return &hgModeAndOptions, nil
}

func CreateHostGroup(storageSetting sanmodel.StorageDeviceSettings, reqBody interface{}) (*string, *int, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	affRes, err := PostCall(storageSetting, "objects/host-groups", reqBody)
	if err != nil {
		return nil, nil, err
	}

	arr := strings.Split(*affRes, ",")
	portID := arr[0]
	hgNum, _ := strconv.Atoi(arr[1])
	fmt.Println("portID=%d hgNum=%v\n", portID, hgNum)
	return &portID, &hgNum, nil
}

func AddWwnToHG(storageSetting sanmodel.StorageDeviceSettings, reqBody interface{}) (err error) {
	_, err = PostCall(storageSetting, "objects/host-wwns", reqBody)
	return
}

func AddLdevToHG(storageSetting sanmodel.StorageDeviceSettings, reqBody interface{}) (err error) {
	_, err = PostCall(storageSetting, "objects/luns", reqBody)
	return
}
