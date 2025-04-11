package vssbstorage

import (
	provisonermodel "terraform-provider-hitachi/hitachi/storage/vosb/provisioner/model"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"
)

type VssbStorageManager interface {
	// COMPUTE NODE
	GetComputeNode(serverID string) (*vssbmodel.ComputeNodeWithPathDetails, error)
	GetAllComputeNodes() (*vssbmodel.Servers, error)
	ReconcileComputeNode(inputCompute *vssbmodel.ComputeResource) (*vssbmodel.ComputeNodeWithPathDetails, error)
	DeleteComputeNodeResource(serverId string) error
	GetComputeNodeInformationByName(computeName string, id string) (*provisonermodel.ComputeNodeWithPathDetails, error)
	// VOLUME
	GetAllVolumes(computeNodeName string) (*vssbmodel.Volumes, error)
	GetVolumeDetails(volumeName string) (*vssbmodel.Volume, error)
	ReconcileVolume(postData *vssbmodel.CreateVolume) (*vssbmodel.Volume, error)
	DeleteVolumeResource(volumeID *string) error
	// STORAGE
	GetStorageVersionInfo() (*vssbmodel.StorageVersionInfo, error)
	GetDashboardInfo() (*vssbmodel.Dashboard, error)
	// STORAGE POOLS
	GetAllStoragePools() (*vssbmodel.StoragePools, error)
	GetStoragePoolsByPoolNames(poolNames []string) (*vssbmodel.StoragePools, error)
	GetStoragePoolByPoolName(poolName string) (*vssbmodel.StoragePool, error)
	// STORAGE PORTS
	GetStoragePorts() (*vssbmodel.StoragePorts, error)
	GetPort(portName string) (*vssbmodel.StoragePort, *vssbmodel.PortAuthSettings, error)

	GetChapUsersAllowedToAccessPort(portID string) (*vssbmodel.ChapUsers, error)

	//CHAP USERS
	GetAllChapUsers() (*vssbmodel.ChapUsers, error)
	GetChapUserInfoById(chapUserId string) (*vssbmodel.ChapUser, error)
	GetChapUserInfoByName(targetChapUserName string) (*vssbmodel.ChapUser, error)
	//CreateChapUser(chapUserResource *vssbmodel.ChapUserReq) error
	ReconcileChapUser(inputChapUser *vssbmodel.ChapUserReq) (*vssbmodel.ChapUser, error)
	DeleteChapUser(chapUserId string) error

	// COMPUTE PORT
	AllowChapUsersToAccessComputePort(portId string, authMode string, inputChapUsers []string) error
	GetPortInfoByID(portId string) (*vssbmodel.PortDetailSettings, error)
	//GetIscsiPortAuthInfo(portId string) (*vssbmodel.PortDetailSettings, error)

	// STORAGE USER
	ChangeUserPassword(userID string, reqBody *vssbmodel.ChangeUserPasswordReq) (*vssbmodel.User, error) 
}
