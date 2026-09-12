variable "vpc_cidr" {
  type        = string
  description = "VPC CIDR block"
}

variable "availability_zones" {
  type        = list(string)
  description = "AZs"
}

variable "environment" {
  type        = string
  default     = "prod"
  description = "deploy environment"
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

variable "project_name" {
  type        = string
  description = "Project name"
}

variable "app_port" {
  type        = number
  description = "apschool backend port"
}

variable "jwt_secret" {
  type        = string
  sensitive   = true
  description = "JWT secret for token signing"
}

variable "github_client_id" {
  type        = string
  default     = ""
  description = "GitHub OAuth Client ID"
}

variable "github_client_secret" {
  type        = string
  sensitive   = true
  default     = ""
  description = "GitHub OAuth Client Secret"
}

variable "alert_email" {
  type        = string
  description = "Email address for alerts"
}
