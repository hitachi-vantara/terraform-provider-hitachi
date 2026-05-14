#
# Hitachi VSP Virtual Clone Parent Volume Discovery
#
# This data source identifies all LDEVs within the storage system that are 
# currently acting as "Parent Volumes" for Virtual Clones (vClones).
#
# Unlike a Snapshot Group, which tracks a defined set of pairs, this data source 
# scans the array to find volumes that have been promoted to vClone status and 
# lists their associated root P-VOLs. This is essential for auditing pool 
# consumption and identifying volumes involved in Redirect-on-Write (ROW) chains.
#

#####################
data "hitachi_vsp_vclone_parent_vols" "vcloneparents" {
  serial = var.serial_number
}

output "vclone_parents" {
  value = data.hitachi_vsp_vclone_parent_vols.vcloneparents.parent_volumes
}

output "vclone_parents_count" {
  value = data.hitachi_vsp_vclone_parent_vols.vcloneparents.parent_count
}