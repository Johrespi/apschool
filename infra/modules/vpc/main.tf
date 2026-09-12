locals {
  public_subnets  = { for idx, az in var.availability_zones : az => cidrsubnet(var.vpc_cidr, 8, idx + 1) }
  private_subnets = { for idx, az in var.availability_zones : az => cidrsubnet(var.vpc_cidr, 8, idx + 11) }
}

resource "aws_vpc" "this" {
  cidr_block           = var.vpc_cidr
  enable_dns_support   = true
  enable_dns_hostnames = true


  tags = {
    Name = "apschool-${var.environment}-vpc"
  }
}

resource "aws_internet_gateway" "this" {
  vpc_id = aws_vpc.this.id

  tags = {
    Name = "apschool-${var.environment}-igw"
  }
}

resource "aws_subnet" "public" {
  vpc_id                  = aws_vpc.this.id
  for_each                = local.public_subnets
  availability_zone       = each.key
  cidr_block              = each.value
  map_public_ip_on_launch = true
  tags = {
    Name = "apschool-${var.environment}-public-subnet-${each.key}"
  }
}

resource "aws_subnet" "private" {
  vpc_id                  = aws_vpc.this.id
  for_each                = local.private_subnets
  availability_zone       = each.key
  cidr_block              = each.value
  map_public_ip_on_launch = false
  tags = {
    Name = "apschool-${var.environment}-private-subnet-${each.key}"
  }
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.this.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.this.id
  }
}

resource "aws_route_table_association" "public" {
  for_each       = aws_subnet.public
  subnet_id      = each.value.id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table" "private" {
  vpc_id = aws_vpc.this.id
}

resource "aws_route_table_association" "private" {
  for_each       = aws_subnet.private
  subnet_id      = each.value.id
  route_table_id = aws_route_table.private.id
}
