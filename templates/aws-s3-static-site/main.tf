resource "aws_s3_bucket" "website" {
  bucket = var.bucket_name
  force_destroy = var.force_destroy

  tags = {
    Name = var.bucket_name
    Project = var.name_prefix 
  }
}

resource "aws_s3_bucket_website_configuration" "website" {
  bucket = aws_s3_bucket.website.id

  index_document {
    suffix = "index.html"
  }
  error_document {
    key = "error.html"
  }
}

# Ownership controls (recommended for modern S3 behavior)
resource "aws_s3_bucket_ownership_controls" "website" {
  bucket = aws_s3_bucket.website.id

  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

# Public access block configuration - we disable some protections
# to allow the static website to be publicly accessible.


