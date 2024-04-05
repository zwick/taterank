provider "aws" {
  region = var.region
  default_tags {
    tags = {
      Project = var.app_name
    }
  }
}

resource "aws_dynamodb_table" "table" {
  name      = local.dynamo_table_name
  hash_key  = "PK"
  range_key = "SK"

  table_class                 = "STANDARD"
  billing_mode                = "PAY_PER_REQUEST"
  deletion_protection_enabled = true
  stream_enabled              = false

  attribute {
    name = "PK"
    type = "S"
  }

  attribute {
    name = "SK"
    type = "S"
  }

  point_in_time_recovery {
    enabled = false
  }
}
