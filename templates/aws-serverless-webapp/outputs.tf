output "frontend_bucket_name" {
  description = "S3 bucket name for the frontend"
  value       = aws_s3_bucket.frontend.bucket
}

output "frontend_domain_name" {
  description = "CloudFront domain name for the frontend"
  value       = aws_cloudfront_distribution.frontend.domain_name
}

output "frontend_url" {
  description = "Public URL for the frontend"
  value       = "https://${aws_cloudfront_distribution.frontend.domain_name}"
}

output "lambda_function_name" {
  description = "Lambda function name"
  value       = aws_lambda_function.this.function_name
}

output "api_id" {
  description = "HTTP API ID"
  value       = aws_apigatewayv2_api.this.id
}

output "api_endpoint" {
  description = "HTTP API endpoint"
  value       = aws_apigatewayv2_api.this.api_endpoint
}

output "dynamodb_table_name" {
  description = "DynamoDB table name"
  value       = aws_dynamodb_table.this.name
}
