variable "name_prefix" {
  description = "Name prefix for the all aws resources"
  type        = string
  default     = "aws-vpc"
}
variable "aws_region" {
  description = "aws region"
  type        = string
  default     = "us-east-1"
}

variable "vpc_cidr_block" {
  description = "VPC cidr block"
  type        = string
  default     = "10.0.0.0/16"
}

variable "public_subnet_cidr_block" {
  description = "public subnet cidr block"
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24"]
}

variable "availability_zone" {
  description = "list of azs"
  type        = list(string)
  default     = ["us-east-1a", "us-east-1b"]
}

