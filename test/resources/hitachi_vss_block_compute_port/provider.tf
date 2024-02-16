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
    vss_block_address = "10.76.47.55"
    username          = "admin"
    password          = "vssb-789"
  }

}
