provider "aws" {
  region                      = var.region
  access_key                  = var.localstack ? "test" : null
  secret_key                  = var.localstack ? "test" : null
  s3_use_path_style           = var.localstack ? false : true
  skip_credentials_validation = var.localstack ? true : false
  skip_metadata_api_check     = var.localstack ? true : false
  skip_requesting_account_id  = var.localstack ? true : false

  default_tags {
    tags = {
      Project = var.app_name
    }
  }

  endpoints {
    dynamodb = var.localstack ? "http://localhost:4566" : null
    iam      = var.localstack ? "http://localhost:4566" : null
    sts      = var.localstack ? "http://localhost:4566" : null
  }
}

data "aws_caller_identity" "current" {}

resource "aws_dynamodb_table" "table" {
  name      = local.dynamo_table_name
  hash_key  = "PK"
  range_key = "ID"

  table_class                 = "STANDARD"
  billing_mode                = "PAY_PER_REQUEST"
  deletion_protection_enabled = false
  stream_enabled              = false

  attribute {
    name = "PK"
    type = "S"
  }

  attribute {
    name = "ID"
    type = "S"
  }

  point_in_time_recovery {
    enabled = false
  }
}
