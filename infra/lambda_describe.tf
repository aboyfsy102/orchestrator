locals {
  lambda_function_name_describe = "${var.project}-${data.aws_region.current.name}-lambda-describe-ec2"
}

# IAM role for the Lambda function
resource "aws_iam_role" "describe" {
  name        = "${var.project}-describe"
  description = "Role for the ${var.project}-describe lambda function"

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
resource "aws_cloudwatch_log_group" "lambda_create_logs" {
  name              = "/aws/lambda/${local.lambda_function_name_describe}"
  retention_in_days = 1
}


# Define the Lambda function
resource "aws_lambda_function" "lambda_describe" {
  description   = "${var.project}: lambda fucntion to process CREATE requests to launch EC2 instances"
  filename      = "../dist/lambda_describe.zip"
  function_name = local.lambda_function_name_describe
  role          = aws_iam_role.describe.arn
  handler       = "lambda_function.lambda_handler"
  runtime       = "provided.al2023"
  memory_size   = 1024

  source_code_hash = filebase64sha256("../dist/lambda_query.zip")

  environment {
    variables = {
      LOG_LEVEL = "INFO"
    }
  }

  depends_on = [aws_cloudwatch_log_group.lambda_log_group]
}

# IAM policy for the Lambda function
resource "aws_iam_role_policy" "lambda_describe_policy" {
  name = "lambda_describe_policy"
  role = aws_iam_role.describe.id

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
# SQS trigger for Lambda function
resource "aws_lambda_event_source_mapping" "lambda_describe_sqs_trigger" {
  event_source_arn = aws_sqs_queue.spot_fleet_requests.arn # Replace with your SQS queue resource
  function_name    = aws_lambda_function.lambda_describe.arn
  batch_size       = 1
}
