variable "project_name" {
  type        = string
  description = "Name of the project"
}

variable "environment" {
  type        = string
  description = "Environment name"
}

variable "alb_arn_suffix" {
  type        = string
  description = "Suffix of the ALB ARN"
}

variable "target_group_arn_suffix" {
  type        = string
  description = "Suffix of the Target Group ARN"
}

variable "asg_name" {
  type        = string
  description = "Name of the Auto Scaling Group"
}

variable "alert_email" {
  type        = string
  description = "Email address for alerts"
}
