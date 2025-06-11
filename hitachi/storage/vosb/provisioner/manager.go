package vssbstorage

import (
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/provisioner/model"
)

type VssbStorageManager interface {
	// COMPUTE NODE
	GetComputeNode(serverID string) (*vssbmodel.ComputeNodeWithPathDetails, error)
	GetAllComputeNodes() (*vssbmodel.Servers, error)
	GetConnectionInfoBtwnVolumeAndServerByServerID(serverID string) (*vssbmodel.VolumeServerConnectionsInfo, error)
	CreateComputeResource(computeResource *vssbmodel.ComputeResource) error
	GetComputeResourceInfo(computeNodeName string) (*vssbmodel.ComputeResourceOutput, error)
	DeleteComputeNodeResource(serverId string) error
	GetComputeNodeIdByName(computeNodeName string) (string, error)
	UpdateComputeNode(computeNodeId, computeNodeName, osType string) error
	AddStoragePortToComputeNodeHbaByHbaIdAndPortId(computeNodeId, computeNodeHbaId, storagePortId string) error
	RemoveStoragePortFromComputeNodeHbaByHbaIdAndPortId(computeNodeId, computeNodeHbaId, storagePortId string) error
	AddIscsiHbaToComputeNode(computeNodeId, iqn string) error
	DeleteIscsiHbaFromComputeNode(computeNodeId, iqn string) error
	GetInitiatorIdByServerId(serverId string, initiatorName string) (string, error)
	GetPortsIdsByName(ports []string) ([]string, error)
	RemoveFCportsWWN(wwns []string, serverId string) error
	AddFCportsWWN(wwns []string, serverId string) error

	// VOLUME INFO
	GetAllVolumes(computeNodeName string) (*vssbmodel.Volumes, error)
	GetVolumeDetails(volumeName string) (*vssbmodel.Volume, error)
	GetVolumeDetailsByName(volumeName string) (*vssbmodel.Volume, error)
	CreateVolume(name string, nickName string, poolName string, capacity float32) (*int, error)
	AddVolumeToComputeNode(volumeName string, computeNodeName string) (*int, error)
	UpdateVolumeNickName(serverId string, nickName string) error
	ExpandVolume(serverId string, additionalCapacity *int32) error
	RemoveVolumeFromComputeNode(volumeId string, serverId string) error
	DeleteVolume(volumeId string) error
	// STORAGE
	GetStorageVersionInfo() (*vssbmodel.StorageVersionInfo, error)
	GetHealthStatuses() (*vssbmodel.HealthStatuses, error)
	GetStorageClusterInfo() (*vssbmodel.StorageClusterInfo, error)
	GetDrivesInfo(status string) (*vssbmodel.Drives, error)
	// STORAGE POOLS
	GetAllStoragePools() (*vssbmodel.StoragePools, error)
	GetStoragePoolsByPoolNames(poolNames []string) (*vssbmodel.StoragePools, error)
	GetStoragePoolByPoolName(poolName string) (*vssbmodel.StoragePool, error)
	ExpandStoragePool(storagePoolName string, driveIds []string) error
	AddOfflineDrivesToStoragePool(storagePoolName string) error

	// STORAGE PORTS
	GetStoragePorts() (*vssbmodel.StoragePorts, error)
	GetPort(portId string) (*vssbmodel.StoragePort, error)
	GetPortAuthSettings(portId string) (*vssbmodel.PortAuthSettings, error)
	GetChapUsersAllowedToAccessPort(portID string) (*vssbmodel.ChapUsers, error)

	//CHAP USERS
	GetAllChapUsers() (*vssbmodel.ChapUsers, error)
	GetChapUserInfoById(chapUserId string) (*vssbmodel.ChapUser, error)
	GetChapUserInfoByName(chapUserTargetName string) (*vssbmodel.ChapUser, error)
	CreateChapUser(chapUserResource *vssbmodel.ChapUserReq) error
	DeleteChapUser(chapUserID string) error
	UpdateChapUserById(chapUserID string, reqBody *vssbmodel.ChapUserReq) error
	UpdateChapUserInfoByName(reqBody *vssbmodel.ChapUserReq) error

	//COMPUTE PORT
	DeleteAllChapUsersFromComputePort(portId string) error
	UpdatePortAuthSettings(portID string, reqBody *vssbmodel.PortAuthSettings) error
	AddChapUsersToComputePort(portId string, chapUserIds []string) error

	// STORAGE CREDENTIAL
	ChangeUserPassword(userId string, reqBody *vssbmodel.ChangeUserPasswordReq) (*vssbmodel.User, error)

	// CONFIGURATION FILE
	RestoreConfigurationDefinitionFile(createConfigParam *vssbmodel.CreateConfigurationFileParam) error
	DownloadConfigurationFile(saveDir string) (string, error)

}
