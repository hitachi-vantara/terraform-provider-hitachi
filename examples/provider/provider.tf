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
    vss_block_address = ""
    username          = "username"
    password          = "password"
  }
  san_storage_system {
    serial        = 40014
    management_ip = ""
    username      = "username"
    password      = "password"
  }
}
