terraform {
  required_providers {
    hitachi = {
      version = "2.1"
      source  = "localhost/hitachi-vantara/hitachi"
    }
  }
}

provider "hitachi" {
  san_storage_system {
    serial        = 40014
    management_ip = "172.25.47.112"
    #Use secret from secret.tfvars (terraform apply -var-file="secret.tfvars")
    username = var.hitachi_storage_user
    password = var.hitachi_storage_password
  }

  # hitachi_vosb_provider {
  #   vosb_address = ""
  #   username          = "YWRtaW4="
  #   password          = "dnNzYi03ODk="
  # }
}

