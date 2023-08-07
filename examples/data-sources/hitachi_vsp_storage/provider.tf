terraform {
  required_providers {
    hitachi = {
      version = "2.0"
      source  = "localhost/hitachi-vantara/hitachi"
    }
  }
}

provider "hitachi" {
  san_storage_system {
    serial        = 40014
    management_ip = ""
    #Use secret from secret.tfvars (terraform apply -var-file="secret.tfvars")
    username = var.hitachi_storage_user
    password = var.hitachi_storage_password
  }

}

