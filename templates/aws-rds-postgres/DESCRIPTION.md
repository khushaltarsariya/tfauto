AWS RDS PostgreSQL

What this template creates:
- 1 PostgreSQL RDS instance
- 1 DB subnet group for private database subnets
- 1 security group for database access
- Optional Multi-AZ deployment
- Optional automated backups and deletion protection

Recommended use:
- Use this in an existing VPC after creating private database subnets
- Pair it with `aws-three-tier-vpc` or another VPC template that exposes DB subnet IDs

Notes:
- This template is intentionally opinionated toward private database deployment
- Restrict `allowed_security_group_ids` to application-tier security groups only
