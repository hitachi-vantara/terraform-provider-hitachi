##############################################################
##############################################################
# Hitachi VSP - Unassign S-VOL Workflow
##############################################################
##############################################################

##############################################################
# VARIABLES
##############################################################
variable "snapshot_group_name" {
  type    = string
  default = "SG_FLOATING_WORKFLOW"
}
# Existing P-VOL and S-VOL
variable "pvol_ldev_id" {
  type    = number
  default = 2239
}
variable "svol_ldev_id" {
  type    = number
  default = 2240
}
variable "hdp_pool_id" {
  type    = number
  default = 0
}

##############################################################
# 1. UNASSIGN S-VOL (STATE=create with svol omitted)
# This breaks the link between the P-VOL and the S-VOL.
##############################################################
resource "hitachi_vsp_snapshot" "unassign_svol" {
  serial              = var.serial_number
  state               = "create"
  snapshot_group_name = var.snapshot_group_name
  snapshot_pool_id    = var.hdp_pool_id
  pvol_ldev_id        = var.pvol_ldev_id
  mirror_unit_id      = 3
  # Omit svol_ldev_id to trigger the unassign logic

  # Optional: mu_number if you are using specific mirror units
  # mu_number         = 3 
}

##############################################################
# OUTPUTS
##############################################################

output "a_unassigned_snapshot" {
  description = "Status of the snapshot after unassignment"
  value       = hitachi_vsp_snapshot.unassign_svol.snapshot
}

output "b_released_svol_id" {
  description = "The S-VOL LDEV ID that is unassigned and deleted"
  value       = var.svol_ldev_id
}
