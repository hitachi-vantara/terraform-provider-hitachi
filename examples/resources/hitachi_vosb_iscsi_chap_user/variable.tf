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

variable "target_chap_user_secret" {
  description = "Secret/password for the iSCSI CHAP user."
  type        = string
  sensitive   = true
}
