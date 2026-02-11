package admin

import (
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/admin/reconciler/model"
)

type AdminStorageManager interface {
	GetStorageAdminInfo(configurable_capacities bool) (*reconcilermodel.StorageAdminInfo, error)

	ReconcileReadAdminVolumes(volumeIDs []int) ([]gwymodel.VolumeInfoByID, []int, error)
	ReconcileDeleteAdminVolumes(volumeIDs []int) error
	ReconcileCreateAdminVolumes(params gwymodel.CreateVolumeParams) ([]int, error)
	ReconcileUpdateAdminVolume(volumeID int, params gwymodel.CreateVolumeParams) error

	// Admin Server CRUD operations
	ReconcileCreateAdminServer(params gwymodel.CreateAdminServerParams, addHgParams gwymodel.AddHostGroupsToServerParam) (int, error)
	ReconcileReadAdminServer(serverID int) (*gwymodel.AdminServerInfo, error)
	ReconcileUpdateAdminServer(serverID int, params gwymodel.UpdateAdminServerParams, addHgParams gwymodel.AddHostGroupsToServerParam) (int, error)
	ReconcileDeleteAdminServer(serverID int, params gwymodel.DeleteAdminServerParams) error
	ReconcileSetAdminServerPath(serverID int, params gwymodel.SetAdminServerPathParams) error
	ReconcileDeleteAdminServerPath(serverID int, params gwymodel.DeleteAdminServerPathParams) error

	ReconcileReadVolumeServerConnections(pairs []reconcilermodel.VolumeServerPair) ([]gwymodel.VolumeServerConnectionDetail, []reconcilermodel.VolumeServerPair, error)
	ReconcileDeleteVolumeServerConnections(pairs []reconcilermodel.VolumeServerPair) error
	ReconcileUpdateVolumeServerConnections(existingPairs, desiredPairs []reconcilermodel.VolumeServerPair) error

	// Admin Server data source operations
	GetAdminServerList(params gwymodel.AdminServerListParams) (*gwymodel.AdminServerListResponse, error)
	GetAdminServerInfo(serverID int) (*gwymodel.AdminServerInfo, error)
	GetAdminServerPath(params gwymodel.AdminServerPathParams) (*gwymodel.AdminServerPathInfo, error)
	ReconcileUpdateAdminPort(portID string, params gwymodel.UpdatePortParams) error
	ReconcileReadAdminPort(portID string) (*gwymodel.PortInfo, error)

	// Server HBA management operations
	ReconcileCreateAdminServerHBAs(serverID int, params gwymodel.CreateServerHBAParams) (*gwymodel.ServerHBAList, error)
	ReconcileDeleteAdminServerHBA(serverID int, initiatorName string) (*gwymodel.ServerHBAList, error)

	// Pool management operations
	ReconcileCreateAdminPool(params gwymodel.CreateAdminPoolParams) (int, error)
	ReconcileReadAdminPool(poolID int) (*gwymodel.AdminPool, bool, error)
	ReconcileUpdateAdminPool(poolID int, params gwymodel.UpdateAdminPoolParams) error
	ReconcileExpandAdminPool(poolID int, params gwymodel.ExpandAdminPoolParams) error
	ReconcileDeleteAdminPool(poolID int) error
	GetAdminPoolList(params gwymodel.AdminPoolListParams) (*gwymodel.AdminPoolListResponse, error)
	GetAdminPoolInfo(poolID int) (*gwymodel.AdminPool, error)
}
