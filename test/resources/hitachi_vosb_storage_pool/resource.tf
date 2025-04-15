// Hitachi VOS Block Storage Pool Resource
//
// This section defines a Terraform resource block for managing storage pool
// on a Hitachi VSP One SDS Block (VOSB) using using HashiCorp Configuration Language (HCL).
//
// The resource "hitachi_vosb_storage_pool" represents the storage pool on a Hitachi VSP One SDS Block
// (VOSB) using its block interface and allows you to manage its configuration
// using Terraform.
//
// Customize the values of the parameters (vosb_address, storage_pool_name, add_offline_drives_to_pool, drive_ids) 
// as needed to match your desired storage pool configuration.
//
// - Set "storage_pool_name" to the name of the storage pool to be expanded.
// - Set "add_all_offline_drives" to true to expand the storage pool by adding all offline drives.
// - Use "drive_ids" to specify specific drives to be added to the storage pool.
//
// **Important**: 
// - You cannot set both "add_all_offline_drives" and "drive_ids" at the same time. 
//   If "add_all_offline_drives" is true, do not provide any "drive_ids", and vice versa. 
//   Set only one of these options, not both.
//
// Parameters:
// - storage_pool_name: the name of the storage pool.
// - add_all_offline_drives: A flag to indicate whether to add all offline drives to the storage pool for expansion.
//   If set to true, no specific drive IDs should be provided.
// - drive_ids: A list of specific drive IDs to be used for the expansion of the storage pool.
//   If provided, the "add_all_offline_drives" flag must be set to false.


resource "hitachi_vosb_storage_pool" "pool" {
  vosb_address = var.vosb_address
  storage_pool_name = "SP01"
  add_all_offline_drives = true
  # drive_ids = [
  #   "0437c9f8-ec5a-4527-900b-300519321f1d",
  #   "cbf7b144-593e-451d-9a49-d62e6b7e1334"
  # ]
}

output "pool_output" {
  value = resource.hitachi_vosb_storage_pool.pool
}
