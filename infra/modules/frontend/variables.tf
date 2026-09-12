variable "project_name" {
  type = string
}

variable "environment" {
  type    = string
  default = "prod"
}

variable "alb_dns_name" {
  type = string
}
