variable "aws_region" {
  type    = string
  default = "us-east-2"
}

variable "dyndb_deletion_protection_enabled" {
  type = bool
}
