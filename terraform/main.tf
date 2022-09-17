// requirements

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.1.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.2.0"
    }
  }

  required_version = "~> 1.0"
}

provider "aws" {
  region = "eu-west-1"
}

// bucket
resource "random_pet" "lambda_bucket_name" {
  prefix = "learn-terraform-functions"
  length = 4
}

resource "aws_s3_bucket" "lambda_bucket" {
  bucket = random_pet.lambda_bucket_name.id
}

resource "aws_s3_bucket_acl" "bucket_acl" {
  bucket = aws_s3_bucket.lambda_bucket.id
  acl    = "private"
}

// role
resource "aws_iam_role" "lambda_exec" {
  name = "serverless_lambda"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Sid    = ""
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_policy" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_apigatewayv2_api" "lambda" {
  name          = "serverless_lambda_gw"
  protocol_type = "HTTP"
  cors_configuration  {
    allow_origins = ["*"]
    allow_headers = ["*"]
    allow_methods = ["*"]
    expose_headers = ["*"]
  }
}

resource "aws_apigatewayv2_stage" "lambda" {
  api_id = aws_apigatewayv2_api.lambda.id

  name        = "serverless_lambda_stage"
  auto_deploy = true

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api_gw.arn

    format = jsonencode({
      requestId               = "$context.requestId"
      sourceIp                = "$context.identity.sourceIp"
      requestTime             = "$context.requestTime"
      protocol                = "$context.protocol"
      httpMethod              = "$context.httpMethod"
      resourcePath            = "$context.resourcePath"
      routeKey                = "$context.routeKey"
      status                  = "$context.status"
      responseLength          = "$context.responseLength"
      integrationErrorMessage = "$context.integrationErrorMessage"
    }
    )
  }
}


resource "aws_cloudwatch_log_group" "api_gw" {
  name = "/aws/api_gw/${aws_apigatewayv2_api.lambda.name}"

  retention_in_days = 30
}


module "user_get" {
  source = "./lambda"

  function_name = "user_get"
  endpoint = "/user/get"
  s3_bucket = aws_s3_bucket.lambda_bucket.id
  role = aws_iam_role.lambda_exec.arn
  api_id = aws_apigatewayv2_api.lambda.id
  execution_arn = aws_apigatewayv2_stage.lambda.execution_arn
}

module "user_auth" {
  source = "./lambda"

  function_name = "user_auth"
  endpoint = "/user/auth"
  s3_bucket = aws_s3_bucket.lambda_bucket.id
  role = aws_iam_role.lambda_exec.arn
  api_id = aws_apigatewayv2_api.lambda.id
  execution_arn = aws_apigatewayv2_stage.lambda.execution_arn
}

module "user_logout" {
  source = "./lambda"

  function_name = "user_logout"
  endpoint = "/user/logout"
  s3_bucket = aws_s3_bucket.lambda_bucket.id
  role = aws_iam_role.lambda_exec.arn
  api_id = aws_apigatewayv2_api.lambda.id
  execution_arn = aws_apigatewayv2_stage.lambda.execution_arn
}

module "exam_create" {
  source = "./lambda"

  function_name = "exam_create"
  endpoint = "/exam/create"
  s3_bucket = aws_s3_bucket.lambda_bucket.id
  role = aws_iam_role.lambda_exec.arn
  api_id = aws_apigatewayv2_api.lambda.id
  execution_arn = aws_apigatewayv2_stage.lambda.execution_arn
}

module "exam_get" {
  source = "./lambda"

  function_name = "exam_get"
  endpoint = "/exam/get"
  s3_bucket = aws_s3_bucket.lambda_bucket.id
  role = aws_iam_role.lambda_exec.arn
  api_id = aws_apigatewayv2_api.lambda.id
  execution_arn = aws_apigatewayv2_stage.lambda.execution_arn
}

module "exam_upload" {
  source = "./lambda"

  function_name = "exam_upload"
  endpoint = "/exam/upload"
  s3_bucket = aws_s3_bucket.lambda_bucket.id
  role = aws_iam_role.lambda_exec.arn
  api_id = aws_apigatewayv2_api.lambda.id
  execution_arn = aws_apigatewayv2_stage.lambda.execution_arn
}

module "exam_files" {
  source = "./lambda"

  function_name = "exam_files"
  endpoint = "/exam/files"
  s3_bucket = aws_s3_bucket.lambda_bucket.id
  role = aws_iam_role.lambda_exec.arn
  api_id = aws_apigatewayv2_api.lambda.id
  execution_arn = aws_apigatewayv2_stage.lambda.execution_arn
}

module "exam_tags_get" {
  source = "./lambda"

