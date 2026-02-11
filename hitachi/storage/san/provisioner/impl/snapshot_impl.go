package sanstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/storage/san/gateway/impl"

	model "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/san/provisioner/message-catalog"
)

func (psm *sanStorageManager) GetSnapshots(params model.GetSnapshotsParams) (*model.SnapshotListResponse, error) {
    log := commonlog.GetLogger()
    log.WriteEnter()
    defer log.WriteExit()

    log.WriteInfo(mc.GetMessage(mc.INFO_GET_SNAPSHOTS_BEGIN), psm.storageSetting.Serial)

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

    snapshotsList, err := gatewayObj.GetSnapshots(params)
    if err != nil {
        log.WriteError(mc.GetMessage(mc.ERR_GET_SNAPSHOTS_FAILED), psm.storageSetting.Serial)
        return nil, err
    }

    log.WriteInfo(mc.GetMessage(mc.INFO_GET_SNAPSHOTS_END), psm.storageSetting.Serial)
    return snapshotsList, nil
}

func (psm *sanStorageManager) GetSnapshot(pvolLdevID int, muNumber int) (*model.Snapshot, error) {
    log := commonlog.GetLogger()
    log.WriteEnter()
    defer log.WriteExit()

    log.WriteInfo(mc.GetMessage(mc.INFO_GET_SNAPSHOT_BEGIN), pvolLdevID, muNumber, psm.storageSetting.Serial)

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

    snapshot, err := gatewayObj.GetSnapshot(pvolLdevID, muNumber)
    if err != nil {
        log.WriteError(mc.GetMessage(mc.ERR_GET_SNAPSHOT_FAILED), pvolLdevID, muNumber, psm.storageSetting.Serial)
        return nil, err
    }

    log.WriteInfo(mc.GetMessage(mc.INFO_GET_SNAPSHOT_END), pvolLdevID, muNumber, psm.storageSetting.Serial)
    return snapshot, nil
}

func (psm *sanStorageManager) GetSnapshotReplicationsRange(params model.GetSnapshotReplicationsRangeParams) (*model.SnapshotAllListResponse, error) {
    log := commonlog.GetLogger()
    log.WriteEnter()
    defer log.WriteExit()

    log.WriteInfo(mc.GetMessage(mc.INFO_GET_SNAPSHOT_REPLICATIONS_BEGIN), psm.storageSetting.Serial)

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

    replications, err := gatewayObj.GetSnapshotReplicationsRange(params)
    if err != nil {
        log.WriteError(mc.GetMessage(mc.ERR_GET_SNAPSHOT_REPLICATIONS_FAILED), psm.storageSetting.Serial)
        return nil, err
    }

    log.WriteInfo(mc.GetMessage(mc.INFO_GET_SNAPSHOT_REPLICATIONS_END), psm.storageSetting.Serial)
    return replications, nil
}

func (psm *sanStorageManager) CreateSnapshot(request model.CreateSnapshotParams) (string, error) {
    log := commonlog.GetLogger()
    log.WriteEnter()
    defer log.WriteExit()

    log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_SNAPSHOT_BEGIN), psm.storageSetting.Serial)

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

    resIDs, err := gatewayObj.CreateSnapshot(request)
    if err != nil {
        log.WriteError(mc.GetMessage(mc.ERR_CREATE_SNAPSHOT_FAILED), psm.storageSetting.Serial)
        return "", err
    }

    log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_SNAPSHOT_END), psm.storageSetting.Serial)
    return resIDs, nil
}

func (psm *sanStorageManager) SplitSnapshot(pvolLdevID int, muNumber int, request model.SplitSnapshotRequest) (string, error) {
    log := commonlog.GetLogger()
    log.WriteEnter()
    defer log.WriteExit()

    log.WriteInfo(mc.GetMessage(mc.INFO_SPLIT_SNAPSHOT_BEGIN), pvolLdevID, muNumber, psm.storageSetting.Serial)

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

    resIDs, err := gatewayObj.SplitSnapshot(pvolLdevID, muNumber, request)
    if err != nil {
        log.WriteError(mc.GetMessage(mc.ERR_SPLIT_SNAPSHOT_FAILED), pvolLdevID, muNumber, psm.storageSetting.Serial)
        return "", err
    }

    log.WriteInfo(mc.GetMessage(mc.INFO_SPLIT_SNAPSHOT_END), pvolLdevID, muNumber, psm.storageSetting.Serial)
    return resIDs, nil
}

func (psm *sanStorageManager) ResyncSnapshot(pvolLdevID int, muNumber int, request model.ResyncSnapshotRequest) (string, error) {
    log := commonlog.GetLogger()
    log.WriteEnter()
    defer log.WriteExit()

    log.WriteInfo(mc.GetMessage(mc.INFO_RESYNC_SNAPSHOT_BEGIN), pvolLdevID, muNumber, psm.storageSetting.Serial)

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

    resIDs, err := gatewayObj.ResyncSnapshot(pvolLdevID, muNumber, request)
    if err != nil {
        log.WriteError(mc.GetMessage(mc.ERR_RESYNC_SNAPSHOT_FAILED), pvolLdevID, muNumber, psm.storageSetting.Serial)
        return "", err
    }

    log.WriteInfo(mc.GetMessage(mc.INFO_RESYNC_SNAPSHOT_END), pvolLdevID, muNumber, psm.storageSetting.Serial)
    return resIDs, nil
}

