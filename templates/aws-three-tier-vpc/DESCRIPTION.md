AWS Three-Tier VPC

What this template creates:
- 1 VPC with DNS support enabled
- 2 public subnets for ingress components such as ALB or bastion hosts
- 2 private application subnets for EC2, ECS, or internal services
- 2 private database subnets for RDS or other stateful tiers
- 1 internet gateway
- 1 or 2 NAT gateways depending on configuration
- Public and private route tables with subnet associations
- 1 DB subnet group for use by RDS
- Security groups for ALB, application, and database tiers

Recommended use:
- Start here when building a production-style AWS network layout
- Pair with application templates that need separate public, app, and DB tiers

Notes:
- `single_nat_gateway = true` reduces cost but lowers AZ-level resilience
- Database subnets are private and have no internet route by default
