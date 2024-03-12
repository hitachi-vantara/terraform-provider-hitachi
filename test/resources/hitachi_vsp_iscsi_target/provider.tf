terraform {
  required_providers {
    hitachi = {
      version = "2.5"
      source  = "localhost/hitachi-vantara/hitachi"
    }
  }
}

provider "hitachi" {
  # san_storage_system {
  #   serial        = 30078
  #   management_ip = ""
  #   username      = var.hitachi_storage_user
  #   password      = var.hitachi_storage_password
  # }


   hitachi_infrastructure_gateway_provider {
    address = "172.25.22.81"
    username      = "ucpadmin"
    password      = "Passw0rd!"
  }

}
