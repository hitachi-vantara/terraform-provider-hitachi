# resource "hitachi_vss_block_volume" "volumecreate" {
#   vss_block_address = ""
#   name              = "test-volume-vss"
#   capacity_gb       = 0.65
#   storage_pool      = "SP01"

# }

# resource "hitachi_vss_block_volume" "volumecreate" {
#   vss_block_address = ""
#   name              = "test-volume-vss"
#   capacity_gb       = 0.85
#   nick_name         = "test-volume-vss-nick"

# }

# resource "hitachi_vss_block_volume" "volumecreate" {
#   vss_block_address = ""
#   name              = "test-volume-vss"
#   capacity_gb       = 0.95
#   storage_pool      = "SP01"
#   # compute_nodes = ["MongoNode1","MongoNode2","MongoNode3"]

# }

resource "hitachi_vss_block_volume" "volumecreate" {
  vss_block_address = "172.25.58.151"
  name              = "test-volume-newCol"
  capacity_gb       = 1.9
  storage_pool      = "SP01"
  compute_nodes     = []
  nick_name         = "Vss_volume_changesnk"


}

output "volumecreateData" {
  value = resource.hitachi_vss_block_volume.volumecreate
}
