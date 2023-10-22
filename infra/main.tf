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

module "vpc" {
  source = "./modules/vpc"
}

resource "aws_ssm_parameter" "security_group_id" {
  name  = "/terraform/security-group-id"
  type  = "String"
  value = module.vpc.security_group_id
}

module "dynamodb" {
  source                      = "./modules/dyndb"
  deletion_protection_enabled = var.dyndb_deletion_protection_enabled
}

module "lambda" {
  source = "./modules/lambda"
}

module "api_gateway" {
  source               = "./modules/apigw"
  lambda_invoke_arn    = module.lambda.invoke_arn
  lambda_function_name = module.lambda.name
}