  function_name = "exam_tags_get"
  endpoint = "/exam/tags/get"
  s3_bucket = aws_s3_bucket.lambda_bucket.id
  role = aws_iam_role.lambda_exec.arn
  api_id = aws_apigatewayv2_api.lambda.id
  execution_arn = aws_apigatewayv2_stage.lambda.execution_arn
}

module "exam_tags_set" {
  source = "./lambda"

  function_name = "exam_tags_set"
  endpoint = "/exam/tags/set"
  s3_bucket = aws_s3_bucket.lambda_bucket.id
  role = aws_iam_role.lambda_exec.arn
  api_id = aws_apigatewayv2_api.lambda.id
  execution_arn = aws_apigatewayv2_stage.lambda.execution_arn
}

module "tag_get" {
  source = "./lambda"

  function_name = "tag_get"
  endpoint = "/tag/get"
  s3_bucket = aws_s3_bucket.lambda_bucket.id
  role = aws_iam_role.lambda_exec.arn
  api_id = aws_apigatewayv2_api.lambda.id
  execution_arn = aws_apigatewayv2_stage.lambda.execution_arn
}

module "study_create" {
  source = "./lambda"

  function_name = "study_create"
  endpoint = "/study/create"
  s3_bucket = aws_s3_bucket.lambda_bucket.id
  role = aws_iam_role.lambda_exec.arn
  api_id = aws_apigatewayv2_api.lambda.id
  execution_arn = aws_apigatewayv2_stage.lambda.execution_arn
}

module "study_get" {
  source = "./lambda"

  function_name = "study_get"
  endpoint = "/study/get"
  s3_bucket = aws_s3_bucket.lambda_bucket.id
  role = aws_iam_role.lambda_exec.arn
  api_id = aws_apigatewayv2_api.lambda.id
  execution_arn = aws_apigatewayv2_stage.lambda.execution_arn
}

module "study_tags_get" {
  source = "./lambda"

  function_name = "study_tags_get"
  endpoint = "/study/tags/get"
  s3_bucket = aws_s3_bucket.lambda_bucket.id
  role = aws_iam_role.lambda_exec.arn
  api_id = aws_apigatewayv2_api.lambda.id
  execution_arn = aws_apigatewayv2_stage.lambda.execution_arn
}

module "study_tags_set" {
  source = "./lambda"

  function_name = "study_tags_set"
  endpoint = "/study/tags/set"
  s3_bucket = aws_s3_bucket.lambda_bucket.id
  role = aws_iam_role.lambda_exec.arn
  api_id = aws_apigatewayv2_api.lambda.id
  execution_arn = aws_apigatewayv2_stage.lambda.execution_arn
}
#
#/study/status

module "study_status" {
  source = "./lambda"

  function_name = "study_status"
  endpoint = "/study/status"
  s3_bucket = aws_s3_bucket.lambda_bucket.id
  role = aws_iam_role.lambda_exec.arn
  api_id = aws_apigatewayv2_api.lambda.id
  execution_arn = aws_apigatewayv2_stage.lambda.execution_arn
}

#/study/match/invites

module "study_match_invites" {
  source = "./lambda"

  function_name = "study_match_invites"
  endpoint = "/study/match/invites"
  s3_bucket = aws_s3_bucket.lambda_bucket.id
  role = aws_iam_role.lambda_exec.arn
  api_id = aws_apigatewayv2_api.lambda.id
  execution_arn = aws_apigatewayv2_stage.lambda.execution_arn
}
#/study/match/list
module "study_match_list" {
  source = "./lambda"

  function_name = "study_match_list"
  endpoint = "/study/match/list"
  s3_bucket = aws_s3_bucket.lambda_bucket.id
  role = aws_iam_role.lambda_exec.arn
  api_id = aws_apigatewayv2_api.lambda.id
  execution_arn = aws_apigatewayv2_stage.lambda.execution_arn
}
#/study/match/response
module "study_match_response" {
  source = "./lambda"

  function_name = "study_match_response"
  endpoint = "/study/match/response"
  s3_bucket = aws_s3_bucket.lambda_bucket.id
  role = aws_iam_role.lambda_exec.arn
  api_id = aws_apigatewayv2_api.lambda.id
  execution_arn = aws_apigatewayv2_stage.lambda.execution_arn
}
#/study/download

module "study_match_download" {
  source = "./lambda"

  function_name = "study_match_download"
  endpoint = "/study/match/download"
  s3_bucket = aws_s3_bucket.lambda_bucket.id
  role = aws_iam_role.lambda_exec.arn
  api_id = aws_apigatewayv2_api.lambda.id
  execution_arn = aws_apigatewayv2_stage.lambda.execution_arn
}