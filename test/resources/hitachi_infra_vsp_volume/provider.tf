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
    address = "172.25.20.56"
    username      = "ucpadmin"
    password      = "Passw0rd!"
  }

}
