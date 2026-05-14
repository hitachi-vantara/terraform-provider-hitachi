#
# Hitachi VSP Snapshot Group Data Retrieval
#
# This data source retrieves detailed configuration and status information for a 
# specific Hitachi Thin Image (TI) or Thin Image Advanced (TIA) Snapshot Group.
#
# A Snapshot Group organizes multiple snapshot pairs into a single consistency unit, 
# ensuring that all member volumes are captured at the exact same point-in-time.
#

#####################
data "hitachi_vsp_snapshot_group" "snapshotgroup" {
  serial              = 12345
  snapshot_group_name = "snapshotgroupname"
}

# Returns the detailed attributes of the specified group, including member pair statuses.
output "snapshot_group" {
  value = data.hitachi_vsp_snapshot_group.snapshotgroup.snapshot_group
}

# Returns the number of P-VOL/S-VOL pairs currently belonging to this group.
output "snapshot_count" {
  value = data.hitachi_vsp_snapshot_group.snapshotgroup.snapshot_count
}