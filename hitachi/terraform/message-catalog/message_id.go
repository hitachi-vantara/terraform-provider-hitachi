package messagecatalog

// MessageID .
type MessageID uint64

const sidx = 2000

const (
	// Default1 .
	Default1 MessageID = iota + sidx
	// STORAGE SYSTEM
	INFO_GET_STORAGE_SYSTEM_BEGIN
	INFO_GET_STORAGE_SYSTEM_END
	ERR_GET_STORAGE_SYSTEM_FAILED

	// VOLUME
	INFO_GET_LUN_BEGIN
	INFO_GET_LUN_END
	ERR_GET_LUN_FAILED
	ERR_DELETE_LUN_FAILED
	INFO_DELETE_LUN_BEGIN
	INFO_DELETE_LUN_END
	INFO_GET_LUN_RANGE_BEGIN
	INFO_GET_LUN_RANGE_END
	ERR_UPDATE_LUN_FAILED
	INFO_UPDATE_LUN_BEGIN
	INFO_UPDATE_LUN_END

	// VSSB - VOLUME
	INFO_GET_ALL_VOLUME_INFO_BEGIN
	ERR_GET_ALL_VOLUME_INFO_FAILED
	INFO_GET_ALL_VOLUME_INFO_END
	INFO_CREATE_VOLUME_BEGIN
	INFO_CREATE_VOLUME_END
	ERR_CREATE_VOLUME_FAILED
	INFO_DELETE_VOLUME_BEGIN
	INFO_DELETE_VOLUME_END
	ERR_DELETE_VOLUME_FAILED
	ERR_DELETE_VOLUME_FAILED_MSG

	//VSSB - COMPUTE NODES
	INFO_GET_ALL_SERVERS_BEGIN
	ERR_GET_ALL_SERVERS_FAILED
	INFO_GET_ALL_SERVERS_END
	INFO_GET_SERVER_BEGIN
	ERR_GET_SERVER_FAILED
	INFO_GET_SERVER_END
	INFO_GET_CONNECTION_BY_SERVER_BEGIN
	ERR_GET_CONNECTION_BY_SERVER_FAILED
	INFO_GET_CONNECTION_BY_SERVER_END
	INFO_DELETE_SERVER_BEGIN
	ERR_DELETE_SERVER_FAILED
	INFO_DELETE_SERVER_END
	INFO_CREATE_COMPUTE_NODE_BEGIN
	ERR_CREATE_COMPUTE_NODE_FAILED
	INFO_CREATE_COMPUTE_NODE_END
	INFO_UPDATE_COMPUTE_NODE_BEGIN
	ERR_UPDATE_COMPUTE_NODE_FAILED
	INFO_UPDATE_COMPUTE_NODE_END

	// HOSTGROUP
	INFO_GET_HOSTGROUP_BEGIN
	INFO_GET_HOSTGROUP_END
	ERR_GET_HOSTGROUP_FAILED
	INFO_GET_ALL_HOSTGROUP_BEGIN
	INFO_GET_ALL_HOSTGROUP_END
	ERR_GET_ALL_HOSTGROUP_FAILED
	INFO_DELETE_HOSTGROUP_BEGIN
	INFO_DELETE_HOSTGROUP_END
	ERR_DELETE_HOSTGROUP_FAILED
	INFO_CREATE_HOSTGROUP_BEGIN
	INFO_CREATE_HOSTGROUP_END
	ERR_CREATE_HOSTGROUP_FAILED
	INFO_UPDATE_HOSTGROUP_BEGIN
	INFO_UPDATE_HOSTGROUP_END
	ERR_UPDATE_HOSTGROUP_FAILED

	// ISCSI TARGET
	INFO_GET_ISCSITARGET_BEGIN
	INFO_GET_ISCSITARGET_END
	ERR_GET_ISCSITARGET_FAILED
	INFO_GET_ALL_ISCSITARGET_BEGIN
	INFO_GET_ALL_ISCSITARGET_END
	ERR_GET_ALL_ISCSITARGET_FAILED
	INFO_UPDATE_ISCSITARGET_BEGIN
	INFO_UPDATE_ISCSITARGET_END
	ERR_UPDATE_ISCSITARGET_FAILED
	INFO_CREATE_ISCSITARGET_BEGIN
	INFO_CREATE_ISCSITARGET_END
	ERR_CREATE_ISCSITARGET_FAILED
	INFO_DELETE_ISCSITARGET_BEGIN
	INFO_DELETE_ISCSITARGET_END
	ERR_DELETE_ISCSITARGET_FAILED

	// ISCSI TARGET CHAP USER
	INFO_GET_ISCSITARGET_CHAPUSER_BEGIN
	INFO_GET_ISCSITARGET_CHAPUSER_END
	ERR_GET_ISCSITARGET_CHAPUSER_FAILED
	INFO_GET_ISCSITARGET_CHAPUSERS_BEGIN
	INFO_GET_ISCSITARGET_CHAPUSERS_END
	ERR_GET_ISCSITARGET_CHAPUSERS_FAILED
	INFO_CREATE_ISCSITARGET_CHAPUSER_BEGIN
	INFO_CREATE_ISCSITARGET_CHAPUSER_END
	ERR_CREATE_ISCSITARGET_CHAPUSER_FAILED
	INFO_SET_ISCSITARGET_CHAPUSERNAME_BEGIN
	INFO_SET_ISCSITARGET_CHAPUSERNAME_END
	ERR_SET_ISCSITARGET_CHAPUSERNAME_FAILED
	INFO_SET_ISCSITARGET_CHAPUSERSECRET_BEGIN
	INFO_SET_ISCSITARGET_CHAPUSERSECRET_END
	ERR_SET_ISCSITARGET_CHAPUSERSECRET_FAILED
	INFO_DELETE_ISCSITARGET_CHAPUSER_BEGIN
	INFO_DELETE_ISCSITARGET_CHAPUSER_END
	ERR_DELETE_ISCSITARGET_CHAPUSER_FAILED
	INFO_CHANGE_ISCSITARGET_CHAPUSER_BEGIN
	INFO_CHANGE_ISCSITARGET_CHAPUSER_END
	ERR_CHANGE_ISCSITARGET_CHAPUSER_FAILED

	// STORAGE PORTS
	INFO_GET_STORAGE_PORTS_BEGIN
	INFO_GET_STORAGE_PORTS_END
	ERR_GET_STORAGE_PORTS_FAILED
	INFO_GET_STORAGE_PORTS_PORTID_BEGIN
	INFO_GET_STORAGE_PORTS_PORTID_END
	ERR_GET_STORAGE_PORTS_PORTID_FAILED

	// DYNAMIC POOL
	INFO_GET_DYNAMIC_POOLS_BEGIN
	INFO_GET_DYNAMIC_POOLS_END
	ERR_GET_DYNAMIC_POOLS_FAILED
	INFO_GET_DYNAMIC_POOL_ID_BEGIN
	INFO_GET_DYNAMIC_POOL_ID_END
	ERR_GET_DYNAMIC_POOL_ID_FAILED

	// STORAGE POOLS
	INFO_GET_ALL_STORAGE_POOLS_BEGIN
	ERR_GET_ALL_STORAGE_POOLS_FAILED
	INFO_GET_ALL_STORAGE_POOLS_END
	INFO_GET_STORAGE_POOL_BEGIN
	ERR_GET_STORAGE_POOL_FAILED
	INFO_GET_STORAGE_POOL_END

	// VSSB - STORAGE PORTS
	INFO_GET_ALL_STORAGE_PORTS_BEGIN
	ERR_GET_ALL_STORAGE_PORTS_FAILED
	INFO_GET_ALL_STORAGE_PORTS_END

	INFO_GET_PORT_BEGIN
	ERR_GET_PORT_FAILED
	INFO_GET_PORT_END

	//CHAP USERS
	INFO_GET_ALL_CHAPUSERS_BEGIN
	ERR_GET_ALL_CHAPUSERS_FAILED
	INFO_GET_ALL_CHAPUSERS_END
	INFO_GET_CHAP_USER_BEGIN
	ERR_GET_CHAP_USER_FAILED
	INFO_GET_CHAP_USER_END
	INFO_CREATE_CHAP_USER_BEGIN
	ERR_CREATE_CHAP_USER_FAILED
	INFO_CREATE_CHAP_USER_END
	INFO_DELETE_CHAP_USER_BEGIN
	ERR_DELETE_CHAP_USER_FAILED
	INFO_DELETE_CHAP_USER_END
	INFO_UPDATE_CHAP_USER_BEGIN
	ERR_UPDATE_CHAP_USER_FAILED
	INFO_UPDATE_CHAP_USER_END

	// PARITY GROUP
	INFO_GET_PARITY_GROUP_BEGIN
	INFO_GET_PARITY_GROUP_END
	ERR_GET_PARITY_GROUP_FAILED

	// Dashboard
	INFO_GET_DASHBOARD_BEGIN
	ERR_GET_DASHBOARD_FAILED
	INFO_GET_DASHBOARD_END
)

