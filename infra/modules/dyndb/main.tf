resource "aws_dynamodb_table" "security_groups" {
  name           = "security-group-requests"
  billing_mode   = "PAY_PER_REQUEST"
  stream_enabled = false
  hash_key       = "id"

  deletion_protection_enabled = var.deletion_protection_enabled

  server_side_encryption {
    enabled = true
  }

  point_in_time_recovery {
    enabled = true
  }

  attribute {
    name = "id"
    type = "S"
  }
}
