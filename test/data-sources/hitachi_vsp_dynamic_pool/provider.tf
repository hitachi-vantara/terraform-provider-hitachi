terraform {
  required_providers {
    hitachi = {
      version = "~> 2.1"
      source  = "localhost/hitachi-vantara/hitachi"
    }
  }
}

provider "hitachi" {
  san_storage_system {
    serial        = var.serial_number
    management_ip = var.vsp_address 
    username      = var.hitachi_storage_user
    password      = var.hitachi_storage_password
  }

}
