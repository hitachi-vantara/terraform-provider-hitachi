#
# Hitachi VSP Snapshot Family Data Retrieval
#
# This data source retrieves the snapshot lineage (Family Table) for a specific 
# Logical Device (LDEV). It provides a list of all related volumes, including 
# P-VOLs, S-VOLs, and Virtual Clones (vClones).
#
# The "Family" represents the relationship hierarchy of a volume, which is 
# especially critical for tracking members after a Snapshot Group promotion (vClone) 
# where traditional group metadata is replaced by standalone LDEV attributes.
#

#####################
data "hitachi_vsp_snapshot_family" "vclonefamily" {
  serial  = var.serial_number
  ldev_id = 681
  # or
  # ldev_id_hex = "0x557"
}

output "vclone_family" {
  value = data.hitachi_vsp_snapshot_family.vclonefamily.family_members
}

output "vclone_family_count" {
  value = data.hitachi_vsp_snapshot_family.vclonefamily.total_members
}