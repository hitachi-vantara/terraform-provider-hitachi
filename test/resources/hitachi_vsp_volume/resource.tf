/*
// Hitachi VSP Volume Resource
//
// This section defines a Terraform resource block to create a Hitachi VSP volume.
// The resource "hitachi_vsp_volume" represents a volume on a Hitachi Virtual Storage
// Platform (VSP) and allows you to manage its configuration using Terraform.
//
// Customize the values of the parameters (serial, size_gb, pool_id) to match your
// desired volume configuration.
// it supports both terraform provider named san_storage and hitachi_infrastructure_gateway_provider
// for more information about the parameters please refer terraform documentation in the docs/resources/vsp_volume.md file



// Parameter details between direct connect and gateway provider



*/


# resource "hitachi_vsp_volume" "mylun" {
#   serial  = 40014
#   size_gb = 1.2
#   pool_id = 0
#   # deduplication_compression_mode = "COMPRESSION_DEDUPLICATION"
#   # name = "voltest_direct1vbdd"
#    # system ="Logical-UCP-611"
#   # subscriber_id = "46519299-c43c-4c6e-a680-81dce45a3fcb"
# }


resource "hitachi_vsp_volume" "mylun1" {
  serial  = 40015
  size_gb = 0.3
  pool_id = 0
  # # ldev_id = 2522
  # name = "testVolume1"
  # deduplication_compression_mode = "DISABLED"
   # system ="Logical-UCP-611"
  # subscriber_id = "46519299-c43c-4c6e-a680-81dce45a3fcb"
}


output "volumesData" {
  value = resource.hitachi_vsp_volume.mylun1
}
