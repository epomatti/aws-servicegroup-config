terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.22.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_security_group" "default" {
  name   = "terraform-sg"
  vpc_id = aws_vpc.main.id
}

resource "aws_ssm_parameter" "default" {
  name  = "/terraform/security-group-id"
  type  = "String"
  value = aws_security_group.default.id
}
