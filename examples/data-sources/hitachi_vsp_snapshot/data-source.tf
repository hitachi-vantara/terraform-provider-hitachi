#
# Hitachi VSP Snapshot Data Retrieval
#
# This section defines a data source block to fetch information about specific
# Thin Image snapshots from a Hitachi Virtual Storage Platform (VSP) using
# HashiCorp Configuration Language (HCL).
#
# The data source block "hitachi_vsp_snapshot" retrieves details about snapshots
# associated with the provided parameters. This allows you to access configuration,
# pair status, and relationship properties for the specified snapshots.
#
# Provide the storage system serial number and a valid combination of search patterns.
# Search Patterns (Mutually Exclusive):
#   1. Group Search: [pvol_ldev_id/pvol_ldev_id_hex + snapshot_group_name]
#   2. Pair Search:  [pvol_ldev_id/pvol_ldev_id_hex + mirror_unit_id]
#   3. P-VOL Search: [pvol_ldev_id/pvol_ldev_id_hex only] -> returns all MU numbers
#   4. S-VOL Search: [svol_ldev_id/svol_ldev_id_hex only] -> finds the parent P-VOL
#

#####################
# Example: Group Search
data "hitachi_vsp_snapshot" "snapshot1" {
  serial              = 12345
  snapshot_group_name = "groupname"
  pvol_ldev_id        = 281
  # or
  # pvol_ldev_id_hex = "0X119"
}

output "snapshot1" {
  value = data.hitachi_vsp_snapshot.snapshot1.snapshots
}

output "snapshot1_count" {
  value = data.hitachi_vsp_snapshot.snapshot1.snapshot_count
}

#####################
# Example: Pair Search
data "hitachi_vsp_snapshot" "snapshot2" {
  serial       = 12345
  pvol_ldev_id = 281
  # or
  # pvol_ldev_id_hex = "0X119"
  mirror_unit_id = 3
}

output "snapshot2" {
  value = data.hitachi_vsp_snapshot.snapshot2.snapshots
}

output "snapshot2_count" {
  value = data.hitachi_vsp_snapshot.snapshot2.snapshot_count
}

#####################
# Example: P-VOL Search
data "hitachi_vsp_snapshot" "snapshot3" {
  serial       = 12345
  pvol_ldev_id = 281
  # or
  # pvol_ldev_id_hex = "0X119"
}

output "snapshot3" {
  value = data.hitachi_vsp_snapshot.snapshot3.snapshots
}

output "snapshot3_count" {
  value = data.hitachi_vsp_snapshot.snapshot3.snapshot_count
}

#####################
# Example: S-VOL Search
data "hitachi_vsp_snapshot" "snapshot4" {
  serial       = 12345
  svol_ldev_id = 100
  # or
  # svol_ldev_id_hex = "0X40"
}

output "snapshot4" {
  value = data.hitachi_vsp_snapshot.snapshot4.snapshots
}

output "snapshot4_count" {
  value = data.hitachi_vsp_snapshot.snapshot4.snapshot_count
}
