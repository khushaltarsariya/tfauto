output "alb_dns_name" {
  description = "Public DNS name of the application load balancer"
  value       = aws_lb.this.dns_name
}

output "alb_zone_id" {
  description = "Route 53 hosted zone ID of the ALB"
  value       = aws_lb.this.zone_id
}

output "target_group_arn" {
  description = "Target group ARN attached to the Auto Scaling group"
  value       = aws_lb_target_group.this.arn
}

output "autoscaling_group_name" {
  description = "Name of the Auto Scaling group"
  value       = aws_autoscaling_group.this.name
}

output "launch_template_id" {
  description = "Launch template ID used by the Auto Scaling group"
  value       = aws_launch_template.this.id
}

output "alb_security_group_id" {
  description = "Security group ID for the ALB"
  value       = aws_security_group.alb.id
}

output "app_security_group_id" {
  description = "Security group ID for application instances"
  value       = aws_security_group.app.id
}
