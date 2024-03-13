variable "hitachi_storage_user" {
  type        = string
  description = "Username of the Hitachi storage system."
  sensitive   = true
}

variable "hitachi_storage_password" {
  type        = string
  description = "Password of the Hitachi storage system."
  sensitive   = true
}


variable "hitachi_gateway_user" {
  type        = string
  description = "Username of the Hitachi UAI gateway system."
  sensitive   = true
}

variable "hitachi_gateway_password" {
  type        = string
  description = "Password of the Hitachi UAI gateway system."
  sensitive   = true
}