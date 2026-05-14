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

variable "target_chap_user_id" {
  type        = string
  description = "ID of the iSCSI CHAP user to retrieve."
}

variable "target_chap_user_name" {
  type        = string
  description = "Name of the iSCSI CHAP user to retrieve."
}
