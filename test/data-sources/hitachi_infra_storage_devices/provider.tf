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
  #   address = "172.25.22.81"
  #   username      = "ucpadmin"
  #   password      = "Passw0rd!"
  # }

  hitachi_infrastructure_gateway_provider {
    address = "172.25.58.50"
    username      = "ucpadmin"
    password      = "overrunsurveysroutewarnssent"
  }

}
