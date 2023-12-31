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
    username          = var.hitachi_storage_user
    password          = var.hitachi_storage_password
  }
}
