data "hitachi_vsp_storage_ports" "storageports" {
  serial  = 30078
  port_id = "CL4-C"
}

output "storageports" {
  value = data.hitachi_vsp_storage_ports.storageports
}
