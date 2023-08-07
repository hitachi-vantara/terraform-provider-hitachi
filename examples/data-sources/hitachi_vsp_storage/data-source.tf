data "hitachi_vsp_storage" "s40014" {
  serial = 40014
}

output "s40014" {
  value = data.hitachi_vsp_storage.s40014
}




