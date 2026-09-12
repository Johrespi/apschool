resource "aws_security_group" "alb_sg" {
  vpc_id      = var.vpc_id
  name_prefix = "${var.project_name}-${var.environment}-alb-sg-"
  description = "Security group for Application Load Balancer"
}

resource "aws_vpc_security_group_ingress_rule" "allow_http" {
  security_group_id = aws_security_group.alb_sg.id
  ip_protocol       = "tcp"
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 80
  to_port           = 80
}

resource "aws_vpc_security_group_egress_rule" "allow_app_port" {
  security_group_id = aws_security_group.alb_sg.id
  ip_protocol       = "tcp"
  cidr_ipv4         = var.cidr_ipv4
  from_port         = var.app_port
  to_port           = var.app_port
}

resource "aws_lb" "alb" {
  name                       = "apschool-load-balancer"
  internal                   = false
  load_balancer_type         = "application"
  security_groups            = [aws_security_group.alb_sg.id]
  subnets                    = var.public_subnet_ids
  enable_deletion_protection = false
}

resource "aws_lb_target_group" "alb_tg" {
  name        = "appschool-alb-tg"
  port        = var.app_port
  protocol    = "HTTP"
  vpc_id      = var.vpc_id
  target_type = "instance"

  health_check {
    enabled             = true
    path                = var.health_check_path
    port                = "traffic-port"
    protocol            = "HTTP"
    interval            = 30
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 3
    matcher             = "200"
  }

}

resource "aws_lb_listener" "lb_listener" {
  load_balancer_arn = aws_lb.alb.arn
  port              = 80
  protocol          = "HTTP"
  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.alb_tg.arn
  }
}
