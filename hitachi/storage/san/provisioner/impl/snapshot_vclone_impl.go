package sanstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/storage/san/gateway/impl"
	model "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/san/provisioner/message-catalog"
)

// --- VIRTUAL CLONE ---

func (psm *sanStorageManager) CreateSnapshotVClone(pvolLdevID int, muNumber int) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_VCLONE_BEGIN), pvolLdevID, muNumber, psm.storageSetting.Serial)

	objStorage := model.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return "", err
	}

	jobID, err := gatewayObj.CreateSnapshotVClone(pvolLdevID, muNumber)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_VCLONE_FAILED), pvolLdevID, muNumber, psm.storageSetting.Serial)
		return "", err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_VCLONE_END), pvolLdevID, muNumber, psm.storageSetting.Serial)
	return jobID, nil
}

func (psm *sanStorageManager) ConvertSnapshotVClone(pvolLdevID int, muNumber int) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_CONVERT_VCLONE_BEGIN), pvolLdevID, muNumber, psm.storageSetting.Serial)

	objStorage := model.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return "", err
	}

	jobID, err := gatewayObj.ConvertSnapshotVClone(pvolLdevID, muNumber)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_CONVERT_VCLONE_FAILED), pvolLdevID, muNumber, psm.storageSetting.Serial)
		return "", err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CONVERT_VCLONE_END), pvolLdevID, muNumber, psm.storageSetting.Serial)
	return jobID, nil
}

func (psm *sanStorageManager) RestoreSnapshotFromVClone(pvolLdevID int, muNumber int) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_RESTORE_VCLONE_BEGIN), pvolLdevID, muNumber, psm.storageSetting.Serial)

	objStorage := model.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return "", err
	}

	jobID, err := gatewayObj.RestoreSnapshotFromVClone(pvolLdevID, muNumber)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_RESTORE_VCLONE_FAILED), pvolLdevID, muNumber, psm.storageSetting.Serial)
		return "", err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_RESTORE_VCLONE_END), pvolLdevID, muNumber, psm.storageSetting.Serial)
	return jobID, nil
}

func (psm *sanStorageManager) CreateSnapshotGroupVClone(snapshotGroupName string) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_GROUP_VCLONE_BEGIN), snapshotGroupName, psm.storageSetting.Serial)

	objStorage := model.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return "", err
	}

	jobID, err := gatewayObj.CreateSnapshotGroupVClone(snapshotGroupName)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_GROUP_VCLONE_FAILED), snapshotGroupName, psm.storageSetting.Serial)
		return "", err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_GROUP_VCLONE_END), snapshotGroupName, psm.storageSetting.Serial)
	return jobID, nil
}

func (psm *sanStorageManager) ConvertSnapshotGroupVClone(snapshotGroupName string) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_CONVERT_GROUP_VCLONE_BEGIN), snapshotGroupName, psm.storageSetting.Serial)

	objStorage := model.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return "", err
	}

	jobID, err := gatewayObj.ConvertSnapshotGroupVClone(snapshotGroupName)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_CONVERT_GROUP_VCLONE_FAILED), snapshotGroupName, psm.storageSetting.Serial)
		return "", err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CONVERT_GROUP_VCLONE_END), snapshotGroupName, psm.storageSetting.Serial)
	return jobID, nil
}

func (psm *sanStorageManager) RestoreSnapshotGroupFromVClone(snapshotGroupName string) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_RESTORE_GROUP_VCLONE_BEGIN), snapshotGroupName, psm.storageSetting.Serial)

	objStorage := model.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return "", err
	}

	jobID, err := gatewayObj.RestoreSnapshotGroupFromVClone(snapshotGroupName)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_RESTORE_GROUP_VCLONE_FAILED), snapshotGroupName, psm.storageSetting.Serial)
		return "", err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_RESTORE_GROUP_VCLONE_END), snapshotGroupName, psm.storageSetting.Serial)
	return jobID, nil
}

//

func (psm *sanStorageManager) GetSnapshotFamily(ldevID int) (*model.SnapshotFamilyListResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_SNAPSHOT_FAMILY_BEGIN), ldevID, psm.storageSetting.Serial)

	objStorage := model.StorageDeviceSettings{
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

	resp, err := gatewayObj.GetSnapshotFamily(ldevID)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_SNAPSHOT_FAMILY_FAILED), ldevID, psm.storageSetting.Serial)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_SNAPSHOT_FAMILY_END), ldevID, psm.storageSetting.Serial)
	return resp, nil
}

func (psm *sanStorageManager) GetVirtualCloneParentVolumes() (*model.VirtualCloneParentVolumeList, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_VCLONE_PARENTS_BEGIN), psm.storageSetting.Serial)

	objStorage := model.StorageDeviceSettings{
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

	resp, err := gatewayObj.GetVirtualCloneParentVolumes()
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_VCLONE_PARENTS_FAILED), psm.storageSetting.Serial)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_VCLONE_PARENTS_END), psm.storageSetting.Serial)
	return resp, nil
}
