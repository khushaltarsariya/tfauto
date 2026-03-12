output "vpc_id" {
  description = "ID of the VPC"
  value       = aws_vpc.this.id
}

output "public_subnet_ids" {
  description = "IDs of the public subnets"
  value       = [for subnet in aws_subnet.public : subnet.id]
}

output "app_subnet_ids" {
  description = "IDs of the private application subnets"
  value       = [for subnet in aws_subnet.app : subnet.id]
}

output "db_subnet_ids" {
  description = "IDs of the private database subnets"
  value       = [for subnet in aws_subnet.db : subnet.id]
}

output "nat_gateway_ids" {
  description = "IDs of the NAT gateways"
  value       = [for nat in aws_nat_gateway.this : nat.id]
}

output "db_subnet_group_name" {
  description = "Database subnet group name"
  value       = aws_db_subnet_group.this.name
}

output "alb_security_group_id" {
  description = "Security group ID for public load balancers"
  value       = aws_security_group.alb.id
}

output "app_security_group_id" {
  description = "Security group ID for the application tier"
  value       = aws_security_group.app.id
}

output "db_security_group_id" {
  description = "Security group ID for the database tier"
  value       = aws_security_group.db.id
}
