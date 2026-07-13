output "vpc_id" {
  description = "ID of the VPC"
  value       = aws_vpc.this.id
}

output "frontend_url" {
  description = "Public URL for the frontend load balancer"
  value       = "http://${aws_lb.this.dns_name}"
}

output "alb_dns_name" {
  description = "DNS name of the frontend load balancer"
  value       = aws_lb.this.dns_name
}

output "ecs_cluster_name" {
  description = "Name of the ECS cluster"
  value       = aws_ecs_cluster.this.name
}

output "ecs_service_name" {
  description = "Name of the ECS service"
  value       = aws_ecs_service.this.name
}

output "db_endpoint" {
  description = "Database endpoint"
  value       = aws_db_instance.this.endpoint
}

output "db_address" {
  description = "Database hostname"
  value       = aws_db_instance.this.address
}

output "db_security_group_id" {
  description = "Security group ID for the database"
  value       = aws_security_group.db.id
}

output "app_security_group_id" {
  description = "Security group ID for the backend"
  value       = aws_security_group.app.id
}

output "alb_security_group_id" {
  description = "Security group ID for the frontend load balancer"
  value       = aws_security_group.alb.id
}
