package vssbstorage

import (
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/gateway/model"
)

type VssbStorageManager interface {
	// COMPUTE NODE
	GetComputeNode(serverID string) (*vssbmodel.Server, error)
	GetAllComputeNodes() (*vssbmodel.Servers, error)
	RegisterComputeNode(reqBody *vssbmodel.ComputeNodeCreateReq) error
	DeleteComputeNode(serverID string) (err error)
	EditComputeNode(serverID string, reqBody *vssbmodel.ComputeNodeInformation) error
	RegisterInitiatorInfoForComputeNode(serverID string, reqBody *vssbmodel.RegisterInitiator) error
	RegisterHbaInfoForComputeNode(serverID string, reqBody *vssbmodel.RegisterHba) error
	ConfigureHbaPortsForComputeNode(serverID string) error
	DeleteInitiatorInfoForComputeNode(serverID, hbaID string) error
	GetInitiatorInformationForComputeNode(serverID, hbaID string) (*vssbmodel.Initiator, error)
	GetInitiatorsInformationForComputeNode(serverID string) (*vssbmodel.Initiators, error)
	GetPathInfoForComputeNode(serverID string, reqBody *vssbmodel.ComputeNodePathReq) (*vssbmodel.ComputeNodePath, error)
	GetPathsInfoForComputeNode(serverID string) (*vssbmodel.ComputeNodePaths, error)
	AddPathInfoToComputeNode(serverID string, reqBody *vssbmodel.ComputeNodePathReq) error
	DeleteComputeNodePath(serverID string, reqBody *vssbmodel.ComputeNodePathReq) error
	GetConnectionInfoBtwnVolumeAndServerByVolumeID(volumeID string) (*vssbmodel.VolumeServerConnectionsInfo, error)
	GetConnectionInfoBtwnVolumeAndServerByServerID(serverID string) (*vssbmodel.VolumeServerConnectionsInfo, error)
	GetConnectionInfoBtwnVolumeAndServerBoth(volumeID, serverID string) (*vssbmodel.VolumeServerConnectionInfo, error)
	SetPathBtwnVolumeAndServer(reqBody *vssbmodel.SetPathVolumeServerReq) error
	ReleaseMultipleConnectionsBtwnVolumeAndServer(reqBody *vssbmodel.ReleaseMultiConVolumeServerReq) error
	ReleaseConnectionBtwnVolumeAndServer(reqBody *vssbmodel.SetPathVolumeServerReq) error
	// VOLUME INFO
	GetAllVolumes() (*vssbmodel.Volumes, error)
	CreateVolume(reqBody *vssbmodel.CreateVolumeRequestGwy) (*int, error)
	AddVolumeToComputeNode(reqBody *vssbmodel.AddVolumeToComputeNodeReq) (*int, error)
	UpdateVolume(volumeId *string, reqbody *vssbmodel.UpdateVolumeReq) error
	ExtendVolumeSize(volumeId *string, capacity *vssbmodel.UpdateVolumeSizeReq) error
	RemoveVolumeFromComputeNode(volumeId *string, serverId *string) error
	DeleteVolume(volumeId *string) error
	// STORAGE
	GetStorageVersionInfo() (*vssbmodel.StorageVersionInfo, error)
	GetHealthStatuses() (*vssbmodel.HealthStatuses, error)
	GetStorageClusterInfo() (*vssbmodel.StorageClusterInfo, error)
	GetDrivesInfo(status string) (*vssbmodel.Drives, error)

	// STORAGE POOLS
	GetAllStoragePools() (*vssbmodel.StoragePools, error)
	GetStoragePoolsByPoolNames(poolNames []string) (*vssbmodel.StoragePools, error)
	ExpandStoragePool(poolId string, reqBody *vssbmodel.ExpandStoragePoolReq) error

	// STORAGE NODES
	GetStorageNodes() (*vssbmodel.StorageNodes, error)
	GetStorageNode(node string) (*vssbmodel.StorageNode, error)
	AddStorageNode(configurationFile string, exportedConfigurationFile string, setupUserPassword string, expectedCloudProvider string, vmConfigFileS3URI string) error

	// STORAGE PORTS
	GetStoragePorts() (*vssbmodel.StoragePorts, error)
	GetPort(portId string) (*vssbmodel.StoragePort, error)
	GetPortAuthSettings(portID string) (*vssbmodel.PortAuthSettings, error)

	// CHAP USERS
	GetAllChapUsers() (*vssbmodel.ChapUsers, error)
	CreateChapUser(reqBody *vssbmodel.ChapUserReq) error
	DeleteChapUser(chapUserID string) error
	GetChapUserInfo(chapUserID string) (*vssbmodel.ChapUser, error)
	UpdateChapUser(chapUserID string, reqBody *vssbmodel.ChapUserReq) error

	GetChapUsersAllowedToAccessPort(portID string) (*vssbmodel.ChapUsers, error)
	DeletePortAccessForChapUser(portId, chapUserId string) error
	GetChapUserInfoAllowedToAccessPort(portId, chapUserId string) (*vssbmodel.ChapUsers, error)
	UpdatePortAuthSettings(portID string, reqBody *vssbmodel.PortAuthSettings) error
	AllowChapUserToAccessPort(portID string, reqBody *vssbmodel.ChapUserIdReq) error

	// STORAGE CREDENTIAL
	ChangeUserPassword(userId string, reqBody *vssbmodel.ChangeUserPasswordReq) (*vssbmodel.User, error)

	// CONFIGURATION FILE
	RestoreConfigurationDefinitionFile(createConfigParam *vssbmodel.CreateConfigurationFileParam) error
	DownloadConfigurationFile(saveDir string) (string, error)
}
