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
    address = "172.25.22.81"
    username      = "ucpadmin"
    password      = "Passw0rd!"
  }
  # hitachi_infrastructure_gateway_provider {
  #   address = "172.25.20.54"
  #   username      = "ucpadmin"
  #   password      = "Passw0rd!"
  # }
  # hitachi_infrastructure_gateway_provider {
  #   address = "172.25.22.81"
  #   username      = "ucpadmin"
  #   password      = "Passw0rd!"
  # }

  #  hitachi_infrastructure_gateway_provider {
  #   address = "172.25.58.50"
  #   username      = "ucpadmin"
  #   password      = "overrunsurveysroutewarnssent"
  

  # san_storage_system {
  #   serial        = 40014
  #   management_ip = "172.25.47.115"
  #   username      = "maintenance"
  #   password      ="raid-maintenance"
  # }
  
# san_storage_system {
#    serial        = 611039
#    management_ip = "172.25.44.107"
#    username      = "maintenance"
#    password      = "raid-maintenance"
# }
}
