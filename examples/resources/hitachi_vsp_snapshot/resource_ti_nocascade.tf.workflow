##############################################################
##############################################################
# Hitachi VSP Snapshot - Thin Image No Cascade Example Workflow
##############################################################
##############################################################

##############################################################
# VARIABLES
##############################################################
variable "snapshot_group_name" {
  type    = string
  default = "SG_TI_WORKFLOW"
}
variable "hdp_pool_id" {
  type    = number
  default = 0
}
variable "hti_pool_id" {
  type    = number
  default = 3
}
variable "vol_size" {
  type    = number
  default = 1
}

##############################################################
# 1. PRIMARY VOLUME (P-VOL)
# Standard HDP volume.
##############################################################
resource "hitachi_vsp_volume" "pvol" {
  serial  = var.serial_number
  pool_id = var.hdp_pool_id
  size_gb = var.vol_size
}

##############################################################
# 2. SECONDARY VOLUME (S-VOL)
# Created as a Snapshot V-VOL (pool_id = -1).
##############################################################
resource "hitachi_vsp_volume" "svol" {
  serial  = var.serial_number
  pool_id = -1 # Forces Snapshot V-VOL attribute
  size_gb = var.vol_size
}

##############################################################
# 3. DUMMY HOST GROUP (ANCHOR PATHS)
##############################################################
resource "hitachi_vsp_hostgroup" "dummy_path" {
  serial         = var.serial_number
  hostgroup_name = "TI_HG_DUMMY"
  port_id        = "CL1-A"
  host_mode      = "Standard"

  # Map P-VOL
  lun {
    ldev_id = hitachi_vsp_volume.pvol.volume[0].ldev_id
    lun     = 11
  }

  # Map S-VOL
  lun {
    ldev_id = hitachi_vsp_volume.svol.volume[0].ldev_id
    lun     = 12
  }

  wwn {
    host_wwn = "1234567890123456"
  }

  depends_on = [
    hitachi_vsp_volume.pvol,
    hitachi_vsp_volume.svol,
  ]
}

##############################################################
# 4. THIN IMAGE SNAPSHOT (NO CASCADE)
##############################################################
resource "hitachi_vsp_snapshot" "snapshot_no_cascade" {
  serial              = var.serial_number
  state               = "create"
  snapshot_group_name = var.snapshot_group_name
  snapshot_pool_id    = var.hti_pool_id

  pvol_ldev_id = hitachi_vsp_volume.pvol.volume[0].ldev_id
  svol_ldev_id = hitachi_vsp_volume.svol.volume[0].ldev_id

  # Non-Cascade
  can_cascade = false

  depends_on = [
    hitachi_vsp_hostgroup.dummy_path
  ]
}

##############################################################
# OUTPUTS
# These provide immediate visibility into the assigned IDs.
##############################################################

output "pvol_id" {
  description = "Assigned LDEV ID for the Primary Volume"
  value       = hitachi_vsp_volume.pvol.volume[0].ldev_id
}

output "svol_id" {
  description = "Assigned LDEV ID for the Secondary Volume"
  value       = hitachi_vsp_volume.svol.volume[0].ldev_id
}

output "snapshot_info" {
  description = "Summary of the created snapshot pair"
  value       = hitachi_vsp_snapshot.snapshot_no_cascade
}
