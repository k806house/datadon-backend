// ${function_name} is the name of the function to be called
// ${endpoint} is the endpoint to be called


/// function
resource "aws_lambda_function" lambda {
  function_name = var.function_name

  s3_bucket = aws_s3_object.lambda.bucket
  s3_key    = aws_s3_object.lambda.key

  runtime = "go1.x"
  handler = "main"

  source_code_hash = data.archive_file.lambda.output_base64sha256
  role = var.role
}

resource "aws_cloudwatch_log_group" "log_group" {
  name = "/aws/lambda/${aws_lambda_function.lambda.function_name}"

  retention_in_days = 30
}

resource "aws_apigatewayv2_integration" "lambda" {
  api_id = var.api_id

  integration_uri    = aws_lambda_function.lambda.invoke_arn
  integration_type   = "AWS_PROXY"
  integration_method = "POST"
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "lambda" {
  api_id = var.api_id

  route_key = "POST ${var.endpoint}"
  target    = "integrations/${aws_apigatewayv2_integration.lambda.id}"
}

resource "aws_lambda_permission" "lambda" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${var.execution_arn}/*/*"
}

// archive
data "archive_file" "lambda" {
  type = "zip"

  source_dir  = "${path.module}/../../build${var.endpoint}"
  output_path = "${path.module}/../../build/${var.function_name}.zip"
}

// upload
resource "aws_s3_object" "lambda" {
  bucket = var.s3_bucket

  key    = "lambda_${var.function_name}.zip"
  source = data.archive_file.lambda.output_path

  etag = filemd5(data.archive_file.lambda.output_path)
}