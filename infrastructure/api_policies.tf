resource "aws_iam_user" "api_user" {
  name = join("-", ["TaterankApi", var.environment])
}

resource "aws_iam_group" "api_user_group" {
  name = join("-",["TaterankApi", var.environment])
}

resource "aws_iam_policy" "api_access_policy" {
  name        = join("-",["TaterankApi", var.environment])
  description = "Policy for Taterank API access"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "TaterankDynamoDBAccess"
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:GetRecords",
          "dynamodb:Scan",
          "dynamodb:Query",
          "dynamodb:PartiQLSelect",
          "dynamodb:DescribeTable",
          "dynamodb:BatchGetItem",
          "dynamodb:BatchWriteItem",
          "dynamodb:DeleteItem",
          "dynamodb:PartiQLDelete",
          "dynamodb:PartiQLInsert",
          "dynamodb:PartiQLUpdate",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem"
        ]
        Resource = [
          "arn:aws:dynamodb:${var.region}:${data.aws_caller_identity.current.account_id}:table/${local.dynamo_table_name}"
        ]
      }
    ]
  })
}

resource "aws_iam_group_policy_attachment" "api_access_policy_attachment" {
  group      = aws_iam_group.api_user_group.name
  policy_arn = aws_iam_policy.api_access_policy.arn
}
