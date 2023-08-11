terraform {
  required_providers {
    hitachi = {
      version = "2.0"
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
