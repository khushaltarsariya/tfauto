variable "aws_region" {
  description = "AWS region to deploy resources"
  type        = string
  default     = "us-east-1"
}

variable "name_prefix" {
  description = "Name of the EC2 instance"
  type        = string
  default     = "aws-basic"
}
variable "ami_id" {
  description = "ami id for the ec2"
  type = string
  default = "ami-0ecb62995f68bb549"
  
}
variable "instance_type" {
  description = "Type of EC2 instance"
  type        = string
  default     = "t2.micro"
}

variable "vpc_cidr_block" {
  description = "specifing the cidr block"
  type = string
  default = "10.0.0.0/16"
}

variable "subnet_cidr_block" {
  description = "specify the public subnet"
  type = string
  default =  "10.0.1.0/24"
}

variable "sg-cidr_blocks" {
  description = "Sg cidr blocks"
  type = list(string)
  default = [ "0.0.0.0/0" ]
}

variable "azs" {
  description = "availability zone"
  type = string
  default = "us-east-1a"
}

