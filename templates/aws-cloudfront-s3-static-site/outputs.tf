output "bucket_name" {
  description = "Name of the S3 bucket"
  value       = aws_s3_bucket.site.bucket
}

output "cloudfront_distribution_id" {
  description = "CloudFront distribution ID"
  value       = aws_cloudfront_distribution.site.id
}

output "cloudfront_domain_name" {
  description = "Public CloudFront domain name"
  value       = aws_cloudfront_distribution.site.domain_name
}

output "cloudfront_hosted_zone_id" {
  description = "Hosted zone ID for Route 53 alias records"
  value       = aws_cloudfront_distribution.site.hosted_zone_id
}
