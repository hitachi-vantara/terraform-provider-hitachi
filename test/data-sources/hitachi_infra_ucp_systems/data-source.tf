data "hitachi_infra_ucp_systems" "ucp_systems" {
  serial_number = "UCP-CI-12035"
  #name = "ucp-20-35"
}

output "ucp_systems" {
  value = data.hitachi_infra_ucp_systems.ucp_systems
}
 