package sanstorage

import (
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
)

// SanStorageManager interface
type SanStorageManager interface {
	// STORAGE
	GetStorageSystemInfo(detailInfoType ...string) (*sanmodel.StorageSystemInfo, error)
	GetStorageCapacity() (*sanmodel.StorageCapacity, error)
	// VOLUME
	GetLun(ldevID int) (*sanmodel.LogicalUnit, error)
	GetRangeOfLunsWithOptions(startLdevID int, endLdevID int, isUndefinedLdev bool, filterOption string, detailInfoType string) (*[]sanmodel.LogicalUnit, error)
	GetAllLun() (*sanmodel.LogicalUnits, error)
	CreateLun(reqBody sanmodel.CreateLunRequestGwy) (*int, error)
	UpdateLun(reqBody sanmodel.UpdateLunRequestGwy, ldevID int) (*int, error)
	SetEseVolume(ldevID int, isEseVolume bool) error
	ExpandLun(reqBody sanmodel.ExpandLunRequestGwy, ldevID int) (*int, error)
	FormatLdev(reqBody sanmodel.FormatLdevRequestGwy, ldevID int) (*int, error)
	BlockLun(ldevID int) error
	UnblockLun(ldevID int) error
	StopAllVolumeFormat() error
	LockResources(req sanmodel.LockResourcesReq) error
	UnlockResources() error
	DeleteLun(ldevID int, capacitySaving bool) error
	// HOSTGROUP
	GetHostGroup(portID string, hostGroupNumber int) (*sanmodel.HostGroupGwy, error)
	GetHostGroupWwns(portID string, hostGroupNumber int) (*[]sanmodel.HostWwnDetail, error)
	GetHostGroupLuPaths(portID string, hostGroupNumber int) (*[]sanmodel.HostLuPath, error)
	RemoveLdevFromHG(portID string, hostGroupNumber int, lunID int) (err error)
	AddLdevToHG(reqBody sanmodel.AddLdevToHgReqGwy) (err error)
	AddWwnToHG(reqBody sanmodel.AddWwnToHgReqGwy) (err error)
	GetAllHostGroups() (*sanmodel.HostGroups, error)
	DeleteHostGroup(portID string, hostGroupNumber int) (err error)
	CreateHostGroup(reqBody sanmodel.CreateHostGroupReqGwy) (portid *string, hgnumber *int, err error)
	SetHostGroupModeAndOptions(portID string, hostGroupNumber int, reqBody sanmodel.SetHostModeAndOptions) error
	GetHostGroupModeAndOptions() (*sanmodel.HostModeAndOptions, error)
	GetSupportedHostModes() (*sanmodel.HostModeAndOptions, error)
	SetHostWwnNickName(portID string, hostGroupNumber int, hostWwn string, wwnNickname string) error
	DeleteWwn(portID string, hostGroupNumber int, wwn string) (err error)
	//ISCSI TARGET
	GetIscsiTarget(portID string, iscsiTargetNumber int) (*sanmodel.IscsiTargetGwy, error)
	GetIscsiTargetsByPortId(portID string) (*sanmodel.IscsiTargets, error)
	GetAllIscsiTargets() (*sanmodel.IscsiTargets, error)
	CreateIscsiTarget(reqBody sanmodel.CreateIscsiTargetReq) (portId *string, itNum *int, err error)
	GetIscsiNameInformation(portID string, iscsiTargetNumber int) (*[]sanmodel.IscsiNameInformation, error)
	GetIscsiTargetGroupLuPaths(portID string, iscsiTargetNumber int) (*[]sanmodel.IscsiTargetLuPath, error)
	SetIScsiTargetHostModeAndHostModeOptions(portID string, hostGroupNumber int, reqBody sanmodel.SetIscsiHostModeAndOptions) error
	SetIscsiNameForIscsiTarget(reqBody sanmodel.SetIscsiNameReq) error
	SetNicknameForIscsiName(portID string, iscsiTargetNumber int, iscsiName string, reqBody sanmodel.SetNicknameIscsiReq) error
	DeleteIscsiNameFromIscsiTarget(portID string, iscsiTargetNumber int, iscsiName string) error
	DeleteIscsiTarget(portID string, iscsiTargetNumber int) error
	//ISCSI TARGET CHAP USER
	GetChapUsers(portID string, iscsiTargetNumber int) (*sanmodel.IscsiTargetChapUsers, error)
	GetChapUser(portID string, iscsiTargetNumber int, chapUserName string, wayOfChapUser string) (*sanmodel.IscsiTargetChapUser, error)
	SetChapUserName(portID string, iscsiTargetNumber int, wayOfChapUser, chapUserName string) error
	SetChapUserSecret(portID string, iscsiTargetNumber int, wayOfChapUser, chapUserName, capUserPassword string) error
	DeleteChapUser(portID string, iscsiTargetNumber int, wayOfChapUser, chapUserName string) error
	// STORAGE PORTS
	GetStoragePorts(detailInfoTypes []string, portType, portAttributes string) (*[]sanmodel.StoragePort, error)
	// DYNAMIC POOL
	GetDynamicPools(isMainframe *bool, poolType string, detailInfoType ...string) (*[]sanmodel.DynamicPool, error)
	GetDynamicPoolById(poolId int) (*sanmodel.DynamicPool, error)
	// PARITY GROUP
	GetParityGroups(detailInfoType string, driveTypeName string, clprId *int) (*[]sanmodel.ParityGroup, error)
	GetParityGroup(pgid string) (*sanmodel.ParityGroup, error)

	// PAV ALIAS
	GetPavAliases(cuNumber *int) (*[]sanmodel.PavAlias, error)
	AssignPavAlias(baseLdevID int, aliasLdevIDs []int) error
	UnassignPavAlias(aliasLdevIDs []int) error

	GetPools() (*[]sanmodel.Pool, error)

	// SNAPSHOT
	GetSnapshots(params sanmodel.GetSnapshotsParams) (*sanmodel.SnapshotListResponse, error)
	GetSnapshot(pvolLdevID int, muNumber int) (*sanmodel.Snapshot, error)
	GetSnapshotReplicationsRange(params sanmodel.GetSnapshotReplicationsRangeParams) (*sanmodel.SnapshotAllListResponse, error)
	CreateSnapshot(request sanmodel.CreateSnapshotParams) (string, error)
	SplitSnapshot(pvolLdevID int, muNumber int, request sanmodel.SplitSnapshotRequest) (string, error)
	ResyncSnapshot(pvolLdevID int, muNumber int, request sanmodel.ResyncSnapshotRequest) (string, error)
	RestoreSnapshot(pvolLdevID int, muNumber int, request sanmodel.RestoreSnapshotRequest) (string, error)
	DeleteSnapshot(pvolLdevID int, muNumber int) (string, error)
	CloneSnapshot(pvolLdevID int, muNumber int, request sanmodel.CloneSnapshotRequest) (string, error)
	AssignSnapshotVolume(pvolLdevID int, muNumber int, request sanmodel.AssignSnapshotVolumeRequest) (string, error)
	UnassignSnapshotVolume(pvolLdevID int, muNumber int) (string, error)
	SetSnapshotRetentionPeriod(pvolLdevID int, muNumber int, request sanmodel.SetSnapshotRetentionPeriodRequest) (string, error)

	// SNAPSHOT GROUP
	GetSnapshotGroups(params sanmodel.GetSnapshotGroupsParams) (*sanmodel.SnapshotGroupListResponse, error)
	GetSnapshotGroup(snapshotGroupID string, params sanmodel.GetSnapshotGroupsParams) (*sanmodel.SnapshotGroup, error)
	SplitSnapshotGroup(snapshotGroupID string, request sanmodel.SplitSnapshotRequest) (string, error)
	ResyncSnapshotGroup(snapshotGroupID string, request sanmodel.ResyncSnapshotRequest) (string, error)
	RestoreSnapshotGroup(snapshotGroupID string, request sanmodel.RestoreSnapshotRequest) (string, error)
	DeleteSnapshotGroup(snapshotGroupID string) (string, error)
	CloneSnapshotGroup(snapshotGroupID string, request sanmodel.CloneSnapshotRequest) (string, error)
	SetSnapshotGroupRetentionPeriod(snapshotGroupID string, request sanmodel.SetSnapshotRetentionPeriodRequest) (string, error)

	// VIRTUAL CLONE
	CreateSnapshotVClone(pvolLdevID int, muNumber int) (string, error)
	ConvertSnapshotVClone(pvolLdevID int, muNumber int) (string, error)
	CreateSnapshotGroupVClone(snapshotGroupName string) (string, error)
	ConvertSnapshotGroupVClone(snapshotGroupName string) (string, error)
	RestoreSnapshotFromVClone(pvolLdevID int, muNumber int) (string, error)
	RestoreSnapshotGroupFromVClone(snapshotGroupName string) (string, error)
	GetSnapshotFamily(ldevID int) (*sanmodel.SnapshotFamilyListResponse, error)
	GetVirtualCloneParentVolumes() (*sanmodel.VirtualCloneParentVolumeList, error)

	// SNAPSHOT TREE
	DeleteSnapshotTree(request sanmodel.DeleteSnapshotTreeRequest) (string, error)
	DeleteGarbageData(request sanmodel.DeleteGarbageDataRequest) (string, error)
}
