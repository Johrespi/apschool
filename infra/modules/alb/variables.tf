variable "project_name" {
  type        = string
  description = "Project name"
}

variable "environment" {
  type        = string
  default     = "prod"
  description = "deploy environment"
}

variable "vpc_id" {
  type        = string
  description = "VPC id"
}

variable "cidr_ipv4" {
  type        = string
  description = "VPC CIDR Block"
}

variable "public_subnet_ids" {
  type        = list(string)
  description = "Public subnet list"
}

variable "app_port" {
  type        = number
  description = "Port"
}

variable "health_check_path" {
  type        = string
  default     = "/api/health"
  description = "health ennpoint"
}
