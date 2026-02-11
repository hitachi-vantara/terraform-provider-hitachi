//
// Hitachi VSP Volume Resource
//
// This section defines a Terraform resource block to create or manage a volume
// on a Hitachi Virtual Storage Platform (VSP). The resource "hitachi_vsp_volume"
// allows you to provision new volumes or manage existing ones using Terraform.
//
// Read-Only (Discovery) Mode:
//   - Triggers when an `ldev_id` or `ldev_id_hex` is provided without capacity 
//     (size_gb, cylinder), placement (pool, parity group), or name parameters.
//   - The provider fetches existing hardware metadata to populate the Terraform 
//     state instead of provisioning a new volume.
//   - To manage or update this volume later, configuration must be updated to 
//     match the discovered state (e.g., adding pool_id, size_gb).
//   - NOTE: Running 'terraform destroy' will physically delete the volume 
//     from the storage system.
//
// Required Parameters (for Create/Update):
//   - `serial`: The serial number of the target storage system.
//   - Capacity (choose one):
//       - Block (open systems) volume: `size_gb` (GB, supports decimals)
//       - Mainframe volume: `cylinder` (capacity in cylinders)
//
// Storage Location (exactly one must be provided):
//   - `pool_id`: ID of the target pool. For snapshot vvol, use -1.
//   - `pool_name`: Name of the pool.
//   - `paritygroup_id`: Parity Group ID.
//   - `external_paritygroup_id`: External parity group ID.
//
// Optional Parameters:
//   - `ldev_id/ldev_id_hex`: Logical Device ID in decimal.
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
//   - `volume_format_type`: Optional LDEV format behavior for the volume.
//                        Example values: "QUICK", "NORMAL". Default: NONE. If set, the
//                        provider may trigger a format action during provisioning
//                        according to this type. This is optional and provider-
//                        specific.
//
// Mainframe Volume Notes:
//   - Mainframe mode is enabled by specifying `cylinder`.
//   - Mainframe-only fields: `emulation_type`, `ssid`, `mp_blade_id`, `clpr_id`,
//     `is_tse_volume`, `is_ese_volume`.
//   - For mainframe volume creation:
//       - `size_gb` is not allowed.
//       - `ldev_id` / `ldev_id_hex` are not supported.
//       - `external_paritygroup_id` is not supported.
//       - Use `emulation_type = "3390-A"` with a pool (dynamic pool)
//         or `emulation_type = "3390-V"` with a parity group.
//   - For mainframe updates, only `is_ese_volume` can be changed; other
//     mainframe-specific fields are immutable.
//
// 

// Example Snapshot Vvol
resource "hitachi_vsp_volume" "myVvol" {
  serial  = 12345
  size_gb = 1
  pool_id = -1

  # optionals
  # name     = "newName1"
  # ldev_id  = 6640
  # # or
  # ldev_id_hex = "0X19"
  # volume_format_type = "NORMAL" # optional: QUICK or NORMAL
    # volume_format_type = "NORMAL" # optional: QUICK or NORMAL (default: NONE)
}

output "vvol" {
  value = hitachi_vsp_volume.myVvol
}


// Example DP Volume
resource "hitachi_vsp_volume" "myDpVol" {
  serial  = 12345
  size_gb = 1
  pool_id = 1 # or 
  # pool_name = "poolName"

  # optionals
  # name     = "newName1"
  # ldev_id  = 6640
  # # or
  # ldev_id_hex = "0X19"

  # capacity_saving                         = "compression" # compression, compression_deduplication or disabled
  # is_data_reduction_shared_volume_enabled = true          # ignored on update
  # is_compression_acceleration_enabled     = true

  # only for update
  # is_alua_enabled             = true
  # data_reduction_process_mode = "inline" # or process
  # volume_format_type = "NORMAL" # optional: QUICK or NORMAL
    # volume_format_type = "NORMAL" # optional: QUICK or NORMAL (default: NONE)
}

output "dpvol" {
  value = hitachi_vsp_volume.myDpVol
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
  # volume_format_type = "NORMAL" # optional: QUICK or NORMAL
    # volume_format_type = "NORMAL" # optional: QUICK or NORMAL (default: NONE)
}

output "pgvol" {
  value = hitachi_vsp_volume.myPgVol
}


// Example Mainframe DP Volume (3390-A / pool)
resource "hitachi_vsp_volume" "myMainframePoolVol" {
  serial         = 12345
  cylinder       = 65525
  emulation_type = "3390-A"
  pool_id        = 0 # or
  # pool_name = "mainframePoolName"

  # mainframe optionals
  name         = "TF-MF-POOL-001"
  is_ese_volume = true
  # is_tse_volume = true
  # ssid         = "00"
  # mp_blade_id  = 0
  # clpr_id      = 0
  # volume_format_type = "NORMAL" # optional: QUICK or NORMAL
    # volume_format_type = "NORMAL" # optional: QUICK or NORMAL (default: NONE)
}

output "mainframe_pool_vol" {
  value = hitachi_vsp_volume.myMainframePoolVol
}


// Example Mainframe PG Volume (3390-V / parity group)
resource "hitachi_vsp_volume" "myMainframePgVol" {
  serial         = 12345
  cylinder       = 65525
  emulation_type = "3390-V"
  paritygroup_id = "1-2"

  # mainframe optionals
  name = "TF-MF-PG-001"
  # is_ese_volume = true
  # volume_format_type = "NORMAL" # optional: QUICK or NORMAL
    # volume_format_type = "NORMAL" # optional: QUICK or NORMAL (default: NONE)
}

output "mainframe_pg_vol" {
  value = hitachi_vsp_volume.myMainframePgVol
}
