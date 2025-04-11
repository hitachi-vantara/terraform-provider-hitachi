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

variable "vosb_block_address" {
  description = "VOSB block address"
  type        = string
}

variable "current_password" {
  type      = string
  sensitive = true
}

variable "new_password" {
  type      = string
  sensitive = true
}
