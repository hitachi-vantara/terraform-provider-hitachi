package vssbstorage

import (
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

type InfraGwManager interface {

	// Storage Device Management
	GetStorageDevices() (*model.StorageDevices, error)
	GetStorageDevice(storageId string) (*model.StorageDevice, error)
	AddStorageDevice(reqBody model.CreateStorageDeviceParam) (*string, error)
	UpdateStorageDevice(storageId string, reqBody model.PatchStorageDeviceParam) (*string, error)
	AddMTStorageDevice(reqBody model.CreateStorageDeviceParam) (*string, error)
	// Storage Port Management
	GetStoragePorts(storageId string) (*model.StoragePorts, error)
	GetParityGroups(storageId string) (*model.ParityGroups, error)

	// Hostgroups Management
	GetHostGroups(storageId string, port string) (*model.HostGroups, error)
	GetHostGroup(storageId, hostGroupId string) (*model.HostGroup, error)
	CreateHostGroup(storageId string, reqBody model.CreateHostGroupParam) (*string, error)
	UpdateHostGroup(storageId, hostGroupId string, reqBody model.PatcheHostGroupParam) (*string, error)
	UpdateHostMode(storageId, iscsiTargetId string, reqBody model.UpdateHostModeParam) (*string, error)

	// Storage Pool Management
	GetStoragePools(storageId string) (*model.StoragePools, error)
	GetStoragePool(storageId, poolId string) (*model.StoragePool, error)

	// Iscsi Management
	GetIscsiTargets(storageId string, port string) (*model.IscsiTargets, error)
	GetIscsiTarget(storageId string, iscsiTargetId string) (*model.IscsiTarget, error)
	CreateIscsiTarget(storageId string, reqBody model.CreateIscsiTargetParam) (*string, error)

	//Volume management
	GetVolumes(storageId string) (*model.Volumes, error)
	GetVolumeByID(storageId string, volumeID string) (*model.Volume, error)
	CreateVolume(storageId string, reqBody *model.CreateVolumeParams) (*string, error)
	CreateMTVolume(storageId string, reqBody *model.CreateVolumeParams) (*string, error)
	UpdateVolume(storageId string, volumeID string, reqBody *model.UpdateVolumeParams) (*string, error)
	DeleteVolume(storageId string, volumeID string) error
	DeleteMTVolume(storageId string, volumeID string) error
	GetVolumesByPartnerSubscriberID(storageId string, fromLdevId int, toLdevId int) (*model.MTVolumes, error)
	GetVolumesFromLdevIds(id string, fromLdevId int, toLdevId int) (*model.Volumes, error)

	//UCP System Management
	GetUcpSystems() (*model.UcpSystems, error)
	GetUcpSystemById(id string) (*model.UcpSystem, error)
	CreateUcpSystem(reqBody model.CreateUcpSystemParam) (*string, error)

	//User Managements
	GetAllUsers() (*model.UserWithDetails, error)
	GetUsersRoles(userId string) (*model.RoleWithDetails, error)

	//Multi-tenancy Management
	GetAllPartners() (*[]model.Partner, error)
	GetPartner(partnerId string) (*model.Partner, error)
	GetAllSubscribers(partnerId string) (*[]model.Subscriber, error)
	GetSubscriber(partnerId string, subscriberId string) (*model.Subscriber, error)
	GetSubscriberResources(partnerId string, subscriberId string) (*model.SubscriberDetails, error)
	RegisterSubscriber(reqBody *model.RegisterSubscriberReq) (*string, error)
	RegisterPartner(reqBody *model.RegisterPartnerReq) (*string, error)
	UnRegisterSubscriber(subscriberId string) (*string, error)
}
