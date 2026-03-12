AWS ECS Fargate Service

What this template creates:
- 1 ECS cluster
- 1 CloudWatch log group
- 1 task execution role and task role
- 1 ECS task definition and Fargate service
- 1 internet-facing application load balancer with target group and listener
- Security groups for the ALB and ECS service

Recommended use:
- Use this in an existing VPC with at least two public subnets for the ALB
- Place tasks in private application subnets when possible

Notes:
- The default container image is `nginx:stable-alpine`
- Extend this template with autoscaling, HTTPS certificates, and service discovery as needed
