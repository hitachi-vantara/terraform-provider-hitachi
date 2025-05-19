# Hitachi VOS Block: Add Drives to Storage Pool
#
# This resource allows you to expand a storage pool on a Hitachi VSP One SDS Block (VOSB)
# system by adding either all offline drives or a specified list of drive IDs.
#
# The `hitachi_vosb_add_drives_to_pool` resource interfaces with the VOSB block API
# and enables drive additions via Terraform using HashiCorp Configuration Language (HCL).
#
# ## Usage
# Configure one of the following options to add drives:
#
# - Set `add_all_offline_drives = true` to add all available offline drives to the pool.
# - Or specify a list of `drive_ids` to add specific drives.
#
# **Important:** You must choose one method — do not set both `add_all_offline_drives` and `drive_ids`.
# If `add_all_offline_drives` is `true`, `drive_ids` must not be provided, and vice versa.
#
# ## Parameters
# - `vosb_address`: The address (IP or hostname) of the VOSB system's REST API.
# - `storage_pool_name`: The name of the storage pool to be expanded.
# - `add_all_offline_drives`: Boolean flag to add all offline drives. Mutually exclusive with `drive_ids`.
# - `drive_ids`: List of drive IDs to add to the pool. Mutually exclusive with `add_all_offline_drives`.

resource "hitachi_vosb_add_drives_to_pool" "pool" {
  vosb_address           = "10.10.12.13"
  storage_pool_name      = "SP01"
  add_all_offline_drives = true
  # drive_ids = [
  #   "0437c9f8-ec5a-4527-900b-300519321f1d",
  #   "cbf7b144-593e-451d-9a49-d62e6b7e1334"
  # ]
}

output "pool_output" {
  value = resource.hitachi_vosb_add_drives_to_pool.pool
}
