package sanstorage

import (
	"fmt"
	"strconv"
	"strings"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/san/gateway/http"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
)

// GetHostGroup .
func (psm *sanStorageManager) GetHostGroup(portID string, hostGroupNumber int) (*sanmodel.HostGroupGwy, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var hostGroup sanmodel.HostGroupGwy
	apiSuf := fmt.Sprintf("objects/host-groups/%v,%v", portID, hostGroupNumber)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &hostGroup)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &hostGroup, nil
}

// GetHostGroupWwns .
func (psm *sanStorageManager) GetHostGroupWwns(portID string, hostGroupNumber int) (*[]sanmodel.HostWwnDetail, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var hgwwns sanmodel.HostWwnDetails
	apiSuf := fmt.Sprintf("objects/host-wwns?portId=%v&hostGroupNumber=%v", portID, hostGroupNumber)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &hgwwns)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &hgwwns.Data, nil
}

// GetHostGroupLuPaths .
func (psm *sanStorageManager) GetHostGroupLuPaths(portID string, hostGroupNumber int) (*[]sanmodel.HostLuPath, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var hglupaths sanmodel.HostLuPaths
	apiSuf := fmt.Sprintf("objects/luns?portId=%v&hostGroupNumber=%v", portID, hostGroupNumber)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &hglupaths)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &hglupaths.Data, nil
}

// GetHostGroupModeAndOptions .
func (psm *sanStorageManager) GetHostGroupModeAndOptions() (*sanmodel.HostModeAndOptions, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var hgModeAndOptions sanmodel.HostModeAndOptions
	apiSuf := "objects/supported-host-modes/instance"
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &hgModeAndOptions)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &hgModeAndOptions, nil
}

// SetHostGroupModeAndOptions .
func (psm *sanStorageManager) SetHostGroupModeAndOptions(portID string, hostGroupNumber int, reqBody sanmodel.SetHostModeAndOptions) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/host-groups/%v,%v", portID, hostGroupNumber)
	_, err := httpmethod.PatchCall(psm.storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	return nil
}

// CreateHostGroup .
func (psm *sanStorageManager) CreateHostGroup(reqBody sanmodel.CreateHostGroupReqGwy) (portid *string, hgnumber *int, err error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	affRes, err := httpmethod.PostCall(psm.storageSetting, "objects/host-groups", reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in CreateHostGroup - objects/host-groups API call, err: %v", err)
		return nil, nil, err
	}

	arr := strings.Split(*affRes, ",")
	portID := arr[0]
	hgNum, _ := strconv.Atoi(arr[1])
	log.WriteDebug("TFDebug| portID=%d hgNum=%v", portID, hgNum)
	return &portID, &hgNum, nil
}

// DeleteHostGroup .
func (psm *sanStorageManager) DeleteHostGroup(portID string, hostGroupNumber int) (err error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/host-groups/%v,%v", portID, hostGroupNumber)
	_, err = httpmethod.DeleteCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	return nil
}

// DeleteWwn .
func (psm *sanStorageManager) DeleteWwn(portID string, hostGroupNumber int, wwn string) (err error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/host-wwns/%s,%d,%s", portID, hostGroupNumber, wwn)
	_, err = httpmethod.DeleteCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	return nil
}

// GetAllHostGroups .
func (psm *sanStorageManager) GetAllHostGroups() (*sanmodel.HostGroups, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var hostGroups sanmodel.HostGroups
	apiSuf := "objects/host-groups"
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &hostGroups)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &hostGroups, nil
}

// AddWwnToHG .
func (psm *sanStorageManager) AddWwnToHG(reqBody sanmodel.AddWwnToHgReqGwy) (err error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	_, err = httpmethod.PostCall(psm.storageSetting, "objects/host-wwns", reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in AddWwnToHG - objects/host-wwns API call, err: %v", err)
		return err
	}
	return nil
}

// AddLdevToHG  is used to Setting the LU path
func (psm *sanStorageManager) AddLdevToHG(reqBody sanmodel.AddLdevToHgReqGwy) (err error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	_, err = httpmethod.PostCall(psm.storageSetting, "objects/luns", reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in AddLdevToHG - objects/luns API call, err: %v", err)
		return err
	}
	return nil
}

// RemoveLdevFromHG  is used for Deleting a LU path
func (psm *sanStorageManager) RemoveLdevFromHG(portID string, hostGroupNumber int, lunID int) (err error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/luns/%s,%d,%d", portID, hostGroupNumber, lunID)
	_, err = httpmethod.DeleteCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	return nil
}

// SetHostWwnNickName used to set nick name of wwn for hostgroup
func (psm *sanStorageManager) SetHostWwnNickName(portID string, hostGroupNumber int, hostWwn string, wwnNickname string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	reqBody := sanmodel.HostWwnNickName{WwnNickname: wwnNickname}
	apiSuf := fmt.Sprintf("objects/host-wwns/%s,%d,%s", portID, hostGroupNumber, hostWwn)
	_, err := httpmethod.PatchCall(psm.storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	return nil
}
