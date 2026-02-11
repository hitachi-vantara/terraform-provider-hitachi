terraform {
  required_providers {
    hitachi = {
      version = "~> 2.3.0"
      source  = "localhost/hitachi-vantara/hitachi"
    }
  }
}

provider "hitachi" {
  hitachi_vsp_one_provider {
    serial        = var.serial_number
    management_ip = var.vsp_address
    username      = var.hitachi_storage_user
    password      = var.hitachi_storage_password
  }

}
