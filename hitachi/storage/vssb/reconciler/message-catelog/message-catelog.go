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
	INFO_GET_ALL_SERVERS_BEGIN: "Reading all compute node information.",
	ERR_GET_ALL_SERVERS_FAILED: "Failed to read all compute node information.",
	INFO_GET_ALL_SERVERS_END:   "Successfully read compute node information.",
	INFO_GET_SERVER_BEGIN:      "Reading compute node information for server id %s.",
	ERR_GET_SERVER_FAILED:      "Failed to read compute node information for server id %s.",
	INFO_GET_SERVER_END:        "Successfully read compute node information for server id %s.",
	INFO_CREATE_SERVER_BEGIN:   "Creating compute resource for server name %s.",
	ERR_CREATE_SERVER_FAILED:   "Failed to create compute resource for server name %s.",
	INFO_CREATE_SERVER_END:     "Successfully created compute resource for server name %s.",
	INFO_DELETE_SERVER_BEGIN:   "Deleting compute node for server %s.",
	ERR_DELETE_SERVER_FAILED:   "Failed to delete compute node server %s.",
	INFO_DELETE_SERVER_END:     "Successfully deleted compute node server %s.",

	// VOLUME
	INFO_GET_ALL_VOLUME_INFO_BEGIN:             "Reading all volume information.",
	ERR_GET_ALL_VOLUME_INFO_FAILED:             "Failed to read all volume information.",
	INFO_GET_ALL_VOLUME_INFO_END:               "Successfully read all volume information.",
	INFO_CREATE_VOLUME_BEGIN:                   "Creating a volume %s.",
	ERR_CREATE_VOLUME_FAILED:                   "Failed to create a volume %s.",
	INFO_CREATE_VOLUME_END:                     "Successfully created volume %s.",
	INFO_DELETE_VOLUME_BEGIN:                   "Deleting volume %s.",
	ERR_DELETE_VOLUME_FAILED:                   "Failed to delete volume %s.",
	INFO_DELETE_VOLUME_END:                     "Successfully deleted volume %s.",
	INFO_ADD_VOLUME_TO_COMPUTE_NODE_BEGIN:      "Adding volume to compute node %s.",
	ERR_ADD_VOLUME_TO_COMPUTE_NODE_FAILED:      "Failed to add volume to compute node %s.",
	INFO_ADD_VOLUME_TO_COMPUTE_NODE_END:        "Successfully added volume to compute node %s.",
	INFO_REMOVE_VOLUME_FROM_COMPUTE_NODE_BEGIN: "Removing volume from compute node %s.",
	ERR_REMOVE_VOLUME_FROM_COMPUTE_NODE_FAILED: "Failed to remove volume from compute node %s.",
	INFO_REMOVE_VOLUME_FROM_COMPUTE_NODE_END:   "Successfully removed volume from compute node %s.",

	// STORAGE POOLS
	INFO_GET_ALL_STORAGE_POOLS_BEGIN: "Reading all storage pool information.",
	ERR_GET_ALL_STORAGE_POOLS_FAILED: "Failed to read all storage pool information.",
	INFO_GET_ALL_STORAGE_POOLS_END:   "Successfully read storage pool information.",
	INFO_GET_STORAGE_POOL_BEGIN:      "Reading storage pool information for pool names %s.",
	ERR_GET_STORAGE_POOL_FAILED:      "Failed to read storage pool information for pool names %s.",
	INFO_GET_STORAGE_POOL_END:        "Successfully read storage pool information for pool names %s.",

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
}
