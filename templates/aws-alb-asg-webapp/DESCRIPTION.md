AWS ALB + Auto Scaling Web App

What this template creates:
- 1 internet-facing application load balancer
- 1 target group and listener
- 1 launch template for EC2 web instances
- 1 Auto Scaling group across multiple subnets
- Security groups for the ALB and application instances
- IAM role and instance profile with SSM access

Recommended use:
- Use this in an existing VPC after creating networking with `aws-three-tier-vpc`
- Put the ALB in public subnets and the EC2 instances in private application subnets

Notes:
- The default user data installs and starts NGINX
- Pass your own AMI ID if you do not want the latest Amazon Linux image
