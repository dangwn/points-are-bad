provider "aws" {
  region = var.aws_region
}

#---------------------------------------------
# Variables
#---------------------------------------------
variable "aws_region" {
  type = string
}

variable "aws_az_1" {
  type = string
}

variable "aws_az_2" {
  type = string
}

variable "aws_az_3" {
  type = string
}

variable "vpc_cidr_range" {
  type = string
}

variable "private_subnet_1_cidr_range" {
  type = string
}

variable "private_subnet_2_cidr_range" {
  type = string
}

variable "private_subnet_3_cidr_range" {
  type = string
}

variable "eks_endpoint_subnet_1_cidr_range" {
  type = string
}

variable "eks_endpoint_subnet_2_cidr_range" {
  type = string
}

variable "eks_endpoint_subnet_3_cidr_range" {
  type = string
}

#---------------------------------------------
# Resources
#---------------------------------------------
resource "aws_vpc" "pab_vpc" {
  cidr_block = var.vpc_cidr_range
  enable_dns_hostnames = true
  enable_dns_support = true

  tags = {
    Name = "points-are-bad-vpc"
    AppName = "points-are-bad"
  }
}

resource "aws_route_table" "private_route_table" {
  vpc_id = aws_vpc.pab_vpc.id

  tags = {
    Name = "points-are-bad-private-route-table"
    AppName = "points-are-bad"
  }
}

resource "aws_subnet" "private_subnet_1" {
  availability_zone = var.aws_az_1
  cidr_block = var.private_subnet_1_cidr_range
  vpc_id = aws_vpc.pab_vpc.id

  tags = {
    Name = "points-are-bad-private-subnet-1"
    AppName = "points-are-bad"
  }
}

resource "aws_subnet" "private_subnet_2" {
  availability_zone = var.aws_az_2
  cidr_block = var.private_subnet_2_cidr_range
  vpc_id = aws_vpc.pab_vpc.id

  tags = {
    Name = "points-are-bad-private-subnet-2"
    AppName = "points-are-bad"
  }
}

resource "aws_subnet" "private_subnet_3" {
  availability_zone = var.aws_az_3
  cidr_block = var.private_subnet_3_cidr_range
  vpc_id = aws_vpc.pab_vpc.id

  tags = {
    Name = "points-are-bad-private-subnet-3"
    AppName = "points-are-bad"
  }
}

resource "aws_route_table_association" "private_subnet_1" {
  route_table_id = aws_route_table.private_route_table.id
  subnet_id = aws_subnet.private_subnet_1.id
}

resource "aws_route_table_association" "private_subnet_2" {
  route_table_id = aws_route_table.private_route_table.id
  subnet_id = aws_subnet.private_subnet_2.id
}

resource "aws_route_table_association" "private_subnet_3" {
  route_table_id = aws_route_table.private_route_table.id
  subnet_id = aws_subnet.private_subnet_3.id
}

resource "aws_subnet" "eks_endpoint_subnet_1" {
  availability_zone = var.aws_az_1
  cidr_block = var.eks_endpoint_subnet_1_cidr_range
  vpc_id = aws_vpc.pab_vpc.id

  tags = {
    Name = "points-are-bad-eks-endpoint-subnet-1"
    AppName = "points-are-bad"
  }
}

resource "aws_subnet" "eks_endpoint_subnet_2" {
  availability_zone = var.aws_az_2
  cidr_block = var.eks_endpoint_subnet_2_cidr_range
  vpc_id = aws_vpc.pab_vpc.id

  tags = {
    Name = "points-are-bad-eks-endpoint-subnet-2"
    AppName = "points-are-bad"
  }
}

resource "aws_subnet" "eks_endpoint_subnet_3" {
  availability_zone = var.aws_az_3
  cidr_block = var.eks_endpoint_subnet_3_cidr_range
  vpc_id = aws_vpc.pab_vpc.id

  tags = {
    Name = "points-are-bad-eks-endpoint-subnet-3"
    AppName = "points-are-bad"
  }
} 

resource "aws_route_table_association" "eks_endpoint_subnet_1" {
  route_table_id = aws_route_table.private_route_table.id
  subnet_id = aws_subnet.eks_endpoint_subnet_1.id
}

resource "aws_route_table_association" "eks_endpoint_subnet_2" {
  route_table_id = aws_route_table.private_route_table.id
  subnet_id = aws_subnet.eks_endpoint_subnet_2.id
}

resource "aws_route_table_association" "eks_endpoint_subnet_3" {
  route_table_id = aws_route_table.private_route_table.id
  subnet_id = aws_subnet.eks_endpoint_subnet_3.id
}

#---------------------------------------------
# Outputs
#---------------------------------------------
output "vpc_id" {
  value = aws_vpc.pab_vpc.id
}

output "private_subnet_1_id" {
  value = aws_subnet.private_subnet_1.id
}

output "private_subnet_2_id" {
  value = aws_subnet.private_subnet_2.id
}

output "private_subnet_3_id" {
  value = aws_subnet.private_subnet_3.id
}