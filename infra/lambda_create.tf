locals {
  lambda_function_name_create = "${var.project}-create"
}

# IAM role for the Lambda function
resource "aws_iam_role" "create" {
  name        = "${var.project}-create"
  description = "Role for the ${var.project}-create lambda function"

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
  name              = "/aws/lambda/${local.lambda_function_name_create}"
  retention_in_days = 1
}


# Define the Lambda function
resource "aws_lambda_function" "lambda_create" {
  description   = "${var.project}: lambda fucntion to process CREATE requests to launch EC2 instances"
  filename      = "../dist/lambda_create.zip"
  function_name = local.lambda_function_name_create
  role          = aws_iam_role.create.arn
  handler       = "lambda_function.lambda_handler"
  runtime       = "python3.12"
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
resource "aws_iam_role_policy" "lambda_create_policy" {
  name = "lambda_create_policy"
  role = aws_iam_role.create.id

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
          "sqs:ReceiveMessage",
          "sqs:DeleteMessage",
          "sqs:GetQueueAttributes"
        ]
        Resource = "*" # Replace with specific SQS queue ARN if possible
      },
      {
        Effect = "Allow"
        Action = [
          "ec2:DescribeImages",
          "ec2:CreateLaunchTemplate",
          "ec2:CreateFleet"
        ]
        Resource = "*"
      }
    ]
  })
}
# SQS trigger for Lambda function
resource "aws_lambda_event_source_mapping" "lambda_create_sqs_trigger" {
  event_source_arn = aws_sqs_queue.spot_fleet_requests.arn # Replace with your SQS queue resource
  function_name    = aws_lambda_function.lambda_create.arn
  batch_size       = 1
}

# SQS queue (if not already defined elsewhere)
resource "aws_sqs_queue" "spot_fleet_requests" {
  name = "spot-fleet-requests"
}