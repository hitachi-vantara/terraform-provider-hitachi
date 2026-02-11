package admin

import (
	model "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
)

type AdminStorageManager interface {
	GetStorageAdminInfo(configurable_capacities bool) (*model.StorageAdminInfo, error)

	GetVolumeQosAdminInfo(volumeID int) (*model.VolumeQosResponse, error)
	SetVolumeQosAdminThreshold(volumeID int, threshold model.VolumeQosThreshold) error
	SetVolumeQosAdminAlertSetting(volumeID int, alert model.VolumeQosAlertSetting) error

	GetVolumes(params model.GetVolumeParams) (*model.VolumeInfoList, error)
	GetVolumeByID(id int) (*model.VolumeInfoByID, error)
	CreateVolume(params model.CreateVolumeParams) (string, error)
	DeleteVolume(volumeID int) error
	ExpandVolume(volumeID int, params model.ExpandVolumeParams) error
	UpdateVolumeNickname(volumeID int, params model.UpdateVolumeNicknameParams) error
	UpdateVolumeReductionSettings(volumeID int, params model.UpdateVolumeReductionParams) error

	GetVolumeServerConnections(params model.GetVolumeServerConnectionsParams) (*model.VolumeServerConnectionsResponse, error)
	GetOneVolumeServerConnection(volumeId, serverId int) (*model.VolumeServerConnectionDetail, error)
	AttachVolumeToServers(params model.AttachVolumeServerConnectionParam) (string, error)
	DetachVolumeToServers(volumeId, serverId int) error

	// Server Management methods
	GetAdminServerList(params model.AdminServerListParams) (*model.AdminServerListResponse, error)
	GetAdminServerInfo(serverID int) (*model.AdminServerInfo, error)
	CreateAdminServer(params model.CreateAdminServerParams) error
	UpdateAdminServer(serverID int, params model.UpdateAdminServerParams) error
	DeleteAdminServer(serverID int, params model.DeleteAdminServerParams) error

	GetIscsiTargets(serverId int) (*model.IscsiTargetInfoList, error)
	GetIscsiTargetByPort(serverId int, portId string) (*model.IscsiTargetInfoByPort, error)
	ChangeIscsiTargetName(serverId int, portId string, targetIscsiName string) error

	AddHostGroupsToServer(serverId int, params model.AddHostGroupsToServerParam) error
	SyncHostGroupsWithServer(serverId int) error

	SetAdminServerPath(serverID int, params model.SetAdminServerPathParams) error
	DeleteAdminServerPath(serverID int, params model.DeleteAdminServerPathParams) error
	GetAdminServerPath(params model.AdminServerPathParams) (*model.AdminServerPathInfo, error)

	// Port management methods
	GetPorts(params model.GetPortParams) (*model.PortInfoList, error)
	GetPortByID(id string) (*model.PortInfo, error)
	UpdatePort(id string, params model.UpdatePortParams) error

	// Server HBA management methods
	GetServerHBAs(serverID int) (*model.ServerHBAList, error)
	GetServerHBAByWwn(serverID int, hbaWwn string) (*model.ServerHBA, error)
	CreateServerHBAs(serverID int, params model.CreateServerHBAParams) (*model.ServerHBAList, error)
	DeleteServerHBA(serverID int, initiatorName string) (*model.ServerHBAList, error)

	// Pool management methods
	GetAdminPoolList(params model.AdminPoolListParams) (*model.AdminPoolListResponse, error)
	GetAdminPoolInfo(poolID int) (*model.AdminPool, error)
	CreateAdminPool(params model.CreateAdminPoolParams) error
	UpdateAdminPool(poolID int, params model.UpdateAdminPoolParams) error
	ExpandAdminPool(poolID int, params model.ExpandAdminPoolParams) error
	DeleteAdminPool(poolID int) error
}