func (psm *sanStorageManager) RestoreSnapshot(pvolLdevID int, muNumber int, request model.RestoreSnapshotRequest) (string, error) {
    log := commonlog.GetLogger()
    log.WriteEnter()
    defer log.WriteExit()

    log.WriteInfo(mc.GetMessage(mc.INFO_RESTORE_SNAPSHOT_BEGIN), pvolLdevID, muNumber, psm.storageSetting.Serial)

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

    resIDs, err := gatewayObj.RestoreSnapshot(pvolLdevID, muNumber, request)
    if err != nil {
        log.WriteError(mc.GetMessage(mc.ERR_RESTORE_SNAPSHOT_FAILED), pvolLdevID, muNumber, psm.storageSetting.Serial)
        return "", err
    }

    log.WriteInfo(mc.GetMessage(mc.INFO_RESTORE_SNAPSHOT_END), pvolLdevID, muNumber, psm.storageSetting.Serial)
    return resIDs, nil
}

func (psm *sanStorageManager) DeleteSnapshot(pvolLdevID int, muNumber int) (string, error) {
    log := commonlog.GetLogger()
    log.WriteEnter()
    defer log.WriteExit()

    log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_SNAPSHOT_BEGIN), pvolLdevID, muNumber, psm.storageSetting.Serial)

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

    resIDs, err := gatewayObj.DeleteSnapshot(pvolLdevID, muNumber)
    if err != nil {
        log.WriteError(mc.GetMessage(mc.ERR_DELETE_SNAPSHOT_FAILED), pvolLdevID, muNumber, psm.storageSetting.Serial)
        return "", err
    }

    log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_SNAPSHOT_END), pvolLdevID, muNumber, psm.storageSetting.Serial)
    return resIDs, nil
}

func (psm *sanStorageManager) CloneSnapshot(pvolLdevID int, muNumber int, request model.CloneSnapshotRequest) (string, error) {
    log := commonlog.GetLogger()
    log.WriteEnter()
    defer log.WriteExit()

    log.WriteInfo(mc.GetMessage(mc.INFO_CLONE_SNAPSHOT_BEGIN), pvolLdevID, muNumber, psm.storageSetting.Serial)

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

    resIDs, err := gatewayObj.CloneSnapshot(pvolLdevID, muNumber, request)
    if err != nil {
        log.WriteError(mc.GetMessage(mc.ERR_CLONE_SNAPSHOT_FAILED), pvolLdevID, muNumber, psm.storageSetting.Serial)
        return "", err
    }

    log.WriteInfo(mc.GetMessage(mc.INFO_CLONE_SNAPSHOT_END), pvolLdevID, muNumber, psm.storageSetting.Serial)
    return resIDs, nil
}

func (psm *sanStorageManager) AssignSnapshotVolume(pvolLdevID int, muNumber int, request model.AssignSnapshotVolumeRequest) (string, error) {
    log := commonlog.GetLogger()
    log.WriteEnter()
    defer log.WriteExit()

    log.WriteInfo(mc.GetMessage(mc.INFO_ASSIGN_SNAPSHOT_VOLUME_BEGIN), pvolLdevID, muNumber, psm.storageSetting.Serial)

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

    resIDs, err := gatewayObj.AssignSnapshotVolume(pvolLdevID, muNumber, request)
    if err != nil {
        log.WriteError(mc.GetMessage(mc.ERR_ASSIGN_SNAPSHOT_VOLUME_FAILED), pvolLdevID, muNumber, psm.storageSetting.Serial)
        return "", err
    }

    log.WriteInfo(mc.GetMessage(mc.INFO_ASSIGN_SNAPSHOT_VOLUME_END), pvolLdevID, muNumber, psm.storageSetting.Serial)
    return resIDs, nil
}

func (psm *sanStorageManager) UnassignSnapshotVolume(pvolLdevID int, muNumber int) (string, error) {
    log := commonlog.GetLogger()
    log.WriteEnter()
    defer log.WriteExit()

    log.WriteInfo(mc.GetMessage(mc.INFO_UNASSIGN_SNAPSHOT_VOLUME_BEGIN), pvolLdevID, muNumber, psm.storageSetting.Serial)

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

    resIDs, err := gatewayObj.UnassignSnapshotVolume(pvolLdevID, muNumber)
    if err != nil {
        log.WriteError(mc.GetMessage(mc.ERR_UNASSIGN_SNAPSHOT_VOLUME_FAILED), pvolLdevID, muNumber, psm.storageSetting.Serial)
        return "", err
    }

    log.WriteInfo(mc.GetMessage(mc.INFO_UNASSIGN_SNAPSHOT_VOLUME_END), pvolLdevID, muNumber, psm.storageSetting.Serial)
    return resIDs, nil
}

func (psm *sanStorageManager) SetSnapshotRetentionPeriod(pvolLdevID int, muNumber int, request model.SetSnapshotRetentionPeriodRequest) (string, error) {
    log := commonlog.GetLogger()
    log.WriteEnter()
    defer log.WriteExit()

    log.WriteInfo(mc.GetMessage(mc.INFO_SNAPSHOT_RETENTION_BEGIN), pvolLdevID, muNumber, psm.storageSetting.Serial)

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

    resIDs, err := gatewayObj.SetSnapshotRetentionPeriod(pvolLdevID, muNumber, request)
    if err != nil {
        log.WriteError(mc.GetMessage(mc.ERR_SNAPSHOT_RETENTION_FAILED), pvolLdevID, muNumber, psm.storageSetting.Serial)
        return "", err
    }

    log.WriteInfo(mc.GetMessage(mc.INFO_SNAPSHOT_RETENTION_END), pvolLdevID, muNumber, psm.storageSetting.Serial)
    return resIDs, nil
}
