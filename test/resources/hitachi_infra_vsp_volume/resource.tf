resource "hitachi_infra_vsp_volume" "newvol2" {
  name = "VolumeTest1111121"
  capacity = "1GB"
  pool_id = 8
  parity_group_id = "1-5"
  system ="UCP-CI-12035"
  serial = 611039


  # pool_name  = "Terraform_Pool"
  # pool_id = 0
  # paritygroup_id = "group_id"
}

