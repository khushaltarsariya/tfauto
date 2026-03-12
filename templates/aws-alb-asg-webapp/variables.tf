variable "aws_region" {
  description = "AWS region for deployment"
  type        = string
  default     = "us-east-1"
}

variable "name_prefix" {
  description = "Prefix used for naming and tags"
  type        = string
  default     = "webapp"
}

variable "vpc_id" {
  description = "VPC ID where the load balancer and instances will be deployed"
  type        = string
}

variable "alb_subnet_ids" {
  description = "Public subnet IDs for the load balancer"
  type        = list(string)
}

variable "instance_subnet_ids" {
  description = "Subnet IDs for EC2 instances, usually private application subnets"
  type        = list(string)
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t3.micro"
}

variable "ami_id" {
  description = "Optional AMI ID. Leave null to use the latest Amazon Linux 2023 image"
  type        = string
  default     = null
}

variable "instance_port" {
  description = "Port exposed by the web application on the EC2 instances"
  type        = number
  default     = 80
}

variable "health_check_path" {
  description = "Health check path for the target group"
  type        = string
  default     = "/"
}

variable "min_size" {
  description = "Minimum number of instances in the Auto Scaling group"
  type        = number
  default     = 2
}

variable "desired_capacity" {
  description = "Desired number of instances in the Auto Scaling group"
  type        = number
  default     = 2
}

variable "max_size" {
  description = "Maximum number of instances in the Auto Scaling group"
  type        = number
  default     = 4
}

variable "root_volume_size" {
  description = "Size of the root EBS volume in GiB"
  type        = number
  default     = 20
}

variable "allow_ssh_cidrs" {
  description = "Optional CIDR blocks allowed to SSH to the instances"
  type        = list(string)
  default     = []
}

variable "user_data" {
  description = "Optional user data script. Leave empty to use the default NGINX bootstrap"
  type        = string
  default     = ""
}

variable "tags" {
  description = "Additional tags applied to all resources"
  type        = map(string)
  default     = {}
}
