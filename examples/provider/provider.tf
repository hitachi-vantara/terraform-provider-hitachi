#
# Hitachi VSP Terraform Provider
#
# The Hitachi Terraform provider enables automation and infrastructure-as-code
# management of Hitachi storage systems, including:
#
# - Hitachi VSP One SDS Block (Baremetal, AWS, GCP, Azure)
# - Hitachi VSP One B20 series
# - Hitachi VSP 5000 series
#
# This provider allows you to retrieve and manage storage configuration such as
# volumes, hostgroups, ports, and replication settings. It supports both software-defined
# and SAN-based storage deployments, giving administrators consistent tooling for provisioning
# and lifecycle operations across different storage types.
#
# The provider supports two connection types:
#
# - `hitachi_vosb_provider`: For VSP One SDS Block systems.
# - `san_storage_system`: For VSP One B20 series and VSP 5000 series.
#
# Example configuration for VSP One SDS Block:
#
# provider "hitachi" {
#   hitachi_vosb_provider {
#     vosb_address = "10.10.12.13"
#     username     = var.hitachi_storage_user
#     password     = var.hitachi_storage_password
#   }
# }
#
# Example configuration for SAN storage systems:
#
# provider "hitachi" {
#   san_storage_system {
#     serial        = 12345
#     management_ip = "10.10.11.12"
#     username      = var.hitachi_storage_user
#     password      = var.hitachi_storage_password
#   }
# }
#

terraform {
  required_providers {
    hitachi = {
      version = "2.1.2"
      source  = "localhost/hitachi-vantara/hitachi"
    }
  }
}

provider "hitachi" {
  hitachi_vosb_provider {
    vosb_address = "10.10.12.13"
    username     = var.hitachi_storage_user
    password     = var.hitachi_storage_password
  }
}
