output "bucket_name" {
  description = "Name of the S3 bucket hosting the website"
  value       = aws_s3_bucket.website.bucket
}

output "bucket_arn" {
  description = "ARN of the S3 bucket"
  value       = aws_s3_bucket.website.arn
}

output "website_endpoint" {
  description = "Website endpoint for the static site"
  value       = aws_s3_bucket_website_configuration.website.website_endpoint
}
