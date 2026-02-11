//
// Hitachi VSP Volume Resource
//
// This section defines a Terraform resource block to create or manage a volume
// on a Hitachi Virtual Storage Platform (VSP). The resource "hitachi_vsp_volume"
// allows you to provision new volumes or manage existing ones using Terraform.
//
// Required Parameters:
//   - `serial`: The serial number of the target storage system.
//   - `size_gb`: The size of the volume to be created in GB (supports decimals).
//
// Storage Location (exactly one must be provided):
//   - `pool_id`: ID of the target pool.
//   - `pool_name`: Name of the pool.
//   - `paritygroup_id`: Parity Group ID.
//   - `external_paritygroup_id`: External parity group ID.
//
// Optional Parameters:
//   - `ldev_id`: Logical Device ID in decimal.
//   - `name`: Name of the volume.
//   - `capacity_saving`: Capacity-saving mode: compression,
//                        compression_deduplication, or disabled.
//   - `is_data_reduction_shared_volume_enabled`: Enables TI Advanced data
//                        reduction shared volumes. Requires a pool and
//                        capacity saving enabled. Ignored on update.
//   - `is_compression_acceleration_enabled`: Enables compression accelerator.
//                        If omitted, it is automatically enabled when
//                        capacity saving is active and hardware supports it.
//   - `is_alua_enabled`: Enables ALUA mode. Allowed only during update.
//   - `data_reduction_process_mode`: inline or post_process. Only valid when
//                        capacity saving is enabled; allowed only during update.
//

// Example DP Volume
resource "hitachi_vsp_volume" "myDpVol" {
  serial  = 12345
  size_gb = 1
  pool_id = 1 # or 
  # pool_name = "poolName"

  # optionals
  # name     = "newName1"
  # ldev_id  = 6640

  # capacity_saving                         = "compression" # compression, compression_deduplication or disabled
  # is_data_reduction_shared_volume_enabled = true          # ignored on update
  # is_compression_acceleration_enabled     = true

  # only for update
  # is_alua_enabled             = true
  # data_reduction_process_mode = "inline" # or process
}

output "voloutput" {
  value = resource.hitachi_vsp_volume.myDpVol
}


// Example PG Volume
resource "hitachi_vsp_volume" "myPgVol" {
  serial         = 12345
  size_gb        = 2.5
  paritygroup_id = "1-2" # or 
  # external_paritygroup_id = "E1-1" 

  # optionals
  # name = "dummyPG"
  # ldev_id = 243
}

output "voloutput" {
  value = resource.hitachi_vsp_volume.myPgVol
}
