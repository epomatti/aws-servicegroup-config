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

  tags = {
    Name = "vpc-terraform"
  }
}

resource "aws_security_group" "default" {
  name   = "terraform-sg"
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "sg-terraform"
  }
}

resource "aws_security_group_rule" "default" {
  type              = "ingress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["10.0.0.0/32", "10.0.0.1/32"]
  security_group_id = aws_security_group.default.id
}

resource "aws_ssm_parameter" "default" {
  name  = "/terraform/security-group-id"
  type  = "String"
  value = aws_security_group.default.id
}
