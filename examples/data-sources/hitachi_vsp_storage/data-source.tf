data "hitachi_vsp_storage" "s12345" {
  serial = 12345
}

output "s12345" {
  value = data.hitachi_vsp_storage.s12345
}




