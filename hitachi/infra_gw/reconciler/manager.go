package infra_gw

import (
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

type InfraGwManager interface {
	GetStorageDevices() (*model.StorageDevices, error)
	GetStorageDevice(storageId string) (*model.StorageDevice, error)
	GetStoragePorts(storageId string) (*model.StoragePorts, error)
	GetParityGroups(storageId string) (*model.ParityGroups, error)
	GetHostGroups(storageId string, port string) (*model.HostGroups, error)
	GetHostGroup(storageId, port, hostGroupName string) (*model.HostGroup, error, bool)
	GetHostGroupById(storageId string, hostGrId string) (*model.HostGroup, error)
	GetStoragePools(storageId string) (*model.StoragePools, error)
	GetStoragePool(storageId, poolId string) (*model.StoragePool, error)
	GetIscsiTargets(storageId string, port string) (*model.IscsiTargets, error)
	GetIscsiTarget(storageId string, iscsiTargetId string) (*model.IscsiTarget, error)
	GetVolumes(storageId string) (*model.Volumes, error)
	GetUcpSystems() (*model.UcpSystems, error)

	ReconcileHostGroup(storageId string, createInput *model.CreateHostGroupParam) (*model.HostGroup, error)
	ReconcileStorageDevice(storageId string, createInput *model.CreateStorageDeviceParam) (*model.StorageDevice, error)
	// CreateHostGroup(storageId string, reqBody model.CreateHostGroupParam) (task *model.TaskResponse, err error)
}
