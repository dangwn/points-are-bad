output "eks_endpoint_subnet_1_id" {
  value = aws_subnet.eks_endpoint_subnet_1.id
}

output "eks_endpoint_subnet_2_id" {
  value = aws_subnet.eks_endpoint_subnet_2.id
}

output "eks_endpoint_subnet_3_id" {
  value = aws_subnet.eks_endpoint_subnet_3.id
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

output "public_subnet_1_id" {
  value = aws_subnet.public_subnet_1.id
}

output "public_subnet_2_id" {
  value = aws_subnet.public_subnet_2.id
}

output "public_subnet_3_id" {
  value = aws_subnet.public_subnet_3.id
}

output "vpc_id" {
  value = aws_vpc.pab_vpc.id
}