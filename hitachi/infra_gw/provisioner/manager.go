package infra_gw

import (
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

type InfraGwManager interface {

	// Storage Device Management
	GetStorageDevices() (*model.StorageDevices, error)
	GetStorageDevice(storageId string) (*model.StorageDevice, error)
	GetStoragePorts(storageId string) (*model.StoragePorts, error)
	AddStorageDevice(reqBody model.CreateStorageDeviceParam) (*string, error)
	UpdateStorageDevice(storageId string, reqBody model.PatchStorageDeviceParam) (*string, error)
	DeleteStorageDevice(storageId string) error

	//Parity Group Management
	GetParityGroups(storageId string) (*model.ParityGroups, error)

	// Host group Management
	GetHostGroups(storageId string, port string) (*model.HostGroups, error)
	GetHostGroup(storageId string, hostGrId string) (*model.HostGroup, error)
	CreateHostGroup(storageId string, reqBody model.CreateHostGroupParam) (*string, error)
	UpdateHostGroup(storageId, hostGroupId string, reqBody model.CreateHostGroupParam) (*string, error)

	//Storage Pool Management
	GetStoragePools(storageId string) (*model.StoragePools, error)
	GetStoragePool(storageId, poolId string) (*model.StoragePool, error)

	//Iscsi Target Management
	GetIscsiTargets(storageId string, port string) (*model.IscsiTargets, error)
	GetIscsiTarget(storageId string, iscsiTargetId string) (*model.IscsiTarget, error)
	CreateIscsiTarget(storageId string, reqBody model.CreateIscsiTargetParam) (*string, error)
	UpdateIscsiTarget(storageId, hostGroupId string, reqBody model.CreateIscsiTargetParam) (*string, error)

	//Volume Management
	GetVolumes(storageId string) (*model.Volumes, error)
	GetVolumeByID(storageId string, volumeId string) (*model.VolumeInfo, error)
	CreateVolume(storageId string, redBody *model.CreateVolumeParams) (*string, error)
	UpdateVolume(storageId string, volumeId string, redBody *model.UpdateVolumeParams) (*string, error)
	DeleteVolume(storageId string, volumeId string) error
	GetVolumesFromLdevIds(id string, fromLdevId int, toLdevId int) (*model.Volumes, error)
	GetVolumesByPartnerSubscriberID(id string, fromLdevId *int, toLdevId *int) (*model.MTVolumes, error)

	//UCP System Management
	GetUcpSystems() (*model.UcpSystems, error)
	GetUcpSystemById(id string) (*model.UcpSystem, error)
	CreateUcpSystem(reqBody model.CreateUcpSystemParam) (*string, error)

	//MT management
	GetPartnerAndSubscriberId(username string) (bool, *string, *string, error)
}
