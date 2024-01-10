resource "hitachi_vsp_volume" "mylun2" {
  serial  = 30595
  size_gb = 1
  # pool_name  = "Terraform_Pool"
  # pool_id = 0
  # paritygroup_id = "group_id"
}