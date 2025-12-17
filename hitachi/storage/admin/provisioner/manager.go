package admin

import (
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	provmodel "terraform-provider-hitachi/hitachi/storage/admin/provisioner/model"
)

type AdminStorageManager interface {
	GetStorageAdminInfo(configurable_capacities bool) (*provmodel.StorageAdminInfo, error)

	GetVolumeQosAdminInfo(volumeID int) (*provmodel.VolumeQosResponse, error)
	SetVolumeQosAdminThreshold(volumeID int, threshold provmodel.VolumeQosThreshold) error
	SetVolumeQosAdminAlertSetting(volumeID int, alert provmodel.VolumeQosAlertSetting) error

	GetVolumes(queryParams gwymodel.GetVolumeParams) (*gwymodel.VolumeInfoList, error)
	GetVolumeByID(volumeID int) (*gwymodel.VolumeInfoByID, error)
	CreateVolume(params gwymodel.CreateVolumeParams) (string, error)
	DeleteVolume(volumeID int) error
	ExpandVolume(volumeID int, params gwymodel.ExpandVolumeParams) error
	UpdateVolumeNickname(volumeID int, params gwymodel.UpdateVolumeNicknameParams) error
	UpdateVolumeReductionSettings(volumeID int, params gwymodel.UpdateVolumeReductionParams) error

	GetVolumeServerConnections(params gwymodel.GetVolumeServerConnectionsParams) (*gwymodel.VolumeServerConnectionsResponse, error)
	GetOneVolumeServerConnection(volumeId, serverId int) (*gwymodel.VolumeServerConnectionDetail, error)
	AttachVolumeToServers(params gwymodel.AttachVolumeServerConnectionParam) (string, error)
	DetachVolumeFromServer(volumeId, serverId int) error

	GetIscsiTargets(serverId int) (*gwymodel.IscsiTargetInfoList, error)
	GetIscsiTargetByPort(serverId int, portId string) (*gwymodel.IscsiTargetInfoByPort, error)
	ChangeIscsiTargetName(serverId int, portId string, targetIscsiName string) error

	AddHostGroupsToServer(serverId int, params gwymodel.AddHostGroupsToServerParam) error
	SyncHostGroupsWithServer(serverId int) error

	GetAdminServerList(params gwymodel.AdminServerListParams) (*gwymodel.AdminServerListResponse, error)
	GetAdminServerInfo(serverID int) (*gwymodel.AdminServerInfo, error)
	CreateAdminServer(params gwymodel.CreateAdminServerParams) error
	UpdateAdminServer(serverID int, params gwymodel.UpdateAdminServerParams) error
	DeleteAdminServer(serverID int, params gwymodel.DeleteAdminServerParams) error

	SetAdminServerPath(serverID int, params gwymodel.SetAdminServerPathParams) error
	DeleteAdminServerPath(serverID int, params gwymodel.DeleteAdminServerPathParams) error
	GetAdminServerPath(params gwymodel.AdminServerPathParams) (*gwymodel.AdminServerPathInfo, error)

	// Port management methods
	GetPorts(queryParams gwymodel.GetPortParams) (*gwymodel.PortInfoList, error)
	GetPortByID(portID string) (*gwymodel.PortInfo, error)
	UpdatePort(portID string, params gwymodel.UpdatePortParams) error

	// Server HBA management methods
	GetServerHBAs(serverID int) (*gwymodel.ServerHBAList, error)
	GetServerHBAByWwn(serverID int, hbaWwn string) (*gwymodel.ServerHBA, error)
	CreateServerHBAs(serverID int, params gwymodel.CreateServerHBAParams) (*gwymodel.ServerHBAList, error)
	DeleteServerHBA(serverID int, initiatorName string) (*gwymodel.ServerHBAList, error)

	// Pool management methods
	GetAdminPoolList(params gwymodel.AdminPoolListParams) (*gwymodel.AdminPoolListResponse, error)
	GetAdminPoolInfo(poolID int) (*gwymodel.AdminPool, error)
	CreateAdminPool(params gwymodel.CreateAdminPoolParams) error
	UpdateAdminPool(poolID int, params gwymodel.UpdateAdminPoolParams) error
	ExpandAdminPool(poolID int, params gwymodel.ExpandAdminPoolParams) error
	DeleteAdminPool(poolID int) error
}
