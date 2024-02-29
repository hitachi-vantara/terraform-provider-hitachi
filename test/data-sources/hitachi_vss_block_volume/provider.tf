terraform {
  required_providers {
    hitachi = {
      version = "2.5"
      source  = "localhost/hitachi-vantara/hitachi"
    }
  }
}

provider "hitachi" {
  # hitachi_vss_block_provider {
  #   vss_block_address = "172.25.45.108"
  #   username          = "ucpa"
  #   password          = "Hitachi1"
  # }
  hitachi_vss_block_provider {
    vss_block_address = "172.25.58.151"
    username          = "admin"
    password          = "vssb-789"
  }

}
