locals {
  dynamo_table_name = join("-", [var.app_name, var.environment])
}
