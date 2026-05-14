##############################################################
##############################################################
# Hitachi VSP Snapshot - Thin Image Standard Example Workflow
##############################################################
##############################################################

variable "pool_id" {
  type        = number
  default     = 30
  description = "Pool ID for primary and secondary volumes."
}

variable "snapshot_pool_id" {
  type        = number
  default     = 103
  description = "Pool ID for the snapshot (Thin Image pool)."
}

##############################################################
# 1. PRIMARY VOLUME (P-VOL)
# Standard volume with Capacity Saving disabled.
##############################################################

resource "hitachi_vsp_volume" "pvol" {
  serial          = var.serial_number
  pool_id         = var.pool_id
  size_gb         = 1
  capacity_saving = "disabled"
}

##############################################################
# 2. SECONDARY VOLUME (S-VOL)
# Target volume for the snapshot; matches P-VOL size.
##############################################################

resource "hitachi_vsp_volume" "svol" {
  serial          = var.serial_number
  pool_id         = var.pool_id
  size_gb         = hitachi_vsp_volume.pvol.size_gb  # same as pvol
  capacity_saving = "disabled"
}

##############################################################
# 3. THIN IMAGE STANDARD SNAPSHOT
# Establishes the relationship between P-VOL and S-VOL.
##############################################################

resource "hitachi_vsp_snapshot" "snapshot_ti_std" {
  # --- Identity ---
  serial              = var.serial_number
  state               = "create"
  snapshot_group_name = "SG_TI_WORKFLOW"
  snapshot_pool_id    = var.snapshot_pool_id

  # --- Logical Mapping (Linking Volumes) ---
  # References the first element [0] of the computed volume list
  pvol_ldev_id = hitachi_vsp_volume.pvol.volume[0].ldev_id
  svol_ldev_id = hitachi_vsp_volume.svol.volume[0].ldev_id

  # --- TI Standard Rules ---
  # Force Copy is not required since capacity_saving is disabled
  is_data_reduction_force_copy = false
  
  # Standard TI allows cloning and cascading
  # is_clone    = true
  can_cascade = true

  # mu_number = 3 # Uncomment to specify a Mirror Unit

  # check output
  lifecycle {
    # Ensure this doesn't accidentally become a TIA (Redirect-on-Write) pair
    postcondition {
      condition     = self.snapshot[0].is_redirect_on_write == false
      error_message = "Logic Error: Snapshot created as TIA (Redirect-on-Write), but HTI Standard was expected."
    }

    # Ensure the pair reached a healthy state (PAIR or PSUS)
    postcondition {
      condition     = contains(["PAIR", "COPY"], self.snapshot[0].status)
      error_message = "Snapshot Failure: The pair is in an unhealthy status: ${self.snapshot[0].status}."
    }
  }
}

##############################################################
# 4. OUTPUTS
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
  value = hitachi_vsp_snapshot.snapshot_ti_std
}