#
# Hitachi VSP Snapshot Groups List Retrieval
#
# This data source retrieves a list of Hitachi Thin Image (TI) and Thin Image 
# Advanced (TIA) Snapshot Groups (Consistency Groups) on a specific VSP system.
#
# The 'include_pairs' attribute controls the depth of data retrieval:
# - false (default): Returns snapshot group names and id only (fast).
# - true: Performs a detailed scan of every snapshot pair within each group. 
#   This is required if you need the 'snapshot_count' or detailed 'snapshots' list.
#

#####################
# Example Usage:
data "hitachi_vsp_snapshot_groups" "all_groups" {
  serial        = var.serial_number
  include_pairs = true
}

# Output the list of all snapshot group objects
output "snapshot_groups" {
  value = data.hitachi_vsp_snapshot_groups.all_groups.snapshot_groups
}

# Output the total number of distinct snapshot groups found
output "snapshotgroup_count" {
  value = data.hitachi_vsp_snapshot_groups.all_groups.snapshotgroup_count
}

# Output the cumulative number of snapshot pairs across all retrieved groups
# Note: This will be 0 if 'include_pairs' is set to false.
output "total_snapshot_pairs" {
  value = data.hitachi_vsp_snapshot_groups.all_groups.snapshot_count
}