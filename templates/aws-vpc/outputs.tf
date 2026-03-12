output "vpc_id" {
  value = aws_vpc.aws-vpc.id
}

output "subnet_id_a" {
  value = aws_subnet.aws-vpc-public-subnet-a.id
}

output "subnet_id_b" {
  value = aws_subnet.aws-vpc-public-subnet-b.id
}

output "aws_internet_gateway" {
  value = aws_internet_gateway.aws-vpc-igw.id
}
