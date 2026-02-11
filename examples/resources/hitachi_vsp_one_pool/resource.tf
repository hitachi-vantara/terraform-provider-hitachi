#
# Hitachi VSP One Pool Management (Create / Read / Update / Delete)
#
# This section defines a resource block for creating and managing storage pools
# in the Hitachi Virtual Storage Platform One Block Administrator.
#
# The "hitachi_vsp_one_pool" resource supports full lifecycle management (CRUD)
# of storage pools — allowing you to provision new pools, modify existing ones, and
# remove them when no longer needed. It provides detailed control over pool
# configuration, drive allocation, encryption settings, and capacity thresholds,
# enabling consistent automation of storage pool provisioning through Terraform.
#
# ------------------------------
# About Create Operations
# ------------------------------
# The create operation provisions a new storage pool in the storage system.
# During create:
#   - "pool_id" must not be specified (it will be auto-assigned).
#   - "name" is required and must be unique within the storage system.
#   - "drives" configuration is required and cannot be changed after creation.
#   - "encryption" setting is optional and cannot be changed after creation.
#   - The created pool ID and all attributes are returned.
#
# ------------------------------
# About Update Operations
# ------------------------------
# The update operation allows modification of pool properties:
#   - "name" can be updated.
#   - "threshold_warning" and "threshold_depletion" can be modified.
#   - Drive configuration cannot be changed after pool creation.
#   - Encryption setting cannot be changed after pool creation.
#
# ------------------------------
# About Delete Operations
# ------------------------------
# The delete operation removes the storage pool from the storage system.
# All associated volumes must be deleted before the pool can be removed.
#



# Create a pool
resource "hitachi_vsp_one_pool" "backup_pool" {
  serial     = 12345
  name       = "new-pool-1"
  encryption = false

  drive_configuration {
    drive_type_code   = "SNB5B-R1R9NC"
    data_drive_count  = 6
    raid_level        = "RAID5"
    parity_group_type = "DDP"
  }
  # For UPDATE operation : Adjusted capacity thresholds for backup pool
  # # More conservative thresholds for backup storage
  # threshold_warning   = 80
  # threshold_depletion = 90
}

output "backup_pool_info" {
  description = "Backup pool details including ID and configuration"
  value = hitachi_vsp_one_pool.backup_pool.data
}
