package messagecatalog

// GetMessage .
func GetMessage(id interface{}) string {
	str := MessageCatalog[id]
	return str
}

// for localization, one option is to replace this file with localized resource
var MessageCatalog = map[interface{}]string{
	Default1: "%v",
	// STORAGE SYSTEM
	INFO_GET_STORAGE_SYSTEM_BEGIN: "Reading storage information for storage system : %s.",
	INFO_GET_STORAGE_SYSTEM_END:   "Successfully read storage information for storage system : %s.",
	ERR_GET_STORAGE_SYSTEM_FAILED: "Failed to read storage information for storage system : %s.",

	// VOLUME
	INFO_GET_LUN_BEGIN:       "Reading lun information for lun id %d.",
	INFO_GET_LUN_END:         "Successfully read lun information for lun id %d.",
	ERR_GET_LUN_FAILED:       "Failed to read lun information for lun id %d.",
	ERR_DELETE_LUN_FAILED:    "Failed to delete lun with id %d on storage serial %d.",
	INFO_DELETE_LUN_BEGIN:    "Deleting lun with id %d on storage serial %d.",
	INFO_DELETE_LUN_END:      "Successfully deleted lun with id %d on storage serial %d.",
	INFO_GET_LUN_RANGE_BEGIN: "Reading lun information from lun id %d to %d.",
	INFO_GET_LUN_RANGE_END:   "Successfully read lun information from lun id %d to %d.",
	ERR_UPDATE_LUN_FAILED:    "Failed to update lun with id %d on storage serial %d.",
	INFO_UPDATE_LUN_BEGIN:    "Updating lun with id %d on storage serial %d.",
	INFO_UPDATE_LUN_END:      "Successfully updated lun with id %d on storage serial %d.",

	// VOSB - VOLUME
	INFO_GET_ALL_VOLUME_INFO_BEGIN: "Reading all volume information.",
	ERR_GET_ALL_VOLUME_INFO_FAILED: "Failed to read all volume information.",
	INFO_GET_ALL_VOLUME_INFO_END:   "Successfully read all volume information.",
	ERR_DELETE_VOLUME_FAILED_MSG:   `"The specified volume cannot be deleted because it is connected to the compute node.","Disconnect the compute node from the specified volume, and then retry the operation.`,

	//VOSB - COMPUTE
	INFO_GET_ALL_SERVERS_BEGIN:     "Reading all compute node information.",
	ERR_GET_ALL_SERVERS_FAILED:     "Failed to read all compute node information.",
	INFO_GET_ALL_SERVERS_END:       "Successfully read compute node information.",
	INFO_GET_SERVER_BEGIN:          "Reading compute node information for server id %s.",
	ERR_GET_SERVER_FAILED:          "Failed to read compute node information for server id %s.",
	INFO_GET_SERVER_END:            "Successfully read compute node information for server id %s.",
	INFO_DELETE_SERVER_BEGIN:       "Deleting compute node for server %s.",
	ERR_DELETE_SERVER_FAILED:       "Failed to delete compute node server %s.",
	INFO_DELETE_SERVER_END:         "Successfully deleted compute node server %s.",
	INFO_CREATE_COMPUTE_NODE_BEGIN: "Creating compute resource for server name %s.",
	ERR_CREATE_COMPUTE_NODE_FAILED: "Failed to create compute resource for server name %s.",
	INFO_CREATE_COMPUTE_NODE_END:   "Successfully created compute resource for server name %s.",
	INFO_UPDATE_COMPUTE_NODE_BEGIN: "Updating compute node for server name %s.",
	ERR_UPDATE_COMPUTE_NODE_FAILED: "Failed to update compute node for server name %s.",
	INFO_UPDATE_COMPUTE_NODE_END:   "Successfully updated compute node for server name %s.",

	// HOSTGROUP
	INFO_GET_HOSTGROUP_BEGIN:     "Reading hostgroup information for portid : %s and hostgroupnumber: %d.",
	INFO_GET_HOSTGROUP_END:       "Successfully read hostgroup information for portid : %s and hostgroupnumber: %d.",
	ERR_GET_HOSTGROUP_FAILED:     "Failed to read hostgroup information for portid : %s and hostgroupnumber: %d.",
	INFO_GET_ALL_HOSTGROUP_BEGIN: "Reading all hostgroup information for serial %d.",
	INFO_GET_ALL_HOSTGROUP_END:   "Successfully read all hostgroup information for serial %d.",
	ERR_GET_ALL_HOSTGROUP_FAILED: "Failed to read all hostgroup information for serial %d.",
	INFO_DELETE_HOSTGROUP_BEGIN:  "Deleting hostgroup for port id %s and hostgroup number %d.",
	INFO_DELETE_HOSTGROUP_END:    "Successfully deleted hostgroup for port id %s and hostgroup number %d.",
	ERR_DELETE_HOSTGROUP_FAILED:  "Failed to delete hostgroup for port id %s and hostgroup number %d.",
	INFO_CREATE_HOSTGROUP_BEGIN:  "Creating hostgroup for port id %s and hostgroup number %d.",
	INFO_CREATE_HOSTGROUP_END:    "Successfully created hostgroup for port id %s and hostgroup number %d.",
	ERR_CREATE_HOSTGROUP_FAILED:  "Failed to create hostgroup for port id %s and hostgroup number %d.",
	INFO_UPDATE_HOSTGROUP_BEGIN:  "Updating hostgroup for port id %s and hostgroup number %d.",
	INFO_UPDATE_HOSTGROUP_END:    "Successfully updated hostgroup for port id %s and hostgroup number %d.",
	ERR_UPDATE_HOSTGROUP_FAILED:  "Failed to update hostgroup for port id %s and hostgroup number %d.",

	// ISCSI TARGET
	INFO_GET_ISCSITARGET_BEGIN:     "Reading iscsi target information for port id %s and iscsi target number %d.",
	INFO_GET_ISCSITARGET_END:       "Successfully read iscsi target information for port id %s and iscsi target number %d.",
	ERR_GET_ISCSITARGET_FAILED:     "Failed to read iscsi target information for port id %s and iscsi target number %d.",
	INFO_GET_ALL_ISCSITARGET_BEGIN: "Reading all iscsi target information for serial %d.",
	INFO_GET_ALL_ISCSITARGET_END:   "Successfully read all iscsi target information for serial %d.",
	ERR_GET_ALL_ISCSITARGET_FAILED: "Failed to read all iscsi target information for serial %d.",
	INFO_UPDATE_ISCSITARGET_BEGIN:  "Updating iscsi target for port id %s and iscsi target number %d.",
	INFO_UPDATE_ISCSITARGET_END:    "Successfully updated iscsi target for port id %s and iscsi target number %d.",
	ERR_UPDATE_ISCSITARGET_FAILED:  "Failed to update iscsi target for port id %s and iscsi target number %d.",
	INFO_CREATE_ISCSITARGET_BEGIN:  "Creating iscsi target for port id %s and iscsi target number %d.",
	INFO_CREATE_ISCSITARGET_END:    "Successfully created iscsi target for port id %s and iscsi target number %d.",
	ERR_CREATE_ISCSITARGET_FAILED:  "Failed to create iscsi target for port id %s and iscsi target number %d.",
	INFO_DELETE_ISCSITARGET_BEGIN:  "Deleting iscsi target for port id %s and iscsi target number %d.",
	INFO_DELETE_ISCSITARGET_END:    "Successfully deleted iscsi target for port id %s and iscsi target number %d.",
	ERR_DELETE_ISCSITARGET_FAILED:  "Failed to delete iscsi target for port id %s and iscsi target number %d.",

	// ISCSI TARGET CHAP USER
	INFO_GET_ISCSITARGET_CHAPUSER_BEGIN:       "Reading iscsi target chap user information for port id %s, iscsi target number %d, chap user name %s and way of chap user %s.",
	INFO_GET_ISCSITARGET_CHAPUSER_END:         "Successfully added iscsi target chap user information for port id %s, iscsi target number %d, chap user name %s and way of chap user %s.",
	ERR_GET_ISCSITARGET_CHAPUSER_FAILED:       "Failed to add iscsi target chap user information for port id %s, iscsi target number %d, chap user name %s and way of chap user %s.",
	INFO_GET_ISCSITARGET_CHAPUSERS_BEGIN:      "Reading iscsi target chap user information for port id %s and iscsi target number %d.",
	INFO_GET_ISCSITARGET_CHAPUSERS_END:        "Successfully added iscsi target chap user information for port id %s and iscsi target number %d.",
	ERR_GET_ISCSITARGET_CHAPUSERS_FAILED:      "Failed to add iscsi target chap user information for port id %s and iscsi target number %d.",
	INFO_CREATE_ISCSITARGET_CHAPUSER_BEGIN:    "Creating iscsi target chap user information for port id %s, iscsi target number %d, chap user name %s and way of chap user %s.",
	INFO_CREATE_ISCSITARGET_CHAPUSER_END:      "Successfully created iscsi target chap user information for port id %s, iscsi target number %d, chap user name %s and way of chap user %s.",
	ERR_CREATE_ISCSITARGET_CHAPUSER_FAILED:    "Failed to create iscsi target chap user information for port id %s, iscsi target number %d, chap user name %s and way of chap user %s.",
	INFO_SET_ISCSITARGET_CHAPUSERNAME_BEGIN:   "Setting iscsi target chap user name for port id %s, iscsi target number %d, chap user name %s and way of chap user %s.",
	INFO_SET_ISCSITARGET_CHAPUSERNAME_END:     "Successfully set iscsi target chap user name for port id %s, iscsi target number %d, chap user name %s and way of chap user %s.",
	ERR_SET_ISCSITARGET_CHAPUSERNAME_FAILED:   "Failed to set iscsi target chap user name for port id %s, iscsi target number %d, chap user name %s and way of chap user %s.",
	INFO_SET_ISCSITARGET_CHAPUSERSECRET_BEGIN: "Setting iscsi target chap user secret for port id %s, iscsi target number %d, chap user name %s, way of chap user %s and secret %s.",
	INFO_SET_ISCSITARGET_CHAPUSERSECRET_END:   "Successfully set iscsi target chap user secret for port id %s, iscsi target number %d, chap user name %s, way of chap user %s and secret %s.",
	ERR_SET_ISCSITARGET_CHAPUSERSECRET_FAILED: "Failed to set iscsi target chap user secret for port id %s, iscsi target number %d, chap user name %s, way of chap user %s and secret %s.",
	INFO_DELETE_ISCSITARGET_CHAPUSER_BEGIN:    "Setting iscsi target chap user name for port id %s, iscsi target number %d, chap user name %s and way of chap user %s.",
	INFO_DELETE_ISCSITARGET_CHAPUSER_END:      "Successfully deleted iscsi target chap user name for port id %s, iscsi target number %d, chap user name %s and way of chap user %s.",
	ERR_DELETE_ISCSITARGET_CHAPUSER_FAILED:    "Failed to delete iscsi target chap user name for port id %s, iscsi target number %d, chap user name %s and way of chap user %s.",
	INFO_CHANGE_ISCSITARGET_CHAPUSER_BEGIN:    "Changing iscsi target chap user secret for port id %s, iscsi target number %d, chap user name %s, way of chap user %s and secret %s.",
	INFO_CHANGE_ISCSITARGET_CHAPUSER_END:      "Successfully changed iscsi target chap user secret for port id %s, iscsi target number %d, chap user name %s, way of chap user %s and secret %s.",
	ERR_CHANGE_ISCSITARGET_CHAPUSER_FAILED:    "Failed to change iscsi target chap user secret for port id %s, iscsi target number %d, chap user name %s, and way of chap user %s .",

	// STORAGE PORTS
	INFO_GET_STORAGE_PORTS_BEGIN:        "Reading storage ports for storage serial %d.",
	INFO_GET_STORAGE_PORTS_END:          "Successfully read storage ports for storage serial %d.",
	ERR_GET_STORAGE_PORTS_FAILED:        "Failed to read storage ports for storage serial %d.",
	INFO_GET_STORAGE_PORTS_PORTID_BEGIN: "Reading storage portId %s for storage serial %d.",
	INFO_GET_STORAGE_PORTS_PORTID_END:   "Successfully read storage portId %s for storage serial %d.",
	ERR_GET_STORAGE_PORTS_PORTID_FAILED: "Failed to read storage portId %s for storage serial %d.",

	// DYNAMIC POOL
	INFO_GET_DYNAMIC_POOLS_BEGIN:   "Reading dynamic pools for storage serial %d.",
	INFO_GET_DYNAMIC_POOLS_END:     "Successfully read dynamic pools for storage serial %d.",
	ERR_GET_DYNAMIC_POOLS_FAILED:   "Failed to read dynamic pools for storage serial %d.",
	INFO_GET_DYNAMIC_POOL_ID_BEGIN: "Reading dynamic pool information with id %d for storage serial %d.",
	INFO_GET_DYNAMIC_POOL_ID_END:   "Successfully read dynamic pool information with id %d for storage serial %d.",
	ERR_GET_DYNAMIC_POOL_ID_FAILED: "Failed to read dynamic pool information with id %d for storage serial %d.",

	// STORAGE POOLS
	INFO_GET_ALL_STORAGE_POOLS_BEGIN: "Reading all storage pool information.",
	ERR_GET_ALL_STORAGE_POOLS_FAILED: "Failed to read all storage pool information.",
	INFO_GET_ALL_STORAGE_POOLS_END:   "Successfully read storage pool information.",
	INFO_GET_STORAGE_POOL_BEGIN:      "Reading storage pool information for pool names %s.",
	ERR_GET_STORAGE_POOL_FAILED:      "Failed to read storage pool information for pool names %s.",
	INFO_GET_STORAGE_POOL_END:        "Successfully read storage pool information for pool names %s.",

	// VOSB - STORAGE PORTS
	INFO_GET_ALL_STORAGE_PORTS_BEGIN: "Reading all storage ports information.",
	ERR_GET_ALL_STORAGE_PORTS_FAILED: "Failed to read all storage ports information.",
	INFO_GET_ALL_STORAGE_PORTS_END:   "Successfully read storage ports information.",

	INFO_GET_PORT_BEGIN: "Reading port information for port id %s",
	ERR_GET_PORT_FAILED: "Failed to read port information for port id %s",
	INFO_GET_PORT_END:   "Successfully read port information for port id %s",

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

	// PARITY GROUP
	INFO_GET_PARITY_GROUP_BEGIN: "Reading parity groups for storage serial %d.",
	INFO_GET_PARITY_GROUP_END:   "Successfully read parity groups for storage serial %d.",
	ERR_GET_PARITY_GROUP_FAILED: "Failed to read parity groups for storage serial %d.",

	//DASHBOARD
	INFO_GET_DASHBOARD_BEGIN: "Reading dashboard information.",
	ERR_GET_DASHBOARD_FAILED: "Failed to read dashboard information.",
	INFO_GET_DASHBOARD_END:   "Successfully read dashboard information.",

	//STORAGE CREDENTIAL
	INFO_CHANGE_USER_PASSWORD_BEGIN: "Changing password for user id %s.",
	ERR_CHANGE_USER_PASSWORD_FAILED: "Failed to change password for user id %s.",
	INFO_CHANGE_USER_PASSWORD_END:   "Successfully changed password for user id %s.",

}
