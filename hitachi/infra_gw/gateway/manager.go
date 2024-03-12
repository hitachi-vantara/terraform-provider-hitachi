package vssbstorage

import (
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

type InfraGwManager interface {

	// Storage Device Management
	GetStorageDevices() (*model.StorageDevices, error)
	//GetMTStorageDevices() (*model.MTStorageDevices, error)
	GetStorageDevice(storageId string) (*model.StorageDevice, error)
	GetMTStorageDevice(storageId string) (*model.MTStorageDevice, error)
	GetMTStorageDevices() (*[]model.MTStorageDevice, error)
	AddStorageDevice(reqBody model.CreateStorageDeviceParam) (*string, error)
	AddMTStorageDevice(reqBody model.CreateMTStorageDeviceParam) (*string, error)
	UpdateStorageDevice(storageId string, reqBody model.PatchStorageDeviceParam) (*string, error)
	DeleteStorageDevice(storageId string) error
	DeleteStorageDeviceFromUcp(id, storageId string) error
	DeleteMTStorageDevice(storageId string) error

	// Storage Port Management
	GetStoragePorts(storageId string) (*model.StoragePorts, error)
	GetStoragePortsByPartnerIdOrSubscriberId(id string) (*model.MTPorts, error)
	GetParityGroups(storageId string) (*model.ParityGroups, error)

	// Hostgroups Management
	GetHostGroups(storageId string, port string) (*model.HostGroups, error)
	GetHostGroupsByPartnerIdOrSubscriberID(storageId string) (*model.MTHostGroups, error)
	GetHostGroup(storageId, hostGroupId string) (*model.HostGroup, error)
	CreateHostGroup(storageId string, reqBody model.CreateHostGroupParam) (*string, error)
	CreateMTHostGroup(storageId string, reqBody model.CreateHostGroupParam) (*string, error)
	UpdateHostGroup(storageId, hostGroupId string, reqBody model.PatcheHostGroupParam) (*string, error)
	AddVolumesToHostGroup(storageId, hostGroupId string, reqBody model.AddVolumesToHostGroupParam) (*string, error)
	AddVolumesToHostGroupToSubscriber(storageId, hostGroupId string, reqBody model.AddVolumesToHostGroupParam) (*string, error)
	DeleteVolumesFromHostGroup(storageId, hostGroupId string, reqBody model.DeleteVolumesToHostGroupParam) error
	DeleteVolumesFromHostGroupFromSubscriber(storageId, hostGroupId string, reqBody model.DeleteVolumesToHostGroupParam) error
	DeleteHostGroup(storageId, hostGroupId string) error
	DeleteMTHostGroup(storageId, hostGroupId string) error

	// Storage Pool Management
	GetStoragePools(storageId string) (*model.StoragePools, error)
	GetStoragePool(storageId, poolId string) (*model.StoragePool, error)

	// Iscsi Management
	GetIscsiTargets(storageId string, port string) (*model.IscsiTargets, error)
	GetIscsiTarget(storageId string, iscsiTargetId string) (*model.IscsiTarget, error)
	CreateIscsiTarget(storageId string, reqBody model.CreateIscsiTargetParam) (*string, error)
	UpdateHostMode(storageId, iscsiTargetId string, reqBody model.UpdateHostModeParam) (*string, error)
	GetMTIscsiTargets(id string) (*model.IscsiTargets, error)
	GetMTIscsiTarget(id string, iscsiTargetId string) (*model.IscsiTarget, error)
	CreateMTIscsiTarget(storageId string, reqBody model.CreateIscsiTargetParam) (*string, error)
	AddVolumesToMTIscsiTarget(storageId, iscsiTargetId string, reqBody model.AddVolumesToIscsiTargetParam) (*string, error)
	RemoveVolumesFromMTscsiTarget(storageId, iscsiTargetId string, reqBody model.RemoveVolumesFromIscsiTargetParam) (*string, error)
	AddIqnInitiatorsToIscsiMTTarget(storageId, iscsiTargetId string, reqBody model.AddIqnInitiatorsToIscsiTargetParam) (*string, error)
	RemoveIqnInitiatorsFromIscsiMTTarget(storageId, iscsiTargetId string, reqBody model.RemoveIqnInitiatorsFromIscsiTargetParam) (*string, error)
	DeleteMTIscsiTarget(storageId, iscsiTargetId string) (*string, error)
	DeleteIscsiTarget(storageId, iscsiTargetId string) (*string, error)
	UpdateTargetIqnInIscsiTarget(storageId, iscsiTargetId string, reqBody model.UpdateTargetIqnInIscsiTargetParam) (*string, error)
	RemoveIqnInitiatorsFromIscsiTarget(storageId, iscsiTargetId string, reqBody model.RemoveIqnInitiatorsFromIscsiTargetParam) (*string, error)
	AddVolumesToIscsiTarget(storageId, iscsiTargetId string, reqBody model.AddVolumesToIscsiTargetParam) (*string, error)
	RemoveVolumesFromIscsiTarget(storageId, iscsiTargetId string, reqBody model.RemoveVolumesFromIscsiTargetParam) (*string, error)
	AddIqnInitiatorsToIscsiTarget(storageId, iscsiTargetId string, reqBody model.AddIqnInitiatorsToIscsiTargetParam) (*string, error)
	

	//Volume management
	GetVolumes(storageId string) (*model.Volumes, error)
	GetVolumeByID(storageId string, volumeID string) (*model.Volume, error)
	CreateVolume(storageId string, reqBody *model.CreateVolumeParams) (*string, error)
	CreateMTVolume(storageId string, reqBody *model.CreateVolumeParams) (*string, error)
	UpdateVolume(storageId string, volumeID string, reqBody *model.UpdateVolumeParams) (*string, error)
	DeleteVolume(storageId string, volumeID string) error
	DeleteMTVolume(storageId string, volumeID string) error
	GetVolumesByPartnerSubscriberID(storageId string, fromLdevId int, toLdevId int) (*model.MTVolumes, error)
	GetVolumesDetailsByPartnerSubscriberID(storageId string, fromLdevId int, toLdevId int) (*model.MTVolumes, error)
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
	UpdateSubscriber(subscriberId string, partnerId string, reqBody *model.UpdateSubscriberReq) (*string, error)
}
