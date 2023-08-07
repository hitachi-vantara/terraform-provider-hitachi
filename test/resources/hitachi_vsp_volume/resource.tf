
/*
resource "hitachi_vsp_volume" "mylun" {
  serial  = 40014
  size_gb = 1
  pool_id = 1
}
*/

/*
hitachi_vsp_volume.mylun: Creating...
hitachi_vsp_volume.mylun: Creation complete after 4s [id=226]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
*/

 resource "hitachi_vsp_volume" "mylun_parity" {
   serial         = 40014
   size_gb        = 1
   paritygroup_id = "1-1"
   #dedup_mode     = "compression"
 }


# resource "hitachi_vsp_volume" "mylun" {
#   serial  = 40014
#   size_gb = 1
#   pool_id = 1
# }

# resource "hitachi_vsp_volume" "mylun" {
#   serial         = 40014
#   size_gb        = 4
#   paritygroup_id = "1-1"
#   dedup_mode     = "compression"
# }

# resource "hitachi_vsp_volume" "mylun" {
#   serial         = 40014
#   ldev_id        = 90
#   size_gb        = 3
#   paritygroup_id = "1-1"
#   dedup_mode     = "compression"
# }

# resource "hitachi_vsp_volume" "mylun" {
#   serial  = 40014
#   ldev_id = 65
#   size_gb = 2
#   pool_id = 1
# }

# resource "hitachi_vsp_volume" "mylun" {
#   serial     = 611032
#   size_gb    = 1
#   name       = "terr1"
#   pool_id    = 9
#   dedup_mode = "compression"
# }

# resource "hitachi_vsp_volume" "mylun" {
#   # ldev_id            = 897
#   serial  = 611032
#   size_gb = 3
#   # name           = "terr5"
#   pool_id = 9
#   # dedup_mode        = "compression"
# }

# terraform destroy -target hitachi_vsp_volume.mylun
# terraform apply -target=hitachi_vsp_volume.mylun
