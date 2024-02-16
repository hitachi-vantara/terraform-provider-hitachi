terraform {
  required_providers {
    hitachi = {
      version = "2.5"
      source  = "localhost/hitachi-vantara/hitachi"
    }
  }
}

provider "hitachi" {
  hitachi_vss_block_provider {
    vss_block_address = "172.25.58.151"
    username          = var.hitachi_storage_user
    password          = var.hitachi_storage_password
  }

}
