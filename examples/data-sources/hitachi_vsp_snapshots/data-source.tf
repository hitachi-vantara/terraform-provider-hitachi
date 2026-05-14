#
#
# Hitachi VSP Snapshot Range Data Retrieval
#
# This section defines a data source block to fetch information about a range of
# Thin Image snapshots from a Hitachi Virtual Storage Platform (VSP) using
# HashiCorp Configuration Language (HCL).
#
# The data source block "hitachi_vsp_snapshots" retrieves a list of snapshots
# within a specified LDEV range. This allows you to discover and audit multiple
# snapshot relationships across a span of Primary Volumes.
#
# Provide the storage system serial number and the search boundaries using 
# start_pvol_ldev_id/start_pvol_ldev_id_hex and end_pvol_ldev_id/end_pvol_ldev_id_hex.
#
# Only valid for VSP 5000 series
#

data "hitachi_vsp_snapshots" "snapshot" {
  serial  = 12345
  start_pvol_ldev_id = 5828
  # or
  # start_pvol_ldev_id_hex = "0X16C4"
  end_pvol_ldev_id = 5829
  # or
  # end_pvol_ldev_id_hex = "0X16C5"
}

output "snapshots" {
  value = data.hitachi_vsp_snapshots.snapshot.snapshots
}

output "snapshot_count" {
  value = data.hitachi_vsp_snapshots.snapshot.snapshot_count
}
