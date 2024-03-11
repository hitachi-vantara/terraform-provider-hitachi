data "hitachi_infra_systems" "systems" {
  # serial_number = "UCP-CI-12035"
  #name = "ucp-20-35"
}

output "systems" {
  value = data.hitachi_infra_systems.systems
}
 