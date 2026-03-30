# main.tf

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

# Networking

resource "aws_vpc" "proxy_vpc" {
  cidr_block           = var.vpc_cidr
  enable_dns_support   = true
  enable_dns_hostnames = true
  tags                 = { Name = "proxy-vpc" }
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.proxy_vpc.id
  tags   = { Name = "proxy-igw" }
}

data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_subnet" "public_1" {
  vpc_id                  = aws_vpc.proxy_vpc.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = data.aws_availability_zones.available.names[0]
  map_public_ip_on_launch = true
  tags                    = { Name = "proxy-public-1" }
}

resource "aws_subnet" "public_2" {
  vpc_id                  = aws_vpc.proxy_vpc.id
  cidr_block              = "10.0.2.0/24"
  availability_zone       = data.aws_availability_zones.available.names[1]
  map_public_ip_on_launch = true
  tags                    = { Name = "proxy-public-2" }
}

resource "aws_subnet" "private_1" {
  vpc_id                  = aws_vpc.proxy_vpc.id
  cidr_block              = "10.0.3.0/24"
  availability_zone       = data.aws_availability_zones.available.names[0]
  map_public_ip_on_launch = false
  tags                    = { Name = "proxy-private-1" }
}

resource "aws_subnet" "private_2" {
  vpc_id                  = aws_vpc.proxy_vpc.id
  cidr_block              = "10.0.4.0/24"
  availability_zone       = data.aws_availability_zones.available.names[1]
  map_public_ip_on_launch = false
  tags                    = { Name = "proxy-private-2" }
}

resource "aws_route_table" "public_rt" {
  vpc_id = aws_vpc.proxy_vpc.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }
}

resource "aws_route_table_association" "public_1_assoc" {
  subnet_id      = aws_subnet.public_1.id
  route_table_id = aws_route_table.public_rt.id
}

resource "aws_route_table_association" "public_2_assoc" {
  subnet_id      = aws_subnet.public_2.id
  route_table_id = aws_route_table.public_rt.id
}

# Security Group for ALB

resource "aws_security_group" "alb_sg" {
  name   = "proxy-alb-sg"
  vpc_id = aws_vpc.proxy_vpc.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Load Balancer and Target Group

resource "aws_lb" "proxy_alb" {
  name               = "proxy-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb_sg.id]
  subnets            = [aws_subnet.public_1.id, aws_subnet.public_2.id]
}

resource "aws_lb_target_group" "proxy_tg" {
  name     = "proxy-tg"
  port     = var.processor_port
  protocol = "HTTP"
  vpc_id   = aws_vpc.proxy_vpc.id
}

resource "aws_lb_listener" "http_listener" {
  load_balancer_arn = aws_lb.proxy_alb.arn
  port              = "80"
  protocol          = "HTTP"

  # Default action: Drop everything else with a 401
  default_action {
    type = "fixed-response"
    fixed_response {
      content_type = "text/plain"
      message_body = "Unauthorized"
      status_code  = "401"
    }
  }
}

resource "aws_secretsmanager_secret" "proxy_auth_token" {
  name                    = "liteproxy/proxy-auth-token"
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret_version" "proxy_auth_token_version" {
  secret_id     = aws_secretsmanager_secret.proxy_auth_token.id
  secret_string = var.proxy_auth_token
}

resource "aws_lb_listener_rule" "auth_rule" {
  listener_arn = aws_lb_listener.http_listener.arn
  priority     = 100

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.proxy_tg.arn
  }

  condition {
    http_header {
      http_header_name = "X-Proxy-Auth"
      values           = [aws_secretsmanager_secret_version.proxy_auth_token_version.secret_string]
    }
  }
}
