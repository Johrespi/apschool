output "rds_endpoint" {
  description = "Connection endpoint for the RDS instance (host:port)"
  value       = module.rds.db_endpoint
}

output "rds_address" {
  description = "Hostname of the RDS instance"
  value       = module.rds.db_address
}

output "rds_port" {
  description = "Port of the RDS instance"
  value       = module.rds.db_port
}

output "rds_security_group_id" {
  description = "Security group ID of the RDS instance"
  value       = module.rds.security_group_id
}

output "alb_dns_name" {
  value = module.alb.alb_dns_name
}

output "ecr_repository_url" {
  value = module.ecr.repository_url
}

output "asg_name" {
  value = module.compute.asg_name
}

output "s3_bucket_name" {
  value = module.frontend.s3_bucket_name
}

output "cloudfront_distribution_id" {
  value = module.frontend.cloudfront_distribution_id
}

output "cloudfront_domain_name" {
  value = module.frontend.cloudfront_domain_name
}

output "sns_topic_arn" {
  value = module.monitoring.sns_topic_arn
}
