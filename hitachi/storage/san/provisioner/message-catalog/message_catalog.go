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
	ERR_GET_ALL_LUN_FAILED:   "Failed to get all lun information.",
	ERR_CREATE_LUN_FAILED:    "Failed to create lun on storage serial %d.",
	INFO_CREATE_LUN_BEGIN:    "Creating lun on storage serial %d.",
	INFO_CREATE_LUN_END:      "Successfully created lun with id %d on storage serial %d.",
	ERR_EXPAND_LUN_FAILED:    "Failed to expand lun with id %d on storage serial %d.",
	INFO_EXPAND_LUN_BEGIN:    "Expanding lun with id %d on storage serial %d.",
	INFO_EXPAND_LUN_END:      "Successfully expanded lun with id %d on storage serial %d.",
	ERR_DELETE_LUN_FAILED:    "Failed to delete lun with id %d on storage serial %d.",
	INFO_DELETE_LUN_BEGIN:    "Deleting lun with id %d on storage serial %d.",
	INFO_DELETE_LUN_END:      "Successfully deleted lun with id %d on storage serial %d.",
	INFO_GET_LUN_RANGE_BEGIN: "Reading lun information from lun id %d to %d.",
	INFO_GET_LUN_RANGE_END:   "Successfully read lun information from lun id %d to %d.",
	ERR_UPDATE_LUN_FAILED:    "Failed to update lun with id %d on storage serial %d.",
	INFO_UPDATE_LUN_BEGIN:    "Updating lun with id %d on storage serial %d.",
	INFO_UPDATE_LUN_END:      "Successfully updated lun with id %d on storage serial %d.",

	// HOSTGROUP
	INFO_GET_HOSTGROUP_BEGIN:             "Reading hostgroup information for port id %s and hostgroup number %d.",
	INFO_GET_HOSTGROUP_END:               "Successfully read hostgroup information for port id %s and hostgroup number %d.",
	ERR_GET_HOSTGROUP_FAILED:             "Failed to read hostgroup information for port id %s and hostgroup number %d.",
	INFO_GET_ALL_HOSTGROUP_BEGIN:         "Reading all hostgroup information for serial %d.",
	INFO_GET_ALL_HOSTGROUP_END:           "Successfully read all hostgroup information for serial %d.",
	ERR_GET_ALL_HOSTGROUP_FAILED:         "Failed to read all hostgroup information for serial %d.",
	INFO_CREATE_HOSTGROUP_BEGIN:          "Creating hostgroup for port id %s and hostgroup number %d.",
	INFO_CREATE_HOSTGROUP_END:            "Successfully created hostgroup for port id %s and hostgroup number %d.",
	ERR_CREATE_HOSTGROUP_FAILED:          "Failed to create hostgroup for port id %s and hostgroup number %d.",
	ERR_ADD_LDEV_TO_HOSTGROUP_FAILED:     "Failed to add Ldev to hostgroup for ldev id %d.",
	INFO_ADD_WWN_TO_HOSTGROUP_BEGIN:      "Adding wwn %s to hostgroup.",
	INFO_ADD_WWN_TO_HOSTGROUP_END:        "Successfully added wwn %s to hostgroup.",
	ERR_ADD_WWN_TO_HOSTGROUP_FAILED:      "Failed to add wwn %s to hostgroup.",
	INFO_SET_NICKNAME_HOSTGROUP_BEGIN:    "Starting to add nickname %s to wwn %s.",
	INFO_SET_NICKNAME_HOSTGROUP_END:      "Successfully added nickname %s to wwn %s.",
	ERR_SET_NICKNAME_TO_WWN_FAILED:       "Failed to add nickname %s to wwn %s.",
	INFO_DELETE_HOSTGROUP_BEGIN:          "Deleting hostgroup for port id %s and hostgroup number %d.",
	INFO_DELETE_HOSTGROUP_END:            "Successfully deleted hostgroup for port id %s and hostgroup number %d.",
	ERR_DELETE_HOSTGROUP_FAILED:          "Failed to delete hostgroup for port id %s and hostgroup number %d.",
	INFO_SET_MODE_OPTION_HOSTGROUP_BEGIN: "Set hostgroup mode and options for port id %s and hostgroup number %d.",
	INFO_SET_MODE_OPTION_HOSTGROUP_END:   "Successfully set hostgroup mode and options for port id %s and hostgroup number %d.",
	ERR_SET_MODE_OPTION_HOSTGROUP_FAILED: "Failed to set hostgroup mode and options for port id %s and hostgroup number %d.",
	INFO_DELETE_WWN_BEGIN:                "Deleting wwn %s for port id %s and hostgroup number %d.",
	INFO_DELETE_WWN_END:                  "Successfully deleted wwn %s for port id %s and hostgroup number %d.",
	ERR_DELETE_WWN_FAILED:                "Failed to delete wwn %s for port id %s and hostgroup number %d.",
	INFO_REMOVE_LUN_HOSTGROUP_BEGIN:      "Removing lun %d for port id %s and hostgroup number %d.",
	INFO_REMOVE_LUN_HOSTGROUP_END:        "Successfully removed lun %d for port id %s and hostgroup number %d.",
	ERR_REMOVE_LUN_HOSTGROUP_FAILED:      "Failed to remove lun %d for port id %s and hostgroup number %d.",
	INFO_ADD_LDEV_TO_HOSTGROUP_BEGIN:     "Adding ldev to hostgroup for ldev id %d.",
	INFO_ADD_LDEV_TO_HOSTGROUP_END:       "Successfully added ldev to hostgroup for ldev id %d.",

	// ISCSI TARGET
	INFO_GET_ISCSITARGET_BEGIN:             "Reading iscsi target information for port id %s and iscsi target number %d.",
	INFO_GET_ISCSITARGET_END:               "Successfully read iscsi target information for port id %s and iscsi target number %d.",
	ERR_GET_ISCSITARGET_FAILED:             "Failed to read iscsi target information for port id %s and iscsi target number %d.",
	INFO_GET_ISCSITARGET_BY_PORTID_BEGIN:   "Reading iscsi targets information for port id %s.",
	INFO_GET_ISCSITARGET_BY_PORTID_END:     "Successfully read iscsi targets information for port id %s.",
	ERR_GET_ISCSITARGET_BY_PORTID_FAILED:   "Failed to read iscsi targets information for port id %s.",
	INFO_GET_ALL_ISCSITARGET_BEGIN:         "Reading all iscsi target information for serial %d.",
	INFO_GET_ALL_ISCSITARGET_END:           "Successfully read all iscsi target information for serial %d.",
	ERR_GET_ALL_ISCSITARGET_FAILED:         "Failed to read all iscsi target information for serial %d.",
	INFO_SET_MODE_OPTION_ISCSITARGET_BEGIN: "Set hostgroup mode and options for port id %s and iscsi target number %d.",
	INFO_SET_MODE_OPTION_ISCSITARGET_END:   "Successfully set hostgroup mode and options for port id %s and iscsi target number %d.",
	ERR_SET_MODE_OPTION_ISCSITARGET_FAILED: "Failed to set hostgroup mode and options for port id %s and iscsi target number %d.",
	INFO_CREATE_ISCSITARGET_BEGIN:          "Creating iscsi target for port id %s and iscsi target number %d.",
	INFO_CREATE_ISCSITARGET_END:            "Successfully created iscsi target for port id %s and iscsi target number %d.",
	ERR_CREATE_ISCSITARGET_FAILED:          "Failed to create iscsi target for port id %s and iscsi target number %d.",
	INFO_ADD_LDEV_TO_ISCSITARGET_BEGIN:     "Adding Ldev %d to iscsi target.",
	INFO_ADD_LDEV_TO_ISCSITARGET_END:       "Successfully added Ldev %d to iscsi target.",
	ERR_ADD_LDEV_TO_ISCSITARGET_FAILED:     "Failed to add Ldev %d to iscsi target.",
	INFO_ADD_IQN_NAME_TO_ISCSITARGET_BEGIN: "Adding iscsi name %s to iscsi target.",
	INFO_ADD_IQN_NAME_TO_ISCSITARGET_END:   "Successfully added iscsi namev %s to iscsi target.",
	ERR_ADD_IQN_NAME_TO_ISCSITARGET_FAILED: "Failed to add iscsi name %s to iscsi target.",
	INFO_ADD_NICKNAME_TO_ISCSITARGET_BEGIN: "Adding nickname %s to iscsi name %s.",
	INFO_ADD_NICKNAME_TO_ISCSITARGET_END:   "Successfully added nickname %s to iscsi name %s.",
	ERR_ADD_NICKNAME_TO_ISCSITARGET_FAILED: "Failed to add nickname %s to iscsi name %s.",
	INFO_DELETE_IQN_NAME_BEGIN:             "Deleting iscsi name %s for port id %s and iscsi target number %d.",
	INFO_DELETE_IQN_NAME_END:               "Successfully deleted iscsi name %s for port id %s and iscsi target number %d.",
	ERR_DELETE_IQN_NAME_FAILED:             "Failed to delete iscsi name %s for port id %s and iscsi target number %d.",
	INFO_DELETE_ISCSITARGET_BEGIN:          "Deleting iscsi target for port id %s and iscsi target number %d.",
	INFO_DELETE_ISCSITARGET_END:            "Successfully deleted iscsi target for port id %s and iscsi target number %d.",
	ERR_DELETE_ISCSITARGET_FAILED:          "Failed to delete iscsi target for port id %s and iscsi target number %d.",

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
	ERR_CHANGE_ISCSITARGET_CHAPUSER_FAILED:    "Faile to change iscsi target chap user secret for port id %s, iscsi target number %d, chap user name %s, way of chap user %s and secret %s.",

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

	// PARITY GROUP
	INFO_GET_PARITY_GROUP_BEGIN: "Reading parity groups for storage serial %d.",
	INFO_GET_PARITY_GROUP_END:   "Successfully read parity groups for storage serial %d.",
	ERR_GET_PARITY_GROUP_FAILED: "Failed to read parity groups for storage serial %d.",

	// SNAPSHOT
	INFO_GET_SNAPSHOTS_BEGIN: "Reading snapshots for storage serial %v.",
	INFO_GET_SNAPSHOTS_END:   "Successfully read snapshots for storage serial %v.",
	ERR_GET_SNAPSHOTS_FAILED: "Failed to read snapshots for storage serial %v.",

	INFO_GET_SNAPSHOT_BEGIN: "Reading snapshot %v,%v for storage serial %v.",
	INFO_GET_SNAPSHOT_END:   "Successfully read snapshot %v,%v for storage serial %v.",
	ERR_GET_SNAPSHOT_FAILED: "Failed to read snapshot %v,%v for storage serial %v.",

	INFO_GET_SNAPSHOT_REPLICATIONS_BEGIN: "Reading Thin Image replications for storage serial %v.",
	INFO_GET_SNAPSHOT_REPLICATIONS_END:   "Successfully read Thin Image replications for storage serial %v.",
	ERR_GET_SNAPSHOT_REPLICATIONS_FAILED: "Failed to read Thin Image replications for storage serial %v.",

	INFO_CREATE_SNAPSHOT_BEGIN: "Creating snapshot for storage serial %v.",
	INFO_CREATE_SNAPSHOT_END:   "Successfully created snapshot for storage serial %v.",
	ERR_CREATE_SNAPSHOT_FAILED: "Failed to create snapshot for storage serial %v.",

	INFO_SPLIT_SNAPSHOT_BEGIN: "Splitting snapshot %v,%v for storage serial %v.",
	INFO_SPLIT_SNAPSHOT_END:   "Successfully split snapshot %v,%v for storage serial %v.",
	ERR_SPLIT_SNAPSHOT_FAILED: "Failed to split snapshot %v,%v for storage serial %v.",

	INFO_RESYNC_SNAPSHOT_BEGIN: "Resyncing snapshot %v,%v for storage serial %v.",
	INFO_RESYNC_SNAPSHOT_END:   "Successfully resynced snapshot %v,%v for storage serial %v.",
	ERR_RESYNC_SNAPSHOT_FAILED: "Failed to resync snapshot %v,%v for storage serial %v.",

	INFO_RESTORE_SNAPSHOT_BEGIN: "Restoring snapshot %v,%v for storage serial %v.",
	INFO_RESTORE_SNAPSHOT_END:   "Successfully restored snapshot %v,%v for storage serial %v.",
	ERR_RESTORE_SNAPSHOT_FAILED: "Failed to restore snapshot %v,%v for storage serial %v.",

	INFO_DELETE_SNAPSHOT_BEGIN: "Deleting snapshot %v,%v for storage serial %v.",
	INFO_DELETE_SNAPSHOT_END:   "Successfully deleted snapshot %v,%v for storage serial %v.",
	ERR_DELETE_SNAPSHOT_FAILED: "Failed to delete snapshot %v,%v for storage serial %v.",

	INFO_CLONE_SNAPSHOT_BEGIN: "Cloning snapshot %v,%v for storage serial %v.",
	INFO_CLONE_SNAPSHOT_END:   "Successfully cloned snapshot %v,%v for storage serial %v.",
	ERR_CLONE_SNAPSHOT_FAILED: "Failed to clone snapshot %v,%v for storage serial %v.",

	INFO_ASSIGN_SNAPSHOT_VOLUME_BEGIN: "Assigning volume to snapshot %v,%v for storage serial %v.",
	INFO_ASSIGN_SNAPSHOT_VOLUME_END:   "Successfully assigned volume to snapshot %v,%v for storage serial %v.",
	ERR_ASSIGN_SNAPSHOT_VOLUME_FAILED: "Failed to assign volume to snapshot %v,%v for storage serial %v.",

	INFO_UNASSIGN_SNAPSHOT_VOLUME_BEGIN: "Unassigning volume for snapshot %v,%v for storage serial %v.",
	INFO_UNASSIGN_SNAPSHOT_VOLUME_END:   "Successfully unassigned volume for snapshot %v,%v for storage serial %v.",
	ERR_UNASSIGN_SNAPSHOT_VOLUME_FAILED: "Failed to unassign volume for snapshot %v,%v for storage serial %v.",

	INFO_SNAPSHOT_RETENTION_BEGIN: "Setting retention period for snapshot %v,%v for storage serial %v.",
	INFO_SNAPSHOT_RETENTION_END:   "Successfully set retention period for snapshot %v,%v for storage serial %v.",
	ERR_SNAPSHOT_RETENTION_FAILED: "Failed to set retention period for snapshot %v,%v for storage serial %v.",

	// SNAPSHOT_GROUP
	INFO_GET_SNAPSHOT_GROUPS_BEGIN: "Reading snapshot groups for storage serial %v.",
	INFO_GET_SNAPSHOT_GROUPS_END:   "Successfully read snapshot groups for storage serial %v.",
	ERR_GET_SNAPSHOT_GROUPS_FAILED: "Failed to read snapshot groups for storage serial %v.",

	INFO_GET_SNAPSHOT_GROUP_BEGIN: "Reading snapshot group %v for storage serial %v.",
	INFO_GET_SNAPSHOT_GROUP_END:   "Successfully read snapshot group %v for storage serial %v.",
	ERR_GET_SNAPSHOT_GROUP_FAILED: "Failed to read snapshot group %v for storage serial %v.",

	INFO_SPLIT_SNAPSHOT_GROUP_BEGIN: "Splitting snapshot group %v for storage serial %v.",
	INFO_SPLIT_SNAPSHOT_GROUP_END:   "Successfully split snapshot group %v for storage serial %v. Job ID: %v.",
	ERR_SPLIT_SNAPSHOT_GROUP_FAILED: "Failed to split snapshot group %v for storage serial %v.",

	INFO_RESYNC_SNAPSHOT_GROUP_BEGIN: "Resyncing snapshot group %v for storage serial %v.",
	INFO_RESYNC_SNAPSHOT_GROUP_END:   "Successfully resynced snapshot group %v for storage serial %v. Job ID: %v.",
	ERR_RESYNC_SNAPSHOT_GROUP_FAILED: "Failed to resync snapshot group %v for storage serial %v.",

	INFO_RESTORE_SNAPSHOT_GROUP_BEGIN: "Restoring snapshot group %v for storage serial %v.",
	INFO_RESTORE_SNAPSHOT_GROUP_END:   "Successfully restored snapshot group %v for storage serial %v. Job ID: %v.",
	ERR_RESTORE_SNAPSHOT_GROUP_FAILED: "Failed to restore snapshot group %v for storage serial %v.",

	INFO_DELETE_SNAPSHOT_GROUP_BEGIN: "Deleting snapshot group %v for storage serial %v.",
	INFO_DELETE_SNAPSHOT_GROUP_END:   "Successfully deleted snapshot group %v for storage serial %v. Job ID: %v.",
	ERR_DELETE_SNAPSHOT_GROUP_FAILED: "Failed to delete snapshot group %v for storage serial %v.",

	INFO_CLONE_SNAPSHOT_GROUP_BEGIN: "Cloning snapshot group %v for storage serial %v.",
	INFO_CLONE_SNAPSHOT_GROUP_END:   "Successfully cloned snapshot group %v for storage serial %v. Job ID: %v.",
	ERR_CLONE_SNAPSHOT_GROUP_FAILED: "Failed to clone snapshot group %v for storage serial %v.",

	INFO_SNAPSHOT_GROUP_RETENTION_BEGIN: "Setting retention period for snapshot group %v for storage serial %v.",
	INFO_SNAPSHOT_GROUP_RETENTION_END:   "Successfully set retention period for snapshot group %v for storage serial %v.",
	ERR_SNAPSHOT_GROUP_RETENTION_FAILED: "Failed to set retention period for snapshot group %v for storage serial %v.",

	INFO_DELETE_SNAPSHOT_TREE_BEGIN: "Deleting snapshot tree for root LDEV %v for storage serial %v.",
	INFO_DELETE_SNAPSHOT_TREE_END:   "Successfully deleted snapshot tree for root LDEV %v for storage serial %v. Job ID: %v.",
	ERR_DELETE_SNAPSHOT_TREE_FAILED: "Failed to delete snapshot tree for root LDEV %v for storage serial %v.",

	INFO_DELETE_GARBAGE_DATA_BEGIN: "Executing garbage data %v for root LDEV %v for storage serial %v.",
	INFO_DELETE_GARBAGE_DATA_END:   "Successfully executed garbage data %v for root LDEV %v for storage serial %v. Job ID: %v.",
	ERR_DELETE_GARBAGE_DATA_FAILED: "Failed to execute garbage data %v for root LDEV %v for storage serial %v.",

	// Virtual Clone Parent Volumes
	INFO_GET_VCLONE_PARENTS_BEGIN: "Retrieving virtual clone parent volumes for storage serial %v.",
	INFO_GET_VCLONE_PARENTS_END:   "Successfully retrieved virtual clone parent volumes for storage serial %v.",
	ERR_GET_VCLONE_PARENTS_FAILED: "Failed to retrieve virtual clone parent volumes for storage serial %v.",

	// Snapshot Family
	INFO_GET_SNAPSHOT_FAMILY_BEGIN: "Retrieving snapshot family for LDEV %v for storage serial %v.",
	INFO_GET_SNAPSHOT_FAMILY_END:   "Successfully retrieved snapshot family for LDEV %v for storage serial %v.",
	ERR_GET_SNAPSHOT_FAMILY_FAILED: "Failed to retrieve snapshot family for LDEV %v for storage serial %v.",

	// Snapshot VClone Actions
	INFO_CREATE_VCLONE_BEGIN: "Executing create virtual clone for LDEV %v MU %v for storage serial %v.",
	INFO_CREATE_VCLONE_END:   "Successfully executed create virtual clone for LDEV %v MU %v for storage serial %v. Job ID: %v.",
	ERR_CREATE_VCLONE_FAILED: "Failed to execute create virtual clone for LDEV %v MU %v for storage serial %v.",

	INFO_CONVERT_VCLONE_BEGIN: "Executing convert virtual clone for LDEV %v MU %v for storage serial %v.",
	INFO_CONVERT_VCLONE_END:   "Successfully executed convert virtual clone for LDEV %v MU %v for storage serial %v. Job ID: %v.",
	ERR_CONVERT_VCLONE_FAILED: "Failed to execute convert virtual clone for LDEV %v MU %v for storage serial %v.",

	INFO_RESTORE_VCLONE_BEGIN: "Executing restore from virtual clone for LDEV %v MU %v for storage serial %v.",
	INFO_RESTORE_VCLONE_END:   "Successfully executed restore from virtual clone for LDEV %v MU %v for storage serial %v. Job ID: %v.",
	ERR_RESTORE_VCLONE_FAILED: "Failed to execute restore from virtual clone for LDEV %v MU %v for storage serial %v.",

	// Snapshot Group VClone Actions
	INFO_CREATE_GROUP_VCLONE_BEGIN: "Executing create virtual clone for snapshot group %v for storage serial %v.",
	INFO_CREATE_GROUP_VCLONE_END:   "Successfully executed create virtual clone for snapshot group %v for storage serial %v. Job ID: %v.",
	ERR_CREATE_GROUP_VCLONE_FAILED: "Failed to execute create virtual clone for snapshot group %v for storage serial %v.",

	INFO_CONVERT_GROUP_VCLONE_BEGIN: "Executing convert virtual clone for snapshot group %v for storage serial %v.",
	INFO_CONVERT_GROUP_VCLONE_END:   "Successfully executed convert virtual clone for snapshot group %v for storage serial %v. Job ID: %v.",
	ERR_CONVERT_GROUP_VCLONE_FAILED: "Failed to execute convert virtual clone for snapshot group %v for storage serial %v.",

	INFO_RESTORE_GROUP_VCLONE_BEGIN: "Executing restore from virtual clone for snapshot group %v for storage serial %v.",
	INFO_RESTORE_GROUP_VCLONE_END:   "Successfully executed restore from virtual clone for snapshot group %v for storage serial %v. Job ID: %v.",
	ERR_RESTORE_GROUP_VCLONE_FAILED: "Failed to execute restore from virtual clone for snapshot group %v for storage serial %v.",

	// NVM Subsystem
	INFO_GET_NVM_SUBSYSTEMS_BEGIN: "Reading NVM subsystems for storage serial %v.",
	INFO_GET_NVM_SUBSYSTEMS_END:   "Successfully read NVM subsystems for storage serial %v.",
	ERR_GET_NVM_SUBSYSTEMS_FAILED: "Failed to read NVM subsystems for storage serial %v.",

	INFO_GET_NVM_SUBSYSTEM_BEGIN: "Reading NVM subsystem %v for storage serial %v.",
	INFO_GET_NVM_SUBSYSTEM_END:   "Successfully read NVM subsystem %v for storage serial %v.",
	ERR_GET_NVM_SUBSYSTEM_FAILED: "Failed to read NVM subsystem %v for storage serial %v.",

	INFO_CREATE_NVM_SUBSYSTEM_BEGIN: "Creating NVM subsystem for storage serial %v.",
	INFO_CREATE_NVM_SUBSYSTEM_END:   "Successfully created NVM subsystem for storage serial %v.",
	ERR_CREATE_NVM_SUBSYSTEM_FAILED: "Failed to create NVM subsystem for storage serial %v.",

	INFO_UPDATE_NVM_SUBSYSTEM_BEGIN: "Updating NVM subsystem %v for storage serial %v.",
	INFO_UPDATE_NVM_SUBSYSTEM_END:   "Successfully updated NVM subsystem %v for storage serial %v.",
	ERR_UPDATE_NVM_SUBSYSTEM_FAILED: "Failed to update NVM subsystem %v for storage serial %v.",

	INFO_DELETE_NVM_SUBSYSTEM_BEGIN: "Deleting NVM subsystem %v for storage serial %v.",
	INFO_DELETE_NVM_SUBSYSTEM_END:   "Successfully deleted NVM subsystem %v for storage serial %v.",
	ERR_DELETE_NVM_SUBSYSTEM_FAILED: "Failed to delete NVM subsystem %v for storage serial %v.",

	// NVM Subsystem Port
	INFO_GET_NVM_SUBSYSTEM_PORT_BEGIN: "Reading NVM subsystem port %v in subsystem %v for storage serial %v.",
	INFO_GET_NVM_SUBSYSTEM_PORT_END:   "Successfully read NVM subsystem port %v in subsystem %v for storage serial %v.",
	ERR_GET_NVM_SUBSYSTEM_PORT_FAILED: "Failed to read NVM subsystem port %v in subsystem %v for storage serial %v.",

	INFO_GET_NVM_SUBSYSTEM_PORTS_BEGIN: "Reading all NVM subsystem ports for subsystem %v for storage serial %v.",
	INFO_GET_NVM_SUBSYSTEM_PORTS_END:   "Successfully read all NVM subsystem ports for subsystem %v for storage serial %v.",
	ERR_GET_NVM_SUBSYSTEM_PORTS_FAILED: "Failed to read NVM subsystem ports for subsystem %v for storage serial %v.",

	INFO_ADD_NVM_SUBSYSTEM_PORT_BEGIN: "Adding port %v to NVM subsystem %v for storage serial %v.",
	INFO_ADD_NVM_SUBSYSTEM_PORT_END:   "Successfully added port %v to NVM subsystem %v for storage serial %v.",
	ERR_ADD_NVM_SUBSYSTEM_PORT_FAILED: "Failed to add port %v to NVM subsystem %v for storage serial %v.",

	INFO_DELETE_NVM_SUBSYSTEM_PORT_BEGIN: "Deleting port %v from NVM subsystem %v for storage serial %v.",
	INFO_DELETE_NVM_SUBSYSTEM_PORT_END:   "Successfully deleted port %v from NVM subsystem %v for storage serial %v.",
	ERR_DELETE_NVM_SUBSYSTEM_PORT_FAILED: "Failed to delete port %v from NVM subsystem %v for storage serial %v.",

	// Host NQN
	INFO_GET_ALL_HOST_NQNS_BEGIN: "Reading all Host NQNs for NVM subsystem %v for storage serial %v.",
	INFO_GET_ALL_HOST_NQNS_END:   "Successfully read all Host NQNs for NVM subsystem %v for storage serial %v.",
	ERR_GET_ALL_HOST_NQNS_FAILED: "Failed to read Host NQNs for NVM subsystem %v for storage serial %v.",

	INFO_GET_HOST_NQN_BEGIN: "Reading Host NQN %v in NVM subsystem %v for storage serial %v.",
	INFO_GET_HOST_NQN_END:   "Successfully read Host NQN %v in NVM subsystem %v for storage serial %v.",
	ERR_GET_HOST_NQN_FAILED: "Failed to read Host NQN %v in NVM subsystem %v for storage serial %v.",

	INFO_REGISTER_HOST_NQN_BEGIN: "Registering Host NQN %v in NVM subsystem %v for storage serial %v.",
	INFO_REGISTER_HOST_NQN_END:   "Successfully registered Host NQN %v in NVM subsystem %v for storage serial %v.",
	ERR_REGISTER_HOST_NQN_FAILED: "Failed to register Host NQN %v in NVM subsystem %v for storage serial %v.",

	INFO_SET_HOST_NQN_NICKNAME_BEGIN: "Setting nickname for Host NQN %v in subsystem %v for storage serial %v.",
	INFO_SET_HOST_NQN_NICKNAME_END:   "Successfully set nickname for Host NQN %v in subsystem %v for storage serial %v.",
	ERR_SET_HOST_NQN_NICKNAME_FAILED: "Failed to set nickname for Host NQN %v in subsystem %v for storage serial %v.",

	INFO_DELETE_LOGIN_HOST_NQN_BEGIN: "Deleting login for Host NQN at port %v for storage serial %v.",
	INFO_DELETE_LOGIN_HOST_NQN_END:   "Successfully deleted login for Host NQN at port %v for storage serial %v.",
	ERR_DELETE_LOGIN_HOST_NQN_FAILED: "Failed to delete login for Host NQN at port %v for storage serial %v.",

	INFO_DELETE_HOST_NQN_BEGIN: "Deleting Host NQN %v from NVM subsystem %v for storage serial %v.",
	INFO_DELETE_HOST_NQN_END:   "Successfully deleted Host NQN %v from NVM subsystem %v for storage serial %v.",
	ERR_DELETE_HOST_NQN_FAILED: "Failed to delete Host NQN %v from NVM subsystem %v for storage serial %v.",

	// Namespace
	INFO_GET_NAMESPACE_BEGIN: "Reading namespace %v in subsystem %v for storage serial %v.",
	INFO_GET_NAMESPACE_END:   "Successfully read namespace %v in subsystem %v for storage serial %v.",
	ERR_GET_NAMESPACE_FAILED: "Failed to read namespace %v in subsystem %v for storage serial %v.",

	INFO_GET_ALL_NAMESPACES_BEGIN: "Reading all namespaces for subsystem %v for storage serial %v.",
	INFO_GET_ALL_NAMESPACES_END:   "Successfully read all namespaces for subsystem %v for storage serial %v.",
	ERR_GET_ALL_NAMESPACES_FAILED: "Failed to read namespaces for subsystem %v for storage serial %v.",

	INFO_CREATE_NAMESPACE_BEGIN: "Creating namespace in NVM subsystem %v for storage serial %v.",
	INFO_CREATE_NAMESPACE_END:   "Successfully created namespace in NVM subsystem %v for storage serial %v.",
	ERR_CREATE_NAMESPACE_FAILED: "Failed to create namespace in NVM subsystem %v for storage serial %v.",

	INFO_SET_NAMESPACE_NICKNAME_BEGIN: "Setting nickname for namespace %v in subsystem %v for storage serial %v.",
	INFO_SET_NAMESPACE_NICKNAME_END:   "Successfully set nickname for namespace %v in subsystem %v for storage serial %v.",
	ERR_SET_NAMESPACE_NICKNAME_FAILED: "Failed to set nickname for namespace %v in subsystem %v for storage serial %v.",

	INFO_DELETE_NAMESPACE_BEGIN: "Deleting namespace %v from subsystem %v for storage serial %v.",
	INFO_DELETE_NAMESPACE_END:   "Successfully deleted namespace %v from subsystem %v for storage serial %v.",
	ERR_DELETE_NAMESPACE_FAILED: "Failed to delete namespace %v from subsystem %v for storage serial %v.",

	// Namespace Path
	INFO_GET_NAMESPACE_PATHS_BEGIN: "Reading namespace paths for subsystem %v for storage serial %v.",
	INFO_GET_NAMESPACE_PATHS_END:   "Successfully read namespace paths for subsystem %v for storage serial %v.",
	ERR_GET_NAMESPACE_PATHS_FAILED: "Failed to read namespace paths for subsystem %v for storage serial %v.",

	INFO_GET_NAMESPACE_PATH_DETAIL_BEGIN: "Reading path detail for namespace %v, Host NQN %v in subsystem %v for storage serial %v.",
	INFO_GET_NAMESPACE_PATH_DETAIL_END:   "Successfully read path detail for namespace %v, Host NQN %v in subsystem %v for storage serial %v.",
	ERR_GET_NAMESPACE_PATH_DETAIL_FAILED: "Failed to read path detail for namespace %v, Host NQN %v in subsystem %v for storage serial %v.",

	INFO_REGISTER_NAMESPACE_PATH_BEGIN: "Registering path for namespace %v to Host NQN %v for storage serial %v.",
	INFO_REGISTER_NAMESPACE_PATH_END:   "Successfully registered path for namespace %v to Host NQN %v for storage serial %v.",
	ERR_REGISTER_NAMESPACE_PATH_FAILED: "Failed to register path for namespace %v to Host NQN %v for storage serial %v.",

	INFO_DELETE_NAMESPACE_PATH_BEGIN: "Deleting path for namespace %v from Host NQN %v for storage serial %v.",
	INFO_DELETE_NAMESPACE_PATH_END:   "Successfully deleted path for namespace %v from Host NQN %v for storage serial %v.",
	ERR_DELETE_NAMESPACE_PATH_FAILED: "Failed to delete path for namespace %v from Host NQN %v for storage serial %v.",
}
