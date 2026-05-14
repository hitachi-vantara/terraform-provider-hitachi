##############################################################
##############################################################
# Hitachi VSP Snapshot - Thin Image Advanced Example Workflow
##############################################################
##############################################################

##############################################################
# VARIABLES
##############################################################

variable "tia_hdp_pool_id" {
  type        = number
  default     = 17
  description = "The HDP Pool ID where TIA volumes and snapshots will reside"
}

variable "serial_number" {
  type        = number
  default     = 810045
  description = "Serial number of the Hitachi VSP storage system."
}

##############################################################
# 1. PRIMARY VOLUME (P-VOL)
##############################################################

resource "hitachi_vsp_volume" "pvol" {
  serial  = var.serial_number
  pool_id = var.tia_hdp_pool_id
  size_gb = 1

  # TIA Requirements: Capacity Saving and DRS must be enabled
  capacity_saving                         = "compression_deduplication"
  is_data_reduction_shared_volume_enabled = true

  # check output
  lifecycle {
    # Verify the backend actually reports 'DRS' in the attributes list
    postcondition {
      condition     = contains(self.volume[0].attributes, "DRS")
      error_message = "P-VOL hardware verification failed: The 'DRS' attribute is missing from the storage system response."
    }

    # Verify Capacity Saving mode matches TIA requirements
    postcondition {
      condition     = contains(["compression", "compression_deduplication"], self.volume[0].data_reduction_mode)
      error_message = "P-VOL must have compression or deduplication active for TIA."
    }
  }
}

##############################################################
# 2. SECONDARY VOLUME (S-VOL)
##############################################################

resource "hitachi_vsp_volume" "svol" {
  serial  = var.serial_number
  pool_id = var.tia_hdp_pool_id
  size_gb = hitachi_vsp_volume.pvol.size_gb  # same as pvol

  capacity_saving                         = "compression_deduplication"
  is_data_reduction_shared_volume_enabled = true

  # check output
  lifecycle {
    # Check the backend attributes for the 'DRS' string
    postcondition {
      condition     = contains(self.volume[0].attributes, "DRS")
      error_message = "S-VOL hardware verification failed: The 'DRS' attribute is missing from the storage system response."
    }

    # TIA Requirement: S-VOL and P-VOL must share the same Pool ID
    postcondition {
      condition     = self.volume[0].pool_id == hitachi_vsp_volume.pvol.volume[0].pool_id
      error_message = "TIA Compliance Error: S-VOL pool ID must match P-VOL pool ID for Redirect-on-Write."
    }
  }
}

##############################################################
# 3. THIN IMAGE ADVANCED (TIA) SNAPSHOT
##############################################################

resource "hitachi_vsp_snapshot" "snapshot_tia" {
  serial              = var.serial_number
  state              = "create"
  snapshot_group_name = "SG_TIA_WORKFLOW"
  snapshot_pool_id    = var.tia_hdp_pool_id

  pvol_ldev_id = hitachi_vsp_volume.pvol.volume[0].ldev_id
  svol_ldev_id = hitachi_vsp_volume.svol.volume[0].ldev_id

  is_data_reduction_force_copy = true
  can_cascade                  = true
  is_clone                     = false

  # check output
  lifecycle {
    # Ensure a TIA (Redirect-on-Write) pair
    postcondition {
      condition     = self.snapshot[0].is_redirect_on_write == true
      error_message = "Hardware Mismatch: The pair was created, but the storage system did not flag it as a TIA (Redirect-on-Write) pair."
    }
    
    # Ensure the pair reached a healthy state (PAIR or PSUS)
    postcondition {
      condition     = contains(["PAIR", "PSUS"], self.snapshot[0].status)
      error_message = "Snapshot Creation Failure: The pair is in ${self.snapshot[0].status} status."
    }
  }
}

##############################################################
# 4. OUTPUTS
##############################################################

output "tia_resource_summary" {
  value = {
    pvol_ldev = hitachi_vsp_volume.pvol.volume[0].ldev_id
    svol_ldev = hitachi_vsp_volume.svol.volume[0].ldev_id
    pool_id   = var.tia_hdp_pool_id
  }
}

output "snapshot_info" {
  description = "Summary of the created snapshot pair"
  value       = hitachi_vsp_snapshot.snapshot_tia
}
