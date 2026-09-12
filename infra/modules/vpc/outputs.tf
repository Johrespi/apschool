output "vpc_id" {
  value       = aws_vpc.this.id
  description = "VPC ID"
}

output "cidr_ipv4" {
  value       = aws_vpc.this.cidr_block
  description = "VPC CIDR Block"
}

output "public_subnet_ids" {
  value       = [for s in aws_subnet.public : s.id]
  description = "Public subnets IDs"
}

output "private_subnet_ids" {
  value       = [for s in aws_subnet.private : s.id]
  description = "Private subnets IDs"
}
