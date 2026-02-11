package sanstorage

import (
	sangatewaymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"
)

// SanStorageManager interface
type SanStorageManager interface {
	// STORAGE
	GetStorageSystemInfo(detailInfoType ...string) (*sanmodel.StorageSystem, error)
	GetStorageSystem(detailInfoType ...string) (*sanmodel.StorageSystem, error)
	// VOLUME
	GetLun(ldevID int) (*sangatewaymodel.LogicalUnit, error)
	GetRangeOfLuns(startLdevID int, endLdevID int, IsUndefinedLdev bool, filterOption string, detailInfoType string) (*[]sangatewaymodel.LogicalUnit, error)
	CreateLun(reqBody sangatewaymodel.CreateLunRequestGwy) (*int, error)
	ExpandLun(ldevId int, newSize string) (*int, error)
	FormatLdev(ldevId int, req sangatewaymodel.FormatLdevRequestGwy) (*int, error)
	// BlockLun requests the gateway to change LDEV status to blocked
	BlockLun(ldevId int) error
	// UnblockLun requests the gateway to change LDEV status back to normal
	UnblockLun(ldevId int) error
	// StopAllVolumeFormat invokes the appliance-wide stop-format action.
	StopAllVolumeFormat() error
	DeleteLun(ldevId int) error
	UpdateLun(ldevId int, updReq sangatewaymodel.UpdateLunRequestGwy) (*int, error)
	SetEseVolume(ldevId int, isEseVolume bool) error
	// HOSTGROUP
	GetHostGroup(portID string, hostGroupNumber int) (*sanmodel.HostGroup, error)
	GetAllHostGroups() (*sanmodel.HostGroups, error)
	CreateHostGroup(hgBody sanmodel.CreateHostGroupRequest) (*sanmodel.HostGroup, error)
	DeleteHostGroup(portID string, hostGroupNumber int) error
	SetHostGroupModeAndOptions(portID string, hostGroupNumber int, reqBody sanmodel.SetHostModeAndOptions) error
	GetSupportedHostModes() (*sanmodel.HostModeAndOptions, error)
	GetPavAliases(cuNumber *int) (*[]sangatewaymodel.PavAlias, error)
	AssignPavAlias(baseLdevID int, aliasLdevIDs []int) error
	UnassignPavAlias(aliasLdevIDs []int) error
	AddWwnToHG(reqBody sanmodel.AddWwnToHg) (err error)
	SetHostWwnNickName(portID string, hostGroupNumber int, hostWwn string, wwnNickname string) error
	DeleteWwn(portID string, hostGroupNumber int, wwn string) (err error)
	AddLdevToHG(reqBody sanmodel.AddLdevToHg) (err error)
	RemoveLdevFromHG(portID string, hostGroupNumber int, lunID int) (err error)
	// ISCSI TARGET
	GetIscsiTarget(portID string, iscsiTargetNumber int) (*sanmodel.IscsiTarget, error)
	GetIscsiTargetsByPortId(portID string) (*sanmodel.IscsiTargets, error)
	GetAllIscsiTargets() (*sanmodel.IscsiTargets, error)
	SetIscsiHostGroupModeAndOptions(portID string, hostGroupNumber int, reqBody sanmodel.SetIscsiHostModeAndOptions) error
	CreateIscsiTarget(reqBody sanmodel.CreateIscsiTargetReq) (*sanmodel.IscsiTarget, error)
	AddLdevToIscsiTarget(reqBody sanmodel.CreateIscsiTargetReq) error
	AddInitiatorsToIscsiTarget(reqBody sanmodel.CreateIscsiTargetReq) error
	SetIscsiNameForIscsiTarget(reqBody sanmodel.SetIscsiNameReq) error
	SetNicknameForIscsiName(portID string, iscsiTargetNumber int, iscsiName string, reqBody sanmodel.SetNicknameIscsiReq) error
	DeleteIscsiNameFromIscsiTarget(portID string, iscsiTargetNumber int, iscsiName string) error
	DeleteIscsiTarget(portID string, iscsiTargetNumber int) error
	// ISCSI TARGET CHAP USER
	GetChapUser(portID string, iscsiTargetNumber int, chapUserName string, wayOfChapUser string) (*sanmodel.IscsiTargetChapUser, error)
	GetChapUsers(portID string, iscsiTargetNumber int) (*sanmodel.IscsiTargetChapUsers, error)
	CreateChapUser(portID string, iscsiTargetNumber int, wayOfChapUser, chapUserName string, chapUserSecret string) error
	DeleteChapUser(portID string, iscsiTargetNumber int, wayOfChapUser, chapUserName string) error
	ChangeChapUserSecret(portID string, iscsiTargetNumber int, wayOfChapUser, chapUserName string, chapUserSecret string) error
	// STORAGE PORTS
	GetStoragePorts(detailInfoTypes []string, portType, portAttributes string) (*[]sanmodel.StoragePort, error)
	GetStoragePortByPortId(portId string) (*sanmodel.StoragePort, error)
	// DYNAMIC POOL
	GetDynamicPools(isMainframe *bool, poolType string, detailInfoType ...string) (*[]sanmodel.DynamicPool, error)
	GetDynamicPoolById(poolId int) (*sanmodel.DynamicPool, error)
	// PARITY GROUP
	GetParityGroups(detailInfoType string, driveTypeName string, clprId *int, parityGroupIds ...[]string) (*[]sanmodel.ParityGroup, error)
	GetParityGroup(parityGroupId string) (*sanmodel.ParityGroup, error)

	GetPools() (*[]sanmodel.Pool, error)

	// SNAPSHOT
	GetSnapshots(params sangatewaymodel.GetSnapshotsParams) (*sangatewaymodel.SnapshotListResponse, error)
	GetSnapshot(pvolLdevID int, muNumber int) (*sangatewaymodel.Snapshot, error)
	GetSnapshotReplicationsRange(params sangatewaymodel.GetSnapshotReplicationsRangeParams) (*sangatewaymodel.SnapshotAllListResponse, error)
	CreateSnapshot(request sangatewaymodel.CreateSnapshotParams) (string, error)
	SplitSnapshot(pvolLdevID int, muNumber int, request sangatewaymodel.SplitSnapshotRequest) (string, error)
	ResyncSnapshot(pvolLdevID int, muNumber int, request sangatewaymodel.ResyncSnapshotRequest) (string, error)
	RestoreSnapshot(pvolLdevID int, muNumber int, request sangatewaymodel.RestoreSnapshotRequest) (string, error)
	DeleteSnapshot(pvolLdevID int, muNumber int) (string, error)
	CloneSnapshot(pvolLdevID int, muNumber int, request sangatewaymodel.CloneSnapshotRequest) (string, error)
	AssignSnapshotVolume(pvolLdevID int, muNumber int, request sangatewaymodel.AssignSnapshotVolumeRequest) (string, error)
	UnassignSnapshotVolume(pvolLdevID int, muNumber int) (string, error)
	SetSnapshotRetentionPeriod(pvolLdevID int, muNumber int, request sangatewaymodel.SetSnapshotRetentionPeriodRequest) (string, error)

	// SNAPSHOT GROUP
	GetSnapshotGroups(params sangatewaymodel.GetSnapshotGroupsParams) (*sangatewaymodel.SnapshotGroupListResponse, error)
	GetSnapshotGroup(snapshotGroupID string, params sangatewaymodel.GetSnapshotGroupsParams) (*sangatewaymodel.SnapshotGroup, error)
	SplitSnapshotGroup(snapshotGroupID string, request sangatewaymodel.SplitSnapshotRequest) (string, error)
	ResyncSnapshotGroup(snapshotGroupID string, request sangatewaymodel.ResyncSnapshotRequest) (string, error)
	RestoreSnapshotGroup(snapshotGroupID string, request sangatewaymodel.RestoreSnapshotRequest) (string, error)
	DeleteSnapshotGroup(snapshotGroupID string) (string, error)
	CloneSnapshotGroup(snapshotGroupID string, request sangatewaymodel.CloneSnapshotRequest) (string, error)
	SetSnapshotGroupRetentionPeriod(snapshotGroupID string, request sangatewaymodel.SetSnapshotRetentionPeriodRequest) (string, error)

	// VIRTUAL CLONE
	CreateSnapshotVClone(pvolLdevID int, muNumber int) (string, error)
	ConvertSnapshotVClone(pvolLdevID int, muNumber int) (string, error)
	RestoreSnapshotFromVClone(pvolLdevID int, muNumber int) (string, error)
	CreateSnapshotGroupVClone(snapshotGroupName string) (string, error)
	ConvertSnapshotGroupVClone(snapshotGroupName string) (string, error)
	RestoreSnapshotGroupFromVClone(snapshotGroupName string) (string, error)
	GetSnapshotFamily(ldevID int) (*sangatewaymodel.SnapshotFamilyListResponse, error)
	GetVirtualCloneParentVolumes() (*sangatewaymodel.VirtualCloneParentVolumeList, error)

	// SNAPSHOT TREE
	DeleteSnapshotTree(request sangatewaymodel.DeleteSnapshotTreeRequest) (string, error)
	DeleteGarbageData(request sangatewaymodel.DeleteGarbageDataRequest) (string, error)
}
