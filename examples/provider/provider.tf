/*
** Terraform Configuration for Hitachi Provider and VSS Block / SAN Storage System
**
** This Terraform configuration defines the required provider blocks for interacting with Hitachi
** resources, specifically for VSS block and SAN storage system.
**
** The "required_providers" block specifies the version and source of the Hitachi provider to be used.
** You can customize the version and source accordingly.
**
** The "provider" block configures the Hitachi provider with necessary authentication details for the
** VSS block and SAN storage system.
** Customize the values of "vss_block_address", "management_ip", "username", and "password" to match your
** environment's configuration.
**
*/

terraform {
  required_providers {
    hitachi = {
      version = "2.5"
      source  = "localhost/hitachi-vantara/hitachi"
    }
  }
}

provider "hitachi" {
  hitachi_vss_block_provider {
    vss_block_address = "10.10.12.13"
    username          = "username"
    password          = "password"
  }
  
  san_storage_system {
    serial        = 12345
    management_ip = "10.10.11.12"
    username      = "username"
    password      = "password"
  }
}
