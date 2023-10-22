locals {
  filename = "${path.module}/init.zip"
}

resource "aws_lambda_function" "main" {
  function_name    = "function-terraform"
  role             = aws_iam_role.default.arn
  filename         = local.filename
  source_code_hash = filebase64sha256(local.filename)
  runtime          = "go1.x"
  handler          = "main"

  memory_size = 128
  timeout     = 10

  environment {
    variables = {
      TEST = ""
    }
  }

  lifecycle {
    ignore_changes = [
      filename,
      source_code_hash
    ]
  }

  depends_on = [
    aws_iam_role_policy_attachment.s3_full_access,
    aws_iam_role_policy_attachment.ec2,
    aws_iam_role_policy_attachment.lambda_basic_exec
  ]
}

resource "aws_iam_role" "default" {
  name = "TerraformLambdaRole"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      },
    ]
  })
}

# resource "aws_iam_role_policy_attachment" "lambdainvocation_dynamodb" {
#   role       = aws_iam_role.lambda.name
#   policy_arn = "arn:aws:iam::aws:policy/AWSLambdaInvocation-DynamoDB"
# }

resource "aws_iam_role_policy_attachment" "s3_full_access" {
  role       = aws_iam_role.default.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonS3FullAccess"
}

resource "aws_iam_role_policy_attachment" "ec2" {
  role       = aws_iam_role.default.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2FullAccess"
}

resource "aws_iam_role_policy_attachment" "lambda_basic_exec" {
  role       = aws_iam_role.default.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}