var enumStrings = map[interface{}]string{
	Default1: "Default1",
	// STORAGE SYSTEM
	INFO_GET_STORAGE_SYSTEM_BEGIN: "INFO_GET_STORAGE_SYSTEM_BEGIN",
	INFO_GET_STORAGE_SYSTEM_END:   "INFO_GET_STORAGE_SYSTEM_END",
	ERR_GET_STORAGE_SYSTEM_FAILED: "ERR_GET_STORAGE_SYSTEM_FAILED",

	// VOLUME
	INFO_GET_LUN_BEGIN:       "INFO_GET_LUN_BEGIN",
	INFO_GET_LUN_END:         "INFO_GET_LUN_END",
	ERR_GET_LUN_FAILED:       "ERR_GET_LUN_FAILED",
	ERR_DELETE_LUN_FAILED:    "ERR_DELETE_LUN_FAILED",
	INFO_DELETE_LUN_BEGIN:    "INFO_DELETE_LUN_BEGIN",
	INFO_DELETE_LUN_END:      "INFO_DELETE_LUN_END",
	INFO_GET_LUN_RANGE_BEGIN: "INFO_GET_LUN_RANGE_BEGIN",
	INFO_GET_LUN_RANGE_END:   "INFO_GET_LUN_RANGE_END",
	ERR_UPDATE_LUN_FAILED:    "ERR_UPDATE_LUN_FAILED",
	INFO_UPDATE_LUN_BEGIN:    "INFO_UPDATE_LUN_BEGIN",
	INFO_UPDATE_LUN_END:      "INFO_UPDATE_LUN_END",

	// VSSB - VOLUME
	INFO_GET_ALL_VOLUME_INFO_BEGIN: "INFO_GET_ALL_VOLUME_INFO_BEGIN",
	ERR_GET_ALL_VOLUME_INFO_FAILED: "ERR_GET_ALL_VOLUME_INFO_FAILED",
	INFO_GET_ALL_VOLUME_INFO_END:   "INFO_GET_ALL_VOLUME_INFO_END",
	INFO_CREATE_VOLUME_BEGIN:       "INFO_CREATE_VOLUME_BEGIN",
	INFO_CREATE_VOLUME_END:         "INFO_CREATE_VOLUME_END",
	ERR_CREATE_VOLUME_FAILED:       "ERR_CREATE_VOLUME_FAILED",
	INFO_DELETE_VOLUME_BEGIN:       "INFO_DELETE_VOLUME_BEGIN",
	INFO_DELETE_VOLUME_END:         "INFO_DELETE_VOLUME_END",
	ERR_DELETE_VOLUME_FAILED:       "ERR_DELETE_VOLUME_FAILED",
	ERR_DELETE_VOLUME_FAILED_MSG:   "ERR_DELETE_VOLUME_FAILED_MSG",

	// COMPUTE NODE - SERVER
	INFO_GET_ALL_SERVERS_BEGIN:          "INFO_GET_ALL_SERVERS_BEGIN",
	ERR_GET_ALL_SERVERS_FAILED:          "ERR_GET_ALL_SERVERS_FAILED",
	INFO_GET_ALL_SERVERS_END:            "INFO_GET_ALL_SERVERS_END",
	INFO_GET_SERVER_BEGIN:               "INFO_GET_SERVER_BEGIN",
	ERR_GET_SERVER_FAILED:               "ERR_GET_SERVER_FAILED",
	INFO_GET_SERVER_END:                 "INFO_GET_SERVER_END",
	INFO_GET_CONNECTION_BY_SERVER_BEGIN: "INFO_GET_CONNECTION_BY_SERVER_BEGIN",
	ERR_GET_CONNECTION_BY_SERVER_FAILED: "ERR_GET_CONNECTION_BY_SERVER_FAILED",
	INFO_GET_CONNECTION_BY_SERVER_END:   "INFO_GET_CONNECTION_BY_SERVER_END",
	INFO_DELETE_SERVER_BEGIN:            "INFO_DELETE_SERVER_BEGIN",
	ERR_DELETE_SERVER_FAILED:            "ERR_DELETE_SERVER_FAILED",
	INFO_DELETE_SERVER_END:              "INFO_DELETE_SERVER_END",
	INFO_CREATE_COMPUTE_NODE_BEGIN:      "INFO_CREATE_COMPUTE_NODE_BEGIN",
	ERR_CREATE_COMPUTE_NODE_FAILED:      "ERR_CREATE_COMPUTE_NODE_FAILED",
	INFO_CREATE_COMPUTE_NODE_END:        "INFO_CREATE_COMPUTE_NODE_END",
	INFO_UPDATE_COMPUTE_NODE_BEGIN:      "INFO_UPDATE_COMPUTE_NODE_BEGIN",
	ERR_UPDATE_COMPUTE_NODE_FAILED:      "ERR_UPDATE_COMPUTE_NODE_FAILED",
	INFO_UPDATE_COMPUTE_NODE_END:        "INFO_UPDATE_COMPUTE_NODE_END",

	// HOSTGROUP
	INFO_GET_HOSTGROUP_BEGIN:     "INFO_GET_HOSTGROUP_BEGIN",
	INFO_GET_HOSTGROUP_END:       "INFO_GET_HOSTGROUP_END",
	ERR_GET_HOSTGROUP_FAILED:     "ERR_GET_HOSTGROUP_FAILED",
	INFO_GET_ALL_HOSTGROUP_BEGIN: "INFO_GET_ALL_HOSTGROUP_BEGIN",
	INFO_GET_ALL_HOSTGROUP_END:   "INFO_GET_ALL_HOSTGROUP_END",
	ERR_GET_ALL_HOSTGROUP_FAILED: "ERR_GET_ALL_HOSTGROUP_FAILED",
	INFO_DELETE_HOSTGROUP_BEGIN:  "INFO_DELETE_HOSTGROUP_BEGIN",
	INFO_DELETE_HOSTGROUP_END:    "INFO_DELETE_HOSTGROUP_END",
	ERR_DELETE_HOSTGROUP_FAILED:  "ERR_DELETE_HOSTGROUP_FAILED",
	INFO_CREATE_HOSTGROUP_BEGIN:  "INFO_CREATE_HOSTGROUP_BEGIN",
	INFO_CREATE_HOSTGROUP_END:    "INFO_CREATE_HOSTGROUP_END",
	ERR_CREATE_HOSTGROUP_FAILED:  "ERR_CREATE_HOSTGROUP_FAILED",
	INFO_UPDATE_HOSTGROUP_BEGIN:  "INFO_UPDATE_HOSTGROUP_BEGIN",
	INFO_UPDATE_HOSTGROUP_END:    "INFO_UPDATE_HOSTGROUP_END",
	ERR_UPDATE_HOSTGROUP_FAILED:  "ERR_UPDATE_HOSTGROUP_FAILED",

	// ISCSI TARGET
	INFO_GET_ISCSITARGET_BEGIN:     "INFO_GET_ISCSITARGET_BEGIN",
	INFO_GET_ISCSITARGET_END:       "INFO_GET_ISCSITARGET_END",
	ERR_GET_ISCSITARGET_FAILED:     "ERR_GET_ISCSITARGET_FAILED",
	INFO_GET_ALL_ISCSITARGET_BEGIN: "INFO_GET_ALL_ISCSITARGET_BEGIN",
	INFO_GET_ALL_ISCSITARGET_END:   "INFO_GET_ALL_ISCSITARGET_END",
	ERR_GET_ALL_ISCSITARGET_FAILED: "ERR_GET_ALL_ISCSITARGET_FAILED",
	INFO_UPDATE_ISCSITARGET_BEGIN:  "INFO_UPDATE_ISCSITARGET_BEGIN",
	INFO_UPDATE_ISCSITARGET_END:    "INFO_UPDATE_ISCSITARGET_END",
	ERR_UPDATE_ISCSITARGET_FAILED:  "ERR_UPDATE_ISCSITARGET_FAILED",
	INFO_CREATE_ISCSITARGET_BEGIN:  "INFO_CREATE_ISCSITARGET_BEGIN",
	INFO_CREATE_ISCSITARGET_END:    "INFO_CREATE_ISCSITARGET_END",
	ERR_CREATE_ISCSITARGET_FAILED:  "ERR_CREATE_ISCSITARGET_FAILED",
	INFO_DELETE_ISCSITARGET_BEGIN:  "INFO_DELETE_ISCSITARGET_BEGIN",
	INFO_DELETE_ISCSITARGET_END:    "INFO_DELETE_ISCSITARGET_END",
	ERR_DELETE_ISCSITARGET_FAILED:  "ERR_DELETE_ISCSITARGET_FAILED",

	// ISCSI TARGET CHAP USERS
	INFO_GET_ISCSITARGET_CHAPUSER_BEGIN:       "INFO_GET_ISCSITARGET_CHAPUSER_BEGIN",
	INFO_GET_ISCSITARGET_CHAPUSER_END:         "INFO_GET_ISCSITARGET_CHAPUSER_END",
	ERR_GET_ISCSITARGET_CHAPUSER_FAILED:       "ERR_GET_ISCSITARGET_CHAPUSER_FAILED",
	INFO_GET_ISCSITARGET_CHAPUSERS_BEGIN:      "INFO_GET_ISCSITARGET_CHAPUSERS_BEGIN",
	INFO_GET_ISCSITARGET_CHAPUSERS_END:        "INFO_GET_ISCSITARGET_CHAPUSERS_END",
	ERR_GET_ISCSITARGET_CHAPUSERS_FAILED:      "ERR_GET_ISCSITARGET_CHAPUSERS_FAILED",
	INFO_CREATE_ISCSITARGET_CHAPUSER_BEGIN:    "INFO_CREATE_ISCSITARGET_CHAPUSER_BEGIN",
	INFO_CREATE_ISCSITARGET_CHAPUSER_END:      "INFO_CREATE_ISCSITARGET_CHAPUSER_END",
	ERR_CREATE_ISCSITARGET_CHAPUSER_FAILED:    "ERR_CREATE_ISCSITARGET_CHAPUSER_FAILED",
	INFO_SET_ISCSITARGET_CHAPUSERNAME_BEGIN:   "INFO_SET_ISCSITARGET_CHAPUSERNAME_BEGIN",
	INFO_SET_ISCSITARGET_CHAPUSERNAME_END:     "INFO_SET_ISCSITARGET_CHAPUSERNAME_END",
	ERR_SET_ISCSITARGET_CHAPUSERNAME_FAILED:   "ERR_SET_ISCSITARGET_CHAPUSERNAME_FAILED",
	INFO_SET_ISCSITARGET_CHAPUSERSECRET_BEGIN: "INFO_SET_ISCSITARGET_CHAPUSERSECRET_BEGIN",
	INFO_SET_ISCSITARGET_CHAPUSERSECRET_END:   "INFO_SET_ISCSITARGET_CHAPUSERSECRET_END",
	ERR_SET_ISCSITARGET_CHAPUSERSECRET_FAILED: "ERR_SET_ISCSITARGET_CHAPUSERSECRET_FAILED",
	INFO_DELETE_ISCSITARGET_CHAPUSER_BEGIN:    "INFO_DELETE_ISCSITARGET_CHAPUSER_BEGIN",
	INFO_DELETE_ISCSITARGET_CHAPUSER_END:      "INFO_DELETE_ISCSITARGET_CHAPUSER_END",
	ERR_DELETE_ISCSITARGET_CHAPUSER_FAILED:    "ERR_DELETE_ISCSITARGET_CHAPUSER_FAILED",
	INFO_CHANGE_ISCSITARGET_CHAPUSER_BEGIN:    "INFO_CHANGE_ISCSITARGET_CHAPUSER_BEGIN",
	INFO_CHANGE_ISCSITARGET_CHAPUSER_END:      "INFO_CHANGE_ISCSITARGET_CHAPUSER_END",
	ERR_CHANGE_ISCSITARGET_CHAPUSER_FAILED:    "ERR_CHANGE_ISCSITARGET_CHAPUSER_FAILED",

	// STORAGE PORTS
	INFO_GET_STORAGE_PORTS_BEGIN:        "INFO_GET_STORAGE_PORTS_BEGIN",
	INFO_GET_STORAGE_PORTS_END:          "INFO_GET_STORAGE_PORTS_END",
	ERR_GET_STORAGE_PORTS_FAILED:        "ERR_GET_STORAGE_PORTS_FAILED",
	INFO_GET_STORAGE_PORTS_PORTID_BEGIN: "INFO_GET_STORAGE_PORTS_PORTID_BEGIN",
	INFO_GET_STORAGE_PORTS_PORTID_END:   "INFO_GET_STORAGE_PORTS_PORTID_END",
	ERR_GET_STORAGE_PORTS_PORTID_FAILED: "ERR_GET_STORAGE_PORTS_PORTID_FAILED",

	// DYNAMIC POOL
	INFO_GET_DYNAMIC_POOLS_BEGIN:   "INFO_GET_DYNAMIC_POOLS_BEGIN",
	INFO_GET_DYNAMIC_POOLS_END:     "INFO_GET_DYNAMIC_POOLS_END",
	ERR_GET_DYNAMIC_POOLS_FAILED:   "ERR_GET_DYNAMIC_POOLS_FAILED",
	INFO_GET_DYNAMIC_POOL_ID_BEGIN: "INFO_GET_DYNAMIC_POOL_ID_BEGIN",
	INFO_GET_DYNAMIC_POOL_ID_END:   "INFO_GET_DYNAMIC_POOL_ID_END",
	ERR_GET_DYNAMIC_POOL_ID_FAILED: "ERR_GET_DYNAMIC_POOL_ID_FAILED",

	// STORAGE POOLS
	INFO_GET_ALL_STORAGE_POOLS_BEGIN: "INFO_GET_ALL_STORAGE_POOLS_BEGIN",
	ERR_GET_ALL_STORAGE_POOLS_FAILED: "ERR_GET_ALL_STORAGE_POOLS_FAILED",
	INFO_GET_ALL_STORAGE_POOLS_END:   "INFO_GET_ALL_STORAGE_POOLS_END",
	INFO_GET_STORAGE_POOL_BEGIN:      "INFO_GET_STORAGE_POOL_BEGIN",
	ERR_GET_STORAGE_POOL_FAILED:      "ERR_GET_STORAGE_POOL_FAILED",
	INFO_GET_STORAGE_POOL_END:        "INFO_GET_STORAGE_POOL_END",

	// VSSB - STORAGE PORTS
	INFO_GET_ALL_STORAGE_PORTS_BEGIN: "INFO_GET_ALL_STORAGE_PORTS_BEGIN",
	ERR_GET_ALL_STORAGE_PORTS_FAILED: "ERR_GET_ALL_STORAGE_PORTS_FAILED",
	INFO_GET_ALL_STORAGE_PORTS_END:   "INFO_GET_ALL_STORAGE_PORTS_END",

	INFO_GET_PORT_BEGIN: "INFO_GET_PORT_BEGIN",
	ERR_GET_PORT_FAILED: "ERR_GET_PORT_FAILED",
	INFO_GET_PORT_END:   "INFO_GET_PORT_END",

	//CHAP USERS
	INFO_GET_ALL_CHAPUSERS_BEGIN: "INFO_GET_ALL_CHAPUSERS_BEGIN",
	ERR_GET_ALL_CHAPUSERS_FAILED: "ERR_GET_ALL_CHAPUSERS_FAILED",
	INFO_GET_ALL_CHAPUSERS_END:   "INFO_GET_ALL_CHAPUSERS_END",
	INFO_GET_CHAP_USER_BEGIN:     "INFO_GET_CHAP_USER_BEGIN",
	ERR_GET_CHAP_USER_FAILED:     "ERR_GET_CHAP_USER_FAILED",
	INFO_GET_CHAP_USER_END:       "INFO_GET_CHAP_USER_END",
	INFO_CREATE_CHAP_USER_BEGIN:  "INFO_CREATE_CHAP_USER_BEGIN",
	ERR_CREATE_CHAP_USER_FAILED:  "ERR_CREATE_CHAP_USER_FAILED",
	INFO_CREATE_CHAP_USER_END:    "INFO_CREATE_CHAP_USER_END",
	INFO_DELETE_CHAP_USER_BEGIN:  "INFO_DELETE_CHAP_USER_BEGIN",
	ERR_DELETE_CHAP_USER_FAILED:  "ERR_DELETE_CHAP_USER_FAILED",
	INFO_DELETE_CHAP_USER_END:    "INFO_DELETE_CHAP_USER_END",
	INFO_UPDATE_CHAP_USER_BEGIN:  "INFO_UPDATE_CHAP_USER_BEGIN",
	ERR_UPDATE_CHAP_USER_FAILED:  "ERR_UPDATE_CHAP_USER_FAILED",
	INFO_UPDATE_CHAP_USER_END:    "INFO_UPDATE_CHAP_USER_END",

	// PARITY GROUP
	INFO_GET_PARITY_GROUP_BEGIN: "INFO_GET_PARITY_GROUP_BEGIN",
	INFO_GET_PARITY_GROUP_END:   "INFO_GET_PARITY_GROUP_END",
	ERR_GET_PARITY_GROUP_FAILED: "ERR_GET_PARITY_GROUP_FAILED",

	//DASHBOARD
	INFO_GET_DASHBOARD_BEGIN: "INFO_GET_DASHBOARD_BEGIN",
	ERR_GET_DASHBOARD_FAILED: "ERR_GET_DASHBOARD_FAILED",
	INFO_GET_DASHBOARD_END:   "INFO_GET_DASHBOARD_END",
}

func (s MessageID) String() string { return enumStrings[s] }

// GetEnumString .
func GetEnumString(m interface{}) string {
	if m, ok := m.(MessageID); ok {
		return m.String()
	}

	return "UNKNOWN"
}
