AWS Startup Web App

What this template creates:
- 1 VPC with DNS support enabled
- 2 public subnets for a public-facing ALB and NAT gateway
- 2 private application subnets for ECS Fargate backend tasks
- 2 private database subnets for PostgreSQL
- 1 internet gateway
- 1 NAT gateway with 1 Elastic IP for private outbound access
- 1 Application Load Balancer for the frontend entry point
- 1 ECS cluster and service for the backend application tier
- 1 PostgreSQL RDS instance in private subnets
- Security groups for frontend, backend, and database traffic

Recommended use:
- Start here when you want a real startup-style web stack without overprovisioning
- Good for founders and small teams that want a backend, database, and public entry point

Cost notes:
- Defaults keep the stack lean with a single NAT gateway and a single RDS instance
- ECS desired count defaults to 1 to keep compute cost low
- Database is single-AZ by default to reduce spend

Customization ideas:
- Swap the default container image for your backend service
- Add ACM and HTTPS when you are ready for production traffic
- Increase ECS desired count or RDS instance size as traffic grows
