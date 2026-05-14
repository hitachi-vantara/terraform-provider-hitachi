##############################################################
##############################################################
# Hitachi VSP - Assign S-VOL Workflow
##############################################################
##############################################################

##############################################################
# VARIABLES
##############################################################
variable "snapshot_group_name" {
  type    = string
  default = "SG_FLOATING_WORKFLOW"
}
variable "hdp_pool_id" {
  type    = number
  default = 0
}
variable "vol_size" {
  type    = number
  default = 1
}
variable "is_drs" {
  type    = bool
  default = true
}

##############################################################
# 1. PRIMARY VOLUME (P-VOL)
##############################################################
resource "hitachi_vsp_volume" "pvol" {
  serial  = var.serial_number
  pool_id = var.hdp_pool_id
  size_gb = var.vol_size

  capacity_saving                         = var.is_drs ? "compression_deduplication" : null
  is_data_reduction_shared_volume_enabled = var.is_drs
}

##############################################################
# 2. FLOATING SNAPSHOT
# Created without S-VOL LDEV ID.
##############################################################
resource "hitachi_vsp_snapshot" "floating_snap" {
  serial              = var.serial_number
  state               = "create"
  snapshot_group_name = var.snapshot_group_name
  snapshot_pool_id    = hitachi_vsp_volume.pvol.volume[0].pool_id

  pvol_ldev_id = hitachi_vsp_volume.pvol.volume[0].ldev_id

  # Note: svol_ldev_id is omitted to keep it "floating"

  is_data_reduction_force_copy = var.is_drs
}

##############################################################
# 3. SECONDARY VOLUME (S-VOL)
# Created using the pool ID and size of P-VOL.
##############################################################
resource "hitachi_vsp_volume" "svol" {
  serial  = var.serial_number
  pool_id = var.hdp_pool_id
  size_gb = var.vol_size

  capacity_saving                         = var.is_drs ? "compression_deduplication" : null
  is_data_reduction_shared_volume_enabled = var.is_drs
}

##############################################################
# 4. ASSIGN S-VOL (STATE=CREATE)
# Links the previously created S-VOL to the floating snapshot.
##############################################################
resource "hitachi_vsp_snapshot" "assign_svol" {
  serial              = var.serial_number
  state               = "create" # Triggers the assignment logic
  snapshot_group_name = hitachi_vsp_snapshot.floating_snap.snapshot_group_name
  snapshot_pool_id    = hitachi_vsp_volume.pvol.volume[0].pool_id
  mirror_unit_id      = hitachi_vsp_snapshot.floating_snap.snapshot[0].mirror_unit_id

  pvol_ldev_id = hitachi_vsp_volume.pvol.volume[0].ldev_id
  svol_ldev_id = hitachi_vsp_volume.svol.volume[0].ldev_id

  is_data_reduction_force_copy = var.is_drs

  # Ensures the S-VOL and floating snapshot exist before assignment
  depends_on = [
    hitachi_vsp_volume.svol,
    hitachi_vsp_snapshot.floating_snap
  ]
}

##############################################################
# OUTPUTS
##############################################################

output "pvol_ldev" {
  value = hitachi_vsp_volume.pvol.volume[0].ldev_id
}

output "svol_ldev" {
  value = hitachi_vsp_volume.svol.volume[0].ldev_id
}

output "floating_snaphot" {
  value = hitachi_vsp_snapshot.floating_snap.snapshot
}

output "assigned_snapshot" {
  value = hitachi_vsp_snapshot.assign_svol.snapshot
}
