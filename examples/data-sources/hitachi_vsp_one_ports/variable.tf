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

variable "serial_number" {
  type        = number
  description = "Serial number of the Hitachi storage system."
  default     = 12345
}

variable "management_ip" {
  type        = string
  description = "Management IP address of the Hitachi storage system."
  default     = "10.10.11.12"
}
