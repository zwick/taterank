provider "aws" {
  access_key                  = "test"
  secret_key                  = "test"
  region                      = var.region
  s3_use_path_style           = false
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
  default_tags {
    tags = {
      Project = var.app_name
    }
  }

  endpoints {
    dynamodb = "http://localhost:4566"
  }
}

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
