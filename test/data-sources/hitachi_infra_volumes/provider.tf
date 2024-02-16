terraform {
  required_providers {
    hitachi = {
      version = "2.5"
      source  = "localhost/hitachi-vantara/hitachi"
    }
  }
}

provider "hitachi" {
  hitachi_infrastructure_gateway_provider {
    #address = "172.25.22.64"
    address = "172.25.22.81"
    username      = "ucpadmin"
    password      = "Passw0rd!"
  }

}
