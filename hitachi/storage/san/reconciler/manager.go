package sanstorage

import (
	gatewaymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
)

// SanStorageManager interface
type SanStorageManager interface {
	// STORAGE
	GetStorageSystemInfo(detailInfoType ...string) (*reconcilermodel.StorageSystem, error)
	GetStorageSystem(detailInfoType ...string) (*reconcilermodel.StorageSystem, error)
	// VOLUME
	GetLun(ldevID int) (*gatewaymodel.LogicalUnit, error)
	GetRangeOfLuns(startLdevID int, endLdevID int, IsUndefinedLdev bool, filterOption string, detailInfoType string) (*[]gatewaymodel.LogicalUnit, error)
	SetLun(lunRequest *reconcilermodel.LunRequest) (*int, error)
	DeleteLun(ldevId int) error
	UpdateLun(lunUpdateRequest *reconcilermodel.UpdateLunRequest) (*int, error)
	SetEseVolume(ldevId int, isEseVolume bool) error
	// HOSTGROUP
	GetHostGroup(portID string, hostGroupNumber int) (*reconcilermodel.HostGroup, error)
	GetAllHostGroups() (*reconcilermodel.HostGroups, error)
	FormatLdev(ldevID int, req reconcilermodel.FormatLdevRequest) (*int, error)
	// StopAllVolumeFormat invokes the appliance-wide stop-format action via the provisioner.
	StopAllVolumeFormat() error
	// BlockLun requests the provisioner to change LDEV status to blocked.
	BlockLun(ldevID int) error
	// UnblockLun requests the provisioner to change LDEV status back to normal.
	UnblockLun(ldevID int) error
	ReconcileHostGroup(createInput *reconcilermodel.CreateHostGroupRequest) (*reconcilermodel.HostGroup, error)
	DeleteHostGroup(portID string, hostGroupNumber int) error
	GetHostGroupsByPortIds(portIds []string) (*reconcilermodel.HostGroups, error)
	GetSupportedHostModes() (*reconcilermodel.HostModeAndOptions, error)
	GetPavAliases(cuNumber *int) (*[]gatewaymodel.PavAlias, error)
	AssignPavAlias(baseLdevID int, aliasLdevIDs []int) error
	UnassignPavAlias(aliasLdevIDs []int) error
	// ISCSITARGET
	GetIscsiTarget(portID string, iscsiTargetNumber int) (*reconcilermodel.IscsiTarget, error)
	GetAllIscsiTargets() (*reconcilermodel.IscsiTargets, error)
	ReconcileIscsiTarget(createInput *reconcilermodel.CreateIscsiTargetReq) (*reconcilermodel.IscsiTarget, error)
	DeleteIscsiTarget(portID string, iscsiTargetNumber int) error
	GetIscsiTargetsByPortIds(portIds []string) (*reconcilermodel.IscsiTargets, error)
	// ISCSI TARGET CHAP USER
	GetChapUsers(portID string, iscsiTargetNumber int) (*reconcilermodel.IscsiTargetChapUsers, error)
	GetChapUser(portID string, iscsiTargetNumber int, chapUserName, wayOfChapUser string) (*reconcilermodel.IscsiTargetChapUser, error)
	ReconcileChapUser(createInput *reconcilermodel.ChapUserRequest) (*reconcilermodel.IscsiTargetChapUser, error)
	DeleteChapUser(portID string, iscsiTargetNumber int, chapUserName, wayOfChapUser string) error
	// STORAGE PORTS
	GetStoragePorts(detailInfoTypes []string, portType, portAttributes string) (*[]reconcilermodel.StoragePort, error)
	GetStoragePortByPortId(portId string) (*reconcilermodel.StoragePort, error)
	// DYNAMIC POOL
	GetDynamicPools(isMainframe *bool, poolType string, detailInfoType ...string) (*[]reconcilermodel.DynamicPool, error)
	GetDynamicPoolById(poolId int) (*reconcilermodel.DynamicPool, error)
	//PARITY GROUP
	GetParityGroups(detailInfoType string, driveTypeName string, clprId *int, parityGroupIds ...[]string) (*[]reconcilermodel.ParityGroup, error)
	GetParityGroup(parityGroupId string) (*reconcilermodel.ParityGroup, error)

	GetPools() (*[]reconcilermodel.Pool, error)

	// SNAPSHOT
	ReconcileGetSnapshot(pvolID *int, mu *int) (*gatewaymodel.Snapshot, error)
	ReconcileGetMultipleSnapshots(input reconcilermodel.SnapshotGetMultipleInput) (*gatewaymodel.SnapshotListResponse, error)
	ReconcileGetMultipleSnapshotsRange(input reconcilermodel.SnapshotGetMultipleRangeInput) (*gatewaymodel.SnapshotListResponse, error)
	ReconcileSnapshotVclone(input reconcilermodel.SnapshotReconcilerInput) (*reconcilermodel.ReconcileSnapshotResult, error)
	ReconcileReadExistingSnapshotVclone(input reconcilermodel.SnapshotReconcilerInput) (*reconcilermodel.ReconcileSnapshotResult, error)
	ReconcileGetFamily(ldevID int) ([]gatewaymodel.SnapshotFamily, error)
	ReconcileGetVirtualCloneParentVolumes() ([]int, error)

	ReconcileGetSnapshotGroup(snapshotGroupID string) (*gatewaymodel.SnapshotGroup, error)
	ReconcileGetMultipleSnapshotGroups(includePairs bool) (*gatewaymodel.SnapshotGroupListResponse, error)
	ReconcileSnapshotGroup(input reconcilermodel.SnapshotGroupReconcilerInput) (*gatewaymodel.SnapshotGroup, error)
	ReconcileSnapshotGroupVFamily(input reconcilermodel.SnapshotGroupReconcilerInput) (*gatewaymodel.SnapshotGroup, []gatewaymodel.SnapshotFamily, error)
}
