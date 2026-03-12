variable "aws_region" {
  description = "aws region to deploy the s3 website"
  type        = string
  default     = "us-east-1"
}

variable "name_prefix" {
  description = "name prefix use for taging"
  type        = string
  default     = "s3-static-site"
}

variable "bucket_name" {
  description = "name of the static bucket"
  type        = string
  default     = "s3-static-hosting"
}

variable "force_destroy" {
  description = "Allow bucket to destroy even if it contain objects"
  type        = bool
  default     = false
}