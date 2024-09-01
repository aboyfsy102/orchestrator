locals {
  lambda_function_name_fleet_http = "${var.project}-${data.aws_region.current.name}-lambda-fleet-http"
}

# IAM role for the Lambda function
resource "aws_iam_role" "fleet_http" {
  name        = "${var.project}-${data.aws_region.current.name}-fleet-http"
  description = "Role for the ${var.project}-fleet-http lambda function"

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
resource "aws_cloudwatch_log_group" "lambda_fleet_http_logs" {
  name              = "/aws/lambda/${local.lambda_function_name_fleet_http}"
  retention_in_days = 1
}


# Define the Lambda function
resource "aws_lambda_function" "lambda_fleet_http" {
  description   = "${var.project}: lambda function to process CREATE requests to launch EC2 instances"
  filename      = "../dist/lambda_fleet_http.zip"
  function_name = local.lambda_function_name_fleet_http
  role          = aws_iam_role.fleet_http.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  memory_size   = 1024

  source_code_hash = filebase64sha256("../dist/lambda_fleet_http.zip")

  environment {
    variables = {
      LOG_LEVEL         = "INFO"
      DATABASE_URL      = "postgresql://${aws_db_instance.rds.username}:${aws_db_instance.rds.password}@${aws_db_instance.rds.endpoint}/${aws_db_instance.rds.db_name}?sslmode=require"
      DATABASE_USERNAME = "j5v3_lambda"
      DATABASE_PASSWORD = ""
    }
  }

  depends_on = [aws_cloudwatch_log_group.lambda_fleet_http_logs]
}

# IAM policy for the Lambda function
resource "aws_iam_role_policy" "lambda_fleet_http_policy" {
  name = "lambda_fleet_http_policy"
  role = aws_iam_role.fleet_http.id

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
      }
    ]
  })
}