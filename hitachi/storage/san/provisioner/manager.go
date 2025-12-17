package sanstorage

import (
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"
	sangatewaymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
)

// SanStorageManager interface
type SanStorageManager interface {
	// STORAGE
	GetStorageSystemInfo() (*sanmodel.StorageSystem, error)
	GetStorageSystem() (*sanmodel.StorageSystem, error)
	// VOLUME
	GetLun(ldevID int) (*sangatewaymodel.LogicalUnit, error)
	GetRangeOfLuns(startLdevID int, endLdevID int, IsUndefinedLdev bool) (*[]sangatewaymodel.LogicalUnit, error)
	CreateLun(reqBody sangatewaymodel.CreateLunRequestGwy) (*int, error)
	ExpandLun(ldevId int, newSize string) (*int, error)
	DeleteLun(ldevId int) error
	UpdateLun(ldevId int, updReq sangatewaymodel.UpdateLunRequestGwy) (*int, error)
	// HOSTGROUP
	GetHostGroup(portID string, hostGroupNumber int) (*sanmodel.HostGroup, error)
	GetAllHostGroups() (*sanmodel.HostGroups, error)
	CreateHostGroup(hgBody sanmodel.CreateHostGroupRequest) (*sanmodel.HostGroup, error)
	DeleteHostGroup(portID string, hostGroupNumber int) error
	SetHostGroupModeAndOptions(portID string, hostGroupNumber int, reqBody sanmodel.SetHostModeAndOptions) error
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
	GetStoragePorts() (*[]sanmodel.StoragePort, error)
	GetStoragePortByPortId(portId string) (*sanmodel.StoragePort, error)
	// DYNAMIC POOL
	GetDynamicPools() (*[]sanmodel.DynamicPool, error)
	GetDynamicPoolById(poolId int) (*sanmodel.DynamicPool, error)
	// PARITY GROUP
	GetParityGroups(parityGroupIds ...[]string) (*[]sanmodel.ParityGroup, error)
	GetParityGroup(parityGroupId string) (*sanmodel.ParityGroup, error)

	GetPools() (*[]sanmodel.Pool, error)
}
