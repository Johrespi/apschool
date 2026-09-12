data "aws_ami" "amazon_linux_2023" {
  most_recent = true
  owners      = ["amazon"]
  filter {
    name   = "name"
    values = ["al2023-ami-2023.*-kernel-6.1-x86_64"]
  }
}

resource "aws_iam_role" "this" {
  name = "${var.project_name}-${var.environment}-ec2-role"
  assume_role_policy = jsonencode({
    Version : "2012-10-17"
    Statement : [
      {
        Effect : "Allow"
        Action : "sts:AssumeRole"
        Principal : { Service = "ec2.amazonaws.com" }
      }
    ]
  })
}


resource "aws_iam_role_policy_attachment" "this" {
  role       = aws_iam_role.this.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_role_policy_attachment" "ecr_read" {
  role       = aws_iam_role.this.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
}

resource "aws_iam_instance_profile" "this" {
  name = "${var.project_name}-${var.environment}-instance-profile"
  role = aws_iam_role.this.name
}

resource "aws_security_group" "this" {
  name_prefix = "${var.project_name}-${var.environment}-app-sg-"
  description = "Security group for APSchool EC2 instances"
  vpc_id      = var.vpc_id
}


resource "aws_vpc_security_group_ingress_rule" "allow_http" {
  security_group_id            = aws_security_group.this.id
  ip_protocol                  = "tcp"
  referenced_security_group_id = var.alb_security_group_id
  from_port                    = var.app_port
  to_port                      = var.app_port
}

resource "aws_vpc_security_group_egress_rule" "allow_all_traffic" {
  security_group_id = aws_security_group.this.id
  ip_protocol       = "-1"
  cidr_ipv4         = "0.0.0.0/0"
}


resource "aws_vpc_security_group_ingress_rule" "allow_db" {
  security_group_id            = var.rds_security_group_id
  ip_protocol                  = "tcp"
  referenced_security_group_id = aws_security_group.this.id
  from_port                    = 5432
  to_port                      = 5432
}

resource "aws_launch_template" "this" {
  name_prefix   = "${var.project_name}-${var.environment}-lt-"
  image_id      = data.aws_ami.amazon_linux_2023.id
  instance_type = "t3.micro"

  iam_instance_profile {
    arn = aws_iam_instance_profile.this.arn
  }

  network_interfaces {
    associate_public_ip_address = true
    security_groups             = [aws_security_group.this.id]
  }

  user_data = base64encode(<<-EOF
    #!/bin/bash
    dnf update -y && dnf install -y docker awscli 
    systemctl start docker && systemctl enable docker
    aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin ${var.ecr_repository_url}
    docker pull ${var.ecr_repository_url}:latest
    docker run -d --name apschool-api --restart always -p 8080:8080 -e APSCHOOL_DB_HOST="${var.db_host}" -e APSCHOOL_DB_PORT="${var.db_port}" -e APSCHOOL_DB_DATABASE="${var.db_name}" -e APSCHOOL_DB_USERNAME="${var.db_username}" -e APSCHOOL_DB_PASSWORD="${var.db_password}" -e APSCHOOL_DB_SCHEMA="public" -e APSCHOOL_DB_SSLMODE="disable" -e PORT="${var.app_port}" -e JWT_SECRET="${var.jwt_secret}" -e GITHUB_CLIENT_ID="${var.github_client_id}" -e GITHUB_CLIENT_SECRET="${var.github_client_secret}" ${var.ecr_repository_url}:latest
    sleep 5 && docker exec apschool-api ./seed 
  EOF
  )

}


resource "aws_autoscaling_group" "this" {
  name_prefix               = "${var.project_name}-${var.environment}-asg-"
  vpc_zone_identifier       = var.public_subnet_ids
  target_group_arns         = [var.target_group_arn]
  health_check_type         = "ELB"
  health_check_grace_period = 300
  min_size                  = 1
  max_size                  = 1
  desired_capacity          = 1

  launch_template {
    id      = aws_launch_template.this.id
    version = "$Latest"
  }

  tag {
    key                 = "Name"
    value               = "${var.project_name}-${var.environment}-backend"
    propagate_at_launch = true
  }

}
