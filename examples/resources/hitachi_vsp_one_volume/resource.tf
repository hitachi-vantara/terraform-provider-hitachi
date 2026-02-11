#
# Hitachi VSP One Volume Management (Create / Read / Update / Delete)
#
# This section defines a resource block for creating and managing volumes
# in the Hitachi Virtual Storage Platform One Block Administrator.
#
# The "hitachi_vsp_one_volume" resource supports full lifecycle management (CRUD)
# of volumes — allowing you to provision new volumes, modify existing ones, and
# remove them when no longer needed. It provides detailed control over capacity,
# naming, pool assignment, and data reduction settings, enabling consistent
# automation of storage provisioning through Terraform.
#
# ------------------------------
# About Create Operations
# ------------------------------
# The create operation provisions one or more new volumes in the specified pool.
# When "number_of_volumes" is greater than 1, multiple volumes are created with
# automatically generated nicknames based on the "nickname_param" settings.
#
# During create:
#   - "volume_id" must not be specified.
#   - "compression_acceleration" is ignored.
#   - "number_of_volumes" defaults to 1 if omitted.
#   - "nickname_param" defines the base name and optional numeric suffix.
#   - "capacity_saving" and "is_data_reduction_share_enabled" configure data reduction.
#   - All created volume IDs and their attributes are returned in "volumes_info".
#
# ------------------------------
# About Update Operations
# ------------------------------
# Update operations are designed to modify **only a single volume** at a time.
# This limitation exists to avoid complexity when dealing with multiple volumes.
#
# During an update:
#   - "volume_id" must be specified to identify the target volume, and must be one of the volumes created
#   - "number_of_volumes" must not be set.
#   - The following fields can be updated:
#       • capacity (increase only)
#       • capacity_saving and compression_acceleration (optional)
#       • nickname_param (for nickname changes)
#   - "compression_acceleration" is optional during update and ignored during create.
#   - "is_data_reduction_share_enabled" cannot be updated.
#   - Update modifies one volume but returns information for all volumes in the state.
#
# Note:
#   The Terraform "destroy" command operates on the Terraform state (tfstate),
#   not the .tf file. Therefore, even after performing an update on a single
#   volume, running "terraform destroy" will delete all volumes that exist
#   in the state file.
#
# ------------------------------
# Key configuration parameters
# ------------------------------
#   - serial: Specifies the storage system’s serial number.
#   - pool_id: Identifies the pool where the volume is created.
#   - capacity: Sets the volume size, specified with a unit (M, G, or T).
#   - number_of_volumes: Creates multiple volumes with sequential nicknames.
#   - nickname_param: Defines the volume base name and optional numeric suffix.
#   - capacity_saving: Enables or disables deduplication and compression.
#   - is_data_reduction_share_enabled: Enables shared data reduction (if allowed).
#
# The resource outputs information about created or updated volumes, including:
#   - id: Comma-separated list of created volume IDs.
#   - volume_count: Number of volumes successfully created.
#   - volumes_info: Detailed attributes of each created or updated volume.
#
# Adjust the parameters to match your environment and desired configuration.
#

# -------------------------------------
# Example: Create
# -------------------------------------
resource "hitachi_vsp_one_volume" "volume" {
  serial   = var.serial_number
  pool_id  = 0
  capacity = "10G"

  nickname_param {
    base_name        = "data_vol"
    start_number     = 0
    number_of_digits = 3
  }

  # number_of_volumes = 2  ##### if not specified, defaults to 1
  capacity_saving                 = "COMPRESSION"
  is_data_reduction_share_enabled = true
}

output "volume_info" {
  value = hitachi_vsp_one_volume.volume.volumes_info
}

output "volume_count" {
  value = hitachi_vsp_one_volume.volume.volume_count
}


# -------------------------------------
# Example: Update
# -------------------------------------
# Specify "volume_id" (required) and omit "number_of_volumes". Only one volume can be updated at a time.
# You can increase capacity, change capacity_saving/compression_acceleration, or nickname parameters.
resource "hitachi_vsp_one_volume" "volume" {
  serial   = var.serial_number
  pool_id  = 0
  capacity = "15G"

  nickname_param {
    base_name        = "newbasename"
    start_number     = 5
    number_of_digits = 2
  }

  volume_id = 6770 ##### required for update
  # or
  # volume_id_hex = "0X1A6"   ##### optional for update  
  # number_of_volumes = 2  ##### must not be specified for update
  capacity_saving                 = "DEDUPLICATION_AND_COMPRESSION"
  compression_acceleration        = true ##### optional for update, ignored for create
  is_data_reduction_share_enabled = true
}

output "volume_info" {
  value = hitachi_vsp_one_volume.volume.volumes_info
}

output "volume_count" {
  value = hitachi_vsp_one_volume.volume.volume_count
}
