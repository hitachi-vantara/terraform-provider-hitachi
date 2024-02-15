package infra_gw

import (
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

type InfraGwManager interface {
	// Storage Device Management
	GetStorageDevices() (*model.StorageDevices, error)
	GetStorageDevice(storageId string) (*model.StorageDevice, error)
	ReconcileStorageDevice(storageId string, createInput *model.CreateStorageDeviceParam) (*model.StorageDevice, error)

	// Storage port management
	GetStoragePorts(storageId string) (*model.StoragePorts, error)

	// Parity group management
	GetParityGroups(storageId string) (*model.ParityGroups, error)

	// Host Group management
	GetHostGroups(storageId string, port string) (*model.HostGroups, error)
	GetHostGroup(storageId, port, hostGroupName string) (*model.HostGroup, error, bool)
	GetHostGroupById(storageId string, hostGrId string) (*model.HostGroup, error)
	ReconcileHostGroup(storageId string, createInput *model.CreateHostGroupParam) (*model.HostGroup, error)

	// Storage pool management
	GetStoragePools(storageId string) (*model.StoragePools, error)
	GetStoragePool(storageId, poolId string) (*model.StoragePool, error)

	//Iscsi target management
	GetIscsiTargets(storageId string, port string) (*model.IscsiTargets, error)
	GetIscsiTarget(storageId, port, iscsiName string) (*model.IscsiTarget, error, bool)
	GetIscsiTargetById(storageId string, iscsiTargetId string) (*model.IscsiTarget, error)
	ReconcileIscsiTarget(storageId string, createInput *model.CreateIscsiTargetParam) (*model.IscsiTarget, error)

	// Volume management
	GetVolumes(storageId string) (*model.Volumes, error)
	GetVolumeByName(storageId string, volumeName string) (*model.VolumeInfo, bool)
	DeleteVolume(storageId string, volumeName string) error
	GetVolumeByID(storageId string, volumeId string) (*model.VolumeInfo, error)
	ReconcileVolume(storageId string, createInput *model.CreateVolumeParams, volumeID *string) (*model.VolumeInfo, error)
	GetVolumesFromLdevIds(id string, fromLdevId *int, toLdevId *int) (*model.Volumes, error)
	GetVolumesByPartnerSubscriberID(id string, fromLdevId int, toLdevId int) (*model.MTVolumes, error)

	
	// UCP System Management
	GetUcpSystems() (*model.UcpSystems, error)

	// CreateHostGroup(storageId string, reqBody model.CreateHostGroupParam) (task *model.TaskResponse, err error)

	//MT Management

	GetPartnerAndSubscriberId(userName string) (*model.MTDetails, error)
}
