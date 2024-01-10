data "hitachi_vsp_storage" "s40014" {
  serial = 40014
}

output "s40014" {
  value = data.hitachi_vsp_storage.s40014
}

## USE SECRET FILE
# terraform plan -var-file="secret.tfvars"
# terraform apply -var-file="secret.tfvars"
# IF Not provided var file then it will ask in command line

# data "hitachi_vsp_storage" "s611039" {
#   serial = 611039
# }

# data "hitachi_vsp_storage" "s611032" {
#   serial = 611032
# }

# data "hitachi_vsp_storage" "s40014" {
#   serial = 40014
# }

# output "s611039" {
#   value = data.hitachi_vsp_storage.s611039
# }

# output "s611032" {
#   value = data.hitachi_vsp_storage.s611032
# }

# output "s40014" {
#   value = data.hitachi_vsp_storage.s40014
# }



