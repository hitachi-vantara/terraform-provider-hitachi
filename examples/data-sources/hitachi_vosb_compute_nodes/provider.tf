terraform {
  required_providers {
    hitachi = {
      version = "2.1"
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
