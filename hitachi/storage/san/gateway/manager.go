package sanstorage

import (
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
)

// SanStorageManager interface
type SanStorageManager interface {
	// STORAGE
	GetStorageSystemInfo() (*sanmodel.StorageSystemInfo, error)
	GetStorageCapacity() (*sanmodel.StorageCapacity, error)
	// VOLUME
	GetLun(ldevID int) (*sanmodel.LogicalUnit, error)
	GetAllLun() (*sanmodel.LogicalUnits, error)
	CreateLun(reqBody sanmodel.CreateLunRequestGwy) (*int, error)
	UpdateLun(reqBody sanmodel.UpdateLunRequestGwy, ldevID int) (*int, error)
	ExpandLun(reqBody sanmodel.ExpandLunRequestGwy, ldevID int) (*int, error)
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
	GetStoragePorts() (*[]sanmodel.StoragePort, error)
	// DYNAMIC POOL
	GetDynamicPools() (*[]sanmodel.DynamicPool, error)
	GetDynamicPoolById(poolId int) (*sanmodel.DynamicPool, error)
	// PARITY GROUP
	GetParityGroups() (*[]sanmodel.ParityGroup, error)
	GetParityGroup(pgid string) (*sanmodel.ParityGroup, error)

	GetPools() (*[]sanmodel.Pool, error)
}
