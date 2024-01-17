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
    serial        = 12345
    management_ip = "10.10.11.12"
    username      = var.hitachi_storage_user
    password      = var.hitachi_storage_password
  }

}
