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
    address = "172.25.94.145"
    username      = var.hitachi_gateway_user
    password      = var.hitachi_gateway_password
  }

}
