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

	INFO_GET_INFRA_HOSTGROUPS_BEGIN: "Reading hostgroups information for storage system : %s.",
	ERR_GET_INFRA_HOSTGROUPS_FAILED: "Failed to read hostgroups information for storage system : %s.",
	INFO_GET_INFRA_HOSTGROUPS_END:   "Successfully read hostgroups information for storage system : %s.",

	INFO_GET_INFRA_HOSTGROUP_BEGIN: "Reading hostgroup information for storage id : %s port id : %s hostgroup name %s.",
	ERR_GET_INFRA_HOSTGROUP_FAILED: "Failed to read hostgroup information for storage id : %s port id : %s hostgroup name %s.",
	INFO_GET_INFRA_HOSTGROUP_END:   "Successfully read hostgroup information for storage id : %s port id : %s hostgroup name %s.",

	// ISCSI TARGET

	INFO_GET_INFRA_ISCSI_TARGETS_BEGIN: "Reading iscsi targets information for storage system : %s.",
	ERR_GET_INFRA_ISCSI_TARGETS_FAILED: "Failed to read iscsi targets information for storage system : %s.",
	INFO_GET_INFRA_ISCSI_TARGETS_END:   "Successfully read iscsi targets information for storage system : %s.",

	INFO_GET_INFRA_ISCSI_TARGET_BEGIN: "Reading iscsi target information for storage id : %s port id : %s iscsi name %s.",
	ERR_GET_INFRA_ISCSI_TARGET_FAILED: "Failed to read iscsi target information for storage id : %s port id : %s iscsi name %s.",
	INFO_GET_INFRA_ISCSI_TARGET_END:   "Successfully read iscsi target information for storage id : %s port id : %s iscsi name %s.",
}
