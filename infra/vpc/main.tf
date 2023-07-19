provider "aws" {
  region = var.aws_region
}

resource "aws_eip" "eip" {
  domain = "vpc"

  tags = {
    Name = "points-are-bad-eip"
    AppName = "points-are-bad"
  }
}

resource "aws_internet_gateway" "internet_gw" {
  tags = {
    Name = "points-are-bad-internet-gateway"
    AppName = "points-are-bad"
  }
}

resource "aws_internet_gateway_attachment" "internet_gw_attachment" {
  internet_gateway_id = aws_internet_gateway.internet_gw.id
  vpc_id = aws_vpc.pab_vpc.id
}

resource "aws_nat_gateway" "nat_gw" {
  allocation_id = aws_eip.eip.allocation_id
  connectivity_type = "public"
  subnet_id = aws_subnet.public_subnet_1.id

  tags = {
    Name = "points-are-bad-nat-gateway"
    AppName = "points-are-bad"
  }
}

resource "aws_route" "public_internet_route" {
  route_table_id = aws_route_table.public_route_table.id
  destination_cidr_block = "0.0.0.0/0" # Public internet baby
  gateway_id = aws_internet_gateway.internet_gw.id

  depends_on = [aws_internet_gateway_attachment.internet_gw_attachment]
}

resource "aws_route" "public_route_table_route" {
  route_table_id = aws_route_table.private_route_table.id
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id = aws_nat_gateway.nat_gw.id
}

resource "aws_route_table" "private_route_table" {
  vpc_id = aws_vpc.pab_vpc.id

  tags = {
    Name = "points-are-bad-private-route-table"
    AppName = "points-are-bad"
  }
}

resource "aws_route_table" "public_route_table" {
  vpc_id = aws_vpc.pab_vpc.id

  tags = {
    Name = "points-are-bad-public-route-table"
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

resource "aws_route_table_association" "public_subnet_1" {
  route_table_id = aws_route_table.public_route_table.id
  subnet_id = aws_subnet.public_subnet_1.id
}

resource "aws_route_table_association" "public_subnet_2" {
  route_table_id = aws_route_table.public_route_table.id
  subnet_id = aws_subnet.public_subnet_2.id
}

resource "aws_route_table_association" "public_subnet_3" {
  route_table_id = aws_route_table.public_route_table.id
  subnet_id = aws_subnet.public_subnet_3.id
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

resource "aws_subnet" "public_subnet_1" {
  availability_zone = var.aws_az_1
  cidr_block = var.public_subnet_1_cidr_range
  vpc_id = aws_vpc.pab_vpc.id

  tags = {
    Name = "points-are-bad-public-subnet-1"
    AppName = "points-are-bad"
  }
}

resource "aws_subnet" "public_subnet_2" {
  availability_zone = var.aws_az_2
  cidr_block = var.public_subnet_2_cidr_range
  vpc_id = aws_vpc.pab_vpc.id

  tags = {
    Name = "points-are-bad-public-subnet-2"
    AppName = "points-are-bad"
  }
}

resource "aws_subnet" "public_subnet_3" {
  availability_zone = var.aws_az_3
  cidr_block = var.public_subnet_3_cidr_range
  vpc_id = aws_vpc.pab_vpc.id

  tags = {
    Name = "points-are-bad-public-subnet-3"
    AppName = "points-are-bad"
  }
}

resource "aws_vpc" "pab_vpc" {
  cidr_block = var.vpc_cidr_range
  enable_dns_hostnames = true
  enable_dns_support = true

  tags = {
    Name = "points-are-bad-vpc"
    AppName = "points-are-bad"
  }
}
