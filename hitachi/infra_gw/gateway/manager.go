package vssbstorage

import (
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

type InfraGwManager interface {
	GetStorageDevices() (*model.StorageDevices, error)
	GetStorageDevice(storageId string) (*model.StorageDevice, error)
	GetStoragePorts(storageId string) (*model.StoragePorts, error)
	GetParityGroups(storageId string) (*model.ParityGroups, error)
	GetHostGroups(storageId string, port string) (*model.HostGroups, error)
	GetHostGroup(storageId, hostGroupId string) (*model.HostGroup, error)
	GetStoragePools(storageId string) (*model.StoragePools, error)
	GetStoragePool(storageId, poolId string) (*model.StoragePool, error)
	GetIscsiTargets(storageId string, port string) (*model.IscsiTargets, error)
	GetIscsiTarget(storageId string, iscsiTargetId string) (*model.IscsiTarget, error)
	GetVolumes(storageId string) (*model.Volumes, error)
	GetUcpSystems() (*model.UcpSystems, error)

	CreateHostGroup(storageId string, reqBody model.CreateHostGroupParam) (*string, error)
	UpdateHostGroup(storageId, hostGroupId string, reqBody model.PatcheHostGroupParam) (*string, error)

	CreateUcpSystem(reqBody model.CreateUcpSystemParam) (*string, error)
	AddStorageDevice(storageId string, reqBody model.CreateStorageDeviceParam) (*string, error)
}
