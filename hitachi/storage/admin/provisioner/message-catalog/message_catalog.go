package messagecatalog

// GetMessage .
func GetMessage(id interface{}) string {
	str := MessageCatalog[id]
	return str
}

// for localization, one option is to replace this file with localized resource
var MessageCatalog = map[interface{}]string{
	Default1: "%v",

	// Messages for GetStorageAdminInfo operation
	INFO_GET_STORAGEADMININFO_BEGIN: "Getting storage admin information",
	INFO_GET_STORAGEADMININFO_END:   "Sucessfully getting storage admin information",
	ERR_GET_STORAGEADMININFO_FAILED: "Failed to get storage admin information",

	// Messages for GetVolumeQosAdminInfo operation
	INFO_GET_VOLUME_QOS_ADMIN_BEGIN: "Getting volume QoS admin information",
	INFO_GET_VOLUME_QOS_ADMIN_END:   "Successfully got volume QoS admin information",
	ERR_GET_VOLUME_QOS_ADMIN_FAILED: "Failed to get volume QoS admin information",

	// VOLUMES
	INFO_GET_VOLUMES_BEGIN: "Reading volumes from storage system %v.",
	INFO_GET_VOLUMES_END:   "Successfully read volumes from storage system %v.",
	ERR_GET_VOLUMES_FAILED: "Failed to read volumes from storage system %v.",

	INFO_GET_VOLUME_BY_ID_BEGIN: "Reading volume information for volume %v from storage system %v.",
	INFO_GET_VOLUME_BY_ID_END:   "Successfully read volume information for volume %v from storage system %v.",
	ERR_GET_VOLUME_BY_ID_FAILED: "Failed to read volume information for volume %v from storage system %v.",

	INFO_CREATE_VOLUMES_BEGIN: "Creating %v volume(s) in pool %v on storage system %v.",
	INFO_CREATE_VOLUMES_END:   "Successfully created %v volume(s) in pool %v on storage system %v.",
	ERR_CREATE_VOLUMES_FAILED: "Failed to create %v volume(s) in pool %v on storage system %v.",

	INFO_DELETE_VOLUME_BEGIN: "Deleting volume %v from storage system %v.",
	INFO_DELETE_VOLUME_END:   "Successfully deleted volume %v from storage system %v.",
	ERR_DELETE_VOLUME_FAILED: "Failed to delete volume %v from storage system %v.",

	INFO_EXPAND_VOLUME_BEGIN: "Expanding volume %d on storage system %v.",
	INFO_EXPAND_VOLUME_END:   "Successfully expanded volume %d on storage system %v.",
	ERR_EXPAND_VOLUME_FAILED: "Failed to expand volume %d on storage system %v.",

	INFO_UPDATE_VOLUME_NICKNAME_BEGIN: "Updating nickname for volume %d on storage system %v.",
	INFO_UPDATE_VOLUME_NICKNAME_END:   "Successfully updated nickname for volume %d on storage system %v.",
	ERR_UPDATE_VOLUME_NICKNAME_FAILED: "Failed to update nickname for volume %d on storage system %v.",

	INFO_UPDATE_VOLUME_REDUCTION_BEGIN: "Updating capacity reduction settings for volume %d on storage system %v.",
	INFO_UPDATE_VOLUME_REDUCTION_END:   "Successfully updated capacity reduction settings for volume %d on storage system %v.",
	ERR_UPDATE_VOLUME_REDUCTION_FAILED: "Failed to update capacity reduction settings for volume %d on storage system %v.",

	INFO_GET_VOLUME_SERVER_CONNECTIONS_BEGIN: "Reading volume-server connections from storage system %v.",
	INFO_GET_VOLUME_SERVER_CONNECTIONS_END:   "Successfully read volume-server connections from storage system %v.",
	ERR_GET_VOLUME_SERVER_CONNECTIONS_FAILED: "Failed to read volume-server connections from storage system %v.",

	INFO_GET_VOLUME_SERVER_CONNECTION_BEGIN: "Reading volume-server connection (volume %v, server %v) from storage system %v.",
	INFO_GET_VOLUME_SERVER_CONNECTION_END:   "Successfully read volume-server connection (volume %v, server %v) from storage system %v.",
	ERR_GET_VOLUME_SERVER_CONNECTION_FAILED: "Failed to read volume-server connection (volume %v, server %v) from storage system %v.",

	INFO_ATTACH_VOLUME_SERVER_CONNECTION_BEGIN: "Attaching volume(s) to server(s) on storage system %v.",
	INFO_ATTACH_VOLUME_SERVER_CONNECTION_END:   "Successfully attached volume(s) to server(s) on storage system %v.",
	ERR_ATTACH_VOLUME_SERVER_CONNECTION_FAILED: "Failed to attach volume(s) to server(s) on storage system %v.",

	INFO_DETACH_VOLUME_SERVER_CONNECTION_BEGIN: "Detaching volume %v from server %v on storage system %v.",
	INFO_DETACH_VOLUME_SERVER_CONNECTION_END:   "Successfully detached volume %v from server %v on storage system %v.",
	ERR_DETACH_VOLUME_SERVER_CONNECTION_FAILED: "Failed to detach volume %v from server %v on storage system %v.",

	// ISCSI TARGET
	INFO_GET_ISCSI_TARGETS_BEGIN: "Retrieving iSCSI targets for server %d on storage system %v.",
	INFO_GET_ISCSI_TARGETS_END:   "Successfully retrieved iSCSI targets for server %d on storage system %v.",
	ERR_GET_ISCSI_TARGETS_FAILED: "Failed to retrieve iSCSI targets for server %d on storage system %v.",

	INFO_GET_ISCSI_TARGET_BY_PORT_BEGIN: "Retrieving iSCSI target on server %d, port %s, on storage system %v.",
	INFO_GET_ISCSI_TARGET_BY_PORT_END:   "Successfully retrieved iSCSI target on server %d, port %s, on storage system %v.",
	ERR_GET_ISCSI_TARGET_BY_PORT_FAILED: "Failed to retrieve iSCSI target on server %d, port %s, on storage system %v.",

	INFO_CHANGE_ISCSI_TARGET_NAME_BEGIN: "Changing iSCSI target name for server %d, port %s, to %s on storage system %v.",
	INFO_CHANGE_ISCSI_TARGET_NAME_END:   "Successfully changed iSCSI target name for server %d, port %s, on storage system %v.",
	ERR_CHANGE_ISCSI_TARGET_NAME_FAILED: "Failed to change iSCSI target name for server %d, port %s, on storage system %v.",

	INFO_ADD_HOSTGROUPS_TO_SERVER_BEGIN: "Adding hostgroups to server %d on storage system %v.",
	INFO_ADD_HOSTGROUPS_TO_SERVER_END:   "Successfully added hostgroups to server %d on storage system %v.",
	ERR_ADD_HOSTGROUPS_TO_SERVER_FAILED: "Failed to add hostgroups to server %d on storage system %v.",

	INFO_SYNC_HOSTGROUPS_WITH_SERVER_BEGIN: "Synchronizing hostgroups names with server nickname for server %d on storage system %v.",
	INFO_SYNC_HOSTGROUPS_WITH_SERVER_END:   "Successfully synchronized hostgroups names with server nickname for server %d on storage system %v.",
	ERR_SYNC_HOSTGROUPS_WITH_SERVER_FAILED: "Failed to synchronize hostgroups names with server nickname for server %d on storage system %v.",

	// Server
	INFO_GET_SERVERS_BEGIN: "Reading admin servers from Mgmt system %v.",
	INFO_GET_SERVERS_END:   "Successfully read admin servers from Mgmt system %v.",
	ERR_GET_SERVERS_FAILED: "Failed to read admin servers from Mgmt system %v.",

	// Messages for GetAdminServerInfo operation
	INFO_GET_SERVER_INFO_BEGIN: "Getting admin server information for server %v.",
	INFO_GET_SERVER_INFO_END:   "Successfully got admin server information for server %v.",
	ERR_GET_SERVER_INFO_FAILED: "Failed to get admin server information for server %v.",

	// Messages for CreateAdminServer operation
	INFO_CREATE_SERVER_BEGIN: "Creating admin server %s on storage system %v.",
	INFO_CREATE_SERVER_END:   "Successfully created admin server %s on storage system %v with ID %d.",
	ERR_CREATE_SERVER_FAILED: "Failed to create admin server %s on storage system %v: %v.",

	// Messages for UpdateAdminServer operation
	INFO_UPDATE_SERVER_BEGIN: "Updating admin server %d on storage system %v.",
	INFO_UPDATE_SERVER_END:   "Successfully updated admin server %d on storage system %v.",
	ERR_UPDATE_SERVER_FAILED: "Failed to update admin server %d on storage system %v: %v.",

	// Messages for DeleteAdminServer operation
	INFO_DELETE_SERVER_BEGIN: "Deleting admin server %d on storage system %v.",
	INFO_DELETE_SERVER_END:   "Successfully deleted admin server %d on storage system %v.",
	ERR_DELETE_SERVER_FAILED: "Failed to delete admin server %d on storage system %v: %v.",

	// Messages for SetAdminServerPath operation
	INFO_SET_SERVER_PATH_BEGIN: "Setting admin server path for server %d on storage system %v.",
	INFO_SET_SERVER_PATH_END:   "Successfully set admin server path for server %d on storage system %v.",
	ERR_SET_SERVER_PATH_FAILED: "Failed to set admin server path for server %d on storage system %v: %v.",

	// Messages for DeleteAdminServerPath operation
	INFO_DELETE_SERVER_PATH_BEGIN: "Deleting admin server path for server %d on storage system %v.",
	INFO_DELETE_SERVER_PATH_END:   "Successfully deleted admin server path for server %d on storage system %v.",
	ERR_DELETE_SERVER_PATH_FAILED: "Failed to delete admin server path for server %d on storage system %v: %v.",

	// Messages for GetAdminServerPath operation
	INFO_GET_SERVER_PATH_BEGIN: "Getting admin server path for server %d on storage system %v.",
	INFO_GET_SERVER_PATH_END:   "Successfully retrieved admin server path for server %d on storage system %v.",
	ERR_GET_SERVER_PATH_FAILED: "Failed to get admin server path for server %d on storage system %v: %v.",

	// PORTS
	INFO_GET_PORTS_BEGIN: "Reading ports from storage system %v.",
	INFO_GET_PORTS_END:   "Successfully read ports from storage system %v.",
	ERR_GET_PORTS_FAILED: "Failed to read ports from storage system %v.",

	INFO_GET_PORT_BY_ID_BEGIN: "Reading port information for port %v from storage system %v.",
	INFO_GET_PORT_BY_ID_END:   "Successfully read port information for port %v from storage system %v.",
	ERR_GET_PORT_BY_ID_FAILED: "Failed to read port information for port %v from storage system %v.",

	INFO_UPDATE_PORT_BEGIN: "Updating port %v on storage system %v.",
	INFO_UPDATE_PORT_END:   "Successfully updated port %v on storage system %v.",
	ERR_UPDATE_PORT_FAILED: "Failed to update port %v on storage system %v.",

	// SERVER HBAS
	INFO_GET_SERVER_HBAS_BEGIN: "Reading server HBAs for server %v from storage system %v.",
	INFO_GET_SERVER_HBAS_END:   "Successfully read server HBAs for server %v from storage system %v.",
	ERR_GET_SERVER_HBAS_FAILED: "Failed to read server HBAs for server %v from storage system %v.",

	INFO_GET_SERVER_HBA_BEGIN: "Reading server HBA information for server %v, HBA WWN %v from storage system %v.",
	INFO_GET_SERVER_HBA_END:   "Successfully read server HBA information for server %v, HBA WWN %v from storage system %v.",
	ERR_GET_SERVER_HBA_FAILED: "Failed to read server HBA information for server %v, HBA WWN %v from storage system %v.",

	// POOLS
	INFO_GET_POOLS_BEGIN: "Reading pools from storage system %v.",
	INFO_GET_POOLS_END:   "Successfully read pools from storage system %v.",
	ERR_GET_POOLS_FAILED: "Failed to read pools from storage system %v.",

	INFO_GET_POOL_INFO_BEGIN: "Reading pool information for pool %v from storage system %v.",
	INFO_GET_POOL_INFO_END:   "Successfully read pool information for pool %v from storage system %v.",
	ERR_GET_POOL_INFO_FAILED: "Failed to read pool information for pool %v from storage system %v.",

	INFO_CREATE_POOL_BEGIN: "Creating pool %v on storage system %v.",
	INFO_CREATE_POOL_END:   "Successfully created pool %v on storage system %v.",
	ERR_CREATE_POOL_FAILED: "Failed to create pool %v on storage system %v: %v.",

	INFO_UPDATE_POOL_BEGIN: "Updating pool %v on storage system %v.",
	INFO_UPDATE_POOL_END:   "Successfully updated pool %v on storage system %v.",
	ERR_UPDATE_POOL_FAILED: "Failed to update pool %v on storage system %v: %v.",

	INFO_DELETE_POOL_BEGIN: "Deleting pool %v on storage system %v.",
	INFO_DELETE_POOL_END:   "Successfully deleted pool %v on storage system %v.",
	ERR_DELETE_POOL_FAILED: "Failed to delete pool %v on storage system %v: %v.",
}
