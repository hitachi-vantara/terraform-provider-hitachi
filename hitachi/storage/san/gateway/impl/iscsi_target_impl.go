package sanstorage

import (
	"fmt"
	"strconv"
	"strings"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/san/gateway/http"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
)

// GetIscsiTarget .
func (psm *sanStorageManager) GetIscsiTarget(portID string, iscsiTargetNumber int) (*sanmodel.IscsiTargetGwy, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var iscsiTarget sanmodel.IscsiTargetGwy
	apiSuf := fmt.Sprintf("objects/host-groups/%v,%v", portID, iscsiTargetNumber)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &iscsiTarget)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &iscsiTarget, nil
}

// GetIscsiTargetByPortId
func (psm *sanStorageManager) GetIscsiTargetsByPortId(portID string) (*sanmodel.IscsiTargets, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var iscsiTargets sanmodel.IscsiTargets
	apiSuf := fmt.Sprintf("objects/host-groups?portId=%v", portID)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &iscsiTargets)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return &iscsiTargets, nil
}

// GetAllIscsiTargets
func (psm *sanStorageManager) GetAllIscsiTargets() (*sanmodel.IscsiTargets, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var iscsiTargets sanmodel.IscsiTargets
	apiSuf := "objects/host-groups"
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &iscsiTargets)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &iscsiTargets, nil
}

// GetIscsiTargetGroupLuPaths .
func (psm *sanStorageManager) GetIscsiTargetGroupLuPaths(portID string, iscsiTargetNumber int) (*[]sanmodel.IscsiTargetLuPath, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var itlupaths sanmodel.IscsiTargetLuPaths
	apiSuf := fmt.Sprintf("objects/luns?portId=%v&hostGroupNumber=%v", portID, iscsiTargetNumber)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &itlupaths)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &itlupaths.Data, nil
}

// GetIscsiNameInformation
func (psm *sanStorageManager) GetIscsiNameInformation(portID string, iscsiTargetNumber int) (*[]sanmodel.IscsiNameInformation, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var iscsiNameInfo sanmodel.AllIscsiNameInformation
	apiSuf := fmt.Sprintf("objects/host-iscsis?portId=%v&hostGroupNumber=%v", portID, iscsiTargetNumber)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &iscsiNameInfo)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &iscsiNameInfo.Data, nil
}

// CreateIscsiTarget
func (psm *sanStorageManager) CreateIscsiTarget(reqBody sanmodel.CreateIscsiTargetReq) (portId *string, itNum *int, err error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	affRes, err := httpmethod.PostCall(psm.storageSetting, "objects/host-groups", reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in CreateIscsiTarget - objects/host-groups API call, err: %v", err)
		return nil, nil, err
	}

	arr := strings.Split(*affRes, ",")
	portID := arr[0]
	iscsiTargetNum, _ := strconv.Atoi(arr[1])
	log.WriteDebug("TFDebug| portID=%d iscsiTargetNum=%v", portID, iscsiTargetNum)
	return &portID, &iscsiTargetNum, nil
}

// SetIScsiTargetHostModeAndHostModeOptions .
func (psm *sanStorageManager) SetIScsiTargetHostModeAndHostModeOptions(portID string, hostGroupNumber int, reqBody sanmodel.SetIscsiHostModeAndOptions) error {
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

// SetIscsiNameForIscsiTarget
func (psm *sanStorageManager) SetIscsiNameForIscsiTarget(reqBody sanmodel.SetIscsiNameReq) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	_, err := httpmethod.PostCall(psm.storageSetting, "objects/host-iscsis", reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in SetIscsiNameForIscsiTarget - objects/host-iscsis API call, err: %v", err)
		return err
	}

	return nil
}
// SetNicknameForIscsiName
func (psm *sanStorageManager) SetNicknameForIscsiName(portID string, iscsiTargetNumber int, iscsiName string, reqBody sanmodel.SetNicknameIscsiReq) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/host-iscsis/%v,%v,%v", portID, iscsiTargetNumber, iscsiName)
	_, err := httpmethod.PatchCall(psm.storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}

	return nil
}

// DeleteIscsiNameFromIscsiTarget
func (psm *sanStorageManager) DeleteIscsiNameFromIscsiTarget(portID string, iscsiTargetNumber int, iscsiName string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/host-iscsis/%v,%v,%v", portID, iscsiTargetNumber, iscsiName)
	_, err := httpmethod.DeleteCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}

	return nil
}

// DeleteIscsiTarget
func (psm *sanStorageManager) DeleteIscsiTarget(portID string, iscsiTargetNumber int) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/host-groups/%v,%v", portID, iscsiTargetNumber)
	_, err := httpmethod.DeleteCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	return nil
}
