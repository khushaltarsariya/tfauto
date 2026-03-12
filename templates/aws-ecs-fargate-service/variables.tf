variable "aws_region" {
  description = "AWS region for the ECS service"
  type        = string
  default     = "us-east-1"
}

variable "name_prefix" {
  description = "Prefix used for names and tags"
  type        = string
  default     = "ecs-app"
}

variable "vpc_id" {
  description = "VPC ID for the ALB and service"
  type        = string
}

variable "alb_subnet_ids" {
  description = "Subnet IDs for the application load balancer"
  type        = list(string)
}

variable "service_subnet_ids" {
  description = "Subnet IDs for the ECS tasks"
  type        = list(string)
}

variable "container_image" {
  description = "Container image for the ECS task"
  type        = string
  default     = "nginx:stable-alpine"
}

variable "container_port" {
  description = "Port exposed by the container"
  type        = number
  default     = 80
}

variable "cpu" {
  description = "Task CPU units"
  type        = number
  default     = 256
}

variable "memory" {
  description = "Task memory in MiB"
  type        = number
  default     = 512
}

variable "desired_count" {
  description = "Desired number of ECS tasks"
  type        = number
  default     = 2
}

variable "assign_public_ip" {
  description = "Assign public IPs to ECS tasks"
  type        = bool
  default     = false
}

variable "health_check_path" {
  description = "ALB target group health check path"
  type        = string
  default     = "/"
}

variable "environment" {
  description = "Environment variables for the container"
  type        = map(string)
  default     = {}
}

variable "tags" {
  description = "Additional tags applied to all resources"
  type        = map(string)
  default     = {}
}
