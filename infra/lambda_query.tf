resource "aws_iam_role" "query" {
  name        = "${var.project}-query"
  description = "Role for the ${var.project}-query lambda function"

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

resource "aws_lambda_function" "query" {
  description   = "${var.project}: lambda fucntion to process all CRUD requests coming from ALB"
  filename      = "../dist/lambda_query.zip"
  function_name = "${var.project}-query"
  role          = aws_iam_role.query.arn
  handler       = "lambda_function.lambda_handler"
  runtime       = "python3.12"
  memory_size = 1024

  source_code_hash = filebase64sha256("../dist/lambda_query.zip")

  environment {
    variables = {
      ENV_VAR_1 = "value1"
      ENV_VAR_2 = "value2"
    }
  }
}

# Add necessary IAM role policies here
resource "aws_iam_role_policy_attachment" "lambda_basic_execution" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  role       = aws_iam_role.query.name
}