
module "vpc" {
  source             = "../../modules/vpc"
  vpc_cidr           = var.vpc_cidr
  availability_zones = var.availability_zones
  environment        = var.environment
}

module "rds" {
  source             = "../../modules/rds"
  vpc_id             = module.vpc.vpc_id
  private_subnet_ids = module.vpc.private_subnet_ids
  db_name            = var.db_name
  db_username        = var.db_username
  db_password        = var.db_password
  project_name       = var.project_name
  environment        = var.environment
}

module "alb" {
  source            = "../../modules/alb"
  vpc_id            = module.vpc.vpc_id
  public_subnet_ids = module.vpc.public_subnet_ids
  app_port          = var.app_port
  cidr_ipv4         = module.vpc.cidr_ipv4
  environment       = var.environment
  project_name      = var.project_name
}

module "ecr" {
  source       = "../../modules/ecr"
  project_name = var.project_name
  environment  = var.environment
}

module "compute" {
  source                = "../../modules/compute"
  project_name          = var.project_name
  environment           = var.environment
  vpc_id                = module.vpc.vpc_id
  public_subnet_ids     = module.vpc.public_subnet_ids
  app_port              = var.app_port
  alb_security_group_id = module.alb.sg_id
  target_group_arn      = module.alb.tg_arn
  rds_security_group_id = module.rds.security_group_id
  ecr_repository_url    = module.ecr.repository_url
  db_host               = module.rds.db_address
  db_port               = module.rds.db_port
  db_name               = var.db_name
  db_username           = var.db_username
  db_password           = var.db_password
  jwt_secret            = var.jwt_secret
  github_client_id      = var.github_client_id
  github_client_secret  = var.github_client_secret

}

module "frontend" {
  source       = "../../modules/frontend"
  project_name = var.project_name
  environment  = var.environment
  alb_dns_name = module.alb.alb_dns_name
}

module "monitoring" {
  source                  = "../../modules/monitoring"
  project_name            = var.project_name
  environment             = var.environment
  alb_arn_suffix          = module.alb.alb_arn_suffix
  target_group_arn_suffix = module.alb.target_group_arn_suffix
  asg_name                = module.compute.asg_name
  alert_email             = var.alert_email
}
