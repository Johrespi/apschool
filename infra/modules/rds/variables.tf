variable "project_name" {
  type        = string
  description = "Project name"
}

variable "environment" {
  type        = string
  description = "Project environment"
}

variable "vpc_id" {
  type        = string
  description = "VPC id"
}

variable "private_subnet_ids" {
  type        = list(string)
  description = "Private subnet list"
}

variable "db_name" {
  type        = string
  description = "DB name"
}

variable "db_username" {
  type        = string
  description = "DB username"
}

variable "db_password" {
  type        = string
  description = "DB password"
  sensitive   = true
}
