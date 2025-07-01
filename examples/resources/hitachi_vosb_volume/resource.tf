//
// Hitachi VOS Block Volume Resource
//
// This section defines a Terraform resource block to create a Hitachi VOS Block volume.
// The resource "hitachi_vosb_volume" represents a volume on a Hitachi VSP One SDS Block (VOSB)
// using its block interface and allows you to manage its configuration using Terraform.
//
// Customize the values of the parameters (vosb_address, name, capacity_gb, storage_pool, compute_nodes,
// nick_name) to match your desired volume configuration.
//

resource "hitachi_vosb_volume" "volumecreate" {
  vosb_address  = "10.10.12.13"
  name          = "test-volume-newCol"
  capacity_gb   = 1.9
  storage_pool  = "SP01"
  compute_nodes = []
  nick_name     = "Vss_volume_changesnk"
}

output "volumecreateData" {
  value = resource.hitachi_vosb_volume.volumecreate
}
