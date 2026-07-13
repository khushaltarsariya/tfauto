data "aws_iam_policy_document" "lambda_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

locals {
  common_tags = merge(
    {
      Project   = var.name_prefix
      ManagedBy = "tfauto"
    },
    var.tags
  )

  frontend_html = <<-HTML
    <!doctype html>
    <html lang="en">
      <head>
        <meta charset="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <title>${var.name_prefix}</title>
        <style>
          body { font-family: Arial, sans-serif; margin: 0; background: #0f172a; color: #e2e8f0; }
          .wrap { max-width: 720px; margin: 8vh auto; padding: 32px; }
          .card { background: #111827; border: 1px solid #334155; border-radius: 16px; padding: 32px; }
          code { background: #1f2937; padding: 2px 6px; border-radius: 6px; }
          a { color: #60a5fa; }
        </style>
      </head>
      <body>
        <div class="wrap">
          <div class="card">
            <h1>${var.name_prefix}</h1>
            <p>Static frontend delivered by CloudFront and backed by Lambda + API Gateway.</p>
            <p>Use this as a low-cost startup starter.</p>
            <p>Backend API: <code>${aws_apigatewayv2_api.this.api_endpoint}</code></p>
          </div>
        </div>
      </body>
    </html>
  HTML
}

resource "aws_s3_bucket" "frontend" {
  bucket = lower(replace("${var.name_prefix}-frontend", "_", "-"))

  tags = merge(local.common_tags, {
    Name = "${var.name_prefix}-frontend"
  })
}

resource "aws_s3_bucket_ownership_controls" "frontend" {
  bucket = aws_s3_bucket.frontend.id

  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_public_access_block" "frontend" {
  bucket = aws_s3_bucket.frontend.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_cloudfront_origin_access_control" "frontend" {
  name                              = "${var.name_prefix}-oac"
  description                       = "Origin access control for the frontend bucket"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

resource "aws_s3_object" "index" {
  bucket       = aws_s3_bucket.frontend.id
  key          = "index.html"
  content_type = "text/html"
  content      = local.frontend_html
  etag         = md5(local.frontend_html)
}

resource "aws_cloudfront_distribution" "frontend" {
  enabled             = true
  default_root_object = "index.html"
  is_ipv6_enabled     = true
  comment             = "${var.name_prefix} frontend"

  origin {
    domain_name              = aws_s3_bucket.frontend.bucket_regional_domain_name
    origin_id                = "frontend-s3"
    origin_access_control_id = aws_cloudfront_origin_access_control.frontend.id
  }

  default_cache_behavior {
    target_origin_id       = "frontend-s3"
    viewer_protocol_policy = "redirect-to-https"
    allowed_methods        = ["GET", "HEAD", "OPTIONS"]
    cached_methods         = ["GET", "HEAD", "OPTIONS"]
    compress               = true

    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }

  tags = local.common_tags
}

resource "aws_s3_bucket_policy" "frontend" {
  bucket = aws_s3_bucket.frontend.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "AllowCloudFrontRead"
        Effect = "Allow"
        Principal = {
          Service = "cloudfront.amazonaws.com"
        }
        Action   = "s3:GetObject"
        Resource = "${aws_s3_bucket.frontend.arn}/*"
        Condition = {
          StringEquals = {
            "AWS:SourceArn" = aws_cloudfront_distribution.frontend.arn
          }
        }
      }
    ]
  })

  depends_on = [aws_s3_bucket_public_access_block.frontend]
}

resource "aws_dynamodb_table" "this" {
  name         = "${var.name_prefix}-table"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "pk"
  range_key    = "sk"

  attribute {
    name = "pk"
    type = "S"
  }

  attribute {
    name = "sk"
    type = "S"
  }

  tags = merge(local.common_tags, {
    Name = "${var.name_prefix}-table"
  })
}

resource "aws_iam_role" "lambda" {
  name               = "${var.name_prefix}-lambda-role"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume_role.json

  tags = merge(local.common_tags, {
    Name = "${var.name_prefix}-lambda-role"
  })
}

resource "aws_iam_role_policy_attachment" "basic_execution" {
  role       = aws_iam_role.lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy" "dynamodb" {
  name = "${var.name_prefix}-dynamodb-access"
  role = aws_iam_role.lambda.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem",
          "dynamodb:Query",
          "dynamodb:Scan"
        ]
        Resource = [
          aws_dynamodb_table.this.arn,
          "${aws_dynamodb_table.this.arn}/index/*"
        ]
      }
    ]
  })
}

data "archive_file" "lambda_zip" {
  type        = "zip"
  output_path = "${path.module}/lambda_payload.zip"

  source {
    content  = <<-PYTHON
      import json
      import os

      def lambda_handler(event, context):
          return {
              "statusCode": 200,
              "headers": {"content-type": "application/json"},
              "body": json.dumps({
                  "message": "Hello from tfauto",
                  "table": os.environ.get("TABLE_NAME", "")
              })
          }
    PYTHON
    filename = "handler.py"
  }
}

resource "aws_lambda_function" "this" {
  function_name    = "${var.name_prefix}-function"
  role             = aws_iam_role.lambda.arn
  handler          = var.handler
  runtime          = var.runtime
  filename         = data.archive_file.lambda_zip.output_path
  source_code_hash = data.archive_file.lambda_zip.output_base64sha256
  timeout          = var.timeout
  memory_size      = var.memory_size

  environment {
    variables = merge(var.environment, {
      TABLE_NAME = aws_dynamodb_table.this.name
    })
  }

  tags = merge(local.common_tags, {
    Name = "${var.name_prefix}-function"
  })
}

resource "aws_apigatewayv2_api" "this" {
  name          = "${var.name_prefix}-http-api"
  protocol_type = "HTTP"

  tags = local.common_tags
}

resource "aws_apigatewayv2_integration" "lambda" {
  api_id                 = aws_apigatewayv2_api.this.id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_function.this.invoke_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "this" {
  api_id    = aws_apigatewayv2_api.this.id
  route_key = var.route_key
  target    = "integrations/${aws_apigatewayv2_integration.lambda.id}"
}

resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.this.id
  name        = "$default"
  auto_deploy = true

  tags = local.common_tags
}

resource "aws_lambda_permission" "api_gateway" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.this.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.this.execution_arn}/*/*"
}
