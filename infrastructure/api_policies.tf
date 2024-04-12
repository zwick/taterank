resource "aws_iam_user" "api_user" {
  name = join("-", ["TaterankApi", var.environment])
}

resource "aws_iam_group" "api_user_group" {
  name = join("-",["TaterankApi", var.environment])
}

resource "aws_iam_role" "api_user_role" {
  name = join("-", ["TaterankApi", var.environment])
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "api_user_role_policy_attachment" {
  policy_arn = aws_iam_policy.api_access_policy.arn
  role       = aws_iam_role.api_user_role.name
}

resource "aws_iam_policy" "api_access_policy" {
  name        = join("-",["TaterankApi", var.environment])
  description = "Policy for Taterank API access"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        "Action": [
          "logs:CreateLogStream",
          "logs:CreateLogGroup",
          "logs:TagResource"
        ],
        "Resource": [
          "arn:aws:logs:${var.region}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/taterank-api-${var.environment}*:*"
        ],
        "Effect": "Allow"
      },
      {
        "Action": [
          "logs:PutLogEvents"
        ],
        "Resource": [
          "arn:aws:logs:${var.region}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/taterank-api-${var.environment}*:*:*"
        ],
        "Effect": "Allow"
      },
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
