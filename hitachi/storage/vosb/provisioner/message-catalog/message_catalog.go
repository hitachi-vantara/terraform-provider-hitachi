package messagecatalog

// GetMessage .
func GetMessage(id interface{}) string {
	str := MessageCatalog[id]
	return str
}

// for localization, one option is to replace this file with localized resource
var MessageCatalog = map[interface{}]string{
	Default1: "%v",
	// COMPUTE NODE
	INFO_GET_ALL_SERVERS_BEGIN:          "Reading all compute node information.",
	ERR_GET_ALL_SERVERS_FAILED:          "Failed to read all compute node information.",
	INFO_GET_ALL_SERVERS_END:            "Successfully read compute node information.",
	INFO_GET_SERVER_BEGIN:               "Reading compute node information for server id %s.",
	ERR_GET_SERVER_FAILED:               "Failed to read compute node information for server id %s.",
	INFO_GET_SERVER_END:                 "Successfully read compute node information for server id %s.",
	INFO_GET_CONNECTION_BY_SERVER_BEGIN: "Reading connection information between volume and server by server id %s.",
	ERR_GET_CONNECTION_BY_SERVER_FAILED: "Failed to read connection information between volume and server by server id %s.",
	INFO_GET_CONNECTION_BY_SERVER_END:   "Successfully read connection information between volume and server by server id %s.",
	INFO_CREATE_SERVER_BEGIN:            "Creating compute resource for server name %s.",
	ERR_CREATE_SERVER_FAILED:            "Failed to create compute resource for server name %s.",
	INFO_CREATE_SERVER_END:              "Successfully created compute resource for server name %s.",
	INFO_DELETE_SERVER_BEGIN:            "Deleting compute node for server %s.",
	ERR_DELETE_SERVER_FAILED:            "Failed to delete compute node server %s.",
	INFO_DELETE_SERVER_END:              "Successfully deleted compute node server %s.",
	INFO_UPDATE_NAME_AND_OS_TYPE_BEGIN:  "Updating compute node name %s and os type %s.",
	ERR_UPDATE_NAME_AND_OS_TYPE_FAILED:  "Failed to update compute node name %s and os type %s.",
	INFO_UPDATE_NAME_AND_OS_TYPE_END:    "Successfully updated compute node name %s and os type %s.",
	INFO_ADD_COMPUTE_PATH_INFO_BEGIN:    "Adding compute node path information for hba id %s and port id %s.",
	ERR_ADD_COMPUTE_PATH_INFO_FAILED:    "Failed to add compute node path information for hba id %s and port id %s.",
	INFO_ADD_COMPUTE_PATH_INFO_END:      "Successfully added compute node path information for hba id %s and port id %s.",
	INFO_REMOVE_COMPUTE_PATH_INFO_BEGIN: "Removing compute node path information for hba id %s and port id %s.",
	ERR_REMOVE_COMPUTE_PATH_INFO_FAILED: "Failed to remove compute node path information for hba id %s and port id %s.",
	INFO_REMOVE_COMPUTE_PATH_INFO_END:   "Successfully removed compute node path information for hba id %s and port id %s.",
	INFO_ADD_COMPUTE_IQN_INFO_BEGIN:     "Adding compute node iqn information for iqn id %s.",
	ERR_ADD_COMPUTE_IQN_INFO_FAILED:     "Failed to add compute node iqn information for iqn id %s.",
	INFO_ADD_COMPUTE_IQN_INFO_END:       "Successfully added compute node iqn information for iqn id %s.",
	INFO_REMOVE_COMPUTE_IQN_INFO_BEGIN:  "Removing compute node iqn information for iqn id %s.",
	ERR_REMOVE_COMPUTE_IQN_INFO_FAILED:  "Failed to remove compute node iqn information for iqn id %s.",
	INFO_REMOVE_COMPUTE_IQN_INFO_END:    "Successfully removed compute node iqn information for iqn id %s.",

	// VOLUME
	INFO_GET_ALL_VOLUME_INFO_BEGIN:            "Reading all volume information.",
	ERR_GET_ALL_VOLUME_INFO_FAILED:            "Failed to read all volume information.",
	INFO_GET_ALL_VOLUME_INFO_END:              "Successfully read all volume information.",
	INFO_ADD_VOLUME_TO_COMPUTE_NODES_BEGIN:    "Adding volume %s to compute nodes.",
	ERR_ADD_VOLUME_TO_COMPUTE_NODES_FAILED:    "Failed to add volume %s to compute nodes.",
	INFO_ADD_VOLUME_TO_COMPUTE_NODES_END:      "Successfully added volume %s to compute nodes.",
	INFO_DELETE_VOLUME_BEGIN:                  "Deleting volume %s.",
	INFO_DELETE_VOLUME_END:                    "Successfully deleted volume %s.",
	ERR_DELETE_VOLUME_FAILED:                  "Failed to delete volume %s.",
	INFO_UPDATE_VOLUME_NICKNAME_BEGIN:         "Updating volume nickname %s.",
	ERR_UPDATE_VOLUME_NICKNAME_FAILED:         "Failed to update volume nickname %s.",
	INFO_UPDATE_VOLUME_NICKNAME_END:           "Successfully updated volume nickname %s.",
	INFO_EXPAND_VOLUME_SIZE_BEGIN:             "Expanding volume size %s.",
	ERR_EXPAND_VOLUME_SIZE_FAILED:             "Failed to expand volume size %s.",
	INFO_EXPAND_VOLUME_SIZE_END:               "Successfully expanded volume size %s.",
	INFO_REMOVE_VOLUME_TO_COMPUTE_NODES_BEGIN: "Removing volume %s from compute nodes.",
	ERR_REMOVE_VOLUME_TO_COMPUTE_NODES_FAILED: "Failed to remove volume %s from compute nodes.",
	INFO_REMOVE_VOLUME_TO_COMPUTE_NODES_END:   "Successfully removed volume %s from compute nodes.",

	// STORAGE POOLS
	INFO_GET_ALL_STORAGE_POOLS_BEGIN: "Reading all storage pool information.",
	ERR_GET_ALL_STORAGE_POOLS_FAILED: "Failed to read all storage pool information.",
	INFO_GET_ALL_STORAGE_POOLS_END:   "Successfully read storage pool information.",
	INFO_GET_STORAGE_POOL_BEGIN:      "Reading storage pool information for pool names %s.",
	ERR_GET_STORAGE_POOL_FAILED:      "Failed to read storage pool information for pool names %s.",
	INFO_GET_STORAGE_POOL_END:        "Successfully read storage pool information for pool names %s.",
	INFO_EXPAND_STORAGE_POOL_BEGIN:   "Expanding storage pool.",
	ERR_EXPAND_STORAGE_POOL_FAILED:   "Failed to expand storage pool.",
	INFO_EXPAND_STORAGE_POOL_END:     "Successfully expanded storage pool.",

	ERR_NO_OFFLINE_DRIVES:              "No offline drives available.",
	INFO_ADD_DRIVES_STORAGE_POOL_BEGIN: "Adding offline drives to storage pool.",
	ERR_ADD_DRIVES_STORAGE_POOL_FAILED: "Failed adding offline drives to storage pool.",
	INFO_ADD_DRIVES_STORAGE_POOL_END:   "Successfully added offline drives to storage pool.",

	// STORAGE PORTS
	INFO_GET_ALL_STORAGE_PORTS_BEGIN: "Reading all storage ports information.",
	ERR_GET_ALL_STORAGE_PORTS_FAILED: "Failed to read all storage ports information.",
	INFO_GET_ALL_STORAGE_PORTS_END:   "Successfully read storage ports information.",

	INFO_GET_PORT_BEGIN:               "Reading port information for port id %s",
	ERR_GET_PORT_FAILED:               "Failed to read port information for port id %s",
	INFO_GET_PORT_END:                 "Successfully read port information for port id %s",
	INFO_GET_PORT_AUTH_SETTINGS_BEGIN: "Reading port auth settings information for port id %s",
	ERR_GET_PORT_AUTH_SETTINGS_FAILED: "Failed to read port auth settings information for port id %s",
	INFO_GET_PORT_AUTH_SETTINGS_END:   "Successfully read port auth settings information for port id %s",

	// STORAGE NODES
	INFO_GET_ALL_STORAGE_NODES_BEGIN: "Reading all storage nodes information.",
	ERR_GET_ALL_STORAGE_NODES_FAILED: "Failed to read all storage nodes information.",
	INFO_GET_ALL_STORAGE_NODES_END:   "Successfully read storage nodes information.",

	INFO_GET_NODE_BEGIN:               "Reading node information for node id %s",
	ERR_GET_NODE_FAILED:               "Failed to read node information for node id %s",
	INFO_GET_NODE_END:                 "Successfully read node information for node id %s",

	//CHAP USERS
	INFO_GET_ALL_CHAPUSERS_BEGIN: "Reading all chap users information.",
	ERR_GET_ALL_CHAPUSERS_FAILED: "Failed to read all chp users information.",
	INFO_GET_ALL_CHAPUSERS_END:   "Successfully read all chap users information.",
	INFO_GET_CHAP_USER_BEGIN:     "Reading chap user information for target chap user name %s.",
	ERR_GET_CHAP_USER_FAILED:     "Failed to read chap user information for target chap user name %s.",
	INFO_GET_CHAP_USER_END:       "Successfully read chap user information for target chap user name %s.",
	INFO_CREATE_CHAP_USER_BEGIN:  "Creating chap user resource for target chap user name %s.",
	ERR_CREATE_CHAP_USER_FAILED:  "Failed to create chap user resource for target chap user name %s.",
	INFO_CREATE_CHAP_USER_END:    "Successfully created chap user resource for target chap user name %s.",
	INFO_DELETE_CHAP_USER_BEGIN:  "Deleting chap user for target chap user %s.",
	ERR_DELETE_CHAP_USER_FAILED:  "Failed to delete chap user for target chap user %s.",
	INFO_DELETE_CHAP_USER_END:    "Successfully deleted chap user for target chap user %s.",
	INFO_UPDATE_CHAP_USER_BEGIN:  "Updating chap user for chap user id  %s and target chap user name %s.",
	ERR_UPDATE_CHAP_USER_FAILED:  "Failed to update chap user for chap user id  %s and target chap user name %s.",
	INFO_UPDATE_CHAP_USER_END:    "Successfully updated chap user for chap user id  %s and target chap user name %s.",

	//STORAGE CREDENTIAL
	INFO_CHANGE_USER_PASSWORD_BEGIN: "Changing password for user id %s.",
	ERR_CHANGE_USER_PASSWORD_FAILED: "Failed to change password for user id %s.",
	INFO_CHANGE_USER_PASSWORD_END:   "Successfully changed password for user id %s.",

	//CONFIGURATION FILE
	INFO_RESTORE_CONFIG_BEGIN:   "Restoring configuration file.",
	ERR_RESTORE_CONFIG_FAILED:   "Failed to restore configuration file.",
	INFO_RESTORE_CONFIG_END:     "Successfully restored configuration file.",
	INFO_DOWNLOAD_CONFIG_BEGIN:  "Downloading configuration file.",
	ERR_DOWNLOAD_CONFIG_FAILED:  "Failed to download configuration file.",
	INFO_DOWNLOAD_CONFIG_END:    "Successfully downloaded configuration file.",

}
