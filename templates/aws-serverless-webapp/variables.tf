variable "aws_region" {
  description = "AWS region for the serverless web app"
  type        = string
  default     = "us-east-1"
}

variable "name_prefix" {
  description = "Prefix used for names and tags"
  type        = string
  default     = "serverless-webapp"
}

variable "runtime" {
  description = "Lambda runtime"
  type        = string
  default     = "python3.12"
}

variable "handler" {
  description = "Lambda handler name"
  type        = string
  default     = "handler.lambda_handler"
}

variable "timeout" {
  description = "Lambda timeout in seconds"
  type        = number
  default     = 10
}

variable "memory_size" {
  description = "Lambda memory size in MiB"
  type        = number
  default     = 256
}

variable "route_key" {
  description = "HTTP API route key"
  type        = string
  default     = "GET /"
}

variable "environment" {
  description = "Environment variables for the Lambda function"
  type        = map(string)
  default     = {}
}

variable "tags" {
  description = "Additional tags applied to all resources"
  type        = map(string)
  default     = {}
}
