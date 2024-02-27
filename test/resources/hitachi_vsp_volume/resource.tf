resource "hitachi_vsp_volume" "mylun2" {
  serial  = 30595
  size_gb = 0.5
  pool_name  = "Terraform_Pool"
  pool_id = 0
  paritygroup_id = "group_id"
}
