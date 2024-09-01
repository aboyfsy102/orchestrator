locals {
  lambda_function_name_describe_http = "${var.project}-${data.aws_region.current.name}-lambda-describe-http"
}

# IAM role for the Lambda function
resource "aws_iam_role" "describe_http" {
  name        = "${var.project}-${data.aws_region.current.name}-describe-http"
  description = "Role for the ${var.project}-describe-http lambda function"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

# CloudWatch Log Group for Lambda function
resource "aws_cloudwatch_log_group" "lambda_describe_http_logs" {
  name              = "/aws/lambda/${local.lambda_function_name_describe_http}"
  retention_in_days = 1
}


# Define the Lambda function
resource "aws_lambda_function" "lambda_describe_http" {
  description   = "${var.project}: lambda function to process CREATE requests to launch EC2 instances"
  filename      = "../dist/lambda_describe_http.zip"
  function_name = local.lambda_function_name_describe_http
  role          = aws_iam_role.describe_http.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  memory_size   = 1024

  source_code_hash = filebase64sha256("../dist/lambda_describe_http.zip")

  environment {
    variables = {
      LOG_LEVEL = "INFO"
    }
  }

    depends_on = [aws_cloudwatch_log_group.lambda_describe_http_logs]
}

# IAM policy for the Lambda function
resource "aws_iam_role_policy" "lambda_describe_policy" {
  name = "lambda_describe_policy"
  role = aws_iam_role.describe_http.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Resource = "arn:aws:logs:*:*:*"
      },
      {
        Effect = "Allow"
        Action = [
          "ec2:DescribeInstances"
        ]
        Resource = "*"
      }
    ]
  })
}