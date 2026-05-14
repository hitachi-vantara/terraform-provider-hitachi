//
// Hitachi VSP One SDS Block Volume Resource
//
// This section defines a Terraform resource block to create a Hitachi VSP One SDS Block volume.
// The resource "hitachi_vosb_volume" represents a volume on a Hitachi VSP One SDS Block
// using its block interface and allows you to manage its configuration using Terraform.
//
// Customize the values of the parameters (vosb_address, name, capacity_gb, storage_pool, 
// compute_nodes, nick_name, storage_controller_name) to match your desired volume configuration.

////////////////////////////////////////////////////////////////////////////////
// Example: Import (existing volume)
// -----------------------------------------------------------------------------
// Use terraform import when the volume already exists on storage and
// you want to bring it under Terraform management without re-creating it.
// The import reads the current state from storage and writes it to tfstate.
//
// Import ID format: <vosb_address>/<volume_name>
//
// terraform import hitachi_vosb_volume.volumecreate 10.10.12.13/test-volume-newCol

// Minimal skeleton block required for `terraform import`.
// Create this block in your .tf before running the import command.
resource "hitachi_vosb_volume" "imported" {}

output "imported_vosb_volume_id" {
  value = hitachi_vosb_volume.imported.id
}
////////////////////////////////////////////////////////////////////////////////
//
// Notes:
// - storage_pool and capacity_gb are required only for create (not for import/update).
// - capacity_gb can be increased, but shrinking is not supported.
// - To detach all compute nodes, set compute_nodes = [].


resource "hitachi_vosb_volume" "volumecreate" {
  vosb_address  = "10.10.12.13"
  name          = "test-volume-newCol"
  capacity_gb   = 1
  storage_pool  = "SP01"
  compute_nodes = []
  nick_name     = "Vss_volume_changesnk"
  #storage_controller_name = "SC-01"
}

output "volumecreateData" {
  value = resource.hitachi_vosb_volume.volumecreate
}
