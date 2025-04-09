terraform {
  required_providers {
    hitachi = {
      version = "2.0"
      source  = "localhost/hitachi-vantara/hitachi"
    }
  }
}

provider "hitachi" {
  hitachi_vosb_block_provider {
    vosb_block_address = "10.76.47.55"
    username          = var.hitachi_storage_user
    password          = var.hitachi_storage_password
  }
}
