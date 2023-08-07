package sanstorage

import reconcilermodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"

// SanStorageManager interface
type SanStorageManager interface {
	// STORAGE
	GetStorageSystemInfo() (*reconcilermodel.StorageSystem, error)
	GetStorageSystem() (*reconcilermodel.StorageSystem, error)
	// VOLUME
	GetLun(ldevID int) (*reconcilermodel.LogicalUnit, error)
	GetRangeOfLuns(startLdevID int, endLdevID int, IsUndefinedLdev bool) (*[]reconcilermodel.LogicalUnit, error)
	SetLun(lunRequest *reconcilermodel.LunRequest) (*reconcilermodel.LogicalUnit, error)
	DeleteLun(ldevId int) error
	UpdateLun(lunUpdateRequest *reconcilermodel.UpdateLunRequest) (*reconcilermodel.LogicalUnit, error)
	// HOSTGROUP
	GetHostGroup(portID string, hostGroupNumber int) (*reconcilermodel.HostGroup, error)
	GetAllHostGroups() (*reconcilermodel.HostGroups, error)
	ReconcileHostGroup(createInput *reconcilermodel.CreateHostGroupRequest) (*reconcilermodel.HostGroup, error)
	DeleteHostGroup(portID string, hostGroupNumber int) error
	GetHostGroupsByPortIds(portIds []string) (*reconcilermodel.HostGroups, error)
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
	GetStoragePorts() (*[]reconcilermodel.StoragePort, error)
	GetStoragePortByPortId(portId string) (*reconcilermodel.StoragePort, error)
	// DYNAMIC POOL
	GetDynamicPools() (*[]reconcilermodel.DynamicPool, error)
	GetDynamicPoolById(poolId int) (*reconcilermodel.DynamicPool, error)
	//PARITY GROUP
	GetParityGroups(parityGroupIds ...[]string) (*[]reconcilermodel.ParityGroup, error)
}
