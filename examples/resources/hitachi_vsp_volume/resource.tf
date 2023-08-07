resource "hitachi_vsp_volume" "mylun" {
  serial  = 40014
  size_gb = 1
  pool_id = 1
}
