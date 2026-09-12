variable "project_name" {
  type        = string
  description = "Project name"
}

variable "environment" {
  type        = string
  default     = "prod"
  description = "Deployment environment"
}

variable "vpc_id" {
  type        = string
  description = "VPC ID where the compute resources will be deployed"
}

variable "public_subnet_ids" {
  type        = list(string)
  description = "Public subnet IDs where EC2 instances will run"
}

variable "app_port" {
  type        = number
  default     = 8080
  description = "Port where the backend Go application listens"
}

variable "alb_security_group_id" {
  type        = string
  description = "Security Group ID of the ALB to allow inbound traffic"
}

variable "target_group_arn" {
  type        = string
  description = "ARN of the ALB Target Group to attach the Auto Scaling Group"
}

variable "rds_security_group_id" {
  type        = string
  description = "Security Group ID of RDS to authorize EC2 inbound traffic"
}

variable "ecr_repository_url" {
  type        = string
  description = "ECR repo URL"
}

# --- Database Connection Variables ---

variable "db_host" {
  type        = string
  description = "Database host address"
}

variable "db_port" {
  type        = number
  default     = 5432
  description = "Database port"
}

variable "db_name" {
  type        = string
  description = "Database name"
}

variable "db_username" {
  type        = string
  description = "Database username"
}

variable "db_password" {
  type        = string
  sensitive   = true
  description = "Database password"
}

# --- Application Secret Variables ---

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
