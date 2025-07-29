//
// Hitachi VSP Volume Resource
//
// This section defines a Terraform resource block to create a Hitachi VSP volume.
// The resource "hitachi_vsp_volume" represents a volume on a Hitachi Virtual Storage
// Platform (VSP) and allows you to manage its configuration using Terraform.
//
// Parameters:
//   - `serial`: The serial number of the target storage system.
//   - `size_gb`: The size of the volume to be created (in GB).
//   - `pool_id`, `pool_name`, or `paritygroup_id`: Specify the storage location.
//      At least one of these must be provided.
//   - `name`: Optional name for the volume.
//   - `ldev_id`: Optional logical device ID.
//

resource "hitachi_vsp_volume" "mylun" {
  serial  = 12345
  size_gb = 1
  pool_id = 1

  //Optional parameters
  name = "hitachi_vsp_volume"
  ldev_id = 0
}

output "voloutput" {
  value = resource.hitachi_vsp_volume.mylun
}
