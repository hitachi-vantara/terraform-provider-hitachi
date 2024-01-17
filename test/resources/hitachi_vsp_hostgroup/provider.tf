terraform {
  required_providers {
    hitachi = {
      version = "2.5"
      source  = "localhost/hitachi-vantara/hitachi"
    }
  }
}

provider "hitachi" {
  san_storage_system {
    serial        = 30595
    management_ip = "172.25.47.112"
    username      = var.hitachi_storage_user
    password      = var.hitachi_storage_password
  }

}
