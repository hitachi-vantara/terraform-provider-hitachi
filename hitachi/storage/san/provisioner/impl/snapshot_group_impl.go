package sanstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/storage/san/gateway/impl"
	model "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/san/provisioner/message-catalog"
)

func (psm *sanStorageManager) GetSnapshotGroups(params model.GetSnapshotGroupsParams) (*model.SnapshotGroupListResponse, error) {
    log := commonlog.GetLogger()
    log.WriteEnter()
    defer log.WriteExit()

    log.WriteInfo(mc.GetMessage(mc.INFO_GET_SNAPSHOT_GROUPS_BEGIN), psm.storageSetting.Serial)

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

    resp, err := gatewayObj.GetSnapshotGroups(params)
    if err != nil {
        log.WriteError(mc.GetMessage(mc.ERR_GET_SNAPSHOT_GROUPS_FAILED), psm.storageSetting.Serial)
        return nil, err
    }

    log.WriteInfo(mc.GetMessage(mc.INFO_GET_SNAPSHOT_GROUPS_END), psm.storageSetting.Serial)
    return resp, nil
}

func (psm *sanStorageManager) GetSnapshotGroup(snapshotGroupID string, params model.GetSnapshotGroupsParams) (*model.SnapshotGroup, error) {
    log := commonlog.GetLogger()
    log.WriteEnter()
    defer log.WriteExit()

    log.WriteInfo(mc.GetMessage(mc.INFO_GET_SNAPSHOT_GROUP_BEGIN), snapshotGroupID, psm.storageSetting.Serial)

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

    resp, err := gatewayObj.GetSnapshotGroup(snapshotGroupID, params)
    if err != nil {
        log.WriteError(mc.GetMessage(mc.ERR_GET_SNAPSHOT_GROUP_FAILED), snapshotGroupID, psm.storageSetting.Serial)
        return nil, err
    }

    log.WriteInfo(mc.GetMessage(mc.INFO_GET_SNAPSHOT_GROUP_END), snapshotGroupID, psm.storageSetting.Serial)
    return resp, nil
}

func (psm *sanStorageManager) SplitSnapshotGroup(snapshotGroupID string, request model.SplitSnapshotRequest) (string, error) {
    log := commonlog.GetLogger()
    log.WriteEnter()
    defer log.WriteExit()

    log.WriteInfo(mc.GetMessage(mc.INFO_SPLIT_SNAPSHOT_GROUP_BEGIN), snapshotGroupID, psm.storageSetting.Serial)

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

    jobID, err := gatewayObj.SplitSnapshotGroup(snapshotGroupID, request)
    if err != nil {
        log.WriteError(mc.GetMessage(mc.ERR_SPLIT_SNAPSHOT_GROUP_FAILED), snapshotGroupID, psm.storageSetting.Serial)
        return "", err
    }

    log.WriteInfo(mc.GetMessage(mc.INFO_SPLIT_SNAPSHOT_GROUP_END), snapshotGroupID, psm.storageSetting.Serial)
    return jobID, nil
}

func (psm *sanStorageManager) ResyncSnapshotGroup(snapshotGroupID string, request model.ResyncSnapshotRequest) (string, error) {
    log := commonlog.GetLogger()
    log.WriteEnter()
    defer log.WriteExit()

    log.WriteInfo(mc.GetMessage(mc.INFO_RESYNC_SNAPSHOT_GROUP_BEGIN), snapshotGroupID, psm.storageSetting.Serial)

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

    jobID, err := gatewayObj.ResyncSnapshotGroup(snapshotGroupID, request)
    if err != nil {
        log.WriteError(mc.GetMessage(mc.ERR_RESYNC_SNAPSHOT_GROUP_FAILED), snapshotGroupID, psm.storageSetting.Serial)
        return "", err
    }

    log.WriteInfo(mc.GetMessage(mc.INFO_RESYNC_SNAPSHOT_GROUP_END), snapshotGroupID, psm.storageSetting.Serial)
    return jobID, nil
}

func (psm *sanStorageManager) RestoreSnapshotGroup(snapshotGroupID string, request model.RestoreSnapshotRequest) (string, error) {
    log := commonlog.GetLogger()
    log.WriteEnter()
    defer log.WriteExit()

    log.WriteInfo(mc.GetMessage(mc.INFO_RESTORE_SNAPSHOT_GROUP_BEGIN), snapshotGroupID, psm.storageSetting.Serial)

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

    jobID, err := gatewayObj.RestoreSnapshotGroup(snapshotGroupID, request)
    if err != nil {
        log.WriteError(mc.GetMessage(mc.ERR_RESTORE_SNAPSHOT_GROUP_FAILED), snapshotGroupID, psm.storageSetting.Serial)
        return "", err
    }

    log.WriteInfo(mc.GetMessage(mc.INFO_RESTORE_SNAPSHOT_GROUP_END), snapshotGroupID, psm.storageSetting.Serial)
    return jobID, nil
}

