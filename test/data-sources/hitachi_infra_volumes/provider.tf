terraform {
  required_providers {
    hitachi = {
      version = "2.5"
      source  = "localhost/hitachi-vantara/hitachi"
    }
  }
}

provider "hitachi" {
  # hitachi_infrastructure_gateway_provider {
  #   address = "172.25.20.56"
  #   # address = "172.25.22.81"
  #   username      = "ucpadmin"
  #   password      = "Passw0rd!"
  # }


   san_storage_system {
    serial        = 611039
    management_ip = "172.25.44.107"
    username      = "maintenance"
    password      ="raid-maintenance"
  }

}


