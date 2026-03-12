output "ecs_cluster_name" {
  description = "Name of the ECS cluster"
  value       = aws_ecs_cluster.this.name
}

output "ecs_service_name" {
  description = "Name of the ECS service"
  value       = aws_ecs_service.this.name
}

output "alb_dns_name" {
  description = "DNS name of the application load balancer"
  value       = aws_lb.this.dns_name
}

output "target_group_arn" {
  description = "Target group ARN for the ECS service"
  value       = aws_lb_target_group.this.arn
}

output "service_security_group_id" {
  description = "Security group ID for ECS tasks"
  value       = aws_security_group.service.id
}
