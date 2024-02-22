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
    address = "10.10.10.10"
    username      = var.hitachi_gateway_user
    password      = var.hitachi_gateway_password
  }

}
