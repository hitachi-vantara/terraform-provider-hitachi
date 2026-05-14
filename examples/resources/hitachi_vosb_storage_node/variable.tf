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

variable "setup_user_password" {
  description = "Password used to log into the storage node being added."
  type        = string
  sensitive   = true
}
