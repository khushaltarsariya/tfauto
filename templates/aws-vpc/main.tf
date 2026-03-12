resource "aws_vpc" "aws-vpc" {
  cidr_block           = var.vpc_cidr_block
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "${var.name_prefix}-main"
  }
}

resource "aws_internet_gateway" "aws-vpc-igw" {
  vpc_id = aws_vpc.aws-vpc.id

  tags = {
    Name = "${var.name_prefix}-igw"
  }
}

resource "aws_subnet" "aws-vpc-public-subnet-a" {
  vpc_id                  = aws_vpc.aws-vpc.id
  cidr_block              = var.public_subnet_cidr_block[0]
  availability_zone       = var.availability_zone[0]
  map_public_ip_on_launch = true

  tags = {
    Name = "${var.name_prefix}-public-subnet-a"
  }
}

resource "aws_subnet" "aws-vpc-public-subnet-b" {
  vpc_id                  = aws_vpc.aws-vpc.id
  cidr_block              = var.public_subnet_cidr_block[1]
  availability_zone       = var.availability_zone[1]
  map_public_ip_on_launch = true

  tags = {
    Name = "${var.name_prefix}-public-subnet-b"
  }
}

resource "aws_route_table" "aws-vpc-route-table" {
  vpc_id = aws_vpc.aws-vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.aws-vpc-igw.id
  }

  tags = {
    Name = "${var.name_prefix}-route-table"
  }
}

resource "aws_route_table_association" "aws-vpc-route-assoc-a" {
  subnet_id      = aws_subnet.aws-vpc-public-subnet-a.id
  route_table_id = aws_route_table.aws-vpc-route-table.id
}

resource "aws_route_table_association" "aws-vpc-route-assoc-b" {
  subnet_id      = aws_subnet.aws-vpc-public-subnet-b.id
  route_table_id = aws_route_table.aws-vpc-route-table.id
}