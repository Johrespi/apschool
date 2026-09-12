resource "aws_db_subnet_group" "this" {
  name_prefix = "${var.project_name}-${var.environment}-subnet-group-"
  subnet_ids  = var.private_subnet_ids
  tags = {
    "Name" = "${var.project_name}-${var.environment}"
  }
}

resource "aws_security_group" "db_sg" {
  name_prefix = "${var.project_name}-${var.environment}-rds-sg-"
  description = "Security group for APSchool DB"
  vpc_id      = var.vpc_id
}

resource "aws_db_instance" "db" {
  identifier_prefix       = "${var.project_name}-${var.environment}-"
  engine                  = "postgres"
  engine_version          = "16.3"
  instance_class          = "db.t4g.micro"
  allocated_storage       = 20
  storage_type            = "gp3"
  max_allocated_storage   = 0
  multi_az                = false
  publicly_accessible     = false
  db_subnet_group_name    = aws_db_subnet_group.this.name
  vpc_security_group_ids  = [aws_security_group.db_sg.id]
  db_name                 = var.db_name
  username                = var.db_username
  password                = var.db_password
  skip_final_snapshot     = true
  backup_retention_period = 0
  deletion_protection     = false
}