func (psm *sanStorageManager) DeleteSnapshotGroup(snapshotGroupID string) (string, error) {
    log := commonlog.GetLogger()
    log.WriteEnter()
    defer log.WriteExit()

    log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_SNAPSHOT_GROUP_BEGIN), snapshotGroupID, psm.storageSetting.Serial)

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

    jobID, err := gatewayObj.DeleteSnapshotGroup(snapshotGroupID)
    if err != nil {
        log.WriteError(mc.GetMessage(mc.ERR_DELETE_SNAPSHOT_GROUP_FAILED), snapshotGroupID, psm.storageSetting.Serial)
        return "", err
    }

    log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_SNAPSHOT_GROUP_END), snapshotGroupID, psm.storageSetting.Serial)
    return jobID, nil
}

func (psm *sanStorageManager) CloneSnapshotGroup(snapshotGroupID string, request model.CloneSnapshotRequest) (string, error) {
    log := commonlog.GetLogger()
    log.WriteEnter()
    defer log.WriteExit()

    log.WriteInfo(mc.GetMessage(mc.INFO_CLONE_SNAPSHOT_GROUP_BEGIN), snapshotGroupID, psm.storageSetting.Serial)

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

    jobID, err := gatewayObj.CloneSnapshotGroup(snapshotGroupID, request)
    if err != nil {
        log.WriteError(mc.GetMessage(mc.ERR_CLONE_SNAPSHOT_GROUP_FAILED), snapshotGroupID, psm.storageSetting.Serial)
        return "", err
    }

    log.WriteInfo(mc.GetMessage(mc.INFO_CLONE_SNAPSHOT_GROUP_END), snapshotGroupID, psm.storageSetting.Serial)
    return jobID, nil
}

func (psm *sanStorageManager) SetSnapshotGroupRetentionPeriod(snapshotGroupID string, request model.SetSnapshotRetentionPeriodRequest) (string, error) {
    log := commonlog.GetLogger()
    log.WriteEnter()
    defer log.WriteExit()

    log.WriteInfo(mc.GetMessage(mc.INFO_SNAPSHOT_GROUP_RETENTION_BEGIN), snapshotGroupID, psm.storageSetting.Serial)

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

    resIDs, err := gatewayObj.SetSnapshotGroupRetentionPeriod(snapshotGroupID, request)
    if err != nil {
        log.WriteError(mc.GetMessage(mc.ERR_SNAPSHOT_GROUP_RETENTION_FAILED), snapshotGroupID, psm.storageSetting.Serial)
        return "", err
    }

    log.WriteInfo(mc.GetMessage(mc.INFO_SNAPSHOT_GROUP_RETENTION_END), snapshotGroupID, psm.storageSetting.Serial)
    return resIDs, nil
}

func (psm *sanStorageManager) DeleteSnapshotTree(request model.DeleteSnapshotTreeRequest) (string, error) {
    log := commonlog.GetLogger()
    log.WriteEnter()
    defer log.WriteExit()

    log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_SNAPSHOT_TREE_BEGIN), request.Parameters.LdevID, psm.storageSetting.Serial)

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

    jobID, err := gatewayObj.DeleteSnapshotTree(request)
    if err != nil {
        log.WriteError(mc.GetMessage(mc.ERR_DELETE_SNAPSHOT_TREE_FAILED), request.Parameters.LdevID, psm.storageSetting.Serial)
        return "", err
    }

    log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_SNAPSHOT_TREE_END), request.Parameters.LdevID, psm.storageSetting.Serial)
    return jobID, nil
}

func (psm *sanStorageManager) DeleteGarbageData(request model.DeleteGarbageDataRequest) (string, error) {
    log := commonlog.GetLogger()
    log.WriteEnter()
    defer log.WriteExit()

    log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_GARBAGE_DATA_BEGIN), request.Parameters.LdevID, psm.storageSetting.Serial)

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

    jobID, err := gatewayObj.DeleteGarbageData(request)
    if err != nil {
        log.WriteError(mc.GetMessage(mc.ERR_DELETE_GARBAGE_DATA_FAILED), request.Parameters.LdevID, psm.storageSetting.Serial)
        return "", err
    }

    log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_GARBAGE_DATA_END), request.Parameters.LdevID, psm.storageSetting.Serial)
    return jobID, nil
}
