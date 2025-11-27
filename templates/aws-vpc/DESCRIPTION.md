What This Template Creates:
Core VPC Infrastructure:
✓ 1 VPC (10.0.0.0/16) with DNS support enabled
✓ 1 Internet Gateway for internet access
✓ 2 Public Subnets across different Availability Zones:

Public Subnet A in us-east-1a (10.0.1.0/24)

Public Subnet B in us-east-1b (10.0.2.0/24)
✓ 1 Route Table with default route to Internet Gateway
✓ 2 Route Table Associations connecting subnets to the route table

Network Architecture:
Public Subnets with auto-assign public IP enabled

Internet-facing resources can be deployed in either subnet

High availability across two Availability Zones

Full internet access via Internet Gateway