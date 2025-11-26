terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "~>6.0"
    }
  }
  required_version = ">=1.0.0"
}

provider "aws" {
  region = var.aws_region
}

resource "aws_vpc" "aws-basic-vpc" {
  cidr_block = var.vpc_cidr_block
  enable_dns_hostnames = true
  enable_dns_support = true

  tags = {
    Name = "${var.name_prefix}-vpc"
  }
}

resource "aws_internet_gateway" "aws-basic-igw" {
  vpc_id = aws_vpc.aws-basic-vpc.id

  tags = {
    Name = "${var.name_prefix}-igw"
  }
}

resource "aws_subnet" "aws-basic-subnet" {
  vpc_id = aws_vpc.aws-basic-vpc.id
  cidr_block = var.subnet_cidr_block
  availability_zone = var.azs
  map_public_ip_on_launch = true
  
  tags = {
    Name = "${var.name_prefix}-subnet"
  }
}

resource "aws_route_table" "aws-basic-route" {
  vpc_id = aws_vpc.aws-basic-vpc.id
  
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.aws-basic-igw.id
  }
  

  tags = {
    Name = "${var.name_prefix}-route-table"
  }
}

resource "aws_route_table_association" "aws-basic-route-assoc" {
  subnet_id = aws_subnet.aws-basic-subnet.id
  route_table_id = aws_route_table.aws-basic-route.id
}

resource "aws_security_group" "aws-basic-sg" {
  name = "${var.name_prefix}-aws-basic-sg"
  vpc_id = aws_vpc.aws-basic-vpc.id
  description = "sg for the ec2 instance"

  #ssh & http
  ingress {
    description = "SSH"
    from_port = 22
    to_port = 22
    protocol = "tcp"
    cidr_blocks = var.sg-cidr_blocks
  }

  ingress {
    description = "HTTP"
    from_port = 80
    to_port = 80
    protocol = "tcp"
    cidr_blocks = var.sg-cidr_blocks
  }

  egress {
    from_port = 0
    to_port = 0
    protocol = -1
    cidr_blocks = [ "0.0.0.0/0" ]
  }
  tags = {
    Name = "${var.name_prefix}-aws-basic-sg"
  }
}

resource "aws_instance" "aws-basic-instance" {
  ami = var.ami_id
  instance_type = var.instance_type
  subnet_id = aws_subnet.aws-basic-subnet.id
  vpc_security_group_ids = [aws_security_group.aws-basic-sg.id]

  tags = {
    Name= "${var.name_prefix}-aws-basic-instance"
  }
}
