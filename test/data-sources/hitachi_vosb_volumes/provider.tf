terraform {
  required_providers {
    hitachi = {
      version = "2.0"
      source  = "localhost/hitachi-vantara/hitachi"
    }
  }
}

provider "hitachi" {
  hitachi_vosb_provider {
    vosb_address = ""
    username          = var.hitachi_storage_user
    password          = var.hitachi_storage_password
  }

}
