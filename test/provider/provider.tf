terraform {
  required_providers {
    hitachi = {
      version = "2.0"
      source  = "localhost/hitachi-vantara/hitachi"
      #source  = "hitachi-vantara/storage-systems"
    }
  }
}

provider "hitachi" {
  hitachi_vosb_provider {
    vosb_address = ""
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
