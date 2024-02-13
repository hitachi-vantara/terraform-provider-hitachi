resource "hitachi_infra_vsp_volume" "newvol2" {
  serial = 40014
  name = "volumneNewName"
  capacity = "100MB"
  pool_id = 4
  parity_group_id = "1-1"
  system ="Logical-UCP-30595"
}

